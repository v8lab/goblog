package provider

import (
	"kandaoni.com/anqicms/dao"
	"kandaoni.com/anqicms/model"
)

func GetFinanceList(page, pageSize int) ([]*model.Finance, int64) {
	var finances []*model.Finance
	var total int64
	offset := (page - 1) * pageSize
	dao.DB.Model(&model.Finance{}).Count(&total).Order("id desc").Limit(pageSize).Offset(offset).Find(&finances)
	if len(finances) > 0 {
		var userIds = make([]uint, 0, len(finances))
		for i := range finances {
			userIds = append(userIds, finances[i].UserId)
		}
		users := GetUsersInfoByIds(userIds)
		for i := range finances {
			for u := range users {
				if finances[i].UserId == users[u].Id {
					finances[i].UserName = users[u].UserName
				}
			}
		}
	}
	return finances, total
}

func GetFinanceById(id uint) (*model.Finance, error) {
	var finance model.Finance
	err := dao.DB.Where(&model.Finance{}).Where("`id` = ?", id).Take(&finance).Error
	if err != nil {
		return nil, err
	}

	return &finance, nil
}

