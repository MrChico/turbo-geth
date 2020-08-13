package rpctest

import (
	"fmt"
	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/core/state"
	"net/http"
	"time"
)

// bench9 tests eth_getProof
func Bench9(needCompare bool) {
	var client = &http.Client{
		Timeout: time.Second * 600,
	}

	var res CallResult
	reqGen := &RequestGenerator{
		client: client,
	}

	reqGen.reqID++
	var blockNumber EthBlockNumber
	res = reqGen.TurboGeth("eth_blockNumber", reqGen.blockNumber(), &blockNumber)
	if res.Err != nil {
		fmt.Printf("Could not get block number: %v\n", res.Err)
		return
	}
	if blockNumber.Error != nil {
		fmt.Printf("Error getting block number: %d %s\n", blockNumber.Error.Code, blockNumber.Error.Message)
		return
	}
	lastBlock := blockNumber.Number
	fmt.Printf("Last block: %d\n", lastBlock)
	// Go back 256 blocks
	bn := int(lastBlock) - 256
	page := common.Hash{}.Bytes()

	for len(page) > 0 {
		accRangeTG := make(map[common.Address]state.DumpAccount)
		var sr DebugAccountRange
		reqGen.reqID++
		res = reqGen.TurboGeth("debug_accountRange", reqGen.accountRange(bn, page, 256), &sr)

		if res.Err != nil {
			fmt.Printf("Could not get accountRange (turbo-geth): %v\n", res.Err)
			return
		}

		if sr.Error != nil {
			fmt.Printf("Error getting accountRange (turbo-geth): %d %s\n", sr.Error.Code, sr.Error.Message)
			break
		} else {
			page = sr.Result.Next
			for k, v := range sr.Result.Accounts {
				accRangeTG[k] = v
			}
		}
		for address, dumpAcc := range accRangeTG {
			var proof EthGetProof
			reqGen.reqID++
			var storageList []common.Hash
			// And now with the storage, if present
			if len(dumpAcc.Storage) > 0 {
				for key := range dumpAcc.Storage {
					storageList = append(storageList, common.HexToHash(key))
					if len(storageList) > 100 {
						break
					}
				}
			}
			res = reqGen.TurboGeth("eth_getProof", reqGen.getProof(bn, address, storageList), &proof)
			if res.Err != nil {
				fmt.Printf("Could not get getProof (turbo-geth): %v\n", res.Err)
				return
			}
			if proof.Error != nil {
				fmt.Printf("Error getting getProof (turbo-geth): %d %s\n", proof.Error.Code, proof.Error.Message)
				break
			}
			if needCompare {
				var gethProof EthGetProof
				reqGen.reqID++
				res = reqGen.Geth("eth_getProof", reqGen.getProof(bn, address, storageList), &gethProof)
				if res.Err != nil {
					fmt.Printf("Could not get getProof (geth): %v\n", res.Err)
					return
				}
				if gethProof.Error != nil {
					fmt.Printf("Error getting getProof (geth): %d %s\n", gethProof.Error.Code, gethProof.Error.Message)
					break
				}
				if !compareProofs(&proof, &gethProof) {
					fmt.Printf("Proofs are different\n")
					break
				}
			}
		}
	}
}