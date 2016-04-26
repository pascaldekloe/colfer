package testdata

// Fuzz is a test for the generated code.
// See https://github.com/dvyukov/go-fuzz
func Fuzz(data []byte) int {
	o := new(O)
	err := o.UnmarshalBinary(data)
	if err != nil {
		return 0
	}

	_, err = o.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return 1
}
