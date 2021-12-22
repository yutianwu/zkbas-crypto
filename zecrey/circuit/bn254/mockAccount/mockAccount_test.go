package mockAccount

import (
	"fmt"
	"math/big"
	"testing"
	curve "zecrey-crypto/ecc/ztwistededwards/tebn254"
	"zecrey-crypto/elgamal/twistededwards/tebn254/twistedElgamal"
)

func TestMock(t *testing.T) {
	C, _ := twistedElgamal.Enc(big.NewInt(int64(JasonAssetABalance)), curve.RandomValue(), JasonPk)
	fmt.Println(C.String())
	C, _ = twistedElgamal.Enc(big.NewInt(int64(JasonAssetBBalance)), curve.RandomValue(), JasonPk)
	fmt.Println(C.String())
	C, _ = twistedElgamal.Enc(big.NewInt(int64(JasonAssetCBalance)), curve.RandomValue(), JasonPk)
	fmt.Println(C.String())
	C, _ = twistedElgamal.Enc(big.NewInt(int64(JasonAssetDBalance)), curve.RandomValue(), JasonPk)
	fmt.Println(C.String())
	C, _ = twistedElgamal.Enc(big.NewInt(int64(JasonLpBalance)), curve.RandomValue(), JasonPk)
	fmt.Println(C.String())
}
