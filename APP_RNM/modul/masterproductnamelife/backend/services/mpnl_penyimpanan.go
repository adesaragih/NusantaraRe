package services

// Penyimpanan berkas STUB (P5) - BAWAAN selama `PELAKSANA_STORAGE` bukan `nyata`; penyimpanan nyata
// (`ServiceGoogle` + `LinkService` + token) ada di `mpnl_storage.go` (keputusan work owner 03-10-2026). Berkas
// ditahan di folder lokal `UNGGAHAN_DIR/master-product-name-life/antre/<IMAGEID>`; "kirim" memindahnya ke
// `…/simpan/<IMAGEID>` (meniru objek di penyimpanan) dan objeknya dicatat apa adanya (URLPUBLIC kosong). ⛔ Nol
// klien HTTP di berkas ini.

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"nusantarare/inti/backend/unggah"
	"nusantarare/modul/masterproductnamelife/backend/models"
)

const folderStub = "master-product-name-life"

// polaImageID - kunci objek aman sebagai nama berkas (MD5 hex `GenerateImageID_SQL`).
var polaImageID = regexp.MustCompile(`^[0-9A-Za-z]{1,64}$`)

var errImageIDTidakSah = errors.New("services: storage ID is not a valid object key")

// PenyimpananLokal menyusun penyimpanan stub di bawah folder unggahan;
// folder kosong = setiap operasi gagal terang (503).
func PenyimpananLokal(dirUnggahan string) PenyimpananBerkas {
	if strings.TrimSpace(dirUnggahan) == "" {
		return penyimpananBelumDisetel{}
	}
	return penyimpananLokal{akar: filepath.Join(dirUnggahan, folderStub)}
}

type penyimpananLokal struct{ akar string }

func (p penyimpananLokal) jalur(bagian, imageID string) (string, error) {
	if !polaImageID.MatchString(imageID) {
		return "", fmt.Errorf("%w: %q", errImageIDTidakSah, imageID)
	}
	return filepath.Join(p.akar, bagian, imageID), nil
}

func (p penyimpananLokal) SimpanAntrean(_ context.Context, imageID string, isi io.Reader) error {
	j, err := p.jalur("antre", imageID)
	if err != nil {
		return err
	}
	return unggah.TulisBerkas(j, isi, unggah.BatasUkuranUnggahan)
}

func (p penyimpananLokal) BuangAntrean(_ context.Context, imageID string) {
	if j, err := p.jalur("antre", imageID); err == nil {
		_ = os.Remove(j)
	}
}

func (p penyimpananLokal) Kirim(_ context.Context, o models.ObjekPenyimpanan, _, _ string) (models.ObjekPenyimpanan, error) {
	sumber, err := p.jalur("antre", o.ImageID)
	if err != nil {
		return models.ObjekPenyimpanan{}, err
	}
	tujuan, _ := p.jalur("simpan", o.ImageID)
	if _, err := os.Stat(tujuan); err == nil {
		return o, nil // sudah terkirim - idempoten
	}
	if _, err := os.Stat(sumber); errors.Is(err, os.ErrNotExist) {
		return models.ObjekPenyimpanan{}, ErrBerkasSumberHilang
	}
	if err := os.MkdirAll(filepath.Dir(tujuan), 0o750); err != nil {
		return models.ObjekPenyimpanan{}, fmt.Errorf("services: preparing the storage stub folder: %w", err)
	}
	if err := os.Rename(sumber, tujuan); err != nil {
		return models.ObjekPenyimpanan{}, fmt.Errorf("services: sending to the storage stub: %w", err)
	}
	return o, nil
}

func (p penyimpananLokal) Buka(_ context.Context, o models.ObjekPenyimpanan) (io.ReadCloser, *models.ObjekPenyimpanan, error) {
	j, err := p.jalur("simpan", o.ImageID)
	if err != nil {
		return nil, nil, err
	}
	f, err := os.Open(j)
	if err != nil {
		return nil, nil, err
	}
	return f, nil, nil
}

func (p penyimpananLokal) Hapus(_ context.Context, imageID string, o *models.ObjekPenyimpanan) error {
	if o != nil && strings.TrimSpace(o.URLPublic) != "" {
		// Objek di penyimpanan nyata (unggahan Pega, atau mode nyata sebelumnya): stub tidak dapat menghapusnya, dan
		// menghapus rekamnya saja meninggalkan objek yatim.
		return galatStorage{jenis: ErrStorageBelumSiap, layar: ErrStorageBelumSiap.Error() + ": the attachment file is in " +
			"storage; deleting it needs the storage service (PELAKSANA_STORAGE=nyata)"}
	}
	for _, bagian := range []string{"simpan", "antre"} {
		j, err := p.jalur(bagian, imageID)
		if err != nil {
			return err
		}
		if err := os.Remove(j); err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("services: deleting from the storage stub: %w", err)
		}
	}
	return nil
}

// penyimpananBelumDisetel - `UNGGAHAN_DIR` kosong: gagal terang, bukan folder kerja diam-diam.
type penyimpananBelumDisetel struct{}

func (penyimpananBelumDisetel) SimpanAntrean(context.Context, string, io.Reader) error {
	return ErrPenyimpananBelumDisetel
}
func (penyimpananBelumDisetel) BuangAntrean(context.Context, string) {}
func (penyimpananBelumDisetel) Kirim(context.Context, models.ObjekPenyimpanan, string, string) (models.ObjekPenyimpanan, error) {
	return models.ObjekPenyimpanan{}, ErrPenyimpananBelumDisetel
}
func (penyimpananBelumDisetel) Buka(context.Context, models.ObjekPenyimpanan) (io.ReadCloser, *models.ObjekPenyimpanan, error) {
	return nil, nil, ErrPenyimpananBelumDisetel
}
func (penyimpananBelumDisetel) Hapus(context.Context, string, *models.ObjekPenyimpanan) error {
	return ErrPenyimpananBelumDisetel
}
