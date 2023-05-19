package main

import (
	"flag"
	"os"
)

var (
	width         int
	height        int
	input         string
	output        string
	pngNameMark   string
	decode        bool
	imageDataSize int
)

func init() {
	flag.IntVar(&width, "w", 512, "generate image width")
	flag.IntVar(&height, "h", 512, "generate image height")
	flag.StringVar(&input, "i", "", "input")
	flag.StringVar(&output, "o", "output", "output dir")
	flag.StringVar(&pngNameMark, "n", "out%d.png", "png名称规则，%d为数字占位符")
	flag.BoolVar(&decode, "d", false, "decode")
	flag.Parse()
}
func main() {
	if decode {
		if err := uncompress(input, output); err != nil {
			panic(err)
		}
	} else {
		if err := os.MkdirAll(output, 0755); err != nil {
			panic(err)
		}
		if err := file2image(input, output); err != nil {
			panic(err)
		}
	}

}
