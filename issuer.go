package main

import (
	"log"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/borjaperez10/tech_challenge/communication"
)

func checkError(err error, handler func(e error)) {
    if err != nil {
        handler(err)
    }
}


func main() {
	//Handler for erorrs
	handler := func(e error) {
        log.Fatal(e)
	}
	var companyNIF, companyName, invoiceName string
	var conn *grpc.ClientConn
	var total, amount float64
	
	//Create connection with gRPC Server and defer it: when finishes it will be closed
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	checkError(err, handler)
	defer conn.Close()
	c := communication.NewCommunicationServiceClient(conn)

	
	//CheckConnectivity first of all
	_, err = c.CheckConnectivity(context.Background(),  &communication.EmptyRequest{})
	checkError(err, handler)



	choose:= menuChooseOption()
	switch choose {
	case 1://Register a new Issuer in the database (with the NIF)
		companyNIF,companyName = menuRegisterIssuer()		
		_, err := c.IntroduceIssuerToDatabase(context.Background(), &communication.Issuer{Nif: companyNIF, Name: companyName})
		checkError(err, handler)

	case 2://Sell an invoice

		companyNIF, invoiceName, total, amount = menuFinanceInvoice()	
		_, err := c.SellInvoice(context.Background(), &communication.Invoice{CompanyID: companyNIF, InvoiceName: invoiceName, TotalPrice: total, AmountToReceive: amount})
		checkError(err, handler)

	case 3://Read market ledger

		response, err := c.ReadMarketLedger(context.Background(), &communication.EmptyRequest{})
		checkError(err, handler)
		fmt.Println(response.Body)

	default:
		log.Fatal("Invalid option")
	}
	 
	 
	


}

func menuRegisterIssuer()(nif string, name string){
		//This function is used to get the data in order to register an issuer
		// Needed the NIF of the Issuer and the name of the issuer
		fmt.Println("Please, enter the NIF of the issuer to register")
		fmt.Scanln(&nif)
		fmt.Println("Please, enter the name of the issuer")
		fmt.Scanln(&name)
		return 
}

func menuFinanceInvoice()(nif string, invoice string, total float64, amount float64){
		//This function is used to gett the data for a new sell request
		//The needed data is the NIF, The invice name, the total price of the invoice, and the amount that is requested
		fmt.Println("Please, enter a registered NIF")
		fmt.Scanln(&nif)
		fmt.Println("Please, enter the name of the invoice")
		fmt.Scanln(&invoice)
		fmt.Println("Please, enter the total invoice price")
		fmt.Scanln(&total)	
		fmt.Println("Please, enter the amount you want to receive")
		fmt.Scanln(&amount)
		return 
}
func menuChooseOption()(option int){
		// A menu will be displayed to choose the option for the Issuer: 
		// He must be registered in the database with the NIF in order to sell invoices'
		// Options: 1: register an issuer
		//          2: Sell an invoice
		fmt.Println("**********************************************************************")
		fmt.Println("Choose an option:")
		fmt.Println("1.-Register an issuer in the database:")
		fmt.Println("2.-Sell an invoice. This means that the issuer is already registered:")
		fmt.Println("3.-Read the Market Ledger")
		fmt.Println("**********************************************************************")
		fmt.Scanln(&option)
		return 
}

