package models

import "time"

type ResumeTokenState struct {
	Db                     string    `json:"db" bson:"db"`
	Coll                   string    `json:"coll" bson:"coll"`
	LastMongoResumeToken   string    `json:"last_mongo_resume_token_raw" bson:"last_mongo_resume_token"`
	LastMongoProcessedTime time.Time `json:"last_mongo_processed_time" bson:"last_mongo_processed_time"`
}

type MongoResumeToken struct {
	Token string `bson:"_data"`
}
