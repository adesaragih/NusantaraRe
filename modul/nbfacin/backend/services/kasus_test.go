package services

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// kasusTiruan - PenyimpanKasus tiruan berkeadaan: SimpanGeneral menulis ke peta,
// BacaKasus membacanya kembali (bentuk teks Pega apa adanya).
type kasusTiruan struct {
	ada   map[string]*models.Kasus
	galat error
}

func (k kasusTiruan) BacaKasus(_ context.Context, id string) (models.Kasus, error) {
	if k.galat != nil {
		return models.Kasus{}, k.galat
	}
	c, ok := k.ada[id]
	if !ok {
		return models.Kasus{}, repository.ErrKasusTidakAda
	}
	return *c, nil
}

func (k kasusTiruan) SimpanGeneral(_ context.Context, _ *db.Tx, id string, g models.General) error {
	c, ok := k.ada[id]
	if !ok {
		return repository.ErrKasusTidakAda
	}
	g.SourceOfBusiness, g.CedingCoName, g.GroupName, g.OldPolicyNumber = c.General.SourceOfBusiness,
		c.General.CedingCoName, c.General.GroupName, c.General.OldPolicyNumber
	g.Tersimpan = true
	c.General = g
	return nil
}

func tanpaTx(_ context.Context, fn func(*db.Tx) error) error { return fn(nil) }

// TestTanggalPega - uji instrumen dengan jawaban yang diketahui (butir 78.1): bentuk
// fixture -> kabel DD-MM-YYYY lewat WIB, dan sebaliknya; teks asing apa adanya (A89).
func TestTanggalPega(t *testing.T) {
	for masuk, mau := range map[string]string{
		"20240101":   "01-01-2024",
		"20240229":   "29-02-2024",
		"20230229":   "20230229", // bukan tanggal - apa adanya
		"2024-01-01": "2024-01-01",
		"":           "",
	} {
		if got := tanggalKeKabel(masuk); got != mau {
			t.Errorf("tanggalKeKabel(%q) = %q, mau %q", masuk, got, mau)
		}
	}
	for masuk, mau := range map[string]string{
		"20231231T170000.000 GMT": "01-01-2024", // 00:00 WIB
		"20240101T050000.000 GMT": "01-01-2024", // 12:00 WIB
		"20240101T165959.999 GMT": "01-01-2024", // 23:59 WIB
		"UJI teks":                "UJI teks",
	} {
		if got := waktuKeKabel(masuk); got != mau {
			t.Errorf("waktuKeKabel(%q) = %q, mau %q", masuk, got, mau)
		}
	}
	g, err := periksaGeneral(IsianGeneral{BeginDate: "01-01-2024", OfferingDate: "15-12-2023", EndDate: "31-12-2024"})
	if err != nil || g.StartDateTime != "20240101T050000.000 GMT" || g.OfferingDate != "20231215" ||
		g.EndDateTime != "20241231T050000.000 GMT" {
		t.Errorf("ke Pega: %+v (%v)", g, err)
	}
	if waktuKeKabel(g.EndDateTime) != "31-12-2024" || tanggalKeKabel(g.OfferingDate) != "15-12-2023" {
		t.Error("pulang-pergi tanggal tidak kembali ke nilai semula")
	}
}

func kasusUji() map[string]*models.Kasus {
	return map[string]*models.Kasus{"UJI-NB-1": {CaseID: "UJI-NB-1", Position: "Offer", StatusWork: "Pending-Policy",
		InsuredName: "UJI TERTANGGUNG",
		Opportunity: models.Opportunity{EstimatedClosingDate: time.Date(2026, 10, 31, 0, 0, 0, 0, time.UTC),
			BusinessProspectName: "UJI Prospek", TypeOfFacultative: "Facultative In"},
		General: models.General{SourceOfBusiness: "UJI SOB", OldPolicyNumber: "UJI-LAMA"}}}
}

// TestBacaKasus - kabel case: tanggal opportunity DD-MM-YYYY; OfferingDate kosong = hari
// ini WIB (PreDT, A91) - jam 18:30 UTC sudah tanggal berikutnya di WIB; 404 / 503.
func TestBacaKasus(t *testing.T) {
	jam := func() time.Time { return time.Date(2026, 10, 2, 18, 30, 0, 0, time.UTC) }
	svc := Baru(nil).DenganKasus(kasusTiruan{ada: kasusUji()}).DenganTransaksi(tanpaTx).DenganJam(jam)
	k, err := svc.BacaKasus(context.Background(), "UJI-NB-1")
	if err != nil || k.Position != "Offer" || k.StatusWork != "Pending-Policy" || k.InsuredName != "UJI TERTANGGUNG" ||
		k.Opportunity.EstimatedClosingDate != "31-10-2026" || k.Opportunity.TypeOfFacultative != "Facultative In" {
		t.Fatalf("%+v (%v)", k, err)
	}
	if k.General.OfferingDate != "03-10-2026" || k.General.BeginDate != "" || k.General.SourceOfBusiness != "UJI SOB" ||
		k.General.OldPolicyNumber != "UJI-LAMA" {
		t.Errorf("general %+v", k.General)
	}
	for _, id := range []string{"UJI-NB-TIDAK-ADA", "", "   ", strings.Repeat("N", 33)} {
		if _, err := svc.BacaKasus(context.Background(), id); !errors.Is(err, ErrKasusTidakAda) {
			t.Errorf("%q: %v, mau 404", id, err)
		}
	}
	if _, err := Baru(nil).BacaKasus(context.Background(), "UJI-NB-1"); !errors.Is(err, ErrKasusTanpaDatabase) {
		t.Errorf("tanpa DB: %v", err)
	}
	galatUji := errors.New("ORA-UJI")
	if _, err := Baru(nil).DenganKasus(kasusTiruan{galat: galatUji}).DenganTransaksi(tanpaTx).BacaKasus(context.Background(), "UJI-NB-2"); !errors.Is(err, galatUji) {
		t.Errorf("galat repository: %v", err)
	}
}

// TestSimpanGeneral - Save for later: simpan sebagian boleh (A90), teks apa adanya
// (78.2/78.3), dibaca kembali; OfferingDate tersimpan mengalahkan PreDT; 401/400/503/404.
func TestSimpanGeneral(t *testing.T) {
	ctx := context.Background()
	akun := inti.Pelaku{AkunID: "UJI-USER"}
	jam := func() time.Time { return time.Date(2026, 10, 3, 1, 0, 0, 0, time.UTC) }
	svc := Baru(nil).DenganKasus(kasusTiruan{ada: kasusUji()}).DenganTransaksi(tanpaTx).DenganJam(jam)
	if _, err := svc.SimpanGeneral(ctx, akun, "UJI-NB-1", IsianGeneral{}); err != nil {
		t.Fatalf("simpan kosong (sebagian): %v", err)
	}
	isian := IsianGeneral{ReffNumber: "UJI-REF", QQName: "UJI QQ", BeginDate: "01-11-2026", OfferingDate: "02-10-2026",
		EndDate: "01-11-2027", PolicyType: "Individual Policy", MarketingID: "UJI-MO", Day: "365", TypeFacultative: "Facultative In"}
	k, err := svc.SimpanGeneral(ctx, akun, "UJI-NB-1", isian)
	g := k.General
	if err != nil || g.ReffNumber != "UJI-REF" || g.QQName != "UJI QQ" || g.BeginDate != "01-11-2026" || g.OfferingDate != "02-10-2026" ||
		g.EndDate != "01-11-2027" || g.PolicyType != "Individual Policy" || g.MarketingID != "UJI-MO" || g.Day != "365" ||
		g.TypeFacultative != "Facultative In" || g.SourceOfBusiness != "UJI SOB" {
		t.Fatalf("%+v (%v)", g, err)
	}
	if _, err := svc.SimpanGeneral(ctx, inti.Pelaku{}, "UJI-NB-1", isian); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("tanpa identitas: %v", err)
	}
	for nama, u := range map[string]struct {
		isian IsianGeneral
		pesan string
	}{
		"tanggal 31-02":  {IsianGeneral{BeginDate: "31-02-2026"}, "beginDate bukan"},
		"tanggal ISO":    {IsianGeneral{OfferingDate: "2026-10-03"}, "offeringDate bukan"},
		"akhir tanpa 0":  {IsianGeneral{EndDate: "1-11-2027"}, "endDate bukan"},
		"ref 51 bita":    {IsianGeneral{ReffNumber: strings.Repeat("U", 51)}, "reffNumber paling banyak 50"},
		"qq 501 bita":    {IsianGeneral{QQName: strings.Repeat("U", 501)}, "qqName paling banyak 500"},
		"MOID 51 bita":   {IsianGeneral{MarketingID: strings.Repeat("U", 51)}, "marketingId paling banyak 50"},
		"policy 51 bita": {IsianGeneral{PolicyType: strings.Repeat("U", 51)}, "policyType paling banyak 50"},
	} {
		if _, err := svc.SimpanGeneral(ctx, akun, "UJI-NB-1", u.isian); !errors.Is(err, ErrMasukanGeneral) || !strings.Contains(err.Error(), u.pesan) {
			t.Errorf("%s: %v, mau 400 %q", nama, err, u.pesan)
		}
	}
	if _, err := Baru(nil).SimpanGeneral(ctx, akun, "UJI-NB-1", isian); !errors.Is(err, ErrKasusTanpaDatabase) {
		t.Errorf("tanpa DB: %v", err)
	}
	if _, err := svc.SimpanGeneral(ctx, akun, "UJI-NB-TIDAK-ADA", isian); !errors.Is(err, ErrKasusTidakAda) {
		t.Errorf("case tidak ada: %v", err)
	}
}

type marketingTiruan struct {
	baris []models.MarketingOfficer
	err   error
}

func (m marketingTiruan) DaftarMarketingOfficer(context.Context) ([]models.MarketingOfficer, error) {
	return m.baris, m.err
}

// TestDaftarMarketingOfficer - urutan repository dipertahankan; 503 tanpa DB.
func TestDaftarMarketingOfficer(t *testing.T) {
	b := []models.MarketingOfficer{{ID: "UJI-2", Nama: "UJI B"}, {ID: "UJI-1", Nama: "UJI A"}}
	got, err := Baru(nil).DenganMarketingOfficer(marketingTiruan{baris: b}).DaftarMarketingOfficer(context.Background())
	if err != nil || len(got) != 2 || got[0] != b[0] {
		t.Errorf("%v (%v)", got, err)
	}
	if _, err := Baru(nil).DaftarMarketingOfficer(context.Background()); !errors.Is(err, ErrMarketingTanpaDatabase) {
		t.Errorf("tanpa DB: %v", err)
	}
}
