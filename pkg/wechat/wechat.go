package wechat

import (
	"encoding/json"
	"fmt"
	"miniapp-backend/internal/config"
	"net/http"
)

type WeChatService struct {
	AppID     string
	AppSecret string
}

func NewWeChatService(cfg config.WeChatConfig) *WeChatService {
	return &WeChatService{
		AppID:     cfg.AppID,
		AppSecret: cfg.AppSecret,
	}
}

type Code2SessionResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

func (s *WeChatService) Code2Session(code string) (*Code2SessionResponse, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		s.AppID, s.AppSecret, code)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result Code2SessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.ErrCode != 0 {
		return nil, fmt.Errorf("wechat api error: %d %s", result.ErrCode, result.ErrMsg)
	}

	return &result, nil
}
