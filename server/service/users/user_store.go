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

func (s *UserStore) CreateSession(session *Session) (*Session, error) {
	stmt := `INSERT INTO sessions (id, user_email, refresh_token, is_revoked, expires_at) VALUES (?, ?, ?, ?, ?)`

	_, err := s.db.Exec(
		stmt, session.ID,
		session.UserEmail,
		session.RefreshToken,
		session.IsRevoked,
		session.ExpiresAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error inserting session: %w", err)
	}

	return session, nil
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
