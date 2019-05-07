package handler

import (
	"context"

	"github.com/micro/go-log"

	example "sss/GetSmscd/proto/example"
	"github.com/astaxie/beego"
	"sss/IhomeWeb/utils"
	"github.com/astaxie/beego/orm"
	"sss/IhomeWeb/models"
	"encoding/json"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/garyburd/redigo/redis"
	"github.com/astaxie/beego/cache"
	"github.com/garyburd/redigo/redis"
	"math/rand"
	"time"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetSmscd(ctx context.Context, req *example.Request, rsp *example.Response) error {

	beego.Info("GetSmscd -- server")
	rsp.Error = utils.RECODE_OK
	rsp.ErrorMsg = utils.RecodeText(rsp.Error)

	// 验证手机号是否存在
	o := orm.NewOrm()
	// 使用手机号作为查询条件
	user := models.User{Mobile:req.Mobile}
	err := o.Read(&user,"mobile")
	// 如果不报错就说明查找到了
	// 如果找到就说明手机号存在
	if err == nil {
		beego.Info("用户已存在")
		rsp.Error = utils.RECODE_MOBILEERR
		rsp.ErrorMsg = utils.RecodeText(rsp.Error)
		return nil
	}


	// 验证图片验证码是否正确

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
		rsp.ErrorMsg = utils.RecodeText(rsp.Error)
		return nil
	}

	// 通过uuid查找图片验证码的值进行对比
	value := bm.Get(req.Uuid)
	if value == nil {
		beego.Info("Redis中取值失败",err)
		rsp.Error = utils.RECODE_DBERR
		rsp.ErrorMsg = utils.RecodeText(rsp.Error)
		return nil
	}

	// reflect.TypeOf(value) 会返回当前数据的变量类型
	//beego.Info(reflect.TypeOf(value))
	value_str, _ := redis.String(value,nil)
	if value_str != req.Imagestr {
		beego.Info("请输入正确的图片验证码",err)
		rsp.Error = utils.RECODE_DATAERR
		rsp.ErrorMsg = utils.RecodeText(rsp.Error)
		return nil
	}
	r :=rand.New(rand.NewSource(time.Now().UnixNano()))
	size := r.Intn(9999)+1001
	beego.Info("验证码",size)
	bm.Put(req.Mobile, size, time.Second * 300)
	if err!=nil {
		beego.Info("redis创建失败",err)
		rsp.Error = utils.RECODE_DBERR
		rsp.ErrorMsg = utils.RecodeText(rsp.Error)
		return nil
	}
	//// 调用短信接口发送短信
	//messageconfig := make(map[string]string)
	//
	////发送短信的配置信息
	//messageconfig := make(map[string]string)
	////预先创建好的appid
	//messageconfig["appid"] = "29672"
	////预先获得的app的ket
	//messageconfig["appkey"] = "89d90165cbea8cae80137d7584179bdb"
	////加密方式默认
	//messageconfig["signtype"] = "md5"
	//
	////messagexsend
	////创建短信发送的句柄
	//messagexsend := submail.CreateMessageXSend()
	////短信发送的手机号
	//submail.MessageXSendAddTo(messagexsend, req.Mobile)
	////短信发送的模板
	//submail.MessageXSendSetProject(messagexsend, "NQ1J94")
	////验证码
	//submail.MessageXSendAddVar(messagexsend, "code", strconv.Itoa(size))
	////发送短信的请求
	//fmt.Println("MessageXSend ", submail.MessageXSendRun(submail.MessageXSendBuildRequest(messagexsend), messageconfig))


	// 将短信验证码存入短信数据库

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
