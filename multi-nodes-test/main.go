package main

import (
	multi_nodes "GoCache/multi-nodes"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func creatGroup() *multi_nodes.Group {
	return multi_nodes.NewGroup("scores", 2<<10,
		multi_nodes.GetterFunc(func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exisit", key)
		}))
}

//startCacheServer()用来启动缓存服务器：创建HTTPPool，添加节点信息，注册到gee中，启动HTTP服务
func startCacheServer(addr string, addrs []string, gee *multi_nodes.Group) {
	peers := multi_nodes.NewHTTPPool(addr)
	peers.Set(addrs...)
	gee.RegisterPeers(peers)
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}
func startAPIServer(apiAddr string, gee *multi_nodes.Group) {
	//启动api路由，根据api做响应的处理
	http.Handle("/api", http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			key := request.URL.Query().Get("key")
			view, err := gee.Get(key)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			writer.Header().Set("Content-Type", "application/octet-stream")
			writer.Write(view.ByteSlice())
		}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}
func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "Geecache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()
	apiAddr := "http://lochaost:9999"
	addrMap := map[int]string{
		8001: "http://localhsot:8001",
		8002: "http://localhsot:8002",
		8003: "http://localhsot:8003",
	}
	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}
	gee := creatGroup()
	if api {
		go startAPIServer(apiAddr, gee)
	}
	startCacheServer(addrMap[port], []string(addrs), gee)
}
