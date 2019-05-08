package utils

import (
	"crypto/md5"
	"encoding/hex"
	"pra/IhomeWeb/utils"
	"encoding/json"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/garyburd/redigo/redis"
	"github.com/astaxie/beego/cache"
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
		"key": utils.G_server_name,
		"conn": utils.G_redis_addr + ":"+utils.G_redis_port,
		"dbNum": utils.G_redis_dbnum,
	}

	// 将map转化成json
	redis_conf_json, _ := json.Marshal(redis_conf)
	// 创建Redis句柄
	bm, err := cache.NewCache("redis",string(redis_conf_json))

	return bm, err
}
