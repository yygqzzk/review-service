package data

import (
	"fmt"
	"strings"

	"github.com/yygqzzk/review-service/internal/conf"
	"github.com/yygqzzk/review-service/internal/data/query"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewDB, NewData, NewGreeterRepo, NewReviewRepo)

// Data .
type Data struct {
	// TODO wrapped database client
	dbClient *query.Query
	log      *log.Helper
}

// NewData .
func NewData(db *gorm.DB, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	// 设置默认的数据库连接对象
	query.SetDefault(db)

	return &Data{
		dbClient: query.Q,
		log:      log.NewHelper(logger),
	}, cleanup, nil
}

func NewDB(cfg *conf.Data) (*gorm.DB, error) {
	switch strings.ToLower(cfg.Database.GetDriver()) {
	case "mysql":
		db, err := gorm.Open(mysql.Open(cfg.Database.GetSource()))
		return db, err
	default:
		panic(fmt.Errorf("unsupported driver: %s", cfg.Database.GetDriver()))
	}
}
