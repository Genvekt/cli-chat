# API library
API library contains `.proto` files which declare gRPC API contracts used by services in the system.
Additionally, there are tools for `.pb.go` generation.

## Structure
Library is designed according to the following structure:
```
api/
├── bin/                        - executables
│ 
├── api_1/                      - folder for API #1
│   ├── v1/                     - 1 version of API #1
│   │   ├── some_file_1.proto   - proto files
│   │   ├── some_file_2.proto   
│   │   ├── some_file_1.pb.go   - generated .go files
│   │   └── some_file_1.pb.go
│   └── v2/                     - 2 version of API #1
│       └── ...
├── api_2/                      - folder for API #2
│   └── ...
│ 
├── Makefile                    - scripts
└── go.mod                      - dependencies
```

## How to use
### Install dependencies
To install dependencies required for code generation use the following commands:

```shell
make install-deps
make get-deps
```

It will create local `./bin` directory with protoc-gen-go executables and install all required go packages

### Generate go code
To generate `.pb.go` files after change in `.proto` files, use the following command:
```shell
make generate
```

### Use in service
To use generated code in service:
1. Paste the following library reference into `go.mod`
```
replace github.com/Genvekt/cli-chat/libraries/api => ../../libraries/api
```
2. Import api structs with 
```
import "github.com/Genvekt/cli-chat/libraries/api/<API>/<VERSION>"
```
3. Use structs from generated files!


## Versioning and expanding
To add new api contract or contract version, add required folders according to structure specification and
update Makefile with new command, expanding `generate` command.

## Deploy note
Due to monorepo project structure it is important to deploy only changed services
and services that use changed libraries.

However `api` library contains reference to several APIs and may trigger unnecessary deploys.

It is recommended to set deploy trigger to change of specific APIs in CI/CD, not for the whole library.