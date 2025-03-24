package config

import (
	"cloud.google.com/go/firestore"
	"context"
	"time"
)

var Ctx context.Context
var Client *firestore.Client

var Starttime time.Time
