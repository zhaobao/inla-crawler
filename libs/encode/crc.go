package encode

import (
	"fmt"
	"hash/crc32"
)

const tablePloy uint32 = 0xD5828281

func CrcEncode(cdn string) string {
	crc32q := crc32.MakeTable(tablePloy)
	return fmt.Sprintf("%08x", crc32.Checksum([]byte(cdn), crc32q))
}
