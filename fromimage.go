package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/Rehtt/Kit/file/files"
	"image/color"
	"image/png"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func image2file(input, output string) error {
	return nil
}
func uncompress(input string, output string) (err error) {
	var pngList []string
	var index int
	filepath.WalkDir(input, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		n, err := fmt.Sscanf(d.Name(), pngNameMark, &index)
		if err != nil {
			return err
		}
		if n != 0 {
			if index >= len(pngList) {
				pngList = append(pngList, make([]string, index+1-len(pngList))...)
			}
			pngList[index] = path
		}
		return nil
	})

	var filesReader = files.NewReader(pngList)
	filesReader.AfterReadFile(image2Data)
	z, err := gzip.NewReader(filesReader)
	if err != nil {
		return err
	}
	t := tar.NewReader(z)
	for {
		h, err := t.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		path := filepath.Join(output, h.Name)
		err = os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			return err
		}
		w, err := os.Create(path)
		if err != nil {
			return err
		}
		io.Copy(w, t)
		w.Close()
	}
	return
}
func image2Data(imageData io.Reader, out io.Writer) error {
	i, err := png.Decode(imageData)
	if err != nil {
		return err
	}

	bounds := i.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	data := make([]byte, 0, (width*height*8)-8)

	// 读取第一个像素的数据长度
	c := i.At(0, 0).(color.NRGBA64)
	n := uint64(c.R) | (uint64(c.G) << 16) | (uint64(c.B) << 32) | (uint64(c.A) << 48)

	defer func() {
		if err != nil {
			return
		}
		_, err = out.Write(data[:n])
	}()

	index := uint64(0)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			c := i.At(x, y).(color.NRGBA64)
			if x == 0 && y == 0 {
				continue
			}
			data = append(data, uint16ToBytes(c.R, c.G, c.B, c.A)...)
			index += 8
			if index >= n {
				return nil
			}
		}
	}
	return nil
}
func uint16ToBytes(u ...uint16) []byte {
	var tmp bytes.Buffer
	for _, v := range u {
		tmp.WriteByte(uint8(v >> 8))
		tmp.WriteByte(uint8(v))
	}
	return tmp.Bytes()
}
