module github.com/Genvekt/cli-chat/services/chat-client

replace github.com/Genvekt/cli-chat/libraries/api => ../../libraries/api

replace github.com/Genvekt/cli-chat/libraries/closer => ../../libraries/closer

replace github.com/Genvekt/cli-chat/libraries/logger => ../../libraries/logger

go 1.22.5

require (
	github.com/Genvekt/cli-chat/libraries/api v0.0.0-00010101000000-000000000000
	github.com/Genvekt/cli-chat/libraries/closer v0.0.0-00010101000000-000000000000
	github.com/Genvekt/cli-chat/libraries/logger v0.0.0-00010101000000-000000000000
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/joho/godotenv v1.5.1
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/spf13/cobra v1.8.1
	go.uber.org/zap v1.27.0
	google.golang.org/grpc v1.66.2
	google.golang.org/protobuf v1.34.2
	gopkg.in/ini.v1 v1.67.0
)

require (
	github.com/BurntSushi/toml v1.4.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240723171418-e6d459c13d2a // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
