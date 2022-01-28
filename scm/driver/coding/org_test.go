package coding

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jenkins-x/go-scm/scm"
	"gopkg.in/h2non/gock.v1"
)

func TestOrgList(t *testing.T) {
	defer gock.Off()

	// gock.New("https://e.coding.net/open-api").
	// 	Post("").
	// 	Reply(200).
	// 	Type("application/json").
	// 	File("testdata/user.json")

	client := NewDefaultWithToken("")
	got, res, err := client.Organizations.List(context.Background(), scm.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	org, res, err := client.Organizations.Find(context.Background(), "tw-test")

	t.Log(org)

	want := new(scm.User)
	raw, _ := ioutil.ReadFile("testdata/user.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
}
