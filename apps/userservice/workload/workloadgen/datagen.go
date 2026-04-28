package workloadgen

import (
	"math/rand/v2"

	"github.com/vaastav/raglan/apps/userservice/workflow"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.IntN(len(letters))]
	}
	return string(b)
}

var na_list []string
var eu_list []string
var sa_list []string
var as_list []string
var af_list []string
var oc_list []string

func initialize_arrays() {
	for k := range workflow.AFCountries {
		af_list = append(af_list, k)
	}
	for k := range workflow.ASCountries {
		as_list = append(as_list, k)
	}
	for k := range workflow.OCCountries {
		oc_list = append(oc_list, k)
	}
	for k := range workflow.NACountries {
		na_list = append(na_list, k)
	}
	for k := range workflow.SACountries {
		sa_list = append(sa_list, k)
	}
	for k := range workflow.EUCountries {
		eu_list = append(eu_list, k)
	}
}

func gen_random_country(region string) string {
	if region == "NA" {
		v := rand.IntN(len(na_list))
		return na_list[v]
	} else if region == "SA" {
		v := rand.IntN(len(sa_list))
		return sa_list[v]
	} else if region == "EU" {
		v := rand.IntN(len(eu_list))
		return eu_list[v]
	} else if region == "AS" {
		v := rand.IntN(len(as_list))
		return as_list[v]
	} else if region == "AF" {
		v := rand.IntN(len(af_list))
		return af_list[v]
	} else if region == "OC" {
		v := rand.IntN(len(oc_list))
		return oc_list[v]
	}
	return ""
}

func gen_user_data() (string, string, string, string, string) {
	fname := randSeq(10)
	lname := randSeq(10)
	password := randSeq(16)
	email := randSeq(10) + "@blueprint.github.io"
	address := randSeq(40)
	return fname, lname, password, email, address
}
