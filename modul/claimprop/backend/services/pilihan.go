package services

// Untuk apa berkas ini: PILIHAN POP-UP DAN AUTOCOMPLETE - isi grid harness (MasterTreatyIn, ListPolicyNoTreaty_Harness,
// CauseofLoss_Harness, SummaryOutSClaim, CatastrofeList) dan sumber daftar sel (adjuster, provinsi, rekening, allocation,
// RNM Share, Limits, Spreading). Baca-saja; pilihan yang dipilih dikirim balik sebagai aksi dan DIBACA ULANG di server.

import (
	"context"
	"fmt"
	"strings"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/claimprop/backend/models"
)

// Jenis pilihan.
const (
	PilihanMaster     = "master"
	PilihanPolis      = "polis"
	PilihanSebab      = "sebab"
	PilihanKatastrofe = "katastrofe"
	PilihanRingkasan  = "ringkasanOS"
)

// BarisRingkasan - satu baris grid "Summary Outstanding Claim" (Section SummaryOutsClaim).
type BarisRingkasan struct {
	ClaimNo      string `json:"claimNo"`      // .CARI5
	PolicyNo     string `json:"policyNo"`     // .CARI6
	AcceptedNo   string `json:"acceptedNo"`   // .CARI7
	Currency     string `json:"currency"`     // .CARI1
	TreatyGross  string `json:"treatyGross"`  // .CARI12
	RNMShare     string `json:"rnmShare"`     // .CARI4
	KursValue    string `json:"kursValue"`    // .CARI9
	ConvertValue string `json:"convertValue"` // .CARI10
}

// Ringkasan - isi harness SummaryOutSClaim beserta kaki gridnya ("Total Outstanding": TotalData.CARI28 / CARI27).
type Ringkasan struct {
	Baris      []BarisRingkasan `json:"baris"`
	TotalLabel string           `json:"totalLabel"`
	TotalNilai string           `json:"totalNilai"`
}

// BerkasPolis - tombol View: berkas NB / EDM Treaty In terbaru untuk nomor polis (layar membukanya di tab baru).
func (l *Layanan) BerkasPolis(ctx context.Context, p inti.Pelaku, nopolis string) (models.BerkasPolis, error) {
	if err := l.periksaPelaku(p); err != nil {
		return models.BerkasPolis{}, err
	}
	nopolis = strings.TrimSpace(nopolis)
	if nopolis == "" {
		return models.BerkasPolis{}, fmt.Errorf("%w: nomor polis kosong", ErrPermintaanTidakSah)
	}
	b, ada, err := l.a.BerkasPolis(ctx, nopolis)
	if err != nil {
		return models.BerkasPolis{}, err
	}
	if !ada {
		return models.BerkasPolis{}, fmt.Errorf("%w: polis belum punya berkas NB / EDM Treaty In", ErrKasusTidakAda)
	}
	return b, nil
}

// Pilihan membaca satu daftar pilihan kasus.
func (l *Layanan) Pilihan(ctx context.Context, p inti.Pelaku, id, jenis string, n int, cari string,
	sm models.SaringanMaster) (any, error) {
	if err := l.periksaPelaku(p); err != nil {
		return nil, err
	}
	k, h, err := l.muat(ctx, nil, id)
	if err != nil {
		return nil, err
	}
	kt, err := l.konteks(ctx, p, l.jam())
	if err != nil {
		return nil, err
	}
	if err := l.turunkan(ctx, kt, h); err != nil {
		return nil, err
	}
	switch jenis {
	case PilihanMaster: // GetMasterTreaty_Act (defer load) + RD BrowseCLAIM_MASTER_TREATY
		return l.a.DaftarMaster(ctx, sm)
	case PilihanPolis: // SetMasterID + SetPolicyTreatyProp
		return l.a.DaftarPolis(ctx, models.AwalanMaster(h.Ambil(models.CD+"IDMaster")), h.Ambil(models.CD+"TreatyGroupName"))
	case PilihanSebab:
		return l.a.DaftarSebab(ctx, cari)
	case PilihanKatastrofe:
		return l.a.DaftarKatastrofe(ctx, cari)
	case PilihanRingkasan:
		return l.ringkasan(ctx, k, h)
	case models.SumberAdjuster:
		return l.a.DaftarAdjuster(ctx, cari)
	case models.SumberProvinsi:
		return l.a.DaftarProvinsi(ctx, cari)
	case models.SumberRekening:
		return models.PilihanRekening(kt, h, n)
	case models.SumberAllocation:
		var out []models.Pilihan
		for _, b := range h.AmbilDaftar(models.JalurAdj(n, models.AnakLossAllocation)) {
			out = append(out, models.Pilihan{Nilai: b["TreatyName"], Label: b["TreatyName"],
				Tambahan: map[string]string{"SharePercentage": b["SharePercentage"]}})
		}
		return out, nil
	case models.SumberMUAdj:
		var out []models.Pilihan
		for _, b := range h.AmbilDaftar(models.JalurAdj(n, "CurencyAdjustment")) {
			out = append(out, models.Pilihan{Nilai: b["CurrencyID"], Label: b["Currency"]})
		}
		return out, nil
	case models.SumberSpreading:
		// Spreading.pxResults (penulisnya tidak diekspor). KOREKSI 08-10-2026 atas `[dugaan]` SpreadingList master:
		// `[data DEV]` TreatyType SpreadingClaim klaim CLMP lama = INDUK (spreading polis), sedangkan SpreadingList
		// master = anaknya (isi SpreadingBreakQS). Opsi = spreading polis klaim (keputusan work owner 08-10-2026:
		// spreading diambil dari polis), berlabel nama treaty REINSURANCETYPE.
		sp, err := l.a.SpreadingPolis(ctx, h.Ambil(models.CD+"PolicyData.PolicyNo"))
		if err != nil {
			return nil, err
		}
		out := []models.Pilihan{}
		sudah := map[string]bool{}
		for _, s := range sp {
			if sudah[s.TreatyType] {
				continue
			}
			sudah[s.TreatyType] = true
			nama, err := l.a.NamaJenisReasuransi(ctx, s.TreatyType)
			if err != nil {
				return nil, err
			}
			if nama == "" {
				nama = s.TreatyType
			}
			out = append(out, models.Pilihan{Nilai: s.TreatyType, Label: nama})
		}
		return out, nil
	case models.SumberShareRNM, models.SumberLimits:
		m, _, err := l.a.MasterTreaty(ctx, h.Ambil(models.CD+"IDMaster"))
		if err != nil {
			return nil, err
		}
		var out []models.Pilihan
		switch jenis {
		case models.SumberShareRNM:
			for _, s := range models.PilihanShareRNM(models.HalamanBaru(), m) {
				out = append(out, models.Pilihan{Nilai: s, Label: s})
			}
		default: // SumberLimits
			for _, li := range m.Limits {
				out = append(out, models.Pilihan{Nilai: li.TreatyType, Label: li.TreatyType})
			}
		}
		return out, nil
	}
	return nil, fmt.Errorf("%w: pilihan %q", ErrPermintaanTidakSah, jenis)
}

// ringkasan = GetDataOustanding (Outstanding Claim) / GetDataOutsClaim_act (Input Acceptation).
func (l *Layanan) ringkasan(ctx context.Context, k models.Kasus, h *models.Halaman) (Ringkasan, error) {
	rows, err := l.a.RingkasanOS(ctx, h.Ambil(models.CD+"NoClaim"))
	if err != nil {
		return Ringkasan{}, err
	}
	var kal models.Kalkulator
	out := Ringkasan{Baris: []BarisRingkasan{}}
	if k.Tahap == models.TahapOutstanding {
		// GetDataOustanding 3.1: CARI10 = kurs x nilai; CARI27 = toDecimal(CARI28 lama) + gross; CARI28 = mata uang
		// (teks - toDecimal teks mata uang = 0, sehingga totalnya gross baris TERAKHIR; ditiru apa adanya).
		for _, r := range rows {
			conv := kal.Kali(kal.Teks("KursValue", r.KursValue), kal.Teks("Value", r.Value))
			out.Baris = append(out.Baris, BarisRingkasan{ClaimNo: r.NoClaim, PolicyNo: r.NoPolis, AcceptedNo: r.AcceptedNo,
				Currency: r.Currency, TreatyGross: r.GrossValue, RNMShare: r.Value, KursValue: r.KursValue,
				ConvertValue: models.Teks(conv)})
			prev := apd.New(0, 0)
			if models.AdalahDesimal(out.TotalLabel) {
				prev = kal.Teks("CARI28", out.TotalLabel)
			}
			out.TotalNilai = models.Teks(kal.Tambah(prev, kal.Teks("GrossValue", r.GrossValue)))
			out.TotalLabel = r.Currency
		}
	} else {
		// GetDataOutsClaim_act 3.1-3.4: nilai baris ber-Type selain "0" dinegatifkan; tanpa nomor akseptasi gross =
		// GrossValue; konversi = kurs x nilai; total = jumlah konversi, label "IDR".
		tot := apd.New(0, 0)
		for _, r := range rows {
			v := kal.Teks("Value", r.Value)
			if r.Type != "0" {
				v = kal.Neg(v)
			}
			gross := r.TotalGross
			if r.AcceptedNo == "" {
				gross = r.GrossValue
			}
			conv := kal.Kali(kal.Teks("KursValue", r.KursValue), v)
			tot = kal.Tambah(tot, conv)
			out.Baris = append(out.Baris, BarisRingkasan{ClaimNo: r.NoClaim, PolicyNo: r.NoPolis, AcceptedNo: r.AcceptedNo,
				Currency: r.Currency, TreatyGross: gross, RNMShare: models.Teks(v), KursValue: r.KursValue,
				ConvertValue: models.Teks(conv)})
		}
		out.TotalLabel, out.TotalNilai = "IDR", models.Teks(tot)
	}
	return out, kal.Galat()
}
