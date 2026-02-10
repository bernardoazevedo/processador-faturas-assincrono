package dates

import (
	"fmt"
	"time"
)

func FormatDateYMD(data time.Time) string {
	return fmt.Sprintf("%d-%d-%d", data.Year(), data.Month(), data.Day())
}

func ActualDateYMD() string {
	return FormatDateYMD(time.Now())
}

func ActualDateHMS() string {
	hour, minute, second := time.Now().Clock()
	return fmt.Sprintf("%d:%d:%d", hour, minute, second)
}
