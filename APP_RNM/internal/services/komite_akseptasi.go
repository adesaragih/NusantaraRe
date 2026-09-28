package services

// Penyelesai tingkat akhir Komite - tiket 04a (nomor akseptasi).
//
// Untuk apa berkas ini: `KomitePostAdjustment` langkah 4 "Approve Last
// Komite" - DI DALAM transaksi keputusan tingkat akhir (`Putuskan`).
//
//	nomor    4.7-4.12  awalan `KODE_PRODUKSI` → penghitung
//	                   `(ASM-FW-GCNMFW-Work-KomiteLife, <awalan>A)` → rakit
//	stempel  4.11/4.12 baris adjustment `ACCEPTEDNO`, `STS_REJECT = 1`,
//	                   `ACCEPTATION_DATE`; peserta `STS_REJECT = 1`
//
// ⛔ PENGHITUNGNYA DITIRU, prosedurnya TIDAK dipanggil - keputusan o (o2/o3),
// sama dengan penomoran PremiumList: `repository.Penomor.UrutNomorBerikut`.
//
// ⛔ Stempel status memakai fungsi Claim Life yang SUDAH ADA
// (`PerbaruiStatusBaris`, `CerminkanHeader`) - dipanggil, tidak disalin (km3).
// Penjaga `STS_REJECT = 0` di sana sekaligus menjadi gerbang "lahir sekali"
// (`ACCEPTEDNO == ""` di XML): baris yang sudah diaksep tidak diaksep lagi.
//
// Tolak di tingkat akhir (langkah 5) - tiket 05, di berkas yang sama supaya
// rekam akhirnya SATU jalur (tiket 04b).
//
// Dibaca sesudah: komite_keputusan.go, models/komite_nomor.go.

import (
	"context"
	"fmt"
	"strings"
	"time"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
)

// penyelesaiAkhirOracle adalah penyelesai langkah 4/5 yang memakai Oracle.
type penyelesaiAkhirOracle struct {
	svc   *Service
	jejak Jejak
}

// PenyelesaiAkhirKomiteOracle menyusunnya, dengan jejak Oracle.
func PenyelesaiAkhirKomiteOracle(svc *Service) PenyelesaiAkhirKomite {
	return penyelesaiAkhirOracle{svc: svc, jejak: PerekamJejakOracle(svc)}
}

// Akseptasi - langkah 4.
func (p penyelesaiAkhirOracle) Akseptasi(ctx context.Context, tx *repository.Tx,
	kasus repository.KasusKomite, pelaku Pelaku, saat time.Time) (string, error) {

	// ⛔ Gerbang wewenang DIULANG di sini, bukan hanya di `Putuskan`: penulis
	// status wajib memegang gerbangnya sendiri (penjaga
	// `TestSetiapPenulisStatusBergerbangPeran`). Snapshot kasusnya sama.
	if err := periksaGiliran(kasus, pelaku.AkunID); err != nil {
		return "", err
	}
	klaimID := kasus.Baris.KlaimID
	if strings.TrimSpace(kasus.AdjID) == "" || strings.TrimSpace(kasus.PesertaID) == "" {
		return "", fmt.Errorf("%w: kasus %q tanpa baris adjustment/peserta", ErrPermintaanTidakSah,
			kasus.Baris.KasusID)
	}
	baca := repository.NewKlaimLife(p.svc.db)
	// `TempOpenPage.PolicyDataLife.Type` / `.BusinessCode` - klaim induk.
	tipe, err := baca.TypeKlaim(ctx, klaimID)
	if err != nil {
		return "", err
	}
	hdr, err := baca.AmbilHeader(ctx, klaimID)
	if err != nil {
		return "", err
	}
	if hdr == nil || strings.TrimSpace(hdr.KodeBisnis) == "" {
		return "", fmt.Errorf("%w: klaim %q", ErrKodeBisnisBelumTersimpan, klaimID)
	}
	nomor, err := p.terbitkanNomor(ctx, tx, tipe, hdr.KodeBisnis, saat)
	if err != nil {
		return "", err
	}
	// 4.11/4.12 - stempel, lewat mesin status Claim Life yang ada.
	if err := baca.PerbaruiStatusBaris(ctx, tx, kasus.PesertaID, kasus.AdjID,
		models.KodeOutstanding, models.KodeAksep, nomor, saat); err != nil {
		return "", err
	}
	if err := baca.CerminkanHeader(ctx, tx, klaimID, models.KodeAksep, nomor); err != nil {
		return "", err
	}
	// 4.15 "Insert ke OS" - tiket 04b, jalur tunggal berparameter status.
	if err := p.rekamAkhir(ctx, tx, kasus, models.KodeAksep, nomor, saat); err != nil {
		return "", err
	}
	// ADR-U-0007: transisi status baris Outstanding → Aksep punya jejaknya
	// sendiri, bentuknya sama dengan jalur akseptasi Claim Life.
	if err := p.jejak.Rekam(ctx, tx, CatatanJejak{
		AdjustmentID: kasus.AdjID,
		KlaimID:      klaimID,
		Dari:         models.KodeOutstanding,
		Ke:           models.KodeAksep,
		AkunID:       pelaku.AkunID,
		Waktu:        saat,
	}); err != nil {
		return "", err
	}
	return nomor, nil
}

// terbitkanNomor - 4.7-4.12.
//
// ⛔ KEUNIKAN LINTAS JALUR, dan kebijakannya TERBUKA. Bentuk nomor Komite sama
// dengan nomor jalur Claim Life, penghitungnya berbeda. Korpus Komite tidak
// memeriksa tabrakan (nol `GetAcceptedNoCL`). Di sini diperiksa DUA tempat -
// tabel datar warisan (`NomorAkseptasiDipakai`, cara Claim Life) dan
// `T_CLAIMLF_ADJUSTMENT` (tempat kedua jalur kini menulis) - dan tabrakan
// GAGAL TERANG: transaksinya batal, penghitungnya ikut batal.
// `[terbuka — work owner]` OQ-K-04a: bila tabrakan terbukti terjadi, apakah
// nomor dilewati, atau seri dipisah (mis. awalan berbeda).
func (p penyelesaiAkhirOracle) terbitkanNomor(ctx context.Context, tx *repository.Tx,
	tipe, kodeBisnis string, saat time.Time) (string, error) {

	penghitung := repository.NewPenomor(p.svc.db)
	awalan, err := penghitung.AwalanProduksi(ctx, tx, repository.TipeKodeProduksiLife)
	if err != nil {
		return "", err
	}
	hari, err := penghitung.HariClosing(ctx, tx)
	if err != nil {
		return "", err
	}
	periode, err := repository.HitungPeriodeNomor(saat, hari)
	if err != nil {
		return "", err
	}
	urut, err := penghitung.UrutNomorBerikut(ctx, tx, models.ClassPenghitungKomiteLife,
		models.JenisPenghitungKomite(awalan), periode, saat)
	if err != nil {
		return "", err
	}
	nomor, err := models.NomorAkseptasiKomite(awalan, tipe, kodeBisnis, periode.MMYYYY, urut)
	if err != nil {
		return "", err
	}
	dipakai, err := repository.NewPohonKlaim(p.svc.db).NomorAkseptasiDipakai(ctx, tx, nomor)
	if err != nil {
		return "", err
	}
	if !dipakai {
		dipakai, err = repository.NewInboxKomite(p.svc.db).NomorAkseptasiDipakaiDiAdjustment(ctx, tx, nomor)
		if err != nil {
			return "", err
		}
	}
	if dipakai {
		return "", fmt.Errorf("%w: %q (tabrakan jalur Komite dengan jalur Claim Life)",
			repository.ErrNomorAkseptasiBerganda, nomor)
	}
	return nomor, nil
}

// rekamAkhir adalah SATU jalur simpan rekam akseptasi - tiket 04b.
//
// `[terverifikasi]` `KomitePostAdjustment` langkah 4.15 (aksep) dan 5.6 (tolak)
// sama-sama "Insert ke OS" lewat `UpdateOsAkseptasiClaimLife_sql`, didahului
// precondition yang SAMA dan mengisi properti yang SAMA - salin-tempel yang
// `[keputusan work owner]` disatukan. Di sini keduanya memanggil fungsi ini;
// yang berbeda hanya `status` (dan nomor, yang hanya ada pada aksep).
func (p penyelesaiAkhirOracle) rekamAkhir(ctx context.Context, tx *repository.Tx,
	kasus repository.KasusKomite, status, nomor string, saat time.Time) error {
	return repository.NewInboxKomite(p.svc.db).RekamAkhirWarisan(ctx, tx, kasus.AdjID,
		status, nomor, saat)
}

// Tolak - langkah 5 "Reject" (gerbang b8119), tiket 05 Komite.
//
// `[terverifikasi]` `KomitePostAdjustment` langkah 5:
//
//	5.1 "Set Reject Komite berjenjang"  SELURUH KomiteList: Aproval = 2,
//	                                    Comment = "", DateApprove = now  - TANPA syarat
//	5.3 "Set Nilai Akseptasi"           AdjustmentList.STS_REJECT = 2,
//	                                    PremiumListDetail.STS_REJECT = 2, .IsCheck = "false"
//	5.6 "Insert ke OS"                  UpdateOsAkseptasiClaimLife_sql
//
// ⚠️ 5.1 TIDAK DITIRU - OQ-K-05 `[terbuka — work owner]`. Ia menimpa keputusan
// `1` dan komentar SETIAP tingkat sebelumnya dengan `2` dan teks kosong:
// riwayat tangga (tiket 09) dan jejak per tingkat (ADR-0007) kehilangan
// siapa yang menyetujui sebelum tingkat akhir menolak. Yang ditulis di sini
// hanya keputusan tingkat akhir itu sendiri (`CatatKeputusan`, tiket 02).
//
// ⛔ `AcceptStatus` TIDAK diteruskan ke Claim Life (ADR-0001): ia dipetakan di
// batas ini menjadi `STS_REJECT = 2` - kode Claim Life `KodeDitolak`.
func (p penyelesaiAkhirOracle) Tolak(ctx context.Context, tx *repository.Tx,
	kasus repository.KasusKomite, pelaku Pelaku, saat time.Time) error {

	// Gerbang penulis status - lihat Akseptasi.
	if err := periksaGiliran(kasus, pelaku.AkunID); err != nil {
		return err
	}
	klaimID := kasus.Baris.KlaimID
	if strings.TrimSpace(kasus.AdjID) == "" || strings.TrimSpace(kasus.PesertaID) == "" {
		return fmt.Errorf("%w: kasus %q tanpa baris adjustment/peserta", ErrPermintaanTidakSah,
			kasus.Baris.KasusID)
	}
	baca := repository.NewKlaimLife(p.svc.db)
	// 5.3 - dua tingkat baris, nilai yang sama, satu operasi (AC 15 spec).
	// Penjaga `STS_REJECT = 0` di `PerbaruiStatusBaris` = gerbang "masih
	// Outstanding"; nomor dan tanggal akseptasi dikosongkan (nol pada Tolak).
	if err := baca.PerbaruiStatusBaris(ctx, tx, kasus.PesertaID, kasus.AdjID,
		models.KodeOutstanding, models.KodeDitolak, "", time.Time{}); err != nil {
		return err
	}
	// 5.3 `.IsCheck = "false"` - peserta dapat dipilih ulang dengan baris
	// pengganti di Claim Life (AC 2 tiket ini). Fungsi Claim Life yang ada.
	if err := baca.CabutPenandaDipilih(ctx, tx, kasus.PesertaID); err != nil {
		return err
	}
	if err := baca.CerminkanHeader(ctx, tx, klaimID, models.KodeDitolak, ""); err != nil {
		return err
	}
	// 5.6 - jalur tunggal rekam akhir, status 2, tanpa nomor.
	if err := p.rekamAkhir(ctx, tx, kasus, models.KodeDitolak, "", saat); err != nil {
		return err
	}
	return p.jejak.Rekam(ctx, tx, CatatanJejak{
		AdjustmentID: kasus.AdjID,
		KlaimID:      klaimID,
		Dari:         models.KodeOutstanding,
		Ke:           models.KodeDitolak,
		AkunID:       pelaku.AkunID,
		Waktu:        saat,
	})
}
