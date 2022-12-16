package console

import (
	"context"
	"github.com/irvankadhafi/user-balance-transfer-service/internal/db"
	"github.com/irvankadhafi/user-balance-transfer-service/internal/helper"
	"github.com/irvankadhafi/user-balance-transfer-service/internal/model"
	"github.com/irvankadhafi/user-balance-transfer-service/internal/repository"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var seedCmd = &cobra.Command{
	Use:   "seed-user",
	Short: "run seed-user",
	Long:  `This subcommand seeding user`,
	Run:   seedUser,
}

func init() {
	RootCmd.AddCommand(seedCmd)
}

func seedUser(cmd *cobra.Command, args []string) {
	// Initiate all connection like db, redis, etc
	db.InitializePostgresConn()

	userRepo := repository.NewUserRepository(db.PostgreSQL, nil)

	cipherPwd, err := helper.HashString("123456")
	if err != nil {
		logrus.Error(err)
	}

	user1 := &model.User{
		Username: "johndoe",
		Email:    "johndoe@mail.com",
		Password: cipherPwd,
	}

	err = userRepo.Create(context.Background(), user1)
	if err != nil {
		return
	}

	user2 := &model.User{
		Username: "irvan",
		Email:    "irvan@mail.com",
		Password: cipherPwd,
	}
	err = userRepo.Create(context.Background(), user2)
	if err != nil {
		return
	}

	logrus.Warn("DONE!")
}
