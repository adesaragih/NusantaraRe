package services

// Tab Inw Fac Cedant Panels kasus FIRE (tiket 49). Port `[terverifikasi]` `D:\migrasi\RNM\NB FacIn\` (dibaca 05-10-2026):
//   - Activity\SetShareOfCeding.xml: ShareCedantType==0 -> QuotationData.ShareOfCeding = "100%"; ==1 -> PercentShare + "%"
//     (salinan ke pyWorkPage.Quotation.ShareOfCeding tidak dibawa - halaman Quotation work tidak disimpan modul ini, K49-4);
//   - Activity\ProtectShareCedant_Act.xml (dari Protection_Act): gerbang langkah 1 + pesan verbatim (K49-2);
//   - Activity\SetTSIPremiCedant_Act.xml: ShareCeding > 100 -> "Value can't be more than 100 or less than 0!" (tanpa
//     gerbang); CurrencyList per cedant = langkah 4 atas `pyWorkPage.OfferFacIn.CurrencyList` (langkah 5 atas
//     TotalSpreadAll ber-label `//` = mati), dasarnya per mata uang dikirim ke layar (K49-8);
//   - Activity\CountPremiAndTSIRNMFireMBU_ACT.xml langkah 6.2.2.1.2 / 6.2.2.1.3.3 / 6.2.3: TotalTSINusaReSpreading =
//     PercentShare × Σ CoverageList(1).TSILiability tiap item / 100, TotalPremiNusaRe = Σ PremiNusantaraRe (K49-3) -
//     MENYIMPANG (C-8, K49-11): dijumlah per mata uang lalu dikonversi ke IDR (Pega menjumlah angka lintas mata uang).

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

var (
	// ErrMasukanCedant - isian tab Cedant tidak sah / proteksi Save. 400.
	ErrMasukanCedant = errors.New("services: isian tab Cedant tidak sah")
	// ErrCedantTanpaDatabase - tabel tab Cedant tidak terbaca. 503.
	ErrCedantTanpaDatabase = errors.New("services: basis data tidak dikonfigurasi, tab Cedant tidak terbaca")
)

const (
	// Pesan verbatim Pega.
	pesanShareCedant  = "% Share cedant more than 100 or less than 0!"            // ProtectShareCedant_Act 3.1
	pesanCedingKosong = "Ceding on list inward facultative cedant can't be null!" // ProtectShareCedant_Act 3.2
	pesanCedantKosong = "List Inward facultative cedant can't be null!"           // ProtectShareCedant_Act 6
	pesanShareCeding  = "Value can't be more than 100 or less than 0!"            // SetTSIPremiCedant_Act
	pesanTotalPremi   = "Total premi share cedant must be equal total premi RNM!" // ProtectShareCedant_Act 5
	// pesanTipeCedantWajib - pesan agent: Share Cedant Type wajib (Section `.FlagOnGoingPolicy ==2 || IsGroup`, K49-2).
	pesanTipeCedantWajib = "Share Cedant Type wajib diisi"
	// pesanTotalShareCedant - pesan agent C-7 (work owner 05-10-2026 "ini masa bisa lebih dari 100%"): Σ ShareCeding > 100.
	pesanTotalShareCedant = "Total % Share cedant tidak boleh lebih dari 100."
	tipeCedantGross       = "0" // DDL\ShareCedantType.xml
	tipeCedantShareRNM    = "1"
	banyakCedant          = 100
	lebarKodeCeding       = 50 // T_CEDINGCEDANTLIST.CEDING_CO
)

// DenganCedant memasang penyimpan tab Cedant (tiket 49).
func (s *Service) DenganCedant(p repository.PenyimpanCedant) *Service {
	s.cedant = p
	return s
}

// TampilanCedant - jawaban GET / PUT tab Cedant. TotalTSIRNM / TotalPremiRNM dalam IDR (C-8); kosong bila % Share RNM
// kosong atau KursTidakAda berisi.
type TampilanCedant struct {
	SobName, ShareCedantType, ShareOfCeding  string
	PercentShare, TotalTSIRNM, TotalPremiRNM *apd.Decimal
	Wajib                                    bool
	Cedant, CedingUmum                       []models.BarisCedant
	DasarMataUang                            []DasarMataUangCedant
	KursTidakAda                             []string
}

// mataUangKosong - penanda KursTidakAda untuk item tanpa mata uang (tidak dapat dikonversi).
const mataUangKosong = "(kosong)"

// jumlahMataUang - Total TSI / Premi RNM satu mata uang asli (sebelum konversi IDR).
type jumlahMataUang struct {
	Currency      string
	TSIRNM, Premi *apd.Decimal
}

// DasarMataUangCedant - satu `OfferFacIn.CurrencyList(n)` sebagai dasar `CedingCedantList(i).CurrencyList` pada ShareCeding
// = 100 (SetTSIPremiCedant_Act langkah 4): TSI = Σ TSIObjectItem item bermata uang itu, Premi = Σ PremiNusantaraRe
// coverage-nya (CountPremiAndTSIRNMFireMBU_ACT 6.2.2.1.2 / 6.2.2.1.3.3 -> 6.2.5 `.TSI` / `.Policy.Payment.Premium`).
//   - Share RNM: TSIRNM = TSI × PercentShare / 100, PremiumRNM = Premi;
//   - Gross: TSIGross = TSI, PremiumGross = Premi × 100 / PercentShare (nil bila PercentShare 0).
type DasarMataUangCedant struct {
	Currency                                   string
	TSIRNM, PremiumRNM, TSIGross, PremiumGross *apd.Decimal
}

// gerbangCedant - masukan langkah 1 ProtectShareCedant_Act.
type gerbangCedant struct {
	Life, EmailBinding4, Group, OnGoing, EdmAdjShareCedant bool
}

// proteksiCedantBerlaku - langkah 1 ProtectShareCedant_Act, urutan When verbatim: IsLife -> 6, EmailTypeBinding "4" -> 6,
// IsGroup "Group" -> 5, FlagOnGoingPolicy 2 -> 5, IsEdmAdjShareCedant true -> 2 / false -> 6. `[dugaan]` 6 = keluar
// activity, 5 = lewati When sisanya dan jalankan, 2 = lanjut ke When berikut (arti kode belum terverifikasi).
func proteksiCedantBerlaku(g gerbangCedant) bool {
	switch {
	case g.Life, g.EmailBinding4:
		return false
	case g.Group, g.OnGoing:
		return true
	default:
		return g.EdmAdjShareCedant
	}
}

// gerbangKasus - sumber gerbang di sistem baru (K49-2): FlagOnGoingPolicy = T_WORK_POLIS.FLAG_ONGOING_POLICY (milik
// premiumlistlife, migrasi 057; case NB ditulis "0"), `pyWorkPage.FlagOnGoingPolicy==2`. Kasus FIRE bukan Life;
// EmailTypeBinding dan IsEdmAdjShareCedant (EDM) tidak berlaku untuk NB; IsGroup (login pembuat case, A14 / butir 33)
// BELUM punya sumber -> false.
func gerbangKasus(k models.KasusCedant) gerbangCedant {
	return gerbangCedant{OnGoing: strings.TrimSpace(k.FlagOnGoingPolicy) == "2"}
}

// periksaProteksiCedant - ProtectShareCedant_Act langkah 2-10. Urutan pesan = urutan Page-Set-Messages Pega (langkah 7
// total premi, 8 daftar kosong, 9 share, 10 ceding). Total premi: @Math.divide(TotalPremiNusaRe, 1, 4) ==
// @Math.divide(Σ CedingCedantList(i).CurrencyList(n).Premium, 1, 4) (langkah 3.3, 4, 5).
func periksaProteksiCedant(baris []models.BarisCedant, tipe string, dasar []DasarMataUangCedant, totalPremiRNM *apd.Decimal) []string {
	var m []string
	satu := apd.New(1, 0)
	if bagiSkala(totalPremiRNM, satu, 4).Cmp(bagiSkala(premiCedant(baris, tipe, dasar), satu, 4)) != 0 {
		m = append(m, pesanTotalPremi)
	}
	if len(baris) == 0 {
		m = append(m, pesanCedantKosong)
	}
	for _, b := range baris {
		if b.ShareCeding == nil || b.ShareCeding.Sign() <= 0 || lebihBesar(b.ShareCeding, seratus) {
			m = append(m, pesanShareCedant)
			break
		}
	}
	for _, b := range baris { // Pega `.CedingCoName==""`; kode kosong = belum Choose Ceding (nama ditulis dari AGENT)
		if strings.TrimSpace(b.CedingCo) == "" || b.CedingCoName == "" {
			m = append(m, pesanCedingKosong)
			break
		}
	}
	return m
}

// shareOfCeding - SetShareOfCeding: nil bila Share Cedant Type kosong (tidak ada langkah yang jalan).
func shareOfCeding(tipe string, ps *apd.Decimal) *string {
	var v string
	switch tipe {
	case tipeCedantGross:
		v = "100%"
	case tipeCedantShareRNM:
		v = "%" // PercentShare kosong -> "" + "%" (verbatim)
		if ps != nil {
			v = utils.FormatDecimal(simpan8(ps)) + "%"
		}
	default:
		return nil
	}
	return &v
}

// ringkasanRNM - dari pohon objek (NR dihitung ulang dari % Share RNM tersimpan, seperti tiap Save Pega):
//   - premiMentah = TotalPremiNusaRe Pega (Σ PremiNusantaraRe lintas mata uang TANPA kurs) - pembanding proteksi K49-9;
//   - jumlah = per mata uang asli: TSIRNM = PercentShare × Σ CoverageList(1).TSILiability item / 100, Premi = Σ premi NR
//     (item tanpa mata uang = mataUangKosong) - dasar Total RNM IDR (C-8);
//   - dasar = CurrencyList cedant pada ShareCeding 100 (K49-8; item tanpa mata uang tidak masuk).
//
// Kosong bila % Share RNM kosong. Skala hitung 20, dibulatkan 8 desimal saat tampil. Mata uang urut kemunculan pertama
// (Pega: OfferFacIn.CurrencyList, `[dugaan]` sama).
func ringkasanRNM(k models.KasusSpreading) (premiMentah *apd.Decimal, jumlah []jumlahMataUang, dasar []DasarMataUangCedant) {
	dasar, jumlah = []DasarMataUangCedant{}, []jumlahMataUang{}
	if k.PercentShare == nil {
		return nil, jumlah, dasar
	}
	hitungNR(&k)
	premiMentah = nol()
	var urut []string
	tsiObjek, tsiLiab, premiUang := map[string]*apd.Decimal{}, map[string]*apd.Decimal{}, map[string]*apd.Decimal{}
	for _, l := range k.Lokasi {
		for _, it := range l.Items {
			u := it.Currency
			if u == "" {
				u = mataUangKosong
			}
			if _, ada := premiUang[u]; !ada {
				urut = append(urut, u)
				tsiObjek[u], tsiLiab[u], premiUang[u] = nol(), nol(), nol()
			}
			for _, c := range it.Coverages {
				premiUang[u] = tambah(premiUang[u], c.PremiNusantaraRe)
				premiMentah = tambah(premiMentah, c.PremiNusantaraRe)
			}
			if len(it.Coverages) > 0 {
				tsiLiab[u] = tambah(tsiLiab[u], it.Coverages[0].TSILiability)
			}
			tsiObjek[u] = tambah(tsiObjek[u], it.TSI)
		}
	}
	for _, u := range urut {
		jumlah = append(jumlah, jumlahMataUang{Currency: u, TSIRNM: persenDari(k.PercentShare, tsiLiab[u]), Premi: premiUang[u]})
		if u == mataUangKosong {
			continue
		}
		d := DasarMataUangCedant{Currency: u, TSIRNM: persenDari(tsiObjek[u], k.PercentShare), PremiumRNM: premiUang[u],
			TSIGross: tsiObjek[u]}
		if !k.PercentShare.IsZero() {
			d.PremiumGross = bagiSkala(kali(premiUang[u], seratus), k.PercentShare, 20)
		}
		dasar = append(dasar, d)
	}
	return premiMentah, jumlah, dasar
}

// totalKeIDR - Total TSI / Premi RNM IDR (C-8, K49-11) = Σ per mata uang × kurs, dengan pembaca kurs ringkasan Spreading
// (SumTSIPremiSpreadedRNM_Act langkah 19): ID mata uang dari CURRENCY; ID "10026" (IDR) -> kurs 1; selain itu
// GetKursLimitSpreading_SQL pada Begin date (KursPada) - tanpa jatuh ke kurs terbaru, sama dengan Spreading. Satu saja mata
// uang tanpa kurs -> kedua total kosong dan mata uangnya dilaporkan (tidak menjumlah sebagian).
func totalKeIDR(ctx context.Context, src sumberKurs, mulai string, jumlah []jumlahMataUang) (tsi, premi *apd.Decimal, tidakAda []string, err error) {
	tidakAda = []string{}
	tsi, premi = nol(), nol()
	for _, j := range jumlah {
		var kurs *apd.Decimal
		if j.Currency != mataUangKosong {
			id, err := src.IDMataUang(ctx, j.Currency)
			if err != nil {
				return nil, nil, nil, err
			}
			switch {
			case id == idMataUangIDR:
				kurs = apd.New(1, 0)
			case id != "":
				if kurs, err = src.KursPada(ctx, id, mulai); err != nil {
					return nil, nil, nil, err
				}
			}
		}
		if kurs == nil {
			tidakAda = append(tidakAda, j.Currency)
			continue
		}
		tsi, premi = tambah(tsi, kali(nolBila(j.TSIRNM), kurs)), tambah(premi, kali(nolBila(j.Premi), kurs))
	}
	if len(tidakAda) > 0 || len(jumlah) == 0 {
		return nil, nil, tidakAda, nil
	}
	return tsi, premi, tidakAda, nil
}

// premiCedant - Σ CedingCedantList(i).CurrencyList(n).Premium (SetTSIPremiCedant_Act langkah 4): Share RNM = ShareCeding ×
// PremiumRNM / 100; Gross = ShareCeding × PremiumGross / 100; tipe kosong = tidak ada langkah yang jalan (0).
func premiCedant(baris []models.BarisCedant, tipe string, dasar []DasarMataUangCedant) *apd.Decimal {
	total := nol()
	for _, b := range baris {
		for _, d := range dasar {
			switch tipe {
			case tipeCedantShareRNM:
				total = tambah(total, persenDari(b.ShareCeding, d.PremiumRNM))
			case tipeCedantGross:
				total = tambah(total, persenDari(b.ShareCeding, d.PremiumGross))
			}
		}
	}
	return total
}

// TampilanCedant - GET /api/nbfacin/kasus/{caseId}/cedant.
func (s *Service) TampilanCedant(ctx context.Context, id string) (TampilanCedant, error) {
	if s.cedant == nil || s.spreading == nil {
		return TampilanCedant{}, ErrCedantTanpaDatabase
	}
	if !idKasusSah(id) {
		return TampilanCedant{}, ErrKasusTidakAda
	}
	k, err := s.cedant.BacaCedant(ctx, id)
	if err != nil {
		return TampilanCedant{}, galatKasus(err)
	}
	sp, err := s.spreading.BacaSpreading(ctx, nil, id)
	if err != nil {
		return TampilanCedant{}, galatKasus(err)
	}
	_, jumlah, dasar := ringkasanRNM(sp)
	tsi, premi, tidakAda, err := totalKeIDR(ctx, s.spreading, sp.StartDateTime, jumlah)
	if err != nil {
		return TampilanCedant{}, err
	}
	for i := range dasar {
		d := &dasar[i]
		d.TSIRNM, d.PremiumRNM, d.TSIGross, d.PremiumGross = simpan8(d.TSIRNM), simpan8(d.PremiumRNM), simpan8(d.TSIGross), simpan8(d.PremiumGross)
	}
	return TampilanCedant{SobName: k.SobName, ShareCedantType: k.ShareCedantType, ShareOfCeding: k.ShareOfCeding,
		PercentShare: simpan8(sp.PercentShare), TotalTSIRNM: simpan8(tsi), TotalPremiRNM: simpan8(premi),
		Wajib: proteksiCedantBerlaku(gerbangKasus(k)), Cedant: k.Cedant, CedingUmum: k.CedingUmum, DasarMataUang: dasar,
		KursTidakAda: tidakAda}, nil
}

// SimpanCedant - PUT …/cedant: periksa isian, proteksi (bila berlaku), SetShareOfCeding, simpan, jawab baca ulang.
func (s *Service) SimpanCedant(ctx context.Context, pelaku inti.Pelaku, id, tipe string, baris []models.BarisCedant) (TampilanCedant, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return TampilanCedant{}, err
	}
	if err := s.wajibBolehUbah(ctx, pelaku, id); err != nil {
		return TampilanCedant{}, err
	}
	if tipe != "" && tipe != tipeCedantGross && tipe != tipeCedantShareRNM {
		return TampilanCedant{}, fmt.Errorf("%w: shareCedantType harus kosong, \"0\", atau \"1\"", ErrMasukanCedant)
	}
	if len(baris) > banyakCedant {
		return TampilanCedant{}, fmt.Errorf("%w: cedant paling banyak %d baris", ErrMasukanCedant, banyakCedant)
	}
	var masalah []string
	totalShare := nol()
	for i, b := range baris {
		if len(b.CedingCo) > lebarKodeCeding {
			masalah = append(masalah, fmt.Sprintf("cedant[%d].cedingCo paling banyak %d bita", i, lebarKodeCeding))
		}
		if b.ShareCeding != nil && lebihBesar(b.ShareCeding, seratus) {
			masalah = append(masalah, fmt.Sprintf("cedant[%d].shareCeding: %s", i, pesanShareCeding))
		}
		totalShare = tambah(totalShare, b.ShareCeding)
	}
	if lebihBesar(totalShare, seratus) { // C-7, tanpa gerbang (Pega tidak memeriksa jumlah)
		masalah = append(masalah, pesanTotalShareCedant)
	}
	if len(masalah) > 0 {
		return TampilanCedant{}, fmt.Errorf("%w: %s", ErrMasukanCedant, strings.Join(masalah, "; "))
	}
	if s.cedant == nil || s.spreading == nil || s.transaksi == nil {
		return TampilanCedant{}, ErrCedantTanpaDatabase
	}
	if !idKasusSah(id) {
		return TampilanCedant{}, ErrKasusTidakAda
	}
	k, err := s.cedant.BacaCedant(ctx, id)
	if err != nil {
		return TampilanCedant{}, galatKasus(err)
	}
	if proteksiCedantBerlaku(gerbangKasus(k)) {
		sp, err := s.spreading.BacaSpreading(ctx, nil, id)
		if err != nil {
			return TampilanCedant{}, galatKasus(err)
		}
		// K49-11: dibandingkan dalam angka asli tanpa kurs, seperti Pega (TotalPremiNusaRe mentah vs Σ premi cedant mentah).
		premiMentah, _, dasar := ringkasanRNM(sp)
		m := periksaProteksiCedant(baris, tipe, dasar, premiMentah)
		if tipe == "" {
			m = append([]string{pesanTipeCedantWajib}, m...)
		}
		if len(m) > 0 {
			return TampilanCedant{}, fmt.Errorf("%w: %s", ErrMasukanCedant, strings.Join(m, "; "))
		}
	}
	// % Share RNM untuk SetShareOfCeding dibaca di transaksi yang sama dengan penulisan (bukan dari baca sebelumnya).
	err = s.transaksi(ctx, func(tx *db.Tx) error {
		sp, err := s.spreading.BacaSpreading(ctx, tx, id)
		if err != nil {
			return err
		}
		isi := models.SimpanCedant{ShareCedantType: tipe, ShareOfCeding: shareOfCeding(tipe, sp.PercentShare), Cedant: baris}
		return s.cedant.TulisCedant(ctx, tx, id, isi)
	})
	if errors.Is(err, repository.ErrCedingTidakSah) || errors.Is(err, repository.ErrCedingTerlaluPanjang) {
		return TampilanCedant{}, fmt.Errorf("%w: %v", ErrMasukanCedant, err)
	}
	if err != nil {
		return TampilanCedant{}, galatKasus(err)
	}
	return s.TampilanCedant(ctx, id)
}

// galatKasus - repository.ErrKasusTidakAda -> ErrKasusTidakAda (404); galat lain diteruskan.
func galatKasus(err error) error {
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return ErrKasusTidakAda
	}
	return err
}
