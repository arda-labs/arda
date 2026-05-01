package worker

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	v1 "github.com/arda-labs/arda/apps/backend-go/bpm-service/api/crm/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type EventPublisher interface {
	Publish(ctx context.Context, topic string, event interface{}) error
}

type GenericWorker struct {
	httpClient *http.Client
	publisher  EventPublisher
	crmClient  v1.CRMClient
}

func NewGenericWorker(pub EventPublisher) *GenericWorker {
	// CRM_GRPC_ADDR defaults to thinkcenter:9010 if not set
	crmAddr := os.Getenv("CRM_GRPC_ADDR")
	if crmAddr == "" {
		crmAddr = "thinkcenter:9010"
	}

	conn, err := grpc.Dial(crmAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to CRM service at %s: %v", crmAddr, err)
	}

	return &GenericWorker{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		publisher:  pub,
		crmClient:  v1.NewCRMClient(conn),
	}
}

func (w *GenericWorker) HandleServiceTask(ctx context.Context, taskVariables map[string]interface{}) (map[string]interface{}, error) {
	protocol := taskVariables["protocol"] // "HTTP" or "GRPC"

	if protocol == "GRPC" {
		return w.handleGRPC(ctx, taskVariables)
	}
	return w.handleHTTP(ctx, taskVariables)
}

func (w *GenericWorker) handleHTTP(ctx context.Context, vars map[string]interface{}) (map[string]interface{}, error) {
	log.Println("BPM Generic Worker: Handling HTTP request...")
	// Existing HTTP logic...
	return map[string]interface{}{"status": 200}, nil
}

func (w *GenericWorker) handleGRPC(ctx context.Context, vars map[string]interface{}) (map[string]interface{}, error) {
	method := fmt.Sprintf("%v", vars["method"])
	log.Printf("BPM Generic Worker: Calling gRPC method %s\n", method)

	if method == "FinalizeCustomer" {
		customerID := fmt.Sprintf("%v", vars["customerId"])
		status := fmt.Sprintf("%v", vars["status"])

		reply, err := w.crmClient.FinalizeCustomer(ctx, &v1.FinalizeCustomerRequest{
			CustomerId: customerID,
			Status:     status,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to finalize customer: %w", err)
		}

		return map[string]interface{}{
			"success": reply.Success,
			"message": reply.Message,
		}, nil
	}

	return nil, fmt.Errorf("unknown gRPC method: %s", method)
}
