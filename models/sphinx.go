package models

import (
	"github.com/astaxie/beego"
	"github.com/yunge/sphinx"
)

func SphinxSearch(keyword string, page, pageSize int) (interface{}, error) {
	SphinxClient := sphinx.NewClient().SetServer(beego.AppConfig.String("SphinxHost"), 0).SetConnectTimeout(5000).SetMatchMode(2).SetLimits((page-1)*pageSize, pageSize, 10000000, 0)
	if err := SphinxClient.Error(); err != nil {
		return nil, err
	}
	defer SphinxClient.Close()

	// 查询，第一个参数是我们要查询的关键字，第二个是索引名称test1，第三个是备注
	res, err := SphinxClient.Query(keyword, "main", "search article!")
	if err != nil {
		return nil, err
	}
	var articleMap []interface{}
	for _, match := range res.Matches {
		articleMap = append(articleMap, match)
	}
	return articleMap, nil
}
