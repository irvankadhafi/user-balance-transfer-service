package model

import "context"

// BankBalance menyimpan data saldo bank.
type BankBalance struct {
	ID             int
	Balance        float64
	BalanceAchieve float64
	Code           string
	Enable         bool
}

// BankBalanceRepository menyediakan akses ke data saldo bank.
type BankBalanceRepository interface {
	// AddBankBalance menambahkan data saldo bank baru.
	AddBankBalance(ctx context.Context, bankBalance *BankBalance) error

	// GetBankBalanceByID mengambil data saldo bank berdasarkan ID.
	GetBankBalanceByID(ctx context.Context, id int) (*BankBalance, error)
}

type BankBalanceUsecase interface {
	AddBankBalance(ctx context.Context, balance, balanceAchieve float64, code string, enable bool) (*BankBalance, error)
	GetBankBalanceByID(ctx context.Context, bankBalanceID int) (*BankBalance, error)
	TransferUserBalance(ctx context.Context, userIDFrom, userIDTo int, balance float64, code string) error
}
