package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/sbxb/loyalty/internal/logger"
	"github.com/sbxb/loyalty/models"
	"github.com/sbxb/loyalty/storage"
	"golang.org/x/crypto/bcrypt"
)

type contextKey string

var ContextUserKey = contextKey("user")

type AuthService struct {
	store storage.Storage
}

func NewAuthService(st storage.Storage) *AuthService {
	return &AuthService{store: st}
}

type AuthError struct {
	msg  string
	Code int
}

func (ae AuthError) Error() string {
	return ae.msg
}

func NewAuthError(msg string, code int) *AuthError {
	return &AuthError{
		msg:  msg,
		Code: code,
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func checkPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}

func (as *AuthService) RegisterUser(ctx context.Context, user *models.User) *AuthError {
	var err error

	if user.Hash, err = hashPassword(user.Password); err != nil {
		return NewAuthError("Crypto failed to process password", http.StatusInternalServerError)
	}

	if err = as.store.AddUser(ctx, user); err != nil {
		if errors.Is(err, storage.ErrLoginAlreadyExists) {
			return NewAuthError(err.Error(), http.StatusConflict)
		}
		return NewAuthError("Server failed to register user", http.StatusInternalServerError)
	}

	return nil
}

func (as *AuthService) LoginUser(ctx context.Context, user *models.User) (*models.User, *AuthError) {
	dbUser, err := as.store.GetUser(ctx, user)

	if err != nil {
		if errors.Is(err, storage.ErrLoginMissing) {
			return nil, NewAuthError("Wrong login and/or password", http.StatusUnauthorized)
		}
		// here could be server internal error when real db is used
		return nil, NewAuthError("Server failed to login user", http.StatusInternalServerError)
	}

	if !checkPassword(user.Password, dbUser.Hash) {
		return nil, NewAuthError("Wrong login and/or password", http.StatusUnauthorized)
	}

	return dbUser, nil
}

func (as *AuthService) SetCookie(w http.ResponseWriter, user *models.User) error {
	b, err := json.Marshal(models.UserAuth{Login: user.Login, ID: user.ID})
	if err != nil {
		return NewAuthError("Auth: SetCookie: JSON.Marshal() failed to serialize user", http.StatusInternalServerError)
	}

	encUserInfo, err := encryptString(string(b), secretKey)
	if err != nil {
		return NewAuthError("Auth: SetCookie: encryptString failed to encrypt user", http.StatusInternalServerError)
	}

	signedEncMessage := GetSignedString(encUserInfo, signatureKey)
	cookie := http.Cookie{
		Name:    "user",
		Value:   signedEncMessage,
		Expires: time.Now().Add(1 * time.Hour),
	}
	http.SetCookie(w, &cookie)

	return nil
}

func GetUserID(ctx context.Context) int {
	UserID, _ := ctx.Value(ContextUserKey).(int)
	if UserID == 0 {
		logger.Warning("User ID not found, check if authMW middleware was enabled")
	}

	return UserID
}
