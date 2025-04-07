package paginate

import (
	"go-job/internal/model"
	"gorm.io/gorm"
)

func PaginateList[T any](db *gorm.DB, pageNum, pageSize int) (model.Page, error) {
	const (
		defaultPageNum  = 1
		defaultPageSize = 20
		maxPageSize     = 100
	)

	if pageNum < 1 {
		pageNum = defaultPageNum
	}
	if pageSize < 1 {
		pageSize = defaultPageSize
	} else if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	var (
		result model.Page
		list   []T
		total  int64
	)

	offset := (pageNum - 1) * pageSize

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return result, err
	}
	if err := db.Limit(pageSize).Offset(offset).Find(&list).Error; err != nil {
		return result, err
	}
	result = model.Page{
		Total:    total,
		PageSize: pageSize,
		PageNum:  pageNum,
		Data:     list,
	}
	return result, nil
}
