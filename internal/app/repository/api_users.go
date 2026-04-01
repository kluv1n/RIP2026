package repository

import (
	"errors"
	"fmt"
	"os"
	"time"

	"RIP2026/internal/app/models"
	"RIP2026/internal/app/serializer"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

func (r *Repository) GetUserByID(id int) (models.User, error) {
	if id <= 0 {
		return models.User{}, fmt.Errorf("неверный id: должен быть > 0")
	}
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *Repository) GetUserByLogin(login string) (models.User, error) {
	if login == "" {
		return models.User{}, errors.New("логин не может быть пустым")
	}
	var user models.User
	err := r.db.Where("login = ?", login).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, fmt.Errorf("%w: пользователь с логином %s не найден", ErrNotFound, login)
		}
		return models.User{}, err
	}
	return user, nil
}

func (r *Repository) CreateUserAPI(j serializer.SignUpRequest) (models.User, error) {
	user := serializer.SignUpRequestToUser(j)
	if user.Login == "" {
		return models.User{}, errors.New("логин обязателен для заполнения")
	}
	if user.Password == "" {
		return models.User{}, errors.New("пароль обязателен для заполнения")
	}
	_, err := r.GetUserByLogin(user.Login)
	if err == nil {
		return models.User{}, fmt.Errorf("%w: пользователь с логином %s уже существует", ErrAlreadyExists, user.Login)
	}
	if !errors.Is(err, ErrNotFound) {
		return models.User{}, err
	}
	if err := r.db.Create(&user).Error; err != nil {
		return models.User{}, fmt.Errorf("ошибка при создании пользователя: %w", err)
	}
	return user, nil
}

func (r *Repository) SignInAPI(j serializer.SignInRequest) (models.User, string, error) {
	if j.Login == "" {
		return models.User{}, "", errors.New("логин обязателен для заполнения")
	}
	if j.Password == "" {
		return models.User{}, "", errors.New("пароль обязателен для заполнения")
	}
	user, err := r.GetUserByLogin(j.Login)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return models.User{}, "", errors.New("неверный логин или пароль")
		}
		return models.User{}, "", err
	}
	if user.Password != j.Password {
		return models.User{}, "", errors.New("неверный логин или пароль")
	}
	r.SetUserID(int(user.ID))
	token, err := GenerateToken(user.ID, user.IsModerator)
	if err != nil {
		return models.User{}, "", err
	}
	return user, token, nil
}

func GenerateToken(userID uint, isModerator bool) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["user_id"] = fmt.Sprintf("%d", userID)
	claims["is_moderator"] = isModerator
	claims["exp"] = time.Now().Add(time.Hour).Unix()

	jwtKey := os.Getenv("JWT_KEY")
	if jwtKey == "" {
		jwtKey = "default-secret-key-change-in-production"
	}
	return token.SignedString([]byte(jwtKey))
}
