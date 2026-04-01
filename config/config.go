package config

import (
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/framework/bunnymq"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/framework/common/logger"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/framework/fiber"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/framework/otel"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/framework/postgres"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/framework/redis"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/libraries/consumer"
)

type ApplicationConfig struct {
	AppConfig InternalAppConfig `config:",squash"`
	Postgres  postgres.Config   `config:",squash"`
	Redis     redis.Config      `config:",squash"`
	RabbitMQ  bunnymq.Config    `config:",squash"`
	Fiber     fiber.Config      `config:",squash"`
	Consumer  consumer.Config   `config:",squash"`
	Otel      otel.Config       `config:",squash"`
	Logger    logger.Config     `config:",squash"`
}

type InternalAppConfig struct {
	RedirectUrl string `config:"redirectUrl"`
}

func Postgres(app *ApplicationConfig) *postgres.Config {
	return &app.Postgres
}

func Redis(app *ApplicationConfig) *redis.Config {
	return &app.Redis
}

func RabbitMQ(app *ApplicationConfig) *bunnymq.Config {
	return &app.RabbitMQ
}

func Fiber(app *ApplicationConfig) *fiber.Config {
	return &app.Fiber
}

func Consumer(app *ApplicationConfig) *consumer.Config {
	return &app.Consumer
}

func Otel(app *ApplicationConfig) *otel.Config {
	return &app.Otel
}

func Logger(app *ApplicationConfig) *logger.Config {
	return &app.Logger
}

func InternalApp(app *ApplicationConfig) *InternalAppConfig {
	return &app.AppConfig
}
