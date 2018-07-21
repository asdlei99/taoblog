package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"./internal/post_translators"
	"./internal/utils/datetime"
)

type PostForArchiveQuery struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

type xPostManager struct {
	db *sql.DB
}

func newPostManager(db *sql.DB) *xPostManager {
	return &xPostManager{
		db: db,
	}
}

// returns nil if post presents
func (me *xPostManager) has(id int64) error {
	query := `SELECT id FROM posts WHERE id=` + strconv.FormatInt(id, 10)
	rows := me.db.QueryRow(query)
	pid := 0
	err := rows.Scan(&pid)
	return err
}

func (me *xPostManager) update(id int64, typ string, source string) error {
	var tr post_translators.PostTranslator
	var content string
	var err error

	switch typ {
	case "html":
		tr = &post_translators.HTMLTranslator{}
	case "markdown":
		tr = &post_translators.MarkdownTranslator{}
	}

	if tr == nil {
		return errors.New("no translator found for " + typ)
	}

	content, err = tr.Translate(source)
	if err != nil {
		return err
	}

	modTime := time.Now().UTC().Format("2006:01:02 15:04:05")

	ret, err := me.db.Exec(
		"UPDATE posts SET content=?,source=?,source_type=?,modified=? WHERE id=? LIMIT 1",
		content,
		source,
		typ,
		modTime,
		id,
	)

	if err != nil {
		return err
	}

	/*
		if n, err := ret.RowsAffected(); err != nil || n != 1 {
			return errors.New("affected rows != 1: n=" + strconv.FormatInt(n, 10))
		}
	*/

	_ = ret

	return nil
}

func (me *xPostManager) getCommentCount(pid int64) (count int) {
	query := `SELECT comments FROM posts WHERE id=` + fmt.Sprint(pid) + ` LIMIT 1`
	row := me.db.QueryRow(query)
	switch row.Scan(&count) {
	case sql.ErrNoRows:
		count = -1
	case nil:
		break
	default:
		count = -1
	}
	return
}

func (z *xPostManager) BeforeQuery(q map[string]interface{}) map[string]interface{} {
	if _, ok := q["where"]; !ok {
		q["where"] = []string{}
	}
	ws := q["where"].([]string)
	ws = append(ws, "status='public'")
	q["where"] = ws
	return q
}

func (z *xPostManager) getRowPosts(q map[string]interface{}) ([]*PostForArchiveQuery, error) {
	q = z.BeforeQuery(q)
	s := BuildQueryString(q)
	log.Println(s)

	rows, err := z.db.Query(s)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	ps := make([]*PostForArchiveQuery, 0)

	for rows.Next() {
		p := &PostForArchiveQuery{}
		if err = rows.Scan(&p.ID, &p.Title); err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}

	return ps, rows.Err()
}

func (z *xPostManager) GetPostsByCategory(catID int64) ([]*PostForArchiveQuery, error) {
	q := make(map[string]interface{})

	q["select"] = "id,title"
	q["from"] = "posts"
	q["where"] = []string{
		fmt.Sprintf("taxonomy=%d", catID),
		"type='post'",
	}
	q["orderby"] = "date DESC"

	return z.getRowPosts(q)
}

func (z *xPostManager) GetPostsByTags(tag string) ([]*PostForArchiveQuery, error) {
	id := tagmgr.getTagID(tag)
	ids := tagmgr.getAliasTagsAll([]int64{id})

	q := make(map[string]interface{})

	q["select"] = "posts.id,posts.title"
	q["from"] = "posts,post_tags"
	q["where"] = []string{
		"posts.id=post_tags.post_id",
		fmt.Sprintf("post_tags.tag_id in (%s)", joinInts(ids, ",")),
	}

	return z.getRowPosts(q)
}

func (z *xPostManager) GetPostsByDate(yy, mm int64) ([]*PostForArchiveQuery, error) {
	q := make(map[string]interface{})
	q["select"] = "id,title"
	q["from"] = "posts"
	q["where"] = []string{
		"type='post'",
	}
	if yy > 1970 {
		var start, end string
		if 1 <= mm && mm <= 12 {
			start, end = datetime.MonthStartEnd(int(yy), int(mm))
		} else {
			start, end = datetime.YearStartEnd(int(yy))
		}
		q["where"] = append(q["where"].([]string), fmt.Sprintf("date>='%s' AND date<='%s'", start, end))
	}

	q["orderby"] = "date DESC"

	return z.getRowPosts(q)
}

func (z *xPostManager) GetVars(fields string, wheres string, outs ...interface{}) error {
	q := make(map[string]interface{})
	q["select"] = fields
	q["from"] = "posts"
	q["where"] = []string{
		wheres,
	}
	q["limit"] = 1

	query := BuildQueryString(q)

	row := z.db.QueryRow(query)

	return row.Scan(outs...)
}
