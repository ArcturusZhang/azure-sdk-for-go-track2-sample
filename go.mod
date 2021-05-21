module github.com/ArcturusZhang/track2-test-program

go 1.16

require (
	github.com/Azure/azure-sdk-for-go/sdk/armcore v0.7.1
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.16.1
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v0.8.0
	github.com/Azure/azure-sdk-for-go/sdk/compute/armcompute v0.1.0
	github.com/Azure/azure-sdk-for-go/sdk/network/armnetwork v0.1.0
	github.com/Azure/azure-sdk-for-go/sdk/resources/armresources v0.1.0
	github.com/Azure/azure-sdk-for-go/sdk/to v0.1.4
)

replace github.com/Azure/azure-sdk-for-go/sdk/compute/armcompute v0.1.0 => github.com/ArcturusZhang/azure-sdk-for-go/sdk/compute/armcompute v0.0.0-20210521061855-4bffc32e2ffa

replace github.com/Azure/azure-sdk-for-go/sdk/resources/armresources v0.1.0 => github.com/ArcturusZhang/azure-sdk-for-go/sdk/resources/armresources v0.0.0-20210521021048-736bb0adf0f3

replace github.com/Azure/azure-sdk-for-go/sdk/network/armnetwork v0.1.0 => github.com/ArcturusZhang/azure-sdk-for-go/sdk/network/armnetwork v0.0.0-20210521063755-e6e567d04b64
