package serializer

import "RIP2026/internal/app/models"

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
	Login       string `json:"login"`
	IsModerator bool   `json:"is_moderator"`
}

type UserJSON struct {
	ID          uint   `json:"id"`
	Login       string `json:"login"`
	Password    string `json:"password"`
	IsModerator bool   `json:"is_moderator"`
}

func SignUpResponseFromUser(user models.User) SignUpResponse {
	return SignUpResponse{
		Login:       user.Login,
		IsModerator: user.IsModerator,
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
		ID:          user.ID,
		Login:       user.Login,
		Password:    user.Password,
		IsModerator: user.IsModerator,
	}
}

func UserFromJSON(j UserJSON) models.User {
	return models.User{
		Login:       j.Login,
		Password:    j.Password,
		IsModerator: j.IsModerator,
	}
}
