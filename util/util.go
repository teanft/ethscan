package util

import "strconv"

func ParseHexStrToInt(hexStr string) (int64, error) {
	stripped := hexStr[2:]
	blockNumber, err := strconv.ParseInt(stripped, 16, 32)
	if err != nil {
		return 0, err
	}
	return blockNumber, nil
}
