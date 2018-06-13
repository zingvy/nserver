# nserver

A lightweight and easy to use micro server framework written in Go based on [gnatsd](https://github.com/nats-io/gnatsd)

# Quickstart

- install and run [gnatsd](https://github.com/nats-io/gnatsd), in this guide the gnatsd address is nats://127.0.0.1:4222
- edit the server code, and run it by passing `--app=nserver --id=yourserverid`
``` go
  type TestHandler struct{}

  func (*TestHandler) World(rc *nserver.ReqContext) {
    rc.Json(rc.Param)
  }

  func main() {
    app := flag.String("app", "nserver", "app identifier")
    id := flag.String("id", "1", "server identifier")
    flag.Parse()

    nats_addr := "nats://127.0.0.1:4222"
    s := nserver.New(*app, *id, nats_addr)
    s.Router("hello", &TestHandler{})
    s.Serving()
  }
```

- then request it
``` go
  func main() {
    nc, _ := nats.Connect("nats://127.0.0.1:4222")
    now := time.Now()
    p := map[string]interface{}{
        "param_a": "value_a",
        "num_1":   1,
        "req_id":  "yourrandomreqid",
    }
    b, err := json.Marshal(p)
    if err != nil {
        fmt.Println(err)
    }
    resp, err := nc.Request("nserver.hello.world", b, 3*time.Second)
    if err != nil {
        fmt.Println(err)
        return
    }
    elapse := time.Now().Sub(now).Seconds()

    fmt.Printf("[%f]receive -- %s\n", elapse, string(resp.Data))
  }
```

- the response should be: 
```
[0.002454]response---{"code":0,"errmsg":"","result":{"num_1":1,"param_a":"value_a","req_id":"yourrandomreqid"}}
```

easy! isn't it?

# internal
1. what's the app and id parameter passing to the server?
```
app is used to specify the micro server name, and it's the prefix in the nats message queue, it's optional and the default value is "nserver"(that means you should request to the queue "nserver.x.xx").
id is a server id that specify a single server instance, it's optional and the default value is 1
```
2, response format?
``` go
  type Response struct {
    Code   int         `json:"code"`
    Errmsg string      `json:"errmsg"`
    Result interface{} `json:"result"`
  }
```
By default, the server would response the data by json.Marshal(Result) if you call rc.Json(data), you can response your own data format by call rc.Write(\[\]byte("the response string")) 

