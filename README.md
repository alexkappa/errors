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
fmt.Println(err) // Example failed. write error [github.com/alexkappa/errors.Wrap(errors.go:76),github.com/alexkappa/errors.ExampleWrap(errors_test.go:56),testing.runExample(example.go:99),testing.RunExamples(example.go:36),testing.(*M).Run(testing.go:486),main.main(_testmain.go:58)]
```

You can also access the stack trace and print it out yourself.

```Go
err := errors.New("error with stack trace")
for _, frame := range err.Stack() {
    fmt.Printf("%s\n", frame.Func)
}
// github.com/alexkappa/errors.New
// main.main
```