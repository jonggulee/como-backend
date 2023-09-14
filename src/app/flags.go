package app

import (
	"github.com/urfave/cli/v2"
)

var (
	ConfigListenPort = cli.IntFlag{
		Name:    "listenPort",
		Usage:   "set listening port",
		Aliases: []string{"p"},
		Value:   8080,
	}

	// kakao
	ConfigKakaoClientId = cli.StringFlag{
		Name:    "kakaoClientId",
		Usage:   "Kakao Client Id",
		Aliases: []string{"ki"},
	}

	ConfigKakaoClientSecret = cli.StringFlag{
		Name:    "kakaoClientSecret",
		Usage:   "Kakao Client Secret",
		Aliases: []string{"ks"},
	}
)
