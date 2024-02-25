# Mattermost HTML2MD

## Usage

0. Up infrastructure `make infrastructure.local.up`
1. Go to http://localhost:8065 create access token and put it to all `.env`.
2. Run tests `go test ./... -v -coverprofile=coverage.txt -covermode atomic -timeout 30s -run ^TestMain$ mattermost-html2md/tests`
3. Run service `make app.run`
4. Do request
```
curl --location 'http://localhost:8080/api/v1/webhook' \
--header 'X-API-KEY: test' \
--header 'Content-Type: application/json' \
--data '{
    "text": "<h1>Hello World</h1><p>This is a simple HTML document.</p>"
    "channel": "pyx1obq8e7ympkm4eitq3uq89c"
}'
```
5. See result `204 No Content`
6. Go to http://localhost:8065 and see result in channel
7. Go to OpanAPI docs http://localhost:8080/docs or http://localhost:8080/redoc
