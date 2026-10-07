package models

// Spreading peserta ke retro - Calculate CSV Type QR (keputusan work owner 06-10-2026).
//
// Untuk apa berkas ini: aturan MURNI (tanpa Oracle) yang membagi RETROCEDED_SHARE
// setiap peserta ke lapisan Contract Retro Life (kaskade Maximum Limit), lalu ke
// reinsurer tiap lapisan dengan rate retro, discount, dan overriding commission.
//
// Acuan: `SetRateLIfePremium_Act` langkah 10.8-10.12, DITULIS ULANG menurut aturan
// work owner 06-10-2026 (aturan ini yang berlaku bila berbeda dari XML):
//
//	S1 treaty year = TREATYYEAR_LIFE yang periodenya memuat BEGIN_DATE produk
//	   (bukan Outward List / ORS; Underwriting/Transaction Year tidak dipakai)
//	S2-S5 lapisan = seluruh kontrak treaty year itu, kaskade Maximum Limit menaik
//	   menurut mata uang peserta (IDR / USD); melebihi lapisan terakhir DITOLAK
//	S6 lapisan terpakai: tepat satu Business = Class of Business polis, total
//	   PCTSHARE reinsurer = 100, master rate seragam
//	S7 rate retro = pencarian yang SAMA dengan rate premi (MasterRate.Cari)
//	S8 per reinsurer: Amount, GROSS (TANPA FACTOR), Discount/Comm/NET; berhenti di
//	   reinsurer (TREATYSECURITYREINSURER_LIFE tidak dipakai)
//
// ⛔ NOL FLOAT (ADR-U-0003): apd presisi penuh; dibulatkan half-up 4 desimal hanya
// saat disimpan (repository).

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/utils"
)

// ErrSyaratSpreading - data Contract Retro Life / produk tidak memenuhi syarat
// spreading (S1, S3 limit kosong, S6). Seluruh unggahan ditolak (409), nol tulisan.
var ErrSyaratSpreading = errors.New("models: syarat spreading retro tidak terpenuhi")

// HasilBarisQR - bahan hitung Type QR satu peserta yang dipakai spreading:
// nilai presisi penuh (sebelum dibulatkan) dan kunci pencarian rate.
type HasilBarisQR struct {
	Retro        *apd.Decimal // RETROCEDED_SHARE = SUM_AT_RISK_GROSS, presisi penuh
	Umur         int64        // ENTRY_AGE
	Kontrak      *apd.Decimal // CONTRACT dari PERIOD_MM
	Sex          string       // F / M / lainnya (huruf besar)
	EM           *apd.Decimal // EM_PERCENT pecahan; kosong = 0
	TahunPertama bool         // BEGIN_DATE = GROSS_VALUATION_BEGIN_DATE
}

// Bahan Contract Retro Life - teks apa adanya dari Oracle (angka teks desimal,
// tanggal YYYY-MM-DD).
type (
	// TahunRetro - satu baris TREATYYEAR_LIFE.
	TahunRetro struct{ ID, Mulai, Akhir string }
	// KontrakRetro - satu baris TREATYCONTRACT_LIFE (lapisan).
	KontrakRetro struct {
		ID, IDTahun, ReinsTypeID, ReinsTypeName, Mulai, Akhir string
		IDR, USD, BIDR, BUSD, IDRSelisih, USDSelisih          string
	}
	// ReinsurerRetro - satu baris TREATYREINSURER_LIFE (persen 0-100).
	ReinsurerRetro struct{ KontrakID, Nama, PctShare, Komisi, OvrComm string }
	// BusinessRetro - satu baris TREATYBUSINESS_LIFE.
	BusinessRetro struct{ KontrakID, BizName, RIRateID, RIRate string }
)

// ReinsurerSpreading - hasil satu reinsurer di satu lapisan (presisi penuh).
type ReinsurerSpreading struct {
	Nama                  string
	PctShare, Komisi, Ovr *apd.Decimal // persen
	Amount, Rate          *apd.Decimal
	Gross, Discount, Comm *apd.Decimal
	Net                   *apd.Decimal
}

// LapisanSpreading - bagian satu peserta di satu lapisan (bagian > 0 saja).
type LapisanSpreading struct {
	TahunID   string
	Kontrak   KontrakRetro
	Bagian    *apd.Decimal
	Reinsurer []ReinsurerSpreading
}

// tanggalRetro - YYYY-MM-DD (atau YYYY-MM-DD HH:MM:SS) → tanggal; kosong/rusak → zero.
func tanggalRetro(s string) time.Time {
	s = strings.TrimSpace(s)
	if len(s) >= 10 {
		s = s[:10]
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}
	}
	return t
}

// PilihTahunRetro - S1: tahun treaty yang STARTDATE ≤ beginProduk ≤ ENDDATE.
// Nol atau lebih dari satu → ErrSyaratSpreading berpesan.
func PilihTahunRetro(tahun []TahunRetro, beginProduk time.Time) (TahunRetro, error) {
	b := time.Date(beginProduk.Year(), beginProduk.Month(), beginProduk.Day(), 0, 0, 0, 0, time.UTC)
	var cocok []TahunRetro
	for _, t := range tahun {
		m, a := tanggalRetro(t.Mulai), tanggalRetro(t.Akhir)
		if m.IsZero() || a.IsZero() {
			continue
		}
		if !b.Before(m) && !b.After(a) {
			cocok = append(cocok, t)
		}
	}
	tgl := b.Format("02/01/2006")
	switch len(cocok) {
	case 0:
		return TahunRetro{}, fmt.Errorf("%w: No Contract Retro Life covers begin date %s", ErrSyaratSpreading, tgl)
	case 1:
		return cocok[0], nil
	}
	id := make([]string, len(cocok))
	for i, t := range cocok {
		id[i] = t.ID
	}
	return TahunRetro{}, fmt.Errorf("%w: More than one Contract Retro Life covers begin date %s (%s) - their periods overlap",
		ErrSyaratSpreading, tgl, strings.Join(id, ", "))
}

// KaskadeRetro - S5: bagian lapisan k = min(nilai, Max_k) − min(nilai, Max_{k-1}),
// Max_0 = 0; \`maks\` urut menaik. lebih = nilai melebihi Max lapisan terakhir.
func KaskadeRetro(nilai *apd.Decimal, maks []*apd.Decimal) (bagian []*apd.Decimal, lebih bool) {
	sebelum := apd.New(0, 0)
	min := func(a, b *apd.Decimal) *apd.Decimal {
		if a.Cmp(b) < 0 {
			return a
		}
		return b
	}
	for _, m := range maks {
		var d apd.Decimal
		_, _ = ctxQR.Sub(&d, min(nilai, m), min(nilai, sebelum))
		if d.Sign() < 0 {
			d.SetInt64(0)
		}
		bagian = append(bagian, &d)
		if m.Cmp(sebelum) > 0 {
			sebelum = m
		}
	}
	return bagian, len(maks) == 0 || nilai.Cmp(sebelum) > 0
}

// reinsurerLapisan - reinsurer satu lapisan, persen sudah diurai.
type reinsurerLapisan struct {
	nama             string
	pct, komisi, ovr *apd.Decimal
}

// lapisanRetro - satu lapisan siap pakai; galat S6 ditunda sampai lapisan terpakai.
type lapisanRetro struct {
	k         KontrakRetro
	maks      map[string]*apd.Decimal // "IDR" / "USD"; nil = kosong
	reinsurer []reinsurerLapisan
	rate      MasterRate
	namaRate  string
	salah     string // S6 - kalimat galat bila lapisan ini terpakai
}

// MasterSpreading - bahan spreading satu polis, siap dipakai SpreadingQR.
type MasterSpreading struct {
	tahunID string
	proRate string
	lapisan []lapisanRetro
}

// SiapkanSpreading mengurai bahan satu treaty year. Angka yang tidak terurai →
// ErrSyaratSpreading; syarat S6 dicatat per lapisan dan baru menolak bila lapisan
// itu terpakai. \`rates\` = baris RATE_LIFE per RIRATEID.
func SiapkanSpreading(tahun TahunRetro, kontrak []KontrakRetro, reinsurer []ReinsurerRetro, bisnis []BusinessRetro,
	rates map[string][]BarisRateQR, bisnisPolis, proRate string) (MasterSpreading, error) {

	gagal := func(format string, a ...any) (MasterSpreading, error) {
		return MasterSpreading{}, fmt.Errorf("%w: %s", ErrSyaratSpreading, fmt.Sprintf(format, a...))
	}
	m := MasterSpreading{tahunID: tahun.ID, proRate: strings.TrimSpace(proRate)}
	namaBisnis := strings.ToUpper(strings.TrimSpace(bisnisPolis))
	for _, k := range kontrak {
		l := lapisanRetro{k: k, maks: map[string]*apd.Decimal{}}
		for mu, teks := range map[string]string{"IDR": k.IDR, "USD": k.USD} {
			d, err := angkaMaster(teks)
			if err != nil {
				return gagal("Maximum Limit (%s) %q of %s (Contract Retro Life %s) is not a number", mu, teks, k.ReinsTypeName, tahun.ID)
			}
			l.maks[mu] = d
		}
		// Reinsurer + total share.
		total := apd.New(0, 0)
		for _, r := range reinsurer {
			if r.KontrakID != k.ID {
				continue
			}
			rl := reinsurerLapisan{nama: strings.TrimSpace(r.Nama)}
			for _, x := range []struct {
				teks, label string
				ke          **apd.Decimal
			}{{r.PctShare, "PCT SHARE", &rl.pct}, {r.Komisi, "COMMISION", &rl.komisi}, {r.OvrComm, "OVR COMM", &rl.ovr}} {
				d, err := angkaMaster(x.teks)
				if err != nil {
					return gagal("Reinsurer %s of %s (Contract Retro Life %s): %s %q is not a number", rl.nama, k.ReinsTypeName, tahun.ID, x.label, x.teks)
				}
				if d == nil {
					d = apd.New(0, 0)
				}
				*x.ke = d
			}
			_, _ = ctxQR.Add(total, total, rl.pct)
			l.reinsurer = append(l.reinsurer, rl)
		}
		// S6 - business tepat satu.
		var cocok []BusinessRetro
		for _, b := range bisnis {
			if b.KontrakID == k.ID && namaBisnis != "" && strings.ToUpper(strings.TrimSpace(b.BizName)) == namaBisnis {
				cocok = append(cocok, b)
			}
		}
		switch {
		case len(cocok) != 1:
			l.salah = fmt.Sprintf("Business %s is not in the Business List of %s (Contract Retro Life %s)", strings.TrimSpace(bisnisPolis), k.ReinsTypeName, tahun.ID)
			if len(cocok) > 1 {
				l.salah = fmt.Sprintf("Business %s appears more than once in the Business List of %s (Contract Retro Life %s)", strings.TrimSpace(bisnisPolis), k.ReinsTypeName, tahun.ID)
			}
		case len(l.reinsurer) == 0 || total.Cmp(seratus) != 0:
			l.salah = fmt.Sprintf("Reinsurer shares of %s (Contract Retro Life %s) total %s%%, must be 100%%", k.ReinsTypeName, tahun.ID, utils.FormatDecimal(total))
		default:
			l.namaRate = strings.TrimSpace(cocok[0].RIRate)
			id := strings.TrimSpace(cocok[0].RIRateID)
			mr, salah := UraiMasterRate(id, rates[id])
			switch {
			case id == "":
				l.salah = fmt.Sprintf("Business %s of %s (Contract Retro Life %s) has no R/I Rate", strings.TrimSpace(bisnisPolis), k.ReinsTypeName, tahun.ID)
			case salah != "":
				l.salah = fmt.Sprintf("%s (%s, Contract Retro Life %s)", salah, k.ReinsTypeName, tahun.ID)
			case mr.gender && !ProRateSah(m.proRate):
				l.salah = fmt.Sprintf("Premium Payment Method (1, 2 or 3) must be set on the policy for gendered R/I Rate %s (%s, Contract Retro Life %s)", l.namaRate, k.ReinsTypeName, tahun.ID)
			default:
				l.rate = mr
			}
		}
		m.lapisan = append(m.lapisan, l)
	}
	if len(m.lapisan) == 0 {
		return gagal("Contract Retro Life %s has no reins type", tahun.ID)
	}
	return m, nil
}

// SpreadingQR - S2-S8 atas setiap baris yang terhitung QR (\`b.QR\` terisi, tidak
// di \`lewati\`), menulis \`b.Spreading\`. Penolakan baris dikembalikan; syarat
// lapisan terpakai yang gagal (S3 limit kosong, S6) → ErrSyaratSpreading.
func SpreadingQR(m MasterSpreading, baris []BarisUnggah, lewati map[int]bool) ([]Penolakan, error) {
	tolak := []Penolakan{}
	for i := range baris {
		b := &baris[i]
		if lewati[b.Nomor] || b.QR == nil || b.QR.Retro == nil {
			continue
		}
		p, err := spreadingSatu(m, b)
		if err != nil {
			return nil, err
		}
		if p != nil {
			tolak = append(tolak, *p)
		}
	}
	return tolak, nil
}

func spreadingSatu(m MasterSpreading, b *BarisUnggah) (*Penolakan, error) {
	tolak := func(kolom, pesan string) (*Penolakan, error) {
		return &Penolakan{Baris: b.Nomor, Kolom: kolom, Pesan: pesan, Sebab: pesan}, nil
	}
	mu := strings.ToUpper(strings.TrimSpace(b.Nilai["CURRENCY"]))
	if mu != "IDR" && mu != "USD" {
		return tolak("CURRENCY", fmt.Sprintf("CURRENCY %q is not supported for retro spreading (IDR or USD)", b.Nilai["CURRENCY"]))
	}
	// S2/S3 - lapisan urut Maximum Limit mata uang ini menaik.
	urut := make([]lapisanRetro, len(m.lapisan))
	copy(urut, m.lapisan)
	for _, l := range urut {
		if l.maks[mu] == nil {
			return nil, fmt.Errorf("%w: Maximum Limit (%s) is empty on %s (Contract Retro Life %s)", ErrSyaratSpreading, mu, l.k.ReinsTypeName, m.tahunID)
		}
	}
	sort.SliceStable(urut, func(i, j int) bool { return urut[i].maks[mu].Cmp(urut[j].maks[mu]) < 0 })
	maks := make([]*apd.Decimal, len(urut))
	for i, l := range urut {
		maks[i] = l.maks[mu]
	}
	bagian, lebih := KaskadeRetro(b.QR.Retro, maks)
	if lebih {
		return tolak("RETROCEDED_SHARE", fmt.Sprintf("Retro value %s exceeds the maximum limit %s of Contract Retro Life %s (row %d)",
			AngkaTampil(bulatDesimal(b.QR.Retro, 4)), AngkaTampil(maks[len(maks)-1]), m.tahunID, b.Nomor))
	}
	var hasil []LapisanSpreading
	for i, l := range urut {
		if bagian[i].Sign() == 0 {
			continue // lapisan berbagian 0 tidak diproses dan tidak disimpan
		}
		if l.salah != "" {
			return nil, fmt.Errorf("%w: %s", ErrSyaratSpreading, l.salah)
		}
		// S7 - rate retro, pencarian yang sama dengan rate premi.
		if l.rate.gender && b.QR.Sex != "F" && b.QR.Sex != "M" {
			return tolak("SEX", fmt.Sprintf("SEX must be F or M (retro R/I Rate %s of %s is gendered), found %q", l.namaRate, l.k.ReinsTypeName, b.QR.Sex))
		}
		rate, cocok, dicari := l.rate.Cari(m.proRate, b.QR.Umur, b.QR.Kontrak, b.QR.Sex)
		if cocok != 1 {
			awal := "No retro rate found"
			if cocok > 1 {
				awal = "More than one retro rate found"
			}
			return tolak("RATE", fmt.Sprintf("%s for %s in R/I Rate %s (%s, Contract Retro Life %s)", awal, dicari, l.namaRate, l.k.ReinsTypeName, m.tahunID))
		}
		lap := LapisanSpreading{TahunID: m.tahunID, Kontrak: l.k, Bagian: bagian[i]}
		for _, r := range l.reinsurer {
			lap.Reinsurer = append(lap.Reinsurer, HitungReinsurerRetro(r.nama, r.pct, r.komisi, r.ovr, bagian[i], rate, b.QR.EM, b.QR.TahunPertama))
		}
		hasil = append(hasil, lap)
	}
	b.Spreading = hasil
	return nil, nil
}

// HitungReinsurerRetro - S8 satu reinsurer:
//
//	Amount = bagian × PCTSHARE / 100
//	GROSS  = Rate / 1000 × (1 + EM) × Amount            (TANPA FACTOR)
//	tahun pertama: Discount = GROSS × COMMISION / 100; Comm = (GROSS − Discount) × OVR_COMM / 100
//	bukan:         Discount = 0;                        Comm = GROSS × OVR_COMM / 100
//	NET = GROSS − Discount − Comm
func HitungReinsurerRetro(nama string, pct, komisi, ovr, bagian, rate, em *apd.Decimal, tahunPertama bool) ReinsurerSpreading {
	if em == nil {
		em = apd.New(0, 0)
	}
	var amount, gross, emSatu, disc, dasar, comm, net apd.Decimal
	_, _ = ctxQR.Mul(&amount, bagian, pct)
	_, _ = ctxQR.Quo(&amount, &amount, seratus)
	_, _ = ctxQR.Quo(&gross, rate, seribu)
	_, _ = ctxQR.Add(&emSatu, em, satu)
	_, _ = ctxQR.Mul(&gross, &gross, &emSatu)
	_, _ = ctxQR.Mul(&gross, &gross, &amount)
	if tahunPertama {
		_, _ = ctxQR.Mul(&disc, &gross, komisi)
		_, _ = ctxQR.Quo(&disc, &disc, seratus)
	}
	_, _ = ctxQR.Sub(&dasar, &gross, &disc)
	if !tahunPertama {
		dasar.Set(&gross)
	}
	_, _ = ctxQR.Mul(&comm, &dasar, ovr)
	_, _ = ctxQR.Quo(&comm, &comm, seratus)
	_, _ = ctxQR.Sub(&net, &gross, &disc)
	_, _ = ctxQR.Sub(&net, &net, &comm)
	return ReinsurerSpreading{Nama: nama, PctShare: pct, Komisi: komisi, Ovr: ovr, Amount: &amount, Rate: rate,
		Gross: &gross, Discount: &disc, Comm: &comm, Net: &net}
}

// bulatDesimal - half-up \`digit\` desimal (salinan, untuk tampilan pesan).
func bulatDesimal(d *apd.Decimal, digit int32) *apd.Decimal {
	var r apd.Decimal
	_, _ = ctxQR.Quantize(&r, d, -digit)
	return &r
}

// TeksSimpan - half-up 4 desimal untuk ditulis ke Oracle; nil → "".
func TeksSimpan(d *apd.Decimal) string {
	if d == nil {
		return ""
	}
	return bulatTeks(d, 4)
}
