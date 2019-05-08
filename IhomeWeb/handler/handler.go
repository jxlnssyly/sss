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
	GETSMSCD "sss/GetSmscd/proto/example"
	POSTRET "sss/PostRet/proto/example"
	GETSSION "sss/GetSession/proto/example"
	POSTLOGIN "sss/PostLogin/proto/example"
	DELETESESSION "sss/DeleteSession/proto/example"
	GETUSERINFO "sss/GetUserInfo/proto/example"
	"github.com/julienschmidt/httprouter"
	"github.com/micro/go-grpc"
	"sss/IhomeWeb/models"
	"github.com/astaxie/beego"
	"sss/IhomeWeb/utils"
	"image"
	"github.com/afocus/captcha"
	"image/png"
)

func ExampleCall(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {
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
		beego.Info(err)
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

	w.Header().Set("Content-Type","application/json")

	cookie, err := r.Cookie("userlogin")
	if err != nil {
		// we want to augment the response
		response := map[string]interface{}{
			"errno": utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	service := grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := GETSSION.NewExampleService("go.micro.srv.GetSession", service.Client())
	rsp, err := exampleClient.GetSession(context.TODO(), &GETSSION.Request{
		Sessionid: cookie.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	data := make(map[string]string)
	data["name"] = rsp.UserName
	// we want to augment the response
	response := map[string]interface{}{
		"errno": rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data": data,
	}

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

// 获取短信验证码
func GetSmscd(w http.ResponseWriter, r *http.Request,  ps httprouter.Params) {
	beego.Info("GetSmscd")
	w.Header().Set("Content-Type","application/json")

	// URL里面的请求参数
	text := r.URL.Query()["text"][0]
	id := r.URL.Query()["id"][0]
	mobile := ps.ByName("mobile")

	/*
	// 通过正则进行手机号的判断
	mobile_re := regexp.MustCompile(`0?(13|14|15|17|18|19)[0-9]{9}`)
	// 通过条件判断字符串是否匹配规则，返回正确或失败
	bl := mobile_re.MatchString(mobile)
	if !bl {
		// we want to augment the response
		response := map[string]interface{}{
			"error": utils.RECODE_MOBILEERR,
			"errmsg": utils.RecodeText(utils.RECODE_MOBILEERR),
		}

		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	*/

	// 创建并初始化服务
	service := grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := GETSMSCD.NewExampleService("go.micro.srv.GetSmscd", service.Client())
	rsp, err := exampleClient.GetSmscd(context.TODO(), &GETSMSCD.Request{
		Mobile: mobile,
		Imagestr: text,
		Uuid: id,
	})
	if err != nil {
		beego.Info(err)
		http.Error(w, err.Error(), 500)
		return
	}
	// we want to augment the response
	response := map[string]interface{}{
		"error": rsp.Error,
		"errmsg": rsp.ErrorMsg,
	}


	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// 注册
func PostRet(w http.ResponseWriter, r *http.Request,  ps httprouter.Params) {

	beego.Info("PostRet")

	service := grpc.NewService()
	service.Init()

	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		beego.Info(err)
		return
	}

	w.Header().Set("Content-Type","application/json")

	if request["mobile"].(string) == "" || request["password"].(string) == "" || request["sms_code"].(string) == "" {
		response := map[string]interface{}{
			"errno": utils.RECODE_DATAERR,
			"errmsg": utils.RecodeText(utils.RECODE_DATAERR),
		}

		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	// call the backend service
	exampleClient := POSTRET.NewExampleService("go.micro.srv.PostRet", service.Client())
	rsp, err := exampleClient.PostRet(context.TODO(), &POSTRET.Request{
		Mobile: request["mobile"].(string),
		Password:request["password"].(string),
		SmsCode:request["sms_code"].(string),
	})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 读取cookie 统一cookie
	cookie, err := r.Cookie("userlogin")
	if err != nil || "" == cookie.Value {
		// 创建1个cookie对象
		cookie := http.Cookie{Name:"userlogin",Value:rsp.SessionId,Path:"/",MaxAge:3600}
		// 对浏览器的cookie进行设置
		http.SetCookie(w, &cookie)
	}

	// we want to augment the response
	response := map[string]interface{}{
		"errno": rsp.Error,
		"errmsg": rsp.Errmsg,
	}

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
//登录
func PostLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	beego.Info("PostLogin ")


	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type","application/json")

	if request["mobile"].(string) == "" || request["password"].(string) == ""  {
		response := map[string]interface{}{
			"errno": utils.RECODE_DATAERR,
			"errmsg": utils.RecodeText(utils.RECODE_DATAERR),
		}

		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	service := grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := POSTLOGIN.NewExampleService("go.micro.srv.PostLogin", service.Client())
	rsp, err := exampleClient.PostLogin(context.TODO(), &POSTLOGIN.Request{
		Mobile:request["mobile"].(string),
		Password:request["password"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	cookie, err := r.Cookie("userlogin")

	// 设置cookie
	if err != nil || cookie.Value == "" {
		cookie := http.Cookie{Name:"userlogin",Value:rsp.Sessionid, Path:"/",MaxAge:600}
		http.SetCookie(w, &cookie)
	}


	// we want to augment the response
	response := map[string]interface{}{
		"errno": rsp.Errno,
		"errmsg": rsp.Errmsg,
	}

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
// 登出
func DeleteSession(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {

	beego.Info("DeleteSession")


	w.Header().Set("Content-Type","application/json")



	service := grpc.NewService()
	service.Init()
	// call the backend service
	exampleClient := DELETESESSION.NewExampleService("go.micro.srv.DeleteSession", service.Client())


	// 获取sessionid
	cookie, err := r.Cookie("userlogin")

	if cookie.Value == "" || err != nil {
		response := map[string]interface{}{
			"errno": utils.RECODE_DATAERR,
			"errmsg": utils.RecodeText(utils.RECODE_DATAERR),
		}

		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(response); err != nil {
			beego.Info(err)
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	rsp, err := exampleClient.DeleteSession(context.TODO(), &DELETESESSION.Request{
		Sessionid:cookie.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 删除sessionID
	cookie, err = r.Cookie("userlogin")

	if err == nil || cookie.Value != "" {
		cookie := http.Cookie{Name:"userlogin",Path:"/", MaxAge:-1, Value:""}
		http.SetCookie(w,&cookie)
	}

	// we want to augment the response
	response := map[string]interface{}{
		"errno": rsp.Errno,
		"errmsg": rsp.Errmsg,
	}

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// 获取用户信息
func GetUserInfo(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {

	beego.Info("GetUserInfo")

	w.Header().Set("Content-Type","application/json")

	service := grpc.NewService()
	service.Init()
	// call the backend service
	exampleClient := GETUSERINFO.NewExampleService("go.micro.srv.GetUserInfo", service.Client())

	// 获取sessionid
	cookie, err := r.Cookie("userlogin")

	if cookie.Value == "" || err != nil {
		response := map[string]interface{}{
			"errno": utils.RECODE_DATAERR,
			"errmsg": utils.RecodeText(utils.RECODE_DATAERR),
		}

		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(response); err != nil {
			beego.Info(err)
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	rsp, err := exampleClient.GetUserInfo(context.TODO(), &GETUSERINFO.Request{
		Sessionid: cookie.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	data := make(map[string]interface{})
	data["name"] = rsp.Name
	data["mobile"] = rsp.Mobile
	data["user_id"] = rsp.UserId
	data["real_name"] = rsp.RealName
	data["id_card"] = rsp.IdCard
	data["avatar_url"] = utils.AddDomain2Url(rsp.AvatarUrl)

	// we want to augment the response
	response := map[string]interface{}{
		"errno": rsp.Errorno,
		"errmsg": utils.RecodeText(rsp.Errorno),
		"data": data,
	}

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}



