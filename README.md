# gobridge
##### GoBridge is not production ready

#### What is GoBridge
GoBridge is a naive attempt at doing client code generation for HTTP in situations where gRPC is not desired.

#### Supported client side languages:
##### Beta:
1. Typescript

##### Roadmap:
1. Swift

#### Example: 
```shell script
./bin/gobridge --go_target="./example/example.go" --go_mod="gobridge" --ts_output="./frontend/services/example.ts"
```