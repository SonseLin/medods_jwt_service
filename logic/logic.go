package logic

import (
	"fmt"
	model "medods_jwt_service/model"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Секреты для подписи Access токенов
var accessSecret = []byte("52 52 52")

// isValidUser checks if the provided user object has valid data.
//
// The function validates the following fields of the user object:
// - Email: Must not be an empty string.
// - Name: Must not be an empty string.
// - GUID: Must be greater than 0.
//
// Parameters:
// - u: A model.User object to be validated.
//
// Returns:
// - A boolean value indicating whether the user object is valid or not.
func isValidUser(u model.User) bool {
	return u.Email != "" && u.Name != "" && u.GUID > 0
}

/**
 * JWT_generator generates a JSON Web Token (JWT) for the provided user and session ID.
 *
 * @param user A model.User object representing the user for whom the token is being generated.
 * @param sessionID A string representing the session ID associated with the user.
 *
 * @return tokenString A string representing the generated JWT token.
 * @return err An error, if any, encountered during the token generation process.
 */
func JWT_generator(user model.User, sessionID string) (tokenString string, err error) {
	if !isValidUser(user) {
		return "", model.NewError("invalid user", "empty fields")
	}
	payload := model.Payload{
		UserID:    user.GUID,
		Email:     user.Email,
		SessionID: sessionID,
		IP:        user.IP,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "medods_jwt_service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, payload)

	tokenString, err = token.SignedString(accessSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

/**
 * CheckIfPossibleGetNewTokens checks if it's possible to generate new tokens for the provided user and refresh token.
 *
 * This function validates if the refresh token has not been used yet and if it matches the hash of the user's refresh token.
 * If both conditions are met, it calls the GenerateTokens function to generate new tokens.
 *
 * @param user A model.User object representing the user for whom the tokens are being generated.
 * @param sessionID A string representing the session ID associated with the user.
 * @param refreshToken A model.Refresh_token object representing the refresh token to be validated.
 *
 * @return tokenString1, tokenString2, tokenString3 A slice of strings representing the generated JWT token, the new refresh token, and the hashed refresh token, respectively.
 * @return err An error, if any, encountered during the token generation process.
 */
func CheckIfPossibleGetNewTokens(user model.User, sessionID string, refreshToken model.Refresh_token, current_IP string) (string, string, string, error) {
	if refreshToken.Times_to_use < 1 {
		return "", "", "", model.NewError("token issue", "this key was already in use")
	}
	if user.Token.JWT_token.HashRefreshToken == refreshToken.Hash {
		return GenerateTokens(user, sessionID, current_IP)
	}
	return "", "", "", model.NewError("token issue", "this jwt token wasnt generated in pair with refresh token")
}

/**
 * RefreshToken_generator generates a new refresh token for the provided user.
 *
 * This function creates a new refresh token by concatenating the user's creation timestamp, user ID, and the current time.
 * The refresh token is then returned as a string.
 *
 * @param user A model.User object representing the user for whom the refresh token is being generated.
 *
 * @return token A string representing the generated refresh token.
 */
func RefreshToken_generator(user model.User) string {
	created := user.Created_at
	id := user.GUID
	new_str := created.String() + fmt.Sprintf("%d", id) + time.Now().String()
	return new_str
}

/**
 * GetHashRefreshToken generates a bcrypt hash of the provided refresh token.
 *
 * This function takes a refresh token as input and returns a bcrypt hash of the token along with any encountered errors.
 * The bcrypt.GenerateFromPassword function is used to generate the hash with a default cost factor.
 *
 * @param token A string representing the refresh token to be hashed.
 *
 * @return hash A string representing the bcrypt hash of the refresh token.
 * @return err An error, if any, encountered during the hash generation process.
 */
func GetHashRefreshToken(token string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

/**
 * GenerateTokens generates a JWT token and a refresh token for the provided user and session ID.
 *
 * This function creates a new JWT token and a new refresh token for the given user. It uses the JWT_generator function to generate the JWT token and the RefreshToken_generator function to generate the refresh token. The refresh token is then hashed using the GetHashRefreshToken function.
 *
 * The function takes a user object, an access key (which is not used in this implementation), and the user's IP address as input. It associates the generated tokens with the provided session ID and hashed refresh token, and then assigns them to the user's Token field.
 *
 * @param user A model.User object representing the user for whom the tokens are being generated.
 * @param accessKey A string representing the access key associated with the user. This parameter is not used in this implementation.
 * @param IP A string representing the IP address of the user.
 *
 * @return jwtToken A string representing the generated JWT token.
 * @return refreshToken A string representing the generated refresh token.
 * @return hashedRefresh A string representing the bcrypt hash of the refresh token.
 * @return err An error, if any, encountered during the token generation process.
 */
func GenerateTokens(user model.User, accessKey string, IP string) (string, string, string, error) {
	sessionID := time.Now().String() + time.October.String()
	jwtToken, err := JWT_generator(user, sessionID)
	if err != nil {
		return "", "", "", err
	}
	refreshToken := RefreshToken_generator(user)
	hashedRefresh, err := GetHashRefreshToken(refreshToken)
	if err != nil {
		return "", "", "", err
	}
	if IP != user.IP {
		IPChangeIssue(user, IP)
	}
	return jwtToken, refreshToken, hashedRefresh, nil
}

/**
 * connectTokensToUser connects the provided tokens to the user object.
 *
 * This function takes a user object, a JWT token, a refresh token, a hashed refresh token, and a session ID as input.
 * It creates a new JWT token object and a new refresh token object, associates them with the provided session ID and hashed refresh token, and then assigns them to the user's Token field.
 *
 * @param user A pointer to a model.User object representing the user for whom the tokens are being connected.
 * @param jwt A string representing the generated JWT token.
 * @param refresh A string representing the generated refresh token.
 * @param hashedRefresh A string representing the bcrypt hash of the refresh token.
 * @param sessionID A string representing the session ID associated with the user.
 */
func connectTokensToUser(user *model.User, jwt, refresh, hashedRefresh, sessionID string) {
	jwt_token := model.JWT_token{}
	jwt_token.SessionID = sessionID
	jwt_token.Value = jwt
	jwt_token.HashRefreshToken = hashedRefresh

	refresh_token := model.Refresh_token{}
	refresh_token.SessionID = sessionID
	refresh_token.Value = refresh
	refresh_token.Hash = hashedRefresh
	refresh_token.Times_to_use = 1

	token_pair := model.TokenPair{}
	token_pair.JWT_token = jwt_token
	token_pair.Refresh_token = refresh_token

	user.Token = token_pair
}
