package request_test

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
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

type SignResponse struct {
	Success bool `json:"success"`
	Data    Sign `json:"data"`
}
type Sign struct {
	Nonce     string `json:"nonce"`
	Sign      string `json:"sign"`
	Timestamp int64  `json:"timestamp"`
}

func TestDangBeiChatComplete(t *testing.T) {
	fs, _ := os.Open("/Users/fwc/Projects/go/OpenAICompatible/a.json")
	defer fs.Close()
	ctx := context.Background()
	deviceID := "f94745995ac935e16ffa37b2dd449f2b_kM4BW_OO_8LY8zJKgUXE"
	reqBody, _ := io.ReadAll(fs)
	var data1 map[string]interface{}
	json.Unmarshal(reqBody, &data1)

	reqBody, _ = json.Marshal(data1)
	fmt.Println("111", data1)
	resp, err := request.NewRequestBuilder("https://ai-dangbei.deno.dev", http.MethodPost).WithJson(bytes.NewReader(reqBody)).Do(ctx, nil)
	if err != nil {
		slog.Error("Failed to sign request", "error", err)
	}
	var signResp SignResponse
	err = json.NewDecoder(resp.Body).Decode(&signResp)
	if err != nil {
		slog.Error("Failed to decode sign response", "error", err)
	}

	fmt.Println(signResp)

	headers := map[string]string{
		"content-type": "application/json",
		"deviceId":     deviceID,
		"nonce":        signResp.Data.Nonce,
		"sign":         signResp.Data.Sign,
		"timestamp":    strconv.FormatInt(signResp.Data.Timestamp, 10),
	}

	resp, _ = request.NewRequestBuilder("https://ai-api.dangbei.net/ai-search/chatApi/v2/chat", http.MethodPost).WithHeaders(headers).WithJson(bytes.NewReader(reqBody)).Do(ctx, nil)
	fmt.Println(resp.StatusCode)
	data, _ := resp.Text()
	fmt.Println(data)

}

func TestHash(t *testing.T) {
	// 1. 读取并解析 JSON 文件
	fs, _ := os.Open("/Users/fwc/Projects/go/OpenAICompatible/a.json")
	defer fs.Close()
	reqBody, _ := io.ReadAll(fs)
	var data interface{}
	json.Unmarshal(reqBody, &data)

	// 2. 使用 Encoder 进行自定义序列化
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	// 关键步骤：禁用 HTML 转义，使其行为与 JS 的 JSON.stringify 一致
	encoder.SetEscapeHTML(false)

	// 3. 执行序列化
	if err := encoder.Encode(&data); err != nil {
		t.Fatalf("JSON encoding error: %v", err)
	}

	// 4. 【重要】移除 Encoder.Encode 添加的末尾换行符
	//    这是确保与 JSON.stringify/json.Marshal 结果完全一致的关键
	canonicalJSONBytes := bytes.TrimRight(buffer.Bytes(), "\n")

	// 5. 计算哈希值
	sum := md5.Sum(canonicalJSONBytes)
	hash := hex.EncodeToString(sum[:])

	// 6. 与JS的哈希值进行比较
	expectedHash := "697bc8709eaee7a01ccf3cf5e0d86d19"
	if hash != expectedHash {
		t.Errorf("Hash mismatch!\nGot: %s\nExpected: %s", hash, expectedHash)
	}
}
