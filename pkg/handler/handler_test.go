package handler

import (
	"fmt"
	"testing"
)

func TestSendPackageRequest(t *testing.T) {
	pkg := "google.golang.org/genproto/googleapis/api"

	resp, _ := ParseResponse(SendPackageRequest(pkg))
	for _, art := range resp.Artifacts {
		fmt.Println(art.Go.Version, art.State.RequestTime)
		if art.Go.Version == "v0.0.0-20240318140521-94a12d6c2237" || art.Go.Version == "v0.0.0-20240123012728-ef4313101c80" {
			fmt.Println()
		}
	}
}
