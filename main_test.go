package rfe

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/stretchr/testify/assert"
)

type DocTest struct {
	Spam string `firestore:"spam"`
	Eggs bool   `firestore:"eggs"`
}

func TestFirestoreEmulator(t *testing.T) {
	t.Run("emulator run properly", func(t *testing.T) {
		assert := assert.New(t)
		firestoreEmulator := FirestoreEmulator{}

		firestoreEmulator.Start()
		defer firestoreEmulator.Shutdown()

		// Test Client Connection
		ctx := context.Background()
		client, err := firestore.NewClient(ctx, "test")
		defer client.Close()

		assert.NoError(err)

		// Test Writing
		want := DocTest{
			Spam: "baz",
			Eggs: true,
		}

		docTest := client.Doc("Foo/Bar")
		wr, err := docTest.Create(ctx, want)

		assert.NoError(err)
		assert.NotEmpty(wr)

		// Test Reading
		docsnap, err := docTest.Get(ctx)

		assert.NoError(err)

		var got DocTest
		err = docsnap.DataTo(&got)

		assert.NoError(err)
		assert.True(reflect.DeepEqual(want, got), fmt.Sprintf("got: %+v\nwant: %+v", got, want))
	})
}
