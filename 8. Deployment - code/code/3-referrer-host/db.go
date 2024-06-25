package main

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
		Addr: []string{"127.0.0.1:9000"},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "",
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
			Referrer_domain String NOT NULL,
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
			nowToInt(),
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

func nowToInt() uint32 {
	now := time.Now().Format("20060102")
	i, err := strconv.ParseInt(now, 10, 32)
	// this should never happen
	if err != nil {
		log.Fatal(err)
	}
	return uint32(i)
}
