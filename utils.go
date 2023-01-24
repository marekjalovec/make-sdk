package makesdk

import (
	"math"
	"net/url"
)

func GetMaxAndLimit(maxItems int) (int, int) {
	if maxItems < 0 {
		maxItems = math.MaxInt
	}

	return maxItems, int(math.Min(float64(defaultPageSize), float64(maxItems)))
}

func ColumnsToParams(params *url.Values, columns []string) {
	for _, c := range columns {
		params.Add("cols", c)
	}
}
