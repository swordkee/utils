package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"hash/crc32"
	"strconv"
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

// Base64Encode base64_encode()
func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// Base64Decode base64_decode()
func Base64Decode(str string) (string, error) {
	switch len(str) % 4 {
	case 2:
		str += "=="
	case 3:
		str += "="
	}

	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// VersionCompare version_compare()
// The possible operators are: <, lt, <=, le, >, gt, >=, ge, ==, =, eq, !=, <>, ne respectively.
// special version strings these are handled in the following order,
// (any string not found) < dev < alpha = a < beta = b < RC = rc < # < pl = p
// Usage:
// VersionCompare("1.2.3-alpha", "1.2.3RC7", '>=')
// VersionCompare("1.2.3-beta", "1.2.3pl", 'lt')
// VersionCompare("1.1_dev", "1.2any", 'eq')
func VersionCompare(version1, version2, operator string) bool {
	var vcompare func(string, string) int
	var canonicalize func(string) string
	var special func(string, string) int

	// version compare
	vcompare = func(origV1, origV2 string) int {
		if origV1 == "" || origV2 == "" {
			if origV1 == "" && origV2 == "" {
				return 0
			}
			if origV1 == "" {
				return -1
			}
			return 1
		}

		ver1, ver2, compare := "", "", 0
		if origV1[0] == '#' {
			ver1 = origV1
		} else {
			ver1 = canonicalize(origV1)
		}
		if origV2[0] == '#' {
			ver2 = origV2
		} else {
			ver2 = canonicalize(origV2)
		}
		n1, n2 := 0, 0
		for {
			p1, p2 := "", ""
			n1 = strings.IndexByte(ver1, '.')
			if n1 == -1 {
				p1, ver1 = ver1[:], ""
			} else {
				p1, ver1 = ver1[:n1], ver1[n1+1:]
			}
			n2 = strings.IndexByte(ver2, '.')
			if n2 == -1 {
				p2, ver2 = ver2, ""
			} else {
				p2, ver2 = ver2[:n2], ver2[n2+1:]
			}

			if (p1[0] >= '0' && p1[0] <= '9') && (p2[0] >= '0' && p2[0] <= '9') { // all is digit
				l1, _ := strconv.Atoi(p1)
				l2, _ := strconv.Atoi(p2)
				if l1 > l2 {
					compare = 1
				} else if l1 == l2 {
					compare = 0
				} else {
					compare = -1
				}
			} else if !(p1[0] >= '0' && p1[0] <= '9') && !(p2[0] >= '0' && p2[0] <= '9') { // all digit
				compare = special(p1, p2)
			} else { // part is digit
				if p1[0] >= '0' && p1[0] <= '9' { // is digit
					compare = special("#N#", p2)
				} else {
					compare = special(p1, "#N#")
				}
			}

			if compare != 0 || n1 == -1 || n2 == -1 {
				break
			}
		}

		if compare == 0 {
			if ver1 != "" {
				if ver1[0] >= '0' && ver1[0] <= '9' {
					compare = 1
				} else {
					compare = vcompare(ver1, "#N#")
				}
			} else if ver2 != "" {
				if ver2[0] >= '0' && ver2[0] <= '9' {
					compare = -1
				} else {
					compare = vcompare("#N#", ver2)
				}
			}
		}

		return compare
	}

	// canonicalize
	canonicalize = func(version string) string {
		ver := []byte(version)
		l := len(ver)
		if l == 0 {
			return ""
		}
		var buf = make([]byte, l*2)
		j := 0
		for i, v := range ver {
			next := uint8(0)
			if i+1 < l { // Have the next one
				next = ver[i+1]
			}
			if v == '-' || v == '_' || v == '+' { // replace '-', '_', '+' to '.'
				if j > 0 && buf[j-1] != '.' {
					buf[j] = '.'
					j++
				}
			} else if (next > 0) &&
				(!(next >= '0' && next <= '9') && (v >= '0' && v <= '9')) ||
				(!(v >= '0' && v <= '9') && (next >= '0' && next <= '9')) { // Insert '.' before and after a non-digit
				buf[j] = v
				j++
				if v != '.' && next != '.' {
					buf[j] = '.'
					j++
				}
				continue
			} else if !((v >= '0' && v <= '9') ||
				(v >= 'a' && v <= 'z') || (v >= 'A' && v <= 'Z')) { // Non-letters and numbers
				if j > 0 && buf[j-1] != '.' {
					buf[j] = '.'
					j++
				}
			} else {
				buf[j] = v
				j++
			}
		}

		return string(buf[:j])
	}

	// compare special version forms
	special = func(form1, form2 string) int {
		found1, found2, len1, len2 := -1, -1, len(form1), len(form2)
		// (Any string not found) < dev < alpha = a < beta = b < RC = rc < # < pl = p
		forms := map[string]int{
			"dev":   0,
			"alpha": 1,
			"a":     1,
			"beta":  2,
			"b":     2,
			"RC":    3,
			"rc":    3,
			"#":     4,
			"pl":    5,
			"p":     5,
		}

		for name, order := range forms {
			if len1 < len(name) {
				continue
			}
			if strings.Compare(form1[:len(name)], name) == 0 {
				found1 = order
				break
			}
		}
		for name, order := range forms {
			if len2 < len(name) {
				continue
			}
			if strings.Compare(form2[:len(name)], name) == 0 {
				found2 = order
				break
			}
		}

		if found1 == found2 {
			return 0
		} else if found1 > found2 {
			return 1
		} else {
			return -1
		}
	}

	compare := vcompare(version1, version2)

	switch operator {
	case "<", "lt":
		return compare == -1
	case "<=", "le":
		return compare != 1
	case ">", "gt":
		return compare == 1
	case ">=", "ge":
		return compare != -1
	case "==", "=", "eq":
		return compare == 0
	case "!=", "<>", "ne":
		return compare != 0
	default:
		panic("operator: invalid")
	}
}
