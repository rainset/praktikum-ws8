package tracer

import (
	"net/http"

	otelcontrib "go.opentelemetry.io/contrib"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		provider := otel.GetTracerProvider()
		tracer := provider.Tracer(
			"",
			oteltrace.WithInstrumentationVersion(otelcontrib.SemVersion()),
		)
		propagators := otel.GetTextMapPropagator()
		propagators.Extract(r.Context(), propagation.HeaderCarrier(r.Header))
		routePath := r.URL.Path
		spanName := r.Method + " " + routePath
		ctx, span := tracer.Start(
			r.Context(), spanName,
			oteltrace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", r)...),
			oteltrace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(r)...),
			oteltrace.WithAttributes(semconv.HTTPRouteKey.String(routePath)),
			oteltrace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest("nanomart", routePath, r)...),
			oteltrace.WithSpanKind(oteltrace.SpanKindServer),
		)
		defer span.End()

		sniffer := &statusRecorder{w, http.StatusOK}

		r = r.WithContext(ctx)
		next(sniffer, r)

		span.SetAttributes(semconv.HTTPStatusCodeKey.Int(sniffer.statusCode))
		spanStatus, spanMessage := semconv.SpanStatusFromHTTPStatusCode(sniffer.statusCode)
		span.SetStatus(spanStatus, spanMessage)
	}
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (ss *statusRecorder) WriteHeader(code int) {
	ss.statusCode = code
	ss.ResponseWriter.WriteHeader(code)
}
