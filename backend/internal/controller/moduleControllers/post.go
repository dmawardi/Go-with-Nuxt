package modulecontrollers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dmawardi/Go-Template/internal/helpers/request"
	schemamodels "github.com/dmawardi/Go-Template/internal/models/schemaModels"
	moduleservices "github.com/dmawardi/Go-Template/internal/service/module"
	"github.com/go-chi/chi/v5"
)

type PostController interface {
	FindAll(w http.ResponseWriter, r *http.Request)
	Find(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type postController struct {
	service moduleservices.PostService
}

func NewPostController(service moduleservices.PostService) PostController {
	return &postController{service}
}

// Used to init the query params for easy extraction in controller
// Returns: map[string]string{"age": "int", "name": "string", "active": "bool"}
func PostConditionQueryParams() map[string]string {
	return map[string]string{
		"title": "string",
		"body":  "string",
	}
}

// API/POSTS
// @Summary      Finds a list of posts
// @Description  Accepts limit, offset, order, search (added as non-case sensitive LIKE), title, body as query parameters
// @Tags         Post
// @Accept       json
// @Produce      json
// @Param        limit   query      int  true  "limit"
// @Param        offset   query      int  false  "offset"
// @Param        order   query      int  false  "order by eg. (asc) "id" (desc) "id_desc" )"
// @Param        search   query      string  false  "search (added to all string conditions as LIKE SQL search)"
// @Param        title   query      string  false  "title"
// @Param        body   query      string  false  "body"
// @Success      200 {object} models.PaginatedPosts
// @Failure      400 {string} string "Can't find posts"
// @Failure      400 {string} string "Must include limit parameter with a max value of 50"
// @Failure 	400 {string} string "Error extracting query params"
// @Router       /posts [get]
// @Security BearerToken
func (c postController) FindAll(w http.ResponseWriter, r *http.Request) {
	// Grab basic query params
	baseQueryParams, err := request.ExtractBasicFindAllQueryParams(r)
	if err != nil {
		http.Error(w, "Error extracting query params", http.StatusBadRequest)
		return
	}

	// Generate query params to extract
	queryParamsToExtract := PostConditionQueryParams()
	// Extract query params
	extractedConditionParams, err := request.ExtractSearchAndConditionParams(r, queryParamsToExtract)
	if err != nil {
		http.Error(w, "Error extracting query params", http.StatusBadRequest)
		return
	}

	// Check that limit is present as requirement
	if (baseQueryParams.Limit == 0) || (baseQueryParams.Limit > 50) {
		http.Error(w, "Must include limit parameter with a max value of 50", http.StatusBadRequest)
		return
	}

	// Query database for all users using query params
	found, err := c.service.FindAll(baseQueryParams.Limit, baseQueryParams.Offset, baseQueryParams.Order, extractedConditionParams)
	if err != nil {
		http.Error(w, "Can't find posts", http.StatusBadRequest)
		return
	}
	err = request.WriteAsJSON(w, found)
	if err != nil {
		http.Error(w, "Can't find posts", http.StatusBadRequest)
		fmt.Println("error writing users to response: ", err)
		return
	}
}

// @Summary      Find post
// @Description  Find a post by ID
// @Tags         Post
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Post ID"
// @Success      200 {object} db.Post
// @Failure      400 {string} string "Can't find post with ID: {id}"
// @Router       /posts/{id} [get]
// @Security BearerToken
func (c postController) Find(w http.ResponseWriter, r *http.Request) {
	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, err := strconv.Atoi(stringParameter)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	found, err := c.service.FindById(idParameter)
	if err != nil {
		http.Error(w, fmt.Sprintf("Can't find post with ID: %v\n", idParameter), http.StatusBadRequest)
		return
	}
	err = request.WriteAsJSON(w, found)
	if err != nil {
		http.Error(w, fmt.Sprintf("Can't find post with ID: %v\n", idParameter), http.StatusBadRequest)
		return
	}
}

// @Summary      Create Post
// @Description  Creates a new post
// @Tags         Post
// @Accept       json
// @Produce      plain
// @Param        post body schemamodels.CreatePost true "New Post"
// @Failure      400 {object} models.ValidationError "Validation Errors"
// @Success      201 {string} string "Post creation successful!"
// @Failure      400 {string} string "Post creation failed."
// @Router       /posts [post]
// @Security BearerToken
func (c postController) Create(w http.ResponseWriter, r *http.Request) {
	// Init
	var toCreate schemamodels.CreatePost
	// Decode request body as JSON and store in login
	err := json.NewDecoder(r.Body).Decode(&toCreate)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}

	// Validate the incoming DTO
	pass, valErrors := request.GoValidateStruct(&toCreate)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		request.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes and allow through

	// Create post
	_, createErr := c.service.Create(&toCreate)
	if createErr != nil {
		http.Error(w, "Post creation failed.", http.StatusBadRequest)
		return
	}

	// Set status to created
	w.WriteHeader(http.StatusCreated)
	// Send post success message in body
	w.Write([]byte("Post creation successful!"))
}

// @Summary      Update Post
// @Description  Updates an existing post
// @Tags         Post
// @Accept       json
// @Produce      json
// @Param        post body schemamodels.UpdatePost true "Update Post"
// @Param        id   path      int  true  "Post ID"
// @Success      200 {object} db.Post
// @Failure      400 {object} models.ValidationError "Validation Errors"
// @Failure      400 {string} string "Failed post update"
// @Failure      403 {string} string "Authentication Token not detected"
// @Router       /posts/{id} [put]
// @Security BearerToken
func (c postController) Update(w http.ResponseWriter, r *http.Request) {
	// grab id parameter
	var toUpdate schemamodels.UpdatePost
	// Decode request body as JSON and store in login
	err := json.NewDecoder(r.Body).Decode(&toUpdate)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}

	// Validate the incoming DTO
	pass, valErrors := request.GoValidateStruct(&toUpdate)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		request.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes and allow through

	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, _ := strconv.Atoi(stringParameter)

	// Update post
	updated, createErr := c.service.Update(idParameter, &toUpdate)
	if createErr != nil {
		http.Error(w, fmt.Sprintf("Failed post update: %s", createErr), http.StatusBadRequest)
		return
	}
	// Write post to output
	err = request.WriteAsJSON(w, updated)
	if err != nil {
		fmt.Printf("Error encountered when writing to JSON. Err: %s", err)
	}
}

// @Summary      Delete Post
// @Description  Deletes an existing post
// @Tags         Post
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Post ID"
// @Success      200 {string} string "Deletion successful!"
// @Failure      400 {string} string "Failed post deletion"
// @Router       /posts/{id} [delete]
// @Security BearerToken
func (c postController) Delete(w http.ResponseWriter, r *http.Request) {
	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, _ := strconv.Atoi(stringParameter)

	// Attampt to delete post using id
	err := c.service.Delete(idParameter)

	// If error detected
	if err != nil {
		http.Error(w, "Failed post deletion", http.StatusBadRequest)
		return
	}
	// Else write success
	w.Write([]byte("Deletion successful!"))
}
