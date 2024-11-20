# tcp-state
> Do tcp connections survive lambda invocations? How long can they be suspended without breaking?

Let's investigate.

TL;DR

TCP connections survive without issue being suspended for at least 2 minutes on my laptop. See limitations and the bottom.

# method
Run local simulation with otel collector, trace generator and tcp connection tracker. Simulate lambda invocations by suspending/resuming the trace generator.

````sh
# tcp connection tracker
$ sudo conntrack -E -p tcp --dport 4317
# otel collector
$ otelcol --config conf.yaml
# trace generator
$ go run lambda.go
# suspend/resume
$ kill -STOP <pid>
$ kill -CONT <pid>
````

# result
The tcp connection remains `ESTABLISHED` throughout the test. Traces are sent to the collector without issue as soon as the trace generator is resumed.

# limitations
- No proxies or load balancers involved
- No monitoring on tcp keepalive settings
- No validation in what type of resume/suspend operation that the AWS lambda runtime actually use. Would be interesting to replace the trace generator with [AWS firecracker](https://firecracker-microvm.github.io/)
