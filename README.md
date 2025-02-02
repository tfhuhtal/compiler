# compiler for University of Helsinki compilers course

## Installation

**Prequisites**: Go 1.23.5

**Clone the repository**

```bash
git clone git:@github.com:tfhuhtal/compiler.git
cd compiler
```

**Install dependencies**

```bash
go mod tidy
```

## Testing

Run all tests
```bash
go test ./src/...
```
Or run a specific test
```bash
go test -v ./src/parsern
```

## Running

Run the compiler

```bash
go run main.go compile <input-file> <output-file>
```

Run the compiler as server

```bash
go run main.go serve
```
Then you can send a POST request to `http://localhost:3000/` with the following body:
```json
{
    "code": "int main() { return 0; }"
}
```