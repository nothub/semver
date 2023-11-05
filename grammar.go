package semver

func isLetter(r rune) bool {
	return r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z'
}

func isPositiveDigit(r rune) bool {
	return r >= '1' && r <= '9'
}

func isDigit(r rune) bool {
	return r == '0' || isPositiveDigit(r)
}

func isDigits(s string) bool {
	for _, r := range s {
		if !isDigit(r) {
			return false
		}
	}
	return true
}

func isNonDigit(r rune) bool {
	return isLetter(r) || r == '-'
}

func isIdentifierCharacter(r rune) bool {
	return isDigit(r) || isNonDigit(r)
}

func isIdentifierCharacters(s string) bool {
	for _, r := range s {
		if !isIdentifierCharacter(r) {
			return false
		}
	}
	return true
}

func isNumericIdentifier(s string) bool {
	if s == "0" {
		return true
	}
	if !isPositiveDigit([]rune(s)[0]) {
		return false
	}
	return isDigits(s)
}

func isAlphanumericIdentifier(s string) bool {
	/*
	   TODO
	   <alphanumeric identifier> ::= <non-digit>
	                               | <non-digit> <identifier characters>
	                               | <identifier characters> <non-digit>
	                               | <identifier characters> <non-digit> <identifier characters>
	*/
	return false
}

/*
TODO
<build identifier> ::= <alphanumeric identifier>
                     | <digits>
*/

/*
TODO
<pre-release identifier> ::= <alphanumeric identifier>
                           | <numeric identifier>
*/

/*
TODO
<dot-separated build identifiers> ::= <build identifier>
                                    | <build identifier> "." <dot-separated build identifiers>
*/

/*
TODO
<build> ::= <dot-separated build identifiers>
*/

/*
TODO
<dot-separated pre-release identifiers> ::= <pre-release identifier>
                                          | <pre-release identifier> "." <dot-separated pre-release identifiers>
*/

/*
TODO
<pre-release> ::= <dot-separated pre-release identifiers>
*/

/*
TODO
<patch> ::= <numeric identifier>
*/

/*
TODO
<minor> ::= <numeric identifier>
*/

/*
TODO
<major> ::= <numeric identifier>
*/

/*
TODO
<version core> ::= <major> "." <minor> "." <patch>
*/

/*
TODO
<valid semver> ::= <version core>
                 | <version core> "-" <pre-release>
                 | <version core> "+" <build>
                 | <version core> "-" <pre-release> "+" <build>
*/
