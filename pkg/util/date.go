package util

import "time"

const (
	FORMAT_DATE_y4Md        = "20060102"
	FORMAT_DATETIME_y4Md    = "2006-01-02"
	FORMAT_DATETIME_Y4MDHMS = "2006-01-02 15:04:05"
)

func GetDate(format string) string {
	return time.Now().Format(format)
}
