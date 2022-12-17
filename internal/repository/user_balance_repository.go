package repository

import (
	"context"
	"github.com/irvankadhafi/user-balance-transfer-service/internal/model"
	"github.com/irvankadhafi/user-balance-transfer-service/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type userBalanceRepository struct {
	db *gorm.DB
}

func NewUserBalanceRepository(
	db *gorm.DB,
) model.UserBalanceRepository {
	return &userBalanceRepository{
		db: db,
	}
}

func (u userBalanceRepository) CreateWithTransaction(ctx context.Context, tx *gorm.DB, userBalance *model.UserBalance) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":         utils.DumpIncomingContext(ctx),
		"userBalance": utils.Dump(userBalance),
	})

	err := tx.WithContext(ctx).Create(userBalance).Error
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

func (u userBalanceRepository) UpsertWithTransaction(ctx context.Context, tx *gorm.DB, userBalance *model.UserBalance) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":         utils.DumpIncomingContext(ctx),
		"userBalance": utils.Dump(userBalance),
	})
	switch {
	case userBalance.ID > 0:
		err := tx.WithContext(ctx).Updates(userBalance).Error
		if err != nil {
			logger.Error(err)
			return err
		}
	default:
		err := tx.WithContext(ctx).Create(userBalance).Error
		if err != nil {
			logger.Error(err)
			return err
		}
	}

	return nil
}

func (u userBalanceRepository) GetCurrentUserBalanceByUserID(ctx context.Context, userID int) (*model.UserBalance, error) {
	var userBalance model.UserBalance
	err := u.db.WithContext(ctx).Order("created_at desc").Take(&userBalance, "user_id = ?", userID).Error
	switch err {
	case nil:
		return &userBalance, nil
	case gorm.ErrRecordNotFound:
		return &userBalance, nil
	default:
		logrus.WithFields(logrus.Fields{
			"ctx":    utils.DumpIncomingContext(ctx),
			"userID": userID,
		}).Error(err)
		return nil, err
	}
}
