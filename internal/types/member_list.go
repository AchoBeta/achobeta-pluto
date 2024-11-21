package types

import "time"

// 查询指定团队成员列表（入参）
type MemberlistReq struct {
	TeamID  int64 `json:"team_id" binding:"required"`
	Page    int   `json:"page" binding:"required"`
	Perpage int   `json:"perpage" binding:"required"`
}

type Members struct {
	ID           int64    `json:"id" gorm:"id"`
	Name         string   `json:"name" gorm:"name"`
	Grade        uint     `json:"grade" gorm:"grade"`
	Major        string   `json:"major" gorm:"major"`
	Status       string   `json:"status" gorm:"status"`
	Positions    string   `json:"-"`         // 原始字段，用于接收 GROUP_CONCAT
	PositionList []string `json:"positions"` // 最终切片
}

// 查询指定团队成员列表（出参）
type MemberlistResp struct {
	Members []Members
}

// 新增成员（入参）
type CreateMemberReq struct {
	Name       string    `json:"name"`
	Sex        string    `json:"sex"`
	CreateDate time.Time `json:"create_date"  binding:"required,datetime=2006/01/02"` //日期格式校验
	IdCard     string    `json:"id_card"`
	PhoneNum   uint64    `json:"phone_num" binding:"required,len = 11"`
	Email      string    `json:"email"`
	Grade      uint      `json:"grade"`
	Major      string    `json:"major"`
	StudentID  uint64    `json:"student_id"`
	Experience string    `json:"experience"`
	Status     string    `json:"status"`

	MemberPositions []MemberPositions `json:"member_position"`
}

// 新增成员（出参）
type CreateMembersResp struct {
	//
}

// 删除成员（入参）
type DeleteMemberReq struct {
	TeamId   int64 `json:"team_id" binding:"required"`
	MemberId int64 `json:"member_id" binding:"required"`
}

// 删除成员（出参
type DeleteMembersResp struct{}
