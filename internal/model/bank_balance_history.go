package model

import (
	"context"
	"gorm.io/gorm"
)

// BankBalanceHistory menyimpan data riwayat saldo bank.
type BankBalanceHistory struct {
	ID            int             `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	BankBalanceID int             `json:"bank_balance_id"`
	BalanceBefore int64           `json:"balance_before"`
	BalanceAfter  int64           `json:"balance_after"`
	Activity      string          `json:"activity"`
	Type          TransactionType `json:"type"`
	IPAddress     string          `json:"ip_address"`
	Location      string          `json:"location"`
	UserAgent     string          `json:"user_agent"`
	Author        string          `json:"author"`
}

// BankBalanceHistoryRepository menyediakan akses ke data riwayat saldo bank.
type BankBalanceHistoryRepository interface {
	CreateWithTransaction(ctx context.Context, tx *gorm.DB, input *BankBalanceHistory) error
}
