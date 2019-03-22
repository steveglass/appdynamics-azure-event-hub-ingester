package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/Azure/azure-amqp-common-go/aad"
	eventhub "github.com/Azure/azure-event-hubs-go"
	mgmt "github.com/Azure/azure-sdk-for-go/services/eventhub/mgmt/2017-04-01/eventhub"
	"github.com/Azure/go-autorest/autorest/azure"
	azauth "github.com/Azure/go-autorest/autorest/azure/auth"
)

const (
	schemaURL   = "/events/schema/"
	publishURL  = "/events/publish/"
	typeGateway = 0 // Used for schema definitions. GATEWAY is APIM
	typeCosmos  = 1 // COSMOSDB

)

// poor man's debug
var debug = false

var wg sync.WaitGroup

// TODO - Need proper logging

func main() {
	hubCount := 0

	// Read config file for AppDynamics + Azure info
	conf := initConfig()
	if conf.AzureCosmosHubName != "nil" {
		hubCount++
		// Checking to see if AppD analytics schema exists yet
		exists := checkSchema(conf.AnalyticsCosmosSchema, conf)
		if exists {
			fmt.Println("Analytics schema for CosmosDB exists")
		} else {
			fmt.Println("Analytics schema for CosmosDB does not exist. Creating Now")
			createSchema(conf.AnalyticsCosmosSchema, conf, typeCosmos)
		}
		go startListener(conf, conf.AzureCosmosHubName, typeCosmos)
	}
	if conf.AzureGatewayHubName != "nil" {
		hubCount++
		// Checking to see if AppD analytics schema exists yet
		exists := checkSchema(conf.AnalyticsGatewaySchema, conf)
		if exists {
			fmt.Println("Analytics schema for Gateway (APIM) exists")
		} else {
			fmt.Println("Analytics schema for Gateway (APIM) does not exist. Creating Now")
			createSchema(conf.AnalyticsGatewaySchema, conf, typeGateway)
		}
		go startListener(conf, conf.AzureGatewayHubName, typeGateway)
	}

	wg.Add(hubCount)

	wg.Wait()
}

func startListener(conf appdConfig, name string, deftype int) {
	exit := make(chan struct{})
	var num int
	hub, partitions := initHub(conf, name)

	handler := func(ctx context.Context, event *eventhub.Event) error {
		switch deftype {
		case typeGateway:
			num = serializeGatewayRecord(event.Data, conf)
		case typeCosmos:
			num = serializeCosmosRecord(event.Data, conf)
		default:
			fmt.Printf("Unkown type in event handler - shouldn't actually be able to get here. [type] %d\n", deftype)
		}
		fmt.Printf("Number of Records Sent to Analytics(%s): %d\n", name, num)
		return nil
	}

	fmt.Printf("Initializing Event Hub Listener for %s with callback...\n", name)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	for _, partitionID := range partitions {
		hub.Receive(ctx, partitionID, handler, eventhub.ReceiveWithLatestOffset())
	}
	cancel()

	fmt.Printf("Initialization Complete. Listening to Event Hub %s\n", name)

	select {
	case <-exit:
		fmt.Println("closing after 2 seconds")
		select {
		case <-time.After(2 * time.Second):
			defer wg.Done()
			return
		}
	}
}

func initHub(conf appdConfig, name string) (*eventhub.Hub, []string) {
	// Set env var used by Azure libs
	os.Setenv("AZURE_CLIENT_ID", conf.AzureClientID)
	os.Setenv("AZURE_CLIENT_SECRET", conf.AzureClientSecret)
	os.Setenv("AZURE_TENANT_ID", conf.AzureTenantID)
	os.Setenv("AZURE_SUBSCRIPTION_ID", conf.AzureSubscriptionID)
	os.Setenv("EVENTHUB_NAMESPACE", conf.AzureEventHubNameSpace)
	os.Setenv("EVENTHUB_CONNECTION_STRING", conf.AzureEventHubConnString)
	namespace := mustGetenv("EVENTHUB_NAMESPACE")
	hubMgmt, err := ensureEventHub(context.Background(), name, conf)
	if err != nil {
		log.Fatal(err)
	}

	provider, err := aad.NewJWTProvider(aad.JWTProviderWithEnvironmentVars())
	if err != nil {
		log.Fatal(err)
	}
	hub, err := eventhub.NewHub(namespace, name, provider)
	if err != nil {
		panic(err)
	}
	return hub, *hubMgmt.PartitionIds
}

func mustGetenv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic("Environment variable '" + key + "' required for integration tests.")
	}
	return v
}

func ensureEventHub(ctx context.Context, name string, conf appdConfig) (*mgmt.Model, error) {
	namespace := mustGetenv("EVENTHUB_NAMESPACE")
	client := getEventHubMgmtClient()
	hub, err := client.Get(ctx, conf.AzureResourceGroup, namespace, name)

	partitionCount := int64(4)
	if err != nil {
		newHub := &mgmt.Model{
			Name: &name,
			Properties: &mgmt.Properties{
				PartitionCount: &partitionCount,
			},
		}

		hub, err = client.CreateOrUpdate(ctx, conf.AzureResourceGroup, namespace, name, *newHub)
		if err != nil {
			return nil, err
		}
	}
	return &hub, nil
}

func getEventHubMgmtClient() *mgmt.EventHubsClient {
	subID := mustGetenv("AZURE_SUBSCRIPTION_ID")
	client := mgmt.NewEventHubsClientWithBaseURI(azure.PublicCloud.ResourceManagerEndpoint, subID)
	a, err := azauth.NewAuthorizerFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	client.Authorizer = a
	return &client
}
