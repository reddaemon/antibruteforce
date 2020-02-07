package service

import (
	"context"
	"errors"
	"net"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/reddaemon/antibruteforce/internal/database/postgres"

	"github.com/reddaemon/antibruteforce/config"
	"github.com/reddaemon/antibruteforce/internal/service/api/bucket"
	"go.uber.org/zap"
)

type Usage interface {
	Auth(ctx context.Context, login string, password string, ip string) error
	Drop(ctx context.Context, login string, ip string) error
	AddToBlacklist(ctx context.Context, subnet string) error
	RemoveFromBlacklist(ctx context.Context, subnet string) error
	AddToWhitelist(ctx context.Context, subnet string) error
	RemoveFromWhitelist(ctx context.Context, subnet string) error
}

type UsageStr struct {
	repo       postgres.Repository
	bucketRepo bucket.Storage
	logger     *zap.Logger
	config     *config.Config
}

func NewUsage(repo postgres.Repository, bucketRepo bucket.Bucket, logger *zap.Logger, config *config.Config) *UsageStr {
	return &UsageStr{repo: repo, bucketRepo: bucketRepo, logger: logger, config: config}

}

func (u *UsageStr) Auth(ctx context.Context, login string, password string, ip string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(a.config.ContextTimeout)*time.Millisecond)
	defer cancel()

	if login == "" {
		return errors.New("login is required")
	}

	if password == "" {
		return errors.New("password is required ")
	}

	if net.ParseIP(ip) == nil {
		return errors.New("incorrect IP")
	}

	iplist, err := u.repo.FindIP(ctx, ip)
	if err != nil {
		return err
	}
	switch iplist {
	case "whitelist":
		return nil
	case "blacklist":
		return errors.New("IP in blacklist")
	}

	var g errgroup.Group

	checks := []struct {
		key      string
		capacity uint
	}{
		{"lgn_" + login,
			u.config.LoginCapacity},
		{"pwd_" + password,
			u.config.PasswordCapacity},
		{"ip" + ip,
			u.config.IPCapacity},
	}
	for _, check := range checks {
		//check := check
		g.Go(func() error {
			return u.bucketRepo.Add(ctx, check.key, check.capacity, time.Duration(u.config.Rate)*time.Second)
		})
	}
	return g.Wait()
}

func (u *UsageStr) Drop(ctx context.Context, login string, ip string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(u.config.ContextTimeout)*time.Millisecond)
	defer cancel()

	if login == "" {
		return errors.New("login is required")
	}

	if net.ParseIP(ip) == nil {
		return errors.New("incorrect IP")
	}

	return u.bucketRepo.Drop(ctx, []string{"lgn_" + login, "ip" + ip})
}

func (u *UsageStr) AddToBlacklist(ctx context.Context, subnet string) error {
	сtx, cancel := context.WithTimeout(ctx, time.Duration(u.config.ContextTimeout)*time.Millisecond)
	defer cancel()

	_, _, err := net.ParseCIDR(subnet)
	if err != nil {
		return err
	}
	return u.repo.AddToBlacklist(ctx, subnet)
}

func (u *UsageStr) RemoveFromBlacklist(ctx context.Context, subnet string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(u.config.ContextTimeout)*time.Millisecond)
	defer cancel()

	_, _, err := net.ParseCIDR(subnet)
	if err != nil {
		return err
	}
	return u.repo.RemoveFromBlacklist(ctx, subnet)
}

func (u *UsageStr) AddToWhitelist(ctx context.Context, subnet string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(u.config.ContextTimeout)*time.Millisecond)
	defer cancel()

	_, _, err := net.ParseCIDR(subnet)
	if err != nil {
		return err
	}
	return u.repo.AddToWhitelist(ctx, subnet)
}

func (u *UsageStr) RemoveFromWhitelist(ctx context.Context, subnet string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(u.config.ContextTimeout)*time.Millisecond)
	defer cancel()
	_, _, err := net.ParseCIDR(subnet)
	if err != nil {
		return err
	}
	return u.repo.RemoveFromWhitelist(ctx, subnet)
}
