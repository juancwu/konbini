package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/juancwu/konbini/server/db"
	"github.com/juancwu/konbini/server/middlewares"
	"github.com/juancwu/konbini/server/permission"
	"github.com/juancwu/konbini/server/utils"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type ingredient struct {
	Name  string `json:"name" validate:"required,min=1,printascii"`
	Value string `json:"value"`
}

type NewBentoRequest struct {
	Name        string       `json:"name" validate:"required,min=3,printascii"`
	Ingredients []ingredient `json:"ingredients,omitempty" validate:"omitnil,omitempty,dive"`
}

func NewBento(connector *db.DBConnector) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := middlewares.GetUser(c)
		if err != nil {
			return err
		}
		if !user.EmailVerified {
			return APIError{
				Code:          http.StatusBadRequest,
				PublicMessage: "Email must be verified before creating a new bento",
			}
		}

		body, err := middlewares.GetJsonBody[NewBentoRequest](c)
		if err != nil {
			return err
		}

		logger := middlewares.GetLogger(c)

		// shouldn't take more than 5 seconds to run
		ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
		defer cancel()

		conn, err := connector.Connect()
		if err != nil {
			return err
		}
		defer conn.Close()
		tx, err := conn.Begin()
		if err != nil {
			return err
		}

		q := db.New(tx)

		exists, err := q.ExistsBentoWithNameOwnedByUser(
			ctx,
			db.ExistsBentoWithNameOwnedByUserParams{
				Name:   body.Name,
				UserID: user.ID,
			},
		)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if exists == 1 {
			return APIError{
				Code:          http.StatusBadRequest,
				PublicMessage: fmt.Sprintf("Bento with name %s already exists.", body.Name),
			}
		}

		// create bento
		bentoID, err := q.NewBento(
			ctx,
			db.NewBentoParams{
				Name:      body.Name,
				UserID:    user.ID,
				CreatedAt: utils.FormatRFC3339NanoFixed(time.Now()),
				UpdatedAt: utils.FormatRFC3339NanoFixed(time.Now()),
			},
		)
		if err != nil {
			db.RollabackWithLog(tx, logger)
			return err
		}

		if body.Ingredients != nil && len(body.Ingredients) > 0 {
			timestamp := utils.FormatRFC3339NanoFixed(time.Now())
			for _, ing := range body.Ingredients {
				err = q.AddIngredientToBento(
					ctx,
					db.AddIngredientToBentoParams{
						BentoID:   bentoID,
						Name:      ing.Name,
						Value:     []byte(ing.Value),
						CreatedAt: timestamp,
						UpdatedAt: timestamp,
					},
				)
				if err != nil {
					db.RollabackWithLog(tx, logger)
					return err
				}
			}
		}

		// create new bento permissions for the owner
		timestamp := utils.FormatRFC3339NanoFixed(time.Now())
		err = q.NewBentoPermission(
			ctx,
			db.NewBentoPermissionParams{
				UserID:    user.ID,
				BentoID:   bentoID,
				Bytes:     permission.ToBytes(permission.GetBentoOwnerPermissions()),
				CreatedAt: timestamp,
				UpdatedAt: timestamp,
			},
		)
		if err != nil {
			db.RollabackWithLog(tx, logger)
			return err
		}

		err = tx.Commit()
		if err != nil {
			db.RollabackWithLog(tx, logger)
			return err
		}

		return c.JSON(http.StatusCreated, map[string]string{"bento_id": bentoID})
	}
}

type AddIngredientsToBentoRequest struct {
	BentoID     string       `json:"bento_id" validate:"required,uuid4"`
	Ingredients []ingredient `json:"ingredients,omitempty" validate:"omitnil,omitempty,dive"`
}

// AddIngredientsToBento add the ingridients in the request body to the bento
func AddIngredientsToBento(cnt *db.DBConnector) echo.HandlerFunc {
	return func(c echo.Context) error {
		replace := c.QueryParam("replace") == "true"

		ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
		defer cancel()

		user, err := middlewares.GetUser(c)
		if err != nil {
			return err
		}
		body, err := middlewares.GetJsonBody[AddIngredientsToBentoRequest](c)
		if err != nil {
			return err
		}

		conn, err := cnt.Connect()
		if err != nil {
			return err
		}
		defer conn.Close()

		q := db.New(conn)

		bento, err := q.GetBentoWithIDOwnedByUser(
			ctx,
			db.GetBentoWithIDOwnedByUserParams{
				ID:     body.BentoID,
				UserID: user.ID,
			},
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return APIError{
					Code:           http.StatusBadRequest,
					PublicMessage:  "No bento found",
					PrivateMessage: "No bento found with given ID owned by requesting user",
					InternalError:  err,
				}
			}
			return err
		}

		tx, err := conn.Begin()
		if err != nil {
			return err
		}

		q = q.WithTx(tx)

		for _, ing := range body.Ingredients {
			if replace {
				err = q.SetBentoIngredient(
					ctx,
					db.SetBentoIngredientParams{
						BentoID:   bento.ID,
						Name:      ing.Name,
						Value:     []byte(ing.Value),
						CreatedAt: utils.FormatRFC3339NanoFixed(time.Now()),
						UpdatedAt: utils.FormatRFC3339NanoFixed(time.Now()),
					},
				)
			} else {
				err = q.AddIngredientToBento(
					ctx,
					db.AddIngredientToBentoParams{
						BentoID:   bento.ID,
						Name:      ing.Name,
						Value:     []byte(ing.Value),
						CreatedAt: utils.FormatRFC3339NanoFixed(time.Now()),
						UpdatedAt: utils.FormatRFC3339NanoFixed(time.Now()),
					},
				)
			}
			if err != nil {
				tx.Rollback()
				if !replace && utils.IsUniqueViolationErr(err) {
					return APIError{
						Code:          http.StatusBadRequest,
						PublicMessage: fmt.Sprintf("Ingridient with name '%s' already exists. To replace set the query parameter 'replace=true'.", ing.Name),
						InternalError: err,
					}
				}
				return err
			}
		}

		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			return err
		}

		return c.NoContent(http.StatusOK)
	}
}

type RemoveIngredientsFromBentoRequest struct {
	BentoID     string   `json:"bento_id" validate:"required,uuid4"`
	Ingredients []string `json:"ingredients" validate:"required,gt=0,dive,uuid4"`
}

// RemoveIngredientsFromBento removes ingridients by id from the bento
func RemoveIngredientsFromBento(cnt *db.DBConnector) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
		defer cancel()

		user, err := middlewares.GetUser(c)
		if err != nil {
			return err
		}
		body, err := middlewares.GetJsonBody[RemoveIngredientsFromBentoRequest](c)
		if err != nil {
			return err
		}

		conn, err := cnt.Connect()
		if err != nil {
			return err
		}
		defer conn.Close()

		tx, err := conn.Begin()
		if err != nil {
			return err
		}

		q := db.New(tx)

		bento, err := q.GetBentoWithIDOwnedByUser(
			ctx,
			db.GetBentoWithIDOwnedByUserParams{
				UserID: user.ID,
				ID:     body.BentoID,
			},
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return APIError{
					Code:          http.StatusBadRequest,
					PublicMessage: "No bento found",
				}
			}
			return err
		}

		deleted := []string{}
		notDeleted := []string{}

		for _, ingID := range body.Ingredients {
			n, err := q.RemoveIngredientFromBento(
				ctx,
				db.RemoveIngredientFromBentoParams{
					BentoID: bento.ID,
					ID:      ingID,
				},
			)
			if err != nil {
				tx.Rollback()
				return err
			}
			if n != 1 {
				notDeleted = append(notDeleted, ingID)
			} else {
				deleted = append(deleted, ingID)
			}
		}

		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			return err
		}

		return c.JSON(http.StatusOK, map[string][]string{
			"deleted":     deleted,
			"not_deleted": notDeleted,
		})
	}
}

// GetBento gets the bento info and ingridients
func GetBento(cnt *db.DBConnector) echo.HandlerFunc {
	return func(c echo.Context) error {
		bentoID := c.QueryParam("bento_id")
		if bentoID == "" {
			return APIError{
				Code:          http.StatusBadRequest,
				PublicMessage: "Missing bento_id query parameter.",
			}
		}

		user, err := middlewares.GetUser(c)
		if err != nil {
			return err
		}

		logger := middlewares.GetLogger(c)

		ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
		defer cancel()

		conn, err := cnt.Connect()
		if err != nil {
			return err
		}
		defer conn.Close()

		q := db.New(conn)

		bento, err := q.GetBentoByIDWithPermissions(
			ctx,
			db.GetBentoByIDWithPermissionsParams{
				UserID: user.ID,
				ID:     bentoID,
			},
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return APIError{
					Code:          http.StatusNotFound,
					PublicMessage: "Bento not found",
					InternalError: err,
				}
			}
			return err
		}

		// check if they have permission to read the bento
		u64Perms, err := permission.FromBytes(bento.Bytes)
		if err != nil {
			return err
		}
		if u64Perms&permission.Read == 0 {
			logger.Debug().Uint64("perms", u64Perms).Send()
			return APIError{
				Code:           http.StatusNotFound,
				PublicMessage:  "Bento not found",
				PrivateMessage: "No permissions to read bento",
			}
		}

		// get bento ingredients
		rows, err := q.GetBentoIngredients(
			ctx,
			bento.ID,
		)
		if err != nil && err != sql.ErrNoRows {
			return err
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"bento_id":    bento.ID,
			"name":        bento.Name,
			"ingredients": rows,
		})
	}
}

type ListBentosResponse struct {
	OwnerID    string `json:"owner_id"`
	BentoID    string `json:"bento_id"`
	BentoName  string `json:"bento_name"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	UserPerms  string `json:"user_perms"`
	GroupPerms string `json:"group_perms"`
}

// ListBentos gets a list of all the user's bentos. The list only contains
// basic information of the bentos, the same as the non-extended version of the metadata.
func ListBentos(cnt *db.DBConnector) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := middlewares.GetUser(c)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(c.Request().Context(), FiveSeconds)
		defer cancel()

		conn, err := cnt.Connect()
		if err != nil {
			return err
		}
		defer conn.Close()

		q := db.New(conn)

		rows, err := q.ListBentosWithAccess(ctx, user.ID)
		if err != nil && err != sql.ErrNoRows {
			return err
		}

		res := make([]ListBentosResponse, len(rows))
		for i := 0; i < len(rows); i++ {
			uPerms, err := permission.BytesToString(rows[i].UserPerms)
			if err != nil {
				return err
			}
			gPerms, err := permission.BytesToString(rows[i].GroupPerms)
			if err != nil {
				return err
			}
			res[i] = ListBentosResponse{
				OwnerID:    rows[i].OwnerID,
				BentoID:    rows[i].BentoID,
				BentoName:  rows[i].BentoName,
				UserPerms:  uPerms,
				GroupPerms: gPerms,
				CreatedAt:  rows[i].CreatedAt,
				UpdatedAt:  rows[i].UpdatedAt,
			}
		}

		return c.JSON(http.StatusOK, res)
	}
}
