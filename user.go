package dexp

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type Role struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID        int64     `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Avatar    string    `json:"avatar"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	IsOnline  bool      `json:"is_online"`
	Verified  bool      `json:"verified"`
	Deleted   bool      `json:"deleted"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserFindQuery struct {
	_where   []string
	_orderBy []string
	_args    []interface{}
	_limit   int
	_offset  int
}

func (q *UserFindQuery) Where(condition string, args ...interface{}) *UserFindQuery {
	q._where = append(q._where, condition)

	if len(args) > 0 {
		for _, arg := range args {
			q._args = append(q._args, arg)
		}
	}
	return q
}

func (q *UserFindQuery) OrderBy(field string, value string) *UserFindQuery {

	value = strings.ToUpper(value)
	if value != "" && (value == "DESC" || value == "ASC") {
		q._orderBy = append(q._orderBy, field+" "+value)
	} else {
		q._orderBy = append(q._orderBy, field)
	}

	return q
}

func (q *UserFindQuery) Range(limit, offset int) *UserFindQuery {

	q._limit = limit
	q._offset = offset

	return q
}

func (q *UserFindQuery) Execute() ([]*User, error) {

	s := "SELECT id, first_name, last_name, avatar, email, password, phone, address, is_online, verified, deleted, created_at, updated_at FROM `user`"

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

	var res []*User

	for rows.Next() {
		var u User

		if rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Avatar, &u.Email, &u.Password, &u.Phone, &u.Address, &u.IsOnline, &u.Verified, &u.Deleted, &u.CreatedAt, &u.UpdatedAt) == nil {
			res = append(res, &u)
		}
	}

	return res, nil

}

func FindUsers() *UserFindQuery {
	return &UserFindQuery{}
}

func NewUserID() int64 {
	var id int64
	_ = DB.QueryRow("SELECT max(id) FROM user").Scan(&id)
	return id + 1
}

func (u *User) Insert() error {

	id := NewUserID()
	_, err := DB.Insert("user").Fields(map[string]interface{}{
		"id":         id,
		"first_name": u.FirstName,
		"last_name":  u.LastName,
		"avatar":     u.Avatar,
		"email":      u.Email,
		"password":   u.Password,
		"phone":      u.Phone,
		"address":    u.Address,
		"verified":   u.Verified,
		"deleted":    0,
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}).Execute()

	if err != nil {
		return err
	}

	u.ID = id

	DB.SetUserCache(id, u)

	return nil
}

func (u *User) Update() error {

	if u.ID == 0 {
		return errors.New("invalid user id")
	}
	_, err := DB.Update("user").Fields(map[string]interface{}{
		"first_name": u.FirstName,
		"last_name":  u.LastName,
		"avatar":     u.Avatar,
		"email":      u.Email,
		"password":   u.Password,
		"phone":      u.Phone,
		"address":    u.Address,
		"verified":   u.Verified,
		"deleted":    u.Deleted,
		"created_at": u.CreatedAt,
		"updated_at": u.UpdatedAt,
	}).Condition("id", u.ID, "=").Execute()

	if err != nil {
		return err
	}

	DB.SetUserCache(u.ID, u)

	return nil
}

func (u *User) Delete() error {

	u.Deleted = true
	_, err := DB.Exec("UPDATE `user` SET deleted = 1 WHERE id = ?", u.ID)
	if err == nil {
		DB.SetUserCache(u.ID, u)
	}

	return err
}

func (u *User) Roles() []*Role {

	rows, err := DB.Select("user_role", "ur").Join("role", "r", "ur.role_id = r.id").Fields("r", []string{
		"id",
		"name",
	}).Condition("ur.user_id", u.ID, "=").FetchAll()

	if err != nil {
		return nil
	}

	defer rows.Close()

	var res [] *Role

	for rows.Next() {
		var r Role
		if rows.Scan(&r.ID, &r.Name) == nil {
			res = append(res, &r)
		}
	}

	return res
}

func UserByEmail(email string) (*User, error) {

	if len(email) == 0 {
		return nil, errors.New("email is required")
	}

	var u User

	err := DB.Select("user", "u").Fields("u", []string{
		"id",
		"first_name",
		"last_name",
		"avatar",
		"email",
		"password",
		"phone",
		"address",
		"is_online",
		"verified",
		"deleted",
		"created_at",
		"updated_at",
	}).Condition("u.email", email, "=").FetchOne().Scan(&u.ID, &u.FirstName, &u.LastName, &u.Avatar, &u.Email, &u.Password, &u.Phone, &u.Address, &u.IsOnline, &u.Verified, &u.Deleted, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &u, nil
}

func Login(email, password string) (*User, error) {

	user, err := UserByEmail(email)
	if err != nil {

		return nil, errors.New("user does not exist")
	}

	if !ComparePassword(password, user.Password) {
		return nil, errors.New("wrong password")
	}

	return user, nil
}

func Logout(token string) error {

	_, err := DB.Insert("jwt_blacklist").Fields(map[string]interface{}{
		"token": token,
	}).Execute()

	return err
}

func TokenIsBlackList(token string) bool {
	var count int
	if DB.QueryRow("SELECT EXISTS(SELECT id FROM jwt_blacklist WHERE token = ?)", token).Scan(&count) != nil {
		return true
	}

	if count > 0 {
		return true
	}
	return false

}
