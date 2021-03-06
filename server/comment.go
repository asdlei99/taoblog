package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"unicode/utf8"

	"./internal/utils/datetime"
)

// hand write regex, not tested well.
var regexpValidEmail = regexp.MustCompile(`^[+-_.a-zA-Z0-9]+@[[:alnum:]]+(\.[[:alnum:]]+)+$`)

type Comment struct {
	ID       int64
	Parent   int64
	Ancestor int64
	PostID   int64
	Author   string
	EMail    string
	URL      string
	IP       string
	Date     string
	Content  string
	Children []*Comment
}

type CommentManager struct {
}

func newCommentManager() *CommentManager {
	return &CommentManager{}
}

func (o *CommentManager) GetAllCount(tx Querier) (count uint) {
	query := `SELECT count(*) as size FROM comments`
	row := tx.QueryRow(query)
	err := row.Scan(&count)
	if err != nil {
		log.Println(err)
	}
	return
}

// DeleteComments deletes comment whose id is id
// It also deletes its children.
func (o *CommentManager) DeleteComments(tx Querier, id int64) error {
	query := fmt.Sprintf(`DELETE FROM comments WHERE id=%d OR ancestor=%d`, id, id)
	_, err := tx.Exec(query)
	return err
}

// GetComment returns the specified comment object.
func (o *CommentManager) GetComment(tx Querier, id int64) (*Comment, error) {
	query := fmt.Sprintf(`SELECT id,parent,ancestor,post_id,author,email,url,ip,date,content FROM comments WHERE id=%d LIMIT 1`, id)
	row := tx.QueryRow(query)
	cmt := &Comment{}
	err := row.Scan(&cmt.ID, &cmt.Parent, &cmt.Ancestor, &cmt.PostID, &cmt.Author, &cmt.EMail, &cmt.URL, &cmt.IP, &cmt.Date, &cmt.Content)
	return cmt, err
}

// GetRecentComments gets the recent comments
// TODO Not tested
func (o *CommentManager) GetRecentComments(tx Querier, num int) ([]*Comment, error) {
	var err error
	query := `SELECT id,parent,ancestor,post_id,author,email,url,ip,date,content FROM comments ORDER BY date DESC LIMIT ` + fmt.Sprint(num)
	rows, err := tx.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var cmts []*Comment

	for rows.Next() {
		var cmt Comment
		if err = rows.Scan(&cmt.ID, &cmt.Parent, &cmt.Ancestor, &cmt.PostID, &cmt.Author, &cmt.EMail, &cmt.URL, &cmt.IP, &cmt.Date, &cmt.Content); err != nil {
			return nil, err
		}
		cmts = append(cmts, &cmt)
	}

	return cmts, rows.Err()
}

// GetChildren gets all children comments of an ancestor
// TODO Not tested
func (o *CommentManager) GetChildren(tx Querier, id int64) ([]*Comment, error) {
	var err error

	query := `SELECT id,parent,ancestor,post_id,author,email,url,ip,date,content FROM comments WHERE ancestor=` + fmt.Sprint(id)
	rows, err := tx.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	cmts := make([]*Comment, 0)

	for rows.Next() {
		var cmt Comment
		if err = rows.Scan(&cmt.ID, &cmt.Parent, &cmt.Ancestor, &cmt.PostID, &cmt.Author, &cmt.EMail, &cmt.URL, &cmt.IP, &cmt.Date, &cmt.Content); err != nil {
			return nil, err
		}
		cmts = append(cmts, &cmt)
	}

	return cmts, rows.Err()
}

// GetAncestor returns the ancestor of a comment
func (o *CommentManager) GetAncestor(tx Querier, id int64) (int64, error) {
	query := `SELECT ancestor FROM comments WHERE id=` + fmt.Sprint(id) + ` LIMIT 1`
	row := tx.QueryRow(query)
	var aid int64
	if err := row.Scan(&aid); err != nil {
		return -1, err
	}

	if aid != 0 {
		return aid, nil
	}

	return 0, nil
}

func (o *CommentManager) GetCommentAndItsChildren(tx Querier, cid int64, offset int64, count int64, pid int64, ascent bool) ([]*Comment, error) {
	var query string

	if cid > 0 {
		query += `SELECT id,parent,ancestor,post_id,author,email,url,ip,date,content FROM comments WHERE id=` + fmt.Sprint(cid)
	} else {
		query += `SELECT id,parent,ancestor,post_id,author,email,url,ip,date,content FROM comments WHERE parent=0`
		if pid > 0 {
			query += ` AND post_id=` + fmt.Sprint(pid)
		}

		if ascent {
			query += ` ORDER BY id ASC`
		} else {
			query += ` ORDER BY id DESC`
		}

		if count > 0 {
			if offset >= 0 {
				query += fmt.Sprintf(" LIMIT %d,%d", offset, count)
			} else {
				query += fmt.Sprintf(" LIMIT %d", count)
			}
		}
	}

	cmts := make([]*Comment, 0)

	rows, err := tx.Query(query)
	if err != nil {
		return cmts, err
	}

	defer rows.Close()

	for rows.Next() {
		var cmt Comment
		if err = rows.Scan(&cmt.ID, &cmt.Parent, &cmt.Ancestor, &cmt.PostID, &cmt.Author, &cmt.EMail, &cmt.URL, &cmt.IP, &cmt.Date, &cmt.Content); err != nil {
			return cmts, err
		}
		cmts = append(cmts, &cmt)
	}

	for _, cmt := range cmts {
		cmt.Children, err = o.GetChildren(tx, cmt.ID)
		if err != nil {
			return nil, err
		}
	}

	return cmts, rows.Err()
}

// BeforeCreate hooks
func (o *CommentManager) beforeCreateComment(tx Querier, c *Comment) error {
	var err error

	// ID
	if c.ID != 0 {
		return errors.New("评论ID必须为0")
	}

	// Ancestor
	if c.Ancestor != 0 {
		return errors.New("不能指定祖先ID")
	}

	// Author
	if len(c.Author) == 0 || utf8.RuneCountInString(c.Author) > 32 {
		return errors.New("昵称不能为空或超出最大长度")
	}

	// Email
	if !regexpValidEmail.MatchString(c.EMail) {
		return errors.New("邮箱不正确")
	}

	// TODO: URL
	c.URL = strings.TrimSpace(c.URL)

	// Content
	if len(c.Content) == 0 || utf8.RuneCountInString(c.Content) > 4096 {
		return errors.New("评论不能为空或超出最大长度")
	}

	// Parent
	if c.Parent > 0 {
		if _, err = o.GetComment(tx, c.Parent); err != nil {
			return err
		}
	}

	return nil
}

// CreateComment creates a comment.
func (o *CommentManager) CreateComment(tx Querier, c *Comment) error {
	var err error

	c.Date = datetime.Local2My(c.Date)
	defer func() {
		c.Date = datetime.My2Local(c.Date)
	}()

	if err = o.beforeCreateComment(tx, c); err != nil {
		return err
	}

	c.Ancestor = 0
	if c.Parent != 0 {
		if c.Ancestor, err = o.GetAncestor(tx, c.Parent); err != nil {
			return err
		}
		if c.Ancestor == 0 {
			c.Ancestor = c.Parent
		}
	}

	query := `INSERT INTO comments (post_id,author,email,url,ip,date,content,parent,ancestor) VALUES (?,?,?,?,?,?,?,?,?)`
	ret, err := tx.Exec(query, c.PostID, c.Author, c.EMail, c.URL, c.IP, c.Date, c.Content, c.Parent, c.Ancestor)
	if err != nil {
		return err
	}

	id, err := ret.LastInsertId()
	c.ID = id
	return err
}

func (o *CommentManager) GetVars(tx Querier, fields string, wheres string, outs ...interface{}) error {
	q := make(map[string]interface{})
	q["select"] = fields
	q["from"] = "comments"
	q["where"] = []string{
		wheres,
	}
	q["limit"] = 1

	query := BuildQueryString(q)

	row := tx.QueryRow(query)

	return row.Scan(outs...)
}
