package handler

import (
	"context"

	"github.com/micro/go-log"

	example "sss/GetUserInfo/proto/example"
	"sss/IhomeWeb/utils"
	"github.com/astaxie/beego"
	"sss/IhomeWeb/models"
	"github.com/astaxie/beego/orm"
	"strconv"
	"github.com/garyburd/redigo/redis"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetUserInfo(ctx context.Context, req *example.Request, rsp *example.Response) error {

	beego.Info("GetUserInfo server")
	// 初始化错误码
	rsp.Errorno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errorno)

	// 获取sessionid
	sessionid := req.Sessionid

	// 连接redis
	bm, err := utils.GetRedisServer()
	if err != nil {
		beego.Info("Redis连接失败",err)
		rsp.Errorno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errorno)
		return nil
	}

	// 拼接Key
	sessionuser_id := sessionid + "user_id"

	// 通过key获取到user_id
	user_id := bm.Get(sessionuser_id)
	userid_str, _ := redis.String(user_id, nil)
	id, _ := strconv.Atoi(userid_str)
	// 通过user_id获取都用户表信息
	// 创建1个user对象
	user := models.User{Id:id}
	o := orm.NewOrm()
	err = o.Read(&user)
	if err != nil {
		beego.Info("数据获取失败",err)
		rsp.Errorno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Errorno)
		return nil
	}


	// 将信息返回
	rsp.UserId = strconv.Itoa(user.Id)
	rsp.Name = user.Name
	rsp.RealName = user.Real_name
	rsp.IdCard = user.Id_card
	rsp.Mobile = user.Mobile
	rsp.AvatarUrl = user.Avatar_url

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
