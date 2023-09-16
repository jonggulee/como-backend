package app

import (
	"fmt"
	"os"
	"time"

	"github.com/jonggulee/go-login.git/src/config"
	"github.com/jonggulee/go-login.git/src/constants"
	"github.com/jonggulee/go-login.git/src/logger"
	"github.com/urfave/cli/v2"
)

func readKakaoLoginClientId(ctx *cli.Context, cfg *config.Config) {
	cfg.KakaoClientId = ctx.String("kakaoClientId")

	if cfg.KakaoClientId == "" {
		logger.Logger.Fatal("Kakao Client Id is empty")
		os.Exit(1)
	}
}

func readKakaoLoginClientSecret(ctx *cli.Context, cfg *config.Config) {
	cfg.KakaoClientSecret = ctx.String("kakaoClientSecret")

	if cfg.KakaoClientSecret == "" {
		logger.Logger.Fatal("Kakao Client Secret is empty")
		os.Exit(1)
	}
}

func readListenPort(ctx *cli.Context, cfg *config.Config) error {
	port := ctx.Int("listenPort")

	if port <= 0 {
		port = 8080
	}

	cfg.ListenPort = port
	return nil
}

func printConfig(cfg *config.Config) {
	fmt.Println("==================================================")
	fmt.Println("                COMO API SERVER")
	fmt.Println("==================================================")
	fmt.Printf("Server is starting...\n\n")
	fmt.Printf("Version:		%s\n", constants.APPVERSION)
	fmt.Printf("Listening port: 	%d\n", cfg.ListenPort)
	fmt.Printf("Started at: 		%s\n", time.Now().Format(time.RFC3339))
	fmt.Println("==================================================")
	logger.Debugf(constants.NOREQID, "KakaoClientId: %s", cfg.KakaoClientId)
	logger.Debugf(constants.NOREQID, "KakaoClientSecret: %s", cfg.KakaoClientSecret)
}
