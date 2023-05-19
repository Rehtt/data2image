package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// 使用png，nrgba64，每像素可以存8byte。每张图片第一个像素存储数据长度
func file2image(filePath string, outPath string) error {
	imageDataSize = (width * height * 8) - 8
	var index int
	err := compress(filePath, func(b []byte) error {
		index += 1
		data := data2Image(b)
		return os.WriteFile(filepath.Join(outPath, fmt.Sprintf(pngNameMark, index)), data, 0644)
	})
	if err != nil {
		return err
	}
	return nil
}
func compress(filePath string, fu func(b []byte) error) error {
	var tmp bytes.Buffer
	var cache = make([]byte, imageDataSize)
	var err error
	var n int
	defer func() {
		if err != nil {
			return
		}
		for tmp.Len() > 0 {
			n, err = tmp.Read(cache)
			if err != nil {
				return
			}
			if err = fu(cache[:n]); err != nil {
				return
			}
		}
	}()

	z := gzip.NewWriter(&tmp)
	defer func() {
		fmt.Println(z.Close())
	}()
	t := tar.NewWriter(z)
	defer func() {
		fmt.Println(t.Close())
	}()

	err = filepath.WalkDir(filePath, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		info, err := f.Stat()
		if err != nil {
			return err
		}
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = strings.ReplaceAll(path, "\\", "/")

		err = t.WriteHeader(header)
		if err != nil {
			return err
		}
		// 分段写入，缓解写入大文件时的内存压力
		var i int64
		var cacheSize = int64(1024 * 1024 * 10) // 10M
		for i < header.Size {
			wn, err := io.CopyN(t, f, cacheSize)
			if err != nil && wn == 0 {
				return err
			}
			i += wn
			if tmp.Len() >= imageDataSize {
				for tmp.Len() >= imageDataSize {
					n, err = tmp.Read(cache)
					if err != nil {
						return err
					}
					if err = fu(cache[:n]); err != nil {
						return err
					}
				}
				if tmp.Len() != 0 {
					n, err = tmp.Read(cache)
					if err != nil {
						return err
					}
					tmp.Reset()
					tmp.Write(cache[:n])
				}
			}
		}

		return err
	})
	if err != nil {
		return err
	}

	return nil
}
func data2Image(data []byte) []byte {
	i := image.NewNRGBA64(image.Rect(0, 0, width, height))
	var index int
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			c := color.NRGBA64{}
			if x == 0 && y == 0 {
				n := len(data)
				c.R = uint16(n & 0xFFFF)
				c.G = uint16((n >> 16) & 0xFFFF)
				c.B = uint16((n >> 32) & 0xFFFF)
				c.A = uint16((n >> 48) & 0xFFFF)
			} else {
				c.R = (uint16(getBytes(data, index)) << 8) | uint16(getBytes(data, index+1))
				c.G = (uint16(getBytes(data, index+2)) << 8) | uint16(getBytes(data, index+3))
				c.B = (uint16(getBytes(data, index+4)) << 8) | uint16(getBytes(data, index+5))
				c.A = (uint16(getBytes(data, index+6)) << 8) | uint16(getBytes(data, index+7))
				index += 8
			}
			i.SetNRGBA64(x, y, c)
		}
	}
	var tmp bytes.Buffer
	png.Encode(&tmp, i)
	return tmp.Bytes()
}
func getBytes(b []byte, index int) byte {
	if index < len(b) {
		return b[index]
	}
	return 0
}
