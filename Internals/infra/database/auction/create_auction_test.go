package auction

import (
	"context"
	"os"
	"renebizelli/go-leilao/Internals/entity/auction_entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestCreateAuction(t *testing.T) {

	os.Setenv("AUCTION_INTERVAL", "1s")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.TODO())

	database := client.Database("testdb")
	coll := database.Collection("auctions")

	coll.DeleteMany(context.Background(), bson.M{})

	auction := &auction_entity.Auction{
		ID:          "test_auction_id",
		ProductName: "Test Product",
		Category:    "Test Category",
		Description: "Test Description",
		Condition:   auction_entity.New,
		Status:      auction_entity.Active,
		Timestamp:   time.Now(),
	}

	repo := NewAuctionRepository(database)
	errCreate := repo.CreateAuction(context.Background(), auction)
	if errCreate != nil {
		t.Fatalf("Failed to create auction: %v", errCreate)
	}

	time.Sleep(2 * time.Second)

	var result AuctionEntityMongo
	errFind := coll.FindOne(context.Background(), bson.M{"_id": auction.ID}).Decode(&result)
	if errFind != nil {
		t.Fatalf("Failed to find auction in database: %v", errFind)
	}

	assert.Equal(t, auction_entity.AuctionStatus(1), result.Status)

}
