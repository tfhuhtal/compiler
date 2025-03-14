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
go test ./...
```
Or run a specific test
```bash
go test -v ./parser
```

## Running

Run the compiler

```bash
go run main.go compile --input=<input> --output=<output-file>
```

Run the compiler as server

```bash
go run main.go serve --host=0.0.0.0
```
Then you can send a POST request to `http://localhost:3000/` with the following body:
```json
{
    "code": "var n: Int = read_int();
          print_int(n);
          while n > 1 do {
            if n % 2 == 0 then {
              n = n / 2;
            } else {
              n = 3*n + 1;
            }
            print_int(n);
          }"
}
```
