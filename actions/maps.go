package actions

import (
	"database/sql"
	"github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/danigunawan/rpg-api/models"
	"net/http"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Map)
// DB Table: Plural (maps)
// Resource: Plural (Maps)
// Path: Plural (/maps)
// View Template Folder: Plural (/templates/maps/)

// MapsResource is the resource for the Map model
type MapsResource struct {
	buffalo.Resource
}

// List gets all Maps. This function is mapped to the path
// GET /maps/{quest_id}
func (v MapsResource) List(c buffalo.Context) error {
	claims := c.Value("claims").(jwt.MapClaims)
	user_id, err := uuid.FromString(claims["id"].(string))
	if err != nil {
		return c.Error(http.StatusInternalServerError, errors.New("bad user id"))
	}

	quest_id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return c.Error(http.StatusInternalServerError, errors.New("bad quest id"))
	}

	quest := &models.Quest{}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}

	// Retrieve all Maps from the DB
	if err := tx.Eager("Maps").Where("user_id = ?", user_id).Find(quest, quest_id); err != nil {
		return c.Error(http.StatusNotFound, errors.New("maps not found"))
	}

	return c.Render(http.StatusOK, r.JSON(map[string]models.Maps{
		"maps": quest.Maps,
	}))
}

// Show gets the data for one Map. This function is mapped to
// the path GET /maps/{map_id}
func (v MapsResource) Show(c buffalo.Context) error {
	claims := c.Value("claims").(jwt.MapClaims)
	user_id, err := uuid.FromString(claims["id"].(string))
	if err != nil {
		return c.Error(http.StatusInternalServerError, errors.New("bad user id"))
	}

	map_id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return c.Error(http.StatusInternalServerError, errors.New("bad map id"))
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}

	// Allocate an empty Map
	rpg_map := &models.Map{}

	// To find the Map the parameter map_id is used.
	if err := tx.Eager().Where("user_id = ?", user_id).Find(rpg_map, map_id); err != nil {
		return c.Error(http.StatusNotFound, errors.New("map not found"))
	}

	return c.Render(200, r.JSON(rpg_map))
}

// Create adds a Map to the DB. This function is mapped to the
// path POST /maps
func (v MapsResource) Create(c buffalo.Context) error {
	claims := c.Value("claims").(jwt.MapClaims)
	user_id, err := uuid.FromString(claims["id"].(string))
	if err != nil {
		return c.Error(http.StatusInternalServerError, errors.New("bad user id"))
	}

	// Allocate an empty Map
	rpg_map := &models.Map{}

	// Bind map to the html form elements
	if err := c.Bind(rpg_map); err != nil {
		return err
	}

	rpg_map.UserID = user_id

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}

	// Get Map of Highest Sort Order
	last_map := &models.Map{}
	err = tx.Where("user_id = ?", user_id).Where("quest_id = ?", rpg_map.QuestID).Order("sort_order desc").First(last_map)
	if err != nil {
		if errors.Cause(err) != sql.ErrNoRows {
			c.Logger().Error(err)
			return errors.New("transaction error")
		}
	}

	// last_map.SortOrder defaults to 0 if no maps found
	rpg_map.SortOrder = last_map.SortOrder + 1

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(rpg_map)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return c.Error(http.StatusUnprocessableEntity, verrs)
	}

	return c.Render(http.StatusCreated, r.JSON(rpg_map))
}

// Update changes a Map in the DB. This function is mapped to
// the path PUT /maps/{map_id}
func (v MapsResource) Update(c buffalo.Context) error {
	claims := c.Value("claims").(jwt.MapClaims)
	user_id, err := uuid.FromString(claims["id"].(string))
	if err != nil {
		return c.Error(http.StatusInternalServerError, errors.New("bad user id"))
	}

	map_id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return c.Error(http.StatusInternalServerError, errors.New("bad map id"))
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}

	// Allocate an empty Map
	rpg_map := &models.Map{}

	rpg_map.UserID = user_id

	if err := tx.Where("user_id = ?", user_id).Find(rpg_map, map_id); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Bind Map to the html form elements
	if err := c.Bind(rpg_map); err != nil {
		return err
	}

	verrs, err := tx.ValidateAndUpdate(rpg_map)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return c.Error(http.StatusUnprocessableEntity, verrs)
	}

	return c.Render(http.StatusOK, r.JSON(rpg_map))
}

// Destroy deletes a Map from the DB. This function is mapped
// to the path DELETE /maps/{map_id}
func (v MapsResource) Destroy(c buffalo.Context) error {
	claims := c.Value("claims").(jwt.MapClaims)
	user_id, err := uuid.FromString(claims["id"].(string))
	if err != nil {
		return c.Error(http.StatusInternalServerError, errors.New("bad user id"))
	}

	map_id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return c.Error(http.StatusInternalServerError, errors.New("bad map id"))
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}

	// Allocate an empty Map
	rpg_map := &models.Map{}

	// To find the Map the parameter map_id is used.
	if err := tx.Where("user_id = ?", user_id).Find(rpg_map, map_id); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	levels := models.Levels{}
	if err := tx.Where("user_id = ?", user_id).Where("map_id = ?", rpg_map.ID).All(&levels); err != nil {
		if errors.Cause(err) != sql.ErrNoRows {
			return c.Error(http.StatusInternalServerError, errors.New("cannot select levels"))
		}
	}
	if err := tx.Destroy(levels); err != nil {
		return c.Error(http.StatusInternalServerError, errors.New("could not destroy levels"))
	}

	if err := tx.Destroy(rpg_map); err != nil {
		return err
	}

	return c.Render(http.StatusOK, r.JSON(rpg_map))
}
