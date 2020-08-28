# Constrained Application Protocol Client and Server for go

You can read more about CoAP in [RFC 7252][coap].  I also did
some preliminary work on `SUBSCRIBE` support from
[an early draft][shelby].

[shelby]: http://tools.ietf.org/html/draft-shelby-core-coap-01
[coap]: http://tools.ietf.org/html/rfc7252

# Benchmark
```
goos: darwin
goarch: amd64
BenchmarkQps
BenchmarkQps-8   	   96454	     12558 ns/op	     360 B/op	       9 allocs/op
PASS
```
