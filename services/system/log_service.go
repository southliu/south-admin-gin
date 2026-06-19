package services

import (
	"south-admin-gin/database"
	"south-admin-gin/models/system"
)

// CreateLog 创建日志
func CreateLog(dto models.CreateLogDto) error {
	log := models.Log{
		Username:  dto.Username,
		IP:        dto.IP,
		Method:    dto.Method,
		URL:       dto.URL,
		Params:    dto.Params,
		UserAgent: dto.UserAgent,
		Status:    dto.Status,
		Error:     dto.Error,
		Latency:   dto.Latency,
		Type:      dto.Type,
	}
	return database.DB.Create(&log).Error
}

// GetLogPage 获取日志分页列表
func GetLogPage(page, pageSize int, username string, logType int) (*models.PageResult, error) {
	var logs []models.Log
	var total int64

	query := database.DB.Model(&models.Log{})

	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if logType >= 0 {
		query = query.Where("type = ?", logType)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, err
	}

	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&logs).Error
	if err != nil {
		return nil, err
	}

	// 转换为结果格式
	var items []models.LogPageResult
	for _, log := range logs {
		items = append(items, models.LogPageResult{
			ID:        log.ID,
			Username:  log.Username,
			IP:        log.IP,
			Method:    log.Method,
			URL:       log.URL,
			Status:    log.Status,
			Error:     log.Error,
			Type:      log.Type,
			CreatedAt: log.CreatedAt,
		})
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &models.PageResult{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

// DeleteLog 删除日志
func DeleteLog(id int64) error {
	return database.DB.Delete(&models.Log{}, id).Error
}

// BatchDeleteLog 批量删除日志
func BatchDeleteLog(ids []int64) error {
	return database.DB.Delete(&models.Log{}, ids).Error
}

// ClearLog 清空所有日志
func ClearLog() error {
	return database.DB.Where("1 = 1").Delete(&models.Log{}).Error
}
