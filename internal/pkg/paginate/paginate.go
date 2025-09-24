package paginate

import (
	"fmt"
	"go-job/internal/model"
	"gorm.io/gorm"
	"strings"
)

const (
	defaultPageNum  = 1
	defaultPageSize = 20
	maxPageSize     = 50
)

func PaginateList[T any](db *gorm.DB, page model.Page) (model.Page, error) {
	if page.PageNum < 1 {
		page.PageNum = defaultPageNum
	}
	if page.PageSize < 1 {
		page.PageSize = defaultPageSize
	} else if page.PageSize > maxPageSize {
		page.PageSize = maxPageSize
	}

	var (
		m      T
		result model.Page
		list   []T
		total  int64
	)

	offset := (page.PageNum - 1) * page.PageSize

	// Count total
	if err := db.Model(&m).Count(&total).Error; err != nil {
		return result, err
	}
	tx := db.Model(&m).Limit(page.PageSize).Offset(offset)

	// order
	if page.Sort != "" {
		var orderBy string
		if strings.ToLower(page.Order) == "desc" {
			orderBy = fmt.Sprintf("%s desc", page.Sort)
		} else {
			orderBy = fmt.Sprintf("%s asc", page.Sort)
		}
		tx = tx.Order(orderBy)
	}
	if err := tx.Find(&list).Error; err != nil {
		return result, err
	}
	page.Data = list
	return page, nil
}

func PaginateListV2[T any](db *gorm.DB, page model.Page, fns ...func(*gorm.DB) *gorm.DB) (model.Page, error) {
	if !page.UnlimitedPageSize {
		if page.PageNum < 1 {
			page.PageNum = defaultPageNum
		}
		if page.PageSize < 1 {
			page.PageSize = defaultPageSize
		} else if page.PageSize > maxPageSize {
			page.PageSize = maxPageSize
		}
	}

	var (
		result model.Page
		list   []T
	)

	offset := (page.PageNum - 1) * page.PageSize

	query := db.Model(new(T))
	if fns != nil {
		for _, fn := range fns {
			query = fn(query)
		}
	}

	// Count total
	if err := query.Count(&page.Total).Error; err != nil {
		return result, err
	}

	// order
	if page.Sort != "" {
		var orderBy string
		if strings.ToLower(page.Order) == "desc" {
			orderBy = fmt.Sprintf("%s desc", page.Sort)
		} else {
			orderBy = fmt.Sprintf("%s asc", page.Sort)
		}
		query = query.Order(orderBy)
	}

	// Query data
	if err := query.Limit(page.PageSize).Offset(offset).Find(&list).Error; err != nil {
		return result, err
	}

	page.Data = list
	return page, nil
}
