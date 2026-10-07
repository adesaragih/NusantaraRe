package services

// Tab Inw Fac Cedant Panels kasus FIRE (tiket 49). Port `[terverifikasi]` `D:\migrasi\RNM\NB FacIn\` (dibaca 05-10-2026):
//   - Activity\SetShareOfCeding.xml: ShareCedantType==0 -> QuotationData.ShareOfCeding = "100%"; ==1 -> PercentShare + "%"
//     (salinan ke pyWorkPage.Quotation.ShareOfCeding tidak dibawa - halaman Quotation work tidak disimpan modul ini, K49-4);
//   - Activity\ProtectShareCedant_Act.xml (dari Protection_Act): gerbang langkah 1 + pesan verbatim (K49-2);
//   - Activity\SetTSIPremiCedant_Act.xml: ShareCeding > 100 -> "Value can't be more than 100 or less than 0!" (tanpa
//     gerbang; CurrencyList per cedant belum, C-1);
//   - Activity\CountPremiAndTSIRNMFireMBU_ACT.xml langkah 6.2.2.1.2 / 6.2.2.1.3.3 / 6.2.3: TotalTSINusaReSpreading =
//     PercentShare × Σ CoverageList(1).TSILiability tiap item / 100, TotalPremiNusaRe = Σ PremiNusantaraRe (K49-3).

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
	// pesanTipeCedantWajib - pesan agent: Share Cedant Type wajib (Section `.FlagOnGoingPolicy ==2 || IsGroup`, K49-2).
	pesanTipeCedantWajib = "Share Cedant Type wajib diisi"
	tipeCedantGross      = "0" // DDL\ShareCedantType.xml
	tipeCedantShareRNM   = "1"
	banyakCedant         = 100
	lebarKodeCeding      = 50 // T_CEDINGCEDANTLIST.CEDING_CO
)

// DenganCedant memasang penyimpan tab Cedant (tiket 49).
func (s *Service) DenganCedant(p repository.PenyimpanCedant) *Service {
	s.cedant = p
	return s
}

// TampilanCedant - jawaban GET / PUT tab Cedant.
type TampilanCedant struct {
	SobName, ShareCedantType                 string
	PercentShare, TotalTSIRNM, TotalPremiRNM *apd.Decimal
	Wajib                                    bool
	Cedant, CedingUmum                       []models.BarisCedant
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

// periksaProteksiCedant - ProtectShareCedant_Act langkah 2-10 (tanpa "Total premi share cedant must be equal total premi
// RNM!", ditunda C-1). Urutan pesan = urutan Page-Set-Messages Pega (langkah 8 daftar kosong, 9 share, 10 ceding).
func periksaProteksiCedant(baris []models.BarisCedant) []string {
	var m []string
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

// totalRNM - TotalTSINusaReSpreading / TotalPremiNusaRe dari pohon objek (NR dihitung ulang dari % Share RNM tersimpan,
// seperti tiap Save Pega); nil bila % Share RNM kosong.
func totalRNM(k models.KasusSpreading) (tsi, premi *apd.Decimal) {
	if k.PercentShare == nil {
		return nil, nil
	}
	hitungNR(&k)
	jumlahTSI, jumlahPremi := nol(), nol()
	for _, l := range k.Lokasi {
		for _, it := range l.Items {
			if len(it.Coverages) > 0 {
				jumlahTSI = tambah(jumlahTSI, it.Coverages[0].TSILiability)
			}
			for _, c := range it.Coverages {
				jumlahPremi = tambah(jumlahPremi, c.PremiNusantaraRe)
			}
		}
	}
	return simpan8(persenDari(k.PercentShare, jumlahTSI)), simpan8(jumlahPremi)
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
		return TampilanCedant{}, galatKasusCedant(err)
	}
	sp, err := s.spreading.BacaSpreading(ctx, nil, id)
	if err != nil {
		return TampilanCedant{}, galatKasusCedant(err)
	}
	tsi, premi := totalRNM(sp)
	return TampilanCedant{SobName: k.SobName, ShareCedantType: k.ShareCedantType, PercentShare: simpan8(sp.PercentShare),
		TotalTSIRNM: tsi, TotalPremiRNM: premi, Wajib: proteksiCedantBerlaku(gerbangKasus(k)), Cedant: k.Cedant,
		CedingUmum: k.CedingUmum}, nil
}

// SimpanCedant - PUT …/cedant: periksa isian, proteksi (bila berlaku), SetShareOfCeding, simpan, jawab baca ulang.
func (s *Service) SimpanCedant(ctx context.Context, pelaku inti.Pelaku, id, tipe string, baris []models.BarisCedant) (TampilanCedant, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return TampilanCedant{}, err
	}
	if tipe != "" && tipe != tipeCedantGross && tipe != tipeCedantShareRNM {
		return TampilanCedant{}, fmt.Errorf("%w: shareCedantType harus kosong, \"0\", atau \"1\"", ErrMasukanCedant)
	}
	if len(baris) > banyakCedant {
		return TampilanCedant{}, fmt.Errorf("%w: cedant paling banyak %d baris", ErrMasukanCedant, banyakCedant)
	}
	var masalah []string
	for i, b := range baris {
		if len(b.CedingCo) > lebarKodeCeding {
			masalah = append(masalah, fmt.Sprintf("cedant[%d].cedingCo paling banyak %d bita", i, lebarKodeCeding))
		}
		if b.ShareCeding != nil && lebihBesar(b.ShareCeding, seratus) {
			masalah = append(masalah, fmt.Sprintf("cedant[%d].shareCeding: %s", i, pesanShareCeding))
		}
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
		return TampilanCedant{}, galatKasusCedant(err)
	}
	if proteksiCedantBerlaku(gerbangKasus(k)) {
		m := periksaProteksiCedant(baris)
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
		return TampilanCedant{}, galatKasusCedant(err)
	}
	return s.TampilanCedant(ctx, id)
}

// galatKasusCedant - repository.ErrKasusTidakAda -> ErrKasusTidakAda (404); galat lain diteruskan.
func galatKasusCedant(err error) error {
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return ErrKasusTidakAda
	}
	return err
}
