package flare

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/grpcreflect"
	"github.com/tommyo/flare/proto"
	"github.com/tommyo/flare/proto/protoconnect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var _ protoconnect.SparkConnectServiceHandler = &Server{}

type Server struct {
	conf     *Config
	sessions *SessionStore
}

// TODO revisit
func (s *Server) Build(sessions *SessionStore) *http.Server {
	s.sessions = sessions

	mux := http.NewServeMux()
	mux.Handle(protoconnect.NewSparkConnectServiceHandler(s))

	reflector := grpcreflect.NewStaticReflector(protoconnect.SparkConnectServiceName)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	return &http.Server{
		Addr:    s.conf.String("addr"),
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}
}

// AddArtifacts implements protoconnect.SparkConnectServiceHandler.
// This should mimic the behavior of the SparkConnectAddArtifactsHandler from the Scala codebase.
func (s *Server) AddArtifacts(ctx context.Context, stream *connect.ClientStream[proto.AddArtifactsRequest]) (*connect.Response[proto.AddArtifactsResponse], error) {
	fmt.Printf("%v, %v\n", ctx, stream)

	// stagingDir is a unique temp jetstream.kv path for staging results
	// TODO do we need more than 1 per session?
	// stagingDir := fmt.Sprintf("staging.%s", stream.Co)
	res := &proto.AddArtifactsResponse{}
	return connect.NewResponse(res), nil
}

// AnalyzePlan implements protoconnect.SparkConnectServiceHandler.
func (s *Server) AnalyzePlan(ctx context.Context, req *connect.Request[proto.AnalyzePlanRequest]) (*connect.Response[proto.AnalyzePlanResponse], error) {
	fmt.Printf("%v, %v\n", ctx, req)
	res := &proto.AnalyzePlanResponse{}
	return connect.NewResponse(res), nil
}

// ArtifactStatus implements protoconnect.SparkConnectServiceHandler.
func (s *Server) ArtifactStatus(ctx context.Context, req *connect.Request[proto.ArtifactStatusesRequest]) (*connect.Response[proto.ArtifactStatusesResponse], error) {
	fmt.Printf("%v\n", req)
	res := &proto.ArtifactStatusesResponse{}
	return connect.NewResponse(res), nil
}

// Config implements protoconnect.SparkConnectServiceHandler.
func (s *Server) Config(ctx context.Context, req *connect.Request[proto.ConfigRequest]) (*connect.Response[proto.ConfigResponse], error) {
	fmt.Printf("%v, %v\n", ctx, req)
	res := &proto.ConfigResponse{}
	return connect.NewResponse(res), nil
}

// ExecutePlan implements protoconnect.SparkConnectServiceHandler.
func (s *Server) ExecutePlan(ctx context.Context, req *connect.Request[proto.ExecutePlanRequest], stream *connect.ServerStream[proto.ExecutePlanResponse]) error {
	fmt.Printf("%v, %v, %v\n", ctx, req, stream)
	return nil
}

// FetchErrorDetails implements protoconnect.SparkConnectServiceHandler.
func (s *Server) FetchErrorDetails(ctx context.Context, req *connect.Request[proto.FetchErrorDetailsRequest]) (*connect.Response[proto.FetchErrorDetailsResponse], error) {
	fmt.Printf("%v, %v\n", ctx, req)
	res := &proto.FetchErrorDetailsResponse{}
	return connect.NewResponse(res), nil
}

// Interrupt implements protoconnect.SparkConnectServiceHandler.
func (s *Server) Interrupt(ctx context.Context, req *connect.Request[proto.InterruptRequest]) (*connect.Response[proto.InterruptResponse], error) {
	fmt.Printf("%v, %v\n", ctx, req)
	res := &proto.InterruptResponse{}
	return connect.NewResponse(res), nil
}

// ReattachExecute implements protoconnect.SparkConnectServiceHandler.
func (s *Server) ReattachExecute(ctx context.Context, req *connect.Request[proto.ReattachExecuteRequest], stream *connect.ServerStream[proto.ExecutePlanResponse]) error {
	fmt.Printf("%v, %v, %v\n", ctx, req, stream)
	return nil
}

// ReleaseExecute implements protoconnect.SparkConnectServiceHandler.
func (s *Server) ReleaseExecute(ctx context.Context, req *connect.Request[proto.ReleaseExecuteRequest]) (*connect.Response[proto.ReleaseExecuteResponse], error) {
	fmt.Printf("%v, %v\n", ctx, req)
	res := &proto.ReleaseExecuteResponse{}
	return connect.NewResponse(res), nil
}

// ReleaseSession implements protoconnect.SparkConnectServiceHandler.
func (s *Server) ReleaseSession(ctx context.Context, req *connect.Request[proto.ReleaseSessionRequest]) (*connect.Response[proto.ReleaseSessionResponse], error) {
	fmt.Printf("%v, %v\n", ctx, req)

	serverSessionId := s.sessions.CloseSession(req.Msg.GetUserContext().UserId, req.Msg.SessionId)

	res := &proto.ReleaseSessionResponse{
		SessionId:           req.Msg.SessionId,
		ServerSideSessionId: serverSessionId,
	}
	return connect.NewResponse(res), nil
}

func NewServer(conf *Config) *Server {
	return &Server{conf: conf}
}
