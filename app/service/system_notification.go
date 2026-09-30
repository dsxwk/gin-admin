package service

import (
	"context"
	"errors"

	"gin/app/enum"
	"gin/app/model"
	"gin/app/request"
	"gin/common/base"
	"gin/pkg/errcode"

	"gorm.io/gorm"
)

const (
	notificationBatchSize = 500
)

// SystemNotificationService 系统通知服务
type SystemNotificationService struct {
	base.BaseService
}

// NotificationItem 通知列表项
type NotificationItem struct {
	ID             int64           `gorm:"column:id" json:"id"`                          // 通知ID
	FromUserId     int64           `gorm:"column:from_user_id" json:"fromUserId"`        // 发送用户ID
	Type           int64           `gorm:"column:type" json:"type"`                      // 通知类型 1=用户 2=部门
	IsToAll        int64           `gorm:"column:is_to_all" json:"isToAll"`              // 是否推送所有 1=是 2=否
	Title          string          `gorm:"column:title" json:"title"`                    // 标题
	Content        string          `gorm:"column:content" json:"content"`                // 内容
	Status         int64           `gorm:"column:status" json:"status"`                  // 状态 1=正常 2=撤回
	IsRead         int64           `gorm:"column:is_read" json:"isRead"`                 // 是否已读 1=是 2=否
	RecipientCount int64           `gorm:"column:recipient_count" json:"recipientCount"` // 接收人数
	CreatedAt      *model.DateTime `gorm:"column:created_at" json:"createdAt"`           // 创建时间
	UpdatedAt      *model.DateTime `gorm:"column:updated_at" json:"updatedAt"`           // 更新时间
}

// NotificationPage 通知分页
type NotificationPage struct {
	Total    int64              `json:"total"`    // 总条数
	Unread   int64              `json:"unread"`   // 未读数量
	Page     int                `json:"page"`     // 当前页
	PageSize int                `json:"pageSize"` // 每页数量
	List     []NotificationItem `json:"list"`     // 通知列表
}

// Push 推送通知
func (s *SystemNotificationService) Push(ctx context.Context, fromUserID int64, req request.SystemNotificationPush) (m model.SystemNotification, userIDs []int64, err error) {
	if fromUserID <= 0 {
		return m, nil, errcode.NewError(401, "未登录")
	}

	userIDs, err = s.resolveRecipients(ctx, req)
	if err != nil {
		return m, nil, err
	}

	m = model.SystemNotification{
		FromUserId: fromUserID,
		Type:       req.Type,
		IsToAll:    req.IsToAll,
		Title:      req.Title,
		Content:    req.Content,
		Status:     enum.NotificationStatusNormal,
	}

	db := s.DB(ctx, &m)
	tx := db.Begin()
	if tx.Error != nil {
		return m, nil, tx.Error
	}

	if err = tx.Model(&m).Create(&m).Error; err != nil {
		tx.Rollback()
		return m, nil, err
	}

	rows := make([]model.NotificationUsers, 0, len(userIDs))
	for _, userID := range userIDs {
		rows = append(rows, model.NotificationUsers{
			NotificationId: m.ID,
			ToUserId:       userID,
			IsRead:         enum.NotificationIsReadNo,
		})
	}

	if err = tx.Model(&model.NotificationUsers{}).CreateInBatches(&rows, notificationBatchSize).Error; err != nil {
		tx.Rollback()
		return m, nil, err
	}

	if err = tx.Commit().Error; err != nil {
		return m, nil, err
	}

	return m, userIDs, nil
}

// List 获取通知列表
func (s *SystemNotificationService) List(ctx context.Context, userID int64, page, pageSize int, sent bool) (NotificationPage, error) {
	if userID <= 0 {
		return NotificationPage{}, errcode.NewError(401, "未登录")
	}
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if sent {
		return s.listSent(ctx, userID, page, pageSize)
	}
	return s.listReceived(ctx, userID, page, pageSize)
}

// listReceived 获取当前用户收到的通知
func (s *SystemNotificationService) listReceived(ctx context.Context, userID int64, page, pageSize int) (NotificationPage, error) {
	var (
		total  int64
		unread int64
		items  = make([]NotificationItem, 0)
		m      model.SystemNotification
	)

	db := s.DB(ctx, &m)
	query := func() *gorm.DB {
		return db.Table("notification_users AS nu").
			Joins("JOIN system_notification AS sn ON sn.id = nu.notification_id").
			Where("nu.to_user_id = ? AND sn.status = ?", userID, enum.NotificationStatusNormal)
	}

	if err := query().Count(&total).Error; err != nil {
		return NotificationPage{}, err
	}

	if err := query().Where("nu.is_read = ?", enum.NotificationIsReadNo).Count(&unread).Error; err != nil {
		return NotificationPage{}, err
	}

	offset, limit := request.Pagination(page, pageSize)
	if err := query().Select(
		"sn.id,sn.from_user_id,sn.type,sn.is_to_all,sn.title,sn.content,sn.status,sn.created_at,sn.updated_at,nu.is_read",
	).Order("sn.id DESC").Offset(offset).Limit(limit).Scan(&items).Error; err != nil {
		return NotificationPage{}, err
	}

	return NotificationPage{
		Total:    total,
		Unread:   unread,
		Page:     page,
		PageSize: pageSize,
		List:     items,
	}, nil
}

// listSent 获取当前用户发送的通知
func (s *SystemNotificationService) listSent(ctx context.Context, userID int64, page, pageSize int) (NotificationPage, error) {
	var (
		total int64
		items = make([]NotificationItem, 0)
		m     model.SystemNotification
	)

	db := s.DB(ctx, &m)
	query := func() *gorm.DB {
		return db.Table("system_notification AS sn").Where("sn.from_user_id = ?", userID)
	}

	if err := query().Count(&total).Error; err != nil {
		return NotificationPage{}, err
	}

	offset, limit := request.Pagination(page, pageSize)
	if err := query().Select(
		"sn.id,sn.from_user_id,sn.type,sn.is_to_all,sn.title,sn.content,sn.status,sn.created_at,sn.updated_at,(SELECT COUNT(*) FROM notification_users AS nu WHERE nu.notification_id = sn.id) AS recipient_count",
	).Order("sn.id DESC").Offset(offset).Limit(limit).Scan(&items).Error; err != nil {
		return NotificationPage{}, err
	}

	return NotificationPage{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		List:     items,
	}, nil
}

// Read 标记通知已读
func (s *SystemNotificationService) Read(ctx context.Context, userID int64, ids []int64) (int64, error) {
	if userID <= 0 {
		return 0, errcode.NewError(401, "未登录")
	}

	db := s.DB(ctx, &model.NotificationUsers{})
	query := db.Model(&model.NotificationUsers{}).
		Where("to_user_id = ? AND is_read = ?", userID, enum.NotificationIsReadNo)

	ids = uniqueNotificationIDs(ids)
	if len(ids) > 0 {
		query = query.Where("notification_id IN ?", ids)
	}

	result := query.Update("is_read", enum.NotificationIsReadYes)
	return result.RowsAffected, result.Error
}

// Revoke 撤回通知
func (s *SystemNotificationService) Revoke(ctx context.Context, userID, notificationID int64) (m model.SystemNotification, userIDs []int64, err error) {
	if userID <= 0 {
		return m, nil, errcode.NewError(401, "未登录")
	}
	db := s.DB(ctx, &m)
	tx := db.Begin()
	if tx.Error != nil {
		return m, nil, tx.Error
	}

	if err = tx.First(&m, notificationID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return m, nil, errcode.NewError(404, "通知不存在")
		}
		return m, nil, err
	}
	if m.FromUserId != userID {
		tx.Rollback()
		return m, nil, errcode.NewError(403, "无权撤回该通知")
	}

	if m.Status != enum.NotificationStatusRevoked {
		if err = tx.Model(&m).Update("status", enum.NotificationStatusRevoked).Error; err != nil {
			tx.Rollback()
			return m, nil, err
		}
		m.Status = enum.NotificationStatusRevoked
	}

	if err = tx.Model(&model.NotificationUsers{}).Where("notification_id = ?", notificationID).Pluck("to_user_id", &userIDs).Error; err != nil {
		tx.Rollback()
		return m, nil, err
	}

	if err = tx.Commit().Error; err != nil {
		return m, nil, err
	}

	return m, uniqueNotificationIDs(userIDs), nil
}

// resolveRecipients 解析通知接收用户
func (s *SystemNotificationService) resolveRecipients(ctx context.Context, req request.SystemNotificationPush) ([]int64, error) {
	if req.IsToAll == enum.NotificationIsToAllYes {
		userIDs, err := s.enabledUserIDs(ctx)
		if err != nil {
			return nil, err
		}
		if len(userIDs) == 0 {
			return nil, errcode.NewError(400, "没有可接收通知的用户")
		}
		return userIDs, nil
	}

	switch req.Type {
	case enum.NotificationTypeUser:
		userIDs := uniqueNotificationIDs(req.UserIds)
		if len(userIDs) == 0 {
			return nil, errcode.NewError(400, "接收用户不能为空")
		}
		validUserIDs, err := s.filterEnabledUserIDs(ctx, userIDs)
		if err != nil {
			return nil, err
		}
		if len(validUserIDs) != len(userIDs) {
			return nil, errcode.NewError(400, "接收用户不存在或已停用")
		}
		return validUserIDs, nil
	case enum.NotificationTypeDepartment:
		departmentIDs := uniqueNotificationIDs(req.DepartmentIds)
		if len(departmentIDs) == 0 {
			return nil, errcode.NewError(400, "接收部门不能为空")
		}

		var departmentCount int64
		if err := s.DB(ctx, &model.Department{}).Model(&model.Department{}).
			Where("id IN ? AND status = ?", departmentIDs, enum.DepartmentStatusEnabled).
			Count(&departmentCount).Error; err != nil {
			return nil, err
		}
		if departmentCount != int64(len(departmentIDs)) {
			return nil, errcode.NewError(400, "接收部门不存在或已停用")
		}

		var userIDs []int64
		if err := s.DB(ctx, &model.UserDepartments{}).Model(&model.UserDepartments{}).
			Where("department_id IN ?", departmentIDs).
			Distinct().
			Pluck("user_id", &userIDs).Error; err != nil {
			return nil, err
		}

		userIDs, err := s.filterEnabledUserIDs(ctx, userIDs)
		if err != nil {
			return nil, err
		}
		if len(userIDs) == 0 {
			return nil, errcode.NewError(400, "接收部门没有可接收通知的用户")
		}
		return userIDs, nil
	default:
		return nil, errcode.NewError(400, "通知类型错误")
	}
}

// enabledUserIDs 获取所有启用用户ID
func (s *SystemNotificationService) enabledUserIDs(ctx context.Context) ([]int64, error) {
	var userIDs []int64
	err := s.DB(ctx, &model.User{}).Model(&model.User{}).
		Where("status = ?", enum.UserStatusEnabled).
		Pluck("id", &userIDs).Error
	return uniqueNotificationIDs(userIDs), err
}

// filterEnabledUserIDs 过滤启用用户ID
func (s *SystemNotificationService) filterEnabledUserIDs(ctx context.Context, userIDs []int64) ([]int64, error) {
	userIDs = uniqueNotificationIDs(userIDs)
	if len(userIDs) == 0 {
		return nil, nil
	}

	var validUserIDs []int64
	err := s.DB(ctx, &model.User{}).Model(&model.User{}).
		Where("id IN ? AND status = ?", userIDs, enum.UserStatusEnabled).
		Pluck("id", &validUserIDs).Error
	return uniqueNotificationIDs(validUserIDs), err
}

// uniqueNotificationIDs 通知ID去重
func uniqueNotificationIDs(ids []int64) []int64 {
	if len(ids) == 0 {
		return nil
	}

	seen := make(map[int64]struct{}, len(ids))
	result := make([]int64, 0, len(ids))
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}
	return result
}
