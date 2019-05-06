package main

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
)

//static char TEST_COS_ENDPOINT[] = "cos.ap-guangzhou.myqcloud.com";
//static char TEST_ACCESS_KEY_ID[] = "AKIDSkoeFcXeJjtqY1RExQdrpXy93h4AsZsu";
//static char TEST_ACCESS_KEY_SECRET[] = "uVrxnNhlZR5TcWBbomVyB8pHTP5VfZDe";
//static char TEST_APPID[] = "1255353254";
//static char TEST_BUCKET_NAME[] = "detection-image";

func main() {
	//将<bucketname>、<appid>和<region>修改为真实的信息
	//例如：http://test-1253846586.cos.ap-guangzhou.myqcloud.com
	u, _ := url.Parse("http://opencv-1255353254.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			//如实填写账号和密钥，也可以设置为环境变量
			SecretID:  "AKIDSkoeFcXeJjtqY1RExQdrpXy93h4AsZsu",
			SecretKey: "uVrxnNhlZR5TcWBbomVyB8pHTP5VfZDe",
		},
	})
	//对象键（Key）是对象在存储桶中的唯一标识。
	//例如，在对象的访问域名 ` bucket1-1250000000.cos.ap-guangzhou.myqcloud.com/test/objectPut.go ` 中，对象键为 test/objectPut.go
	//name := "001/opencv-4.0.1.zip"
	name := "00000000e9617414/20190418/0000/0000_2_1555555556.zip"
	//Local file
	f := strings.NewReader("./0000_2_1555555556.zip")
	//f := strings.NewReader("./opencv-4.0.1.zip")

	_, err := c.Object.Put(context.Background(), name, f, nil)
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 10)
}
