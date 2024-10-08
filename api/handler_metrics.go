package api

import (
	"context"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/s-hammon/recipls/internal/database"
)

const maxLimit = 1000

const (
	ErrQueryLimitVal   = "limit must be a positive, non-zero integer"
	ErrQueryLimitThres = "limit must be less than or equal to 1000"
)

type UserForMetrics struct {
	CreatedAt        time.Time `json:"created_at"`
	Name             string    `json:"name"`
	RecipesPublished int       `json:"recipes_published"`
}

type RecipeForMetrics struct {
	CreatedAt  time.Time `json:"created_at"`
	Title      string    `json:"title"`
	Difficulty int       `json:"difficulty"`
	Steps      int       `json:"steps"`
	Category   string    `json:"category"`
}

func (c *config) handlerGetMetrics(w http.ResponseWriter, r *http.Request, user database.User) {
	type response struct {
		Users   []UserForMetrics
		Recipes []RecipeForMetrics
	}

	limit := 100
	reqLimit := r.URL.Query().Get("limit")
	if reqLimit != "" {
		intLimit, err := strconv.Atoi(reqLimit)
		if err != nil || intLimit < 1 {
			respondError(w, http.StatusBadRequest, ErrQueryLimitVal)
			return
		}
		if intLimit > maxLimit {
			respondError(w, http.StatusBadRequest, ErrQueryLimitThres)
			return
		}
		limit = intLimit
	}

	channels := metricsChannels{
		usersCh:   make(chan []UserForMetrics),
		recipesCh: make(chan []RecipeForMetrics),
		errCh:     make(chan error, 2),
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go c.fetchUsers(r.Context(), limit, &channels, &wg)

	wg.Add(1)
	go c.fetchRecipes(r.Context(), limit, &channels, &wg)

	go func() {
		wg.Wait()
		close(channels.usersCh)
		close(channels.recipesCh)
	}()

	var users []UserForMetrics
	var recipes []RecipeForMetrics

	select {
	case err := <-channels.errCh:
		respondError(w, http.StatusInternalServerError, "error fetching data: "+err.Error())
		return
	case users = <-channels.usersCh:
		recipes = <-channels.recipesCh
	}

	respondJSON(w, http.StatusOK, response{
		Users:   users,
		Recipes: recipes,
	})
}

type metricsChannels struct {
	usersCh   chan []UserForMetrics
	recipesCh chan []RecipeForMetrics
	errCh     chan error
}

func (c *config) fetchUsers(ctx context.Context, limit int, channels *metricsChannels, wg *sync.WaitGroup) {
	defer wg.Done()

	usersDB, err := c.DB.GetUsersWithLimit(ctx, int32(limit))
	if err != nil {
		channels.errCh <- err
		return
	}

	users := []UserForMetrics{}
	for _, u := range usersDB {
		userRecipes, err := c.DB.GetRecipesByUser(ctx, u.ID)
		if err != nil {
			channels.errCh <- err
			return
		}
		user := DBToUser(u)
		users = append(users, user.toMetrics(len(userRecipes)))
	}
	channels.usersCh <- users
}

func (c *config) fetchRecipes(ctx context.Context, limit int, channels *metricsChannels, wg *sync.WaitGroup) {
	defer wg.Done()

	recipesDB, err := c.DB.GetRecipesWithLimit(ctx, int32(limit))
	if err != nil {
		channels.errCh <- err
		return
	}

	recipes := []RecipeForMetrics{}
	for _, p := range recipesDB {
		category, err := c.DB.GetCategoryByID(ctx, p.CategoryID)
		if err != nil {
			channels.errCh <- err
			return
		}
		recipe := DBToRecipe(p)
		recipes = append(recipes, recipe.toMetrics(category.Name))
	}
	channels.recipesCh <- recipes

}
