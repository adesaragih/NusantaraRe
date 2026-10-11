package services

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// caseKotakMasuk - satu baris T_WORK_POLIS tiruan; team = team case (MOID -> MARKETINGOFFICER.TEAMGROUP).
type caseKotakMasuk struct {
	id, posisi, team, status string
}

// kotakMasukTiruan - PembacaKotakMasuk yang menerapkan semantik saringan repository atas case dalam memori.
type kotakMasukTiruan struct {
	cases           []caseKotakMasuk
	nama            map[string]string
	teamMkt, teamTL map[string][]string // login -> team
	bacaTeam        *int
}

func (k kotakMasukTiruan) cocok(c caseKotakMasuk, s repository.SaringanKotakMasuk) bool {
	if strings.HasPrefix(c.status, "Resolved-") {
		return false
	}
	for _, t := range s.Tim {
		if dalam(t.Posisi, c.posisi) && dalam(t.Team, c.team) {
			return true
		}
	}
	return dalam(s.Lain, c.posisi)
}

func (k kotakMasukTiruan) HitungKotakMasuk(_ context.Context, s repository.SaringanKotakMasuk) (map[string]int, error) {
	out := map[string]int{}
	for _, c := range k.cases {
		if k.cocok(c, s) {
			out[c.posisi]++
		}
	}
	return out, nil
}

func (k kotakMasukTiruan) DaftarMenunggu(_ context.Context, s repository.SaringanKotakMasuk) ([]models.BarisKotakMasuk, error) {
	out := []models.BarisKotakMasuk{}
	for _, c := range k.cases {
		if k.cocok(c, s) {
			out = append(out, models.BarisKotakMasuk{CaseID: c.id})
		}
	}
	return out, nil
}

func (k kotakMasukTiruan) NamaWorkbasket(_ context.Context, id []string) (map[string]string, error) {
	out := map[string]string{}
	for _, x := range id {
		if n, ada := k.nama[x]; ada {
			out[x] = n
		}
	}
	return out, nil
}

func (k kotakMasukTiruan) TeamKotakMasuk(_ context.Context, login string) ([]string, []string, error) {
	if k.bacaTeam != nil {
		*k.bacaTeam++
	}
	return k.teamMkt[login], k.teamTL[login], nil
}

var kotakMasukUji = kotakMasukTiruan{
	cases: []caseKotakMasuk{
		{"UJI-NB-1", "Offer", "1", "Pending-Policy"},
		{"UJI-NB-2", "ReasFacInMarketing", "1", "Pending-Policy"},
		{"UJI-NB-3", "Offer", "2", "Pending-Policy"},
		{"UJI-NB-4", "ReasFacInTeamLeader", "1", "Pending-Policy"},
		{"UJI-NB-5", "ReasFacInTeamLeader", "3", "Pending-Policy"},
		{"UJI-NB-6", "Decline", "1", "Resolved-Decline"},
		{"UJI-NB-7", "ReasFacInTeamLeader", "1", "Resolved-Decline"},
		{"UJI-NB-8", "ReasFacInSeniorUnderwriting", "2", "Pending-Policy"},
		{"UJI-NB-9", "ReasFacInMarketing", "", "Pending-Policy"},
	},
	nama:    map[string]string{"ReasFacInTeamLeader": "Team Leader Fac In"},
	teamMkt: map[string][]string{"UJI-MKT": {"1"}},
	teamTL:  map[string][]string{"UJI-MKT": {"1"}, "UJI-TL": {"1", "3"}},
}

// TestKotakMasuk - marketing / TL: case se-team akun (Offer + ReasFacInMarketing dijumlah); nama M_WORKBASKET, kosong ->
// kode; akun tanpa workbasket antrean -> []; 401; 503.
func TestKotakMasuk(t *testing.T) {
	ctx := context.Background()
	svc := Baru(nil).DenganKotakMasuk(kotakMasukUji)
	p := inti.Pelaku{AkunID: "UJI-MKT", Peran: []string{"ReasFacInTeamLeader", "ReasFacInMarketing", "WB-MODUL-LAIN"}}
	d, err := svc.KotakMasuk(ctx, p)
	mau := []AntreanKotakMasuk{{"ReasFacInMarketing", "ReasFacInMarketing", 2}, {"ReasFacInTeamLeader", "Team Leader Fac In", 1}}
	if err != nil || !reflect.DeepEqual(d, mau) {
		t.Fatalf("%v %+v", err, d)
	}
	if d, err := svc.KotakMasuk(ctx, inti.Pelaku{AkunID: "UJI", Peran: []string{"WB-MODUL-LAIN"}}); err != nil || len(d) != 0 || d == nil {
		t.Errorf("tanpa workbasket antrean: %v %+v", err, d)
	}
	if _, err := svc.KotakMasuk(ctx, inti.Pelaku{}); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("401: %v", err)
	}
	if _, err := Baru(nil).KotakMasuk(ctx, p); !errors.Is(err, ErrKotakMasukTanpaDatabase) {
		t.Errorf("503: %v", err)
	}
}

// TestKotakMasukPerTeam - keputusan work owner 08-10-2026: matriks marketing / TL / UW × punya team / tanpa team × team case
// sama / beda; hitungan = daftar (satu saringan); UW tanpa saringan team dan tanpa membaca team akun.
func TestKotakMasukPerTeam(t *testing.T) {
	ctx := context.Background()
	n := 0
	k := kotakMasukUji
	k.bacaTeam = &n
	svc := Baru(nil).DenganKotakMasuk(k)
	id := func(b []BarisInbox) []string {
		s := []string{}
		for _, x := range b {
			s = append(s, x.CaseID)
		}
		return s
	}
	for _, c := range []struct {
		akun, wb string
		mau      []string
	}{
		{"UJI-MKT", "ReasFacInMarketing", []string{"UJI-NB-1", "UJI-NB-2"}},   // team 1; NB-3 team 2 dan NB-9 tanpa team tidak
		{"UJI-MKT-LAIN", "ReasFacInMarketing", []string{}},                    // tanpa team -> kosong
		{"UJI-TL", "ReasFacInTeamLeader", []string{"UJI-NB-4", "UJI-NB-5"}},   // team 1 dan 3 (IN)
		{"UJI-MKT", "ReasFacInTeamLeader", []string{"UJI-NB-4"}},              // team 1 saja
		{"UJI-MKT-LAIN", "ReasFacInTeamLeader", []string{}},                   // tanpa team -> kosong
		{"UJI-MKT-LAIN", "ReasFacInSeniorUnderwriting", []string{"UJI-NB-8"}}, // UW tanpa saringan team
	} {
		p := inti.Pelaku{AkunID: c.akun, Peran: []string{c.wb}}
		d, err := svc.DaftarMenunggu(ctx, p, c.wb)
		if err != nil || !reflect.DeepEqual(id(d), c.mau) {
			t.Errorf("%s/%s daftar: %v %v", c.akun, c.wb, err, id(d))
		}
		h, err := svc.KotakMasuk(ctx, p)
		if err != nil || len(h) != 1 || h[0].Jumlah != len(c.mau) {
			t.Errorf("%s/%s hitung != daftar: %v %+v", c.akun, c.wb, err, h)
		}
	}
	n = 0
	if _, err := svc.KotakMasuk(ctx, inti.Pelaku{AkunID: "UJI-UW", Peran: []string{"ReasFacInUnderwriting", "ReasFacInSeniorUnderwriting"}}); err != nil || n != 0 {
		t.Errorf("UW membaca team akun: %v %d", err, n)
	}
}

// TestDaftarMenunggu - workbasket dipegang / kosong = semua yang dipegang; workbasket asing / bukan antrean -> 403; NB
// status kosong -> "NEW NB".
func TestDaftarMenunggu(t *testing.T) {
	ctx := context.Background()
	svc := Baru(nil).DenganKotakMasuk(kotakMasukUji)
	tl := inti.Pelaku{AkunID: "UJI-TL", Peran: []string{"ReasFacInTeamLeader"}}
	id := func(b []BarisInbox) []string {
		var s []string
		for _, x := range b {
			s = append(s, x.CaseID)
		}
		return s
	}
	if d, err := svc.DaftarMenunggu(ctx, tl, "ReasFacInTeamLeader"); err != nil || !reflect.DeepEqual(id(d), []string{"UJI-NB-4", "UJI-NB-5"}) ||
		d[0].StatusInbox != "NEW NB" {
		t.Errorf("TL: %v %+v", err, d)
	}
	mkt := inti.Pelaku{AkunID: "UJI-MKT", Peran: []string{"ReasFacInMarketing", "ReasFacInSeniorUnderwriting"}}
	if d, err := svc.DaftarMenunggu(ctx, mkt, ""); err != nil || !reflect.DeepEqual(id(d), []string{"UJI-NB-1", "UJI-NB-2", "UJI-NB-8"}) {
		t.Errorf("semua dipegang: %v %v", err, id(d))
	}
	for _, wb := range []string{"ReasFacInMarketing", "WB-MODUL-LAIN"} {
		if _, err := svc.DaftarMenunggu(ctx, tl, wb); !errors.Is(err, inti.ErrTanpaWewenang) {
			t.Errorf("403 %s: %v", wb, err)
		}
	}
	if d, err := svc.DaftarMenunggu(ctx, inti.Pelaku{AkunID: "UJI"}, ""); err != nil || len(d) != 0 || d == nil {
		t.Errorf("tanpa workbasket: %v %+v", err, d)
	}
	if _, err := svc.DaftarMenunggu(ctx, inti.Pelaku{}, ""); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("401: %v", err)
	}
}

// TestKeBarisInbox - pemetaan kolom inbox "New Business": tanggal Pega -> DD-MM-YYYY (Begin date 12:00 WIB, Offering
// YYYYMMDD), NB status kosong -> "NEW NB", aksi terakhir = Approval + " by " + PIC (tanpa riwayat -> ""), validityDate "".
func TestKeBarisInbox(t *testing.T) {
	b := models.BarisKotakMasuk{CaseID: "UJI-NB-1", InsuredName: "UJI TERTANGGUNG", BusinessName: "UJI COB", SourceOfBusiness: "UJI SOB",
		CedingCompany: "UJI A;UJI B", StartDateTime: "20261001T050000.000 GMT", ShareRNM: "12.5", Marketing: "UJI MO",
		FlagOnGoingPolicy: "0", CreatedDate: "07-10-2026", OfferingDate: "20261002", NBStatus: "NB IS IN TEAM LEADER'S INBOX",
		ApprovalTerakhir: "Accept", PICAkhir: "UJI PIC", TglUpdate: "2026-10-07 22:04:05"}
	mau := BarisInbox{CaseID: "UJI-NB-1", InsuredName: "UJI TERTANGGUNG", BusinessName: "UJI COB", SourceOfBusiness: "UJI SOB",
		CedingCompany: "UJI A;UJI B", InceptionDate: "01-10-2026", ShareRNM: "12.5", Marketing: "UJI MO", Status: "Offer",
		CreatedDate: "07-10-2026", OfferingDate: "02-10-2026", StatusInbox: "NB IS IN TEAM LEADER'S INBOX", AksiTerakhir: "Accept by UJI PIC",
		TglUpdate: "2026-10-07 22:04:05"}
	if got := keBarisInbox(b); got != mau {
		t.Errorf("%+v", got)
	}
	if got := keBarisInbox(models.BarisKotakMasuk{CaseID: "UJI-NB-2"}); got.StatusInbox != "NEW NB" || got.AksiTerakhir != "" || got.ValidityDate != "" {
		t.Errorf("kosong: %+v", got)
	}
	// submitToUw (tiket 58, "Time in Underwriting"; sumber SUBMIT_TO_UW): teks Pega DateTime UTC -> YYYY-MM-DD HH:MM:SS WIB; kosong / bukan
	// BentukWaktuPega -> "" (bukan teks mentah).
	for asal, mau := range map[string]string{"20261006T020000.000 GMT": "2026-10-06 09:00:00", "20261006T200000.000 GMT": "2026-10-07 03:00:00",
		"": "", "bukan-waktu": "", "20261006": "", "2026-10-06 09:00:00": ""} {
		if got := keBarisInbox(models.BarisKotakMasuk{SubmitToUW: asal}).SubmitToUW; got != mau {
			t.Errorf("submitToUw %q: %q, mau %q", asal, got, mau)
		}
	}
	// Status = padanan .IsCedingConfirm dari FLAG_ONGOING_POLICY: 0 Offer, 2 Binding, lainnya "" (tanpa tebakan Policy / Accepted).
	for flag, mau := range map[string]string{"0": "Offer", "2": "Binding", "1": "", "": "", "9": ""} {
		if got := keBarisInbox(models.BarisKotakMasuk{FlagOnGoingPolicy: flag}).Status; got != mau {
			t.Errorf("flag %q: %q, mau %q", flag, got, mau)
		}
	}
}

// TestAntreanKasusNB - setiap posisi yang dapat dituju Submit (Team Leader, tujuan Banding) ada di antrean kotak masuk.
func TestAntreanKasusNB(t *testing.T) {
	for _, tj := range tujuanBanding {
		if !dalam(antreanKasusNB, tj.posisi) {
			t.Errorf("%s tidak di antrean", tj.posisi)
		}
	}
	if !dalam(antreanKasusNB, posisiTeamLeader) || !dalam(antreanKasusNB, "ReasFacInMarketing") {
		t.Error("marketing / Team Leader")
	}
	// Tiket 62: tujuan Accept tahap uw (statusInboxPosisi) dan posisi tahap Submit uw (Decision24 / jenjang lanjutan).
	for _, daftar := range []map[string]bool{posisiDecision24, posisiJenjangLanjut} {
		for pos := range daftar {
			if !dalam(antreanKasusNB, pos) {
				t.Errorf("%s tidak di antrean", pos)
			}
		}
	}
	for pos := range statusInboxPosisi {
		if !dalam(antreanKasusNB, pos) {
			t.Errorf("tujuan %s tidak di antrean", pos)
		}
	}
}
