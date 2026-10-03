// Package services memuat layanan modul NB Fac In yang dipanggil handlers (tiket
// 20). Perilakunya tetap mesin yang sudah diport - premium, acceptance, rules -
// lewat penyedia kontrak `kontrakfacin`; paket ini hanya menyambung repository
// dan memilah galat.
package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/inti/backend/uang"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
	"nusantarare/modul/nbfacin/backend/services/acceptance"
	"nusantarare/modul/nbfacin/backend/services/kontrakfacin"
	"nusantarare/modul/nbfacin/backend/services/premium"
	"nusantarare/modul/nbfacin/backend/services/rules"
)

// ErrTidakDapatDiproses - masukan atau kasus yang mesin tolak (angka tak terbaca,
// rumus belum diport, keadaan tangga tak pasti, predikat belum diport). Handler
// menjawabnya 422. ErrTanpaDatabase dan ErrTabelLimitTakTersedia 503; galat lain
// (repository, bug program) 500.
var ErrTidakDapatDiproses = errors.New("services: masukan tidak dapat diproses mesin NB")

// ErrTanpaDatabase - layanan dirakit tanpa basis data: tangga akseptasi butuh tabel
// limit Oracle. Hitung premi tidak terdampak. Handler menjawabnya 503.
var ErrTanpaDatabase = errors.New("services: basis data tidak dikonfigurasi, tabel limit akseptasi tidak terbaca")

// ErrTabelLimitTakTersedia - keenam tabel limit terbaca tetapi kosong: tangga tidak
// boleh "selesai tanpa approver" karena datanya belum tersambung. Handler 503.
var ErrTabelLimitTakTersedia = errors.New("services: tabel limit akseptasi kosong atau belum tersambung")

// DariDasar merakit layanan dari dasar aplikasi: pembaca tabel limit Oracle bila
// basis data tersedia, tanpa pembaca bila tidak.
func DariDasar(d *inti.Dasar) *Service {
	if !d.PunyaDatabase() {
		return Baru(nil)
	}
	return Baru(repository.NewLimitOracle(d.DB())).DenganAkun(repository.NewAkunOracle(d.DB())).
		DenganKelasBisnis(repository.NewKelasBisnisOracle(d.DB())).
		DenganTransaksi(d.DalamTransaksi).DenganCaseNB(repository.NewCaseNBOracle(d.DB())).
		DenganKasus(repository.NewKasusOracle(d.DB())).
		DenganMarketingOfficer(repository.NewMarketingOfficerOracle(d.DB())).
		DenganPortal(repository.NewPortalOracle(d.DB())).
		DenganSOB(repository.NewSOBOracle(d.DB())).
		DenganObjek(repository.NewObjekOracle(d.DB())).
		DenganRisk(repository.NewRiskOracle(d.DB())).
		DenganRW(repository.NewRWOracle(d.DB())).
		DenganOccupation(repository.NewOccupationOracle(d.DB())).
		DenganPilihanItem(repository.NewPilihanItemOracle(d.DB())).
		DenganTableOfLimit(repository.NewTableOfLimitOracle(d.DB()))
}

// Service - layanan NB Fac In.
type Service struct {
	limit repository.PembacaLimit
	// tangga - tangga yang dirakit pemanggil (BaruDenganTangga); nil = disusun dari
	// tabel limit repository setiap permintaan.
	tangga kontrak.TanggaAkseptasiFacIn
	// akun - pembaca POOLDATA.T_M_ACCOUNT (tiket 27); nil = tanpa basis data (503).
	akun repository.PembacaAkun
	// kelasBisnis - pembaca POOLDATA.BUSINESS (tiket 28); nil = tanpa basis data (503).
	kelasBisnis repository.PembacaKelasBisnis
	// caseNB, transaksi - pembuat case NB (tiket 29) dan pembuka transaksinya
	// (inti.Dasar.DalamTransaksi); nil = tanpa basis data (503).
	caseNB    repository.PenulisCaseNB
	transaksi Transaksi
	// kasus, marketing, jam - layar Inward Facultative (tiket 31); nil = tanpa basis
	// data (503); jam nil = time.Now.
	kasus     repository.PenyimpanKasus
	marketing repository.PembacaMarketingOfficer
	// portal - daftar case NB portal Opportunity (tiket 32); nil = tanpa basis data (503).
	portal repository.PembacaPortal
	// sob - pilihan SOB popup Change SOB (tiket 33); nil = tanpa basis data (503).
	sob repository.PembacaSOB
	// objek - tab Object FIRE (tiket 35); nil = tanpa basis data (503).
	objek repository.PenyimpanObjek
	// risk - pencarian alamat risiko (tiket 36); nil = tanpa basis data (503).
	risk repository.PembacaRisk
	// rw - saran Zip Code dan simpan alamat baru (tiket 37); nil = tanpa basis data (503).
	rw repository.PenyimpanRW
	// occupation - saran Occupation Surrounding Risk (tiket 38); nil = tanpa basis data (503).
	occupation repository.PembacaOccupation
	// jenisItem, mataUang - pilihan Object Item Type / Currency (tiket 39); nil = tanpa basis data (503).
	jenisItem repository.PembacaJenisItem
	mataUang  repository.PembacaMataUang
	// tableOfLimit - pilihan Class of Construction (tiket 40); nil = tanpa basis data (503).
	tableOfLimit repository.PembacaTableOfLimit
	jam          func() time.Time
}

// DenganAkun memasang pembaca tabel akun (tiket 27).
func (s *Service) DenganAkun(a repository.PembacaAkun) *Service { s.akun = a; return s }

// DenganKelasBisnis memasang pembaca tabel bisnis (tiket 28).
func (s *Service) DenganKelasBisnis(k repository.PembacaKelasBisnis) *Service {
	s.kelasBisnis = k
	return s
}

// BaruDenganTangga merakit layanan atas tangga yang sudah dirakit (uji HTTP, atau
// pemakai yang memegang tabel limit sendiri) - tanpa repository.
func BaruDenganTangga(tangga kontrak.TanggaAkseptasiFacIn) *Service { return &Service{tangga: tangga} }

// Baru merakit layanan.
func Baru(limit repository.PembacaLimit) *Service { return &Service{limit: limit} }

// HasilPremi - premi satu coverage dan rule asal rumusnya.
type HasilPremi struct {
	Premi     uang.Money
	AsalRumus string
}

// HitungPremi - premium.Calculate lewat kontrakfacin.Premi. Lini yang tidak ada di
// peta K-018 ditolak lebih dulu (SatuanRate akan panic).
//
// Panic mesin dipilah lewat TIPE (galatDariPanic): `premium.PanikLini` = masukan (422),
// selain itu bug program (500). Pra-cek LiniDikenal tetap ada sebagai pesan yang jelas.
func (s *Service) HitungPremi(in kontrak.MasukanPremiFacIn) (hasil HasilPremi, err error) {
	defer func() {
		if r := recover(); r != nil {
			hasil, err = HasilPremi{}, galatDariPanic(r)
		}
	}()
	if !premium.LiniDikenal(premium.LiniBisnis(in.LiniBisnis)) {
		return HasilPremi{}, fmt.Errorf("%w: lini bisnis %q tidak dikenal", ErrTidakDapatDiproses, in.LiniBisnis)
	}
	mesin := kontrakfacin.Premi{}
	premi, err := mesin.Hitung(in)
	if err != nil {
		return HasilPremi{}, fmt.Errorf("%w: %w", ErrTidakDapatDiproses, err)
	}
	return HasilPremi{Premi: premi, AsalRumus: mesin.AsalRumus(in)}, nil
}

// LangkahAkseptasi - satu langkah tangga (bentuk A, beralih ke B) atas tabel limit
// dari repository. Predikat bersikap panic (belum diport, termasuk pembanding data
// identitas) menjadi ErrTidakDapatDiproses, bukan jawaban 500 tanpa sebab.
func (s *Service) LangkahAkseptasi(ctx context.Context, k kontrak.KasusFacIn, p kontrak.PenggunaFacIn) (tr kontrak.TransisiFacIn, err error) {
	tangga := s.tangga
	if tangga == nil {
		if s.limit == nil {
			return kontrak.TransisiFacIn{}, ErrTanpaDatabase
		}
		a, b, err := s.limit.MuatLimit(ctx)
		if err != nil {
			return kontrak.TransisiFacIn{}, err
		}
		if tangga, err = susunTangga(a, b); err != nil {
			return kontrak.TransisiFacIn{}, err
		}
	}
	defer func() {
		if r := recover(); r != nil {
			tr, err = kontrak.TransisiFacIn{}, galatDariPanic(r)
		}
	}()
	tr, err = tangga.Langkah(k, p)
	switch {
	case errors.Is(err, kontrakfacin.ErrTabelLimitKosong):
		return kontrak.TransisiFacIn{}, fmt.Errorf("%w: %w", ErrTabelLimitTakTersedia, err)
	case err != nil:
		return kontrak.TransisiFacIn{}, fmt.Errorf("%w: %w", ErrTidakDapatDiproses, err)
	}
	return tr, nil
}

// galatDariPanic - panic yang TIPE-nya menandai keadaan yang diketahui -
// rules.PanikSikap (predikat belum diport) dan premium.PanikLini (lini di luar peta
// skala) - menjadi ErrTidakDapatDiproses (422). Panic lain = bug program → galat
// biasa (500, dicatat handler), bukan disamarkan sebagai galat masukan.
func galatDariPanic(r any) error {
	if e, ok := r.(error); ok {
		var sikap rules.PanikSikap
		var lini premium.PanikLini
		if errors.As(e, &sikap) || errors.As(e, &lini) {
			return fmt.Errorf("%w: %w", ErrTidakDapatDiproses, e)
		}
	}
	return fmt.Errorf("services: panic mesin NB: %v", r)
}

// tabelDikenal - lima tabel bentuk A, ejaan repository = ejaan acceptance.
var tabelDikenal = map[string]acceptance.NamaTabel{
	string(acceptance.TabelProperty):                    acceptance.TabelProperty,
	string(acceptance.TabelPropertyNonPreferred):        acceptance.TabelPropertyNonPreferred,
	string(acceptance.TabelPropertyPreferredCommercial): acceptance.TabelPropertyPreferredCommercial,
	string(acceptance.TabelEngineering):                 acceptance.TabelEngineering,
	string(acceptance.TabelNonPropEng):                  acceptance.TabelNonPropEng,
}

// susunTangga - baris repository → kontrakfacin.Tangga. Tabel tak dikenal = galat
// konfigurasi, bukan diabaikan.
func susunTangga(a []models.BarisLimitA, b []models.BarisLimitB) (kontrakfacin.Tangga, error) {
	t := kontrakfacin.Tangga{BentukA: acceptance.TabelLimit{}}
	for _, r := range a {
		nama, ada := tabelDikenal[r.Tabel]
		if !ada {
			return kontrakfacin.Tangga{}, fmt.Errorf("services: tabel limit %q tidak dikenal tangga", r.Tabel)
		}
		t.BentukA[nama] = append(t.BentukA[nama], acceptance.BarisLimit{Jabatan: acceptance.Jabatan(r.Jabatan),
			TeamGroup: r.TeamGroup, LimitBottom: r.LimitBottom, LimitBottom2: r.LimitBottom2})
	}
	for _, r := range b {
		t.BentukB = append(t.BentukB, acceptance.BarisFinancial{Jabatan: acceptance.Jabatan(r.Jabatan),
			LimitBond: r.LimitBond, LimitCreditCL: r.LimitCreditCL, LimitCreditNCL: r.LimitCreditNCL})
	}
	if len(t.BentukA) == 0 {
		t.BentukA = nil
	}
	return t, nil
}

// --- Lookup akun ChooseAccount (tiket 27) ---

// UkuranHalamanAkun - keputusan work owner 02-10-2026 "15 baris per halaman" (A71 diubah;
// semula 20 = UKURAN_HALAMAN inti frontend).
const UkuranHalamanAkun = 15

// batasCari - A73: kolom terpanjang T_M_ACCOUNT 255 karakter.
const batasCari = 255

// halamanMaks - offset (halaman-1)*ukuran tetap dalam int32 (bind Oracle), tidak meluap.
const halamanMaks = (1<<31-1)/UkuranHalamanAkun + 1

// ErrMasukanAkun - halaman bukan bilangan bulat >= 1 atau cari terlalu panjang. 400.
var ErrMasukanAkun = fmt.Errorf("services: halaman harus bilangan bulat >= 1 dan cari paling banyak %d karakter", batasCari)

// ErrAkunTanpaDatabase - layanan dirakit tanpa basis data: tabel akun tidak terbaca. 503.
var ErrAkunTanpaDatabase = errors.New("services: basis data tidak dikonfigurasi, tabel akun T_M_ACCOUNT tidak terbaca")

// HasilCariAkun - satu halaman lookup akun.
type HasilCariAkun struct {
	Baris          []models.Akun
	Total, Halaman int
	Ukuran         int
}

// CariAkun - halaman ke-`halaman` (mulai 1) baris T_M_ACCOUNT yang "mengandung" `cari`
// (TIDAK peka huruf besar-kecil, butir 75) di INSUREDID, INSUREDNAME, atau
// GROUPBUSINESS (A69); `cari` kosong = semua baris. Urutan INSUREDID, ID (A72).
func (s *Service) CariAkun(ctx context.Context, cari string, halaman int) (HasilCariAkun, error) {
	if halaman < 1 || halaman > halamanMaks || utf8.RuneCountInString(cari) > batasCari {
		return HasilCariAkun{}, ErrMasukanAkun
	}
	if s.akun == nil {
		return HasilCariAkun{}, ErrAkunTanpaDatabase
	}
	baris, total, err := s.akun.CariAkun(ctx, cari, (halaman-1)*UkuranHalamanAkun, UkuranHalamanAkun)
	if err != nil {
		return HasilCariAkun{}, err
	}
	return HasilCariAkun{Baris: baris, Total: total, Halaman: halaman, Ukuran: UkuranHalamanAkun}, nil
}

// --- Class Of Business (tiket 28) ---

// batasGroupBusiness - A76: BUSINESSGROUPID VARCHAR2(4000 BYTE); masukan lebih panjang
// tidak mungkin cocok.
const batasGroupBusiness = 4000

// ErrMasukanKelasBisnis - groupBusinessId kosong atau terlalu panjang. 400.
var ErrMasukanKelasBisnis = fmt.Errorf("services: groupBusinessId wajib diisi, paling banyak %d byte", batasGroupBusiness)

// ErrKelasBisnisTanpaDatabase - layanan dirakit tanpa basis data: tabel bisnis tidak
// terbaca. 503.
var ErrKelasBisnisTanpaDatabase = errors.New("services: basis data tidak dikonfigurasi, tabel bisnis BUSINESS tidak terbaca")

// KelasBisnis - semua pilihan Class Of Business milik group business `groupBusinessID`
// (diteruskan apa adanya, tidak dipangkas); kosong/spasi saja = 400. Urutan dari
// repository (A74, A75).
func (s *Service) KelasBisnis(ctx context.Context, groupBusinessID string) ([]models.KelasBisnis, error) {
	if strings.TrimSpace(groupBusinessID) == "" || len(groupBusinessID) > batasGroupBusiness {
		return nil, ErrMasukanKelasBisnis
	}
	if s.kelasBisnis == nil {
		return nil, ErrKelasBisnisTanpaDatabase
	}
	return s.kelasBisnis.KelasBisnis(ctx, groupBusinessID)
}
