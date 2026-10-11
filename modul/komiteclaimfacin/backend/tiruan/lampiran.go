package tiruan

// Lampiran kasus klaim induk dalam memori - bagian lampiran `kontrak.KlaimFacInKomite` (perintah work owner
// 11-10-2026). Perilaku galat sama dengan penyedia asli (`claimfacin/backend/services/kontrak_komite_lampiran.go`).

import (
	"context"
	"io"
	"strings"

	"nusantarare/inti/backend/kontrak"
	"nusantarare/inti/backend/penyimpanan"
	"nusantarare/inti/backend/unggah"
)

// AwalanTautanOffice - awalan tautan View Office Online tiruan (diikuti ID dokumen).
const AwalanTautanOffice = "UJI-OFFICE-"

// ekstensiOffice - berkas yang boleh View Office Online (sama dengan penyedia asli).
var ekstensiOffice = map[string]bool{"xls": true, "xlsx": true, "doc": true, "docx": true, "ppt": true, "pptx": true}

// SetelLampiran menyimpan kategori + dokumen kasus klaim `klaimID` dan isi berkas per ID dokumen.
func (p *KlaimPalsu) SetelLampiran(klaimID string, d kontrak.DaftarLampiranKlaim, isi map[string]string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.lampiran == nil {
		p.lampiran, p.isiLampiran = map[string]kontrak.DaftarLampiranKlaim{}, map[string]string{}
	}
	p.lampiran[klaimID] = kontrak.DaftarLampiranKlaim{Kategori: append([]kontrak.KategoriLampiranKlaim{}, d.Kategori...),
		Lampiran: append([]kontrak.LampiranKlaim{}, d.Lampiran...)}
	for id, s := range isi {
		p.isiLampiran[id] = s
	}
}

// LampiranKlaimFacIn - lihat `kontrak.KlaimFacInKomite`.
func (p *KlaimPalsu) LampiranKlaimFacIn(_ context.Context, klaimID string) (kontrak.DaftarLampiranKlaim, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.klaim[klaimID]; !ok {
		return kontrak.DaftarLampiranKlaim{}, kontrak.ErrKlaimFacInTidakAda
	}
	d := p.lampiran[klaimID]
	return kontrak.DaftarLampiranKlaim{Kategori: append([]kontrak.KategoriLampiranKlaim{}, d.Kategori...),
		Lampiran: append([]kontrak.LampiranKlaim{}, d.Lampiran...)}, nil
}

func (p *KlaimPalsu) satuLampiran(klaimID, lampiranID string) (kontrak.LampiranKlaim, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.klaim[klaimID]; !ok {
		return kontrak.LampiranKlaim{}, kontrak.ErrKlaimFacInTidakAda
	}
	for _, a := range p.lampiran[klaimID].Lampiran {
		if a.ID == strings.TrimSpace(lampiranID) {
			return a, nil
		}
	}
	return kontrak.LampiranKlaim{}, kontrak.ErrLampiranKlaimTidakAda
}

// BukaLampiranKlaimFacIn - lihat `kontrak.KlaimFacInKomite`.
func (p *KlaimPalsu) BukaLampiranKlaimFacIn(_ context.Context, klaimID, lampiranID, _ string) (kontrak.IsiLampiranKlaim,
	error) {
	a, err := p.satuLampiran(klaimID, lampiranID)
	if err != nil {
		return kontrak.IsiLampiranKlaim{}, err
	}
	if !a.AdaObjek {
		return kontrak.IsiLampiranKlaim{}, penyimpanan.ErrObjekTidakAda
	}
	p.mu.Lock()
	isi := p.isiLampiran[a.ID]
	p.mu.Unlock()
	return kontrak.IsiLampiranKlaim{Isi: io.NopCloser(strings.NewReader(isi)), Nama: a.NamaFile,
		Mime: unggah.MimeDariNamaFile(a.NamaFile)}, nil
}

// TautanOfficeLampiranKlaimFacIn - lihat `kontrak.KlaimFacInKomite`.
func (p *KlaimPalsu) TautanOfficeLampiranKlaimFacIn(_ context.Context, klaimID, lampiranID, _ string) (string, error) {
	a, err := p.satuLampiran(klaimID, lampiranID)
	if err != nil {
		return "", err
	}
	if !ekstensiOffice[strings.ToLower(a.MIME)] {
		return "", kontrak.ErrLampiranKlaimBukanOffice
	}
	if !a.AdaObjek {
		return "", penyimpanan.ErrObjekTidakAda
	}
	return AwalanTautanOffice + a.ID, nil
}
