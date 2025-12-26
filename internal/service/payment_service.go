package service

import (
	"encoding/xml"
	"fmt"
	"io"
	"miniapp-backend/internal/config"
	"miniapp-backend/internal/model"
	"miniapp-backend/internal/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/pay"
	payConfig "github.com/silenceper/wechat/v2/pay/config"
	"github.com/silenceper/wechat/v2/pay/notify"
	"github.com/silenceper/wechat/v2/pay/order"
)

type PaymentService struct {
	pay       *pay.Pay
	orderRepo *repository.OrderRepository
	userRepo  *repository.UserRepository
	cfg       config.WeChatConfig
}

func NewPaymentService(cfg config.WeChatConfig, orderRepo *repository.OrderRepository, userRepo *repository.UserRepository) *PaymentService {
	wc := wechat.NewWechat()
	wc.SetCache(cache.NewMemory())

	payCfg := &payConfig.Config{
		AppID:     cfg.AppID,
		MchID:     cfg.MchID,
		Key:       cfg.MchKey,
		NotifyURL: cfg.NotifyURL,
	}
	
	payment := wc.GetPay(payCfg)

	return &PaymentService{
		pay:       payment,
		orderRepo: orderRepo,
		userRepo:  userRepo,
		cfg:       cfg,
	}
}

func (s *PaymentService) CreateOrder(userID int64, openID string, planID string, ip string) (*order.Config, error) {
	// 0. Fetch OpenID if not provided
	if openID == "" {
		user, err := s.userRepo.FindByID(userID)
		if err != nil {
			return nil, fmt.Errorf("failed to find user: %v", err)
		}
		if user.OpenID == nil || *user.OpenID == "" {
			return nil, fmt.Errorf("user has no openid")
		}
		openID = *user.OpenID
	}

	// 1. Determine amount based on plan
	var amount int64
	var description string
	
	switch planID {
	case "monthly_plan":
		amount = 990 // 9.90 CNY
		description = "Monthly VIP Subscription"
	case "yearly_plan":
		amount = 9900 // 99.00 CNY
		description = "Yearly VIP Subscription"
	default:
		return nil, fmt.Errorf("invalid plan_id")
	}

	// 2. Create local order record
	orderNo := fmt.Sprintf("%d_%d", time.Now().UnixNano(), userID)
	localOrder := &model.Order{
		UserID:      userID,
		OrderNo:     orderNo,
		Amount:      amount,
		Description: description,
		Status:      "pending",
		PlanID:      planID,
	}

	if err := s.orderRepo.Create(localOrder); err != nil {
		return nil, err
	}

	// 3. Call WeChat BridgeConfig (which calls UnifiedOrder)
	params := &order.Params{
		Body:       description,
		OutTradeNo: orderNo,
		TotalFee:   strconv.FormatInt(amount, 10),
		CreateIP:   ip,
		NotifyURL:  s.cfg.NotifyURL,
		TradeType:  "JSAPI",
		OpenID:     openID,
	}

	payParams, err := s.pay.GetOrder().BridgeConfig(params)
	if err != nil {
		return nil, err
	}

	return &payParams, nil
}

func (s *PaymentService) HandleNotify(req *http.Request, respWriter http.ResponseWriter) error {
	// 1. Parse notification
	var result notify.PaidResult
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	if err := xml.Unmarshal(body, &result); err != nil {
		return err
	}

	// 2. Verify Sign
	if !s.pay.GetNotify().PaidVerifySign(result) {
		return fmt.Errorf("signature verification failed")
	}

	// 3. Check if order exists and status
	if result.OutTradeNo == nil {
		return fmt.Errorf("out_trade_no is missing")
	}
	orderNo := *result.OutTradeNo
	
	order, err := s.orderRepo.FindByOrderNo(orderNo)
	if err != nil {
		return err
	}

	if order.Status == "paid" {
		return nil // Already processed
	}

	// 4. Update Order Status
	if err := s.orderRepo.UpdateStatus(orderNo, "paid"); err != nil {
		return err
	}

	// 5. Update User VIP Status
	// Determine duration based on plan_id from order
	var duration time.Duration
	switch order.PlanID {
	case "monthly_plan":
		duration = 30 * 24 * time.Hour
	case "yearly_plan":
		duration = 365 * 24 * time.Hour
	default:
		duration = 30 * 24 * time.Hour // Default
	}

	user, err := s.userRepo.FindByID(order.UserID)
	if err != nil {
		return err
	}

	now := time.Now()
	var newExpireAt time.Time
	if user.VIPExpireAt != nil && user.VIPExpireAt.After(now) {
		newExpireAt = user.VIPExpireAt.Add(duration)
	} else {
		newExpireAt = now.Add(duration)
	}

	updates := map[string]interface{}{
		"is_vip":        true,
		"vip_expire_at": newExpireAt,
	}
	if err := s.userRepo.UpdateSettings(order.UserID, updates); err != nil {
		return err
	}

	return nil
}

func (s *PaymentService) GetSubscription(userID int64) (bool, string, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return false, "", err
	}

	isVIP := user.IsVIP
	expireDate := ""
	
	if user.VIPExpireAt != nil {
		if user.VIPExpireAt.Before(time.Now()) {
			isVIP = false // Expired
		}
		expireDate = user.VIPExpireAt.Format("2006-01-02")
	}

	return isVIP, expireDate, nil
}
