package broker

type Config struct {
	Port  string
	Grpc  GrpcConfig
	Redis RedisConfig
}

type RedisConfig struct {
	Addr     string
	Username string
	Password string
	DB       int
}
type GrpcConfig struct {
	Addr string
}

func DefaultConfig() *Config {
	return &Config{
		Port: ":8080",
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
