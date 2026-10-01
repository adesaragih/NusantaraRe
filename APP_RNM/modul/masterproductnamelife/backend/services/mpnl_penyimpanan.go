package services

// Penyimpanan berkas STUB (P5, OQ-MPNL-10) - pengganti `ServiceGoogle` +
// `LinkService` + `GetTokenStorage_SQL`. Berkas ditahan di folder lokal
// `UNGGAHAN_DIR/master-product-name-life/antre/<IMAGEID>`; "kirim" memindahnya
// ke `…/simpan/<IMAGEID>` (meniru objek di penyimpanan). ⛔ Nol klien HTTP, nol
// alamat layanan di kode maupun berkas apa pun.

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

func (p penyimpananLokal) Kirim(_ context.Context, imageID string) error {
	sumber, err := p.jalur("antre", imageID)
	if err != nil {
		return err
	}
	tujuan, _ := p.jalur("simpan", imageID)
	if _, err := os.Stat(tujuan); err == nil {
		return nil // sudah terkirim - idempoten
	}
	if _, err := os.Stat(sumber); errors.Is(err, os.ErrNotExist) {
		return ErrBerkasSumberHilang
	}
	if err := os.MkdirAll(filepath.Dir(tujuan), 0o750); err != nil {
		return fmt.Errorf("services: preparing the storage stub folder: %w", err)
	}
	if err := os.Rename(sumber, tujuan); err != nil {
		return fmt.Errorf("services: sending to the storage stub: %w", err)
	}
	return nil
}

func (p penyimpananLokal) Buka(_ context.Context, imageID string) (io.ReadCloser, error) {
	j, err := p.jalur("simpan", imageID)
	if err != nil {
		return nil, err
	}
	return os.Open(j)
}

func (p penyimpananLokal) Hapus(_ context.Context, imageID string) error {
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
func (penyimpananBelumDisetel) Kirim(context.Context, string) error {
	return ErrPenyimpananBelumDisetel
}
func (penyimpananBelumDisetel) Buka(context.Context, string) (io.ReadCloser, error) {
	return nil, ErrPenyimpananBelumDisetel
}
func (penyimpananBelumDisetel) Hapus(context.Context, string) error {
	return ErrPenyimpananBelumDisetel
}
