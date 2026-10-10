package services

// Rumus tab SHARE cabang NON-PROPORSIONAL — disalin dari Activity
// `D:\XML_NURE\Treaty In`, langkah demi langkah:
//
//	TreatyInNonAddItem(share)      Update Summary: baris Share / Facultative
//	                               Share disusun ulang SATU per layer Limits
//	TreatyInXOLAddSpreading        bagian OR / R/I tiap baris + Total OR / R/I Limit
//	FetchQSfromMasterXOL           spreading dari susunan treaty master
//	SetSpreadingXOL                spreading manual (Spreading Type kosong)
//	TreatyInSetBrokerage           Brokerage fee → Deduction → Net; Net OR / R/I
//	CalculateDeduction             Deduction Details (`HitungDeduksi`)
//	TreatyInNPSetTotal(share)      Total All Layers RNM Share
//	TreatyInSummaryLimitShare      Summarry of RNM Share
//	TreatyInSummaryLimitFacShare   ringkasan Facultative Share (tak tampil di tab ini)
//	TreatyInXOLAddSpreadingDetail  % RNM Share satu baris (panel)
//	AddSpreadingXOL                Add / Delete baris spreading manual
//	TreatyInShareListValue         Update Value in Share
//
// ⭐ DIUKUR terhadap data Pega yang tersimpan (7 Oktober 2026, seluruh
// `T_TREATY_SHARE*` sisi baru):
//   - baris Share ↔ layer Limits berpasangan menurut `URUTAN`: 3.892/3.892,
//     3.816 dengan layer identik;
//   - Gross Premium = MDP layer × `@divide(RNMShare,100,4)`: rasio Gross/MDP
//     SERAGAM per kontrak di 999 dari 1.001 kontrak, sama dengan
//     `RNMShare` baris di 311 dari 314 yang menyimpannya;
//   - Deduction `Brokerage fee` = `@divide(pct,100,4)` × Gross(1):
//     3.290/3.290;
//   - Net = Gross − Σ Deduction bermata uang sama: 4.090/4.123;
//   - `SpreadingTypeXOL` tersimpan ada di daftar RD
//     `ParentReinsMasterTrt` kontraknya: 3.751/3.788.
//
// ⚠️ SIMPANGAN, dan alasannya — semuanya JALUR, bukan rumus:
//  S1 `TreatyInNonAddItem` membuang larik Share lalu menyusunnya ulang, dan
//     `SpreadingTypeXOL` baris baru KOSONG; dropdown Spreading Type sendiri
//     hanya tampil bila nilainya TIDAK kosong. Teks ekspor itu tak dapat
//     menghasilkan data yang terukur (97% baris bernama). Di sini Spreading
//     Type baris lama DIBAWA ke baris baru berindeks sama.
//  S2 `TreatyInXOLAddSpreading` memanggil `FetchQSfromMasterXOL` dengan
//     induk PERTAMA dari RD (tanpa urutan) untuk SEMUA baris bernama. Di
//     sini dipakai induk milik baris itu sendiri — sama persis bila RD
//     mengembalikan satu induk (1.557 dari 3.892 baris).
//  S3 Larik bagian Net / Deduction OR-R/I di-`APPEND` Pega tanpa
//     dikosongkan di dalam satu rantai tombol (`FetchQSfromMasterXOL` [28]
//     `<LAST>`, `TreatyInSetBrokerage` [7], `…SpreadingDetail` [5.4]),
//     sehingga tiap klik menggandakan isinya dan totalnya. Di sini larik itu
//     DISUSUN ULANG per langkah. `GetNilaiTotal` (total menurut POSISI larik)
//     dihitung menurut MATA UANG — sama bila urutan mata uang tiap baris sama.
//  S4 `TreatyInSetBrokerage` [4] menambah `Overiding Commision` ke Facultative
//     Share tanpa membuang yang lama; di sini disusun ulang (S3).
//
// ⚠️ Yang DITIRU APA ADANYA walau janggal, sebab ia rumus Pega:
//   - `TreatyInNPSetTotal(share)`: `TotalShareGrossMinNP` tidak pernah
//     terisi. (⛔ RALAT 8 Oktober 2026: totalnya menjumlah larik SETIAP
//     baris Share — tangkapan Pega kontrak 1001855; lihat `NPSetTotalShare`.)
//   - `% RNM Share` satu baris (panel) memanggil `FetchQSfromMasterXOL`
//     TANPA induk — Spreading Type baris itu terhapus dan bagian OR-nya 0;
//   - `TreatyInSetBrokerage` memakai `QS (OR)` saja untuk Net OR (ORS tidak),
//     dan persennya terbawa dari baris sebelumnya;
//   - kontrak `1000951`: tambahan ORS 15% dan faktor 85% (kode keras Pega).

import (
	"context"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/repository"
)

// Aksi tab Share Non-Prop — satu per tombol / isian ekspor.
const (
	AksiShareRNM       = "rnm"       // % RNM Share
	AksiShareBrokerage = "brokerage" // % Brokerage
	// AksiShareSetBrokerage - `TreatyInSetBrokerage` SAJA. Kotak centang
	// Share Across The Board menjalankan `TreatyInXOLAddSpreading` LALU
	// `TreatyInSetBrokerage` sebagai dua langkah bersyarat terpisah
	// (`pyActionConditions`); layar Adjustment memanggil satu aksi per
	// langkah, jadi SetBrokerage tidak boleh menyeret XOLAddSpreading.
	AksiShareSetBrokerage    = "set-brokerage"
	AksiShareCentang         = "centang"          // Share Across The Board
	AksiShareFac             = "fac"              // Share to Other Retro · Brokerage From Other Retro
	AksiShareSummary         = "summary"          // Update Summary
	AksiShareTotal           = "total"            // Update Total
	AksiShareNilai           = "nilai-share"      // Update Value in Share
	AksiShareSpreadingType   = "spreading-type"   // panel: Spreading Type
	AksiShareRNMBaris        = "rnm-baris"        // panel: % RNM Share baris
	AksiShareDeduksi         = "deduksi"          // panel: Deduction Details
	AksiShareSpreadingTambah = "spreading-tambah" // panel: Add spreading manual
	AksiShareSpreadingHapus  = "spreading-hapus"  // panel: Delete spreading manual
	AksiShareSpreadingPct    = "spreading-pct"    // panel: Pct spreading manual
)

// Pesan Activity, apa adanya.
const (
	PesanShareNilaiKosong = "Value Cannot Be Empty"
	PesanShareSpreading   = "Total share must equal RNM share.!!"
)

// Kunci kesembilan grid "Total All Layers RNM Share", urutan layar.
var KunciTotalShareNP = []string{
	"TotalShareRnmNP", "TotalSpreadedRnmProp", "TotalSpreadedRnmRIProp",
	"TotalShareGrossMinNP", "TotalShareGrossNP", "TotalShareDeductionNP",
	"TotalShareNetNP", "TotalSpreadedNetPremi", "TotalSpreadedNetPremiRI",
}

const (
	namaQSOR        = "QS (OR)"
	namaORS         = "ORS"
	idKontrakKhusus = "1000951" // kode keras `FetchQSfromMasterXOL` / `TreatyInSetBrokerage`
	komentarBroker  = "Brokerage fee"
	komentarFac     = "Overiding Commision" // ejaan Pega
	benar           = "true"
)

// SumberSpreading — kedua Report Definition spreading.
type SumberSpreading interface {
	// `BrowseTreatyArrangement_ParentReinsMasterTrt`.
	Induk(treatyGroupID, tanggalMulai string) []models.SusunanSpreading
	// `BrowseTreatyArrangement_Limit_RD`.
	Anak(parentReinsTypeID, treatyYearID string) []models.SusunanSpreading
}

// MasukanShareNP - satu aksi tab Share beserta seluruh isinya.
type MasukanShareNP struct {
	Aksi  string         `json:"aksi"`
	Share models.ShareNP `json:"share"`
	// Layer Limits Non-Prop — sumber `TreatyInNonAddItem` dan
	// `…SpreadingDetail`.
	Layers []LayerNP `json:"layers"`
	// Indeks baris Share (mulai 0) dan baris deduksi / spreading-nya.
	Indeks int    `json:"indeks"`
	Baris  int    `json:"baris"`
	Sts    string `json:"sts"`
	// `TreatyIn.ID` dan `TreatyIn.Commencement` (YYYYMMDD).
	IDKontrak    string `json:"idKontrak"`
	Commencement string `json:"commencement"`
}

// HasilShareNP - isi tab sesudah aksi, dan pesan Activity.
type HasilShareNP struct {
	Share models.ShareNP `json:"share"`
	// Pesan tingkat tab (mis. `Value Cannot Be Empty` pada % RNM Share).
	Pesan []string `json:"pesan"`
	// Pesan yang Pega tempelkan pada medan SATU baris Share — `Total share
	// must equal RNM share.!!` (`.SpreadingTotalPctXOL`, `SetSpreadingXOL`
	// [9]) dan `Gross Premium (MDP) is still empty` (`.Layer`,
	// `CalculateDeduction` [2]). Di Pega keduanya tampil di panel rincian
	// baris itu, bukan di kepala tab.
	PesanBaris []PesanBarisShare `json:"pesanBaris"`
}

// PesanBarisShare - satu pesan Activity milik baris Share ke-`Indeks`.
type PesanBarisShare struct {
	Indeks int    `json:"indeks"`
	Pesan  string `json:"pesan"`
}

// sumberGudang membaca kedua RD dari gudang sekali per kunci; galat pertama
// disimpan dan dikembalikan sesudah hitungan.
type sumberGudang struct {
	ctx   context.Context
	g     Gudang
	induk map[string][]models.SusunanSpreading
	anak  map[string][]models.SusunanSpreading
	err   error
}

// DescSpreading — `Param.TreatyDescID` yang SELURUH pemanggil RD susunan
// kirim ketika mereka mengirimnya: `"10001"` (`FetchQSfromMaster[2]`,
// `FetchQSfromMasterXOL[4]`, `SetSpreadName[3]`, `Section/DetailShare.xml`).
//
// ⚠️ Satu-satunya yang TIDAK mengirimnya: dropdown `Spreading Type` XOL
// (`Section/Share.xml`) — dan di sana filternya karena itu dilewati.
const DescSpreading = "10001"

func (s *sumberGudang) Induk(grup, mulai string) []models.SusunanSpreading {
	k := grup + "|" + mulai
	if v, ada := s.induk[k]; ada || s.err != nil {
		return v
	}
	// ⛔ ACTIVITY BERBEDA DARI DROPDOWN, dan bedanya baru terbaca setelah
	// langkahnya dibuka satu per satu 8 Oktober 2026:
	//
	//	Section/Share.xml           grup —      desc —       (dropdown)
	//	Activity/FetchQSfromMasterXOL[4]  grup .TreatyGroupList(1).TreatyGroupID
	//	                            desc "10001"
	//
	// Jadi `TreatyDescID` TETAP dikirim di jalur Activity; yang tidak dikirim
	// hanya parameter H. Menghapusnya di sini akan melonggarkan pencarian
	// yang di Pega memang disaring.
	v, err := s.g.BacaIndukSpreading(s.ctx, grup, DescSpreading, mulai, "")
	if err != nil {
		s.err = err
		return nil
	}
	s.induk[k] = v
	return v
}

func (s *sumberGudang) Anak(induk, tahun string) []models.SusunanSpreading {
	k := induk + "|" + tahun
	if v, ada := s.anak[k]; ada || s.err != nil {
		return v
	}
	v, err := s.g.BacaAnakSpreading(s.ctx, induk, tahun)
	if err != nil {
		s.err = err
		return nil
	}
	s.anak[k] = v
	return v
}

// HitungShareNP — bentuk ber-pelaku untuk handler.
func (l *Layanan) HitungShareNP(ctx context.Context, p inti.Pelaku, m MasukanShareNP) (HasilShareNP, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return HasilShareNP{}, err
	}
	src := &sumberGudang{ctx: ctx, g: l.gudang, induk: map[string][]models.SusunanSpreading{}, anak: map[string][]models.SusunanSpreading{}}
	h := HitungShareNP(m, src)
	if src.err != nil {
		return HasilShareNP{}, src.err
	}
	return h, nil
}

// HitungShareNP menjalankan rantai Activity satu aksi. Pure terhadap
// masukan (disalin dulu); bacaan RD lewat `src`.
func HitungShareNP(m MasukanShareNP, src SumberSpreading) HasilShareNP {
	s := salinShare(m.Share)
	pesanBaris := []PesanBarisShare{}
	k := konteksShare{src: src, idKontrak: m.IDKontrak, mulai: m.Commencement, pesanBaris: &pesanBaris}
	pesan := []string{}
	ada := m.Indeks >= 0 && m.Indeks < len(s.Share)
	switch m.Aksi {
	case AksiShareRNM, AksiShareFac:
		pesan = append(pesan, k.xolAddSpreading(&s)...)
	case AksiShareBrokerage:
		// Urutan ekspor sel `% Brokerage`: SetBrokerage LALU XOLAddSpreading.
		pesan = append(pesan, k.setBrokerage(&s)...)
		pesan = append(pesan, k.xolAddSpreading(&s)...)
	case AksiShareSetBrokerage:
		pesan = append(pesan, k.setBrokerage(&s)...)
	case AksiShareCentang:
		if s.RNMShareAcrossTheBoard == benar {
			pesan = append(pesan, k.xolAddSpreading(&s)...)
		}
		if angka(s.BrokeragePercent).Sign() > 0 {
			pesan = append(pesan, k.setBrokerage(&s)...)
		}
	case AksiShareSummary:
		// setValue `FacultativeShare = 0` bila kosong; langkah EDM ber-
		// syarat `EDMMaterialType = '7897987'` — tak pernah benar.
		if s.FacultativeShare == "" {
			s.FacultativeShare = "0"
		}
		pesan = append(pesan, NonAddItemShare(&s, m.Layers)...)
		pesan = append(pesan, k.xolAddSpreading(&s)...)
		pesan = append(pesan, k.setBrokerage(&s)...)
	case AksiShareTotal:
		NPSetTotalShare(&s)
		s.LimitShareSummaryList = SummaryLimitShare(s.Share)
		s.LimitFacShareSummaryList = SummaryLimitShare(s.FacultativeShareList)
	case AksiShareNilai:
		// `TreatyInShareListValue`.
		pesan = append(pesan, NonAddItemShare(&s, m.Layers)...)
		pesan = append(pesan, k.xolAddSpreading(&s)...)
		pesan = append(pesan, k.setBrokerage(&s)...)
		NPSetTotalShare(&s)
		s.LimitShareSummaryList = SummaryLimitShare(s.Share)
		s.LimitFacShareSummaryList = SummaryLimitShare(s.FacultativeShareList)
	case AksiShareSpreadingType:
		if ada {
			b := &s.Share[m.Indeks]
			k.fetchQS(b, b.SpreadingTypeXOL, false)
			totalSebaranOR(&s, true)
		}
	case AksiShareRNMBaris:
		if ada {
			k.catatBaris(m.Indeks, k.xolAddSpreadingDetail(&s, m.Indeks, m.Layers))
		}
	case AksiShareDeduksi:
		if ada {
			k.catatBaris(m.Indeks, hitungDeduksiShare(&s.Share[m.Indeks], m.Sts, m.Baris))
		}
	case AksiShareSpreadingTambah, AksiShareSpreadingHapus:
		if ada {
			AddSpreadingXOL(&s.Share[m.Indeks], m.Aksi == AksiShareSpreadingTambah, m.Baris)
		}
	case AksiShareSpreadingPct:
		if ada {
			k.catatBaris(m.Indeks, k.setSpreadingXOL(&s, &s.Share[m.Indeks]))
		}
	}
	lengkapiShare(&s)
	return HasilShareNP{Share: s, Pesan: pesan, PesanBaris: pesanBaris}
}

type konteksShare struct {
	src        SumberSpreading
	idKontrak  string
	mulai      string
	pesanBaris *[]PesanBarisShare
}

// catatBaris menempelkan pesan Activity pada baris Share ke-`i`.
func (k konteksShare) catatBaris(i int, pesan []string) {
	if k.pesanBaris == nil {
		return
	}
	for _, p := range pesan {
		*k.pesanBaris = append(*k.pesanBaris, PesanBarisShare{Indeks: i, Pesan: p})
	}
}

// --- TreatyInNonAddItem(share) ---------------------------------------------

// NonAddItemShare — `TreatyInNonAddItem` param `share`, langkah 7–14.
func NonAddItemShare(s *models.ShareNP, layers []LayerNP) []string {
	lama := s.Share
	// [7] buang Share dan FacultativeShareList.
	s.Share = []models.BarisShareNP{}
	s.FacultativeShareList = []models.BarisShareNP{}
	// [8] RNMShare == 0 → pesan, keluar.
	if angka(s.RNMShare).IsZero() {
		return []string{PesanShareNilaiKosong}
	}
	// [9]–[10]
	calc := angka(s.RNMShare)
	if fac := angka(s.FacultativeShare); fac.Sign() > 0 {
		calc = kurang(calc, fac)
		s.RnmShareDeducted = teks(calc)
		// [11] Facultative Share — `@divide(FacultativeShare,100,4)`.
		r := bagiBulat(fac, 100, 4)
		for _, l := range layers {
			b := barisDariLayer(l, r)
			b.Limit, b.Limit2 = l.Limit, l.Limit2
			// [11.4] Net := Gross.
			b.NetPremiumList = salinNilai(b.GrossPremiumList)
			s.FacultativeShareList = append(s.FacultativeShareList, b)
		}
	}
	// [13] Share — `@divide(Local.ShareCalc,100,4)`.
	r := bagiBulat(calc, 100, 4)
	for i, l := range layers {
		b := barisDariLayer(l, r)
		// [13.4] Gross Min dari `MDPMinList`. [13.5] Net := Gross DIKOMENTARI.
		for _, v := range l.MDPMinList {
			b.GrossPremiumMinList = append(b.GrossPremiumMinList, models.NilaiMataUang{Currency: v.Currency, Value: teks(kali(angka(v.Value), r))})
		}
		// [14.1] `.RNMShare := TreatyIn.RNMShare`.
		b.RNMShare = s.RNMShare
		// S1 — Spreading Type baris lama berindeks sama dibawa.
		if i < len(lama) {
			b.SpreadingTypeXOL = lama[i].SpreadingTypeXOL
			b.SpreadingTypeIDXOL = lama[i].SpreadingTypeIDXOL
		}
		s.Share = append(s.Share, b)
	}
	return nil
}

// barisDariLayer — langkah 11.1–11.3 / 13.1–13.3: medan layer, 100% Limit
// dan Gross bagian `r`.
func barisDariLayer(l LayerNP, r *apd.Decimal) models.BarisShareNP {
	b := barisShareKosong()
	b.LayerType, b.Layer, b.LayerPartType, b.LayerPart, b.Cover = l.LayerType, l.Layer, l.LayerPartType, l.LayerPart, l.Cover
	for _, g := range l.TreatyGroupList {
		b.TreatyGroupList = append(b.TreatyGroupList, models.GrupShareNP{TreatyGroup: g.TreatyGroup, TreatyGroupID: g.TreatyGroupID})
	}
	b.RnmLimitList = rnmLimitDariLayer(l, r)
	for _, v := range l.MDPList {
		b.GrossPremiumList = append(b.GrossPremiumList, models.NilaiMataUang{Currency: v.Currency, Value: teks(kali(angka(v.Value), r))})
	}
	return b
}

// rnmLimitDariLayer — `Limit × r`, lalu `Limit2 × r` kecuali `.Limit2 < 1`.
func rnmLimitDariLayer(l LayerNP, r *apd.Decimal) []models.NilaiMataUang {
	out := []models.NilaiMataUang{{Currency: l.Currency, Value: teks(kali(angka(l.Limit), r))}}
	if angka(l.Limit2).Cmp(apd.New(1, 0)) >= 0 {
		out = append(out, models.NilaiMataUang{Currency: l.Currency2, Value: teks(kali(angka(l.Limit2), r))})
	}
	return out
}

// --- TreatyInXOLAddSpreading ------------------------------------------------

func (k konteksShare) xolAddSpreading(s *models.ShareNP) []string {
	pesan := []string{}
	// [1]
	for _, n := range []string{"TotalSpreadedRnmProp", "TotalSpreadedRnmRIProp", "TotalSpreadedNetPremiRI", "TotalSpreadedNetPremi"} {
		s.Total[n] = []models.NilaiMataUang{}
	}
	// [7]
	for i := range s.Share {
		s.Share[i].RNMSpreadedListXOL = []models.NilaiMataUang{}
		s.Share[i].RNMSpreadedListRIXOL = []models.NilaiMataUang{}
	}
	// [9]–[10] baris bernama → FetchQS (S2), kosong → SetSpreadingXOL.
	// Pesan baris Share menempel pada barisnya; baris Facultative Share
	// tidak tampil di tab ini, jadi pesannya ke tingkat tab.
	for li, larik := range [][]models.BarisShareNP{s.Share, s.FacultativeShareList} {
		for i := range larik {
			b := &larik[i]
			if b.SpreadingTypeXOL != "" {
				k.fetchQS(b, b.SpreadingTypeXOL, true)
				continue
			}
			p := k.setSpreadingXOL(s, b)
			if li == 0 {
				k.catatBaris(i, p)
			} else {
				pesan = append(pesan, p...)
			}
		}
	}
	// [11] Total OR / R/I Limit dari baris Share.
	s.Total["TotalSpreadedRnmProp"] = []models.NilaiMataUang{}
	s.Total["TotalSpreadedRnmRIProp"] = []models.NilaiMataUang{}
	for _, b := range s.Share {
		s.Total["TotalSpreadedRnmProp"] = jumlahkan(s.Total["TotalSpreadedRnmProp"], b.RNMSpreadedListXOL)
		s.Total["TotalSpreadedRnmRIProp"] = jumlahkan(s.Total["TotalSpreadedRnmRIProp"], b.RNMSpreadedListRIXOL)
	}
	return pesan
}

// --- FetchQSfromMasterXOL -----------------------------------------------------

// fetchQS — `FetchQSfromMasterXOL(ParentReinsTypeID, indexShare, IsUpdate)`.
func (k konteksShare) fetchQS(b *models.BarisShareNP, induk string, perbarui bool) {
	// [3]–[4]
	b.SpreadingTypeXOL = induk
	b.SpreadingListXOL = []models.BarisSpreadingNP{}
	// [5]–[11] induk cocok NAMA → pengenal & tahun → anak susunan.
	// ⛔ FILTER TREATY GROUP SENGAJA DIBUANG — PENYIMPANGAN DARI PEGA, dan ini
	// pernyataannya.
	//
	// `Section/Share.xml` dan `Section/DetailShare.xml` MENGIRIM
	// `TreatyGroupID`, begitu pula `FetchQSfromMaster(XOL)`. Mengikutinya berarti
	// hanya induk yang terdaftar di Treaty Group baris itu yang dapat dipilih —
	// dan untuk kontrak yang dilaporkan 8 Oktober 2026 daftarnya kosong, padahal
	// susunannya ADA di `PROPORTIONALARRG` (dinyatakan pemilik proses).
	//
	// Keputusan pemilik proses, 8 Oktober 2026: *"gimana pun caranya asal itu ada
	// isinya"*.
	//
	// ⚠️ DIBUANG DI KEDUA TEMPAT SEKALIGUS, dan itu syaratnya. Membuangnya hanya
	// di dropdown — yang sempat terjadi — membuat layar menawarkan induk yang
	// pencariannya sendiri tidak dapat menemukan: Spreading Type terpilih,
	// grid spreading tetap kosong. Setengah penyimpangan lebih buruk daripada
	// keduanya, sebab ia terbaca seperti berhasil.
	//
	// ⭐ `TreatyDescID = "10001"` TETAP dikirim: nol pemanggil yang menghilangkannya,
	// dan ia tidak pernah menjadi sebab daftar kosong.
	//
	// Pasangannya di layar: `TabShareNonProp.tsx`.
	if induk != "" {
		idInduk, tahun := "", ""
		for _, p := range k.src.Induk("", k.mulai) {
			if p.ReinsTypeName == induk {
				b.SpreadingTypeIDXOL = p.ReinsTypeID
				idInduk, tahun = p.ReinsTypeID, p.TreatyYearID
			}
		}
		// ⚠️ [11.1] menyaring `Local.ParentID == .ParentReinsTypeID`, dan
		// `Local.ParentID` tak pernah diisi. RD-nya sendiri sudah menyaring
		// induk; data tersimpan memuat anak induk itu (`QS (OR)` 40 /
		// `QS (R/I)` 60 di bawah 10227), jadi seluruh jawaban RD dipakai.
		if idInduk != "" {
			for _, a := range k.src.Anak(idInduk, tahun) {
				b.SpreadingListXOL = append(b.SpreadingListXOL, models.BarisSpreadingNP{
					ReinsTypeName: a.ReinsTypeName, ReinsTypeID: a.ReinsTypeID, ParentReinsTypeID: a.ParentReinsTypeID,
					Pct: a.Pct, Rp: a.Rp, Usd: a.Usd,
				})
			}
		}
	}
	// [12] kode keras kontrak 1000951.
	if k.idKontrak == idKontrakKhusus {
		b.SpreadingListXOL = append(b.SpreadingListXOL, models.BarisSpreadingNP{ReinsTypeName: namaORS, ReinsTypeID: "10007", ParentReinsTypeID: "00", Pct: "15.00"})
	}
	for i := range b.SpreadingListXOL {
		lengkapiSpreading(&b.SpreadingListXOL[i])
	}
	sebarLama(b, perbarui, k.idKontrak)
}

// sebarLama — `FetchQSfromMasterXOL` langkah 13–28 atas `SpreadingListXOL`
// yang sudah ada.
func sebarLama(b *models.BarisShareNP, perbarui bool, idKontrak string) {
	// [13]
	b.SpreadingTotalPctXOL = teks(totalPct(b.SpreadingListXOL))
	// [14], [19]
	q := apd.New(0, 0)
	for i := range b.SpreadingListXOL {
		sp := &b.SpreadingListXOL[i]
		if sp.ReinsTypeName == namaQSOR {
			q = angka(sp.Pct)
		}
		if sp.ReinsTypeName == namaORS {
			q = apd.New(100, 0)
		}
		if b.SpreadingTypeXOL == namaORS {
			sp.ReinsTypeName = namaORS
		}
	}
	or9 := bagiBulat(q, 100, 9)
	ri9 := bagiBulat(kurang(seratus, q), 100, 9)
	khusus := idKontrak == idKontrakKhusus
	// [20] Limit — kontrak 1000951: × @divide(85,100,8) pada baris terakhir.
	b.RNMSpreadedListXOL, b.RNMSpreadedListRIXOL = []models.NilaiMataUang{}, []models.NilaiMataUang{}
	for _, v := range b.RnmLimitList {
		x := angka(v.Value)
		if khusus {
			x = kali(bagiBulat(apd.New(85, 0), 100, 8), x)
		}
		b.RNMSpreadedListXOL = append(b.RNMSpreadedListXOL, models.NilaiMataUang{Currency: v.Currency, Value: teks(kali(or9, x))})
		b.RNMSpreadedListRIXOL = append(b.RNMSpreadedListRIXOL, models.NilaiMataUang{Currency: v.Currency, Value: teks(kali(ri9, x))})
	}
	// [22] Gross, [24] Gross Min — `@divide(…,100,9)`.
	b.RNMSpreadedListGrossXOL, b.RNMSpreadedListGrossRIXOL = bagiDua(b.GrossPremiumList, or9, ri9)
	b.RNMSpreadedListGrossMinXOL, b.RNMSpreadedListGrossRIMinXOL = bagiDua(b.GrossPremiumMinList, or9, ri9)
	// [26] Deduction — pembagian biasa `(QSPCT/100)`.
	orP, riP := bagiPolos(q, seratus), bagiPolos(kurang(seratus, q), seratus)
	b.RNMSpreadedListDeductXOL, b.RNMSpreadedListDeductRIXOL = bagiDua(b.DeductionTotalList, orP, riP)
	// [28] Net — hanya bila bukan `IsUpdate == 1` (S3: disusun ulang).
	if !perbarui {
		b.RNMSpreadedListNetXOL, b.RNMSpreadedListNetRIXOL = bagiDua(b.NetPremiumList, orP, riP)
	}
}

// --- SetSpreadingXOL ----------------------------------------------------------

// setSpreadingXOL — spreading MANUAL (Spreading Type kosong).
func (k konteksShare) setSpreadingXOL(s *models.ShareNP, b *models.BarisShareNP) []string {
	// [3] keempat total spreaded dibuang.
	for _, n := range []string{"TotalSpreadedRnmProp", "TotalSpreadedRnmRIProp", "TotalSpreadedNetPremi", "TotalSpreadedNetPremiRI"} {
		s.Total[n] = []models.NilaiMataUang{}
	}
	// [4]
	kosongkanSebaran(b)
	for i := range b.SpreadingListXOL {
		sp := &b.SpreadingListXOL[i]
		// [5.3.4] nama induk dari RD bila `ReinsTypeID` terisi — parameter
		// `TreatyGroupID := TempSprd.TreatyGroupID`, yang tak pernah diisi
		// untuk baris Non-Prop: filter B dilewati, induk SEMUA grup.
		// Rincian anaknya ([5.3.7]) bersyarat `Primary.SpreadingTypeID ==
		// ""` → LEWATI, dan `SpreadingTypeID` tidak pernah diisi di
		// Non-Prop: rincian itu tak tercapai, jadi bagian OR / R/I-nya kosong.
		if sp.ReinsTypeID != "" {
			for _, p := range k.src.Induk("", k.mulai) {
				if p.ReinsTypeID == sp.ReinsTypeID {
					sp.ReinsTypeName = p.ReinsTypeName
				}
			}
		}
	}
	sebarManual(b)
	// [6] `CountTotalPctSpreadXOL`.
	b.SpreadingTotalPctXOL = teks(totalPct(b.SpreadingListXOL))
	// [7] menjumlah larik spreaded TINGKAT SPREADING — selalu kosong (lihat
	// [5.3.7]), jadi keempat total tetap kosong.
	// [9]
	if angka(b.SpreadingTotalPctXOL).Cmp(angka(b.RNMShare)) != 0 {
		return []string{PesanShareSpreading}
	}
	return nil
}

// sebarManual — `SetSpreadingXOL` langkah 5.1–5.8: tiap baris spreading
// mendapat bagian `@divide(.Pct, RNMShare, 20)` dari larik baris Share.
func sebarManual(b *models.BarisShareNP) {
	for i := range b.SpreadingListXOL {
		sp := &b.SpreadingListXOL[i]
		f := bagiBulatDes(angka(sp.Pct), angka(b.RNMShare), 20)
		// [5.4] `.RnmLimitList(<APPEND>)` dari `.Currency`/`.Value` baris
		// spreading — keduanya hanya diisi rincian anak yang tak tercapai,
		// jadi barisnya bermata uang kosong dan tak pernah dijumlah.
		sp.RnmLimitList = []models.NilaiMataUang{}
		sp.GrossPremiumList = kaliLarik(b.GrossPremiumList, f)
		sp.GrossPremiumMinList = kaliLarik(b.GrossPremiumMinList, f)
		sp.DeductionTotalList = kaliLarik(b.DeductionTotalList, f)
		sp.NetPremiumList = kaliLarik(b.NetPremiumList, f)
	}
}

// AddSpreadingXOL — `add=true`: baris baru `Pct 0`; `add=false`: hapus baris
// `idx`. Lalu `CountTotalPctSpreadXOL`.
func AddSpreadingXOL(b *models.BarisShareNP, tambahBaris bool, idx int) {
	if tambahBaris {
		sp := models.BarisSpreadingNP{Pct: "0"}
		lengkapiSpreading(&sp)
		b.SpreadingListXOL = append(b.SpreadingListXOL, sp)
	} else if idx >= 0 && idx < len(b.SpreadingListXOL) {
		b.SpreadingListXOL = append(b.SpreadingListXOL[:idx:idx], b.SpreadingListXOL[idx+1:]...)
	}
	b.SpreadingTotalPctXOL = teks(totalPct(b.SpreadingListXOL))
}

// --- TreatyInSetBrokerage -----------------------------------------------------

func (k konteksShare) setBrokerage(s *models.ShareNP) []string {
	pesan := []string{}
	// [1]–[3] Brokerage fee → CalculateDeduction(pct, 1).
	for i := range s.Share {
		s.Share[i].DeductionList = []models.BarisDeduksiShare{{Comment: komentarBroker, DeductionPct: s.BrokeragePercent, DeductionPctCalculate: benar}}
		k.catatBaris(i, hitungDeduksiShare(&s.Share[i], "pct", 0))
	}
	// [4]–[5] Overiding Commision (S4: disusun ulang).
	for i := range s.FacultativeShareList {
		s.FacultativeShareList[i].DeductionList = []models.BarisDeduksiShare{{Comment: komentarFac, DeductionPct: s.FacultativeShareBrokerage, DeductionPctCalculate: benar}}
		pesan = append(pesan, hitungDeduksiShare(&s.FacultativeShareList[i], "pct", 0)...)
	}
	// [6] buang SATU baris deduksi bernilai < 1 — yang TERAKHIR ditemukan.
	// Total dan Net tidak dihitung ulang (Activity tidak memanggilnya).
	for i := range s.Share {
		hapus := -1
		for j, d := range s.Share[i].DeductionList {
			if angka(d.Deduction).Cmp(apd.New(1, 0)) < 0 {
				hapus = j
			}
		}
		if hapus >= 0 {
			dl := s.Share[i].DeductionList
			s.Share[i].DeductionList = append(dl[:hapus:hapus], dl[hapus+1:]...)
		}
	}
	// [7]–[8] Net & Deduction OR / R/I. `local.QSPCT` hanya dari `QS (OR)`
	// dan TERBAWA antar baris (Share lalu Facultative Share).
	q := apd.New(0, 0)
	for _, larik := range [][]models.BarisShareNP{s.Share, s.FacultativeShareList} {
		for i := range larik {
			b := &larik[i]
			for _, sp := range b.SpreadingListXOL {
				if sp.ReinsTypeName == namaQSOR {
					q = angka(sp.Pct)
				}
			}
			sebarBrokerage(b, q, k.idKontrak)
		}
	}
	// [9.3]–[9.4] total Net OR / R/I dari baris Share — DITAMBAHKAN ke yang
	// sudah ada (Activity tidak membuangnya).
	for _, b := range s.Share {
		s.Total["TotalSpreadedNetPremi"] = jumlahkan(s.Total["TotalSpreadedNetPremi"], b.RNMSpreadedListNetXOL)
		s.Total["TotalSpreadedNetPremiRI"] = jumlahkan(s.Total["TotalSpreadedNetPremiRI"], b.RNMSpreadedListNetRIXOL)
	}
	return pesan
}

// sebarBrokerage — `TreatyInSetBrokerage` 7.3–7.4 untuk satu baris (S3).
func sebarBrokerage(b *models.BarisShareNP, q *apd.Decimal, idKontrak string) {
	orP, riP := bagiPolos(q, seratus), bagiPolos(kurang(seratus, q), seratus)
	net := b.NetPremiumList
	if idKontrak == idKontrakKhusus {
		net = kaliLarik(net, bagiBulat(apd.New(85, 0), 100, 8))
	}
	b.RNMSpreadedListNetXOL, b.RNMSpreadedListNetRIXOL = bagiDua(net, orP, riP)
	ded := b.DeductionTotalList
	if idKontrak == idKontrakKhusus {
		ded = kaliLarik(ded, bagiBulat(apd.New(85, 0), 100, 8))
	}
	b.RNMSpreadedListDeductXOL, b.RNMSpreadedListDeductRIXOL = bagiDua(ded, orP, riP)
}

// hitungDeduksiShare — `CalculateDeduction(sts, index)` atas satu baris.
func hitungDeduksiShare(b *models.BarisShareNP, sts string, idx int) []string {
	h := HitungDeduksi(MasukanDeduksi{Sts: sts, Indeks: idx, DeductionList: b.DeductionList, GrossPremiumList: b.GrossPremiumList, BenderaPct: true})
	b.DeductionList = h.DeductionList
	b.DeductionTotalList = h.DeductionTotalList
	b.NetPremiumList = h.NetPremiumList
	return h.Pesan
}

// --- TreatyInXOLAddSpreadingDetail ------------------------------------------

func (k konteksShare) xolAddSpreadingDetail(s *models.ShareNP, idx int, layers []LayerNP) []string {
	b := &s.Share[idx]
	// [1]
	for _, n := range []string{"TotalSpreadedRnmProp", "TotalSpreadedRnmRIProp", "TotalSpreadedNetPremi", "TotalSpreadedNetPremiRI"} {
		s.Total[n] = []models.NilaiMataUang{}
	}
	b.RnmLimitList, b.GrossPremiumList = []models.NilaiMataUang{}, []models.NilaiMataUang{}
	// [2] dari `Limits(idx)` dengan `@divide(RNMShare baris,100,4)`.
	if idx < len(layers) {
		l := layers[idx]
		r := bagiBulat(angka(b.RNMShare), 100, 4)
		b.RnmLimitList = rnmLimitDariLayer(l, r)
		for _, v := range l.MDPList {
			b.GrossPremiumList = append(b.GrossPremiumList, models.NilaiMataUang{Currency: v.Currency, Value: teks(kali(angka(v.Value), r))})
		}
	}
	// [2.4] Net := Gross.
	b.NetPremiumList = salinNilai(b.GrossPremiumList)
	// [4.1] FetchQS TANPA induk — Spreading Type baris ini terhapus.
	k.fetchQS(b, "", false)
	// [4.2]
	pesan := hitungDeduksiShare(b, "pct", 0)
	// [5] QSPCT dari `QS (OR)` (kini kosong → 0); Net & Deduction OR / R/I.
	q := apd.New(0, 0)
	for _, sp := range b.SpreadingListXOL {
		if sp.ReinsTypeName == namaQSOR {
			q = angka(sp.Pct)
		}
	}
	orP, riP := bagiPolos(q, seratus), bagiPolos(kurang(seratus, q), seratus)
	b.RNMSpreadedListNetXOL, b.RNMSpreadedListNetRIXOL = bagiDua(b.NetPremiumList, orP, riP)
	b.RNMSpreadedListDeductXOL, b.RNMSpreadedListDeductRIXOL = bagiDua(b.DeductionTotalList, orP, riP)
	// [6]
	totalSebaranOR(s, true)
	return pesan
}

// totalSebaranOR — total OR / R/I Limit (dan Net bila `net`) dari seluruh
// baris Share, per mata uang.
func totalSebaranOR(s *models.ShareNP, net bool) {
	s.Total["TotalSpreadedRnmProp"] = []models.NilaiMataUang{}
	s.Total["TotalSpreadedRnmRIProp"] = []models.NilaiMataUang{}
	if net {
		s.Total["TotalSpreadedNetPremi"] = []models.NilaiMataUang{}
		s.Total["TotalSpreadedNetPremiRI"] = []models.NilaiMataUang{}
	}
	for _, b := range s.Share {
		s.Total["TotalSpreadedRnmProp"] = jumlahkan(s.Total["TotalSpreadedRnmProp"], b.RNMSpreadedListXOL)
		s.Total["TotalSpreadedRnmRIProp"] = jumlahkan(s.Total["TotalSpreadedRnmRIProp"], b.RNMSpreadedListRIXOL)
		if net {
			s.Total["TotalSpreadedNetPremi"] = jumlahkan(s.Total["TotalSpreadedNetPremi"], b.RNMSpreadedListNetXOL)
			s.Total["TotalSpreadedNetPremiRI"] = jumlahkan(s.Total["TotalSpreadedNetPremiRI"], b.RNMSpreadedListNetRIXOL)
		}
	}
}

// --- TreatyInNPSetTotal(share) ------------------------------------------------

// NPSetTotalShare — `TreatyInNPSetTotal(share)`: jumlah per mata uang dari
// larik SETIAP baris Share. Mata uang kosong tidak ditambahkan.
//
// ---------------------------------------------------------------------
// ⛔ RALAT 8 Oktober 2026 — HASIL PEGA MENGALAHKAN BACAAN TEKS EKSPOR
// ---------------------------------------------------------------------
// Bacaan lama: hanya [24], hanya baris ber-Spreading Type KOSONG, dari
// larik tingkat SPREADING. Untuk kontrak 1001855 (tiga baris ber-Spreading
// Type, sebaran QS 40/60) itu membuat Total RNM Limit, Gross, Deduction,
// dan Net "No items". Tangkapan layar Pega pemakai untuk kontrak yang SAMA:
//
//	Total RNM Limit   35.250.000.000 = 4 M + 6,25 M + 25 M  (RnmLimitList)
//	Total Gross (MDP)  1.271.036.250 = Σ GrossPremiumList ketiga baris
//	Total Deduction                0 · Total Net 1.271.036.250
//	Total Gross Min Premium   No items
//
// — jumlah larik baris Share SENDIRI, semua baris. Itu juga bentuk
// kembarannya yang sudah diport, `TreatyInNPSetTotalActualShare`
// (`hitung_aktual.go` `totalShareAktual`).
//
// ⚠️ Yang TETAP: Total Gross Min Premium kosong (tangkapan yang sama), dan
// `FacultativeShareList` tidak ikut.
//
// ---------------------------------------------------------------------
// ⛔ LANGKAH [23] ADA DI EKSPOR DAN TIDAK DIBANGUN — IA DIKOMENTARI
// ---------------------------------------------------------------------
// Di antara [22] (mengosongkan) dan [24] (yang dibangun di bawah) ada satu
// blok lagi yang menjumlah larik baris Share SENDIRI, berikut
// `GrossPremiumMinList`. Ia TIDAK berjalan:
//
//	[22] Property-Remove  pyStepsBlockName `Share`    praAktif true
//	[23] (blok)           pyStepsBlockName `//`       praAktif false
//	[24] (blok)           pyStepsBlockName ``         praAktif true
//	[25] Property-Remove  pyStepsBlockName `FShare`   praAktif true
//
// ⭐ `pyStepsBlockName` = `//` adalah tanda Pega untuk langkah yang
// DIKOMENTARI. Ia mematikan langkah terlepas dari `pyStepsPreCondition`.
//
// ⚠️ `praAktif=false` BUKAN tandanya, dan mengira begitu adalah kekeliruan
// yang sudah pernah terjadi di berkas ini: [26] juga `praAktif=false` dengan
// `blockName` kosong, dan ia BERJALAN — ia yang mengisi `TotalFacShare*`.
// `praAktif=false` hanya berarti "tanpa pra-syarat".
//
// ⛔ Akibat yang harus diketahui: `TotalShareGrossMinNP` karena itu SELALU
// kosong, di Pega maupun di sini. Grid `Total Gross Min Premium` memang
// selalu kosong — itu bukan cacat port.
func NPSetTotalShare(s *models.ShareNP) {
	for _, n := range []string{"TotalShareRnmNP", "TotalShareGrossNP", "TotalShareNetNP", "TotalShareDeductionNP", "TotalShareGrossMinNP"} {
		s.Total[n] = []models.NilaiMataUang{}
	}
	for _, b := range s.Share {
		s.Total["TotalShareGrossNP"] = jumlahkanBerMataUang(s.Total["TotalShareGrossNP"], b.GrossPremiumList)
		s.Total["TotalShareNetNP"] = jumlahkanBerMataUang(s.Total["TotalShareNetNP"], b.NetPremiumList)
		s.Total["TotalShareDeductionNP"] = jumlahkanBerMataUang(s.Total["TotalShareDeductionNP"], b.DeductionTotalList)
		s.Total["TotalShareRnmNP"] = jumlahkanBerMataUang(s.Total["TotalShareRnmNP"], b.RnmLimitList)
	}
}

// --- TreatyInSummaryLimitShare / FacShare -------------------------------------

// SummaryLimitShare — satu baris per layer (kunci empat medan layer):
// kolom IDR dan USD dari 100% Limit, Gross (MDP), Deduction, Net. Mata uang
// lain tidak masuk. Layer yang 100% Limit-nya kosong tidak masuk sama
// sekali (langkah 4.3.3 satu-satunya yang menambah baris).
func SummaryLimitShare(rows []models.BarisShareNP) []models.RingkasanShareNP {
	out := []models.RingkasanShareNP{}
	sama := func(r models.RingkasanShareNP, b models.BarisShareNP) bool {
		return r.LayerType == b.LayerType && r.Layer == b.Layer && r.LayerPartType == b.LayerPartType && r.LayerPart == b.LayerPart
	}
	tambahKe := func(b models.BarisShareNP, f func(r *models.RingkasanShareNP)) {
		for i := range out {
			if sama(out[i], b) {
				f(&out[i])
			}
		}
	}
	tambahNilai := func(isi string, v models.NilaiMataUang, mau string) string {
		if v.Currency != mau {
			return isi
		}
		return teks(tambah(angka(isi), angka(v.Value)))
	}
	for _, b := range rows {
		// [4.2]
		flag := false
		for _, r := range out {
			if sama(r, b) {
				flag = true
			}
		}
		// [4.3]
		for _, v := range b.RnmLimitList {
			if !flag {
				r := models.RingkasanShareNP{LayerType: b.LayerType, Layer: b.Layer, LayerPartType: b.LayerPartType, LayerPart: b.LayerPart, Limit: "0", Limit2: "0"}
				if v.Currency == "IDR" {
					r.Limit = teks(angka(v.Value))
				}
				if v.Currency == "USD" {
					r.Limit2 = teks(angka(v.Value))
				}
				out = append(out, r)
			} else {
				tambahKe(b, func(r *models.RingkasanShareNP) {
					r.Limit = tambahNilai(r.Limit, v, "IDR")
					r.Limit2 = tambahNilai(r.Limit2, v, "USD")
				})
			}
			flag = true
		}
		if !flag {
			continue
		}
		// [4.4]–[4.6]
		for _, d := range b.DeductionList {
			v := models.NilaiMataUang{Currency: d.Currency, Value: d.Deduction}
			tambahKe(b, func(r *models.RingkasanShareNP) {
				r.Deductible = tambahNilai(r.Deductible, v, "IDR")
				r.Deductible2 = tambahNilai(r.Deductible2, v, "USD")
			})
		}
		for _, v := range b.NetPremiumList {
			tambahKe(b, func(r *models.RingkasanShareNP) {
				r.NetPremi = tambahNilai(r.NetPremi, v, "IDR")
				r.NetPremi2 = tambahNilai(r.NetPremi2, v, "USD")
			})
		}
		for _, v := range b.GrossPremiumList {
			tambahKe(b, func(r *models.RingkasanShareNP) {
				r.MDP = tambahNilai(r.MDP, v, "IDR")
				r.MDP2 = tambahNilai(r.MDP2, v, "USD")
			})
		}
	}
	// [6] `.LayerType + .Layer + " of " + .LayerPartType + .LayerPart`.
	for i := range out {
		r := &out[i]
		r.Note = r.LayerType + r.Layer + " of " + r.LayerPartType + r.LayerPart
	}
	return out
}

// --- pembantu ---------------------------------------------------------------

func bagiPolos(a, b *apd.Decimal) *apd.Decimal {
	r := new(apd.Decimal)
	if b.IsZero() {
		return r
	}
	_, _ = konteksLimit.Quo(r, a, b)
	return r
}

func totalPct(daftar []models.BarisSpreadingNP) *apd.Decimal {
	t := apd.New(0, 0)
	for _, sp := range daftar {
		t = tambah(t, angka(sp.Pct))
	}
	return t
}

func kaliLarik(daftar []models.NilaiMataUang, f *apd.Decimal) []models.NilaiMataUang {
	out := []models.NilaiMataUang{}
	for _, v := range daftar {
		out = append(out, models.NilaiMataUang{Currency: v.Currency, CurrencyID: v.CurrencyID, Value: teks(kali(f, angka(v.Value)))})
	}
	return out
}

func bagiDua(daftar []models.NilaiMataUang, a, b *apd.Decimal) ([]models.NilaiMataUang, []models.NilaiMataUang) {
	x, y := []models.NilaiMataUang{}, []models.NilaiMataUang{}
	for _, v := range daftar {
		x = append(x, models.NilaiMataUang{Currency: v.Currency, Value: teks(kali(a, angka(v.Value)))})
		y = append(y, models.NilaiMataUang{Currency: v.Currency, Value: teks(kali(b, angka(v.Value)))})
	}
	return x, y
}

// jumlahkan — pola `Local.test` seluruh Activity total: tambah ke baris
// bermata uang sama, atau tambah baris.
func jumlahkan(tujuan, sumber []models.NilaiMataUang) []models.NilaiMataUang {
	for _, v := range sumber {
		tujuan = tambahPerMataUang(tujuan, v.Currency, "", angka(v.Value))
	}
	return tujuan
}

// jumlahkanBerMataUang — `TreatyInNPSetTotal`: baris BARU hanya bila mata
// uangnya terisi (`.Currency != ""`).
func jumlahkanBerMataUang(tujuan, sumber []models.NilaiMataUang) []models.NilaiMataUang {
	for _, v := range sumber {
		ada := false
		for i := range tujuan {
			if tujuan[i].Currency == v.Currency {
				tujuan[i].Value = teks(tambah(angka(tujuan[i].Value), angka(v.Value)))
				ada = true
			}
		}
		if !ada && v.Currency != "" {
			tujuan = append(tujuan, models.NilaiMataUang{Currency: v.Currency, Value: teks(angka(v.Value))})
		}
	}
	return tujuan
}

func salinNilai(v []models.NilaiMataUang) []models.NilaiMataUang {
	return append([]models.NilaiMataUang{}, v...)
}

func kosongkanSebaran(b *models.BarisShareNP) {
	b.RNMSpreadedListXOL, b.RNMSpreadedListRIXOL = []models.NilaiMataUang{}, []models.NilaiMataUang{}
	b.RNMSpreadedListGrossXOL, b.RNMSpreadedListGrossRIXOL = []models.NilaiMataUang{}, []models.NilaiMataUang{}
	b.RNMSpreadedListDeductXOL, b.RNMSpreadedListDeductRIXOL = []models.NilaiMataUang{}, []models.NilaiMataUang{}
	b.RNMSpreadedListNetXOL, b.RNMSpreadedListNetRIXOL = []models.NilaiMataUang{}, []models.NilaiMataUang{}
}

// DaftarIndukSpreading — isi kedua dropdown spreading panel Share:
// `BrowseTreatyArrangement_ParentReinsMasterTrt` dengan `TreatyIn.Commencement`
// dan Treaty Group — `.TreatyGroupList(1).TreatyGroupID` untuk Spreading Type,
// KOSONG untuk Reins Type spreading manual (`TempSprd.TreatyGroupID`, tak
// pernah diisi → filter B dilewati, induk semua grup).
func (l *Layanan) DaftarIndukSpreading(ctx context.Context, p inti.Pelaku, treatyGroupID, treatyDescID, mulai string) ([]models.SusunanSpreading, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return nil, err
	}
	// Kedua dropdown mengirim `ReinsTypeID = "10246"` (filter H). `TreatyDescID`
	// datang dari layar: Prop `"10001"`, XOL kosong — lihat
	// `SaringanIndukSpreading`.
	return l.gudang.BacaIndukSpreading(ctx, treatyGroupID, treatyDescID, mulai, repository.IndukDikecualikanDropdown)
}

// DaftarReasuradurShare — isi autocomplete `Reinsurer Name` /
// `Facultative Reinsurers`.
func (l *Layanan) DaftarReasuradurShare(ctx context.Context, p inti.Pelaku) ([]models.PilihanWarisan, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return nil, err
	}
	return l.gudang.BacaDaftarReasuradurShare(ctx)
}
