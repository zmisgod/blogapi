// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/zmisgod/blogApi/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/tag",
			beego.NSInclude(
				&controllers.TagController{},
			),
		),
		beego.NSNamespace("/article",
			beego.NSInclude(
				&controllers.ArticleController{},
			),
		),
		beego.NSNamespace("/category",
			beego.NSInclude(
				&controllers.CategoryController{},
			),
		),
		beego.NSNamespace("/home",
			beego.NSInclude(
				&controllers.HomeController{},
			),
		),
		beego.NSNamespace("/search",
			beego.NSInclude(
				&controllers.SearchController{},
			),
		),
		beego.NSNamespace("/comment",
			beego.NSInclude(
				&controllers.CommentController{},
			),
		),
		beego.NSNamespace("/link",
			beego.NSInclude(
				&controllers.LinkController{},
			),
		),
		beego.NSNamespace("/topics",
			beego.NSInclude(
				&controllers.TopicsController{},
			),
		),
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/badge",
			beego.NSInclude(
				&controllers.BadgeController{},
			),
		),
		beego.NSNamespace("/crh",
			beego.NSInclude(
				&controllers.CrhController{},
			),
		),
	)
	beego.AddNamespace(ns)

	admin := beego.NewNamespace("/admin",
		beego.NSNamespace("/wps/article",
			beego.NSInclude(
				&controllers.AdminWpsArticleController{},
			),
		),
		beego.NSNamespace("/wps/tag",
			beego.NSInclude(
				&controllers.AdminWpsTagController{},
			),
		),
		beego.NSNamespace("/wps/topic",
			beego.NSInclude(
				&controllers.AdminWpsTopicController{},
			),
		),
		beego.NSNamespace("/upload",
			beego.NSInclude(
				&controllers.AdminUploadController{},
			),
		),
	)
	beego.AddNamespace(admin)
	beego.ErrorController(&controllers.ErrorController{})
}
