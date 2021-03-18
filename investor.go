package main

import (
	"log"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/borjaperez10/grpc_example1/communication"
	"strconv"

	"io"
)


type InvoiceStruct struct {
  ID       	 int32
  Name		 string
  Issuer     string
  Total  	 float64
  Amount     float64
  Closed	 string
}

var invoices []InvoiceStruct


func checkError(err error, handler func(e error)) {
    if err != nil {
        handler(err)
    }
}


func main() {
	


	var conn *grpc.ClientConn
	var investorDNI, investorName string
	var availableMoney, retainedMoney float64 
	var bidTotal, bidAmount string
		//Handler for erorrs
	handler := func(e error) {
        log.Fatal(e)
	}



	//Create connection with gRPC Server and defer it: when finishes it will be closed
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	checkError(err, handler)
	defer conn.Close()
	c:= communication.NewCommunicationServiceClient(conn)


	//CheckConnectivity first of all
	_, err = c.CheckConnectivity(context.Background(),  &communication.EmptyRequest{})
	checkError(err, handler)


	choose:= menuChooseOption()
	switch choose {
	case 1: //Register a new Investor in the database (with the DNI)
		
		investorDNI,investorName, availableMoney= menuRegisterInvestor()
		retainedMoney = float64(0)
		_, err := c.IntroduceInvestorToDatabase(context.Background(), &communication.Investor{Dni: investorDNI, Name: investorName, AvailableMoney:availableMoney, RetainedMoney: retainedMoney}) //at first retained = 0
		checkError(err, handler)


	case 2: //Buy an invoice from the available Invoices in the database 

		investorDNI = menuIdentification()	

		// We will create a stream
		// Iterating from the stream we will get the responses (Invoices) from the server
		// The responses that we get will be appended obtaining the variable invoices with all the available invoices
		stream, err := c.GetAvailableInvoices(context.Background(),  &communication.EmptyRequest{})
		checkError(err, handler)
		
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				//means stream is finished
				break
			}
			checkError(err, handler)  
			i := InvoiceStruct{
				ID: resp.Id,
				Name:resp.InvoiceName,
				Issuer:resp.CompanyID,
				Total:resp.TotalPrice,
				Amount: resp.AmountToReceive,
				Closed:resp.Closed,
			}
			invoices = append(invoices, i)
		}
		
		// Once we have the array with the invoices, we will print them. 
		// Investor must select which one wants to buy selecting its identifier
		// If there are no invoices system will return -1
		// If a bid is selected, the investor will be requested to introduce the total bid and the amount of bid 
		option:=printOptions()
		if (option != -1){
			option = option -1 //to take the index from the array 
			bidTotal, bidAmount = menuMakeBid()

			bidTotalF, _ := strconv.ParseFloat(bidTotal, 64);
			bidAmountF, _ := strconv.ParseFloat(bidAmount, 64);

			//We will try to modify the investor's bank account,putting the bid amount as retained money
			response, err := c.TryToModifyInvestorMoney(context.Background(), &communication.Bid{Dni: investorDNI, Total:bidTotalF, Amount:bidAmountF, Action:"rm"}) 
			checkError(err, handler)

			if response.Body == "yes"{
				fmt.Println("Bid has been placed. The amount of money has been retained in your bank account.")

				// We will create an Invoice Part with the details of the transaction and introduce it in the database. 
				response, err = c.TryToCreateInvoicePart(context.Background(), &communication.InvoicePart{OriginalId: invoices[option].ID, Total: bidTotalF, Amount:bidAmountF, Buyer:investorDNI }) 
		
				if response.Body == "0"{
					fmt.Println("Bids have been placed for 100% of the invoice")
				}else if response.Body == "1"{
					fmt.Println("Bid has been rejected due to discout difference")
					break
				}else if response.Body == "2"{
					fmt.Println("Bid has been placed but the invoice has not been 100% financed yet")
					break
				}
			}else{
				fmt.Println("Bid can not be placed: you dont have enough balance")
			}

			
		}	

	default:
		log.Fatal("Invalid option")
	}




	

}






func printOptions()(option int ){
		//This function prints all the available invoices from the object invoices
		// The investor must select the one that wants to buy
		count := 1

		fmt.Println("**********************************************************************")
		fmt.Println("Select the identifier of the invoice you want to finance from the currently available ones:")
		for _,invoice := range invoices {
			if invoice.Closed == "No" { 
				idVal := strconv.Itoa(int(invoice.ID))
				line :=  idVal +"- " + fmt.Sprintf("Seller %s wants to sell invoice %s. The price is %f (the seller wants to receive %f)", invoice.Name, invoice.Issuer, invoice.Total, invoice.Amount)
				fmt.Println(line)
				count++;
			}
		}
		fmt.Println("**********************************************************************")
		if count == 1{
			fmt.Println("There are no available invoices")
			option = -1
		}else{
			fmt.Scanln(&option)
		}
		return 
}

func menuChooseOption()(option int){
		// A menu will be displayed to choose the option for the Investor: 
		// He must be registered in the database with the DNI in order to buy invoices'
		// Options: 1: register an investor
		//          2: Buy an invoice from the available ones
		fmt.Println("**********************************************************************")
		fmt.Println("Choose an option:")
		fmt.Println("1.-Register an inversor in the database:")
		fmt.Println("2.-Place a bid over the available invoices. This means that the investor is already registered:")
		fmt.Println("**********************************************************************")
		fmt.Scanln(&option)
		return 
}

func menuMakeBid()(tot, amo string){
		//This function requires the investor to introduce the total bid and the amount of bid prices
			//Total bid is the amount of money he wants to receive back
			//Amount bid is the amount of money he will pay for the invoice
		fmt.Println("**********************************************************************")
		fmt.Println("Place your bid:")
		fmt.Scanln(&tot)
		fmt.Println("Place the amount:")
		fmt.Scanln(&amo)
		return 
}

func menuIdentification()(dni string){
		//Menu for identificating investor
		fmt.Println("Please, enter a registered DNI")
		fmt.Scanln(&dni)
		return 
}

func menuRegisterInvestor()(dni string, name string, money float64){
		//A menu that gets the data for registering a new investor
		// DNI 
		// name
		// Money in bank account
		fmt.Println("Please, enter the DNI of the inversor to register")
		fmt.Scanln(&dni)
		fmt.Println("Please, enter the name of the issuer")
		fmt.Scanln(&name)
		fmt.Println("Please, enter the amount of money in your account. This will be added as available money")
		fmt.Scanln(&money)
		return 
}