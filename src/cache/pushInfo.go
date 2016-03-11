package cache

import (
	"log"
	"net/http"
	"net/url"
)

import (
//	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
//	"log"
)

// 推送消息的管道长度
const PUSH_INFO_LEN = 2000
const PUSH_SERVER = "http://apns.em-brain.cn/relay.aspx"

// 推送消息结构
//type PushInfo url.Values

// 推送管道
var G_PushInfo = make(chan url.Values, PUSH_INFO_LEN)

// 推送消息
func PushToServer(chInfo <-chan url.Values) {	
	log.Println("begin push module:")
	
	// 发布消息
	for v := range chInfo {
		//data := make(url.Values)
		

		/*switch v[type].(string){
			case "position_share"
				data["type"] = []string{"position_share"}
				data["destuid"] = []string{v["destuid"]}
				data[""]
		}

		switch v.Topic {
		case "send_message":
			data["type"] = []string{"send_message"}
			data["uid"] = []string{v.Content["uid"].(string)}
			data["nickname"] = []string{v.Content["nickname"].(string)}
			data["tel"] = []string{v.Content["tel"].(string)}
			data["message"] = []string{v.Content["message"].(string)}
			data["time"] = []string{v.Content["time"].(string)}
		}*/
		ret, err := http.PostForm(PUSH_SERVER, v)
		
		if err != nil {
			log.Println("Push2Servr: post error!")
			return
		}
		log.Println("ret: ", ret)
		log.Println("push obj:", v)
	}
}
