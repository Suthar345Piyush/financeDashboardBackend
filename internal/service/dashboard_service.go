// dashboard service part logic

package service

import (
	"github.com/suthar345Piyush/financedashboard/internal/models"
	"github.com/suthar345Piyush/financedashboard/internal/repository"
)

type DashboardService struct {
	repo *repository.DashboardRepository
}

func NewDashboardService(repo *repository.DashboardRepository) *DashboardService {
	return &DashboardService{repo: repo}
}

func (s *DashboardService) GetSummary() (*models.DashboardSummary, error) {
	return s.repo.GetSummary()
}

func (s *DashboardService) GetCategoryTotals() ([]models.CategoryTotal, error) {
	return s.repo.GetCategoryTotals()
}

func (s *DashboardService) GetMonthlyTrends(months int) ([]models.MonthlyTrend, error) {
	return s.repo.GetMonthlyTrends(months)
}

func (s *DashboardService) GetRecentActivity(limit int) ([]models.FinancialRecord, error) {
	return s.repo.GetRecentActivity(limit)
}

func (s *DashboardService) GetSummaryByDateRange(from, to string) (*models.DashboardSummary, error) {
	return s.repo.GetSummaryByDateRange(from, to)
}
