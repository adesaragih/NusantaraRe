package repository

// Ubah `JSONDATA` tanpa fungsi JSON Oracle versi baru.
//
// RALAT R3 (06-10-2026): `JSON_MERGEPATCH(... RETURNING CLOB)` ditolak Oracle DEV - log server
// `repository: mengubah ringkasan: ORA-00907: missing right parenthesis` saat Save Rate Detail (fungsi itu baru ada
// sejak 19c; `JSON_OBJECT` untuk sisip berjalan). Pola penggantinya = pola yang sudah berjalan di DEV
// (`masterproductnamelife`: `UPDATE ... SET JSONDATA = :1`): baca JSONDATA `FOR UPDATE` di dalam transaksi layanan,
// ganti kuncinya di Go (TerapkanKunci - kunci Pega lain dan angka apa adanya), tulis kembali utuh.

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

// TerapkanKunci mengganti kunci `ubah` pada objek JSON `lama`: nilai kosong = kunci dibuang (sama dengan sisip ABSENT
// ON NULL), selain itu ditulis sebagai teks seperti kolom view. Kunci lain dipertahankan apa adanya (angka tidak
// diubah bentuknya); JSONDATA kosong / NULL = objek kosong.
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
		var b bytes.Buffer
		e := json.NewEncoder(&b)
		e.SetEscapeHTML(false)
		if err := e.Encode(v); err != nil {
			return "", err
		}
		obj[k] = json.RawMessage(bytes.TrimSpace(b.Bytes()))
	}
	var b bytes.Buffer
	e := json.NewEncoder(&b)
	e.SetEscapeHTML(false)
	if err := e.Encode(obj); err != nil {
		return "", err
	}
	return strings.TrimSpace(b.String()), nil
}

// TeksKunci - nilai kunci JSON sebagai teks (angka atau teks); kunci tidak ada / bukan skalar = "".
func TeksKunci(jsonData, kunci string) string {
	obj := map[string]any{}
	d := json.NewDecoder(strings.NewReader(jsonData))
	d.UseNumber()
	if d.Decode(&obj) != nil {
		return ""
	}
	switch v := obj[kunci].(type) {
	case string:
		return strings.TrimSpace(v)
	case json.Number:
		return v.String()
	}
	return ""
}

// SqlBacaJSON - JSONDATA satu baris, dikunci sampai transaksi selesai.
func SqlBacaJSON(t string) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE ID = :1 FOR UPDATE`, KolomJSON, t)
}

// SqlTulisJSON - JSONDATA utuh hasil TerapkanKunci.
func SqlTulisJSON(t string) string {
	return fmt.Sprintf(`UPDATE %s SET %s = :1 WHERE ID = :2`, t, KolomJSON)
}

// SqlIDRateMilik - ID baris rate milik satu ringkasan (lewat view, seperti hapus).
func SqlIDRateMilik(v string) string { return fmt.Sprintf(`SELECT ID FROM %s WHERE IDUSEDBY = :1`, v) }

// ubahJSON - baca-ubah-tulis JSONDATA baris id tabel `objek` (nama berskema `t`). `syarat` (boleh nil) memeriksa
// JSON lama; false = ErrTidakAda. Baris tidak ada = ErrTidakAda.
func (g *Gudang) ubahJSON(ctx context.Context, tx *db.Tx, objek, t, id string, ubah map[string]string,
	syarat func(lama string) bool) error {
	q := SqlBacaJSON(t)
	if err := siap(objek, q); err != nil {
		return err
	}
	var p db.PindaiTeksPanjang
	if err := g.dari(tx).QueryRowContext(ctx, q, id).Scan(&p); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrTidakAda
		}
		return bungkus(err, "membaca JSONDATA")
	}
	if syarat != nil && !syarat(p.Teks()) {
		return ErrTidakAda
	}
	baru, err := TerapkanKunci(p.Teks(), ubah)
	if err != nil {
		return fmt.Errorf("%w (%s ID %s)", err, objek, id)
	}
	j, err := g.tulis(ctx, tx, objek, SqlTulisJSON(t), "menulis JSONDATA", baru, id)
	if err == nil && j == 0 {
		return ErrTidakAda
	}
	return err
}
