package services

// Pendaftaran klaim Life - tiket 02.
//
// Untuk apa berkas ini: membentuk klaim baru - SATU case per PL number + mata uang (keputusan work owner
// 10-10-2026) - dari peserta yang dipilih pengguna, semua case dalam SATU transaksi. Tiga tempat ditulis sekaligus - baris work
// object, header klaim, dan peserta terpilih - ditambah baris datar warisan,
// sebab dua sistem masih berjalan berdampingan (ADR-U-0042).
//
// Dibaca sesudah: services.go dan repository/pohonklaim.go.
//
// Istilah:
//   - work object : satu baris T_WORK_CLAIM; akar tangga klaim di Pega.
//   - nomor klaim : nomor BISNIS yang dibaca manusia, berbeda dari pengenal.

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/inti/backend/penomor"
	"nusantarare/inti/backend/uang"
	"nusantarare/modul/claimlife/backend/models"
	"nusantarare/modul/claimlife/backend/repository"
)

// ErrPenomorBelumDiputuskan menandai penomoran yang belum boleh ditulis.
var ErrPenomorBelumDiputuskan = errors.New(
	"services: cara membentuk nomor klaim belum diputuskan work owner")

// Penomor membentuk NOMOR BISNIS klaim - bukan pengenal barisnya.
//
// ⛔ Kenapa ini antarmuka dan bukan fungsi: keputusan work owner o
// (26-09-2026) melarang memanggil stored procedure, sedangkan AC 2, 3, dan
// 7-11 tiket 02 masih MENUNTUT procedure itu. Sampai teks AC-nya ditulis
// ulang, menulis logikanya berarti memilih salah satu pihak diam-diam.
//
// Yang dilakukan di sini: memisahkan tempatnya, supaya seluruh jalur lain
// dapat dibangun dan diuji, dan supaya yang belum diputuskan terlihat sebagai
// satu galat terang - bukan sebagai nomor karangan yang tampak benar.
type Penomor interface {
	NomorBerikut(ctx context.Context, tx *db.Tx, kodeBisnis string, saat time.Time) (string, error)
}

// PenomorBelumDiputuskan adalah implementasi bawaan sampai butir o dijawab.
type PenomorBelumDiputuskan struct{}

// NomorBerikut selalu gagal, dengan pesan yang menyebut apa yang ditunggu.
func (PenomorBelumDiputuskan) NomorBerikut(context.Context, *db.Tx, string, time.Time) (string, error) {
	return "", fmt.Errorf("%w: butir o melarang memanggil "+
		"POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER, sedangkan AC 2, 3, dan 7-11 tiket 02 "+
		"masih menuntutnya. Teks AC itu perlu ditulis ulang work owner lebih dulu",
		ErrPenomorBelumDiputuskan)
}

// PermintaanDaftar adalah isi satu pendaftaran klaim.
//
// ⛔ BANYAK KELOMPOK sejak keputusan work owner 10-10-2026: *"saat user ceklis 2 orang dari PL A, lalu 4 orang dari
// PL B, claim membuat 2 case (CLMLF) berdasarkan PL number"*. Satu pendaftaran boleh memuat peserta dari beberapa
// premium list; server membuat satu case per PL number + mata uang (menyimpang dari Pega: satu polis per register).
//
// Dua bentuk diterima, keduanya bermakna sama:
//   - `Kelompok` terisi - bentuk baru layar Register;
//   - `Kelompok` kosong - medan tunggal NomorPremiList/NomorPolis/MataUang/Sertifikat = SATU kelompok (bentuk lama,
//     dipertahankan supaya pemanggil dan uji lama tetap sah).
//
// ⛔ Type dan kode bisnis TIDAK lagi dari klien (keputusan work owner 10-10-2026: "kolom type dan kode bisnis
// hilangkan saja"): per case dibaca dari header IDPEGA pesertanya (RingkasPolisIdPega) - Type = PolisRingkas.Type (jendela DOL),
// kode bisnis = PolisRingkas.BusinessCode (prefix nomor, sama dengan Save to RNM). Medan `type`/`kodeBisnis` di badan
// JSON lama diabaikan.
type PermintaanDaftar struct {
	NomorPremiList string
	NomorPolis     string
	MataUang       string
	// Sertifikat adalah NOMOR SERTIFIKAT peserta yang dipilih pengguna.
	//
	// ⛔ Hanya nomornya. Nilai polis - uang, tanggal valuasi, share - dibaca
	// server sendiri dari sumbernya, sebab nilai itu menentukan angka klaim
	// dan jendela DOL. Menerimanya dari badan permintaan berarti siapa pun
	// yang dapat mengirim permintaan dapat menentukannya.
	Sertifikat []string
	// Kelompok - pilihan per premium list (bentuk baru). ⛔ Pengelompokan klien TIDAK dipercaya: yang dipakai
	// hanya pasangan PL + sertifikat (untuk membaca peserta) dan NomorPolis (untuk dicocokkan). Case dibentuk ulang
	// server dari PL_NUMBER dan CURRENCY baris sumber - lihat KelompokkanPeserta.
	Kelompok []KelompokDaftar
}

// KelompokDaftar - pilihan peserta satu premium list dari klien.
type KelompokDaftar struct {
	NomorPremiList string
	// NomorPolis - DIABAIKAN sejak 10-10-2026 (dulu dicocokkan dengan POLICY_NO sumber): kunci polis = IDPEGA peserta,
	// dan POLICY_NO ceding header diturunkan dari baris sumber.
	NomorPolis string
	// MataUang - DIABAIKAN: mata uang case selalu dari CURRENCY baris sumber.
	MataUang   string
	Sertifikat []string
}

// Batas satu pendaftaran - di atasnya ditolak dengan pesan (bukan dipotong diam-diam).
const (
	// MaksKelompokDaftar - case (PL + mata uang) paling banyak sekali daftar, juga kelompok kiriman klien.
	MaksKelompokDaftar = 50
	// MaksPesertaDaftar - peserta paling banyak sekali daftar (satu query ber-index per peserta).
	MaksPesertaDaftar = 500
)

// ErrKelompokTidakSah - pilihan tidak dapat dibentuk menjadi case: peserta tidak ada / batal, nomor polis tidak
// cocok dengan PL, atau satu PL + mata uang memuat lebih dari satu polis. Handler: 422.
var ErrKelompokTidakSah = errors.New("services: pilihan peserta tidak dapat didaftarkan")

// ErrPesertaSudahDiklaim - peserta terpilih masih aktif (Outstanding / Accepted) di klaim lain - aturan in-app
// Register (keputusan work owner 10-10-2026). Handler: 409. Kalimatnya di GalatPesertaSudahDiklaim.
var ErrPesertaSudahDiklaim = errors.New("services: peserta masih aktif di klaim lain")

// GalatPesertaSudahDiklaim - SEMUA peserta yang bentrok, satu kalimat per peserta, tanpa awalan teknis.
type GalatPesertaSudahDiklaim struct {
	Bentrok []repository.PesertaAktifInApp
}

// Error - `Certificate No "X" (PL Y) is already in claim N (Outstanding|Accepted)`, dipisah "; ".
func (g *GalatPesertaSudahDiklaim) Error() string {
	baris := make([]string, len(g.Bentrok))
	for i, b := range g.Bentrok {
		nomor := strings.TrimSpace(b.NomorKlaim)
		if nomor == "" {
			nomor = b.KlaimID // klaim tanpa CLAIM_NO: pengenal CLMLF-nya
		}
		status := "Outstanding"
		if strings.TrimSpace(b.KodeStatus) == kontrak.KodeAksep {
			status = "Accepted"
		}
		baris[i] = fmt.Sprintf("Certificate No %q (PL %s) is already in claim %s (%s)", b.Sertifikat, b.NomorPremiList, nomor, status)
	}
	return strings.Join(baris, "; ")
}

// Unwrap membuat `errors.Is(err, ErrPesertaSudahDiklaim)` benar.
func (g *GalatPesertaSudahDiklaim) Unwrap() error { return ErrPesertaSudahDiklaim }

// periksaPesertaAktif - aturan in-app Register: SATU kueri per PL (kelompok ber-PL sama digabung), seluruh bentrok
// dikumpulkan. Dijalankan SEBELUM peserta dibaca dan sebelum nomor diambil.
func periksaPesertaAktif(ctx context.Context, cek func(context.Context, string, []string) ([]repository.PesertaAktifInApp, error),
	masuk []KelompokDaftar) error {
	var urutPL []string
	perPL := map[string][]string{}
	for _, k := range masuk {
		pl := strings.TrimSpace(k.NomorPremiList)
		if _, ada := perPL[pl]; !ada {
			urutPL = append(urutPL, pl)
		}
		for _, s := range k.Sertifikat {
			perPL[pl] = append(perPL[pl], strings.TrimSpace(s))
		}
	}
	var bentrok []repository.PesertaAktifInApp
	for _, pl := range urutPL {
		b, err := cek(ctx, pl, perPL[pl])
		if err != nil {
			return err
		}
		bentrok = append(bentrok, b...)
	}
	if len(bentrok) > 0 {
		return &GalatPesertaSudahDiklaim{Bentrok: bentrok}
	}
	return nil
}

// kelompokMasuk - kelompok kiriman klien; bentuk lama = satu kelompok.
func (p PermintaanDaftar) kelompokMasuk() []KelompokDaftar {
	if len(p.Kelompok) > 0 {
		return p.Kelompok
	}
	return []KelompokDaftar{{NomorPremiList: p.NomorPremiList, NomorPolis: p.NomorPolis, MataUang: p.MataUang,
		Sertifikat: p.Sertifikat}}
}

// Periksa memastikan permintaan cukup untuk membentuk klaim.
//
// Dipisah dari Daftar supaya dapat diuji tanpa Oracle, dan supaya jalur yang
// DITOLAK punya ujinya sendiri.
func (p PermintaanDaftar) Periksa() error {
	kelompok := p.kelompokMasuk()
	if len(kelompok) > MaksKelompokDaftar {
		return fmt.Errorf("%w: %d kelompok dikirim; paling banyak %d sekali daftar", galat.ErrPermintaanTidakSah,
			len(kelompok), MaksKelompokDaftar)
	}
	jumlah := 0
	// ⛔ Ganda DI DALAM permintaan: sertifikat yang sama pada PL yang sama dipilih dua kali (dalam satu kelompok
	// atau di dua kelompok) akan menjadi dua baris peserta atas satu tertanggung.
	sudah := map[string]bool{}
	for _, k := range kelompok {
		pl := strings.TrimSpace(k.NomorPremiList)
		switch {
		case pl == "":
			return fmt.Errorf("%w: nomor premium list wajib terisi", galat.ErrPermintaanTidakSah)
		case len(k.Sertifikat) == 0:
			return fmt.Errorf("%w: nol peserta dipilih pada premium list %s; klaim tanpa peserta tidak berarti apa-apa",
				galat.ErrPermintaanTidakSah, pl)
		}
		for _, s := range k.Sertifikat {
			s = strings.TrimSpace(s)
			if s == "" {
				return fmt.Errorf("%w: nomor sertifikat kosong pada premium list %s", galat.ErrPermintaanTidakSah, pl)
			}
			kunci := pl + "\x00" + s
			if sudah[kunci] {
				return fmt.Errorf("%w: sertifikat %s pada premium list %s dipilih lebih dari sekali",
					galat.ErrPermintaanTidakSah, s, pl)
			}
			sudah[kunci] = true
			jumlah++
		}
	}
	if jumlah > MaksPesertaDaftar {
		return fmt.Errorf("%w: %d peserta dipilih; paling banyak %d sekali daftar", galat.ErrPermintaanTidakSah,
			jumlah, MaksPesertaDaftar)
	}
	return nil
}

// KasusDaftar - satu case yang akan dibuat: satu PL number + mata uang, dengan peserta dari sumber.
type KasusDaftar struct {
	NomorPremiList string
	MataUang       string
	// NomorPolis - POLICY_NO ceding untuk header (arti kolom tidak berubah): peserta pertama yang terisi.
	NomorPolis string
	Peserta    []models.Peserta
	// Polis - header polis dari IDPEGA peserta case (RingkasPolisIdPega); sumber Type dan kode bisnis case.
	Polis kontrak.PolisRingkas
}

// KelompokkanPeserta membentuk case dari peserta yang dibaca ULANG dari sumber (murni, tanpa Oracle).
//
// `dibaca[i]` = peserta hasil AmbilUntukKlaim untuk `masuk[i]` (urutan sama).
//
//   - Kunci case = PL_NUMBER + CURRENCY baris SUMBER (keputusan work owner 10-10-2026); kelompok, mata uang, dan nomor
//     polis kiriman klien diabaikan (sumber menang). Urutan case = urutan kemunculan pertama.
//   - ⛔ Aturan "satu case satu POLICY_NO" dan cross-check nomorPolis klien <-> POLICY_NO sumber DIBUANG 10-10-2026:
//     kunci polis kini IDPEGA peserta (kuncipolis.go), jadi POLICY_NO ceding tidak lagi menentukan polis yang dibaca.
//     Header menyimpan POLICY_NO peserta pertama yang terisi (arti kolom tidak berubah); peserta satu PL dengan
//     POLICY_NO berbeda tidak memblokir.
func KelompokkanPeserta(masuk []KelompokDaftar, dibaca [][]models.Peserta) ([]KasusDaftar, error) {
	if len(masuk) != len(dibaca) {
		return nil, fmt.Errorf("services: %d kelompok, %d hasil baca", len(masuk), len(dibaca))
	}
	var kasus []KasusDaftar
	letak := map[string]int{}
	for i := range masuk {
		for _, ps := range dibaca[i] {
			pl, uang, pol := strings.TrimSpace(ps.NomorPremiList), strings.TrimSpace(ps.MataUang), strings.TrimSpace(ps.NomorPolis)
			kunci := pl + "\x00" + uang
			j, ada := letak[kunci]
			if !ada {
				j = len(kasus)
				letak[kunci] = j
				kasus = append(kasus, KasusDaftar{NomorPremiList: pl, MataUang: uang})
			}
			if kasus[j].NomorPolis == "" {
				kasus[j].NomorPolis = pol
			}
			kasus[j].Peserta = append(kasus[j].Peserta, ps)
		}
	}
	if len(kasus) > MaksKelompokDaftar {
		return nil, fmt.Errorf("%w: pilihan ini membentuk %d case; paling banyak %d sekali daftar",
			ErrKelompokTidakSah, len(kasus), MaksKelompokDaftar)
	}
	return kasus, nil
}

// lengkapiPolisKasus membaca header polis SETIAP case dari IDPEGA pesertanya (RingkasPolisIdPega, keputusan work owner
// 10-10-2026) dan menuntut Type dan kode bisnis - seluruh kekurangan dilaporkan sekaligus (ErrKelompokTidakSah, 422).
// Sebelum nomor diambil. Peserta tanpa IDPEGA (baris warisan) ditolak dengan namanya - tidak ditebak dari PL/NO_POLIS.
func lengkapiPolisKasus(ctx context.Context, ringkas func(context.Context, []string) (kontrak.PolisRingkas, string, error),
	kasus []KasusDaftar) error {
	var kurang []string
	for i := range kasus {
		k := &kasus[i]
		for _, ps := range k.Peserta {
			if strings.TrimSpace(ps.IdPega) == "" {
				kurang = append(kurang, fmt.Sprintf("Insured Certificate No %q (PL %s) has no IDPEGA in its policy row.",
					ps.NomorSertifikat, k.NomorPremiList))
			}
		}
		idPega := IdPegaPeserta(k.Peserta)
		if len(idPega) == 0 {
			continue
		}
		polis, id, err := ringkas(ctx, idPega)
		if errors.Is(err, kontrak.ErrPolisNomorTakDitemukan) {
			kurang = append(kurang, fmt.Sprintf("Policy data not found for this insured (%s).", id))
			continue
		}
		if err != nil {
			return err
		}
		// ⛔ Keputusan work owner 10-10-2026: "tidak, polis harus sudah resolve completed". Status = status efektif header
		// IDPEGA (PremiumList `sqlPolisRingkasIdPega`: T_WORK_POLIS.STATUS_WORK kasus NB, selain itu T_PREMIUM_LIST.STATUSS
		// yang ditulis EDM Confirm). Nilai pyWorkStatus yang sama dengan klaim: kontrak.StatusWorkSelesai.
		if status := strings.TrimSpace(polis.Status); status != kontrak.StatusWorkSelesai {
			kurang = append(kurang, fmt.Sprintf("Policy for PL %s is not Resolved-Completed yet (status: %s).", k.NomorPremiList, status))
		}
		if strings.TrimSpace(polis.Type) == "" {
			kurang = append(kurang, fmt.Sprintf("Policy data for this insured (%s, PL %s) has no Type (it decides the DOL window).", id, k.NomorPremiList))
		}
		if strings.TrimSpace(polis.BusinessCode) == "" {
			kurang = append(kurang, fmt.Sprintf("Policy data for this insured (%s, PL %s) has no business code (claim number prefix).", id, k.NomorPremiList))
		}
		k.Polis = polis
	}
	if len(kurang) > 0 {
		return fmt.Errorf("%w: %s", ErrKelompokTidakSah, strings.Join(kurang, " "))
	}
	return nil
}

// Pendaftaran adalah layanan pendaftaran klaim baru.
type Pendaftaran struct {
	svc *Service
	// jahit - jahitan uji tanpa Oracle (nil = Oracle lewat svc).
	jahit *jahitanDaftar
}

// jahitanDaftar - empat pintu Oracle pendaftaran, diganti uji supaya aturan "semua atau tidak sama sekali" dapat
// dibuktikan tanpa basis data.
type jahitanDaftar struct {
	bacaPeserta func(ctx context.Context, nomorPremiList string, sertifikat []string) ([]models.Peserta, error)
	// pesertaAktif - aturan in-app Register (repository.PesertaAktifDiKlaimLain).
	pesertaAktif func(ctx context.Context, nomorPremiList string, sertifikat []string) ([]repository.PesertaAktifInApp, error)
	transaksi    func(ctx context.Context, fn func(tx *db.Tx) error) error
	pengenal     func(ctx context.Context, tx *db.Tx) (string, error)
	simpan       func(ctx context.Context, tx *db.Tx, p models.PohonKlaim) error
	// ringkas - header polis dari IDPEGA peserta case (RingkasPolisIdPega).
	ringkas func(ctx context.Context, idPega []string) (kontrak.PolisRingkas, string, error)
}

func (p *Pendaftaran) pintu() jahitanDaftar {
	if p.jahit != nil {
		return *p.jahit
	}
	return jahitanDaftar{
		bacaPeserta: func(ctx context.Context, pl string, s []string) ([]models.Peserta, error) {
			return repository.NewPesertaPolis(p.svc.DB()).AmbilUntukKlaim(ctx, pl, s)
		},
		pesertaAktif: func(ctx context.Context, pl string, s []string) ([]repository.PesertaAktifInApp, error) {
			return repository.NewKlaimLife(p.svc.DB()).PesertaAktifDiKlaimLain(ctx, pl, s)
		},
		transaksi: p.svc.DalamTransaksi,
		pengenal: func(ctx context.Context, tx *db.Tx) (string, error) {
			return repository.NewPohonKlaim(p.svc.DB()).PengenalWorkBerikut(ctx, tx, repository.AwalanKlaim)
		},
		simpan: func(ctx context.Context, tx *db.Tx, pohon models.PohonKlaim) error {
			return repository.NewPohonKlaim(p.svc.DB()).Simpan(ctx, tx, pohon)
		},
		ringkas: func(ctx context.Context, idPega []string) (kontrak.PolisRingkas, string, error) {
			return RingkasPolisIdPega(ctx, p.svc.PembacaPolis(), idPega)
		},
	}
}

// Pendaftaran menyusun layanan itu dengan penomor bawaan.
func (s *Service) Pendaftaran() *Pendaftaran {
	// RALAT 10-10-2026: tanpa penomor - nomor klaim lahir di Save to RNM.
	return &Pendaftaran{svc: s}
}

// Daftar membentuk SATU case per PL number + mata uang dari pilihan peserta - semua atau tidak sama sekali.
//
// Urutannya disengaja:
//
//  1. pelaku, peran, dan permintaan diperiksa;
//  2. peserta SETIAP kelompok dibaca ulang dari sumber, dibulatkan, baris adjustment pertamanya dilahirkan, lalu
//     dikelompokkan, dicek aktif di klaim lain, dan polisnya dibaca (Type, Resolved-Completed);
//  3. SATU TRANSAKSI untuk SEMUA case: pengenal, T_WORK_CLAIM, T_GENERAL_CLAIM, T_CLAIMLF_*, baris datar warisan.
//     Satu case gagal ditulis = seluruh transaksi dibatalkan, nol case tersimpan.
//
// ⛔ RALAT 10-10-2026 (keputusan work owner "claim no tidak generate saat register claim, tapi nanti setelah save to
// rnm"): Register TIDAK lagi mengambil nomor klaim - CLAIM_NO kosong sampai Save to RNM (`simpanrnm.go` langkah
// 16-20, transaksi penomoran tersendiri). Akibatnya "transaksi pertama nomor" hilang: semua-atau-tidak kini MURNI,
// tanpa nomor hangus. `SEQ_WORK_CLAIM` (pengenal CLMLF-) tetap non-transaksional - pengenal yang diambil transaksi
// gagal tidak kembali.
//
// ⚠️ Pendaftaran TIDAK menulis jejak (T_CLAIMLF_JEJAK) dan TIDAK menjalankan pemeriksa klaim ganda warisan
// (`gandawarisan.go`, langkah 11 `Save to RNM`, tidak berubah). Ganda di dalam permintaan ditolak Periksa; peserta
// yang masih aktif di klaim APLIKASI lain ditolak periksaPesertaAktif (aturan in-app 10-10-2026, 409).
func (p *Pendaftaran) Daftar(ctx context.Context, pelaku inti.Pelaku, minta PermintaanDaftar) (
	[]models.PohonKlaim, error) {

	// ⛔ FAIL-CLOSED atas pelaku anonim. Tanpa ini, permintaan tanpa pelaku -
	// keadaan BAWAAN saat stub mati, yaitu keadaan produksi - tetap menulis
	// klaim, hanya dengan CREATE_OP_NAME kosong. Itu bukan pagar yang gagal
	// tertutup melainkan pagar yang gagal ANONIM: jejaknya hilang, dan tidak
	// ada yang tahu siapa mendaftarkan klaim itu.
	//
	// ⛔ Identitas, lalu PERAN. Gerbang perannya dipasang tiket 07 - komentar
	// di sini sebelumnya menandainya `[terbuka - tiket 07]`, dan tiket itu
	// sekarang sudah lewat. `[terverifikasi]` `Flow/Register_Flow.xml`:
	// `Assignment2` (Input Register) dipegang `ReasLifeAdmin`.
	//
	// ⚠️ Pendaftaran MENULIS status Outstanding lewat TandaiOutstanding, jadi
	// ia salah satu jalur pengubah status - dan jalur pengubah status yang
	// tidak bergerbang adalah lubang, bukan kelonggaran.
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return nil, fmt.Errorf("%w: klaim mencatat siapa yang membuatnya", err)
	}
	if err := inti.WajibPeran(pelaku, PeranInputRegister); err != nil {
		return nil, err
	}
	if err := minta.Periksa(); err != nil {
		return nil, err
	}
	if p.jahit == nil && !p.svc.PunyaDatabase() {
		return nil, db.ErrTanpaOracle
	}
	pintu := p.pintu()

	// Langkah 2 - baca ulang, siapkan, kelompokkan; nol nomor diambil sebelum semuanya sah.
	masuk := minta.kelompokMasuk()
	// ⛔ Aturan in-app (keputusan work owner 10-10-2026): peserta yang masih Outstanding / Accepted di klaim lain
	// ditolak; Rejected dan baris tercabut boleh didaftarkan lagi. Seluruh bentrok dilaporkan sekaligus.
	if err := periksaPesertaAktif(ctx, pintu.pesertaAktif, masuk); err != nil {
		return nil, err
	}
	dibaca := make([][]models.Peserta, len(masuk))
	for i, k := range masuk {
		sertifikat := make([]string, len(k.Sertifikat))
		for j, s := range k.Sertifikat {
			sertifikat[j] = strings.TrimSpace(s)
		}
		// ⛔ Peserta dibaca ULANG dari sumbernya, bukan diterima dari klien.
		peserta, err := pintu.bacaPeserta(ctx, strings.TrimSpace(k.NomorPremiList), sertifikat)
		if err != nil {
			if errors.Is(err, repository.ErrPesertaKlaimTidakAda) {
				return nil, fmt.Errorf("%w: %w", ErrKelompokTidakSah, err)
			}
			return nil, err
		}
		// ⛔ BUTIR bp (GILIRAN-14): baris adjustment pertama LAHIR DI SINI,
		// seperti `SavePesertaClaim` 7.8 - bukan kelak lewat `Add`. Di
		// transaksi yang sama dengan pesertanya. ⛔ RALAT 10-10-2026: `Simpan` TIDAK
		// lagi menulis baris datar `OS_AKSEPTASI_KLAIM_LIFE` - baris itu lahir di
		// Save to RNM bersama nomor klaim (keputusan work owner).
		// ⛔ OQ-N10 (GILIRAN-15): nilai peserta dibulatkan seperti 7.7 LEBIH
		// DULU - baris 7.8 menyalin nilai dari sumber yang sama dengan rumus
		// yang sama, jadi urutannya tidak mengubah angka.
		if err := BulatkanPesertaPendaftaran(peserta); err != nil {
			return nil, err
		}
		if err := LahirkanBarisPendaftaran(peserta); err != nil {
			return nil, err
		}
		dibaca[i] = peserta
	}
	kasus, err := KelompokkanPeserta(masuk, dibaca)
	if err != nil {
		return nil, err
	}
	// Type dan kode bisnis per case dari header IDPEGA pesertanya (10-10-2026) - masih sebelum nomor diambil.
	if err := lengkapiPolisKasus(ctx, pintu.ringkas, kasus); err != nil {
		return nil, err
	}

	// SATU jam untuk seluruh pendaftaran (TGL_CREATE = TGL_UPDATE seluruh case).
	saat := time.Now()

	// SATU TRANSAKSI - SEMUA case, semua atau tidak sama sekali (tanpa nomor klaim; Save to RNM yang menomori).
	hasil := make([]models.PohonKlaim, len(kasus))
	err = pintu.transaksi(ctx, func(tx *db.Tx) error {
		for i, k := range kasus {
			pengenal, err := pintu.pengenal(ctx, tx)
			if err != nil {
				return err
			}
			pohon := models.PohonKlaim{
				Work: models.WorkClaim{
					// `[keputusan work owner 26-09-2026, butir ae1]` CASEID =
					// pengenal work object - sejak migrasi 023 (keputusan work
					// owner 01-10-2026) kolom CASE_ID dibuang dan ID-nya sendiri
					// yang dipakai: baris datar warisan klaim baru tetap
					// dikelompokkan per CASEID = ID.
					ID:           pengenal,
					Lini:         inti.LiniLife,
					CreateOpName: pelaku.AkunID,
					// ⛔ BUTIR au: `CREATE_OP` adalah padanan `pxCreateOperator`,
					// dan ITULAH yang worklist Pega rutekan (`Register_Flow.xml`
					// 1351 `ToWorklist`, 1508 `Current operator`). Tanpa kolom ini
					// terisi, kotak masuk Admin tidak dapat menyaring "kasus
					// milik saya" sama sekali.
					CreateOp: pelaku.AkunID,
					// ⛔ BUTIR at: kasus BARU berada di Outstanding Claim, bukan
					// Input Register. `Assignment2` (Input Register) SELESAI pada
					// saat borang register disimpan - yang tersimpan di sini
					// adalah hasilnya. Kasus kembali ke Input Register hanya lewat
					// `Send Back to Register` (`InputOSClaimLife.xml:21404`).
					Tahap: models.TahapOutstanding.String(),
					// ⛔ Peran pemegangnya tetap ditulis: `TAHAP` mengatakan di anak
					// tangga mana, `POSITION` (`pyPosition`) mengatakan peran siapa.
					Position: models.PeranAdminLife,
					// BUTIR au: waktu LAHIR, terpisah dari waktu ubah.
					TglCreate: saat,
					TglUpdate: saat,
				},
				Klaim: models.Klaim{
					ID: pengenal,
					// POLICY_NO ceding dari baris sumber (arti kolom tidak berubah, 10-10-2026).
					NomorPolis: k.NomorPolis,
					// TYPE di header klaim sejak migrasi 023 - dari polis PL case (keputusan work owner 10-10-2026).
					Type: strings.TrimSpace(k.Polis.Type),
					// Mata uang dari CURRENCY baris sumber (butir z1), bukan dari klien.
					ClaimRetro: uang.Money{Currency: k.MataUang},
					Peserta:    k.Peserta,
					// Bawaan tab Register Claim (keputusan work owner 11-10-2026: "by default datanya terisi Non-Catastrophe
					// dan Non-KPR").
					StsKatastrofe: NonKatastrofe,
					IsKPR:         NonKPR,
				},
			}
			if err := pintu.simpan(ctx, tx, pohon); err != nil {
				return fmt.Errorf("case premium list %s (%s): %w", k.NomorPremiList, k.MataUang, err)
			}
			hasil[i] = pohon
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return hasil, nil
}

// Kunci penghitung nomor klaim Life.
//
// `[data DBA]` `GENERATE_SEQUENCE_NUMBER` berkunci `(CLASS, JENIS, TAHUN)`.
// `CLASS` dan `JENIS` menentukan seri; `TAHUN` memisahkan tahun buku.
const (
	ClassNomorKlaimLife = "CLAIMLIFE"
	JenisNomorKlaimLife = "K"
)

// ⛔ KONSTANTA `AwalanNomorKlaim` DIBUANG, ronde 27-09-2026.
//
// Ia bernilai `"RNML-"` - benar untuk lingkungan yang kebetulan dipakai
// saat kode ditulis, dan AC tiket 02 melarangnya dengan kalimat yang tidak
// dapat disalahartikan: *"Prefix diperoleh lewat lookup ke
// `POOLDATA.KODE_PRODUKSI` (`TYPE='LIFE'`), tidak ditanam sebagai konstanta
// di kode."* Awalan itu MILIK basis data.
//
// `[terverifikasi]` `Claim Life/RDBList/GetKodeProdLife_SQL.xml` baris 85.
// Bentuk `<awalan>K<kode bisnis>.MM.YYYY.<5 digit>` sejajar dengan nomor
// akseptasi `RNML-A…`/`RNML-AR…` (`Generate_NoAccept_Life.xml` baris 85);
// huruf tengahnya yang membedakan seri - `K` klaim, `A` akseptasi.
//
// ⚠️ Nomor AKSEPTASI masih merakit awalannya sendiri (`AwalanAkseptasiGross`
// / `AwalanAkseptasiRetro` di akseptasi.go). Itu dinyatakan, bukan
// dilewatkan: awalan akseptasi TIDAK ada di `KODE_PRODUKSI` - tabel itu
// hanya punya `RNML-`, sedangkan akseptasi menuntut `RNML-A` dan `RNML-AR`.
// Menyatukannya menuntut keputusan work owner tentang dari mana huruf
// `A`/`AR` datang, dan itu `[terbuka]`.

// RakitNomorKlaim menyusun nomor bisnis klaim.
//
// ⚠️ LIMA digit seperti `LPAD(v_seq,5,'0')`, dan urut yang sudah lebih panjang
// TIDAK dipotong - memotongnya menerbitkan nomor yang bertabrakan.
//
// ⚠️ `awalan` DITERIMA, tidak diambil sendiri: fungsi ini murni supaya
// bentuk nomornya dapat diuji tanpa Oracle. Yang mengambilnya penomornya.
func RakitNomorKlaim(awalan, kodeBisnis, mmYYYY string, urut int) string {
	u := strconv.Itoa(urut)
	if len(u) < 5 {
		u = strings.Repeat("0", 5-len(u)) + u
	}
	return awalan + JenisNomorKlaimLife + kodeBisnis + "." + mmYYYY + "." + u
}

// penomorCounter menulis ulang `PROC_GENERATE_SEQUENCE_NUMBER` di Go.
//
// ⛔ `[keputusan work owner]` butir **o1**: procedure tidak dipanggil. Yang
// tetap di Oracle hanya `SELECT … FOR UPDATE`, sebab kunci baris memang milik
// basis data.
type penomorCounter struct{ penghitung *penomor.Penomor }

// PenomorCounterOracle menyusun penomor yang memakai penghitung Oracle.
func PenomorCounterOracle(svc *Service) Penomor {
	return penomorCounter{penghitung: penomor.NewPenomor(svc.DB())}
}

// NomorBerikut menerbitkan satu nomor klaim baru.
//
// Urutannya: awalan → hari tutup buku → periode → kunci baris → naikkan →
// rakit. Empat langkah tengahnya persis procedure-nya; awalan di depan
// sebab nomor tanpa awalan bukan nomor.
func (p penomorCounter) NomorBerikut(ctx context.Context, tx *db.Tx,
	kodeBisnis string, saat time.Time) (string, error) {

	if strings.TrimSpace(kodeBisnis) == "" {
		return "", fmt.Errorf("%w: kode bisnis kosong; nomor klaim memuatnya",
			galat.ErrPermintaanTidakSah)
	}
	// ⛔ Awalannya di-LOOKUP, bukan konstanta - lihat catatan di atas.
	awalan, err := p.penghitung.AwalanProduksi(ctx, tx, penomor.TipeKodeProduksiLife)
	if err != nil {
		return "", err
	}
	hariClosing, err := p.penghitung.HariClosing(ctx, tx)
	if err != nil {
		return "", err
	}
	periode, err := penomor.HitungPeriodeNomor(saat, hariClosing)
	if err != nil {
		return "", err
	}
	urut, err := p.penghitung.UrutNomorBerikut(ctx, tx,
		ClassNomorKlaimLife, JenisNomorKlaimLife, periode, saat)
	if err != nil {
		return "", err
	}
	return RakitNomorKlaim(awalan, kodeBisnis, periode.MMYYYY, urut), nil
}
