package database

import (
	"fmt"
	"log"
	"serve-wechat-gin/config"
	"serve-wechat-gin/models/system"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB(cfg *config.Config) error {
	dsn := cfg.MySQL.GetDSN()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(cfg.MySQL.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MySQL.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MySQL.ConnMaxLifetime) * time.Second)

	DB = db
	log.Println("数据库连接成功")
	return nil
}

// AutoMigrate 自动迁移数据库表
func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}

	err := DB.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Menu{},
		&models.Permission{},
		&models.Log{},
	)
	if err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	log.Println("数据库表迁移成功")
	return nil
}

// SeedData 初始化基础数据
func SeedData() {
	if DB == nil {
		log.Println("数据库未初始化，跳过数据初始化")
		return
	}

	log.Println("跳过数据初始化（不创建默认数据）")
}
