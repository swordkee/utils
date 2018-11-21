package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"hash/crc32"
	"strings"
)

// convert like this: "HelloWorld" to "hello_world"
func SnakeCasedName(s string) string {
	newStr := make([]rune, 0)
	firstTime := true

	for _, chr := range string(s) {
		if isUpper := 'A' <= chr && chr <= 'Z'; isUpper {
			if firstTime == true {
				firstTime = false
			} else {
				newStr = append(newStr, '_')
			}
			chr -= 'A' - 'a'
		}
		newStr = append(newStr, chr)
	}

	return string(newStr)
}

// convert like this: "hello_world" to "HelloWorld"
func TitleCasedName(s string) string {
	newStr := make([]rune, 0)
	upNextChar := true

	for _, chr := range string(s) {
		switch {
		case upNextChar:
			upNextChar = false
			chr -= 'a' - 'A'
		case chr == '_':
			upNextChar = true
			continue
		}

		newStr = append(newStr, chr)
	}

	return string(newStr)
}

func PluralizeString(s string) string {
	str := string(s)
	if strings.HasSuffix(str, "y") {
		str = str[:len(str)-1] + "ie"
	}
	return str + "s"
}

// md5()
func Md5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

// sha1()
func Sha1(str string) string {
	hash := sha1.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

// crc32()
func Crc32(str string) uint32 {
	return crc32.ChecksumIEEE([]byte(str))
}

// str_replace()
func StrReplace(search, replace, subject string, count int) string {
	return strings.Replace(subject, search, replace, count)
}
func StringBuilder(str []string, cap int) string {
	var b strings.Builder
	l := len(str)
	b.Grow(cap)
	for i := 0; i < l; i++ {
		b.WriteString(str[i])
	}
	return b.String()
}
func Substring(source string, start int, end int) string {
	var r = []rune(source)
	length := len(r)

	if start < 0 || end > length || start > end {
		return ""
	}

	if start == 0 && end == length {
		return source
	}

	return string(r[start:end])
}

// trim()
func Trim(str string, characterMask ...string) string {
	mask := ""
	if len(characterMask) == 0 {
		mask = " \\t\\n\\r\\0\\x0B"
	} else {
		mask = characterMask[0]
	}
	return strings.Trim(str, mask)
}

// ltrim()
func Ltrim(str string, characterMask ...string) string {
	mask := ""
	if len(characterMask) == 0 {
		mask = " \\t\\n\\r\\0\\x0B"
	} else {
		mask = characterMask[0]
	}
	return strings.TrimLeft(str, mask)
}

// rtrim()
func Rtrim(str string, characterMask ...string) string {
	mask := ""
	if len(characterMask) == 0 {
		mask = " \\t\\n\\r\\0\\x0B"
	} else {
		mask = characterMask[0]
	}
	return strings.TrimRight(str, mask)
}
