package mongdb

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/pkg/log"
	"github.com/go-kratos/kratos/pkg/net/netutil/breaker"
	"github.com/go-kratos/kratos/pkg/net/trace"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	_family          = "mongo_client"
	_slowLogDuration = time.Millisecond * 250
)

// DB database.
type DB struct {
	*mongo.Client
	opts    *options.ClientOptions
	breaker breaker.Breaker
	conf    *Config
}

func Open(c *Config) (*DB, error) {
	opts := options.Client().ApplyURI(c.URI)
	opts.SetConnectTimeout(time.Duration(c.ConnectTimeout))

	d, err := mongo.NewClient(opts)
	if err != nil {
		return nil, err
	}

	brkGroup := breaker.NewGroup(c.Breaker)
	brk := brkGroup.Get(opts.Hosts[0])
	if err = d.Connect(context.Background()); err != nil {
		return nil, err
	}

	db := &DB{Client: d, breaker: brk, conf: c, opts: opts}
	return db, nil
}

func (db *DB) UpdateOne(c context.Context, collection *mongo.Collection, filter, update interface{},
	opts ...*options.UpdateOptions) (result *mongo.UpdateResult, err error) {

	now := time.Now()
	defer slowLog(fmt.Sprintf("UpdateOne filter(%+v) update(%+v)", filter, update), now)
	if t, ok := trace.FromContext(c); ok {
		t = t.Fork(_family, "UpdateOne")
		t.SetTag(trace.String(trace.TagAddress, db.opts.Hosts[0]))
		defer t.Finish(&err)
	}
	if err = db.breaker.Allow(); err != nil {
		_metricReqErr.Inc(db.opts.Hosts[0], db.opts.Hosts[0], "UpdateOne", "breaker")
		return
	}
	_, c, cancel := db.conf.ExecTimeout.Shrink(c)
	result, err = collection.UpdateOne(c, filter, update, opts...)
	cancel()
	db.onBreaker(&err)
	_metricReqDur.Observe(int64(time.Since(now)/time.Millisecond), db.opts.Hosts[0], db.opts.Hosts[0], "UpdateOne")
	if err != nil {
		err = errors.Wrapf(err, "UpdateOne filter(%+v) update(%+v)", filter, update)
	}
	return
}

func (db *DB) DeleteOne(c context.Context, collection *mongo.Collection, filter interface{},
	opts ...*options.DeleteOptions) (result *mongo.DeleteResult, err error) {

	now := time.Now()
	defer slowLog(fmt.Sprintf("DeleteOne filter(%+v)", filter), now)
	if t, ok := trace.FromContext(c); ok {
		t = t.Fork(_family, "DeleteOne")
		t.SetTag(trace.String(trace.TagAddress, db.opts.Hosts[0]))
		defer t.Finish(&err)
	}
	if err = db.breaker.Allow(); err != nil {
		_metricReqErr.Inc(db.opts.Hosts[0], db.opts.Hosts[0], "DeleteOne", "breaker")
		return
	}
	_, c, cancel := db.conf.ExecTimeout.Shrink(c)
	result, err = collection.DeleteOne(c, filter, opts...)
	cancel()
	db.onBreaker(&err)
	_metricReqDur.Observe(int64(time.Since(now)/time.Millisecond), db.opts.Hosts[0], db.opts.Hosts[0], "DeleteOne")
	if err != nil {
		err = errors.Wrapf(err, "DeleteOne filter(%+v)", filter)
	}
	return
}

func (db *DB) UpdateMany(c context.Context, collection *mongo.Collection, filter, update interface{},
	opts ...*options.UpdateOptions) (result *mongo.UpdateResult, err error) {

	now := time.Now()
	defer slowLog(fmt.Sprintf("UpdateMany filter(%+v) update(%+v)", filter, update), now)
	if t, ok := trace.FromContext(c); ok {
		t = t.Fork(_family, "UpdateMany")
		t.SetTag(trace.String(trace.TagAddress, db.opts.Hosts[0]))
		defer t.Finish(&err)
	}
	if err = db.breaker.Allow(); err != nil {
		_metricReqErr.Inc(db.opts.Hosts[0], db.opts.Hosts[0], "UpdateMany", "breaker")
		return
	}
	_, c, cancel := db.conf.ExecTimeout.Shrink(c)
	result, err = collection.UpdateMany(c, filter, update, opts...)
	cancel()
	db.onBreaker(&err)
	_metricReqDur.Observe(int64(time.Since(now)/time.Millisecond), db.opts.Hosts[0], db.opts.Hosts[0], "UpdateMany")
	if err != nil {
		err = errors.Wrapf(err, "UpdateMany filter(%+v) update(%+v)", filter, update)
	}

	return
}

func (db *DB) Find(c context.Context, collection *mongo.Collection, filter interface{},
	opts ...*options.FindOptions) (cursor *mongo.Cursor, err error) {

	now := time.Now()
	defer slowLog(fmt.Sprintf("Find filter(%+v)", filter), now)
	if t, ok := trace.FromContext(c); ok {
		t = t.Fork(_family, "Find")
		t.SetTag(trace.String(trace.TagAddress, db.opts.Hosts[0]))
		defer t.Finish(&err)
	}
	if err = db.breaker.Allow(); err != nil {
		_metricReqErr.Inc(db.opts.Hosts[0], db.opts.Hosts[0], "Find", "breaker")
		return
	}
	_, c, cancel := db.conf.QueryTimeout.Shrink(c)
	cursor, err = collection.Find(c, filter, opts...)
	cancel()
	db.onBreaker(&err)
	_metricReqDur.Observe(int64(time.Since(now)/time.Millisecond), db.opts.Hosts[0], db.opts.Hosts[0], "Find")
	if err != nil {
		err = errors.Wrapf(err, "Find filter(%+v)", filter)
	}

	return
}

func (db *DB) FindOne(c context.Context, collection *mongo.Collection, filter interface{},
	opts ...*options.FindOneOptions) (result *mongo.SingleResult, err error) {

	now := time.Now()
	defer slowLog(fmt.Sprintf("FindOne filter(%+v)", filter), now)
	if t, ok := trace.FromContext(c); ok {
		t = t.Fork(_family, "FindOne")
		t.SetTag(trace.String(trace.TagAddress, db.opts.Hosts[0]))
		defer t.Finish(&err)
	}
	if err = db.breaker.Allow(); err != nil {
		_metricReqErr.Inc(db.opts.Hosts[0], db.opts.Hosts[0], "FindOne", "breaker")
		return
	}
	_, c, cancel := db.conf.QueryTimeout.Shrink(c)
	result = collection.FindOne(c, filter, opts...)
	cancel()
	db.onBreaker(&err)
	_metricReqDur.Observe(int64(time.Since(now)/time.Millisecond), db.opts.Hosts[0], db.opts.Hosts[0], "FindOne")
	if err != nil {
		err = errors.Wrapf(err, "FindOne filter(%+v)", filter)
	}
	return
}

func (db *DB) InsertOne(ctx context.Context, collection *mongo.Collection, document interface{}, opts ...*options.InsertOneOptions) (result *mongo.InsertOneResult, err error) {
	now := time.Now()
	defer slowLog(fmt.Sprintf("InsertOne document(%+v)", document), now)
	if t, ok := trace.FromContext(ctx); ok {
		t = t.Fork(_family, "InsertOne")
		t.SetTag(trace.String(trace.TagAddress, db.opts.Hosts[0]))
		defer t.Finish(&err)
	}
	if err = db.breaker.Allow(); err != nil {
		_metricReqErr.Inc(db.opts.Hosts[0], db.opts.Hosts[0], "InsertOne", "breaker")
		return
	}
	_, ctx, cancel := db.conf.QueryTimeout.Shrink(ctx)
	defer cancel()

	result, err = collection.InsertOne(ctx, document, opts...)

	db.onBreaker(&err)
	_metricReqDur.Observe(int64(time.Since(now)/time.Millisecond), db.opts.Hosts[0], db.opts.Hosts[0], "InsertOne")
	if err != nil {
		err = errors.Wrapf(err, "InsertOne filter(%+v)", document)
	}
	return
}

func (db *DB) InsertMany(ctx context.Context, collection *mongo.Collection, documents []interface{}, opts ...*options.InsertManyOptions) (result *mongo.InsertManyResult, err error) {
	now := time.Now()
	defer slowLog(fmt.Sprintf("InsertMany documents(%+v)", documents), now)
	if t, ok := trace.FromContext(ctx); ok {
		t = t.Fork(_family, "InsertMany")
		t.SetTag(trace.String(trace.TagAddress, db.opts.Hosts[0]))
		defer t.Finish(&err)
	}
	if err = db.breaker.Allow(); err != nil {
		_metricReqErr.Inc(db.opts.Hosts[0], db.opts.Hosts[0], "InsertMany", "breaker")
		return
	}
	_, ctx, cancel := db.conf.QueryTimeout.Shrink(ctx)
	defer cancel()

	result, err = collection.InsertMany(ctx, documents, opts...)

	db.onBreaker(&err)
	_metricReqDur.Observe(int64(time.Since(now)/time.Millisecond), db.opts.Hosts[0], db.opts.Hosts[0], "InsertMany")
	if err != nil {
		err = errors.Wrapf(err, "InsertMany filter(%+v)", documents)
	}
	return
}

func slowLog(statement string, now time.Time) {
	du := time.Since(now)
	if du > _slowLogDuration {
		log.Warn("%s slow log statement: %s time: %v", _family, statement, du)
	}
}

func (db *DB) onBreaker(err *error) {
	if err != nil && *err != nil {
		db.breaker.MarkFailed()
	} else {
		db.breaker.MarkSuccess()
	}
}
