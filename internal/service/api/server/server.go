package server

import (
	"context"

	"github.com/reddaemon/antibruteforce/internal/service/api/usage"
	api "github.com/reddaemon/antibruteforce/protofiles"
	"go.uber.org/zap"
)

type Server struct {
	usage  usage.Usage
	logger *zap.Logger
}

func NewServer(usage usage.Usage, logger *zap.Logger) *Server {
	return &Server{usage: usage, logger: logger}
}

// Auth method authorising request
func (s *Server) Auth(ctx context.Context, req *api.AuthRequest) (*api.AuthResponse, error) {
	err := s.usage.Auth(ctx, req.Login, req.Password, req.Ip)

	return &api.AuthResponse{Ok: err == nil}, err
}

// Drop method dropping request
func (s *Server) Drop(ctx context.Context, req *api.DropRequest) (*api.DropResponse, error) {
	err := s.usage.Drop(ctx, req.Login, req.Ip)

	return &api.DropResponse{Ok: err == nil}, err
}

// AddToBlacklist method adding suspicious subnet to blacklist
func (s *Server) AddToBlacklist(ctx context.Context,
	req *api.AddToBlacklistRequest) (*api.AddToBlacklistResponse, error) {
	err := s.usage.AddToBlacklist(ctx, req.Subnet)
	return &api.AddToBlacklistResponse{
		Ok: err == nil}, err
}

// RemoveFromBlacklist method removing subnet from blacklist
func (s *Server) RemoveFromBlacklist(ctx context.Context,
	req *api.RemoveFromBlacklistRequest) (*api.RemoveFromBlacklistResponse, error) {
	err := s.usage.RemoveFromBlacklist(ctx, req.Subnet)
	return &api.RemoveFromBlacklistResponse{Ok: err == nil}, err
}

// AddToWhitelist method adding subnet to whitelist
func (s *Server) AddToWhitelist(ctx context.Context,
	req *api.AddToWhitelistRequest) (*api.AddToWhitelistResponse, error) {
	err := s.usage.AddToWhitelist(ctx, req.Subnet)
	return &api.AddToWhitelistResponse{Ok: err == nil}, err
}

// RemoveFromWhitelist method removing subnet from whitelist
func (s *Server) RemoveFromWhitelist(ctx context.Context,
	req *api.RemoveFromWhitelistRequest) (*api.RemoveFromWhitelistResponse, error) {
	err := s.usage.RemoveFromWhitelist(ctx, req.Subnet)
	return &api.RemoveFromWhitelistResponse{Ok: err == nil}, err
}
