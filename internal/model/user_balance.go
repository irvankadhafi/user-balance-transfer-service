package model

import (
	"context"
	"gorm.io/gorm"
	"time"
)

// UserBalance menyimpan data saldo user.
type UserBalance struct {
	ID             int       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	UserID         int       `json:"user_id"`
	Balance        int64     `json:"balance"`
	BalanceAchieve int64     `json:"balance_achieve"`
	CreatedAt      time.Time `json:"created_at" sql:"DEFAULT:'now()':::STRING::TIMESTAMP" gorm:"->;<-:create"`
}

type AddUserBalanceInput struct {
	UserID    int    `json:"user_id"`
	Balance   int64  `json:"balance"`
	SessionID int    `json:"session_id"`
	Author    string `json:"author"`
}

type TransferUserBalanceInput struct {
	FromUserID int    `json:"from_user_id"`
	ToUserID   int    `json:"to_user_id"`
	Balance    int64  `json:"balance"`
	SessionID  int    `json:"session_id"`
	Author     string `json:"author"`
}

// UserBalanceRepository menyediakan akses ke data saldo user.
type UserBalanceRepository interface {
	CreateWithTransaction(ctx context.Context, tx *gorm.DB, userBalance *UserBalance) error
	UpsertWithTransaction(ctx context.Context, tx *gorm.DB, userBalance *UserBalance) error
	GetCurrentUserBalanceByUserID(ctx context.Context, userID int) (*UserBalance, error)
}

// UserBalanceUsecase menyediakan fungsi-fungsi yang berkaitan dengan model UserBalance.
type UserBalanceUsecase interface {
	AddUserBalance(ctx context.Context, input AddUserBalanceInput) error
	TransferUserBalance(ctx context.Context, input TransferUserBalanceInput) error
	GetCurrentUserBalanceByUserID(ctx context.Context, userID int) (*UserBalance, error)
}
