package handler

import (
	"context"
	"github.com/micro/go-log"
	example "sss/GetArea/proto/example"
	"github.com/astaxie/beego"
	"sss/IhomeWeb/utils"
	"github.com/astaxie/beego/orm"
	"sss/IhomeWeb/models"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/garyburd/redigo/redis"
	"encoding/json"
	"time"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetArea(ctx context.Context, req *example.Request, rsp *example.Response) error {

	beego.Info("请求地区信息 GetArea api/v1.0/areas")

	// 初始化 错误码
	rsp.Error = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Error)

	bm, err := utils.GetRedisServer()

	if err != nil {
		beego.Info("Redis连接失败",err)
		rsp.Error = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Error)
	}

	// 获取数据 在这里需要定制1个key 就算用来area查询
	area_value := bm.Get("area_info")
	if area_value != nil {
		 beego.Info("获取到地域信息缓存")
		 area_map := []map[string]interface{}{}
		 json.Unmarshal(area_value.([]byte),&area_map)

		//beego.Info("从缓存中提取的area数据",area_map)
		//将查询到的数据按照proto发送给web服务
		for _, value := range area_map {
			tmp := example.Response_Area{Aid:int32(value["aid"].(float64)),Aname:value["aname"].(string)}
			rsp.Data = append(rsp.Data, &tmp)
		}
		return nil
	}

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

	// 将查出来的数据存到Redis中
	// 将获取到的数据转化为json
	area_json, _:= json.Marshal(area)

	// 操作Redis将数据存入
	err = bm.Put("area_info",area_json,time.Second * 3600)
	if err != nil {
		beego.Info("数据缓存失败",err)
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
