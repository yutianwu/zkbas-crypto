package legend

import (
	"encoding/json"
	curve "github.com/bnb-chain/zkbas-crypto/ecc/ztwistededwards/tebn254"
	"github.com/bnb-chain/zkbas-crypto/wasm/legend/legendTxTypes"
	"log"
)

func SignCancelOffer(seed string, segmentInfo string) (txInfo string, err error) {
	// parse segmentInfo
	sk, err := curve.GenerateEddsaPrivateKey(seed)
	if err != nil {
		return "", err
	}
	oTxInfo, err := legendTxTypes.ConstructCancelOfferTxInfo(sk, segmentInfo)
	if err != nil {
		return "", err
	}
	txInfoBytes, err := json.Marshal(oTxInfo)
	if err != nil {
		log.Println("unable to marshal:", err)
		return "", err
	}
	return string(txInfoBytes), nil
}
