package usecase

import (
	"context"
	"fmt"
	"github.com/irvankadhafi/user-balance-transfer-service/internal/model"
	"github.com/irvankadhafi/user-balance-transfer-service/internal/repository"
	"github.com/irvankadhafi/user-balance-transfer-service/utils"
	"github.com/sirupsen/logrus"
)

type bankBalanceUsecase struct {
	bankBalanceRepo        model.BankBalanceRepository
	bankBalanceHistoryRepo model.BankBalanceHistoryRepository
	gormTransactioner      repository.GormTransactioner
	sessionRepo            model.SessionRepository
}

func NewBankBalanceUsecase(
	bankBalanceRepo model.BankBalanceRepository,
	bankBalanceHistoryRepo model.BankBalanceHistoryRepository,
	gormTransactioner repository.GormTransactioner,
	sessionRepo model.SessionRepository,
) model.BankBalanceUsecase {
	return &bankBalanceUsecase{
		bankBalanceRepo:        bankBalanceRepo,
		bankBalanceHistoryRepo: bankBalanceHistoryRepo,
		gormTransactioner:      gormTransactioner,
		sessionRepo:            sessionRepo,
	}
}

func (b *bankBalanceUsecase) CreateBankAccount(ctx context.Context, sessionID int, code string) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":  utils.DumpIncomingContext(ctx),
		"code": code,
	})

	// Get IP, UserAgent, etc..
	session, err := b.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		logger.Error(err)
		return err
	}

	tx := b.gormTransactioner.Begin(ctx)
	bankBalance := &model.BankBalance{
		Balance:        0,
		BalanceAchieve: 0,
		Code:           code,
		Enable:         true,
	}
	err = b.bankBalanceRepo.CreateWithTransaction(ctx, tx, bankBalance)
	if err != nil {
		logger.Error(err)
		b.gormTransactioner.Rollback(tx)
		return err
	}

	bankBalanceAudit := &model.BankBalanceHistory{
		BankBalanceID: bankBalance.ID,
		BalanceBefore: 0,
		BalanceAfter:  0,
		Activity:      fmt.Sprintf("Create Bank Account by userID: %d", session.UserID),
		Type:          model.CREDIT,
		IPAddress:     session.Location,
		Location:      session.IPAddress,
		UserAgent:     session.UserAgent,
		Author:        "system",
	}
	err = b.bankBalanceHistoryRepo.CreateWithTransaction(ctx, tx, bankBalanceAudit)
	if err != nil {
		logger.Error(err)
		b.gormTransactioner.Rollback(tx)
		return err
	}

	if err = b.gormTransactioner.Commit(tx); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (b *bankBalanceUsecase) AddBankBalance(ctx context.Context, input model.AddBankBalanceInput) error {
	if input.Code == "" || input.Balance <= 0 {
		return ErrFailedPrecondition
	}

	logger := logrus.WithFields(logrus.Fields{
		"ctx":   utils.DumpIncomingContext(ctx),
		"input": utils.Dump(input),
	})

	// Get IP, UserAgent, etc..
	session, err := b.sessionRepo.FindByID(ctx, input.SessionID)
	if err != nil {
		logger.Error(err)
		return err
	}

	// Get the current bank balance
	balance, err := b.bankBalanceRepo.GetCurrentBankBalanceByCode(ctx, input.Code)
	if err != nil {
		logger.Error(err)
		return err
	}
	activity := "Add Balance"

	tx := b.gormTransactioner.Begin(ctx)
	balance.Balance += input.Balance
	balance.BalanceAchieve += input.Balance
	balance.Code = input.Code
	err = b.bankBalanceRepo.UpsertWithTransaction(ctx, tx, balance)
	if err != nil {
		logger.Error(err)
		b.gormTransactioner.Rollback(tx)
		return err
	}

	bankBalanceAudit := &model.BankBalanceHistory{
		BankBalanceID: balance.ID,
		BalanceBefore: balance.BalanceAchieve - input.Balance,
		BalanceAfter:  balance.BalanceAchieve,
		Activity:      activity,
		Type:          model.CREDIT,
		IPAddress:     session.Location,
		Location:      session.IPAddress,
		UserAgent:     session.UserAgent,
		Author:        input.Author,
	}
	err = b.bankBalanceHistoryRepo.CreateWithTransaction(ctx, tx, bankBalanceAudit)
	if err != nil {
		logger.Error(err)
		b.gormTransactioner.Rollback(tx)
		return err
	}

	if err = b.gormTransactioner.Commit(tx); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (b *bankBalanceUsecase) GetBankBalanceByID(ctx context.Context, bankBalanceID int) (*model.BankBalance, error) {
	//TODO implement me
	panic("implement me")
}

func (b *bankBalanceUsecase) TransferUserBalance(ctx context.Context, userIDFrom, userIDTo int, balance float64, code string) error {
	//TODO implement me
	panic("implement me")
}
