# SDFS
***

# Distributed File System


Program to create a fault-tolerant distributed file service across distributed processes

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
