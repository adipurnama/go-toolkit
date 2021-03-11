package helper

import (
	"github.com/adipurnama/go-toolkit/errors"
)

var errGuaranteed = errors.New("guarantee error")

// GuaranteedErr ...
func GuaranteedErr() error {
	return errors.WrapFuncMsg(errGuaranteed, "at three.go", "reason", "custom_reason")
}
