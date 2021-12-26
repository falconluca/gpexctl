package test

import (
	"encoding/base64"
	"fmt"
	"math"
	"net/url"
	"testing"
	"time"
)

func TestRFC3339(t *testing.T) {
	// RFC3339
	// https://medium.com/easyread/understanding-about-rfc-3339-for-datetime-formatting-in-software-engineering-940aa5d5f68a

	// https://stackoverflow.com/questions/21814874/how-do-i-format-an-unix-timestamp-to-rfc3339-golang
	o := time.Unix(time.Now().Unix(), 0).Format(time.RFC3339)
	fmt.Printf("o: %v", o)
	// o: 2021-11-18T23:48:52+08:00

	// https://forum.golangbridge.org/t/how-to-convert-datetime-string-to-rfc3339/22200
	createdOn, _ := time.Parse("2006-01-02 15:04:05", "2021-11-18 12:00:59")
	fmt.Printf("\no: %v", createdOn.Format(time.RFC3339))
	// o: 2006-01-02T15:04:05Z

	// https://gobyexample.com/time-formatting-parsing
}

func TestTime(t *testing.T) {
	period := 12
	date := time.Now().AddDate(0, 0, -period)
	// https://stackoverflow.com/questions/20234104/how-to-format-current-time-using-a-yyyymmddhhmmss-format
	createdOn, _ := time.Parse("2006-01-02 15:04:05", date.Format("2006-01-02 15:04:05"))
	fmt.Println(createdOn.Format(time.RFC3339))
}

func TestUrlEncode(t *testing.T) {
	s := "守望先锋"
	fmt.Println(encodeParam(s))
	fmt.Println(encodeStringBase64(s))
}

func encodeParam(s string) string {
	return url.QueryEscape(s)
}

func encodeStringBase64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func TestParseTime(t *testing.T) {
	publishedAt, _ := time.Parse(time.RFC3339, "2021-11-18T13:00:40Z")
	fmt.Println(publishedAt.Format("2006-01-02 15:04:05"))
	sub := time.Now().Sub(publishedAt)
	days := math.Floor(sub.Hours() / 24)
	days2 := math.Ceil(sub.Hours() / 24)
	fmt.Println(days)
	fmt.Println(days2)
}
