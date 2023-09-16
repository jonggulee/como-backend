package app

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jonggulee/go-login.git/src/api/controller"
	"github.com/jonggulee/go-login.git/src/config"
	"github.com/jonggulee/go-login.git/src/constants"
	dbController "github.com/jonggulee/go-login.git/src/db/controller"
	"github.com/jonggulee/go-login.git/src/logger"
	"github.com/urfave/cli/v2"
)

var (
	App = &cli.App{}
)

func init() {
	App.Commands = []*cli.Command{
		&RunCommand,
		&VersionCommand,
	}
}

func Configure(ctx *cli.Context, cfg *config.Config) (*config.Config, error) {
	readListenPort(ctx, cfg)

	// OAuth Kakao
	readKakaoLoginClientId(ctx, cfg)
	readKakaoLoginClientSecret(ctx, cfg)

	// DB
	readDbAddress(ctx, cfg)
	readDbPort(ctx, cfg)
	readDbName(ctx, cfg)
	readDbUser(ctx, cfg)
	readDbPassword(ctx, cfg)

	return cfg, nil
}

func VersionFunc(ctx *cli.Context) error {
	fmt.Fprintf(os.Stdout, "%s-%s\n", constants.APPNAME, constants.APPVERSION)
	return nil
}

func RunFunc(ctx *cli.Context) error {
	cfg, err := Configure(ctx, config.AppCtx.Cfg)
	if err != nil {
		return err
	}

	printConfig(cfg)

	controller.ReadKakaoConfig(cfg)
	dbConn, err := dbController.Connection(cfg.DbName, cfg.DbUser, cfg.DbPassword, cfg.DbAddress, cfg.DbPort)
	if err != nil {
		logger.Errorf(constants.NOREQID, "Failed to connect DB, %s", err)
		os.Exit(1)
	}

	config.AppCtx.Db = dbConn

	router := controller.NewRouter()

	logger.Debugf(constants.NOREQID, "%s", http.ListenAndServe(fmt.Sprintf(":%d", cfg.ListenPort), router))

	return nil
}
