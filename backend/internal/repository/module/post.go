package modulerepositories

import (
	"fmt"

	"github.com/dmawardi/Go-Template/internal/db"
	data "github.com/dmawardi/Go-Template/internal/helpers/data"
	"github.com/dmawardi/Go-Template/internal/models"
	"gorm.io/gorm"
)

type PostRepository interface {
	// Find a list of all users in the Database
	FindAll(limit int, offset int, order string, conditions []models.QueryConditionParameters) (*models.BasicPaginatedResponse[db.Post], error)
	FindById(int) (*db.Post, error)
	Create(post *db.Post) (*db.Post, error)
	Update(int, *db.Post) (*db.Post, error)
	Delete(int) error
	BulkDelete([]int) error
}

type postRepository struct {
	DB *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db}
}

// Creates a post in the database
func (r *postRepository) Create(post *db.Post) (*db.Post, error) {
	// Create above post in database
	result := r.DB.Create(&post)
	if result.Error != nil {
		return nil, fmt.Errorf("failed creating post: %w", result.Error)
	}

	return post, nil
}

// Find a list of posts in the database
func (r *postRepository) FindAll(limit int, offset int, order string, conditions []models.QueryConditionParameters) (*models.BasicPaginatedResponse[db.Post], error) {
	// Build meta data for posts
	metaData, err := data.BuildMetaData(r.DB, db.Post{}, limit, offset, order, conditions)
	if err != nil {
		fmt.Printf("Error building meta data: %s", err)
		return nil, err
	}

	// Query all posts based on the received parameters
	var posts []db.Post
	err = data.QueryAll(r.DB, &posts, limit, offset, order, conditions, []string{"User"})
	if err != nil {
		fmt.Printf("Error querying db for list of posts: %s", err)
		return nil, err
	}

	return &models.BasicPaginatedResponse[db.Post]{
		Data: &posts,
		Meta: *metaData,
	}, nil
}

// Find post in database by ID
func (r *postRepository) FindById(id int) (*db.Post, error) {
	// Create an empty ref object of type post
	post := db.Post{}
	// Check if post exists in db
	result := r.DB.First(&post, id)
	// If error detected
	if result.Error != nil {
		return nil, result.Error
	}
	// else
	return &post, nil
}

// Delete post in database
func (r *postRepository) Delete(id int) error {
	// Create an empty ref object of type post
	post := db.Post{}
	// Check if post exists in db
	result := r.DB.Delete(&post, id)

	// If error detected
	if result.Error != nil {
		fmt.Println("error in deleting post: ", result.Error)
		return result.Error
	}
	// else
	return nil
}

// Bulk delete posts in database
func (r *postRepository) BulkDelete(ids []int) error {
	// Delete users with specified IDs
	err := data.BulkDeleteByIds(db.Post{}, ids, r.DB)
	if err != nil {
		fmt.Println("error in deleting posts: ", err)
		return err
	}
	// else
	return nil
}

// Updates post in database
func (r *postRepository) Update(id int, post *db.Post) (*db.Post, error) {
	// Init
	var err error
	// Find post by id
	found, err := r.FindById(id)
	if err != nil {
		fmt.Println("Post to update not found: ", err)
		return nil, err
	}
	// Set post user id (gorm requires this as it does not automatically set the foreign key)
	if post.UserID == 0 {
		post.UserID = post.User.ID
	}

	// Update post using found post
	updateResult := r.DB.Model(&found).Updates(post)
	if updateResult.Error != nil {
		fmt.Println("Post update failed: ", updateResult.Error)
		return nil, updateResult.Error
	}

	// Retrieve changed post by id
	updated, err := r.FindById(id)
	if err != nil {
		fmt.Println("Post to update not found: ", err)
		return nil, err
	}
	return updated, nil
}
