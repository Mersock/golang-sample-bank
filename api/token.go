package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type renewAccessTokenUserReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenUserRes struct {
	AccessToken         string    `json:"access_token"`
	AccessTokenExpireAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) renewAccessTokenUser(ctx *gin.Context) {
	var req renewAccessTokenUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}

	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errRes(err))
		return
	}

	session, err := server.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errRes(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}

	if session.IsBlocked {
		err := fmt.Errorf("blocked session")
		ctx.JSON(http.StatusUnauthorized, errRes(err))
		return
	}

	if session.Username != refreshPayload.Username {
		err := fmt.Errorf("incorrect session user")
		ctx.JSON(http.StatusUnauthorized, errRes(err))
		return
	}

	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("Mismatched session token")
		ctx.JSON(http.StatusUnauthorized, errRes(err))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("expire session")
		ctx.JSON(http.StatusUnauthorized, errRes(err))
		return
	}

	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(refreshPayload.Username, server.config.AccessTokenDuration)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errRes(err))
		return
	}
	res := renewAccessTokenUserRes{
		AccessToken:         accessToken,
		AccessTokenExpireAt: accessTokenPayload.ExpireAt,
	}

	ctx.JSON(http.StatusOK, res)
}
