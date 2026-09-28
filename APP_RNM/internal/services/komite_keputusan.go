package services

// Keputusan satu tingkat Komite - tiket 02 Komite Claim Life.
//
// Untuk apa berkas ini: tombol `Submit` layar `ShowTransfer` (b34722 →
// `finishAssignment` → activity flow action `KomitePostAdjustment`). Satu
// anggota BERJALAN memutuskan Setuju/Tolak; tangganya naik satu tingkat atau
// berhenti. Aturannya murni di `models/komite_tangga.go`.
//
// ⛔ TINGKAT AKHIR DIGERBANG, dengan sengaja. `KomitePostAdjustment` langkah 4
// (akseptasi: nomor + rekam, tiket 04a/04b) dan langkah 5 (tolak ke baris
// klaim, tiket 05) berjalan DI DALAM activity yang sama dengan pencatatan
// tingkatnya. Keduanya belum dibangun; `PenyelesaiAkhirBelumAda` karena itu
// MENOLAK keputusan tingkat akhir, dan transaksinya batal utuh - supaya tidak
// ada kasus yang "selesai disetujui" tanpa nomor akseptasi, atau "ditolak"
// tanpa baris klaimnya tahu.
//
// Dibaca sesudah: models/komite_tangga.go, repository/komite_keputusan.go.

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
)

var (
	// ErrTanggaKomiteBerhenti - kasus tidak lagi di `Assignment1`.
	ErrTanggaKomiteBerhenti = errors.New(
		"services: tangga komite kasus ini sudah berhenti; tidak ada yang menunggu keputusan")
	// ErrPenyelesaianAkhirBelumAda - langkah 4/5 belum dibangun.
	ErrPenyelesaianAkhirBelumAda = errors.New(
		"services: keputusan tingkat akhir komite menuntut akseptasi (tiket 04a/04b) " +
			"atau penolakan ke baris klaim (tiket 05), yang belum tersedia")
	// ErrKeputusanKomiteTidakDikenal - nilai di luar {1, 2}.
	ErrKeputusanKomiteTidakDikenal = models.ErrKeputusanKomiteTidakDikenal
	// ErrKeputusanKomiteBersamaan - tangga berubah sejak dibaca.
	ErrKeputusanKomiteBersamaan = repository.ErrKeputusanKomiteBersamaan
)

// PenyelesaiAkhirKomite menjalankan langkah 4 dan 5 `KomitePostAdjustment`,
// di dalam transaksi keputusannya.
type PenyelesaiAkhirKomite interface {
	// Akseptasi - langkah 4 "Approve Last Komite" (tiket 04a/04b).
	Akseptasi(ctx context.Context, tx *repository.Tx, kasus repository.KasusKomite,
		pelaku Pelaku, saat time.Time) error
	// Tolak - langkah 5 "Reject" (tiket 05).
	Tolak(ctx context.Context, tx *repository.Tx, kasus repository.KasusKomite,
		pelaku Pelaku, saat time.Time) error
}

// PenyelesaiAkhirBelumAda gagal terang untuk kedua langkah.
type PenyelesaiAkhirBelumAda struct{}

// Akseptasi selalu gagal.
func (PenyelesaiAkhirBelumAda) Akseptasi(context.Context, *repository.Tx,
	repository.KasusKomite, Pelaku, time.Time) error {
	return ErrPenyelesaianAkhirBelumAda
}

// Tolak selalu gagal.
func (PenyelesaiAkhirBelumAda) Tolak(context.Context, *repository.Tx,
	repository.KasusKomite, Pelaku, time.Time) error {
	return ErrPenyelesaianAkhirBelumAda
}

// HasilKeputusanKomite adalah jawaban `Putuskan`.
type HasilKeputusanKomite struct {
	TingkatDiputus int    `json:"tingkatDiputus"`
	Keputusan      string `json:"keputusan"`
	KataKeputusan  string `json:"kataKeputusan"`
	Berlanjut      bool   `json:"berlanjut"`
	TingkatBerikut int    `json:"tingkatBerikut"`
	AkseptasiAkhir bool   `json:"akseptasiAkhir"`
	TolakAkhir     bool   `json:"tolakAkhir"`
}

// KeputusanKomite melayani keputusan satu tingkat.
type KeputusanKomite struct {
	svc   *Service
	jejak Jejak
	akhir PenyelesaiAkhirKomite
}

// KeputusanKomite menyusunnya dengan jejak dan penyelesai akhir yang gagal terang.
func (s *Service) KeputusanKomite() *KeputusanKomite {
	return &KeputusanKomite{svc: s, jejak: JejakBelumDiputuskan{}, akhir: PenyelesaiAkhirBelumAda{}}
}

// DenganJejak mengganti perekamnya.
func (k *KeputusanKomite) DenganJejak(j Jejak) *KeputusanKomite {
	salin := *k
	salin.jejak = j
	return &salin
}

// DenganPenyelesaiAkhir mengganti penyelesai langkah 4/5.
func (k *KeputusanKomite) DenganPenyelesaiAkhir(p PenyelesaiAkhirKomite) *KeputusanKomite {
	salin := *k
	salin.akhir = p
	return &salin
}

// anggotaBerjalan mencari anak tangga `urut` - murni.
func anggotaBerjalan(k repository.KasusKomite, urut int) (repository.AnggotaKasus, bool) {
	for _, a := range k.Tangga {
		if a.Urut == urut {
			return a, true
		}
	}
	return repository.AnggotaKasus{}, false
}

// periksaGiliran memastikan pelaku adalah anggota berjalan yang belum memutuskan.
//
// ⛔ Anggota BERJALAN = anak tangga ber-`KOMITE_URUT` = `KomiteCount`
// (`KomitePostAdjustment` langkah 2 `Local.Komite = KomiteCount`). ADR-0014:
// anggota lain - termasuk tingkat lebih tinggi - tidak memutuskan atas nama
// tingkat yang sedang berjalan (eskalasi manual = tiket 03).
func periksaGiliran(k repository.KasusKomite, akunID string) error {
	b := k.Baris
	if models.KasusTertutup(b.StatusWork) ||
		!models.KasusDiTangga(b.AcceptStatus, b.KomiteCount, b.KomiteLoop) {
		return fmt.Errorf("%w: kasus %q", ErrTanggaKomiteBerhenti, b.KasusID)
	}
	a, ada := anggotaBerjalan(k, b.KomiteCount)
	if !ada {
		return fmt.Errorf("%w: kasus %q tanpa anak tangga %d", ErrTanggaKomiteBerhenti,
			b.KasusID, b.KomiteCount)
	}
	if strings.TrimSpace(a.OperatorID) == "" || a.OperatorID != akunID {
		return fmt.Errorf("%w: tingkat %d kasus %q bukan milik pelaku",
			ErrTanpaWewenang, b.KomiteCount, b.KasusID)
	}
	if a.Approval != repository.ApprovalKomiteMenunggu {
		return fmt.Errorf("%w: tingkat %d sudah diputuskan", ErrKeputusanKomiteBersamaan, b.KomiteCount)
	}
	return nil
}

// Putuskan mencatat keputusan anggota berjalan dan memajukan tangganya.
func (k *KeputusanKomite) Putuskan(ctx context.Context, pelaku Pelaku,
	kasusID, keputusan, komentar string, saat time.Time) (HasilKeputusanKomite, error) {

	if err := WajibIdentitas(pelaku); err != nil {
		return HasilKeputusanKomite{}, err
	}
	// Bentuk keputusan diperiksa SEBELUM basis data - nilai asing ditolak
	// terang di mesin mana pun.
	if _, err := models.TerapkanKeputusanKomite(keputusan, 1, 1); err != nil {
		return HasilKeputusanKomite{}, err
	}
	if k == nil || k.svc == nil || !k.svc.PunyaDatabase() {
		return HasilKeputusanKomite{}, repository.ErrTanpaOracle
	}
	if strings.TrimSpace(kasusID) == "" {
		return HasilKeputusanKomite{}, fmt.Errorf("%w: id kasus komite kosong", ErrPermintaanTidakSah)
	}
	baca := repository.NewInboxKomite(k.svc.db)
	kasus, err := baca.Kasus(ctx, kasusID)
	if err != nil {
		return HasilKeputusanKomite{}, err
	}
	// ⛔ Butir bb: klaim induk yang sudah ditutup tidak menerima keputusan
	// komite lagi. Diperiksa lewat gerbang Claim Life yang sama.
	klaimID := kasus.Baris.KlaimID
	if err := k.svc.PastikanKasusTerbuka(ctx, klaimID); err != nil {
		return HasilKeputusanKomite{}, err
	}
	if err := periksaGiliran(kasus, pelaku.AkunID); err != nil {
		return HasilKeputusanKomite{}, err
	}
	akibat, err := models.TerapkanKeputusanKomite(keputusan, kasus.Baris.KomiteCount,
		kasus.Baris.KomiteLoop)
	if err != nil {
		return HasilKeputusanKomite{}, err
	}

	err = k.svc.DalamTransaksi(ctx, func(tx *repository.Tx) error {
		// Langkah 3 + 13 - satu tulisan bersyarat untuk dua baris.
		if err := baca.CatatKeputusan(ctx, tx, kasusID, akibat.TingkatDiputus,
			pelaku.AkunID, akibat.Keputusan, komentar, saat); err != nil {
			return err
		}
		// Langkah 4 / 5 - tingkat akhir.
		switch {
		case akibat.AkseptasiAkhir:
			if err := k.akhir.Akseptasi(ctx, tx, kasus, pelaku, saat); err != nil {
				return err
			}
		case akibat.TolakAkhir:
			if err := k.akhir.Tolak(ctx, tx, kasus, pelaku, saat); err != nil {
				return err
			}
		}
		// ADR-0007: satu jejak per tingkat, di transaksi yang sama.
		return k.jejak.Rekam(ctx, tx, CatatanJejak{
			AdjustmentID: kasus.AdjID,
			KlaimID:      klaimID,
			Dari:         "Komite tingkat " + strconv.Itoa(akibat.TingkatDiputus),
			Ke:           models.KataKeputusanKomite(akibat.Keputusan) + " (" + kasusID + ")",
			AkunID:       pelaku.AkunID,
			Waktu:        saat,
		})
	})
	if err != nil {
		return HasilKeputusanKomite{}, err
	}
	hasil := HasilKeputusanKomite{
		TingkatDiputus: akibat.TingkatDiputus,
		Keputusan:      akibat.Keputusan,
		KataKeputusan:  models.KataKeputusanKomite(akibat.Keputusan),
		Berlanjut:      akibat.Berlanjut,
		AkseptasiAkhir: akibat.AkseptasiAkhir,
		TolakAkhir:     akibat.TolakAkhir,
	}
	if akibat.Berlanjut {
		hasil.TingkatBerikut = akibat.CountBaru
	}
	return hasil, nil
}
