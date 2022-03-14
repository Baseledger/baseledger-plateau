module github.com/Baseledger/baseledger

go 1.16

require (
	github.com/cosmos/cosmos-sdk v0.44.5
	github.com/gogo/protobuf v1.3.3
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/spf13/cast v1.4.1
	github.com/spf13/cobra v1.2.1
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/starport v0.19.2
	github.com/tendermint/tendermint v0.34.14
	github.com/tendermint/tm-db v0.6.4
	google.golang.org/genproto v0.0.0-20220302033224-9aa15565e42a
	google.golang.org/grpc v1.44.0
)

require (
	github.com/cosmos/ibc-go v1.2.2
	github.com/golang/protobuf v1.5.2
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/google/uuid v1.3.0
	github.com/rs/zerolog v1.23.0
	github.com/spf13/viper v1.10.1
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4
	google.golang.org/grpc => google.golang.org/grpc v1.33.2
)

// this is done in cosmos-sdk go.mod as well, with comment:
// Fix upstream GHSA-h395-qcrw-5vmq vulnerability.
// TODO Remove it: https://github.com/cosmos/cosmos-sdk/issues/10409
replace github.com/gin-gonic/gin => github.com/gin-gonic/gin v1.7.0

// next 2 replacements fix lower priority alerts
replace github.com/opencontainers/runc => github.com/opencontainers/runc v1.0.3

replace github.com/opencontainers/image-spec => github.com/opencontainers/image-spec v1.0.2
