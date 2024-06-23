package flare

import (
	"context"

	"connectrpc.com/connect"
	"github.com/tommyo/flare/proto"
	"github.com/tommyo/flare/proto/protoconnect"
)

var _ protoconnect.SparkConnectServiceHandler = &Server{}

type Server struct {
	conf *Config
}

// AddArtifacts implements protoconnect.SparkConnectServiceHandler.
func (s *Server) AddArtifacts(context.Context, *connect.ClientStream[proto.AddArtifactsRequest]) (*connect.Response[proto.AddArtifactsResponse], error) {
	panic("unimplemented")
}

// AnalyzePlan implements protoconnect.SparkConnectServiceHandler.
func (s *Server) AnalyzePlan(context.Context, *connect.Request[proto.AnalyzePlanRequest]) (*connect.Response[proto.AnalyzePlanResponse], error) {
	panic("unimplemented")
}

// ArtifactStatus implements protoconnect.SparkConnectServiceHandler.
func (s *Server) ArtifactStatus(context.Context, *connect.Request[proto.ArtifactStatusesRequest]) (*connect.Response[proto.ArtifactStatusesResponse], error) {
	panic("unimplemented")
}

// Config implements protoconnect.SparkConnectServiceHandler.
func (s *Server) Config(context.Context, *connect.Request[proto.ConfigRequest]) (*connect.Response[proto.ConfigResponse], error) {
	panic("unimplemented")
}

// ExecutePlan implements protoconnect.SparkConnectServiceHandler.
func (s *Server) ExecutePlan(context.Context, *connect.Request[proto.ExecutePlanRequest], *connect.ServerStream[proto.ExecutePlanResponse]) error {
	panic("unimplemented")
}

// FetchErrorDetails implements protoconnect.SparkConnectServiceHandler.
func (s *Server) FetchErrorDetails(context.Context, *connect.Request[proto.FetchErrorDetailsRequest]) (*connect.Response[proto.FetchErrorDetailsResponse], error) {
	panic("unimplemented")
}

// Interrupt implements protoconnect.SparkConnectServiceHandler.
func (s *Server) Interrupt(context.Context, *connect.Request[proto.InterruptRequest]) (*connect.Response[proto.InterruptResponse], error) {
	panic("unimplemented")
}

// ReattachExecute implements protoconnect.SparkConnectServiceHandler.
func (s *Server) ReattachExecute(context.Context, *connect.Request[proto.ReattachExecuteRequest], *connect.ServerStream[proto.ExecutePlanResponse]) error {
	panic("unimplemented")
}

// ReleaseExecute implements protoconnect.SparkConnectServiceHandler.
func (s *Server) ReleaseExecute(context.Context, *connect.Request[proto.ReleaseExecuteRequest]) (*connect.Response[proto.ReleaseExecuteResponse], error) {
	panic("unimplemented")
}

// ReleaseSession implements protoconnect.SparkConnectServiceHandler.
func (s *Server) ReleaseSession(context.Context, *connect.Request[proto.ReleaseSessionRequest]) (*connect.Response[proto.ReleaseSessionResponse], error) {
	panic("unimplemented")
}

func NewServer(conf *Config) *Server {
	return &Server{conf: conf}
}
