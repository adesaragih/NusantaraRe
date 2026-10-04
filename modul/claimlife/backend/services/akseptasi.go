package services

// Akseptasi baris adjustment dari Claim Life sendiri - audit A0.
//
// Untuk apa berkas ini: jalur akseptasi yang TIDAK lewat Komite. Ia ada di
// korpus sejak awal dan terlewat sampai audit A0.
//
// `[terverifikasi]` `Claim Life/Activity/SaveAdjustment_Act.xml` pecahan baris
// 1833, 1879, 1899 menulis `.AdjustmentList(<LAST>).ACCEPTEDNO`,
// `.STS_REJECT = 1`, dan `.ACCEPTATION_DATE = @CurrentDateTime()`. Dipicu
// tombol berlabel "Save Adjustment" di `Section/ClaimLifeDetailGCNM.xml`
// (pecahan baris 22641 → 22665).
//
// ⛔ RALAT BESAR atas tiket 04, 07, dan 08. Ketiganya berdiri di atas bacaan
// bahwa akseptasi HANYA milik Komite - `WajibPeranPengubahStatus` bahkan
// menolak tujuan Aksep dengan `ErrAksepBukanDariModulIni`. Menurut XML itu
// keliru: akseptasi punya DUA jalur, dan modul ini hanya mengenal satu.
//
// Dibaca sesudah: wewenang.go, statusbaris.go.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/jejak"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/claimlife/backend/models"
	"nusantarare/modul/claimlife/backend/repository"
)

var (
	// ErrPesertaTidakDipilih - `.IsCheck` peserta bukan "true".
	ErrPesertaTidakDipilih = errors.New(
		"services: peserta tidak dipilih untuk diklaim (IS_CHECK)")
	// ErrBarisSudahBernomorAkseptasi - `.ACCEPTEDNO` sudah terisi.
	ErrBarisSudahBernomorAkseptasi = errors.New(
		"services: baris sudah punya nomor akseptasi")
)

var (
	// ErrNomorAkseptasiBerganda dipaparkan ULANG di sini.
	//
	// ⛔ Supaya `handlers` tidak perlu mengimpor `repository` hanya demi satu
	// sentinel - arah `handlers -> services -> repository` dijaga penjaga
	// statik, dan melanggarnya demi kenyamanan satu baris adalah harga yang
	// salah.
	ErrNomorAkseptasiBerganda = kontrak.ErrNomorAkseptasiBerganda
)

// Awalan nomor akseptasi per Type.
//
// `[terverifikasi]` `Claim Life/RDBList/Generate_NoAccept_Life.xml` baris 85
// dan `Generate_NoAccept_LifeRetro.xml` baris 85 - keduanya
// `SELECT '<awalan>'||{BusinessCode}||'.'||{CARI1}||'.'||{CARI2}||'.'||
// LPAD(TO_CHAR(ACCEPTATIONNOLIFE_SEQ.NEXTVAL),5,'0') FROM dual`.
const (
	AwalanAkseptasiGross = "RNML-A"
	AwalanAkseptasiRetro = "RNML-AR"
)

// ambangGeserBulan adalah tanggal yang, bila dilewati, memakai bulan berikut.
//
// `[terverifikasi]` `SaveAdjustment_Act.xml` pecahan baris 521:
// `@if(@toDecimal(Local.currentdate)>25, CurrentMonth+1, CurrentMonth)`.
// Ini menutup OQ-030.
const ambangGeserBulan = 25

// AwalanNomorAkseptasi memilih awalan menurut Type.
//
// ⛔ Type asing TIDAK ditebak. Awalan menentukan seri nomor yang tercetak di
// dokumen; menebaknya berarti menerbitkan nomor di seri yang salah.
func AwalanNomorAkseptasi(tipe string) (string, error) {
	switch tipe {
	case TypeQP, TypeQR:
		return AwalanAkseptasiGross, nil
	case TypeTP, TypeTR:
		return AwalanAkseptasiRetro, nil
	default:
		return "", fmt.Errorf("%w: %q", ErrTypeTidakDikenal, tipe)
	}
}

// PeriodeNomorAkseptasi menghitung MM dan YY nomor akseptasi.
//
// `[terverifikasi]` `SaveAdjustment_Act.xml`: baris 521 menggeser bulan bila
// tanggal > 25; baris 542 memadankannya menjadi dua digit; baris 563 mengubah
// `"13"` menjadi `"01"`.
//
// ⛔ PENYIMPANGAN SADAR - DILAPORKAN, bukan disembunyikan. `[terverifikasi]`
// baris 605 berbunyi:
//
//	@if(MM=="12" && NextMonth=="01", @toDecimal(@CurrentDate("YY")),
//	                                 @toDecimal(@CurrentDate("YY")))
//
// KEDUA cabangnya IDENTIK. Pega karena itu TIDAK menggeser tahun ketika bulan
// berguling 12 → 01, sehingga nomor Januari memakai tahun Desember - dua
// nomor di tahun berbeda dapat bertabrakan. Go menggeser tahunnya.
//
// Bila paritas persis yang diinginkan, ubah di satu tempat ini; test
// `TestPeriodeNomorAkseptasiAmbangDuaPuluhLima` yang merah lebih dulu.
func PeriodeNomorAkseptasi(saat time.Time) (mm, yy string) {
	bulan := int(saat.Month())
	tahun := saat.Year()
	if saat.Day() > ambangGeserBulan {
		bulan++
		if bulan == 13 {
			bulan = 1
			// ⛔ Inilah penyimpangannya: Pega tidak menaikkan tahun di sini.
			tahun++
		}
	}
	return fmt.Sprintf("%02d", bulan), fmt.Sprintf("%02d", tahun%100)
}

// RakitNomorAkseptasi menyusun nomor utuhnya.
//
// ⚠️ `urut` dipadankan ke LIMA digit seperti `LPAD(...,5,'0')`, tetapi urut
// yang sudah lebih panjang TIDAK dipotong: memotongnya menerbitkan nomor yang
// bertabrakan dengan nomor lama.
func RakitNomorAkseptasi(awalan, kodeBisnis, mm, yy, urut string) string {
	if len(urut) < 5 {
		urut = strings.Repeat("0", 5-len(urut)) + urut
	}
	return awalan + kodeBisnis + "." + mm + "." + yy + "." + urut
}

// PeriksaBolehAksep menegakkan prasyarat XML jalur akseptasi.
//
// `[terverifikasi]` `SaveAdjustment_Act.xml` pecahan baris 854 dan 1048:
// `.ACCEPTEDNO=="" && .IsCheck=true && .STS_REJECT=="0"`.
//
// ⛔ PENYIMPANGAN SADAR - DILAPORKAN. Prasyarat XML-nya ditulis TANPA KURUNG:
// `.IsCheck=true && Type=="QP" || Type=="QR"` (baris 837 dan 1031). Karena
// `&&` mengikat lebih erat daripada `||`, bacaan harfiahnya meloloskan `QR`
// dan `TR` TANPA memeriksa `IsCheck` sama sekali. Go memakai bacaan yang
// DIMAKSUD - `IsCheck && (QP || QR)` - sebab meloloskan baris peserta yang
// tidak dipilih ke akseptasi adalah cacat, bukan fitur.
func PeriksaBolehAksep(p models.Peserta, b models.BarisAdjustment) error {
	if !PesertaDipilih(p) {
		return fmt.Errorf("%w: peserta %q ber-IS_CHECK %q",
			ErrPesertaTidakDipilih, p.ID, p.IsCheck)
	}
	if strings.TrimSpace(b.NomorAkseptasi) != "" {
		return fmt.Errorf("%w: baris %q bernomor %q",
			ErrBarisSudahBernomorAkseptasi, b.ID, b.NomorAkseptasi)
	}
	if b.KodeStatus != kontrak.KodeOutstanding {
		return fmt.Errorf("%w: baris %q berkode %q",
			ErrBarisBukanOutstanding, b.ID, b.KodeStatus)
	}
	return nil
}

// WajibPemegangTahap menuntut pelaku memegang tahap klaim saat ini.
//
// `[terverifikasi — pohon XML]` Tombol "Save Adjustment" TIDAK dibungkus
// gerbang peran mana pun di `ClaimLifeDetailGCNM`: seluruh leluhurnya
// bervisibilitas ALWAYS, dan `pyCondition 1=2` yang tampak bertetangga adalah
// `pyContainerVisibleWhen` layout LAIN. Yang menentukan siapa boleh
// menekannya adalah SIAPA MEMEGANG TAHAP - sectionnya diluncurkan dari layar
// tahap lewat flow action `ViewClaimDetailLifeGCNM`.
//
// `[terverifikasi]` `Flow/Register_Flow.xml`: Outstanding Claim →
// `ReasLifeAdmin`; Medical Check → workbasket `ReasLifeMedicalAdvisor`;
// Claim Analis → workbasket `ReasLifeSPV`.
func WajibPemegangTahap(pelaku inti.Pelaku, tahap models.Tahap) error {
	peran, ada := models.PeranPemegangTahap(tahap)
	if !ada {
		return fmt.Errorf("%w: tahap %v tidak punya pemegang",
			ErrTahapTidakDikenal, tahap)
	}
	return inti.WajibPeran(pelaku, peran)
}

// PenerbitNomorAkseptasi menerbitkan satu nomor akseptasi baru.
//
// ⛔ Antarmuka, bukan query. Ia membaca sequence Oracle dan memeriksa
// keunikannya terhadap `OS_AKSEPTASI_KLAIM_LIFE.NO_ACCEPTATION` seperti
// `GetAcceptedNoCL` - keduanya menuntut Oracle.
type PenerbitNomorAkseptasi interface {
	Terbitkan(ctx context.Context, tx *db.Tx, awalan, kodeBisnis string,
		saat time.Time) (string, error)
}

// Akseptasi menyimpan keputusan aksep dari Claim Life sendiri.
type Akseptasi struct {
	svc      *Service
	penerbit PenerbitNomorAkseptasi
	jejak    jejak.Jejak
}

// Akseptasi menyusun layanan itu.
func (s *Service) Akseptasi() *Akseptasi {
	return &Akseptasi{svc: s, jejak: jejak.JejakBelumDiputuskan{}}
}

// DenganPenerbit mengganti penerbit nomornya.
func (a *Akseptasi) DenganPenerbit(p PenerbitNomorAkseptasi) *Akseptasi {
	salin := *a
	salin.penerbit = p
	return &salin
}

// DenganJejak mengganti perekam jejaknya.
func (a *Akseptasi) DenganJejak(j jejak.Jejak) *Akseptasi {
	salin := *a
	salin.jejak = j
	return &salin
}

// SimpanAdjustment mengaksep seluruh baris peserta yang layak.
//
// Meniru `SaveAdjustment_Act` langkah 1: ia MENGULANG peserta
// (`pyWorkPage.PremiumListSummary.PremiumListDetail`), dan untuk tiap peserta
// yang lolos prasyarat menerbitkan nomor lalu menstempel baris TERAKHIR-nya.
//
// ⛔ Urutan gerbangnya bukan gaya:
//
//	identitas → pengenal → Oracle → Type → tahap → pemegang tahap →
//	prasyarat XML per peserta → nomor → stempel
//
// Peran diperiksa SESUDAH tahap dibaca, sebab pemegangnya bergantung tahap -
// itulah yang pohon XML tetapkan, bukan daftar peran datar.
func (a *Akseptasi) SimpanAdjustment(ctx context.Context, pelaku inti.Pelaku,
	klaimID, pesertaID string, saat time.Time) (string, error) {

	if err := inti.WajibIdentitas(pelaku); err != nil {
		return "", err
	}
	if strings.TrimSpace(klaimID) == "" || strings.TrimSpace(pesertaID) == "" {
		return "", fmt.Errorf("%w: pengenal klaim dan peserta wajib diisi",
			galat.ErrPermintaanTidakSah)
	}
	if !a.svc.PunyaDatabase() {
		return "", db.ErrTanpaOracle
	}

	// ⛔ BUTIR bb: kasus yang sudah ditutup tidak dapat diubah lagi.
	// Satu pintu untuk seluruh rute pengubah - lihat
	// services.PastikanKasusTerbuka, yang pemanggilannya ditagih penjaga
	// statik. Diperiksa SESUDAH wewenang: pemanggil yang tidak berhak tidak
	// berhak pula tahu keadaan kasusnya.
	if err := a.svc.PastikanKasusTerbuka(ctx, klaimID); err != nil {
		return "", err
	}
	if a.penerbit == nil {
		return "", ErrPenomorBelumDiputuskan
	}

	baca := repository.NewKlaimLife(a.svc.DB())
	tipe, err := baca.TypeKlaim(ctx, klaimID)
	if err != nil {
		return "", err
	}
	awalan, err := AwalanNomorAkseptasi(tipe)
	if err != nil {
		return "", err
	}

	// ⛔ Pemegang TAHAP, bukan daftar peran. `[terverifikasi — pohon XML]`
	// tombolnya tidak bergerbang peran; yang menentukan dari layar tahap mana
	// ia diluncurkan.
	peranAsal, err := baca.TahapKlaim(ctx, klaimID)
	if err != nil {
		return "", err
	}
	tahap := models.TahapDariPeran(peranAsal)
	if !tahap.Diketahui() {
		return "", fmt.Errorf("%w: peran pemegang %q", ErrTahapTidakDikenal, peranAsal)
	}
	if err := WajibPemegangTahap(pelaku, tahap); err != nil {
		return "", err
	}

	hdr, err := baca.AmbilHeader(ctx, klaimID)
	if err != nil {
		return "", err
	}
	if hdr == nil {
		return "", fmt.Errorf("%w: klaim %q tidak ada", galat.ErrPermintaanTidakSah, klaimID)
	}
	hdrKodeBisnis := hdr.KodeBisnis

	peserta, err := baca.AmbilPeserta(ctx, klaimID)
	if err != nil {
		return "", err
	}
	ps, ada := cariPeserta(peserta, pesertaID)
	if !ada {
		return "", fmt.Errorf("%w: peserta %q bukan milik klaim %q",
			galat.ErrPermintaanTidakSah, pesertaID, klaimID)
	}
	perBaris, err := baca.AmbilBaris(ctx, klaimID)
	if err != nil {
		return "", err
	}
	daftar := perBaris[pesertaID]
	if len(daftar) == 0 {
		return "", fmt.Errorf("%w: peserta %q belum punya baris adjustment",
			galat.ErrPermintaanTidakSah, pesertaID)
	}
	// `.AdjustmentList(<LAST>)` - baris TERAKHIR peserta, seperti XML.
	baris := daftar[len(daftar)-1]
	if err := PeriksaBolehAksep(ps, baris); err != nil {
		return "", err
	}

	// ⛔ Kode bisnis belum tersimpan di mana pun - lihat ErrKodeBisnisBelumTersimpan.
	kodeBisnis := strings.TrimSpace(hdrKodeBisnis)
	if kodeBisnis == "" {
		return "", fmt.Errorf("%w: klaim %q", kontrak.ErrKodeBisnisBelumTersimpan, klaimID)
	}

	var nomor string
	err = a.svc.DalamTransaksi(ctx, func(tx *db.Tx) error {
		n, err := a.penerbit.Terbitkan(ctx, tx, awalan, kodeBisnis, saat)
		if err != nil {
			return err
		}
		nomor = n
		// Mesin transisi tiket 04 dipakai APA ADANYA - kefinalan, pencerminan
		// tiga tingkat, dan jejaknya. Nol aturan status ditulis ulang di sini.
		if err := baca.PerbaruiStatusBaris(ctx, tx, pesertaID, baris.ID,
			kontrak.KodeOutstanding, kontrak.KodeAksep, nomor, saat); err != nil {
			return err
		}
		if err := baca.CerminkanHeader(ctx, tx, klaimID, kontrak.KodeAksep, nomor); err != nil {
			return err
		}
		return a.jejak.Rekam(ctx, tx, jejak.CatatanJejak{
			AdjustmentID: baris.ID,
			KlaimID:      klaimID,
			Dari:         baris.KodeStatus,
			Ke:           kontrak.KodeAksep,
			AkunID:       pelaku.AkunID,
			Waktu:        saat,
		})
	})
	if err != nil {
		return "", err
	}
	return nomor, nil
}

// cariPeserta menunjuk satu peserta dari daftarnya.
func cariPeserta(daftar []models.Peserta, id string) (models.Peserta, bool) {
	for _, p := range daftar {
		if p.ID == id {
			return p, true
		}
	}
	return models.Peserta{}, false
}

// penerbitOracle merakit nomor akseptasi dari urut yang Oracle terbitkan.
//
// ⛔ Perakitannya DI SINI, bukan di repository: bentuk nomor adalah aturan
// dagang, dan arah `services -> repository` tidak dibalik demi kenyamanan.
// Yang diminta dari repository hanya dua hal yang hanya Oracle tahu - urut
// berikut, dan apakah nomornya sudah dipakai.
type penerbitOracle struct{ pohon *repository.PohonKlaim }

// PenerbitAkseptasiOracle menyusun penerbit yang memakai sequence Oracle.
func PenerbitAkseptasiOracle(svc *Service) PenerbitNomorAkseptasi {
	return penerbitOracle{pohon: repository.NewPohonKlaim(svc.DB())}
}

// Terbitkan menerbitkan satu nomor akseptasi yang belum pernah dipakai.
func (p penerbitOracle) Terbitkan(ctx context.Context, tx *db.Tx,
	awalan, kodeBisnis string, saat time.Time) (string, error) {

	urut, err := p.pohon.UrutAkseptasiBerikut(ctx, tx)
	if err != nil {
		return "", err
	}
	mm, yy := PeriodeNomorAkseptasi(saat)
	nomor := RakitNomorAkseptasi(awalan, kodeBisnis, mm, yy, urut)

	dipakai, err := p.pohon.NomorAkseptasiDipakai(ctx, tx, nomor)
	if err != nil {
		return "", err
	}
	if dipakai {
		return "", fmt.Errorf("%w: %q", kontrak.ErrNomorAkseptasiBerganda, nomor)
	}
	return nomor, nil
}

// PesertaDipilih menjawab apakah peserta itu ditandai untuk diklaim.
//
// ⛔ SATU rumah bagi aturannya. Sebelum ini perbandingannya tertulis di
// tempat pemakaiannya, dan jalur penulis di repository memakai nilai LAIN -
// tak satu pun uji memaksa keduanya bertemu.
//
// ⚠️ Longgar saat MEMBACA (EqualFold + TrimSpace) tetapi ketat saat
// MENULIS (models.PenandaDipilih). Data warisan boleh saja menyimpan "TRUE"
// atau " true "; yang tidak boleh adalah kita sendiri menambah ragam baru.
func PesertaDipilih(p models.Peserta) bool {
	return strings.EqualFold(strings.TrimSpace(p.IsCheck), models.PenandaDipilih)
}
