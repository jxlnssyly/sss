package handler

import (
	"context"

	"github.com/micro/go-log"

	example "sss/PostLogin/proto/example"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"sss/IhomeWeb/models"

	"sss/IhomeWeb/utils"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/garyburd/redigo/redis"
	"time"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PostLogin(ctx context.Context, req *example.Request, rsp *example.Response) error {

	beego.Info("PostLogin -- server")
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	// 查询数据
	o := orm.NewOrm()
	user := models.User{}
	qs := o.QueryTable(&user)
	// 通过QS句柄进行查询
	err := qs.Filter("mobile",req.Mobile).One(&user)
	if err != nil {
		beego.Info("查询数据失败")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}

	// 密码的校验
	if utils.MD5String(req.Password) != user.Password_hash {
		beego.Info("密码错误")

		rsp.Errno = utils.RECODE_PWDERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_PWDERR)
		return nil
	}

	// 创建sessionid
	sessionid := utils.MD5String(req.Mobile + req.Password)
	rsp.Sessionid = sessionid


	// 拼接Key
	// 准备连接Redis信息
	bm, err := utils.GetRedisServer()

	if err != nil {
		beego.Info("Redis连接失败",err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	// 将登录信息缓存
	sessionuser_id := sessionid + "user_id"
	bm.Put(sessionuser_id, user.Id,time.Second * 3600)

	// name
	sessionname := sessionid+"name"
	bm.Put(sessionname, user.Name,time.Second * 3600)

	// mobile
	sessionmobile := sessionid + "mobile"

	bm.Put(sessionmobile, user.Mobile,time.Second * 3600)



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
