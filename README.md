# errors [![GoDoc](https://godoc.org/github.com/alexkappa/errors?status.svg)](http://godoc.org/github.com/alexkappa/errors)

Similar to standard library errors but with some stack trace goodness.

## Usage

```Go
err := errors.New("whoops!") // looks familiar?
```

You can even wrap an error to provide context.

```Go
_, err := w.Write(p)
if err != nil {
    err = errors.Wrap(err, "Example failed")
}
fmt.Println(err) // Output: Example failed. write error
```

It also provides a stack trace.

```Go
err := errors.New("error with stack trace")
for _, frame := range err.Stack() {
    fmt.Printf("%s\n", frame.Func)
}
// Output:
// github.com/alexkappa/errors.New
// github.com/alexkappa/errors.ExampleStack
// testing.runExample
// testing.RunExamples
// testing.(*M).Run
// main.main
```