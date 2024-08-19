package usecase

import (
	"context"
	"errors"

	"pragusga/internal/domain"
	"pragusga/internal/events"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type AuthUseCase struct {
	userRepo           domain.UserRepository
	userEventPublisher *events.UserEventPublisher
}

func NewAuthUseCase(userRepo domain.UserRepository, userEventPublisher *events.UserEventPublisher) *AuthUseCase {
	return &AuthUseCase{userRepo: userRepo, userEventPublisher: userEventPublisher}
}

func (uc *AuthUseCase) SignUp(ctx context.Context, email, password string) error {
	tenantId := "public"
	ep, err := emailpassword.SignUp(tenantId, email, password)

	if ep.EmailAlreadyExistsError != nil {
		return errors.New("email already exists")
	}

	if err != nil {
		return errors.New("supertokens failed to sign up")
	}

	user := &domain.User{Email: email, ID: ep.OK.User.ID}
	err = uc.userRepo.Create(ctx, user)
	if err != nil {
		// Rollback user creation
		supertokens.DeleteUser(ep.OK.User.ID)
		return errors.New("failed to create user")
	}

	// Publish USER_CREATED event
	return uc.userEventPublisher.PublishUserCreated(ctx, user)
}

func (uc *AuthUseCase) SignIn(ctx context.Context, email, password string) (*domain.User, error) {
	tenantId := "public"
	ep, err := emailpassword.SignIn(tenantId, email, password)
	if ep.WrongCredentialsError != nil {
		return nil, errors.New("wrong credentials")
	}

	if err != nil {
		return nil, errors.New("supertokens failed to sign in")
	}

	return uc.userRepo.GetById(ctx, ep.OK.User.ID)

}

func (uc *AuthUseCase) GetUserInfo(ctx context.Context, userId string) (*domain.User, error) {
	return uc.userRepo.GetById(ctx, userId)
}
