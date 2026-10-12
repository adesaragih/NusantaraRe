package services_test

import (
	"os"
	"strings"
	"testing"
)

// 12-10-2026 - Reject Outstanding hanya di inbox Admin: peran Admin (403) DAN tahap Outstanding Claim (409).
func TestTolakBergerbangTahapOutstanding(t *testing.T) {
	isi, err := os.ReadFile("tolak.go")
	if err != nil {
		t.Fatal(err)
	}
	s := string(isi)
	tolak := s[strings.Index(s, "func (st *Status) Tolak("):]
	for _, mau := range []string{"inti.WajibPeran(pelaku, PeranRejectOutstanding)", "tahap, err := tahapKasus(ctx, baca, klaimID)",
		"if tahap != models.TahapOutstanding {", "ErrTolakBukanOutstanding"} {
		if !strings.Contains(tolak, mau) {
			t.Errorf("Tolak tidak lagi memuat %s", mau)
		}
	}
	h, err := os.ReadFile("../handlers/tolak.go")
	if err != nil {
		t.Fatal(err)
	}
	i := strings.Index(string(h), "case errors.Is(err, services.ErrTolakBukanOutstanding):")
	if i < 0 || !strings.Contains(string(h)[i:i+250], "http.StatusConflict") {
		t.Error("ErrTolakBukanOutstanding tidak dijawab 409")
	}
}
