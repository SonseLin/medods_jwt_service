package logic

import (
	"database/sql"
	"medods_jwt_service/model"
)

func GetUserById(id int, db *sql.DB) (model.User, error) {
	if id < 0 {
		return model.User{}, model.NewError("invalid id", "id can't be lower than 0")
	}

	query := "SELECT name, ip, id_int, email, created_at FROM users WHERE id_int = $1"
	row := db.QueryRow(query, id)

	user := model.User{}
	err := row.Scan(&user.Name, &user.IP, &user.GUID, &user.Email, &user.Created_at)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, model.NewError("not found", "no user found with the provided ID")
		}
		return model.User{}, err
	}

	return user, nil
}

func GetUserToken(user model.User) model.TokenPair {
	return user.Token
}

func UpdateUserTokens(user model.User, sessionID, current_IP string) {
	jwt, refresh, hash, err := CheckIfPossibleGetNewTokens(user, sessionID, GetUserToken(user).Refresh_token, current_IP)
	if err != nil {
		return
	}
	connectTokensToUser(&user, jwt, refresh, hash, sessionID)
}

func IPChangeIssue(user model.User, current_IP string) {
	user.IP = current_IP
	SendEmail(user, EmailTemplateMessage(user))
}
