// How do I connect to MongoDB?
package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect establishes a connection to MongoDB and verifies it with a ping.
// It returns the connected MongoDB client or an error.
func Connect(uri string) (*mongo.Client, error) {

	// Create a context with a 10-second timeout so the connection
	// attempt doesn't hang indefinitely.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a MongoDB client using the provided connection URI.
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Verify that MongoDB is reachable.
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}
