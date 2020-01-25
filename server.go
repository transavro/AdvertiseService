package main

import (
	"fmt"
	codecs "github.com/amsokol/mongo-go-driver-protobuf"
	"github.com/transavro/AdvertiseService/apihandler"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "github.com/transavro/AdvertiseService/proto"
	"time"
)

const (
	atlasMongoHost          = "mongodb://nayan:tlwn722n@cluster0-shard-00-00-8aov2.mongodb.net:27017,cluster0-shard-00-01-8aov2.mongodb.net:27017,cluster0-shard-00-02-8aov2.mongodb.net:27017/test?ssl=true&replicaSet=Cluster0-shard-0&authSource=admin&retryWrites=true&w=majority"
	developmentMongoHost = "mongodb://dev-uni.cloudwalker.tv:6592"
	RedisHost   = ":6379"
	grpc_port        = ":7775"
	rest_port		 = ":7776"
)

var advertiseCollection *mongo.Collection

// Multiple init() function
func init() {
	fmt.Println("Welcome to init() function")
	advertiseCollection = getMongoCollection("cloudwalker", "advertise", atlasMongoHost)
}

func getMongoCollection(dbName, collectionName, mongoHost string) *mongo.Collection {
	// Register custom codecs for protobuf Timestamp and wrapper types
	reg := codecs.Register(bson.NewRegistryBuilder()).Build()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoHost), options.Client().SetRegistry(reg))
	if err != nil {
		log.Fatal(err)
	}
	return mongoClient.Database(dbName).Collection(collectionName)
}



func startGRPCServer(address string) error {
	// create a listener on TCP port
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	} // create a server instance
	s := apihandler.Server{
		advertiseCollection,
	}

	// attach the Ping service to the server
	grpcServer := grpc.NewServer()                    // attach the Ping service to the server
	pb.RegisterAdvertiseServiceServer(grpcServer, &s) // start the server
	log.Printf("starting HTTP/2 gRPC server on %s", address)
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %s", err)
	}
	return nil
}

func main() {

	// fire the gRPC server in a goroutine
	go func() {
		err := startGRPCServer(grpc_port)
		if err != nil {
			log.Fatalf("failed to start gRPC server: %s", err)
		}
	}()

	// infinite loop
	//log.Printf("Entering infinite loop")
	select {}
}