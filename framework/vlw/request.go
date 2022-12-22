package vlw

import (
	"fmt"
	"os"
	"time"

	"github.com/imroc/req/v3"
)

type MessageResp struct {
	Code      int    `json:"Code"`
	Result    string `json:"Result"`
	ReturnStr string `json:"ReturnStr"`
	ReturnInt string `json:"ReturnInt"`
}

type Client struct {
	*req.Client
}

func NewRequest() *Client {
	c := req.C()
	c.SetTimeout(10 * time.Second)
	c.SetCommonError(&MessageResp{})
	c.OnAfterResponse(func(client *req.Client, resp *req.Response) error {
		if resp.Err != nil {
			if dump := resp.Dump(); dump != "" {
				resp.Err = fmt.Errorf("%s\nraw content:\n%s", resp.Err.Error(), resp.Dump())
			}
			return nil
		}
		if err, ok := resp.Error().(*MessageResp); ok {
			if err.Code != 0 {
				resp.Err = fmt.Errorf(resp.String())
			}
			return nil
		}
		if !resp.IsSuccess() {
			resp.Err = fmt.Errorf("bad response, raw content:\n%s", resp.Dump())
			return nil
		}
		return nil
	})
	if os.Getenv("DEBUG") == "true" {
		c.DevMode()
	}
	return &Client{c}
}
