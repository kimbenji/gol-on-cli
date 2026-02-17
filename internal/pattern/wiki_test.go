package pattern

import "testing"

func TestShouldValidateConwayLifeWikiHTTPSURL(t *testing.T) {
	if !ValidateWikiURL("https://conwaylife.com/wiki/Glider") {
		t.Fatalf("expected https://conwaylife.com/wiki/... to be valid")
	}
}

func TestShouldRejectNonHTTPSWikiURL(t *testing.T) {
	if ValidateWikiURL("http://conwaylife.com/wiki/Glider") {
		t.Fatalf("expected non-https URL to be invalid")
	}
}

func TestShouldRejectNonConwayLifeHostURL(t *testing.T) {
	if ValidateWikiURL("https://example.com/wiki/Glider") {
		t.Fatalf("expected non-conwaylife host to be invalid")
	}
}

func TestShouldRejectNonWikiPathURL(t *testing.T) {
	if ValidateWikiURL("https://conwaylife.com/forums/Glider") {
		t.Fatalf("expected non-/wiki/ path to be invalid")
	}
}

func TestShouldPreferRLEWhenAllFormatsExist(t *testing.T) {
	content := `
#Life 1.06
0 0

!Name: glider
.O.
..O
OOO

x = 3, y = 3
bo$2bo$3o!
`

	format, _, err := SelectPreferredPattern(content)
	if err != nil {
		t.Fatalf("expected format selection to succeed, got error: %v", err)
	}
	if format != FormatRLE {
		t.Fatalf("expected RLE priority, got %s", format)
	}
}

func TestShouldPreferPlainTextWhenRLEIsMissing(t *testing.T) {
	content := `
#Life 1.06
0 0

!Name: glider
.O.
..O
OOO
`

	format, _, err := SelectPreferredPattern(content)
	if err != nil {
		t.Fatalf("expected format selection to succeed, got error: %v", err)
	}
	if format != FormatPlainText {
		t.Fatalf("expected PlainText priority, got %s", format)
	}
}

func TestShouldSelectLife106WhenOnlyLife106Exists(t *testing.T) {
	content := `
#Life 1.06
0 0
1 1
`

	format, _, err := SelectPreferredPattern(content)
	if err != nil {
		t.Fatalf("expected format selection to succeed, got error: %v", err)
	}
	if format != FormatLife106 {
		t.Fatalf("expected Life1.06 selection, got %s", format)
	}
}

func TestShouldReturnErrorForRLEOverflowRunLength(t *testing.T) {
	_, err := ParseToBoard(FormatRLE, "x = 1, y = 1\n9223372036854775808o!", 3, 3)
	if err == nil {
		t.Fatalf("expected overflow run-length to return error")
	}
}
