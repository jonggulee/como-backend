package app

import "github.com/urfave/cli/v2"

var (
	RunCommand = cli.Command{
		Action: RunFunc,
		Name:   "run",
		Usage:  "Run the server",
		Flags: []cli.Flag{
			&ConfigListenPort,

			// OAuth Kakao
			&ConfigKakaoClientId,
			&ConfigKakaoClientSecret,

			// DB
			&ConfigDbAddressFlag,
			&ConfigDbPortFlag,
			&ConfigDbNameFlag,
			&ConfigDbUserFlag,
			&ConfigDbPasswordFlag,
		},
	}
	VersionCommand = cli.Command{
		Action: VersionFunc,
		Name:   "version",
		Usage:  "print a version",
	}
)
