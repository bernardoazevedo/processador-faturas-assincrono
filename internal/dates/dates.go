package dates

import (
	"fmt"
	"time"
)

func FormatDateYMD(data time.Time) string {
	return fmt.Sprintf("%04d-%02d-%02d", data.Year(), data.Month(), data.Day())
}

func ActualDateYMD() string {
	return FormatDateYMD(time.Now())
}

func ActualDateHMS() string {
	hour, minute, second := time.Now().Clock()
	return fmt.Sprintf("%d:%d:%d", hour, minute, second)
}
