package openfga_test

import (
	"testing"

	openfga "github.com/goropikari/gha_openfga"
	fga "github.com/openfga/go-sdk"
	"github.com/stretchr/testify/require"
)

func TestCheck(t *testing.T) {
	t.Run("example test", func(t *testing.T) {
		cl, err := openfga.NewClient()
		require.NoError(t, err)

		err = cl.CreateStore(t.Context())
		require.NoError(t, err)

		err = cl.WriteAuthorizationModel(t.Context())
		require.NoError(t, err)

		err = cl.WriteTuple(t.Context(), fga.TupleKey{
			User:     "user:alice",
			Relation: "reader",
			Object:   "document:doc1",
		})
		require.NoError(t, err)

		allow, err := cl.Check(t.Context(), fga.TupleKey{
			User:     "user:alice",
			Relation: "reader",
			Object:   "document:doc1",
		})
		require.NoError(t, err)
		require.True(t, allow)

		allow2, err := cl.Check(t.Context(), fga.TupleKey{
			User:     "user:bob",
			Relation: "reader",
			Object:   "document:doc1",
		})
		require.NoError(t, err)
		require.False(t, allow2)
	})
}
