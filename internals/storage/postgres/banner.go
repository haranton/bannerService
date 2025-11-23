package postgres

import (
	"bannerService/internals/dto"
	"bannerService/internals/mapper"
	"bannerService/internals/models"
	"bannerService/internals/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

func (st *PostgresStorage) Banner(ctx context.Context, params dto.BannerQuery) (*models.Banner, error) {
	query := `
        SELECT b.id, b.content, b.is_active, b.created_at, b.updated_at
        FROM banners b
        JOIN feature_tag_banner ftb ON b.id = ftb.banner_id
        WHERE ftb.feature_id = $1 
          AND ftb.tag_id = $2 
          AND b.is_active = TRUE
        LIMIT 1
    `

	var banner models.Banner

	err := st.db.QueryRowContext(ctx, query, params.Feature_id, params.Tag_id).Scan(
		&banner.ID,
		&banner.Content,
		&banner.IsActive,
		&banner.CreatedAt,
		&banner.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrBannerNotFound
		}
		return nil, err
	}

	return &banner, nil

}

func (st *PostgresStorage) Banners(ctx context.Context, params dto.BannersQuery) ([]*models.Banner, error) {

	conditions := make([]string, 0, 3)
	args := make([]any, 0, 4)

	prefixFtbTable := "ftb."

	conditions = append(conditions, "b.is_active = TRUE")

	if params.Feature_id > 0 {
		conditions = append(conditions, fmt.Sprintf(prefixFtbTable+"feature_id = $%d", len(args)+1))
		args = append(args, params.Feature_id)
	}

	if params.Tag_id > 0 {
		conditions = append(conditions, fmt.Sprintf(prefixFtbTable+"tag_id = $%d", len(args)+1))
		args = append(args, params.Tag_id)
	}

	whereClause := strings.Join(conditions, " AND ")

	query :=
		"SELECT b.id, b.content, b.is_active, b.created_at, b.updated_at " +
			"FROM banners b " +
			"JOIN feature_tag_banner ftb ON b.id = ftb.banner_id " +
			"WHERE " + whereClause

	var limitPart string
	if params.Limit > 0 {
		args = append(args, params.Limit)
		limitPart = fmt.Sprintf(" LIMIT $%d", len(args)+1)
		query = query + limitPart

	}

	var Offset string
	if params.Offset > 0 {
		args = append(args, params.Offset)
		Offset = fmt.Sprintf(" OFFSET $%d", len(args)+1)
		query = query + Offset
	}

	rows, err := st.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	banners := []*models.Banner{}
	for rows.Next() {
		b := &models.Banner{}
		if err := rows.Scan(&b.ID, &b.Content, &b.IsActive, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		banners = append(banners, b)
	}

	return banners, nil

}

func (st *PostgresStorage) CreateBanner(
	ctx context.Context,
	banner *models.Banner,
	featureTags *dto.FeatureTags) (*models.Banner, error) {

	tx, err := st.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	queryCreateBanner := ` INSERT INTO banners (content, is_active, created_at, updated_at)
        VALUES (:content, :is_active, NOW(), NOW())
        RETURNING id`

	rows, err := tx.NamedQuery(queryCreateBanner, banner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		if rows.Err() != nil {
			return nil, fmt.Errorf("query failed before scanning banner id: %w", rows.Err())
		}
		return nil, fmt.Errorf("no rows returned for banner id")
	}

	if err := rows.Scan(&banner.ID); err != nil {
		return nil, fmt.Errorf("scan returned id failed: %w", err)
	}
	rows.Close()

	stmt, err := tx.PrepareNamed(`
    INSERT INTO feature_tag_banner (tag_id, feature_id, banner_id)
    VALUES (:tag_id, :feature_id, :banner_id)
`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	featureTagBanner := mapper.FeatureTagsBanner(featureTags, banner.ID)

	for _, rec := range featureTagBanner {

		_, err := stmt.Exec(rec)

		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return nil, fmt.Errorf("%w: constraint %s", storage.ErrDuplicateFeatureTag, pqErr.Constraint)
			}
		}
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return banner, nil

}

func (st *PostgresStorage) UpdateBanner(ctx context.Context,
	banner *models.Banner,
	featureTagBanner []*models.FeatureTagBanner) error {

	tx, err := st.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	queryDeleteOldLinks

	queryDeleteOldLinks := `DELETE FROM feature_tag_banner WHERE banner_id = $1`

	resDeleted, err := tx.Exec(queryDeleteOldLinks, banner.ID)
	if err != nil {
		return err
	}

	affected, err := resDeleted.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return storage.ErrBannerNotFound
	}

	query := `
	INSERT INTO feature_tag_banner (tag_id, feature_id, banner_id) 
	VALUES (:tag_id, :feature_id, :banner_id)
	`

	if _, err := tx.NamedExec(query, featureTagBanner); err != nil {
		return err
	}

	updateBannerQuery := `
        UPDATE banners
        SET content = :content,
            is_active = :is_active,
            updated_at = NOW()
        WHERE id = :id
    `

	_, err = tx.NamedExec(updateBannerQuery, banner)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (st *PostgresStorage) DeleteBanner(ctx context.Context, idBanner int) error {

	tx, err := st.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		DELETE from feature_tag_banner 
		WHERE banner_id = $1
	`
	if _, err := tx.Exec(query, idBanner); err != nil {
		return err
	}

	query = `
		DELETE from banners 
		WHERE id = $1
	`
	if _, err := tx.Exec(query, idBanner); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil

}
