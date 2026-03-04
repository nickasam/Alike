package main

import (
	"fmt"
	"log"

	"github.com/Alike/internal/api/routes"
	"github.com/Alike/internal/config"
	"github.com/Alike/internal/domain"
	"github.com/Alike/internal/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	log.Println("✅ 配置加载成功")

	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 连接数据库
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	log.Println("✅ 数据库连接成功")

	// 运行数据库迁移
	if err := runMigrations(db); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 初始化 repositories
	userRepo := repository.NewUserRepository(db)
	globalChatRepo := repository.NewGlobalChatRepository(db)

	log.Println("✅ Repositories 初始化成功")

	// 创建默认全局聊天室
	if err := createDefaultGlobalRoom(db); err != nil {
		log.Printf("警告: 创建默认聊天室失败: %v", err)
	}

	// 设置路由
	r := gin.Default()
	routes.SetupRoutes(r, userRepo, globalChatRepo, cfg.JWT.Secret)

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("✅ 服务器启动在 %s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}

// runMigrations 运行数据库迁移
func runMigrations(db *gorm.DB) error {
	log.Println("运行数据库迁移...")

	// 自动迁移所有模型
	if err := db.AutoMigrate(
		&domain.User{},
		&domain.Match{},
		&domain.Chat{},
		&domain.Message{},
		&domain.Notification{},
		&domain.GlobalChatRoom{},
		&domain.GlobalMessage{},
	); err != nil {
		return err
	}

	log.Println("✅ 数据库迁移完成")
	return nil
}

// createDefaultGlobalRoom 创建默认全局聊天室
func createDefaultGlobalRoom(db *gorm.DB) error {
	defaultRoom := &domain.GlobalChatRoom{
		ID:          "global",
		Name:        "Alike大家庭",
		Description: "欢迎来到Alike大家庭！",
		MaxMembers:  1000,
	}

	result := db.FirstOrCreate(defaultRoom, "id = ?", "global")
	if result.Error != nil {
		return result.Error
	}

	log.Println("✅ 默认全局聊天室已就绪")
	return nil
}
