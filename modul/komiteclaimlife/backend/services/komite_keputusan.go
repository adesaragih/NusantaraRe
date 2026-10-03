package services

// Keputusan satu tingkat Komite - tiket 02 Komite Claim Life.
//
// Untuk apa berkas ini: tombol `Submit` layar `ShowTransfer` (b34722 →
// `finishAssignment` → activity flow action `KomitePostAdjustment`). Satu
// anggota BERJALAN memutuskan Setuju/Tolak; tangganya naik satu tingkat atau
// berhenti. Aturannya murni di `models/komite_tangga.go`.
//
// ⛔ TINGKAT AKHIR: `KomitePostAdjustment` langkah 4 (akseptasi, tiket
// 04a/04b) dan 5 (tolak ke baris klaim, tiket 05) berjalan DI DALAM
// transaksi keputusan lewat `PenyelesaiAkhirKomite` - kini
// `PenyelesaiAkhirKomiteOracle` (komite_akseptasi.go). Bawaan
// `PenyelesaiAkhirBelumAda` tetap MENOLAK, supaya layanan yang disusun tanpa
// penyelesai tidak pernah menutup tangga setengah jalan.
//
// Dibaca sesudah: models/komite_tangga.go, repository/komite_keputusan.go.

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
	"nusantarare/inti/backend/jejak"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/komiteclaimlife/backend/models"
	"nusantarare/modul/komiteclaimlife/backend/repository"
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
	// Akseptasi - langkah 4 "Approve Last Komite" (tiket 04a/04b). Mengembalikan
	// nomor akseptasi yang lahir.
	Akseptasi(ctx context.Context, tx *db.Tx, kasus repository.KasusKomite,
		pelaku inti.Pelaku, saat time.Time) (string, error)
	// Tolak - langkah 5 "Reject" (tiket 05).
	Tolak(ctx context.Context, tx *db.Tx, kasus repository.KasusKomite,
		pelaku inti.Pelaku, saat time.Time) error
}

// PenyelesaiAkhirBelumAda gagal terang untuk kedua langkah.
type PenyelesaiAkhirBelumAda struct{}

// Akseptasi selalu gagal.
func (PenyelesaiAkhirBelumAda) Akseptasi(context.Context, *db.Tx,
	repository.KasusKomite, inti.Pelaku, time.Time) (string, error) {
	return "", ErrPenyelesaianAkhirBelumAda
}

// Tolak selalu gagal.
func (PenyelesaiAkhirBelumAda) Tolak(context.Context, *db.Tx,
	repository.KasusKomite, inti.Pelaku, time.Time) error {
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
	// NomorAkseptasi terisi hanya pada Setuju di tingkat akhir (tiket 04a).
	NomorAkseptasi string `json:"nomorAkseptasi"`
	// EfekTertunda - efek yang DIANTRE, belum tuntas (tiket 06, ADR-0015).
	EfekTertunda []string `json:"efekTertunda"`
}

// KeputusanKomite melayani keputusan satu tingkat.
type KeputusanKomite struct {
	svc   *Service
	jejak jejak.Jejak
	akhir PenyelesaiAkhirKomite
}

// KeputusanKomite menyusunnya dengan jejak dan penyelesai akhir yang gagal terang.
func (s *Service) KeputusanKomite() *KeputusanKomite {
	return &KeputusanKomite{svc: s, jejak: jejak.JejakBelumDiputuskan{}, akhir: PenyelesaiAkhirBelumAda{}}
}

// DenganJejak mengganti perekamnya.
func (k *KeputusanKomite) DenganJejak(j jejak.Jejak) *KeputusanKomite {
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
	if kontrak.KasusTertutup(b.StatusWork) ||
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
			inti.ErrTanpaWewenang, b.KomiteCount, b.KasusID)
	}
	if a.Approval != repository.ApprovalKomiteMenunggu {
		return fmt.Errorf("%w: tingkat %d sudah diputuskan", ErrKeputusanKomiteBersamaan, b.KomiteCount)
	}
	return nil
}

// Putuskan mencatat keputusan anggota berjalan dan memajukan tangganya.
func (k *KeputusanKomite) Putuskan(ctx context.Context, pelaku inti.Pelaku,
	kasusID, keputusan, komentar string, saat time.Time) (HasilKeputusanKomite, error) {

	if err := inti.WajibIdentitas(pelaku); err != nil {
		return HasilKeputusanKomite{}, err
	}
	// Bentuk keputusan diperiksa SEBELUM basis data - nilai asing ditolak
	// terang di mesin mana pun.
	if _, err := models.TerapkanKeputusanKomite(keputusan, 1, 1); err != nil {
		return HasilKeputusanKomite{}, err
	}
	if k == nil || k.svc == nil || !k.svc.PunyaDatabase() {
		return HasilKeputusanKomite{}, db.ErrTanpaOracle
	}
	if strings.TrimSpace(kasusID) == "" {
		return HasilKeputusanKomite{}, fmt.Errorf("%w: id kasus komite kosong", galat.ErrPermintaanTidakSah)
	}
	baca := repository.NewInboxKomite(k.svc.DB())
	kasus, err := baca.Kasus(ctx, kasusID)
	if err != nil {
		return HasilKeputusanKomite{}, err
	}
	// ⛔ Butir bb: klaim induk yang sudah ditutup tidak menerima keputusan
	// komite lagi. Diperiksa lewat gerbang Claim Life yang sama.
	klaimID := kasus.Baris.KlaimID
	if err := k.svc.Klaim().PastikanKasusTerbuka(ctx, klaimID); err != nil {
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

	var nomorAksep string
	var tertunda []string
	err = k.svc.DalamTransaksi(ctx, func(tx *db.Tx) error {
		// Langkah 3 + 13 - satu tulisan bersyarat untuk dua baris.
		if err := baca.CatatKeputusan(ctx, tx, kasusID, akibat.TingkatDiputus,
			pelaku.AkunID, akibat.Keputusan, komentar, saat); err != nil {
			return err
		}
		// Langkah 4 / 5 - tingkat akhir.
		switch {
		case akibat.AkseptasiAkhir:
			n, err := k.akhir.Akseptasi(ctx, tx, kasus, pelaku, saat)
			if err != nil {
				return err
			}
			nomorAksep = n
		case akibat.TolakAkhir:
			if err := k.akhir.Tolak(ctx, tx, kasus, pelaku, saat); err != nil {
				return err
			}
		}
		// Tiket 06 - efek keluar DIANTRE di transaksi yang sama (ADR-0015):
		// keputusan tanpa antreannya, atau antrean tanpa keputusannya, tidak
		// mungkin.
		t, err := antreEfekKomite(ctx, k.svc, tx, kasus, akibat, pelaku, saat)
		if err != nil {
			return err
		}
		tertunda = t
		// ADR-0007: satu jejak per tingkat, di transaksi yang sama.
		return k.jejak.Rekam(ctx, tx, jejak.CatatanJejak{
			AdjustmentID: kasus.AdjID,
			KlaimID:      klaimID,
			Dari:         awalanJejakTingkat + strconv.Itoa(akibat.TingkatDiputus),
			Ke:           models.KataKeputusanKomite(akibat.Keputusan) + " (" + kasusID + ")",
			AkunID:       pelaku.AkunID,
			Waktu:        saat,
			// Migrasi 021 (GILIRAN-17): komentar keputusan ikut tercatat.
			Komentar: komentar,
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
		NomorAkseptasi: nomorAksep,
		EfekTertunda:   tertunda,
	}
	if akibat.Berlanjut {
		hasil.TingkatBerikut = akibat.CountBaru
	}
	return hasil, nil
}

// HasilEskalasiKomite adalah jawaban `Eskalasi`.
type HasilEskalasiKomite struct {
	DariTingkat int `json:"dariTingkat"`
	KeTingkat   int `json:"keTingkat"`
}

// Eskalasi memindahkan kasus NAIK SATU tingkat - tiket 03, ADR-0014.
//
// ⚠️ `[asumsi — OQ-007/OQ-021]` "admin komite" = peran `ReasLifeAdmin`
// (`PeranAdmin`). Korpus tidak memuat satu pun rule otorisasi (Identity &
// Access ABSENT, `discovery/context-map.md`), dan tidak ada peran admin
// komite yang terekspor. Peran itu dipilih karena ia yang MENYERAHKAN kasus
// ke Komite (`Penyerahan.Serahkan`); bila model RBAC memutuskan lain, hanya
// konstanta ini yang berubah.
//
// ⛔ Eskalasi BUKAN pintu belakang keputusan: ia tidak mencatat `1`/`2`
// untuk siapa pun, dan admin yang sama tetap tidak dapat memutuskan atas nama
// tingkat mana pun (`periksaGiliran`).
func (k *KeputusanKomite) Eskalasi(ctx context.Context, pelaku inti.Pelaku,
	kasusID string, saat time.Time) (HasilEskalasiKomite, error) {

	if err := inti.WajibIdentitas(pelaku); err != nil {
		return HasilEskalasiKomite{}, err
	}
	if !pelaku.PunyaPeran(inti.PeranAdmin) {
		return HasilEskalasiKomite{}, fmt.Errorf("%w: eskalasi komite menuntut peran %s",
			inti.ErrTanpaWewenang, inti.PeranAdmin)
	}
	if k == nil || k.svc == nil || !k.svc.PunyaDatabase() {
		return HasilEskalasiKomite{}, db.ErrTanpaOracle
	}
	if strings.TrimSpace(kasusID) == "" {
		return HasilEskalasiKomite{}, fmt.Errorf("%w: id kasus komite kosong", galat.ErrPermintaanTidakSah)
	}
	baca := repository.NewInboxKomite(k.svc.DB())
	kasus, err := baca.Kasus(ctx, kasusID)
	if err != nil {
		return HasilEskalasiKomite{}, err
	}
	klaimID := kasus.Baris.KlaimID
	if err := k.svc.Klaim().PastikanKasusTerbuka(ctx, klaimID); err != nil {
		return HasilEskalasiKomite{}, err
	}
	b := kasus.Baris
	if kontrak.KasusTertutup(b.StatusWork) || !models.KasusDiTangga(b.AcceptStatus, b.KomiteCount, b.KomiteLoop) {
		return HasilEskalasiKomite{}, fmt.Errorf("%w: kasus %q", ErrTanggaKomiteBerhenti, kasusID)
	}
	ke, err := models.EskalasiNaik(b.KomiteCount, b.KomiteLoop)
	if err != nil {
		return HasilEskalasiKomite{}, err
	}
	err = k.svc.DalamTransaksi(ctx, func(tx *db.Tx) error {
		if err := baca.Eskalasi(ctx, tx, kasusID, b.KomiteCount); err != nil {
			return err
		}
		// AC 12: siapa, kapan, dari tingkat mana ke mana.
		return k.jejak.Rekam(ctx, tx, jejak.CatatanJejak{
			AdjustmentID: kasus.AdjID,
			KlaimID:      klaimID,
			Dari:         awalanJejakTingkat + strconv.Itoa(b.KomiteCount),
			Ke:           awalanJejakEskalasi + strconv.Itoa(ke) + " (" + kasusID + ")",
			AkunID:       pelaku.AkunID,
			Waktu:        saat,
		})
	})
	if err != nil {
		return HasilEskalasiKomite{}, err
	}
	return HasilEskalasiKomite{DariTingkat: b.KomiteCount, KeTingkat: ke}, nil
}

// ErrEskalasiTanpaTingkatAtas - eskalasi dari tingkat akhir.
var ErrEskalasiTanpaTingkatAtas = models.ErrEskalasiTanpaTingkatAtas
