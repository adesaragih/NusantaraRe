package services

// Perbaikan laporan DEV work owner 10-10-2026 (case hasil Register) - tanpa Oracle.

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/claimlife/backend/repository"
)

// pembacaGalat - pembaca polis yang menjawab galat basis data (bukan "tidak ditemukan").
type pembacaGalat struct{}

func (pembacaGalat) Ringkas(context.Context, string) (kontrak.PolisRingkas, error) {
	return kontrak.PolisRingkas{}, errors.New("UJI: ORA-00942 tiruan")
}

func (pembacaGalat) RingkasDariIdPega(context.Context, string) (kontrak.PolisRingkas, error) {
	return kontrak.PolisRingkas{}, errors.New("UJI: ORA-00942 tiruan")
}

// Galat membaca polis untuk penanda MAX CLAIM RECEIVED menjadi ALASAN, tidak menggagalkan layar Detail (bug 3).
func TestAmbangGagalBacaTidakMenggagalkanDetail(t *testing.T) {
	k := &KlaimLife{polis: pembacaGalat{}}
	ambang, alasan, err := k.ambangTerimaKlaim(context.Background(), []string{"NBLF-UJI-1"})
	if err != nil || ambang != "" || !strings.Contains(alasan, "tidak dihitung") {
		t.Errorf("ambang %q alasan %q galat %v", ambang, alasan, err)
	}
}

// Case ID baru berawalan CLMLF- (keputusan work owner 10-10-2026); pengenal lama CLM- tidak diurai siapa pun.
func TestAwalanKlaimCLMLF(t *testing.T) {
	if repository.AwalanKlaim != "CLMLF-" {
		t.Fatalf("AwalanKlaim = %q", repository.AwalanKlaim)
	}
	if got := repository.RakitPengenalWork(repository.AwalanKlaim, "10"); got != "CLMLF-000010" || len(got) > 32 {
		t.Errorf("pengenal %q (kolom ID VARCHAR2(32))", got)
	}
}

// Handler detail mencatat galat asli ke log dan menjawab kalimat Inggris tanpa rincian.
func TestHandlerDetailMencatatGalat(t *testing.T) {
	b, err := os.ReadFile("../handlers/klaimlife.go")
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	if !strings.Contains(s, `log.Printf("claimlife: membaca klaim %q: %v", r.PathValue("id"), err)`) ||
		!strings.Contains(s, `"Failed to load the claim."`) || strings.Contains(s, `StatusInternalServerError, "gagal membaca klaim")`) {
		t.Error("handler detail tidak mencatat galat / kalimat layar belum diganti")
	}
}
