package config

import (
	"github.com/gofiber/fiber/v2/log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Server struct {
	Addr         string
	Prefork      bool
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DB struct {
	DSN             string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

type Auth struct {
	JWTSecret string
	AccessTTL time.Duration
}

type Logger struct {
	Level string
}

type Config struct {
	Server Server
	DB     DB
	Auth   Auth
	Logger Logger
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")

	// Включаем ENV и маппинг "a.b" -> "A_B"
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Defaults
	v.SetDefault("server.addr", ":8080")
	v.SetDefault("server.prefork", true)
	v.SetDefault("server.readTimeout", "5s")
	v.SetDefault("server.writeTimeout", "10s")

	// Файл опционален — при отсутствии используем ENV/дефолты
	if err := v.ReadInConfig(); err != nil {
		log.Warnf("config file not found, will use env/defaults: %v", err)
	}

	c := &Config{}
	// Можно распарсить в структуру, но чувствительные поля извлекаем явно
	if err := v.Unmarshal(&c); err != nil {
		return nil, err
	}

	// Durations и точные значения из ENV/файла
	c.Server.ReadTimeout = v.GetDuration("server.readTimeout")
	c.Server.WriteTimeout = v.GetDuration("server.writeTimeout")

	c.DB.DSN = v.GetString("db.dsn")
	c.DB.MaxConns = int32(v.GetInt("db.maxConns"))
	c.DB.MinConns = int32(v.GetInt("db.minConns"))
	c.DB.MaxConnLifetime = v.GetDuration("db.maxConnLifetime")
	c.DB.MaxConnIdleTime = v.GetDuration("db.maxConnIdleTime")

	c.Auth.JWTSecret = v.GetString("auth.jwtSecret")
	c.Auth.AccessTTL = v.GetDuration("auth.accessTTL")

	c.Logger.Level = v.GetString("logger.level")
	return c, nil
}
