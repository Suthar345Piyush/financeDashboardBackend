// dashboard handler function

package handler

import (
	"net/http"
	"strconv"

	"github.com/suthar345Piyush/financedashboard/internal/service"
	"github.com/suthar345Piyush/financedashboard/pkg/response"
)

type DashboardHandler struct {
	dashSvc *service.DashboardService
}

func NewDashboardHandler(svc *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashSvc: svc}
}

func (h *DashboardHandler) Summary(w http.ResponseWriter, r *http.Request) {

	from := r.URL.Query().Get("date_from")
	to := r.URL.Query().Get("date_to")

	var (
		data any
		err  error
	)

	if from != "" && to != "" {
		data, err = h.dashSvc.GetSummaryByDateRange(from, to)
	} else {
		data, err = h.dashSvc.GetSummary()
	}

	if err != nil {
		response.InternalError(w)
		return
	}

	response.OK(w, data)
}

func (h *DashboardHandler) CategoryTotals(w http.ResponseWriter, r *http.Request) {
	data, err := h.dashSvc.GetCategoryTotals()

	if err != nil {
		response.InternalError(w)
		return
	}

	response.OK(w, data)
}

func (h *DashboardHandler) MonthlyTrends(w http.ResponseWriter, r *http.Request) {
	months := 6

	if m := r.URL.Query().Get("months"); m != "" {
		if v, err := strconv.Atoi(m); err == nil && v > 0 {
			months = v
		}
	}

	data, err := h.dashSvc.GetMonthlyTrends(months)

	if err != nil {
		response.InternalError(w)
		return
	}

	response.OK(w, data)
}

// recent activity

func (h *DashboardHandler) RecentActivity(w http.ResponseWriter, r *http.Request) {
	limit := 10

	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}

	data, err := h.dashSvc.GetRecentActivity(limit)

	if err != nil {
		response.InternalError(w)
		return
	}

	response.OK(w, data)
}
