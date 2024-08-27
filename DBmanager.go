package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (repoUsers *MongoDBRepository) UsersReadData() ([]UserRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Find all records
	cursor, err := repoUsers.collection.Find(ctx, bson.D{{"phoneNumber", bson.D{{"$exists", true}}}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var records []UserRecord
	for cursor.Next(ctx) {
		var record UserRecord
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
	var existingRecord UserRecord
	err := repoUsers.collection.FindOne(ctx, bson.M{"phoneNumber": phoneNumber, "city": city}).Decode(&existingRecord)
	if err == nil {
		// Record exists, so delete it
		repoUsers.UsersDeleteRecord(repoOutboundMessages, existingRecord.ID, phoneNumber, city)
	} else {

		//TODO searching for the highest ID in the collection is not necessary no more. Let's just  use the 'current ID'
		// holder in IDCounterDocument
		//Get the current highest ID in the collection
		var iDCounterDocument IDCounterDocument
		filter := bson.D{{"currentID", bson.D{{"$exists", true}}}}
		opts := options.FindOne().SetSort(bson.D{{"currentID", -1}})
		err := repoUsers.collection.FindOne(ctx, filter, opts).Decode(&iDCounterDocument)
		if err != nil && err != mongo.ErrNoDocuments {
			return err
		}
		newID := iDCounterDocument.CurrentID + 1

		// Add the new record with the current date as SubscriptionDate
		_, err = repoUsers.collection.InsertOne(ctx, UserRecord{
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

	return err
}

func (repoInboundOutboundSMS *MongoDBRepository) InboundOutboundSMSDeleteProcessedSMS() error {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	count, err := repoInboundOutboundSMS.collection.CountDocuments(ctx, bson.M{"processed": true})
	if err != nil {
		return err
	}

	if count != 0 {
		_, err := repoInboundOutboundSMS.collection.DeleteMany(ctx, bson.M{"processed": true})
		if err != nil {
			return err
		}
	}

	return err
}

func (repoInboundSMS *MongoDBRepository) InboundOutboundMakeUnreadMongoSMSarrayFromDB() ([]MongoSMS, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{"creationDate", 1}})
	cursor, err := repoInboundSMS.collection.Find(ctx, bson.D{{"processed", false}}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var UnreadMongoSMSarray []MongoSMS
	for cursor.Next(ctx) {
		var UnreadSMS MongoSMS
		if err := cursor.Decode(&UnreadSMS); err != nil {
			return nil, err
		}
		UnreadMongoSMSarray = append(UnreadMongoSMSarray, UnreadSMS)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return UnreadMongoSMSarray, nil
}

func (repoInboundSMS *MongoDBRepository) InboundOutboundUpdateReadMongoSMSProcessedField(UnreadMongoSMSarray []MongoSMS) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Create a slice to store the IDs of the documents to be deleted.
	var ids []primitive.ObjectID
	for _, sms := range UnreadMongoSMSarray {
		ids = append(ids, sms.ID)
	}

	// Perform the deletion using the IDs
	filter := bson.M{"_id": bson.M{"$in": ids}}

	update := bson.M{
		"$set": bson.M{
			"processed": true,
		},
	}

	result, err := repoInboundSMS.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}

	fmt.Printf("Deleted %d documents\n", result)
	return nil
}
