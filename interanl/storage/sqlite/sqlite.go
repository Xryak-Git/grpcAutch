package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"grpcAuth/interanl/domain/models"
	"grpcAuth/interanl/storage"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const fn = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const fn = "storage.sqlite.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO users (email, pass_hash) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	res, err := stmt.ExecContext(ctx, email, passHash)
	if err != nil {
		var sqliteErr sqlite3.Error

		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", fn, storage.ErrUserAlreadyExists)
		}
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const fn = "storage.sqlite.User"

	stmt, err := s.db.Prepare("SELECT id, email, pass_hash FROM users WHERE email = ?")
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", fn, err)
	}

	var user models.User
	row := stmt.QueryRowContext(ctx, email)
	err = row.Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", fn, storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", fn, err)
	}

	return user, nil
}
func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const fn = "storage.sqlite.IsAdmin"

	stmt, err := s.db.Prepare("SELECT * FROM users WHERE id = ?")
	if err != nil {
		return false, fmt.Errorf("%s: %w", fn, err)
	}

	var isAdmin bool
	row := stmt.QueryRowContext(ctx, userID)
	err = row.Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", fn, storage.ErrUserNotFound)
		}
		return false, fmt.Errorf("%s: %w", fn, err)
	}
	return isAdmin, nil
}
func (s *Storage) App(ctx context.Context, appID int) (models.App, error) {
	const fn = "storage.sqlite.App"

	stmt, err := s.db.Prepare("SELECT * FROM apps WHERE id = ?")
	if err != nil {
		return models.App{}, fmt.Errorf("%s: %w", fn, err)
	}

	var app models.App
	row := stmt.QueryRowContext(ctx, appID)
	err = row.Scan(&app.AppID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", fn, storage.ErrAppNotFound)
		}
		return models.App{}, fmt.Errorf("%s: %w", fn, err)
	}
	return app, nil
}
