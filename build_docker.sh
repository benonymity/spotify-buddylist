# /bin/sh
# Build Go binary:
cd activity_api_go
CGO_ENABLED=0 go build server.go
# Build Vue JS:
cd ../activity_frontend
npm i
npm run build
cd ..
