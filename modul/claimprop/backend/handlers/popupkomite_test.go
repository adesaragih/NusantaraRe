package handlers_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"nusantarare/modul/claimprop/backend/handlers"
	"nusantarare/modul/claimprop/backend/models"
	"nusantarare/modul/claimprop/backend/services"
)

// Laporan work owner 09-10-2026 pada popup "Komite klaim Treaty": Date dan Initial kosong (Initial jadi PIC), tombol
// Send tidak ada sesudah Circumstanses / Remarks diisi. XML: `ProteksiInitialandDate_Act` langkah 2 mengisi
// `TreatyExchangeYearly.UserName = OperatorID.pyLabel` dan `TempCommiteClaim.DateOfComitee` = hari ini; Remarks /
// Circumstanses ber-event change "Refresh this section"; kontainer tombol tampil bila `Protect.CARI1 = 1 &&
// Protect.CARI2 = 1` - halaman sementara itu hidup selama popup terbuka, jadi layar membawanya di `mode`.
func TestPopupKomiteDateDanPICTerisiDanKirimTampilSesudahRefresh(t *testing.T) {
	u, _ := baruUjiLampiran(t)
	id := sampaiAdjustment(u)
	u.a.Nama = map[string]string{teknik: "UJI NAMA TEKNIK"}
	var mode map[string]string
	langkah := func(aksi string, m map[string]string) map[string]any {
		t.Helper()
		kode, out := u.minta(http.MethodPost, handlers.Prefix+"/kasus/"+id+"/aksi", teknik, models.WorkbasketAcceptation,
			services.PermintaanAksi{Aksi: aksi, Indeks: 1, Masukan: m, Mode: mode})
		u.wajib(kode, http.StatusOK, out, aksi)
		mode = map[string]string{}
		if mm, ok := out["mode"].(map[string]any); ok {
			for k, v := range mm {
				mode[k] = v.(string)
			}
		}
		return out
	}
	adj := func(p string) string { return models.JalurAnak(models.DaftarAdjustment, 1, p) }
	langkah("SetNameCurrency", map[string]string{adj("CurrencyID"): matauang})
	langkah("CountGrossAdjTreaty", map[string]string{adj("Type"): "1", adj("GrossAdjustment"): "100"})
	for _, k := range []string{"LOD", "DLA", "SPGR"} {
		kode, out := u.unggah(id, teknik, models.WorkbasketAcceptation, k, map[string]string{"UJI-" + k + ".pdf": "%PDF"})
		u.wajib(kode, http.StatusOK, out, "unggah "+k)
	}
	komite := func(out map[string]any) []models.Tata {
		t.Helper()
		b, _ := json.Marshal(out["modal"].(map[string]any)["komite:1"])
		var ts []models.Tata
		if err := json.Unmarshal(b, &ts); err != nil {
			t.Fatal(err)
		}
		return ts
	}
	nilai := func(out map[string]any, j string) any {
		return out["halaman"].(map[string]any)["nilai"].(map[string]any)[j]
	}
	kepala := func(out map[string]any, apa string) {
		t.Helper()
		if d, p := nilai(out, "TempCommiteClaim.DateOfComitee"), nilai(out, "TreatyExchangeYearly.UserName"); d != "2026-05-20" ||
			p != "UJI NAMA TEKNIK" {
			t.Fatalf("%s: Date %v / PIC %v, mau 2026-05-20 / UJI NAMA TEKNIK (nama tampilan, bukan akun)", apa, d, p)
		}
	}

	out := langkah("BukaKomite", map[string]string{models.CD + "Payable": "2", adj("NameOfBank"): "UJI BANK",
		adj("NoAccount"): "UJI-REK"})
	if nilai(out, "Protect.CARI1") != "1" || nilai(out, "Protect.CARI2") != "1" {
		t.Fatalf("proteksi harus lolos: CARI1 %v CARI2 %v pesan %v", nilai(out, "Protect.CARI1"),
			nilai(out, "Protect.CARI2"), out["pesan"])
	}
	kepala(out, "popup dibuka")
	ts := komite(out)
	if cari(ts, func(x models.Tata) bool { return x.Label == "PIC" && x.Jalur == "TreatyExchangeYearly.UserName" }) == nil {
		t.Fatal("medan TreatyExchangeYearly.UserName harus berlabel PIC")
	}
	if cari(ts, func(x models.Tata) bool { return x.Label == "Initial" }) != nil {
		t.Fatal("label Initial diganti PIC")
	}
	for _, p := range []string{"DataCommitteeTreaty.Remarks", "DataCommitteeTreaty.CircumCauseOfLoss"} {
		m := cari(ts, func(x models.Tata) bool { return x.Jalur == adj(p) })
		if m == nil || m.Aksi != models.AksiSegarKomite {
			t.Fatalf("%s harus ber-refresh (event change, Refresh this section): %+v", p, m)
		}
	}
	if adaTombolKirim(ts) {
		t.Fatal("Remarks / Circumstanses kosong: Send Claim to Committee belum tampil")
	}

	out = langkah(models.AksiSegarKomite, map[string]string{adj("DataCommitteeTreaty.CircumCauseOfLoss"): "UJI",
		adj("DataCommitteeTreaty.Remarks"): "UJI"})
	if !adaTombolKirim(komite(out)) {
		t.Fatalf("sesudah Remarks / Circumstanses diisi, Send Claim to Committee harus tampil (mode %v)", mode)
	}
	kepala(out, "sesudah refresh")
}

func cari(ts []models.Tata, f func(models.Tata) bool) *models.Tata {
	for i := range ts {
		if f(ts[i]) {
			return &ts[i]
		}
		if x := cari(ts[i].Anak, f); x != nil {
			return x
		}
	}
	return nil
}

func adaTombolKirim(ts []models.Tata) bool {
	return cari(ts, func(x models.Tata) bool { return x.ID == "SendClaimToCommittee" }) != nil
}
