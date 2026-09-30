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

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/jejak"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/inti/backend/penomor"
	"nusantarare/modul/komiteclaimlife/models"
	"nusantarare/modul/komiteclaimlife/repository"
)

// penyelesaiAkhirOracle adalah penyelesai langkah 4/5 yang memakai Oracle.
type penyelesaiAkhirOracle struct {
	svc   *Service
	jejak jejak.Jejak
}

// PenyelesaiAkhirKomiteOracle menyusunnya, dengan jejak Oracle.
func PenyelesaiAkhirKomiteOracle(svc *Service) PenyelesaiAkhirKomite {
	return penyelesaiAkhirOracle{svc: svc, jejak: jejak.PerekamJejakOracle(svc)}
}

// Akseptasi - langkah 4.
func (p penyelesaiAkhirOracle) Akseptasi(ctx context.Context, tx *db.Tx,
	kasus repository.KasusKomite, pelaku inti.Pelaku, saat time.Time) (string, error) {

	// ⛔ Gerbang wewenang DIULANG di sini, bukan hanya di `Putuskan`: penulis
	// status wajib memegang gerbangnya sendiri (penjaga
	// `TestSetiapPenulisStatusBergerbangPeran`). Snapshot kasusnya sama.
	if err := periksaGiliran(kasus, pelaku.AkunID); err != nil {
		return "", err
	}
	klaimID := kasus.Baris.KlaimID
	if strings.TrimSpace(kasus.AdjID) == "" || strings.TrimSpace(kasus.PesertaID) == "" {
		return "", fmt.Errorf("%w: kasus %q tanpa baris adjustment/peserta", galat.ErrPermintaanTidakSah,
			kasus.Baris.KasusID)
	}
	// Refactor bentuk B: baris dan header klaim milik Claim Life, dibaca dan
	// ditulis lewat inti/backend/kontrak.KlaimKomite di dalam transaksi ini.
	baca := p.svc.Klaim()
	// `TempOpenPage.PolicyDataLife.Type` / `.BusinessCode` - klaim induk.
	tipe, err := baca.TypeKlaim(ctx, klaimID)
	if err != nil {
		return "", err
	}
	kodeBisnis, err := baca.KodeBisnisKlaim(ctx, klaimID)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(kodeBisnis) == "" {
		return "", fmt.Errorf("%w: klaim %q", kontrak.ErrKodeBisnisBelumTersimpan, klaimID)
	}
	nomor, err := p.terbitkanNomor(ctx, tx, tipe, kodeBisnis, saat)
	if err != nil {
		return "", err
	}
	// 4.11/4.12 - stempel, lewat mesin status Claim Life yang ada.
	if err := baca.PerbaruiStatusBaris(ctx, tx, kasus.PesertaID, kasus.AdjID,
		kontrak.KodeOutstanding, kontrak.KodeAksep, nomor, saat); err != nil {
		return "", err
	}
	if err := baca.CerminkanHeader(ctx, tx, klaimID, kontrak.KodeAksep, nomor); err != nil {
		return "", err
	}
	// 4.15 "Insert ke OS" - tiket 04b, jalur tunggal berparameter status.
	if err := p.rekamAkhir(ctx, tx, kasus, kontrak.KodeAksep, nomor, saat); err != nil {
		return "", err
	}
	// ADR-U-0007: transisi status baris Outstanding → Aksep punya jejaknya
	// sendiri, bentuknya sama dengan jalur akseptasi Claim Life.
	if err := p.jejak.Rekam(ctx, tx, jejak.CatatanJejak{
		AdjustmentID: kasus.AdjID,
		KlaimID:      klaimID,
		Dari:         kontrak.KodeOutstanding,
		Ke:           kontrak.KodeAksep,
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
func (p penyelesaiAkhirOracle) terbitkanNomor(ctx context.Context, tx *db.Tx,
	tipe, kodeBisnis string, saat time.Time) (string, error) {

	penghitung := penomor.NewPenomor(p.svc.DB())
	awalan, err := penghitung.AwalanProduksi(ctx, tx, penomor.TipeKodeProduksiLife)
	if err != nil {
		return "", err
	}
	hari, err := penghitung.HariClosing(ctx, tx)
	if err != nil {
		return "", err
	}
	periode, err := penomor.HitungPeriodeNomor(saat, hari)
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
	dipakai, err := p.svc.Klaim().NomorAkseptasiDipakai(ctx, tx, nomor)
	if err != nil {
		return "", err
	}
	if !dipakai {
		dipakai, err = repository.NewInboxKomite(p.svc.DB()).NomorAkseptasiDipakaiDiAdjustment(ctx, tx, nomor)
		if err != nil {
			return "", err
		}
	}
	if dipakai {
		return "", fmt.Errorf("%w: %q (tabrakan jalur Komite dengan jalur Claim Life)",
			kontrak.ErrNomorAkseptasiBerganda, nomor)
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
func (p penyelesaiAkhirOracle) rekamAkhir(ctx context.Context, tx *db.Tx,
	kasus repository.KasusKomite, status, nomor string, saat time.Time) error {
	return repository.NewInboxKomite(p.svc.DB()).RekamAkhirWarisan(ctx, tx, kasus.AdjID,
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
// ⛔ 5.1 DITIRU - OQ-K-05 DITUTUP 29-09-2026 (GILIRAN-17) `[keputusan work
// owner]`: satu UPDATE bersyarat menimpa keputusan seluruh tingkat yang
// memutus dengan `2`, komentar kosong, `DateApprove` sekarang. SEBELUM
// menimpa, keputusan dan komentar lama tiap tingkat dicatat di jejak
// (ADR-0007, kolom `KOMENTAR` migrasi 021), sehingga riwayat tangga tiket 09
// tetap dapat membacanya (`susunRiwayat`, `Asli`).
//
// ⛔ `AcceptStatus` TIDAK diteruskan ke Claim Life (ADR-0001): ia dipetakan di
// batas ini menjadi `STS_REJECT = 2` - kode Claim Life `KodeDitolak`.
func (p penyelesaiAkhirOracle) Tolak(ctx context.Context, tx *db.Tx,
	kasus repository.KasusKomite, pelaku inti.Pelaku, saat time.Time) error {

	// Gerbang penulis status - lihat Akseptasi.
	if err := periksaGiliran(kasus, pelaku.AkunID); err != nil {
		return err
	}
	klaimID := kasus.Baris.KlaimID
	if strings.TrimSpace(kasus.AdjID) == "" || strings.TrimSpace(kasus.PesertaID) == "" {
		return fmt.Errorf("%w: kasus %q tanpa baris adjustment/peserta", galat.ErrPermintaanTidakSah,
			kasus.Baris.KasusID)
	}
	// 5.1 - tangga dibaca DI DALAM transaksi (keputusan tingkat akhir langkah
	// 3 ikut terbaca), keputusan lamanya dijejaki, lalu ditimpa.
	komite := repository.NewInboxKomite(p.svc.DB())
	tangga, err := komite.TanggaSebelumDitimpa(ctx, tx, kasus.Baris.KasusID)
	if err != nil {
		return err
	}
	for _, c := range jejakTimpaTangga(kasus, tangga, pelaku, saat) {
		if err := p.jejak.Rekam(ctx, tx, c); err != nil {
			return err
		}
	}
	if err := komite.TimpaTanggaTolakAkhir(ctx, tx, kasus.Baris.KasusID, saat); err != nil {
		return err
	}
	baca := p.svc.Klaim()
	// 5.3 - dua tingkat baris, nilai yang sama, satu operasi (AC 15 spec).
	// Penjaga `STS_REJECT = 0` di `PerbaruiStatusBaris` = gerbang "masih
	// Outstanding"; nomor dan tanggal akseptasi dikosongkan (nol pada Tolak).
	if err := baca.PerbaruiStatusBaris(ctx, tx, kasus.PesertaID, kasus.AdjID,
		kontrak.KodeOutstanding, kontrak.KodeDitolak, "", time.Time{}); err != nil {
		return err
	}
	// 5.3 `.IsCheck = "false"` - peserta dapat dipilih ulang dengan baris
	// pengganti di Claim Life (AC 2 tiket ini). Fungsi Claim Life yang ada.
	if err := baca.CabutPenandaDipilih(ctx, tx, kasus.PesertaID); err != nil {
		return err
	}
	if err := baca.CerminkanHeader(ctx, tx, klaimID, kontrak.KodeDitolak, ""); err != nil {
		return err
	}
	// 5.6 - jalur tunggal rekam akhir, status 2, tanpa nomor.
	if err := p.rekamAkhir(ctx, tx, kasus, kontrak.KodeDitolak, "", saat); err != nil {
		return err
	}
	return p.jejak.Rekam(ctx, tx, jejak.CatatanJejak{
		AdjustmentID: kasus.AdjID,
		KlaimID:      klaimID,
		Dari:         kontrak.KodeOutstanding,
		Ke:           kontrak.KodeDitolak,
		AkunID:       pelaku.AkunID,
		Waktu:        saat,
	})
}
