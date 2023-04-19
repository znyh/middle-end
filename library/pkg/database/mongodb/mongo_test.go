package mongdb

import (
	"context"
	"sync"
	"testing"
	"time"

	xtime "github.com/go-kratos/kratos/pkg/time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	_MongoURL = "mongodb://127.0.0.1:27017"
)

func TestConnectPing(t *testing.T) {
	// Create a Client to a MongoDB server and use Ping to verify that the server is running.

	clientOpts := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			t.Fatal(err)
		}
	}()

	// Call Ping to verify that the deployment is up and the Client was configured successfully.
	// As mentioned in the Ping documentation, this reduces application resiliency as the server may be
	// temporarily unavailable when Ping is called.
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		t.Fatal(err)
	}
}

func TestInsertOne(t *testing.T) {
	// Create a Client and execute a ListDatabases operation.

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			t.Fatal(err)
		}
	}()

	collection := client.Database("db").Collection("coll")
	result, err := collection.InsertOne(context.TODO(), bson.D{{"x", 1}, {"y", 2}})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("inserted ID: %v\n", result.InsertedID)
}

func TestConnect(t *testing.T) {
	client := GetClient()

	collection := client.Database("db").Collection("coll")

	res, err := collection.InsertOne(context.Background(), bson.M{"key": "10001", "value": "{\"login_status\":true}"})
	if err != nil {
		panic(err)
	}

	id := res.InsertedID

	t.Logf("insert id %v", id)
}

func GetClient() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(_MongoURL))
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	return client
}

func TestUpdate(t *testing.T) {
	db := NewMongoDB(&Config{
		URI:            _MongoURL,
		ConnectTimeout: xtime.Duration(5 * time.Second),
		QueryTimeout:   xtime.Duration(1 * time.Second),
		ExecTimeout:    xtime.Duration(2 * time.Second),
	})

	collection := db.Database("db").Collection("coll")

	// update failed, insert it.
	opts := new(options.UpdateOptions)
	opts.SetUpsert(true)

	result, err := db.UpdateOne(context.Background(), collection,
		bson.M{"key": "10001.login"},
		bson.D{{"$set", bson.M{"value": "1234567"}}}, opts)

	if err != nil {
		t.Errorf("%v", err)
	}

	t.Logf("updateone result %+v", result)
}

func TestFind(t *testing.T) {
	db := NewMongoDB(&Config{
		URI:            _MongoURL,
		ConnectTimeout: xtime.Duration(5 * time.Second),
		QueryTimeout:   xtime.Duration(1 * time.Second),
		ExecTimeout:    xtime.Duration(2 * time.Second),
	})

	collection := db.Database("db").Collection("coll")

	ctx := context.Background()
	cur, err := db.Find(ctx, collection, bson.D{{"key", "10001.login"}})
	if err != nil {
		panic(err)
	}

	for cur.Next(ctx) {
		doc := cur.Current

		//方式一
		elem, err := doc.Elements()
		if err != nil {
			panic(err)
		}

		for _, v := range elem {
			t.Logf("key: %s  value: %s", v.Key(), v.Value().String())
		}

		//方式二
		t.Logf("==> %s %s %s",
			doc.Lookup("key").String(),
			doc.Lookup("value").String(),
			doc.Lookup("notExist").String())
	}
}

func TestDelete(t *testing.T) {
	db := NewMongoDB(&Config{
		URI:            _MongoURL,
		ConnectTimeout: xtime.Duration(5 * time.Second),
		QueryTimeout:   xtime.Duration(1 * time.Second),
		ExecTimeout:    xtime.Duration(2 * time.Second),
	})
	collection := db.Database("db").Collection("coll")

	result, err := db.DeleteOne(context.Background(), collection, bson.M{"Id": "ff38e8a1-411f-4aad-a3ba-32b16273e4a4"})
	if err != nil {
		t.Errorf("%v", err)
	}

	t.Logf("DeleteOne result %+v", result)
}

func TestCollection(t *testing.T) {
	db := NewMongoDB(&Config{
		URI:            _MongoURL,
		ConnectTimeout: xtime.Duration(5 * time.Second),
		QueryTimeout:   xtime.Duration(1 * time.Second),
		ExecTimeout:    xtime.Duration(2 * time.Second),
	})
	ctx := context.Background()

	wg := &sync.WaitGroup{}
	start := time.Now().UnixNano() / 1e6
	collection := db.Database("db").Collection("coll")

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(w *sync.WaitGroup) {
			for i := 0; i < 100; i++ {
				_, err := collection.Find(ctx, bson.D{{"key", "10001.login"}})
				if err != nil {
					t.Logf("one collection error : %s", err.Error())
				}
			}
			w.Done()
		}(wg)
	}
	wg.Wait()
	end := time.Now().UnixNano() / 1e6
	t.Logf("one collection finished time : %d\n", end-start)

	start = time.Now().UnixNano() / 1e6
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(w *sync.WaitGroup) {
			coll := db.Database("sharedata").Collection("game_1007")
			for i := 0; i < 100; i++ {
				_, err := coll.Find(ctx, bson.D{{"key", "10001.login"}})
				if err != nil {
					t.Logf("mult collection error : %s", err.Error())
				}
			}
			w.Done()
		}(wg)
	}
	wg.Wait()

	end = time.Now().UnixNano() / 1e6
	t.Logf("mult collection finished time : %d\n", end-start)

}
