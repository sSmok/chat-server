package prettier

import (
	"fmt"
	"strconv"
	"strings"
)

// PlaceholderDollar - символ, который используется текущей БД при подставлении аргументов запроса
const PlaceholderDollar = "$"

// Prettier - функция, форматирующая запрос в БД под формат логов
func Prettier(query string, placeholder string, args ...interface{}) string {
	for i, arg := range args {
		var value string
		switch v := arg.(type) {
		case string:
			value = fmt.Sprintf("%q", v)
		case []byte:
			value = fmt.Sprintf("%q", string(v))
		default:
			value = fmt.Sprintf("%v", v)
		}
		query = strings.Replace(query, fmt.Sprintf("%s%s", placeholder, strconv.Itoa(i+1)), value, -1)
	}
	query = strings.ReplaceAll(query, "\t", "")
	query = strings.ReplaceAll(query, "\n", " ")
	return strings.TrimSpace(query)
}
