package database

import (
	"github.com/google/uuid"
	"github.com/segmentio/ksuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"reflect"
)

var client *mongo.Client
var db *mongo.Database
var accounts *mongo.Collection

func Init() {
	tUUID := reflect.TypeOf(uuid.UUID{})

	registry := bson.NewRegistryBuilder().
		RegisterTypeEncoder(tUUID, bsoncodec.ValueEncoderFunc(encodeUUID)).
		RegisterTypeDecoder(tUUID, bsoncodec.ValueDecoderFunc(decodeUUID)).
		Build()

	var err error
	client, err = mongo.NewClient(options.Client().SetRegistry(registry).ApplyURI(os.Getenv("MONGO_URL")))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(nil)
	if err != nil {
		log.Fatal(err)
	}

	db = client.Database("wormhole", options.Database())
	accounts = db.Collection("accounts", options.Collection())
}

func Close() {
	err := client.Disconnect(nil)
	if err != nil {
		log.Fatal(err)
	}
}

func GetAccount(id ksuid.KSUID) (*Account, error) {
	var acc Account
	err := accounts.FindOne(nil, bson.M{"id": id}).Decode(&acc)
	if err != nil {
		return nil, err
	}
	return &acc, nil
}

func GetAccountWithUsername(username string) (*Account, error) {
	var acc Account
	err := accounts.FindOne(nil, bson.M{"username": username}).Decode(&acc)
	if err != nil {
		return nil, err
	}
	return &acc, nil
}
