package config

import (
	"cloud.google.com/go/firestore"
	"context"
)

var Ctx context.Context
var Client *firestore.Client
