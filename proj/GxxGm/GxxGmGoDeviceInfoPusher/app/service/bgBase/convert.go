package bgBase

func Uint8Array_2_String(data []uint8) string {
	ba := []byte{}
	for _, b := range data {
		ba = append(ba, byte(b))
	}
	return string(ba)
}
