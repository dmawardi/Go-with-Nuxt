package moduleservices

import (
	"fmt"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/models"
	schemamodels "github.com/dmawardi/Go-Template/internal/models/schemaModels"
	modulerepositories "github.com/dmawardi/Go-Template/internal/repository/module"
)

type PostService interface {
	FindAll(limit int, offset int, order string, conditions []models.QueryConditionParameters) (*models.BasicPaginatedResponse[db.Post], error)
	FindById(int) (*db.Post, error)
	Create(post *schemamodels.CreatePost) (*db.Post, error)
	Update(int, *schemamodels.UpdatePost) (*db.Post, error)
	Delete(int) error
	BulkDelete([]int) error
}

type postService struct {
	repo modulerepositories.PostRepository
}

func NewPostService(repo modulerepositories.PostRepository) PostService {
	return &postService{repo: repo}
}

// Creates a post in the database
func (s *postService) Create(post *schemamodels.CreatePost) (*db.Post, error) {
	// Create a new user of type db User
	toCreate := db.Post{
		Title: post.Title,
		Body:  post.Body,
		User:  post.User,
	}

	// Create above post in database
	created, err := s.repo.Create(&toCreate)
	if err != nil {
		return nil, fmt.Errorf("failed creating post: %w", err)
	}

	return created, nil
}

// Find a list of posts in the database
func (s *postService) FindAll(limit int, offset int, order string, conditions []models.QueryConditionParameters) (*models.BasicPaginatedResponse[db.Post], error) {
	posts, err := s.repo.FindAll(limit, offset, order, conditions)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

// Find post in database by ID
func (s *postService) FindById(id int) (*db.Post, error) {
	// Search cache
	// Define a key with a naming convention
	cacheKey := fmt.Sprintf("post:%d", id)
	// Check if post is in cache
	cachedPost, found := app.Cache.Load(cacheKey)
	if found {
		// If found, return cached post
		return cachedPost.(*db.Post), nil
	}

	// Find post by id
	post, err := s.repo.FindById(id)
	// If error detected
	if err != nil {
		return nil, err
	}

	// Store post in cache
	app.Cache.Store(cacheKey, post)

	return post, nil
}

// Delete post in database
func (s *postService) Delete(id int) error {
	err := s.repo.Delete(id)
	// If error detected
	if err != nil {
		fmt.Println("error in deleting post: ", err)
		return err
	}
	// else
	cacheKey := fmt.Sprintf("post:%d", id)
	app.Cache.Delete(cacheKey)
	return nil
}

// Deletes multiple users in database
func (s *postService) BulkDelete(ids []int) error {
	err := s.repo.BulkDelete(ids)
	// If error detected
	if err != nil {
		fmt.Println("error in bulk deleting users: ", err)
		return err
	}
	// else
	for _, id := range ids {
		cacheKey := fmt.Sprintf("post:%d", id)
		app.Cache.Delete(cacheKey)
	}
	return nil
}

// Updates post in database
func (s *postService) Update(id int, post *schemamodels.UpdatePost) (*db.Post, error) {
	// Create db Post type from incoming DTO
	toUpdate := &db.Post{
		Title: post.Title,
		Body:  post.Body,
		User:  post.User,
	}

	// Update using repo
	updated, err := s.repo.Update(id, toUpdate)
	if err != nil {
		return nil, err
	}

	// Store updated post in cache
	cacheKey := fmt.Sprintf("post:%d", id)
	app.Cache.Store(cacheKey, updated)

	return updated, nil
}
