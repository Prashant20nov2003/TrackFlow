package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/mileusna/useragent"
)

type Events struct {
	DB *pgx.Conn
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
	conn, err := pgx.Connect(
		context.Background(),
		"postgres://postgres:password@localhost:5432/postgres",
	)
	if err != nil {
		return err
	} else if err := conn.Ping(context.Background()); err != nil {
		return err
	}

	e.DB = conn
	return nil
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

	_, err := e.DB.Exec(context.Background(), qry,
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
