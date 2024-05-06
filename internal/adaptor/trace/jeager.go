package trace

import (
	"context"
	fmt "fmt"
	"giftcard/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"io"
	"log"
	"os"
)

type ITracer interface {
	SpanFromContext(ctx context.Context, name string, spanType string) (trace.Span, context.Context)
}

var (
	T *tracer
)

type tracer struct {
	oteTracer trace.Tracer
}

func InitGlobalTracer(lc fx.Lifecycle) {
	config := config.C()
	jaegerCloser, err := connect(config)
	if err != nil {
		log.Fatalf("Error initializing Jaeger: %v", err)
	}
	lc.Append(fx.Hook{
		OnStop: func(c context.Context) error {
			log.Printf("Jaeger and OpenTelemetry shutdown\n")
			jaegerCloser.Close()
			return nil
		},
	})
}

func connect(confs *config.Config) (io.Closer, error) {
	os.Setenv("OTEL_INSTRUMENTATION_CAPTURE_REQUEST_BODY", "true")
	log.Println(confs.Jaeger)
	exp, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(confs.Jaeger.HostPort)))
	if err != nil {
		log.Fatalf(err.Error())
	}

	resources := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(config.C().Service.Name),
		semconv.DeploymentEnvironmentKey.String("dev"))

	tp := sdk.NewTracerProvider(
		sdk.WithSampler(sdk.AlwaysSample()),
		sdk.WithResource(resources),
		sdk.WithBatcher(exp),
		//sdk.WithBatcher(sexp),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return nil, nil
}

func (e *tracer) SpanFromContext(ctx context.Context, name string, spanType string) (trace.Span, context.Context) {
	c, span := otel.GetTracerProvider().Tracer("giftcard").Start(ctx, fmt.Sprintf("%s.%s", name, spanType))
	return span, c
}

//func (e *tracer) SetAttribute(span trace.Span,atr ...attribute.KeyValue) (trace.Span, context.Context) {
//	span.SetAttributes(atr...)
//}
