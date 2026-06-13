package models

// Log 日志模型
type Log struct {
	ID        int64  `json:"id" gorm:"primaryKey"`
	Username  string `json:"username" gorm:"type:varchar(100)"`
	IP        string `json:"ip" gorm:"type:varchar(50)"`
	Method    string `json:"method" gorm:"type:varchar(10)"`
	URL       string `json:"url" gorm:"type:varchar(500)"`
	Params    string `json:"params" gorm:"type:text"`
	UserAgent string `json:"userAgent" gorm:"type:varchar(500)"`
	Status    string `json:"status" gorm:"type:varchar(20)"`
	Error     string `json:"error" gorm:"type:text"`
	Latency   int64  `json:"latency" gorm:"default:0"`
	Type      int        `json:"type" gorm:"default:0;comment:0=后端错误 1=成功 2=警告 3=前端错误"`
	CreatedAt CustomTime `json:"createdAt" gorm:"autoCreateTime;comment:创建时间"`
}

// TableName 指定表名
func (Log) TableName() string {
	return "sys_log"
}

// CreateLogDto 创建日志请求
type CreateLogDto struct {
	Username  string `json:"username"`
	IP        string `json:"ip"`
	Method    string `json:"method"`
	URL       string `json:"url"`
	Params    string `json:"params"`
	UserAgent string `json:"userAgent"`
	Status    string `json:"status"`
	Error     string `json:"error"`
	Latency   int64  `json:"latency"`
	Type      int    `json:"type"`
}

// LogPageResult 日志分页结果
type LogPageResult struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	IP        string `json:"ip"`
	Method    string `json:"method"`
	URL       string `json:"url"`
	Status    string `json:"status"`
	Error     string     `json:"error"`
	Type      int        `json:"type"`
	CreatedAt CustomTime `json:"createdAt"`
}
