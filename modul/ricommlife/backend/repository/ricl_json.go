package repository

// Ubah `M_RICOMM_LIFE_SUMMARY.JSONDATA` tanpa fungsi JSON Oracle versi baru.
//
// Keputusan work owner 06-10-2026 butir 6 (+ RALAT R3 riratelife): Oracle DEV menolak `JSON_MERGEPATCH` (ORA-00907).
// Pola yang sudah berjalan di DEV: baca JSONDATA `FOR UPDATE` di dalam transaksi layanan, ganti kuncinya di Go
// (TerapkanKunci - kunci Pega lain dan angka apa adanya), tulis kembali utuh `UPDATE ... SET JSONDATA = :1`.

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
)

// ErrJSONRusak - JSONDATA baris itu bukan objek JSON; baris tidak diubah.
var ErrJSONRusak = errors.New("repository: JSONDATA bukan objek JSON")

func teksJSON(v any) (json.RawMessage, error) {
	var b bytes.Buffer
	e := json.NewEncoder(&b)
	e.SetEscapeHTML(false)
	if err := e.Encode(v); err != nil {
		return nil, err
	}
	return json.RawMessage(bytes.TrimSpace(b.Bytes())), nil
}

// TerapkanKunci mengganti kunci `ubah` pada objek JSON `lama`: nilai kosong = kunci dibuang (sama dengan sisip ABSENT
// ON NULL), selain itu ditulis sebagai teks seperti kolom view. Kunci lain dipertahankan apa adanya; JSONDATA kosong /
// NULL = objek kosong.
func TerapkanKunci(lama string, ubah map[string]string) (string, error) {
	obj := map[string]json.RawMessage{}
	if t := strings.TrimSpace(lama); t != "" {
		if err := json.Unmarshal([]byte(t), &obj); err != nil || obj == nil {
			return "", ErrJSONRusak
		}
	}
	for k, v := range ubah {
		if v == "" {
			delete(obj, k)
			continue
		}
		r, err := teksJSON(v)
		if err != nil {
			return "", err
		}
		obj[k] = r
	}
	r, err := teksJSON(obj)
	if err != nil {
		return "", err
	}
	return string(r), nil
}

// TeksKunci - nilai kunci JSON sebagai teks (teks dipangkas, angka apa adanya); kunci tidak ada / bukan skalar /
// JSON rusak = "", ada=false bila kuncinya tidak ada atau null.
func TeksKunci(jsonData, kunci string) (teks string, ada bool) {
	obj := map[string]any{}
	d := json.NewDecoder(strings.NewReader(jsonData))
	d.UseNumber()
	if d.Decode(&obj) != nil {
		return "", false
	}
	switch v := obj[kunci].(type) {
	case string:
		return strings.TrimSpace(v), true
	case json.Number:
		return v.String(), true
	case nil:
		return "", false
	default:
		return fmt.Sprint(v), true
	}
}

// SqlBacaJSON - JSONDATA satu baris, dikunci sampai transaksi selesai.
func SqlBacaJSON(t string) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE ID = :1 FOR UPDATE`, KolomJSON, t)
}

// SqlTulisJSON - JSONDATA utuh hasil TerapkanKunci.
func SqlTulisJSON(t string) string {
	return fmt.Sprintf(`UPDATE %s SET %s = :1 WHERE ID = :2`, t, KolomJSON)
}

// ubahJSON - baca-ubah-tulis JSONDATA ringkasan id (nama berskema `t`). Baris tidak ada = ErrTidakAda.
func (g *Gudang) ubahJSON(ctx context.Context, tx *db.Tx, t, id string, ubah map[string]string) error {
	q := SqlBacaJSON(t)
	if err := siap(TabelRingkasan, q); err != nil {
		return err
	}
	var p db.PindaiTeksPanjang
	if err := g.dari(tx).QueryRowContext(ctx, q, id).Scan(&p); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrTidakAda
		}
		return bungkus(err, "membaca JSONDATA")
	}
	baru, err := TerapkanKunci(p.Teks(), ubah)
	if err != nil {
		return fmt.Errorf("%w (%s ID %s)", err, TabelRingkasan, id)
	}
	j, err := g.tulis(ctx, tx, TabelRingkasan, SqlTulisJSON(t), "menulis JSONDATA", baru, id)
	if err == nil && j == 0 {
		return ErrTidakAda
	}
	return err
}
