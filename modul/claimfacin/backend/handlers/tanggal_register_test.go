package handlers_test

// Proteksi tanggal Input Register (`CheckDate_Act` / `CheckDateReport_Act` / `CheckDateReceived_Act`, aksi ubah medan)
// lewat jalur HTTP. Polis fixture 2026-01-01 s.d. 2026-12-31, hari uji 2026-03-05.

import (
	"net/http"
	"reflect"
	"testing"

	"nusantarare/modul/claimfacin/backend/models"
	"nusantarare/modul/claimfacin/backend/services"
)

func TestProteksiTanggalRegister(t *testing.T) {
	u := baruUji(t)
	id := u.buat()
	kode, out := u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "CopyNB", Konteks: services.ModalPilihPolis,
		Param: polisUji + "|1"})
	u.wajib(kode, http.StatusOK, out, "pilih polis")
	kode, out = u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "SetInputParam", Konteks: services.ModalPilihPolis})
	u.wajib(kode, http.StatusOK, out, "submit pilih polis")
	const dol, lapor, terima = models.CD + "DateOfLoss", models.CD + "ReportDate", models.CD + "DateReceived"
	for _, c := range []struct {
		apa, aksi, d, r, e string
		mau                map[string][]string
	}{
		// keputusan work owner 11-10-2026 (screenshot DOL 01-01-2023): DOL sebelum awal polis ditolak
		{"DOL sebelum awal polis", "CheckDate", "2023-01-01", "2025-01-01", "2025-01-01",
			map[string][]string{dol: {models.PesanDOLLebihMulaiPolis}}},
		{"DOL tepat di awal polis", "CheckDate", "2026-01-01", "2026-01-02", "2026-01-02", nil},
		{"DOL di tengah periode", "CheckDate", "2026-02-01", "2026-02-02", "2026-02-03", nil},
		{"DOL masa depan + lebih dari lapor / terima", "CheckDate", "2026-04-01", "2026-03-01", "2026-03-01",
			map[string][]string{dol: {models.PesanDOLLebihHariIni, models.PesanDOLLebihTerima, models.PesanDOLLebihLapor}}},
		{"DOL sesudah akhir polis", "CheckDate", "2027-02-01", "2027-02-01", "2027-02-01",
			map[string][]string{dol: {models.PesanDOLLebihHariIni, models.PesanDOLLebihAkhirPolis}}},
		{"Report Date < DOL", "CheckDateReport", "2026-02-10", "2026-02-01", "2026-02-15",
			map[string][]string{lapor: {models.PesanLaporKurangDOL}}},
		{"Report Date > Received Date", "CheckDateReport", "2026-02-01", "2026-02-10", "2026-02-05",
			map[string][]string{lapor: {models.PesanLaporLebihTerima}}},
		{"Received Date < Report Date", "CheckDateReceived", "2026-02-01", "2026-02-10", "2026-02-05",
			map[string][]string{terima: {models.PesanTerimaKurangLapor}}},
		{"Received Date masa depan", "CheckDateReceived", "2026-02-01", "2026-02-10", "2026-04-05",
			map[string][]string{terima: {models.PesanTerimaLebihHariIni}}},
	} {
		kode, out := u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: c.aksi,
			Masukan: map[string]string{dol: c.d, lapor: c.r, terima: c.e}})
		u.wajib(kode, http.StatusOK, out, c.apa)
		got := map[string][]string{}
		if p, ok := out["halaman"].(map[string]any)["pesan"].(map[string]any); ok {
			for k, v := range p {
				for _, s := range v.([]any) {
					got[k] = append(got[k], s.(string))
				}
			}
		}
		if len(got) == 0 && len(c.mau) == 0 {
			continue
		}
		if !reflect.DeepEqual(got, c.mau) {
			t.Errorf("%s: pesan %v, mau %v", c.apa, got, c.mau)
		}
	}
}
