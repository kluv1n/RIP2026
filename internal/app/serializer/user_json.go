package serializer

import "RIP2026/internal/app/models"

type UserJSON struct {
	ID          uint   `json:"id"`
	Login       string `json:"login"`
	Password    string `json:"password"`
	IsModerator bool   `json:"is_moderator"`
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
