package repository

import (
	"github.com/irvankadhafi/user-balance-transfer-service/cacher"
	"github.com/sirupsen/logrus"
)

func storeNil(ck cacher.CacheManager, key string) {
	err := ck.StoreNil(key)
	if err != nil {
		logrus.Error(err)
	}
}
