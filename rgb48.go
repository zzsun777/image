// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image

import (
	"image"
	"image/color"
)

var (
	_ image.Image = (*RGB48Image)(nil)
)

type RGB48Image struct {
	Pix    []uint8
	Stride int
	Rect   image.Rectangle
}

func (p *RGB48Image) ColorModel() color.Model { return color.RGBA64Model }

func (p *RGB48Image) Bounds() image.Rectangle { return p.Rect }

func (p *RGB48Image) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.RGBA64{}
	}
	i := p.PixOffset(x, y)
	return color.RGBA64{
		R: uint16(p.Pix[i+0])<<8 | uint16(p.Pix[i+1]),
		G: uint16(p.Pix[i+2])<<8 | uint16(p.Pix[i+3]),
		B: uint16(p.Pix[i+4])<<8 | uint16(p.Pix[i+5]),
		A: 0xffff,
	}
}

func (p *RGB48Image) RGB48At(x, y int) [3]uint16 {
	if !(image.Point{x, y}.In(p.Rect)) {
		return [3]uint16{}
	}
	i := p.PixOffset(x, y)
	return [3]uint16{
		uint16(p.Pix[i+0])<<8 | uint16(p.Pix[i+1]),
		uint16(p.Pix[i+2])<<8 | uint16(p.Pix[i+3]),
		uint16(p.Pix[i+4])<<8 | uint16(p.Pix[i+5]),
	}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *RGB48Image) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*3
}

func (p *RGB48Image) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := color.RGBA64Model.Convert(c).(color.RGBA64)
	p.Pix[i+0] = uint8(c1.R >> 8)
	p.Pix[i+1] = uint8(c1.R)
	p.Pix[i+2] = uint8(c1.G >> 8)
	p.Pix[i+3] = uint8(c1.G)
	p.Pix[i+4] = uint8(c1.B >> 8)
	p.Pix[i+5] = uint8(c1.B)
	return
}

func (p *RGB48Image) SetRGB48(x, y int, c [3]uint16) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	p.Pix[i+0] = uint8(c[0] >> 8)
	p.Pix[i+1] = uint8(c[0])
	p.Pix[i+2] = uint8(c[1] >> 8)
	p.Pix[i+3] = uint8(c[1])
	p.Pix[i+4] = uint8(c[2] >> 8)
	p.Pix[i+5] = uint8(c[2])
	return
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *RGB48Image) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be inside
	// either r1 or r2 if the intersection is empty. Without explicitly checking for
	// this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &RGB48Image{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &RGB48Image{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (p *RGB48Image) Opaque() bool {
	return true
}

// NewRGB48Image returns a new RGB48Image with the given bounds.
func NewRGB48Image(r image.Rectangle) *RGB48Image {
	w, h := r.Dx(), r.Dy()
	pix := make([]uint8, 6*w*h)
	return &RGB48Image{
		Pix:    pix,
		Stride: 6 * w,
		Rect:   r,
	}
}

func NewRGB48ImageFrom(m image.Image) *RGB48Image {
	if m, ok := m.(*RGB48Image); ok {
		return m
	}

	// convert to RGB48Image
	b := m.Bounds()
	rgb := NewRGB48Image(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			pr, pg, pb, _ := m.At(x, y).RGBA()
			rgb.SetRGB48(x, y, [3]uint16{
				uint16(pr),
				uint16(pg),
				uint16(pb),
			})
		}
	}
	return rgb
}
