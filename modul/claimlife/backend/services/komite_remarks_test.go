package services_test

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/claimlife/backend/services"
)

// 12-10-2026 (opsi A) - Remarks pop-up Send Claim to Committee: wajib, ≤ 4000 byte, masuk KOMENTAR jejak penyerahan.

func TestPeriksaRemarksKomite(t *testing.T) {
	for _, u := range []struct {
		masuk, mau string
		galat      error
	}{
		{"", "", services.ErrRemarksKomiteKosong},
		{"   ", "", services.ErrRemarksKomiteKosong},
		{" UJI alasan ", "UJI alasan", nil},
		{strings.Repeat("a", 4000), strings.Repeat("a", 4000), nil},
		{strings.Repeat("a", 4001), "", services.ErrRemarksKomiteKepanjangan},
		{strings.Repeat("é", 2001), "", services.ErrRemarksKomiteKepanjangan}, // byte, bukan karakter
	} {
		got, err := services.PeriksaRemarksKomite(u.masuk)
		if (u.galat == nil && err != nil) || (u.galat != nil && !errors.Is(err, u.galat)) || got != u.mau {
			t.Errorf("%q: (%q, %v), mau (%q, %v)", u.masuk[:min(len(u.masuk), 12)], got, err, u.mau, u.galat)
		}
	}
}

func TestSerahkanMenolakRemarksKosongSebelumOracle(t *testing.T) {
	svc := services.New(nil)
	err := svc.Komite().Serahkan(context.Background(), pelakuBerperan(inti.PeranSPV), "CLM-1", "P-1", "A-1", " ", saatUji)
	if !errors.Is(err, services.ErrRemarksKomiteKosong) {
		t.Errorf("Remarks kosong: galat = %v, mau ErrRemarksKomiteKosong", err)
	}
}

func TestPratinjauKomiteMenjagaPagarnya(t *testing.T) {
	svc := services.New(nil)
	ctx := context.Background()
	if _, err := svc.Komite().Pratinjau(ctx, inti.Pelaku{}, "CLM-1", "P-1", "A-1", saatUji); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("tanpa identitas: %v", err)
	}
	if _, err := svc.Komite().Pratinjau(ctx, pelakuBerperan(inti.PeranSPV), "CLM-1", "P-1", "A-1", saatUji); !errors.Is(err, db.ErrTanpaOracle) {
		t.Errorf("tanpa Oracle: %v", err)
	}
}

// Pratinjau menjalankan gerbang Serahkan yang SAMA sampai sebelum menulis; Remarks hanya ke jejak, tidak ke efek keluar.
func TestPratinjauDanSerahkanBerbagiGerbang(t *testing.T) {
	isi, err := os.ReadFile("komite.go")
	if err != nil {
		t.Fatal(err)
	}
	s := string(isi)
	badan := func(kepala string) string {
		i := strings.Index(s, kepala)
		if i < 0 {
			t.Fatalf("tidak ada %s", kepala)
		}
		j := strings.Index(s[i+len(kepala):], "\nfunc ")
		if j < 0 {
			return s[i:]
		}
		return s[i : i+len(kepala)+j]
	}
	serahkan := badan("func (p *Penyerahan) Serahkan(")
	pratinjau := badan("func (p *Penyerahan) Pratinjau(")
	for _, g := range []string{"WajibIdentitas(pelaku)", "PastikanKasusTerbuka(ctx, klaimID)", "WajibWewenangKomite(pelaku, tipe)",
		"if tahap != models.TahapClaimAnalis {", "p.siapkanPenyerahan(ctx, baca, klaimID, pesertaID, adjID)"} {
		if !strings.Contains(serahkan, g) || !strings.Contains(pratinjau, g) {
			t.Errorf("gerbang %q tidak ada di kedua fungsi", g)
		}
	}
	if strings.Contains(pratinjau, "DalamTransaksi(") || strings.Contains(pratinjau, "Salurkan(") {
		t.Error("Pratinjau menulis / menyalurkan efek")
	}
	if !strings.Contains(serahkan, "Komentar: remarks,") {
		t.Error("Remarks tidak masuk KOMENTAR jejak")
	}
	salur := serahkan[strings.Index(serahkan, "p.penyalur.Salurkan("):]
	salur = salur[:strings.Index(salur, "})")]
	if strings.Contains(salur, "remarks") || strings.Contains(salur, "Remarks") {
		t.Error("Remarks ikut ke MuatanEfek (email/outbox) tanpa keputusan")
	}
	siap := badan("func (p *Penyerahan) siapkanPenyerahan(")
	for _, g := range []string{"PeriksaSatuMataUang(", "PeriksaBolehDiserahkan(baris)", "PeriksaRekening(baris)", "AmbangRoster(",
		"p.roster.AmbilAnggota(", "PeriksaTingkatKomite(tingkat)"} {
		if !strings.Contains(siap, g) {
			t.Errorf("siapkanPenyerahan tanpa %s", g)
		}
	}
}

// 12-10-2026 - gerbang rekening penyerahan Komite MATI sementara (menunggu master bank DBA), dapat dinyalakan kembali.
func TestSaklarGerbangRekeningKomite(t *testing.T) {
	if services.GerbangRekeningKomiteAktif {
		t.Error("GerbangRekeningKomiteAktif bawaan harus false sampai master bank DBA tersedia (keputusan WO 12-10-2026)")
	}
	isi, err := os.ReadFile("komite.go")
	if err != nil {
		t.Fatal(err)
	}
	s := string(isi)
	siap := s[strings.Index(s, "func (p *Penyerahan) siapkanPenyerahan("):]
	siap = siap[:strings.Index(siap, "\n}\n")]
	if !strings.Contains(siap, "if GerbangRekeningKomiteAktif {\n\t\tif err := PeriksaRekening(baris); err != nil {") {
		t.Error("PeriksaRekening tidak lagi di balik saklar di siapkanPenyerahan (dipakai Pratinjau DAN Serahkan)")
	}
	// Kode gerbang dan pesannya TETAP.
	if err := services.PeriksaRekening(barisLengkap(t)); err != nil {
		t.Errorf("rekening lengkap ditolak: %v", err)
	}
	b := barisLengkap(t)
	b.IDBank = ""
	if err := services.PeriksaRekening(b); !errors.Is(err, services.ErrRekeningBelumLengkap) || err.Error() != "Name of bank cannot be empty" {
		t.Errorf("rekening tak lengkap: %v", err)
	}
}
