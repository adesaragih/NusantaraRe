package handlers_test

// Rantai wilayah Input Register (perintah work owner 11-10-2026 "PERBAIKI DEPAN DAN BELAKANG. PERBAIKI LOGIC NYA. GIMANA
// SEHARUSNYA MEMILIH INI"): Country -> Province -> City -> District -> Region berurutan, daftar menurut NAMA induk
// terpilih, mengganti induk mengosongkan tingkat bawahnya, tingkat tanpa induk nonaktif.

import (
	"net/http"
	"reflect"
	"testing"

	"nusantarare/modul/claimfacin/backend/handlers"
	"nusantarare/modul/claimfacin/backend/models"
	"nusantarare/modul/claimfacin/backend/services"
)

func TestRantaiWilayahRegister(t *testing.T) {
	u := baruUji(t)
	u.a.Negara = []models.Pilihan{{Nilai: "INDONESIA", Label: "INDONESIA", Tambahan: map[string]string{"CountryID": "UJI-ID"}}}
	u.a.Provinsi = []models.Pilihan{{Nilai: "UJI PROV A", Label: "UJI PROV A"}, {Nilai: "UJI PROV B", Label: "UJI PROV B"}}
	u.a.Kota = []models.Pilihan{{Nilai: "UJI KOTA", Label: "UJI KOTA"}}
	u.a.Distrik = []models.Pilihan{{Nilai: "UJI KEC", Label: "UJI KEC"}}
	u.a.RW = []models.Pilihan{{Nilai: "UJI KEL", Label: "UJI KEL", Tambahan: map[string]string{"RWID": "UJI-RW1",
		"PostalCode": "11111"}}}
	id := u.buat()
	kode0, out0 := u.minta(http.MethodGet, handlers.Prefix+"/kasus/"+id, admin, "", nil)
	u.wajib(kode0, http.StatusOK, out0, "buka kasus baru")
	if c, _ := out0["halaman"].(map[string]any)["nilai"].(map[string]any)[models.CD+"Country"].(string); c != models.NegaraBawaan {
		t.Fatalf("Country awal %q, mau %q (sel 924 pyDefaultValue)", c, models.NegaraBawaan)
	}
	// pilih meniru kiriman layar (`LayarKasus` -> `nilai.masukan`): SEMUA isian wilayah layar terakhir + nilai baru medan
	// yang dipilih, plus Param. Uji lama hanya mengirim Param sehingga luput: server menggabungkan kiriman SEBELUM aksi
	// (temuan work owner 11-10-2026, Province BALI + City JAKARTA TIMUR).
	layar := map[string]any{}
	jalurAksi := map[string]string{"PilihNegara": "Country", "PilihProvinsi": "Province", "PilihKota": "City",
		"PilihDistrik": "District", "PilihRW": "RW"}
	pilih := func(aksi, nilai string) map[string]any {
		t.Helper()
		kirim := map[string]string{}
		for _, p := range []string{"Country", "Province", "City", "District", "RW", "PostalCode"} {
			v, _ := layar[models.CD+p].(string)
			kirim[models.CD+p] = v
		}
		kirim[models.CD+jalurAksi[aksi]] = nilai
		kode, out := u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: aksi, Param: nilai, Masukan: kirim})
		u.wajib(kode, http.StatusOK, out, aksi+" "+nilai)
		layar = out["halaman"].(map[string]any)["nilai"].(map[string]any)
		return layar
	}
	nilai := func(n map[string]any, p string) string { s, _ := n[models.CD+p].(string); return s }
	pilih("PilihNegara", "INDONESIA")
	pilih("PilihProvinsi", "UJI PROV A")
	pilih("PilihKota", "UJI KOTA")
	pilih("PilihDistrik", "UJI KEC")
	if got, mau := u.a.IndukWilayah["distrik"], []string{"UJI PROV A", "UJI KOTA"}; !reflect.DeepEqual(got, mau) {
		t.Fatalf("daftar District disaring %v, mau nama provinsi + kota terpilih %v", got, mau)
	}
	n := pilih("PilihRW", "UJI KEL")
	if got, mau := u.a.IndukWilayah["rw"], []string{"UJI PROV A", "UJI KOTA", "UJI KEC"}; !reflect.DeepEqual(got, mau) {
		t.Fatalf("daftar Region disaring %v, mau %v", got, mau)
	}
	if nilai(n, "RW") != "UJI KEL" || nilai(n, "RWID") != "UJI-RW1" || nilai(n, "PostalCode") != "11111" {
		t.Fatalf("Region: %q %q %q", nilai(n, "RW"), nilai(n, "RWID"), nilai(n, "PostalCode"))
	}
	// kiriman layar lewat aksi lain (Save) tidak dapat mengganti wilayah yang sudah dipilih (hanya lewat aksi pemilihan)
	kodeS, outS := u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "Simpan", Masukan: map[string]string{
		models.CD + "Province": "UJI PALSU", models.CD + "City": "UJI KOTA LAIN", models.CD + "District": "UJI KEC LAIN"}})
	u.wajib(kodeS, http.StatusOK, outS, "Save berwilayah palsu")
	nS := outS["halaman"].(map[string]any)["nilai"].(map[string]any)
	if nilai(nS, "Province") != "UJI PROV A" || nilai(nS, "City") != "UJI KOTA" || nilai(nS, "District") != "UJI KEC" {
		t.Fatalf("Save menimpa wilayah: %q / %q / %q", nilai(nS, "Province"), nilai(nS, "City"), nilai(nS, "District"))
	}
	layar = nS
	// pilih ulang provinsi yang sama: tingkat bawah tetap
	n = pilih("PilihProvinsi", "UJI PROV A")
	if nilai(n, "City") != "UJI KOTA" || nilai(n, "PostalCode") != "11111" {
		t.Fatalf("provinsi sama dipilih ulang mengosongkan bawahnya: %v", n)
	}
	// ganti provinsi: City / District / Region / Zip Code kosong, lalu City nonaktif berpetunjuk sampai dipilih
	n = pilih("PilihProvinsi", "UJI PROV B")
	for _, p := range []string{"City", "CityID", "District", "DistrictID", "RW", "RWID", "PostalCode"} {
		if v := nilai(n, p); v != "" {
			t.Errorf("ganti provinsi: %s masih %q", p, v)
		}
	}
	if nilai(n, "Province") != "UJI PROV B" {
		t.Fatalf("Province %q", nilai(n, "Province"))
	}
	pilih("PilihKota", "UJI KOTA")
	n = pilih("PilihKota", "UJI KOTA")
	if nilai(n, "District") != "" {
		t.Fatal("District terisi tanpa dipilih")
	}
	// District kosong -> Region nonaktif: aksinya ditolak server
	kode, out := u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "PilihRW", Param: "UJI KEL",
		Masukan: map[string]string{models.CD + "RW": "UJI KEL"}})
	if kode == http.StatusOK {
		t.Fatalf("PilihRW tanpa District diterima: %v", out)
	}
	_, out = u.minta(http.MethodGet, handlers.Prefix+"/kasus/"+id, admin, "", nil)
	if out == nil {
		t.Fatal("buka kasus")
	}
}

func TestMedanWilayahTanpaIndukNonaktifBerpetunjuk(t *testing.T) {
	h := models.HalamanBaru()
	h.Setel(models.CD+"Country", "INDONESIA")
	h.Setel(models.CD+"Province", "UJI PROV")
	ts := models.Evaluasi(h, models.LayarRegister(), false)
	var kota, distrik []models.Tata
	var cari func([]models.Tata)
	cari = func(xs []models.Tata) {
		for _, x := range xs {
			if x.Jenis == models.JenisMedan && x.Jalur == models.CD+"City" {
				kota = append(kota, x)
			}
			if x.Jenis == models.JenisMedan && x.Jalur == models.CD+"District" {
				distrik = append(distrik, x)
			}
			cari(x.Anak)
		}
	}
	cari(ts)
	if len(kota) != 1 || kota[0].Nonaktif || kota[0].Aksi != "PilihKota" {
		t.Fatalf("City (Province terisi) harus aktif beraksi: %+v", kota)
	}
	if len(distrik) != 1 || !distrik[0].Nonaktif || distrik[0].Aksi != "" || distrik[0].Catatan != models.PetunjukWilayah("City") {
		t.Fatalf("District (City kosong) harus nonaktif berpetunjuk: %+v", distrik)
	}
}
