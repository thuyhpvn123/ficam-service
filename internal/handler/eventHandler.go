package handler

import (
	"math/big"
	"time"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"encoding/hex"
	"encoding/json"
	"meta-node-ficam/internal/model"
	"meta-node-ficam/utils"
	e_common "github.com/ethereum/go-ethereum/common"
	"github.com/meta-node-blockchain/meta-node/pkg/logger"
	"github.com/meta-node-blockchain/meta-node/types"
)

type EventHandler struct {
	ficamABI *abi.ABI
}


func NewEventHandler(
	ficamABI *abi.ABI,
) *EventHandler {
	return &EventHandler{
		ficamABI: ficamABI,
	}
}
func (h *EventHandler) HandleEvent(events types.EventLogs) {
	for _, event := range events.EventLogList() {
		switch event.Topics()[0] {
		case h.ficamABI.Events["EmailOrder"].ID.String()[2:]:
			h.handleEmailOrder(event.Topics(), event.Data())
		}
	}
}
func (h *EventHandler) handleEmailOrder(topics []string, data string) {
	result := make(map[string]interface{})
	err := h.ficamABI.UnpackIntoMap(result, "EmailOrder", e_common.FromHex(data))
	if err != nil {
		logger.Error("can't unpack to map handleEmailOder", err)
	}
	email := result["email"].(string)
	price := uint(result["totalPrice"].(*big.Int).Uint64())
	order := result["order"]
	jsonData, _ := json.Marshal(order)
	orderContent := model.EmailOrder{}
	json.Unmarshal([]byte(jsonData), &orderContent)
	orderContent.HexID = hex.EncodeToString(orderContent.ID)
	for i, _ := range orderContent.Products {
		orderContent.Products[i].HexIdProduct = hex.EncodeToString(orderContent.Products[i].ID)
	}
	orderContent.CreateAtDate = time.Unix(int64(orderContent.CreateAt), 0).Format("2006-01-02 15:04:05")
	dataEmail := model.Data{}
	dataEmail.Order = orderContent
	dataEmail.PaymentOrder = price
	go func() {
		utils.ReplyEmailOrder(email, dataEmail,"Email Order Comfirmation")
	}()
}
