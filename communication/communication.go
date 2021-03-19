package communication

import (
	"log"
	"fmt"
	"golang.org/x/net/context"
	"github.com/jinzhu/gorm"
	    "database/sql"
	_ "github.com/lib/pq" 
	
)


type Server struct {
}



type InvoiceStruct struct {
  gorm.Model
  ID       	 int32
  Name		 string
  Issuer     string
  Total  	 float64
  Amount     float64
  Closed	 string
}

var invoic InvoiceStruct
var db *gorm.DB
var err error


func checkError(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

// GetDB Returns a reference to the database
func OpenDB() (*sql.DB, error) {
	//Open database
	psqlconn := fmt.Sprintf("host=localhost port=5432 user=postgres dbname=TechChallenge sslmode=disable password=borjius10")
    db, err := sql.Open("postgres", psqlconn)
    return db, err
}
func CloseDB(db *sql.DB) error {
    return db.Close()
}


func (s *Server) CheckConnectivity(ctx context.Context, req *EmptyRequest) (*EmptyRequest, error) {	
	return &EmptyRequest{}, nil
}


func (s *Server) SellInvoice(ctx context.Context, invoice *Invoice) (*Message, error) {
	var returnMessage string
	var issId, issNam string

	db, err := OpenDB();
	defer CloseDB(db);
	checkError(err)

	
	//Check if issuer exists in the Database
	sqlStatement := `SELECT * FROM issuer WHERE nif=$1;`
	row := db.QueryRow(sqlStatement, invoice.CompanyID)
	err = row.Scan(&issId, &issNam)
	
	switch err {
	case sql.ErrNoRows:
	  returnMessage = "Issuer does not exist in the database. It must be added in order to sell an invoice."
	case nil:
		//Check that there is no any invoice with that same name name and company:
		sqlStatement = `SELECT * FROM invoice WHERE name=$1 and issuer=$2;`
		row := db.QueryRow(sqlStatement, invoice.InvoiceName, invoice.CompanyID)
		err = row.Scan(&invoic.ID, &invoic.Issuer)
		if err == sql.ErrNoRows{
			insertStatement := `insert into "invoice"("name", "issuer", "total", "toreceive", "closed") values($1, $2, $3, $4, $5)`
			_, err := db.Exec(insertStatement, invoice.InvoiceName, invoice.CompanyID, invoice.TotalPrice, invoice.AmountToReceive, "No")
			checkError(err)
			returnMessage = "Introduced in the database that issuer " + invoice.CompanyID + "sells an invoice " + invoice.InvoiceName + "of total value " +fmt.Sprintf("%f", invoice.TotalPrice) + " asking for " + fmt.Sprintf("%f", invoice.AmountToReceive)
		}else{
			returnMessage = "There is an invoice in the database containing " + invoice.InvoiceName + " and " + invoice.CompanyID
		}
		
	default:
	 	checkError(err)
	}
	return &Message{Body: returnMessage}, nil
}




func ModifyBankAccount(db *sql.DB, dni_investor string, toModify float64, action string)(result bool){
	//This function modifies the bank account quantity of the requested investor
	//If the action is add, we will add the specified quantity to the total amount, and remove it from the retained money, as it means that it was a retained bid that must be returned.
	//If the action is rm, we will remove the specified quantity to the total amount, and add it to the retained money, as it means that it was a bid that has been placed and money must be retained
	//If the action is removeRet, the quantity will be removed from the retained amount, as it means that bid has been placed, and this money should go for the seller
	// The function will return error, if the investor does not have enough money for the required operation
	
	var totalM, retainedM, newMoney, retMoney float64;

	//Check that there is no any invoice with that same name name and company:
	sqlStatement := `SELECT total_money, retained_money FROM investor WHERE dni=$1;`
	row := db.QueryRow(sqlStatement, dni_investor)
	err = row.Scan(&totalM, &retainedM)
	checkError(err)

	if action == "removeRet"{
		retMoney = retainedM - toModify
		updateStatement := `UPDATE investor SET retained_money=$1 WHERE dni=$2;`
		_, err := db.Exec(updateStatement, retMoney, dni_investor)
		checkError(err)
		result = true
	}else{
		if action == "add"{
			newMoney = totalM + toModify
			retMoney = retainedM - toModify
		}else if action == "rm"{
			newMoney = totalM - toModify
			retMoney = retainedM + toModify
		}
		if newMoney >= 0 {
			updateStatement := `UPDATE investor SET total_money=$1, retained_money=$2 WHERE dni=$3;`
			_, err := db.Exec(updateStatement,newMoney, retMoney, dni_investor)
			checkError(err)
			result = true
		}else{
			result = false
		}
	}
	return 
}

func (s *Server) TryToModifyInvestorMoney(ctx context.Context, bid *Bid) (*Message, error) {
	// This function receives a Bid object from grpc request
	// It will modify the bank account of the investor, removing the amount for the bid from the total money
	var returnMessage string

	db, err := OpenDB();
	defer CloseDB(db);
	checkError(err)

	res:=ModifyBankAccount(db, bid.Dni, bid.Amount, bid.Action)
	if res == true{
		returnMessage="yes"
	}else{
		returnMessage="no"
	}
	return &Message{Body: returnMessage}, nil
}






func (s *Server) TryToCreateInvoicePart(ctx context.Context, part *InvoicePart) (*Message, error) {
	var returnMessage string
	var total, toreceive float64
	var nameInv, issuerName string 


	db, err := OpenDB();
	defer CloseDB(db);
	checkError(err)
	
	//Get the details of the original Invoice from the database
	sqlStatement := `SELECT name, issuer, total, toreceive FROM invoice WHERE id=$1;`
	row := db.QueryRow(sqlStatement, part.OriginalId)
	err = row.Scan(&nameInv, &issuerName, &total, &toreceive)
  	checkError(err)

	// We will execute the matching algorithm
	// The result of the algorithm will be sent to the client
	returnMessage = matchingAlgorithm(db, part.Buyer, issuerName, part.OriginalId, total, toreceive, part.Total, part.Amount)
	return &Message{Body: returnMessage}, nil
}




func matchingAlgorithm(db *sql.DB, idinvestor string, idissuer string, idInvoice int32, totalInvoice float64, amountInvoice float64, totalBid float64, amountBid float64)(result string){
	//This function implements the matching algorithm. It will take the incoming parts and decice if:
	// Result codes:
		//0: Bids are placed for 100% of the invoice
		//1: Bid is rejected (discount is bigger) 
		//2: Invoice is not yet financed
	
	// In case that bid is rejected, (case 1), the user's bank acount will be released with the bid amount
	// In case that bids are placed but invoice is not yet fully financed (case 2) a new object of "part invoice" will be introduced in the database. This will contain all the info of the transaction.
	// In case that bids are placed for 100% of the invoice (case 0), new bids will be calculated, and a new object of "part invoice" will be introduced in the database.
	// In this last case,  the bid will be set as closed, and an object of type debt will be also registered in the database

	var discountInvoice float64
	var discountBid float64
	var sum float64
	var newTotalBid float64
	var newAmountBid float64
	var previousTotalBid float64

	discountInvoice = (1-(amountInvoice / totalInvoice))
	discountBid = (1-(amountBid / totalBid))

	if discountBid > discountInvoice{
		result = "1" //Means that bid has been rejected and money must be returned to the investor
		ModifyBankAccount(db, idinvestor, amountBid , "add")

	}else{		
		sqlStatement := `SELECT total FROM part_invoice WHERE invoice_id=$1;`
		rows, _:= db.Query(sqlStatement, idInvoice)
		defer rows.Close()
		sum = 0
		for rows.Next() {
			err := rows.Scan(&previousTotalBid)
			checkError(err)
			sum = sum + previousTotalBid
		}
		err := rows.Err();
		checkError(err)

		if sum + totalBid <= totalInvoice{
			insertStatement := `insert into "part_invoice"("invoice_id", "total", "amount", "buyer", "seller") values($1, $2, $3, $4, $5)`
			_, err := db.Exec(insertStatement, idInvoice, totalBid, amountBid, idinvestor, idissuer)
			checkError(err)
			result = "2"//Means that invoice is not yet financed
		}else{
			//Recalculate new bids
			newTotalBid = totalInvoice - sum
			newAmountBid = newTotalBid * (1-discountBid)

			insertStatement := `insert into "part_invoice"("invoice_id", "total", "amount", "buyer", "seller") values($1, $2, $3, $4, $5)`
			_, err := db.Exec(insertStatement, idInvoice, newTotalBid, newAmountBid, idinvestor, idissuer)
			checkError(err)
			ModifyBankAccount(db, idinvestor, amountBid - newAmountBid , "add")

			//Set bid as closed
			updateStatement := `UPDATE invoice SET closed=$1 WHERE id=$2;`
			_, err = db.Exec(updateStatement,"Yes", idInvoice)
			checkError(err)

			//remove money from retained and create debts in database and return 
			PayAndCreateDebts(db, idInvoice)
			result = "0" //Bids are placed for 100% of the invoice
		}
	}
	
	return result
}

func PayAndCreateDebts(db *sql.DB, idInvoice int32)(result bool){
	var total, amount float64
	var partid int32
	var buyer, seller string
	result = false


	//Get all invoice parts with ID of invoice
	sqlStatement := `SELECT invoice_part, total, amount, buyer, seller FROM part_invoice WHERE invoice_id=$1;`
	rows, err := db.Query(sqlStatement, idInvoice)
	checkError(err)
    defer rows.Close()

	//Iterate: for each invoice part we must remove the retained amount of money from account and create debt 
    for rows.Next() {
        err := rows.Scan(&partid, &total, &amount, &buyer,  &seller)
		checkError(err)
		ModifyBankAccount(db, buyer, amount , "removeRet")

		//Create debt
		insertStatement := `insert into "debt"("invoice_part", "amount", "creditor", "debtor") values($1, $2, $3, $4)`
		_, err = db.Exec(insertStatement, partid, total, buyer, seller)
		result = true
		checkError(err)
    }
	err = rows.Err()
	checkError(err)
	return 
}


func (s *Server) IntroduceIssuerToDatabase(ctx context.Context, message *Issuer) (*Message, error) {
	// This function introduces an Issuer in a database
	// It takes as parameter an Issuer grpc object
	db, err := OpenDB();
	defer CloseDB(db);
	checkError(err)
	
	insertStatement := `insert into "issuer"("nif", "name") values($1, $2)`
	_, err = db.Exec(insertStatement, message.Nif, message.Name)
	checkError(err)
	return &Message{Body: "Introduced in the database " + message.Name + " issuer with NIF " + message.Nif}, nil
}



func (s *Server) IntroduceInvestorToDatabase(ctx context.Context, message *Investor) (*Message, error) {
	// This function introduces a Investor in a database
	// It takes as parameter an Investor grpc object
	db, err := OpenDB();
	defer CloseDB(db);
	checkError(err)

	insertStatement := `insert into "investor"("dni", "name", "total_money", "retained_money") values($1, $2, $3, $4)`
	_, err = db.Exec(insertStatement, message.Dni, message.Name, message.AvailableMoney, message.RetainedMoney)
	checkError(err)
	return &Message{Body: "Introduced in the database " + message.Name + " Investor with DNI " + message.Dni}, nil
}


	

func (s *Server)GetAvailableInvoices(message *EmptyRequest, srv CommunicationService_GetAvailableInvoicesServer) error {
	//This function returns all the available Invoices in the Database as a stream
	//It does not take any parameter, and returns these values as grpc objects in a stream

	var invoiceissuer, invoicename, invoiceclosed string
	var invoicetotal, invoiceamount float64
	var invoiceid int32

	db, err := OpenDB();
	defer CloseDB(db);
	checkError(err)


	//Get all the available Invoices in the Database
	sqlStatement := `SELECT * FROM invoice;`
	rows, err := db.Query(sqlStatement)
    defer rows.Close()
	checkError(err)

    for rows.Next() {
        err := rows.Scan(&invoiceid, &invoiceissuer, &invoicename,  &invoicetotal,  &invoiceamount, &invoiceclosed)
		checkError(err)
		resp := Invoice{Id: invoiceid, CompanyID: invoiceissuer, InvoiceName:invoicename , TotalPrice: invoicetotal, AmountToReceive: invoiceamount, Closed: invoiceclosed }
		err = srv.Send(&resp);
		checkError(err)
    }
	err = rows.Err();
	checkError(err)
	return nil
}



                
