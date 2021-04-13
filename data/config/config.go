package config

import "time"

// Config - конфигурация запуска приложения.
type Config struct {
	// StartTimeout - время на запуск всех модулей приложения.
	StartTimeout time.Duration
	// ShutdownTimeout - время на остановку приложения после поступления системных сигналов.
	ShutdownTimeout time.Duration
	// ServerAddr - адрес и порт сервера.
	ServerAddr string
	// ServerTimeouts - время на обработку отдельных пользовательских запросов к серверу.
	ServerTimeouts time.Duration
	// EventBusTimeouts - время на обработку внутренних команд поступающих в шину приложения.
	EventBusTimeouts time.Duration
	// ISSMaxCons - максимальное число соединений с сервером ISS MOEX. Рекомендуется ограничиться 20.
	ISSMaxCons int
	// MongoURI - адрес сервера MongoDB.
	MongoURI string
	// MongoDB - наименование базы данных с таблицами.
	MongoDB string
}
