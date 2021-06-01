package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/network/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

const (
	interval = 10 * time.Second
)

var (
	ctx               context.Context
	subscriptionId    string
	location          = "westus2"
	resourceGroupName = "dapzhang-track2"
	vnetName          = "dapzhang-vnet"
	subnetName        = "internal"
	nicName           = "dapzhang-nic"
	vmName            = "dapzhang-vm"

	resourceGroupID string
	vnetID          string
	subnetID        string
	nicID           string
	vmID            string
)

func init() {
	ctx = context.Background()
	subscriptionId = os.Getenv("SUBSCRIPTION_ID")
}

func main() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}
	conn := armcore.NewDefaultConnection(cred, &armcore.ConnectionOptions{
		Logging: azcore.LogOptions{
			IncludeBody: true,
		},
	})

	defer cleanup(conn)

	if err := createResourceGroup(conn); err != nil {
		panic(err)
	}
	addCleanupFunction(deleteResourceGroup)

	if err := createVirtualNetwork(conn); err != nil {
		panic(err)
	}
	addCleanupFunction(deleteVirtualNetwork)

	if err := createSubnet(conn); err != nil {
		panic(err)
	}
	addCleanupFunction(deleteSubnet)

	if err := createNIC(conn); err != nil {
		panic(err)
	}
	addCleanupFunction(deleteNIC)

	if err := createVirtualMachine(conn); err != nil {
		panic(err)
	}
	addCleanupFunction(deleteVirtualMachine)
}

var cleanupFuncs []cleanupFunc

type cleanupFunc func(connection *armcore.Connection) error

func addCleanupFunction(f cleanupFunc) {
	cleanupFuncs = append(cleanupFuncs, f)
}

func cleanup(connection *armcore.Connection) {
	for i := len(cleanupFuncs) - 1; i >= 0; i-- {
		f := cleanupFuncs[i]
		_ = f(connection)
	}
}

func createResourceGroup(connection *armcore.Connection) error {
	rgClient := armresources.NewResourceGroupsClient(connection, subscriptionId)

	param := armresources.ResourceGroup{
		Location: to.StringPtr(location),
	}

	resp, err := rgClient.CreateOrUpdate(ctx, resourceGroupName, param, nil)
	if err != nil {
		return err
	}
	b, err := json.MarshalIndent(*resp.ResourceGroup, "", "  ")
	if err != nil {
		return err
	}

	resourceGroupID = *resp.ResourceGroup.ID
	fmt.Printf("Resource Group '%s' created: \n%s\n", resourceGroupID, string(b))
	return nil
}

func deleteResourceGroup(connection *armcore.Connection) error {
	rgClient := armresources.NewResourceGroupsClient(connection, subscriptionId)

	poller, err := rgClient.BeginDelete(ctx, resourceGroupName, nil)
	if err != nil {
		return err
	}
	if _, err := poller.PollUntilDone(ctx, interval); err != nil {
		return err
	}
	fmt.Printf("Resource Group '%s' deleted.\n", resourceGroupID)
	return nil
}

func createVirtualNetwork(connection *armcore.Connection) error {
	vnetClient := armnetwork.NewVirtualNetworksClient(connection, subscriptionId)

	param := armnetwork.VirtualNetwork{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(location),
		},
		Properties: &armnetwork.VirtualNetworkPropertiesFormat{
			AddressSpace: &armnetwork.AddressSpace{
				AddressPrefixes: []*string{
					to.StringPtr("10.0.0.0/16"),
				},
			},
		},
	}
	poller, err := vnetClient.BeginCreateOrUpdate(ctx, resourceGroupName, vnetName, param, nil)
	if err != nil {
		return err
	}
	resp, err := poller.PollUntilDone(ctx, interval)
	if err != nil {
		return err
	}
	b, err := json.MarshalIndent(*resp.VirtualNetwork, "", "  ")
	if err != nil {
		return err
	}

	vnetID = *resp.VirtualNetwork.ID
	fmt.Printf("Virtual Network '%s' created: \n%s\n", vnetID, string(b))
	return nil
}

func deleteVirtualNetwork(connection *armcore.Connection) error {
	vnetClient := armnetwork.NewVirtualNetworksClient(connection, subscriptionId)

	poller, err := vnetClient.BeginDelete(ctx, resourceGroupName, vnetName, nil)
	if err != nil {
		return err
	}
	if _, err := poller.PollUntilDone(ctx, interval); err != nil {
		return err
	}
	fmt.Printf("Virtual Network '%s' deleted.\n", vnetID)
	return nil
}

func createSubnet(connection *armcore.Connection) error {
	subnetClient := armnetwork.NewSubnetsClient(connection, subscriptionId)

	param := armnetwork.Subnet{
		Properties: &armnetwork.SubnetPropertiesFormat{
			AddressPrefix: to.StringPtr("10.0.2.0/24"),
		},
	}
	poller, err := subnetClient.BeginCreateOrUpdate(ctx, resourceGroupName, vnetName, subnetName, param, nil)
	if err != nil {
		return err
	}
	resp, err := poller.PollUntilDone(ctx, interval)
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(*resp.Subnet, "", "  ")
	if err != nil {
		return err
	}

	subnetID = *resp.Subnet.ID
	fmt.Printf("Subnet '%s' created: \n%s\n", subnetID, string(b))
	return nil
}

func deleteSubnet(connection *armcore.Connection) error {
	subnetClient := armnetwork.NewSubnetsClient(connection, subscriptionId)

	poller, err := subnetClient.BeginDelete(ctx, resourceGroupName, vnetName, subnetName, nil)
	if err != nil {
		return err
	}
	if _, err := poller.PollUntilDone(ctx, interval); err != nil {
		return err
	}
	fmt.Printf("Subnet '%s' deleted.\n", subnetID)
	return nil
}

func createNIC(connection *armcore.Connection) error {
	nicClient := armnetwork.NewNetworkInterfacesClient(connection, subscriptionId)

	param := armnetwork.NetworkInterface{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(location),
		},
		Properties: &armnetwork.NetworkInterfacePropertiesFormat{
			IPConfigurations: []*armnetwork.NetworkInterfaceIPConfiguration{
				{
					Name: to.StringPtr("internal"),
					Properties: &armnetwork.NetworkInterfaceIPConfigurationPropertiesFormat{
						PrivateIPAllocationMethod: armnetwork.IPAllocationMethodDynamic.ToPtr(),
						Subnet: &armnetwork.Subnet{
							SubResource: armnetwork.SubResource{
								ID: to.StringPtr(subnetID),
							},
						},
					},
				},
			},
		},
	}
	poller, err := nicClient.BeginCreateOrUpdate(ctx, resourceGroupName, nicName, param, nil)
	if err != nil {
		return err
	}
	resp, err := poller.PollUntilDone(ctx, interval)
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(*resp.NetworkInterface, "", "  ")
	if err != nil {
		return err
	}

	nicID = *resp.NetworkInterface.ID
	fmt.Printf("Network Interface '%s' created: \n%s\n", nicID, string(b))
	return nil
}

func deleteNIC(connection *armcore.Connection) error {
	nicClient := armnetwork.NewNetworkInterfacesClient(connection, subscriptionId)

	poller, err := nicClient.BeginDelete(ctx, resourceGroupName, nicName, nil)
	if err != nil {
		return err
	}
	if _, err := poller.PollUntilDone(ctx, interval); err != nil {
		return err
	}

	fmt.Printf("NIC '%s' deleted.\n", nicID)
	return nil
}

func createVirtualMachine(connection *armcore.Connection) error {
	vmClient := armcompute.NewVirtualMachinesClient(connection, subscriptionId)

	param := armcompute.VirtualMachine{
		Resource: armcompute.Resource{
			Location: to.StringPtr(location),
		},
		Identity: &armcompute.VirtualMachineIdentity{
			Type: armcompute.ResourceIdentityTypeSystemAssigned.ToPtr(),
		},
		Properties: &armcompute.VirtualMachineProperties{
			HardwareProfile: &armcompute.HardwareProfile{
				VMSize: armcompute.VirtualMachineSizeTypesStandardF2.ToPtr(),
			},
			OSProfile: &armcompute.OSProfile{
				AdminPassword:        to.StringPtr("P@$$w0rd1234!"),
				AdminUsername:        to.StringPtr("adminuser"),
				ComputerName:         to.StringPtr("arcturus"),
				WindowsConfiguration: &armcompute.WindowsConfiguration{},
			},
			NetworkProfile: &armcompute.NetworkProfile{
				NetworkInterfaces: []*armcompute.NetworkInterfaceReference{
					{
						SubResource: armcompute.SubResource{
							ID: to.StringPtr(nicID),
						},
					},
				},
			},
			StorageProfile: &armcompute.StorageProfile{
				ImageReference: &armcompute.ImageReference{
					Offer:     to.StringPtr("WindowsServer"),
					Publisher: to.StringPtr("MicrosoftWindowsServer"),
					SKU:       to.StringPtr("2016-Datacenter"),
					Version:   to.StringPtr("latest"),
				},
				OSDisk: &armcompute.OSDisk{
					CreateOption: armcompute.DiskCreateOptionTypesFromImage.ToPtr(),
					Caching:      armcompute.CachingTypesReadWrite.ToPtr(),
					ManagedDisk: &armcompute.ManagedDiskParameters{
						StorageAccountType: armcompute.StorageAccountTypesStandardLRS.ToPtr(),
					},
					OSType: armcompute.OperatingSystemTypesWindows.ToPtr(),
				},
			},
		},
	}

	poller, err := vmClient.BeginCreateOrUpdate(ctx, resourceGroupName, vmName, param, nil)
	if err != nil {
		return err
	}

	// we cannot use the resp returned by the service because this response does not returned with a final polling URL in its header
	if _, err := poller.PollUntilDone(ctx, interval); err != nil {
		return err
	}

	resp, err := vmClient.Get(ctx, resourceGroupName, vmName, nil)
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(*resp.VirtualMachine, "", "  ")
	if err != nil {
		return err
	}

	vmID = *resp.VirtualMachine.ID
	fmt.Printf("Virtual Machine '%s' created: \n%s\n", vmID, string(b))
	return nil
}

func deleteVirtualMachine(connection *armcore.Connection) error {
	vmClient := armcompute.NewVirtualMachinesClient(connection, subscriptionId)

	poller, err := vmClient.BeginDelete(ctx, resourceGroupName, vmName, &armcompute.VirtualMachinesBeginDeleteOptions{
		ForceDeletion: to.BoolPtr(true),
	})
	if err != nil {
		return err
	}
	if _, err := poller.PollUntilDone(ctx, interval); err != nil {
		return err
	}

	fmt.Printf("Virtual Machine '%s' deleted.\n", vmID)
	return nil
}
