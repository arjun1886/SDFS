# MP1
***

# Distributed Log Querier


Program to query distributed log files on multiple machines.

- Uses GRPCs to communicate between servers
- Client, Coordinator & Worker Architecture
- Produces results of grep -[E]c command


## Tech

Dillinger uses a number of open source projects to work properly:

- [AngularJS] - HTML enhanced for web apps!
- [Ace Editor] - awesome web-based text editor
- [markdown-it] - Markdown parser done right. Fast and easy to extend.
- [Twitter Bootstrap] - great UI boilerplate for modern web apps
- [node.js] - evented I/O for the backend
- [Express] - fast node.js network app framework [@tjholowaychuk]
- [Gulp] - the streaming build system
- [Breakdance](https://breakdance.github.io/breakdance/) - HTML
to Markdown converter
- [jQuery] - duh

And of course Dillinger itself is open source with a [public repository][dill]
 on GitHub.

## Installation

This project requires installation of Go to run.

Install the dependencies and devDependencies and start the server.

```sh
cd dillinger
npm i
node app
```

For production environments...

```sh
npm install --production
NODE_ENV=production node app
```

## Plugins

Dillinger is currently extended with the following plugins.
Instructions on how to use them in your own application are linked below.

| Plugin | README |
| ------ | ------ |
| Dropbox | [plugins/dropbox/README.md][PlDb] |
| GitHub | [plugins/github/README.md][PlGh] |
| Google Drive | [plugins/googledrive/README.md][PlGd] |
| OneDrive | [plugins/onedrive/README.md][PlOd] |
| Medium | [plugins/medium/README.md][PlMe] |
| Google Analytics | [plugins/googleanalytics/README.md][PlGa] |

## Testing


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


## License

UIUC
