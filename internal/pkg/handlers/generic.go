package handlers

import "strconv"

func ConvertStrToUint(variable string) (uint, error) {

	ui64, err := strconv.ParseUint(variable, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(ui64), nil

}
