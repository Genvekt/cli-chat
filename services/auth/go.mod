module github.com/Genvekt/cli-chat/services/auth

replace github.com/Genvekt/cli-chat/libraries/api => ../../libraries/api

replace github.com/Genvekt/cli-chat/libraries/db_client => ../../libraries/db_client

replace github.com/Genvekt/cli-chat/libraries/cache_client => ../../libraries/cache_client

replace github.com/Genvekt/cli-chat/libraries/closer => ../../libraries/closer

go 1.22.5

require (
	github.com/Genvekt/cli-chat/libraries/api v0.0.0-00010101000000-000000000000
	github.com/Genvekt/cli-chat/libraries/cache_client v0.0.0-00010101000000-000000000000
	github.com/Genvekt/cli-chat/libraries/closer v0.0.0-00010101000000-000000000000
	github.com/Genvekt/cli-chat/libraries/db_client v0.0.0-00010101000000-000000000000
	github.com/Masterminds/squirrel v1.5.4
	github.com/brianvoe/gofakeit/v7 v7.0.4
	github.com/gojuno/minimock/v3 v3.3.13
	github.com/gomodule/redigo v1.9.2
	github.com/jackc/pgx/v4 v4.18.3
	github.com/joho/godotenv v1.5.1
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.9.0
	google.golang.org/grpc v1.65.0
	google.golang.org/protobuf v1.34.2
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/georgysavva/scany v1.2.2 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.3 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgtype v1.14.0 // indirect
	github.com/jackc/puddle v1.3.0 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.23.0 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240528184218-531527333157 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
