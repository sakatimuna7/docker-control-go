package configs

import (
	logger "docker-control-go/src/log"
	"log"

	"github.com/casbin/casbin/v2"
	xormadapter "github.com/casbin/xorm-adapter/v3"
	"xorm.io/xorm"
)

var Enforcer *casbin.Enforcer

func InitCasbin(db *xorm.Engine) {
	adapter, err := xormadapter.NewAdapterByEngine(db)
	if err != nil {
		log.Fatalf("Failed to create Casbin adapter: %v", err)
		logger.Log.Fatalf("Failed to create Casbin adapter: %v", err)
	}

	Enforcer, err = casbin.NewEnforcer("src/policies/model.conf", adapter)
	if err != nil {
		log.Fatalf("Failed to create Casbin enforcer: %v", err)
		logger.Log.Fatalf("Failed to create Casbin enforcer: %v", err)
	}

	// Load policies
	err = Enforcer.LoadPolicy()
	if err != nil {
		log.Fatalf("Failed to load Casbin policies: %v", err)
		logger.Log.Fatalf("Failed to load Casbin policies: %v", err)
	}

	log.Println("Casbin initialized successfully!")
	logger.Log.Info("Casbin initialized successfully!")
}
