package types

// 查询团队架构(入参)
type TeamStructReq struct {
	TeamId int64 `form:"team_id" binding:"required"`
}

// 单个团队架构记录
type TeamStructure struct {
	TeamId    int64  `json:"team_id" binding:"required"`    // 团队 ID
	MyselfId  int64  `json:"myself_id"`                     // 当前节点 ID
	FatherId  int64  `json:"father_id" binding:"required"`  // 父节点 ID
	NodeName  string `json:"node_name"  binding:"required"` // 架构名称（职位名称等）
	IsDeleted int8   `json:"is_deleted" binding:"required"` //查看是否被删除

}

// 查询团队架构（出参）
type TeamStructResp struct {
	RootOfTeam     int64           `json:"root_of_team"`
	TeamStructures []TeamStructure `json:"team_structures"` // 团队架构记录的数组
}

// 保存已经更改好的团队架构信息(入参）
type PutTeamNodeReq struct {
	TeamStructures []TeamStructure `json:"team_structures"` // 团队架构记录的数组
}

// 保存已经更改好的团队架构信息(出参）
type PutTeamNodeResp struct {
	// msg,code
}

// 新增团队，初始化团队架构（入参）
type CreateTeamReq struct {
	Name string `json:"team_name" binding:"required"`
}

// 新增团队，初始化团队架构(出参）
type CreateTeamResp struct {
}
