// ⛔⛔ BERKAS INI TIDAK LAGI DIPAKAI JALUR BACA, sejak 6 Oktober 2026.
//
// Keputusan pemilik proses melarang keras menarik nilai dari `JSONDATA`.
// Keempat tab kini dibaca dari tabel pendaratan lewat `pendaratan_layer.go`;
// NOL kode produksi memanggil fungsi di berkas ini, hanya ujinya sendiri.
//
// ⚠️ Ia DITINGGALKAN satu ronde dengan sengaja, bukan karena terlupa:
// ketiga belas tabel pendaratan masih NOL BARIS, sehingga jalur barunya
// belum pernah terbukti di layar. Membuang jalur lama sebelum penggantinya
// terbukti berarti membuang satu-satunya pembanding ketika hasilnya
// berselisih.
//
// ⭐ Begitu tabelnya terisi dan layarnya terbukti, berkas ini beserta
// ujinya DIBUANG — dua sumber untuk satu tab berarti salah satunya akan
// basi tanpa suara, dan itu doktrin modul ini sendiri.

package repository

// Baris LAYER dari `M_TREATY_IN.JSONDATA` — empat tab, satu sumber.
//
// ⛔⛔ `M_TREATY_IN2` DICABUT sebagai sumber, 5 Oktober 2026. Keputusan
// pemilik proses, dinyatakan tiga kali. Berkas ini yang menggantikan
// `warisan_in2.go`, dan sesudahnya **nol kueri di modul ini menyebut tabel
// itu** — dijaga `TestNolKueriMTreatyIn2`.
//
// ⚠️ Tabelnya TETAP ADA di Oracle dan TETAP terlarang disentuh. Yang dicabut
// jalur bacanya, bukan tabelnya; `TestWarisanHanyaDibaca` tetap menamainya.
//
// # Bentuk pohonnya — DIUKUR sebelum dibangun
//
// Sapuan seluruh 1.854 dokumen, 5 Oktober 2026, atas setiap jalur bersarang
// di bawah `Limits[]`:
//
//	Limits[]                      ← satu elemen per LAYER        (4.210 elemen)
//	  .Layer .LayerType .Cover .Currency .Currency2 .Limit .Limit2
//	  .Deductible .Deductible2 .AgregateLimit .AgregateLimit2
//	  .AdjRate .ROLPct .MDPPct .CurrencyRelation .TreatyType
//	  .MDPList[] .PremiumEarnedList[] .EgnpiTotalList[] .MDPMinList[]
//	  .Reinstatement_List[] .TreatyGroupList[] .LayerList[]
//	  └ .Detail[]                 ← satu elemen per TREATY GROUP  (5.716 silang)
//	      .TreatyGroup .TreatyGroupID .TreatyType .CessionPct .QSPct
//	      .Brokerage .RNMShare .RIOGR .Earthquake .FloodJab .FloodNation
//	      .RSMDLimit .SpreadingType .SpreadingTotalPct .ProfitCommision
//	      .IOOLimitList[] .RetentionList[] .CessionList[] .EPIList[]
//	      .PLAList[] .CashLossList[] .ClaimCoopList[] .DeductionList[]
//	      .ReserveList[] .RNMShareList[] .RNMSpreadedList[]
//	      .RNMSpreadedListRI[] .SpreadingList[] .AchievementLists[]
//	      └ .COBList[]            ← satu elemen per CLASS OF BUSINESS
//	          .ClassOfBusiness .ClassOfBusinessID .TreatyGroup .TreatyGroupID
//
// ⭐ POHON TIGA TINGKAT ITU ADA, dan di sinilah ia — bukan di kunci puncak.
// Ronde sebelumnya menyatakan `Kind of Treaty`/`Treaty Type`/
// `Class of Business` nol jejaknya; yang benar adalah nol di PUNCAK. Mereka
// ada di `Limits[].Detail[]` dan `Limits[].Detail[].COBList[]`.
//
// # Pemetaan kolom → jalur dokumen: DIBUKTIKAN, bukan ditebak
//
// Nilai ke-32 kolom `M_TREATY_IN2` diadu satu per satu dengan nilai di jalur
// dokumen, atas 1.340 kontrak yang ada di KEDUA sumber, disejajarkan menurut
// NILAI `Layer` (bukan indeks baris):
//
//	LAYER            ← .Layer                       100,0%
//	MDP_RATIO        ← .MDPPct                      100,0%
//	CURRENCYRELATION ← .CurrencyRelation             99,8%
//	BASIS_COVER      ← .Cover                        98,6%
//	LAYERTYPE        ← .LayerType                    91,6%
//	ADJ_RATE         ← .AdjRate                      90,3%
//	ROL              ← .ROLPct                       87,7%
//	TREATYTYPE       ← .TreatyType                   81,2%
//	EARN_PREMIUM     ← .PremiumEarnedList[].Value    80,0%
//	CEDANT_RETENTION ← .Deductible                   70,6%
//
// ⭐ Baris terakhir menutup pertanyaan lama: `PEMETAAN-M-TREATY-IN2.md`
// menduga `CEDANT_RETENTION` ↔ `.Deductible` dengan keyakinan **51%**, dan
// pengadu nilai memberi **70,6%** — ditambah gambar `30` dokumen desain yang
// memberi kolom `Deductible ( IDR )` dengan nama itu di layar. Dugaan
// menjadi padanan.
//
// ⚠️ KOLOM BER-ANGKA RENDAH (12–45%) SELURUHNYA di tingkat `Detail[]`, dan
// itu menjelaskan dirinya sendiri: tabel memipihkan LAYER × TREATY GROUP
// menjadi satu baris, sementara pengadu hanya membaca `Detail[0]`. Cocok
// ketika barisnya kebetulan treaty group pertama, meleset selainnya. Jadi
// angka rendah itu BUKAN tanda padanan yang salah — ia tanda bahwa tabelnya
// pipih dan dokumennya bersarang.
//
// # Yang pencabutan ini selesaikan
//
//	dokumen dengan `Limits[]` berisi        1.850 kontrak
//	di antaranya ada di `M_TREATY_IN2`      1.340
//	LUBANG                                    510 kontrak · 1.210 elemen
//
// Dengan dokumen sebagai sumber, jangkauannya 1.340 → **1.850**, dan kosong
// kembali punya SATU arti. Petunjuk kosong yang menjelaskan lubang 510 itu
// dicabut bersama ini.

import (
	"encoding/json"
	"strings"

	"nusantarare/modul/treatyin/backend/models"
)

// jsonLimit - satu elemen `Limits[]`, yaitu satu LAYER.
//
// ⚠️ Seluruh medan skalarnya `json.RawMessage`, bukan `*string`. Sapuan
// menemukan tipe yang BERCAMPUR di jalur yang sama - `.Layer` muncul sebagai
// angka JSON pada sebagian dokumen dan sebagai string pada sebagian lain,
// dan `.IsCombineMDP` sebagai boolean. `*string` menolak keduanya dan
// menggugurkan seluruh dokumen; `RawMessage` menerima apa adanya dan
// `teksJSON()` yang memutuskan.
type jsonLimit struct {
	Layer      json.RawMessage `json:"Layer"`
	LayerType  json.RawMessage `json:"LayerType"`
	Cover      json.RawMessage `json:"Cover"`
	TreatyType json.RawMessage `json:"TreatyType"`
	Currency   json.RawMessage `json:"Currency"`
	Currency2  json.RawMessage `json:"Currency2"`
	Limit      json.RawMessage `json:"Limit"`
	Deductible json.RawMessage `json:"Deductible"`
	// ⭐ Pasangan MATA UANG KEDUA — ada di dokumen sejak awal dan tidak
	// pernah dibaca sampai 6 Oktober 2026. Terukur: `Limit2` 2.441
	// kemunculan, `Deductible2` 2.436.
	Limit2           json.RawMessage `json:"Limit2"`
	Deductible2      json.RawMessage `json:"Deductible2"`
	AdjRate          json.RawMessage `json:"AdjRate"`
	ROLPct           json.RawMessage `json:"ROLPct"`
	MDPPct           json.RawMessage `json:"MDPPct"`
	CurrencyRelation json.RawMessage `json:"CurrencyRelation"`

	MDPList           []jsonNilaiMataUang `json:"MDPList"`
	PremiumEarnedList []jsonNilaiMataUang `json:"PremiumEarnedList"`

	Detail []jsonLimitDetail `json:"Detail"`
}

// jsonLimitDetail - satu elemen `Limits[].Detail[]`, yaitu satu TREATY GROUP
// di dalam layer.
type jsonLimitDetail struct {
	TreatyGroup json.RawMessage `json:"TreatyGroup"`
	TreatyType  json.RawMessage `json:"TreatyType"`
	CessionPct  json.RawMessage `json:"CessionPct"`
	Brokerage   json.RawMessage `json:"Brokerage"`
	RNMShare    json.RawMessage `json:"RNMShare"`
	RIOGR       json.RawMessage `json:"RIOGR"`
	Earthquake  json.RawMessage `json:"Earthquake"`
	// ⭐ Tiga batas yang `M_TREATY_IN2` tidak punya — gambar 28.
	RSMDLimit          json.RawMessage `json:"RSMDLimit"`
	FloodJab           json.RawMessage `json:"FloodJab"`
	FloodNation        json.RawMessage `json:"FloodNation"`
	CurrencyRSMD       json.RawMessage `json:"CurrencyRSMD"`
	CurrencyEarthquake json.RawMessage `json:"CurrencyEarthquake"`
	CurrencyFloodJab   json.RawMessage `json:"CurrencyFloodJab"`
	CurrencyFloodNat   json.RawMessage `json:"CurrencyFloodNat"`
	SpreadingType      json.RawMessage `json:"SpreadingType"`

	CessionList       []jsonNilaiMataUang `json:"CessionList"`
	EPIList           []jsonNilaiMataUang `json:"EPIList"`
	RNMShareList      []jsonNilaiMataUang `json:"RNMShareList"`
	RNMSpreadedList   []jsonNilaiMataUang `json:"RNMSpreadedList"`
	RNMSpreadedListRI []jsonNilaiMataUang `json:"RNMSpreadedListRI"`
	SpreadingList     []jsonSpreading     `json:"SpreadingList"`
	COBList           []jsonCOB           `json:"COBList"`
}

// jsonNilaiMataUang - bentuk yang BERULANG di belasan tempat: satu nilai
// beserta mata uangnya.
type jsonNilaiMataUang struct {
	Currency json.RawMessage `json:"Currency"`
	Value    json.RawMessage `json:"Value"`
}

// jsonSpreading - satu baris `SpreadingList[]`; `ReinsTypeName` yang
// memisahkan `QS (OR)` dari `QS (R/I)`.
type jsonSpreading struct {
	ReinsTypeName json.RawMessage `json:"ReinsTypeName"`
	Pct           json.RawMessage `json:"Pct"`
	Value         json.RawMessage `json:"Value"`
}

// jsonCOB - satu kelas bisnis di dalam satu treaty group.
type jsonCOB struct {
	ClassOfBusiness json.RawMessage `json:"ClassOfBusiness"`
	TreatyGroup     json.RawMessage `json:"TreatyGroup"`
}

// teksJSON mengubah nilai JSON apa pun menjadi teks tampil.
//
// ⛔ Angka JSON dikembalikan APA ADANYA dari sumbernya - `json.Number`, bukan
// `float64`. Melewatkan `9007199254740993` melalui `float64` mengembalikan
// `9007199254740992`, dan digit yang hilang itu tidak akan pernah terlihat
// sebagai galat; ia hanya menjadi angka yang salah di layar.
func teksJSON(r json.RawMessage) string {
	s := strings.TrimSpace(string(r))
	if s == "" || s == "null" {
		return ""
	}
	// Teks JSON: buang tanda kutipnya lewat pengurai, supaya escape (`\n`,
	// `\"`, `\uXXXX`) ikut terurai alih-alih tampil mentah.
	if s[0] == '"' {
		var v string
		if json.Unmarshal(r, &v) == nil {
			return strings.TrimSpace(v)
		}
	}
	// Angka dan boolean: apa adanya.
	return s
}

// nilaiKe mengambil `.Value` elemen ke-`i`, atau kosong bila tidak ada.
//
// ⛔⛔ PENGGANTI `nilaiPertama`, yang hanya membaca elemen [0] dan karena itu
// MENYEMBUNYIKAN mata uang kedua pada 106 layer (`MDPList`) dan 100 layer
// (`PremiumEarnedList`). Cacat itu ditemukan lewat pengukuran, bukan lewat
// keluhan — dan uji positifnya memakai salah satu layer itu supaya ia tidak
// dapat kembali diam-diam.
//
// ⚠️ Layer bermata uang TUNGGAL lulus bahkan dengan cacatnya; itu sebabnya
// uji positifnya WAJIB memakai layer bermata uang DUA.
//
// ⛔ Lariknya TIDAK dijumlahkan. Terukur: setiap larik berpanjang dua
// bermata uang BERBEDA, jadi menjumlahkannya memberi satu bilangan yang
// bukan uang apa pun.
func nilaiKe(l []jsonNilaiMataUang, i int) string {
	if i < 0 || i >= len(l) {
		return ""
	}
	return teksJSON(l[i].Value)
}

// pctSpreading mencari persen satu jenis reasuransi di dalam `SpreadingList`.
//
// ⛔ Dicocokkan menurut `ReinsTypeName`, bukan menurut urutan. Sapuan
// menemukan `QS (OR)` dan `QS (R/I)` dalam urutan yang BERBEDA antar
// dokumen; membaca elemen ke-0 sebagai OR akan menukar keduanya diam-diam.
func pctSpreading(l []jsonSpreading, nama string) string {
	for _, s := range l {
		if strings.EqualFold(strings.TrimSpace(teksJSON(s.ReinsTypeName)), nama) {
			return teksJSON(s.Pct)
		}
	}
	return ""
}

// LayerDariDokumen mengubah `Limits[]` menjadi baris layer siap tampil.
//
// ⛔ SATU baris per (LAYER × TREATY GROUP), bukan per layer. Itu bentuk yang
// keempat tab harapkan, dan ia sama dengan cara `M_TREATY_IN2` memipihkannya
// - 5.716 silang terukur di dokumen lawan 7.281 baris di tabel.
//
// ⚠️ Layer TANPA `Detail[]` tetap memberi SATU baris. Terukur: sebagian
// dokumen punya elemen `Limits[]` berisi medan layer tetapi nol treaty
// group, dan menjatuhkannya akan menyembunyikan layer yang di Pega terlihat.
func LayerDariDokumen(lim []jsonLimit) []models.BarisLayerWarisan {
	// ⛔ Irisan KOSONG, bukan nil - layar membedakan "nol baris" dari
	// "belum dibaca", dan `null` di JSON jawaban merender berbeda.
	hasil := []models.BarisLayerWarisan{}
	for i := range lim {
		l := &lim[i]
		dasar := models.BarisLayerWarisan{
			Layer:              teksJSON(l.Layer),
			JenisLayer:         teksJSON(l.LayerType),
			DasarCover:         teksJSON(l.Cover),
			JenisTreaty:        teksJSON(l.TreatyType),
			MataUang:           teksJSON(l.Currency),
			MataUangLimit:      teksJSON(l.Currency2),
			Limit100:           teksJSON(l.Limit),
			RetensiCedant:      teksJSON(l.Deductible),
			AdjRate:            teksJSON(l.AdjRate),
			ROL:                teksJSON(l.ROLPct),
			RasioMDP:           teksJSON(l.MDPPct),
			RelasiMataUang:     teksJSON(l.CurrencyRelation),
			Limit100Kedua:      teksJSON(l.Limit2),
			RetensiCedantKedua: teksJSON(l.Deductible2),
			MDP:                nilaiKe(l.MDPList, 0),
			MDPKedua:           nilaiKe(l.MDPList, 1),
			PremiEarned:        nilaiKe(l.PremiumEarnedList, 0),
			PremiEarnedKedua:   nilaiKe(l.PremiumEarnedList, 1),
		}
		if len(l.Detail) == 0 {
			dasar.KelasBisnis = []string{}
			hasil = append(hasil, dasar)
			continue
		}
		for j := range l.Detail {
			d := &l.Detail[j]
			b := dasar
			b.KelompokTreaty = teksJSON(d.TreatyGroup)
			// ⛔ Irisan KOSONG, bukan nil — grid `Class of Business` yang
			// kosong berbunyi `No items`, dan `null` merender berbeda.
			b.KelasBisnis = []string{}
			for _, c := range d.COBList {
				if v := teksJSON(c.ClassOfBusiness); v != "" {
					b.KelasBisnis = append(b.KelasBisnis, v)
				}
			}
			// `TreatyType` ada di DUA tingkat; yang di `Detail[]` lebih
			// khusus, jadi ia menang bila terisi.
			if v := teksJSON(d.TreatyType); v != "" {
				b.JenisTreaty = v
			}
			b.PersenCession = teksJSON(d.CessionPct)
			b.PersenBrokerage = teksJSON(d.Brokerage)
			b.RNMShare = teksJSON(d.RNMShare)
			b.RIOGR = teksJSON(d.RIOGR)
			b.Gempa = teksJSON(d.Earthquake)
			b.BatasRSMD = teksJSON(d.RSMDLimit)
			b.BatasBanjirJab = teksJSON(d.FloodJab)
			b.BatasBanjirNas = teksJSON(d.FloodNation)
			b.MataUangRSMD = teksJSON(d.CurrencyRSMD)
			b.MataUangGempa = teksJSON(d.CurrencyEarthquake)
			b.MataUangBanjirJab = teksJSON(d.CurrencyFloodJab)
			b.MataUangBanjirNas = teksJSON(d.CurrencyFloodNat)
			b.JenisPenyebaran = teksJSON(d.SpreadingType)
			b.CessionKeRI = nilaiKe(d.CessionList, 0)
			b.EPI100 = nilaiKe(d.EPIList, 0)
			b.LiabilityRNM = nilaiKe(d.RNMShareList, 0)
			b.LiabilityQSOR = nilaiKe(d.RNMSpreadedList, 0)
			b.LiabilityQSRI = nilaiKe(d.RNMSpreadedListRI, 0)
			b.QSOR = pctSpreading(d.SpreadingList, "QS (OR)")
			b.QSRI = pctSpreading(d.SpreadingList, "QS (R/I)")
			hasil = append(hasil, b)
		}
	}
	return hasil
}
