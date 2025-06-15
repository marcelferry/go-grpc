package main

import (
	"context"
	"grpc-go/pb"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"
)

func main() {
	log.Printf("connecting to MongoDB...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Create MongoDB client
	// NOTE: "mongodb+srv" protocol means connect to Altas cloud MongoDB server
	//       use just "mongodb" if you connect to on-premise MongoDB server
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://USER:PASSWORD@SERVER/experiments"))
	if err != nil {
		log.Fatalf("failed to create new MongoDB client: %#v", err)
	}

	// Connect client
	if err = client.Connect(ctx); err != nil {
		log.Fatalf("failed to connect to MongoDB: %#v", err)
	}

	log.Printf("connected successfully")

	// Get collection from database
	coll := client.Database("experiments").Collection("proto")

	// Create protobuf Timestamp value from golang Time
	t := time.Now()
	ts, err := ptypes.TimestampProto(t)
	if err != nil {
		log.Fatalf("failed to convert golang Time to protobuf Timestamp: %#v", err)
	}

	// Fill in data structure
	in := pb.Data{
		BoolValue:   true,
		Int64Value:  12345,
		DoubleValue: 123.45,
		StringValue: "qwerty",

		TimestampValue: ts,

		BoolWrappedValue:   &wrappers.BoolValue{Value: true},
		Int64WrappedValue:  &wrappers.Int64Value{Value: 12345},
		DoubleWrappedValue: &wrappers.DoubleValue{Value: 123.45},
		StringWrappedValue: &wrappers.StringValue{Value: "qwerty"},
	}

	log.Printf("insert data into collection <experiments.proto>...")

	// Insert data into the collection
	res, err := coll.InsertOne(ctx, &in)
	if err != nil {
		log.Fatalf("insert data into collection <experiments.proto>: %#v", err)
	}
	id := res.InsertedID
	log.Printf("inserted new item with id=%v successfully", id)

	// Create filter and output structure to read data from collection
	var out Data
	filter := bson.D{{Key: "_id", Value: id}}

	// Read data from collection
	err = coll.FindOne(ctx, filter).Decode(&out)
	if err != nil {
		log.Fatalf("failed to read data (id=%v) from collection <experiments.proto>: %#v", id, err)
	}

	log.Print("read successfully")
}
