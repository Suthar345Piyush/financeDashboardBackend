// repository contains the dashboard functions

package repository

import (
	"database/sql"
	"fmt"

	"github.com/suthar345Piyush/financedashboard/internal/models"
)

type DashboardRepository struct {
	db *sql.DB
}

func NewDashboardRepository(db *sql.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

// get summary on the dashboard function

func (r *DashboardRepository) GetSummary() (*models.DashboardSummary, error) {

	query := `SELECT COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) AS total_income,
	   
	    COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) AS total_expenses,

			COUNT(*) AS record_count FROM financial_records`

	summary := &models.DashboardSummary{}

	err := r.db.QueryRow(query).Scan(
		&summary.TotalIncome,
		&summary.TotalExpenses,
		&summary.RecordCount,
	)

	if err != nil {
		return nil, fmt.Errorf("get summary: %w", err)
	}

	summary.NetBalance = summary.TotalIncome - summary.TotalExpenses

	return summary, err
}

// category total function

func (r *DashboardRepository) GetCategoryTotals() ([]models.CategoryTotal, error) {

	// get category total sql  query

	query := `SELECT category, type, COALESCE(SUM(amount), 0) AS total, COUNT(*) as count  FROM financial_records GROUP BY category, type ORDER BY total DESC`

	rows, err := r.db.Query(query)

	if err != nil {
		return nil, fmt.Errorf("category totals: %w", err)
	}

	defer rows.Close()

	var results []models.CategoryTotal

	for rows.Next() {

		var ct models.CategoryTotal

		if err := rows.Scan(&ct.Category, &ct.Type, &ct.Total, &ct.Count); err != nil {
			return nil, err
		}

		results = append(results, ct)
	}

	return results, rows.Err()
}

// monthly trends on the dashboard

func (r *DashboardRepository) GetMonthlyTrends(months int) ([]models.MonthlyTrend, error) {

	query := `SELECT strftime('%Y-%m', date) AS month, COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) AS income, 
	  
	   COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) AS expense

		 FROM financial_records
		 WHERE date >= date('now', ? || 'months')
		 GROUP BY month
		 ORDER BY ASC
	 `

	rows, err := r.db.Query(query, fmt.Sprintf("-%d", months))

	if err != nil {
		return nil, fmt.Errorf("monthly trends: %w", err)
	}

	defer rows.Close()

	var trends []models.MonthlyTrend

	for rows.Next() {
		var t models.MonthlyTrend

		if err := rows.Scan(&t.Month, &t.Income, &t.Expense); err != nil {
			return nil, err
		}

		t.Net = t.Income - t.Expense

		trends = append(trends, t)
	}

	return trends, rows.Err()
}

// function for recent activity

func (r *DashboardRepository) GetRecentActivity(limit int) ([]models.FinancialRecord, error) {

	query := `SELECT id, user_id, amount, type, category, date, notes, created_at, updated_at FROM financial_records ORDER BY created_at DESC LIMIT ?`

	rows, err := r.db.Query(query, limit)

	if err != nil {
		return nil, fmt.Errorf("recent activity: %w", err)
	}

	defer rows.Close()

	repo := &RecordRepository{db: r.db}

	var records []models.FinancialRecord

	for rows.Next() {
		rec, err := repo.scanRecordRow(rows)

		if err != nil {
			return nil, err
		}

		records = append(records, *rec)
	}

	return records, rows.Err()
}

// function to get summary on specific date (filter)
// getting this from dashboard summary

func (r *DashboardRepository) GetSummaryByDateRange(from, to string) (*models.DashboardSummary, error) {

	query := `SELECT 

	    COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0),
			COUNT(*)  FROM financial_records WHERE date BETWEEN ? AND ?
	`

	summary := &models.DashboardSummary{}

	err := r.db.QueryRow(query, from, to).Scan(&summary.TotalIncome, &summary.TotalExpenses, &summary.RecordCount)

	if err != nil {
		return nil, fmt.Errorf("summary by date range: %w", err)
	}

	summary.NetBalance = summary.TotalIncome - summary.TotalExpenses

	return summary, nil
}

// helper function to reuse from the record repository

// func nullableString(ns sql.NullString) string {
// 	if ns.Valid {
// 		return ns.String
// 	}
// 	return ""
// }
