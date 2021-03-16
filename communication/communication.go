package communication

import (
	"log"

	"golang.org/x/net/context"
)

type Server struct {
}

func (s *Server) SellInvoice(ctx context.Context, invoice *Invoice) (*Message, error) {
	log.Printf("Seller %s wants to finance an invoice of price %f ", invoice.CompanyID)
	log.Printf("Received request from client: %s", invoice.CompanyID)
	return &Message{Body: "Response from the server!"}, nil
}
