package pggood

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/vaberof/hezzl-backend/internal/domain/good"
	"github.com/vaberof/hezzl-backend/internal/infra/storage"
	"github.com/vaberof/hezzl-backend/pkg/domain"
)

type PgGoodStorage struct {
	db *sqlx.DB
}

func NewPgGoodStorage(db *sqlx.DB) *PgGoodStorage {
	return &PgGoodStorage{
		db: db,
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
		return nil, err
	}
	return toDomainGood(&postgresGood), nil
}

func (gs *PgGoodStorage) Update(id domain.GoodId, projectId domain.ProjectId, name domain.GoodName, description *domain.GoodDescription) (*good.Good, error) {
	tx, err := gs.db.Begin()
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec("LOCK TABLE goods IN SHARE ROW EXCLUSIVE MODE")
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return nil, err
		}
	}

	var postgresGood Good

	query := `
		UPDATE goods 
		SET name=$1, 
		    description=CASE 
		        			WHEN $2 IS NOT NULL THEN $2
							ELSE description
						END
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
		if err = tx.Rollback(); err != nil {
			return nil, err
		}
	}

	return toDomainGood(&postgresGood), nil
}

func (gs *PgGoodStorage) Delete(id domain.GoodId, projectId domain.ProjectId) error {
	query := `
		DELETE FROM goods WHERE id=$1 AND project_id=$2 
	`
	result, err := gs.db.Exec(query, id, projectId)
	if err != nil {
		return err
	}
	if count, _ := result.RowsAffected(); count == 0 {
		return storage.ErrPostgresGoodNotFound
	}
	return nil

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
		ORDER BY created_at DESC
		` + limitOffsetParams

	rows, err := gs.db.Query(query)
	if err != nil {
		return nil, err
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
			return nil, err
		}

		postgresGoods = append(postgresGoods, &postgresGood)
	}

	return toDomainGoods(postgresGoods), nil
}

func (gs *PgGoodStorage) ChangePriority(id domain.GoodId, projectId domain.ProjectId, newPriority int) ([]*good.Good, error) {
	tx, err := gs.db.Begin()
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec("LOCK TABLE goods IN SHARE ROW EXCLUSIVE MODE")
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return nil, err
		}
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
		if err = tx.Rollback(); err != nil {
			return nil, err
		}
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
			if err = tx.Rollback(); err != nil {
				return nil, err
			}
		}

		postgresGoods = append(postgresGoods, &postgresGood)
	}

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
		return false, err
	}
	return true, nil
}
