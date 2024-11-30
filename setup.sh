#!/bin/bash

# Variables
SUBSCRIPTION_ID="<Your-Subscription-ID>"
SP_NAME="custom-role-sp"
CUSTOM_ROLE_NAME="RetrieveBootDiagnosticsRole"

# Login to Azure and set the subscription
echo "Logging into Azure..."
az login
az account set --subscription "$SUBSCRIPTION_ID"

# Create a custom role definition JSON
CUSTOM_ROLE_FILE="custom-role-definition.json"
cat <<EOF > "$CUSTOM_ROLE_FILE"
{
  "Name": "$CUSTOM_ROLE_NAME",
  "IsCustom": true,
  "Description": "Custom role for retrieving boot diagnostics data",
  "Actions": [
    "Microsoft.Compute/virtualMachines/retrieveBootDiagnosticsData/action"
  ],
  "AssignableScopes": [
    "/subscriptions/$SUBSCRIPTION_ID"
  ]
}
EOF

# Create the custom role
echo "Creating custom role..."
az role definition create --role-definition "$CUSTOM_ROLE_FILE"

# Create the service principal
echo "Creating service principal..."
SP=$(az ad sp create-for-rbac --name "$SP_NAME" --role "$CUSTOM_ROLE_NAME" --scopes "/subscriptions/$SUBSCRIPTION_ID" --query '{appId:appId, tenantId:tenant, clientSecret:password}' -o json)

# Extract service principal details
APP_ID=$(echo "$SP" | jq -r '.appId')
TENANT_ID=$(echo "$SP" | jq -r '.tenantId')
CLIENT_SECRET=$(echo "$SP" | jq -r '.clientSecret')

# Output the service principal details
echo "Service Principal created successfully!"
echo "App ID: $APP_ID"
echo "Tenant ID: $TENANT_ID"
echo "Client Secret: $CLIENT_SECRET"
echo "Save these details securely. The client secret will not be retrievable again."

# Cleanup
rm -f "$CUSTOM_ROLE_FILE"
