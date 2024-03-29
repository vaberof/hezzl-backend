package pggood

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/vaberof/hezzl-backend/internal/domain/good"
	"github.com/vaberof/hezzl-backend/internal/infra/messagebroker/nats/publisher"
	"github.com/vaberof/hezzl-backend/internal/infra/storage"
	"github.com/vaberof/hezzl-backend/pkg/domain"
	"log"
)

type PgGoodStorage struct {
	db               *sqlx.DB
	goodLogPublisher publisher.Publisher
}

func NewPgGoodStorage(db *sqlx.DB, goodLogPublisher publisher.Publisher) *PgGoodStorage {
	return &PgGoodStorage{
		db:               db,
		goodLogPublisher: goodLogPublisher,
	}
}

func (gs *PgGoodStorage) Create(projectId domain.ProjectId, name domain.GoodName) (*good.Good, error) {
	var postgresGood Good
	query := `
			INSERT INTO goods(
			                  project_id,
			                  name
			) VALUES ($1, $2)
			RETURNING 
			    id, 
			    project_id,
			    name,
			    description,
			    priority,
			    removed,
			    created_at	    
	`
	row := gs.db.QueryRow(query, projectId, name)
	if err := row.Scan(
		&postgresGood.Id,
		&postgresGood.ProjectId,
		&postgresGood.Name,
		&postgresGood.Description,
		&postgresGood.Priority,
		&postgresGood.Removed,
		&postgresGood.CreatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to create good in database: %w", err)
	}

	if err := gs.goodLogPublisher.PublishGoodLog(
		postgresGood.Id,
		postgresGood.ProjectId,
		postgresGood.Name,
		postgresGood.Description.String,
		postgresGood.Priority,
		postgresGood.Removed,
		postgresGood.CreatedAt,
	); err != nil {
		log.Println("Failed to publish good log:", err)
	}

	return toDomainGood(&postgresGood), nil
}

func (gs *PgGoodStorage) Update(id domain.GoodId, projectId domain.ProjectId, name domain.GoodName, description *domain.GoodDescription) (*good.Good, error) {
	tx, err := gs.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction while updating good: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec("LOCK TABLE goods IN SHARE ROW EXCLUSIVE MODE")
	if err != nil {
		return nil, fmt.Errorf("failed to lock table while updating good: %w", err)
	}

	var postgresGood Good

	query := `
		UPDATE goods 
		SET name=$1, 
		    description=COALESCE($2, description)
		WHERE id=$3 AND project_id=$4
		RETURNING 
			    id, 
			    project_id,
			    name,
			    description,
			    priority,
			    removed,
			    created_at
	`

	row := tx.QueryRow(query, name, description, id, projectId)
	if err = row.Scan(
		&postgresGood.Id,
		&postgresGood.ProjectId,
		&postgresGood.Name,
		&postgresGood.Description,
		&postgresGood.Priority,
		&postgresGood.Removed,
		&postgresGood.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to update good in database: %w", storage.ErrPostgresGoodNotFound)
		}
		return nil, fmt.Errorf("failed to update good in database: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction while updating good: %w", err)
	}

	if err = gs.goodLogPublisher.PublishGoodLog(
		postgresGood.Id,
		postgresGood.ProjectId,
		postgresGood.Name,
		postgresGood.Description.String,
		postgresGood.Priority,
		postgresGood.Removed,
		postgresGood.CreatedAt,
	); err != nil {
		log.Println("Failed to publish good log:", err)
	}

	return toDomainGood(&postgresGood), nil
}

func (gs *PgGoodStorage) Delete(id domain.GoodId, projectId domain.ProjectId) (*good.Good, error) {
	tx, err := gs.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction while deleting good: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec("LOCK TABLE goods IN SHARE ROW EXCLUSIVE MODE")
	if err != nil {
		return nil, fmt.Errorf("failed to lock table while deleting good: %w", err)
	}

	var postgresGood Good

	query := `
		UPDATE goods SET removed=TRUE 
		             WHERE id=$1 AND project_id=$2
		RETURNING 
			    id, 
			    project_id,
			    name,
			    description,
			    priority,
			    removed,
			    created_at
	`

	row := tx.QueryRow(query, id, projectId)
	err = row.Scan(
		&postgresGood.Id,
		&postgresGood.ProjectId,
		&postgresGood.Name,
		&postgresGood.Description,
		&postgresGood.Priority,
		&postgresGood.Removed,
		&postgresGood.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to delete good: %w", storage.ErrPostgresGoodNotFound)
		}
		return nil, fmt.Errorf("failed to delete good: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction while deleting good: %w", err)
	}

	if err = gs.goodLogPublisher.PublishGoodLog(
		postgresGood.Id,
		postgresGood.ProjectId,
		postgresGood.Name,
		postgresGood.Description.String,
		postgresGood.Priority,
		postgresGood.Removed,
		postgresGood.CreatedAt,
	); err != nil {
		log.Println("Failed to publish good log:", err)
	}

	return toDomainGood(&postgresGood), nil
}

func (gs *PgGoodStorage) List(limit, offset int) ([]*good.Good, error) {
	limitOffsetParams := fmt.Sprintf(" LIMIT %d OFFSET %d ", limit, offset)

	query := `
			SELECT 
				id, 
				project_id,
				name,
				description,
				priority,
				removed,
				created_at
			FROM goods
			ORDER BY id
			` + limitOffsetParams

	rows, err := gs.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list goods: %w", err)
	}
	defer rows.Close()

	var postgresGoods []*Good

	for rows.Next() {
		var postgresGood Good

		err = rows.Scan(
			&postgresGood.Id,
			&postgresGood.ProjectId,
			&postgresGood.Name,
			&postgresGood.Description,
			&postgresGood.Priority,
			&postgresGood.Removed,
			&postgresGood.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan while listing goods: %w", err)
		}

		postgresGoods = append(postgresGoods, &postgresGood)
	}

	return toDomainGoods(postgresGoods), nil
}

func (gs *PgGoodStorage) ChangePriority(id domain.GoodId, projectId domain.ProjectId, newPriority domain.GoodPriority) ([]*good.Good, error) {
	tx, err := gs.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction while changing good priorities: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec("LOCK TABLE goods IN SHARE ROW EXCLUSIVE MODE")
	if err != nil {
		return nil, fmt.Errorf("failed to lock table while changing good priorities: %w", err)
	}

	query := `
		UPDATE goods SET priority=$1
		WHERE (id=$2 AND project_id=$3) OR id>$2
		RETURNING 
			    id, 
			    project_id,
			    name,
			    description,
			    priority,
			    removed,
			    created_at
	`

	rows, err := tx.Query(query, newPriority, id, projectId)
	if err != nil {
		return nil, fmt.Errorf("failed to change good priorities: %w", err)
	}
	defer rows.Close()

	var postgresGoods []*Good

	for rows.Next() {
		var postgresGood Good

		err = rows.Scan(
			&postgresGood.Id,
			&postgresGood.ProjectId,
			&postgresGood.Name,
			&postgresGood.Description,
			&postgresGood.Priority,
			&postgresGood.Removed,
			&postgresGood.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan changing good priorities: %w", err)
		}

		postgresGoods = append(postgresGoods, &postgresGood)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction while changing good priorities: %w", err)
	}

	go func() {
		for _, postgresGood := range postgresGoods {
			if err = gs.goodLogPublisher.PublishGoodLog(
				postgresGood.Id,
				postgresGood.ProjectId,
				postgresGood.Name,
				postgresGood.Description.String,
				postgresGood.Priority,
				postgresGood.Removed,
				postgresGood.CreatedAt,
			); err != nil {
				log.Println("Failed to publish good log:", err)
			}
		}
	}()

	return toDomainGoods(postgresGoods), nil
}

func (gs *PgGoodStorage) IsExists(id domain.GoodId, projectId domain.ProjectId) (bool, error) {
	query := `
			SELECT id FROM goods
			WHERE id=$1 AND project_id=$2
	`
	var goodId int64
	err := gs.db.QueryRow(query, id, projectId).Scan(&goodId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check whether the good exists or not: %w", err)
	}
	return true, nil
}
