# SDFS
***

# Distributed File System


Program to create a fault-tolerant distributed file service across distributed processes

- Uses streaming to transfer chunks of files across members
- Introducer & Nodes Architecture

#### Building for source


Generating executable from source code:

```sh
go build -o process
```

Running source file:

```sh
go run process.go
```

Running executable:
```sh
./process
```

## Execution


Open Terminal and run these commands.

To run introducer (typically on VM2):

```sh
cd cs-425-mp1/src/cmd/introducer/
```


```sh
./introducer
```
***
To run node process (on all VMs):
```sh
cd cs-425-mp1/src/cmd/process/
```


```sh
./process
```
***
To run commands on node process (on all VMs):
```sh
GET
get <sdfsfilename> <localfilename> 
```

```sh
GET_VERSIONS
get-versions <sdfsfilename> <num_versions> <localfilename>
```
```sh
PUT
put <localfilename> <sdfsfilename>
```
```sh
DELETE
delete <sdfsfilename>
```

```sh
LS
ls <sdfsfilename>
```

```sh
STORE
```




# MP2
***

# Group Membership Service


Program to create a group membership service across distributed processes with fault detection

- Uses TCP to join network and UDP to communicate with other members
- Introducer & Nodes Architecture

## Notes

This code builds on top of the MP1 code - the new files are inside 
- src/introducer
- src/membership
- src/conf
- src/cmd/introducer
- src/cmd/process

#### Building for source


Generating executable from source code:

```sh
go build -o process
```

Running source file:

```sh
go run process.go
```

Running executable:
```sh
./process
```

## Execution


Open Terminal and run these commands.

To run introducer (typically on VM2):

```sh
cd cs-425-mp1/src/cmd/introducer/
```


```sh
./introducer
```
***
To run node process (on all VMs):
```sh
cd cs-425-mp1/src/cmd/process/
```


```sh
./process
```
***
To run commands on node process (on all VMs):
```sh
JOIN
```


```sh
LIST_MEM
```
```sh
LIST_SELF
```
```sh
LEAVE
```

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
