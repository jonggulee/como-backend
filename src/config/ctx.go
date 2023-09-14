package config

type Ctx struct {
	Cfg *Config
}

var (
	AppCtx = &Ctx{
		Cfg: &Config{},
	}
)
