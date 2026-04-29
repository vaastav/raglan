package workflow

import (
	"context"
	"errors"
	"log"
	"sort"
	"time"

	"github.com/blueprint-uservices/blueprint/runtime/core/backend"
	"github.com/vaastav/raglan/iridescent_rt/autotune"
	"github.com/vaastav/raglan/iridescent_rt/pass"
	"go.mongodb.org/mongo-driver/bson"
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
	GetUserInfoHard1(ctx context.Context, user User) (Info, error)
	GetUserInfoHard2(ctx context.Context, user User) (Info, error)
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
	Pass          *pass.ReorderIfPass
}

func NewUserServieImpl(ctx context.Context, nadb backend.NoSQLDatabase, eudb backend.NoSQLDatabase, apdb backend.NoSQLDatabase, sadb backend.NoSQLDatabase, afdb backend.NoSQLDatabase, ocdb backend.NoSQLDatabase) (UserService, error) {
	impl := &UserServiceImpl{NaDB: nadb, EuDB: eudb, ApDB: apdb, SaDB: sadb, AfDB: afdb, OcDB: ocdb}
	// Wait for initialization to complete!
	impl.init_service()
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
		// Add a custom if reorder pass
		p := pass.NewReorderIfPass()
		rt.SpecRT.AddSpecializationPass(p)
		go u.Policy()
		u.Pass = p
	}
	log.Println("Service initialization complete")
	u.IsInitialized = true
}

func (u *UserServiceImpl) Policy() {
	rt := autotune.GetRuntime().SpecRT
	pt := rt.PtsMap["userinfo"]
	for {
		// Sleep for 5 seconds
		time.Sleep(5 * time.Second)
		// Update the if-else reorder
		indices := make([]int, len(pt.Counter))
		for i := range indices {
			indices[i] = i
		}
		sort.Slice(indices, func(i, j int) bool {
			return pt.Counter[indices[i]] > pt.Counter[indices[j]]
		})
		log.Println("Selected Order:", indices)
		u.Pass.SetOrder("userinfo", indices)
		// Reset the stats for this point
		pt.ResetStats()
		// Update the plugin
		err := rt.UpdatePlugin()
		if err != nil {
			log.Fatal(err)
		}
	}
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

func find_user_helper(ctx context.Context, db backend.NoSQLDatabase, username string) (Info, error) {
	coll, err := db.GetCollection(ctx, "user", "user")
	if err != nil {
		return Info{}, err
	}
	query := bson.D{{"username", username}}
	var info Info
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

func (u *UserServiceImpl) GetUserInfoHard1(ctx context.Context, user User) (Info, error) {
	username := user.Username
	if info, err := find_user_helper(ctx, u.NaDB, username); err == nil {
		return info, err
	} else if info, err := find_user_helper(ctx, u.EuDB, username); err == nil {
		return info, err
	} else if info, err := find_user_helper(ctx, u.ApDB, username); err == nil {
		return info, err
	} else if info, err := find_user_helper(ctx, u.SaDB, username); err == nil {
		return info, err
	} else if info, err := find_user_helper(ctx, u.AfDB, username); err == nil {
		return info, err
	} else if info, err := find_user_helper(ctx, u.OcDB, username); err == nil {
		return info, err
	}
	return Info{}, errors.New("User does not exist")
}

func (u *UserServiceImpl) GetUserInfoHard2(ctx context.Context, user User) (Info, error) {
	username := user.Username
	if info, err := find_user_helper(ctx, u.OcDB, username); err == nil {
		return info, err
	} else if info, err := find_user_helper(ctx, u.AfDB, username); err == nil {
		return info, err
	} else if info, err := find_user_helper(ctx, u.SaDB, username); err == nil {
		return info, err
	} else if info, err := find_user_helper(ctx, u.ApDB, username); err == nil {
		return info, err
	} else if info, err := find_user_helper(ctx, u.EuDB, username); err == nil {
		return info, err
	} else if info, err := find_user_helper(ctx, u.NaDB, username); err == nil {
		return info, err
	}
	return Info{}, errors.New("User does not exist")
}
