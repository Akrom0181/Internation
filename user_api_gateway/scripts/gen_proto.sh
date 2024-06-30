CURRENT_DIR=$1

sudo rm -rf ./genproto/*

for x in $(find ${CURRENT_DIR}/protos/* -type d); do
  sudo protoc --plugin="protoc-gen-go=${GOPATH}/bin/protoc-gen-go" --plugin="protoc-gen-go-grpc=${GOPATH}/bin/protoc-gen-go-grpc" -I=${x} -I=${CURRENT_DIR}/protos -I /usr/local/include --go_out=${CURRENT_DIR} \
   --go-grpc_out=require_unimplemented_servers=false:${CURRENT_DIR} ${x}/*.proto
done


