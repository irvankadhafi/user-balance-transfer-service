package repository

import (
	"context"
	"github.com/irvankadhafi/user-balance-transfer-service/internal/model"
	"github.com/irvankadhafi/user-balance-transfer-service/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type userBalanceHistoryRepository struct {
	db *gorm.DB
}

func NewUserBalanceHistoryRepository(
	db *gorm.DB,
) model.UserBalanceHistoryRepository {
	return &userBalanceHistoryRepository{
		db: db,
	}
}

func (u userBalanceHistoryRepository) CreateWithTransaction(ctx context.Context, tx *gorm.DB, input *model.UserBalanceHistory) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":                utils.DumpIncomingContext(ctx),
		"userBalanceHistory": utils.Dump(input),
	})

	err := tx.WithContext(ctx).Create(input).Error
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
