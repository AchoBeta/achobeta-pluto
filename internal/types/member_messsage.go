package types

import "time"

// 查询成员详细信息(入参）
type GetMemberDetailReq struct {
	UserID int64 `json:"user_id" binging:"required" `
}

// 职位信息
type StructureNodes struct {
	ID   int64  `json:"position_id"`
	Name string `json:"position_name"`
}

type MemberPositions struct {
	//团队id
	TeamId int64 `json:"team_id"`
	//团队名称
	TeamName string `json:"team_name"`
	//组别 + 职位
	StructNodes []StructureNodes `json:"position_node"`
	//权限级别
	Level int `json:"level"`
}

// 查询成员详细信息(出参）
type GetMemberDetailResp struct {
	Name           string            `json:"name"`
	Sex            string            `json:"sex"`
	CreateDate     time.Time         `json:"create_date"`
	IdCard         string            `json:"id_card"`
	PhoneNum       uint64            `json:"phone_num"`
	Email          string            `json:"email"`
	Grade          uint              `json:"grade"`
	Major          string            `json:"major"`
	StudentID      uint64            `json:"student_id"`
	Experience     string            `json:"experience"`
	Status         string            `json:"status"`
	LikeCount      uint64            `json:"like_count"`
	MemberPosition []MemberPositions `json:"member_position"`
}

// 给用户点赞/取消赞(入参）
type LikeCountReq struct {
	MemberID int64 `json:"member_id" binging:"required" `
}

// 给用户点赞/取消赞(出参）
type LikeCountResp struct {
	LikeCount uint64 `json:"like_count"`
}

// 编辑成员详细信息(入参）
type PutTeamMemberReq struct {
	ID         int64     `json:"id" binging:"required" `
	Name       string    `json:"name"`
	Sex        string    `json:"sex"`
	CreateDate time.Time `json:"create_date"`
	IdCard     string    `json:"id_card"`
	PhoneNum   uint64    `json:"phone_num" binging:"required,len = 11"`
	Email      string    `json:"email"`
	Grade      uint      `json:"grade"`
	Major      string    `json:"major"`
	StudentID  uint64    `json:"student_id"`
	Experience string    `json:"experience"`
	Status     string    `json:"status"`
	//组别 + 职位
	MemberPosition []MemberPositions `json:"member_position"`
}

// 编辑成员详细信息(出参）
type PutTeamMemberResp struct{}
