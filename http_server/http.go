package http_server

import (
	"GoCache/single_node"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/_geecache"

type HTTPPool struct {
	self     string
	basePath string
}

//首先我们创建一个结构体HTTPool，作为承载节点间HTTP通信的核心数据结构
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,            //用来记录自己的地址(主机和ip和端口号)
		basePath: defaultBasePath, //作为节点通信的前缀
	}
}

// Log info with server name
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

//serverHttp handle all http requests
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path:" + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	groupName := parts[0]
	key := parts[1]
	group := single_node.GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group:"+groupName, http.StatusNotFound)
		return
	}
	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())

}
