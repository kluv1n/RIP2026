package serializer

import "RIP2026/internal/app/models"

// UserRole возвращает строковую роль для API и JWT.
func UserRole(isModerator bool) string {
	if isModerator {
		return "moderator"
	}
	return "creator"
}

// SignInRequest — запрос на вход.
type SignInRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// SignUpRequest — запрос на регистрацию.
type SignUpRequest struct {
	Login       string `json:"login" binding:"required"`
	Password    string `json:"password" binding:"required"`
	IsModerator bool   `json:"is_moderator"`
}

// SignUpResponse — ответ после регистрации без id и пароля.
type SignUpResponse struct {
	Login string `json:"login"`
	Role  string `json:"role"` // moderator | creator
}

// SignInResponse — ответ после входа.
type SignInResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"` // moderator | creator
}

type UserJSON struct {
	ID       uint   `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Role     string `json:"role"` // moderator | creator
}

func SignUpResponseFromUser(user models.User) SignUpResponse {
	return SignUpResponse{
		Login: user.Login,
		Role:  UserRole(user.IsModerator),
	}
}

func SignUpRequestToUser(j SignUpRequest) models.User {
	return models.User{
		Login:       j.Login,
		Password:    j.Password,
		IsModerator: j.IsModerator,
	}
}

func UserToJSON(user models.User) UserJSON {
	return UserJSON{
		ID:       user.ID,
		Login:    user.Login,
		Password: user.Password,
		Role:     UserRole(user.IsModerator),
	}
}

func UserFromJSON(j UserJSON) models.User {
	return models.User{
		Login:       j.Login,
		Password:    j.Password,
		IsModerator: j.Role == "moderator",
	}
}
