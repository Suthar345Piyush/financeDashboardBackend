// record respository

package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/suthar345Piyush/financedashboard/internal/models"
)

type RecordRepository struct {
	db *sql.DB
}

// function for new record repository

func NewRecordRepository(db *sql.DB) *RecordRepository {
	return &RecordRepository{db: db}
}

// create function for record

func (r *RecordRepository) Create(rec *models.FinancialRecord) error {

	// query to run on financial_records table

	query := `INSERT INTO financial_records (id, user_id, amount, type, category, date, notes, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()

	_, err := r.db.Exec(query, rec.ID, rec.UserID, rec.Amount, rec.Type, rec.Category, rec.Date, rec.Notes, now, now)

	if err != nil {
		return fmt.Errorf("create record: %w", err)
	}

	return nil

}

// get by ID function

func (r *RecordRepository) GetByID(id string) (*models.FinancialRecord, error) {

	query := `SELECT id, user_id, amount, type, category, date, notes, created_at, updated_at FROM financial_records WHERE id = ?`

	row := r.db.QueryRow(query, id)

	return r.scanRecord(row)
}

// listing all the records (int -> total records)

func (r *RecordRepository) List(filter models.RecordFilter) ([]models.FinancialRecord, int, error) {

	// dynamic where clause

	conditions := []string{}
	args := []any{}

	if filter.Type != "" {
		conditions = append(conditions, "type = ?")
		args = append(args, filter.Type)
	}

	if filter.Category != "" {
		conditions = append(conditions, "category = ?")
		args = append(args, filter.Category)
	}

	if filter.DateFrom != "" {
		conditions = append(conditions, "date >= ?")
		args = append(args, filter.DateFrom)
	}

	if filter.DateTo != "" {
		conditions = append(conditions, "date <= ?")
		args = append(args, filter.DateTo)
	}

	where := ""

	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	// counting the total no of records for pagination

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM financial_records %s", where)

	var total int

	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count records: %w", err)
	}

	// for pagination

	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	if filter.Page <= 0 {
		filter.Page = 1
	}

	offset := (filter.Page - 1) * filter.PageSize

	dataQuery := fmt.Sprintf(`SELECT id, user_id, amount, type, category, date, notes, created_at, updated_at FROM financial_records %s ORDER BY date DESC, created_at DESC LIMIT ? OFFSET ?`, where)

	args = append(args, filter.PageSize, offset)

	rows, err := r.db.Query(dataQuery, args...)

	if err != nil {
		return nil, 0, fmt.Errorf("list records: %w", err)
	}

	defer rows.Close()

	var records []models.FinancialRecord

	for rows.Next() {
		rec, err := r.scanRecordRow(rows)

		if err != nil {
			return nil, 0, err
		}

		records = append(records, *rec)
	}

	return records, total, rows.Err()

}

// update the record

func (r *RecordRepository) Update(rec *models.FinancialRecord) error {

	// query to perform update

	query := `UPDATE financial_records SET amount=?, type=?, category=?, date=?, notes=?, updated_at=? WHERE id=?`

	result, err := r.db.Exec(query, rec.Amount, rec.Type, rec.Category, rec.Date, rec.Notes, time.Now(), rec.ID)

	if err != nil {
		return fmt.Errorf("update record: %w", err)
	}

	return checkRowsAffected(result)
}

// delete the record

func (r *RecordRepository) Delete(id string) error {

	result, err := r.db.Exec("DELETE FROM financial_records WHERE id=?", id)

	if err != nil {
		return fmt.Errorf("delete record: %w", err)
	}

	return checkRowsAffected(result)
}

// some helper function like scanRecord and scanRecordRow

func (r *RecordRepository) scanRecord(row *sql.Row) (*models.FinancialRecord, error) {

	rec := &models.FinancialRecord{}

	err := row.Scan(&rec.ID, &rec.UserID, &rec.Amount, &rec.Type, &rec.Category, &rec.Date, &rec.Notes, &rec.CreatedAt, &rec.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("scan record: %w", err)
	}

	return rec, nil
}

// function for scanRecordRow

func (r *RecordRepository) scanRecordRow(rows *sql.Rows) (*models.FinancialRecord, error) {

	rec := &models.FinancialRecord{}
	err := rows.Scan(&rec.ID, &rec.UserID, &rec.Amount, &rec.Type, &rec.Category, &rec.Date, &rec.Notes, &rec.CreatedAt, &rec.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("scan record row: %w", err)
	}

	return rec, nil
}
