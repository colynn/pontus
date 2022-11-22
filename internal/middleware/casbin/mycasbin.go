package mycasbin

import (
	"github.com/colynn/pontus/config"
	"github.com/colynn/pontus/internal/db"
	"github.com/colynn/pontus/tools"

	log "unknwon.dev/clog/v2"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/go-kit/kit/endpoint"
	_ "github.com/go-sql-driver/mysql"
)

// Em ..
var Em endpoint.Middleware

// NewCasbin ..
func NewCasbin() (*casbin.Enforcer, error) {
	c := config.GetConfig()
	host := c.GetString("database.host")
	port := c.GetInt("database.port")
	dbType := c.GetString("database.dbtype")
	database := c.GetString("database.database")
	username := c.GetString("database.username")
	password := c.GetString("database.password")

	conn := db.GetMysqlConn(username, password, host, database, port)
	// TODO: current dbType only support mysql
	Apter, err := gormadapter.NewAdapter(dbType, conn, true)
	if err != nil {
		return nil, err
	}

	rbacConf := tools.EnsureAbs("config/rbac_model.conf")
	e, err := casbin.NewEnforcer(rbacConf, Apter)
	if err != nil {
		log.Error("casbin new enforcer error: %s", err.Error())
		return nil, err
	}
	if err := e.LoadPolicy(); err == nil {
		return e, err
	}
	e.AddPolicy()
	log.Error("casbin rbac_model or policy init error, message: %v", err)
	return nil, err
}
