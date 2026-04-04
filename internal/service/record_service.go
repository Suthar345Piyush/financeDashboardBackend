// logic for record service

package service

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/suthar345Piyush/financedashboard/internal/models"
	"github.com/suthar345Piyush/financedashboard/internal/repository"
)

var ErrRecordNotFound = errors.New("record not found")

type RecordService struct {
	repo *repository.RecordRepository
}

func NewRecordService(repo *repository.RecordRepository) *RecordService {
	return &RecordService{repo: repo}
}

func (s *RecordService) Create(userID string, req CreateRecordRequest) (*models.FinancialRecord, error) {

	rec := &models.FinancialRecord{
		ID:       uuid.NewString(),
		UserID:   userID,
		Amount:   req.Amount,
		Type:     models.RecordType(req.Type),
		Category: req.Category,
		Date:     req.Date,
		Notes:    req.Notes,
	}

	if err := s.repo.Create(rec); err != nil {
		return nil, fmt.Errorf("create record: %w", err)
	}

	return rec, nil
}

type CreateRecordRequest struct {
	Amount   string
	Type     string
	Category string
	Date     string
	Notes    string
}

// listing all the records function

func (s *RecordService) List(filter models.RecordFilter) ([]models.FinancialRecord, int, error) {
	return s.repo.List(filter)
}

// getting the record by id

func (s *RecordService) GetByID(id string) (*models.FinancialRecord, error) {
	rec, err := s.repo.GetByID(id)

	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrRecordNotFound
	}

	return rec, err
}

// update record function

func (s *RecordService) Update(id string, req UpdateRecordRequest) (*models.FinancialRecord, error) {

	// taking the record which is going to update

	rec, err := s.GetByID(id)

	if err != nil {
		return nil, err
	}

	rec.Amount = req.Amount
	rec.Type = models.RecordType(req.Type)
	rec.Category = req.Category
	rec.Date = req.Date
	rec.Notes = req.Notes

	if err = s.repo.Update(rec); err != nil {
		return nil, fmt.Errorf("update record: %w", err)
	}

	return rec, nil

}

type UpdateRecordRequest = CreateRecordRequest

// delete record function

func (s *RecordService) Delete(id string) error {
	err := s.repo.Delete(id)

	if errors.Is(err, repository.ErrNotFound) {
		return ErrRecordNotFound
	}
	return err
}
