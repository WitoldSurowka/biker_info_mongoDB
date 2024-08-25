package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// NewMongoDBRepository creates a new instance of MongoDBRepository
func NewMongoDBRepository(uri, dbName, collectionName string) (*MongoDBRepository, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	collection := client.Database(dbName).Collection(collectionName)

	return &MongoDBRepository{
		client:     client,
		collection: collection,
	}, nil
}

func (repo *MongoDBRepository) ReadData() ([]Record, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Find all records
	cursor, err := repo.collection.Find(ctx, bson.D{{"phoneNumber", bson.D{{"$exists", true}}}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var records []Record
	for cursor.Next(ctx) {
		var record Record
		if err := cursor.Decode(&record); err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

// AddRecord adds a new record to the MongoDB collection
func (repo *MongoDBRepository) AddRecord(phoneNumber, city string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Check if the record already exists
	var existingRecord Record
	err := repo.collection.FindOne(ctx, bson.M{"phoneNumber": phoneNumber, "city": city}).Decode(&existingRecord)
	if err == nil {
		// Record exists, so delete it
		_, err := repo.collection.DeleteOne(ctx, bson.M{"id": existingRecord.ID})
		if err != nil {
			return err
		}
		InfoStatusDeleted(phoneNumber, city)
	} else {

		// Get the current highest ID in the collection
		var iDCounterDocument IDCounterDocument
		filter := bson.D{{"currentID", bson.D{{"$exists", true}}}}
		opts := options.FindOne().SetSort(bson.D{{"currentID", -1}})
		err := repo.collection.FindOne(ctx, filter, opts).Decode(&iDCounterDocument)
		if err != nil && err != mongo.ErrNoDocuments {
			return err
		}
		newID := iDCounterDocument.CurrentID + 1

		// Add the new record with the current date as SubscriptionDate
		_, err = repo.collection.InsertOne(ctx, Record{
			ID:               newID,
			PhoneNumber:      phoneNumber,
			City:             city,
			SubscriptionDate: time.Now(), // Assign current date and time
		})
		if err != nil {
			return err
		} else {
			InfoStatusAdded(phoneNumber, city)
		}

		// Update CurrentID
		_, err = repo.collection.UpdateOne(
			ctx,
			bson.M{},
			bson.D{{"$set", bson.D{{"currentID", newID}}}},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteRecord deletes a record by ID from the MongoDB collection
func (repo *MongoDBRepository) DeleteRecord(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := repo.collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}
