package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrVersionConflict   = errors.New("resource version conflict")
	ErrDataConflict      = errors.New("resource data conflict")
	ErrDuplicateEmail    = errors.New("dublicate email")
	ErrDuplicateUsername = errors.New("dublicate username")
)

var (
	QueryTimeOutDuration = time.Second * 5
)

type Storage struct {
	Posts interface {
		GetByID(context.Context, int64) (*Post, error)
		Create(context.Context, *Post) error
		Update(context.Context, *Post) error
		Delete(context.Context, int64) error
		GetUserFeed(context.Context, int64, PaginatedFeedQuery) ([]*PostWithMetaData, error)
	}
	Users interface {
		Create(context.Context, *sql.Tx, *User) error
		GetUserByID(context.Context, int64) (*User, error)
		GetByEmail(context.Context, string) (*User, error)
		CreateAndInvite(context.Context, *User, string, time.Duration) error
		Activate(context.Context, string) error
		Delete(context.Context, int64) error
	}
	Comments interface {
		GetByPostID(context.Context, int64) ([]*Comment, error)
		Create(context.Context, *Comment) error
	}
	Followers interface {
		Follow(context.Context, int64, int64) error
		Unfollow(context.Context, int64, int64) error
	}
	Roles interface {
		GetByName(context.Context, string) (*Role, error)
	}
}

func NewStore(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostStore{db},
		Users:     &UserStore{db},
		Comments:  &CommentStore{db},
		Followers: &FollowerStore{db},
		Roles:     &RoleStore{db},
	}
}

func withTX(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {

	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = fn(tx)
	if err != nil {
		return err
	}

	return tx.Commit()
}
