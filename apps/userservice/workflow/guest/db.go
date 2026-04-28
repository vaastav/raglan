package main

import (
	"context"
	"errors"

	"github.com/blueprint-uservices/blueprint/runtime/core/backend"
	"github.com/vaastav/raglan/apps/userservice/workflow"
	"go.mongodb.org/mongo-driver/bson"
)

func find_user(ctx context.Context, db backend.NoSQLDatabase, username string) (workflow.Info, error) {
	coll, err := db.GetCollection(ctx, "user", "user")
	if err != nil {
		return workflow.Info{}, err
	}
	query := bson.D{{"username", username}}
	var info workflow.Info
	res, err := coll.FindOne(ctx, query)
	if err != nil {
		return info, err
	}
	ok, err := res.One(ctx, &info)
	if err != nil {
		return info, err
	}
	if !ok {
		return info, errors.New("Unable to deserialize database result!")
	}
	return info, nil
}

func GetUserInfo(ctx context.Context, u *workflow.UserServiceImpl, username string) (workflow.Info, error) {

	if info, err := find_user(ctx, u.NaDB, username); err == nil {
		return info, err
	} else if info, err := find_user(ctx, u.EuDB, username); err == nil {
		return info, err
	} else if info, err := find_user(ctx, u.ApDB, username); err == nil {
		return info, err
	} else if info, err := find_user(ctx, u.SaDB, username); err == nil {
		return info, err
	} else if info, err := find_user(ctx, u.AfDB, username); err == nil {
		return info, err
	} else if info, err := find_user(ctx, u.OcDB, username); err == nil {
		return info, err
	}
	return workflow.Info{}, errors.New("User does not exist!")
}
