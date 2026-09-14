package core

import "testing"

func TestIsChatGPTWindowAcceptsChatGPTTitle(t *testing.T) {
	cases := []string{
		"ChatGPT - Microsoft Edge",
		"ALO Music - ChatGPT",
		"ChatGPT",
	}
	for _, title := range cases {
		if !IsChatGPTWindow(title) {
			t.Fatalf("expected %q to match ChatGPT", title)
		}
	}
}

func TestIsChatGPTWindowRejectsOtherTitles(t *testing.T) {
	cases := []string{
		"GitHub Desktop",
		"YouTube - Microsoft Edge",
		"PowerShell",
	}
	for _, title := range cases {
		if IsChatGPTWindow(title) {
			t.Fatalf("did not expect %q to match ChatGPT", title)
		}
	}
}
