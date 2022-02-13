package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
)

//地图块
type Cell struct {
	Flag bool
}

//地图类
type Map struct {
	Cells  [][]*Cell
	width  int64
	height int64
	c      chan string
}

func (m *Map) laodMap(mapPath string) {
	m.c = make(chan string)
	f, err := os.OpenFile(mapPath, os.O_RDWR, 0)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	readBuff22 := make([]byte, 22)
	readBuff2 := make([]byte, 2)
	readBuff1 := make([]byte, 1)
	readBuff13 := make([]byte, 13)

	f.Read(readBuff22)
	//读取2字节
	count, err := f.Read(readBuff2)
	if err != nil {
		fmt.Println(err)
	}
	str := hex.EncodeToString(readBuff2[:count])
	width, err := strconv.ParseInt(str[2:]+str[:2], 16, 0)
	if err != nil {
		fmt.Println(err)
	}
	m.width = width
	fmt.Println(width)
	//读取2字节
	count, err = f.Read(readBuff2)
	if err != nil {
		fmt.Println(err)
	}
	str = hex.EncodeToString(readBuff2[:count])
	height, err := strconv.ParseInt(str[2:]+str[:2], 16, 0)
	m.height = height
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(height)
	m.Cells = make([][]*Cell, width)
	for i := range m.Cells {
		m.Cells[i] = make([]*Cell, height)
	}
	//offset 28字节
	f.Seek(28, 0)
	var i, j int64
	for i = 0; i < width; i++ {
		for j = 0; j < height; j++ {
			m.Cells[i][j] = new(Cell)
		}
	}
	for i = 0; i < width/2; i++ {
		for j = 0; j < height/2; j++ {
			//读取1字节
			f.Read(readBuff1)
			//读取2字节
			f.Read(readBuff2)
		}
	}
	for i = 0; i < width; i++ {
		for j = 0; j < height; j++ {
			//读取1字节
			count, err = f.Read(readBuff1)
			if err != nil {
				fmt.Println(err)
			}
			flag := hex.EncodeToString(readBuff1[:count])
			newflag, err := strconv.ParseInt(flag, 16, 0)
			if err != nil {
				fmt.Println(err)
			}
			//读取13字节
			f.Read(readBuff13)
			//
			par_1, err := strconv.ParseInt("01", 16, 0)
			if err != nil {
				fmt.Println(err)
			}
			par_2, err := strconv.ParseInt("02", 16, 0)
			if err != nil {
				fmt.Println(err)
			}
			var res bool
			if ((newflag & par_1) != 1) || ((newflag & par_2) != 2) {
				res = true
			} else {
				res = false
			}
			m.Cells[i][j].Flag = res
		}
	}
}

func main() {
	maps := new(Map)
	maps.laodMap("11.map")
	handle, err := os.OpenFile("map.txt", os.O_RDONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		fmt.Println(err)
	}
	var i, j int64
	var flg int
	defer handle.Close()
	for i = 0; i < maps.width; i++ {
		for j = 0; j < maps.height; j++ {
			buff := bufio.NewWriter(handle)
			if maps.Cells[i][j].Flag {
				flg = 1
			} else {
				flg = 0
			}
			buff.WriteString(strconv.Itoa(flg))
			if j == maps.height-1 {
				buff.WriteString("\n")
			}
			err = buff.Flush()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
