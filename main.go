package main

import (
	"dfrobot.com/stumSuper/stum"
	"fmt"
)

// InitSys 系统初始化,构建StuMan对象
func InitSys() *stum.Stm {
	stm := &stum.Stm{}
	stm.LoadAllStu()
	stm.LoadAllRoom()
	return stm
}

func main() {
	stm := InitSys()
	for {
		stm.ShowMenu()
		var choice int
		_, err := fmt.Scanln(&choice)
		if err != nil {
			fmt.Println("输入不正确,请重新选择！")
			continue
		}
		switch choice {
		case 1:
			stm.AddStudent()
		case 2:
			stm.ChangeRoom()
		case 3:
			stm.PrintAllStu()
		case 4:
			stm.PrintAllRoom()
		case 5:
			stm.PrintStuById()
		case 6:
			if r := stm.Quit(); r {
				return
			}
		default:
			fmt.Println("输入错误,请重新输入")
		}
	}
}
