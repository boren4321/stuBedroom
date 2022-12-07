package stum

import "fmt"

type Student struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
	Age    int    `json:"age"`
	RoomId int    `json:"roomId"`
}

// CheckStudentIdRight 校验输入的学生Id是否正确
func (stu *Student) CheckStudentIdRight() bool {
	fmt.Print("请输入学生Id：")
	_, err := fmt.Scanln(&stu.Id)
	if err != nil {
		fmt.Println("输入的学生Id不正确！")
		return false
	}
	if stu.Id <= 0 {
		fmt.Println("输入的学生Id不能为负数！")
		return false
	}
	return true
}

// CheckStudentNameRight 校验输入的学生姓名是否正确
func (stu *Student) CheckStudentNameRight() bool {
	fmt.Print("请输入学生姓名：")
	_, err := fmt.Scanln(&stu.Name)
	if err != nil {
		fmt.Println("输入的学生姓名不正确！")
		return false
	}
	return true
}

// CheckStudentGenderRight 校验输入的学生性别是否正确
func (stu *Student) CheckStudentGenderRight() bool {
	fmt.Print("请输入学生性别:男/女")
	_, err := fmt.Scanln(&stu.Gender)
	if err != nil {
		fmt.Println("学生性别不正确！")
		return false
	}
	if stu.Gender != "男" && stu.Gender != "女" {
		fmt.Println("输入性别信息不正确,请重新录入")
		return false
	}
	return true
}

// CheckStudentAgeRight 校验输入的学生年龄是否正确
func (stu *Student) CheckStudentAgeRight() bool {
	fmt.Print("请输入学生年龄:")
	_, err := fmt.Scanln(&stu.Age)
	if err != nil {
		fmt.Println("年龄信息不正确！")
		return false
	}
	if stu.Age < 0 {
		fmt.Println("年龄不能为负数！")
		return false
	}
	return true
}

// CheckStudentRoomIdRight 校验输入的寝室Id是否正确
func (stu *Student) CheckStudentRoomIdRight() bool {
	fmt.Print("请输入学生寝室Id:")
	_, err := fmt.Scanln(&stu.RoomId)
	if err != nil || stu.RoomId < 0 {
		fmt.Println("寝室信息不正确！")
		return false
	}
	return true
}

func (stu *Student) PrintStuInfo() {
	fmt.Printf("%2d %8s %8s %9d %10d\n", stu.Id, stu.Name, stu.Gender, stu.Age, stu.RoomId)
}

// GetRoomInfo 获取学生寝室信息
func (stu *Student) GetRoomInfo(stm *Stm) *Room {
	for i, room := range stm.AllRoom {
		if i == stu.RoomId {
			return room
		}
	}
	return nil
}
