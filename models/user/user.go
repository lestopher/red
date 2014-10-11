package user

import (
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	password             string
	passwordConfirmation string

	ID                  bson.ObjectId `bson:"_id, omitempty"`
	Username            string        `bson:"username"`
	Email               string        `bson:"email"`
	EncryptedPassword   string        `bson:"encrypted_password"`
	FirstLogin          bool          `bson:"first_login"`
	LegacyID            int           `bson:"legacy_id"`
	ResetPasswordToken  string        `bson:"reset_password_token"`
	ResetPasswordSentAt string        `bson:"reset_password_sent_at"`
	SignInCount         int           `bson:"sign_in_count"`
	CurrentSignInAt     time.Time     `bson:"current_sign_in_at"`
	LastSignInAt        time.Time     `bson:"last_sign_in_at"`
	CurrentSignInIP     string        `bson:"current_sign_in_ip"`
	LastSignInIP        string        `bson:"last_sign_in_ip"`
	PasswordSalt        string        `bson:"password_salt"`
	InvitationKey       string        `bson:"invitation_key"`
	AuthenticationToken string        `bson:"authentication_token"`
}

func EnsureIndicies(c *mgo.Collection) error {
	usernameIndex := mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(usernameIndex)
	if err != nil {
		return err
	}

	emailIndex := mgo.Index{
		Key:        []string{"email"},
		Unique:     true,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	}
	err = c.EnsureIndex(emailIndex)
	if err != nil {
		return err
	}

	authenticationTokenIndex := mgo.Index{
		Key:        []string{"authentication_token"},
		Unique:     true,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	}
	err = c.EnsureIndex(authenticationTokenIndex)

	if err != nil {
		return err
	}

	return nil
}
