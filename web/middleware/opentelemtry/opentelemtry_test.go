package opentelemetry

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"goStudy/web"
	"testing"
)

func initJeager(t *testing.T) {
	url := "http://localhost:14268/api/traces"
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		t.Fatal(err)
	}
	tp := sdktrace.NewTracerProvider(
		// Always be sure to batch in production.
		sdktrace.WithBatcher(exp),
		// Record information about this application in a Resource.
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("opentelemetry-demo"),
			attribute.String("environment", "dev"),
			attribute.Int64("ID", 1),
		)),
	)

	otel.SetTracerProvider(tp)
}

func TestMiddlewareBuilder_Build(t *testing.T) {
	tracer := otel.GetTracerProvider().Tracer("")
	initJeager(t)
	s := web.NewHttpServer()
	s.GET("/", func(ctx *web.Context) {
		ctx.Resp.Write([]byte("hello, world"))
	})
	s.GET("/user", func(ctx *web.Context) {
		c, span := tracer.Start(ctx.Req.Context(), "first_layer")
		defer span.End()

		c, second := tracer.Start(c, "second_layer")

		c, third1 := tracer.Start(c, "third_layer_1")

		third1.End()
		c, third2 := tracer.Start(c, "third_layer_1")
		
		third2.End()
		second.End()
		ctx.RespStatusCode = 200
		ctx.RespData = []byte("hello, world")
	})

	s.Use((&MiddlewareBuilder{Tracer: tracer}).Build())
	s.Start(":8081")
}
