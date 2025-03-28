package database

import (
	"assignment-2/utils"
)

// CreateWebhook stores a new webhook in Firestore
func CreateWebhook(hook utils.Webhook) (string, error) {
	docRef, _, err := Client.Collection("webhooks").Add(Ctx, hook)
	if err != nil {
		return "", err
	}
	return docRef.ID, nil
}

// GetWebhook retrieves a single webhook by ID
func GetWebhook(id string) (*utils.Webhook, error) {
	docSnap, err := Client.Collection("webhooks").Doc(id).Get(Ctx)
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
func GetAllWebhooks() ([]utils.Webhook, error) {
	var hooks []utils.Webhook
	iter := Client.Collection("webhooks").Documents(Ctx)
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
func DeleteWebhook(id string) error {
	_, err := Client.Collection("webhooks").Doc(id).Delete(Ctx)
	return err
}
