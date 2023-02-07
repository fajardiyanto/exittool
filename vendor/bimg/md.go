package bimg

import "runtime"

func RemoveMetaData(buf []byte, o Options) ([]byte, error) {
	defer runtime.KeepAlive(buf)
	return ProcessRemoveMetadata(buf, o)
}
