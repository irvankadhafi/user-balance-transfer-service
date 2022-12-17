package usecase

import (
	"context"
	"fmt"
	"github.com/irvankadhafi/user-balance-transfer-service/internal/model"
	"github.com/irvankadhafi/user-balance-transfer-service/internal/repository"
	"github.com/irvankadhafi/user-balance-transfer-service/utils"
	"github.com/sirupsen/logrus"
)

type userBalanceUsecase struct {
	userRepo               model.UserRepository
	userBalanceRepo        model.UserBalanceRepository
	userBalanceHistoryRepo model.UserBalanceHistoryRepository
	gormTransactioner      repository.GormTransactioner
	sessionRepo            model.SessionRepository
}

func NewUserBalanceUsecase(
	userRepo model.UserRepository,
	userBalanceRepo model.UserBalanceRepository,
	userBalanceHistoryRepo model.UserBalanceHistoryRepository,
	gormTransactioner repository.GormTransactioner,
	sessionRepo model.SessionRepository,
) model.UserBalanceUsecase {
	return &userBalanceUsecase{
		userRepo:               userRepo,
		userBalanceRepo:        userBalanceRepo,
		userBalanceHistoryRepo: userBalanceHistoryRepo,
		gormTransactioner:      gormTransactioner,
		sessionRepo:            sessionRepo,
	}
}

func (u *userBalanceUsecase) AddUserBalance(ctx context.Context, input model.AddUserBalanceInput) error {
	if input.UserID <= 0 || input.Balance <= 0 {
		return ErrFailedPrecondition
	}
	logger := logrus.WithFields(logrus.Fields{
		"ctx":   utils.DumpIncomingContext(ctx),
		"input": utils.Dump(input),
	})

	// Get IP, UserAgent, etc..
	session, err := u.sessionRepo.FindByID(ctx, input.SessionID)
	if err != nil {
		logger.Error(err)
		return err
	}

	// Get the current user balance
	balance, err := u.userBalanceRepo.GetCurrentUserBalanceByUserID(ctx, input.UserID)
	if err != nil {
		logger.Error(err)
		return err
	}

	activity := "Add Balance"

	tx := u.gormTransactioner.Begin(ctx)
	balance.Balance += input.Balance
	balance.BalanceAchieve += input.Balance
	balance.UserID = input.UserID
	err = u.userBalanceRepo.UpsertWithTransaction(ctx, tx, balance)
	if err != nil {
		logger.Error(err)
		u.gormTransactioner.Rollback(tx)
		return err
	}

	userBalanceAudit := &model.UserBalanceHistory{
		UserBalanceID: balance.ID,
		BalanceBefore: balance.BalanceAchieve - input.Balance,
		BalanceAfter:  balance.BalanceAchieve,
		Activity:      activity,
		Type:          model.CREDIT,
		IPAddress:     session.Location,
		Location:      session.IPAddress,
		UserAgent:     session.UserAgent,
		Author:        input.Author,
	}
	err = u.userBalanceHistoryRepo.CreateWithTransaction(ctx, tx, userBalanceAudit)
	if err != nil {
		logger.Error(err)
		u.gormTransactioner.Rollback(tx)
		return err
	}

	if err = u.gormTransactioner.Commit(tx); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (u *userBalanceUsecase) TransferUserBalance(ctx context.Context, input model.TransferUserBalanceInput) error {
	if input.FromUserID <= 0 || input.ToUserID <= 0 || input.Balance <= 0 {
		return ErrFailedPrecondition
	}

	logger := logrus.WithFields(logrus.Fields{
		"ctx":   utils.DumpIncomingContext(ctx),
		"input": utils.Dump(input),
	})

	session, err := u.sessionRepo.FindByID(ctx, input.SessionID)
	if err != nil {
		logger.Error(err)
		return err
	}

	// Check user tujuan transfer
	toUser, err := u.userRepo.FindByID(ctx, input.ToUserID)
	if err != nil {
		logger.Error(err)
		return err
	}
	if toUser == nil {
		return ErrNotFound
	}

	fromUser, err := u.userRepo.FindByID(ctx, input.FromUserID)
	if err != nil {
		logger.Error(err)
		return err
	}
	if fromUser == nil {
		return ErrNotFound
	}

	// get the current balance for the from user
	fromBalance, err := u.userBalanceRepo.GetCurrentUserBalanceByUserID(ctx, fromUser.ID)
	if err != nil {
		logger.Error(err)
		return err
	}

	if fromBalance.BalanceAchieve < input.Balance {
		return ErrBalanceNotEnough
	}

	// get the current balance for the to user
	toBalance, err := u.userBalanceRepo.GetCurrentUserBalanceByUserID(ctx, toUser.ID)
	if err != nil {
		logger.Error(err)
		return err
	}

	tx := u.gormTransactioner.Begin(ctx)
	fromBalance.Balance -= input.Balance
	err = u.userBalanceRepo.UpsertWithTransaction(ctx, tx, fromBalance)
	if err != nil {
		logger.Error(err)
		u.gormTransactioner.Rollback(tx)
		return err
	}

	fromUserBalanceAudit := &model.UserBalanceHistory{
		UserBalanceID: fromBalance.ID,
		BalanceBefore: fromBalance.BalanceAchieve,
		BalanceAfter:  fromBalance.Balance,
		Activity:      fmt.Sprintf("Transfer to %s", toUser.Username),
		Type:          model.DEBIT,
		IPAddress:     session.Location,
		Location:      session.IPAddress,
		UserAgent:     session.UserAgent,
		Author:        input.Author,
	}
	err = u.userBalanceHistoryRepo.CreateWithTransaction(ctx, tx, fromUserBalanceAudit)
	if err != nil {
		logger.Error(err)
		u.gormTransactioner.Rollback(tx)
		return err
	}

	toBalance.Balance += input.Balance
	toBalance.BalanceAchieve += input.Balance
	toBalance.UserID = toUser.ID
	err = u.userBalanceRepo.UpsertWithTransaction(ctx, tx, toBalance)
	if err != nil {
		logger.Error(err)
		u.gormTransactioner.Rollback(tx)
		return err
	}

	toUserBalanceAudit := &model.UserBalanceHistory{
		UserBalanceID: toBalance.ID,
		BalanceBefore: toBalance.BalanceAchieve,
		BalanceAfter:  toBalance.Balance,
		Activity:      fmt.Sprintf("Transfer from %s", fromUser.Username),
		Type:          model.CREDIT,
		IPAddress:     session.Location,
		Location:      session.IPAddress,
		UserAgent:     session.UserAgent,
		Author:        input.Author,
	}
	err = u.userBalanceHistoryRepo.CreateWithTransaction(ctx, tx, toUserBalanceAudit)
	if err != nil {
		logger.Error(err)
		u.gormTransactioner.Rollback(tx)
		return err
	}

	if err = u.gormTransactioner.Commit(tx); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (u *userBalanceUsecase) GetCurrentUserBalanceByUserID(ctx context.Context, userID int) (*model.UserBalance, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":    utils.DumpIncomingContext(ctx),
		"userID": userID,
	})
	// Get the current user balance
	balance, err := u.userBalanceRepo.GetCurrentUserBalanceByUserID(ctx, userID)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return balance, nil
}
