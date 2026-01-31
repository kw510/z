package srv_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"connectrpc.com/connect"
	"github.com/kw510/z/pkg/db"
	apiv1 "github.com/kw510/z/pkg/gen/z/api/v1"
	"github.com/kw510/z/pkg/srv"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	if err := db.Init(ctx); err != nil {
		log.Fatalln(fmt.Errorf("failed to initialize database: %w", err))
	}
	defer db.Close()

	code := m.Run()

	os.Exit(code)
}

func TestApiServerCheckDirect(t *testing.T) {
	namespace := "co.acme." + xid.New().String()
	tuple := &apiv1.Tuple{
		Parent: &apiv1.Set{
			Namespace: namespace + ".Note",
			Id:        xid.New().String(),
			Relation:  "owner",
		},
		Child: &apiv1.Set{
			Namespace: namespace + ".User",
			Id:        xid.New().String(),
		},
	}

	req := connect.NewRequest(
		&apiv1.CheckRequest{
			Parent: tuple.Parent,
			Child:  tuple.Child,
		})
	t.Run("Denied", func(t *testing.T) {
		ctx := context.Background()
		srv := &srv.ApiServer{}
		res, err := srv.Check(ctx, req)
		require.NoError(t, err)
		require.False(t, res.Msg.Allowed)
	})

	t.Run("Allowed", func(t *testing.T) {
		ctx := context.Background()
		srv := &srv.ApiServer{}
		writeReq := connect.NewRequest(&apiv1.WriteRequest{
			Adds: []*apiv1.Tuple{
				tuple,
			},
		})

		_, err := srv.Write(ctx, writeReq)
		require.NoError(t, err)

		res, err := srv.Check(ctx, req)
		require.NoError(t, err)
		require.True(t, res.Msg.Allowed)
	})
}

func TestApiServerCheckParentSubset(t *testing.T) {
	namespace := "co.acme." + xid.New().String()
	groupMember := &apiv1.Tuple{
		Parent: &apiv1.Set{
			Namespace: namespace + ".Group",
			Id:        xid.New().String(),
			Relation:  "member",
		},
		Child: &apiv1.Set{
			Namespace: namespace + ".User",
			Id:        xid.New().String(),
		},
	}

	groupAccess := &apiv1.Tuple{
		Parent: &apiv1.Set{
			Namespace: namespace + ".Note",
			Id:        xid.New().String(),
			Relation:  "owner",
		},
		Child: groupMember.Parent,
	}

	req := connect.NewRequest(
		&apiv1.CheckRequest{
			Parent: groupAccess.Parent,
			Child:  groupMember.Child,
		})

	t.Run("Denined", func(t *testing.T) {
		ctx := context.Background()
		srv := &srv.ApiServer{}
		res, err := srv.Check(ctx, req)
		require.NoError(t, err)
		require.False(t, res.Msg.Allowed)
	})

	t.Run("Denined if group not linked to post", func(t *testing.T) {
		ctx := context.Background()
		srv := &srv.ApiServer{}
		_, err := srv.Write(ctx, connect.NewRequest(&apiv1.WriteRequest{
			Adds: []*apiv1.Tuple{
				groupMember,
			},
		}))
		require.NoError(t, err)

		res, err := srv.Check(ctx, req)
		require.NoError(t, err)
		require.False(t, res.Msg.Allowed)

		// Cleanup
		_, err = srv.Write(ctx, connect.NewRequest(&apiv1.WriteRequest{
			Removes: []*apiv1.Tuple{
				groupMember,
			},
		}))
		require.NoError(t, err)
	})

	t.Run("Denied if user not part of group", func(t *testing.T) {
		ctx := context.Background()
		srv := &srv.ApiServer{}
		_, err := srv.Write(ctx, connect.NewRequest(&apiv1.WriteRequest{
			Adds: []*apiv1.Tuple{
				groupAccess,
			},
		}))
		require.NoError(t, err)

		res, err := srv.Check(ctx, req)
		require.NoError(t, err)
		require.False(t, res.Msg.Allowed)

		// Cleanup
		_, err = srv.Write(ctx, connect.NewRequest(&apiv1.WriteRequest{
			Removes: []*apiv1.Tuple{
				groupAccess,
			},
		}))
		require.NoError(t, err)
	})

	t.Run("Allowed if user is part of group and group has access", func(t *testing.T) {
		ctx := context.Background()
		srv := &srv.ApiServer{}
		_, err := srv.Write(ctx, connect.NewRequest(&apiv1.WriteRequest{
			Adds: []*apiv1.Tuple{
				groupAccess,
				groupMember,
			},
		}))
		require.NoError(t, err)

		res, err := srv.Check(ctx, req)
		require.NoError(t, err)
		require.True(t, res.Msg.Allowed)
	})
}

func TestApiServerNamespaces(t *testing.T) {
	namespace := "co.acme." + xid.New().String()
	ctx := context.Background()
	srv := &srv.ApiServer{}

	srv.WriteNamespaceRelations(ctx, connect.NewRequest(&apiv1.WriteNamespaceRelationsRequest{
		Adds: []*apiv1.NamespaceRelation{
			{
				Namespace:  namespace + ".Note",
				Relation:   "owner",
				Permission: "read",
			},
			{
				Namespace:  namespace + ".Note",
				Relation:   "owner",
				Permission: "write",
			},
			{
				Namespace:  namespace + ".Note",
				Relation:   "owner",
				Permission: "delete",
			},
		},
	}))

	req := connect.NewRequest(&apiv1.NamespacesRequest{})
	res, err := srv.Namespaces(ctx, req)

	require.NoError(t, err)
	require.ElementsMatch(t, []string{"read", "write", "delete"}, res.Msg.Namespaces[namespace+".Note"].Relations["owner"].Permissions)

}

func TestApiServerParentRelations(t *testing.T) {
	ctx := context.Background()
	srv := &srv.ApiServer{}

	namespace := "co.acme." + xid.New().String()
	tuple := &apiv1.Tuple{
		Parent: &apiv1.Set{
			Namespace: namespace + ".Note",
			Id:        xid.New().String(),
			Relation:  "owner",
		},
		Child: &apiv1.Set{
			Namespace: namespace + ".User",
			Id:        xid.New().String(),
		},
	}

	_, err := srv.Write(ctx, connect.NewRequest(&apiv1.WriteRequest{
		Adds: []*apiv1.Tuple{
			tuple,
		},
	}))
	require.NoError(t, err)

	_, err = srv.WriteNamespaceRelations(ctx, connect.NewRequest(&apiv1.WriteNamespaceRelationsRequest{
		Adds: []*apiv1.NamespaceRelation{
			{
				Namespace:  namespace + ".Note",
				Relation:   "owner",
				Permission: "read",
			},
			{
				Namespace:  namespace + ".Note",
				Relation:   "owner",
				Permission: "write",
			},
			{
				Namespace:  namespace + ".Note",
				Relation:   "owner",
				Permission: "delete",
			},
		},
	}))
	require.NoError(t, err)

	req := connect.NewRequest(&apiv1.ParentRelationsRequest{
		ParentNamespace: tuple.Parent.Namespace,
		ParentId:        tuple.Parent.Id,
		Child:           tuple.Child,
	})
	res, err := srv.ParentRelations(ctx, req)
	require.NoError(t, err)
	require.ElementsMatch(t, []string{"read", "write", "delete"}, res.Msg.Relations["owner"].Permissions)
}
