package stum

import (
	"encoding/json"
	"fmt"
	"os"
)

type Stm struct {
	AllStudent map[int]*Student
	AllRoom    map[int]*Room
}

func (stm *Stm) AddStudent() {
	stu := &Student{}
	if ok := stm.CheckCanAddStu(); !ok {
		fmt.Println("当前所有寝室都已满员！不能再录入学生了")
		return
	}
	fmt.Println("您选择了添加学生")
	if ok := stu.CheckStudentIdRight(); !ok {
		stm.AddStudent()
		return
	}
	if _, ok := stm.CheckIfStuExist(stu.Id); ok {
		fmt.Println("学生Id已存在,请重新输入")
		stm.AddStudent()
		return
	}
	if ok := stu.CheckStudentNameRight(); !ok {
		stm.AddStudent()
		return
	}
	if ok := stu.CheckStudentGenderRight(); !ok {
		stm.AddStudent()
		return
	}
	if ok := stu.CheckStudentAgeRight(); !ok {
		stm.AddStudent()
		return
	}
	if ok := stu.CheckStudentRoomIdRight(); !ok {
		stm.AddStudent()
		return
	}
	room, ok := stm.CheckRoomExist(stu.RoomId)
	if !ok {
		fmt.Println("寝室不存在,请重新录入")
		stm.AddStudent()
		return
	}
	if ok := room.CheckStuGenderRight(stu); !ok {
		stm.AddStudent()
		return
	}
	roomId, ok := stm.AutoDistributeRoom(stu)
	if !ok {
		fmt.Println("寝室分配出现错误!请重新录入")
		stm.AddStudent()
		return
	}
	fmt.Printf("学生添加成功！\n")
	stm.AllRoom[roomId].SettleInStu()
	stu.RoomId = roomId
	stm.AllStudent[stu.Id] = stu
	stm.writeAllStuToFile()
	stm.writeAllRoomToFile()
}

func (stm *Stm) ChangeRoom() {
	fmt.Println("您选择了学生寝室信息调整")
	stu := &Student{}
	if ok := stu.CheckStudentIdRight(); !ok {
		stm.ChangeRoom()
		return
	}
	if student, ok := stm.CheckIfStuExist(stu.Id); ok {
		stu = student
		fmt.Printf("学生id:%d 姓名%s的当前寝室编号为:%d\n", stu.Id, stu.Name, stu.RoomId)
		room := stu.GetRoomInfo(stm)
		if room == nil {
			fmt.Printf("出错了！id为%d的学生不存在寝室信息", stu.Id)
			stm.ChangeRoom()
			return
		}
		if targetRoom, ok := stm.CheckChangeInfo(stu); ok {
			//判断目标寝室是否满员
			if targetRoom.CheckRoomIsFull() {
				fmt.Print("目标寝室目前已满员，是否与该寝室成员互换?y/n:")
				var choice string
				_, err := fmt.Scanln(&choice)
				if err != nil {
					fmt.Println("输入不正确")
					stm.ChangeRoom()
					return
				}
				if choice == "y" || choice == "Y" {
					//列出目前寝室当前成员
					fmt.Println("当前目标寝室成员信息如下：")
					stm.PrintStuInRoom(targetRoom.Id)
					targetStu := &Student{}
					if ok := targetStu.CheckStudentIdRight(); !ok {
						stm.ChangeRoom()
						return
					}
					targetStu, ok := stm.CheckIfStuExist(targetStu.Id)
					if !ok {
						fmt.Println("目标学生不存在！")
						stm.ChangeRoom()
						return
					}
					//判断目标学生是否是在目标寝室
					if targetRoom.Id != targetStu.RoomId {
						fmt.Printf("id为%d的学生不在寝室 %d,请重新操作！\n", targetStu.Id, targetRoom.Id)
						stm.ChangeRoom()
						return
					}
					curRoomId := stu.RoomId
					stm.AllStudent[stu.Id].RoomId = targetRoom.Id
					stm.AllStudent[targetStu.Id].RoomId = curRoomId
					fmt.Printf("学生 %s 与学生 %s 互换寝室成功\n", stu.Name, targetStu.Name)
				} else {
					fmt.Println("请重新操作！")
					stm.ChangeRoom()
					return
				}
			} else {
				room.RemoveStu()
				stu.RoomId = targetRoom.Id
				targetRoom.SettleInStu()
				stm.AllStudent[stu.Id] = stu
				fmt.Printf("学生 %s 已成功迁移至寝室 %d \n", stu.Name, targetRoom.Id)
			}
		} else {
			stm.ChangeRoom()
			return
		}
	} else {
		fmt.Println("学生Id不存在,请重新输入")
		stm.ChangeRoom()
		return
	}
	stm.writeAllRoomToFile()
	stm.writeAllStuToFile()
}

// AutoDistributeRoom 分配寝室
func (stm *Stm) AutoDistributeRoom(stu *Student) (int, bool) {
	fmt.Println("开始寝室分配")
	if ok := stm.AllRoom[stu.RoomId].CheckRoomIsFull(); !ok {
		return stu.RoomId, true
	} else {
		for _, room := range stm.AllRoom {
			if room.Gender == stu.Gender && !room.CheckRoomIsFull() {
				fmt.Printf("%d寝室已满员,已将学生%s自动分配至寝室%d\n", stu.RoomId, stu.Name, room.Id)
				return room.Id, true
			}
		}
	}
	fmt.Printf("所有%s寝室都已满员！\n", stu.Gender)
	return 0, false
}

func (stm *Stm) LoadAllStu() {
	jsonFile, err := os.Open("data.json")
	if err != nil {
		fmt.Println("打开文件错误")
		return
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			fmt.Println("关闭文件错误！")
		}
	}(jsonFile)
	decoder := json.NewDecoder(jsonFile)
	//读取json的学生数据
	err = decoder.Decode(&stm.AllStudent)
	if err != nil {
		fmt.Println("error opening json file")
	}
}

func (stm *Stm) LoadAllRoom() {
	jsonFile, err := os.Open("room.json")
	if err != nil {
		fmt.Println("打开文件错误")
		return
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			fmt.Println("关闭文件错误！")
		}
	}(jsonFile)
	decoder := json.NewDecoder(jsonFile)
	err = decoder.Decode(&stm.AllRoom)
	if err != nil {
		fmt.Println("error opening json file")
	}
}

// 将所有寝室数据写入文件
func (stm *Stm) writeAllRoomToFile() {
	jsonFile, err := os.Create("room.json") // 创建 json 文件
	if err != nil {
		fmt.Println("error", err)
		return
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			fmt.Println("关闭文件错误！")
		}
	}(jsonFile)
	encode := json.NewEncoder(jsonFile) // 创建编码器
	err = encode.Encode(stm.AllRoom)    // 编码
	if err != nil {
		fmt.Printf("encode error [ %v ]", err)
		return
	}
}

// 将所有学生数据写入文件
func (stm *Stm) writeAllStuToFile() {
	jsonFile, err := os.Create("data.json") // 创建 json 文件
	if err != nil {
		fmt.Println("error", err)
		return
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			fmt.Println("关闭文件错误！")
		}
	}(jsonFile)
	encode := json.NewEncoder(jsonFile) // 创建编码器
	err = encode.Encode(stm.AllStudent) // 编码
	if err != nil {
		fmt.Printf("encode error [ %v ]", err)
		return
	}

}

// ShowMenu 显示菜单
func (stm *Stm) ShowMenu() {
	fmt.Print(`
*************************************
*  欢迎使用学生寝室管理系统             
*	1. 录入学生信息                   
*	2. 学生寝室调整                   
*	3. 查看所有学生信息                
*	4. 查看所有寝室信息                
*	5. 查看特定学生信息                
*	6. 退出                          
*************************************
请输入您的选择：`)
}

// CheckCanAddStu 校验是否可以录入学生
func (stm *Stm) CheckCanAddStu() bool {
	for _, room := range stm.AllRoom {
		if room.AvaiCount > 0 {
			return true
		}
	}
	return false
}

// PrintStuInRoom 打印寝室中的学生信息
func (stm *Stm) PrintStuInRoom(id int) {
	fmt.Printf("-------------------%d 寝室信息表-------------------\n", id)
	fmt.Println(" 寝室Id     学生Id     姓名       性别      年龄   ")
	for _, stu := range stm.AllStudent {
		if stu.RoomId == id {
			fmt.Printf("%6d %10d %9s %8s %9d \n", stu.RoomId, stu.Id, stu.Name, stu.Gender, stu.Age)
		}
	}
}

func (stm *Stm) CheckIfStuExist(id int) (*Student, bool) {
	if stm.AllStudent[id] != nil {
		return stm.AllStudent[id], true
	}
	return nil, false
}

// CheckRoomExist 判断寝室是否存在
func (stm *Stm) CheckRoomExist(id int) (*Room, bool) {
	if stm.AllRoom[id] != nil {
		return stm.AllRoom[id], true
	}
	fmt.Printf("id为%d的寝室不存在,请重新操作！\n", id)
	return nil, false
}

// CheckStuInRoom 校验学生是否在寝室中
func (stm *Stm) CheckStuInRoom(stu *Student, roomId int) (*Student, bool) {
	if roomId == stu.Id {
		return stu, true
	}
	return nil, false
}

func (stm *Stm) PrintAllStudent() {
	fmt.Println("-------------------学生信息表-------------------")
	fmt.Println(" Id    姓名       性别      年龄    寝室编号")
	for _, student := range stm.AllStudent {
		student.PrintStuInfo()
	}
}

// PrintAllStu 打印所有学生
func (stm *Stm) PrintAllStu() {
	fmt.Println("-------------------学生信息表-------------------")
	fmt.Println(" Id    姓名       性别      年龄    寝室编号")
	for _, student := range stm.AllStudent {
		student.PrintStuInfo()
	}
}

// PrintStuById 根据Id打印学生信息
func (stm *Stm) PrintStuById() {
	fmt.Print("您选择了查看特定学生信息,请输入要查看的学生Id：")
	var id int
	_, err := fmt.Scanln(&id)
	if err != nil {
		fmt.Println("学生Id不正确！")
		stm.PrintStuById()
		return
	}
	if stu, ok := stm.CheckIfStuExist(id); ok {
		fmt.Println("-------------------学生信息表-------------------")
		fmt.Println(" Id    姓名       性别      年龄    寝室编号")
		stu.PrintStuInfo()
		fmt.Print("是否继续查看?y/n:")
		var choice string
		_, err2 := fmt.Scanln(&choice)
		if err2 != nil {
			fmt.Println("输入错误！")
			return
		}
		if choice == "y" || choice == "Y" {
			stm.PrintStuById()
		}
	} else {
		fmt.Println("学生不存在,请重新输入")
		stm.PrintStuById()
	}
}

// PrintAllRoom 打印所有寝室信息
func (stm *Stm) PrintAllRoom() {
	fmt.Println("-------------------寝室信息表-------------------")
	fmt.Println(" 寝室Id     学生Id     姓名       性别      年龄   当前剩余床位")
	for key, stubs := range stm.AllRoom {
		for _, stu := range stm.AllStudent {
			if stu.RoomId == key {
				fmt.Printf("%6d %10d %9s %8s %9d %9d\n", key, stu.Id, stu.Name, stu.Gender, stu.Age, stubs.AvaiCount)
			}
		}
	}
}

// Quit 退出
func (stm *Stm) Quit() bool {
	var answer string
	fmt.Print("确认退出系统吗？y/n:")
	_, err := fmt.Scanln(&answer)
	if err != nil {
		fmt.Println("感谢您的使用,下次再见！")
		return true
	}
	if answer == "y" || answer == "Y" {
		fmt.Println("感谢您的使用,下次再见！")
		return true
	}
	return false
}

// CheckChangeInfo 校验目标寝室信息是否正确
func (stm *Stm) CheckChangeInfo(stu *Student) (*Room, bool) {
	targetRoom := &Room{}
	fmt.Print("请输入目标寝室Id:")
	_, err := fmt.Scanln(&targetRoom.Id)
	if err != nil {
		fmt.Println("目标寝室Id不正确！")
		return nil, false
	}
	//判断自己寝室不能换到自己寝室
	if targetRoom.Id == stu.RoomId {
		fmt.Printf("学生%s已经在寝室%d了,不能这样操作！\n", stu.Name, stu.RoomId)
		return nil, false
	}
	//判断目标寝室是否存在
	if targetRoom, ok := stm.CheckRoomExist(targetRoom.Id); ok {
		//判断目标寝室男女正确
		if ok := targetRoom.CheckStuGenderRight(stu); !ok {
			return nil, false
		}
		return targetRoom, true
	}
	return nil, false
}
