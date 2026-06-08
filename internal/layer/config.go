package layer

import (
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/config"
	"github.com/ilyakaznacheev/cleanenv"
)

func NewConfig(flags string) config.Config {
	var cfg config.Config

	err := cleanenv.ReadConfig(flags, &cfg)
	if err != nil {
		// Тут пусть будет паника, ибо если мы не смогли инициализировать конфиги,
		// то о какой вообще дальнейшей логике может идти речь?
		panic(err)
	}

	return cfg
}
