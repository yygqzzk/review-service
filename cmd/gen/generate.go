package main

// gorm gen configure

import (
	"flag"
	"fmt"
	"strings"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/yygqzzk/review-service/internal/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"gorm.io/gen"
)

var (
	flagconf string
)

func connectDB(cfg *conf.Data_Database) *gorm.DB {
	if cfg == nil {
		panic("database config is nil")
	}
	switch strings.ToLower(cfg.GetDriver()) {
	case "mysql":
		db, err := gorm.Open(mysql.Open(cfg.GetSource()))
		if err != nil {
			panic(fmt.Errorf("connect db fail: %w", err))
		}
		return db
	default:
		panic(fmt.Errorf("unsupported driver: %s", cfg.GetDriver()))
	}
}

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func main() {
	// 从配置文件中读取数据源配置信息
	flag.Parse()

	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	// 指定生成代码的具体相对目录(相对当前文件)，默认为：./query
	// 默认生成需要使用WithContext之后才可以查询的代码，但可以通过设置gen.WithoutContext禁用该模式
	g := gen.NewGenerator(gen.Config{
		// 默认会在 OutPath 目录生成CRUD代码，并且同目录下生成 model 包
		// 所以OutPath最终package不能设置为model，在有数据库表同步的情况下会产生冲突
		// 若一定要使用可以通过ModelPkgPath单独指定model package的名称
		OutPath: "../../internal/data/query",
		/* ModelPkgPath: "dal/model"*/

		// gen.WithoutContext：禁用WithContext模式
		// gen.WithDefaultQuery：生成一个全局Query对象Q
		// gen.WithQueryInterface：生成Query接口
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true, // 允许字段为null时, 自动生成为指针类型
	})

	// 通常复用项目中已有的SQL连接配置db(*gorm.DB)
	// 非必需，但如果需要复用连接时的gorm.Config或需要连接数据库同步表信息则必须设置
	g.UseDB(connectDB(bc.Data.Database))

	// 从连接的数据库为所有表生成Model结构体和CRUD代码
	// 也可以手动指定需要生成代码的数据表
	g.ApplyBasic(g.GenerateAllTable()...)

	// 执行并生成代码
	g.Execute()
}
