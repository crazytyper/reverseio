package reverseio

func prepend(dst []byte, src []byte) []byte {
	ldst := len(dst)
	lsrc := len(src)
	l := ldst + lsrc
	if l > cap(dst) {
		tmp := make([]byte, l) // re-allocate
		copy(tmp[lsrc:], dst)
		copy(tmp, src)
		return tmp
	}
	dst = dst[0:l] // make room
	copy(dst[lsrc:], dst[:ldst])
	copy(dst, src)
	return dst
}
