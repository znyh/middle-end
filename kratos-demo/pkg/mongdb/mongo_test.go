package mongdb

//package mongdb
//
//import (
//    "context"
//    "fmt"
//    "sync"
//    "testing"
//    "time"
//    xtime "time"
//
//    "go.mongodb.org/mongo-driver/bson"
//    "go.mongodb.org/mongo-driver/mongo"
//    "go.mongodb.org/mongo-driver/mongo/options"
//)
//
//func TestConnect(t *testing.T) {
//    client := GetClient()
//
//    collection := client.Database("gamedata").Collection("game_1001")
//
//    res, err := collection.InsertOne(context.Background(), bson.M{"key": "10001", "value": "{\"login_status\":true}"})
//    if err != nil {
//        panic(err)
//    }
//
//    id := res.InsertedID
//
//    t.Logf("insert id %v", id)
//}
//
//func GetClient() *mongo.Client {
//    client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://gamer:EZoqsFxF8eDfW6xg@172.13.3.141:27117"))
//    if err != nil {
//        panic(err)
//    }
//    ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
//
//    defer cancel()
//    err = client.Connect(ctx)
//    if err != nil {
//        panic(err)
//    }
//    return client
//}
//
//func TestDelete(t *testing.T) {
//    db := NewMongoDB(&Config{
//        URI:            "mongodb://gamer:EZoqsFxF8eDfW6xg@172.13.3.141:27117",
//        ConnectTimeout: xtime.Duration(5 * time.Second),
//        QueryTimeout:   xtime.Duration(1 * time.Second),
//        ExecTimeout:    xtime.Duration(2 * time.Second),
//    })
//    collection := db.Database("pay").Collection("Order")
//
//    result, err := db.DeleteOne(context.Background(), collection, bson.M{"Id": "ff38e8a1-411f-4aad-a3ba-32b16273e4a4"})
//    if err != nil {
//        t.Errorf("%v", err)
//    }
//
//    t.Logf("DeleteOne result %+v", result)
//}
//
//func TestUpdate(t *testing.T) {
//    db := NewMongoDB(&Config{
//        URI:            "mongodb://gamer:EZoqsFxF8eDfW6xg@172.13.3.141:27117",
//        ConnectTimeout: xtime.Duration(5 * time.Second),
//        QueryTimeout:   xtime.Duration(1 * time.Second),
//        ExecTimeout:    xtime.Duration(2 * time.Second),
//    })
//
//    collection := db.Database("sharedata").Collection("game_1007")
//
//    // update failed, insert it.
//    opts := new(options.UpdateOptions)
//    opts.SetUpsert(true)
//
//    result, err := db.UpdateOne(context.Background(), collection,
//        bson.M{"key": "10001.login"},
//        bson.D{{"$set", bson.M{"value": "1234567"}}}, opts)
//
//    if err != nil {
//        t.Errorf("%v", err)
//    }
//
//    t.Logf("updateone result %+v", result)
//}
//
//func TestFind(t *testing.T) {
//    db := NewMongoDB(&Config{
//        URI:            "mongodb://gamer:EZoqsFxF8eDfW6xg@172.13.3.141:27117",
//        ConnectTimeout: xtime.Duration(5 * time.Second),
//        QueryTimeout:   xtime.Duration(1 * time.Second),
//        ExecTimeout:    xtime.Duration(2 * time.Second),
//    })
//
//    collection := db.Database("sharedata").Collection("game_1007")
//
//    ctx := context.Background()
//    cur, err := db.Find(ctx, collection, bson.D{{"key", "10001.login"}})
//    if err != nil {
//        panic(err)
//    }
//
//    for cur.Next(ctx) {
//        doc := cur.Current
//
//        //方式一
//        elem, err := doc.Elements()
//        if err != nil {
//            panic(err)
//        }
//
//        for _, v := range elem {
//            t.Logf("key: %s  value: %s", v.Key(), v.Value().String())
//        }
//
//        //方式二
//        t.Logf("==> %s %s %s",
//            doc.Lookup("key").String(),
//            doc.Lookup("value").String(),
//            doc.Lookup("notExist").String())
//    }
//}
//
//func TestCollection(t *testing.T) {
//    db := NewMongoDB(&Config{
//        URI:            "mongodb://gamer:EZoqsFxF8eDfW6xg@172.13.3.141:27117",
//        ConnectTimeout: xtime.Duration(5 * time.Second),
//        QueryTimeout:   xtime.Duration(1 * time.Second),
//        ExecTimeout:    xtime.Duration(2 * time.Second),
//    })
//    ctx := context.Background()
//
//    wg := &sync.WaitGroup{}
//    start := time.Now().UnixNano() / 1e6
//    collection := db.Database("sharedata").Collection("game_1007")
//    for i := 0; i < 50; i++ {
//        wg.Add(1)
//        go func(w *sync.WaitGroup) {
//            for i := 0; i < 100; i++ {
//                _, err := collection.Find(ctx, bson.D{{"key", "10001.login"}})
//                if err != nil {
//                    t.Logf("one collection error : %s", err.Error())
//                }
//            }
//            w.Done()
//        }(wg)
//    }
//    wg.Wait()
//    end := time.Now().UnixNano() / 1e6
//    fmt.Printf("one collection finished time : %d\n", end-start)
//
//    start = time.Now().UnixNano() / 1e6
//    for i := 0; i < 50; i++ {
//        wg.Add(1)
//        go func(w *sync.WaitGroup) {
//            coll := db.Database("sharedata").Collection("game_1007")
//            for i := 0; i < 100; i++ {
//                _, err := coll.Find(ctx, bson.D{{"key", "10001.login"}})
//                if err != nil {
//                    t.Logf("mult collection error : %s", err.Error())
//                }
//            }
//            w.Done()
//        }(wg)
//    }
//    wg.Wait()
//
//    end = time.Now().UnixNano() / 1e6
//    fmt.Printf("mult collection finished time : %d\n", end-start)
//
//}
