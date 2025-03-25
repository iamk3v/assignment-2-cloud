package database

import (
	"assignment-2/utils"
	"context"

	"cloud.google.com/go/firestore"
)

// CreateWebhook stores a new webhook in Firestore
func CreateWebhook(ctx context.Context, client *firestore.Client, hook utils.Webhook) (string, error) {
	docRef, _, err := client.Collection("webhooks").Add(ctx, hook)
	if err != nil {
		return "", err
	}
	return docRef.ID, nil
}

// GetWebhook retrieves a single webhook by ID
func GetWebhook(ctx context.Context, client *firestore.Client, id string) (*utils.Webhook, error) {
	docSnap, err := client.Collection("webhooks").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var hook utils.Webhook
	if err := docSnap.DataTo(&hook); err != nil {
		return nil, err
	}
	hook.ID = docSnap.Ref.ID
	return &hook, nil
}

// GetAllWebhooks retrieves all webhooks
func GetAllWebhooks(ctx context.Context, client *firestore.Client) ([]utils.Webhook, error) {
	var hooks []utils.Webhook
	iter := client.Collection("webhooks").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		var hook utils.Webhook
		if err := doc.DataTo(&hook); err != nil {
			continue
		}
		hook.ID = doc.Ref.ID
		hooks = append(hooks, hook)
	}
	return hooks, nil
}

// DeleteWebhook deletes a single webhook
func DeleteWebhook(ctx context.Context, client *firestore.Client, id string) error {
	_, err := client.Collection("webhooks").Doc(id).Delete(ctx)
	return err
}
