package auth

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"grpcAuth/interanl/domain/models"
	"grpcAuth/interanl/lib/jwt"
	"grpcAuth/interanl/services"
	"grpcAuth/interanl/storage"
	"log/slog"
	"time"
)

type Auth struct {
	log      *slog.Logger
	tokenTTL time.Duration

	userSaver    UserSaver
	userProvider UserProveder
	appProvider  AppProvider
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error)
}

type UserProveder interface {
	User(ctx context.Context, email string) (user models.User, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

// New returns new instance of Auth service
func New(
	log *slog.Logger,
	tokenTTL time.Duration,
	userSaver UserSaver,
	userProvider UserProveder,
	appProvider AppProvider) *Auth {
	return &Auth{
		log:          log,
		tokenTTL:     tokenTTL,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
	}
}

func (a *Auth) Login(ctx context.Context, email, password string, appID int) (string, error) {
	fn := "auth.Login"

	log := a.log.With(
		slog.String("fn", fn),
		slog.String("email", email),
		slog.Int("appID", appID),
	)
	log.Info("login user")

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", slog.String("err", err.Error()))
			return "", fmt.Errorf("%s: %w", fn, services.ErrInvalidCredentials)
		}

		log.Error("failed to save user", slog.String("err", err.Error()))

		return "", fmt.Errorf("%s: %w", fn, services.ErrInvalidCredentials)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Error("invalid credentials", slog.String("err", err.Error()))
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		log.Error("app dose not exists", slog.String("err", err.Error()))
		return "", fmt.Errorf("%s: %w", fn, storage.ErrAppNotFound)
	}

	log.Info("user logged in sucessfully")

	token, err := jwt.NewToken(user, app, a.tokenTTL)

	if err != nil {
		log.Error("faild to create token", slog.String("err", err.Error()))

		return "", fmt.Errorf("%s: %w", err)
	}

	return token, nil
}

func (a *Auth) Register(ctx context.Context, email, password string) (int64, error) {
	fn := "auth.RegiserNewUser"

	log := a.log.With(
		slog.String("fn", fn),
		slog.String("email", email),
	)
	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		log.Error("faild to generate pass hash", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	uid, err := a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserAlreadyExists) {
			log.Warn("user already exists", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", fn, services.ErrEmailAlreadyExists)
		}

		log.Error("failed to save user", slog.String("err", err.Error()))

		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return uid, nil

}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	fn := "auth.IsAdmin"

	log := a.log.With(
		slog.String("fn", fn),
		slog.Int("email", int(userID)),
	)
	log.Info("checking if user is admin")

	isAdmin, err := a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", slog.String("err", err.Error()))
			return false, fmt.Errorf("%s: %w", fn, services.ErrUserNotFound)
		}

		return false, fmt.Errorf("%s: %w", fn, err)
	}

	log.Info("checked if user if admin", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}
