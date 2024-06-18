package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"net/http/httputil"
)

func GinRun() {
	ginEngine := gin.New()
	UnaryServerInterceptorOtelp("gin-service")
	ginEngine.Use(otelgin.Middleware("gin-service"))
	ginEngine.GET("/test", func(c *gin.Context) {
		otel.GetTextMapPropagator().Inject(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))
		//http.Get("http://127.0.0.1:21114/foo/api/v1/get_address")
		ctx := otel.GetTextMapPropagator().Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))
		fmt.Println("gin header: ", c.Request.Header)
		fmt.Println("gin trace: ", trace.SpanFromContext(ctx).SpanContext().TraceID().String())
		c.Request = c.Request.WithContext(ctx)
		u := &httputil.ReverseProxy{
			Director: func(req *http.Request) {
				req.URL.Scheme = "http"
				req.URL.Host = "127.0.0.1:21114"
				req.URL.Path = "/foo/api/v1/get_address"
			},
		}

		u.ServeHTTP(c.Writer, c.Request)
	})
	go func() {
		fmt.Println("gin run")
		http.ListenAndServe(":21115", ginEngine)
	}()

}
