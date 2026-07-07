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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var books []Book

	if err := cursor.All(ctx, &books); err != nil {
		return nil, err
	}

	return books, nil
}

func (r *Repository) GetBookByID(ctx context.Context, id string) (*Book, error) {

	// Convert the hexadecimal string from the URL
	// into MongoDB's ObjectID type.
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var book Book

	err = r.collection.FindOne(
		ctx,
		bson.M{"_id": objectID},
	).Decode(&book)

	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (r *Repository) UpdateBook(ctx context.Context, id string, book *Book) (*Book, error) {
	// MongoDB UpdateOne()
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id": objectID,
	}

	update := bson.M{
		"$set": bson.M{
			"title":  book.Title,
			"author": book.Author,
			"year":   book.Year,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return r.GetBookByID(ctx, id)
}

func (r *Repository) DeleteBook(ctx context.Context, id string) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": objectID,
	}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
