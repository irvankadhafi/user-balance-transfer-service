package usecase

import (
	"github.com/irvankadhafi/user-balance-transfer-service/internal/helper"
	"github.com/irvankadhafi/user-balance-transfer-service/internal/model"
	"github.com/irvankadhafi/user-balance-transfer-service/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type userUsecase struct {
	userRepo model.UserRepository
}

func NewUserUsecase(userRepo model.UserRepository) model.UserUsecase {
	return &userUsecase{userRepo: userRepo}
}

func (u *userUsecase) Create(ctx context.Context, input model.CreateUserInput) (*model.User, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":   utils.DumpIncomingContext(ctx),
		"input": utils.Dump(input),
	})

	input.Email = helper.FormatEmail(input.Email)
	if err := input.Validate(); err != nil {
		logger.Error(err)
		return nil, err
	}

	// check if user already exist
	existingUser, err := u.userRepo.FindByEmail(ctx, input.Email)
	switch {
	case existingUser != nil:
		return nil, ErrDuplicateUser
	case err == ErrNotFound:
	default:
		logger.Error(err)
		return nil, err
	}

	user := &model.User{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	}

	err = u.userRepo.Create(ctx, user)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return u.FindByID(ctx, user.ID)

}

func (u *userUsecase) FindByID(ctx context.Context, userID int) (*model.User, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		logrus.WithField("userID", userID).Error(err)
		return nil, err
	}

	return user, nil
}
