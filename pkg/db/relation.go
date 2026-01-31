package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Relation struct {
	Namespace  string
	Relation   string
	Permission string
}

func (r Relation) CreateTx(ctx context.Context, tx pgx.Tx) error {
	q := `
		INSERT INTO relations (namespace, relation, permission) VALUES ($1, $2, $3)
		ON CONFLICT DO NOTHING
	`
	_, err := tx.Exec(ctx, q, r.Namespace, r.Relation, r.Permission)
	if err != nil {
		return fmt.Errorf("failed to create relation: %w", err)
	}
	return nil
}

func (r Relation) DeleteTx(ctx context.Context, tx pgx.Tx) error {
	q := `
		DELETE FROM relations WHERE namespace = $1 AND relation = $2 AND permission = $3
	`
	_, err := tx.Exec(ctx, q, r.Namespace, r.Relation, r.Permission)
	if err != nil {
		return fmt.Errorf("failed to delete relation: %w", err)
	}
	return nil
}

func (r Relation) List(ctx context.Context) ([]Relation, error) {
	q := `
		SELECT namespace, relation, permission FROM relations
	`
	rows, err := pg.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("failed to list relations: %w", err)
	}
	defer rows.Close()

	relations := []Relation{}
	for rows.Next() {
		var relation Relation
		err := rows.Scan(&relation.Namespace, &relation.Relation, &relation.Permission)
		if err != nil {
			return nil, fmt.Errorf("failed to scan relation: %w", err)
		}
		relations = append(relations, relation)
	}
	return relations, nil
}

func (r Relation) ListParentRelations(ctx context.Context, parentNamespace string, parentId string, child Set) ([]Relation, error) {
	q := `
		SELECT 
			namespace, relation, permission
		FROM tuples, relations
		WHERE
			relations.relation = tuples.parent_relation AND
			relations.namespace = tuples.parent_namespace AND
			tuples.parent_namespace = $1 AND
			tuples.parent_id = $2 AND
			tuples.child_namespace = $3 AND
			tuples.child_id = $4 AND
			tuples.child_relation = $5
	`

	rows, err := pg.Query(ctx, q, parentNamespace, parentId, child.Namespace, child.Id, child.Relation)
	if err != nil {
		return nil, fmt.Errorf("failed to list relations: %w", err)
	}
	defer rows.Close()

	relations := []Relation{}
	for rows.Next() {
		var relation Relation
		err := rows.Scan(&relation.Namespace, &relation.Relation, &relation.Permission)
		if err != nil {
			return nil, fmt.Errorf("failed to scan relation: %w", err)
		}
		relations = append(relations, relation)
	}
	return relations, nil
}
