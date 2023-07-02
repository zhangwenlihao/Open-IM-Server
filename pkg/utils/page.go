package utils

import "github.com/xiaoyiEdu/Open-IM-Server/pkg/proto/sdkws"

func GetPage(pagination *sdkws.RequestPagination) (pageNumber, showNumber int32) {
	if pagination != nil {
		return pagination.PageNumber, pagination.ShowNumber
	}
	return
}
