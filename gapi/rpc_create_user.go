package gapi

import (
	"context"

	db "github.com/Mersock/golang-sample-bank/db/sqlc"
	"github.com/Mersock/golang-sample-bank/pb"
	"github.com/Mersock/golang-sample-bank/util"
	"github.com/Mersock/golang-sample-bank/val"
	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserReq) (*pb.CreateuserRes, error) {
	violations := validateCreateuserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	hashPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)

	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username already exists: %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)

	}

	res := &pb.CreateuserRes{
		User: convertUser(user),
	}
	return res, nil
}

func validateCreateuserRequest(req *pb.CreateUserReq) (violation []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violation = append(violation, fieldViolation("username", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violation = append(violation, fieldViolation("password", err))
	}

	if err := val.ValidateFullname(req.GetFullName()); err != nil {
		violation = append(violation, fieldViolation("full_name", err))
	}

	if err := val.ValidateEmail(req.GetFullName()); err != nil {
		violation = append(violation, fieldViolation("email", err))
	}
	return violation
}
