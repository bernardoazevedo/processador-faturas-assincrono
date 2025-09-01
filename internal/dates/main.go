package dates

import (
	"fmt"
	"time"
)

func FormatDateYMD(data time.Time) string {
	year := data.Year()
	month := data.Month()
	day := data.Day()

	return fmt.Sprintf("%d-%d-%d", year, month, day)
}

func ActualDateYMD() string {
	return FormatDateYMD(time.Now())
}

func ActualDateHMS() string {
	hour, minute, second := time.Now().Clock()
	return fmt.Sprintf("%d:%d:%d", hour, minute, second)
}