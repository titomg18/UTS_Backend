package repository

import (
	"crud-app/app/models"
	"database/sql"
	"fmt"
)

type UserRepository interface {
	GetByUsername(username string) (*models.User, error)
	GetByID(id int) (*models.User, error)
	Create(user *models.User) error
	GetUser(search, sortBy, order string, page, limit int) ([]models.User, int, error)
	SoftDelete(id int) error
}

type userRepository struct{ DB *sql.DB }

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{DB: db}
}

func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	var u models.User
	err := r.DB.QueryRow("SELECT id, username, password_hash, role FROM users WHERE username=$1", username).
		Scan(&u.ID, &u.Username, &u.Password, &u.Role)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) GetByID(id int) (*models.User, error) {
	var u models.User
	err := r.DB.QueryRow("SELECT id, username, role FROM users WHERE id=$1", id).
		Scan(&u.ID, &u.Username, &u.Role)
	if err != nil {
		return nil, err
	}
	return &u, nil
}


func (r *userRepository) Create(u *models.User) error {
    query := `INSERT INTO users (email, username, password_hash, role) VALUES ($1, $2, $3, $4)`
    _, err := r.DB.Exec(query, u.Email, u.Username, u.Password, u.Role)
    return err
}

func (r *userRepository) GetUser(search, sortBy, order string, page, limit int) ([]models.User, int, error) {
	var users []models.User
	var total int

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	// validasi sort
	allowedSort := map[string]bool{"id": true, "username": true, "email": true}
	if !allowedSort[sortBy] {
		sortBy = "id"
	}
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	// base query
	query := `SELECT id, username, email, password_hash FROM users WHERE 1=1`
	args := []interface{}{}

	if search != "" {
		query += " AND (username ILIKE $1 OR email ILIKE $2)"
		args = append(args, "%"+search+"%", "%"+search+"%")
	}

	// hitung total
	countQuery := "SELECT COUNT(*) FROM (" + query + ") AS sub"
	if len(args) > 0 {
		if err := r.DB.QueryRow(countQuery, args...).Scan(&total); err != nil {
			return nil, 0, err
		}
	} else {
		if err := r.DB.QueryRow(countQuery).Scan(&total); err != nil {
			return nil, 0, err
		}
	}

	// tambahkan sorting dan limit/offset
	query = fmt.Sprintf("%s ORDER BY %s %s LIMIT %d OFFSET %d", query, sortBy, order, limit, offset)

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Password); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}

	return users, total, nil
}


func (r *userRepository) SoftDelete(id int) error {
	query := `UPDATE users SET is_delete = NOW() WHERE id = $1`
	 res, err := r.DB.Exec(query, id)
	 
	 if err != nil {
		return err
	 }

	 rows, _ := res.RowsAffected()
	 if rows == 0 {
		return fmt.Errorf("user not found")
	 }

    return err
}



