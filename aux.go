package main

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func AuthenticateClient() (*azidentity.ClientSecretCredential, error) {
	var creds ClientSecretInfo
	f, err := os.Open("./appcreds.json")
	if err != nil {
		log.Fatalf("failed to open file %s\n", err)
		return nil, err
	}
	data, err := io.ReadAll(f)
	if err != nil {
		log.Fatalf("failed to read file %s\n", err)
		return nil, err
	}
	err = json.Unmarshal(data, &creds)
	if err != nil {
		log.Fatal("failed to unmarshal data %s\n", err)
		return nil, err
	}
	cred, err := azidentity.NewClientSecretCredential(creds.TenantID, creds.ClientID, creds.ClientSecret, nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
		return nil, err
	}
	return cred, nil
}
