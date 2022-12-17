package repository

import (
	"context"
	"github.com/irvankadhafi/user-balance-transfer-service/internal/model"
	"github.com/irvankadhafi/user-balance-transfer-service/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type bankBalanceRepository struct {
	db *gorm.DB
}

func NewBankBalanceRepository(
	db *gorm.DB,
) model.BankBalanceRepository {
	return &bankBalanceRepository{
		db: db,
	}
}

func (b *bankBalanceRepository) CreateWithTransaction(ctx context.Context, tx *gorm.DB, bankBalance *model.BankBalance) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":         utils.DumpIncomingContext(ctx),
		"bankBalance": utils.Dump(bankBalance),
	})

	err := tx.WithContext(ctx).Create(bankBalance).Error
	if err != nil {
		logger.Error(err)
		return err
	}
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (b *bankBalanceRepository) UpsertWithTransaction(ctx context.Context, tx *gorm.DB, bankBalance *model.BankBalance) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":         utils.DumpIncomingContext(ctx),
		"bankBalance": utils.Dump(bankBalance),
	})
	switch {
	case bankBalance.ID > 0:
		err := tx.WithContext(ctx).Updates(bankBalance).Error
		if err != nil {
			logger.Error(err)
			return err
		}
	default:
		err := tx.WithContext(ctx).Create(bankBalance).Error
		if err != nil {
			logger.Error(err)
			return err
		}
	}

	return nil
}

func (b *bankBalanceRepository) GetCurrentBankBalanceByCode(ctx context.Context, code string) (*model.BankBalance, error) {
	var bankBalance model.BankBalance
	err := b.db.WithContext(ctx).Order("created_at desc").Take(&bankBalance, "code = ?", code).Error
	switch err {
	case nil:
		return &bankBalance, nil
	case gorm.ErrRecordNotFound:
		return &bankBalance, nil
	default:
		logrus.WithFields(logrus.Fields{
			"ctx":  utils.DumpIncomingContext(ctx),
			"code": code,
		}).Error(err)
		return nil, err
	}
}
