package model

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type TransactionType string

const (
	DEBIT  TransactionType = "DEBIT"
	CREDIT TransactionType = "CREDIT"
)

// UserBalanceHistory menyimpan data riwayat saldo user.
type UserBalanceHistory struct {
	ID            int             `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	UserBalanceID int             `json:"user_balance_id"`
	BalanceBefore int64           `json:"balance_before"`
	BalanceAfter  int64           `json:"balance_after"`
	Activity      string          `json:"activity"`
	Type          TransactionType `json:"type"`
	IPAddress     string          `json:"ip_address"`
	Location      string          `json:"location"`
	UserAgent     string          `json:"user_agent"`
	Author        string          `json:"author"`
	CreatedAt     time.Time       `json:"created_at" sql:"DEFAULT:'now()':::STRING::TIMESTAMP" gorm:"->;<-:create"`
}

// UserBalanceHistoryRepository menyediakan akses ke data riwayat saldo user.
type UserBalanceHistoryRepository interface {
	CreateWithTransaction(ctx context.Context, tx *gorm.DB, input *UserBalanceHistory) error
}
