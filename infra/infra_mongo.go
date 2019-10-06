package infra

import (
	"context"
	"fmt"

	"github.com/art-frela/blog/domain"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoPostRepo implementation of domain post repository
type MongoPostRepo struct {
	mongoURL       string
	database       string
	collectionName string
	session        *mongo.Client
	log            *logrus.Entry
}

// NewMongoPostRepo builder of MongoDB post repository implementation
func NewMongoPostRepo(databaseURL, database string, logger *logrus.Entry, countExamplePosts int, clearStorage bool) *MongoPostRepo {
	repo := &MongoPostRepo{}
	p := &domain.PostInBlog{}
	repo.mongoURL = databaseURL
	repo.collectionName = p.TableCollectionName()
	repo.database = database
	repo.log = logger.WithField("database", database)
	session, err := repo.connDB()
	if err != nil {
		repo.log.Fatalf("connect to mongoDB error, %v", err)
	}
	repo.session = session
	if clearStorage {
		c := repo.collection(repo.collectionName)
		c.Drop(context.TODO())
	}
	repo.fillExampleData(countExamplePosts)
	return repo
}

// FindByID returns one post from MongoDB,
// implement FindByID method of post repository
func (mpr *MongoPostRepo) FindByID(id string) (domain.PostInBlog, error) {
	post := domain.PostInBlog{}
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return post, err
	}

	// use filter by ID
	var filter = bson.D{{"_id", objectID}}
	err = mpr.collection(mpr.collectionName).FindOne(context.TODO(), filter, options.FindOne()).Decode(&post)
	if err != nil {
		return post, err
	}
	post.ID = post.ID.(primitive.ObjectID).Hex()
	return post, nil
}

// Find returns slice of posts from MongoDB,
// implement Find method of post repository
func (mpr *MongoPostRepo) Find(limit, offset int) ([]domain.PostInBlog, error) {
	posts := make([]domain.PostInBlog, 0, 16)
	filter := bson.D{}
	cur, err := mpr.collection(mpr.collectionName).Find(context.TODO(), filter, options.Find().SetLimit(int64(limit)).SetSkip(int64(offset)))
	if err != nil {
		return posts, err
	}
	defer cur.Close(context.Background())
	for cur.Next(nil) {
		post := domain.PostInBlog{}
		err := cur.Decode(&post)
		if err != nil {
			return posts, err
		}
		post.ID = post.ID.(primitive.ObjectID).Hex()
		posts = append(posts, post)
	}
	return posts, nil
}

// Save returns id of saved post in the MongoDB,
// implement Save method of post repository
func (mpr *MongoPostRepo) Save(p domain.PostInBlog) (string, error) {
	insertResult, err := mpr.collection(mpr.collectionName).InsertOne(context.TODO(), &p)
	if err != nil {
		return "-1", err
	}
	postID := insertResult.InsertedID.(primitive.ObjectID)
	mpr.log.Debugf("insertedID=%s", postID.Hex())
	return postID.Hex(), nil
}

// Update replace fields values of post in the MongoDB,
// implement Update method of post repository
func (mpr *MongoPostRepo) Update(p domain.PostInBlog) error {
	objectID, err := primitive.ObjectIDFromHex(p.ID.(string))
	if err != nil {
		return err
	}
	filter := bson.D{{"_id", objectID}}
	// implement update only for some fields
	update := bson.D{}
	update = append(update, bson.E{"title", p.Title})
	update = append(update, bson.E{"content", p.Content})
	p.CountOfViews++
	update = append(update, bson.E{"count_of_views", p.CountOfViews})
	update = bson.D{{"$set", update}}
	_, err = mpr.collection(mpr.collectionName).UpdateOne(context.TODO(), filter, update)
	return err
}

// connDB - connects to mongoDB and sets session propertie
func (mpr *MongoPostRepo) connDB() (*mongo.Client, error) {
	// make session
	session, err := mongo.NewClient(options.Client().SetAppName("blog:v1").ApplyURI(mpr.mongoURL))
	if err != nil {
		return nil, fmt.Errorf("Mongo.NewClient (%s) error, %v", mpr.mongoURL, err)
	}
	// Create connect
	err = session.Connect(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("Mongo.Client.Connect to (%s) error, %v", mpr.mongoURL, err)
	}
	// Check the connection
	err = session.Ping(context.TODO(), nil)
	if err != nil {
		return nil, fmt.Errorf("Mongo.Client.Ping to (%s) error, %v", mpr.mongoURL, err)
	}
	mpr.log.Debugf("Mongo.Client.Connect+Ping to (%s) success", mpr.mongoURL)
	return session, nil
}

// collection - returns new collection
func (mpr *MongoPostRepo) collection(name string) *mongo.Collection {
	return mpr.session.Database(mpr.database).Collection(name)
}

// fillExampleData fills SimplePostRepo with fake posts exactly N pieces,
// but no more 3rd
func (mpr *MongoPostRepo) fillExampleData(n int) {
	if n > 3 || n <= 0 { // simple fuse
		n = 3
	}
	//newID := uuid.Must(uuid.NewV4()).String()
	postTmpl := domain.PostInBlog{
		Title: "Example post #%d",
		Author: domain.User{
			ID:    "someUserID_%d",
			Name:  "anonimous#%d",
			EMail: "anonim@example.com",
		},
		Rubric: domain.Rubric{
			ID:    "rubricID_%d",
			Title: "Go for fun",
		},
		Content: `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore
            magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo
            consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla
            pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est
			laborum.`,
	}

	for i := 1; i <= n; i++ {
		post := domain.PostInBlog{
			Title:   fmt.Sprintf(postTmpl.Title, i),
			Author:  postTmpl.Author,
			Rubric:  postTmpl.Rubric,
			Content: postTmpl.Content,
		}
		post.Author.ID = domain.AnonimousID
		post.Author.Name = fmt.Sprintf(postTmpl.Author.Name, i)
		post.Rubric.ID = "00000000-0000-0000-00000000"
		_, err := mpr.Save(post)
		if err != nil {
			mpr.log.Errorf("for post=%+v, error %v", post, err)
		}
	}
}
