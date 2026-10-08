package services

// Tab Share Non-Prop saat kontrak DIBUKA — dari tabel pendaratan, lalu nilai
// turunan dihitung dengan rumus Activity yang sama (`hitung_share_np.go`).
//
// ⭐ SKALAR AKAR — urutan sumber, keputusan pemakai 7 Oktober 2026 (Save
// menyimpan ke tabel sendiri: skema v2, akar `TreatyIn` → `T_TREATY_REVISION`):
//
//	1. kolom `T_TREATY_REVISION` (migrasi 445) — `TerapkanAkarShare`;
//	2. salinan yang Pega TULIS sendiri — `ShareDariPendaratan`:
//	   RNMShare          baris Share pertama yang berisi (`TreatyInNonAddItem`
//	                     [14.1] `.RNMShare := TreatyIn.RNMShare`)
//	   BrokeragePercent  deduksi `Brokerage fee` baris Share pertama
//	                     (`TreatyInSetBrokerage` [2])
//	3. `TREATYINDETAIL` (`SaveTreatyInDetail_Act` [5.1]) — `TerapkanAkarShare`.
//
//	RNMShareAcrossTheBoard tanpa kolom → `pyDefaultValue = true` (@1781134).
//	RnmShareDeducted       tanpa kolom → `RNMShare − FacultativeShare` bila > 0.
//
// Nol tebakan dari rasio Gross/MDP: kontrak tanpa satu pun sumber tampil
// KOSONG sampai diisi lalu Save.
//
// ⚠️ Larik yang TIDAK didaratkan dan karena itu DIHITUNG saat dibuka:
// `RnmLimitList` (layer Limits berindeks sama × bagian RNM), `TreatyGroupList`
// baris Share (salinan layer), `DeductionTotalList`, seluruh bagian OR / R/I,
// Summary, dan kesembilan Total. Hasilnya = yang Pega simpan sesudah tombol
// terakhir ditekan, SELAMA nilai tersimpannya tidak basi. `GrossPremiumMinList`
// (dari `MDPMinList`, juga tak didaratkan) kosong sampai Update Summary.

import (
	"encoding/json"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/modul/treatyin/backend/models"
)

// ShareDariPendaratan menafsirkan pohon pendaratan + medan revisi menjadi
// isi tab. `medan` = `RevisiPendaratan.Medan` (ejaan dokumen).
func ShareDariPendaratan(sp models.SharePendaratan, medan map[string]string) models.ShareNP {
	s := models.ShareNP{
		FacultativeShare:          medan["FacultativeShare"],
		FacultativeShareBrokerage: medan["FacultativeShareBrokerage"],
		IsProRate:                 medan["IsProRate"],
		RNMShareAcrossTheBoard:    benar,
	}
	dariPeta(sp.Share, &s.Share)
	dariPeta(sp.FacultativeShareList, &s.FacultativeShareList)
	dariPeta(sp.ShareReins, &s.ShareReins)
	dariPeta(sp.ShareFacultativeReinsurers, &s.ShareFacultativeReinsurers)
	for _, b := range s.Share {
		if strings.TrimSpace(b.RNMShare) != "" {
			s.RNMShare = b.RNMShare
			break
		}
	}
	if len(s.Share) > 0 {
		for _, d := range s.Share[0].DeductionList {
			if d.Comment == komentarBroker {
				s.BrokeragePercent = d.DeductionPct
				break
			}
		}
	}
	if fac := angka(s.FacultativeShare); fac.Sign() > 0 && s.RNMShare != "" {
		s.RnmShareDeducted = teks(kurang(angka(s.RNMShare), fac))
	}
	lengkapiShare(&s)
	return s
}

// SiapkanShareNP mengisi larik turunan saat kontrak dibuka.
func SiapkanShareNP(s *models.ShareNP, layers []LayerNP, idKontrak string) {
	lengkapiShare(s)
	calc := angka(s.RNMShare)
	fac := angka(s.FacultativeShare)
	if fac.Sign() > 0 {
		calc = kurang(calc, fac)
	}
	rShare := bagiBulat(calc, 100, 4)
	rFac := bagiBulat(fac, 100, 4)
	isi := func(larik []models.BarisShareNP, r *apd.Decimal, adaBagian bool) {
		for i := range larik {
			b := &larik[i]
			if i < len(layers) {
				l := layers[i]
				if len(b.TreatyGroupList) == 0 {
					for _, g := range l.TreatyGroupList {
						b.TreatyGroupList = append(b.TreatyGroupList, models.GrupShareNP{TreatyGroup: g.TreatyGroup, TreatyGroupID: g.TreatyGroupID})
					}
				}
				if adaBagian {
					b.RnmLimitList = rnmLimitDariLayer(l, r)
				}
			}
			b.DeductionTotalList = []models.NilaiMataUang{}
			for _, d := range b.DeductionList {
				b.DeductionTotalList = tambahPerMataUang(b.DeductionTotalList, d.Currency, "", angka(d.Deduction))
			}
			if b.SpreadingTypeXOL != "" {
				sebarLama(b, true, idKontrak)
			} else {
				sebarManual(b)
			}
		}
	}
	isi(s.Share, rShare, strings.TrimSpace(s.RNMShare) != "" && !calc.IsZero())
	isi(s.FacultativeShareList, rFac, fac.Sign() > 0)
	// `TreatyInSetBrokerage` [7]–[8].
	q := apd.New(0, 0)
	for _, larik := range [][]models.BarisShareNP{s.Share, s.FacultativeShareList} {
		for i := range larik {
			for _, sp := range larik[i].SpreadingListXOL {
				if sp.ReinsTypeName == namaQSOR {
					q = angka(sp.Pct)
				}
			}
			sebarBrokerage(&larik[i], q, idKontrak)
		}
	}
	totalSebaranOR(s, true)
	NPSetTotalShare(s)
	s.LimitShareSummaryList = SummaryLimitShare(s.Share)
	s.LimitFacShareSummaryList = SummaryLimitShare(s.FacultativeShareList)
	lengkapiShare(s)
}

// LayerDariPohon — simpul pohon Limits Non-Prop → `LayerNP` (ejaan sama).
func LayerDariPohon(pohon []map[string]any) []LayerNP {
	var out []LayerNP
	dariPeta(pohon, &out)
	if out == nil {
		out = []LayerNP{}
	}
	return out
}

// dariPeta menyalin simpul berkunci ejaan dokumen ke struktur berejaan sama.
// Medan yang tak dikenal diabaikan; galat berarti bentuk tak cocok dan
// hasilnya dibiarkan kosong.
func dariPeta(src any, dst any) {
	b, err := json.Marshal(src)
	if err != nil {
		return
	}
	_ = json.Unmarshal(b, dst)
}

// salinShare — salinan dalam; masukan aksi tidak diubah.
func salinShare(s models.ShareNP) models.ShareNP {
	var out models.ShareNP
	dariPeta(s, &out)
	lengkapiShare(&out)
	return out
}

func barisShareKosong() models.BarisShareNP {
	b := models.BarisShareNP{}
	lengkapiBaris(&b)
	return b
}

func lengkapiSpreading(sp *models.BarisSpreadingNP) {
	for _, l := range []*[]models.NilaiMataUang{&sp.RnmLimitList, &sp.GrossPremiumList, &sp.GrossPremiumMinList, &sp.DeductionTotalList, &sp.NetPremiumList} {
		if *l == nil {
			*l = []models.NilaiMataUang{}
		}
	}
}

func lengkapiBaris(b *models.BarisShareNP) {
	if b.TreatyGroupList == nil {
		b.TreatyGroupList = []models.GrupShareNP{}
	}
	if b.SpreadingListXOL == nil {
		b.SpreadingListXOL = []models.BarisSpreadingNP{}
	}
	for i := range b.SpreadingListXOL {
		lengkapiSpreading(&b.SpreadingListXOL[i])
	}
	if b.DeductionList == nil {
		b.DeductionList = []models.BarisDeduksiShare{}
	}
	for _, l := range []*[]models.NilaiMataUang{
		&b.RnmLimitList, &b.GrossPremiumList, &b.GrossPremiumMinList, &b.DeductionTotalList, &b.NetPremiumList,
		&b.RNMSpreadedListXOL, &b.RNMSpreadedListRIXOL, &b.RNMSpreadedListGrossXOL, &b.RNMSpreadedListGrossRIXOL,
		&b.RNMSpreadedListGrossMinXOL, &b.RNMSpreadedListGrossRIMinXOL, &b.RNMSpreadedListDeductXOL,
		&b.RNMSpreadedListDeductRIXOL, &b.RNMSpreadedListNetXOL, &b.RNMSpreadedListNetRIXOL,
	} {
		if *l == nil {
			*l = []models.NilaiMataUang{}
		}
	}
}

// lengkapiShare — larik nihil menjadi larik KOSONG (layar membaca
// `.length`), dan kesembilan kunci Total selalu hadir.
func lengkapiShare(s *models.ShareNP) {
	if s.ShareReins == nil {
		s.ShareReins = []models.BarisReinsShare{}
	}
	if s.ShareFacultativeReinsurers == nil {
		s.ShareFacultativeReinsurers = []models.BarisReinsShare{}
	}
	if s.Share == nil {
		s.Share = []models.BarisShareNP{}
	}
	if s.FacultativeShareList == nil {
		s.FacultativeShareList = []models.BarisShareNP{}
	}
	for i := range s.Share {
		lengkapiBaris(&s.Share[i])
	}
	for i := range s.FacultativeShareList {
		lengkapiBaris(&s.FacultativeShareList[i])
	}
	if s.LimitShareSummaryList == nil {
		s.LimitShareSummaryList = []models.RingkasanShareNP{}
	}
	if s.LimitFacShareSummaryList == nil {
		s.LimitFacShareSummaryList = []models.RingkasanShareNP{}
	}
	if s.Total == nil {
		s.Total = map[string][]models.NilaiMataUang{}
	}
	for _, k := range KunciTotalShareNP {
		if s.Total[k] == nil {
			s.Total[k] = []models.NilaiMataUang{}
		}
	}
}

// TerapkanAkarShare menetapkan skalar akar menurut urutan sumber yang
// diputuskan pemakai 7 Oktober 2026:
//
//  1. kolom `T_TREATY_REVISION` (migrasi 445) — tempat Save menyimpan;
//  2. salinan yang Pega tulis ke baris Share / deduksi `Brokerage fee`
//     (sudah diisi `ShareDariPendaratan`);
//  3. `TREATYINDETAIL` (`SaveTreatyInDetail_Act`).
//
// `RNMShareAcrossTheBoard` tanpa kolom → bawaan Pega `true`.
// `RnmShareDeducted` tanpa kolom → `RNMShare − FacultativeShare` bila > 0.
// Tidak ada sumber yang menebak dari rasio Gross/MDP.
func TerapkanAkarShare(s *models.ShareNP, akar map[string]string, d models.ShareDetailWarisan) {
	pilih := func(kolom, kini, cadangan string) string {
		for _, v := range []string{kolom, kini, cadangan} {
			if strings.TrimSpace(v) != "" {
				return v
			}
		}
		return ""
	}
	s.RNMShare = pilih(akar["RNMShare"], s.RNMShare, d.RNMShare)
	s.BrokeragePercent = pilih(akar["BrokeragePercent"], s.BrokeragePercent, d.BrokeragePercent)
	if v := strings.TrimSpace(akar["RNMShareAcrossTheBoard"]); v != "" {
		s.RNMShareAcrossTheBoard = v
	}
	if v := strings.TrimSpace(akar["RnmShareDeducted"]); v != "" {
		s.RnmShareDeducted = v
	} else if fac := angka(s.FacultativeShare); fac.Sign() > 0 && s.RNMShare != "" {
		s.RnmShareDeducted = teks(kurang(angka(s.RNMShare), fac))
	}
}
