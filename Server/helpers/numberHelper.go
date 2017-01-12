package helpers

import (
	"strconv"
	"strings"

	_ "github.com/WOWLABS/LSPTool/Server/log"
	"github.com/dustin/go-humanize"
)

func ParceNumberAndLocalize(source string) string {
	var str = strings.TrimSpace(source)
	var value, err = strconv.ParseInt(str, 10, 64)
	if err != nil {
		return source
	}
	return humanize.Comma(value)
}
