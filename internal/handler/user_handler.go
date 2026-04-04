// user handler function

package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/suthar345Piyush/financedashboard/internal/models"
	"github.com/suthar345Piyush/financedashboard/internal/repository"
	"github.com/suthar345Piyush/financedashboard/internal/service"
	"github.com/suthar345Piyush/financedashboard/pkg/response"
	"github.com/suthar345Piyush/financedashboard/pkg/validator"
)

type UserHandler struct {
	userSvc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{userSvc: svc}
}

/// listing all the users

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {

	users, err := h.userSvc.List()

	if err != nil {
		response.InternalError(w)
		return
	}

	out := make([]models.UserResponse, 0, len(users))

	for _, u := range users {
		out = append(out, u.ToResponse())
	}

	response.OK(w, out)

}

// getting user by id

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	user, err := h.userSvc.GetByID(id)

	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.NotFound(w, "user")
			return
		}
		response.InternalError(w)
		return
	}

	response.OK(w, user.ToResponse())
}

// role update function

func (h *UserHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var body struct {
		Role string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.BadRequest(w, "invalid JSON")
		return
	}

	v := validator.New()

	v.OneOf(body.Role, "role", "viewer", "analyst", "admin")

	if !v.Valid() {
		response.ValidationError(w, v.Errors)
		return
	}

	if err := h.userSvc.UpdateRole(id, models.Role(body.Role)); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.NotFound(w, "user")
			return
		}

		response.InternalError(w)
		return
	}

	response.OK(w, map[string]string{"message": "role updated"})

}

// final update status  function

func (h *UserHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	var body struct {
		IsActive bool `json:"is_active"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.BadRequest(w, "invalid JSON")
		return
	}

	if err := h.userSvc.UpdateStatus(id, body.IsActive); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.NotFound(w, "user")
			return
		}

		response.InternalError(w)
		return
	}

	response.OK(w, map[string]string{"message": "status updated"})
}
