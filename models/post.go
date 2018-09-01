package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

//Post 文章基础结构
type Post struct {
	ID            int    `json:"id"`
	UserID        int    `json:"user_id"`
	UserName      string `json:"user_name"`
	CategoryName  string `json:"category_name"`
	CategoryID    int    `json:"category_id"`
	PostTitle     string `json:"post_title"`
	PostIntro     string `json:"post_intro"`
	CoverURL      string `json:"cover_url"`
	createdAt     int
	CreatedAt     string `json:"created_at"`
	Contents      string `json:"contents"`
	Tags          []Tag  `json:"tags"`
	NumInfo       Num    `json:"num_info"`
	CommentStatus int    `json:"comment_status"`
}

//PostDetail 文章详情
type PostDetail struct {
	Post
	UserInfo UserInfo `json:"user_info"`
}

//PostList 文章列表
type PostList struct {
	Post
}

//GetArticleLists 获取文章列表
func GetArticleLists(page, pageSize int) ([]PostList, error) {
	var (
		rows *sql.Rows
		err  error
	)
	// var Post CommentLists
	postList := []PostList{}
	offset := (page - 1) * pageSize
	rows, err = dbConn.Query(fmt.Sprintf("select p.id,p.user_id, p.post_title,u.name as user_name,c.c_name as category_name,p.post_title,p.post_intro,p.created_at from wps_posts as p left join wps_users as u on p.user_id = u.id left join wps_post_categories as c on c.id = p.cat_id where p.post_status = 1 order by p.created_at desc limit %d,%d", offset, pageSize))
	defer rows.Close()
	if err != nil {
		return postList, err
	}

	for rows.Next() {
		var aPost PostList
		err = rows.Scan(
			&aPost.ID,
			&aPost.UserID,
			&aPost.PostTitle,
			&aPost.UserName,
			&aPost.CategoryName,
			&aPost.PostTitle,
			&aPost.PostIntro,
			&aPost.createdAt,
		)
		tm := time.Unix(int64(aPost.createdAt), 0)
		aPost.CreatedAt = tm.Format("2006-01-02 15:04")

		tags, _ := GetPostTagLists(aPost.ID)
		aPost.Tags = tags

		num, _ := GetArticleNumsByPost(aPost.ID)
		aPost.NumInfo = num

		postList = append(postList, aPost)
	}
	return postList, nil
}

//GetArticleDetail 获取文章详情
func GetArticleDetail(postID int) (PostDetail, error) {
	var post PostDetail
	err := dbConn.QueryRow(
		fmt.Sprintf("select p.comment_status,p.id,p.user_id, p.post_title,u.name as user_name,c.c_name as category_name,p.post_title,p.post_intro,p.created_at,pc.contents,p.cat_id from wps_posts as p left join wps_users as u on p.user_id = u.id left join wps_post_categories as c on c.id = p.cat_id left join wps_post_contents as pc on pc.id = p.id where p.post_status = 1 and p.id  = %d", postID)).
		Scan(
			&post.CommentStatus,
			&post.ID,
			&post.UserID,
			&post.PostTitle,
			&post.UserName,
			&post.CategoryName,
			&post.PostTitle,
			&post.PostIntro,
			&post.createdAt,
			&post.Contents,
			&post.CategoryID,
		)
	if err != nil {
		return post, errors.New("empty")
	}
	tm := time.Unix(int64(post.createdAt), 0)
	post.CreatedAt = tm.Format("2006-01-02 15:04")

	tags, _ := GetPostTagLists(post.ID)
	post.Tags = tags

	num, _ := GetArticleNumsByPost(post.ID)
	post.NumInfo = num

	userInfo, _ := GetUserInfo(post.UserID)
	post.UserInfo = userInfo
	return post, nil
}

//AutoSubPostView 文章详情浏览加一操作
func AutoSubPostView(postID int) int64 {
	var (
		err       error
		affectRow int64
	)
	affectRow = 0
	stmt, err := dbConn.Prepare(`update wps_post_nums set view_num = view_num + 1 where post_id = ?`)
	defer stmt.Close()
	if err == nil {
		res, err := stmt.Exec(postID)
		if err == nil {
			num, err := res.RowsAffected()
			if err == nil {
				return num
			}
		}
	}
	return affectRow
}
