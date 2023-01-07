package tracer

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type Tracer struct {
	provider *tracesdk.TracerProvider
}

func Build(jaegerURL string) (*Tracer, error) {
	exporter, err := jaeger.New(
		jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerURL)),
	)
	if err != nil {
		return nil, fmt.Errorf("create jaeger exporter: %w", err)
	}

	provider := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter, tracesdk.WithBatchTimeout(time.Second)),
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1))),
		tracesdk.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("nanomart"),
				semconv.DeploymentEnvironmentKey.String("local"),
			),
		),
	)

	otel.SetTracerProvider(provider)

	return &Tracer{provider: provider}, nil
}

func (t *Tracer) Shutdown() {
	_ = t.provider.Shutdown(context.Background())
}
