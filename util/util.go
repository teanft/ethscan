package util

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"golang.org/x/crypto/sha3"
	"math/big"
	"reflect"
	"regexp"
)

func ToBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	pending := big.NewInt(-1)
	if number.Cmp(pending) == 0 {
		return "pending"
	}
	finalized := big.NewInt(int64(rpc.FinalizedBlockNumber))
	if number.Cmp(finalized) == 0 {
		return "finalized"
	}
	safe := big.NewInt(int64(rpc.SafeBlockNumber))
	if number.Cmp(safe) == 0 {
		return "safe"
	}
	return hexutil.EncodeBig(number)
}

func IsValidAddress(iaddress any) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	switch v := iaddress.(type) {
	case string:
		return re.MatchString(v)
	case common.Address:
		return re.MatchString(v.Hex())
	default:
		return false
	}
}

func IsValidHash(ihash any) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{64}$")
	switch v := ihash.(type) {
	case string:
		return re.MatchString(v)
	case common.Hash:
		return re.MatchString(v.Hex())
	default:
		return false
	}
}

func IsZeroAddress(iaddress any) bool {
	var address common.Address
	switch v := iaddress.(type) {
	case string:
		address = common.HexToAddress(v)
	case common.Address:
		address = v
	default:
		return false
	}

	zeroAddressBytes := common.FromHex("0x0000000000000000000000000000000000000000")
	addressBytes := address.Bytes()
	return reflect.DeepEqual(addressBytes, zeroAddressBytes)
}

func Hex(h []byte) string {
	if len(h) != 0 {
		return ""
	}
	return hexutil.Encode(h[:])
}

func GetHash(hash any) []byte {
	h := sha3.NewLegacyKeccak256()
	h.Write(hash.([]byte))
	var result []byte
	h.Sum(result[:0])

	return result
}
