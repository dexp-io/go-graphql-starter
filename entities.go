// generated by dexp.io DO NOT EDIT
package dexp

import (
	
		"database/sql"
	
		"encoding/json"
	
		"errors"
	
		"github.com/go-sql-driver/mysql"
	
		"strconv"
	
		"strings"
	
		"time"
	
)




// permission

type Permission struct {
	ID        int64     `json:"id"`
	Type      string    `json:"type"`
	AuthorID  int64     `json:"author_id"`
	LastUpdatedBy int64 `json:"last_updated_by"`
	Status    int8      `json:"status"`
	Deleted   bool      `json:"deleted"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	GrantUserId int64 `json:"grant_user_id"`
	
	GrantEntityId int64 `json:"grant_entity_id"`
	
	GrantEntityType string `json:"grant_entity_type"`
	
	CanView bool `json:"can_view"`
	
	CanInsert bool `json:"can_insert"`
	
	CanUpdate bool `json:"can_update"`
	
	CanDelete bool `json:"can_delete"`
	

}


func (p *Permission) IsEntity() {}



func (p *Permission) ToString() string {

	data, _ := json.Marshal(p)
	return string(data)
}

type PermissionFindQuery struct {
	_where   []string
	_orderBy []string
	_args    []interface{}
	_limit   int
	_offset  int
}

func (q *PermissionFindQuery) Where(condition string, args ...interface{}) *PermissionFindQuery {
	q._where = append(q._where, condition)

	if len(args) > 0 {
		for _, arg := range args {
			q._args = append(q._args, arg)
		}
	}
	return q
}

func (q *PermissionFindQuery) OrderBy(field string, value string) *PermissionFindQuery {

	value = strings.ToUpper(value)
	if value != "" && (value == "DESC" || value == "ASC") {
		q._orderBy = append(q._orderBy, field+" "+value)
	} else {
		q._orderBy = append(q._orderBy, field)
	}

	return q
}

func (q *PermissionFindQuery) Range(limit, offset int) *PermissionFindQuery {

	q._limit = limit
	q._offset = offset

	return q
}

func (q *PermissionFindQuery) Execute() ([]*Permission, error) {

	s := "SELECT entity.*, p.grant_user_id, p.grant_entity_id, p.grant_entity_type, p.can_view, p.can_insert, p.can_update, p.can_delete FROM entity_permission AS p INNER JOIN `entity` ON p.entity_id = entity.id"

	if len(q._where) > 0 {
		s += " WHERE " + strings.Join(q._where, " ")
	}

	if len(q._orderBy) > 0 {
		s += " ORDER BY " + strings.Join(q._orderBy, ", ")
	}
	if q._limit > 0 {
		s += " LIMIT " + strconv.Itoa(q._limit) + " OFFSET " + strconv.Itoa(q._offset)
	}

	rows, err := DB.Query(s, q._args...)



	if err != nil {
		return nil, err
	}

	defer rows.Close()


	var res []*Permission

	for rows.Next() {
		var p Permission

		if rows.Scan(&p.ID, &p.Type, &p.AuthorID, &p.Status, &p.CreatedAt, &p.UpdatedAt, &p.GrantUserId, &p.GrantEntityId, &p.GrantEntityType, &p.CanView, &p.CanInsert, &p.CanUpdate, &p.CanDelete) == nil {
			res = append(res, &p)
		}
	}

	return res, nil

}

func FindPermissions() *PermissionFindQuery {
	return &PermissionFindQuery{}
}


func (p *Permission) Get(ID int64) error {
	if ID == 0 {
		return errors.New("ID is required")
	}
	err := DB.GetEntityCache(ID, p)
	if err == nil {
		return nil
	}

	p.ID = ID
	row := DB.Select("entity", "e").Join("entity_permission", "p", "p.entity_id = e.id").Fields("p", []string{
	
		"grant_user_id",
		
		"grant_entity_id",
		
		"grant_entity_type",
		
		"can_view",
		
		"can_insert",
		
		"can_update",
		
		"can_delete",
		
}).Condition("e.id", p.ID, "=").FetchOne()
	err = row.Scan(&p.ID, &p.Type, &p.AuthorID, &p.Status, &p.CreatedAt, &p.UpdatedAt, &p.GrantUserId, &p.GrantEntityId, &p.GrantEntityType, &p.CanView, &p.CanInsert, &p.CanUpdate, &p.CanDelete)
	if err == nil {
		DB.SetEntityCache(ID, p)
	}
	return err
}


func (p * Permission) Insert() error {

	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	p.Type = "permission"
	p.LastUpdatedBy = p.AuthorID

	res, err := tx.Exec("INSERT INTO `entity` (`type`, `author_id`, `status`, `created_at`, `updated_at`) VALUES (?, ?, ?, ?, ?)",
		p.Type, p.AuthorID, p.Status, p.CreatedAt, p.UpdatedAt)

	if err != nil {

		tx.Rollback()
		return err
	}

	id, err := res.LastInsertId()

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	p.ID = id

	query := "INSERT INTO `entity_permission` (entity_id, grant_user_id, grant_entity_id, grant_entity_type, can_view, can_insert, can_update, can_delete) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"

	_, err = tx.Exec(query, p.ID, p.GrantUserId, p.GrantEntityId, p.GrantEntityType, p.CanView, p.CanInsert, p.CanUpdate, p.CanDelete)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()

	if err != nil{
		return err
	}

	_ = DB.SetEntityCache(id, p)

	activity := &Activity{
		UserID:     p.AuthorID,
		EntityID:   p.ID,
		EntityType: p.Type,
		Action:     "created",
		Payload:    "",
	}
	_ = activity.Create()


	
	if p.Type != "permission" && p.AuthorID > 0 {

		perm := &Permission{
			AuthorID:        p.AuthorID,
			GrantUserId:     p.AuthorID,
			GrantEntityId:   p.ID,
			GrantEntityType: p.Type,
			CanView: 		 true,
			CanInsert:       true,
			CanUpdate:       true,
			CanDelete:       true,
		}
		_ = perm.Insert()
	}



	return nil

}



func (p *Permission) Update() error {

	tx, err := DB.Begin()

	if err != nil {
		return err
	}

	p.UpdatedAt = time.Now()

	_, err = tx.Exec("UPDATE entity SET updated_at = ? WHERE id = ?", p.UpdatedAt, p.ID)

	if err != nil {

		_ = tx.Rollback()

		return err
	}

	res, err := tx.Exec("UPDATE entity_permission SET grant_user_id = ? , grant_entity_id = ? , grant_entity_type = ? , can_view = ? , can_insert = ? , can_update = ? , can_delete WHERE entity_id = ?", p.GrantUserId, p.GrantEntityId, p.GrantEntityType, p.CanView, p.CanInsert, p.CanUpdate, p.CanDelete, p.ID)
	if err != nil {
		tx.Rollback()

		return err
	}


	err =  tx.Commit()

	if err != nil {
		return err
	}

	effectIds , _ := res.LastInsertId()

	if effectIds > 0{
			activity := &Activity{
			UserID:     p.LastUpdatedBy,
			EntityID:   p.ID,
			EntityType: p.Type,
			Action:     "updated",
			Payload:    "",
		}
		_ = activity.Create()
		DB.SetEntityCache(p.ID, p)
	}


	return nil

}



func (p *Permission) Delete() error {

	p.Deleted = true
	res, err := DB.Exec("UPDATE `entity` SET deleted = 1 WHERE id = ?", p.ID)



	if err != nil{
		return nil
	}

	effectIds , _ := res.LastInsertId()

	if effectIds > 0{
		if effectIds > 0{
			activity := &Activity{
                UserID:     p.LastUpdatedBy,
				EntityID:   p.ID,
				EntityType: p.Type,
				Action:     "deleted",
				Payload:    "",
			}
			_ = activity.Create()
			DB.SetEntityCache(p.ID, p)
		}
	}

	return err
}


func (p *Permission) AddSubscriber(userID int64) error {

	_, err := DB.Insert("entity_has_subscriber").Fields(map[string]interface{}{
		"entity_id": p.ID,
		"user_id":   userID,
	}).Execute()

	return err
}

// end permission


// project

type Project struct {
	ID        int64     `json:"id"`
	Type      string    `json:"type"`
	AuthorID  int64     `json:"author_id"`
	LastUpdatedBy int64 `json:"last_updated_by"`
	Status    int8      `json:"status"`
	Deleted   bool      `json:"deleted"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	Title string `json:"title"`
	
	Body string `json:"body"`
	

}


func (p *Project) IsEntity() {}



func (p *Project) ToString() string {

	data, _ := json.Marshal(p)
	return string(data)
}

type ProjectFindQuery struct {
	_where   []string
	_orderBy []string
	_args    []interface{}
	_limit   int
	_offset  int
}

func (q *ProjectFindQuery) Where(condition string, args ...interface{}) *ProjectFindQuery {
	q._where = append(q._where, condition)

	if len(args) > 0 {
		for _, arg := range args {
			q._args = append(q._args, arg)
		}
	}
	return q
}

func (q *ProjectFindQuery) OrderBy(field string, value string) *ProjectFindQuery {

	value = strings.ToUpper(value)
	if value != "" && (value == "DESC" || value == "ASC") {
		q._orderBy = append(q._orderBy, field+" "+value)
	} else {
		q._orderBy = append(q._orderBy, field)
	}

	return q
}

func (q *ProjectFindQuery) Range(limit, offset int) *ProjectFindQuery {

	q._limit = limit
	q._offset = offset

	return q
}

func (q *ProjectFindQuery) Execute() ([]*Project, error) {

	s := "SELECT entity.*, p.title, p.body FROM entity_project AS p INNER JOIN `entity` ON p.entity_id = entity.id"

	if len(q._where) > 0 {
		s += " WHERE " + strings.Join(q._where, " ")
	}

	if len(q._orderBy) > 0 {
		s += " ORDER BY " + strings.Join(q._orderBy, ", ")
	}
	if q._limit > 0 {
		s += " LIMIT " + strconv.Itoa(q._limit) + " OFFSET " + strconv.Itoa(q._offset)
	}

	rows, err := DB.Query(s, q._args...)



	if err != nil {
		return nil, err
	}

	defer rows.Close()


	var res []*Project

	for rows.Next() {
		var p Project

		if rows.Scan(&p.ID, &p.Type, &p.AuthorID, &p.Status, &p.CreatedAt, &p.UpdatedAt, &p.Title, &p.Body) == nil {
			res = append(res, &p)
		}
	}

	return res, nil

}

func FindProjects() *ProjectFindQuery {
	return &ProjectFindQuery{}
}


func (p *Project) Get(ID int64) error {
	if ID == 0 {
		return errors.New("ID is required")
	}
	err := DB.GetEntityCache(ID, p)
	if err == nil {
		return nil
	}

	p.ID = ID
	row := DB.Select("entity", "e").Join("entity_project", "p", "p.entity_id = e.id").Fields("p", []string{
	
		"title",
		
		"body",
		
}).Condition("e.id", p.ID, "=").FetchOne()
	err = row.Scan(&p.ID, &p.Type, &p.AuthorID, &p.Status, &p.CreatedAt, &p.UpdatedAt, &p.Title, &p.Body)
	if err == nil {
		DB.SetEntityCache(ID, p)
	}
	return err
}


func (p * Project) Insert() error {

	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	p.Type = "project"
	p.LastUpdatedBy = p.AuthorID

	res, err := tx.Exec("INSERT INTO `entity` (`type`, `author_id`, `status`, `created_at`, `updated_at`) VALUES (?, ?, ?, ?, ?)",
		p.Type, p.AuthorID, p.Status, p.CreatedAt, p.UpdatedAt)

	if err != nil {

		tx.Rollback()
		return err
	}

	id, err := res.LastInsertId()

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	p.ID = id

	query := "INSERT INTO `entity_project` (entity_id, title, body) VALUES (?, ?, ?)"

	_, err = tx.Exec(query, p.ID, p.Title, p.Body)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()

	if err != nil{
		return err
	}

	_ = DB.SetEntityCache(id, p)

	activity := &Activity{
		UserID:     p.AuthorID,
		EntityID:   p.ID,
		EntityType: p.Type,
		Action:     "created",
		Payload:    "",
	}
	_ = activity.Create()


	
	if p.Type != "permission" && p.AuthorID > 0 {

		perm := &Permission{
			AuthorID:        p.AuthorID,
			GrantUserId:     p.AuthorID,
			GrantEntityId:   p.ID,
			GrantEntityType: p.Type,
			CanView: 		 true,
			CanInsert:       true,
			CanUpdate:       true,
			CanDelete:       true,
		}
		_ = perm.Insert()
	}



	return nil

}



func (p *Project) Update() error {

	tx, err := DB.Begin()

	if err != nil {
		return err
	}

	p.UpdatedAt = time.Now()

	_, err = tx.Exec("UPDATE entity SET updated_at = ? WHERE id = ?", p.UpdatedAt, p.ID)

	if err != nil {

		_ = tx.Rollback()

		return err
	}

	res, err := tx.Exec("UPDATE entity_project SET title = ? , body WHERE entity_id = ?", p.Title, p.Body, p.ID)
	if err != nil {
		tx.Rollback()

		return err
	}


	err =  tx.Commit()

	if err != nil {
		return err
	}

	effectIds , _ := res.LastInsertId()

	if effectIds > 0{
			activity := &Activity{
			UserID:     p.LastUpdatedBy,
			EntityID:   p.ID,
			EntityType: p.Type,
			Action:     "updated",
			Payload:    "",
		}
		_ = activity.Create()
		DB.SetEntityCache(p.ID, p)
	}


	return nil

}



func (p *Project) Delete() error {

	p.Deleted = true
	res, err := DB.Exec("UPDATE `entity` SET deleted = 1 WHERE id = ?", p.ID)



	if err != nil{
		return nil
	}

	effectIds , _ := res.LastInsertId()

	if effectIds > 0{
		if effectIds > 0{
			activity := &Activity{
                UserID:     p.LastUpdatedBy,
				EntityID:   p.ID,
				EntityType: p.Type,
				Action:     "deleted",
				Payload:    "",
			}
			_ = activity.Create()
			DB.SetEntityCache(p.ID, p)
		}
	}

	return err
}


func (p *Project) AddSubscriber(userID int64) error {

	_, err := DB.Insert("entity_has_subscriber").Fields(map[string]interface{}{
		"entity_id": p.ID,
		"user_id":   userID,
	}).Execute()

	return err
}

// end project


// board

type Board struct {
	ID        int64     `json:"id"`
	Type      string    `json:"type"`
	AuthorID  int64     `json:"author_id"`
	LastUpdatedBy int64 `json:"last_updated_by"`
	Status    int8      `json:"status"`
	Deleted   bool      `json:"deleted"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	ProjectId int64 `json:"project_id"`
	
	Title string `json:"title"`
	
	Weight int64 `json:"weight"`
	

}


func (b *Board) IsEntity() {}



func (b *Board) ToString() string {

	data, _ := json.Marshal(b)
	return string(data)
}

type BoardFindQuery struct {
	_where   []string
	_orderBy []string
	_args    []interface{}
	_limit   int
	_offset  int
}

func (q *BoardFindQuery) Where(condition string, args ...interface{}) *BoardFindQuery {
	q._where = append(q._where, condition)

	if len(args) > 0 {
		for _, arg := range args {
			q._args = append(q._args, arg)
		}
	}
	return q
}

func (q *BoardFindQuery) OrderBy(field string, value string) *BoardFindQuery {

	value = strings.ToUpper(value)
	if value != "" && (value == "DESC" || value == "ASC") {
		q._orderBy = append(q._orderBy, field+" "+value)
	} else {
		q._orderBy = append(q._orderBy, field)
	}

	return q
}

func (q *BoardFindQuery) Range(limit, offset int) *BoardFindQuery {

	q._limit = limit
	q._offset = offset

	return q
}

func (q *BoardFindQuery) Execute() ([]*Board, error) {

	s := "SELECT entity.*, b.project_id, b.title, b.weight FROM entity_board AS b INNER JOIN `entity` ON b.entity_id = entity.id"

	if len(q._where) > 0 {
		s += " WHERE " + strings.Join(q._where, " ")
	}

	if len(q._orderBy) > 0 {
		s += " ORDER BY " + strings.Join(q._orderBy, ", ")
	}
	if q._limit > 0 {
		s += " LIMIT " + strconv.Itoa(q._limit) + " OFFSET " + strconv.Itoa(q._offset)
	}

	rows, err := DB.Query(s, q._args...)



	if err != nil {
		return nil, err
	}

	defer rows.Close()


	var res []*Board

	for rows.Next() {
		var b Board

		if rows.Scan(&b.ID, &b.Type, &b.AuthorID, &b.Status, &b.CreatedAt, &b.UpdatedAt, &b.ProjectId, &b.Title, &b.Weight) == nil {
			res = append(res, &b)
		}
	}

	return res, nil

}

func FindBoards() *BoardFindQuery {
	return &BoardFindQuery{}
}


func (b *Board) Get(ID int64) error {
	if ID == 0 {
		return errors.New("ID is required")
	}
	err := DB.GetEntityCache(ID, b)
	if err == nil {
		return nil
	}

	b.ID = ID
	row := DB.Select("entity", "e").Join("entity_board", "b", "b.entity_id = e.id").Fields("b", []string{
	
		"project_id",
		
		"title",
		
		"weight",
		
}).Condition("e.id", b.ID, "=").FetchOne()
	err = row.Scan(&b.ID, &b.Type, &b.AuthorID, &b.Status, &b.CreatedAt, &b.UpdatedAt, &b.ProjectId, &b.Title, &b.Weight)
	if err == nil {
		DB.SetEntityCache(ID, b)
	}
	return err
}


func (b * Board) Insert() error {

	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	b.Type = "board"
	b.LastUpdatedBy = b.AuthorID

	res, err := tx.Exec("INSERT INTO `entity` (`type`, `author_id`, `status`, `created_at`, `updated_at`) VALUES (?, ?, ?, ?, ?)",
		b.Type, b.AuthorID, b.Status, b.CreatedAt, b.UpdatedAt)

	if err != nil {

		tx.Rollback()
		return err
	}

	id, err := res.LastInsertId()

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	b.ID = id

	query := "INSERT INTO `entity_board` (entity_id, project_id, title, weight) VALUES (?, ?, ?, ?)"

	_, err = tx.Exec(query, b.ID, b.ProjectId, b.Title, b.Weight)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()

	if err != nil{
		return err
	}

	_ = DB.SetEntityCache(id, b)

	activity := &Activity{
		UserID:     b.AuthorID,
		EntityID:   b.ID,
		EntityType: b.Type,
		Action:     "created",
		Payload:    "",
	}
	_ = activity.Create()


	
	if b.Type != "permission" && b.AuthorID > 0 {

		perm := &Permission{
			AuthorID:        b.AuthorID,
			GrantUserId:     b.AuthorID,
			GrantEntityId:   b.ID,
			GrantEntityType: b.Type,
			CanView: 		 true,
			CanInsert:       true,
			CanUpdate:       true,
			CanDelete:       true,
		}
		_ = perm.Insert()
	}



	return nil

}



func (b *Board) Update() error {

	tx, err := DB.Begin()

	if err != nil {
		return err
	}

	b.UpdatedAt = time.Now()

	_, err = tx.Exec("UPDATE entity SET updated_at = ? WHERE id = ?", b.UpdatedAt, b.ID)

	if err != nil {

		_ = tx.Rollback()

		return err
	}

	res, err := tx.Exec("UPDATE entity_board SET project_id = ? , title = ? , weight WHERE entity_id = ?", b.ProjectId, b.Title, b.Weight, b.ID)
	if err != nil {
		tx.Rollback()

		return err
	}


	err =  tx.Commit()

	if err != nil {
		return err
	}

	effectIds , _ := res.LastInsertId()

	if effectIds > 0{
			activity := &Activity{
			UserID:     b.LastUpdatedBy,
			EntityID:   b.ID,
			EntityType: b.Type,
			Action:     "updated",
			Payload:    "",
		}
		_ = activity.Create()
		DB.SetEntityCache(b.ID, b)
	}


	return nil

}



func (b *Board) Delete() error {

	b.Deleted = true
	res, err := DB.Exec("UPDATE `entity` SET deleted = 1 WHERE id = ?", b.ID)



	if err != nil{
		return nil
	}

	effectIds , _ := res.LastInsertId()

	if effectIds > 0{
		if effectIds > 0{
			activity := &Activity{
                UserID:     b.LastUpdatedBy,
				EntityID:   b.ID,
				EntityType: b.Type,
				Action:     "deleted",
				Payload:    "",
			}
			_ = activity.Create()
			DB.SetEntityCache(b.ID, b)
		}
	}

	return err
}


func (b *Board) AddSubscriber(userID int64) error {

	_, err := DB.Insert("entity_has_subscriber").Fields(map[string]interface{}{
		"entity_id": b.ID,
		"user_id":   userID,
	}).Execute()

	return err
}

// end board


// room

type Room struct {
	ID        int64     `json:"id"`
	Type      string    `json:"type"`
	AuthorID  int64     `json:"author_id"`
	LastUpdatedBy int64 `json:"last_updated_by"`
	Status    int8      `json:"status"`
	Deleted   bool      `json:"deleted"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	Title string `json:"title"`
	
	AvatarFileId int64 `json:"avatar_file_id"`
	

}


func (r *Room) IsEntity() {}



func (r *Room) ToString() string {

	data, _ := json.Marshal(r)
	return string(data)
}

type RoomFindQuery struct {
	_where   []string
	_orderBy []string
	_args    []interface{}
	_limit   int
	_offset  int
}

func (q *RoomFindQuery) Where(condition string, args ...interface{}) *RoomFindQuery {
	q._where = append(q._where, condition)

	if len(args) > 0 {
		for _, arg := range args {
			q._args = append(q._args, arg)
		}
	}
	return q
}

func (q *RoomFindQuery) OrderBy(field string, value string) *RoomFindQuery {

	value = strings.ToUpper(value)
	if value != "" && (value == "DESC" || value == "ASC") {
		q._orderBy = append(q._orderBy, field+" "+value)
	} else {
		q._orderBy = append(q._orderBy, field)
	}

	return q
}

func (q *RoomFindQuery) Range(limit, offset int) *RoomFindQuery {

	q._limit = limit
	q._offset = offset

	return q
}

func (q *RoomFindQuery) Execute() ([]*Room, error) {

	s := "SELECT entity.*, r.title, r.avatar_file_id FROM entity_room AS r INNER JOIN `entity` ON r.entity_id = entity.id"

	if len(q._where) > 0 {
		s += " WHERE " + strings.Join(q._where, " ")
	}

	if len(q._orderBy) > 0 {
		s += " ORDER BY " + strings.Join(q._orderBy, ", ")
	}
	if q._limit > 0 {
		s += " LIMIT " + strconv.Itoa(q._limit) + " OFFSET " + strconv.Itoa(q._offset)
	}

	rows, err := DB.Query(s, q._args...)



	if err != nil {
		return nil, err
	}

	defer rows.Close()


	var res []*Room

	for rows.Next() {
		var r Room

		if rows.Scan(&r.ID, &r.Type, &r.AuthorID, &r.Status, &r.CreatedAt, &r.UpdatedAt, &r.Title, &r.AvatarFileId) == nil {
			res = append(res, &r)
		}
	}

	return res, nil

}

func FindRooms() *RoomFindQuery {
	return &RoomFindQuery{}
}


func (r *Room) Get(ID int64) error {
	if ID == 0 {
		return errors.New("ID is required")
	}
	err := DB.GetEntityCache(ID, r)
	if err == nil {
		return nil
	}

	r.ID = ID
	row := DB.Select("entity", "e").Join("entity_room", "r", "r.entity_id = e.id").Fields("r", []string{
	
		"title",
		
		"avatar_file_id",
		
}).Condition("e.id", r.ID, "=").FetchOne()
	err = row.Scan(&r.ID, &r.Type, &r.AuthorID, &r.Status, &r.CreatedAt, &r.UpdatedAt, &r.Title, &r.AvatarFileId)
	if err == nil {
		DB.SetEntityCache(ID, r)
	}
	return err
}


func (r * Room) Insert() error {

	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	r.Type = "room"
	r.LastUpdatedBy = r.AuthorID

	res, err := tx.Exec("INSERT INTO `entity` (`type`, `author_id`, `status`, `created_at`, `updated_at`) VALUES (?, ?, ?, ?, ?)",
		r.Type, r.AuthorID, r.Status, r.CreatedAt, r.UpdatedAt)

	if err != nil {

		tx.Rollback()
		return err
	}

	id, err := res.LastInsertId()

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	r.ID = id

	query := "INSERT INTO `entity_room` (entity_id, title, avatar_file_id) VALUES (?, ?, ?)"

	_, err = tx.Exec(query, r.ID, r.Title, r.AvatarFileId)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()

	if err != nil{
		return err
	}

	_ = DB.SetEntityCache(id, r)

	activity := &Activity{
		UserID:     r.AuthorID,
		EntityID:   r.ID,
		EntityType: r.Type,
		Action:     "created",
		Payload:    "",
	}
	_ = activity.Create()


	
	if r.Type != "permission" && r.AuthorID > 0 {

		perm := &Permission{
			AuthorID:        r.AuthorID,
			GrantUserId:     r.AuthorID,
			GrantEntityId:   r.ID,
			GrantEntityType: r.Type,
			CanView: 		 true,
			CanInsert:       true,
			CanUpdate:       true,
			CanDelete:       true,
		}
		_ = perm.Insert()
	}



	return nil

}



func (r *Room) Update() error {

	tx, err := DB.Begin()

	if err != nil {
		return err
	}

	r.UpdatedAt = time.Now()

	_, err = tx.Exec("UPDATE entity SET updated_at = ? WHERE id = ?", r.UpdatedAt, r.ID)

	if err != nil {

		_ = tx.Rollback()

		return err
	}

	res, err := tx.Exec("UPDATE entity_room SET title = ? , avatar_file_id WHERE entity_id = ?", r.Title, r.AvatarFileId, r.ID)
	if err != nil {
		tx.Rollback()

		return err
	}


	err =  tx.Commit()

	if err != nil {
		return err
	}

	effectIds , _ := res.LastInsertId()

	if effectIds > 0{
			activity := &Activity{
			UserID:     r.LastUpdatedBy,
			EntityID:   r.ID,
			EntityType: r.Type,
			Action:     "updated",
			Payload:    "",
		}
		_ = activity.Create()
		DB.SetEntityCache(r.ID, r)
	}


	return nil

}



func (r *Room) Delete() error {

	r.Deleted = true
	res, err := DB.Exec("UPDATE `entity` SET deleted = 1 WHERE id = ?", r.ID)



	if err != nil{
		return nil
	}

	effectIds , _ := res.LastInsertId()

	if effectIds > 0{
		if effectIds > 0{
			activity := &Activity{
                UserID:     r.LastUpdatedBy,
				EntityID:   r.ID,
				EntityType: r.Type,
				Action:     "deleted",
				Payload:    "",
			}
			_ = activity.Create()
			DB.SetEntityCache(r.ID, r)
		}
	}

	return err
}


func (r *Room) AddSubscriber(userID int64) error {

	_, err := DB.Insert("entity_has_subscriber").Fields(map[string]interface{}{
		"entity_id": r.ID,
		"user_id":   userID,
	}).Execute()

	return err
}

// end room


// project_member

type ProjectMember struct {
	ID        int64     `json:"id"`
	Type      string    `json:"type"`
	AuthorID  int64     `json:"author_id"`
	LastUpdatedBy int64 `json:"last_updated_by"`
	Status    int8      `json:"status"`
	Deleted   bool      `json:"deleted"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	UserId int64 `json:"user_id"`
	
	ProjectId int64 `json:"project_id"`
	

}


func (p *ProjectMember) IsEntity() {}



func (p *ProjectMember) ToString() string {

	data, _ := json.Marshal(p)
	return string(data)
}

type ProjectMemberFindQuery struct {
	_where   []string
	_orderBy []string
	_args    []interface{}
	_limit   int
	_offset  int
}

func (q *ProjectMemberFindQuery) Where(condition string, args ...interface{}) *ProjectMemberFindQuery {
	q._where = append(q._where, condition)

	if len(args) > 0 {
		for _, arg := range args {
			q._args = append(q._args, arg)
		}
	}
	return q
}

func (q *ProjectMemberFindQuery) OrderBy(field string, value string) *ProjectMemberFindQuery {

	value = strings.ToUpper(value)
	if value != "" && (value == "DESC" || value == "ASC") {
		q._orderBy = append(q._orderBy, field+" "+value)
	} else {
		q._orderBy = append(q._orderBy, field)
	}

	return q
}

func (q *ProjectMemberFindQuery) Range(limit, offset int) *ProjectMemberFindQuery {

	q._limit = limit
	q._offset = offset

	return q
}

func (q *ProjectMemberFindQuery) Execute() ([]*ProjectMember, error) {

	s := "SELECT entity.*, p.user_id, p.project_id FROM entity_project_member AS p INNER JOIN `entity` ON p.entity_id = entity.id"

	if len(q._where) > 0 {
		s += " WHERE " + strings.Join(q._where, " ")
	}

	if len(q._orderBy) > 0 {
		s += " ORDER BY " + strings.Join(q._orderBy, ", ")
	}
	if q._limit > 0 {
		s += " LIMIT " + strconv.Itoa(q._limit) + " OFFSET " + strconv.Itoa(q._offset)
	}

	rows, err := DB.Query(s, q._args...)



	if err != nil {
		return nil, err
	}

	defer rows.Close()


	var res []*ProjectMember

	for rows.Next() {
		var p ProjectMember

		if rows.Scan(&p.ID, &p.Type, &p.AuthorID, &p.Status, &p.CreatedAt, &p.UpdatedAt, &p.UserId, &p.ProjectId) == nil {
			res = append(res, &p)
		}
	}

	return res, nil

}

func FindProjectMembers() *ProjectMemberFindQuery {
	return &ProjectMemberFindQuery{}
}


func (p *ProjectMember) Get(ID int64) error {
	if ID == 0 {
		return errors.New("ID is required")
	}
	err := DB.GetEntityCache(ID, p)
	if err == nil {
		return nil
	}

	p.ID = ID
	row := DB.Select("entity", "e").Join("entity_project_member", "p", "p.entity_id = e.id").Fields("p", []string{
	
		"user_id",
		
		"project_id",
		
}).Condition("e.id", p.ID, "=").FetchOne()
	err = row.Scan(&p.ID, &p.Type, &p.AuthorID, &p.Status, &p.CreatedAt, &p.UpdatedAt, &p.UserId, &p.ProjectId)
	if err == nil {
		DB.SetEntityCache(ID, p)
	}
	return err
}


func (p * ProjectMember) Insert() error {

	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	p.Type = "project_member"
	p.LastUpdatedBy = p.AuthorID

	res, err := tx.Exec("INSERT INTO `entity` (`type`, `author_id`, `status`, `created_at`, `updated_at`) VALUES (?, ?, ?, ?, ?)",
		p.Type, p.AuthorID, p.Status, p.CreatedAt, p.UpdatedAt)

	if err != nil {

		tx.Rollback()
		return err
	}

	id, err := res.LastInsertId()

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	p.ID = id

	query := "INSERT INTO `entity_project_member` (entity_id, user_id, project_id) VALUES (?, ?, ?)"

	_, err = tx.Exec(query, p.ID, p.UserId, p.ProjectId)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()

	if err != nil{
		return err
	}

	_ = DB.SetEntityCache(id, p)

	activity := &Activity{
		UserID:     p.AuthorID,
		EntityID:   p.ID,
		EntityType: p.Type,
		Action:     "created",
		Payload:    "",
	}
	_ = activity.Create()


	
	if p.Type != "permission" && p.AuthorID > 0 {

		perm := &Permission{
			AuthorID:        p.AuthorID,
			GrantUserId:     p.AuthorID,
			GrantEntityId:   p.ID,
			GrantEntityType: p.Type,
			CanView: 		 true,
			CanInsert:       true,
			CanUpdate:       true,
			CanDelete:       true,
		}
		_ = perm.Insert()
	}



	return nil

}



func (p *ProjectMember) Update() error {

	tx, err := DB.Begin()

	if err != nil {
		return err
	}

	p.UpdatedAt = time.Now()

	_, err = tx.Exec("UPDATE entity SET updated_at = ? WHERE id = ?", p.UpdatedAt, p.ID)

	if err != nil {

		_ = tx.Rollback()

		return err
	}

	res, err := tx.Exec("UPDATE entity_project_member SET user_id = ? , project_id WHERE entity_id = ?", p.UserId, p.ProjectId, p.ID)
	if err != nil {
		tx.Rollback()

		return err
	}


	err =  tx.Commit()

	if err != nil {
		return err
	}

	effectIds , _ := res.LastInsertId()

	if effectIds > 0{
			activity := &Activity{
			UserID:     p.LastUpdatedBy,
			EntityID:   p.ID,
			EntityType: p.Type,
			Action:     "updated",
			Payload:    "",
		}
		_ = activity.Create()
		DB.SetEntityCache(p.ID, p)
	}


	return nil

}



func (p *ProjectMember) Delete() error {

	p.Deleted = true
	res, err := DB.Exec("UPDATE `entity` SET deleted = 1 WHERE id = ?", p.ID)



	if err != nil{
		return nil
	}

	effectIds , _ := res.LastInsertId()

	if effectIds > 0{
		if effectIds > 0{
			activity := &Activity{
                UserID:     p.LastUpdatedBy,
				EntityID:   p.ID,
				EntityType: p.Type,
				Action:     "deleted",
				Payload:    "",
			}
			_ = activity.Create()
			DB.SetEntityCache(p.ID, p)
		}
	}

	return err
}


func (p *ProjectMember) AddSubscriber(userID int64) error {

	_, err := DB.Insert("entity_has_subscriber").Fields(map[string]interface{}{
		"entity_id": p.ID,
		"user_id":   userID,
	}).Execute()

	return err
}

// end project_member


// task

type Task struct {
	ID        int64     `json:"id"`
	Type      string    `json:"type"`
	AuthorID  int64     `json:"author_id"`
	LastUpdatedBy int64 `json:"last_updated_by"`
	Status    int8      `json:"status"`
	Deleted   bool      `json:"deleted"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	ParentTaskId int64 `json:"parent_task_id"`
	
	AssigneeUserId sql.NullInt64 `json:"assignee_user_id"`
	
	BoardId int64 `json:"board_id"`
	
	Title string `json:"title"`
	
	Body string `json:"body"`
	
	Weight int64 `json:"weight"`
	
	DueDate mysql.NullTime `json:"due_date"`
	

}


func (t *Task) IsEntity() {}



func (t *Task) ToString() string {

	data, _ := json.Marshal(t)
	return string(data)
}

type TaskFindQuery struct {
	_where   []string
	_orderBy []string
	_args    []interface{}
	_limit   int
	_offset  int
}

func (q *TaskFindQuery) Where(condition string, args ...interface{}) *TaskFindQuery {
	q._where = append(q._where, condition)

	if len(args) > 0 {
		for _, arg := range args {
			q._args = append(q._args, arg)
		}
	}
	return q
}

func (q *TaskFindQuery) OrderBy(field string, value string) *TaskFindQuery {

	value = strings.ToUpper(value)
	if value != "" && (value == "DESC" || value == "ASC") {
		q._orderBy = append(q._orderBy, field+" "+value)
	} else {
		q._orderBy = append(q._orderBy, field)
	}

	return q
}

func (q *TaskFindQuery) Range(limit, offset int) *TaskFindQuery {

	q._limit = limit
	q._offset = offset

	return q
}

func (q *TaskFindQuery) Execute() ([]*Task, error) {

	s := "SELECT entity.*, t.parent_task_id, t.assignee_user_id, t.board_id, t.title, t.body, t.weight, t.due_date FROM entity_task AS t INNER JOIN `entity` ON t.entity_id = entity.id"

	if len(q._where) > 0 {
		s += " WHERE " + strings.Join(q._where, " ")
	}

	if len(q._orderBy) > 0 {
		s += " ORDER BY " + strings.Join(q._orderBy, ", ")
	}
	if q._limit > 0 {
		s += " LIMIT " + strconv.Itoa(q._limit) + " OFFSET " + strconv.Itoa(q._offset)
	}

	rows, err := DB.Query(s, q._args...)



	if err != nil {
		return nil, err
	}

	defer rows.Close()


	var res []*Task

	for rows.Next() {
		var t Task

		if rows.Scan(&t.ID, &t.Type, &t.AuthorID, &t.Status, &t.CreatedAt, &t.UpdatedAt, &t.ParentTaskId, &t.AssigneeUserId, &t.BoardId, &t.Title, &t.Body, &t.Weight, &t.DueDate) == nil {
			res = append(res, &t)
		}
	}

	return res, nil

}

func FindTasks() *TaskFindQuery {
	return &TaskFindQuery{}
}


func (t *Task) Get(ID int64) error {
	if ID == 0 {
		return errors.New("ID is required")
	}
	err := DB.GetEntityCache(ID, t)
	if err == nil {
		return nil
	}

	t.ID = ID
	row := DB.Select("entity", "e").Join("entity_task", "t", "t.entity_id = e.id").Fields("t", []string{
	
		"parent_task_id",
		
		"assignee_user_id",
		
		"board_id",
		
		"title",
		
		"body",
		
		"weight",
		
		"due_date",
		
}).Condition("e.id", t.ID, "=").FetchOne()
	err = row.Scan(&t.ID, &t.Type, &t.AuthorID, &t.Status, &t.CreatedAt, &t.UpdatedAt, &t.ParentTaskId, &t.AssigneeUserId, &t.BoardId, &t.Title, &t.Body, &t.Weight, &t.DueDate)
	if err == nil {
		DB.SetEntityCache(ID, t)
	}
	return err
}


func (t * Task) Insert() error {

	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	t.Type = "task"
	t.LastUpdatedBy = t.AuthorID

	res, err := tx.Exec("INSERT INTO `entity` (`type`, `author_id`, `status`, `created_at`, `updated_at`) VALUES (?, ?, ?, ?, ?)",
		t.Type, t.AuthorID, t.Status, t.CreatedAt, t.UpdatedAt)

	if err != nil {

		tx.Rollback()
		return err
	}

	id, err := res.LastInsertId()

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	t.ID = id

	query := "INSERT INTO `entity_task` (entity_id, parent_task_id, assignee_user_id, board_id, title, body, weight, due_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"

	_, err = tx.Exec(query, t.ID, t.ParentTaskId, t.AssigneeUserId, t.BoardId, t.Title, t.Body, t.Weight, t.DueDate)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()

	if err != nil{
		return err
	}

	_ = DB.SetEntityCache(id, t)

	activity := &Activity{
		UserID:     t.AuthorID,
		EntityID:   t.ID,
		EntityType: t.Type,
		Action:     "created",
		Payload:    "",
	}
	_ = activity.Create()


	
	if t.Type != "permission" && t.AuthorID > 0 {

		perm := &Permission{
			AuthorID:        t.AuthorID,
			GrantUserId:     t.AuthorID,
			GrantEntityId:   t.ID,
			GrantEntityType: t.Type,
			CanView: 		 true,
			CanInsert:       true,
			CanUpdate:       true,
			CanDelete:       true,
		}
		_ = perm.Insert()
	}



	return nil

}



func (t *Task) Update() error {

	tx, err := DB.Begin()

	if err != nil {
		return err
	}

	t.UpdatedAt = time.Now()

	_, err = tx.Exec("UPDATE entity SET updated_at = ? WHERE id = ?", t.UpdatedAt, t.ID)

	if err != nil {

		_ = tx.Rollback()

		return err
	}

	res, err := tx.Exec("UPDATE entity_task SET parent_task_id = ? , assignee_user_id = ? , board_id = ? , title = ? , body = ? , weight = ? , due_date WHERE entity_id = ?", t.ParentTaskId, t.AssigneeUserId, t.BoardId, t.Title, t.Body, t.Weight, t.DueDate, t.ID)
	if err != nil {
		tx.Rollback()

		return err
	}


	err =  tx.Commit()

	if err != nil {
		return err
	}

	effectIds , _ := res.LastInsertId()

	if effectIds > 0{
			activity := &Activity{
			UserID:     t.LastUpdatedBy,
			EntityID:   t.ID,
			EntityType: t.Type,
			Action:     "updated",
			Payload:    "",
		}
		_ = activity.Create()
		DB.SetEntityCache(t.ID, t)
	}


	return nil

}



func (t *Task) Delete() error {

	t.Deleted = true
	res, err := DB.Exec("UPDATE `entity` SET deleted = 1 WHERE id = ?", t.ID)



	if err != nil{
		return nil
	}

	effectIds , _ := res.LastInsertId()

	if effectIds > 0{
		if effectIds > 0{
			activity := &Activity{
                UserID:     t.LastUpdatedBy,
				EntityID:   t.ID,
				EntityType: t.Type,
				Action:     "deleted",
				Payload:    "",
			}
			_ = activity.Create()
			DB.SetEntityCache(t.ID, t)
		}
	}

	return err
}


func (t *Task) AddSubscriber(userID int64) error {

	_, err := DB.Insert("entity_has_subscriber").Fields(map[string]interface{}{
		"entity_id": t.ID,
		"user_id":   userID,
	}).Execute()

	return err
}

// end task
