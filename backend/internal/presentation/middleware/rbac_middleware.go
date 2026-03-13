package middleware

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"project-mgmt/backend/internal/application/service"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Authorizer struct {
	publicKey *rsa.PublicKey
}

func NewAuthorizer(publicKey *rsa.PublicKey) *Authorizer {
	return &Authorizer{publicKey: publicKey}
}

// Authenticate verifies the RS256 JWT
func (a *Authorizer) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return a.publicKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		c.Set("userID", claims["sub"])
		c.Set("roleID", claims["role"])
		c.Set("isAdmin", claims["admin"])
		c.Set("perms", claims["perms"])
		c.Next()
	}
}

// Authorize checks for specific roles (Role IDs)
func Authorize(allowedRoleIDs ...uint) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, isAdmin := GetUserInfo(c)
		if isAdmin {
			c.Next()
			return
		}

		roleIDRaw, exists := c.Get("roleID")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Không tìm thấy vai trò người dùng"})
			return
		}

		// JWT numbers are floats
		var currentRoleID uint
		if f, ok := roleIDRaw.(float64); ok {
			currentRoleID = uint(f)
		} else if i, ok := roleIDRaw.(int); ok {
			currentRoleID = uint(i)
		} else if u, ok := roleIDRaw.(uint); ok {
			currentRoleID = u
		}

		allowed := false
		for _, rid := range allowedRoleIDs {
			if rid == currentRoleID {
				allowed = true
				break
			}
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Bạn không có quyền thực hiện hành động này"})
			return
		}

		c.Next()
	}
}

// AuthorizePermission checks for a specific named permission in the JWT perms claim.
// Admin users bypass all permission checks.
func AuthorizePermission(requiredPerm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Admin bypasses all checks
		_, isAdmin := GetUserInfo(c)
		if isAdmin {
			c.Next()
			return
		}

		// Read perms from context (set during Authenticate)
		permsRaw, exists := c.Get("perms")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Không có thông tin quyền"})
			return
		}

		perms, ok := permsRaw.([]interface{})
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Định dạng quyền không hợp lệ"})
			return
		}

		for _, p := range perms {
			if pStr, ok := p.(string); ok && pStr == requiredPerm {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Bạn không có quyền: " + requiredPerm})
	}
}

// AuthorizeProjectAccess checks if the user has access to a specific project.
// It extracts the project ID from the URL param "id".
func (a *Authorizer) AuthorizeProjectAccess(pmService service.ProjectMemberService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Admin bypasses all checks
		_, isAdmin := GetUserInfo(c)
		if isAdmin {
			c.Next()
			return
		}

		projectIDStr := c.Param("id")
		if projectIDStr == "" {
			c.Next() // No project ID in param, skip check
			return
		}

		projectID, err := strconv.Atoi(projectIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ID dự án không hợp lệ"})
			return
		}

		userID, _ := GetUserInfo(c)
		isMember, _, err := pmService.IsMember(c.Request.Context(), projectID, userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Lỗi kiểm tra quyền truy cập dự án: " + err.Error()})
			return
		}

		if !isMember {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Bạn không có quyền truy cập dự án này"})
			return
		}

		c.Next()
	}
}

// AuthorizeProjectPermission checks if the user has a specific permission in the project.
// It resolves project ID from URL param "id".
func (a *Authorizer) AuthorizeProjectPermission(pmService service.ProjectMemberService, prService service.ProjectRoleService, requiredPerm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Admin bypasses all checks
		_, isAdmin := GetUserInfo(c)
		if isAdmin {
			c.Next()
			return
		}

		projectIDStr := c.Param("id")
		if projectIDStr == "" {
			c.Next()
			return
		}

		projectID, err := strconv.Atoi(projectIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ID dự án không hợp lệ"})
			return
		}

		userID, _ := GetUserInfo(c)
		isMember, roleID, err := pmService.IsMember(c.Request.Context(), projectID, userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Lỗi kiểm tra quyền dự án: " + err.Error()})
			return
		}

		if !isMember {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Bạn không phải là thành viên của dự án này"})
			return
		}

		// Get permissions for the user's role in this project
		perms, err := prService.GetRolePermissions(c.Request.Context(), uint(roleID))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Lỗi lấy quyền vai trò dự án: " + err.Error()})
			return
		}

		for _, p := range perms {
			if p.Name == requiredPerm {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Bạn không có quyền dự án: " + requiredPerm})
	}
}

// GetUserInfo extracts the current user's ID and admin status from the Gin context.
func GetUserInfo(c *gin.Context) (int, bool) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		return 0, false
	}

	var userID int
	// JWT numbers are float64 by default in golang-jwt/v5
	if f, ok := userIDRaw.(float64); ok {
		userID = int(f)
	} else if i, ok := userIDRaw.(int); ok {
		userID = i
	} else if u, ok := userIDRaw.(uint); ok {
		userID = int(u)
	} else if s, ok := userIDRaw.(string); ok {
		userID, _ = strconv.Atoi(s)
	}

	isAdminRaw, exists := c.Get("isAdmin")
	isAdmin := false
	if exists {
		if val, ok := isAdminRaw.(bool); ok {
			isAdmin = val
		} else if f, ok := isAdminRaw.(float64); ok {
			isAdmin = f != 0
		} else if i, ok := isAdminRaw.(int); ok {
			isAdmin = i != 0
		}
	}

	return userID, isAdmin
}
