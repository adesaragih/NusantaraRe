package loader

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// jenisSimpul - bentuk satu simpul dokumen penawaran.
type jenisSimpul uint8

const (
	simpulObjek jenisSimpul = iota + 1
	simpulLarik
	simpulNilai
)

// simpul - pohon dokumen yang MEMPERTAHANKAN urutan kunci objek. `encoding/json` ke
// map akan mengacaknya, padahal spec 11 menuntut "urutan dipertahankan" dan SEQ_NO
// adalah posisi baris saat muat (V-41).
type simpul struct {
	jenis jenisSimpul
	kunci []string  // objek: nama medan, urut dokumen
	anak  []*simpul // objek: nilai per kunci; larik: unsur
	teks  string    // nilai: teks apa adanya (angka JSON tidak dibulatkan; null = "")
}

// medan - anak objek bernama k, atau nil.
func (s *simpul) medan(k string) *simpul {
	if s == nil || s.jenis != simpulObjek {
		return nil
	}
	for i, n := range s.kunci {
		if n == k {
			return s.anak[i]
		}
	}
	return nil
}

// galatBentukDokumen - galat bentuk yang aman dikutip: menyebut nama medan,
// tidak pernah nilainya.
type galatBentukDokumen string

func (g galatBentukDokumen) Error() string { return string(g) }

// uraiJSON - DATA_JSON menjadi pohon berurutan. Akar wajib objek; kunci ganda di
// satu objek ditolak (nilai mana yang berlaku tidak boleh ditebak).
func uraiJSON(b []byte) (*simpul, error) {
	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()
	s, err := uraiSimpul(d)
	if err != nil {
		return nil, err
	}
	if _, err := d.Token(); !errors.Is(err, io.EOF) {
		return nil, galatBentukDokumen("isi sesudah akar dokumen")
	}
	if s.jenis != simpulObjek {
		return nil, galatBentukDokumen("akar dokumen bukan objek")
	}
	return s, nil
}

func uraiSimpul(d *json.Decoder) (*simpul, error) {
	t, err := d.Token()
	if err != nil {
		return nil, err
	}
	switch v := t.(type) {
	case json.Delim:
		switch v {
		case '{':
			s := &simpul{jenis: simpulObjek}
			ada := map[string]bool{}
			for d.More() {
				kt, err := d.Token()
				if err != nil {
					return nil, err
				}
				k := kt.(string)
				if ada[k] {
					return nil, galatBentukDokumen(fmt.Sprintf("kunci %q ganda di satu objek", k))
				}
				ada[k] = true
				a, err := uraiSimpul(d)
				if err != nil {
					return nil, err
				}
				s.kunci, s.anak = append(s.kunci, k), append(s.anak, a)
			}
			_, err := d.Token()
			return s, err
		case '[':
			s := &simpul{jenis: simpulLarik}
			for d.More() {
				a, err := uraiSimpul(d)
				if err != nil {
					return nil, err
				}
				s.anak = append(s.anak, a)
			}
			_, err := d.Token()
			return s, err
		}
		return nil, galatBentukDokumen(fmt.Sprintf("pembatas %q tak terduga", v))
	case string:
		return &simpul{jenis: simpulNilai, teks: v}, nil
	case json.Number:
		return &simpul{jenis: simpulNilai, teks: v.String()}, nil
	case bool:
		return &simpul{jenis: simpulNilai, teks: fmt.Sprint(v)}, nil
	case nil:
		return &simpul{jenis: simpulNilai}, nil
	}
	return nil, galatBentukDokumen(fmt.Sprintf("token %T tak terduga", t))
}
