package repository

import (
	"context"
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"github.com/irvankadhafi/user-balance-transfer-service/cacher"
	"github.com/irvankadhafi/user-balance-transfer-service/internal/config"
	"github.com/irvankadhafi/user-balance-transfer-service/internal/model"
	"github.com/irvankadhafi/user-balance-transfer-service/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type userRepository struct {
	db           *gorm.DB
	cacheManager cacher.CacheManager
}

func NewUserRepository(
	db *gorm.DB,
	cacheManager cacher.CacheManager,
) model.UserRepository {
	return &userRepository{
		db:           db,
		cacheManager: cacheManager,
	}
}

func (u *userRepository) Create(ctx context.Context, user *model.User) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":  utils.DumpIncomingContext(ctx),
		"user": utils.Dump(user),
	})

	err := u.db.Create(user).Error
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (u *userRepository) FindByID(ctx context.Context, id int) (*model.User, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx": utils.DumpIncomingContext(ctx),
		"id":  id,
	})

	if id <= 0 {
		return nil, nil
	}

	cacheKey := u.newCacheKeyByID(id)

	user := &model.User{}
	err := u.db.WithContext(ctx).Take(user, "id = ?", id).Error
	switch err {
	case nil:
		return user, nil
	case gorm.ErrRecordNotFound:
		storeNil(u.cacheManager, cacheKey)
		return nil, nil
	default:
		logger.Error(err)
		return nil, err
	}
}

func (u *userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":      utils.DumpIncomingContext(ctx),
		"username": username,
	})
	cacheKey := u.newUserCacheKeyByUsername(username)
	var id int
	err := u.db.Model(model.User{}).Select("id").Take(&id, "username = ?", username).Error
	switch err {
	case nil:
		err := u.cacheManager.StoreWithoutBlocking(cacher.NewItem(cacheKey, utils.ToByte(id)))
		if err != nil {
			logger.Error(err)
		}
		return u.FindByID(ctx, id)
	case gorm.ErrRecordNotFound:
		storeNil(u.cacheManager, cacheKey)
		return nil, nil
	default:
		logger.Error(err)
		return nil, err
	}
}

func (u *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":   utils.DumpIncomingContext(ctx),
		"email": email,
	})
	cacheKey := u.newUserCacheKeyByEmail(email)
	var id int
	err := u.db.Model(model.User{}).Select("id").Take(&id, "email = ?", email).Error
	switch err {
	case nil:
		err := u.cacheManager.StoreWithoutBlocking(cacher.NewItem(cacheKey, utils.ToByte(id)))
		if err != nil {
			logger.Error(err)
		}
		return u.FindByID(ctx, id)
	case gorm.ErrRecordNotFound:
		storeNil(u.cacheManager, cacheKey)
		return nil, nil
	default:
		logger.Error(err)
		return nil, err
	}
}

func (u *userRepository) FindPasswordByID(ctx context.Context, id int) ([]byte, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx": utils.DumpIncomingContext(ctx),
		"id":  id,
	})

	cacheKey := u.newPasswordCacheKeyByID(id)
	reply, mu, err := u.findStringValueFromCacheByKey(cacheKey)
	defer cacher.SafeUnlock(mu)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	if mu == nil {
		return []byte(reply), nil
	}

	var pass string
	err = u.db.WithContext(ctx).Model(model.User{}).Select("password").Take(&pass, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	err = u.cacheManager.StoreWithoutBlocking(cacher.NewItem(cacheKey, utils.ToByte(pass)))
	if err != nil {
		logger.Error(err)
	}

	return []byte(pass), err
}

// IncrementLoginByEmailPasswordRetryAttempts increment login by email and password retry attempts by one
func (u *userRepository) IncrementLoginByEmailPasswordRetryAttempts(ctx context.Context, email string) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":   utils.DumpIncomingContext(ctx),
		"email": email,
	})

	key := u.newLoginByEmailPasswordAttemptsCacheKeyByEmail(email)
	if err := u.cacheManager.IncreaseCachedValueByOne(key); err != nil {
		logger.Error(err)
		return err
	}

	// resets the ttl duration everytime the attempts is incremented
	if err := u.cacheManager.Expire(key, config.LoginByUsernamePasswordLockTTL()); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (u *userRepository) IsLoginByEmailPasswordLocked(ctx context.Context, email string) (bool, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":   utils.DumpIncomingContext(ctx),
		"email": email,
	})

	key := u.newLoginByEmailPasswordAttemptsCacheKeyByEmail(email)
	ttl, err := u.cacheManager.GetTTL(key)
	if err != nil {
		logger.Error(err)
		return false, err
	}

	loginAttempts, mu, err := u.findIntValueFromCacheByKey(key)
	defer cacher.SafeUnlock(mu)
	if err != nil {
		logger.Error(err)
		return false, err
	}

	if ttl > int64(0) && loginAttempts >= config.LoginByUsernamePasswordRetryAttempts() {
		return true, nil
	}

	return false, nil
}
func (u *userRepository) findFromCacheByKey(key string) (reply *model.User, mu *redsync.Mutex, err error) {
	var rep interface{}
	rep, mu, err = u.cacheManager.GetOrLock(key)
	if err != nil || rep == nil {
		return
	}

	reply = utils.InterfaceBytesToType[*model.User](rep)
	return
}

func (u *userRepository) findIntValueFromCacheByKey(key string) (reply int, mu *redsync.Mutex, err error) {
	var rep interface{}
	rep, mu, err = u.cacheManager.GetOrLock(key)
	if err != nil || rep == nil {
		return
	}

	reply = utils.InterfaceBytesToType[int](rep)
	return
}

func (u *userRepository) findStringValueFromCacheByKey(key string) (reply string, mu *redsync.Mutex, err error) {
	var rep interface{}
	rep, mu, err = u.cacheManager.GetOrLock(key)
	if err != nil || rep == nil {
		return
	}

	reply = utils.InterfaceBytesToType[string](rep)
	return
}

func (u *userRepository) newLoginByEmailPasswordAttemptsCacheKeyByEmail(email string) string {
	return fmt.Sprintf("cache:login_attempts:email_password:user_email:%s", email)
}

func (u *userRepository) newCacheKeyByID(id int) string {
	return fmt.Sprintf("cache:object:user:id:%d", id)
}

func (u *userRepository) newUserCacheKeyByEmail(email string) string {
	return fmt.Sprintf("cache:id:user_email:%s", email)
}
func (u *userRepository) newUserCacheKeyByUsername(username string) string {
	return fmt.Sprintf("cache:id:username:%s", username)
}

func (u *userRepository) newPasswordCacheKeyByID(id int) string {
	return fmt.Sprintf("cache:password:id:%d", id)
}
