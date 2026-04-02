// structs for financial records
// considering Income and Expenses as record type initially

package models

import "time"

type RecordType string

const (
	RecordTypeIncome  RecordType = "income"
	RecordTypeExpense RecordType = "expense"
)

// validation function , same what we did in user models

func (t RecordType) IsValid() bool {
	return t == RecordTypeIncome || t == RecordTypeExpense
}

// struct for financial records

type FinancialRecord struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	Amount    string     `json:"amount"`
	Type      RecordType `json:"type"`
	Category  string     `json:"category"`
	Date      string     `json:"date"`
	Notes     string     `json:"notes,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// record filter based on category , date , type ( they hold query params for listing/filtering records )

type RecordFilter struct {
	Type     string
	Category string
	DateFrom string
	DateTo   string
	Page     int
	PageSize int
}

// some dashboard related models/ structs

type DashboardSummary struct {
	TotalIncome   float64 `json:"total_income"`
	TotalExpenses float64 `json:"total_expenses"`
	NetBalance    float64 `json:"net_balance"`
	RecordCount   int     `json:"record_count"`
}

// total in category

type CategoryTotal struct {
	Category string  `json:"category"`
	Type     string  `json:"type"`
	Total    float64 `json:"total"`
	Count    int     `json:"count"`
}

// charts or monthly trends  struct

type MonthlyTrend struct {
	Month   string  `json:"month"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	Net     float64 `json:"net"`
}

// recent activity

type RecentActivity struct {
	Records []FinancialRecord `json:"records"`
	Total   int               `json:"total"`
}
