package sharings

import (
	"fmt"
	"testing"

	"github.com/cozy/cozy-stack/client/auth"
	"github.com/cozy/cozy-stack/pkg/consts"
	"github.com/cozy/cozy-stack/pkg/couchdb"
	"github.com/cozy/cozy-stack/pkg/permissions"
	"github.com/cozy/cozy-stack/web/jsonapi"
	"github.com/stretchr/testify/assert"
)

var rec = &Recipient{
	URL:   "",
	Email: "",
}

var recStatus = &RecipientStatus{
	RefRecipient: jsonapi.ResourceIdentifier{
		Type: consts.Recipients,
	},
	Client: &auth.Client{
		ClientID:     "",
		RedirectURIs: []string{},
	},
	recipient: rec,
}

var mailValues = &mailTemplateValues{}

var sharingTest = &Sharing{
	SharingType:      consts.OneShotSharing,
	RecipientsStatus: []*RecipientStatus{recStatus},
	SharingID:        "sparta-id",
	Permissions:      permissions.Set{},
}

var instanceScheme = "http"

func TestGenerateMailMessageWhenRecipientHasNoEmail(t *testing.T) {
	msg, err := generateMailMessage(sharingTest, rec, mailValues)
	assert.Equal(t, ErrRecipientHasNoEmail, err)
	assert.Nil(t, msg)
}

func TestGenerateMailMessageSuccess(t *testing.T) {
	rec.Email = "this@mail.com"
	_, err := generateMailMessage(sharingTest, rec, mailValues)
	assert.NoError(t, err)
}

func TestGenerateOAuthQueryStringWhenThereIsNoOAuthClient(t *testing.T) {
	// Without client id.
	recStatus.Client.RedirectURIs = []string{"redirect.me.to.heaven"}
	oauthQueryString, err := generateOAuthQueryString(sharingTest, recStatus,
		instanceScheme)
	assert.Equal(t, ErrNoOAuthClient, err)
	assert.Equal(t, oauthQueryString, "")

	// Without redirect uri.
	recStatus.Client.ClientID = "sparta"
	recStatus.Client.RedirectURIs = []string{}
	oauthQueryString, err = generateOAuthQueryString(sharingTest, recStatus,
		instanceScheme)
	assert.Equal(t, ErrNoOAuthClient, err)
	assert.Equal(t, oauthQueryString, "")

}

func TestGenerateOAuthQueryStringWhenRecipientHasNoURL(t *testing.T) {
	recStatus.Client.RedirectURIs = []string{"redirect.me.to.sparta"}

	oauthQueryString, err := generateOAuthQueryString(sharingTest, recStatus,
		instanceScheme)
	assert.Equal(t, ErrRecipientHasNoURL, err)
	assert.Equal(t, "", oauthQueryString)
}

func TestGenerateOAuthQueryStringSuccess(t *testing.T) {
	// First test: no scheme in the url.
	rec.URL = "this.is.url"
	expectedStr := "http://this.is.url/sharings/request?client_id=sparta&redirect_uri=redirect.me.to.sparta&response_type=code&scope=&sharing_type=one-shot&state=sparta-id"

	oAuthQueryString, err := generateOAuthQueryString(sharingTest, recStatus,
		instanceScheme)
	assert.NoError(t, err)
	assert.Equal(t, expectedStr, oAuthQueryString)

	// Second test: "http" scheme in the url.
	rec.URL = "http://this.is.url"
	oAuthQueryString, err = generateOAuthQueryString(sharingTest, recStatus,
		instanceScheme)
	assert.NoError(t, err)
	assert.Equal(t, expectedStr, oAuthQueryString)

	// Third test: "https" scheme in the url.
	rec.URL = "https://this.is.url"
	expectedStr = "https://this.is.url/sharings/request?client_id=sparta&redirect_uri=redirect.me.to.sparta&response_type=code&scope=&sharing_type=one-shot&state=sparta-id"
	oAuthQueryString, err = generateOAuthQueryString(sharingTest, recStatus,
		instanceScheme)
	assert.NoError(t, err)
	assert.Equal(t, expectedStr, oAuthQueryString)
}

func TestSendSharingMails(t *testing.T) {
	// We provoke the error that occurrs when a recipient has no URL or no
	// OAuth client by creating an incomplete recipient document.
	rec.URL = ""
	// Add the recipient in the database.
	err := couchdb.CreateDoc(in, rec)
	if err != nil {
		fmt.Printf("%v\n", err)
		t.Fail()
	}
	defer couchdb.DeleteDoc(in, rec)
	// Set the id to the id generated by Couch.
	recStatus.RefRecipient.ID = rec.RID

	err = SendSharingMails(in, sharingTest)
	assert.Error(t, err)

	// The other scenario is when the recipient has no email set.
	rec.URL = "this.is.url"
	rec.Email = ""
	err = couchdb.UpdateDoc(in, rec)
	assert.NoError(t, err)

	err = SendSharingMails(in, sharingTest)
	assert.Error(t, err)
}
