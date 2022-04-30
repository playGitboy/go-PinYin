package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/xuender/oils/base"
)

func main() {
	file := base.Panic1(os.Open("data/dict.txt"))
	reader := bufio.NewReader(file)
	pySet := base.NewSet[string]()
	pyMap := base.NewMap[rune, string]()

	defer file.Close()

	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		str := string(line)
		runes := []rune(str)
		py := string(runes[1:])
		pyMap[runes[0]] = py
		pySet.Add(strings.Split(py, ",")...)
	}

	pySlice := base.NewSlice(pySet.Slice()...)
	sort.Sort(pySlice)
	createTones(pySlice)
	createDict(pyMap, pySlice)
}

func createDict(pyMap base.Map[rune, string], pySlice base.Slice[string]) {
	tones := base.NewMap[int, base.Map[rune, []int]]()

	for key, value := range pyMap {
		values := strings.Split(value, ",")
		length := len(values)

		if !tones.Has(length) {
			tones[length] = base.NewMap[rune, []int]()
		}

		tone := tones[length]
		tone[key] = make([]int, length)

		for sub, pyStr := range values {
			tone[key][sub] = pySlice.Index(pyStr)
		}
	}

	file := base.Panic1(os.Create("dict_gen.go"))
	defer file.Close()
	_, _ = file.WriteString("// Code generated by dict2data. DO NOT EDIT.\n")
	_, _ = file.WriteString("package py\n\n")

	for num, tone := range tones {
		_, _ = file.WriteString("// nolint\n")
		_, _ = file.WriteString(fmt.Sprintf("var dict%d = map[rune][%d]uint16{\n", num, num))

		for han, pys := range tone {
			s := base.NewSlice(pys...)
			_, _ = file.WriteString(fmt.Sprintf("\t0x%x:{%s},\n", han, s.Join(", ")))
		}

		_, _ = file.WriteString("}\n")
	}
}

func createTones(pySlice base.Slice[string]) {
	tonesGo := base.Panic1(os.Create("tones_gen.go"))
	defer tonesGo.Close()

	_, _ = tonesGo.WriteString("// Code generated by dict2data. DO NOT EDIT.\n")
	_, _ = tonesGo.WriteString("package py\n\n")
	_, _ = tonesGo.WriteString("var tones = []string{")

	var (
		one   byte = ' '
		index      = 0
	)

	for _, pinyin := range pySlice {
		if index%10 == 0 || one != pinyin[0] {
			if one != pinyin[0] {
				_, _ = tonesGo.WriteString("\n")
			}

			_, _ = tonesGo.WriteString("\n\t")
			index = 0
			one = pinyin[0]
		}

		index++

		_, _ = tonesGo.WriteString("\"")
		_, _ = tonesGo.WriteString(pinyin)
		_, _ = tonesGo.WriteString("\", ")
	}

	_, _ = tonesGo.WriteString("\n}")
}
