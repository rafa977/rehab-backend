package handlers

import (
	"net/http"

	gcontext "github.com/gorilla/context"
)

func ValidateCompany(compId uint, r *http.Request) (bool, string) {
	compIDs := gcontext.Get(r, "compIDs").([]uint)
	var isOwner = false

	if len(compIDs) == 0 {
		return false, "Please register your company"
	}

	for _, v := range compIDs {
		if v == compId {
			isOwner = true
		}
	}
	if !isOwner {
		return false, "You don't have access to these data"
	}

	return true, ""

}

func GetRole(r *http.Request) uint {

	roleID := gcontext.Get(r, "roleID").(uint)
	return roleID
}

func GetCompany(r *http.Request) []uint {

	compIDs := gcontext.Get(r, "compIDs").([]uint)
	return compIDs
}

func GetAccount(r *http.Request) uint {

	accountID := gcontext.Get(r, "id").(uint)
	return accountID
}
