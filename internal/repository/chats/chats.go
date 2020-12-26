package chats

import (
	"context"
	"strconv"

	"github.com/umarkotak/go-animapu/internal/models"
	firebaseHelper "github.com/umarkotak/go-animapu/internal/pkg/firebase_helper"
)

var ctx = context.Background()

// SetChatMessageToFirebase save message to firebase
func SetChatMessageToFirebase(message models.Message) models.Message {
	firebaseDB := firebaseHelper.GetFirebaseDB()

	ref := firebaseDB.NewRef("")
	messagesRef := ref.Child("chat_message_db")
	messageDataRef := messagesRef.Child(strconv.Itoa(int(message.TimestampUnix)))
	messageDataRef.Set(ctx, &message)

	return message
}
