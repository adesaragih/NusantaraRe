package services_test

// Reject Outstanding oleh Admin - TANPA Oracle.
//
// Pemilik: tiket 05. Dibaca sesudah: tolak.go.

import (
	"context"
	"errors"
	"strings"
	"testing"

	"nusantarare/internal/models"
	"nusantarare/internal/services"
	"nusantarare/inti"
	"nusantarare/inti/db"
	"nusantarare/inti/galat"
)

func pelakuAdmin() inti.Pelaku {
	return inti.Pelaku{AkunID: "UJI-AKUN", Peran: []string{services.PeranRejectOutstanding}}
}

// TestTolakMenuntutPeranAdmin - `[terverifikasi]` gerbang XML
// `pyWorkPage.pyPosition =='ReasLifeAdmin' && … && .STS_REJECT == 0`
// (`AdjustmentDetail_Section.xml` baris 15399 berkas pecahan).
//
// ⛔ Gerbangnya dipasang DI SINI, bukan ditunggu tiket 07: AC tiket ini
// menuntut penolakan oleh peran lain ditolak, dan gerbang yang tidak ada tidak
// dapat diuji. Tiket 07 kelak menggeneralisasinya, bukan memperkenalkannya.
func TestTolakMenuntutPeranAdmin(t *testing.T) {
	svc := services.New(nil)
	ctx := context.Background()

	// Tanpa identitas sama sekali - 401, bukan 403. Ronde pertama menjawab
	// "bukan ReasLifeAdmin" kepada pemanggil yang sebenarnya belum menyebut
	// dirinya; kedua pertanyaan itu berbeda.
	err := svc.Status().Tolak(ctx, inti.Pelaku{
		Peran: []string{services.PeranRejectOutstanding}}, "CLM-1", "A-1", alasanUji, saatUji)
	if !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("tanpa identitas: galat = %v, mau ErrTanpaIdentitas", err)
	}
	// Beridentitas, tanpa peran sama sekali.
	err = svc.Status().Tolak(ctx, inti.Pelaku{AkunID: "UJI-AKUN"},
		"CLM-1", "A-1", alasanUji, saatUji)
	if !errors.Is(err, inti.ErrTanpaWewenang) {
		t.Errorf("tanpa peran: galat = %v, mau ErrTanpaWewenang", err)
	}
	// Peran lain - SPV boleh menyimpan ke Outstanding, tetapi tidak menolak.
	err = svc.Status().Tolak(ctx, inti.Pelaku{
		AkunID: "UJI-AKUN", Peran: []string{services.PeranSimpanOutstanding}},
		"CLM-1", "A-1", alasanUji, saatUji)
	if !errors.Is(err, inti.ErrTanpaWewenang) {
		t.Errorf("peran SPV: galat = %v, mau ErrTanpaWewenang", err)
	}
	// Dengan peran yang benar, ia lolos gerbang peran dan berhenti di Oracle.
	err = svc.Status().Tolak(ctx, pelakuAdmin(), "CLM-1", "A-1", alasanUji, saatUji)
	if !errors.Is(err, db.ErrTanpaOracle) {
		t.Errorf("peran Admin: galat = %v, mau ErrTanpaOracle", err)
	}
}

// TestTolakMenuntutPengenal - pengenal kosong tidak pernah sampai ke SQL.
func TestTolakMenuntutPengenal(t *testing.T) {
	svc := services.New(nil)
	for _, k := range []struct{ klaim, adj string }{
		{"", "A-1"}, {"CLM-1", ""}, {"  ", "  "},
	} {
		err := svc.Status().Tolak(context.Background(), pelakuAdmin(), k.klaim, k.adj, alasanUji, saatUji)
		if !errors.Is(err, galat.ErrPermintaanTidakSah) {
			t.Errorf("klaim %q adj %q: galat = %v, mau ErrPermintaanTidakSah",
				k.klaim, k.adj, err)
		}
	}
}

// TestKlaimBelumBernomorTidakDapatDitolak - padanan `CLAIM_NO != ”` pada
// gerbang XML. Nomor klaim adalah bukti klaim itu sudah benar-benar terdaftar.
func TestKlaimBelumBernomorTidakDapatDitolak(t *testing.T) {
	if err := services.PeriksaKlaimBernomor(models.Klaim{NomorKlaim: ""}); !errors.Is(
		err, services.ErrKlaimBelumBernomor) {
		t.Errorf("nomor kosong: galat = %v, mau ErrKlaimBelumBernomor", err)
	}
	if err := services.PeriksaKlaimBernomor(models.Klaim{NomorKlaim: "   "}); !errors.Is(
		err, services.ErrKlaimBelumBernomor) {
		t.Errorf("nomor spasi: galat = %v, mau ErrKlaimBelumBernomor", err)
	}
	if err := services.PeriksaKlaimBernomor(
		models.Klaim{NomorKlaim: "UJI-CLM-1"}); err != nil {
		t.Errorf("nomor terisi ditolak: %v", err)
	}
}

// TestNilaiTolakSamaDenganKomite - AC: tidak ada nilai khusus yang membedakan
// penolakan Admin dari penolakan Komite. `[terverifikasi]` korpus menulis `2`
// dalam bentuk yang sama persis di kedua jalur.
func TestNilaiTolakSamaDenganKomite(t *testing.T) {
	hasil, err := services.Transisi(
		baris("A", models.KodeOutstanding), models.StatusDitolak, saatUji)
	if err != nil {
		t.Fatalf("Transisi: %v", err)
	}
	if hasil.KodeStatus != models.KodeDitolak {
		t.Errorf("kode = %q, mau %q - satu nilai untuk kedua sumber penolakan",
			hasil.KodeStatus, models.KodeDitolak)
	}
}

// TestTolakBarisFinalDitolak - penolakan hanya mungkin pada baris yang masih
// Outstanding; aturannya dipakai ulang dari Transisi, bukan ditulis dua kali.
func TestTolakBarisFinalDitolak(t *testing.T) {
	for _, kode := range []string{models.KodeAksep, models.KodeDitolak} {
		_, err := services.Transisi(baris("A", kode), models.StatusDitolak, saatUji)
		if !errors.Is(err, services.ErrBarisSudahFinal) {
			t.Errorf("dari %q: galat = %v, mau ErrBarisSudahFinal", kode, err)
		}
	}
}

// alasanUji - isian `Remarks` dialog Reject Outstanding (b1687).
const alasanUji = "UJI alasan penolakan"

// OQ-M5 DITUTUP (GILIRAN-17): `Remarks` WAJIB (b1653, b1695 `always`, b1699),
// dan disimpan di `T_CLAIMLF_JEJAK.KOMENTAR` VARCHAR2(4000) - kosong atau
// melampaui lebar ditolak di pintu, sebelum basis data.
func TestTolakMenuntutAlasan(t *testing.T) {
	svc := services.New(nil)
	for _, alasan := range []string{"", "   ", strings.Repeat("x", services.BatasKomentarJejak+1)} {
		err := svc.Status().Tolak(context.Background(), pelakuAdmin(), "CLM-1", "A-1", alasan, saatUji)
		if !errors.Is(err, galat.ErrPermintaanTidakSah) {
			t.Errorf("alasan %d karakter: galat = %v, mau ErrPermintaanTidakSah", len(alasan), err)
		}
	}
	err := svc.Status().Tolak(context.Background(), pelakuAdmin(), "CLM-1", "A-1",
		strings.Repeat("x", services.BatasKomentarJejak), saatUji)
	if !errors.Is(err, db.ErrTanpaOracle) {
		t.Errorf("alasan tepat di batas: galat = %v, mau ErrTanpaOracle", err)
	}
}
