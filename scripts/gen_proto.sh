#!/bin/sh

# Tentukan direktori output
OUTPUT_DIR=.

# Hapus file lama jika ada (opsional)
# rm -f api/protos/*.pb.go

echo "Generating Go code from .proto files..."

# Perintah protoc yang sebenarnya
protoc --proto_path=. \
       --go_out=${OUTPUT_DIR} \
       --go-grpc_out=${OUTPUT_DIR} \
       api/protos/transaction.proto

echo "Done."