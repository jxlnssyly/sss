package handler

import (
	"context"

	"github.com/micro/go-log"

	example "sss/PostAvatar/proto/example"
	"github.com/astaxie/beego"

	_ "github.com/astaxie/beego/cache/redis"
	 "github.com/garyburd/redigo/redis"
	"sss/IhomeWeb/utils"
	"path"
	"strconv"
	"sss/IhomeWeb/models"
	"github.com/astaxie/beego/orm"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PostAvatar(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("PostAvatar -- server")

	// 初始化返回值
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	// 图片数据验证
	if req.Filesize != int64(len(req.Avatar)) {
		beego.Info("传输数据丢失")
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
	}

	// 获取文件的后缀名
	ext := path.Ext(req.FileExt)
	// 路径..jpg
	// 调用fdfs函数上传到图片服务器
	fileid, err := utils.UploadByBuffer(req.Avatar,ext[1:])
	if err != nil {
		beego.Info("头像上传失败", err)
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
	}

	beego.Info(fileid)
	// 获取sessionid
	bm, err := utils.GetRedisServer()
	if err != nil {
		beego.Info("Redis连接失败",err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	// 连接Redis， 拼接key 获取当前用户的user_id
	sessinid := req.Sessionid
	session_userid := sessinid + "user_id"
	user_id := bm.Get(session_userid)
	user_id_str, _ := redis.String(user_id, nil)
	id , _ := strconv.Atoi(user_id_str)
	user := models.User{Id:id,Avatar_url:fileid}

	// 将图片的存储地址 更新到user表中
	o := orm.NewOrm()
	_, err = o.Update(&user,"avatar_url")
	if err != nil {
		beego.Info("数据更新",err)
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
	}

	// 回传头像url
	rsp.AvatarUrl = fileid
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
