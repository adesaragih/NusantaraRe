package services

// Sub-tab Achievement tab Limits Prop — `Activity/GetAchievement.xml`,
// dijalankan tombol Refresh dan pilihan Quarter Year.
//
//	[3]–[5]  baris produksi `ACHIEVEMENT` + daftar kuartal & tahun kuartal
//	[8]      saringan As At Quarter / Quarter Year (hanya `search`)
//	[11]     tiap baris ditempel ke Detail ber-TreatyType (Kind of Treaty)
//	         DAN TreatyGroup sama; Detail kembar berturutan mendapat baris nol
//	[11.4]   Cash Call, Incured, Total, Loss Ratio, nilai IDR per baris
//	[12.2]   total per Detail, Achievement %, grid Gross/Nett, baris
//	         " Total In IDR"
//	[13]     FlagExcel — Generate Excel / Submit tampil
//
// ⚠️ TIGA YANG TIDAK DITIRU, dan sebabnya:
//  1. Klaim (`GetHistoryClaim_act`/`GetEstimasiClaim_act`) — sumbernya
//     `OS_AKSEPTASI_KLAIM.DATA_JSON`; nilai dari JSON dilarang pemilik proses
//     sampai ada keputusan. Cash Call = Estimation = 0.
//  2. [12.3] menghapus Detail ber-TreatyGroup kosong dan [12.4] menambah
//     Detail tiruan " Total In IDR" ke grid Treaty Group — keduanya
//     MENGUBAH pohon Limits; jumlah tingkat Limits dikembalikan terpisah.
//  3. Achievement % memerlukan `TreatyIn.RNMShareP` (tab Share), yang belum
//     punya tabel pendaratan; tanpa nilainya Achievement % dibiarkan kosong.
//
// ⚠️ Saringan [8.1]: kode langkah dibaca sebagai "hapus bila kuartal > As
// At ATAU tahun > Quarter Year" — kuartal dibandingkan TANPA tahun, persis
// ekspor.

import (
	"context"
	"sort"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/models"
)

// DetailAchievement - medan Detail yang GetAchievement baca.
type DetailAchievement struct {
	TreatyGroup string          `json:"TreatyGroup"`
	EPIList     []NilaiMataUang `json:"EPIList"`
}

// LimitAchievement - satu Kind of Treaty.
type LimitAchievement struct {
	TreatyType string              `json:"TreatyType"`
	Detail     []DetailAchievement `json:"Detail"`
}

// MasukanAchievement - parameter Refresh / Quarter Year.
type MasukanAchievement struct {
	IDKontrak string `json:"idKontrak"`
	// `param.search == "search"` — Quarter Year dipilih.
	Cari bool `json:"cari"`
	// `SearchData.CARI1` (As At Quarter) dan `CARI2` (Quarter Year).
	AsAt  string `json:"asAt"`
	Tahun string `json:"tahun"`
	// `TreatyIn.RNMShareP` — kosong bila tidak diketahui.
	RNMShareP string             `json:"rnmShareP"`
	Limits    []LimitAchievement `json:"limits"`
}

// HasilDetailAchievement - medan Detail yang GetAchievement tulis.
type HasilDetailAchievement struct {
	AchievementLists   []map[string]string `json:"AchievementLists"`
	CurrencyList       []map[string]string `json:"CurrencyList"`
	AchievementPct     string              `json:"AchievementPct"`
	LossRatio          string              `json:"LossRatio"`
	TotalAchPremium    string              `json:"TotalAchPremium"`
	TotalAchPaid       string              `json:"TotalAchPaid"`
	TotalAchOuts       string              `json:"TotalAchOuts"`
	TotalAchIncured    string              `json:"TotalAchIncured"`
	TotalAchNetPremium string              `json:"TotalAchNetPremium"`
}

// HasilAchievement - per Limits per Detail, plus daftar pilihan dan bendera.
type HasilAchievement struct {
	Limits       [][]HasilDetailAchievement `json:"limits"`
	Kuartal      []string                   `json:"kuartal"`
	TahunKuartal []string                   `json:"tahunKuartal"`
	FlagExcel    bool                       `json:"flagExcel"`
	// Ringkasan - isi tab `Achievement In IDR` (`AchievementCombine.xml`),
	// proyeksi atas `Limits` di atas. Lihat `RingkasAchievement`.
	Ringkasan RingkasanAchievement `json:"ringkasan"`
}

// GetAchievement — baca sumbernya lalu hitung.
func (l *Layanan) GetAchievement(ctx context.Context, p inti.Pelaku, m MasukanAchievement) (HasilAchievement, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return HasilAchievement{}, err
	}
	baris, err := l.gudang.BacaAchievement(ctx, m.IDKontrak)
	if err != nil {
		return HasilAchievement{}, err
	}
	var ids []string
	for _, b := range baris {
		ids = append(ids, b.IDCurrency)
	}
	kurs, err := l.gudang.BacaKursKeIDR(ctx, ids)
	if err != nil {
		return HasilAchievement{}, err
	}
	return HitungAchievement(m, baris, kurs), nil
}

func unik(s []string) []string {
	lihat := map[string]bool{}
	out := []string{}
	for _, v := range s {
		if !lihat[v] {
			lihat[v] = true
			out = append(out, v)
		}
	}
	sort.Strings(out) // `SELECT DISTINCT … ORDER BY … ASC`
	return out
}

// HitungAchievement — pure; `kurs` = TOIDR terbaru per IDCURRENCY.
func HitungAchievement(m MasukanAchievement, baris []models.BarisAchievement, kurs map[string]string) HasilAchievement {
	h := HasilAchievement{Kuartal: []string{}, TahunKuartal: []string{}, Limits: [][]HasilDetailAchievement{}}
	var q, y []string
	for _, b := range baris {
		q = append(q, b.Quarter)
		y = append(y, b.QuarterYear)
	}
	h.Kuartal, h.TahunKuartal = unik(q), unik(y)

	// [8] saringan
	pakai := []models.BarisAchievement{}
	for _, b := range baris {
		if m.Cari && (b.Quarter > m.AsAt || b.QuarterYear > m.Tahun) {
			continue
		}
		pakai = append(pakai, b)
	}

	// [11] tempel ke Detail
	lists := make([][][]map[string]string, len(m.Limits))
	for i, l := range m.Limits {
		lists[i] = make([][]map[string]string, len(l.Detail))
	}
	for _, b := range pakai {
		for i, l := range m.Limits {
			if b.TreatyType != l.TreatyType {
				continue
			}
			for k, d := range l.Detail {
				if b.TreatyGroupName != d.TreatyGroup {
					continue
				}
				kembar := k > 0 && l.Detail[k-1].TreatyGroup == d.TreatyGroup
				lists[i][k] = append(lists[i][k], barisAchievement(b, kembar))
			}
		}
	}

	satu := apd.New(1, 0)
	for i, l := range m.Limits {
		hasilL := []HasilDetailAchievement{}
		for k, d := range l.Detail {
			daftar := lists[i][k]
			// [11.1.2.4] per baris
			for _, r := range daftar {
				konv := angka(kurs[r["CurrencyID"]])
				if r["Currency"] == "IDR" {
					konv = satu
				}
				r["Conversion"] = teks(konv)
				// Klaim tidak dibaca — CASHCALL/Estimation 0 (lihat kepala berkas).
				r["CASHCALL"], r["EstimationCASHCALL"], r["OutstandingCASHCALL"] = "0", "0", "0"
				inc := tambah(tambah(angka(r["PaidClaim"]), angka(r["CASHCALL"])), tambah(angka(r["OutstandingClaim"]), angka(r["OutstandingCASHCALL"])))
				net := angka(r["NETPREMIUM"])
				r["IncuredClaim"] = teks(inc)
				r["Total"] = teks(kurang(net, inc))
				if net.IsZero() {
					r["LossRatio"] = "0"
				} else {
					r["LossRatio"] = teks(kali(bagiBulatDes(inc, net, 20), seratus))
				}
				for _, k2 := range [][2]string{
					{"PREMIUMtoIDR", "PREMIUM"}, {"PaidClaimtoIDR", "PaidClaim"}, {"OutstandingClaimtoIDR", "OutstandingClaim"},
					{"IncuredClaimtoIDR", "IncuredClaim"}, {"NETPREMIUMtoIDR", "NETPREMIUM"}, {"CASHCALLtoIDR", "CASHCALL"},
					{"EstimationCASHCALLtoIDR", "EstimationCASHCALL"}, {"OutstandingCASHCALLtoIDR", "OutstandingCASHCALL"},
					{"TotaltoIDR", "Total"},
				} {
					r[k2[0]] = teks(kali(angka(r[k2[1]]), konv))
				}
			}
			// [12.2] total per Detail
			jml := func(kunci string) *apd.Decimal {
				t := apd.New(0, 0)
				for _, r := range daftar {
					if r["TREATYGROUPNAME"] == d.TreatyGroup {
						t = tambah(t, angka(r[kunci]))
					}
				}
				return t
			}
			prem, paid, outs := jml("PREMIUMtoIDR"), jml("PaidClaimtoIDR"), jml("OutstandingClaimtoIDR")
			inc, netp := jml("IncuredClaimtoIDR"), jml("NETPREMIUMtoIDR")
			cc, ecc, occ, oscl, aft := jml("CASHCALLtoIDR"), jml("EstimationCASHCALLtoIDR"), jml("OutstandingCASHCALLtoIDR"), jml("OutstandingClaimtoIDR"), jml("TotaltoIDR")
			epi := ""
			for _, e := range d.EPIList {
				if e.Currency == "IDR" {
					epi = e.Value
				}
			}
			persen := func(pembilang *apd.Decimal) string {
				share := bagiBulatDes(angka(m.RNMShareP), seratus, 20)
				if epi == "" {
					return "0"
				}
				if share.IsZero() || angka(epi).IsZero() {
					return "" // RNMShareP tidak diketahui — lihat kepala berkas
				}
				return teks(kali(bagiBulatDes(bagiBulatDes(pembilang, share, 20), angka(epi), 20), seratus))
			}
			lr := func(pembagi *apd.Decimal) string {
				if prem.IsZero() || pembagi.IsZero() {
					return "0"
				}
				return teks(kali(bagiBulatDes(inc, pembagi, 20), seratus))
			}
			hd := HasilDetailAchievement{
				AchievementPct: persen(prem), LossRatio: lr(prem),
				TotalAchPremium: teks(prem), TotalAchPaid: teks(paid), TotalAchOuts: teks(outs),
				TotalAchIncured: teks(inc), TotalAchNetPremium: teks(netp),
				CurrencyList: []map[string]string{
					{"Parameter": " Based on Gross ", "AchievementPctGross": persen(prem), "LossRatioGross": lr(prem)},
					{"Parameter": "Based on Nett", "AchievementPctGross": persen(netp), "LossRatioGross": lr(netp)},
				},
			}
			// [12.2.6] baris total
			daftar = append(daftar, map[string]string{
				"Currency": " Total In IDR", "PREMIUM": teks(prem), "PaidClaim": teks(paid), "NETPREMIUM": teks(netp),
				"CASHCALL": teks(cc), "EstimationCASHCALL": teks(ecc), "OutstandingCASHCALL": teks(occ),
				"OutstandingClaim": teks(oscl), "Total": teks(aft), "IncuredClaim": teks(inc),
			})
			hd.AchievementLists = daftar
			if len(daftar) > 0 {
				h.FlagExcel = true // [13]
			}
			hasilL = append(hasilL, hd)
		}
		h.Limits = append(h.Limits, hasilL)
	}
	h.Ringkasan = RingkasAchievement(m, h)
	return h
}

// barisAchievement — [11.1.2.2] baris penuh, atau [11.1.2.3] baris nol bagi
// Detail yang grupnya sama dengan Detail sebelumnya.
func barisAchievement(b models.BarisAchievement, nol bool) map[string]string {
	r := map[string]string{
		"IDPEGA": b.IDPega, "NOOFFER": b.NoOffer, "SOBNAME": b.SoBName, "TREATYGROUPNAME": b.TreatyGroupName,
		"TREATYTYPE": b.TreatyType, "Quarter": "Q " + b.Quarter, "QUARTERYEAR": b.QuarterYear,
		"CurrencyID": b.IDCurrency, "Currency": b.Currency, "NOPOLIS": b.NoPolis,
		"PREMIUM": b.Premium, "RICOMM": b.RIComm, "BROKERAGE": b.Brokerage, "NETPREMIUM": b.NetPremium,
		"PaidClaim": b.PaidClaim, "OutstandingClaim": b.OutstandingClaim, "FlagTreaty": "0",
	}
	if nol {
		for _, k := range []string{"PREMIUM", "RICOMM", "BROKERAGE", "NETPREMIUM", "PaidClaim", "OutstandingClaim"} {
			r[k] = "0"
		}
		r["FlagTreaty"] = "1"
	}
	return r
}

// --- Tab `Achievement In IDR` (`Section/AchievementCombine.xml`) -------------

// BarisRingkasAchievement - satu baris grid tab `Achievement In IDR`.
//
// Keenam kolomnya, urut ekspor:
//
//	Treaty Group                        .TreatyGroup
//	Reins Type                          .TreatyType
//	Gross Premium Before Claim in IDR   .TotalAchPremium
//	Net Premium Before Claim in IDR     .TotalAchNetPremium
//	Incured Claim in IDR                .TotalAchIncured
//	Net Loss Ratio                      .LossRatio
//
// ⛔ SELURUHNYA hanya-baca: keenam sel ber-`pyReadOnly = true` TANPA
// `pyReadOnlyCondition`. Dibedakan dengan sengaja dari tab EGNPI dan Maximum
// Retention, yang sel serupanya BERSYARAT dan karena itu aktif di mode Edit.
type BarisRingkasAchievement struct {
	TreatyGroup        string `json:"TreatyGroup"`
	TreatyType         string `json:"TreatyType"`
	TotalAchPremium    string `json:"TotalAchPremium"`
	TotalAchNetPremium string `json:"TotalAchNetPremium"`
	TotalAchIncured    string `json:"TotalAchIncured"`
	LossRatio          string `json:"LossRatio"`
}

// RingkasanAchievement - isi tab `Achievement In IDR`: grid + tiga total kaki.
type RingkasanAchievement struct {
	Baris []BarisRingkasAchievement `json:"baris"`
	// Ketiga sel kaki `AchievementCombine` yang `pyVisible = ALWAYS`.
	SumTotalAchievNetPremium string `json:"SumTotalAchievNetPremium"`
	SumTotalAchievIncured    string `json:"SumTotalAchievIncured"`
	SumLossRatio             string `json:"SumLossRatio"`
}

// RingkasAchievement - `GetAchievement` langkah [12.3]-[12.4], sebagai
// PROYEKSI atas hasil per-Detail.
//
// ⭐ Mengapa proyeksi dan bukan penulisan ke pohon: ekspor menambahkan baris
// tiruan ke `TreatyIn.Limits().Detail()` dan menghapus Detail bergrup kosong
// dari pohon yang SAMA yang dipakai tab Limits. Menirunya apa adanya membuat
// tab Limits ikut berubah hanya karena tab Achievement dibuka. Bentuk dan
// angkanya sama; yang berbeda hanya pohonnya tidak disentuh.
//
// ⛔ [12.3] Detail ber-`TreatyGroup` KOSONG dibuang dari grid.
//
// ⚠️ [12.4] menimpa `.SumLossRatio` dengan pembagian, membuang jumlah yang
// baru saja ditumpuk [12.2.7] (`Local.SumLossRatio + .LossRatio`). Jadi
// total Loss Ratio BUKAN jumlah Loss Ratio per baris, melainkan dihitung
// ulang dari kedua total. Menjumlahkan persen memang tidak bermakna.
func RingkasAchievement(m MasukanAchievement, h HasilAchievement) RingkasanAchievement {
	r := RingkasanAchievement{Baris: []BarisRingkasAchievement{}}
	jmlPrem, jmlNet, jmlInc := apd.New(0, 0), apd.New(0, 0), apd.New(0, 0)

	for i, l := range m.Limits {
		if i >= len(h.Limits) {
			break
		}
		for j, d := range l.Detail {
			if j >= len(h.Limits[i]) {
				break
			}
			// [12.3]
			if d.TreatyGroup == "" {
				continue
			}
			hd := h.Limits[i][j]
			r.Baris = append(r.Baris, BarisRingkasAchievement{
				TreatyGroup: d.TreatyGroup, TreatyType: l.TreatyType,
				TotalAchPremium: hd.TotalAchPremium, TotalAchNetPremium: hd.TotalAchNetPremium,
				TotalAchIncured: hd.TotalAchIncured, LossRatio: hd.LossRatio,
			})
			// [12.2.7] — syaratnya `Local.TreatyType == .TreatyType`, dan
			// `.TreatyType` sebuah Detail DATANG dari Limits induknya. Jadi
			// syarat itu selalu benar; dinyatakan, bukan dibuang diam-diam.
			jmlPrem = tambah(jmlPrem, angka(hd.TotalAchPremium))
			jmlNet = tambah(jmlNet, angka(hd.TotalAchNetPremium))
			jmlInc = tambah(jmlInc, angka(hd.TotalAchIncured))
		}
	}

	r.SumTotalAchievNetPremium = teks(jmlNet)
	r.SumTotalAchievIncured = teks(jmlInc)
	if jmlNet.IsZero() {
		r.SumLossRatio = "0"
	} else {
		r.SumLossRatio = teks(kali(bagiBulatDes(jmlInc, jmlNet, 20), seratus))
	}

	// [12.4] baris " Total In IDR" — `TreatyType`-nya, bukan `TreatyGroup`.
	if len(r.Baris) > 0 {
		r.Baris = append(r.Baris, BarisRingkasAchievement{
			TreatyType:      " Total In IDR",
			TotalAchPremium: teks(jmlPrem), TotalAchNetPremium: r.SumTotalAchievNetPremium,
			TotalAchIncured: r.SumTotalAchievIncured, LossRatio: r.SumLossRatio,
		})
	}
	return r
}
