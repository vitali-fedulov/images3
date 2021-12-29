package images

// yCbCr transforms RGB components to YCbCr. // TODO: Maybe export the var?
func yCbCr(r, g, b float32) (yc, cb, cr float32) {
	yc = 0.299000*r + 0.587000*g + 0.114000*b
	cb = 128 - 0.168736*r - 0.331264*g + 0.500000*b
	cr = 128 + 0.500000*r - 0.418688*g - 0.081312*b
	return yc, cb, cr
}
