package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"RIP2026/internal/app/repository"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const bearerPrefix = "Bearer"

func extractTokenFromRequest(c *gin.Context) string {
	bearerToken := c.GetHeader("Authorization")
	if bearerToken != "" {
		parts := strings.SplitN(bearerToken, " ", 2)
		if len(parts) == 2 && parts[0] == bearerPrefix {
			return parts[1]
		}
	}
	cookieToken, err := c.Cookie("token")
	if err == nil && cookieToken != "" {
		return cookieToken
	}
	return ""
}

func parseToken(tokenString string) (jwt.MapClaims, error) {
	jwtKey := os.Getenv("JWT_KEY")
	if jwtKey == "" {
		jwtKey = "default-secret-key-change-in-production"
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(jwtKey), nil
	})
	if err != nil || token == nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}

func userIDFromClaims(claims jwt.MapClaims) (uint, error) {
	userIDStr, ok := claims["user_id"].(string)
	if !ok || userIDStr == "" {
		return 0, fmt.Errorf("user_id not found")
	}
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(userID), nil
}

// moderatorFromClaims: приоритет role; fallback на is_moderator для старых JWT.
func moderatorFromClaims(claims jwt.MapClaims) bool {
	if role, ok := claims["role"].(string); ok {
		return role == "moderator"
	}
	if b, ok := claims["is_moderator"].(bool); ok {
		return b
	}
	return false
}

func (h *Handler) OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractTokenFromRequest(c)
		if tokenString == "" {
			c.Next()
			return
		}
		blacklisted, err := h.Repository.IsTokenBlacklisted(context.Background(), tokenString)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if blacklisted {
			c.Next()
			return
		}
		claims, err := parseToken(tokenString)
		if err != nil {
			c.Next()
			return
		}
		userID, err := userIDFromClaims(claims)
		if err != nil {
			c.Next()
			return
		}
		isModerator := moderatorFromClaims(claims)
		c.Set("user_id", userID)
		c.Set("is_moderator", isModerator)
		c.Next()
	}
}

func (h *Handler) AuthMiddleware(requireModerator bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractTokenFromRequest(c)
		if tokenString == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		blacklisted, err := h.Repository.IsTokenBlacklisted(context.Background(), tokenString)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if blacklisted {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		claims, err := parseToken(tokenString)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		userID, err := userIDFromClaims(claims)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		isModerator := moderatorFromClaims(claims)
		if requireModerator && !isModerator {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Set("user_id", userID)
		c.Set("is_moderator", isModerator)
		c.Next()
	}
}

func currentUserID(c *gin.Context) (uint, error) {
	v, ok := c.Get("user_id")
	if !ok {
		return 0, repository.ErrNotAllowed
	}
	userID, ok := v.(uint)
	if !ok || userID == 0 {
		return 0, repository.ErrNotAllowed
	}
	return userID, nil
}
