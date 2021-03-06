package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/zmisgod/blogApi/util"
)

// Comment wps_comments
type Comment struct {
	ID          int `json:"id"`
	userID      int
	Content     string `json:"content"`
	AuthorName  string `json:"author_name"`
	authorEmail string
	AuthorImage string `json:"author_image"`
	AuthorURL   string `json:"author_url"`
	createdAt   int
	CreatedAt   string `json:"created_at"`
}

//getArticleCommentLists 获取文章评论列表
func getArticleCommentLists(postID, page, pageSize int, orderby string) ([]Comment, error) {
	commentList := []Comment{}

	offset := (page - 1) * pageSize
	rows, err := dbConn.Query(fmt.Sprintf("select created_at,author_name,author_url,content,id,user_id from wps_post_comments where post_id = %d order by created_at desc limit %d,%d", postID, offset, pageSize))
	defer rows.Close()
	if err != nil {
		return commentList, err
	}
	for rows.Next() {
		var aComment Comment
		err = rows.Scan(
			&aComment.createdAt,
			&aComment.AuthorName,
			&aComment.AuthorURL,
			&aComment.Content,
			&aComment.ID,
			&aComment.userID,
		)
		if err != nil {
			continue
		}
		tm := time.Unix(int64(aComment.createdAt), 0)
		aComment.CreatedAt = tm.Format("2006-01-02 15:04")
		commentList = append(commentList, aComment)
	}
	return commentList, nil
}

//GetArticleCommentLists 获取文章评论列表
func GetArticleCommentLists(postID, page, pageSize int, orderby string) ([]Comment, error) {
	commentLists, err := getArticleCommentLists(postID, page, pageSize, orderby)
	if err != nil {
		return commentLists, err
	}
	userIDs := ""
	for _, comment := range commentLists {
		if comment.userID != 0 {
			userIDs = strconv.Itoa(comment.userID) + ","
		}
	}
	if userIDs != "" {
		userIDs = strings.TrimSuffix(userIDs, ",")
	}
	if len(userIDs) != 0 {
		userImages, err := GetUserHeadImages(userIDs)
		if err == nil {
			for index, comment := range commentLists {
				if _, ok := userImages[comment.userID]; ok {
					commentLists[index].AuthorImage = userImages[comment.userID]
				}
			}
		}
	}

	return commentLists, err
}

//SaveArticleComment 保存评论
func SaveArticleComment(postID, commentID int, authorName, authorEmail, authorURL, content, authorIP, authorAgent string) Comment {
	//检查评论是否超过240字符
	contentTmp := []rune(content)
	if len(contentTmp) > 240 {
		content = string(contentTmp[:240])
	}

	//检查email是否合法
	if !util.CheckEmail(authorEmail) {
		authorEmail = ""
	}
	var aComment Comment
	commentTimeStamp := time.Now()
	stmt, err := dbConn.Prepare("insert into wps_post_comments (post_id, user_id, content, author_name, author_email, author_url, author_ip, comment_agent, comment_parent, created_at) values (?,?,?,?,?,?,?,?,?,?)")
	defer stmt.Close()
	if err == nil {
		result, err := stmt.Exec(postID, 0, content, authorName, authorEmail, authorURL, authorIP, authorAgent, commentID, commentTimeStamp.Unix())
		if err == nil {
			lastID, err := result.LastInsertId()
			if err == nil {
				//评论 + 1
				addArticleCommentNumber(postID)
				aComment.ID = int(lastID)
				aComment.AuthorImage = ""
				aComment.AuthorName = authorName
				aComment.AuthorURL = authorURL
				aComment.Content = content
				aComment.CreatedAt = commentTimeStamp.Format("2006-01-02 15:04")
			} else {
				fmt.Println(err)
			}
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
	return aComment
}

//addArticleCommentNumber 评论 + 1
func addArticleCommentNumber(articleID int) {
	stmt, err := dbConn.Prepare("update wps_post_nums set comment_num = comment_num + 1 where post_id = ?")
	defer stmt.Close()
	if err == nil {
		_, err = stmt.Exec(articleID)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
}
