// before running this test, make sure to start a local minio server
//
//	with the following command inside data directory:
//
// docker compose up s3 -d
package s3

import (
	"context"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
)

const (
	testBucketName = "test-bucket"
	testKey        = "test-key"
)

func initState() (context.Context, *Store) {
	ctx := context.Background()
	creds := NewCredentials("0.0.0.0:9000", "minio", "minio123")
	store, err := New(creds)
	if err != nil {
		panic(err)
	}
	if err := store.s3.MakeBucket(ctx, testBucketName, minio.MakeBucketOptions{
		Region:        "",
		ObjectLocking: false,
	}); err != nil {
		panic(err)
	}
	return ctx, store
}

func tearUpStore(ctx context.Context, store *Store) {
	objectsCh := store.s3.ListObjects(ctx, testBucketName, minio.ListObjectsOptions{
		Recursive: true,
	})
	for object := range objectsCh {
		if object.Err != nil {
			panic(object.Err)
		}
		if err := store.s3.RemoveObject(ctx, testBucketName, object.Key, minio.RemoveObjectOptions{}); err != nil {
			panic(err)
		}
	}
	_ = store.s3.RemoveBucket(ctx, testBucketName)
}

func Test_upload(t *testing.T) {
	ctx, store := initState()
	defer tearUpStore(ctx, store)

	object := NewObject(testBucketName, testKey, []byte("test data"), ContentTypeText)
	err := store.Upload(ctx, object)
	assert.NoError(t, err)

	retrievedObject, err := store.Download(ctx, testBucketName, testKey)
	assert.NoError(t, err)
	assert.Equal(t, object.GetKey(), retrievedObject.GetKey())
	assert.Equal(t, object.GetData(), retrievedObject.GetData())
	assert.Equal(t, object.GetContentType(), retrievedObject.GetContentType())
}

func Test_updates(t *testing.T) {
	ctx, store := initState()
	defer tearUpStore(ctx, store)

	object := NewObject(testBucketName, testKey, []byte("test data"), ContentTypeText)
	err := store.Upload(ctx, object)
	assert.NoError(t, err)

	updatedObject := NewObject(testBucketName, testKey, []byte("updated data"), ContentTypeText)
	err = store.Upload(ctx, updatedObject)
	assert.NoError(t, err)
	retrievedObject, err := store.Download(ctx, testBucketName, testKey)
	assert.NoError(t, err)
	assert.Equal(t, object.GetKey(), retrievedObject.GetKey())
	assert.Equal(t, updatedObject.GetData(), retrievedObject.GetData())
	assert.Equal(t, updatedObject.GetData(), retrievedObject.GetData())
	assert.Equal(t, object.GetContentType(), retrievedObject.GetContentType())
}
