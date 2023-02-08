package useragent

import "strings"

type userAgent struct {
	Product string
	Version string
	Comment string
	userType
}

func (u *userAgent) typeParse() {
	if u.Comment == "" {
		u.CommandLine = true
		return
	}
	u.Browser = true
}

type userType struct {
	Browser     bool
	CommandLine bool
}

func IsCommandLine(ua string) bool {
	return parse(ua).CommandLine
}

func parse(ua string) *userAgent {
	// User Agent Structure:
	uas := &userAgent{}
	if ua == "" {
		uas.CommandLine = true
		return uas
	}
	uaArray := strings.Split(ua, " ")

	// Product parsing
	product := strings.Split(uaArray[0], "/")

	uas.Product = product[0]
	uas.Version = product[1]

	uas.Comment = strings.Join(uaArray[1:], " ")

	// User Type parsing
	uas.typeParse()
	return uas
}
