package main

import (
	"context"

	"github.com/mhwwhu/QuickStone/src/rpc/metadata"
)

type MetadataService struct {
	metadata.MetadataServiceServer
}

func (MetadataService) Init() {

}

func (MetadataService) RegisterUploadingObject(ctx context.Context, req *metadata.RegisterUploadingObjectRequest) (
	resp *metadata.RegisterUploadingObjectResponse, err error) {
	return &metadata.RegisterUploadingObjectResponse{StatusCode: 0, StatusMsg: ""}, nil
	// TODO: 后期完善
}
