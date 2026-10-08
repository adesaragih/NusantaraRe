package services

// Rumus sub-tab Deduction dan Reserve tab Limits PROPORSIONAL
// (`Section/DetailLimits.xml`):
//
//	CalculateDeduction       Deduction Details ↔ Deductions (total)
//	PremiumReserveCalculate  Premium Reserve dari Cession to R/I
//
// ⚠️ `CalculateDeduction` adalah Activity milik Share XOL (kelas
// `TreatyInShare`) yang DetailLimits pinjam. Ia membaca
// `.GrossPremiumList`, dan larik itu TIDAK ADA di rincian Limits Prop —
// sehingga di Pega pesannya selalu muncul dan cabang persen mengisi
// Deduction dengan nol serta MENGOSONGKAN mata uangnya. Yang dibangun di
// sini: pesan yang sama, total per mata uang yang sama, tetapi baris
// TIDAK dirusak ketika Gross kosong — merusak mata uang yang pemakai
// pilih bukan rumus, melainkan efek dari alamat yang tidak ada.

import (
	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/models"
)

// BarisDeduksi - satu baris `DeductionList`.
//
// Alias baris deduksi Share Non-Prop — `CalculateDeduction` milik kelas
// `TreatyInShare`, jadi satu bentuk melayani kedua tab. Bendera
// `DeductionPctCalculate` ber-`omitempty`: rincian Limits Prop tidak
// memilikinya dan jawabannya tidak berubah.
type BarisDeduksi = models.BarisDeduksiShare

// MasukanDeduksi - parameter `CalculateDeduction` + rincian yang dihitung.
type MasukanDeduksi struct {
	// `param.sts` — `val` (nilai diketik), `pct` (persen diketik), atau
	// kosong (mata uang berubah / baris dihapus).
	Sts string `json:"sts"`
	// `param.index` — baris (mulai 0).
	Indeks           int             `json:"indeks"`
	DeductionList    []BarisDeduksi  `json:"DeductionList"`
	GrossPremiumList []NilaiMataUang `json:"GrossPremiumList"`
	// BenderaPct — langkah 4 bertransisi `.DeductionList(param.index).
	// DeductionPctCalculate==true`: hanya bila benar langkah 5 (baris mata
	// uang berikutnya) dijalankan. Dipakai panel Share, yang barisnya
	// membawa bendera itu; rincian Limits Prop tidak, dan perilakunya tetap.
	BenderaPct bool `json:"-"`
}

// HasilDeduksi - larik yang `CalculateDeduction` tulis.
type HasilDeduksi struct {
	DeductionList      []BarisDeduksi  `json:"DeductionList"`
	DeductionTotalList []NilaiMataUang `json:"DeductionTotalList"`
	NetPremiumList     []NilaiMataUang `json:"NetPremiumList"`
	Pesan              []string        `json:"pesan"`
}

// PesanGrossKosong — `CalculateDeduction` langkah 1, apa adanya.
const PesanGrossKosong = "Gross Premium (MDP) is still empty"

// HitungDeduksi — bentuk ber-pelaku.
func (l *Layanan) HitungDeduksi(p inti.Pelaku, m MasukanDeduksi) (HasilDeduksi, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return HasilDeduksi{}, err
	}
	return HitungDeduksi(m), nil
}

// HitungDeduksi — `Activity/CalculateDeduction.xml`.
func HitungDeduksi(m MasukanDeduksi) HasilDeduksi {
	h := HasilDeduksi{
		DeductionList:      append([]BarisDeduksi{}, m.DeductionList...),
		DeductionTotalList: []NilaiMataUang{},
		NetPremiumList:     []NilaiMataUang{},
		Pesan:              []string{},
	}
	adaGross := len(m.GrossPremiumList) > 0
	// [1]–[2]
	if !adaGross {
		h.Pesan = append(h.Pesan, PesanGrossKosong)
	}
	sah := m.Indeks >= 0 && m.Indeks < len(h.DeductionList)
	// [3] nilai diketik → persen dibuang
	if m.Sts == "val" && sah {
		h.DeductionList[m.Indeks].DeductionPct = ""
	}
	// [4]–[5] persen diketik → nilai dari Gross, satu baris per mata uang
	if m.Sts == "pct" && sah && adaGross {
		b := h.DeductionList[m.Indeks]
		f := bagiBulat(angka(b.DeductionPct), 100, 4)
		g1 := m.GrossPremiumList[0]
		h.DeductionList[m.Indeks].Deduction = teks(kali(f, angka(g1.Value)))
		h.DeductionList[m.Indeks].Currency = g1.Currency
		h.DeductionList[m.Indeks].CurrencyID = g1.CurrencyID
		lanjut := !m.BenderaPct || b.DeductionPctCalculate == "true"
		for _, g := range m.GrossPremiumList[1:] {
			if !lanjut {
				break
			}
			h.DeductionList = append(h.DeductionList, BarisDeduksi{
				Comment: b.Comment, DeductionPct: b.DeductionPct,
				Currency: g.Currency, CurrencyID: g.CurrencyID,
				Deduction: teks(kali(f, angka(g.Value))),
			})
		}
	}
	// [6]–[7] total per mata uang, urut kemunculan
	for _, b := range h.DeductionList {
		h.DeductionTotalList = tambahPerMataUang(h.DeductionTotalList, b.Currency, "", angka(b.Deduction))
	}
	// [8]–[9] Net = Gross − Σ Deduction bermata uang sama
	for _, g := range m.GrossPremiumList {
		net := angka(g.Value)
		for _, b := range h.DeductionList {
			if b.Currency == g.Currency {
				net = kurang(net, angka(b.Deduction))
			}
		}
		h.NetPremiumList = append(h.NetPremiumList, NilaiMataUang{Currency: g.Currency, CurrencyID: g.CurrencyID, Value: teks(net)})
	}
	return h
}

// MasukanCadangan - medan yang `PremiumReserveCalculate` baca.
type MasukanCadangan struct {
	PremiumReservePct string          `json:"PremiumReservePct"`
	CessionList       []NilaiMataUang `json:"CessionList"`
}

// HasilCadangan - `ReserveList` baru beserta pesan.
type HasilCadangan struct {
	ReserveList []NilaiMataUang `json:"ReserveList"`
	Pesan       []string        `json:"pesan"`
}

// PesanCadanganLebih — `PremiumReserveCalculate` langkah 2, apa adanya.
const PesanCadanganLebih = "Tidak Boleh Lebih dari 100"

// HitungCadangan — bentuk ber-pelaku.
func (l *Layanan) HitungCadangan(p inti.Pelaku, m MasukanCadangan) (HasilCadangan, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return HasilCadangan{}, err
	}
	return HitungCadangan(m), nil
}

// HitungCadangan — `Activity/PremiumReserveCalculate.xml`.
func HitungCadangan(m MasukanCadangan) HasilCadangan {
	h := HasilCadangan{ReserveList: []NilaiMataUang{}, Pesan: []string{}}
	pct := angka(m.PremiumReservePct)
	lebih := pct.Cmp(seratus) > 0
	// [3]
	if lebih {
		h.Pesan = append(h.Pesan, PesanCadanganLebih)
	}
	// [5] satu baris per baris Cession; [5.2] dikosongkan bila > 100
	for _, c := range m.CessionList {
		if lebih {
			h.ReserveList = append(h.ReserveList, NilaiMataUang{})
			continue
		}
		v := new(apd.Decimal)
		_, _ = konteksLimit.Quo(v, kali(angka(c.Value), pct), seratus)
		h.ReserveList = append(h.ReserveList, NilaiMataUang{Currency: c.Currency, Value: teks(v)})
	}
	return h
}
