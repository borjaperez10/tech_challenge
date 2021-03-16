package main

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/borjaperez10/grpc_example1/communication"
)

func main() {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := communication.NewCommunicationServiceClient(conn)

	response, err := c.SellInvoice(context.Background(), &communication.Invoice{CompanyID: "COMPANYID", Price:float64(1000)})
	if err != nil {
		log.Fatalf("Error when calling SellInvoice: %s", err)
	}
	log.Printf("Response from server: %s", response.Body)


}
