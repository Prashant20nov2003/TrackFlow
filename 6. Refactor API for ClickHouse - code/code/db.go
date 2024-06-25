package gotracker

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/mileusna/useragent"
)

type QueryType int

const (
	QueryPageViews QueryType = iota
	QueryPageViewList
	QueryUniqueVisitors
	QueryReferrerHost
	QueryReferrer
	QueryBrowsers
	QueryOSes
	QueryCountry
)

type TrackingData struct {
	Type          string `json:"type"`
	Identity      string `json:"identity"`
	UserAgent     string `json:"ua"`
	Event         string `json:"event"`
	Category      string `json:"category"`
	Referrer      string `json:"referrer"`
	ReferrerHost  string
	IsTouchDevice bool `json:"isTouchDevice"`
	OccuredAt     uint32
}

type Tracking struct {
	SiteID string       `json:"site_id"`
	Action TrackingData `json:"tracking"`
}

type MetricData struct {
	What   QueryType `json:"what"`
	SiteID string    `json:"siteId"`
	Start  uint32    `json:"start"`
	End    uint32    `json:"end"`
	Extra  string    `json:"extra"`
}

type qdata struct {
	trk Tracking
	ua  useragent.UserAgent
	geo *GeoInfo
}

type Events struct {
	DB   driver.Conn
	ch   chan qdata
	lock sync.RWMutex
	q    []qdata
}

func (e *Events) Open() error {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{config.ClickHouseHost},
		Auth: clickhouse.Auth{
			Database: config.ClickHouseDB,
			Username: config.ClickHouseUser,
			Password: config.ClickHousePassword,
		},
		Debug: false,
		Debugf: func(format string, v ...any) {
			fmt.Printf(format, v)
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		DialTimeout:          time.Second * 30,
		MaxOpenConns:         5,
		MaxIdleConns:         5,
		ConnMaxLifetime:      time.Duration(10) * time.Minute,
		ConnOpenStrategy:     clickhouse.ConnOpenInOrder,
		BlockBufferSize:      10,
		MaxCompressionBuffer: 10240,
		ClientInfo: clickhouse.ClientInfo{ // optional, please see Client info section in the README.md
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "my-app", Version: "0.1"},
			},
		},
	})
	if err != nil {
		return err
	}

	if err := conn.Ping(context.Background()); err != nil {
		return err
	}

	e.DB = conn
	return nil
}

func (e *Events) EnsureTable() error {
	qry := `		
		CREATE TABLE IF NOT EXISTS events (
			site_id String NOT NULL,
			occured_at UInt32 NOT NULL,
			type String NOT NULL,
			user_id String NOT NULL,
			event String NOT NULL,
			category String NOT NULL,
			referrer String NOT NULL,
			referrer_domain String NOT NULL,
			is_touch BOOLEAN NOT NULL,
			browser_name String NOT NULL,
			os_name String NOT NULL,
			device_type String NOT NULL,
			country String NOT NULL,
			region String NOT NULL,
			timestamp DateTime DEFAULT now()
		)
		ENGINE MergeTree
		ORDER BY (site_id, occured_at);
	`

	ctx := context.Background()
	return e.DB.Exec(ctx, qry)
}

func (e *Events) Add(trk Tracking, ua useragent.UserAgent, geo *GeoInfo) {
	e.ch <- qdata{trk, ua, geo}
}

func (e *Events) Run() {
	e.ch = make(chan qdata)

	timer := time.NewTimer(time.Second * 10)
	for {
		select {
		case data := <-e.ch:
			e.lock.Lock()
			e.q = append(e.q, data)
			c := len(e.q)
			e.lock.Unlock()

			if c >= 15 {
				if err := e.Insert(); err != nil {
					fmt.Println("error while inserting data: ", err)
				}
			}
		case <-timer.C:
			timer.Reset(time.Second * 10)

			e.lock.RLock()
			c := len(e.q)
			e.lock.RUnlock()

			if c > 0 {
				if err := e.Insert(); err != nil {
					fmt.Println("error while inserting data: ", err)
				}
			}
		}
	}
}

func (e *Events) Insert() error {
	var tmp []qdata
	e.lock.Lock()
	for _, qd := range e.q {
		tmp = append(tmp, qd)
	}

	e.q = nil
	e.lock.Unlock()

	qry := `
		INSERT INTO events
		(
			site_id,
			occured_at,
			type,
			user_id,
			event,
			category,
			referrer,
			referrer_domain,
			is_touch,
			browser_name,
			os_name,
			device_type,
			country,
			region
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)
	`

	ctx := context.Background()
	batch, err := e.DB.PrepareBatch(ctx, qry)
	if err != nil {
		return err
	}

	for _, qd := range tmp {
		err := batch.Append(
			qd.trk.SiteID,
			TimeToInt(time.Now()),
			qd.trk.Action.Type,
			qd.trk.Action.Identity,
			qd.trk.Action.Event,
			qd.trk.Action.Category,
			qd.trk.Action.Referrer,
			qd.trk.Action.ReferrerHost,
			qd.trk.Action.IsTouchDevice,
			qd.ua.Name,
			qd.ua.OS,
			qd.ua.Device,
			qd.geo.Country,
			qd.geo.RegionName,
		)

		if err != nil {
			return err
		}
	}

	return batch.Send()
}

func TimeToInt(d time.Time) uint32 {
	now := d.Format("20060102")
	i, err := strconv.ParseInt(now, 10, 32)
	// this should never happen
	if err != nil {
		log.Fatal(err)
	}
	return uint32(i)
}

type Metric struct {
	OccuredAt uint32 `json:"occuredAt"`
	Value     string `json:"value"`
	Count     uint64 `json:"count"`
}

func (e *Events) GetStats(data MetricData) ([]Metric, error) {
	qry := e.GenQuery(data)

	rows, err := e.DB.Query(
		context.Background(),
		qry,
		data.SiteID,
		data.Start,
		data.End,
		data.Extra,
	)
	if err != nil {
		return nil, err
	}

	var metrics []Metric
	for rows.Next() {
		var m Metric
		if err := rows.Scan(&m.OccuredAt, &m.Value, &m.Count); err != nil {
			return nil, err
		}

		metrics = append(metrics, m)
	}

	return metrics, rows.Err()
}

func (e *Events) GenQuery(data MetricData) string {
	field := ""
	daily := true
	where := "AND $4 = $4"
	switch data.What {
	case QueryPageViews:
		field = "event"
	case QueryPageViewList:
		field = "event"
		daily = false
	case QueryUniqueVisitors:
		field = "user_id"
	case QueryReferrerHost:
		field = "referrer_domain"
		daily = false
	case QueryReferrer:
		field = "referrer"
		where = "AND referrer_domain = $3 "
		daily = false
	case QueryBrowsers:
		field = "browser_name"
		daily = false
	case QueryOSes:
		field = "os_name"
		daily = false
	case QueryCountry:
		field = "country"
		daily = false
	}

	if daily {
		return fmt.Sprintf(`
		SELECT occured_at, %s, COUNT(*)
		FROM events
		WHERE site_id = $1
		AND category = 'Page views'
		GROUP BY occured_at, %s
		HAVING occured_at BETWEEN $2 AND $3
		ORDER BY 3 DESC;
	`, field, field)
	}

	return fmt.Sprintf(`
		SELECT toUInt32(0), %s, COUNT(*)
		FROM events
		WHERE site_id = $1
		AND occured_at BETWEEN $2 AND $3
		AND category = 'Page views'
		%s 
		GROUP BY %s
		ORDER BY 3 DESC;
	`, field, where, field)
}
