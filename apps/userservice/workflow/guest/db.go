package main

import (
	"context"
	"errors"

	"github.com/blueprint-uservices/blueprint/runtime/core/backend"
	"github.com/vaastav/raglan/apps/userservice/workflow"
	"github.com/vaastav/raglan/iridescent_rt/autotune"
	"go.mongodb.org/mongo-driver/bson"
)

func find_user(ctx context.Context, db backend.NoSQLDatabase, username string, info *workflow.Info) bool {
	coll, err := db.GetCollection(ctx, "user", "user")
	if err != nil {
		return false
	}
	query := bson.D{{"username", username}}
	res, err := coll.FindOne(ctx, query)
	if err != nil {
		return false
	}
	ok, err := res.One(ctx, info)
	if err != nil {
		return false
	}
	return ok
}

func GetUserInfo(ctx context.Context, u *workflow.UserServiceImpl, username string) (workflow.Info, error) {
	iridescent_instr_range_int("userinfo", 0, 5)
	var info workflow.Info
	//reorder:if userinfo
	if find_user(ctx, u.NaDB, username, &info) {
		autotune.GetRuntime().SpecRT.Instrument("userinfo", 0)
		return info, nil
	} else if find_user(ctx, u.EuDB, username, &info) {
		autotune.GetRuntime().SpecRT.Instrument("userinfo", 1)
		return info, nil
	} else if find_user(ctx, u.ApDB, username, &info) {
		autotune.GetRuntime().SpecRT.Instrument("userinfo", 2)
		return info, nil
	} else if find_user(ctx, u.SaDB, username, &info) {
		autotune.GetRuntime().SpecRT.Instrument("userinfo", 3)
		return info, nil
	} else if find_user(ctx, u.AfDB, username, &info) {
		autotune.GetRuntime().SpecRT.Instrument("userinfo", 4)
		return info, nil
	} else if find_user(ctx, u.OcDB, username, &info) {
		autotune.GetRuntime().SpecRT.Instrument("userinfo", 5)
		return info, nil
	}
	return workflow.Info{}, errors.New("User does not exist!")
}
