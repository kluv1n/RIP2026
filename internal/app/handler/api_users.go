package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"RIP2026/internal/app/repository"
	"RIP2026/internal/app/serializer"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// APICreateUser godoc
// @Summary Регистрация пользователя
// @Description Регистрирует нового пользователя. Возвращает login и role (moderator|creator).
// @Tags battery life
// @Accept json
// @Produce json
// @Param user body serializer.SignUpRequest true "Логин и пароль"
// @Success 201 {object} serializer.SignUpResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /users/signup [post]
func (h *Handler) APICreateUser(ctx *gin.Context) {
	var j serializer.SignUpRequest
	if err := ctx.BindJSON(&j); err != nil {
		h.apiError(ctx, http.StatusBadRequest, err)
		return
	}
	if j.Login == "" {
		h.apiError(ctx, http.StatusBadRequest, fmt.Errorf("field 'login' is required"))
		return
	}
	if j.Password == "" {
		h.apiError(ctx, http.StatusBadRequest, fmt.Errorf("field 'password' is required"))
		return
	}
	user, err := h.Repository.CreateUserAPI(j)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			h.apiError(ctx, http.StatusConflict, err)
		} else {
			h.apiError(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	ctx.Header("Location", fmt.Sprintf("/api/users/%d", user.ID))
	ctx.JSON(http.StatusCreated, serializer.SignUpResponseFromUser(user))
}

// APISignIn godoc
// @Summary Вход (получение токена)
// @Description Принимает логин/пароль, возвращает token и role (moderator|creator).
// @Tags battery life
// @Accept json
// @Produce json
// @Param credentials body serializer.SignInRequest true "Логин и пароль"
// @Success 200 {object} serializer.SignInResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/signin [post]
func (h *Handler) APISignIn(ctx *gin.Context) {
	var j serializer.SignInRequest
	if err := ctx.BindJSON(&j); err != nil {
		h.apiError(ctx, http.StatusBadRequest, err)
		return
	}
	if j.Login == "" {
		h.apiError(ctx, http.StatusBadRequest, fmt.Errorf("field 'login' is required"))
		return
	}
	if j.Password == "" {
		h.apiError(ctx, http.StatusBadRequest, fmt.Errorf("field 'password' is required"))
		return
	}
	user, token, err := h.Repository.SignInAPI(j)
	if err != nil {
		if err.Error() == "неверный логин или пароль" {
			h.apiError(ctx, http.StatusUnauthorized, fmt.Errorf("invalid login or password"))
		} else {
			h.apiError(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	ctx.SetCookie("token", token, 3600, "/", "", false, true)
	ctx.JSON(http.StatusOK, serializer.SignInResponse{
		Token: token,
		Role:  serializer.UserRole(user.IsModerator),
	})
}

// APISignOut godoc
// @Summary Выход (удаление токена)
// @Description Удаляет токен текущего пользователя из blacklist. 204 No Content.
// @Tags battery life
// @Produce json
// @Success 204 "Токен добавлен в blacklist"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /users/signout [post]
func (h *Handler) APISignOut(ctx *gin.Context) {
	tokenString := extractTokenFromRequest(ctx)
	if tokenString != "" {
		if claims, err := parseToken(tokenString); err == nil {
			if ttl, err := tokenTTLFromClaims(claims); err == nil {
				_ = h.Repository.AddTokenToBlacklist(context.Background(), tokenString, ttl)
			}
		}
	}
	h.Repository.SignOut()
	ctx.SetCookie("token", "", -1, "/", "", false, true)
	ctx.Status(http.StatusNoContent)
}

func tokenTTLFromClaims(claims jwt.MapClaims) (time.Duration, error) {
	expVal, ok := claims["exp"]
	if !ok {
		return 0, errors.New("exp not present")
	}
	var expUnix int64
	switch v := expVal.(type) {
	case float64:
		expUnix = int64(v)
	case int64:
		expUnix = v
	case json.Number:
		i, err := v.Int64()
		if err != nil {
			return 0, err
		}
		expUnix = i
	default:
		return 0, fmt.Errorf("unsupported exp type %T", v)
	}
	ttl := time.Until(time.Unix(expUnix, 0))
	if ttl < 0 {
		return 0, errors.New("token already expired")
	}
	return ttl, nil
}
