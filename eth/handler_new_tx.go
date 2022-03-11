package eth

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"net/http"
	"os"
	"strings"
	"time"
)

var gatewayUrl = os.Getenv("CS_MEMPOOL_COLLECTOR_GATEWAY")
var defaultMemPoolMsg *MsgTx

func init() {
	hostname, _ := os.Hostname()
	defaultMemPoolMsg = &MsgTx{
		Hostname: hostname,
		ChainId:  1,
		Height:   0, // 0 means mem-pool tx
	}
}

func (h *handler) AsyncPushTransactions(txs types.Transactions) {
	for _, tx := range txs {
		log.Info("try push new tx to mem-pool collector", "hash", tx.Hash().String())
		go h.pushTxToGateway(tx)
	}
}

type MsgTx struct {
	Hostname string      `json:"hostname"`
	ChainId  int64       `json:"chain_id"`
	Height   int64       `json:"height"`
	Time     time.Time   `json:"time"`
	Tx       interface{} `json:"tx"`
}

func (tx *MsgTx) FromJSON(msg string) error {
	return json.Unmarshal([]byte(msg), tx)
}

func (tx *MsgTx) ToJSON() string {
	str, _ := json.Marshal(tx)
	return string(str)
}

func (h *handler) pushTxToGateway(tx *types.Transaction) {
	if gatewayUrl == "" {
		log.Warn("CS_MEMPOOL_COLLECTOR_GATEWAY not set")
		return
	}
	msg := *defaultMemPoolMsg
	msg.Tx = tx
	msg.Time = time.Now()
	resp, err := http.Post(gatewayUrl, "application/json", strings.NewReader(msg.ToJSON()))
	if err != nil {
		log.Warn("request mem-pool collector gateway failed: " + err.Error())
	}
	defer resp.Body.Close()
}
