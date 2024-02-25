# Mattermost HTML2MD

## Usage

0. Up infrastructure `make infrastructure.local.up`
1. Go to http://localhost:8085 create incoming webhook and put url to all `.env`.
2. Run tests `go test -v -coverprofile=coverage.txt -covermode atomic -timeout 30s -run ^TestMain$ mattermost-html2md/tests`
3. Run service `go run cmd/server/cmd/main.go`
4. Do request
```
curl --location 'http://localhost:8080/send' \
--header 'X-API-KEY: test' \
--header 'Content-Type: application/json' \
--data '{
    "text": "<h1>Hello World</h1><p>This is a simple HTML document.</p>"
}'
```
5. See result `204 No Content`
6. Go to http://localhost:8085 and see result in channel
