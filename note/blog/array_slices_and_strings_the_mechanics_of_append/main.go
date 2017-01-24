package main

import (
	"bytes"
	"fmt"
)

type path []byte

// 实际是更新了 p 指向的 sliceHeader
func (p *path) TruncateAtFinalSlash() {
	i := bytes.LastIndex(*p, []byte("/"))
	if i >= 0 {
		*p = (*p)[0:i] // 修改 sliceHeader
	}
}

// 什么也没有变, 因为此时在 TruncateAtFinalSlash 内部修改的是 p 的一个拷贝,不会影响原始的p值
func (p path) TruncateAtFinalSlash2() {
	i := bytes.LastIndex(p, []byte("/"))
	if i >= 0 {
		p = p[0:i]
	}
}

// the method could be a value because the value receiver will still point to the same underlying array.
func (p path) ToUpper() {
	for i, b := range p {
		if 'a' <= b && b <= 'z' {
			p[i] = b + 'A' - 'a'
		}
	}
}

// [Exercise: Convert the ToUpper method to use a pointer receiver and see if its behavior changes.]
// 与 ToUpper 相比没有变化,因为还是同一个 sliceHeader 并且指向底层数组的同一块区域
func (p *path) ToUpper2() {
	for i, b := range *p {
		if 'a' <= b && b <= 'z' {
			// 注意: 这里不能写为 *p[i] = b + 'A' - 'a', 会报错:
			// *p[i] = b + 'A' - 'a' // invalid operation: p[i] (type *path does not support indexing)
			(*p)[i] = b + 'A' - 'a'
		}
	}
}

// [Advanced exercise: Convert the ToUpper method to handle Unicode letters, not just ASCII.]
// ??????????????????? 怎么搞?
func (p *path) ToUpper3() {

}

func AddOneToEachElement(slice []byte) {
	for i := range slice {
		slice[i]++
	}
}

func SubtractOneFromLength(slice []byte) []byte {
	slice = slice[0 : len(slice)-1]
	return slice
}

func PtrSubtractOneFromLength(slicePtr *[]byte) {
	slice := *slicePtr
	*slicePtr = slice[0 : len(slice)-1]
}

func Extend(slice []int, element int) []int {
	n := len(slice)
	slice = slice[0 : n+1]
	slice[n] = element
	return slice
}

func Extend2(slice []int, element int) []int {
	n := len(slice)
	if n == cap(slice) {
		// Slice is full; must grow.
		// We double its size and add 1, so if the size is zero we still grow.
		newSlice := make([]int, len(slice), 2*len(slice)+1)
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[0 : n+1]
	slice[n] = element
	return slice
}

func trySlice() {
	fmt.Println("=====================start trySlice")
	// Create a couple of starter slices.
	slice := []int{1, 2, 3}
	slice2 := []int{55, 66, 77}
	fmt.Println("Start slice: ", slice) // Start slice:  [1 2 3]
	fmt.Println("Start slice2:", slice2) // Start slice2: [55 66 77]

	// Add an item to a slice.
	slice = append(slice, 4)
	fmt.Println("Add one item:", slice) // Add one item: [1 2 3 4]

	// Add one slice to another.
	slice = append(slice, slice2...)
	fmt.Println("Add one slice:", slice) // Add one slice: [1 2 3 4 55 66 77]

	// Make a copy of a slice (of int).
	slice3 := append([]int(nil), slice...)
	fmt.Println("Copy a slice:", slice3) // Copy a slice: [1 2 3 4 55 66 77]

	// Copy a slice to the end of itself.
	fmt.Println("Before append to self:", slice) // Before append to self: [1 2 3 4 55 66 77]
	slice = append(slice, slice...)
	fmt.Println("After append to self:", slice) // After append to self: [1 2 3 4 55 66 77 1 2 3 4 55 66 77]
	fmt.Println("=====================end trySlice")
}

func main() {
	var buffer [256]byte

	slice := buffer[10:20]
	for i := 0; i < len(slice); i++ {
		slice[i] = byte(i)
	}
	fmt.Println("before", slice) // before [0 1 2 3 4 5 6 7 8 9]
	AddOneToEachElement(slice)
	fmt.Println("after", slice) // after [1 2 3 4 5 6 7 8 9 10]

	fmt.Println("Before: len(slice) =", len(slice)) // Before: len(slice) = 10
	newSlice := SubtractOneFromLength(slice)
	fmt.Println("After:  len(slice) =", len(slice))       // After:  len(slice) = 10
	fmt.Println("After:  len(newSlice) =", len(newSlice)) // After:  len(newSlice) = 9

	fmt.Println("Before: len(slice) =", len(slice)) // Before: len(slice) = 10
	PtrSubtractOneFromLength(&slice)
	fmt.Println("After:  len(slice) =", len(slice)) // After:  len(slice) = 9

	pathName := path("/usr/bin/tso") // Conversion from string to path.
	pathName.TruncateAtFinalSlash()
	fmt.Printf("%s\n", pathName) // /usr/bin

	pathName2 := path("/usr/bin/tso") // Conversion from string to path.
	pathName2.TruncateAtFinalSlash2()
	fmt.Printf("%s\n", pathName2) // /usr/bin/tso

	pathName3 := path("/usr/bin/tso")
	pathName3.ToUpper()
	fmt.Printf("%s\n", pathName3) // /USR/BIN/TSO

	pathName4 := path("/usr/bin/tso")
	pathName4.ToUpper2()
	fmt.Printf("%s\n", pathName4) // /USR/BIN/TSO

	// 下面代码会报错
	//var iBuffer [10]int
	//slice2 := iBuffer[0:0]
	//for i := 0; i < 20; i++ {
	//	slice2 = Extend(slice2, i)
	//	fmt.Println(slice2)
	//}

	// ------------------
	// 报错:
	// [0]
	// [0 1]
	// [0 1 2]
	// [0 1 2 3]
	// [0 1 2 3 4]
	// [0 1 2 3 4 5]
	// [0 1 2 3 4 5 6]
	// [0 1 2 3 4 5 6 7]
	// [0 1 2 3 4 5 6 7 8]
	// [0 1 2 3 4 5 6 7 8 9]
	// panic: runtime error: slice bounds out of range
	//	------------------

	slice3 := make([]int, 0, 5)
	for i := 0; i < 10; i++ {
		slice3 = Extend2(slice3, i)
		fmt.Printf("len=%d cap=%d slice=%v\n", len(slice3), cap(slice3), slice3)
		fmt.Println("address of 0th element:", &slice3[0])
	}

	//len=1 cap=5 slice=[0]
	//address of 0th element: 0xc04203fe60
	//len=2 cap=5 slice=[0 1]
	//address of 0th element: 0xc04203fe60
	//len=3 cap=5 slice=[0 1 2]
	//address of 0th element: 0xc04203fe60
	//len=4 cap=5 slice=[0 1 2 3]
	//address of 0th element: 0xc04203fe60
	//len=5 cap=5 slice=[0 1 2 3 4]
	//address of 0th element: 0xc04203fe60
	//len=6 cap=11 slice=[0 1 2 3 4 5]
	//address of 0th element: 0xc04203a180
	//len=7 cap=11 slice=[0 1 2 3 4 5 6]
	//address of 0th element: 0xc04203a180
	//len=8 cap=11 slice=[0 1 2 3 4 5 6 7]
	//address of 0th element: 0xc04203a180
	//len=9 cap=11 slice=[0 1 2 3 4 5 6 7 8]
	//address of 0th element: 0xc04203a180
	//len=10 cap=11 slice=[0 1 2 3 4 5 6 7 8 9]
	//address of 0th element: 0xc04203a180

	trySlice()

}
