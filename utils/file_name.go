package utils

import (
	st "strings"
)

type FileName string

func (f FileName) HasExtension(exts ...string) (_ bool) {
	var i int
	for ; i < len(exts); i++ {
		if st.HasSuffix(string(f), exts[i]) {
			return true
		}
	}
	return
}
