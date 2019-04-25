package util

/*
	Set correspond bit to 1
*/
func SetBit(n byte, pos uint) byte {
	n |= (1 << pos)
	return n
}

/*
	Clear correspond bit to 0
*/
func ClearBit(n byte, pos uint) byte {
	n = n &^ (1 << pos)
	return n
}

/*
	Check whether the correspond bit equal 1
*/
func TestBit(n byte, pos uint) bool {
	return (((n) & (1 << (pos))) > 0)
}

/*
	Get bit value in the specific bit
*/
func GetVal(val byte, pos uint) byte {
	return (val >> pos) & 1
}
