package worker

import (
	"context"
	"fmt"
	"log"

	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
)

type ZeebeClient struct {
	client zbc.Client
}

func NewZeebeClient(address string) (*ZeebeClient, error) {
	client, err := zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress:         address,
		UsePlaintextConnection: true, // For local dev/Arda infra
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Zeebe client: %w", err)
	}
	return &ZeebeClient{client: client}, nil
}

func (z *ZeebeClient) StartInstance(ctx context.Context, processID string, variables map[string]interface{}) (int64, error) {
	request, err := z.client.NewCreateInstanceCommand().
		BPMNProcessId(processID).
		LatestVersion().
		VariablesFromMap(variables)
	if err != nil {
		return 0, err
	}

	result, err := request.Send(ctx)
	if err != nil {
		log.Printf("Zeebe: Error starting instance %s: %v\n", processID, err)
		return 0, err
	}

	log.Printf("Zeebe: Started instance %d for process %s\n", result.ProcessInstanceKey, processID)
	return result.ProcessInstanceKey, nil
}
