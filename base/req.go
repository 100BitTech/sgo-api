package base

import (
	"context"
	"fmt"
	"strings"

	"github.com/imroc/req/v3"
)

func NewReqClient(log Logger, traceContextKey any) *req.Client {
	return req.C().WrapRoundTripFunc(NewReqLogRoundTripFunc(log, traceContextKey))
}

func NewReqLogRoundTripFunc(log Logger, traceContextKey any) req.RoundTripWrapperFunc {
	logReq := func(req *req.Request) {
		log = log.WithTrace(req.Context(), traceContextKey).WithTag("Req")

		sb := strings.Builder{}
		sb.WriteString(fmt.Sprintf("REQ:%v %v\n", req.Method, req.URL))
		for k, v := range req.Headers {
			sb.WriteString(fmt.Sprintf("  %v: %v\n", k, v))
		}
		sb.WriteString(fmt.Sprintf("\n%v\n", string(req.Body)))
		log.Debug(sb.String())
	}

	logResp := func(resp *req.Response, err error) {
		log = log.WithTrace(resp.Request.Context(), traceContextKey).WithTag("Req")

		sb := strings.Builder{}
		sb.WriteString(fmt.Sprintf("RESP:%v %v\n", resp.Request.Method, resp.Request.URL))
		if err != nil {
			sb.WriteString(fmt.Sprintf("\n%+v", err))
			log.Error(sb.String())
		} else {
			sb.WriteString(fmt.Sprintf("%d %v\n", resp.StatusCode, resp.TotalTime()))
			for k, v := range resp.Header {
				sb.WriteString(fmt.Sprintf("  %v: %v\n", k, v))
			}
			sb.WriteString(fmt.Sprintf("\n%v", resp.String()))
			log.Debug(sb.String())
		}
	}

	return func(rt req.RoundTripper) req.RoundTripFunc {
		return func(req *req.Request) (resp *req.Response, err error) {
			if traceID := req.Context().Value(traceContextKey); traceID == nil {
				traceID = NewNanoID()
				req.SetContext(context.WithValue(req.Context(), traceContextKey, traceID))
			}

			logReq(req)
			resp, err = rt.RoundTrip(req)
			logResp(resp, err)
			return
		}
	}
}
