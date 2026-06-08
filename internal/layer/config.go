package layer

import (
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/config"
	"github.com/ilyakaznacheev/cleanenv"
)

func NewConfig(flags string) config.Config {
	var cfg config.Config

	err := cleanenv.ReadConfig(flags, &cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}
