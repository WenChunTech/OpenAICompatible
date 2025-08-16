package request_test

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"testing"
	"time"

	"github.com/WenChunTech/OpenAICompatible/src/request"
	"github.com/WenChunTech/OpenAICompatible/src/util"
)

func TestDangbeiModelList(t *testing.T) {
	timestamp := time.Now().Unix()
	sign := util.Sign(timestamp, "")
	fmt.Println(sign)
	url := "https://ai-api.dangbei.net/ai-search/configApi/v1/getChatModelConfig"
	method := "GET"
	headers := map[string]string{
		"deviceId":  "f94745995ac935e16ffa37b2dd449f2b_kM4BW_OO_8LY8zJKgUXE",
		"nonce":     util.Nonce,
		"sign":      sign,
		"timestamp": strconv.FormatInt(timestamp, 10),
	}

	resp, err := request.NewRequestBuilder(url, method).WithHeaders(headers).Do(context.Background(), nil)
	if err != nil {
		slog.Error(err.Error())
	}

	body, err := resp.Text()
	if err != nil {
		slog.Error(err.Error())
	}
	slog.Info(body)

}
