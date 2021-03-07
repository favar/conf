package reload_err

import (
	"fmt"
	"github.com/favar/conf/core"
)

type fmtError struct {
}

func (f *fmtError) Catch(err error) {
	fmt.Println(err)
}

func FmtReloadError() core.ReloadError {
	return &fmtError{}
}
