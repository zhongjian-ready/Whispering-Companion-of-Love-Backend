package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"miniapp-backend/internal/config"
	"net/http"
	"sync"
	"time"
)

type WeChatService struct {
	AppID       string
	AppSecret   string
	accessToken string
	tokenExpiry time.Time
	mu          sync.Mutex
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

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

func (s *WeChatService) GetAccessToken() (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.accessToken != "" && time.Now().Before(s.tokenExpiry) {
		return s.accessToken, nil
	}

	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		s.AppID, s.AppSecret)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result AccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.ErrCode != 0 {
		return "", fmt.Errorf("wechat api error: %d %s", result.ErrCode, result.ErrMsg)
	}

	s.accessToken = result.AccessToken
	// Expire 5 minutes early to be safe
	s.tokenExpiry = time.Now().Add(time.Duration(result.ExpiresIn-300) * time.Second)

	return s.accessToken, nil
}

type SubscribeMessageRequest struct {
	ToUser           string                 `json:"touser"`
	TemplateID       string                 `json:"template_id"`
	Page             string                 `json:"page,omitempty"`
	MiniprogramState string                 `json:"miniprogram_state,omitempty"` // developer, trial, formal
	Lang             string                 `json:"lang,omitempty"`
	Data             map[string]MessageData `json:"data"`
}

type MessageData struct {
	Value string `json:"value"`
}

func (s *WeChatService) SendSubscribeMessage(req *SubscribeMessageRequest) error {
	token, err := s.GetAccessToken()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/subscribe/send?access_token=%s", token)

	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if result.ErrCode != 0 {
		return fmt.Errorf("wechat api error: %d %s", result.ErrCode, result.ErrMsg)
	}

	return nil
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
