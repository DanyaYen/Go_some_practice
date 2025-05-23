package main

import (
    "encoding/json"
    "github.com/graphql-go/graphql"
    "io/ioutil"
    "log"
    "net/http"
    "strconv" 
)

type Todo struct {
    ID   string `json:"id"`
    Text string `json:"text"`
    Done bool   `json:"done"`
}

var (
    todos = []Todo{
        {ID: "1", Text: "Hello World!", Done: true},
        {ID: "2", Text: "Hello Me!", Done: true},
        {ID: "3", Text: "Hello You!", Done: false},
    }
    nextID = 4
)

var todoType = graphql.NewObject(
    graphql.ObjectConfig{
        Name: "Todo",
        Fields: graphql.Fields{
            "id": &graphql.Field{
                Type: graphql.String,
            },
            "text": &graphql.Field{
                Type: graphql.String,
            },
            "done": &graphql.Field{
                Type: graphql.Boolean,
            },
        },
    },
)

var queryType = graphql.NewObject(
    graphql.ObjectConfig{
        Name: "Query",
        Fields: graphql.Fields{
            "todos": &graphql.Field{
                Type: graphql.NewList(todoType),
                Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                    return todos, nil
                },
            },
            "todo": &graphql.Field{
                Type: todoType,
                Args: graphql.FieldConfigArgument{
                    "id": &graphql.ArgumentConfig{
                        Type: graphql.NewNonNull(graphql.String),
                    },
                },
                Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                    id, ok := p.Args["id"].(string)
                    if ok {
                        for _, todo := range todos {
                            if todo.ID == id {
                                return todo, nil
                            }
                        }
                    }
                    return nil, nil
                },
            },
        },
    },
)

var mutationType = graphql.NewObject(
    graphql.ObjectConfig{
        Name: "Mutation",
        Fields: graphql.Fields{
            "addTodo": &graphql.Field{
                Type: todoType,
                Args: graphql.FieldConfigArgument{
                    "text": &graphql.ArgumentConfig{
                        Type: graphql.NewNonNull(graphql.String),
                    },
                },
                Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                    text, ok := p.Args["text"].(string)
                    if ok {
                        newTodo := Todo{
                            ID:   strconv.Itoa(nextID), 
                            Text: text,
                            Done: false,
                        }
                        todos = append(todos, newTodo)
                        nextID++
                        return newTodo, nil
                    }
                    return nil, nil
                },
            },
        },
    },
)

var schema graphql.Schema 

func init() { 
    var err error
    schema, err = graphql.NewSchema(
        graphql.SchemaConfig{
            Query:    queryType,
            Mutation: mutationType,
        },
    )
    if err != nil {
        log.Fatalf("failed to create schema, error: %v", err)
    }
}

func handler(w http.ResponseWriter, r *http.Request) {
    var queryString string

    if r.Method == http.MethodGet {
        queryString = r.URL.Query().Get("query")
    } else if r.Method == http.MethodPost {
        body, err := ioutil.ReadAll(r.Body)
        if err != nil {
            log.Printf("could not read request body: %v", err)
            http.Error(w, "Could not read request body", http.StatusBadRequest)
            return
        }
        defer r.Body.Close() // Good practice to close request body

        var requestBody struct {
            Query string `json:"query"`
        }
        if err := json.Unmarshal(body, &requestBody); err != nil {
            log.Printf("could not unmarshal request body: %v. Body was: %s", err, string(body)) // Log body on unmarshal error
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }
        queryString = requestBody.Query
    }

    log.Printf("Received GraphQL query string: [%s]", queryString) // Log the query string

    if queryString == "" {
        log.Println("Query string is empty, which will lead to 'Must provide an operation' error.")
    }

    result := graphql.Do(graphql.Params{
        Schema:        schema,
        RequestString: queryString,
    })

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(result); err != nil {
        log.Printf("could not write response: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func main() {
    http.HandleFunc("/graphql", handler)

    log.Println("Server is running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}