package model

import "context"

// UserBalance menyimpan data saldo user.
type UserBalance struct {
	ID             int
	UserID         int
	Balance        int64
	BalanceAchieve int64
}

// UserBalanceRepository menyediakan akses ke data saldo user.
type UserBalanceRepository interface {
	// AddUserBalance menambahkan data saldo user baru.
	AddUserBalance(ctx context.Context, userBalance *UserBalance) error

	// FindUserBalanceByID mengambil data saldo user berdasarkan ID.
	FindUserBalanceByID(ctx context.Context, id int) (*UserBalance, error)
}

// UserBalanceUsecase menyediakan fungsi-fungsi yang berkaitan dengan model UserBalance.
type UserBalanceUsecase interface {
	AddUserBalance(ctx context.Context, userID int, balance, balanceAchieve float64) (*UserBalance, error)
	GetUserBalanceByID(ctx context.Context, userBalanceID int) (*UserBalance, error)
}
