package services

// Untuk apa berkas ini: KONTRAK yang Claim Fac In sediakan untuk Komite Claim Fac In - `kontrak.KlaimFacInKomite`
// (prompt work owner tahap 2 10-10-2026 §2 butir 1; disalin dari pola `modul/claimnonprop/backend/services/
// kontrak_komite.go`, bukan impor). Komite membaca kasus klaim induk dan, lewat `KomitePost_Adjustment` /
// `KomitePost_Reject` / `KomitePost_CloseClaim`, menulis kembali ke kasus itu di DALAM transaksi komite. Inilah
// satu-satunya jalan modul lain menulis tabel klaim Fac In; jalur yang boleh ditulis dibatasi daftar putih di paket
// kontrak.
//
// Bacaan = `muat` (halaman + polis `OfferFacIn` dari JSON_POLIS) + `models.HitungTurunan` (medan turunan layar klaim).
// Tulisan = `BacaHalaman` -> ubahan -> `SimpanHalaman` + `SentuhKasus`, jalur yang sama dengan setiap aksi Claim Fac In;
// kolom milik komite baris adjustment (`Kolom.MilikKomite`, dilewati simpan halaman) ditulis `UbahAdjustmentKomite`.

import (
	"context"
	"errors"
	"fmt"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/claimfacin/backend/models"
	"nusantarare/modul/claimfacin/backend/repository"
)

// KlaimUntukKomite memenuhi `kontrak.KlaimFacInKomite`.
type KlaimUntukKomite struct{ l *Layanan }

// KlaimUntukKomite menyusun penyedia kontrak Komite Claim Fac In di atas layanan ini.
func (l *Layanan) KlaimUntukKomite() KlaimUntukKomite { return KlaimUntukKomite{l: l} }

var _ kontrak.KlaimFacInKomite = KlaimUntukKomite{}

func (k KlaimUntukKomite) siap() error {
	if !k.l.AdaGudang() {
		return ErrTanpaOracle
	}
	return nil
}

// galatKasusKomite menerjemahkan galat repository ke galat kontrak.
func galatKasusKomite(err error) error {
	switch {
	case errors.Is(err, repository.ErrKasusTidakAda):
		return kontrak.ErrKlaimFacInTidakAda
	case errors.Is(err, repository.ErrTahapBerubah):
		return kontrak.ErrKlaimFacInTertutup
	}
	return err
}

// posisiAdjKomite - posisi (objek, item, adjustment; 1..n) baris adjustment ber-ID `adjID`; ok false bila tidak ada.
func posisiAdjKomite(h *models.Halaman, adjID string) (o, i, a int, ok bool) {
	if adjID == "" {
		return 0, 0, 0, false
	}
	for on := range h.AmbilDaftar(models.DaftarObjek) {
		for in := range h.AmbilDaftar(models.DaftarItem(on + 1)) {
			for an, b := range h.AmbilDaftar(models.DaftarAdj(on+1, in+1)) {
				if b[models.PropID] == adjID {
					return on + 1, in + 1, an + 1, true
				}
			}
		}
	}
	return 0, 0, 0, false
}

// BacaKlaimFacIn - lihat `kontrak.KlaimFacInKomite`.
func (k KlaimUntukKomite) BacaKlaimFacIn(ctx context.Context, tx *db.Tx, klaimID, adjID string) (kontrak.KlaimFacIn,
	error) {
	if err := k.siap(); err != nil {
		return kontrak.KlaimFacIn{}, err
	}
	kasus, h, err := k.l.muat(ctx, tx, klaimID)
	if err != nil {
		return kontrak.KlaimFacIn{}, galatKasusKomite(err)
	}
	models.HitungTurunan(h)
	out := kontrak.KlaimFacIn{Nilai: map[string]string{}, Daftar: map[string][]map[string]string{},
		Tertutup: kasus.Tertutup()}
	if adjID != "" {
		o, i, a, ok := posisiAdjKomite(h, adjID)
		if !ok {
			return kontrak.KlaimFacIn{}, kontrak.ErrAdjustmentFacInTidakAda
		}
		out.Objek, out.Item, out.Adjustment = o, i, a
	}
	for j, v := range h.Nilai {
		out.Nilai[j] = v
	}
	// pembuat kasus klaim (`TempOpenPage.pxCreateOperator` / `pxCreateOpName`, KomitePost_Reject S14.1 CARI6-CARI7)
	out.Nilai[kontrak.JalurPembuatKlaim], out.Nilai[kontrak.JalurNamaPembuatKlaim] = kasus.PembuatID, kasus.PembuatNama
	// `ClaimData.ClaimNo := pyID` (CallActivityInputRegister 9, pra-proses register - tidak disimpan kolom sendiri);
	// dibaca ShowTransfer LS50 komite.
	if out.Nilai[models.CD+"ClaimNo"] == "" {
		out.Nilai[models.CD+"ClaimNo"] = kasus.ID
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

// KunciKlaimFacIn - lihat `kontrak.KlaimFacInKomite`.
func (k KlaimUntukKomite) KunciKlaimFacIn(ctx context.Context, tx *db.Tx, klaimID string) error {
	if err := k.siap(); err != nil {
		return err
	}
	kasus, err := k.l.g.Keadaan(ctx, tx, klaimID)
	if err != nil {
		return galatKasusKomite(err)
	}
	if kasus.Tertutup() {
		return kontrak.ErrKlaimFacInTertutup
	}
	if _, err := k.l.g.KunciKasus(ctx, tx, klaimID, kasus.Tahap); err != nil {
		return galatKasusKomite(err)
	}
	return nil
}

// periksaUbahanKomite menolak jalur di luar daftar putih kontrak.
func periksaUbahanKomite(u kontrak.UbahanKlaimFacIn) error {
	for _, x := range []struct {
		nama  string
		ubah  map[string]string
		putih map[string]string
	}{{"header", u.Header, kontrak.JalurHeaderKomiteFacIn}, {"adjustment", u.Adjustment, kontrak.PropAdjustmentKomiteFacIn},
		{"objek", u.Objek, kontrak.PropObjekKomiteFacIn}, {"item", u.Item, kontrak.PropItemKomiteFacIn}} {
		for j := range x.ubah {
			if _, ok := x.putih[j]; !ok {
				return fmt.Errorf("%w: %s %q", kontrak.ErrUbahanKlaimFacInTidakSah, x.nama, j)
			}
		}
	}
	return nil
}

// kolomMilikKomite - properti baris adjustment yang hanya ditulis komite (dilewati simpan halaman).
func kolomMilikKomite(p string) bool {
	for _, kol := range models.TabelAdjustment.Kolom {
		if kol.Properti == p {
			return kol.MilikKomite
		}
	}
	return false
}

// setelBaris menulis ubahan ke satu baris (nilai kosong = properti dibuang, seperti properti Pega yang dikosongkan).
func setelBaris(b models.Baris, u map[string]string) {
	for p, v := range u {
		if v == "" {
			delete(b, p)
			continue
		}
		b[p] = v
	}
}

// TulisBalikKlaimFacIn - lihat `kontrak.KlaimFacInKomite`.
func (k KlaimUntukKomite) TulisBalikKlaimFacIn(ctx context.Context, tx *db.Tx, klaimID, adjID string,
	u kontrak.UbahanKlaimFacIn) error {
	if err := k.siap(); err != nil {
		return err
	}
	if err := periksaUbahanKomite(u); err != nil {
		return err
	}
	kasus, err := k.l.g.Keadaan(ctx, tx, klaimID)
	if err != nil {
		return galatKasusKomite(err)
	}
	if kasus.Tertutup() {
		return kontrak.ErrKlaimFacInTertutup
	}
	h, err := k.l.g.BacaHalaman(ctx, tx, klaimID)
	if err != nil {
		return galatKasusKomite(err)
	}
	komite := map[string]string{}
	if len(u.Adjustment) > 0 || len(u.Objek) > 0 || len(u.Item) > 0 {
		o, i, a, ok := posisiAdjKomite(h, adjID)
		if !ok {
			return kontrak.ErrAdjustmentFacInTidakAda
		}
		ob, err := models.Objek(h, o)
		if err != nil {
			return err
		}
		it, err := models.Item(h, o, i)
		if err != nil {
			return err
		}
		b, err := models.Adj(h, o, i, a)
		if err != nil {
			return err
		}
		setelBaris(ob, u.Objek)
		setelBaris(it, u.Item)
		biasa := map[string]string{}
		for p, v := range u.Adjustment {
			if kolomMilikKomite(p) {
				komite[p] = v
				continue
			}
			biasa[p] = v
		}
		setelBaris(b, biasa)
	}
	for j, v := range u.Header {
		h.Setel(j, v)
	}
	for _, r := range u.Riwayat { // ChronologyInsertion_DT: langkah tanpa label = "Committee"
		kt := &models.Konteks{Ctx: ctx, Pelaku: r.Pelaku, Tingkat: r.Tingkat, Sekarang: r.Saat}
		kt.Kronologi(h, r.Teks)
	}
	if err := k.l.g.SimpanHalaman(ctx, tx, klaimID, h); err != nil {
		return err
	}
	if len(komite) > 0 {
		if err := k.l.g.UbahAdjustmentKomite(ctx, tx, klaimID, adjID, komite); err != nil {
			return err
		}
	}
	return k.l.g.SentuhKasus(ctx, tx, klaimID, k.l.jam())
}

// TutupKlaimFacIn - lihat `kontrak.KlaimFacInKomite`.
func (k KlaimUntukKomite) TutupKlaimFacIn(ctx context.Context, tx *db.Tx, klaimID, komiteID, status string,
	saat time.Time) error {
	if err := k.siap(); err != nil {
		return err
	}
	if status != kontrak.StatusKlaimDitolak && status != kontrak.StatusKlaimSelesai {
		return fmt.Errorf("%w: status %q", kontrak.ErrUbahanKlaimFacInTidakSah, status)
	}
	kasus, err := k.l.g.Keadaan(ctx, tx, klaimID)
	if err != nil {
		return galatKasusKomite(err)
	}
	if kasus.Tertutup() {
		return kontrak.ErrKlaimFacInTertutup
	}
	if err := k.l.g.TutupKasusStatus(ctx, tx, klaimID, kasus.Tahap, status, saat); err != nil {
		return galatKasusKomite(err)
	}
	return k.l.g.TutupKomiteAnak(ctx, tx, klaimID, komiteID, status, saat) // CloseAllSubCases=true
}
