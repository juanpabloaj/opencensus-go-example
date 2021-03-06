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
	"github.com/kelseyhightower/envconfig"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

type envSpecification struct {
	JaegerHost string `default:"0.0.0.0"`
}

func main() {
	enableObservabilityAndExporters()

	client := &http.Client{Transport: &ochttp.Transport{}}
	i := uint64(0)

	for {
		i++
		log.Printf("Performing fetch #%d", i)
		ctx, span := trace.StartSpan(context.Background(), fmt.Sprintf("fetch-%d", i))
		doWork(ctx, client)
		span.End()

		<-time.After(5 * time.Second)
	}

}

func doWork(ctx context.Context, client *http.Client) {
	req, _ := http.NewRequest("GET", "http://server:8888/", nil)

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
	var envs envSpecification

	err := envconfig.Process("client", &envs)
	if err != nil {
		log.Fatalf("%v", err)
	}

	agentEndpointURI := fmt.Sprintf("%s:6831", envs.JaegerHost)
	collectorEndpointURI := fmt.Sprintf("http://%s:14268/api/traces", envs.JaegerHost)

	jaegerExporter, err := jaeger.NewExporter(jaeger.Options{
		AgentEndpoint:     agentEndpointURI,
		CollectorEndpoint: collectorEndpointURI,
		ServiceName:       "demo-client",
	})

	if err != nil {
		log.Printf("failed to create the jaeger exporter %v", err)
	}

	trace.RegisterExporter(jaegerExporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
}
