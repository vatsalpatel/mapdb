# MapDB

## Description

MapDB is a standalone in-memory key-value store application, similar to Redis and KeyDB, but implemented in Go. It provides a fast, lightweight, and efficient in-memory database solution. MapDB is designed to be compatible with Redis libraries, offering high performance and simplicity in its usage and integration. It aims to serve as a powerful tool for developers needing a robust in-memory database for their applications.
## Installation

MapDB can be run locally or using Docker. 

### Local Installation

1. **Install Go**: If you don't have Go installed, you can download it from the [official Go website](https://golang.org/dl/). Follow the instructions provided for your specific operating system.

1. **Clone the repository**: Clone the MapDB repository to your local machine using the following command in your terminal:
```bash
git clone https://github.com/vatsalpatel/mapdb
```

3. **Build the application**
```bash
go build -o mapdb .
```

4. **Run the application**
```bash
./mapdb
```

### Docker Installation

1. **Build the Docker Image**
```bash
docker build -t mapdb .
```

2. **Run MapDB in a Docker container**
```bash
docker run -p 6379:6379 mapdb
```
## Usage

MapDB can be run either directly from the command line or inside a Docker container. It accepts several command-line flags for configuration:

- `-port <port>`: Set the port number on which MapDB will listen for connections. Default is 6379.
- `-server-type <type>`: Set the server type to run. Default is 0. The server types are:
    - 0: TCP Async Server
    - 1: TCP Single Threaded Server
    - 2: TCP Sync Server

### Local Usage

After building the application as described in the Installation section, you can run MapDB with the following command:

```bash
./main -port 6380 -server-type 1
```
### Docker Usage

```bash
docker run -p 6380:6380 mapdb -e PORT=6380 -e SERVER_TYPE=0
```

## Testing

MapDB uses the built-in testing framework of Go. You can run the tests with the `go test` command.

### Running Tests Locally

To run the tests locally, navigate to the project directory and run the following command:

```bash
go test ./...
```

## Contributing

Contributions to MapDB are welcome and appreciated. Here are some ways you can contribute:

- **Bug Reports**: If you find a bug, please create an issue in the GitHub issue tracker describing the problem and including steps to reproduce the issue.
- **Feature Requests**: If you have an idea for a new feature, feel free to create an issue in the GitHub issue tracker describing your idea.
- **Pull Requests**: If you've fixed a bug or implemented a new feature, you can submit a pull request. Please make sure your code follows the existing style and all tests pass. Please follow [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/#summary)

### Development Workflow

Here's how you can set up a development environment for MapDB:

1. Fork the repository on GitHub.
2. Clone your fork to your local machine: `git clone https://github.com/<your-username>/mapdb.git`
3. Create a new branch for your changes: `git checkout -b my-feature-branch`
4. Make your changes and commit them with a descriptive commit message.
5. Push your changes to your fork on GitHub: `git push origin my-feature-branch`
6. Create a pull request from your fork to the main MapDB repository.

Before submitting a pull request, please make sure your code builds successfully and all tests pass.

## License

MapDB is licensed under the [MIT License](LICENSE). The MIT License allows you to freely use, copy, and modify MapDB, as long as you provide attribution and don’t hold us liable. See the LICENSE file for more details.

## Future Enhancements

Here are some features and improvements that are planned for future versions of MapDB:

- **Pipelined Commands**: Support for executing multiple commands in a pipeline to improve performance.
- **Transactions**: Support for transactions to ensure data consistency and integrity.
- **Monitor Command**: Support for streaming the commands ran by clients.
- **Append Only Log**: Support for logging the commands ran by clients.
- **List Support**: Support for storing and manipulating lists of items.

Please note that these are planned features and improvements, and their implementation may change or be delayed. If you have a feature you'd like to see in MapDB, please open an issue in the GitHub issue tracker.
Feel free to implement any of the future enhancements above and create a PR. See [Contributing](#contributing) for more details.