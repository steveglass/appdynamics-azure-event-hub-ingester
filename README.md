# AppDynamics Azure EventHub Ingestion Service

This is a standalone application, no AppDynamics Machine Agent is required to run this service. Transaction analytics is required in AppDynamics, as is EventHub diagnostic monitoring for supported services. 

## Purpose

The purpose of this application is to ingest data from the Azure EventHub and push it to AppDynamics analytics so that users of both AppDynamics and Azure can further correlate application performance with the health of their Azure platform services where agents cannot be deployed. At the moment, the service only supports Azure API Management (APIM) and CosmosDB. 

## Solution, Pre-requisites, and Builds

The service is written in Go and utlizes a handful of third-party integrations from Microsoft and Google. To setup a Go development environment to download and build the service, run the following commands to get the dependencies:

```
  go get -u github.com/Azure/azure-event-hubs-go/...
  go get -u github.com/Azure/azure-amqp-common-go/...
  go get -u github.com/Azure/go-autorest/...
  go get  -u gopkg.in/yaml.v2/...
```

Included in this repository is a simple MakeFile for building an executable as well as a Dockerfile for building the service into a container. To properly build a Docker container, you must download the source, configure the config.yml file, and then build the container. The MakeFile will also build the container by uncommenting the following configurations:

```
 image:
   docker build -t appdynamics-azure-event-hub .
```

## Configuration

To configure the AppDynamics portion of the configuration, start by configuring the AppDynamics analytics information:
```
# AppD analytics endpoint
analyticsEndPoint: 
# AppD global account name on license page
globalAccountName: 
# unique key from analytics
analyticsKey: 
```

The analytics key is configured in the Analytics section of the AppDynamics UI. For more information, please refer to the documentation: https://docs.appdynamics.com/display/PRO45/Managing+API+Keys. For this service, you will need to give the key permissions to "Custom Analytics Events" for Manage, Query, and Publish. 

The next configurations are used to name the custom schemas that will appear in AppDynamics. If you do not intend to collect data for one of the Azure services listed, you can leave it blank. If you do conifgure a corresponding Hub, this cannot be blank

```
# this will be the name of the custom schema for APIM in AppD
analyticsGatewaySchema: 
# this will be the name of the custom schema for CosmosDB in AppD
analyticsCosmosSchema:
```

Next is the Azure conifguration for connecting to your subscription and service principal. Follow the instruction here if you have not done so already - [Azure - Resource Manager - Howto - Control Access - Create Service Principal - Azure Portal](https://docs.microsoft.com/en-us/azure/azure-resource-manager/resource-group-create-service-principal-portal)

```
# subscriptionId is the Azure Subscription ID 
azureSubscriptionId: 
# Application ID for your application via Access Control (IAM)
azureClientID: 
# client secret is available in the keys section for app registrations
azureClientSecret: 
# tenantId is the Azure Active Directory “Directory ID”
azureTenantID: 
```

Finally the EventHub information for each service to pull data from.

```
# EventHub namespace
azureEventHubNameSpace: 
# EventHub conneciton string
azureEventHubConnString: 
# Name of the resource group
azureResourceGroup: 
# if you're not monitoring APIM, set this to nil
azureGatewayHubName: nil
# if you're not monitoring CosmosDB, set this to nil
azureCosmosHubName: nil
```

## Current Service Support

Today the AppDynamics Azure Event Hub Ingester support APIM and CosmosDB. 

### Azure API Management (APIM or Gateway)

| Metric            |   Data type   |
| ----------------- |-------------- |
|apiId	            |   string      |
|apimSubscription	|   string      |
|apiRevision	    |   sring       |
|apiRevision	    |   string      |
|backendMethod	    |   string      |
|backendProtocol    |	string      |
|backendResponseCode|   int         |
|backendTime        |	string      |
|backendUrl 	    |   string      |
|cache	            |   string      |
|callerIpAddress    |	string      |
|category           |	string      |
|clientProtocol     |	string      |
|correlationId      |	string      |
|durationMs         |	int         |
|isRequestSuccess   |	string      |
|level	            |   int         |
|location	        |   string      |
|method	            |   string      |
|message            |	string      |
|operationId        |	string      |
|operationName      |	string      |
|requestSize        |	int         |
|reason	            |   string      |
|resonseSize        |	int         |
|resourceId	        |   string      |
|responseCode	    |   int         |
|section	        |   string      |
|source	            |   string      |
|time	            |   string      |
|url	            |   string      |

### CosmosDB

| Metric            |   Data type   |
| ----------------- |-------------- |
|time	            |   string      |
|resourceId         |	string      |
|category           |	string      |
|operationName	    |   string      |
|activityId         |	string      |
|opCode	            |   string      |
|errorCode	        |   string      |
|duration	        |   int         |
|requestCharge	    |   float       |
|databaseName	    |   string      |
|collectionName	    |   string      |
|retryCount	        |   string      |
