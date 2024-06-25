package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/mileusna/useragent"
)

type Events struct {
	DB driver.Conn
}

type Event struct {
	ID          int64
	SiteID      string
	OccuredAt   int32
	Type        string
	UserID      string
	Event       string
	Category    string
	Referrer    string
	IsTouch     bool
	BrowserName string
	OSName      string
	DeviceType  string
	Country     string
	Region      string
	Timestamp   time.Time
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

func (e *Events) Add(trk Tracking, ua useragent.UserAgent, geo *GeoInfo) error {
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

	err := e.DB.Exec(context.Background(), qry,
		trk.SiteID,
		nowToInt(),
		trk.Action.Type,
		trk.Action.Identity,
		trk.Action.Event,
		trk.Action.Category,
		trk.Action.Referrer,
		trk.Action.IsTouchDevice,
		ua.Name,
		ua.OS,
		ua.Device,
		geo.Country,
		geo.RegionName,
	)

	return err
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
