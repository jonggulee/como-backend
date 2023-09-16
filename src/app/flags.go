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

	// OAuth Kakao
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

	// DB
	ConfigDbAddressFlag = cli.StringFlag{
		Name:    "dbAddress",
		Usage:   "set DB Address",
		Aliases: []string{"d"},
	}

	ConfigDbPortFlag = cli.IntFlag{
		Name:    "dbPort",
		Usage:   "set DB Port",
		Aliases: []string{"t"},
		Value:   3306,
	}

	ConfigDbUserFlag = cli.StringFlag{
		Name:    "dbUser",
		Usage:   "set DB User",
		Aliases: []string{"u"},
	}

	ConfigDbPasswordFlag = cli.StringFlag{
		Name:    "dbPassword",
		Usage:   "set DB Password",
		Aliases: []string{"w"},
	}

	ConfigDbNameFlag = cli.StringFlag{
		Name:    "dbName",
		Usage:   "set Db Name",
		Aliases: []string{"n"},
	}
)
