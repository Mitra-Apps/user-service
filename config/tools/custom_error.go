package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

type ErrorResponse struct {
	Code       string `json:"code"`
	CodeDetail string `json:"code_detail"`
	Message    string `json:"message"`
}

// TODO: need to move it to util service if possible
func CustomErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, writer http.ResponseWriter, request *http.Request, err error) {
	const fallback = `{"success": "false", "message": "failed to marshal error message"}`
	// var res api.BasicErrorResponse

	s := status.Convert(err)
	pb := s.Proto()
	st := runtime.HTTPStatusFromCode(s.Code())

	var errResponse ErrorResponse
	unmarshalErr := json.Unmarshal([]byte(pb.Message), &errResponse)
	if unmarshalErr != nil {
		fmt.Println("Error unmarshaling JSON:", unmarshalErr)
		errResponse.Code = s.Code().String()
		errResponse.CodeDetail = s.Code().String()
		errResponse.Message = pb.Message
	}

	pbApi := &ErrorResponse{
		Code:       errResponse.Code,
		CodeDetail: errResponse.CodeDetail,
		Message:    errResponse.Message,
	}

	writer.Header().Del("Trailer")
	writer.Header().Del("Transfer-Encoding")

	contentType := marshaler.ContentType(pbApi)
	writer.Header().Add("Content-Type", contentType)

	buf, merr := marshaler.Marshal(pbApi)
	if merr != nil {
		grpclog.Infof("Failed to marshal error message %q: %v", s, merr)
		writer.WriteHeader(http.StatusInternalServerError)
		if _, err := io.WriteString(writer, fallback); err != nil {
			grpclog.Infof("Failed to write response: %v", err)
		}
		return
	}

	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		grpclog.Infof("Failed to extract ServerMetadata from context")
	}

	for k, vs := range md.HeaderMD {
		{
			for _, v := range vs {
				writer.Header().Add(k, v)
			}
		}
	}

	// RFC 7230 https://tools.ietf.org/html/rfc7230#section-4.1.2
	// Unless the request includes a TE header field indicating "trailers"
	// is acceptable, as described in Section 4.3, a server SHOULD NOT
	// generate trailer fields that it believes are necessary for the user
	// agent to receive.
	te := request.Header.Get("TE")
	doForwardTrailers := strings.Contains(strings.ToLower(te), "trailers")

	if doForwardTrailers {
		writer.Header().Set("Transfer-Encoding", "chunked")
	}

	writer.WriteHeader(st)
	if _, err := writer.Write(buf); err != nil {
		grpclog.Infof("Failed to write response: %v", err)
	}
}
