package services

// Untuk apa berkas ini: KONTRAK yang Claim Non Prop sediakan untuk Komite Claim Non Prop - `kontrak.KlaimTreatyNonPropKomite`
// (perintah work owner 09-10-2026: komite Non Prop ikut pola Claim Prop; disalin dari
// `modul/claimprop/backend/services/kontrak_komite.go`, bukan impor). Komite membaca kasus klaim induk dan, lewat
// `KomitePostAdjustment`, menulis kembali ke kasus itu di DALAM transaksi Komite. Inilah satu-satunya jalan modul lain
// menulis tabel klaim Non Prop; jalur yang boleh ditulis dibatasi daftar putih di paket kontrak.
//
// Bacaan = `muat` + `turunkan` (halaman persis seperti dibuka layar Claim Non Prop) ditambah
// `OfferFacIn.QuotationData.BusinessOldId` bila belum terisi: properti itu tanpa penulis di korpus (OQ-CNP-37), padahal
// nomor akseptasi `KomitePostAdjustment` S14.11 memakainya - diisi `ClaimData.QuotationData.BusinessOldId` (Save to issue
// RNM) atau OLDID bisnis menurut nama (GetDataBusiness_SQL), pola `PeriksaOldID` Claim Prop. Tulisan = `BacaHalaman` ->
// ubahan -> `SimpanHalaman` + `SentuhKasus`, jalur yang sama dengan setiap aksi Claim Non Prop.

import (
	"context"
	"errors"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/claimnonprop/backend/models"
	"nusantarare/modul/claimnonprop/backend/repository"
)

// KlaimUntukKomite memenuhi `kontrak.KlaimTreatyNonPropKomite`.
type KlaimUntukKomite struct{ l *Layanan }

// KlaimUntukKomite menyusun penyedia kontrak Komite Claim Non Prop di atas layanan ini.
func (l *Layanan) KlaimUntukKomite() KlaimUntukKomite { return KlaimUntukKomite{l: l} }

var _ kontrak.KlaimTreatyNonPropKomite = KlaimUntukKomite{}

func (k KlaimUntukKomite) siap() error {
	if !k.l.AdaGudang() {
		return ErrTanpaOracle
	}
	return nil
}

// galatKasus menerjemahkan galat repository ke galat kontrak.
func galatKasus(err error) error {
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return kontrak.ErrKlaimTreatyTidakAda
	}
	return err
}

// posisiAdjustment - posisi (1..n) baris akseptasi ber-ID `adjID`; 0 bila tidak ada.
func posisiAdjustment(h *models.Halaman, adjID string) int {
	for i, b := range h.AmbilDaftar(models.DaftarAdjustment) {
		if adjID != "" && b[models.PropID] == adjID {
			return i + 1
		}
	}
	return 0
}

// BacaKlaimTreaty - lihat `kontrak.KlaimTreatyNonPropKomite`.
func (k KlaimUntukKomite) BacaKlaimTreaty(ctx context.Context, tx *db.Tx, klaimID, adjID string) (kontrak.KlaimTreaty,
	error) {
	if err := k.siap(); err != nil {
		return kontrak.KlaimTreaty{}, err
	}
	kasus, h, err := k.l.muat(ctx, tx, klaimID)
	if err != nil {
		return kontrak.KlaimTreaty{}, galatKasus(err)
	}
	if err := k.l.turunkan(ctx, h); err != nil {
		return kontrak.KlaimTreaty{}, err
	}
	if h.Ambil(models.OQ+"BusinessOldId") == "" { // OQ-CNP-37: properti tanpa penulis, pola PeriksaOldID Claim Prop
		old := h.Ambil(models.CD + "QuotationData.BusinessOldId")
		if nama := h.Ambil(models.OQ + "BusinessName"); old == "" && nama != "" {
			if old, err = k.l.a.KodeLamaBisnis(ctx, nama); err != nil {
				return kontrak.KlaimTreaty{}, err
			}
		}
		h.Setel(models.OQ+"BusinessOldId", old)
	}
	n := posisiAdjustment(h, adjID)
	if n == 0 {
		return kontrak.KlaimTreaty{}, kontrak.ErrAdjustmentTreatyTidakAda
	}
	out := kontrak.KlaimTreaty{Nilai: map[string]string{}, Daftar: map[string][]map[string]string{}, Adjustment: n,
		Tertutup: kasus.Tertutup()}
	for j, v := range h.Nilai {
		out.Nilai[j] = v
	}
	for j, rows := range h.Daftar {
		salin := make([]map[string]string, 0, len(rows))
		for _, b := range rows {
			m := make(map[string]string, len(b))
			for p, v := range b {
				m[p] = v
			}
			salin = append(salin, m)
		}
		out.Daftar[j] = salin
	}
	return out, nil
}

// KunciKlaimTreaty - lihat `kontrak.KlaimTreatyNonPropKomite`.
func (k KlaimUntukKomite) KunciKlaimTreaty(ctx context.Context, tx *db.Tx, klaimID string) error {
	if err := k.siap(); err != nil {
		return err
	}
	kasus, err := k.l.g.Keadaan(ctx, tx, klaimID)
	if err != nil {
		return galatKasus(err)
	}
	if kasus.Tertutup() {
		return kontrak.ErrKlaimTreatyTertutup
	}
	if _, err := k.l.g.KunciKasus(ctx, tx, klaimID, kasus.Tahap); err != nil {
		if errors.Is(err, repository.ErrTahapBerubah) {
			return kontrak.ErrKlaimTreatyTertutup
		}
		return galatKasus(err)
	}
	return nil
}

// periksaUbahan menolak jalur di luar daftar putih kontrak (Fac Retro tidak ada di Non Prop).
func periksaUbahan(u kontrak.UbahanKlaimTreaty) error {
	for j := range u.Header {
		if _, ok := kontrak.JalurHeaderKomiteNonProp[j]; !ok {
			return fmt.Errorf("%w: header %q", kontrak.ErrUbahanKlaimTreatyTidakSah, j)
		}
	}
	for p := range u.Adjustment {
		if _, ok := kontrak.PropAdjustmentKomiteNonProp[p]; !ok {
			return fmt.Errorf("%w: adjustment %q", kontrak.ErrUbahanKlaimTreatyTidakSah, p)
		}
	}
	if len(u.FacRetro) > 0 {
		return fmt.Errorf("%w: FacRetroList", kontrak.ErrUbahanKlaimTreatyTidakSah)
	}
	return nil
}

// TulisBalikKlaimTreaty - lihat `kontrak.KlaimTreatyNonPropKomite`.
func (k KlaimUntukKomite) TulisBalikKlaimTreaty(ctx context.Context, tx *db.Tx, klaimID, adjID string,
	u kontrak.UbahanKlaimTreaty) error {
	if err := k.siap(); err != nil {
		return err
	}
	if err := periksaUbahan(u); err != nil {
		return err
	}
	kasus, err := k.l.g.Keadaan(ctx, tx, klaimID)
	if err != nil {
		return galatKasus(err)
	}
	if kasus.Tertutup() {
		return kontrak.ErrKlaimTreatyTertutup
	}
	h, err := k.l.g.BacaHalaman(ctx, tx, klaimID)
	if err != nil {
		return galatKasus(err)
	}
	n := posisiAdjustment(h, adjID)
	if n == 0 {
		return kontrak.ErrAdjustmentTreatyTidakAda
	}
	for j, v := range u.Header {
		h.Setel(j, v)
	}
	b := h.AmbilDaftar(models.DaftarAdjustment)[n-1]
	for p, v := range u.Adjustment {
		b[p] = v
	}
	for _, r := range u.Riwayat {
		kt := &models.Konteks{Ctx: ctx, Pelaku: r.Pelaku, Tingkat: r.Tingkat, Sekarang: r.Saat}
		kt.Riwayat(h, r.Teks)
	}
	if err := k.l.g.SimpanHalaman(ctx, tx, klaimID, h); err != nil {
		return err
	}
	return k.l.g.SentuhKasus(ctx, tx, klaimID, k.l.jam())
}
