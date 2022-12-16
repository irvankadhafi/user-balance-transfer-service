package model

import "context"

type TransactionType string

const (
	DEBIT  TransactionType = "DEBIT"
	CREDIT TransactionType = "CREDIT"
)

// UserBalanceHistory menyimpan data riwayat saldo user.
type UserBalanceHistory struct {
	ID            int
	UserBalanceID int
	BalanceBefore float64
	BalanceAfter  float64
	Activity      string
	Type          TransactionType
	IP            string
	Location      string
	UserAgent     string
	Author        string
}

// UserBalanceHistoryRepository menyediakan akses ke data riwayat saldo user.
type UserBalanceHistoryRepository interface {
	// AddUserBalanceHistory menambahkan data riwayat saldo user baru.
	AddUserBalanceHistory(ctx context.Context, userBalanceHistory *UserBalanceHistory) error

	// GetUserBalanceHistoryByID mengambil data riwayat saldo user berdasarkan ID.
	GetUserBalanceHistoryByID(ctx context.Context, id int) (*UserBalanceHistory, error)

	// DeleteUserBalanceHistory menghapus data riwayat saldo user berdasarkan ID.
	DeleteUserBalanceHistory(ctx context.Context, id int) error
}

// UserBalanceHistoryUsecase menyediakan fungsi-fungsi yang berkaitan dengan model UserBalanceHistory.
type UserBalanceHistoryUsecase interface {
	AddUserBalanceHistory(ctx context.Context, userBalanceID int, balanceBefore, balanceAfter float64, activity, transactionType, ip, location, userAgent, author string) (*UserBalanceHistory, error)
	GetUserBalanceHistoryByID(ctx context.Context, userBalanceHistoryID int) (*UserBalanceHistory, error)
	DeleteUserBalanceHistory(ctx context.Context, userBalanceID int) error
}
