package config

import "github.com/ilyakaznacheev/cleanenv"

const (
	LogErrorMsg            string = "error message: "
	LogErrorCtx            string = "Error: "
	LogErrorInitgRPCServer string = "failed create new grpc server"
	LogErrorLoadConfig     string = "errorvid load config"

	LogInitServer string = "Server started on port"

	LogShutdownServer string = "shutting down the server"
	LogServerStopped  string = "Server stopped"
)

type Config struct {
	RestServerPort int    `env:"SERVER_PORT" env-default:"8989"`
	PatternURL     string `env:"API_URL" env-default:"/api/v1/calculate"`
}

// Конструктор
func New() (*Config, error) {
	cfg := Config{}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
