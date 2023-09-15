package config

var (
	AppCtx = &Ctx{
		Cfg: &Config{},
	}
)

type Ctx struct {
	Cfg *Config
}
