package services

// Rumus tab Limits NON-PROPORSIONAL — disalin dari Activity dan Data
// Transform `D:\XML_NURE\Treaty In`, langkah demi langkah:
//
//	TotalEgnpi                 Egnpi this layer
//	SetReinstatementPct        grid Reinstatement (juga saat kontrak dibuka —
//	                           `TreatySetReinstatement`)
//	CalculateReinstatement     Reinstatement Amount IDR/USD
//	CalculateReinstatementPct  % dari Reinstatement Amount IDR
//	ReCalculateReinstatement   Reinstatement Premium Amount IDR/USD
//	DetailCalculation adj/mdp  Premium Earned · MDP · Min Premium
//	DetailCalculationROL       ROL %
//	TreatyInNPSetTotal limits  Total All Layers
//	TreatyInSummaryLimit       Summary of Limit
//	TreatyInLimitsListValue    Update Value in List (kontrak revisi)
//
// ⭐ DIUKUR terhadap data Pega yang tersimpan (6 Oktober 2026):
//   - `TotalEgnpi` — EGNPI dicocokkan lewat `.TreatyGroup` baris
//     `TreatyGroupList`: 5.020 dari 5.069 baris `EgnpiTotalList` (99,0%).
//     Prasyarat grup `.ClassOfBusiness==Local.cob` TIDAK menyaring per
//     baris — 3.650 dari 3.761 baris EGNPI ber-ClassOfBusiness kosong, dan
//     totalnya tetap terisi.
//   - Total Premium Earned / MDP = jumlah biasa per mata uang (prasyarat
//     `.EgnpiCurr` tidak menyaring) — cocok 1.303/2.059 dan 1.314/2.059;
//     sisanya total yang tidak dihitung ulang sesudah nilainya berubah.
//
// ⚠️ TIGA SIMPANGAN dari teks ekspor, ketiganya salah ALAMAT, bukan rumus:
//  1. `ReCalculateReinstatement` menunjuk layer dengan `.pxListSubscript`
//     baris REINSTATEMENT (bukan layer) — mengambil MDP layer lain. Di sini
//     dipakai layer milik baris itu (`Param.ID`, yang DT-nya sendiri isi
//     lalu tidak pakai).
//  2. `TreatyInNPSetTotal(limits)` melompat ke cabang Installment dan
//     MENGHAPUS `TotalInstallmentNP` — efek samping ke tab lain; tidak
//     ditiru.
//  3. `TotalLimitMDPMinNP` (tak tampil) tidak dihitung.
//
// ⚠️ Yang DITIRU APA ADANYA walau janggal, sebab ia rumus, bukan alamat:
//   - relasi mata uang `AND` menjumlahkan limit DUA kali di ROL;
//   - `CalculateReinstatementPct` menulis `ReinstatementPct` (kolom
//     "% Additional Premium") dari Reinstatement Amount IDR ÷ Limit;
//   - ROL = 9.989.998 × 100 bila limit IDR-nya nol;
//   - total 100% Limit/Deductible memakai mata uang LAYER PERTAMA.

import (
	"encoding/json"
	"strings"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
)

// GrupLayerNP - satu baris `TreatyGroupList`.
type GrupLayerNP struct {
	TreatyGroup   string `json:"TreatyGroup"`
	TreatyGroupID string `json:"TreatyGroupID"`
}

// BarisReinstatement - satu baris `Reinstatement_List`.
type BarisReinstatement struct {
	ReinstatementValue   string `json:"ReinstatementValue"`
	ReinstatementPct     string `json:"ReinstatementPct"`
	ReinstatementNote    string `json:"ReinstatementNote"`
	AdditionalAmount1    string `json:"AdditionalAmount1"`
	AdditionalAmount2    string `json:"AdditionalAmount2"`
	AdditionalPct        string `json:"AdditionalPct"`
	ReinstatementAmount1 string `json:"ReinstatementAmount1"`
	ReinstatementAmount2 string `json:"ReinstatementAmount2"`
	ID                   string `json:"ID"`
}

// LayerNP - medan satu `Limits[]` Non-Prop yang rumus tab Limits baca/tulis.
// Ejaan JSON = ejaan dokumen.
type LayerNP struct {
	LayerType          string               `json:"LayerType"`
	Layer              string               `json:"Layer"`
	LayerPartType      string               `json:"LayerPartType"`
	LayerPart          string               `json:"LayerPart"`
	Cover              string               `json:"Cover"`
	Currency           string               `json:"Currency"`
	Currency2          string               `json:"Currency2"`
	CurrencyRelation   string               `json:"CurrencyRelation"`
	Limit              string               `json:"Limit"`
	Limit2             string               `json:"Limit2"`
	AgregateLimit      string               `json:"AgregateLimit"`
	AgregateLimit2     string               `json:"AgregateLimit2"`
	Deductible         string               `json:"Deductible"`
	Deductible2        string               `json:"Deductible2"`
	AdjRate            string               `json:"AdjRate"`
	MDPPct             string               `json:"MDPPct"`
	MDPMinPct          string               `json:"MDPMinPct"`
	ROLPct             string               `json:"ROLPct"`
	ReinstatementValue string               `json:"ReinstatementValue"`
	ReinstatementNote  string               `json:"ReinstatementNote"`
	TreatyGroupList    []GrupLayerNP        `json:"TreatyGroupList"`
	EgnpiTotalList     []NilaiMataUang      `json:"EgnpiTotalList"`
	PremiumEarnedList  []NilaiMataUang      `json:"PremiumEarnedList"`
	MDPList            []NilaiMataUang      `json:"MDPList"`
	MDPMinList         []NilaiMataUang      `json:"MDPMinList"`
	ReinstatementList  []BarisReinstatement `json:"Reinstatement_List"`
}

// EgnpiNP - satu baris tab EGNPI.
//
// ⭐ SATU tipe untuk DUA pemakai, dan itu disengaja: `TotalEgnpi` (tab
// Limits) hanya membaca empat medan pertama, sementara tab EGNPI sendiri
// (`hitung_egnpi.go`) memakai seluruhnya. Dua tipe untuk satu baris layar
// adalah cara termudah keduanya berbeda diam-diam.
//
// ⚠️ `Proportion` DIHITUNG, bukan diisi tangan - `TreatyInNPSetTotal`
// menimpanya setiap kali `Update Total` ditekan. Lihat `NPSetTotalEgnpi`.
type EgnpiNP struct {
	TreatyGroup string `json:"TreatyGroup"`
	Currency    string `json:"Currency"`
	CurrencyID  string `json:"CurrencyID"`
	Amount      string `json:"Amount"`

	// Medan tab EGNPI - `Section/DetailEGNPI.xml` dan grid tab EGNPI.
	ID                string `json:"ID"`
	TreatyGroupID     string `json:"TreatyGroupID"`
	AsDate            string `json:"AsDate"`
	Proportion        string `json:"Proportion"`
	AmountIDR         string `json:"AmountIDR"`
	ClassOfBusiness   string `json:"ClassOfBusiness"`
	ClassOfBusinessID string `json:"ClassOfBusinessID"`
	Note              string `json:"Note"`
}

// KursNP - satu baris `TreatyIn.CurrencyList` (grid Rate of Exchange).
type KursNP struct {
	Currency   string `json:"Currency"`
	Conversion string `json:"Conversion"`
}

// RingkasanLimitNP - satu baris `LimitSummaryList`.
type RingkasanLimitNP struct {
	LayerType       string `json:"LayerType"`
	Layer           string `json:"Layer"`
	LayerPartType   string `json:"LayerPartType"`
	LayerPart       string `json:"LayerPart"`
	Currency        string `json:"Currency"`
	Currency2       string `json:"Currency2"`
	Limit           string `json:"Limit"`
	Limit2          string `json:"Limit2"`
	AggregateLimit  string `json:"AggregateLimit"`
	AggregateLimit2 string `json:"AggregateLimit2"`
	Deductible      string `json:"Deductible"`
	Deductible2     string `json:"Deductible2"`
	MDP             string `json:"MDP"`
	MDP2            string `json:"MDP2"`
	Note            string `json:"Note"`
}

// Aksi tab Limits Non-Prop.
const (
	AksiNPEgnpi          = "egnpi"           // TotalEgnpi
	AksiNPReinstatement  = "reinstatement"   // SetReinstatementPct(Indeks)
	AksiNPReinstJumlah   = "reinst-jumlah"   // CalculateReinstatement(Indeks, Baris)
	AksiNPReinstPersen   = "reinst-persen"   // CalculateReinstatementPct(Indeks, Baris)
	AksiNPReinstTambahan = "reinst-tambahan" // ReCalculateReinstatement(Indeks, Baris)
	AksiNPAdj            = "adj"             // DetailCalculation(adj) → ROL
	AksiNPMDP            = "mdp"             // DetailCalculation(mdp) → ROL
	AksiNPTotal          = "total"           // NPSetTotal(limits) → SummaryLimit
	AksiNPNilaiList      = "nilai-list"      // TreatyInLimitsListValue
)

// MasukanLimitNP - aksi + seluruh `Limits[]` + EGNPI + kurs kontrak.
type MasukanLimitNP struct {
	Aksi   string    `json:"aksi"`
	Layers []LayerNP `json:"layers"`
	EGNPI  []EgnpiNP `json:"egnpi"`
	Kurs   []KursNP  `json:"kurs"`
	// Indeks layer (mulai 0) dan indeks baris reinstatement (mulai 0).
	Indeks int `json:"indeks"`
	Baris  int `json:"baris"`
}

// HasilLimitNP - `Limits[]` sesudah aksi, ditambah total dan ringkasan bila
// aksinya menghitungnya (`total`, `nilai-list`).
type HasilLimitNP struct {
	Layers           []LayerNP                  `json:"layers"`
	Total            map[string][]NilaiMataUang `json:"Total,omitempty"`
	TotalLimitsROL   string                     `json:"TotalLimitsROL,omitempty"`
	LimitSummaryList []RingkasanLimitNP         `json:"LimitSummaryList,omitempty"`
	// Pesan `Property-Set-Messages`, apa adanya.
	Pesan []string `json:"pesan"`
}

// Pesan Activity apa adanya.
const (
	PesanLimitKosong = "100% Limit must not be empty"
	PesanMDPKosong   = "MDP% must not be empty"
)

// HitungLimitNP — bentuk ber-pelaku untuk handler.
func (l *Layanan) HitungLimitNP(p inti.Pelaku, m MasukanLimitNP) (HasilLimitNP, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return HasilLimitNP{}, err
	}
	return HitungLimitNP(m), nil
}

// HitungLimitNP menjalankan satu aksi tab Limits Non-Prop. Pure; masukan
// tidak diubah.
func HitungLimitNP(m MasukanLimitNP) HasilLimitNP {
	h := HasilLimitNP{Layers: salinLayers(m.Layers), Pesan: []string{}}
	ada := m.Indeks >= 0 && m.Indeks < len(h.Layers)
	switch m.Aksi {
	case AksiNPEgnpi:
		TotalEgnpi(h.Layers, m.EGNPI)
	case AksiNPReinstatement:
		if ada {
			SetReinstatementPct(&h.Layers[m.Indeks], m.Indeks+1)
		}
	case AksiNPReinstJumlah, AksiNPReinstPersen, AksiNPReinstTambahan:
		if ada && m.Baris >= 0 && m.Baris < len(h.Layers[m.Indeks].ReinstatementList) {
			hitungBarisReinstatement(&h.Layers[m.Indeks], m.Baris, m.Aksi)
		}
	case AksiNPAdj, AksiNPMDP:
		if ada {
			h.Pesan = append(h.Pesan, DetailCalculation(&h.Layers[m.Indeks], m.Aksi, m.Kurs)...)
			if m.Aksi == AksiNPMDP {
				segarkanReinstatement(&h.Layers[m.Indeks])
			}
		}
	case AksiNPTotal:
		h.Total, h.TotalLimitsROL = NPSetTotalLimits(h.Layers)
		h.LimitSummaryList = SummaryLimit(h.Layers)
	case AksiNPNilaiList:
		// `TreatyInLimitsListValue`: per layer TotalEgnpi + adj + mdp, lalu
		// total dan ringkasan. `TotalEgnpi` menapaki SEMUA layer sekaligus,
		// jadi sekali sudah cukup.
		TotalEgnpi(h.Layers, m.EGNPI)
		for i := range h.Layers {
			h.Pesan = append(h.Pesan, DetailCalculation(&h.Layers[i], AksiNPAdj, m.Kurs)...)
			h.Pesan = append(h.Pesan, DetailCalculation(&h.Layers[i], AksiNPMDP, m.Kurs)...)
			segarkanReinstatement(&h.Layers[i])
		}
		h.Total, h.TotalLimitsROL = NPSetTotalLimits(h.Layers)
		h.LimitSummaryList = SummaryLimit(h.Layers)
	}
	return h
}

func salinLayers(src []LayerNP) []LayerNP {
	out := make([]LayerNP, len(src))
	for i, l := range src {
		k := l
		k.TreatyGroupList = append([]GrupLayerNP{}, l.TreatyGroupList...)
		k.EgnpiTotalList = append([]NilaiMataUang{}, l.EgnpiTotalList...)
		k.PremiumEarnedList = append([]NilaiMataUang{}, l.PremiumEarnedList...)
		k.MDPList = append([]NilaiMataUang{}, l.MDPList...)
		k.MDPMinList = append([]NilaiMataUang{}, l.MDPMinList...)
		k.ReinstatementList = append([]BarisReinstatement{}, l.ReinstatementList...)
		out[i] = k
	}
	return out
}

// tambahPerMataUang menjumlah `nilai` ke baris bermata uang sama, atau
// menambah baris baru — pola `Appendflag` di seluruh Activity total.
func tambahPerMataUang(daftar []NilaiMataUang, cur, curID string, nilai *apd.Decimal) []NilaiMataUang {
	for i := range daftar {
		if daftar[i].Currency == cur {
			daftar[i].Value = teks(tambah(angka(daftar[i].Value), nilai))
			return daftar
		}
	}
	return append(daftar, NilaiMataUang{Currency: cur, CurrencyID: curID, Value: teks(nilai)})
}

// TotalEgnpi — `Activity/TotalEgnpi.xml`: untuk SETIAP layer, jumlah
// `Amount` EGNPI yang `.TreatyGroup`-nya sama dengan baris `TreatyGroupList`
// layer itu, per mata uang, urut kemunculan. Grup yang muncul dua kali
// dihitung dua kali — begitu Activity-nya.
func TotalEgnpi(layers []LayerNP, egnpi []EgnpiNP) {
	for i := range layers {
		layers[i].EgnpiTotalList = []NilaiMataUang{}
		for _, g := range layers[i].TreatyGroupList {
			for _, e := range egnpi {
				if e.TreatyGroup != g.TreatyGroup {
					continue
				}
				layers[i].EgnpiTotalList = tambahPerMataUang(layers[i].EgnpiTotalList, e.Currency, e.CurrencyID, angka(e.Amount))
			}
		}
	}
}

// SetReinstatementPct — `Activity/SetReinstatementPct.xml`. `idLayer` =
// `.pxListSubscript` layer (mulai 1).
func SetReinstatementPct(l *LayerNP, idLayer int) {
	l.ReinstatementList = []BarisReinstatement{}
	n := bulat(angka(l.ReinstatementValue))
	for i := 1; i <= n; i++ {
		b := BarisReinstatement{
			ReinstatementPct:     "100",
			ReinstatementValue:   itoa(i),
			ReinstatementNote:    l.ReinstatementNote,
			ReinstatementAmount1: l.Limit,
			ReinstatementAmount2: l.Limit2,
			ID:                   itoa(idLayer),
			AdditionalPct:        "100",
		}
		switch {
		case len(l.MDPList) == 1:
			b.AdditionalAmount1 = jikaMataUang(l.MDPList[0], "IDR")
			b.AdditionalAmount2 = jikaMataUang(l.MDPList[0], "USD")
		case len(l.MDPList) >= 2:
			b.AdditionalAmount1 = jikaMataUang(l.MDPList[0], "IDR")
			b.AdditionalAmount2 = jikaMataUang(l.MDPList[1], "USD")
		}
		l.ReinstatementList = append(l.ReinstatementList, b)
	}
}

// SiapkanReinstatement — `Activity/TreatySetReinstatement.xml`, dijalankan
// `SetTreatyIn_Act` langkah 13 saat kontrak Non-Prop dibuka: layer yang
// `Reinstatement_List(1).ReinstatementValue` kosong dibangun ulang.
func SiapkanReinstatement(layers []LayerNP) {
	for i := range layers {
		if len(layers[i].ReinstatementList) == 0 || layers[i].ReinstatementList[0].ReinstatementValue == "" {
			SetReinstatementPct(&layers[i], i+1)
		}
	}
}

// @if(MDP.Currency==cur, MDP.Value, 0)
func jikaMataUang(n NilaiMataUang, cur string) string {
	if n.Currency == cur {
		return n.Value
	}
	return "0"
}

// segarkanReinstatement — ⭐ SIMPANGAN WAKTU, BUKAN RUMUS (keputusan pemakai
// 8 Oktober 2026: *"lakukan saja rekomendasi terbaik asal rumusnya jangan
// berbeda"*). Sesudah MDP satu layer dihitung ulang (`DetailCalculation(mdp)`,
// juga di `TreatyInLimitsListValue`), DT `ReCalculateReinstatement` dijalankan
// untuk SETIAP baris Reinstatement — rumus yang SAMA persis
// (`hitungBarisReinstatement`). Di Pega DT itu hanya berjalan saat
// `% Additional Premium` berubah, sehingga grid yang dibangun SEBELUM MDP ada
// tetap kosong (laporan pemakai, tangkapan tab Limits Non-Prop).
// MDP kosong = tidak ada yang disegarkan (`SetReinstatementPct` 3.2/3.3).
func segarkanReinstatement(l *LayerNP) {
	if len(l.MDPList) == 0 {
		return
	}
	for i := range l.ReinstatementList {
		hitungBarisReinstatement(l, i, AksiNPReinstTambahan)
	}
}

// hitungBarisReinstatement — ketiga DT grid Reinstatement.
func hitungBarisReinstatement(l *LayerNP, baris int, aksi string) {
	b := &l.ReinstatementList[baris]
	switch aksi {
	case AksiNPReinstJumlah: // CalculateReinstatement
		f := bagiBulatDes(angka(b.AdditionalPct), seratus, 12)
		b.ReinstatementAmount1 = teks(kali(f, angka(l.Limit)))
		b.ReinstatementAmount2 = teks(kali(f, angka(l.Limit2)))
	case AksiNPReinstPersen: // CalculateReinstatementPct
		lim := angka(l.Limit)
		if lim.IsZero() {
			return // pembagi nol — Pega tidak menjaganya; di sini nilai lama dibiarkan
		}
		b.ReinstatementPct = teks(kali(bagiBulatDes(angka(b.ReinstatementAmount1), lim, 12), seratus))
	case AksiNPReinstTambahan: // ReCalculateReinstatement
		f := bagiBulatDes(angka(b.ReinstatementPct), seratus, 12)
		nilaiJika := func(n NilaiMataUang, cur string) string {
			if n.Currency == cur {
				return teks(kali(f, angka(n.Value)))
			}
			return "0"
		}
		switch {
		case len(l.MDPList) == 1:
			b.AdditionalAmount1 = nilaiJika(l.MDPList[0], "IDR")
			b.AdditionalAmount2 = nilaiJika(l.MDPList[0], "USD")
		case len(l.MDPList) > 1:
			b.AdditionalAmount1 = nilaiJika(l.MDPList[0], "IDR")
			b.AdditionalAmount2 = nilaiJika(l.MDPList[1], "USD")
		}
	}
}

// DetailCalculation — `Activity/DetailCalculation.xml` cabang adj / mdp,
// lalu `DetailCalculationROL`. Mengembalikan pesan Activity.
//
// Cabang `copy` (langkah 1–9) hanya mengisi halaman sementara dan tidak
// punya pemanggil; `Param.TNOP` tidak pernah dikirim pemanggil mana pun.
func DetailCalculation(l *LayerNP, jenis string, kurs []KursNP) []string {
	pesan := []string{}
	switch jenis {
	case AksiNPAdj:
		// [11] adjrate := @divide(.AdjRate,100,20) · [13]–[14]
		r := bagiBulatDes(angka(l.AdjRate), seratus, 20)
		l.PremiumEarnedList = []NilaiMataUang{}
		for _, e := range l.EgnpiTotalList {
			l.PremiumEarnedList = append(l.PremiumEarnedList,
				NilaiMataUang{Currency: e.Currency, Value: teks(kali(angka(e.Value), r))})
		}
	case AksiNPMDP:
		// [17]–[21]
		l.MDPList = []NilaiMataUang{}
		l.MDPMinList = []NilaiMataUang{}
		if angka(l.MDPPct).Cmp(apd.New(1, 0)) < 0 {
			pesan = append(pesan, PesanMDPKosong)
		}
		r := bagiBulatDes(angka(l.MDPPct), seratus, 20)
		rMin := apd.New(0, 0)
		if !angka(l.MDPMinPct).IsZero() {
			rMin = bagiBulatDes(angka(l.MDPMinPct), seratus, 20)
		}
		for _, pe := range l.PremiumEarnedList {
			v := angka(pe.Value)
			l.MDPList = append(l.MDPList, NilaiMataUang{Currency: pe.Currency, Value: teks(kali(v, r))})
			l.MDPMinList = append(l.MDPMinList, NilaiMataUang{Currency: pe.Currency, Value: teks(kali(v, rMin))})
		}
	default:
		return pesan
	}
	return append(pesan, DetailCalculationROL(l, kurs)...)
}

// kursDenganIDR menambah baris IDR = 1 bila kurs kontrak tidak memuatnya.
//
// ⚠️ Grid Rate of Exchange aplikasi ini dibaca dari `TREATYEXCHANGEYEARLY`
// (keputusan pemilik proses 4 & 6 Oktober 2026) — tabel kurs KE IDR, yang
// karena itu tidak punya baris IDR. `DetailCalculationROL` menapaki
// `CurrencyList` dan hanya menghitung mata uang yang ADA di sana; tanpa
// baris IDR, layer berlimit IDR jatuh ke ROL 9.989.998 × 100. IDR → IDR = 1
// adalah identitas, bukan nilai karangan.
func kursDenganIDR(kurs []KursNP) []KursNP {
	for _, k := range kurs {
		if k.Currency == "IDR" {
			return kurs
		}
	}
	return append(append([]KursNP{}, kurs...), KursNP{Currency: "IDR", Conversion: "1"})
}

// DetailCalculationROL — `Activity/DetailCalculationROL.xml`.
func DetailCalculationROL(l *LayerNP, kurs []KursNP) []string {
	kurs = kursDenganIDR(kurs)
	satu := apd.New(1, 0)
	lim1, lim2 := angka(l.Limit), angka(l.Limit2)
	// [2] pesan + keluar
	if lim1.Cmp(satu) < 0 && lim2.Cmp(satu) < 0 {
		return []string{PesanLimitKosong}
	}
	// [4] limit → IDR per mata uang EGNPI layer
	totalLimit := apd.New(0, 0)
	for _, e := range l.EgnpiTotalList {
		for _, k := range kurs {
			if k.Currency != e.Currency {
				continue
			}
			konv := angka(k.Conversion)
			if k.Currency == "IDR" {
				totalLimit = tambah(totalLimit, lim1)
			} else {
				totalLimit = tambah(totalLimit, kali(lim2, konv))
			}
			if l.CurrencyRelation == "AND" {
				pilih := lim2
				if k.Currency == "IDR" {
					pilih = lim1
				}
				totalLimit = tambah(totalLimit, kali(pilih, konv))
			}
		}
	}
	// [6] premi → IDR
	totalPremi := apd.New(0, 0)
	for _, pe := range l.PremiumEarnedList {
		for _, k := range kurs {
			if k.Currency == pe.Currency {
				totalPremi = tambah(totalPremi, kali(angka(pe.Value), angka(k.Conversion)))
			}
		}
	}
	// [7]
	rasio := apd.New(9989998, 0)
	if !totalLimit.IsZero() {
		rasio = bagiBulatDes(totalPremi, totalLimit, 8)
	}
	l.ROLPct = teks(kali(rasio, seratus))
	return nil
}

// Kunci keempat total `Total All Layers`.
const (
	TotalLimitIOO       = "TotalLimitIOONP"
	TotalLimitDeductibl = "TotalLimitDeductblNP"
	TotalLimitPE        = "TotalLimitPremiEarnNP"
	TotalLimitMDP       = "TotalLimitMDPNP"
)

// NPSetTotalLimits — `Activity/TreatyInNPSetTotal.xml` cabang `limits`.
func NPSetTotalLimits(layers []LayerNP) (map[string][]NilaiMataUang, string) {
	t := map[string][]NilaiMataUang{
		TotalLimitIOO: {}, TotalLimitDeductibl: {}, TotalLimitPE: {}, TotalLimitMDP: {},
	}
	rol := apd.New(0, 0)
	if len(layers) == 0 {
		return t, ""
	}
	// [11]–[15] jumlah Limit/Limit2 dan Deductible/Deductible2 — mata uang
	// dari LAYER PERTAMA.
	sl, sl2, sd, sd2 := apd.New(0, 0), apd.New(0, 0), apd.New(0, 0), apd.New(0, 0)
	for _, l := range layers {
		sl = tambah(sl, angka(l.Limit))
		sl2 = tambah(sl2, angka(l.Limit2))
		sd = tambah(sd, angka(l.Deductible))
		sd2 = tambah(sd2, angka(l.Deductible2))
	}
	c1, c2 := layers[0].Currency, layers[0].Currency2
	t[TotalLimitIOO] = []NilaiMataUang{{Currency: c1, Value: teks(sl)}, {Currency: c2, Value: teks(sl2)}}
	t[TotalLimitDeductibl] = []NilaiMataUang{{Currency: c1, Value: teks(sd)}, {Currency: c2, Value: teks(sd2)}}
	// [16] mata uang unik dari MDPList semua layer
	var curs []string
	lihat := map[string]bool{}
	for _, l := range layers {
		for _, m := range l.MDPList {
			if !lihat[m.Currency] {
				lihat[m.Currency] = true
				curs = append(curs, m.Currency)
			}
		}
	}
	// [18]–[20]
	for _, c := range curs {
		pe, mdp := apd.New(0, 0), apd.New(0, 0)
		for _, l := range layers {
			for _, x := range l.PremiumEarnedList {
				if x.Currency == c {
					pe = tambah(pe, angka(x.Value))
				}
			}
			for _, x := range l.MDPList {
				if x.Currency == c {
					mdp = tambah(mdp, angka(x.Value))
				}
			}
		}
		t[TotalLimitPE] = append(t[TotalLimitPE], NilaiMataUang{Currency: c, Value: teks(pe)})
		t[TotalLimitMDP] = append(t[TotalLimitMDP], NilaiMataUang{Currency: c, Value: teks(mdp)})
	}
	// [21] Total ROL = Σ ROL% (jumlah, bukan rata-rata)
	for _, l := range layers {
		rol = tambah(rol, angka(l.ROLPct))
	}
	return t, teks(rol)
}

// SummaryLimit — `Activity/TreatyInSummaryLimit.xml`. Layer berkunci sama
// (LayerType, Layer, LayerPartType, LayerPart) dijumlahkan ke satu baris.
func SummaryLimit(layers []LayerNP) []RingkasanLimitNP {
	out := []RingkasanLimitNP{}
	sama := func(r RingkasanLimitNP, l LayerNP) bool {
		return r.LayerType == l.LayerType && r.Layer == l.Layer &&
			r.LayerPartType == l.LayerPartType && r.LayerPart == l.LayerPart
	}
	for _, l := range layers {
		i := -1
		for j := range out {
			if sama(out[j], l) {
				i = j
				break
			}
		}
		if i < 0 {
			out = append(out, RingkasanLimitNP{
				LayerType: l.LayerType, Layer: l.Layer, LayerPartType: l.LayerPartType, LayerPart: l.LayerPart,
				Currency: l.Currency, Currency2: l.Currency2,
				AggregateLimit: teks(angka(l.AgregateLimit)), AggregateLimit2: teks(angka(l.AgregateLimit2)),
				Limit: teks(angka(l.Limit)), Limit2: teks(angka(l.Limit2)),
				Deductible: teks(angka(l.Deductible)), Deductible2: teks(angka(l.Deductible2)),
			})
			i = len(out) - 1
		} else {
			r := &out[i]
			r.Deductible = teks(tambah(angka(r.Deductible), angka(l.Deductible)))
			r.Deductible2 = teks(tambah(angka(r.Deductible2), angka(l.Deductible2)))
			r.Limit = teks(tambah(angka(r.Limit), angka(l.Limit)))
			r.Limit2 = teks(tambah(angka(r.Limit2), angka(l.Limit2)))
			r.AggregateLimit = teks(tambah(angka(r.AggregateLimit), angka(l.AgregateLimit)))
			r.AggregateLimit2 = teks(tambah(angka(r.AggregateLimit2), angka(l.AgregateLimit2)))
		}
		// [3.5] MDP IDR → MDP, USD → MDP2; yang lain diabaikan.
		for _, m := range l.MDPList {
			switch m.Currency {
			case "IDR":
				out[i].MDP = teks(tambah(angka(out[i].MDP), angka(m.Value)))
			case "USD":
				out[i].MDP2 = teks(tambah(angka(out[i].MDP2), angka(m.Value)))
			}
		}
	}
	for i := range out {
		r := &out[i]
		r.Note = r.LayerType + r.Layer + " of " + r.LayerPartType + r.LayerPart
	}
	return out
}

// --- aritmetika ---------------------------------------------------------

func tambah(a, b *apd.Decimal) *apd.Decimal {
	r := new(apd.Decimal)
	_, _ = konteksLimit.Add(r, a, b)
	return r
}

// bagiBulatDes = `@divide(a, b, skala)` dengan pembagi desimal; pembagi nol
// menghasilkan nol (Pega tidak pernah sampai ke sini dengan pembagi nol
// pada jalur yang dibangun — ROL menjaganya sendiri).
func bagiBulatDes(a, b *apd.Decimal, skala int32) *apd.Decimal {
	r := new(apd.Decimal)
	if b.IsZero() {
		return r
	}
	_, _ = konteksLimit.Quo(r, a, b)
	_, _ = konteksLimit.Quantize(r, r, -skala)
	return r
}

// bulat — bagian bulat sebuah desimal (batas REPEAT Pega).
func bulat(d *apd.Decimal) int {
	r := new(apd.Decimal)
	_, _ = konteksLimit.RoundToIntegralValue(r, d)
	i, err := r.Int64()
	if err != nil || i < 0 {
		return 0
	}
	if i > 1000 { // pagar: jumlah reinstatement yang tak masuk akal
		return 1000
	}
	return int(i)
}

func itoa(i int) string {
	return strings.TrimSpace(apd.New(int64(i), 0).Text('f'))
}

// SiapkanReinstatementPohon — `SiapkanReinstatement` atas pohon layar
// (`[]map[string]any`): hanya `Reinstatement_List` yang ditulis.
func SiapkanReinstatementPohon(pohon []map[string]any) {
	for i, simpul := range pohon {
		if larik, ok := simpul["Reinstatement_List"].([]map[string]any); ok && len(larik) > 0 {
			if v, _ := larik[0]["ReinstatementValue"].(string); v != "" {
				continue
			}
		}
		var l LayerNP
		if b, err := json.Marshal(simpul); err != nil || json.Unmarshal(b, &l) != nil {
			continue
		}
		SetReinstatementPct(&l, i+1)
		baru := []map[string]any{}
		if b, err := json.Marshal(l.ReinstatementList); err == nil {
			_ = json.Unmarshal(b, &baru)
		}
		if baru == nil {
			baru = []map[string]any{}
		}
		simpul["Reinstatement_List"] = baru
	}
}
