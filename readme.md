
we can update the extension for a new version with
<https://developer.chrome.com/docs/extensions/how-to/distribute/host-on-linux#package_through_command_line>
```ps1
chrome.exe --pack-extension=./dist --pack-extension-key=<selfhost-crx-pack-ext-key.pem>
```
NOTE: make sure to version bump, otherwise this may do nothing?

# Azure Container Apps – End-to-End Deployment & Update Flow

## 1. Prerequisites
- Install Azure CLI (`winget install Microsoft.AzureCLI`)
- Install ACR & Container Apps extensions
  ```bash
  az extension add --name containerapp --upgrade
  az provider register --namespace Microsoft.ContainerRegistry    # Required once
  az provider register --namespace Microsoft.OperationalInsights --wait  # Required once
  ```

## 2. Local Build & Test
```bash
# Build and start the container locally
docker build -t selfhost-crx .
docker run -p 8080:8080 selfhost-crx:latest

# (WSL only) get URL-localhost
wsl -d ubuntu hostname -I   # yields <WSL_IP>
curl http://<WSL_IP>:8080/updates.xml          # manual smoke-test
```

## 3. Resource Group
Create once, reuse for everything.
```bash
az group create --name WebExtensionResourceGroup --location eastus
```

## 4. Azure Container Registry (ACR)
```bash
# Create registry
az acr create --resource-group WebExtensionResourceGroup --name hostcrx --sku Basic

# Quick login
az acr login --name hostcrx

# Alternative token login (if PowerShell/docker login prompt is restrictive)
ACR_NAME="hostcrx"
LOGIN_SERVER=$(az acr show --name $ACR_NAME --query loginServer -o tsv | tr -d '\r')
TOKEN=$(az acr login --name $ACR_NAME --expose-token | jq -r '.accessToken')
echo $TOKEN | docker login $LOGIN_SERVER -u 00000000-0000-0000-0000-000000000000 --password-stdin
```

## 5. Tag & Push Image to ACR
```bash
# Tag and push
docker tag selfhost-crx hostcrx.azurecr.io/selfhost-crx:latest
docker push hostcrx.azurecr.io/selfhost-crx:latest
```

## 6. Create Container Apps Environment (once)
```bash
az containerapp env create \
  --name WebExtensionContainerEnv \
  --resource-group WebExtensionResourceGroup \
  --location eastus
```

## 7. Deploy Initial Container App
```bash
az containerapp create \
  --name hostcrx \
  --resource-group WebExtensionResourceGroup \
  --environment WebExtensionContainerEnv \
  --image hostcrx.azurecr.io/selfhost-crx:latest \
  --target-port 8080 \
  --ingress external \
  --registry-server hostcrx.azurecr.io \
  --registry-identity system
```

Result: `https://hostcrx.victoriouswave-1bc4d533.eastus.azurecontainerapps.io/` is live.

---

## 8. Zero-Downtime Updates (rolling)
```bash
# 1. Re-build (or pull) latest image
docker build -t selfhost-crx .      # or use CI/CD pipeline

# 2. Tag & push
docker tag selfhost-crx hostcrx.azurecr.io/selfhost-crx:latest
docker push hostcrx.azurecr.io/selfhost-crx:latest

# 3. Trigger Azure rollout (zero downtime)
az containerapp update \
  --name hostcrx \
  --resource-group WebExtensionResourceGroup \
  --image hostcrx.azurecr.io/selfhost-crx:latest
```
