package config

import (
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
	v.AutomaticEnv()

	v.SetDefault("server.addr", ":8080")
	v.SetDefault("server.prefork", true)
	v.SetDefault("server.readTimeout", "5s")
	v.SetDefault("server.writeTimeout", "10s")

	if err := v.ReadInConfig(); err != nil {
		// allow env-only
	}

	c := &Config{}
	if err := v.Unmarshal(&c); err != nil {
		return nil, err
	}
	// durations from string
	c.Server.ReadTimeout = v.GetDuration("server.readTimeout")
	c.Server.WriteTimeout = v.GetDuration("server.writeTimeout")
	c.DB.MaxConnLifetime = v.GetDuration("db.maxConnLifetime")
	c.DB.MaxConnIdleTime = v.GetDuration("db.maxConnIdleTime")
	c.Auth.AccessTTL = v.GetDuration("auth.accessTTL")
	c.Auth.JWTSecret = v.GetString("auth.jwtSecret")
	return c, nil
}
