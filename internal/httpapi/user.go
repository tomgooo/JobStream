package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/tomgooo/JobStream/internal/common/db"
	"github.com/tomgooo/JobStream/internal/model"
	"golang.org/x/crypto/bcrypt"
)

// 通用响应结构（可以根据需要扩展）
type baseResponse struct {
	Code    int         `json:"code"`           // 0 表示成功，其他表示错误
	Message string      `json:"message"`        // 错误信息或简单提示
	Data    interface{} `json:"data,omitempty"` // 具体数据
}

// WriteJSON 工具函数：返回 JSON
func WriteJSON(w http.ResponseWriter, statusCode int, resp baseResponse) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(resp)
}

// WriteError 工具函数：返回错误
func WriteError(w http.ResponseWriter, statusCode int, msg string) {
	WriteJSON(w, statusCode, baseResponse{
		Code:    statusCode,
		Message: msg,
	})
}

// ========== 注册 ==========

// RegisterRequest 用于解析注册请求体
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// HandleRegister 处理用户注册：POST /api/v1/register
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	if req.Email == "" || req.Password == "" {
		WriteError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	// 检查邮箱是否已存在
	var count int64
	if err := db.DB.Model(&model.User{}).
		Where("email = ?", req.Email).
		Count(&count).Error; err != nil {
		WriteError(w, http.StatusInternalServerError, "db error")
		return
	}
	if count > 0 {
		WriteError(w, http.StatusBadRequest, "email already registered")
		return
	}

	// 生成密码哈希
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	u := &model.User{
		Email:        req.Email,
		PasswordHash: string(hash),
		// CreatedAt/UpdatedAt 由 GORM 自动填充（如果你想手动也可以用 time.Now()）
	}

	if err := db.DB.Create(u).Error; err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	// 返回给前端时不要带密码哈希
	respData := map[string]interface{}{
		"id":    u.ID,
		"email": u.Email,
	}

	WriteJSON(w, http.StatusCreated, baseResponse{
		Code:    0,
		Message: "register ok",
		Data:    respData,
	})
}

// ========== 登录 ==========

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// HandleLogin 处理用户登录：POST /api/v1/login
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	if req.Email == "" || req.Password == "" {
		WriteError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	// 查找用户
	var u model.User
	if err := db.DB.
		Where("email = ?", req.Email).
		First(&u).Error; err != nil {
		// 不暴露“用户不存在”的细节，统一返回 401
		WriteError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	// 比对密码
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		WriteError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	// TODO：这里将来可以生成 JWT / session token，这里先简单返回登录成功
	respData := map[string]interface{}{
		"id":       u.ID,
		"email":    u.Email,
		"login_at": time.Now().Format(time.RFC3339),
		"token":    "", // 预留字段，后面可以填真正的 token
	}

	WriteJSON(w, http.StatusOK, baseResponse{
		Code:    0,
		Message: "login ok",
		Data:    respData,
	})
}
