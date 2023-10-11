package util

import "time"

const (
	OneDay   = time.Hour * 24
	OneWeek  = OneDay * 7
	OneMonth = OneDay * 30
)

// FormatTime 格式化时间
func FormatTime(time time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}

// ParseTime 将时间字符串转为Time
func ParseTime(str string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", str)
}

// UnixMilliTime 将时间转化为毫秒数
func UnixMilliTime(t time.Time) int64 {
	return t.UnixNano() / 1000000
}

// UnunixMilliTime 将毫秒数转为为时间
func UnunixMilliTime(unix int64) time.Time {
	return time.Unix(0, unix*1000000)
}

// 随机时间
func RandDuration(t time.Duration) time.Duration {
	t1 := int64(t.Seconds() / 2)
	t2 := r.Int63n(t1)
	return t + time.Duration(t2)*time.Second
}

func GetDays(start, end int64) (timeSilce []string) {
	startTime := time.Unix(start, 0)
	endTime := time.Unix(end, 0)
	sub := int(endTime.Sub(startTime).Hours())
	days := sub / 24
	if (sub % 24) > 0 {
		days = days + 1
	}

	for i := days; i > 0; i-- {
		timeSilce = append(timeSilce, time.Now().AddDate(0, 0, -i).Format("2006-01-02"))
	}

	return
}

// GetDiffDays 获取两个时间相差的天数，0表同一天，正数表t1>t2，负数表t1<t2
func GetDiffDays(t1, t2 time.Time) int {
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, time.Local)
	t2 = time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, time.Local)

	return int(t1.Sub(t2).Hours() / 24)
}

// GetDiffDaysBySecond 获取t1和t2的相差天数，单位：秒，0表同一天，正数表t1>t2，负数表t1<t2
func GetDiffDaysBySecond(t1, t2 int64) int {
	time1 := time.Unix(t1, 0)
	time2 := time.Unix(t2, 0)

	// 调用上面的函数
	return GetDiffDays(time1, time2)
}

// GetFirstDateOfWeek 获取本周周一的日期
func GetFirstDateOfWeek(t time.Time) time.Time {
	offset := int(time.Monday - t.Weekday())
	if offset > 0 {
		offset = -6
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).
		AddDate(0, 0, offset)
}

// GetLastDateOfWeek 获取本周周日
func GetLastDateOfWeek(t time.Time) time.Time {
	return GetFirstDateOfWeek(t).
		AddDate(0, 0, 7)
}
