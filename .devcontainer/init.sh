
npm config set @buf:registry https://buf.build/gen/npm/v1 && \
npm install -g @buf/connectrpc_eliza.connectrpc_es@1.4.0-20230913231627-233fca715f49.3 @connectrpc/connect@"^v1.0.0" @connectrpc/connect-web@"^v1.0.0" \
 @bufbuild/buf @connectrpc/protoc-gen-connect-es@"^1.0.0" @bufbuild/protoc-gen-es@"^1.0.0" \
 @bufbuild/protobuf@"^1.0.0"

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

# install for ebitengine

sudo apt install -y gcc
sudo apt install -y libc6-dev libgl1-mesa-dev libgl1-mesa-dev libxcursor-dev libxi-dev libxinerama-dev libxrandr-dev libxxf86vm-dev libasound2-dev pkg-config

