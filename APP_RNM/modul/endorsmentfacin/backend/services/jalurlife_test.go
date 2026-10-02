package services_test

// Jalur Life endorsement - tiket E18: rute flow sesudah konfirmasi Marketing
// dan penomoran endorsement.
//
// Dibaca sesudah: services/jalurlife.go.

import (
	"reflect"
	"testing"

	"nusantarare/modul/endorsmentfacin/backend/services"
)

// TestLifeMelewatiTanggaAkseptasi - E18/K-044: sesudah Marketing menyetujui,
// Life langsung ke pemeriksaan galat konversi (lalu simpan polis); non-Life ke
// cabang "Is it group?" yang membuka tangga akseptasi.
func TestLifeMelewatiTanggaAkseptasi(t *testing.T) {
	if got := services.LangkahSesudahKonfirmasiMarketing(true); got != services.LangkahCekGalatKonversi {
		t.Errorf("Life → %q, mau %q", got, services.LangkahCekGalatKonversi)
	}
	if got := services.LangkahSesudahKonfirmasiMarketing(false); got != services.LangkahCekGrup {
		t.Errorf("non-Life → %q, mau %q", got, services.LangkahCekGrup)
	}
	if services.LangkahCekGalatKonversi.MembukaTanggaAkseptasi() {
		t.Error("langkah jalur Life membuka tangga akseptasi")
	}
	if !services.LangkahCekGrup.MembukaTanggaAkseptasi() {
		t.Error("langkah jalur non-Life tidak membuka tangga akseptasi")
	}
}

// TestRencanaNomorEndorsemen - langkah 14-16 `SaveEDMToJsonPolicy_Act` (12-13
// di-remark, tidak diport): pasangan komplementer Life/non-Life, hanya bila
// nomor endorsement kosong.
func TestRencanaNomorEndorsemen(t *testing.T) {
	for _, u := range []struct {
		nama    string
		life    bool
		nomor   string
		rencana []services.PermintaanData
	}{
		{"non-Life, nomor kosong", false, "", []services.PermintaanData{
			services.MintaKodeProdNonLife, services.MintaNomorUrut}},
		{"Life, nomor kosong", true, "", []services.PermintaanData{
			services.MintaKodeProdLife, services.MintaNomorUrut}},
		{"nomor sudah ada", true, "UJI-EDM-1", nil},
		{"nomor sudah ada, non-Life", false, "UJI-EDM-1", nil},
	} {
		if got := services.RencanaNomorEndorsemen(u.life, u.nomor); !reflect.DeepEqual(got, u.rencana) {
			t.Errorf("%s: %v, mau %v", u.nama, got, u.rencana)
		}
	}
}

// TestSusunNomorEndorsemen - langkah 14.2/15.2 + 17:
// (HASIL3 + "E") + BusinessOldId + "." + HASIL1 + "." + HASIL2.
func TestSusunNomorEndorsemen(t *testing.T) {
	got := services.SusunNomorEndorsemen(services.HasilNomorEndorsemen{
		KodeProd: "UJI-K", BusinessOldId: "L9", Hasil1: "10.2026", Hasil2: "00042",
	})
	if want := "UJI-KEL9.10.2026.00042"; got != want {
		t.Errorf("%q, mau %q", got, want)
	}
}
