# Service for Profilbilder

## Forklaring

### Prereqs

Du trenger go, terraform og func (azure function cli) for å kjøre. I tillegg trenger du riktig environment variables fra portal.azure:
<code>
AZURE_STORAGE_ACCOUNT_NAME=<br/>
AZURE_STORAGE_ACCOUNT_KEY=<br/>
AZURE_STORAGE_CONTAINER_NAME=<br/>
</code>

### Images Handler

Dette er her handlers ligger. En fil for hver metode.

### Terraform

Her ligger Azure modulene:

- Blob Storage
- Function App
- Dependencies (resource group, etc., ...)

<p>Her er også <code>apply.sh</code> som applyer endringer.<span style="color: red;"> VÆR FORSIKTIG MED DENNE!</span> Første gang må man kjøre <code>terraform init</code> lokalt hos seg.</p>

### Test Opplasting m/ cURL

<code>POST</code> bilde lokalt til Azure Storage.

```
curl -X POST -F "image=@/path/til/bilde.png" https://torger-function-app.azurewebsites.net/api/images
```

## TODO

- API key / auth
- Vurdere hvilke filer som skal lagres i containeren. Vi kan evt. ha flere containers, 1 for profilblider og andre for annet. I så fall handle hvilke type filer som er tillatt.
- <code>/DELETE</code> endepunkt
