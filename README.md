# MP1
***

# Distributed Log Querier


Program to query distributed log files on multiple machines.

- Uses GRPCs to communicate between servers
- Client, Coordinator & Worker Architecture
- Produces results of grep -[E]c command 

#### Building for source


Generating executable from source code:

```sh
go build -o client
```

Running source file:

```sh
go run client.go grep -c "test pattern" .log
```

Running executable:
```sh
./client grep -c "test pattern" .log
```

## Execution


Open Terminal and run these commands.

To run worker (on all VMs):
```sh
cd cs-425-mp1/src/cmd/worker/
```


```sh
./worker_server
```

To run coordinator (typically on VM2):

```sh
cd cs-425-mp1/src/cmd/coordinator/
```


```sh
./coordinator_server
```

To run client (typically on VM1):

```sh
cd cs-425-mp1/src/cmd/client/
```


```sh
./client grep -c "test pattern" .log
```

To run client tests:

```sh
cd cs-425-mp1/src/cmd/client/
```
```sh
./client_tests
```

## License

UIUC
