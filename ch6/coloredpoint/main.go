// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 161.

// Coloredpoint demonstrates struct embedding.
package main

import (
	"fmt"
	"math"
)

//!+decl
import "image/color"

type Point struct{ X, Y float64 }

type ColoredPoint struct {
	Point
	Color color.RGBA
}

//!-decl

func (p Point) Distance(q Point) float64 {
	dX := q.X - p.X
	dY := q.Y - p.Y
	return math.Sqrt(dX*dX + dY*dY)
}

func (p *Point) ScaleBy(factor float64) {
	p.X *= factor
	p.Y *= factor
}

func main() {
	/**
	We can call methods of the embedded Point field using a receiver of type ColoredPoint, even though
	ColoredPoint has no declared methods.
	The methods of Point have been promoted to ColoredPoint. In this way, embedding allows
	complex types with many methods to be built up by the composition of several fields, each providing a few methods.
	*/
	//!+main
	red := color.RGBA{255, 0, 0, 255}
	blue := color.RGBA{0, 0, 255, 255}
	var p = ColoredPoint{Point{1, 1}, red}
	var q = ColoredPoint{Point{5, 4}, blue}
	fmt.Println(p.Distance(q.Point)) // "5"
	p.ScaleBy(2)
	q.ScaleBy(2)
	fmt.Println(p.Distance(q.Point)) // "10"
	//!-main
}

// 书中有解释此错误的原因
/*
//!+error
	p.Distance(q) // compile error: cannot use q (ColoredPoint) as Point
//!-error
*/

func init() {
	//!+methodexpr
	p := Point{1, 2}
	q := Point{4, 6}

	distance := Point.Distance   // method expression
	fmt.Println(distance(p, q))  // "5"
	fmt.Printf("%T\n", distance) // "func(Point, Point) float64"

	scale := (*Point).ScaleBy
	scale(&p, 2)
	fmt.Println(p)            // "{2 4}"
	fmt.Printf("%T\n", scale) // "func(*Point, float64)"
	//!-methodexpr
}

func init() {
	/**
	The type of an anonymous field may be a pointer to a named type, in which case fields and methods are
	promoted indirectly from the pointed-to object. Adding another level of indirection lets us share common
	structures and vary the relationships between objects dynamically. The declaration of ColoredPoint
	below embeds a *Point:
	*/
	red := color.RGBA{255, 0, 0, 255}
	blue := color.RGBA{0, 0, 255, 255}

	//!+indirect
	type ColoredPoint struct {
		*Point
		Color color.RGBA
	}

	p := ColoredPoint{&Point{1, 1}, red}
	q := ColoredPoint{&Point{5, 4}, blue}
	fmt.Println(p.Distance(*q.Point)) // "5"
	// p 和 q 现在指向同一块内存
	q.Point = p.Point                 // p and q now share the same Point
	p.ScaleBy(2)
	// 修改了 p, q 也被改变了
	fmt.Println(*p.Point, *q.Point) // "{2 2} {2 2}"
	//!-indirect
}
