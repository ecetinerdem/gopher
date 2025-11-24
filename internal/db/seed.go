package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	"github.com/ecetinerdem/gopherSocial/internal/store"
)

func Seed(store store.Storage, db *sql.DB) {

	ctx := context.Background()

	users := generateUsers(100)
	tx, _ := db.BeginTx(ctx, nil)
	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			_ = tx.Rollback()
			log.Println("Error creating user: ", err)
			return
		}
	}

	tx.Commit()

	posts := generatePosts(200, users)

	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating post: ", err)
			return
		}

	}

	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating comment: ", err)
			return
		}

	}

	log.Println("Seeding completed")
}

func generateUsers(amount int) []*store.User {
	users := make([]*store.User, amount)

	//Todo get role by id
	for i := 0; i < amount; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
			RoleID:   1,
		}
	}
	return users
}

func generatePosts(amount int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, amount)

	for i := 0; i < amount; i++ {
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

func generateComments(amount int, users []*store.User, posts []*store.Post) []*store.Comment {
	cms := make([]*store.Comment, amount)

	for i := 0; i < amount; i++ {
		user := users[rand.Intn(len(users))]
		post := posts[rand.Intn(len(posts))]

		cms[i] = &store.Comment{
			UserID:  user.ID,
			PostID:  post.ID,
			Content: comments[rand.Intn(len(comments))],
		}
	}
	return cms
}
