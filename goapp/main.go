package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/codeedu/otel-go/infra/opentel"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func main() {
	ot := opentel.NewOpenTel()
	ot.ServiceName = "GoApp"
	ot.ServiceVersion = "0.1"
	ot.ExporterEndpoint = "http://localhost:9411/api/v2/spans"
	tracer = ot.GetTracer()

	router := mux.NewRouter()
	router.Use(otelmux.Middleware(ot.ServiceName))
	router.HandleFunc("/", homeHandler)
	http.ListenAndServe(":8888", router)
}

func homeHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := baggage.ContextWithoutBaggage(request.Context())

	// rotina 1 - Process File
	ctx, processFile := tracer.Start(ctx, "process-file")
	time.Sleep(time.Millisecond * 100)
	processFile.End()

	// rotina 2 - Fazer Request HTTP
	ctx, httpCall := tracer.Start(ctx, "request-remote-json")
	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:3000/", nil)
	if err != nil {
		log.Fatal(err)
	}
	res, err := client.Do(req) // chamo a requisição

	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	time.Sleep(time.Millisecond * 300)
	httpCall.End()

	// rotina 3 - Exibir resultado
	ctx, renderContent := tracer.Start(ctx, "render-content")
	time.Sleep(time.Millisecond * 100)
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(body))
	renderContent.End()
}
