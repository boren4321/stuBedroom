package stum

import "fmt"

type Room struct {
	Id        int    `json:"id"`
	AvaiCount int    `json:"avaiCount"`
	Gender    string `json:"gender"`
}

func (room *Room) CheckRoomIsFull() bool {
	return room.AvaiCount == 0
}

// SettleInStu 学生入住寝室
func (room *Room) SettleInStu() {
	room.AvaiCount -= 1
}

func (room *Room) RemoveStu() {
	room.AvaiCount += 1
}

// CheckStuGenderRight 校验寝室性别跟入住学生性别是否一致
func (room *Room) CheckStuGenderRight(stu *Student) bool {
	if room.Gender != stu.Gender {
		fmt.Printf(" %d 为 %s寝室,不能入住 %s生\n", room.Id, room.Gender, stu.Gender)
		return false
	}
	return true
}
