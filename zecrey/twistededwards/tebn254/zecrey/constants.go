/*
 * Copyright © 2021 Zecrey Protocol
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package zecrey

import (
	"math/big"
	curve "zecrey-crypto/ecc/ztwistededwards/tebn254"
	"zecrey-crypto/elgamal/twistededwards/tebn254/twistedElgamal"
	"zecrey-crypto/rangeProofs/twistededwards/tebn254/ctrange"
)

type (
	ElGamalEnc = twistedElgamal.ElGamalEnc
	Point      = curve.Point
	RangeProof = ctrange.RangeProof
)

const (
	RangeMaxBits          = ctrange.RangeMaxBits // max bits
	PointSize             = curve.PointSize
	RangeProofSize        = ctrange.RangeProofSize
	WithdrawProofSize     = 10*PointSize + RangeProofSize + 2*EightBytes + AddressSize
	OneMillion            = 1000000
	OneThousand           = 1000
	FourBytes             = 4
	EightBytes            = 8
	TransferSubProofCount = 3
	TransferSubProofSize  = 24*PointSize + RangeProofSize
	TransferProofSize     = TransferSubProofCount*TransferSubProofSize + 6*PointSize + 1*EightBytes

	SwapProofSize            = 35*PointSize + 2*RangeProofSize + 7*EightBytes + 1*FourBytes
	AddLiquidityProofSize    = 32*PointSize + 5*EightBytes + 2*RangeProofSize
	RemoveLiquidityProofSize = 33*PointSize + 6*EightBytes + 1*RangeProofSize
	UnlockProofSize          = 3*PointSize + 2*FourBytes + 2*EightBytes

	AddressSize = 20

	ErrCode = -1
)

var (
	G           = curve.G
	H           = curve.H
	Order       = curve.Order
	MaxRange    = 1099511627775 // 2^{40} - 1
	MaxRangeNeg = -1099511627776
	FixedCurve  = new(big.Int).SetBytes([]byte("ZecreyBN254"))
)
