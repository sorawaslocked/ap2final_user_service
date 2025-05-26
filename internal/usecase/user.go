package usecase

import (
	"context"
	"github.com/sorawaslocked/ap2final_base/pkg/logger"
	"github.com/sorawaslocked/ap2final_base/pkg/security"
	"github.com/sorawaslocked/ap2final_user_service/internal/model"
	"log/slog"
	"time"
)

type User struct {
	log         *slog.Logger
	repo        UserRepository
	tokenRepo   TokenRepository
	jwtProvider *security.JWTProvider
}

func NewUser(
	log *slog.Logger,
	repo UserRepository,
	tokenRepo TokenRepository,
	jwtProvider *security.JWTProvider,
) *User {
	return &User{
		log:         log,
		repo:        repo,
		tokenRepo:   tokenRepo,
		jwtProvider: jwtProvider,
	}
}

func (uc *User) Register(ctx context.Context, user model.User) (model.User, error) {
	const op = "usecase.User.Register"

	log := uc.log.With(slog.String("op", op))

	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = time.Now().UTC()

	hashedPassword, err := security.HashPassword(user.Password)
	if err != nil {
		log.Error("hashing password", logger.Err(err))

		return model.User{}, err
	}

	user.Password = hashedPassword

	createdUser, err := uc.repo.InsertOne(ctx, user)
	if err != nil {
		log.Error("creating user", logger.Err(err))

		return model.User{}, err
	}

	return createdUser, err
}

func (uc *User) Login(ctx context.Context, user model.User) (model.Token, error) {
	const op = "usecase.User.Login"

	log := uc.log.With(slog.String("op", op))

	userFromDb, err := uc.repo.FindOne(ctx, model.UserFilter{Email: &user.Email})
	if err != nil {
		log.Warn("finding user", logger.Err(err))

		return model.Token{}, err
	}

	err = security.CheckPassword(user.Password, userFromDb.PasswordHash)
	if err != nil {
		err := model.ErrPasswordsDoNotMatch
		log.Warn("checking password", logger.Err(err))

		return model.Token{}, err
	}

	accessToken, err := uc.jwtProvider.GenerateAccessToken(userFromDb.ID, userFromDb.Role)
	if err != nil {
		log.Warn("generating access token", logger.Err(err))

		return model.Token{}, err
	}

	refreshToken, err := uc.jwtProvider.GenerateRefreshToken(userFromDb.ID)
	if err != nil {
		log.Warn("generating refresh token", logger.Err(err))

		return model.Token{}, err
	}

	session := model.Session{
		UserID:       userFromDb.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().UTC().Add(uc.jwtProvider.RefreshTokenTTL),
		CreatedAt:    time.Now().UTC(),
	}

	err = uc.tokenRepo.InsertOne(ctx, session)
	if err != nil {
		log.Warn("inserting new token session", logger.Err(err))

		return model.Token{}, err
	}

	return model.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (uc *User) RefreshToken(ctx context.Context, refreshToken string) (model.Token, error) {
	const op = "usecase.User.RefreshToken"

	log := uc.log.With(slog.String("op", op))

	session, err := uc.tokenRepo.FindOneByToken(ctx, refreshToken)
	if err != nil {
		log.Warn("refreshing token", logger.Err(err))

		return model.Token{}, err
	}

	if session.ExpiresAt.Before(time.Now().UTC()) {
		err := model.ErrRefreshTokenExpired

		log.Warn("refreshing token", logger.Err(err))

		return model.Token{}, err
	}

	user, err := uc.repo.FindOne(ctx, model.UserFilter{ID: &session.UserID})
	if err != nil {
		log.Warn("finding user", logger.Err(err))

		return model.Token{}, err
	}

	accessToken, err := uc.jwtProvider.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		log.Warn("generating access token", logger.Err(err))

		return model.Token{}, err
	}

	newRefreshToken, err := uc.jwtProvider.GenerateRefreshToken(user.ID)
	if err != nil {
		log.Warn("generating refresh token", logger.Err(err))

		return model.Token{}, err
	}

	err = uc.tokenRepo.DeleteByToken(ctx, refreshToken)
	if err != nil {
		log.Warn("deleting refresh token", logger.Err(err))

		return model.Token{}, err
	}

	newSession := model.Session{
		UserID:       user.ID,
		RefreshToken: newRefreshToken,
		ExpiresAt:    time.Now().UTC().Add(uc.jwtProvider.RefreshTokenTTL),
		CreatedAt:    time.Now().UTC(),
	}

	err = uc.tokenRepo.InsertOne(ctx, newSession)
	if err != nil {
		log.Warn("inserting new token session", logger.Err(err))

		return model.Token{}, err
	}

	return model.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (uc *User) GetByID(ctx context.Context, token model.Token, id string) (model.User, error) {
	const op = "usecase.User.GetByID"

	log := uc.log.With(slog.String("op", op))

	claims, err := uc.jwtProvider.VerifyAndParseClaims(token.AccessToken)
	if err != nil {
		err := model.ErrInvalidToken
		log.Warn(
			"verifying token and parsing claims",
			logger.Err(err),
			slog.String("accessToken", token.AccessToken),
		)

		return model.User{}, err
	}

	if claims.UserID == nil || claims.Role == nil {
		err := model.ErrEmptyClaims
		log.Warn(
			"checking claims",
			logger.Err(err),
			slog.String("accessToken", token.AccessToken),
		)

		return model.User{}, err
	}

	if *claims.Role != "admin" || *claims.UserID != id {
		err := model.ErrUnauthorized
		log.Warn(
			"checking claims",
			logger.Err(err),
			slog.String("claimsUserID", *claims.UserID),
			slog.String("claimsRole", *claims.Role),
		)

		return model.User{}, err
	}

	user, err := uc.repo.FindOne(ctx, model.UserFilter{ID: &id})
	if err != nil {
		log.Warn(
			"finding user",
			logger.Err(err),
			slog.String("id", id),
		)

		return model.User{}, err
	}

	return user, nil
}

func (uc *User) UpdateByID(
	ctx context.Context,
	token model.Token,
	id string,
	credentialsUpdate model.UserCredentialUpdateData,
	update model.UserUpdateData,
) (model.User, error) {
	const op = "usecase.User.UpdateByID"

	log := uc.log.With(slog.String("op", op))

	update.UpdatedAt = time.Now().UTC()

	claims, err := uc.jwtProvider.VerifyAndParseClaims(token.AccessToken)
	if err != nil {
		err := model.ErrInvalidToken
		log.Warn(
			"verifying token and parsing claims",
			logger.Err(err),
			slog.String("accessToken", token.AccessToken),
		)

		return model.User{}, err
	}

	if claims.UserID == nil || claims.Role == nil {
		err := model.ErrEmptyClaims
		log.Warn(
			"checking claims",
			logger.Err(err),
			slog.String("accessToken", token.AccessToken),
		)

		return model.User{}, err
	}

	if *claims.Role != "admin" || *claims.UserID != id {
		err := model.ErrUnauthorized
		log.Warn(
			"checking claims",
			logger.Err(err),
			slog.String("claimsUserID", *claims.UserID),
			slog.String("claimsRole", *claims.Role),
		)

		return model.User{}, err
	}

	userFromDb, err := uc.repo.FindOne(ctx, model.UserFilter{ID: &id})
	if err != nil {
		log.Warn(
			"finding user",
			logger.Err(err),
			slog.String("id", id),
		)

		return model.User{}, err
	}

	err = security.CheckPassword(credentialsUpdate.CurrentPassword, userFromDb.PasswordHash)
	if err != nil {
		err := model.ErrPasswordsDoNotMatch
		log.Warn("checking passwords", logger.Err(err))

		return model.User{}, err
	}

	hashedPassword, err := security.HashPassword(credentialsUpdate.NewPassword)
	if err != nil {
		log.Warn("hashing password", logger.Err(err))

		return model.User{}, err
	}
	update.PasswordHash = &hashedPassword

	updateUser, err := uc.repo.UpdateOne(
		ctx,
		model.UserFilter{ID: &id},
		update,
	)
	if err != nil {
		log.Warn(
			"updating user",
			logger.Err(err),
			slog.String("id", id),
		)

		return model.User{}, err
	}

	return updateUser, nil
}

func (uc *User) DeleteByID(ctx context.Context, token model.Token, id string) (model.User, error) {
	const op = "usecase.User.DeleteByID"

	log := uc.log.With(slog.String("op", op))

	claims, err := uc.jwtProvider.VerifyAndParseClaims(token.AccessToken)
	if err != nil {
		err := model.ErrInvalidToken
		log.Warn(
			"verifying token and parsing claims",
			logger.Err(err),
			slog.String("accessToken", token.AccessToken),
		)

		return model.User{}, err
	}

	if claims.UserID == nil || claims.Role == nil {
		err := model.ErrEmptyClaims
		log.Warn(
			"checking claims",
			logger.Err(err),
			slog.String("accessToken", token.AccessToken),
		)

		return model.User{}, err
	}

	if *claims.Role != "admin" || *claims.UserID != id {
		err := model.ErrUnauthorized
		log.Warn(
			"checking claims",
			logger.Err(err),
			slog.String("claimsUserID", *claims.UserID),
			slog.String("claimsRole", *claims.Role),
		)

		return model.User{}, err
	}

	deletedUser, err := uc.repo.DeleteOne(ctx, model.UserFilter{ID: &id})
	if err != nil {
		log.Warn(
			"deleting user",
			logger.Err(err),
			slog.String("id", id),
		)

		return model.User{}, err
	}

	return deletedUser, nil
}
