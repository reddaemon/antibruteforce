package postgres

import (
	"context"
	"database/sql"

	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Repository interface contain methods to work with storage
type Repository interface {
	AddToBlacklist(ctx context.Context, subnet string) error
	RemoveFromBlacklist(ctx context.Context, subnet string) error
	AddToWhitelist(ctx context.Context, subnet string) error
	RemoveFromWhitelist(ctx context.Context, subnet string) error
	FindIP(ctx context.Context, ip string) (string, error)
}

type PsqlRepository struct {
	*sqlx.DB
	logger *zap.Logger
}

func NewPsqlRepository(DB *sqlx.DB, logger *zap.Logger) *PsqlRepository {
	return &PsqlRepository{DB: DB, logger: logger}
}

func (p PsqlRepository) AddToBlacklist(ctx context.Context, subnet string) error {
	query := `INSERT INTO blacklist (subnet) VALUES ($1)`
	_, err := p.DB.ExecContext(ctx, query, subnet)

	return err
}

func (p PsqlRepository) RemoveFromBlacklist(ctx context.Context, subnet string) error {
	query := `DELETE FROM blacklist WHERE subnet = $1`
	_, err := p.DB.ExecContext(ctx, query, subnet)

	return err
}

func (p PsqlRepository) AddTowhitelist(ctx context.Context, subnet string) error {
	query := `INSERT INTO whitelist (subnet) VALUES ($1)`
	_, err := p.DB.ExecContext(ctx, query, subnet)

	return err
}

func (p PsqlRepository) RemoveFromWhitelist(ctx context.Context, subnet string) error {
	query := `DELETE FROM whitelist WHERE subnet = $1`
	_, err := p.DB.ExecContext(ctx, query, subnet)

	return err
}

func (p PsqlRepository) FindIP(ctx context.Context, ip string) (string, error) {
	query := `SELECT distinct $2 FROM blacklist where $1::inet <<= subnet
			  union (
					SELECT distinct $3 FROM whitelist where $1::inet <== subnet
			)
	`

	list := make([]string, 0, 2)

	err := p.DB.SelectContext(ctx, &list, query)

	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	switch len(list) {
	case 0:
		return "", nil
	case 1:
		return list[0], nil
	default:
		p.logger.Info(fmt.Sprintf("ip: %s in more than one list. lists: %v", ip, list))
		return "blacklist", nil

	}

}
