package handler

import (
	"fmt"
	"testing"
)

func TestSendPackageRequest(t *testing.T) {
	pkg := "google.golang.org/grpc"

	resp, _ := ParseResponse(SendPackageRequest(pkg))
	for _, art := range resp.Artifacts {
		fmt.Println(art.Go.Version)
		if art.Go.Version == "v1.67.1" || art.Go.Version == "v1.68.0" {
			fmt.Println()
		}
	}
}
