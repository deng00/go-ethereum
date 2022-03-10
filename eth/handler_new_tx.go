package eth

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"net/http"
	"strings"
	"time"
)

func (h *handler) AsyncPushTransactions(txs types.Transactions) {
	for _, tx := range txs {
		log.Info("new tx", "hash", tx.Hash().String())
		go h.pushTxToGateway(tx)
	}
}

type MsgTx struct {
	ChainId int64       `json:"chain_id"`
	Height  int64       `json:"height"`
	Time    time.Time   `json:"time"`
	Tx      interface{} `json:"tx"`
}

func (tx *MsgTx) FromJSON(msg string) error {
	return json.Unmarshal([]byte(msg), tx)
}

func (tx *MsgTx) ToJSON() string {
	str, _ := json.Marshal(tx)
	return string(str)
}

func (h *handler) pushTxToGateway(tx *types.Transaction) {
	msg := &MsgTx{
		ChainId: 1,
		Height:  0, // 0 means mem-pool tx
		Time:    time.Now(),
		Tx:      tx,
	}
	resp, _ := http.Post("http://sodium-test.coinsummer.com/api/mempool-collector-eth/v1/msg/123456", "application/json", strings.NewReader(msg.ToJSON()))
	defer resp.Body.Close()
}
