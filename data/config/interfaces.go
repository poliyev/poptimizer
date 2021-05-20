package config

import "context"

// Module - составная часть приложения.
//
// Модуль запускается в начале работы приложения, останавливается в конце и имеет имя для облегчения логирования.
type Module interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}
