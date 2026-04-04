// record handler function

package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/suthar345Piyush/financedashboard/internal/models"
	"github.com/suthar345Piyush/financedashboard/internal/service"
	"github.com/suthar345Piyush/financedashboard/middleware"
	"github.com/suthar345Piyush/financedashboard/pkg/response"
	"github.com/suthar345Piyush/financedashboard/pkg/validator"
)

type RecordHandler struct {
	recordSvc *service.RecordService
}

func NewRecordHandler(svc *service.RecordService) *RecordHandler {
	return &RecordHandler{recordSvc: svc}
}

func (h *RecordHandler) List(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()

	page, _ := strconv.Atoi(q.Get("page"))
	pageSize, _ := strconv.Atoi(q.Get("page_size"))

	filter := models.RecordFilter{
		Type:     q.Get("type"),
		Category: q.Get("category"),
		DateFrom: q.Get("date_from"),
		DateTo:   q.Get("date_to"),
		Page:     page,
		PageSize: pageSize,
	}

	records, total, err := h.recordSvc.List(filter)

	if err != nil {
		response.InternalError(w)
		return
	}

	if records == nil {
		records = []models.FinancialRecord{}
	}

	response.OK(w, map[string]any{
		"records":   records,
		"total":     total,
		"page":      filter.Page,
		"page_size": filter.PageSize,
	})
}

//  get by id function

func (h *RecordHandler) GetByID(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	rec, err := h.recordSvc.GetByID(id)

	if err != nil {
		if errors.Is(err, service.ErrRecordNotFound) {
			response.NotFound(w, "record")
			return
		}

		response.InternalError(w)
		return
	}

	response.OK(w, rec)
}

/// create function

func (h *RecordHandler) Create(w http.ResponseWriter, r *http.Request) {

	var body struct {
		Amount   string `json:"amount"`
		Type     string `json:"type"`
		Category string `json:"category"`
		Date     string `json:"date"`
		Notes    string `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.BadRequest(w, "invalid JSON")
		return
	}

	v := validator.New()
	v.RequiredString(body.Amount, "amount")
	v.OneOf(body.Type, "type", "income", "expense")
	v.RequiredString(body.Category, "category")
	v.RequiredString(body.Date, "date")
	v.ValidDate(body.Date, "date")

	if !v.Valid() {
		response.ValidationError(w, v.Errors)
		return
	}

	userID := middleware.GetUserID(r)

	rec, err := h.recordSvc.Create(userID, service.CreateRecordRequest{
		Amount: body.Amount, Type: body.Type, Category: body.Category, Date: body.Date, Notes: body.Notes,
	})

	if err != nil {
		response.InternalError(w)
		return
	}

	response.Created(w, rec)

}

// update function

func (h *RecordHandler) Update(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	var body struct {
		Amount   string `json:"amount"`
		Type     string `json:"type"`
		Category string `json:"category"`
		Date     string `json:"date"`
		Notes    string `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.BadRequest(w, "invalid JSON")
		return
	}

	v := validator.New()
	v.RequiredString(body.Amount, "amount")
	v.OneOf(body.Type, "type", "income", "expense")
	v.RequiredString(body.Category, "category")
	v.ValidDate(body.Date, "date")

	if !v.Valid() {
		response.ValidationError(w, v.Errors)
		return
	}

	rec, err := h.recordSvc.Update(id, service.UpdateRecordRequest{
		Amount: body.Amount, Type: body.Type, Category: body.Category, Date: body.Date, Notes: body.Notes,
	})

	if err != nil {
		if errors.Is(err, service.ErrRecordNotFound) {
			response.NotFound(w, "record")
			return
		}

		response.InternalError(w)
		return
	}

	response.OK(w, rec)
}

// delete function

func (h *RecordHandler) Delete(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	if err := h.recordSvc.Delete(id); err != nil {
		if errors.Is(err, service.ErrRecordNotFound) {
			response.NotFound(w, "record")
			return
		}

		response.InternalError(w)
		return
	}

	response.NoContent(w)

}
