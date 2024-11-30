package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s <resource group> <vm name>\n", os.Args[0])
		os.Exit(-1)
	}

	GetBootDiag(os.Args[1], os.Args[2])
}

type ClientSecretInfo struct {
	ClientID     string `json:"ClientID"`
	TenantID     string `json:"TenantID"`
	ClientSecret string `json:"ClientSecret`
}

func GetBootDiag(rg string, vmName string) {

	var creds ClientSecretInfo
	f, err := os.Open("./appcreds.json")
	if err != nil {
		log.Fatalf("failed to open file %s\n", err)
	}
	data, err := io.ReadAll(f)
	if err != nil {
		log.Fatalf("failed to read file %s\n", err)
	}
	err = json.Unmarshal(data, &creds)
	if err != nil {
		log.Fatal("failed to unmarshal data %s\n", err)
	}
	cred, err := azidentity.NewClientSecretCredential(creds.TenantID, creds.ClientID, creds.ClientSecret, nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory(os.Getenv("AZURE_SUBSCRIPTION_ID"), cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewVirtualMachinesClient().RetrieveBootDiagnosticsData(ctx, rg, vmName, &armcompute.VirtualMachinesClientRetrieveBootDiagnosticsDataOptions{SasURIExpirationTimeInMinutes: to.Ptr[int32](60)})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}

	screenshot := res.RetrieveBootDiagnosticsDataResult.ConsoleScreenshotBlobURI

	resp, err := http.Get(*screenshot)
	if err != nil {
		log.Fatalf("failed to get screenshot %s\n", err)
	}
	log.Printf("resp code: %s\n", resp.StatusCode)
	defer resp.Body.Close()

	screenshot_file, err := os.Create("screenshot.jpg")
	buffer, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("failed to get response body %s\n", err)
	}
	n, err := screenshot_file.Write(buffer)
	if err != nil {
		log.Fatalf("failed to save screenshot %s\n", err)
	}
	log.Printf("saved screenshot file with size: %d bytes\n", n)

}
