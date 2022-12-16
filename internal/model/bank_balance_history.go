package model

import "context"

// BankBalanceHistory menyimpan data riwayat saldo bank.
type BankBalanceHistory struct {
	ID            int
	BankBalanceID int
	BalanceBefore float64
	BalanceAfter  float64
	Activity      string
	Type          string
	IP            string
	Location      string
	UserAgent     string
	Author        string
}

// BankBalanceHistoryRepository menyediakan akses ke data riwayat saldo bank.
type BankBalanceHistoryRepository interface {
	// AddBankBalanceHistory menambahkan data riwayat saldo bank baru.
	AddBankBalanceHistory(ctx context.Context, bankBalanceHistory *BankBalanceHistory) error

	// UpdateBankBalanceHistory mengubah data riwayat saldo bank.
	UpdateBankBalanceHistory(ctx context.Context, bankBalanceHistory *BankBalanceHistory) error

	// GetBankBalanceHistoryByID mengambil data riwayat saldo bank berdasarkan ID.
	GetBankBalanceHistoryByID(ctx context.Context, id int) (*BankBalanceHistory, error)

	// DeleteBankBalanceHistory menghapus data riwayat saldo bank berdasarkan ID.
	DeleteBankBalanceHistory(ctx context.Context, id int) error
}

type BankBalanceHistoryUsecase interface {
	AddBankBalanceHistory(ctx context.Context, bankBalanceID int, balanceBefore, balanceAfter float64, activity, transactionType, ip, location, userAgent, author string) (*BankBalanceHistory, error)
	GetBankBalanceHistoryByID(ctx context.Context, bankBalanceHistoryID int) (*BankBalanceHistory, error)
	DeleteBankBalanceHistory(ctx context.Context, bankBalanceID int) error
}
