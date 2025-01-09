package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go_project/src/server/utils"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/rs/cors"
)

func main() {

	// Ensure the "uploads" directory exists
	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		err := os.Mkdir("./uploads", os.ModePerm)
		if err != nil {
			log.Fatal("Error creating uploads directory:", err)
		}
	}

	r := mux.NewRouter()

	// Serve static files from the "uploads" directory
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	r.HandleFunc("/getUsers", GetUsers)
	r.HandleFunc("/getPosts", GetPosts)
	r.HandleFunc("/getFriends", GetFriends)
	r.HandleFunc("/createUser", CreateUser).Methods("POST", "OPTIONS")
	r.HandleFunc("/createPost", CreatePost).Methods("POST", "OPTIONS")
	r.HandleFunc("/addFriend", AddFriend).Methods("POST", "OPTIONS")
	r.HandleFunc("/login", Login).Methods("POST", "OPTIONS")
	r.HandleFunc("/uploadImage", UploadImage).Methods("POST", "OPTIONS")

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

	// Run Query to get the posts
	result, _ := neo4j.ExecuteQuery(ctx, driver,
		"match (u:User) -[c:CREATED]-> (p:Post) return p.title as title, p.body as body, p.image as image, date(c.date) AS date, u.username as username",
		nil, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"))

	posts := []PostRequest{}

	// Loop through results and do something with them
	for _, record := range result.Records {
		var post = PostRequest{
			Username: fmt.Sprint(record.AsMap()["username"]),
			Title:    fmt.Sprint(record.AsMap()["title"]),
			Body:     fmt.Sprint(record.AsMap()["body"]),
			Image:    fmt.Sprint(record.AsMap()["image"]),
			Date:     fmt.Sprint(record.AsMap()["date"]),
		}
		posts = append(posts, post)
	}

	// Summary information
	fmt.Printf("The query `%v` returned %v records in %+v.\n",
		result.Summary.Query().Text(), len(result.Records),
		result.Summary.ResultAvailableAfter())

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
				Message: "Created successful!",
				Result:  true,
			}
			fmt.Println("Success!!")
		} else {
			result = LoginResult{
				Message: "Creation failed",
				Result:  false,
			}
			fmt.Println("Fail!!")
		}

	} else {
		if len(postReq.Username) < 4 {
			fmt.Println(postReq.Username)
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

func UploadImage(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Uploading Image...")

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form with a maximum memory size (e.g., 10MB)
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "Unable to parse form data", http.StatusBadRequest)
		log.Println("Error parsing form data:", err)
		return
	}

	// Retrieve the file from the form data
	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Unable to retrieve file from form data", http.StatusBadRequest)
		log.Println("Error retrieving file:", err)
		return
	}
	defer file.Close()

	// Create a new file on the server to save the uploaded file
	dst, err := os.Create("./uploads/" + fileHeader.Filename)
	if err != nil {
		http.Error(w, "Unable to save the file", http.StatusInternalServerError)
		log.Println("Error saving file:", err)
		return
	}
	defer dst.Close()

	// Copy the file's content to the new file
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		log.Println("Error copying file content:", err)
		return
	}

	// Respond to the client
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", "uploads/"+fileHeader.Filename)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Connecting to Database from GetUsers...")
	ctx, driver := ConnectToDatabase()
	defer driver.Close(ctx)

	// Run get users query
	result, _ := neo4j.ExecuteQuery(ctx, driver,
		"MATCH (p:User) RETURN p.username AS username",
		nil, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"))

	listTest := []any{}

	// Loop through results and do something with them
	for _, record := range result.Records {
		listTest = append(listTest, record.AsMap()["username"])
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

		// Run the Create User query
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
		}

		fmt.Println(result)

		// Use the err
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

func GetFriends(w http.ResponseWriter, r *http.Request) {

	var username string = r.URL.Query().Get("username")

	//Create User in database
	fmt.Println("Connecting to Database from GetFriends...")
	ctx, driver := ConnectToDatabase()
	defer driver.Close(ctx)

	// Create a post based on the given user and post information
	databaseResult, err := neo4j.ExecuteQuery(ctx, driver,
		`match (u:User {username: $username}) -[:Friend]- (p:User) return p.username as username`,
		map[string]any{
			"username": username,
		}, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"))

	if err != nil {
		panic(err)
	}

	summary := databaseResult.Summary
	fmt.Printf("Created %v nodes in %+v.\n",
		summary.Counters().NodesCreated(),
		summary.ResultAvailableAfter())

	posts := []any{}

	// Loop through results and do something with them
	for _, record := range databaseResult.Records {
		posts = append(posts, record.AsMap()["username"])
	}

	jsonBytes, err := utils.StructToJSON(posts)
	if err != nil {
		fmt.Print(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
	return
}

func AddFriend(w http.ResponseWriter, r *http.Request) {

	var friendReq FriendRequest

	// Decode the JSON body into the LoginRequest struct
	err := json.NewDecoder(r.Body).Decode(&friendReq)
	fmt.Println(friendReq)
	fmt.Println(err)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var loginResult LoginResult

	fmt.Println("Connecting to Database from CreateUser...")
	ctx, driver := ConnectToDatabase()
	defer driver.Close(ctx)

	fmt.Println("Remove friend: ", friendReq.Remove)
	var query string = ""
	if friendReq.Remove {
		query = `match (u:User {username: $username}) -[f:Friend]- (p:User {username: $friend})
				delete f`
	} else {
		query = `match (u:User {username: $username}), (p:User {username: $friend})
				where not (u)--(p) and u <> p
				create (u) -[f:Friend]-> (p)`
	}

	result, err := neo4j.ExecuteQuery(ctx, driver,
		query,
		map[string]any{
			"username": friendReq.Username,
			"friend":   friendReq.Friend,
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
		if friendReq.Remove {
			loginResult = LoginResult{
				Message: "Removed Friend",
				Result:  true,
			}
		} else {
			loginResult = LoginResult{
				Message: "Added Friend",
				Result:  true,
			}
		}
		fmt.Println("Success!!")
	} else {
		loginResult = LoginResult{
			Message: "Error",
			Result:  false,
		}
		fmt.Println("Fail!!")
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

type FriendRequest struct {
	Username string `json:"username"`
	Friend   string `json:"friend"`
	Remove   bool   `json:"remove"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResult struct {
	Message string `json:"message"`
	Result  bool   `json:"result"`
}
