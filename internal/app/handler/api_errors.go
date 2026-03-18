package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"RIP2026/internal/app/repository"
)

func (h *Handler) apiError(ctx *gin.Context, statusCode int, err error) {
	logrus.Error(err.Error())
	msg := err.Error()
	switch {
	case errors.Is(err, repository.ErrNotFound):
		msg = "Не найден"
	case errors.Is(err, repository.ErrAlreadyExists):
		msg = "Уже существует"
	case errors.Is(err, repository.ErrNotAllowed):
		msg = "Доступ запрещен"
	case errors.Is(err, repository.ErrNoDraft):
		msg = "Черновик не найден"
	}
	ctx.JSON(statusCode, gin.H{
		"status":      "error",
		"description": msg,
	})
}
