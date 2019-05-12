package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/garyburd/redigo/redis"
	"github.com/astaxie/beego/cache"
	"github.com/weilaihui/fdfs_client"
	"fmt"
)

/* 将url加上 http://IP:PROT/  前缀 */
//http:// + 127.0.0.1 + ：+ 8080 + 请求

func AddDomain2Url(url string) (domain_url string) {
	domain_url = "http://" + G_fastdfs_addr + ":" + G_fastdfs_port + "/" + url

	return domain_url
}

// 加密函数
func MD5String(s string) string  {
	// 创建一个md5对象
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func GetRedisServer() (adapter cache.Cache, err error) {

	// 缓存中获取
	// 准备连接Redis信息
	redis_conf := map[string]string{
		"key": G_server_name,
		"conn": G_redis_addr + ":"+G_redis_port,
		"dbNum": G_redis_dbnum,
	}

	// 将map转化成json
	redis_conf_json, _ := json.Marshal(redis_conf)
	// 创建Redis句柄
	bm, err := cache.NewCache("redis",string(redis_conf_json))

	return bm, err
}

// 上传二进制文件到fdfs
func UploadByBuffer(filebuffer []byte, fileExt string) (fildid string,err error)  {
	fd_client, err := fdfs_client.NewFdfsClient("/Users/user/go/src/sss/IhomeWeb/conf/client.conf")
	if err != nil {
		fmt.Println("创建fdfs句柄失败",err)
		fildid = ""
		return
	}
	fd_rsp, err := fd_client.UploadByBuffer(filebuffer,fileExt)
	if err != nil {
		fmt.Println("上传失败",err)
		fildid = ""
		return
	}
	fmt.Println(fd_rsp)
	fmt.Println(fd_rsp.RemoteFileId)
	fildid = fd_rsp.RemoteFileId

	return fildid, nil
}
