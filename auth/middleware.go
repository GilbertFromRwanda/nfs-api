package auth

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"nfs-api/config"
	"nfs-api/database"
	"nfs-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// jwtErrorMessage maps jwt/v5 sentinel errors to human-readable messages.
func jwtErrorMessage(err error) string {
	switch {
	case errors.Is(err, jwt.ErrTokenExpired):
		return "token has expired, please log in again"
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return "token signature is invalid"
	case errors.Is(err, jwt.ErrTokenMalformed):
		return "token is malformed"
	case errors.Is(err, jwt.ErrTokenNotValidYet):
		return "token is not valid yet"
	default:
		return "invalid or expired token"
	}
}

func JWTMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip JWT validation for CORS preflight requests.
		// Browsers send an OPTIONS request with no Authorization header before
		// every non-simple request (PUT, POST, DELETE, PATCH). Without this
		// skip those methods are blocked before CORS headers can be returned.
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		ip := c.ClientIP()

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("[AUTH] 401 missing Authorization header | ip=%s path=%s", ip, c.FullPath())
			utils.Error(c, http.StatusUnauthorized, "authorization header required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Printf("[AUTH] 401 invalid format (got %q) | ip=%s path=%s", parts[0], ip, c.FullPath())
			utils.Error(c, http.StatusUnauthorized, "invalid authorization format")
			c.Abort()
			return
		}

		if IsBlacklisted(parts[1]) {
			log.Printf("[AUTH] 401 token has been revoked | ip=%s path=%s", ip, c.FullPath())
			utils.Error(c, http.StatusUnauthorized, "token has been revoked, please log in again")
			c.Abort()
			return
		}

		token, err := jwt.Parse(parts[1], func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			msg := jwtErrorMessage(err)
			log.Printf("[AUTH] 401 %s | ip=%s path=%s err=%v", msg, ip, c.FullPath(), err)
			utils.Error(c, http.StatusUnauthorized, msg)
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Printf("[AUTH] 401 invalid token claims | ip=%s path=%s", ip, c.FullPath())
			utils.Error(c, http.StatusUnauthorized, "invalid token claims")
			c.Abort()
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			log.Printf("[AUTH] 401 invalid user_id claim | ip=%s path=%s", ip, c.FullPath())
			utils.Error(c, http.StatusUnauthorized, "invalid user id in token")
			c.Abort()
			return
		}

		c.Set("user_id", uint(userID))

		if officeID, ok := claims["office_id"].(float64); ok {
			c.Set("office_id", uint(officeID))
		}
		if officeName, ok := claims["office_name"].(string); ok {
			c.Set("office_name", officeName)
		}

		c.Next()
	}
}

// HasAccess checks if the authenticated user has at least one of the given permissions.
// Must be used after JWTMiddleware.
// Usage:
//
//	router.GET("/route", auth.JWTMiddleware(cfg), auth.HasAccess("manage_users"), handler)
//	router.GET("/route", auth.JWTMiddleware(cfg), auth.HasAccess("edit_agreements", "view_agreements"), handler)
func HasAccess(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			utils.Error(c, http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		var count int64
		err := database.DB.Table("user_permissions").
			Joins("JOIN permissions ON permissions.id = user_permissions.permission_id").
			Where("user_permissions.user_id = ? AND permissions.name IN ? AND permissions.deleted = 0", userID.(uint), permissions).
			Count(&count).Error

		if err != nil || count == 0 {
			utils.Error(c, http.StatusForbidden, "access denied: missing required permission")
			c.Abort()
			return
		}

		c.Next()
	}
}
