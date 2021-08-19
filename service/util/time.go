package util

import "time"

func DayBeginStamp(now int64) int64 {
	_, offset := time.Now().Zone()
	return now - (now+int64(offset))%int64(3600*24)
}
