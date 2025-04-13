package paginate

import (
	"go-job/internal/model"
	"gorm.io/gorm"
)

const (
	defaultPageNum  = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

func PaginateList[T any](db *gorm.DB, pageNum, pageSize int) (model.Page, error) {
	if pageNum < 1 {
		pageNum = defaultPageNum
	}
	if pageSize < 1 {
		pageSize = defaultPageSize
	} else if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	var (
		m      T
		result model.Page
		list   []T
		total  int64
	)

	offset := (pageNum - 1) * pageSize

	// Count total
	if err := db.Model(&m).Count(&total).Error; err != nil {
		return result, err
	}
	if err := db.Model(&m).Limit(pageSize).Offset(offset).Find(&list).Error; err != nil {
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
