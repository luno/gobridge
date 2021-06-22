# gobridge

![alt text](example/screenshots/how_to_configure.png)

# What is GoBridge
#### GoBridge is a client code generation for HTTP in situations where gRPC is not desired. Currently Angular (typescipt) is only supported client application.

# Tutorial:
#### 1. Clone the repo or copy the binary directly from ./bin/
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

#### GoBridge will output this for Angular (typescript)
![alt text](example/screenshots/ts_output.png)

#### GoBridge will output this as the server side implementation
![alt text](example/screenshots/server_side_code.png)