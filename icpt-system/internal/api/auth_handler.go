// Package api 提供HTTP API处理器
// 包含用户认证、图像管理等核心业务功能
package api

import (
	"log"
	"net/http"

	"icpt-system/internal/models"
	"icpt-system/internal/services"
	"icpt-system/internal/store"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterHandler 处理用户注册请求
// @Summary 用户注册
// @Description 创建新用户账户，返回JWT令牌
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "注册信息"
// @Success 201 {object} models.AuthResponse "注册成功"
// @Failure 400 {object} map[string]interface{} "请求数据错误"
// @Failure 409 {object} map[string]interface{} "用户已存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/auth/register [post]
func RegisterHandler(c *gin.Context) {
	var req models.RegisterRequest

	// 解析并验证请求数据
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "请求数据格式错误",
			"code":    "INVALID_REQUEST",
			"details": err.Error(),
		})
		return
	}

	// 检查用户名或邮箱是否已被使用
	var existingUser models.User
	result := store.DB.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser)
	if result.Error != gorm.ErrRecordNotFound {
		if result.Error == nil {
			// 用户已存在，返回冲突错误
			c.JSON(http.StatusConflict, gin.H{
				"error": "用户名或邮箱已存在",
				"code":  "USER_EXISTS",
			})
			return
		}
		// 数据库查询错误
		log.Printf("数据库查询错误: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "服务器内部错误",
			"code":  "DATABASE_ERROR",
		})
		return
	}

	// 创建新用户实例
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Status:   "active", // 新用户默认为活跃状态
	}

	// 对密码进行安全哈希处理
	if err := user.HashPassword(req.Password); err != nil {
		log.Printf("密码哈希错误: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "服务器内部错误",
			"code":  "PASSWORD_HASH_ERROR",
		})
		return
	}

	// 将用户信息保存到数据库
	if err := store.DB.Create(&user).Error; err != nil {
		log.Printf("创建用户错误: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建用户失败",
			"code":  "CREATE_USER_ERROR",
		})
		return
	}

	// 生成JWT认证令牌
	token, err := services.GenerateToken(user.ID, user.Username)
	if err != nil {
		log.Printf("生成令牌错误: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "生成认证令牌失败",
			"code":  "TOKEN_GENERATION_ERROR",
		})
		return
	}

	// 记录成功注册日志
	log.Printf("用户注册成功: %s (ID: %d)", user.Username, user.ID)

	// 构造并返回成功响应
	response := models.AuthResponse{
		User:  &user,
		Token: token,
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "注册成功",
		"data":    response,
	})
}

// LoginHandler 处理用户登录请求
// @Summary 用户登录
// @Description 验证用户凭据，返回JWT令牌
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "登录信息"
// @Success 200 {object} models.AuthResponse "登录成功"
// @Failure 400 {object} map[string]interface{} "请求数据错误"
// @Failure 401 {object} map[string]interface{} "认证失败"
// @Failure 403 {object} map[string]interface{} "账户被禁用"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/auth/login [post]
func LoginHandler(c *gin.Context) {
	var req models.LoginRequest

	// 解析并验证登录请求数据
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "请求数据格式错误",
			"code":    "INVALID_REQUEST",
			"details": err.Error(),
		})
		return
	}

	// 查找用户（支持用户名或邮箱登录）
	var user models.User
	result := store.DB.Where("username = ? OR email = ?", req.Username, req.Username).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// 用户不存在，返回认证失败（不泄露具体原因）
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "用户名或密码错误",
				"code":  "INVALID_CREDENTIALS",
			})
			return
		}
		// 数据库查询错误
		log.Printf("数据库查询错误: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "服务器内部错误",
			"code":  "DATABASE_ERROR",
		})
		return
	}

	// 检查用户账户状态
	if user.Status != "active" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "账户已被禁用",
			"code":  "ACCOUNT_DISABLED",
		})
		return
	}

	// 验证密码是否正确
	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "用户名或密码错误",
			"code":  "INVALID_CREDENTIALS",
		})
		return
	}

	// 生成新的JWT认证令牌
	token, err := services.GenerateToken(user.ID, user.Username)
	if err != nil {
		log.Printf("生成令牌错误: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "生成认证令牌失败",
			"code":  "TOKEN_GENERATION_ERROR",
		})
		return
	}

	// 记录成功登录日志
	log.Printf("用户登录成功: %s (ID: %d)", user.Username, user.ID)

	// 构造并返回成功响应
	response := models.AuthResponse{
		User:  &user,
		Token: token,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"data":    response,
	})
}

// GetProfileHandler 获取当前用户信息
// @Summary 获取用户信息
// @Description 获取当前认证用户的详细信息
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User "用户信息"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 404 {object} map[string]interface{} "用户不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/profile [get]
func GetProfileHandler(c *gin.Context) {
	// 从上下文中获取用户ID（由认证中间件设置）
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未认证",
			"code":  "UNAUTHENTICATED",
		})
		return
	}

	// 从数据库中查询用户详细信息
	var user models.User
	if err := store.DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "用户未找到",
				"code":  "USER_NOT_FOUND",
			})
			return
		}
		log.Printf("数据库查询错误: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "服务器内部错误",
			"code":  "DATABASE_ERROR",
		})
		return
	}

	// 返回用户信息
	c.JSON(http.StatusOK, gin.H{
		"message": "获取用户信息成功",
		"data":    user,
	})
}
