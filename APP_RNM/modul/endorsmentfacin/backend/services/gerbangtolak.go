package services

// Enam gerbang penolakan endorsement dan empat klep pembatalnya - tiket E04.
//
// Untuk apa berkas ini: `Endorsment Fac In/Activity/SetErrorBatalEndorsement_Act.xml`
// (19 langkah; gerbang di langkah 5, 24 sub-langkah) - pemeriksaan sebelum
// kasus endorsement boleh dibuat.
//
// Dibaca sesudah: bukakasus.go.
//
// ⛔ Hasil query dan hasil predikat lini DITERIMA SEBAGAI MASUKAN. Pelaksanaan
// query milik repository (E17) - termasuk `BrowseOpenProteksiEdm_RD`, yang
// dirujuk sebagai teks runtime lewat `pxRetrieveReportData` dan isi tabelnya
// `OPENPROTEKSI_EDM` di luar korpus (K-049: rule di dalam lingkup, isi tabel
// tidak; ketiadaannya bukan `panic`). Predikat lini milik registry (E01).
//
// ⛔ `[terverifikasi]` Prakondisi yang dieksekusi TIDAK memuat satu pun nomor
// polis. Dua cara, sepakat (di `Endorsment Fac In/Activity`):
//
//	grep -oi 'RNM-[A-Z0-9.-]*' SetErrorBatalEndorsement_Act.xml | wc -l            # 560
//	grep -oi 'RNM-[A-Z0-9.-]*' SetErrorBatalEndorsement_Act.xml | sort -u | wc -l  # 176 unik
//	py docs/alat/langkah.py --lokasi-pola SetErrorBatalEndorsement_Act.xml
//	  # 557 …/pyStepsPreCondParams/pyExpressionGadget/pyExpressionMapNew (cache editor)
//	  #   2 …/pyStepsDescription · 1 …/pyExpressionGadget/…/pyCompany · total 560
//
// Nol di `pyStepsPreCondParamsWhen` dan `PropertiesValue`. Spec alur masuk §4.1
// mencacah 502 di berkas ini - cacah cache, bukan logika.
//
// Langkah 4 berlabel `//` (di-remark, butir 43) - tidak diport.

import (
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/utils"
	"nusantarare/modul/endorsmentfacin/backend/models"
)

// Teks pesan - langkah 2, apa adanya.
const (
	pesanBatal           = pesanSudahBatal
	pesanEDMBelumSelesai = "There's EDM with this policy no that haven't finish yet! "
	pesanKlaim           = "There's already a claim with this policy no"
	pesanPembayaran      = "There's already payment with this policy no"
	pesanRenewal         = "There's RNW with this policy no!"
	pesanFacout          = "This policy is not spreading FACOUT"
	pesanBisnis          = "Invalid business, please contact IT!"
	pesanPolis           = "Invalid policy no!"
)

// ErrTafsirKodeBerbeda - kode yang dibandingkan dengan angka (`CARI20==1`,
// `EdmType==1`) terbaca angka yang sama tetapi teksnya berbeda (mis. "1.0").
// Tipe propertinya belum terverifikasi; ditolak, tidak ditebak (sikap registry
// NB, keputusan work owner 01-10-2026 butir 20).
var ErrTafsirKodeBerbeda = errors.New("endorsement: kode terbaca angka tetapi teksnya berbeda; tipe properti belum terverifikasi")

// KlepPembatal - `Param.Type` laporan `BrowseOpenProteksiEdm_RD`.
type KlepPembatal int

const (
	KlepBatal      KlepPembatal = 1 // OpenProteksiEDMBatal - gerbang 1
	KlepKlaim      KlepPembatal = 2 // OpenProteksiEDMKlaim - gerbang 3
	KlepPembayaran KlepPembatal = 3 // OpenProteksiEDMPembayaran - gerbang 4
	KlepRenewal    KlepPembatal = 4 // OpenProteksiEDMRenewal - gerbang 5
)

// HasilQueryGerbang - hasil query pemeriksa, dari repository.
type HasilQueryGerbang struct {
	// StatusEDM - `GetEDMStatus_SQL` → `OutputData1.pxResults(1).CARI20`.
	StatusEDM string
	// EDMBelumSelesai - `GetListEdm` → `ListEdm.pxResults(*).pyLabel`.
	EDMBelumSelesai []string
	// Klaim - `GetDataClaim_SQL` → `ListClaim.pxResults(*).CARI1`.
	Klaim []string
	// CacahPembayaran - baris `SearcStatusBayarArasaps_SQL`.
	CacahPembayaran int
	// CacahRenewal - baris `GetListRNWbyNopolis_SQL`.
	CacahRenewal int
	// CacahFacout - baris `GetFacoutList_SQL` (dijalankan hanya untuk EDM retro).
	CacahFacout int
	// Klep - `BrowseOpenProteksiEdm_RD` per `Param.Type`: ada baris = terisi.
	// Klep yang tidak ada di peta = kosong → penolakan berlaku penuh.
	Klep map[KlepPembatal]bool
	// JenisBisnis - `GetBusinessType_Sql` → `OutData.pxResults(1).CARI2`,
	// SESUDAH `JenisBisnisUntukPredikat`.
	JenisBisnis string
}

// MasukanGerbang - satu pemeriksaan.
type MasukanGerbang struct {
	PolicyNo                 string
	EdmType                  models.JenisEndorsemen // `.Quotation.EdmType` (K-029)
	EdmTypeNew               string                 // `.Quotation.EdmTypeNew`
	EndorsementInternalRetro string
	Query                    HasilQueryGerbang
	// Predikat - IsFire, IsAneka, IsMBU, IsMarineCargo, IsPA, IsGolfInsurance,
	// IsLife (langkah 10-16), dinilai SESUDAH query jenis bisnis (E01).
	Predikat Predikat
}

// HasilGerbang - keluaran bisnis.
type HasilGerbang struct {
	Pesan []PesanMedan
	// Ditolak - langkah 19: ada pesan → `CARIGROUP = "Error"`.
	Ditolak bool
	// SalahGerbang - langkah 5.24: ada pesan dari gerbang 1-6 → `CARI3 = "SALAH"`.
	SalahGerbang bool
}

// polaPolisLife - langkah 8 `@contains(InputData.CARI17,"RNML")` → jenis bisnis
// "Life". Literal yang sama dengan pola bypass `CheckEDMPolisDate`, tetapi
// PERAN berbeda; sengaja dua konstanta.
const polaPolisLife = "RNML"

// JenisBisnisUntukPredikat - langkah 8: polis berpola "RNML" dibaca "Life".
// Pemanggil menjalankannya SEBELUM predikat lini dinilai - urutan mengikat
// (E01).
func JenisBisnisUntukPredikat(policyNo, cari2 string) string {
	if strings.Contains(policyNo, polaPolisLife) {
		return "Life"
	}
	return cari2
}

// PeriksaGerbangPenolakan - `SetErrorBatalEndorsement_Act`.
func PeriksaGerbangPenolakan(m MasukanGerbang) (HasilGerbang, error) {
	var h HasilGerbang
	pesan := func(teks string) { h.Pesan = append(h.Pesan, PesanMedan{Medan: medanPolicyNo, Teks: teks}) }
	q := m.Query
	// 5.1 / 5.10 `InputData.CARI40=="4" && InputData.CARI41=="4"` - EDM RI slip
	// (CARI40/41 diisi langkah 2 dari `.Quotation.EdmType`/`EdmTypeNew`).
	riSlip := m.EdmType == models.EdmPerubahan && m.EdmTypeNew == "4"

	// Gerbang 1 (5.2-5.5), dilompati RI slip (5.1 → label `jmp`).
	if !riSlip {
		satu, err := samaKode(q.StatusEDM, "1")
		if err != nil {
			return HasilGerbang{}, err
		}
		dua, err := samaKode(q.StatusEDM, "2")
		if err != nil {
			return HasilGerbang{}, err
		}
		if (satu || dua) && !q.Klep[KlepBatal] {
			pesan(pesanBatal)
		}
	}
	// Gerbang 2 (5.6 `jmp` - 5.9), TANPA klep.
	if len(q.EDMBelumSelesai) > 0 {
		pesan(pesanEDMBelumSelesai + q.EDMBelumSelesai[0])
	}
	// 5.10 RI slip → label `jmp2` (5.24): gerbang 3-6 dilompati.
	if !riSlip {
		// Gerbang 3 (5.11-5.14).
		// `ListClaim.pxResults(1).CARI1!="2"` - literal BERKUTIP: pembanding teks.
		if len(q.Klaim) > 0 && q.Klaim[0] != "2" && !q.Klep[KlepKlaim] {
			pesan(pesanKlaim)
		}
		// Gerbang 4 (5.15-5.18) - hanya endorsement pembatalan (EdmType 1/2, K-029).
		// [dugaan] Kondisinya membaca `pyWorkPage.Quotation.EdmType`; halaman
		// primer activity ini (kelas portal) dianggap pyWorkPage, sama dengan
		// `.Quotation.EdmType` langkah 2.
		batal1, err := samaKode(string(m.EdmType), string(models.EdmBatalSejakSemula))
		if err != nil {
			return HasilGerbang{}, err
		}
		batal2, err := samaKode(string(m.EdmType), string(models.EdmBatalProrata))
		if err != nil {
			return HasilGerbang{}, err
		}
		if q.CacahPembayaran > 0 && (batal1 || batal2) && !q.Klep[KlepPembayaran] {
			pesan(pesanPembayaran)
		}
		// Gerbang 5 (5.19-5.22; 5.19 `pyStepsPreCondition=false` tetap jalan).
		if q.CacahRenewal > 0 && !q.Klep[KlepRenewal] {
			pesan(pesanRenewal)
		}
		// Gerbang 6 (5.23), TANPA klep. [dugaan] `@LengthOfPageList(…)=0` -
		// `=` tunggal dibaca kesetaraan.
		retro, err := samaKode(m.EndorsementInternalRetro, "1")
		if err != nil {
			return HasilGerbang{}, err
		}
		if retro && q.CacahFacout == 0 {
			pesan(pesanFacout)
		}
	}
	// 5.24 `jmp2` - baris pertama `@hasMessages` memutuskan (T=5 jalan, F=3 lewati).
	// [dugaan] `@hasMessages(myStepPage)` = pesan yang dipasang activity ini;
	// langkah 3 hanya membersihkan pesan `pyWorkPage`.
	h.SalahGerbang = len(h.Pesan) > 0

	// 8 - polis berpola RNML dibaca "Life" sebelum pemeriksaan berikutnya.
	// Predikat lini (masukan) WAJIB dinilai pemanggil atas nilai yang sama
	// (`JenisBisnisUntukPredikat`).
	jenisBisnis := JenisBisnisUntukPredikat(m.PolicyNo, q.JenisBisnis)
	// 9 jenis bisnis kosong → lompat ke label `TO` (18).
	if jenisBisnis == "" {
		pesan(pesanPolis) // 18
	} else {
		// 10-16 `Local.test = a` bila salah satu predikat lini benar; 17 pesan
		// bila tidak satu pun.
		p := m.Predikat
		if !(p.IsFire || p.IsAneka || p.IsMBU || p.IsMarineCargo || p.IsPA || p.IsGolfInsurance || p.IsLife) {
			pesan(pesanBisnis)
		}
	}
	// 19 - `@hasMessages` → `CARIGROUP = "Error"`, `SPPENumber = "Error"`.
	h.Ditolak = len(h.Pesan) > 0
	return h, nil
}

// samaKode - `<properti> == <angka>`: sama bila teksnya persis; teks yang
// terbaca angka sama tetapi berbeda teks → ErrTafsirKodeBerbeda.
func samaKode(nilai, angka string) (bool, error) {
	if nilai == angka {
		return true, nil
	}
	d, err := utils.ParseDecimal(nilai)
	if err != nil {
		return false, nil
	}
	a, _ := utils.ParseDecimal(angka)
	if d.Cmp(a) == 0 {
		return false, fmt.Errorf("%w: %q", ErrTafsirKodeBerbeda, nilai)
	}
	return false, nil
}
