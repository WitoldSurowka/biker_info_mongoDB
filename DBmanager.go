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

// ------ methods for users collection ------

func (repoUsers *MongoDBRepository) UsersReadData() ([]Record, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Find all records
	cursor, err := repoUsers.collection.Find(ctx, bson.D{{"phoneNumber", bson.D{{"$exists", true}}}})
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
func (repoUsers *MongoDBRepository) UsersAddRecord(repoOutboundMessages *MongoDBRepository, phoneNumber, city string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Check if the record already exists
	var existingRecord Record
	err := repoUsers.collection.FindOne(ctx, bson.M{"phoneNumber": phoneNumber, "city": city}).Decode(&existingRecord)
	if err == nil {
		// Record exists, so delete it
		repoUsers.UsersDeleteRecord(repoOutboundMessages, existingRecord.ID, phoneNumber, city)
	} else {

		// Get the current highest ID in the collection
		var iDCounterDocument IDCounterDocument
		filter := bson.D{{"currentID", bson.D{{"$exists", true}}}}
		opts := options.FindOne().SetSort(bson.D{{"currentID", -1}})
		err := repoUsers.collection.FindOne(ctx, filter, opts).Decode(&iDCounterDocument)
		if err != nil && err != mongo.ErrNoDocuments {
			return err
		}
		newID := iDCounterDocument.CurrentID + 1

		// Add the new record with the current date as SubscriptionDate
		_, err = repoUsers.collection.InsertOne(ctx, Record{
			ID:               newID,
			PhoneNumber:      phoneNumber,
			City:             city,
			SubscriptionDate: time.Now(), // Assign current date and time
		})
		if err != nil {
			return err
		} else {
			InfoStatusAdded(repoOutboundMessages, phoneNumber, city)
		}

		// Update CurrentID
		_, err = repoUsers.collection.UpdateOne(
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
func (repoUsers *MongoDBRepository) UsersDeleteRecord(repoOutboundMessages *MongoDBRepository, id int, phoneNumber string, city string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := repoUsers.collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return err
	}
	InfoStatusDeleted(repoOutboundMessages, phoneNumber, city)

	return err
}

// ------ methods for outboundMessages collection ------

// OutboundSMSAddSMS ads an SMS to the MongoDB collection
func (repoOutboundSMS *MongoDBRepository) OutboundSMSAddSMS(phoneNumber, message string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := repoOutboundSMS.collection.InsertOne(ctx, SMS{
		PhoneNumber:  phoneNumber,
		Message:      message,
		CreationDate: time.Now(),
		Processed:    false,
	})
	if err != nil {
		return err
	}

	err = repoOutboundSMS.OutboundSMSDeleteSMS()

	return err
}

func (repoOutboundSMS *MongoDBRepository) OutboundSMSDeleteSMS() error {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	count, err := repoOutboundSMS.collection.CountDocuments(ctx, bson.M{"processed": true})
	if err != nil {
		return err
	}

	//var processedSMSes []SMS
	//if err := cursor.All(ctx, &processedSMSes); err != nil {
	//	return err
	//}
	// Sprawdź, czy znaleziono przetworzone SMSy i usuń je
	if count != 0 {
		_, err := repoOutboundSMS.collection.DeleteMany(ctx, bson.M{"processed": true})
		if err != nil {
			return err
		}
	}

	return err
}
