1. copy dist.crx into this directory

2. manually update version in main.go to target version
```ps1
rg --ignore-case "version" main.go ../gaafl/gaafl.json
```

```lua
print(vim.fn.system({ 'rg', '--ignore-case', 'version', 'main.go', '../gaafl/gaafl.json' }))
```

3. run go vet just in case

```ps1
git add dist.crx main.go
git commit -m 'version bump to VERSION'
```



... upload using az (wsl)...
1. install docker & azure-cli
<https://docs.docker.com/engine/install/ubuntu/>
<https://learn.microsoft.com/en-us/cli/azure/install-azure-cli-linux?view=azure-cli-latest&pivots=apt>



2. authenticate with azure-cli (DefaultSubscription)
```bash
az acr login --name hostcrx
az acr update -n hostcrx

# Alternative token login (if PowerShell/docker login prompt is restrictive)
ACR_NAME="hostcrx"
LOGIN_SERVER=$(az acr show --name $ACR_NAME --query loginServer -o tsv | tr -d '\r')
TOKEN=$(az acr login --name $ACR_NAME --expose-token | jq -r '.accessToken')
echo $TOKEN | docker login $LOGIN_SERVER -u 00000000-0000-0000-0000-000000000000 --password-stdin
```
?

test with

```sh
docker run selfhost-crx:latest --network host
```


4.
```bash
# 1. Re-build (or pull) latest image
docker build -t selfhost-crx . --build-arg APP_VERSION=v1.0.4 # --no-cache

# 2. Tag & push
docker tag selfhost-crx hostcrx.azurecr.io/selfhost-crx:latest
docker push hostcrx.azurecr.io/selfhost-crx:latest

# 3. Trigger Azure rollout (zero downtime)
az containerapp update \
  --name hostcrx \
  --resource-group WebExtensionResourceGroup \
  --image hostcrx.azurecr.io/selfhost-crx:latest
```
