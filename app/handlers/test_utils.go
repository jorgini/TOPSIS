package handlers

import (
	"io"
	"strings"
)

func ReadResponse(reader io.ReadCloser) (string, error) {
	resp := make([]byte, 100)
	n, err := reader.Read(resp)
	if n == 0 {
		return "", err
	} else if n == 100 {
		for n != 0 {
			expResp := make([]byte, len(resp)*2)
			n, err = reader.Read(expResp)
			resp = append(resp, expResp...)
		}
	}
	return strings.Trim(string(resp), "\x00"), nil
}
