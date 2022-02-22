package config

import (
	"errors"
	"net"
	"regexp"
	"strconv"
)

func ValidateServerAddress(address string) error {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return err
	}

	if !isPortNumberValid(port) {
		return errors.New("invalid port")
	}

	// empty host is OK
	if host == "" {
		return nil
	}

	if isHostValidIP(host) || isHostnameSomewhatValid(host) {
		return nil
	}

	return errors.New("invalid host")
}

// isHostnameSomewhatValid returns true if hostname is in accordance with RFC 3696
// There is no need to make 100% strict validation, simply cutting off obvious typos
// will do
// Go does not support PCRE (particularly lookaheads and lookbehinds), therefore
// longer and plainer regex is used
func isHostnameSomewhatValid(hostname string) bool {
	// RFC 3696
	// ... the labels (words or strings separated by periods) that make up
	// a domain name must consist of only the ASCII alphabetic and numeric
	// characters, plus the hyphen.
	// No other symbols or punctuation characters are permitted, nor is
	// blank space.  If the hyphen is used, it is not permitted to appear
	// at either the beginning or end of a label.  There is an additional
	// rule that essentially requires that top-level domain names not be
	// all-numeric ...
	// ... A DNS label may be no more than 63 octets long.  This is in the
	// form actually stored; if a non-ASCII label is converted to encoded
	// "punycode" form (see Section 5), the length of that form may restrict
	// the number of actual characters (in the original character set) that
	// can be accommodated.  A complete, fully-qualified, domain name must
	// not exceed 255 octets ...
	const hostnameRegexString = `^([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9]\.)*([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9])$`
	const allNumericTLDRegexString = `\.[0-9]+$`
	hostnameRegex := regexp.MustCompile(hostnameRegexString)
	allNumericTLDRegex := regexp.MustCompile(allNumericTLDRegexString)

	return hostnameRegex.MatchString(hostname) &&
		!allNumericTLDRegex.MatchString(hostname) &&
		len(hostname) <= 255
}

// isHostValidIP returns true for string representations of valid IPv4 and IPv6 addresses
// totally relies upon net.ParseIP()
func isHostValidIP(host string) bool {
	return net.ParseIP(host) != nil
}

// isPortNumberValid returns true for string representations of ports 1..65535, false otherwise
func isPortNumberValid(port string) bool {
	portNum, err := strconv.ParseInt(port, 10, 32)
	if err != nil {
		return false
	}
	if portNum < 1 || portNum > 65535 {
		return false
	}
	return true
}
