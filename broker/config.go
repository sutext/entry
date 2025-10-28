package broker

type Config struct {
	Addr  string
	Grpc  GrpcConfig
	Redis RedisConfig
}

type RedisConfig struct {
	DB       int
	Addr     string
	Username string
	Password string
}

type GrpcConfig struct {
	Addr string
}

func DefaultConfig() *Config {
	return &Config{
		Addr: ":8080",
		Grpc: GrpcConfig{
			Addr: ":9090",
		},
		Redis: RedisConfig{
			Addr:     "localhost:6379",
			Username: "",
			Password: "",
			DB:       0,
		},
	}
}
