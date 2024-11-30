package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6"

	"github.com/google/uuid"

	"github.com/fatih/color"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s <resource group> <vm name>\n", os.Args[0])
		os.Exit(-1)
	}

	GetBootDiag(os.Args[1], os.Args[2])
}

func GetBootDiag(rg string, vmName string) {

	cred, err := AuthenticateClient()

	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory(os.Getenv("AZURE_SUBSCRIPTION_ID"), cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewVirtualMachinesClient().RetrieveBootDiagnosticsData(ctx, rg, vmName, &armcompute.VirtualMachinesClientRetrieveBootDiagnosticsDataOptions{SasURIExpirationTimeInMinutes: to.Ptr[int32](60)})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}

	screenshotResponse := res.RetrieveBootDiagnosticsDataResult.ConsoleScreenshotBlobURI

	resp, err := http.Get(*screenshotResponse)
	if err != nil {
		log.Fatalf("failed to get screenshot %s\n", err)
	}
	defer resp.Body.Close()

	uuid, err := uuid.NewUUID()
	filename := fmt.Sprintf("screenshots/screenshot_%s.jpg", uuid.String())
	screenshot_file, err := os.Create(filename)
	buffer, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("failed to get response body %s\n", err)
	}
	_, err = screenshot_file.Write(buffer)
	if err != nil {
		log.Fatalf("failed to save screenshot %s\n", err)
	}

	color.Blue("[*] saved screenshot file: %s\n", color.RedString(strings.Split(filename, "/")[1]))

}
