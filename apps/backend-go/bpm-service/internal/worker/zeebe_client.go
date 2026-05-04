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

// DeployBPMN deploys a BPMN XML to Zeebe and returns the deployment key.
func (z *ZeebeClient) DeployBPMN(ctx context.Context, bpmnXML, resourceName string) (int64, error) {
	cmd, err := z.client.NewDeployResourceCommand().
		AddResource([]byte(bpmnXML), resourceName).
		Send(ctx)
	if err != nil {
		return 0, fmt.Errorf("zeebe deploy: %w", err)
	}

	// The deployment key is embedded in the response - use the key from the first deployment
	key := cmd.GetKey()
	log.Printf("Zeebe: Deployed process %s, key: %d\n", resourceName, key)
	return key, nil
}

