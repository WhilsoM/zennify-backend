package httpapi

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zennify/backend/internal/shared/grpcerr"
)

// writeErrorJSON maps an error to an HTTP response.
//
// Pass any error here — local validation or an upstream gRPC status.
// The upstream services always return gRPC status errors, so we read the code
// from them to pick the right HTTP status. Local errors that are not a gRPC
// status land in the default branch (502).
func writeErrorJSON(w http.ResponseWriter, err error) {
	st := status.Convert(err)
	switch st.Code() {
	case codes.InvalidArgument:
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": grpcerr.MsgInvalidRequest})
	case codes.AlreadyExists:
		writeJSON(w, http.StatusConflict, map[string]string{"error": st.Message()})
	case codes.Unauthenticated:
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": st.Message()})
	case codes.NotFound:
		writeJSON(w, http.StatusNotFound, map[string]string{"error": st.Message()})
	case codes.DeadlineExceeded, codes.Unavailable:
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": grpcerr.MsgUpstreamUnavailable})
	case codes.Internal:
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": grpcerr.MsgInternal})
	default:
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": grpcerr.MsgInternal})
	}
}
