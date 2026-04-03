/*

repository contains the data access of the application

this contains all the data about our user

this repository section has the access to directly talk to our sqlite database

*/

package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/suthar345Piyush/financedashboard/internal/models"
)

var ErrNotFound = errors.New("record not found")
var ErrDuplicateEmail = errors.New("email already exists")

type UserRepository struct {
	db *sql.DB
}

// new user repository

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// crud operation implementation

// creating user requires model of that user

func (r *UserRepository) Create(user *models.User) error {

	// sql query for user

	query := `INSERT INTO users (id, email, password, role, is_active, created_at, updated_at)  VALUES (?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()

	// executing the query

	_, err := r.db.Exec(query, user.ID, user.Email, user.Password, user.Role, user.IsActive, now, now)

	if err != nil {
		if isUniqueConstraintError(err) {
			return ErrDuplicateEmail
		}

		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

// unique constraint error function

func isUniqueConstraintError(err error) bool {

	return err != nil && (contains(err.Error(), "UNIQUE constraint failed") || contains(err.Error(), "unique constraint"))

}

//getting user , using email

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {

	// query

	query := `SELECT id, email, password, role, is_active, created_at, updated_at FROM users WHERE email = ?`

	return r.scanUser(r.db.QueryRow(query, email))
}

// getting user by ID

func (r *UserRepository) GetByID(id string) (*models.User, error) {

	query := `SELECT id, email, password, role, is_active, created_at, updated_at FROM users WHERE id = ?`

	return r.scanUser(r.db.QueryRow(query, id))

}

// listing down all the user , this will return a slice of all the users

func (r *UserRepository) List() ([]models.User, error) {

	query := `SELECT id, email, password, role, is_active, created_at, updated_at FROM users ORDER BY created_at DESC`

	// executing the query  , and it will return the row / getting rows

	rows, err := r.db.Query(query)

	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	defer rows.Close()

	var users []models.User

	for rows.Next() {
		u, err := r.scanUserRow(rows)

		if err != nil {
			return nil, err
		}

		users = append(users, *u)

	}

	return users, rows.Err()
}

// function for updating the role

func (r *UserRepository) UpdateRole(id string, role models.Role) error {

	// query

	query := `UPDATE users SET role = ?, updated_at = ? WHERE id = ?`

	result, err := r.db.Exec(query, role, time.Now(), id)

	if err != nil {
		return fmt.Errorf("update role: %w", err)
	}

	return checkRowsAffected(result)

}

// function to update the status

func (r *UserRepository) updateStatus(id string, isActive bool) error {

	query := `UPDATE users SET is_active = ?, updated_at = ? WHERE id = ?`

	result, err := r.db.Exec(query, isActive, time.Now(), id)

	if err != nil {
		return fmt.Errorf("update status: %w", err)
	}

	return checkRowsAffected(result)

}

// ------------ some helper functions ------------------------

// scans user and returns users

func (r *UserRepository) scanUser(row *sql.Row) (*models.User, error) {

	u := &models.User{}

	err := row.Scan(&u.ID, &u.Email, &u.Password, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("scan user: %w", err)

	}

	return u, nil

}

// scan user row function

func (r *UserRepository) scanUserRow(rows *sql.Rows) (*models.User, error) {

	u := &models.User{}
	err := rows.Scan(&u.ID, &u.Email, &u.Password, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("scan user row: %w", err)
	}

	return u, nil

}

// no. of rows affected  function
// this function , used for checking which rows in db are affected by update, insert and delete operations

//this Result is basically summerizer here, it summerizes the executed sql command

func checkRowsAffected(result sql.Result) error {
	n, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if n == 0 {
		return ErrNotFound
	}

	return nil
}

// contains function for checking some specific string contains the substring or not

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

// contain string function

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
