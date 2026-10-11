package handlers_test

// Tab Loss Record / Correspondence tampilan polis NB (keputusan work owner 11-10-2026 "Ikut NB: hanya akun IT"): rute
// sendiri, ditegakkan di layanan - akun ber-JOB_POSITION "IT" saja (nbfacin tiket 55 / 56), selain itu 403 tanpa isi.

import (
	"net/http"
	"testing"

	"nusantarare/modul/claimfacin/backend/handlers"
	"nusantarare/modul/claimfacin/backend/models"
	"nusantarare/modul/claimfacin/backend/tiruan"
)

func TestPolisNBITHanyaAkunIT(t *testing.T) {
	u := baruUji(t)
	id := u.sampaiSurveyor()
	u.a.PolisIT[tiruan.KunciPolis(polisUji, "1")] = models.PolisNBIT{
		LossRecord: models.LossRecordNB{PremiPar: "1000", Gross1: "12.5", Baris: []models.BarisLossRecordNB{
			{UwYear: "2025", ClassOfBusiness: "UJI FIRE", Currency: "IDR", LossAmount: "250", LossAmountIDR: "250"}}},
		Korespondensi: []models.KorespondensiNB{{Sender: "UJI PENGIRIM", Subject: "UJI SUBJEK", Terkirim: true,
			Penerima: []models.PenerimaNB{{Recipient: "UJI PENERIMA", RecipientMail: "uji@contoh.invalid"}}}},
	}
	jalur := handlers.Prefix + "/kasus/" + id + "/polis-nb-it"

	// bukan IT (jabatan lain / tanpa jabatan): 403, isi tidak dikirim
	u.a.Jabatan[admin] = "UJI-MARKETING"
	kode, out := u.minta(http.MethodGet, jalur, admin, "", nil)
	u.wajib(kode, http.StatusForbidden, out, "bukan IT")
	if _, ada := out["lossRecord"]; ada {
		t.Fatalf("isi tab IT bocor ke akun bukan IT: %v", out)
	}
	// tanpa identitas: 401
	kode, out = u.minta(http.MethodGet, jalur, "", "", nil)
	u.wajib(kode, http.StatusUnauthorized, out, "tanpa identitas")

	// IT (tanpa beda huruf / spasi tepi, seperti nbfacin): isi tersimpan apa adanya
	u.a.Jabatan[admin] = " it "
	kode, out = u.minta(http.MethodGet, jalur, admin, "", nil)
	u.wajib(kode, http.StatusOK, out, "IT")
	lr := out["lossRecord"].(map[string]any)
	if lr["premiPar"] != "1000" || lr["gross1"] != "12.5" || len(lr["baris"].([]any)) != 1 {
		t.Fatalf("lossRecord: %v", lr)
	}
	kor := out["korespondensi"].([]any)
	if len(kor) != 1 || kor[0].(map[string]any)["terkirim"] != true ||
		len(kor[0].(map[string]any)["penerima"].([]any)) != 1 {
		t.Fatalf("korespondensi: %v", kor)
	}

	// kasus tidak ada: 404 (IT)
	kode, out = u.minta(http.MethodGet, handlers.Prefix+"/kasus/CLM-999999/polis-nb-it", admin, "", nil)
	u.wajib(kode, http.StatusNotFound, out, "kasus tidak ada")
}
