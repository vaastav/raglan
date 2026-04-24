package workflow

import (
	"context"
	"errors"
	"log"

	"github.com/blueprint-uservices/blueprint/runtime/core/backend"
	"github.com/vaastav/iridescent/iridescent_rt/autotune"
)

type Info struct {
	Fname    string
	Lname    string
	Address  string
	Email    string
	Country  string
	Username string
	UserID   int64
	Password string
}

type User struct {
	Username string
	UserID   int64
}

func (i Info) remote() {}

type UserService interface {
	RegsiterUser(ctx context.Context, username string, fname string, lname string, password string, userID int64, email string, address string, country string) error
	GetUserInfo(ctx context.Context, user User) (Info, error)
}

type UserServiceImpl struct {
	NaDB          backend.NoSQLDatabase
	EuDB          backend.NoSQLDatabase
	ApDB          backend.NoSQLDatabase
	SaDB          backend.NoSQLDatabase
	AfDB          backend.NoSQLDatabase
	OcDB          backend.NoSQLDatabase
	Fn            func(ctx context.Context, u *UserServiceImpl, username string) (Info, error)
	IsInitialized bool
}

func NewUserServieImpl(ctx context.Context, nadb backend.NoSQLDatabase, eudb backend.NoSQLDatabase, apdb backend.NoSQLDatabase, sadb backend.NoSQLDatabase, afdb backend.NoSQLDatabase, ocdb backend.NoSQLDatabase) (UserService, error) {
	impl := &UserServiceImpl{NaDB: nadb, EuDB: eudb, ApDB: apdb, SaDB: sadb, AfDB: afdb, OcDB: ocdb}
	go impl.init_service()
	return impl, nil
}

func (u *UserServiceImpl) init_service() {
	initialized := false
	for !initialized {
		rt := autotune.GetRuntime()
		if rt == nil {
			continue
		}
		update_fn := func() error {
			srt := autotune.GetRuntime().SpecRT
			if srt == nil {
				return errors.New("Runtime is nil oops!")
			}
			fn, err := srt.Lookup("GetUserInfo")
			if err != nil {
				return err
			}
			var ok bool
			u.Fn, ok = fn.(func(context.Context, *UserServiceImpl, string) (Info, error))
			if !ok {
				return errors.New("Failed to convert loaded symbol into desired type")
			}
			return nil
		}
		err := update_fn()
		if err != nil {
			log.Fatal(err)
		}
		rt.SpecRT.AddCallbackFn(update_fn)
	}
	log.Println("Service initialization complete")
	u.IsInitialized = true
}

func (u *UserServiceImpl) RegsiterUser(ctx context.Context, username string, fname string, lname string, password string, userID int64, email string, address string, country string) error {
	info := Info{Username: username, UserID: userID, Fname: fname, Lname: lname, Password: password, Email: email, Address: address, Country: country}
	if IsEu(country) {
		coll, err := u.EuDB.GetCollection(ctx, "user", "user")
		if err != nil {
			return err
		}
		return coll.InsertOne(ctx, info)
	} else if IsNa(country) {
		coll, err := u.NaDB.GetCollection(ctx, "user", "user")
		if err != nil {
			return err
		}
		return coll.InsertOne(ctx, info)
	} else if IsAs(country) {
		coll, err := u.ApDB.GetCollection(ctx, "user", "user")
		if err != nil {
			return err
		}
		return coll.InsertOne(ctx, info)
	} else if IsSa(country) {
		coll, err := u.EuDB.GetCollection(ctx, "user", "user")
		if err != nil {
			return err
		}
		return coll.InsertOne(ctx, info)
	} else if IsAf(country) {
		coll, err := u.AfDB.GetCollection(ctx, "user", "user")
		if err != nil {
			return err
		}
		return coll.InsertOne(ctx, info)
	} else {
		coll, err := u.OcDB.GetCollection(ctx, "user", "user")
		if err != nil {
			return err
		}
		return coll.InsertOne(ctx, info)
	}
}

func (u *UserServiceImpl) GetUserInfo(ctx context.Context, user User) (Info, error) {
	if !u.IsInitialized {
		return Info{}, errors.New("Function has not initialized yet")
	}
	// Make call into the guest code!
	return u.Fn(ctx, u, user.Username)
}
