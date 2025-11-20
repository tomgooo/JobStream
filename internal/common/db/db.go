package db

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/tomgooo/JobStream/internal/model"
)

// 全局的 GORM DB 句柄，后面别的包可以用 db.DB 访问
var DB *gorm.DB

// Init 初始化数据库连接，并自动建表
func Init() error {
	// 根据 docker-compose.yml 里的配置拼 DSN
	// 用户名: jobstream
	// 密码: jobstreampwd
	// 数据库: jobstream
	dsn := "jobstream:jobstreampwd@tcp(127.0.0.1:3306)/jobstream?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("connect db failed: %w", err)
	}

	log.Println("✅ connected to MySQL OK")

	// 自动迁移：根据 struct 创建 / 更新表结构
	if err := DB.AutoMigrate(&model.User{}, &model.Job{}, &model.Subscription{}); err != nil {
		return fmt.Errorf("auto migrate failed: %w", err)
	}

	log.Println("✅ auto-migrate tables OK")
	return nil
}
