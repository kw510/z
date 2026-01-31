package srv

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5"
	"github.com/kw510/zee/pkg/db"
	apiv1 "github.com/kw510/zee/pkg/gen/zee/api/v1"
)

type ApiServer struct{}

func (srv *ApiServer) Check(ctx context.Context, req *connect.Request[apiv1.CheckRequest]) (*connect.Response[apiv1.CheckResponse], error) {
	t := db.Tuple{
		Parent: db.Set{
			Namespace: req.Msg.GetParent().GetNamespace(),
			Id:        req.Msg.GetParent().GetId(),
			Relation:  req.Msg.GetParent().GetRelation(),
		},
		Child: db.Set{
			Namespace: req.Msg.GetChild().GetNamespace(),
			Id:        req.Msg.GetChild().GetId(),
			Relation:  req.Msg.GetChild().GetRelation(),
		},
	}

	// Check if the tuple exists directly
	exists, err := t.Exists(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if exists {
		return &connect.Response[apiv1.CheckResponse]{
			Msg: &apiv1.CheckResponse{
				Allowed: true,
			},
		}, nil
	}

	// Check if the tuple exists in any of the parent's subsets
	ss, err := t.Parent.Subsets(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	for _, s := range ss {
		res, err := srv.Check(ctx, &connect.Request[apiv1.CheckRequest]{
			Msg: &apiv1.CheckRequest{
				Parent: &apiv1.Set{
					Namespace: s.Namespace,
					Id:        s.Id,
					Relation:  s.Relation,
				},
				Child: &apiv1.Set{
					Namespace: t.Child.Namespace,
					Id:        t.Child.Id,
					Relation:  t.Child.Relation,
				},
			},
		})
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		if res.Msg.Allowed {
			return &connect.Response[apiv1.CheckResponse]{
				Msg: &apiv1.CheckResponse{
					Allowed: true,
				},
			}, nil
		}
	}
	return &connect.Response[apiv1.CheckResponse]{
		Msg: &apiv1.CheckResponse{
			Allowed: false,
		},
	}, nil
}

func (srv *ApiServer) Write(ctx context.Context, req *connect.Request[apiv1.WriteRequest]) (*connect.Response[apiv1.WriteResponse], error) {
	err := db.Tx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		for _, tuple := range req.Msg.GetAdds() {
			t := db.Tuple{
				Parent: db.Set{
					Namespace: tuple.GetParent().GetNamespace(),
					Id:        tuple.GetParent().GetId(),
					Relation:  tuple.GetParent().GetRelation(),
				},
				Child: db.Set{
					Namespace: tuple.GetChild().GetNamespace(),
					Id:        tuple.GetChild().GetId(),
					Relation:  tuple.GetChild().GetRelation(),
				},
			}
			err := t.CreateTx(ctx, tx)
			if err != nil {
				return fmt.Errorf("failed to create tuple: %w", err)
			}
		}
		for _, tuple := range req.Msg.GetRemoves() {
			t := db.Tuple{
				Parent: db.Set{
					Namespace: tuple.GetParent().GetNamespace(),
					Id:        tuple.GetParent().GetId(),
					Relation:  tuple.GetParent().GetRelation(),
				},
				Child: db.Set{
					Namespace: tuple.GetChild().GetNamespace(),
					Id:        tuple.GetChild().GetId(),
					Relation:  tuple.GetChild().GetRelation(),
				},
			}
			err := t.DeleteTx(ctx, tx)
			if err != nil {
				return fmt.Errorf("failed to delete tuple: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return &connect.Response[apiv1.WriteResponse]{}, nil
}

func (srv *ApiServer) WriteNamespaceRelations(ctx context.Context, req *connect.Request[apiv1.WriteNamespaceRelationsRequest]) (*connect.Response[apiv1.WriteNamespaceRelationsResponse], error) {
	err := db.Tx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		for _, namespaceRelation := range req.Msg.GetAdds() {
			r := db.Relation{
				Namespace:  namespaceRelation.GetNamespace(),
				Relation:   namespaceRelation.GetRelation(),
				Permission: namespaceRelation.GetPermission(),
			}
			err := r.CreateTx(ctx, tx)
			if err != nil {
				return fmt.Errorf("failed to create relation: %w", err)
			}
		}
		for _, namespaceRelation := range req.Msg.GetRemoves() {
			r := db.Relation{
				Namespace:  namespaceRelation.GetNamespace(),
				Relation:   namespaceRelation.GetRelation(),
				Permission: namespaceRelation.GetPermission(),
			}
			err := r.DeleteTx(ctx, tx)
			if err != nil {
				return fmt.Errorf("failed to delete relation: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return &connect.Response[apiv1.WriteNamespaceRelationsResponse]{}, nil
}

func (srv *ApiServer) Namespaces(ctx context.Context, req *connect.Request[apiv1.NamespacesRequest]) (*connect.Response[apiv1.NamespacesResponse], error) {
	relations, err := db.Relation{}.List(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	namespaces := map[string]*apiv1.Namespace{}
	for _, r := range relations {
		if _, ok := namespaces[r.Namespace]; !ok {
			namespaces[r.Namespace] = &apiv1.Namespace{}
			namespaces[r.Namespace].Relations = map[string]*apiv1.Relation{}
		}
		if _, ok := namespaces[r.Namespace].Relations[r.Relation]; !ok {
			namespaces[r.Namespace].Relations[r.Relation] = &apiv1.Relation{}
		}
		namespaces[r.Namespace].Relations[r.Relation].Permissions = append(namespaces[r.Namespace].Relations[r.Relation].Permissions, r.Permission)
	}

	return &connect.Response[apiv1.NamespacesResponse]{
		Msg: &apiv1.NamespacesResponse{
			Namespaces: namespaces,
		},
	}, nil
}

func (srv *ApiServer) ParentRelations(ctx context.Context, req *connect.Request[apiv1.ParentRelationsRequest]) (*connect.Response[apiv1.ParentRelationsResponse], error) {
	relations, err := db.Relation{}.ListParentRelations(ctx, req.Msg.GetParentNamespace(), req.Msg.GetParentId(), db.Set{
		Namespace: req.Msg.GetChild().GetNamespace(),
		Id:        req.Msg.GetChild().GetId(),
		Relation:  req.Msg.GetChild().GetRelation(),
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	relationsProto := map[string]*apiv1.Relation{}
	for _, relation := range relations {
		if _, ok := relationsProto[relation.Relation]; !ok {
			relationsProto[relation.Relation] = &apiv1.Relation{}
		}
		relationsProto[relation.Relation].Permissions = append(relationsProto[relation.Relation].Permissions, relation.Permission)
	}
	return &connect.Response[apiv1.ParentRelationsResponse]{
		Msg: &apiv1.ParentRelationsResponse{
			Relations: relationsProto,
		},
	}, nil
}
