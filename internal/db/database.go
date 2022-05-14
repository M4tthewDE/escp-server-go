package db

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

type DatabaseHandler struct {
	app *firebase.App
}

func NewDatabaseHandler() DatabaseHandler {
	ctx := context.Background()

	options := option.WithCredentialsFile("hello-world-338921-b3875418957d.json")
	conf := &firebase.Config{ProjectID: "hello-world-338921"}

	app, err := firebase.NewApp(ctx, conf, options)
	if err != nil {
		log.Fatalln(err)
	}

	return DatabaseHandler{
		app: app,
	}
}

func (dbHandler DatabaseHandler) GetCountries() (*Countries, error) {
	ctx := context.Background()

	client, err := dbHandler.app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	defer client.Close()

	docs, err := client.Collection("countries").Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var countries Countries
	for _, doc := range docs {
		err = doc.DataTo(&countries)
		if err != nil {
			return nil, err
		}
	}

	return &countries, nil
}

func (dbHandler DatabaseHandler) SaveResult(result Result) error {
	ctx := context.Background()

	client, err := dbHandler.app.Firestore(ctx)
	if err != nil {
		return err
	}

	defer client.Close()

	_, err = client.Collection("result").NewDoc().Set(ctx, result)
	if err != nil {
		return err
	}

	return nil
}

func (dbHandler DatabaseHandler) SaveRanking(ranking Ranking) error {
	ctx := context.Background()

	client, err := dbHandler.app.Firestore(ctx)
	if err != nil {
		return err
	}

	defer client.Close()

	docs, err := client.Collection("rankings").Where("Name", "==", ranking.Name).Documents(ctx).GetAll()
	if err != nil {
		return nil
	}

	if len(docs) == 0 {
		_, err = client.Collection("rankings").NewDoc().Set(ctx, ranking)
		if err != nil {
			return err
		}
		return nil
	}

	_, err = docs[0].Ref.Set(ctx, map[string]interface{}{
		"Ranking": ranking.Ranking,
	}, firestore.MergeAll)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (dbHandler DatabaseHandler) GetRanking(user string) (*Ranking, error) {
	ctx := context.Background()

	client, err := dbHandler.app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	defer client.Close()

	docs, err := client.Collection("rankings").Where("Name", "==", user).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	if len(docs) == 0 {
		countries, err := dbHandler.GetCountries()
		if err != nil {
			return nil, err
		}

		return &Ranking{
			user,
			countries.Countries,
		}, nil
	}

	var ranking Ranking
	err = docs[0].DataTo(&ranking)
	if err != nil {
		return nil, err
	}

	return &ranking, nil
}

func (dbHandler DatabaseHandler) GetLock() (bool, error) {
	ctx := context.Background()

	client, err := dbHandler.app.Firestore(ctx)
	if err != nil {
		return false, err
	}

	defer client.Close()

	docs, err := client.Collection("lock").Documents(ctx).GetAll()
	if err != nil {
		return false, err
	}

	if len(docs) == 0 {
		return false, err
	}

	var lock Lock
	err = docs[0].DataTo(&lock)
	if err != nil {
		return false, err
	}

	return lock.Lock, nil
}

func (dbHandler DatabaseHandler) GetDone() (bool, error) {
	ctx := context.Background()

	client, err := dbHandler.app.Firestore(ctx)
	if err != nil {
		return false, err
	}

	defer client.Close()

	docs, err := client.Collection("done").Documents(ctx).GetAll()
	if err != nil {
		return false, err
	}

	if len(docs) == 0 {
		return false, err
	}

	var done Done
	err = docs[0].DataTo(&done)
	if err != nil {
		return false, err
	}

	return done.Done, nil
}

type Countries struct {
	Countries []Country
}

type Country struct {
	Name string
	Flag string
}

// only used for mapping with database.
type Result struct {
	Result []Country
}

type Ranking struct {
	Name    string
	Ranking []Country
}

type Lock struct {
	Lock bool
}

type Done struct {
	Done bool
}
