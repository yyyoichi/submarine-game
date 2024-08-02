
npm config set @buf:registry https://buf.build/gen/npm/v1 && \
npm install -g @buf/connectrpc_eliza.connectrpc_es @connectrpc/connect @connectrpc/connect-web @bufbuild/protoc-gen-es @connectrpc/protoc-gen-connect-es

go install github.com/bufbuild/buf/cmd/buf@latest && \
    go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest

cd web && npm ci

# deploy
sudo apt-get update
sudo apt-get install apt-transport-https ca-certificates gnupg curl
curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo gpg --dearmor -o /usr/share/keyrings/cloud.google.gpg
echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" | sudo tee -a /etc/apt/sources.list.d/google-cloud-sdk.list
sudo apt-get update && sudo apt-get install google-cloud-cli -y

