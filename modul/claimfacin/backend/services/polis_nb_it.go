package services

// Tab Loss Record / Correspondence tampilan polis NB hanya-baca (`models.PolisNBIT`; keputusan work owner 11-10-2026
// "Ikut NB: hanya akun IT"): penegakan di layanan seperti nbfacin `wajibJabatanIT` - tanpa identitas 401, bukan
// JOB_POSITION "IT" 403 (galat baca jabatan diteruskan); polis = polis halaman klaim (POLICY_NO + PRODKE), kasus NB dari
// dokumen penunjuk JSON_POLIS bila ada, selain itu dicari seperti `DokumenPolisDatar`.

import (
	"context"
	"fmt"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/claimfacin/backend/models"
)

// PolisNBIT - GET /kasus/{id}/polis-nb-it.
func (l *Layanan) PolisNBIT(ctx context.Context, p inti.Pelaku, id string) (models.PolisNBIT, error) {
	if err := l.periksaPelaku(p); err != nil {
		return models.PolisNBIT{}, err
	}
	if l.a == nil {
		return models.PolisNBIT{}, ErrTanpaOracle
	}
	jabatan, err := l.a.JabatanPelaku(ctx, p.AkunID)
	if err != nil {
		return models.PolisNBIT{}, err
	}
	if !models.JabatanIT(jabatan) {
		return models.PolisNBIT{}, fmt.Errorf("%w: tab Loss Record / Correspondence hanya untuk JOB_POSITION %s",
			inti.ErrTanpaWewenang, models.JabatanTabIT)
	}
	if _, err := l.g.Keadaan(ctx, nil, id); err != nil {
		return models.PolisNBIT{}, err
	}
	h, err := l.g.BacaHalaman(ctx, nil, id)
	if err != nil {
		return models.PolisNBIT{}, err
	}
	nopolis, prodke := h.Ambil(models.JalurNoPolis), h.Ambil(models.JalurProdke)
	kosong := models.PolisNBIT{LossRecord: models.LossRecordNB{Baris: []models.BarisLossRecordNB{}},
		Korespondensi: []models.KorespondensiNB{}}
	if nopolis == "" {
		return kosong, nil
	}
	kasus := ""
	if dok, ada, err := l.a.DokumenPolis(ctx, nopolis, prodke); err != nil {
		return models.PolisNBIT{}, err
	} else if ada && !models.DokumenPolisBerisi(dok) {
		kasus = models.KasusPenunjukPolis(dok)
	}
	return l.a.PolisNBIT(ctx, nopolis, prodke, kasus)
}
