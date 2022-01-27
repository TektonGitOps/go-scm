// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package coding

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/jenkins-x/go-scm/scm"
	"gopkg.in/h2non/gock.v1"

	"github.com/google/go-cmp/cmp"
)

func TestUserFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://e.coding.net/open-api").
		Post("").
		Reply(200).
		Type("application/json").
		File("testdata/user.json")

	client := NewDefaultWithToken("")
	got, res, err := client.Users.Find(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.User)
	raw, _ := ioutil.ReadFile("testdata/user.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
}
