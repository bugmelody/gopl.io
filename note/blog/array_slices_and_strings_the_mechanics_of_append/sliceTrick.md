见: https://github.com/golang/go/wiki/SliceTricks

#### AppendVector
```go
a = append(a, b...)
```

#### Copy
```go
b = make([]T, len(a))
copy(b, a)
// or, if a is not the empty slice,
b = append([]T(nil), a...)
```

#### Cut

将 a[i:j] 删除掉 

```go
a = append(a[:i], a[j:]...)
```

#### Delete

删除 a[i]

```go
a = append(a[:i], a[i+1:]...)
// or
a = a[:i+copy(a[i:], a[i+1:])]
```

#### Delete without preserving order
```go
a[i] = a[len(a)-1] // 将i设置为a的最后一个元素
a = a[:len(a)-1] // reslice,不包括最后的元素

```
**NOTE** If the type of the element is a _pointer_ or a struct with pointer fields, which need to be garbage collected, the above implementations of ` Cut ` and ` Delete ` have a potential _memory leak_ problem: some elements with values are still referenced by slice ` a ` and thus can not be collected. The following code can fix this problem:
> **Cut**
```go
copy(a[i:], a[j:])
for k, n := len(a)-j+i, len(a); k < n; k++ {
	a[k] = nil // or the zero value of T // 将不会再引用到的元素设置为nil, 设置成nil后,就可以被回收
}
a = a[:len(a)-j+i]
```

> **Delete**
```go
copy(a[i:], a[i+1:])
a[len(a)-1] = nil // or the zero value of T // 将不会再引用到的元素设置为nil, 设置成nil后,就可以被回收
a = a[:len(a)-1]
```

> **Delete without preserving order**
```go
a[i] = a[len(a)-1]
a[len(a)-1] = nil // 之后可以被垃圾回收
a = a[:len(a)-1]
```

#### Expand
```go
a = append(a[:i], append(make([]T, j), a[i:]...)...)
```

#### Extend
```go
a = append(a, make([]T, j)...)
```

#### Insert

将 x 插入原来的 i 的位置上

```go
a = append(a[:i], append([]T{x}, a[i:]...)...)
```
**NOTE** The second ` append ` creates a new slice with its own underlying storage and  copies elements in ` a[i:] ` to that slice, and these elements are then copied back to slice ` a ` (by the first ` append `). The creation of the new slice (and thus memory garbage) and the second copy can be avoided by using an alternative way:
> **Insert**
```go
s = append(s, 0)
copy(s[i+1:], s[i:])
s[i] = x
```

#### InsertVector

b: abc
a: 0123456789
i:     i
insert之后
a: 0123abc456789


```go
a = append(a[:i], append(b, a[i:]...)...)
```

#### Pop
```go
x, a = a[len(a)-1], a[:len(a)-1]
```

#### Push
```go
a = append(a, x)
```

#### Shift
```go
x, a := a[0], a[1:]
```

#### Unshift
```go
a = append([]T{x}, a...)
```

## Additional Tricks
### Filtering without allocating

This trick uses the fact that a slice shares the same backing array and capacity as the original, so the storage is reused for the filtered slice. Of course, the original contents are modified.

```go
b := a[:0] // b和a共享同一个底层数组
for _, x := range a {
	if f(x) { // 如果x满足f
		b = append(b, x) // filter的过程中没有内存分配
	}
}
```

### Reversing

To replace the contents of a slice with the same elements but in reverse order:
```go
for i := len(a)/2-1; i >= 0; i-- {
	opp := len(a)-1-i
	a[i], a[opp] = a[opp], a[i]
}
```
The same thing, except with two indices:
```go
for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
	a[left], a[right] = a[right], a[left]
}
```