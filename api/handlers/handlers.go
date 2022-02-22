package handlers

import (
	"net/http"

	"github.com/sbxb/loyalty/config"
	"github.com/sbxb/loyalty/internal/logger"
)

// URLHandler defines a container for handlers and their dependencies
type URLHandler struct {
	//store  storage.Storage
	config config.Config
}

func NewURLHandler( /*st storage.Storage,*/ cfg config.Config) URLHandler {
	return URLHandler{
		//store:  st,
		config: cfg,
	}
}

// UserRegister process POST /api/user/register request
func (uh URLHandler) UserRegister(w http.ResponseWriter, r *http.Request) {
	logger.Info("UserRegister hit by POST /api/user/register")
}

// UserLogin process POST /api/user/login request
func (uh URLHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
	logger.Info("UserLogin hit by POST /api/user/login")
}
