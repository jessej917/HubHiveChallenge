package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go_project/src/server/utils"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/rs/cors"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/getUsers", GetUsers)
	r.HandleFunc("/getPosts", GetPosts)
	r.HandleFunc("/createUser", CreateUser).Methods("POST", "OPTIONS")
	r.HandleFunc("/createPost", CreatePost).Methods("POST", "OPTIONS")
	r.HandleFunc("/login", Login).Methods("POST", "OPTIONS")

	// Solves Cross Origin Access Issue
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"}, // Allow frontend origin
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true, // Allow cookies or authentication headers
	})
	handler := c.Handler(r)

	srv := &http.Server{
		Handler: handler,
		Addr:    ":" + os.Getenv("PORT"),
	}

	log.Fatal(srv.ListenAndServe())
}

func GetPosts(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Connecting to Database from GetPosts...")
	ctx, driver := ConnectToDatabase()
	defer driver.Close(ctx)

	// map parameter can be nil
	// "MATCH (p:User {username: $username}) RETURN p.username AS username",
	// Get the name of all 42 year-olds
	result, _ := neo4j.ExecuteQuery(ctx, driver,
		"match (u:User) -[c:CREATED]-> (p:Post) return p.title as title, p.body as body, p.image as image, date(c.date) AS date, u.username as username",
		nil, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"))

	posts := []PostRequest{}
	//listTest := []any{}

	// Loop through results and do something with them
	for _, record := range result.Records {
		fmt.Println(record.AsMap())
		fmt.Println(record.AsMap()["username"])
		fmt.Println(record.AsMap()["title"])
		fmt.Println(record.AsMap()["body"])
		fmt.Println(record.AsMap()["image"])
		fmt.Println(record.AsMap()["date"])
		var post = PostRequest{
			Username: fmt.Sprint(record.AsMap()["username"]),
			Title:    fmt.Sprint(record.AsMap()["title"]),
			Body:     fmt.Sprint(record.AsMap()["body"]),
			Image:    fmt.Sprint(record.AsMap()["image"]),
			Date:     fmt.Sprint(record.AsMap()["date"]),
		}
		posts = append(posts, post)
		//listTest = append(listTest, record.AsMap()["username"])
	}

	// For convertion
	fmt.Printf("Thing: %s\n", posts[0].Date)
	fmt.Printf("Type of things!: %T\n", posts[0].Date)

	// Summary information
	fmt.Printf("The query `%v` returned %v records in %+v.\n",
		result.Summary.Query().Text(), len(result.Records),
		result.Summary.ResultAvailableAfter())

	fmt.Println(posts)

	jsonBytes, err := utils.StructToJSON(posts)
	if err != nil {
		fmt.Print(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
	return
}

func CreatePost(w http.ResponseWriter, r *http.Request) {

	var postReq PostRequest

	// Decode the JSON body into the LoginRequest struct
	err := json.NewDecoder(r.Body).Decode(&postReq)
	fmt.Println(postReq)
	fmt.Println(err)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var result LoginResult

	if len(postReq.Username) >= 4 && len(postReq.Title) >= 4 && len(postReq.Body) >= 4 {

		//Create User in database
		fmt.Println("Connecting to Database from CreatePost...")
		ctx, driver := ConnectToDatabase()
		defer driver.Close(ctx)

		// Create a post based on the given user and post information
		databaseResult, err := neo4j.ExecuteQuery(ctx, driver,
			`MATCH (p:User)
		WHERE p.username = $username
		CREATE (p) -[:CREATED {date: datetime()}]-> (post:Post {title: $title, body: $body, image: $image})`,
			map[string]any{
				"username": postReq.Username,
				"title":    postReq.Title,
				"body":     postReq.Body,
				"image":    postReq.Image,
			}, neo4j.EagerResultTransformer,
			neo4j.ExecuteQueryWithDatabase("neo4j"))

		if err != nil {
			panic(err)
		}

		summary := databaseResult.Summary
		fmt.Printf("Created %v nodes in %+v.\n",
			summary.Counters().NodesCreated(),
			summary.ResultAvailableAfter())

		// Use the err
		if err == nil {
			result = LoginResult{
				Message: "Register successful!",
				Result:  true,
			}
			fmt.Println("Success!!")
		} else {
			result = LoginResult{
				Message: "Register failed",
				Result:  false,
			}
			fmt.Println("Fail!!")
		}

	} else {
		if len(postReq.Username) < 4 {
			result = LoginResult{
				Message: "Username must be at least 4 characters!",
				Result:  false,
			}
		} else if len(postReq.Title) < 4 {
			result = LoginResult{
				Message: "Title must be at least 4 characters!",
				Result:  false,
			}
		} else if len(postReq.Body) < 4 {
			result = LoginResult{
				Message: "Body must be at least 4 characters!",
				Result:  false,
			}
		}
	}

	fmt.Println(result)
	fmt.Println(result.Result)

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Print JSON to server logs
	log.Println("JSON response:", string(jsonBytes))

	// Write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Connecting to Database from GetUsers...")
	ctx, driver := ConnectToDatabase()
	defer driver.Close(ctx)

	// map parameter can be nil
	// "MATCH (p:User {username: $username}) RETURN p.username AS username",
	// Get the name of all 42 year-olds
	result, _ := neo4j.ExecuteQuery(ctx, driver,
		"MATCH (p:User) RETURN p.username AS username",
		nil, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"))

	test := []map[string]any{}
	listTest := []any{}

	// Loop through results and do something with them
	for _, record := range result.Records {
		fmt.Println(record.AsMap())
		fmt.Println(record.AsMap()["username"])
		test = append(test, record.AsMap())
		listTest = append(listTest, record.AsMap()["username"])
	}

	// For convertion
	fmt.Printf("Thing: %s\n", test[0]["username"])
	fmt.Printf("Type of things!: %T\n", test[0]["username"])
	m, ok := test[0]["username"].(dbtype.Node)
	if ok {
		fmt.Println("Conversion successful:", m.Props["username"])
	} else {
		fmt.Println("Conversion failed")
	}

	// Summary information
	fmt.Printf("The query `%v` returned %v records in %+v.\n",
		result.Summary.Query().Text(), len(result.Records),
		result.Summary.ResultAvailableAfter())

	fmt.Println(listTest)
	fmt.Println("This is a really important test")

	jsonBytes, err := utils.StructToJSON(listTest)
	if err != nil {
		fmt.Print(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
	return
}

func CreateUser(w http.ResponseWriter, r *http.Request) {

	var loginReq LoginRequest

	// Decode the JSON body into the LoginRequest struct
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	fmt.Println(loginReq)
	fmt.Println(err)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var loginResult LoginResult
	if len(loginReq.Username) >= 4 && len(loginReq.Password) >= 4 {

		fmt.Println("Connecting to Database from CreateUser...")
		ctx, driver := ConnectToDatabase()
		defer driver.Close(ctx)

		// map parameter can be nil
		// "MATCH (p:User {username: $username}) RETURN p.username AS username",
		// Get the name of all 42 year-olds
		result, err := neo4j.ExecuteQuery(ctx, driver,
			"Create (p:User {username: $username, password: $password}) return p",
			map[string]any{
				"username": loginReq.Username,
				"password": loginReq.Password,
			}, neo4j.EagerResultTransformer,
			neo4j.ExecuteQueryWithDatabase("neo4j"))
		fmt.Println("test!")
		if err != nil {
			fmt.Println("Error: ")
			//panic(err)
		}

		fmt.Println(result)

		// Use the err
		//var loginResult LoginResult
		if err == nil {
			loginResult = LoginResult{
				Message: "Registered successful!",
				Result:  true,
			}
			fmt.Println("Success!!")
		} else {
			loginResult = LoginResult{
				Message: "Username already exists!",
				Result:  false,
			}
			fmt.Println("Fail!!")
		}
	} else {
		if len(loginReq.Username) < 4 {
			loginResult = LoginResult{
				Message: "Username must be at least 4 characters!",
				Result:  false,
			}
		} else if len(loginReq.Password) < 4 {
			loginResult = LoginResult{
				Message: "Password must be at least 4 characters!",
				Result:  false,
			}
		}
	}

	fmt.Println(loginResult)
	fmt.Println(loginResult.Result)

	jsonBytes, err := json.Marshal(loginResult)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Print JSON to server logs
	log.Println("JSON response:", string(jsonBytes))

	// Write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

// TODO: Add security (salt, pepper, hash)
func Login(w http.ResponseWriter, r *http.Request) {
	var loginReq LoginRequest

	// Decode the JSON body into the LoginRequest struct
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	fmt.Println(loginReq)
	fmt.Println(err)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fmt.Println("Connecting to Database from Login...")
	ctx, driver := ConnectToDatabase()
	defer driver.Close(ctx)

	// map parameter can be nil
	// "MATCH (p:User {username: $username}) RETURN p.username AS username",
	// Get the name of all 42 year-olds
	result, _ := neo4j.ExecuteQuery(ctx, driver,
		"MATCH (p:User {username: $username, password: $password}) RETURN count(p) as count",
		map[string]any{
			"username": loginReq.Username,
			"password": loginReq.Password,
		}, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"))

	// Use the data (e.g., authentication, logging, etc.)
	var loginRes LoginResult
	if result.Records[0].AsMap()["count"] == int64(1) {
		loginRes = LoginResult{
			Message: "Login successful!",
			Result:  true,
		}
		fmt.Println("Success!!")
	} else {
		loginRes = LoginResult{
			Message: "Invalid username or password",
			Result:  false,
		}
		fmt.Println("Fail!!")
	}

	fmt.Println(loginRes)
	fmt.Println(loginRes.Result)

	jsonBytes, err := json.Marshal(loginRes)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Print JSON to server logs
	log.Println("JSON response:", string(jsonBytes))

	// Write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func ConnectToDatabase() (context.Context, neo4j.DriverWithContext) {
	ctx := context.Background()
	dbUri := "bolt://localhost:7687"
	dbUser := "HubhiveDB"
	dbPassword := "asdf1234"
	driver, err := neo4j.NewDriverWithContext(
		dbUri,
		neo4j.BasicAuth(dbUser, dbPassword, ""))
	if err != nil {
		panic(err)
	}
	//defer driver.Close(ctx)

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connection established.")
	fmt.Printf("Type of driver: %T\n", driver)

	return ctx, driver
}

type PostRequest struct {
	Username string `json:"username"`
	Title    string `json:"title"`
	Body     string `json:"body"`
	Image    string `json:"image"`
	Date     string `json:"date"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResult struct {
	Message string `json:"message"`
	Result  bool   `json:"result"`
}
