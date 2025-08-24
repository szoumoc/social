package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	"github.com/szoumoc/social/internal/store"
)

var usernames = []string{
	"CodeWarrior01", "BackendNinja", "SystemGuru", "APIMaster",
	"SecureCoder", "BugHunterX", "CloudArchitect", "DataForge",
	"GoLangSamurai", "PostgresWizard", "SQLKnight", "PythonPhantom",
	"CryptoGuard", "AlgoSensei", "ThreadMaster", "CacheCrusader",
	"PacketWizard", "KernelKid", "RustyCoder", "StackGuardian",
	"LambdaLord", "BitwiseBoss", "ZeroDayHero", "PortScannerX",
	"MemoryMapper", "HeapHacker", "SocketSurgeon", "RaceCondRider",
	"MutexMonk", "PointerPilot", "ShellShockX", "SecureSession",
	"TLSKnight", "HashHunter", "BackendBeast", "AuthOverlord",
	"JWTJuggler", "DockerDemon", "K8sOverseer", "QueueHandler",
	"ThreadTracer", "APIProtector", "LogAnalyzerX", "IndexInvoker",
	"KeyGenKaiser", "CloudCipher", "VaultKeeper", "ScriptSamurai",
	"BackendWizard", "BitBreaker",
}
var titles = []string{
	"Mastering Go: Concurrency", "Understanding Kubernetes", "Docker Deep Dive",
	"Building REST APIs", "Microservices Architecture", "Cloud Native Go",
	"Advanced SQL Queries", "Data Structures and Algorithms", "Go for DevOps",
	"Secure Coding Practices", "Performance Optimization", "GraphQL in Action",
	"WebAssembly with Go", "Go and Machine Learning", "Blockchain Basics",
	"Real-time Data Processing", "Event-Driven Architecture", "Go for Game Development",
	"Building Scalable Applications", "Go in the Cloud",
}

var contents = []string{
	"Go routines and concurrency patterns", "Kubernetes for beginners", "Docker best practices",
	"REST API design principles", "Microservices with Go", "Cloud native application development",
	"SQL query optimization techniques", "Data structures in Go", "DevOps with Go",
	"Secure coding in Go", "Performance tuning Go applications", "GraphQL API development",
	"WebAssembly use cases", "Machine learning with Go", "Blockchain implementation details",
	"Real-time data processing with Go", "Event-driven architecture patterns", "Game development with Go",
	"Building scalable web applications", "Deploying Go applications in the cloud",
}

var comments = []string{
	"Great post!", "Very informative, thanks!", "I learned something new today.",
	"Can't wait to try this out!", "This is exactly what I needed.", "Well explained!",
	"Looking forward to more content like this.", "This helped me a lot, thank you!",
	"Interesting perspective!", "I appreciate the detailed examples.",
}

var tags = []string{
	"Go", "Kubernetes", "Docker", "REST", "Microservices", "Cloud",
	"SQL", "Data Structures", "DevOps", "Security", "Performance",
	"GraphQL", "WebAssembly", "Machine Learning", "Blockchain",
	"Real-time", "Event-driven", "Game Development", "Scalability",
}

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()
	users := generateUsers(100)
	tx, _ := db.BeginTx(ctx, nil)

	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			_ = tx.Rollback()
			log.Printf("failed to create user: %v", err)
			return
		}
	}
	tx.Commit()
	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Printf("failed to create post: %v", err)
			return
		}
	}
	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Printf("failed to create comment: %v", err)
			return
		}
	}
	log.Printf("Seeding completed successfully")

}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)],
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
			Role: store.Role{
				Name: "user",
			},
		}
	}

	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)
	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]

		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}

	return posts
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	cms := make([]*store.Comment, num)
	for i := 0; i < num; i++ {
		cms[i] = &store.Comment{
			PostID:  posts[rand.Intn(len(posts))].ID,
			UserID:  users[rand.Intn(len(users))].ID,
			Content: comments[rand.Intn(len(comments))],
		}
	}
	return cms
}
