package handler

import (
	"context"

	"github.com/micro/go-log"

	example "sss/GetArea/proto/example"
	"github.com/astaxie/beego"
	"sss/IhomeWeb/utils"
	"github.com/astaxie/beego/orm"
	"sss/IhomeWeb/models"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetArea(ctx context.Context, req *example.Request, rsp *example.Response) error {

	beego.Info("请求地区信息 GetArea api/v1.0/areas")

	// 初始化 错误码
	rsp.Error = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Error)

	// 缓存中获取

	// 缓存中没有，从mysql中获取
	var area []models.Area
	o := orm.NewOrm()
	num, err := o.QueryTable(&models.Area{}).All(&area)
	if err != nil {
		beego.Info("数据库查询失败", err)
		rsp.Error = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Error)
		return nil
	}

	if num == 0 {
		beego.Info("数据库没有数据")

		rsp.Error = utils.RECODE_NODATA
		rsp.Errmsg = utils.RecodeText(rsp.Error)
		return nil
	}

	//将查询到的数据按照proto发送给web服务
	for _, value := range area {
		tmp := example.Response_Area{Aid:int32(value.Id),Aname:value.Name}
		rsp.Data = append(rsp.Data, &tmp)
	}

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
