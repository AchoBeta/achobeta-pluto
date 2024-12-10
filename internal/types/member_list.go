package types

// 查询指定团队成员列表（入参）
type MemberlistReq struct {
	TeamID  int64 `form:"team_id" binding:"required"`
	Page    int   `form:"page" binding:"required"`
	Perpage int   `form:"perpage" binding:"required"`
}

type Members struct {
	ID           int64    `json:"id"   gorm:"id"`
	Name         string   `json:"name" gorm:"name"`
	Grade        string   `json:"grade" gorm:"grade"`
	Major        string   `json:"major" gorm:"major"`
	Status       string   `json:"status" gorm:"status"`
	PhoneNum     string   `json:"phone_num"`
	Positions    string   `json:"-"`         // 原始字段，用于接收 GROUP_CONCAT
	PositionList []string `json:"positions"` // 最终切片
}

// 查询指定团队成员列表（出参）
type MemberlistResp struct {
	Members []Members `json:"members"`
	Total   int64     `json:"total"`
}

// 新增成员（入参）
type CreateMemberReq struct {
	Name            string            `json:"name"`
	Sex             string            `json:"sex"`
	CreateDate      string            `json:"create_date" binding:"omitempty"` //日期格式校验
	IdCard          *string           `json:"id_card"`
	PhoneNum        string            `json:"phone_num" binding:"required,len=11"`
	Email           *string           `json:"email"`
	Grade           string            `json:"grade"`
	Major           string            `json:"major"`
	StudentID       *string           `json:"student_id"`
	Experience      string            `json:"experience"`
	Status          string            `json:"status"`
	MemberPositions []MemberPositions `json:"member_position"`
}

// 新增成员（出参）
type CreateMembersResp struct {
	//
}

// 删除成员（入参）
type DeleteMemberReq struct {
	TeamId   int64 `uri:"team_id" binding:"required"`
	MemberId int64 `uri:"member_id" binding:"required"`
}

// 删除成员（出参）
type DeleteMembersResp struct{}
