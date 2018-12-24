package function

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"
	"strings"

	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"
	"golang.org/x/net/context"

	_ "github.com/micro/go-plugins/broker/rabbitmq"
	//	_ "github.com/micro/go-plugins/registry/kubernetes"
	micro "github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/registry/zookeeper"
)

func HandleJSONRPC(Servicec micro.Service, r *http.Request) []byte {

	body := []byte(`{}`)
	urla := strings.Replace(strings.ToLower(r.URL.Path), "/api/", "/", 1)
	service, method := PathToReceiver(Namespace, urla)

	rbody, _ := ioutil.ReadAll(r.Body)
	if len(rbody) < 2 {
		rbody = []byte("{}")
	}
	request := json.RawMessage(rbody)
	var response json.RawMessage
	req := (Servicec.Client()).NewRequest(service, method, &request, client.WithContentType("application/json"))
	ctx := RequestToContext(r)
	err := (Servicec.Client()).Call(ctx, req, &response)
	if err != nil {
		//(ce.Error()))
		return []byte("{}")
	}

	return body
}

func PathToReceiver(ns, p string) (string, string) {
	p = path.Clean(p)
	p = strings.TrimPrefix(p, "/")
	parts := strings.Split(p, "/")

	if len(parts) <= 2 {
		service := ns + strings.Join(parts[:len(parts)-1], ".")
		method := strings.Title(strings.Join(parts, "."))
		return service, method
	}

	if len(parts) == 3 && VersionRe.Match([]byte(parts[0])) {
		service := ns + strings.Join(parts[:len(parts)-1], ".")
		method := strings.Title(strings.Join(parts[len(parts)-2:], "."))
		return service, method
	}

	service := ns + strings.Join(parts[:len(parts)-2], ".")
	method := strings.Title(strings.Join(parts[len(parts)-2:], "."))
	return service, method
}

func RequestToContext(r *http.Request) context.Context {
	ctx := context.Background()
	md := make(metadata.Metadata)
	for k, v := range r.Header {
		md[k] = strings.Join(v, ",")
	}
	return metadata.NewContext(ctx, md)
}
func Getview(url string) (bool, string) {

	parts := strings.Split(url, "/")
	if len(parts) >= 3 {
		if len(parts[2]) > 2 {
			return true, parts[2] + ".html"
		}
		return false, ""
	}

	return false, ""

}
