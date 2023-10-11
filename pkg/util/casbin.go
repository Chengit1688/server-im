package util

import (
	"im/config"
	"im/pkg/logger"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"

	gormAdapter "github.com/go-admin-team/gorm-adapter/v3"
	redisWatcher "github.com/go-admin-team/redis-watcher/v2"
	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

// Initialize the model from a string.
var text = `
[request_definition]
r = sub, obj, act
[policy_definition]
p = sub, obj, act
[policy_effect]
e = some(where (p.eft == allow))
[matchers]
m = r.sub == p.sub && (keyMatch2(r.obj, p.obj) || keyMatch(r.obj, p.obj)) && (r.act == p.act || p.act == "*")
`

var (
	Enforcer *casbin.SyncedEnforcer
	once     sync.Once
)

func SetupCasbin(db *gorm.DB, _ string) *casbin.SyncedEnforcer {
	once.Do(func() {
		Apter, err := gormAdapter.NewAdapterByDBUseTableName(db, "cms", "casbin_rule")
		if err != nil && err.Error() != "invalid DDL" {
			panic(err)
		}

		m, err := model.NewModelFromString(text)
		if err != nil {
			panic(err)
		}
		Enforcer, err = casbin.NewSyncedEnforcer(m, Apter)
		if err != nil {
			panic(err)
		}
		err = Enforcer.LoadPolicy()
		if err != nil {
			panic(err)
		}
		// set redis watcher if redis config is not nil

		w, err := redisWatcher.NewWatcher(config.Config.Redis.Address, redisWatcher.WatcherOptions{
			Options: redis.Options{
				Network:  "tcp",
				Password: config.Config.Redis.Password,
			},
			Channel:    "/casbin",
			IgnoreSelf: false,
		})
		if err != nil {
			panic(err)
		}

		err = w.SetUpdateCallback(updateCallback)
		if err != nil {
			panic(err)
		}
		err = Enforcer.SetWatcher(w)
		if err != nil {
			panic(err)
		}

		Enforcer.EnableLog(true)
	})

	return Enforcer
}

func updateCallback(msg string) {
	logger.Sugar.Info("casbin updateCallback msg:", GetSelfFuncName(), msg)

	err := Enforcer.LoadPolicy()
	if err != nil {
		logger.Sugar.Error("casbin LoadPolicy err:", GetSelfFuncName(), msg)
	}
}

func GetEnforcer() *casbin.SyncedEnforcer {
	return Enforcer
}
