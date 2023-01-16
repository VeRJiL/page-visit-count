package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"html/template"
	"net/http"
	"os"
	"strconv"
)

type Page struct {
	Visit        int
	templatePath string
	redisClient  *redis.Client
	ctx          *context.Context
}

func (p *Page) handleHomePage(w http.ResponseWriter, r *http.Request) {
	var visitCount int
	
	homePageVisit, err := p.redisClient.Get(*p.ctx, "homePageVisit").Result()
	if err == redis.Nil {
		http.Error(w, "homePageVisit does not exist", http.StatusInternalServerError)
		visitCount = 0
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		visitCount, err = strconv.Atoi(homePageVisit)
		if err != nil {
			http.Error(w, "Page Visit is NOT NUMBER!", http.StatusInternalServerError)
			return
		}
	}
	
	visitCount++
	p.Visit = visitCount
	err = p.redisClient.Set(context.Background(), "homePageVisit", visitCount, 0).Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	
	tmpl, err := template.ParseFiles(p.templatePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if err := tmpl.Execute(w, p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func connectToRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis-server:6379",
		Password: "",
		DB:       0,
	})
	
	pong, err := client.Ping(client.Context()).Result()
	fmt.Println(pong, err)
	
	return client
}

func main() {
	fmt.Println("Server running at http://localhost:8081/")
	
	ctx := context.Background()
	redisClient := connectToRedis()
	
	homePage := Page{
		Visit:        0,
		templatePath: "/templates/home.html",
		ctx:          &ctx,
		redisClient:  redisClient,
	}
	
	http.HandleFunc("/", homePage.handleHomePage)
	
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error Lunching the Server")
	}
}
