package worker

import (
	"context"
	"log"
)

type CustomerWorker struct {
	// Add dependencies like Zeebe client or CRM gRPC client here
}

func NewCustomerWorker() *CustomerWorker {
	return &CustomerWorker{}
}

// HandleRegisterCustomer logic for automatic registration step
func (w *CustomerWorker) HandleRegisterCustomer(ctx context.Context, job interface{}) error {
	log.Println("BPM Worker: Handling Register Customer step...")
	// Logic: Call CRM Service via gRPC to create/activate customer
	return nil
}

// HandleAdjustCustomer logic for adjustment step
func (w *CustomerWorker) HandleAdjustCustomer(ctx context.Context, job interface{}) error {
	log.Println("BPM Worker: Handling Adjust Customer step...")
	return nil
}
