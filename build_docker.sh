# /bin/sh
# Build Go binary:
cd api_go
CGO_ENABLED=0 go build server.go
# Build Vue JS:
cd ../frontend
npm i
npm run build
cd ..
