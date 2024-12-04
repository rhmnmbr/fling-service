package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
	db "github.com/rhmnmbr/fling-service/db/sqlc"
	"github.com/rhmnmbr/fling-service/util"
)

var ErrInvalidEmailOrPassword = errors.New("invalid email/password")

type createUserRequest struct {
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=6"`
	Phone        string `json:"phone" binding:"required,e164"`
	FirstName    string `json:"first_name" binding:"required"`
	BirthDate    string `json:"birth_date" binding:"required,datetime=2006-01-02" time_format:"2006-01-02"`
	Gender       string `json:"gender" binding:"required,oneof=male female"`
	LocationInfo string `json:"location_info" binding:"omitempty"`
	Bio          string `json:"bio" binding:"omitempty"`
}

type userResponse struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	FirstName    string    `json:"first_name"`
	BirthDate    time.Time `json:"birth_date"`
	Gender       string    `json:"gender"`
	LocationInfo *string   `json:"location_info"`
	Bio          *string   `json:"bio"`
	CreatedAt    time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	var locationInfo *string
	if user.LocationInfo.Valid {
		locationInfo = &user.LocationInfo.String
	}

	var bio *string
	if user.Bio.Valid {
		bio = &user.Bio.String
	}

	return userResponse{
		ID:           user.ID,
		Email:        user.Email,
		Phone:        user.Phone,
		FirstName:    user.FirstName,
		BirthDate:    user.BirthDate,
		Gender:       string(user.Gender),
		LocationInfo: locationInfo,
		Bio:          bio,
		CreatedAt:    user.CreatedAt,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	parsedBirthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Email:          req.Email,
		HashedPassword: hashedPassword,
		Phone:          req.Phone,
		FirstName:      req.FirstName,
		BirthDate:      parsedBirthDate,
		Gender:         db.GenderEnum(req.Gender),
		LocationInfo:   sql.NullString{String: req.LocationInfo, Valid: req.LocationInfo != ""},
		Bio:            sql.NullString{String: req.Bio, Valid: req.Bio != ""},
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(ErrInvalidEmailOrPassword))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrInvalidEmailOrPassword))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(
		user.ID,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := loginUserResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, rsp)
}
