package handler

import (
	"context"

	"github.com/micro/go-log"

	example "sss/PostRet/proto/example"
	"github.com/astaxie/beego"
	"sss/IhomeWeb/utils"
	"encoding/json"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/garyburd/redigo/redis"
	"github.com/garyburd/redigo/redis"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/orm"
	"sss/IhomeWeb/models"
	"time"
)

type Example struct{}



// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PostRet(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("PostRet server")

	rsp.Error = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Error)
	// 验证短信验证码
	// redis
	// 准备连接Redis信息
	//
	redis_conf := map[string]string{
		"key": utils.G_server_name,
		"conn": utils.G_redis_addr + ":"+utils.G_redis_port,
		"dbNum": utils.G_redis_dbnum,
	}

	// 将map转化成json
	redis_conf_json, _ := json.Marshal(redis_conf)
	// 创建Redis句柄
	bm, err := cache.NewCache("redis",string(redis_conf_json))
	if err != nil {
		beego.Info("Redis连接失败",err)
		rsp.Error = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Error)
		return nil
	}

	// 通过手机号获取到短信验证码
	sms_code := bm.Get(req.Mobile)
	if sms_code == nil {
		 beego.Info("获取数据失败")
		rsp.Error = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Error)
		return nil
	}

	// 短信验证码对比
	sms_code_str, _ := redis.String(sms_code,nil)
	if sms_code_str != req.SmsCode {
		beego.Info("短信验证码错误")
		rsp.Error = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Error)
		return nil
	}

	// 将数据存入数据库
	o := orm.NewOrm()
	user := models.User{Mobile:req.Mobile,Password_hash:utils.MD5String(req.Password),Name:req.Mobile}
	id, err := o.Insert(&user)
	if err != nil {
		beego.Info("注册失败")
		rsp.Error = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Error)
		return nil
	}
	beego.Info("user_id",id)

	// 创建sessionid(唯一随机码)
	sessionid := utils.MD5String(req.Mobile + req.Password)
	rsp.SessionId = sessionid

	// 以sessionid为key的一部分创建session
	bm.Put(sessionid + "name",user.Mobile, time.Second * 3600)
	bm.Put(sessionid + "user_id",id, time.Second * 3600)
	bm.Put(sessionid + "mobile",user.Mobile, time.Second * 3600)

	// 创建session


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
