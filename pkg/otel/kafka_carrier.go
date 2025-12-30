package otel

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func ExtractTraceContext(ctx context.Context) ([]byte, error) {
	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, &carrier)

	traceContext, err := json.Marshal(carrier)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal(carrier): %w", err)
	}

	return traceContext, nil
}

func CtxFromKafkaHeaders(ctx context.Context, headers []kafka.Header) context.Context {
	carrier := propagation.MapCarrier{}

	for _, h := range headers {
		carrier[h.Key] = string(h.Value)
	}

	return otel.GetTextMapPropagator().Extract(ctx, carrier)
}
