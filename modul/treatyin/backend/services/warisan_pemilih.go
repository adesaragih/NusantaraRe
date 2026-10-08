package services

// Isi pemilih "Choose Ceding" dan "Choose Source of Business" — keputusan
// pemilik proses 4 Oktober 2026.
//
// Kedua tombol itu MATI sejak layar ini lahir, dengan alasan tertulis:
// `ERD.md` §2.8 menempatkan `CEDANT` dan `ASAL_BISNIS` di luar skema modul
// ini, jadi "tanpa modul pemiliknya tidak ada yang dapat dipilih". Separuh
// alasan itu keliru, dan diukur: tabel masternya memang tidak ada, tetapi
// NILAINYA ada — 131 pasangan pengenal+nama cedant dan 96 asal bisnis,
// terisi pada 1.853 dari 1.854 kontrak. Menolak menampilkannya berarti
// menyembunyikan data yang ada.

import (
	"context"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/models"
)

// DaftarCedant mengembalikan isi pemilih "Choose Ceding".
func (l *Layanan) DaftarCedant(ctx context.Context, p inti.Pelaku) ([]models.PilihanWarisan, error) {
	return l.pilihanWarisan(ctx, p, l.gudang.BacaDaftarCedant)
}

// DaftarAsalBisnis mengembalikan isi pemilih "Choose Source of Business".
func (l *Layanan) DaftarAsalBisnis(ctx context.Context, p inti.Pelaku) ([]models.PilihanWarisan, error) {
	return l.pilihanWarisan(ctx, p, l.gudang.BacaDaftarAsalBisnis)
}

// DaftarJenisTreaty mengembalikan isi dropdown `Treaty Type` tab Limits
// (`LimitProportional.xml` → `BrowseReinsuranceType_RD`). Nama kembar
// ditandai seperti pemilih lain: terukur 84 pilihan aktif, 83 nama unik.
func (l *Layanan) DaftarJenisTreaty(ctx context.Context, p inti.Pelaku) ([]models.PilihanWarisan, error) {
	return l.pilihanWarisan(ctx, p, l.gudang.BacaDaftarJenisTreaty)
}

// OpsiLimits mengembalikan isi SELURUH dropdown tab Limits sekaligus —
// satu permintaan untuk satu tab.
func (l *Layanan) OpsiLimits(ctx context.Context, p inti.Pelaku) (models.OpsiLimits, error) {
	var o models.OpsiLimits
	var err error
	if o.JenisTreaty, err = l.pilihanWarisan(ctx, p, l.gudang.BacaDaftarJenisTreaty); err != nil {
		return models.OpsiLimits{}, err
	}
	if o.KelompokTreaty, err = l.pilihanWarisan(ctx, p, l.gudang.BacaDaftarKelompokTreaty); err != nil {
		return models.OpsiLimits{}, err
	}
	if o.MataUang, err = l.pilihanWarisan(ctx, p, l.gudang.BacaDaftarMataUangLimit); err != nil {
		return models.OpsiLimits{}, err
	}
	// ⭐ LABEL dari PROMPT VALUE rule Property, bukan standard value —
	// permintaan pemilik proses 7 Oktober 2026. Yang rule-nya belum
	// diekspor jatuh ke nilainya sendiri; lihat `prompt_value.go`.
	o.JenisLayer = opsiDariProperty("LayerType", domainJenisLayer)
	o.Cover = opsiDariProperty("Cover", domainCover)
	o.RelasiMataUang = opsiDariProperty("CurrencyRelation", domainRelasiMataUang)
	o.CatatanReinstatement = opsiDariProperty("ReinstatementNote", domainCatatanReinstatement)
	return o, nil
}

// Domain dropdown `associated` tab Limits Non-Prop — TERUKUR atas seluruh
// `T_TREATY_LIMITS` 6 Oktober 2026 (6.262 layer). Label tampilnya hidup di
// rule Property yang tidak diekspor, jadi nilainya tampil apa adanya.
// `CurrencyRelation` dibaca Activity sebagai `"AND"` / `"OR"`
// (`DetailCalculationROL`).
var (
	domainJenisLayer           = []string{"layer", "sublayer"}      // LayerType / LayerPartType: 4.457 / 323
	domainCover                = []string{"risk", "cat", "riskcat"} // 1.089 / 1.124 / 2.564
	domainRelasiMataUang       = []string{"OR", "AND"}              // 3.342 / 138
	domainCatatanReinstatement = []string{"asamount", "astime"}     // 1.456 / 16
)

// DaftarKelasBisnisTreaty mengembalikan autocomplete `Class of Business`
// untuk satu Treaty Group.
func (l *Layanan) DaftarKelasBisnisTreaty(ctx context.Context, p inti.Pelaku, treatyGroupID string) ([]models.PilihanWarisan, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return nil, err
	}
	if strings.TrimSpace(treatyGroupID) == "" {
		// Grup belum dipilih — parameter RD-nya kosong, nol baris.
		return []models.PilihanWarisan{}, nil
	}
	return l.gudang.BacaDaftarKelasBisnisTreaty(ctx, treatyGroupID)
}

// pilihanWarisan membaca lalu MENANDAI nama yang dipakai lebih dari satu
// pengenal.
//
// ⛔ Penandaan dikerjakan DI SINI, bukan di repository: ia bukan apa yang
// tersimpan, melainkan kesimpulan atas seluruh himpunan. Repository
// mengembalikan apa adanya; yang menafsirkan lapisan ini.
//
// ⚠️ Nama kembar TIDAK disatukan dan TIDAK dibuang. Terukur pada 4 Oktober
// 2026: **36 nama cedant** dipakai lebih dari satu pengenal.
// `REASURANSI MAIPARK INDONESIA` punya 3; di antara sisanya `ASURANSI ADIRA DINAMIKA`, yang satu
// pengenalnya (`ASM-SFAGIS-WORK-ORG ORG-34`) adalah ID kerja Pega yang
// menyelinap menjadi data. Menyatukannya berarti layar diam-diam memilih
// pengenal mana yang ditulis, dan kontrak nyangkut ke cedant yang keliru
// tanpa seorang pun melihatnya sampai rekonsiliasi.
func (l *Layanan) pilihanWarisan(
	ctx context.Context,
	p inti.Pelaku,
	baca func(context.Context) ([]models.PilihanWarisan, error),
) ([]models.PilihanWarisan, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return nil, err
	}
	daftar, err := baca(ctx)
	if err != nil {
		return nil, err
	}
	cacah := map[string]int{}
	for _, b := range daftar {
		cacah[b.Nama]++
	}
	for i := range daftar {
		daftar[i].Kembar = cacah[daftar[i].Nama] > 1
	}
	return daftar, nil
}
