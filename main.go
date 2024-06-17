package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
)

type FlipperReader struct {
	r              io.Reader
	count, flipIdx int
	hasFlipped     bool
}

func NewFlipperReader(r io.Reader, rSize int) *FlipperReader {
	return &FlipperReader{
		r:       r,
		flipIdx: rand.Int() % rSize,
	}
}

func (f *FlipperReader) Read(p []byte) (int, error) {
	w, err := f.r.Read(p)
	if err != nil && err != io.EOF {
		return w, err
	}

	readStart := f.count
	f.count += w
	if !f.hasFlipped && f.count > f.flipIdx {
		idx := f.flipIdx - readStart
		p[idx] = flipRandBit(p[idx])

		f.hasFlipped = true
	}

	return int(w), err
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("usage: bitflip <in_file> <out_file>")
		os.Exit(1)
	}

	inFilename := os.Args[1]
	inFile, err := os.Open(inFilename)
	if err != nil {
		fmt.Printf("Could not open file \"%s\": %s\n", inFilename, err)
		os.Exit(1)
	}
	defer inFile.Close()

	info, err := inFile.Stat()
	if err != nil {
		fmt.Printf("Could not get file info \"%s\": %s\n", inFilename, err)
		os.Exit(1)
	} else if info.IsDir() {
		fmt.Printf("Could not open file \"%s\": is a directory\n", inFilename)
		os.Exit(1)
	}

	outFilename := os.Args[2]
	outFile, err := os.Create(outFilename)
	if err != nil {
		fmt.Printf("Could not open file \"%s\": %s\n", outFilename, err)
		os.Exit(1)
	}
	defer outFile.Close()

	w, err := io.Copy(outFile, NewFlipperReader(inFile, int(info.Size())))
	if err != nil {
		fmt.Println("Could not copy files")
		os.Exit(1)
	}

	fmt.Printf("Copied %d bytes (and flipped one bit)\n", w)
}

func flipRandBit(b byte) byte {
	n := rand.Int() % 8

	inv := ^b
	flip := byte(1 << n)
	mask := ^flip
	bMask := b & mask
	invFlip := inv & flip
	return bMask | invFlip
}
