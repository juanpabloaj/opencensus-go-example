package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"contrib.go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

func main() {

	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	enableObservabilityAndExporters()

	client := &http.Client{Transport: &ochttp.Transport{}}
	i := uint64(0)

	for {
		i += 1
		log.Printf("Performing fetch #%d", i)
		ctx, span := trace.StartSpan(context.Background(), fmt.Sprintf("fetch-%d", i))
		doWork(ctx, client)
		span.End()

		<-time.After(5 * time.Second)
	}

}

func doWork(ctx context.Context, client *http.Client) {
	req, _ := http.NewRequest("GET", "https://opencensus.io/", nil)

	req = req.WithContext(ctx)

	res, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to make the request: %v", err)
		return
	}

	io.Copy(ioutil.Discard, res.Body)
	_ = res.Body.Close()
}

func enableObservabilityAndExporters() {
	agentEndpointURI := "localhost:6831"
	collectorEndpointURI := "http://localhost:14268/api/traces"

	jaegerExporter, err := jaeger.NewExporter(jaeger.Options{
		AgentEndpoint:     agentEndpointURI,
		CollectorEndpoint: collectorEndpointURI,
		ServiceName:       "demo",
	})

	if err != nil {
		log.Printf("failed to create the jaeger exporter %v", err)
	}

	trace.RegisterExporter(jaegerExporter)
}
