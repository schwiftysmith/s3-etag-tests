# S3-etag-tests

## Quickstart

1. Run `docker-compose up`.
2. Execute tests in /test with `$ go test ./test -v`.

Current output with `SeaweedFS`:
```
$ go test ./test -v

=== RUN   TestEtagOfSmallFile
--- PASS: TestEtagOfSmallFile (0.23s)
=== RUN   TestEtagOfLargeFile
    s3_etag_test.go:80: 
                Error Trace:    s3_etag_test.go:80
                Error:          Not equal: 
                                expected: "630082f05fd91879b0380df1a9d78108-4"
                                actual  : "1374ca29ea88f7989442b66b589c59e1-13"
                            
                                Diff:
                                --- Expected
                                +++ Actual
                                @@ -1 +1 @@
                                -630082f05fd91879b0380df1a9d78108-4
                                +1374ca29ea88f7989442b66b589c59e1-13
                Test:           TestEtagOfLargeFile
--- FAIL: TestEtagOfLargeFile (0.64s)
FAIL
FAIL    s3-etag-tests/test      0.880s
FAIL
```