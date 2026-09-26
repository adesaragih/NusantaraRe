package services

// Mesin status per baris adjustment - tiket 04.
//
// Untuk apa berkas ini: yang diputuskan adalah BARIS, bukan klaim
// (ADR-U-0011). Berkas ini memindahkan satu baris dari Outstanding ke
// keputusan akhirnya, dan menentukan baris mana yang dicerminkan ke atas.
//
// Dibaca sesudah: adjustment.go.
//
// ⛔ Fungsi aturannya MURNI - Transisi, TandaiOutstandingKlaim, dan
// BarisTerakhir tidak menyentuh basis data dan tidak membaca jam. Yang TIDAK
// murni hanya layanan Status.Ubah di bagian bawah: ia membaca Oracle dan
// membuka transaksi, dan jam diserahkan pemanggilnya.
//
// Istilah:
//   - final      : Aksep atau Ditolak. Sekali final, tidak berubah lagi.
//   - pencerminan: menyalin status baris ke peserta dan ke header klaim.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
)

var (
	// ErrBarisSudahFinal - baris yang sudah diputus tidak berubah lagi.
	//
	// `[terverifikasi]` sensus OQ-061: nol rule di `Claim Life` maupun
	// `Komite Claim Life` menulis "0" sesudah "1" atau "2".
	ErrBarisSudahFinal = errors.New("services: baris sudah final dan tidak dapat berubah")
	// ErrTransisiTidakSah - asal atau tujuan transisi tidak diizinkan.
	ErrTransisiTidakSah = errors.New("services: transisi status baris tidak sah")
	// ErrJejakBelumDiputuskan - tempat jejak audit belum ada.
	ErrJejakBelumDiputuskan = errors.New("services: tempat jejak audit belum diputuskan")
)

// CatatanJejak adalah satu baris jejak audit sebuah transisi.
type CatatanJejak struct {
	AdjustmentID string
	Dari         string
	Ke           string
	AkunID       string
	Waktu        time.Time
}

// Jejak merekam SIAPA dan KAPAN untuk setiap transisi status.
//
// ⛔ Kenapa ini antarmuka dan bukan penulisan langsung: ADR-U-0007 menuntut
// setiap transisi terekam, sedangkan TABELNYA belum ada - keputusan membuatnya
// masih `[USULAN]` butir am, dan brief melarang menebaknya. Memisahkannya
// membuat yang belum ada terlihat sebagai satu galat terang, bukan sebagai
// transisi yang diam-diam tak tercatat.
//
// Pola yang sama dengan `Penomor` pada tiket 02.
type Jejak interface {
	Rekam(ctx context.Context, tx *repository.Tx, c CatatanJejak) error
}

// JejakBelumDiputuskan adalah implementasi bawaan; ia selalu gagal.
type JejakBelumDiputuskan struct{}

// Rekam selalu gagal, dengan pesan yang menyebut apa yang ditunggu.
func (JejakBelumDiputuskan) Rekam(context.Context, *repository.Tx, CatatanJejak) error {
	return fmt.Errorf("%w: tabel jejak audit (butir am) belum disahkan work owner, "+
		"sedangkan ADR-U-0007 menuntut setiap transisi merekam siapa dan kapan",
		ErrJejakBelumDiputuskan)
}

// Transisi memindahkan satu baris dari Outstanding ke keputusan akhirnya.
//
// Hanya DARI Outstanding, dan hanya KE Aksep atau Ditolak:
//
//	""            -> apa pun : tidak sah, baris belum pernah disimpan
//	Outstanding   -> Aksep   : sah  (`SaveAdjustment_Act` menulis 1; juga
//	                                 `KomitePostAdjustment` di modul Komite)
//	Outstanding   -> Ditolak : sah  (`RejectOSClaimLife_Act` menulis 2)
//	Aksep/Ditolak -> apa pun : ErrBarisSudahFinal
//
// ⛔ Kode "4" tidak pernah menjadi tujuan. Ia ada di data warisan, artinya
// belum diputuskan work owner, dan `models` sengaja tidak menamainya.
//
// ⭐ Aksep menstempel ACCEPTATION_DATE. `[terverifikasi]` `SaveAdjustment_Act`
// (pecahan 1833-1899) menulis KETIGANYA dalam satu Property-Set:
// `.ACCEPTEDNO`, `.STS_REJECT = 1`, dan `.ACCEPTATION_DATE =
// @CurrentDateTime()`. Itulah tempat tanggal akseptasi lahir - bukan saat
// insert (tiket 03 AC 43), melainkan saat baris benar-benar diaksep.
//
// Jam diserahkan pemanggil, tidak dibaca di sini: fungsi yang membaca jamnya
// sendiri tidak dapat diuji tanpa menunggu waktu berlalu.
//
// Baris dikembalikan sebagai NILAI BARU, bukan diubah di tempat: baris yang
// gagal transisi tidak boleh tertinggal setengah berubah.
func Transisi(b models.BarisAdjustment, ke models.StatusBaris,
	saat time.Time) (models.BarisAdjustment, error) {
	switch b.Status() {
	case models.StatusAksep, models.StatusDitolak:
		return b, fmt.Errorf("%w: baris %s berstatus %q", ErrBarisSudahFinal, b.ID, b.KodeStatus)
	case models.StatusOutstanding:
		// satu-satunya asal yang sah
	default:
		return b, fmt.Errorf("%w: baris %s belum disimpan ke Outstanding (kode %q)",
			ErrTransisiTidakSah, b.ID, b.KodeStatus)
	}

	var kode string
	switch ke {
	case models.StatusAksep:
		kode = models.KodeAksep
	case models.StatusDitolak:
		kode = models.KodeDitolak
	default:
		return b, fmt.Errorf("%w: tujuan %v bukan keputusan akhir", ErrTransisiTidakSah, ke)
	}
	b.KodeStatus = kode
	if ke == models.StatusAksep {
		b.TanggalAkseptasi = saat
	}
	return b, nil
}

// TandaiOutstandingKlaim menuliskan Outstanding pada baris yang belum
// berstatus di dalam satu klaim.
//
// ⛔ SATU tempat aturannya. `TandaiOutstanding` (tiket 03) mendelegasi ke sini
// alih-alih menyalin isinya - ronde pertama tiket 04 menyalinnya, dan dua
// salinan aturan kefinalan berarti dua kesempatan untuk berbeda.
func TandaiOutstandingKlaim(k *models.Klaim) int {
	if k == nil {
		return 0
	}
	n := 0
	for i := range k.Peserta {
		for j := range k.Peserta[i].Baris {
			b := &k.Peserta[i].Baris[j]
			if strings.TrimSpace(b.KodeStatus) != "" {
				continue
			}
			b.KodeStatus = models.KodeOutstanding
			n++
		}
	}
	return n
}

// BarisTerakhir memilih baris yang dicerminkan ke peserta dan ke header.
//
// `[terverifikasi]` `Claim Life/Activity/serviceInsertArasapasClaimLife_act.xml`
// langkah 1.1.1 (berkas pecahan baris 371-418): di dalam putaran BERSARANG -
// peserta (langkah 1) lalu `.AdjustmentList` (langkah 1.1) - ia menyetel
// `pyWorkPage.ClaimData.ACCEPTEDNO = .ACCEPTEDNO` dan
// `pyWorkPage.ClaimData.STS_REJECT = .STS_REJECT` TANPA precondition dan tanpa
// henti. Nilai yang bertahan karena itu milik baris yang diulang TERAKHIR.
//
// ⛔ Pencerminan ini bukan sumber kebenaran. Status klaim yang DILAPORKAN
// tetap `Klaim.StatusTurunan()`; kolom cermin ada untuk pembaca hilir yang
// masih membaca tabel datar (ADR-U-0042).
//
// ⚠️ "Terakhir" di sini berarti terakhir menurut URUTAN BACA
// (`AmbilBaris`: `ORDER BY a.PREMIUM_LIST_DETAIL_ID, a.ID`), yaitu peserta
// terakhir lalu baris ber-ID terbesar. Itu padanan terdekat dari "yang diulang
// terakhir" di Pega; hubungannya dengan "baris yang paling akhir diputus"
// tidak dinyatakan di mana pun. `[terbuka - work owner]`
//
// Penunjuk, bukan nilai: pemanggil memang perlu membaca baris yang sudah
// diubahnya. Parameter pun penunjuk supaya kemampuan menulis itu TERLIHAT -
// bukan menyelinap lewat array yang kebetulan dibagi.
//
// Mengembalikan nil bila klaim belum punya satu pun baris.
func BarisTerakhir(k *models.Klaim) *models.BarisAdjustment {
	if k == nil {
		return nil
	}
	var akhir *models.BarisAdjustment
	for i := range k.Peserta {
		for j := range k.Peserta[i].Baris {
			akhir = &k.Peserta[i].Baris[j]
		}
	}
	return akhir
}

// Status adalah layanan transisi status baris adjustment.
type Status struct {
	svc   *Service
	jejak Jejak
}

// Status menyusun layanan itu dengan jejak bawaan yang gagal terang.
func (s *Service) Status() *Status {
	return &Status{svc: s, jejak: JejakBelumDiputuskan{}}
}

// DenganJejak mengganti perekamnya - dipakai test, dan kelak oleh tiket 09.
func (st *Status) DenganJejak(j Jejak) *Status {
	return &Status{svc: st.svc, jejak: j}
}

// Ubah memindahkan satu baris ke keputusan akhirnya dan mencerminkannya.
//
// Satu transaksi menutup tiga tulisan: baris adjustment, peserta pemiliknya,
// dan header klaim. Ketiganya harus jadi bersama - pencerminan setengah jadi
// membuat pembaca hilir melihat dua jawaban untuk satu pertanyaan.
//
// ⚠️ Sumber nilai header BUKAN baris yang berubah melainkan baris TERAKHIR
// klaim sesudah perubahan itu - lihat BarisTerakhir.
func (st *Status) Ubah(ctx context.Context, pelaku Pelaku,
	klaimID, pesertaID, adjID string, ke models.StatusBaris, saat time.Time) error {

	if strings.TrimSpace(pelaku.AkunID) == "" {
		return fmt.Errorf("%w: perubahan status tanpa identitas pelaku ditolak",
			ErrTanpaWewenang)
	}
	if strings.TrimSpace(klaimID) == "" || strings.TrimSpace(pesertaID) == "" ||
		strings.TrimSpace(adjID) == "" {
		return fmt.Errorf("%w: pengenal klaim, peserta, dan baris wajib diisi",
			ErrPermintaanTidakSah)
	}
	if !st.svc.PunyaDatabase() {
		return repository.ErrTanpaOracle
	}

	baca := repository.NewKlaimLife(st.svc.db)
	peserta, err := baca.AmbilPeserta(ctx, klaimID)
	if err != nil {
		return err
	}
	perBaris, err := baca.AmbilBaris(ctx, klaimID)
	if err != nil {
		return err
	}
	for i := range peserta {
		peserta[i].Baris = perBaris[peserta[i].ID]
	}
	klaim := models.Klaim{ID: klaimID, Peserta: peserta}

	sasaran := cariBaris(klaim, pesertaID, adjID)
	if sasaran == nil {
		return fmt.Errorf("%w: baris %q bukan milik peserta %q pada klaim %q",
			ErrPermintaanTidakSah, adjID, pesertaID, klaimID)
	}
	sasaranKodeLama := sasaran.KodeStatus
	baru, err := Transisi(*sasaran, ke, saat)
	if err != nil {
		return err
	}
	*sasaran = baru

	akhir := BarisTerakhir(&klaim)
	if akhir == nil {
		return fmt.Errorf("%w: klaim %q tanpa baris", ErrPermintaanTidakSah, klaimID)
	}
	return st.svc.DalamTransaksi(ctx, func(tx *repository.Tx) error {
		// ⛔ Kode LAMA ikut dikirim sebagai syarat WHERE. Baris dibaca di luar
		// transaksi, jadi ia dapat berubah di antara baca dan tulis; tanpa
		// syarat itu, kefinalan hanya berlaku di dalam proses ini dan dua
		// permintaan serentak dapat sama-sama menang.
		if err := baca.PerbaruiStatusBaris(ctx, tx, pesertaID, adjID,
			sasaranKodeLama, baru.KodeStatus, baru.NomorAkseptasi,
			baru.TanggalAkseptasi); err != nil {
			return err
		}
		if err := baca.CerminkanHeader(ctx, tx, klaimID,
			akhir.KodeStatus, akhir.NomorAkseptasi); err != nil {
			return err
		}
		// ⛔ Jejak direkam DI DALAM transaksi yang sama. Jejak yang ditulis
		// terpisah dapat hilang sendirian, dan transisi tanpa jejak persis
		// yang ADR-U-0007 larang.
		return st.jejak.Rekam(ctx, tx, CatatanJejak{
			AdjustmentID: adjID,
			Dari:         sasaranKodeLama,
			Ke:           baru.KodeStatus,
			AkunID:       pelaku.AkunID,
			Waktu:        saat,
		})
	})
}

// cariBaris menunjuk baris milik peserta tertentu di dalam klaim.
func cariBaris(k models.Klaim, pesertaID, adjID string) *models.BarisAdjustment {
	for i := range k.Peserta {
		if k.Peserta[i].ID != pesertaID {
			continue
		}
		for j := range k.Peserta[i].Baris {
			if k.Peserta[i].Baris[j].ID == adjID {
				return &k.Peserta[i].Baris[j]
			}
		}
	}
	return nil
}
