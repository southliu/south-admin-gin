package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const TimeFormat = "2006-01-02 15:04:05"

// CustomTime 自定义时间类型，JSON序列化为 YYYY-MM-DD HH:mm:ss 格式
type CustomTime int64

// MarshalJSON 序列化为 YYYY-MM-DD HH:mm:ss 格式
func (t CustomTime) MarshalJSON() ([]byte, error) {
	if t == 0 {
		return []byte(`""`), nil
	}
	tm := time.Unix(int64(t), 0)
	return []byte(`"` + tm.Format(TimeFormat) + `"`), nil
}

// UnmarshalJSON 从字符串或数字反序列化
func (t *CustomTime) UnmarshalJSON(data []byte) error {
	if string(data) == `""` || string(data) == "null" {
		*t = 0
		return nil
	}
	var ts int64
	if _, err := fmt.Sscanf(string(data), "%d", &ts); err == nil {
		*t = CustomTime(ts)
		return nil
	}
	tm, err := time.Parse(`"`+TimeFormat+`"`, string(data))
	if err != nil {
		return err
	}
	*t = CustomTime(tm.Unix())
	return nil
}

// Scan 从数据库读取 DATETIME 字段
func (t *CustomTime) Scan(value interface{}) error {
	if value == nil {
		*t = 0
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		*t = CustomTime(v.Unix())
	case int64:
		*t = CustomTime(v)
	default:
		*t = 0
	}
	return nil
}

// Value 写入数据库时转换为 DATETIME
func (t CustomTime) Value() (driver.Value, error) {
	if t == 0 {
		return nil, nil
	}
	return time.Unix(int64(t), 0), nil
}
