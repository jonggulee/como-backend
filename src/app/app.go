package app

import (
	"net/http"

	"github.com/jonggulee/go-login.git/src/api/controller"
	"github.com/jonggulee/go-login.git/src/config"
	"github.com/jonggulee/go-login.git/src/constants"
	"github.com/jonggulee/go-login.git/src/logger"
	"github.com/urfave/cli/v2"
)

var (
	App = &cli.App{}
)

func init() {
	App.Commands = []*cli.Command{
		&RunCommand,
	}
}

func Configure(ctx *cli.Context, cfg *config.Config) (*config.Config, error) {

	readListenPort(ctx, cfg)

	readKakaoLoginClientId(ctx, cfg)
	readKakaoLoginClientSecret(ctx, cfg)

	return cfg, nil
}

func RunFunc(ctx *cli.Context) error {
	cfg, err := Configure(ctx, config.AppCtx.Cfg)
	if err != nil {
		return err
	}

	printConfig(cfg)

	router := controller.NewRouter()

	logger.Debugf(constants.NOREQID, "%s", http.ListenAndServe(":8080", router))

	return nil
}
