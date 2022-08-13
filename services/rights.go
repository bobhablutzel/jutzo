package main

import (
	"strings"
)

type UserRights struct {
	Grants []string `json:"grants"`
}

func revoke(rights *UserRights, right string) {
	var results []string
	for _, item := range rights.Grants {
		if item != right {
			results = append(results, item)
		}
	}
	rights.Grants = results
}

func grant(rights *UserRights, right string) {
	rights.Grants = append(rights.Grants, right)
}

func hasRight(rights *UserRights, right string) bool {
	for _, item := range rights.Grants {
		if item == right {
			return true
		}
	}
	return false
}

const ADMIN = "admin"
const LOGIN = "login"

func NewRights(rightString string) *UserRights {
	// Break into the Rights array and return the result
	return &UserRights{strings.Split(rightString, ",")}
}

func HasAdmin(rights *UserRights) bool { return hasRight(rights, ADMIN) }
func GrantAdmin(rights *UserRights)    { grant(rights, ADMIN) }
func RevokeAdmin(rights *UserRights)   { revoke(rights, ADMIN) }

func HasLogin(rights *UserRights) bool { return hasRight(rights, LOGIN) }
func GrantLogin(rights *UserRights)    { grant(rights, LOGIN) }
func RevokeLogin(rights *UserRights)   { revoke(rights, LOGIN) }
