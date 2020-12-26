package clips

import (
	"context"

	"github.com/umarkotak/go-animapu/internal/models"
	firebaseHelper "github.com/umarkotak/go-animapu/internal/pkg/firebase_helper"
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
