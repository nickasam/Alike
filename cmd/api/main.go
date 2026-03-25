package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Alike/internal/domain"
	"github.com/Alike/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var jwtSecret = []byte("alike-secret-key-2024")

func main() {
	// 数据库连接
	dsn := "host=localhost user=alike_user password=alike_password dbname=alike_db port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 自动迁移
	log.Println("Running migrations...")
	db.AutoMigrate(&domain.User{}, &domain.Match{}, &domain.Chat{}, &domain.Message{}, &domain.Notification{})
	db.AutoMigrate(&domain.GlobalChatRoom{}, &domain.GlobalMessage{})
	log.Println("✅ Migrations completed")

	// 初始化repository
	userRepo := repository.NewUserRepository(db)
	globalChatRepo := repository.NewGlobalChatRepository(db)

	// 创建默认全局聊天室
	defaultRoom := &domain.GlobalChatRoom{
		ID:          "global",
		Name:        "Alike大家庭",
		Description: "欢迎来到Alike大家庭！",
		MaxMembers:  1000,
	}
	db.FirstOrCreate(defaultRoom, "id = ?", "global")
	log.Println("✅ Global chat room ready")

	// 设置路由
	r := gin.Default()

	// CORS 中间件 - 必须在所有路由之前设置
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "alike-api", "version": "1.0.0"})
	})

	// API v1
	v1 := r.Group("/api/v1")
	{
		// 注册
		v1.POST("/auth/register", func(c *gin.Context) {
			var req struct {
				Phone    string `json:"phone" binding:"required"`
				Password string `json:"password" binding:"required"`
				Nickname string `json:"nickname" binding:"required"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
				return
			}

			existingUser, _ := userRepo.GetByPhone(req.Phone)
			if existingUser != nil && existingUser.ID != "" {
				c.JSON(http.StatusConflict, gin.H{"success": false, "error": gin.H{"message": "用户已存在"}})
				return
			}

			now := time.Now()
			user := &domain.User{
				ID:       uuid.New().String(),
				Phone:    req.Phone,
				Nickname: req.Nickname,
				IsActive: true,
				Timestamps: domain.Timestamps{
					CreatedAt: now,
					UpdatedAt: now,
				},
			}

			if err := userRepo.Create(user); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": "创建用户失败"}})
				return
			}

			accessToken := generateToken(user.ID)
			refreshToken := generateToken(user.ID)

			c.JSON(http.StatusCreated, gin.H{
				"success": true,
				"data": gin.H{
					"user": gin.H{
						"id":       user.ID,
						"phone":    user.Phone,
						"nickname": user.Nickname,
					},
					"tokens": gin.H{
						"access_token":  accessToken,
						"refresh_token": refreshToken,
					},
				},
			})
		})

		// 登录
		v1.POST("/auth/login", func(c *gin.Context) {
			var req struct {
				Phone    string `json:"phone" binding:"required"`
				Password string `json:"password" binding:"required"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
				return
			}

			user, err := userRepo.GetByPhone(req.Phone)
			if err != nil || user == nil || user.ID == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": gin.H{"message": "用户不存在"}})
				return
			}

			accessToken := generateToken(user.ID)
			refreshToken := generateToken(user.ID)

			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data": gin.H{
					"user": gin.H{
						"id":       user.ID,
						"phone":    user.Phone,
						"nickname": user.Nickname,
					},
					"tokens": gin.H{
						"access_token":  accessToken,
						"refresh_token": refreshToken,
					},
				},
			})
		})

		// 认证中间件
		authMiddleware := func(c *gin.Context) {
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": gin.H{"message": "未授权"}})
				c.Abort()
				return
			}

			if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
				c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": gin.H{"message": "授权格式错误"}})
				c.Abort()
				return
			}

			token := authHeader[7:]
			userID, err := validateToken(token)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": gin.H{"message": "无效token"}})
				c.Abort()
				return
			}

			c.Set("user_id", userID)
			c.Next()
		}

		// 获取当前用户
		v1.GET("/auth/me", authMiddleware, func(c *gin.Context) {
			userID := c.GetString("user_id")
			user, _ := userRepo.GetByID(userID)
			c.JSON(http.StatusOK, gin.H{"success": true, "data": user})
		})

		// 获取附近用户
		v1.GET("/users/nearby", authMiddleware, func(c *gin.Context) {
			users, _ := userRepo.List(100, 0)
			c.JSON(http.StatusOK, gin.H{"success": true, "data": users})
		})

		// 获取全局聊天室
		v1.GET("/global/room", authMiddleware, func(c *gin.Context) {
			room, _ := globalChatRepo.GetRoom("global")
			c.JSON(http.StatusOK, gin.H{"success": true, "data": room})
		})

		// 获取全局消息
		v1.GET("/global/messages", authMiddleware, func(c *gin.Context) {
			messages, _ := globalChatRepo.GetMessages("global", 100)
			c.JSON(http.StatusOK, gin.H{"success": true, "data": messages})
		})

		// 发送全局消息
		v1.POST("/global/messages", authMiddleware, func(c *gin.Context) {
			userID := c.GetString("user_id")
			log.Printf("[DEBUG] 发送消息请求 - userID: %s", userID)

			var req struct {
				Content string `json:"content" binding:"required,min=1,max=500"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				log.Printf("[DEBUG] 请求验证失败: %v", err)
				c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": "消息内容不能为空，最多500字"}})
				return
			}

			// 获取用户信息
			user, err := userRepo.GetByID(userID)
			if err != nil {
				log.Printf("[DEBUG] 获取用户失败: %v", err)
				c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": gin.H{"message": "用户不存在"}})
				return
			}
			if user == nil || user.ID == "" {
				log.Printf("[DEBUG] 用户为空")
				c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": gin.H{"message": "用户不存在"}})
				return
			}

			log.Printf("[DEBUG] 用户信息: id=%s, nickname=%s", user.ID, user.Nickname)

			message := &domain.GlobalMessage{
				ID:        time.Now().Format("20060102150405") + "000",
				RoomID:    "global",
				UserID:    userID,
				Username:  user.Nickname,
				Content:   req.Content,
				CreatedAt: time.Now(),
			}

			if err := globalChatRepo.CreateMessage(message); err != nil {
				log.Printf("[DEBUG] 创建消息失败: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": "发送消息失败"}})
				return
			}

			log.Printf("[DEBUG] 消息发送成功: id=%s", message.ID)
			c.JSON(http.StatusCreated, gin.H{"success": true, "data": message})
		})

		// 获取在线用户列表
		v1.GET("/global/online-users", authMiddleware, func(c *gin.Context) {
			userID := c.GetString("user_id")
			userNickname := c.GetString("username")

			// 获取最近活跃的用户（从全局消息中获取）
			messages, _ := globalChatRepo.GetMessages("global", 1000)

			// 去重并构建用户列表
			userMap := make(map[string]*gin.H)
			for _, msg := range messages {
				if _, exists := userMap[msg.UserID]; !exists {
					userMap[msg.UserID] = &gin.H{
						"id":       msg.UserID,
						"nickname": msg.Username,
					}
				}
			}

			// 确保当前用户在列表中
			if _, exists := userMap[userID]; !exists {
				userMap[userID] = &gin.H{
					"id":       userID,
					"nickname": userNickname,
				}
			}

			// 转换为数组
			users := make([]gin.H, 0, len(userMap))
			for _, user := range userMap {
				users = append(users, *user)
			}

			c.JSON(http.StatusOK, gin.H{"success": true, "data": users})
		})

		// 获取在线用户数量
		v1.GET("/global/online-count", authMiddleware, func(c *gin.Context) {
			userID := c.GetString("user_id")

			// 获取最近活跃的用户数量
			messages, _ := globalChatRepo.GetMessages("global", 1000)

			// 去重统计用户数
			userSet := make(map[string]bool)
			for _, msg := range messages {
				userSet[msg.UserID] = true
			}

			// 确保当前用户被计数
			userSet[userID] = true

			c.JSON(http.StatusOK, gin.H{"success": true, "data": len(userSet)})
		})
	}

	// 启动服务器
	log.Println("✅ Server starting on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// 生成JWT token
func generateToken(userID string) string {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtSecret)
	return tokenString
}

// 验证JWT token
func validateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["user_id"].(string)
		return userID, nil
	}

	return "", err
}
