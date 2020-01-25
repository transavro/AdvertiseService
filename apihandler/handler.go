package apihandler

import (
	"fmt"
	pb "github.com/transavro/AdvertiseService/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type Server struct {
	AdvertiseColl *mongo.Collection
}

func(s Server) CreateAdvertise(ctx context.Context, reqAdd *pb.Advertise) (*pb.Advertise, error) {
	log.Print("hi create add ")
	// creating mongo find query
	findFilter := bson.D{{"$and", bson.A{bson.D{{"brand", reqAdd.GetBrand()}}, bson.D{{"title", reqAdd.GetTitle()}}}}}

	// hitting mongo db with the find query
	result := s.AdvertiseColl.FindOne(ctx, findFilter)

	// checking if the result has error
	if result.Err() != nil {
		// if the error is equales to document not found, means there is no document for the related query filter
		if result.Err() == mongo.ErrNoDocuments{

			// inserting the new document in the collection
			_ , err := s.AdvertiseColl.InsertOne(ctx, reqAdd)

			// check if any error while inserting in mongo db
			if err != nil {
				return nil, status.Error(codes.Internal,  fmt.Sprintf("Error while inserting data %s ", result.Err()))
			}
			// response send
			return reqAdd, nil
		}else {
			// if there is some another error than "mongo.ErrNilDocument" throw the error in response.
			return nil, status.Error(codes.Internal,  fmt.Sprintf("Error while finding data in db %s ", result.Err()))
		}
	}else {
		// when advertise already exits.
		return nil, status.Error(codes.AlreadyExists,  fmt.Sprintf("Advertise All ready exits %s ", result.Err()))
	}
}

func(s Server) GetAdvertise(req *pb.GetAdd, stream pb.AdvertiseService_GetAdvertiseServer) error {
	// making pipeling
	pipeline := mongo.Pipeline{}

	var filterArray []bson.E

	// checking if there is any genre
	if len(req.GetGenre()) > 0 {
		filterArray = append(filterArray, bson.E{"genre", bson.D{{"$in", req.GetGenre()}}})
	}

	// checkin if there is any language
	if len(req.GetLanguage()) > 0 {
		filterArray = append(filterArray, bson.E{"language", bson.D{{"$in", req.GetLanguage()}}})
	}

	// checkin if there is any Categories
	if len(req.GetCategories()) > 0 {
		filterArray = append(filterArray, bson.E{"categories", bson.D{{"$in", req.GetCategories()}}})
	}

	// adding first pipeling stage
	pipeline = append(pipeline, bson.D{{"$match", filterArray}})

	// hitting db with aggregation pipeliine
	cur, err := s.AdvertiseColl.Aggregate(context.Background(), pipeline)
	// check if there is any error while getting data from db
	if err != nil {
		return  status.Error(codes.Internal, fmt.Sprintf("Error while fetch data from db  %s ", err.Error()))
	}

	// looping throught the cursor
	for cur.Next(context.Background()){
		var add pb.Advertise
		// decoding the mongo document data in pb object
		err = cur.Decode(&add)

		// check if there is any error while decoding data.
		if err != nil {
			return  status.Error(codes.Internal, fmt.Sprintf("Error while decoding data  %s ", err.Error()))
		}

		// sending data to the response
		err = stream.Send(&add)

		// checking if there is error while sending
		if err != nil {
			return  status.Error(codes.Internal, fmt.Sprintf("Error while sending data  %s ", err.Error()))
		}
	}

	//finally closing the cursor
	cur.Close(context.Background())
	return nil

}

func(s Server) UpdateAdvertise (ctx context.Context, reqAdd *pb.Advertise) (*pb.Advertise, error) {
	// first find the document in db to update it
	findFilter := bson.D{{"$and", bson.A{bson.D{{"brand", reqAdd.GetBrand()}}, bson.D{{"title", reqAdd.GetTitle()}}}}}

	//hitting mongo db with find filter
	result := s.AdvertiseColl.FindOne(ctx, findFilter)

	// checking if the result has error
	if result.Err() != nil {
		return nil, status.Error(codes.NotFound,  fmt.Sprintf("Error while finding data in db %s ", result.Err()))
	}

	// replacing the document in db
	_, err := s.AdvertiseColl.ReplaceOne(ctx, findFilter, reqAdd)

	// check if there is any error while replacing.
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition,  fmt.Sprintf("Error while replacing data in db %s ", err))
	}

	return reqAdd, nil
}

func(s Server) DeleteAdvertise( ctx context.Context, req *pb.DeleteAdvertiseReq) (*pb.DeleteAdvertiseResp, error) {

	// making delete query
	deleteFilter := bson.D{{"$and", bson.A{bson.D{{"brand", req.GetBrand()}}, bson.D{{"title", req.GetTitle()}}}}}

	//making mongo delete reqIsSucessfulluest
	_ , err := s.AdvertiseColl.DeleteOne(ctx, deleteFilter)

	// check if there is any error
	if err != nil {
		return nil, status.Error(codes.NotFound,  fmt.Sprintf("Error while deleting data from db %s ", err))
	}

	//sending resp
	return &pb.DeleteAdvertiseResp{IsSucessfull:true}, nil
}



























