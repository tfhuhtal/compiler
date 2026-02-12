# compiler for University of Helsinki compilers course

The full language spec can be found [here](https://hy-compilers.github.io/spring-2026/language-spec/).

For example Fibonacci series looks like this:
```bash

fun fibonacci(x: Int): Int {
  if x == 0 or x == 1 then {
    return x;
  } else {
    return fibonacci(x - 1) + fibonacci(x - 2);
  }
}

var i: Int = 0;

while i <= 15 do {
  print_int(fibonacci(i));
  i = i + 1;
}
```

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

Run the compiler:

```bash
go run main.go compile --input=<input> --output=<output-file>
```

Run the compiler as server

```bash
go run main.go serve --host=0.0.0.0
```
Then you can send request to the server, for example:
```
echo '{"command":"compile","code":"var x: Int = 10; var y: Int = 20; print_int(x + y);"}' | nc -w 2 -q 1 127.0.0.1 3000
```

Run the interpreter

```bash
go run main.go interpret --input=<input>
```

for example

```bash
go run main.go interpret --input="var a: Int = 0; var b: Int = 1; var next: Int = b; var count: Int = 1; while count <= 50 do { print_int(next); count = count + 1; a = b; b = next; next = a + b;}"
```
