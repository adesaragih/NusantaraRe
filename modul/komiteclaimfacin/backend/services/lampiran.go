package services

// Untuk apa berkas ini: LAMPIRAN KLAIM INDUK di layar komite (perintah work owner 11-10-2026) - dibaca lewat
// `kontrak.KlaimFacInKomite`, hanya-baca. Siapa pun yang boleh membuka kasus komite (`BukaKasus`) boleh melihat dan
// mengunduh; unggah tetap di layar Claim Fac In (tombol attach komite nonaktif).

import (
	"context"
	"errors"
	"fmt"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/komiteclaimfacin/backend/models"
)

// PesanBukanOffice - `ReasViewAttachment` b2951: View Office Online hanya Excel / Word / PowerPoint.
const PesanBukanOffice = "View Office Online is only for Excel, Word, and PowerPoint files"

// kasusLampiran - kepala kasus komite `id` (tanpa tangga / klaim induk) untuk rute lampiran.
func (l *Layanan) kasusLampiran(ctx context.Context, p inti.Pelaku, id string) (models.Kasus, error) {
	if err := l.siap(p); err != nil {
		return models.Kasus{}, err
	}
	if !strings.HasPrefix(id, models.AwalanKomite) {
		return models.Kasus{}, fmt.Errorf("%w: %q", ErrKasusTidakAda, id)
	}
	return l.g.BacaKasus(ctx, nil, id, false)
}

// galatLampiran menerjemahkan galat lampiran kontrak.
func galatLampiran(err error) error {
	switch {
	case errors.Is(err, kontrak.ErrLampiranKlaimTidakAda):
		return fmt.Errorf("%w: %v", ErrKasusTidakAda, err)
	case errors.Is(err, kontrak.ErrLampiranKlaimBukanOffice):
		return &GalatValidasi{Pesan: []string{PesanBukanOffice}}
	}
	return galatKontrak(err)
}

// DaftarLampiran - grid kategori (hanya yang berisi) dan dokumen klaim induk kasus komite `id`.
func (l *Layanan) DaftarLampiran(ctx context.Context, p inti.Pelaku, id string) (models.LampiranKomite, error) {
	k, err := l.kasusLampiran(ctx, p, id)
	if err != nil {
		return models.LampiranKomite{}, err
	}
	d, err := l.klaim.LampiranKlaimFacIn(ctx, k.KlaimID)
	if err != nil {
		return models.LampiranKomite{}, galatLampiran(err)
	}
	return models.SusunLampiranKomite(d), nil
}

// UnduhLampiran - File / View: isi satu dokumen klaim induk (pemanggil menutup `Isi`).
func (l *Layanan) UnduhLampiran(ctx context.Context, p inti.Pelaku, id, lid string) (kontrak.IsiLampiranKlaim, error) {
	k, err := l.kasusLampiran(ctx, p, id)
	if err != nil {
		return kontrak.IsiLampiranKlaim{}, err
	}
	f, err := l.klaim.BukaLampiranKlaimFacIn(ctx, k.KlaimID, lid, p.AkunID)
	if err != nil {
		return kontrak.IsiLampiranKlaim{}, galatLampiran(err)
	}
	return f, nil
}

// TautanOfficeLampiran - View Office Online satu dokumen klaim induk.
func (l *Layanan) TautanOfficeLampiran(ctx context.Context, p inti.Pelaku, id, lid string) (string, error) {
	k, err := l.kasusLampiran(ctx, p, id)
	if err != nil {
		return "", err
	}
	url, err := l.klaim.TautanOfficeLampiranKlaimFacIn(ctx, k.KlaimID, lid, p.AkunID)
	if err != nil {
		return "", galatLampiran(err)
	}
	return url, nil
}
