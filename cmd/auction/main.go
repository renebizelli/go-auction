package main

import (
	"context"
	"fmt"
	"renebizelli/go-leilao/Internals/entity/user_entity"
	"renebizelli/go-leilao/Internals/infra/api/web/controllers/auction_controller"
	"renebizelli/go-leilao/Internals/infra/api/web/controllers/bid_controller"
	"renebizelli/go-leilao/Internals/infra/api/web/controllers/user_controller"
	"renebizelli/go-leilao/Internals/infra/database/auction"
	"renebizelli/go-leilao/Internals/infra/database/bid"
	"renebizelli/go-leilao/Internals/infra/database/user"
	"renebizelli/go-leilao/Internals/usecase/auction_usecase"
	"renebizelli/go-leilao/Internals/usecase/bid_usecase"
	"renebizelli/go-leilao/Internals/usecase/user_usecase"
	"renebizelli/go-leilao/configuration/database/mongodb"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {

	ctx := context.Background()

	if err := godotenv.Load("./.env"); err != nil {
		panic("Error loading .env file")
	}

	database, err := mongodb.NewMongoDBConnection(ctx)

	if err != nil {
		panic("Error connecting to MongoDB")
	}

	user_controller,
		auction_controller,
		bid_controller := initDependencies(ctx, database, true)

	router := gin.Default()

	router.GET("/auctions", auction_controller.FindAuctions)
	router.GET("/auctions/:auctionsId", auction_controller.FindAuctionById)
	router.GET("/auction/winner/:auctionId", auction_controller.FindWinningBidByAuctionId)
	router.POST("/auctions", auction_controller.CreateAuction)

	router.GET("/bid/:auctionId", bid_controller.FindBidByAuctionId)
	router.POST("/bid", bid_controller.CreateBid)

	router.GET("/user/:user", user_controller.FindUserByID)
	router.POST("/user/", user_controller.CreateUser)

	router.Run(":8080")
}

func initDependencies(ctx context.Context, database *mongo.Database, initialData bool) (
	*user_controller.UserController,
	*auction_controller.AuctionController,
	*bid_controller.BidController,
) {

	auctionRepository := auction.NewAuctionRepository(database)

	bidRepository := bid.NewBidRepository(database, auctionRepository)
	bidUseCase := bid_usecase.NewBidUseCase(bidRepository)
	bidController := bid_controller.NewBidController(bidUseCase)

	auctionUseCase := auction_usecase.NewAuctionUseCase(auctionRepository, bidRepository)
	auctionController := auction_controller.NewAuctionController(auctionUseCase)

	userRepositoty := user.NewUserRepository(database)
	userUseCase := user_usecase.NewUserUseCase(userRepositoty)
	userController := user_controller.NewUserController(userUseCase)

	if initialData {
		initData(ctx, userRepositoty)
	}

	return userController, auctionController, bidController
}

func initData(ctx context.Context, userRepositoty user_entity.UserRepositoryInterface) {

	user := &user_entity.User{
		ID:   "c3f62675-895e-4306-a5c1-bcf4df0ae7c6",
		Name: "René Bizelli",
	}

	_, e := userRepositoty.FindUserByID(ctx, user.ID)
	if e != nil {
		if err := userRepositoty.CreateUser(ctx, user); err != nil {
			panic("Error creating user")
		}
	}

	fmt.Println("User created:", user.ID)
}
