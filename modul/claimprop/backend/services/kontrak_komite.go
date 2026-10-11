package services

// Untuk apa berkas ini: KONTRAK yang Claim Prop sediakan untuk Komite Claim Prop - `kontrak.KlaimTreatyKomite`
// (keputusan work owner 08-10-2026, pola `KlaimKomite` Claim Life). Komite membaca kasus klaim induk dan, lewat
// `KomitePostAdjustment`, menulis kembali ke kasus itu di DALAM transaksi Komite. Inilah satu-satunya jalan modul lain
// menulis tabel klaim Prop; jalur yang boleh ditulis dibatasi daftar putih di paket kontrak.
//
// Bacaan = `muat` + `turunkan` (halaman persis seperti dibuka layar Claim Prop) ditambah `BusinessOldId` bila belum
// terisi (pola `PeriksaOldID`, SaveOutstanding_Act langkah 3-5). Tulisan = `BacaHalaman` -> ubahan -> `SimpanHalaman`
// + `SentuhKasus`, jalur yang sama dengan setiap aksi Claim Prop.

import (
	"context"
	"errors"
	"fmt"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/claimprop/backend/models"
	"nusantarare/modul/claimprop/backend/repository"
)

// KlaimUntukKomite memenuhi `kontrak.KlaimTreatyKomite`.
type KlaimUntukKomite struct{ l *Layanan }

// KlaimUntukKomite menyusun penyedia kontrak Komite Claim Prop di atas layanan ini.
func (l *Layanan) KlaimUntukKomite() KlaimUntukKomite { return KlaimUntukKomite{l: l} }

var _ kontrak.KlaimTreatyKomite = KlaimUntukKomite{}

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

// posisiAdjustment - posisi (1..n) baris adjustment ber-ID `adjID`; 0 bila tidak ada.
func posisiAdjustment(h *models.Halaman, adjID string) int {
	for i, b := range h.AmbilDaftar(models.DaftarAdjustment) {
		if adjID != "" && b[models.PropID] == adjID {
			return i + 1
		}
	}
	return 0
}

// BacaKlaimTreaty - lihat `kontrak.KlaimTreatyKomite`.
func (k KlaimUntukKomite) BacaKlaimTreaty(ctx context.Context, tx *db.Tx, klaimID, adjID string) (kontrak.KlaimTreaty,
	error) {
	if err := k.siap(); err != nil {
		return kontrak.KlaimTreaty{}, err
	}
	kasus, h, err := k.l.muat(ctx, tx, klaimID)
	if err != nil {
		return kontrak.KlaimTreaty{}, galatKasus(err)
	}
	kt := &models.Konteks{Ctx: ctx, Acuan: k.l.a, Sekarang: k.l.jam(), Produksi: k.l.produksi}
	if err := k.l.turunkan(ctx, kt, h); err != nil {
		return kontrak.KlaimTreaty{}, err
	}
	if h.Ambil(models.OQ+"BusinessOldId") == "" { // PeriksaOldID: BUSINESS.OLDID bila kosong
		if kode := h.Ambil(models.OQ + "BusinessCode"); kode != "" {
			old, err := k.l.a.KodeLamaBisnis(ctx, kode)
			if err != nil {
				return kontrak.KlaimTreaty{}, err
			}
			h.Setel(models.OQ+"BusinessOldId", old)
		}
	}
	n := posisiAdjustment(h, adjID)
	if n == 0 && adjID != "" { // adjID kosong = kasus komite Close Without Payment (Adjustment 0)
		return kontrak.KlaimTreaty{}, kontrak.ErrAdjustmentTreatyTidakAda
	}
	out := kontrak.KlaimTreaty{Nilai: map[string]string{}, Daftar: map[string][]map[string]string{}, Adjustment: n,
		Tertutup: kasus.Tertutup()}
	for j, v := range h.Nilai {
		out.Nilai[j] = v
	}
	out.Nilai[kontrak.JalurPembuatKlaimTreaty] = kasus.PembuatID
	out.Nilai[kontrak.JalurNamaPembuatKlaimTreaty] = kasus.PembuatNama
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

// KunciKlaimTreaty - lihat `kontrak.KlaimTreatyKomite`.
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

// kolomFacRetro - properti baris `ClaimData.FacRetroList` yang tersimpan (T_CLAIM_FAC_RETRO).
func kolomFacRetro() map[string]bool {
	out := map[string]bool{}
	for _, c := range models.TabelFacRetro.Kolom {
		out[c.Properti] = true
	}
	return out
}

// periksaUbahan menolak jalur di luar daftar putih kontrak.
func periksaUbahan(u kontrak.UbahanKlaimTreaty) error {
	for j := range u.Header {
		if _, ok := kontrak.JalurHeaderKomite[j]; !ok {
			return fmt.Errorf("%w: header %q", kontrak.ErrUbahanKlaimTreatyTidakSah, j)
		}
	}
	for p := range u.Adjustment {
		if _, ok := kontrak.PropAdjustmentKomite[p]; !ok {
			return fmt.Errorf("%w: adjustment %q", kontrak.ErrUbahanKlaimTreatyTidakSah, p)
		}
	}
	kol := kolomFacRetro()
	for _, b := range u.FacRetro {
		for p := range b {
			if !kol[p] {
				return fmt.Errorf("%w: FacRetroList %q", kontrak.ErrUbahanKlaimTreatyTidakSah, p)
			}
		}
	}
	return nil
}

// TulisBalikKlaimTreaty - lihat `kontrak.KlaimTreatyKomite`.
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
	switch {
	case adjID == "" && len(u.Adjustment) > 0: // kasus komite Close Without Payment: tanpa baris adjustment
		return fmt.Errorf("%w: ubahan adjustment tanpa baris adjustment", kontrak.ErrUbahanKlaimTreatyTidakSah)
	case adjID != "" && n == 0:
		return kontrak.ErrAdjustmentTreatyTidakAda
	}
	for j, v := range u.Header {
		h.Setel(j, v)
	}
	if n > 0 {
		b := h.AmbilDaftar(models.DaftarAdjustment)[n-1]
		for p, v := range u.Adjustment {
			b[p] = v
		}
	}
	if len(u.FacRetro) > 0 && len(h.AmbilDaftar(models.DaftarFacRetro)) == 0 {
		rows := make([]models.Baris, 0, len(u.FacRetro))
		for _, r := range u.FacRetro {
			rows = append(rows, models.Baris(r))
		}
		h.SetelDaftar(models.DaftarFacRetro, rows)
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

// TutupKlaimTreaty - lihat `kontrak.KlaimTreatyKomite` (KomitePost_Close S17 pxForceCaseClose, CloseAllSubCases).
func (k KlaimUntukKomite) TutupKlaimTreaty(ctx context.Context, tx *db.Tx, klaimID, komiteID string, saat time.Time) error {
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
	if err := k.l.g.TutupKasus(ctx, tx, klaimID, kasus.Tahap, saat); err != nil {
		if errors.Is(err, repository.ErrTahapBerubah) {
			return kontrak.ErrKlaimTreatyTertutup
		}
		return err
	}
	return k.l.g.TutupKomiteTerbuka(ctx, tx, klaimID, komiteID, saat)
}
