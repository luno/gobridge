# gobridge

![alt text](example/screenshots/how_to_configure.png)
# What is GoBridge
#### GoBridge is a client code generation for HTTP in situations where gRPC is not desired. Currently Angular (typescipt) is only supported client application.

# How to use it:
### 1. Copy the binary
### 2. Add a //go:generate or bash script with the following command
```shell script
./{{PATH_TO_BINARY}}/gobridge --api="./{{YOUR_GO_FILE}}.go" --mod="{{YOUR_GO_MODULE_NAME}}" --ts="{{PATH_AND_NAME_OF_ANGULAR_SERVICE}}" --ts_service="{{ANGULAR_SERVICE_NAME}}" --server="{{PATH_TO_GENERATED_SERVER_CODE}}"
```

# Tutorial:
#### 1. Clone the repo
git clone https://github.com/andrewwormald/gobridge.git

#### 2. Build the binary
```shell script
go build -o bin/gobridge main.go
```

#### 3. Run it with filling out the relevant information in the flags
```shell script
./bin/gobridge --api="./example/backend/example.go" --mod="gobridge" --ts="./frontend/services/example.ts" --ts_service="Example" --server="./backend/example/server/server_gen.go"
```

#### 4. It will take delcarations like this:
![alt text](example/screenshots/how_to_configure.png)

#### 5. And generate this with the supporting server side code
![alt text](example/screenshots/ts_output.png)