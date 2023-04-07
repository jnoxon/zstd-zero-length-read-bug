#### klauspost/compress zero-length read bug example

This is a minimal example of a bug in the klauspost/compress library. When decompressing from a reader that returns zero-length reads, the decompressor may return an "unexpected EOF" error.

To reproduce, run `go test -v` in this directory. The test will fail with an "unexpected EOF" error.

```
$ go test
--- FAIL: TestZeroLengthReads (0.00s)
    poc_test.go:64: unexpected EOF
FAIL
exit status 1
FAIL	poc	0.003s
```

To validate the fix:

```
$ go mod edit -replace=github.com/klauspost/compress=github.com/jnoxon/compress@readbyte-unexpected-eof
$ go mod tidy
$ go test
PASS
ok  	poc	0.003s
```
