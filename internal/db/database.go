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

func (dbHandler DatabaseHandler) GetCountries() ([]CountryDto, error) {
	ctx := context.Background()

	client, err := dbHandler.app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	defer client.Close()

	countries := make([]CountryDto, 0)

	docs, err := client.Collection("countries").Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var country Country
	for _, doc := range docs {
		err = doc.DataTo(&country)
		if err != nil {
			return nil, err
		}

		countries = append(countries, CountryDto(country))
	}

	return countries, nil
}

func (dbHandler DatabaseHandler) SaveResult(result ResultDto) error {
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

func (dbHandler DatabaseHandler) SaveRanking(ranking RankingDto) error {
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

func (dbHandler DatabaseHandler) GetRanking(user string) (*RankingDto, error) {
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

		return &RankingDto{
			user,
			countries,
		}, nil
	}

	var ranking RankingDto
	err = docs[0].DataTo(&ranking)
	if err != nil {
		return nil, err
	}

	return &ranking, nil
}

func (dbHandler DatabaseHandler) SetLock() error {
	return nil
}

// only used for mapping with database.
type Country struct {
	Name string
	Flag string
}

// use for transfer between database and application.
type CountryDto struct {
	Name string
	Flag string
}

// only used for mapping with database.
type Result struct {
	Result []Country
}

// use for transfer between database and application.
type ResultDto struct {
	Result []CountryDto
}

// only used for mapping with database.
type Ranking struct {
	Name    string
	Ranking []Country
}

// use for transfer between database and application.
type RankingDto struct {
	Name    string
	Ranking []CountryDto
}
