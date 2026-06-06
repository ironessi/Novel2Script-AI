package audit

import (
	"context"

	"novel2script-backend/internal/dao"
	"novel2script-backend/internal/model/entity"
	"novel2script-backend/internal/service"
)

type auditImpl struct{}

func init() {
	service.Audit = &auditImpl{}
}

func (a *auditImpl) Log(ctx context.Context, userId, projectId int64, action, resourceType string, resourceId int64, ip, userAgent, requestId string) error {
	return dao.CreateAuditLog(ctx, &entity.AuditLog{
		UserId:       userId,
		ProjectId:    projectId,
		Action:       action,
		ResourceType: resourceType,
		ResourceId:   resourceId,
		IpAddress:    ip,
		UserAgent:    userAgent,
		RequestId:    requestId,
	})
}

func (a *auditImpl) GetList(ctx context.Context, projectId int64, page, pageSize int) ([]entity.AuditLog, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return dao.GetAuditLogList(ctx, projectId, page, pageSize)
}
