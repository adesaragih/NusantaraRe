package services

// Pemeriksaan sebelum Create - irisan I2 (kontrak sesi utama 11-10-2026), nol tulis.
//
// Untuk apa berkas ini: menyusun jawaban `POST /api/endorsmentfacin/periksa` dari hasil query - enam gerbang
// `SetErrorBatalEndorsement_Act` langkah 5 (rincian per gerbang), jenis bisnis langkah 7-18, dan `CheckEDMPolisDate`.
// Logika murni: query milik repository, dijalankan `Layanan.Periksa` (layanan.go).
//
// Dibaca sesudah: gerbangtolak.go, validasitanggal.go.
//
// ⛔ Aditif: pesan dan penolakan TETAP datang dari PeriksaGerbangPenolakan dan ValidasiTanggalEndorsemen yang sudah
// ada; berkas ini hanya menambah rincian per gerbang, dan menolak (galat) bila rinciannya tidak sepakat dengan pesan
// yang dihasilkan fungsi itu.

import (
	"errors"
	"fmt"
	"time"

	"nusantarare/modul/endorsmentfacin/backend/models"
	"nusantarare/modul/endorsmentfacin/backend/services/predikat"
)

// Catatan gerbang (teks untuk layar; bukan pesan Pega).
const (
	// catatanKlepTakTerbaca - OPENPROTEKSI_EDM tidak ada di skema ini (ORA-00942): klep kosong = penolakan penuh
	// (spec 07 §9.3).
	catatanKlepTakTerbaca = "OPENPROTEKSI_EDM not found; treated as no override"
	// catatanDiklep - gerbang yang syarat tolaknya benar tetapi dibuka baris OPENPROTEKSI_EDM (K-049).
	catatanDiklep = "Opened by OPENPROTEKSI_EDM"
	// catatanGetListEdm - `GetListEdm` dipanggil lewat pxRetrieveReportData (langkah 5.6-5.7); body report tidak ada.
	catatanGetListEdm = "Unfinished EDM list not checked: report GetListEdm (ASM-FW-GISFW-Work-Endorsement) is not in corpus"
	// catatanBisnisSistemBaru - A52: lini polis sistem baru belum dapat dihitung.
	catatanBisnisSistemBaru = "Line of business of a policy created in the new system is not checked yet"
	// catatanBisnisDariTabel - A51: lini polis sistem baru dari FACINPRODUCTION → BUSINESS → BusinessType_DeT (nbfacin).
	catatanBisnisDariTabel = "Line of business taken from production tables (policy created in the new system)"
	// catatanG1SistemBaru - A51: polis sistem baru tanpa kasus EDM berproduksi.
	catatanG1SistemBaru = "No completed endorsement in the new system"
	// catatanEDMPega - gerbang 2 lolos: hanya kasus sistem baru yang diperiksa (keputusan work owner W-EDM-12).
	catatanEDMPega = "EDM cases created in Pega are not checked"
	// catatanKasusBelumAda - T_EDM_FACIN_KASUS (migrasi 262) belum dibuat DBA.
	catatanKasusBelumAda = "T_EDM_FACIN_KASUS not found; unfinished EDM cases not checked"
	// catatanArasapas - ARASAPAS.DETAIL_INVOICE tidak terbaca dari koneksi ini.
	catatanArasapas = "Payment not checked: ARASAPAS.DETAIL_INVOICE is not readable from this connection"
	// catatanPengecualianBelumAda - tabel 261 belum dibuat DBA: daftar dianggap kosong.
	catatanPengecualianBelumAda = "M_EDM_PENGECUALIAN_POLIS not found; exception list treated as empty"
)

// pesanPeriodeTidakAda - [keputusan agent] `GetStartDate` tanpa baris (atau BEGINDATE/ENDDATE NULL). Teks BUKAN dari
// korpus: perilaku Pega (`@toDate("")`) belum terverifikasi; dipilih menolak (tombol OK mati) daripada meloloskan.
const pesanPeriodeTidakAda = "Policy period not found in facinproduction"

// ErrRincianTidakSepakat - rincian gerbang berkas ini berbeda dari pesan PeriksaGerbangPenolakan: kesalahan program.
var ErrRincianTidakSepakat = errors.New("endorsement: rincian gerbang tidak sepakat dengan pesan gerbang")

// PeriodeMentah - BEGINDATE / ENDDATE `GetStartDate` sebagai teks `YYYYMMDDHHMMSS` jam dinding DATE Oracle (WIB).
type PeriodeMentah struct {
	Mulai, Akhir string
}

// DataPeriksa - masukan terurai + hasil query.
type DataPeriksa struct {
	PolicyNo                 string
	EdmType                  models.JenisEndorsemen // "" = belum dipilih
	EdmTypeNew               string
	EndorsementInternalRetro string
	// TanggalEndorsement - nil = pemeriksaan tanggal belum jalan.
	TanggalEndorsement *time.Time

	// StatusEDM - GetEDMStatus_SQL CARI20 ("" = tanpa baris).
	StatusEDM string
	// StatusEDMDibaca - StatusEDM sudah diquery (RI slip melompatinya di /periksa; OpenCase tetap membutuhkannya).
	StatusEDMDibaca bool
	// SistemBaru - polis buatan NB FacIn sistem baru (DATA_JSON versi terakhir memuat IDNewBisnis, A51). StatusEDM
	// lalu bersumber kasus EDM sistem baru yang berproduksi, bukan DATA_JSON.
	SistemBaru bool
	// G1TakTerbaca - polis sistem baru, T_EDM_FACIN_KASUS tidak ada (ORA-00942): gerbang 1 belum diperiksa.
	G1TakTerbaca bool
	// EDMBelumSelesai - gerbang 2: ID kasus EDM sistem baru yang belum selesai (nol atau satu).
	EDMBelumSelesai []string
	// G2Diperiksa - gerbang 2 sudah dibaca dari T_EDM_FACIN_KASUS; false = belum diperiksa.
	G2Diperiksa bool
	// G2TakTerbaca - T_EDM_FACIN_KASUS tidak ada (ORA-00942).
	G2TakTerbaca bool
	// Klaim - GetDataClaim_SQL CARI1, nol atau satu nilai.
	Klaim []string
	// Pembayaran - SearcStatusBayarArasaps_SQL ada baris; nil = tidak terbaca (skema lain).
	Pembayaran *bool
	// Renewal - GetListRNWbyNopolis_SQL ada baris.
	Renewal bool
	// Facout - GetFacoutList_SQL ada baris (hanya dibaca bila retro).
	Facout bool
	// JenisBisnis - GetBusinessType_Sql CARI2 mentah, SEBELUM langkah 8.
	JenisBisnis string
	// JenisBisnisDariTabel - JenisBisnis diturunkan dari FACINPRODUCTION → BUSINESS → kontrak JenisBisnisFacIn
	// (polis sistem baru, A51), bukan dari DATA_JSON.
	JenisBisnisDariTabel bool
	// Klep - BrowseOpenProteksiEdm_RD per Param.Type: true = ada baris. Tidak ada di peta = kosong.
	Klep map[KlepPembatal]bool
	// KlepTakTerbaca - OPENPROTEKSI_EDM tidak ada (ORA-00942): seluruh klep kosong.
	KlepTakTerbaca bool

	// Periode - GetStartDate; nil = tanpa baris.
	Periode *PeriodeMentah
	// PengecualianTanggal - M_EDM_PENGECUALIAN_POLIS JENIS='TANGGAL'.
	PengecualianTanggal DaftarPolis
	// PengecualianBelumAda - tabel 261 tidak ada (ORA-00942).
	PengecualianBelumAda bool
}

// predikatLini - langkah 10-16, urutan korpus.
var predikatLini = []string{"IsFire", "IsAneka", "IsMBU", "IsMarineCargo", "IsPA", "IsGolfInsurance", "IsLife"}

// PredikatLiniDari - langkah 10-16 dari registry varian EDM. Hanya tujuh predikat itu yang dinilai: langkah 9
// melompat ke `TO` bila jenis bisnis kosong, jadi pemanggil tidak memanggil ini untuk jenis bisnis "".
//
// ⚠️ `[pertanyaan terbuka]` IsLife varian EDM membaca `pyWorkPage.Quotation.BusinessOldId` (L1…L16). Tidak satu pun
// rule `Endorsment Fac In` mengisi `.Quotation.BusinessOldId` kasus portal (`grep -rl BusinessOldId` → yang ada hanya
// `OfferFacIn.QuotationData.BusinessOldId`), dan layar mengisi `.PolicyNo` sebagai teks bebas - jadi properti itu
// dibaca kosong [dugaan], IsLife salah, dan polis berpola RNML ditolak "Invalid business" oleh langkah 17.
func PredikatLiniDari(jenisBisnis string, properti map[string]string) (Predikat, error) {
	k := predikat.KasusEDM{Properti: properti, JenisBisnisLama: &jenisBisnis}
	var p Predikat
	tuju := map[string]*bool{
		"IsFire": &p.IsFire, "IsAneka": &p.IsAneka, "IsMBU": &p.IsMBU, "IsMarineCargo": &p.IsMarineCargo,
		"IsPA": &p.IsPA, "IsGolfInsurance": &p.IsGolfInsurance, "IsLife": &p.IsLife,
	}
	for _, nama := range predikatLini {
		v, err := predikat.Eval(nama, k)
		if err != nil {
			return Predikat{}, fmt.Errorf("predikat %s: %w", nama, err)
		}
		*tuju[nama] = v
	}
	return p, nil
}

// SusunHasilPeriksa - jawaban pemeriksaan dari data yang sudah terbaca.
func SusunHasilPeriksa(d DataPeriksa) (models.HasilPeriksa, error) {
	h := models.HasilPeriksa{Pesan: []models.PesanPeriksa{}}

	// Langkah 7-18: jenis bisnis dan predikat lini.
	jenisBisnis := JenisBisnisUntukPredikat(d.PolicyNo, d.JenisBisnis)
	// A52: polis sistem baru tanpa BusinessType di DATA_JSON - lini belum dapat dihitung (menunggu kontrak
	// JenisBisnisFacIn); tidak dinilai, bukan "Invalid policy no!".
	bisnisBelum := d.SistemBaru && jenisBisnis == ""
	var p Predikat
	if jenisBisnis != "" {
		var err error
		if p, err = PredikatLiniDari(jenisBisnis, map[string]string{}); err != nil {
			return models.HasilPeriksa{}, err
		}
	}

	m := MasukanGerbang{
		PolicyNo: d.PolicyNo, EdmType: d.EdmType, EdmTypeNew: d.EdmTypeNew,
		EndorsementInternalRetro: d.EndorsementInternalRetro, Predikat: p, BisnisBelumDiperiksa: bisnisBelum,
		Query: HasilQueryGerbang{
			StatusEDM: d.StatusEDM, Klaim: d.Klaim, CacahRenewal: cacah(d.Renewal), CacahFacout: cacah(d.Facout),
			JenisBisnis: d.JenisBisnis, Klep: d.Klep, EDMBelumSelesai: d.EDMBelumSelesai,
		},
	}
	if d.Pembayaran != nil {
		m.Query.CacahPembayaran = cacah(*d.Pembayaran)
	}
	hg, err := PeriksaGerbangPenolakan(m)
	if err != nil {
		return models.HasilPeriksa{}, err
	}
	for _, ps := range hg.Pesan {
		h.Pesan = append(h.Pesan, models.PesanPeriksa{Medan: ps.Medan, Teks: ps.Teks})
	}

	if h.Gerbang, err = rincianGerbang(d); err != nil {
		return models.HasilPeriksa{}, err
	}
	adaTolak := false
	h.Lengkap = true
	for _, g := range h.Gerbang {
		adaTolak = adaTolak || g.Hasil == models.GerbangTolak
		if g.Hasil == models.GerbangBelumDiperiksa {
			h.Lengkap = false
		}
	}
	if adaTolak != hg.SalahGerbang {
		return models.HasilPeriksa{}, fmt.Errorf("%w: rincian tolak=%v, pesan gerbang=%v", ErrRincianTidakSepakat, adaTolak, hg.SalahGerbang)
	}

	h.Bisnis = models.BisnisPeriksa{Hasil: models.GerbangLolos, JenisBisnis: jenisBisnis}
	switch {
	case bisnisBelum:
		h.Bisnis.Hasil, h.Bisnis.Catatan = models.GerbangBelumDiperiksa, catatanBisnisSistemBaru
		h.Lengkap = false
	case jenisBisnis == "" || !(p.IsFire || p.IsAneka || p.IsMBU || p.IsMarineCargo || p.IsPA || p.IsGolfInsurance || p.IsLife):
		h.Bisnis.Hasil = models.GerbangTolak
	}
	if d.JenisBisnisDariTabel && !bisnisBelum {
		h.Bisnis.Catatan = catatanBisnisDariTabel
	}

	if d.TanggalEndorsement != nil {
		t, pesan, err := periksaTanggal(d)
		if err != nil {
			return models.HasilPeriksa{}, err
		}
		h.Tanggal = t
		if pesan != nil {
			h.Pesan = append(h.Pesan, *pesan)
		}
	}
	h.Ditolak = len(h.Pesan) > 0
	return h, nil
}

func cacah(ada bool) int {
	if ada {
		return 1
	}
	return 0
}

// rincianGerbang - hasil per gerbang, kondisinya sama dengan PeriksaGerbangPenolakan (gerbangtolak.go).
func rincianGerbang(d DataPeriksa) ([]models.GerbangPeriksa, error) {
	g := make([]models.GerbangPeriksa, 6)
	for i := range g {
		g[i].No = i + 1
	}
	set := func(no int, hasil models.HasilGerbangPeriksa, catatan string) {
		g[no-1].Hasil, g[no-1].Catatan = hasil, catatan
	}
	// Gerbang 2 berjalan juga untuk RI slip (5.6 label `jmp`), tanpa klep. Sumbernya kasus EDM SISTEM BARU
	// (T_EDM_FACIN_KASUS); kasus era Pega tidak terbaca (W-EDM-12).
	switch {
	case d.G2TakTerbaca:
		set(2, models.GerbangBelumDiperiksa, catatanKasusBelumAda)
	case !d.G2Diperiksa:
		set(2, models.GerbangBelumDiperiksa, catatanGetListEdm)
	case len(d.EDMBelumSelesai) > 0:
		set(2, models.GerbangTolak, "")
	default:
		set(2, models.GerbangLolos, catatanEDMPega)
	}
	// 5.1 / 5.10 RI slip → gerbang 1 dan 3-6 dilompati.
	if d.EdmType == models.EdmPerubahan && d.EdmTypeNew == models.EdmTypeNewRISlip {
		for _, no := range []int{1, 3, 4, 5, 6} {
			set(no, models.GerbangDilewati, "")
		}
		return g, nil
	}
	// gerbangBerklep - baris pertama prakondisi (syarat tolak) T=2 lanjut ke baris kedua; baris kedua
	// `@LengthOfPageList(OpenProteksiEDM….pxResults)>0` T=3 melewati pesan (gerbangtolak.go).
	gerbangBerklep := func(no int, tolak bool, klep KlepPembatal) {
		switch {
		case !tolak:
			set(no, models.GerbangLolos, "")
		case d.Klep[klep]:
			set(no, models.GerbangDiklep, catatanDiklep)
		case d.KlepTakTerbaca:
			set(no, models.GerbangTolak, catatanKlepTakTerbaca)
		default:
			set(no, models.GerbangTolak, "")
		}
	}
	// Gerbang 1 (5.5): CARI20==1 || CARI20==2, klep 1. Polis sistem baru: CARI20 dari kasus EDM berproduksi (A51).
	satu, err := samaKode(d.StatusEDM, "1")
	if err != nil {
		return nil, err
	}
	dua, err := samaKode(d.StatusEDM, "2")
	if err != nil {
		return nil, err
	}
	switch {
	case d.G1TakTerbaca:
		set(1, models.GerbangBelumDiperiksa, catatanKasusBelumAda)
	case d.SistemBaru && d.StatusEDM == "":
		set(1, models.GerbangLolos, catatanG1SistemBaru)
	default:
		gerbangBerklep(1, satu || dua, KlepBatal)
	}
	// Gerbang 3 (5.14): klaim terakhir CARI1 != "2" (teks), klep 2.
	gerbangBerklep(3, len(d.Klaim) > 0 && d.Klaim[0] != "2", KlepKlaim)
	// Gerbang 4 (5.18): hanya EdmType 1/2; klep 3.
	b1, err := samaKode(string(d.EdmType), string(models.EdmBatalSejakSemula))
	if err != nil {
		return nil, err
	}
	b2, err := samaKode(string(d.EdmType), string(models.EdmBatalProrata))
	if err != nil {
		return nil, err
	}
	switch {
	case !b1 && !b2:
		set(4, models.GerbangLolos, "")
	case d.Pembayaran == nil:
		set(4, models.GerbangBelumDiperiksa, catatanArasapas)
	default:
		gerbangBerklep(4, *d.Pembayaran, KlepPembayaran)
	}
	// Gerbang 5 (5.22), klep 4.
	gerbangBerklep(5, d.Renewal, KlepRenewal)
	// Gerbang 6 (5.23): hanya `.EndorsementInternalRetro==1`, tanpa klep.
	retro, err := samaKode(d.EndorsementInternalRetro, "1")
	if err != nil {
		return nil, err
	}
	switch {
	case !retro:
		set(6, models.GerbangDilewati, "")
	case !d.Facout:
		set(6, models.GerbangTolak, "")
	default:
		set(6, models.GerbangLolos, "")
	}
	return g, nil
}

// formatKabel - tanggal kabel `DD-MM-YYYY`.
const formatKabel = "02-01-2006"

// periksaTanggal - `CheckEDMPolisDate` lewat ValidasiTanggalEndorsemen.
//
// Zona: `[keputusan work owner 11-10-2026: "zona waktu samain kaya NB facin"]` - jam dinding DATE
// `FACINPRODUCTION.BEGINDATE/ENDDATE` = WIB (Asia/Jakarta), konvensi NB FacIn yang MENULIS kolom itu
// (`modul/nbfacin/backend/services/binding.go:54` `tanggalFacInOffer`: `t.In(WIB)` "DD/MM/YYYY HH24:MI:SS", dipanggil
// `selesaipolicy.go` untuk BEGINDATE/ENDDATE;
// `SaveOfferProduction_Act` "dd/MM/yyyy hh:mm:ss" Asia/Jakarta). Baris tulisan NB berjam 12:00 WIB = 05:00 GMT,
// tanggal sama. ⛔ K-046 apa adanya (validasitanggal.go): baris berjam 00:00 WIB = 17:00 GMT hari SEBELUMNYA, dan
// ValidasiTanggalEndorsemen membandingkan tanggal kalender GMT - batas mulainya mundur sehari. Mulai/Akhir keluaran =
// batas yang DIPAKAI pembanding (tanggal kalender GMT), bukan tanggal WIB tersimpan.
func periksaTanggal(d DataPeriksa) (*models.TanggalPeriksa, *models.PesanPeriksa, error) {
	t := &models.TanggalPeriksa{}
	if d.PengecualianBelumAda {
		t.Catatan = catatanPengecualianBelumAda
	}
	var periode *PeriodeProduksi
	if d.Periode != nil && d.Periode.Mulai != "" && d.Periode.Akhir != "" {
		mulai, err := uraiDateOracle(d.Periode.Mulai)
		if err != nil {
			return nil, nil, err
		}
		akhir, err := uraiDateOracle(d.Periode.Akhir)
		if err != nil {
			return nil, nil, err
		}
		periode = &PeriodeProduksi{BeginDate: mulai, EndDate: akhir}
		t.Mulai, t.Akhir = mulai.UTC().Format(formatKabel), akhir.UTC().Format(formatKabel)
	}
	hv, err := ValidasiTanggalEndorsemen(MasukanValidasiTanggal{
		PolicyNo: d.PolicyNo, TanggalEndorsement: *d.TanggalEndorsement, Periode: periode,
		Pengecualian: d.PengecualianTanggal,
	})
	switch {
	case errors.Is(err, ErrPeriodeProduksiTidakAda):
		t.Hasil = models.GerbangTolak
		return t, &models.PesanPeriksa{Medan: medanTanggalEndorsement, Teks: pesanPeriodeTidakAda}, nil
	case err != nil:
		return nil, nil, err
	case hv.Dilewati:
		t.Hasil = models.GerbangDilewati
	case hv.Ditolak:
		t.Hasil = models.GerbangTolak
		return t, &models.PesanPeriksa{Medan: hv.Medan, Teks: hv.Pesan}, nil
	default:
		t.Hasil = models.GerbangLolos
	}
	return t, nil, nil
}

// medanTanggalEndorsement - medan pesan `CheckEDMPolisDate` 5.2 (sama dengan HasilValidasiTanggal.Medan).
const medanTanggalEndorsement = "EndorsementDate"

// ErrDateOracleTakTerurai - teks TO_CHAR(…,'YYYYMMDDHH24MISS') rusak.
var ErrDateOracleTakTerurai = errors.New("endorsement: tanggal periode produksi tidak terurai")

// uraiDateOracle - jam dinding DATE sebagai WIB (lihat periksaTanggal).
func uraiDateOracle(s string) (time.Time, error) {
	w, err := time.ParseInLocation("20060102150405", s, jakarta)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %q", ErrDateOracleTakTerurai, s)
	}
	return w, nil
}
