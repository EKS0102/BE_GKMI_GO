package repository

import (
    "context"
    "time"

    "BE_GKMI_NTC_GO/internal/model"
    "github.com/jackc/pgx/v5"
 )

type RefreshTokenRepository struct {
    DB *pgx.Conn
}

func NewRefreshTokenRepository(db *pgx.Conn) *RefreshTokenRepository {
    return &RefreshTokenRepository{DB: db}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, refreshToken model.RefreshToken) (*model.RefreshToken, error) {
    query := `
        INSERT INTO refresh_tokens (user_id, token_hash, expires_at, revoked_at, created_at, replaced_by_token_hash)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, user_id, token_hash, expires_at, revoked_at, created_at, replaced_by_token_hash
    `

    var result model.RefreshToken
    err := r.DB.QueryRow(ctx, query, refreshToken.UserID, refreshToken.TokenHash, refreshToken.ExpiresAt, refreshToken.RevokedAt, refreshToken.CreatedAt, refreshToken.ReplacedByTokenHash).Scan(&result.ID, &result.UserID, &result.TokenHash, &result.ExpiresAt, &result.RevokedAt, &result.CreatedAt, &result.ReplacedByTokenHash)
    if err != nil {
        return nil, err
    }
    return &result, nil
}

func (r *RefreshTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error) {
    query := `
        SELECT id, user_id, token_hash, expires_at, revoked_at, created_at, replaced_by_token_hash
        FROM refresh_tokens
        WHERE token_hash = $1
    `

    var result model.RefreshToken
    err := r.DB.QueryRow(ctx, query, tokenHash).Scan(&result.ID, &result.UserID, &result.TokenHash, &result.ExpiresAt, &result.RevokedAt, &result.CreatedAt, &result.ReplacedByTokenHash)
    if err != nil {
        return nil, err
    }
    return &result, nil
}

func (r *RefreshTokenRepository) Revoke(ctx context.Context, tokenHash string, revokedAt time.Time) error {
    _, err := r.DB.Exec(ctx, `UPDATE refresh_tokens SET revoked_at = $1 WHERE token_hash = $2`, revokedAt, tokenHash)
    return err
}

func (r *RefreshTokenRepository) SetReplacedBy(ctx context.Context, tokenHash string, replacedByTokenHash string) error {
    _, err := r.DB.Exec(ctx, `UPDATE refresh_tokens SET replaced_by_token_hash = $1 WHERE token_hash = $2`, replacedByTokenHash, tokenHash)
    return err
}

func (r *RefreshTokenRepository) Rotate(ctx context.Context, oldTokenHash string, newRefreshToken model.RefreshToken, revokedAt time.Time) (*model.RefreshToken, error) {
    tx, err := r.DB.Begin(ctx)
    if err != nil {
        return nil, err
    }
    defer tx.Rollback(ctx)

    insertQuery := `
        INSERT INTO refresh_tokens (user_id, token_hash, expires_at, revoked_at, created_at, replaced_by_token_hash)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, user_id, token_hash, expires_at, revoked_at, created_at, replaced_by_token_hash
    `

    var result model.RefreshToken
    err = tx.QueryRow(ctx, insertQuery, newRefreshToken.UserID, newRefreshToken.TokenHash, newRefreshToken.ExpiresAt, newRefreshToken.RevokedAt, newRefreshToken.CreatedAt, newRefreshToken.ReplacedByTokenHash).Scan(&result.ID, &result.UserID, &result.TokenHash, &result.ExpiresAt, &result.RevokedAt, &result.CreatedAt, &result.ReplacedByTokenHash)
    if err != nil {
        return nil, err
    }

    _, err = tx.Exec(ctx, `UPDATE refresh_tokens SET revoked_at = $1, replaced_by_token_hash = $2 WHERE token_hash = $3`, revokedAt, result.TokenHash, oldTokenHash)
    if err != nil {
        return nil, err
    }

    if err := tx.Commit(ctx); err != nil {
        return nil, err
    }

    return &result, nil
}