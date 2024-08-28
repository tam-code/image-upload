package repositories

import (
	"testing"
	"time"

	"github.com/tam-code/lrn/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"gotest.tools/assert"
)

func TestCreateUploadLink(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.ClearCollections()

	uploadLink := models.UploadLink{
		ExpirationTime: time.Now().Add(time.Hour),
	}

	tests := []struct {
		name             string
		prepare          func(mt *mtest.T)
		expectError      bool
		expectUploadLink bool
	}{
		{
			name: "create upload link",
			prepare: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateSuccessResponse())
			},
			expectError:      false,
			expectUploadLink: true,
		},
		{
			name: "simple error",
			prepare: func(mt *mtest.T) {
				mt.AddMockResponses(bson.D{{"ok", 0}})
			},
			expectError:      true,
			expectUploadLink: false,
		},
	}

	for _, test := range tests {
		mt.Run(test.name, func(mt *mtest.T) {
			repo := uploadLinkRepository{
				mongoCollection: mt.Coll,
			}

			test.prepare(mt)

			insertedLink, err := repo.CreateUploadLink(uploadLink)
			assert.Equal(t, test.expectError, err != nil)
			assert.Equal(t, test.expectUploadLink, insertedLink != nil)
		})
	}

}

func TestGetUploadLinkByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.ClearCollections()

	uploadLink := models.UploadLink{
		ID:             "5f9f1f1b6f6b589b3f3b3b3b",
		ExpirationTime: time.Now().Add(time.Hour),
	}

	tests := []struct {
		name             string
		prepare          func(mt *mtest.T)
		expectError      bool
		expectUploadLink bool
	}{
		{
			name: "get upload link by id",
			prepare: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
					{"_id", uploadLink.ID},
					{"expiration_time", uploadLink.ExpirationTime.Format(time.RFC3339)},
				}))
			},
			expectError:      false,
			expectUploadLink: true,
		},
		{
			name: "simple error",
			prepare: func(mt *mtest.T) {
				mt.AddMockResponses(bson.D{{"ok", 0}})
			},
			expectError:      true,
			expectUploadLink: false,
		},
	}

	for _, test := range tests {
		mt.Run(test.name, func(mt *mtest.T) {
			repo := uploadLinkRepository{
				mongoCollection: mt.Coll,
			}

			test.prepare(mt)

			getUploadLink, err := repo.GetUploadLinkByID(uploadLink.ID)
			assert.Equal(t, test.expectError, err != nil)
			assert.Equal(t, test.expectUploadLink, getUploadLink != nil)
		})
	}
}
