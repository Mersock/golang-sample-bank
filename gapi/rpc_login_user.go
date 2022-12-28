package gapi

import (
	"context"
	"database/sql"

	db "github.com/Mersock/golang-sample-bank/db/sqlc"
	"github.com/Mersock/golang-sample-bank/pb"
	"github.com/Mersock/golang-sample-bank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginReq) (*pb.LoginRes, error) {
	user, err := server.store.GetUsers(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "Not found user")
		}
		return nil, status.Errorf(codes.Internal, "failed to find user")
	}

	err = util.CheckPassword(req.GetPassword(), user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "incorrect password")

	}

	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create access token error")
	}

	refreshToken, refreshTokenPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create refresh token error")

	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshTokenPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    "",
		ClientIp:     "",
		IsBlocked:    false,
		ExpiresAt:    refreshTokenPayload.ExpireAt,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "create session token error")

	}

	res := &pb.LoginRes{
		User:                  convertUser(user),
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accessTokenPayload.ExpireAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshTokenPayload.ExpireAt),
	}

	return res, nil
}
