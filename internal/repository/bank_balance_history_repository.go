package repository

import (
	"context"
	"github.com/irvankadhafi/user-balance-transfer-service/internal/model"
	"github.com/irvankadhafi/user-balance-transfer-service/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type bankBalanceHistoryRepository struct {
	db *gorm.DB
}

func NewBankBalanceHistoryRepository(
	db *gorm.DB,
) model.BankBalanceHistoryRepository {
	return &bankBalanceHistoryRepository{
		db: db,
	}
}

func (b *bankBalanceHistoryRepository) CreateWithTransaction(ctx context.Context, tx *gorm.DB, input *model.BankBalanceHistory) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":                utils.DumpIncomingContext(ctx),
		"bankBalanceHistory": utils.Dump(input),
	})

	err := tx.WithContext(ctx).Create(input).Error
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
