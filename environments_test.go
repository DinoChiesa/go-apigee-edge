package apigee

import (
	"testing"
)

const ()

func TestEnvList(t *testing.T) {
	client := NewClientForTesting(t)
	namelist, resp, e := client.Environments.List()
	if e != nil {
		t.Errorf("while listing environments, error:\n%#v\n", e)
		return
	}
	defer resp.Body.Close()
	if len(namelist) <= 0 {
		t.Errorf("no environments found")
		return
	}
}
