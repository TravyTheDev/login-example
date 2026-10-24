package auth

import (
	"database/sql"
	"fmt"
)

type AuthStore struct {
	db *sql.DB
}

func NewAuthStore(db *sql.DB) *AuthStore {
	return &AuthStore{
		db: db,
	}
}

func (st *AuthStore) CreateSession(s *Session) (*Session, error) {
	stmt := `INSERT INTO sessions (id, user_email, refresh_token, is_revoked, expires_at) VALUES (?, ?, ?, ?, ?)`

	_, err := st.db.Exec(stmt, s.ID, s.UserEmail, s.RefreshToken, s.IsRevoked, s.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("error inserting session: %w", err)
	}

	return s, nil
}

func (st *AuthStore) GetSession(id string) (*Session, error) {
	stmt := `SELECT * FROM sessions WHERE id = ?`

	rows, err := st.db.Query(stmt, id)
	if err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	defer rows.Close()
	s := new(Session)
	for rows.Next() {
		err := rows.Scan(
			&s.ID,
			&s.UserEmail,
			&s.RefreshToken,
			&s.IsRevoked,
			&s.CreatedAt,
			&s.ExpiresAt,
		)
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (st *AuthStore) RevokeSession(id string) error {

	stmt := `UPDATE sessions SET is_revoked=1 WHERE id = ?`
	_, err := st.db.Exec(stmt, id)
	if err != nil {
		return fmt.Errorf("error revoking session: %w", err)
	}

	return nil
}

func (st *AuthStore) DeleteSession(id string) error {
	_, err := st.db.Exec("DELETE FROM sessions WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("error deleting session: %w", err)
	}

	return nil
}

func (st AuthStore) CreateUser(user User) error {
	stmt := `INSERT INTO users (username, email, password, role) VALUES (?, ?, ?, ?)`

	_, err := st.db.Exec(stmt, user.Username, user.Email, user.Password, user.Role)
	if err != nil {
		return err
	}

	return nil
}

func (st *AuthStore) GetUserByEmail(email string) (*User, error) {
	stmt := `SELECT * FROM users WHERE email = ?`

	rows, err := st.db.Query(stmt, email)
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
