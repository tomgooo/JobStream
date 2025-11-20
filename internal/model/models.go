package model

import "time"

// User 用户表
type User struct {
	ID           uint64 `gorm:"primaryKey;autoIncrement"`
	Email        string `gorm:"size:255;uniqueIndex;not null"`
	PasswordHash string `gorm:"size:255;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Job 职位表
type Job struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement"`
	Title     string `gorm:"size:255;not null"` // 职位名称
	Company   string `gorm:"size:255;not null"` // 公司名称
	City      string `gorm:"size:100;index"`    // 城市
	SalaryMin int    // 最低薪资（单位：k 或者直接元，看你后面定义）
	SalaryMax int    // 最高薪资
	TechStack string `gorm:"type:text"` // 技术栈（先用文本，后续可以改 json）
	Source    string `gorm:"size:100"`  // 来源：boss、拉勾、自建等
	URL       string `gorm:"size:512"`  // 职位详情链接
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Subscription 订阅表：用户配置的订阅条件
type Subscription struct {
	ID           uint64 `gorm:"primaryKey;autoIncrement"`
	UserID       uint64 `gorm:"index;not null"` // 外键：用户 ID
	City         string `gorm:"size:100;index"` // 希望工作的城市
	MinSalary    int    // 希望的最低薪资
	TechKeywords string `gorm:"type:text"` // 技术关键词，用逗号分隔，比如 "go,gin,grpc"
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
