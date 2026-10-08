package services_test

import (
	"bytes"
	"context"
	"encoding/csv"
	"slices"
	"strings"
	"testing"
	"time"

	"nusantarare/modul/claimprop/backend/models"
	"nusantarare/modul/claimprop/backend/services"
	"nusantarare/modul/claimprop/backend/tiruan"
)

var (
	jam1 = time.Date(2026, 5, 1, 8, 0, 0, 0, time.UTC)
	jam2 = time.Date(2026, 5, 3, 9, 0, 0, 0, time.UTC)
)

// gudangLama - tiga kasus lama: UJI-CLMP-7 (dua baris OS, terakhir tutup berkas), UJI-CLMP-8 (akseptasi komite, STS 1),
// UJI-CLMP-9 (satu baris acceptation + halaman JSON_KLAIM).
func gudangLama() *tiruan.Gudang {
	g := tiruan.Baru()
	g.OSLama = []models.BarisOSLama{
		{CaseID: models.KunciPegaLama("UJI-CLMP-7"), NoClaim: "UJI-NO-7", NoPolis: "UJI-P7", InsertOp: "UJI-OP1",
			StsReject: models.StsOSEstimasi, Tanggal: jam1, DataJSON: `{"CauseOfLoss":"UJI-SEBAB","Value":"10"}`},
		{CaseID: models.KunciPegaLama("UJI-CLMP-7"), NoClaim: "UJI-NO-7", InsertOp: "UJI-OP2",
			StsReject: models.StsOSTutupBerkas, Tanggal: jam2, DataJSON: `{"AcceptedNo":"UJI-A1"}`, AcceptedNo: "UJI-A1"},
		{CaseID: models.KunciPegaLama("UJI-CLMP-8"), StsReject: "1", Tanggal: jam1},
		{CaseID: models.KunciPegaLama("UJI-CLMP-9"), NoClaim: "UJI-NO-9", InsertOp: "UJI-OP3",
			StsReject: models.StsOSAkseptasi, Tanggal: jam1},
	}
	g.JSONLama = map[string]string{models.KunciPegaLama("UJI-CLMP-9"): `{"NoClaim":"UJI-NO-9","pxObjClass":"UJI",
		"EstimationList":[{"EstimationDate":"20260501","GrossEstimationPct":"100"}],"DaftarObjek":"UJI-OBJEK"}`}
	return g
}

func jalankanPemuat(t *testing.T, g *tiruan.Gudang, tulis bool) (services.RingkasanPemuat, string, string) {
	t.Helper()
	var arsip, galat bytes.Buffer
	lap := services.LaporanPemuatBaru(&arsip, &galat, tulis)
	if err := services.PemuatBaru(g, nil).Jalankan(context.Background(), tulis, lap); err != nil {
		t.Fatal(err)
	}
	if err := lap.Tutup(); err != nil {
		t.Fatal(err)
	}
	return lap.Ringkasan(), arsip.String(), galat.String()
}

func TestPemuatUjiKeringTidakMenulis(t *testing.T) {
	g := gudangLama()
	r, arsip, galat := jalankanPemuat(t, g, false)
	if r.Kasus != 3 || r.Baris != 4 || r.Riwayat != 1 || r.DariJSON != 1 || r.Siap != 3 || r.Gagal != 0 || r.Dimuat != 0 {
		t.Fatalf("ringkasan uji-kering %+v", r)
	}
	if r.PerTahap[models.StatusSelesai] != 1 || r.PerTahap[models.TahapAcceptation] != 2 {
		t.Fatalf("per tahap %+v", r.PerTahap)
	}
	if len(g.Kasus) != 0 {
		t.Fatalf("uji-kering menulis %d kasus", len(g.Kasus))
	}
	if !adaBaris(t, arsip, "UJI-CLMP-7", "AcceptedNo", "UJI-A1", models.SebabBarisOS) ||
		!adaBaris(t, arsip, "UJI-CLMP-9", "ClaimData.DaftarObjek", "UJI-OBJEK", models.SebabTanpaKolom) {
		t.Fatalf("arsip medan:\n%s", arsip)
	}
	if strings.Contains(galat, "UJI-CLMP-8") || !r.Selesai() { // STS 1 dimuat (keputusan work owner 08-10-2026)
		t.Fatalf("galat:\n%s (selesai %v)", galat, r.Selesai())
	}
	if !strings.Contains(r.Teks(), "uji-kering") {
		t.Fatalf("teks ringkasan tanpa mode:\n%s", r.Teks())
	}
}

func TestPemuatJalankanMenulisLaluMelewatiUlang(t *testing.T) {
	g := gudangLama()
	r, _, _ := jalankanPemuat(t, g, true)
	if r.Dimuat != 3 || r.Gagal != 0 || r.Dilewati != 0 {
		t.Fatalf("ringkasan jalankan %+v", r)
	}
	k7 := g.Kasus["UJI-CLMP-7"]
	if k7.Tahap != models.TahapAcceptation || !k7.Tertutup() || k7.Sumber != models.SumberPega ||
		k7.PembuatID != "UJI-OP1" || !k7.TglCreate.Equal(jam1) || !k7.TglUpdate.Equal(jam2) {
		t.Fatalf("kasus UJI-CLMP-7 %+v", k7)
	}
	if h := g.Halaman("UJI-CLMP-7"); h.Ambil(models.CD+"NoClaim") != "UJI-NO-7" || h.Ambil(models.CD+"CauseOfLoss") != "" {
		t.Fatalf("header UJI-CLMP-7 dari baris BERLAKU (tanpa CauseOfLoss): %+v", h.Nilai)
	}
	k9 := g.Kasus["UJI-CLMP-9"]
	if k9.Tahap != models.TahapAcceptation || k9.Tertutup() || k9.Posisi != models.WorkbasketAcceptation {
		t.Fatalf("kasus UJI-CLMP-9 %+v", k9)
	}
	if n := len(g.Halaman("UJI-CLMP-9").AmbilDaftar(models.DaftarEstimasi)); n != 1 {
		t.Fatalf("EstimationList UJI-CLMP-9 %d baris", n)
	}
	if k8 := g.Kasus["UJI-CLMP-8"]; k8.Tahap != models.TahapAcceptation || k8.Tertutup() {
		t.Fatalf("kasus akseptasi komite (STS 1) di Input Acceptation, terbuka: %+v", k8)
	}
	r, _, _ = jalankanPemuat(t, g, true)
	if r.Dimuat != 0 || r.Dilewati != 3 {
		t.Fatalf("jalankan ulang mau melewati: %+v", r)
	}
}

// adaBaris - CSV `teks` memuat baris persis `kolom`.
func adaBaris(t *testing.T, teks string, kolom ...string) bool {
	t.Helper()
	rows, err := csv.NewReader(strings.NewReader(teks)).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range rows {
		if slices.Equal(r, kolom) {
			return true
		}
	}
	return false
}
