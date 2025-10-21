#!/bin/bash

if [ "$#" -lt 3 ]; then
  echo "Usage: $0 <proto_dir> <out_dir> <file1.proto> [file2.proto] [...]"
  exit 1
fi

PROTO_DIR=$1
OUT_DIR=$2

if ! command -v protoc-gen-go &> /dev/null || ! command -v protoc-gen-go-grpc &> /dev/null; then
  echo "Error: protoc-gen-go or protoc-gen-go-grpc not installed."
  echo "Run: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
  echo "     go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"
  exit 1
fi

for file in "${@:3}"; do
  if [ ! -f "${PROTO_DIR}/$file" ]; then
    echo "File not found: $file"
    continue
  fi

  echo "Compiling $file ..."
  protoc \
    --proto_path="${PROTO_DIR}" \
    --go_out="${OUT_DIR}" --go_opt=paths=source_relative \
    --go-grpc_out="${OUT_DIR}" --go-grpc_opt=paths=source_relative \
    "$file"
done

echo "Compilation finished. Output in ${OUT_DIR}"