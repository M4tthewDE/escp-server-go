package db

import (
	"context"
	"log"

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

	_, err = client.Collection("ranking").NewDoc().Set(ctx, ranking)
	if err != nil {
		return err
	}

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
