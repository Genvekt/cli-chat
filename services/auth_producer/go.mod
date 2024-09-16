module github.com/Genvekt/cli-chat/services/auth_producer

go 1.22.5

replace github.com/Genvekt/cli-chat/libraries/api => ../../libraries/api

replace github.com/Genvekt/cli-chat/libraries/kafka => ../../libraries/kafka

replace github.com/Genvekt/cli-chat/libraries/closer => ../../libraries/closer

replace github.com/Genvekt/cli-chat/libraries/logger => ../../libraries/logger

require (
	github.com/Genvekt/cli-chat/libraries/api v0.0.0-00010101000000-000000000000
	github.com/Genvekt/cli-chat/libraries/closer v0.0.0-00010101000000-000000000000
	github.com/Genvekt/cli-chat/libraries/kafka v0.0.0-00010101000000-000000000000
	github.com/Genvekt/cli-chat/libraries/logger v0.0.0-00010101000000-000000000000
	github.com/IBM/sarama v1.43.3
	github.com/brianvoe/gofakeit/v7 v7.0.4
	github.com/joho/godotenv v1.5.1
	github.com/natefinch/lumberjack v2.0.0+incompatible
	go.uber.org/zap v1.27.0
)

require (
	github.com/BurntSushi/toml v1.4.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/eapache/go-resiliency v1.7.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20230731223053-c322873962e3 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.0.4 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.21.0 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.4 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.21 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/crypto v0.26.0 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.23.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240723171418-e6d459c13d2a // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240723171418-e6d459c13d2a // indirect
	google.golang.org/grpc v1.65.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
)
