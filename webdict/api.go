package webdict

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"github.com/go-martini/martini"
//    "github.com/martini-contrib/render"
	"net/http"
//    "net/http/httputil"
	"time"
)

func createResponse(value string) [] byte {
	type Response struct {
		Result string `json:"result"`;
		Time string   `json:"time"`;
	}
	result, _ := json.Marshal(Response{value, time.Now().Format("2006-01-02 03:04")})
	return result
}

func parseRequest(req *http.Request) (key, value string, err error) {
	type Params struct {
		Key   string `json:"key"`;
		Value string `json:"value"`;
	}
	var params Params
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &params)
	if err != nil {
		return
	}
	key, value = params.Key, params.Value
	return
}

func NewApi(prefixPath string) *martini.ClassicMartini {
	data := make(map[string]string)
	data["test"] = "my value"
	api := martini.Classic()
	api.Get(prefixPath + "/:key", func(params martini.Params, res http.ResponseWriter) {
		key := params["key"]
		if value, ok := data[key]; ok {
			// 200
			res.WriteHeader(http.StatusOK)
			res.Write(createResponse(value))
		} else {
			// 404
			res.WriteHeader(http.StatusNotFound)
			res.Write([]byte(fmt.Sprintf("Can not find key: %v", key)))
		}
	})
	api.Get(prefixPath, func(res http.ResponseWriter) {
		// 400
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Missed key value on GET request"))
	})
	api.Post(prefixPath, func(req *http.Request, res http.ResponseWriter) {
		key, value, err := parseRequest(req)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte(fmt.Sprintf("Can not parse POST request: %v", err)))
			return
		}
		actual, ok := data[key]
		if ok {
			// 409
			res.WriteHeader(http.StatusConflict)
			res.Write(createResponse(actual))
		} else {
			data[key] = value
		}
	})
	api.Put(prefixPath, func(req *http.Request, res http.ResponseWriter) {
		key, value, err := parseRequest(req)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte(fmt.Sprintf("Can not parse POST request: %v", err)))
			return
		}
		_, ok := data[key]
		if !ok {
			// 404
			res.WriteHeader(http.StatusNotFound)
			res.Write([]byte(fmt.Sprintf("Can not find key: %v", key)))
		} else {
			data[key] = value
		}
	})
	api.Delete(prefixPath + "/:key", func(params martini.Params, res http.ResponseWriter) {
		key := params["key"]
		delete(data, key)
		res.WriteHeader(http.StatusOK)
	})
	return api
}

func Run() {
	api := NewApi("/dictionary")
	api.Run()
}
