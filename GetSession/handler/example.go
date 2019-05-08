package handler

import (
	"context"

	"github.com/micro/go-log"

	example "sss/GetSession/proto/example"
	"sss/IhomeWeb/utils"
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/garyburd/redigo/redis"
	"github.com/garyburd/redigo/redis"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetSession(ctx context.Context, req *example.Request, rsp *example.Response) error {

	// 初始化返回值
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	// 准备连接Redis信息

	bm, _ := utils.GetRedisServer()


	// 获取username

	username := bm.Get(req.Sessionid + "name")

	if username == nil {
		beego.Info("获取数据不存在")
		rsp.Errno = utils.RECODE_DBERR

		rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
	}
	rsp.UserName, _ = redis.String(username, nil)


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
