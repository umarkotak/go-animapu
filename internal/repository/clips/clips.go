package clips

import (
	"context"

	firebaseHelper "github.com/umarkotak/go-animapu/internal/lib/firebase_helper"
	"github.com/umarkotak/go-animapu/internal/models"
)

var ctx = context.Background()

func GetClipsFromFirebase() map[string]models.Clip {
	firebaseDB := firebaseHelper.GetFirebaseDB()

	ref := firebaseDB.NewRef("")
	clipRef := ref.Child("clip_db")
	var clipDatas map[string]models.Clip
	clipRef.Get(ctx, &clipDatas)

	return clipDatas
}

func CreateClipsToFirebase(clip models.Clip) models.Clip {
	firebaseDB := firebaseHelper.GetFirebaseDB()

	ref := firebaseDB.NewRef("")
	clipRef := ref.Child("clip_db")
	clipDataRef := clipRef.Child(clip.ID)
	clipDataRef.Set(ctx, clip)
	return clip
}
