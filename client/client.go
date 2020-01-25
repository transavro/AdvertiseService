package main

import (
	"context"
	pb "github.com/transavro/AdvertiseService/proto"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main()  {

	conn, err := grpc.Dial("localhost:7775", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	addServer := pb.NewAdvertiseServiceClient(conn)


	streamResp, err := addServer.GetAdvertise(context.Background(), &pb.GetAdd{Categories:[]string{"LifeStyle"}})
	if err != nil {
		log.Fatal(err)
	}

	for{

		add, err := streamResp.Recv()

		if err != nil {
			if err == io.EOF {
				break
			}else {
				log.Fatal(err)
			}
		}
		log.Println(add)
	}
	log.Println("Out for loop ")


	//resp , err := addServer.DeleteAdvertise(context.Background(), &pb.DeleteAdvertiseReq{Title:"Nirma", Brand:"Sabun"})
	//
	//if err != nil {
	//	log.Println(err)
	//}
	//log.Println(resp)

	//
	//ts, _ := ptypes.TimestampProto(time.Now())
	//
	//newAdd := pb.Advertise{
	//	StartTime:            ts,
	//	EndTime:              ts,
	//	Title:                "Nirma",
	//	Image:                "dummy Image Url",
	//	Video:                "dummy Video Url",
	//	Target:               &pb.Target{
	//		Package:              "com.google.android.youtube.tv",
	//		Url:                  "add target url",
	//	},
	//	Genre:                nil,
	//	Language:             []string{"Telegu"},
	//	Categories:           []string{"LifeStyle", "Family"},
	//	Position:             "1.1.1",
	//	AdversiteType:        pb.AdversiteType_TILE,
	//	ViewCount:            0,
	//	ClickCount:           0,
	//	ViewDuration:         nil,
	//	Brand:                "Sabun",
	//}
	//
	//result, err := addServer.UpdateAdvertise(context.Background(), &newAdd)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Println(result.String())
}