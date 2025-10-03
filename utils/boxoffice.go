package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type BoxOfficeAPIResponse struct {
	Title       string `json:"title"`
	Distributor string `json:"distributor"`
	ReleaseDate string `json:"releaseDate"`
	Budget      int64  `json:"budget"`
	Revenue     struct {
		Worldwide         int64 `json:"worldwide"`
		OpeningWeekendUsa int64 `json:"openingWeekendUSA"`
	} `json:"revenue"`
	MPARating string `json:"mpaRating"`
}

func FetchBoxOfficeData(title string) (*BoxOfficeAPIResponse, error) {
	boxOfficeURL := os.Getenv("BOXOFFICE_URL")
	apiKey := os.Getenv("BOXOFFICE_API_KEY")

	if boxOfficeURL == "" || apiKey == "" {
		return nil, fmt.Errorf("BOXOFFICE_URL or BOXOFFICE_API_KEY not set")
	}

	// 构建请求URL
	url := fmt.Sprintf("%s/boxoffice?title=%s", boxOfficeURL, title)

	// 创建HTTP客户端，设置超时
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 添加API密钥头
	req.Header.Set("X-API-Key", apiKey)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("box office API returned status code: %d", resp.StatusCode)
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// 解析JSON响应
	var boxOfficeData BoxOfficeAPIResponse
	err = json.Unmarshal(body, &boxOfficeData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &boxOfficeData, nil
}
