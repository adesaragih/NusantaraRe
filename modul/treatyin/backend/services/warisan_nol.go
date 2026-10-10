package services

// Satu tempat yang menjamin: TIDAK ADA larik `nil` yang meninggalkan jalur
// baca kontrak warisan.
//
// ---------------------------------------------------------------------
// ⛔ SEBABNYA, dan ia galat nyata yang dilaporkan pemakai 6 Oktober 2026
// ---------------------------------------------------------------------
//   Tekan `Edit`, lalu `View` -> layar berhenti dengan
//
//       Cannot read properties of null (reading 'length')
//
//   Di Go, iris `nil` dan iris KOSONG adalah dua hal yang berbeda ketika
//   dijadikan JSON:
//
//       var a []T            -> `null`
//       a := []T{}           -> `[]`
//
//   `repository.BacaKontrakWarisan` mengisi `Kurs` dengan `append` ke atas
//   medan yang belum pernah disentuh. Kontrak yang dokumennya TIDAK punya
//   `CurrencyList` - dan kontrak yang tidak punya dokumen sama sekali -
//   karena itu mengirim `"kurs": null`, dan `kurs.length` di layar
//   menghentikan seluruh halaman, bukan hanya satu grid.
//
// ---------------------------------------------------------------------
// ⛔ MENGAPA SATU FUNGSI, BUKAN TAMBALAN DI TEMPAT
// ---------------------------------------------------------------------
//   `LimitsPohon` dan `AdaDiJSON` sudah pernah ditambal satu per satu di
//   ekor `BacaKontrakWarisan`, masing-masing sesudah gejalanya terlihat.
//   `Kurs` adalah yang ketiga, dan pola "tambal yang kebetulan ketahuan"
//   menjamin ada yang keempat.
//
//   Jadi penjaganya dipindahkan ke SATU fungsi, dan kelengkapannya
//   dibuktikan `warisan_nol_test.go` - yang menyusuri medan `KontrakWarisan`
//   dengan refleksi dan GAGAL begitu seseorang menambah larik baru tanpa
//   menyebutnya di sini. Penjaga yang hanya menyebut medan yang sudah ada
//   tidak menjaga medan yang belum ditulis.

import "nusantarare/modul/treatyin/backend/models"

// nolkanLarik mengganti setiap larik `nil` pada satu kontrak warisan dengan
// larik KOSONG, dan setiap peta `nil` dengan peta kosong.
//
// ⚠️ Yang diubah hanya BENTUK KOSONGNYA. Larik yang sudah berisi tidak
// disentuh, dan nol nilai dikarang: `[]` berarti "nol baris", sama persis
// dengan arti `null` sebelumnya - bedanya hanya layar dapat membacanya.
func nolkanLarik(k *models.KontrakWarisan) {
	if k.Kurs == nil {
		k.Kurs = []models.BarisKursWarisan{}
	}
	if k.PeriodePelaporan == nil {
		k.PeriodePelaporan = []models.BarisPeriodeWarisan{}
	}
	if k.Portofolio == nil {
		k.Portofolio = []models.BarisPortofolioWarisan{}
	}
	if k.Akumulasi == nil {
		k.Akumulasi = []models.BarisAkumulasiWarisan{}
	}
	if k.Egnpi == nil {
		k.Egnpi = []models.BarisEgnpiWarisan{}
	}
	if k.Retensi == nil {
		k.Retensi = []models.BarisRetensiWarisan{}
	}
	if k.Angsuran == nil {
		k.Angsuran = []models.BarisAngsuranWarisan{}
	}
	if k.Catatan == nil {
		k.Catatan = []models.BarisCatatanWarisan{}
	}
	if k.Layer == nil {
		k.Layer = []models.BarisLayerWarisan{}
	}
	if k.LimitsPohon == nil {
		k.LimitsPohon = []map[string]any{}
	}
	if k.TotalRetensi == nil {
		k.TotalRetensi = []models.BarisTotalRetensiWarisan{}
	}
	// Larik akar tab Limits Non-Prop — `Summary of Limit` dan keempat
	// `TotalLimit…NP`, yang dibaca layar dengan `.length`.
	if k.LimitsAkar.LimitSummaryList == nil {
		k.LimitsAkar.LimitSummaryList = []map[string]any{}
	}
	if k.LimitsAkar.Total == nil {
		k.LimitsAkar.Total = map[string][]map[string]any{}
	}
	for _, nama := range []string{"TotalLimitIOONP", "TotalLimitDeductblNP", "TotalLimitPremiEarnNP", "TotalLimitMDPNP"} {
		if k.LimitsAkar.Total[nama] == nil {
			k.LimitsAkar.Total[nama] = []map[string]any{}
		}
	}
	// Tab Share Non-Prop — larik dan kesembilan kunci Total.
	lengkapiShare(&k.ShareNP)
	if k.PolisProduksi == nil {
		k.PolisProduksi = []models.BarisPolisProduksi{}
	}
	if k.SkalaKoasuransi == nil {
		k.SkalaKoasuransi = []models.BarisSkalaKoasuransiWarisan{}
	}
	if k.Lampiran == nil {
		k.Lampiran = []models.BarisLampiranWarisan{}
	}
	if k.KategoriLampiran == nil {
		k.KategoriLampiran = []models.BarisKategoriLampiran{}
	}
	// ⛔ Larik DI DALAM medan bersarang ikut dijaga. `ejaanLain` dibaca
	// `TabTeksPanjang` dengan `.length` yang sama, dan kedua tab teks dapat
	// berbentuk nol-nilai pada kontrak yang dokumennya tidak ada.
	if k.Pengecualian.EjaanLain == nil {
		k.Pengecualian.EjaanLain = []string{}
	}
	if k.SyaratKhusus.EjaanLain == nil {
		k.SyaratKhusus.EjaanLain = []string{}
	}
	if k.OpsiKepala.Bordereaux == nil {
		k.OpsiKepala.Bordereaux = []models.Opsi{}
	}
	if k.OpsiKepala.CaraPembukuan == nil {
		k.OpsiKepala.CaraPembukuan = []models.Opsi{}
	}
	if k.OpsiKepala.CaraPembukuanNonProp == nil {
		k.OpsiKepala.CaraPembukuanNonProp = []models.Opsi{}
	}
	if k.OpsiKepala.PeriodePelaporan == nil {
		k.OpsiKepala.PeriodePelaporan = []models.Opsi{}
	}
	// ⛔ Larik DI DALAM elemen larik ikut dijaga. `PohonLimits` membaca
	// `b.kelasBisnis.length` per baris layer; satu baris yang kelas
	// bisnisnya `null` menghentikan halaman sama seperti larik akarnya.
	for i := range k.Layer {
		if k.Layer[i].KelasBisnis == nil {
			k.Layer[i].KelasBisnis = []string{}
		}
	}
	// Peta - `adaDiJson[kunci]` atas `null` memberi galat yang sama.
	if k.AdaDiJSON == nil {
		k.AdaDiJSON = map[string]bool{}
	}
	if k.TeksMentah == nil {
		k.TeksMentah = map[string]string{}
	}
	if k.Penampung == nil {
		k.Penampung = map[string]string{}
	}
	if k.PenampungLarik == nil {
		k.PenampungLarik = map[string][]map[string]any{}
	}
}
