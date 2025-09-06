package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Dereference func
func Dereference[T any](i *T) T {
	if i != nil {
		return *i
	}

	return *new(T)
}

func WithPagination(page *int64, limit *int64) map[string]int64 {
	pageData := make(map[string]int64, 2)

	if page == nil {
		pageData["page"] = 1
	} else {
		pageData["page"] = *page
	}

	if limit == nil {
		pageData["limit"] = 20
	} else {
		pageData["limit"] = *limit
	}

	return pageData
}

func CalcTotalPage(total int64, limit int64) int64 {
	t := total / limit
	if (total % limit) != 0 {
		t += 1
	}
	return t
}

func ExtractYoutubeVideoID(url string) string {
	re := regexp.MustCompile(`(?:https?:\/\/)?(?:www\.)?(?:youtube\.com\/(?:[^\/\n\s]+\/\S+\/|(?:v|e(?:mbed)?)\/|\S*?[?&]v=)|youtu\.be\/)([a-zA-Z0-9_-]{11})`)
	matches := re.FindStringSubmatch(url)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func NewPointer[T any](v T) *T {
	return &v
}

func CheckKeyMatch(key1, key2 string) (bool, error) {
	key2 = strings.Replace(key2, "/*", "/.*", -1)

	re := regexp.MustCompile(`:[^/]+`)
	key2 = re.ReplaceAllString(key2, "$1[^/]+$2")

	return regexp.MatchString("^"+key2+"$", key1)
}

func ToInt[T any](value T) (int, error) {
	switch v := any(value).(type) {
	case int:
		return v, nil
	case int8:
		return int(v), nil
	case int16:
		return int(v), nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case uint:
		return int(v), nil
	case uint8:
		return int(v), nil
	case uint16:
		return int(v), nil
	case uint32:
		return int(v), nil
	case uint64:
		return int(v), nil
	case float32:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		parsed, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return 0, err
		}
		return int(parsed), nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("unsupported type: %v", reflect.TypeOf(value))
	}
}

func GenerateSlug(str string) string {
	// Initialize the normalizer for removing diacritics
	normalizer := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

	// Normalize and remove diacritics from non-Chinese characters
	var result []rune
	for _, r := range str {
		// Skip normalization for Chinese characters
		if unicode.Is(unicode.Han, r) {
			result = append(result, r)
		} else {
			s, _, _ := transform.String(normalizer, string(r))
			result = append(result, []rune(s)...)
		}
	}
	s := string(result)

	// Convert to lowercase
	s = strings.ToLower(s)

	// Replace Vietnamese characters
	replacements := map[string]string{
		"đ": "d",
		"Đ": "d",
		"ð": "d",
	}

	for old, new := range replacements {
		s = strings.ReplaceAll(s, old, new)
	}

	// Replace spaces with hyphens
	s = strings.ReplaceAll(s, " ", "-")

	// Remove any character that is not a letter, number, hyphen or Chinese character
	reg := regexp.MustCompile(`[^\p{Han}\p{L}\p{N}-]+`)
	s = reg.ReplaceAllString(s, "")

	// Replace multiple hyphens with a single hyphen
	reg = regexp.MustCompile(`-+`)
	s = reg.ReplaceAllString(s, "-")

	// Trim hyphens from the start and end
	s = strings.Trim(s, "-")

	return s
}

func GetStartOfDay() time.Time {
	local, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		return time.Now()
	}

	now := time.Now().In(local)
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, local)
	return startOfDay
}

func IsEmpty[T any](value T) bool {
	return reflect.DeepEqual(value, reflect.Zero(reflect.TypeOf(value)).Interface())
}
