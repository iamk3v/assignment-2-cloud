package database

import (
	"assignment-2/config"
	"assignment-2/utils"
	"cloud.google.com/go/firestore"
)

/*
CreateWebhook creates and stores a new webhook in the notification database
*/
func CreateWebhook(hook utils.Webhook) (string, error) {
	docRef, _, err := Client.Collection(config.NOTIFICATION_COLLECTION).Add(Ctx, hook)
	if err != nil {
		return "", err
	}
	// Update the document to include its generated ID.
	updateData := map[string]interface{}{
		"id": docRef.ID,
	}
	_, err = docRef.Set(Ctx, updateData, firestore.MergeAll)
	if err != nil {
		return "", err
	}
	return docRef.ID, nil
}

/*
GetWebhook retrieves a single webhook by ID from the notifications database
*/
func GetWebhook(id string) (*utils.Webhook, error) {
	docSnap, err := Client.Collection(config.NOTIFICATION_COLLECTION).Doc(id).Get(Ctx)
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

/*
GetAllWebhooks retrieves all webhooks from the notifications database
*/
func GetAllWebhooks() ([]utils.Webhook, error) {
	var hooks []utils.Webhook
	iter := Client.Collection(config.NOTIFICATION_COLLECTION).Documents(Ctx)
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

/*
DeleteWebhook deletes a single webhook from the notification database
*/
func DeleteWebhook(id string) error {
	_, err := Client.Collection(config.NOTIFICATION_COLLECTION).Doc(id).Delete(Ctx)
	return err
}

/*
UpdateWebhook updates an existing webhook document in notification database by merging the provided data.
*/
func UpdateWebhook(id string, updatedData map[string]interface{}) error {
	// The MergeAll option will update only the fields provided in updatedData.
	_, err := Client.Collection(config.NOTIFICATION_COLLECTION).Doc(id).Set(Ctx, updatedData, firestore.MergeAll)
	return err
}
