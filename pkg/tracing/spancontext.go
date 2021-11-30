package tracing

import (
	"encoding/json"

	"github.com/opentracing/opentracing-go"
)

func GetSpanContextAsJSON(tracer opentracing.Tracer, span opentracing.Span) (string, error) {
	// Debezium tracing guides that tracingspancontext uses 'java.util.Properties' of serialization format.
	// But 'java.util.Properties' is not univarsal across various languages.
	// So serialize span context with JSON.
	// https://debezium.io/documentation/reference/1.8/integrations/tracing.html

	// Inject context to carrier
	carrier := opentracing.TextMapCarrier{}
	tracer.Inject(span.Context(), opentracing.TextMap, carrier)

	spanContext, err := json.Marshal(carrier)
	if err != nil {
		return "", err
	}
	return string(spanContext), nil
}
