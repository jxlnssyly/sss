package handler

import (
	"context"

	"github.com/micro/go-log"

	example "sss/PostMutilImage/proto/example"
	"github.com/astaxie/beego"
	"sss/IhomeWeb/utils"
	"path"
	"strconv"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PostMutilImage(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("PostMutilImage -- server")

	bm, _ := utils.GetRedisServer()

	for index, image := range req.Images{
		// 图片数据验证
		if image.Filesize != int64(len(image.Avatar)) {
			beego.Info("传输数据丢失")

		}

		// 获取文件的后缀名
		ext := path.Ext(image.FileExt)
		// 路径..jpg
		// 调用fdfs函数上传到图片服务器
		fileid, err := utils.UploadByBuffer(image.Avatar,ext[1:])
		if err != nil {
			beego.Info("头像上传失败", err)
		}

		beego.Info(fileid)
		bm.Put("avatar"+strconv.Itoa(index),fileid,3600)

	}


	// 获取sessionid

	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Example) Stream(ctx context.Context, req *example.StreamingRequest, stream example.Example_StreamStream) error {
	log.Logf("Received Example.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&example.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Example) PingPong(ctx context.Context, stream example.Example_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&example.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
