package bootstrap

import "testing"

func TestShouldReportBootstrapReady(t *testing.T) {
	if !IsReady() {
		t.Fatalf("expected bootstrap to be ready")
	}
}
