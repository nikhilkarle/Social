package store

import (
	"context"
	"database/sql"
	"errors"
)

type User struct{
	ID int64 `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"-"`
	CreatedAt string `json:"created_at"`
}

type UserStore struct{
	db *sql.DB 
}

func ( s *UserStore) Create(ctx context.Context, user *User) error{
	query := `
		INSERT INTO users (username, password, email) VALUES($1, $2, $3) RETURNING id,
		created_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query, 
		user.Username, 
		user.Password,
		user.Email,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)

	if err != nil{
		return err
	 }

	 return nil
}

func (s *UserStore) GetByID(ctx context.Context, userID int64) (*User, error){
	query := `
	Select id, username, password, email, created_at from users
	Where ID = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var user User
	err := s.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil{
		switch{
			case errors.Is(err, sql.ErrNoRows):
				return nil, ErrNotFound

			default:
				return nil, err
		} 
	 }
	 return &user, nil
}