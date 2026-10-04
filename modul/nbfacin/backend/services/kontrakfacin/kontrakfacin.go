// Package kontrakfacin menyediakan mesin Fac In NB kepada modul lain lewat
// `inti/backend/kontrak` (facin.go) - keputusan work owner 01-10-2026 (butir 55).
// Isinya hanya penerjemah tipe: perilakunya tetap acceptance dan rules.
package kontrakfacin

import (
	"errors"

	"nusantarare/inti/backend/kontrak"
	"nusantarare/inti/backend/uang"
	"nusantarare/modul/nbfacin/backend/services/acceptance"
	"nusantarare/modul/nbfacin/backend/services/premium"
	"nusantarare/modul/nbfacin/backend/services/rules"
)

// Tangga memenuhi kontrak.TanggaAkseptasiFacIn atas tabel limit yang disuntikkan
// (kelak pemuat repository; kini fixture).
type Tangga struct {
	BentukA acceptance.TabelLimit
	BentukB []acceptance.BarisFinancial
}

var _ kontrak.TanggaAkseptasiFacIn = Tangga{}

// ErrTabelLimitKosong - Tangga dirakit tanpa tabel limit. Tanpa penolakan ini,
// tabel kosong = "tidak ada kandidat" = tangga selesai, dan kasus lolos tanpa
// approver karena datanya belum tersambung.
var ErrTabelLimitKosong = errors.New("kontrakfacin: tabel limit akseptasi belum tersambung")

// Langkah - acceptance.Next; kasus bentuk B (ErrBentukB) dialihkan ke
// acceptance.NextFinancial. Keduanya satu activity di Pega
// (GetLimitAkseptasi_ActFlow), jadi pemakai tidak perlu memilih bentuk.
func (t Tangga) Langkah(k kontrak.KasusFacIn, p kontrak.PenggunaFacIn) (kontrak.TransisiFacIn, error) {
	if len(t.BentukA) == 0 && len(t.BentukB) == 0 {
		return kontrak.TransisiFacIn{}, ErrTabelLimitKosong
	}
	pengguna := acceptance.Pengguna{Jabatan: acceptance.Jabatan(p.Jabatan), AnggotaGrup: p.AnggotaGrup}
	tr, err := acceptance.Next(k, t.BentukA, pengguna)
	if errors.Is(err, acceptance.ErrBentukB) {
		tr, err = acceptance.NextFinancial(k, t.BentukB, pengguna)
	}
	if err != nil {
		return kontrak.TransisiFacIn{}, err
	}
	return kontrak.TransisiFacIn{Selesai: tr.Selesai, JabatanTujuan: kontrak.JabatanFacIn(tr.JabatanTujuan),
		Antrean: kontrak.AntreanFacIn(tr.Antrean), PositionNoteDitulis: tr.PositionNoteDitulis}, nil
}

// Predikat memenuhi kontrak.PenilaiPredikatFacIn: registry `When` NB.
type Predikat struct{}

var _ kontrak.PenilaiPredikatFacIn = Predikat{}

// Eval - rules.Eval (galat untuk keraguan data, panic untuk kesalahan program).
func (Predikat) Eval(nama string, k kontrak.KasusFacIn) (bool, error) { return rules.Eval(nama, k) }

// Premi memenuhi kontrak.MesinPremiFacIn: paket premium NB.
type Premi struct{}

var _ kontrak.MesinPremiFacIn = Premi{}

func masukan(in kontrak.MasukanPremiFacIn) premium.Input {
	return premium.Input{LiniBisnis: premium.LiniBisnis(in.LiniBisnis), CalculateMethod: in.CalculateMethod,
		MataUang: in.MataUang, TSI: in.TSI, Rate: in.Rate, ProRatePercent: in.ProRatePercent,
		PctShortPeriod: in.PctShortPeriod, Discount: in.Discount, DiscountType: in.DiscountType,
		DiscountPercentage: in.DiscountPercentage, PremiSebelumnya: in.PremiSebelumnya, Loading: in.Loading,
		ProRatePercentCoverage: in.ProRatePercentCoverage, IndemnityPercentage: in.IndemnityPercentage,
		LossLimit: in.LossLimit, PctAdjustment: in.PctAdjustment, CoverageBasis: in.CoverageBasis,
		NetRate: in.NetRate, FirstScale: in.FirstScale, MBD: in.MBD, MasterPolicy: in.MasterPolicy}
}

// Hitung - premium.Calculate.
func (Premi) Hitung(in kontrak.MasukanPremiFacIn) (uang.Money, error) {
	return premium.Calculate(masukan(in))
}

// AsalRumus - premium.AsalRumus.
func (Premi) AsalRumus(in kontrak.MasukanPremiFacIn) string { return premium.AsalRumus(masukan(in)) }

// LiniDariPredikat - premium.LiniDariPredikat (panic bila tidak satu gerbang terbuka).
func (Premi) LiniDariPredikat(k kontrak.KasusFacIn) (kontrak.HasilLiniFacIn, error) {
	h, err := premium.LiniDariPredikat(k)
	if err != nil {
		return kontrak.HasilLiniFacIn{}, err
	}
	hasil := kontrak.HasilLiniFacIn{Peringatan: h.Peringatan}
	for _, l := range h.Lini {
		hasil.Lini = append(hasil.Lini, kontrak.LiniFacIn{Lini: string(l.Lini), PembagiRate: l.Satuan.Pembagi(), SimbolSatuan: l.Satuan.String()})
	}
	return hasil, nil
}
