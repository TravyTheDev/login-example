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

func (s *UserStore) getUserByID(id int) (*UserPayload, error) {
	stmt := `SELECT users.id, users.username, users.email FROM users WHERE id = ?`

	rows, err := s.db.Query(stmt, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	user := new(UserPayload)
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

func (s *UserStore) getAllUsers() ([]*UserPayload, error) {
	stmt := `SELECT users.id, users.username, users.email FROM users`

	rows, err := s.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*UserPayload, 0)
	for rows.Next() {
		user, err := scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func scanRowsIntoUser(rows *sql.Rows) (*UserPayload, error) {
	user := new(UserPayload)
	err := rows.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
	)

	if err != nil {
		return nil, err
	}
	return user, nil
}
