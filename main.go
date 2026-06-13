package main

import (
	"fmt"
	"log"
	"serve-wechat-gin/config"
	"serve-wechat-gin/database"
	"serve-wechat-gin/router"
)

func main() {
	// 加载配置文件
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// 打印 MySQL DSN（生产环境应该移除或使用日志库）
	fmt.Printf("MySQL DSN: %s\n", cfg.MySQL.GetDSN())

	// 初始化数据库（如果不存在则自动创建）
	if err := database.InitDB(cfg); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 自动迁移数据库表
	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("数据库表迁移失败: %v", err)
	}

	// 初始化基础数据
	database.SeedData()

	r := router.SetupRouter()

	// 启动服务器，监听 0.0.0.0:8081
	if err := r.Run(":8081"); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
