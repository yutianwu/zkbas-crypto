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

package std

import (
	"errors"
	tedwards "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/std/algebra/twistededwards"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/zecrey-labs/zecrey-crypto/zecrey/twistededwards/tebn254/zecrey"
	"log"
)

type TransferNftProofConstraints struct {
	// commitments
	A_pk Point
	// response
	Z_sk, Z_skInv Variable
	// Commitment Range Proofs
	//GasFeePrimeRangeProof CtRangeProofConstraints
	// common inputs
	Pk                   Point
	TxType               Variable
	NftAccountIndex      Variable
	NftIndex             Variable
	NftContentHash       Variable
	ReceiverAccountIndex Variable
	// gas fee
	A_T_feeC_feeRPrimeInv Point
	Z_bar_r_fee           Variable
	C_fee                 ElGamalEncConstraints
	T_fee                 Point
	GasFeeAssetId         Variable
	GasFee                Variable
	C_fee_DeltaForFrom    ElGamalEncConstraints
	C_fee_DeltaForGas     ElGamalEncConstraints
	IsEnabled             Variable
}

// define tests for verifying the claim nft proof
func (circuit TransferNftProofConstraints) Define(api API) error {
	// first check if C = c_1 \oplus c_2
	// get edwards curve params
	params, err := twistededwards.NewEdCurve(api, tedwards.BN254)
	if err != nil {
		return err
	}
	// verify H
	H := Point{
		X: HX,
		Y: HY,
	}
	// mimc
	hFunc, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}
	tool := NewEccTool(api, params)
	VerifyTransferNftProof(tool, api, &circuit, hFunc, H)
	return nil
}

/*
	VerifyTransferNftProof verify the claim nft proof in circuit
	@api: the constraint system
	@proof: claim nft proof circuit
	@params: params for the curve tebn254
*/
func VerifyTransferNftProof(
	tool *EccTool,
	api API,
	proof *TransferNftProofConstraints,
	hFunc MiMC,
	h Point,
) (c Variable, pkProofs [MaxRangeProofCount]CommonPkProof, tProofs [MaxRangeProofCount]CommonTProof) {
	// check params
	C_fee_DeltaForGas := ElGamalEncConstraints{
		CL: tool.ZeroPoint(),
		CR: tool.ScalarMul(h, proof.GasFee),
	}
	C_fee_DeltaForFrom := ElGamalEncConstraints{
		CL: tool.ZeroPoint(),
		CR: tool.Neg(C_fee_DeltaForGas.CR),
	}
	C_feePrime := tool.EncAdd(proof.C_fee, C_fee_DeltaForFrom)
	C_feePrimeNeg := tool.NegElgamal(C_feePrime)
	hFunc.Write(FixedCurveParam(api))
	WriteEncIntoBuf(&hFunc, proof.C_fee)
	hFunc.Write(proof.GasFeeAssetId)
	hFunc.Write(proof.GasFee)
	WritePointIntoBuf(&hFunc, proof.T_fee)
	WritePointIntoBuf(&hFunc, proof.Pk)
	hFunc.Write(proof.TxType)
	hFunc.Write(proof.NftAccountIndex)
	hFunc.Write(proof.NftIndex)
	hFunc.Write(proof.NftContentHash)
	hFunc.Write(proof.ReceiverAccountIndex)
	WritePointIntoBuf(&hFunc, proof.A_pk)
	WritePointIntoBuf(&hFunc, proof.A_T_feeC_feeRPrimeInv)
	c = hFunc.Sum()
	// Verify balance
	//var l1, r1 Point
	//// verify pk = g^{sk}
	//l1.ScalarMulFixedBase(api, params.BaseX, params.BaseY, proof.Z_sk, params)
	//r1.ScalarMulNonFixedBase(api, &proof.Pk, c, params)
	//r1.AddGeneric(api, &proof.A_pk, &r1, params)
	//IsPointEqual(api, proof.IsEnabled, l1, r1)

	//var l2, r2 Point
	// verify T(C_R - C_R^{\star})^{-1} = (C_L - C_L^{\star})^{-sk^{-1}} g^{\bar{r}}
	//l2 = tool.Add(tool.ScalarBaseMul(proof.Z_bar_r), tool.ScalarMul(CPrimeNeg.CL, proof.Z_skInv))
	//r2 = tool.Add(proof.A_TDivCRprime, tool.ScalarMul(tool.Add(proof.T, CPrimeNeg.CR), c))
	//IsPointEqual(api, proof.IsEnabled, l2, r2)
	// Verify T(C_R - C_R^{\star})^{-1} = (C_L - C_L^{\star})^{-sk^{-1}} g^{\bar{r}}
	//l1 := tool.Add(tool.ScalarBaseMul(proof.Z_bar_r_fee), tool.ScalarMul(C_feePrimeNeg.CL, proof.Z_skInv))
	//r1 := tool.Add(proof.A_T_feeC_feeRPrimeInv, tool.ScalarMul(tool.Add(proof.T_fee, C_feePrimeNeg.CR), c))
	//IsPointEqual(api, proof.IsEnabled, l1, r1)
	// set common parts
	pkProofs[0] = SetPkProof(proof.Pk, proof.A_pk, proof.Z_sk, proof.Z_skInv)
	for i := 1; i < MaxRangeProofCount; i++ {
		pkProofs[i] = pkProofs[0]
	}
	tProofs[0] = SetTProof(C_feePrimeNeg, proof.A_T_feeC_feeRPrimeInv, proof.Z_bar_r_fee, proof.T_fee)
	for i := 1; i < MaxRangeProofCount; i++ {
		tProofs[i] = tProofs[0]
	}
	// set proof deltas
	proof.C_fee_DeltaForGas = C_fee_DeltaForGas
	proof.C_fee_DeltaForFrom = C_fee_DeltaForFrom
	// set proof deltas
	return c, pkProofs, tProofs
}

func SetEmptyTransferNftProofWitness() (witness TransferNftProofConstraints) {
	// commitments
	witness.A_pk, _ = SetPointWitness(BasePoint)
	// response
	witness.Z_sk = ZeroInt
	witness.Z_skInv = ZeroInt
	// common inputs
	witness.Pk, _ = SetPointWitness(BasePoint)
	witness.TxType = ZeroInt
	witness.NftAccountIndex = ZeroInt
	witness.NftIndex = ZeroInt
	witness.NftContentHash = ZeroInt
	witness.ReceiverAccountIndex = ZeroInt
	// gas fee
	witness.A_T_feeC_feeRPrimeInv, _ = SetPointWitness(BasePoint)
	witness.Z_bar_r_fee = ZeroInt
	witness.C_fee, _ = SetElGamalEncWitness(ZeroElgamalEnc)
	witness.T_fee, _ = SetPointWitness(BasePoint)
	witness.GasFeeAssetId = ZeroInt
	witness.GasFee = ZeroInt
	witness.C_fee_DeltaForFrom, _ = SetElGamalEncWitness(ZeroElgamalEnc)
	witness.C_fee_DeltaForGas, _ = SetElGamalEncWitness(ZeroElgamalEnc)
	witness.IsEnabled = SetBoolWitness(false)
	return witness
}

// set the witness for withdraw proof
func SetTransferNftProofWitness(proof *zecrey.TransferNftProof, isEnabled bool) (witness TransferNftProofConstraints, err error) {
	if proof == nil {
		log.Println("[SetTransferNftProofWitness] invalid params")
		return witness, err
	}

	// proof must be correct
	verifyRes, err := proof.Verify()
	if err != nil {
		log.Println("[SetTransferNftProofWitness] invalid proof:", err)
		return witness, err
	}
	if !verifyRes {
		log.Println("[SetTransferNftProofWitness] invalid proof")
		return witness, errors.New("[SetTransferNftProofWitness] invalid proof")
	}
	// commitments
	witness.A_pk, err = SetPointWitness(proof.A_pk)
	if err != nil {
		return witness, err
	}
	witness.Z_sk = proof.Z_sk
	witness.Z_skInv = proof.Z_skInv
	// common inputs
	witness.Pk, err = SetPointWitness(proof.Pk)
	if err != nil {
		return witness, err
	}
	witness.TxType = uint64(proof.TxType)
	witness.NftAccountIndex = proof.NftAccountIndex
	witness.NftIndex = proof.NftIndex
	witness.NftContentHash = proof.NftContentHash
	witness.ReceiverAccountIndex = proof.ReceiverAccountIndex
	// gas fee
	witness.A_T_feeC_feeRPrimeInv, err = SetPointWitness(proof.A_T_feeC_feeRPrimeInv)
	if err != nil {
		return witness, err
	}
	witness.Z_bar_r_fee = proof.Z_bar_r_fee
	witness.C_fee, err = SetElGamalEncWitness(proof.C_fee)
	if err != nil {
		return witness, err
	}
	witness.T_fee, err = SetPointWitness(proof.T_fee)
	if err != nil {
		return witness, err
	}
	witness.GasFeeAssetId = uint64(proof.GasFeeAssetId)
	witness.GasFee = proof.GasFee
	//witness.BPrimeRangeProof, err = SetCtRangeProofWitness(proof.BPrimeRangeProof, isEnabled)
	//if err != nil {
	//	return witness, err
	//}
	// common inputs
	witness.C_fee_DeltaForFrom, _ = SetElGamalEncWitness(ZeroElgamalEnc)
	witness.C_fee_DeltaForGas, _ = SetElGamalEncWitness(ZeroElgamalEnc)
	witness.IsEnabled = SetBoolWitness(isEnabled)
	return witness, nil
}

/*
	VerifyTransferNftTxParams:
	accounts order is:
	- FromAccount
		- Assets
			- AssetGas
		- Nft
			- nft index
	- ToAccount
		- Nft
			- nft index
	- GasAccount
		- Assets
			- AssetGas
*/
func VerifyTransferNftTxParams(api API, flag Variable, nilHash Variable, tx TransferNftProofConstraints, accountsBefore, accountsAfter [NbAccountsPerTx]AccountConstraints) {
	// verify params
	// nft index
	IsVariableEqual(api, flag, tx.NftAccountIndex, accountsBefore[0].NftInfo.NftAccountIndex)
	IsVariableEqual(api, flag, tx.NftIndex, accountsBefore[0].NftInfo.NftIndex)
	IsVariableEqual(api, flag, tx.NftIndex, accountsAfter[1].NftInfo.NftIndex)
	// before account nft should be empty
	IsVariableEqual(api, flag, accountsBefore[0].NftInfo.NftContentHash, nilHash)
	IsVariableEqual(api, flag, accountsBefore[0].NftInfo.AssetId, DefaultInt)
	IsVariableEqual(api, flag, accountsBefore[0].NftInfo.AssetAmount, DefaultInt)
	// gas asset id
	IsVariableEqual(api, flag, tx.GasFeeAssetId, accountsBefore[0].AssetsInfo[1].AssetId)
}
