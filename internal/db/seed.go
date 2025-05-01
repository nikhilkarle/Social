package db

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/nikhilkarle/social/internal/store"
)

var usernames = []string{
	"alice", "bob", "charlie", "dave", "eve", "frank", "grace", "heidi",
	"ivan", "judy", "karl", "laura", "mallory", "nina", "oscar", "peggy",
	"quinn", "rachel", "steve", "trent", "ursula", "victor", "wendy", "xander",
	"yvonne", "zack", "amber", "brian", "carol", "doug", "eric", "fiona",
	"george", "hannah", "ian", "jessica", "kevin", "lisa", "mike", "natalie",
	"oliver", "peter", "queen", "ron", "susan", "tim", "uma", "vicky",
	"walter", "xenia", "yasmin", "zoe",
}

var titles = []string{
	"The Power of Habit", "Embracing Minimalism", "Healthy Eating Tips",
	"Travel on a Budget", "Mindfulness Meditation", "Boost Your Productivity",
	"Home Office Setup", "Digital Detox", "Gardening Basics",
	"DIY Home Projects", "Yoga for Beginners", "Sustainable Living",
	"Mastering Time Management", "Exploring Nature", "Simple Cooking Recipes",
	"Fitness at Home", "Personal Finance Tips", "Creative Writing",
	"Mental Health Awareness", "Learning New Skills",
}

var contents = []string{
	"Building new habits starts with small, consistent actions.",
	"Minimalism can free your mind and simplify your environment.",
	"Healthy eating begins with mindful grocery shopping.",
	"Budget travel opens up adventure without breaking the bank.",
	"Meditation can bring clarity and reduce daily stress.",
	"Productivity improves with fewer distractions and better planning.",
	"A well-designed home office boosts focus and comfort.",
	"Digital detoxing once a week helps reset your attention span.",
	"Gardening is therapeutic and helps you connect with nature.",
	"DIY projects make your home more personalized and creative.",
	"Yoga enhances flexibility, balance, and mental calmness.",
	"Sustainable living starts with reusable products and conscious choices.",
	"Time management is about priorities, not just schedules.",
	"Spending time in nature improves mood and mental health.",
	"Simple home recipes can be both healthy and time-saving.",
	"Staying active at home is easier with short workout routines.",
	"Financial wellness begins with budgeting and tracking expenses.",
	"Creative writing can be both therapeutic and intellectually rewarding.",
	"Raising awareness about mental health helps reduce stigma.",
	"Trying new skills keeps the brain sharp and builds confidence.",
}

var tags = []string{
	"wellness",
	"productivity",
	"mindfulness",
	"fitness",
	"finance",
	"mental-health",
	"self-improvement",
	"lifestyle",
	"minimalism",
	"meditation",
	"yoga",
	"home-office",
	"gardening",
	"travel",
	"creativity",
	"nutrition",
	"eco-friendly",
	"time-management",
	"diy",
	"learning",
}

var comments = []string{
	"Love this post — super helpful!",
	"Ive been trying this and it really works!",
	"Thanks for sharing, very insightful.",
	"Bookmarking this to come back later.",
	"Great tips! I going to apply them today.",
	"This is exactly what I needed to read.",
	"Can you do a follow-up on this topic?",
	"I hadnt thought about it that way before.",
	"Such a well-written piece, thank you!",
	"Totally agree — consistency is key.",
	"This changed my perspective completely.",
	"Ive shared this with a few friends already.",
	"Looking forward to more posts like this.",
	"Simple, clear, and effective advice.",
	"Wow, I never realized how important this is.",
	"Do you have any book recommendations on this?",
	"Would love to hear more examples.",
	"This was so motivating to read.",
	"Please post more on this subject!",
	"Im inspired to make a change now.",
}

func Seed(store store.Storage){
	ctx := context.Background()

	users := generateUsers(100)
	for _,user := range users{
		if err := store.Users.Create(ctx, user); err != nil{
			log.Println("Error creating user:", err)
			return 
		}
	}

	posts := generatePosts(200, users)
	for _, post := range posts{
		if err := store.Posts.Create(ctx, post); err != nil{
			log.Println("Error creating post", err)
			return 
		}
	}
	
	comments := generateComments(1000, users, posts)
	for _, comment := range comments{
		if err := store.Comments.Create(ctx, comment); err != nil{
			log.Println("Error creating comment", err)
			return 
		}
	}

	log.Println("Seeding complete.")
}

func generateUsers(num int)[]*store.User{
	users := make([]*store.User, num)

	for i:=0 ; i < num; i++{
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d",i),
			Email: usernames[i%len(usernames)] + fmt.Sprintf("%d",i) + "@example.com",
			Password: "123123",
		}
	}

	return users
}

func generatePosts(num int, users []*store.User) []*store.Post{
	posts := make([]*store.Post, num)
	for i := 0; i < num; i++{
		user := users[rand.Intn(len(users))]
		
		posts[i] = &store.Post{
			UserID: user.ID,
			Title: titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}

	return posts

}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment{
	newComments := make([]*store.Comment, num)
	for i := 0; i<num; i++{
		newComments[i] = &store.Comment{
			PostID: posts[rand.Intn(len(posts))].ID,
			UserID: users[rand.Intn(len(users))].ID,
			Content: comments[rand.Intn(len(comments))],
		}
	}

	return newComments

}