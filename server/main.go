package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
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

	originalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("get request ...")
		io.Copy(ioutil.Discard, r.Body)

		time.Sleep(time.Duration(rand.Intn(1000)+1) * time.Millisecond)

		w.Write([]byte("Hello, World!"))
	})

	opencensusHandler := &ochttp.Handler{
		Handler: originalHandler,
	}

	mux := http.NewServeMux()
	mux.Handle("/", opencensusHandler)
	log.Fatal(http.ListenAndServe(":8888", mux))

}

func enableObservabilityAndExporters() {
	var envs envSpecification

	err := envconfig.Process("server", &envs)
	if err != nil {
		log.Fatalf("%v", err)
	}

	agentEndpointURI := fmt.Sprintf("%s:6831", envs.JaegerHost)
	collectorEndpointURI := fmt.Sprintf("http://%s:14268/api/traces", envs.JaegerHost)

	jaegerExporter, err := jaeger.NewExporter(jaeger.Options{
		AgentEndpoint:     agentEndpointURI,
		CollectorEndpoint: collectorEndpointURI,
		ServiceName:       "demo-server",
	})

	if err != nil {
		log.Printf("failed to create the jaeger exporter %v", err)
	}

	trace.RegisterExporter(jaegerExporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
}
