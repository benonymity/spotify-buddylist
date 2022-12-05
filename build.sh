# /bin/sh
# Build Go binary:
cd api
CGO_ENABLED=1 go build server.go
# Build Vue JS:
cd ../frontend
npm i
npm run build
cd ..
