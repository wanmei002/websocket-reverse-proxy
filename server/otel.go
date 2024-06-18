package server

import (
	"context"
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/appengine/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

var (
	tp         = sdktrace.NewTracerProvider()
	propagator = propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return otelgrpc.UnaryServerInterceptor(
		otelgrpc.WithTracerProvider(tp),
		otelgrpc.WithPropagators(propagator),
	)
}

func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return otelgrpc.UnaryClientInterceptor(
		otelgrpc.WithTracerProvider(tp),
		otelgrpc.WithPropagators(propagator),
	)
}

func UnaryServerInterceptorOtelp(appName string) grpc.UnaryServerInterceptor {

	ntp, _ := GetTracerProvider(appName, "127.0.0.1:4174")

	return otelgrpc.UnaryServerInterceptor(
		otelgrpc.WithTracerProvider(ntp),
		otelgrpc.WithPropagators(propagator),
	)
}

func GetTracerProvider(appName, address string) (*sdktrace.TracerProvider, error) {
	ntp := sdktrace.NewTracerProvider()
	ctx := context.Background()
	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceName(appName),
		),
	)
	if err != nil {
		log.Errorf(ctx, fmt.Sprintf("failed to create resource: %v", err))
		return ntp, err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	var traceExporter *otlptrace.Exporter
	if address != "" {
		conn, connErr := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if connErr != nil {
			log.Errorf(ctx, fmt.Sprintf("failed to create client: %v", connErr))
		} else {
			traceExporter, err = otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
		}
	} else {
		traceExporter, err = otlptracegrpc.New(ctx)
	}
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	ntp = sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(ntp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return ntp, nil
}
