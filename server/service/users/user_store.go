package users

import (
	"database/sql"
	"fmt"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (s UserStore) CreateUser(user User) error {
	stmt := `INSERT INTO users (username, email, password, role) VALUES (?, ?, ?, ?)`

	_, err := s.db.Exec(stmt, user.Username, user.Email, user.Password, user.Role)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) GetUserByEmail(email string) (*User, error) {
	stmt := `SELECT * FROM users WHERE email = ?`

	rows, err := s.db.Query(stmt, email)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	user := new(User)
	for rows.Next() {
		user, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (s *UserStore) getUserByID(id int) (*User, error) {
	stmt := `SELECT * FROM users WHERE id = ?`

	rows, err := s.db.Query(stmt, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	user := new(User)
	for rows.Next() {
		user, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func scanRowsIntoUser(rows *sql.Rows) (*User, error) {
	user := new(User)
	err := rows.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return user, nil
}
