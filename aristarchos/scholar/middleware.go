package scholar

import (
	"context"
	"fmt"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/service"
	pbar "github.com/odysseia-greek/attike/aristophanes/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"strings"
	"time"
)

func AggregatorInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	var requestId string
	if ok {
		headerValue := md.Get(service.HeaderKey)
		if len(headerValue) > 0 {
			requestId = headerValue[0]
		}

		if requestId == "" {
			return handler(ctx, req)
		}

		splitID := strings.Split(requestId, "+")

		traceCall := false
		var traceID, spanID string

		if len(splitID) >= 3 {
			traceCall = splitID[2] == "1"
		}

		if len(splitID) >= 1 {
			traceID = splitID[0]
		}
		if len(splitID) >= 2 {
			spanID = splitID[1]
		}

		methodName := info.FullMethod

		host := ""
		if p, ok := peer.FromContext(ctx); ok {
			host = p.Addr.String()
		}

		if traceCall {
			traceReceived := &pbar.TraceRequest{
				TraceId:      traceID,
				ParentSpanId: spanID,
				Method:       methodName,
				Url:          methodName,
				Host:         host,
			}

			traceCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			trace, err := Tracer.Trace(traceCtx, traceReceived)
			if err != nil {
				logging.Error(err.Error())
				return handler(ctx, req)
			}
			logging.Trace(fmt.Sprintf("found traceId: %s", trace.CombinedId))
		}

		responseMd := metadata.New(map[string]string{service.HeaderKey: traceID})
		grpc.SendHeader(ctx, responseMd)
	}

	return handler(ctx, req)
}
