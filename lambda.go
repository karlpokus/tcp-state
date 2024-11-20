package main

import (
	"context"
	"os"
	"time"

	"log"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/trace"
)

const LAMBDA_CRON_TIMER_SECS = 10

func main() {
	ctx := context.Background()
	// log pid
	log.Printf("pid: %d", os.Getpid())
	// tracing conf
	os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://localhost:4317")
	tp, err := newTraceProvider(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer tp.Shutdown(ctx)
	// invoke lambda on a cron
	//
	// simulating the AWS SDK lambda.Start
	//
	// force flush on each invocation
	// and maintain tcp connection
	invoke(func() (err error) {
		defer func() {
			err = tp.ForceFlush(ctx)
		}()
		_, span := tp.Tracer("test").Start(ctx, "test")
		defer span.End()
		return
	})
}

type invokeFunc func() error

// invoke runs the passed invokeFunc on a cron simulating lambda.Start.
//
// We bail on invocation errors.
func invoke(f invokeFunc) {
	for {
		if err := f(); err != nil {
			log.Println(err)
			return
		}
		<-time.After(LAMBDA_CRON_TIMER_SECS * time.Second)
	}
}

// newTraceProvider returns a confgured trace provider with a batched gRPC exporter.
func newTraceProvider(ctx context.Context) (*trace.TracerProvider, error) {
	exp, err := otlptracegrpc.New(ctx)
	if err != nil {
		return nil, err
	}
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
	)
	return tp, nil
}
