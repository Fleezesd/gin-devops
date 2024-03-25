package view

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"
)

func ping(c *gin.Context) {
	ctx := c.Request.Context()

	otelgin.HTML(c, http.StatusOK, "ping", gin.H{
		"message":  "pong",
		"traceURL": uptrace.TraceURL(trace.SpanFromContext(ctx)),
	})
}
