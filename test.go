package dexp

import (
	"log"
	"strconv"
	"time"
)

func Test() {
	//CreateUsers()
	//UpdateUser()
	//CreateProjects()
	//Select()
	//entityLoad()

	go func() {
		for i := 0; i < 10000; i++ {
			project := &Project{
				AuthorID:  1,
				Title:     "Project " + time.Now().String(),
				Body:      "Description " + time.Now().String(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			project.Insert()
		}
	}()
}

func entityLoad() {

	q := FindProjects()
	q.Where("id > ?", 1)
	q.OrderBy("title", "ASC")
	q.Range(10, 0)

	res, err := q.Execute()

	log.Println("pro", res[0].Title, err)
}

func CreateUsers() {

	for i := 1; i < 10; i++ {
		_, _ = DB.Insert("user").Fields(map[string]interface{}{
			"id":         i,
			"first_name": "Toan " + strconv.Itoa(i),
			"last_name":  "Nguyen",
			"email":      "toan" + strconv.Itoa(i) + "@gmail.com",
			"password":   "admin",
			"phone":      "0932504043",
			"address":    "50 Hoa Minh 8, Lien Chieu Da nang",
		}).Execute()

		//log.Println("error", err, result)

	}

}

func CreateProjects() {

	for i := 1; i < 10; i++ {

		project := &Project{
			AuthorID:  1,
			Title:     "Project " + strconv.Itoa(i),
			Body:      "Description " + strconv.Itoa(i),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := project.Insert()
		log.Println("project insert", err, project)
	}

}

func UpdateUser() {

	result, err := DB.Update("user").Fields(map[string]interface{}{
		"first_name": "Toan has updated",
	}).Condition("id", 1, "=").Execute()

	log.Println("Update error", err, result)
}

func Select() {

	var id int64
	var title string

	err := DB.Select("entity", "e").
		Fields("e", []string{"id"}).
		Join("entity_project", "p", "p.entity_id = e.id").
		Fields("p", []string{"title"}).
		Condition("e.id", 2, "=").
		FetchOne().
		Scan(&id, &title)

	println("select ", err, id, title)
}
