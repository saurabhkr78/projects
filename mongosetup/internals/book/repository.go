/*
How do I talk to MongoDB?

Responsibilities:

InsertOne
FindOne
Find
UpdateOne
DeleteOne
*/

package book

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(collection *mongo.Collection) *Repository {
	return &Repository{
		collection: collection,
	}
}

func (r *Repository) CreateBook(ctx context.Context, book *Book) (*Book, error) {
	result, err := r.collection.InsertOne(ctx, book)
	if err != nil {
		return nil, err
	}

	// We'll use result.InsertedID once we change Book.ID
	// to primitive.ObjectID.

	_ = result

	return book, nil
}

func (r *Repository) GetBooks(ctx context.Context) ([]Book, error) {
	// MongoDB Find()
}

func (r *Repository) GetBookByID(ctx context.Context, id string) (*Book, error) {
	// MongoDB FindOne()
}

func (r *Repository) UpdateBook(ctx context.Context, id string, book *Book) (*Book, error) {
	// MongoDB UpdateOne()
}

func (r *Repository) DeleteBook(ctx context.Context, id string) error {
	// MongoDB DeleteOne()
}
