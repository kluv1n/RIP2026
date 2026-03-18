package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"RIP2026/internal/app/repository"
	"RIP2026/internal/app/serializer"
)

func (h *Handler) APICreateUser(ctx *gin.Context) {
	var j serializer.UserJSON
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
	ctx.JSON(http.StatusCreated, serializer.UserToJSON(user))
}

func (h *Handler) APISignIn(ctx *gin.Context) {
	var j serializer.UserJSON
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
	user, err := h.Repository.SignInAPI(j)
	if err != nil {
		if err.Error() == "неверный логин или пароль" {
			h.apiError(ctx, http.StatusUnauthorized, fmt.Errorf("invalid login or password"))
		} else {
			h.apiError(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	ctx.JSON(http.StatusOK, serializer.UserToJSON(user))
}

func (h *Handler) APISignOut(ctx *gin.Context) {
	h.Repository.SignOut()
	ctx.JSON(http.StatusOK, gin.H{
		"status": "signed_out",
	})
}
