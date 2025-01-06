# GoCache

GoCache is a distributed caching system implemented in Go. It supports consistent hashing, LRU cache eviction, and single-flight suppression to prevent cache stampede.

## Features

- **Consistent Hashing**: Implemented in [consistenthash/consistenthash.go](consistenthash/consistenthash.go).
- **LRU Cache**: Implemented in [lru/lru.go](lru/lru.go).
- **Single-Flight Suppression**: Implemented in [singleflight/singleflight.go](singleflight/singleflight.go).
- **gRPC Support**: Protobuf definitions and generated code in [gocachepb/](gocachepb/).

## Getting Started

### Prerequisites

- Go 1.22.4 or later

### Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/yourusername/gocache.git
    cd gocache/cmd
    ```

2. Build the project:

    ```sh
    go build -o server
    ```

### Running the Project

To start the cache servers and API server, run the provided script:

```sh
./run.sh
```

This will start three cache servers on ports 8001, 8002, and 8003, and an API server on port 9999.

### Testing

To run the tests, use the following command:

```sh
go test ./...
```

### Usage

API
The API server provides an endpoint to get cached values:

```sh
curl "http://localhost:9999/api?key=Tom"
```

### Consistent Hashing

The consistent hashing implementation can be found in consistenthash.go. It allows for adding and retrieving keys in a distributed manner.

### LRU Cache

The LRU cache implementation can be found in lru.go. It supports adding, retrieving, and evicting cache entries based on the least recently used policy.

### Single-Flight Suppression

The single-flight suppression implementation can be found in singleflight.go. It ensures that only one execution is in-flight for a given key at a time.

Contributing
Contributions are welcome! Please open an issue or submit a pull request.

License
This project is licensed under the MIT License

Feel free to customize the `README.md` as needed for your project.

