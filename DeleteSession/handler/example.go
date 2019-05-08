package handler

import (
	"context"

	"github.com/micro/go-log"

	example "sss/DeleteSession/proto/example"
	"github.com/astaxie/beego"
	"sss/IhomeWeb/utils"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) DeleteSession(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("DeleteSession -- server")

	// 返回值初始化
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	// 准备连接Redis信息
	bm, err := utils.GetRedisServer()

	if err != nil {
		beego.Info("Redis连接失败",err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	sessionid := req.Sessionid

	// 拼接key
	// 将登录信息缓存
	sessionuser_id := sessionid + "user_id"
	bm.Delete(sessionuser_id)
	// name
	sessionname := sessionid+"name"
	bm.Delete(sessionname)

	// mobile
	sessionmobile := sessionid + "mobile"
	bm.Delete(sessionmobile)

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
