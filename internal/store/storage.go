package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("record not found")
	ErrConflict = errors.New("resource already exits")
	QueryTimeoutDuration = time.Second *5
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetByID(context.Context, int64) (*Post, error)
		Update(context.Context, *Post)(error)
		Delete(context.Context, int64) (error)
		GetUserFeed(context.Context, int64, PaginatedFeedQuery) ([]PostWithMetadata, error)
	}

	Users interface {
		Create(context.Context, *User) error
		GetByID(context.Context, int64) (*User, error)
	}

	Followers interface{
		Follow(context.Context, int64, int64 ) error
		Unfollow(context.Context, int64, int64) error
	}

	Comments interface{
		GetByPostID(context.Context, int64) ([]Comment, error)
		Create(context.Context, *Comment) error
	}
}


func NewStorage(db *sql.DB) Storage{
	return Storage{
		Posts: &PostStore{db},
		Users: &UserStore{db},
		Comments: &CommentStore{db},
		Followers: &FollowesStore{db},
	}
}