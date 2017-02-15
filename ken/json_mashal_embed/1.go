/**
json.Marshal 使用 struct 字段:
  如果非匿名, 生成的 json 也是嵌套的
  如果是匿名, 生成的 json 是扁平的
  如果是匿名但是给了tag, 生成的 json 也是嵌套的
*/
package main

import (
	"encoding/json"
	"fmt"
)

type Point struct {
	X int
	Y int
}

type Circle1 struct {
	Point  Point
	Radius int
}

type Circle2 struct {
	Point
	Radius int
}

type Circle3 struct {
	Point  `json:"Point"`
	Radius int
}

func main() {
	p := Point{X: 5, Y: 5}
	c1 := Circle1{
		Point:  p,
		Radius: 1,
	}

	c2 := Circle2{
		Point:  p,
		Radius: 2,
	}

	c3 := Circle3{
		Point:  p,
		Radius: 2,
	}

	c1Json, _ := json.MarshalIndent(c1, "", "  ")
	/**
		{
	  "Point": {
	    "X": 5,
	    "Y": 5
	  },
	  "Radius": 1
	  }
	*/
	c2Json, _ := json.MarshalIndent(c2, "", "  ")
	/**
		{
	  "X": 5,
	  "Y": 5,
	  "Radius": 2
	  }
	*/
	c3Json, _ := json.MarshalIndent(c3, "", "  ")
	/**
	{
	  "Point": {
	    "X": 5,
	    "Y": 5
	  },
	  "Radius": 2
	}
	*/
	fmt.Println(string(c1Json))
	fmt.Println(string(c2Json))
	fmt.Println(string(c3Json))
}
