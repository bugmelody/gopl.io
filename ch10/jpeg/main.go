// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 287.

//!+main

// The jpeg command reads a PNG image from the standard input
// and writes it as a JPEG image to the standard output.
package main

/**
Notice the blank import of image/png. Without that line, the program compiles and links as
usual but can no longer recognize or decode input in PNG format
 */
import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png" // register PNG decoder
	"io"
	"os"
)

func main() {
	if err := toJPEG(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "jpeg: %v\n", err)
		os.Exit(1)
	}
}

func toJPEG(in io.Reader, out io.Writer) error {
	img, kind, err := image.Decode(in)
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, "Input format =", kind)
	return jpeg.Encode(out, img, &jpeg.Options{Quality: 95})
}

//!-main

/*
//!+with
$ go build gopl.io/ch3/mandelbrot
$ go build gopl.io/ch10/jpeg
$ ./mandelbrot | ./jpeg >mandelbrot.jpg
Input format = png
//!-with

//!+without
$ go build gopl.io/ch10/jpeg
$ ./mandelbrot | ./jpeg >mandelbrot.jpg
jpeg: image: unknown format
//!-without
*/


/**
Here’s how it works. The standard library provides decoders for GIF, PNG, and
JPEG, and users may provide others, but to keep executables small, decoders
are not included in an application unless explicitly requested. The image.Decode
function consults a table of supported formats. Each entry in the table specifies
four things: the name of the format; a string that is a prefix of all images
encoded this way, used to detect the encoding; a function Decode that decodes
an encoded image; and another function DecodeConfig that decodes only the image
metadata, such as its size and color space. An entry is added to the table by
calling image.RegisterFormat, typically from within the package initializer of
the supporting package for each format, like this one in image/png:

package png // image/png

func Decode(r io.Reader) (image.Image, error)
func DecodeConfig(r io.Reader) (image.Config, error)

func init() {
    const pngHeader = "\x89PNG\r\n\x1a\n"
    image.RegisterFormat("png", pngHeader, Decode, DecodeConfig)
}

The effect is that an application need only blank-import the package for the
format it needs to make the image.Decode function able to decode it.
 */