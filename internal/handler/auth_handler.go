// handler functions (Auth)

// these handlers  are just HTTP layer, that further connects to the services (core business logic of the app)

package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/suthar345Piyush/financedashboard/internal/service"
	"github.com/suthar345Piyush/financedashboard/pkg/response"
	"github.com/suthar345Piyush/financedashboard/pkg/validator"
)

type AuthHandler struct {
	authSvc *service.AuthService
}

// function to make a new auth handler

func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

// register handler function

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

	// http request body

	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.BadRequest(w, "invalid json")
		return
	}

	v := validator.New()

	v.RequiredString(body.Email, "email")
	v.ValidEmail(body.Email, "email")
	v.RequiredString(body.Password, "password")
	v.Check(len(body.Password) >= 8, "password", "password must be at least 8 characters")

	if !v.Valid() {
		response.ValidationError(w, v.Errors)
		return
	}

	// time to register the user

	user, err := h.authSvc.Register(body.Email, body.Password)

	if err != nil {
		if errors.Is(err, service.ErrDuplicateEmail) {
			response.Error(w, http.StatusConflict, "email already registered")
			return
		}

		response.InternalError(w)
		return
	}

	response.Created(w, user.ToResponse())
}

// function for login

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.BadRequest(w, "invalid JSON")
		return
	}

	v := validator.New()
	v.RequiredString(body.Email, "email")
	v.RequiredString(body.Password, "password")

	if !v.Valid() {
		response.ValidationError(w, v.Errors)
		return
	}

	token, user, err := h.authSvc.Login(body.Email, body.Password)

	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			response.Error(w, http.StatusUnauthorized, "invalid email or password")
			return
		}

		if errors.Is(err, service.ErrUserInactive) {
			response.Error(w, http.StatusForbidden, "account is inactive")
			return
		}

		response.InternalError(w)
		return
	}

	response.OK(w, map[string]any{"token": token, "user": user.ToResponse()})

}
