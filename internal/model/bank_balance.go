package model

import (
	"context"
	"gorm.io/gorm"
)

// BankBalance menyimpan data saldo bank.
type BankBalance struct {
	ID             int    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Balance        int64  `json:"balance"`
	BalanceAchieve int64  `json:"balance_achieve"`
	Code           string `json:"code"`
	Enable         bool   `json:"enable"`
}

type AddBankBalanceInput struct {
	Code      string `json:"code"`
	Balance   int64  `json:"balance"`
	SessionID int    `json:"session_id"`
	Author    string `json:"author"`
}

// BankBalanceRepository menyediakan akses ke data saldo bank.
type BankBalanceRepository interface {
	CreateWithTransaction(ctx context.Context, tx *gorm.DB, bankBalance *BankBalance) error
	UpsertWithTransaction(ctx context.Context, tx *gorm.DB, bankBalance *BankBalance) error
	GetCurrentBankBalanceByCode(ctx context.Context, code string) (*BankBalance, error)
}

type BankBalanceUsecase interface {
	CreateBankAccount(ctx context.Context, sessionID int, code string) error
	AddBankBalance(ctx context.Context, input AddBankBalanceInput) error
	GetBankBalanceByID(ctx context.Context, bankBalanceID int) (*BankBalance, error)
	TransferUserBalance(ctx context.Context, userIDFrom, userIDTo int, balance float64, code string) error
}
