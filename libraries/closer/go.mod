module github.com/Genvekt/cli-chat/libraries/closer

replace github.com/Genvekt/cli-chat/libraries/logger => ../../libraries/logger

go 1.22.5

require github.com/Genvekt/cli-chat/libraries/logger v0.0.0-00010101000000-000000000000

require (
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
)
