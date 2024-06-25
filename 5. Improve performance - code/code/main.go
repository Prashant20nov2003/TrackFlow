package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	TOTALROWS  = 15000000
	TOTALUSERS = 125321
	TOTALPAGES = 500
)

var (
	SITEID    = []string{"site1", "site2", "site3"}
	MONTHS    = []string{"01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12"}
	TYPES     = []string{"page", "event"}
	USERID    = make([]string, TOTALUSERS)
	PAGES     = make([]string, TOTALPAGES)
	REFERRERS = []string{"", "google", "twitter", "reddit", "siteabc.com"}
	DEVICES   = []string{"desktop", "tablet", "phone"}
	BROWSERS  = []string{"chrome", "firefox", "edge"}
	OSNAME    = []string{"linux", "windows", "macos"}
	COUNTRIES = []string{"c1", "c2", "c3", "c4", "c5", "c6", "c7", "c8"}
)

func main() {
	genUserID()
	genPages()

	rows := make([]string, TOTALROWS)
	for i := 0; i < TOTALROWS; i++ {
		rows[i] = genInsert(int64(i + 1))
	}

	s := strings.Join(rows, "")

	if err := os.WriteFile("dump.data", []byte(s), 0644); err != nil {
		fmt.Println(err)
	}
}

func genInsert(id int64) string {
	device := DEVICES[rand.Intn(len(DEVICES))]
	is_touch := "true"
	if device == "desktop" {
		is_touch = "false"
	}

	qry := "%d\t%s\t%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n"

	return fmt.Sprintf(qry,
		id,
		SITEID[rand.Intn(len(SITEID))],
		genOccuredAt(),
		TYPES[rand.Intn(len(TYPES))],
		USERID[rand.Intn(len(USERID))],
		PAGES[rand.Intn(len(PAGES))],
		"Page views",
		REFERRERS[rand.Intn(len(REFERRERS))],
		is_touch,
		BROWSERS[rand.Intn(len(BROWSERS))],
		OSNAME[rand.Intn(len(OSNAME))],
		device,
		COUNTRIES[rand.Intn(len(COUNTRIES))],
		"no need",
		time.Now().Format(time.RFC3339),
	)
}

func genOccuredAt() uint32 {
	year := time.Now().Year() - rand.Intn(3)
	month := MONTHS[rand.Intn(len(MONTHS))]
	day := rand.Intn(27) + 1 // let's play safe...

	d := fmt.Sprintf("%d%s%d", year, month, day)
	i, err := strconv.ParseInt(d, 10, 64)
	if err != nil {
		return 20231205
	}
	return uint32(i)
}

func genUserID() {
	for i := 0; i < TOTALUSERS; i++ {
		USERID[i] = fmt.Sprintf("user_%d", i)
	}
}

func genPages() {
	for i := 0; i < TOTALPAGES; i++ {
		sep := rand.Intn(3)
		if sep == 0 {
			PAGES[i] = "/"
		} else {
			p := "/"
			for j := 0; j < sep; j++ {
				p += pageName() + "/"
			}

			PAGES[i] = p
		}
	}
}

func pageName() string {
	return fmt.Sprint(
		string(rand.Intn(26)+65),
		string(rand.Intn(26)+65),
		string(rand.Intn(26)+65),
		string(rand.Intn(26)+65),
		string(rand.Intn(26)+65),
	)
}
