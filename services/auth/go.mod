module github.com/Genvekt/cli-chat/services/auth

replace github.com/Genvekt/cli-chat/libraries/api => ../../libraries/api

replace github.com/Genvekt/cli-chat/libraries/db_client => ../../libraries/db_client

replace github.com/Genvekt/cli-chat/libraries/cache_client => ../../libraries/cache_client

replace github.com/Genvekt/cli-chat/libraries/closer => ../../libraries/closer

replace github.com/Genvekt/cli-chat/libraries/kafka => ../../libraries/kafka

go 1.22.5

require (
	github.com/Genvekt/cli-chat/libraries/api v0.0.0-00010101000000-000000000000
	github.com/Genvekt/cli-chat/libraries/cache_client v0.0.0-00010101000000-000000000000
	github.com/Genvekt/cli-chat/libraries/closer v0.0.0-00010101000000-000000000000
	github.com/Genvekt/cli-chat/libraries/db_client v0.0.0-00010101000000-000000000000
	github.com/Genvekt/cli-chat/libraries/kafka v0.0.0-00010101000000-000000000000
	github.com/IBM/sarama v1.43.2
	github.com/Masterminds/squirrel v1.5.4
	github.com/brianvoe/gofakeit/v7 v7.0.4
	github.com/gojuno/minimock/v3 v3.3.13
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/gomodule/redigo v1.9.2
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.21.0
	github.com/jackc/pgx/v4 v4.18.3
	github.com/joho/godotenv v1.5.1
	github.com/pkg/errors v0.9.1
	github.com/rakyll/statik v0.1.7
	github.com/rs/cors v1.11.0
	github.com/stretchr/testify v1.9.0
	golang.org/x/crypto v0.24.0
	google.golang.org/grpc v1.65.0
	google.golang.org/protobuf v1.34.2
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/eapache/go-resiliency v1.6.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20230731223053-c322873962e3 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.0.4 // indirect
	github.com/georgysavva/scany v1.2.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.3 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgtype v1.14.0 // indirect
	github.com/jackc/puddle v1.3.0 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.4 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/klauspost/compress v1.17.8 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/pierrec/lz4/v4 v4.1.21 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240723171418-e6d459c13d2a // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240723171418-e6d459c13d2a // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
