package main

import (
	"bytes"
	"context"

	"QuickStone/src/common"
	"QuickStone/src/config"
	"QuickStone/src/constant"
	trans "QuickStone/src/rpc/transmission"
	"QuickStone/src/storage"

	"github.com/sirupsen/logrus"
)

type TransmissionService struct {
	trans.TransmissionServiceServer
}

func (TransmissionService) Init() {

}

func (s TransmissionService) UploadObject(stream trans.TransmissionService_UploadObjectServer) error {
	req, err := stream.Recv()
	if err != nil {
		stream.SendAndClose(&trans.UploadObjectResponse{
			StatusCode: constant.GrpcCommunicationErrorCode,
		})
		return err
	}

	head := req.GetHeader()
	storagePath := common.StoragePath{
		UserId: head.TargetUserId,
		Bucket: head.Bucket,
		Key:    head.Key,
	}

	var buf bytes.Buffer
	frameNum := uint32((head.ObjectSize-1)/config.GrpcStreamUploadSliceSize + 1)
	frameNo := uint32(0)
	for frameNo < frameNum {
		req, err := stream.Recv()
		body := req.GetData()
		if err != nil || frameNo != body.SeriesNo {
			logrus.Errorf("Grpc communication error: err = %v, frameNo = %d, seriesNo = %d", err, frameNo, body.SeriesNo)
			stream.SendAndClose(&trans.UploadObjectResponse{
				StatusCode: constant.GrpcCommunicationErrorCode,
			})
			return err
		}
		frameNo++

		data := body.Data
		if _, err := buf.Write(data); err != nil {
			stream.SendAndClose(&trans.UploadObjectResponse{
				StatusCode: constant.GolangInternalErrorCode,
			})
			return err
		}
	}

	reader := bytes.NewReader(buf.Bytes())

	// TODO: 目前是同步等待存储模块，后续可以改成MQ异步发布任务
	err = storage.StorageClient.UploadObject(stream.Context(), storagePath, reader)
	if err != nil {
		stream.SendAndClose(&trans.UploadObjectResponse{
			StatusCode: constant.StorageUploadErrorCode,
		})
		return err
	}

	stream.SendAndClose(&trans.UploadObjectResponse{
		StatusCode: 0,
	})

	return nil
}

func (s TransmissionService) DownloadObject(ctx context.Context, req *trans.DownloadObjectRequest) (*trans.DownloadObjectResponse, error) {
	return nil, nil
}
