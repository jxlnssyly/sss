package handler

import (
	"context"

	"github.com/micro/go-log"

	example "sss/GetImaged/proto/example"
	"github.com/astaxie/beego"
	"github.com/afocus/captcha"
	"image/color"

	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/garyburd/redigo/redis"
	_ "github.com/garyburd/redigo/redis"
	"encoding/json"
	"github.com/astaxie/beego/cache"
	"sss/IhomeWeb/utils"
	"time"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetImaged(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("GetImageCode -- server ")

	cap := captcha.New()

	if err := cap.SetFont("comic.ttf"); err != nil {
		panic(err.Error())
	}

	cap.SetSize(90, 41)
	cap.SetDisturbance(captcha.NORMAL)
	cap.SetFrontColor(color.RGBA{255, 255, 255, 255})
	cap.SetBkgColor(color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 153, 0, 255})

	// 生成随机的验证码图片
	img, str := cap.Create(4, captcha.NUM)

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
	}

	// 验证码缓存与uuid进行缓存
	bm.Put(req.Uuid,str,time.Second * 300)

	// 图片解引用
	img1 := *img
	img2 := *img1.RGBA

	rsp.Error = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Error)

	// 返回图片拆分
	rsp.Pix = []byte(img2.Pix)
	rsp.Stride = int64(img2.Stride)
	rsp.Max = &example.Response_Point{X:int64(img2.Rect.Max.X),Y:int64(img2.Rect.Max.Y)}
	rsp.Min = &example.Response_Point{X:int64(img2.Rect.Min.X),Y:int64(img2.Rect.Min.Y)}



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
