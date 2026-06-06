package v1

// ScriptGetReq 获取剧本请求
type ScriptGetReq struct {
	ProjectId int64 `json:"project_id" v:"required#项目ID不能为空"`
}

// ScriptGetRes 获取剧本响应
type ScriptGetRes struct {
	Id                int64  `json:"id"`
	VersionNo         int    `json:"version_no"`
	YamlContent       string `json:"yaml_content"`
	ValidationStatus  string `json:"validation_status"`
	HallucinationRisk string `json:"hallucination_risk"`
	SafetyRisk        string `json:"safety_risk"`
	CreatedBy         string `json:"created_by"`
	CreatedAt         string `json:"created_at"`
}

// ScriptUpdateReq 修改剧本请求
type ScriptUpdateReq struct {
	ProjectId   int64  `json:"project_id"   v:"required#项目ID不能为空"`
	YamlContent string `json:"yaml_content" v:"required#YAML内容不能为空"`
}

// ScriptValidateReq 校验YAML请求
type ScriptValidateReq struct {
	ProjectId int64 `json:"project_id" v:"required#项目ID不能为空"`
}

// ScriptValidateRes 校验YAML响应
type ScriptValidateRes struct {
	Valid  bool              `json:"valid"`
	Issues []ValidationIssue `json:"issues"`
}

// ValidationIssue 校验问题
type ValidationIssue struct {
	IssueType    string `json:"issue_type"`
	Severity     string `json:"severity"`
	Message      string `json:"message"`
	LocationPath string `json:"location_path"`
	Suggestion   string `json:"suggestion"`
}

// ScriptExportReq 导出剧本请求
type ScriptExportReq struct {
	ProjectId int64  `json:"project_id" v:"required#项目ID不能为空"`
	Format    string `json:"format"     v:"in:yaml,markdown#导出格式不合法"`
}

// HallucinationCheckRes 幻觉检测响应
type HallucinationCheckRes struct {
	HallucinationRisk string `json:"hallucination_risk"`
}

// SafetyCheckRes 安全审查响应
type SafetyCheckRes struct {
	SafetyRisk string `json:"safety_risk"`
}
