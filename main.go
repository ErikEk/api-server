package main

import (
	//"bytes"
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"syscall"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/ssh/terminal"
)

type Person struct {
	ID        string   `json:"id,omitempty`
	Firstname string   `json:"firstname,omitempty`
	Lastname  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}
type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

var (
	people []Person
	db *sql.DB
	stmt *sql.Stmt
)


/* ENDCODE and send json
u := User{Id: "US123", Balance: 8}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(u)
	res, _ := http.Post("https://httpbin.org/post", "application/json; charset=utf-8", b)
	var body struct {
		// httpbin.org sends back key/value pairs, no map[string][]string
		Headers map[string]string `json:"headers"`
		Origin  string            `json:"origin"`
	}
	json.NewDecoder(res.Body).Decode(&body)
	fmt.Println(body)
*/

/* PRINT FULL DATASET
// query
        rows, err := db.Query("SELECT * FROM userinfo")
        checkErr(err)

        for rows.Next() {
            var uid int
            var username string
            var department string
            var created string
            err = rows.Scan(&uid, &username, &department, &created)
            checkErr(err)
            fmt.Println(uid)
            fmt.Println(username)
            fmt.Println(department)
            fmt.Println(created)
        }
*/

func GetPersonEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Person{})
}
func GetPeopleEndpoint(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(people)
	var person Person

	row := db.QueryRow("select id, firstname, lastname from gotest where id = ?;", 1)
	err := row.Scan(&person.ID, &person.Firstname, &person.Lastname)
	if err != nil {
		// If no results send null
		fmt.Println("No results")
	} else {
		fmt.Println(person.ID)
		fmt.Println(person.Firstname)
		fmt.Println(person.Lastname)
	}
}
func CreatePersonEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	person.ID = params["id"]

	
    res, err := stmt.Exec(6, "axel", "ek")
    fmt.Println("ddd2")
    me, ok := err.(*mysql.MySQLError)
	if !ok {
	    fmt.Print(err)
	    return
	}
	if me.Number == 1062 {
	    fmt.Println("It already exists in a database.")
	    return
	}

	/*
    if err.(*mysql.MySQLError).Number == 1062 { // 1062 duplicate entry
    	fmt.Println("Already exists")
    	return
    }
    if err != nil {
    	fmt.Println("Other eeror:",err)
    }
    */
    fmt.Println("ddd")
    id, err := res.LastInsertId()
    if err != nil {
    	fmt.Println(err)
    }

    fmt.Println(id)


	fmt.Println(person.Firstname)
	fmt.Println(person.Lastname)
	people = append(people, person)
	json.NewEncoder(w).Encode(people)
}
func DeletePersonEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range people {
		if item.ID == params["id"] {
			people = append(people[:index], people[index+1:]...)
			break
		}
		json.NewEncoder(w).Encode(people)
	}
}

func credentials() (string, string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, _ := reader.ReadString('\n')

	fmt.Print("Enter Password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))

	password := string(bytePassword)

	return strings.TrimSpace(username), strings.TrimSpace(password)
}

func main() {
	router := mux.NewRouter()
	people = append(people, Person{ID: "1", Firstname: "John", Lastname: "Doe", Address: &Address{City: "City X", State: "State X"}})
	people = append(people, Person{ID: "2", Firstname: "Koko", Lastname: "Doe", Address: &Address{City: "City Z", State: "State Y"}})
	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/people/{id}", GetPersonEndpoint).Methods("GET")
	router.HandleFunc("/people/{id}", CreatePersonEndpoint).Methods("POST")
	//router.HandleFunc("/people/{id}", DeletePersonEndpoint).Methods("DELETE")

	// Get user and password for MySQL database
	username, password := credentials()
	// Temp
	var err error

	db, err = sql.Open("mysql", "root:"+password+"@tcp(127.0.0.1:3306)/sys")
	if err != nil {
		fmt.Print(err.Error())
	}

	defer db.Close()

	// make sure connection is available
	err = db.Ping()
	if err != nil {
		fmt.Print(err.Error())
	}

	// Prepare a writer
	stmt, err = db.Prepare("INSERT gotest SET id=?,firstname=?,lastname=?")
    if err != nil {
    	fmt.Println(err)
    }

	
	var person Person

	row := db.QueryRow("select id, firstname, lastname from gotest where id = ?;", 1)
	err = row.Scan(&person.ID, &person.Firstname, &person.Lastname)
	if err != nil {
		// If no results send null
		fmt.Println("No results")
	} else {
		fmt.Println(person.ID)
		fmt.Println(person.Firstname)
		fmt.Println(person.Lastname)
	}

	log.Fatal(http.ListenAndServe(":8080", router))
}

/*
func main() {
	db, err := sql.Open("mysql", "root:passapp@tcp(127.0.0.1:3306)/sys")
	if err != nil {
		fmt.Print(err.Error())
	}
	defer db.Close()
	// make sure connection is available
	err = db.Ping()
	if err != nil {
		fmt.Print(err.Error())
	}
	type Person struct {
		Id         int
		First_Name string
		Last_Name  string
	}
	router := mux.NewRouter()

	// GET a person detail
	router.GET("/person/:id", func(c *gin.Context) {

		id := c.Param("id")
		row := db.QueryRow("select id, first_name, last_name from person where id = ?;", id)
		err = row.Scan(&person.Id, &person.First_Name, &person.Last_Name)
		if err != nil {
			// If no results send null
			result = gin.H{
				"result": nil,
				"count":  0,
			}
		} else {
			result = gin.H{
				"result": person,
				"count":  1,
			}
		}
		//c.JSON(http.StatusOK, result)
	})

	// GET all persons
	router.GET("/persons", func(c *gin.Context) {
		var (
			person  Person
			persons []Person
		)
		rows, err := db.Query("select id, first_name, last_name from person;")
		if err != nil {
			fmt.Print(err.Error())
		}
		for rows.Next() {
			err = rows.Scan(&person.Id, &person.First_Name, &person.Last_Name)
			persons = append(persons, person)
			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()
		c.JSON(http.StatusOK, gin.H{
			"result": persons,
			"count":  len(persons),
		})
	})

	// POST new person details
	router.POST("/person", func(c *gin.Context) {
		var buffer bytes.Buffer
		first_name := c.PostForm("first_name")
		last_name := c.PostForm("last_name")
		stmt, err := db.Prepare("insert into person (first_name, last_name) values(?,?);")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(first_name, last_name)

		if err != nil {
			fmt.Print(err.Error())
		}

		// Fastest way to append strings
		buffer.WriteString(first_name)
		buffer.WriteString(" ")
		buffer.WriteString(last_name)
		defer stmt.Close()
		name := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf(" %s successfully created", name),
		})
	})

	// PUT - update a person details
	router.PUT("/person", func(c *gin.Context) {
		var buffer bytes.Buffer
		id := c.Query("id")
		first_name := c.PostForm("first_name")
		last_name := c.PostForm("last_name")
		stmt, err := db.Prepare("update person set first_name= ?, last_name= ? where id= ?;")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(first_name, last_name, id)
		if err != nil {
			fmt.Print(err.Error())
		}

		// Fastest way to append strings
		buffer.WriteString(first_name)
		buffer.WriteString(" ")
		buffer.WriteString(last_name)
		defer stmt.Close()
		name := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully updated to %s", name),
		})
	})

	// Delete resources
	router.DELETE("/person", func(c *gin.Context) {
		id := c.Query("id")
		stmt, err := db.Prepare("delete from person where id= ?;")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(id)
		if err != nil {
			fmt.Print(err.Error())
		}
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully deleted user: %s", id),
		})
	})

	router.Run(":8080")
}*/
