package usecase

import (
	"context"
	"errors"
	"github.com/irvankadhafi/user-balance-transfer-service/internal/config"
	"github.com/irvankadhafi/user-balance-transfer-service/internal/model"
	"github.com/irvankadhafi/user-balance-transfer-service/utils"
	"strconv"
	"strings"
	"time"
)

// generateToken and check uniqueness
func generateToken(sr model.SessionRepository, userID int) (token string, err error) {
	sleep := 10 * time.Millisecond
	ctxTimeout := 50 * time.Millisecond
	err = utils.Retry(3, sleep, func() error {
		sb := strings.Builder{}
		sb.WriteString(strconv.Itoa(userID))
		sb.WriteString("_")

		randomAlphanum := utils.GenerateRandomAlphanumeric(config.DefaultSessionTokenLength)
		sb.WriteString(randomAlphanum)
		token = sb.String()

		ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
		defer cancel()

		exist, err := sr.CheckToken(ctx, token)
		if err != nil {
			return err
		}
		if exist {
			return errors.New("token exists, retry")
		}

		return nil
	})

	return token, err
}
