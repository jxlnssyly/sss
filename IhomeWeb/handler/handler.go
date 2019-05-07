package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/micro/go-micro/client"
	example "github.com/micro/examples/template/srv/proto/example"
	// 调用area的proto
	GETAREA "sss/GetArea/proto/example"
	GETIMAGED "sss/GetImaged/proto/example"
	"github.com/julienschmidt/httprouter"
	"github.com/micro/go-grpc"
	"sss/IhomeWeb/models"
	"github.com/astaxie/beego"
	"pra/IhomeWeb/utils"
	"image"
	"github.com/afocus/captcha"
	"image/png"
)

func ExampleCall(w http.ResponseWriter, r *http.Request) {
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// call the backend service
	exampleClient := example.NewExampleService("go.micro.srv.template", client.DefaultClient)
	rsp, err := exampleClient.Call(context.TODO(), &example.Request{
		Name: request["name"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"msg": rsp.Msg,
		"ref": time.Now().UnixNano(),
	}

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// 获取地区信息
func GetArea(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {
	beego.Info("请求地区信息 GetArea api/v1.0/areas")

	// 创建服务获取句柄
	server := grpc.NewService()

	// 服务初始化
	server.Init()


	// 调用服务，返回句柄
	exampleClient := GETAREA.NewExampleService("go.micro.srv.GetArea", server.Client())
	// 调用服务，返回数据
	rsp, err := exampleClient.GetArea(context.TODO(), &GETAREA.Request{})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}



	// 接收数据
	// 准备接收切片
	area_list := []models.Area{}
	for _,value := range rsp.Data {
		tmp := models.Area{Id:int(value.Aid), Name:value.Aname}
		area_list = append(area_list,tmp)
	}

	// 返回给前端的map
	response := map[string]interface{}{
		"errno": rsp.Error,
		"errmsg": rsp.Errmsg,
		"data": area_list,
		"ref": time.Now().UnixNano(),
	}

	// 设置数据格式
	w.Header().Set("Content-Type","application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// 获取验证图片
func GetImageCode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {


	beego.Info("GetImageCode ")

	// 创建服务
	service := grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := GETIMAGED.NewExampleService("go.micro.srv.GetImaged", service.Client())

	// 获取uuid
	uuid := ps.ByName("uuid")

	rsp, err := exampleClient.GetImaged(context.TODO(), &GETIMAGED.Request{
		Uuid: uuid,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 接收图片信息的图片格式
	var img image.RGBA
	img.Stride = int(rsp.Stride)
	img.Pix = []uint8(rsp.Pix)
	img.Rect.Min.X = int(rsp.Min.X)
	img.Rect.Min.Y = int(rsp.Min.Y)
	img.Rect.Max.Y = int(rsp.Max.Y)
	img.Rect.Max.X = int(rsp.Max.X)

	var capt captcha.Image
	capt.RGBA = &img


	png.Encode(w, capt)
}


// 获取session信息
func GetSession(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {
	beego.Info("GetSession ")

	// we want to augment the response
	response := map[string]interface{}{
		"errno": utils.RECODE_SESSIONERR,
		"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
	}
	// 设置数据格式
	w.Header().Set("Content-Type","application/json")
	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// 获取首页轮播图
func GetIndex(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {

	beego.Info("GetIndex ")

	// we want to augment the response
	response := map[string]interface{}{
		"errno": utils.RECODE_OK,
		"errmsg": utils.RecodeText(utils.RECODE_OK),
	}
	// 设置数据格式
	w.Header().Set("Content-Type","application/json")
	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}


