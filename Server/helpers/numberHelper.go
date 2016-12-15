package helpers

import (
	"strconv"
	"strings"

	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/log"
	"github.com/dustin/go-humanize"
)

func ParceNumberAndLocalize(source string) string {
	var str = strings.TrimSpace(source)
	var value, err = strconv.ParseInt(str, 10, 64)
	if err != nil {
		lspLogger.Errorf("error: %v", err)
		return source
	}
	return humanize.Comma(value)
}
