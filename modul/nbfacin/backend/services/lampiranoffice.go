package services

// View Office Online lampiran NB FacIn (tiket 67, permintaan work owner 11-10-2026 "seperti panel lampiran bersama / Product Name
// Life"): `ReasViewAttachment` b2951 hanya untuk xls / xlsx / doc / docx / ppt / pptx; `DownloadDocumentPolis` 3 ->
// GetUrlGoogleStorage_Act (URL bertanda tangan berumur DurasiURLStorageLampiran) - pola `inti/backend/dokumenpolis` TautanOffice.
// ⛔ URL bertanda tangan TIDAK pernah ditulis ke log atau pesan galat.

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/nbfacin/backend/models"
)

var (
	// ErrLampiranBukanOffice - ekstensi di luar penampil kantor (400).
	ErrLampiranBukanOffice = fmt.Errorf("%w: View Office Online hanya untuk berkas Excel, Word, dan PowerPoint "+
		"(xls, xlsx, doc, docx, ppt, pptx)", ErrMasukanLampiran)
	// ErrLampiranOfficeTanpaStorage - berkas ada di penyimpanan lokal (stub, URLPUBLIC kosong): penampil online butuh URL publik
	// penyimpanan nyata (409).
	ErrLampiranOfficeTanpaStorage = errors.New("services: View Office Online butuh penyimpanan nyata (Google Storage) - berkas ini " +
		"tersimpan di penyimpanan lokal sistem baru; gunakan View / unduh")
)

// ekstensiOfficeLampiran - `ReasViewAttachment` b2951 (sama dengan dokumenpolis.ekstensiOffice).
var ekstensiOfficeLampiran = map[string]bool{"xls": true, "xlsx": true, "doc": true, "docx": true, "ppt": true, "pptx": true}

// penautStorage - klien penyimpanan yang dapat memberi URL bertanda tangan (penyimpanan nyata saja; stub lokal tidak).
type penautStorage interface {
	Tautan(ctx context.Context, link models.LinkStorage) (string, error)
}

// TautanOfficeLampiran - GET …/lampiran/berkas/{id}/office: otorisasi / tahap sama dengan /isi (UnduhLampiran); ekstensi kantor
// (400); belum tersimpan (409 ErrLampiranBelumTersimpan); stub lokal (409 ErrLampiranOfficeTanpaStorage); objek hilang (409
// ErrLampiranTanpaBerkas).
func (s *Service) TautanOfficeLampiran(ctx context.Context, pelaku inti.Pelaku, id, dokID string) (string, error) {
	if _, err := s.tahapLampiran(ctx, pelaku, id); err != nil {
		return "", err
	}
	d, err := s.ambilDokumen(ctx, id, dokID)
	if err != nil {
		return "", err
	}
	if !ekstensiOfficeLampiran[strings.ToLower(strings.TrimSpace(d.Mime))] {
		return "", ErrLampiranBukanOffice
	}
	link, ada, err := s.lampiran.LinkStorage(ctx, nil, d.TStorageID)
	if err != nil {
		return "", err
	}
	if !ada {
		return "", fmt.Errorf("%w: dokumen %s", ErrLampiranBelumTersimpan, dokID)
	}
	p, nyata := s.storageLampiran.(penautStorage)
	if !nyata || strings.TrimSpace(link.URLPublic) == "" {
		return "", ErrLampiranOfficeTanpaStorage
	}
	u, err := p.Tautan(ctx, link)
	if errors.Is(err, ErrBerkasTidakAdaDiPenyimpanan) {
		return "", fmt.Errorf("%w: dokumen %s", ErrLampiranTanpaBerkas, dokID)
	}
	return u, err
}

// Tautan - URL bertanda tangan https objek tercatat (GetUrlGoogleStorage_Act seperti Buka, tanpa mengunduh isinya): EXPDATE kosong
// = tanpa objek; URL tersimpan dipakai selama berlaku; selain itu geturl + Update_T_Storage_SQL.
func (p *storageJarakJauh) Tautan(ctx context.Context, link models.LinkStorage) (string, error) {
	if strings.TrimSpace(link.ExpUTC) == "" {
		return "", ErrBerkasTidakAdaDiPenyimpanan
	}
	if models.LinkBerlaku(link, p.jam(), sisaURLTersimpan) {
		u, err := url.Parse(strings.TrimSpace(link.URLPublic))
		if err == nil && u.Host != "" && u.Scheme == "https" {
			return u.String(), nil
		}
	}
	b, objek, err := p.urlBertanda(ctx, models.LokasiDariLink(link))
	if err != nil {
		return "", err
	}
	p.segarkan(ctx, objek)
	return b, nil
}
