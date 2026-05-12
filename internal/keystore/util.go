package keystore

import "bytes"

// bytesReader wraps a byte slice in an io.Reader for age.ParseIdentities.
func bytesReader(b []byte) *bytes.Reader {
	return bytes.NewReader(b)
}
