package main

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/go-nats"
	"testing"
	"time"
)

var nc *nats.Conn

func init() {
	nc, _ = nats.Connect("nats://127.0.0.1:4222")
}

func req() {
	camp := map[string]string{
		"a":     "b",
		"hello": "world",
	}
	b, err := json.Marshal(camp)
	if err != nil {
		fmt.Println(err)
	}
	_, err = nc.Request("arran.abc.create", b, 3*time.Second)
	if err != nil {
		fmt.Println(err)
	}
}

func Test_single(t *testing.T) {
	now := time.Now()
	camp := map[string]string{
		"a":      "123",
		"hello":  "world",
		"req_id": "abcdefjdffjeorgjverorgjvbsaerg",
	}
	b, err := json.Marshal(camp)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := nc.Request("arran.user.create", b, 3*time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}
	elapse := time.Now().Sub(now).Seconds()

	fmt.Printf("[%f]response---%s\n", elapse, string(resp.Data))
}

func Benchmark_nats(b *testing.B) {
	for i := 0; i < b.N; i++ { //use b.N for looping
		req()
	}
}
