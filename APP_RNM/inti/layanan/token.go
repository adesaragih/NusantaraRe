package layanan

// Token penyimpanan berkas - butir an, A2.
//
// Untuk apa berkas ini: logika `pooldata.GET_TOKEN_STORAGE` DITULIS ULANG DI
// GO - prinsip **o**, `[keputusan work owner]` *"jangan ada lagi pemanggilan
// procedure"*.
//
// `[data DBA — belum dikonfirmasi DBA]` bentuk procedure-nya: bila `APPNAME`
// kosong → galat; ambil `KODEAKSES` terbaru dari `GCP_IMAGE` yang
// `INPUTDATE > SYSDATE`; bila tidak ada → terbitkan token baru
// `RAWTOHEX(STANDARD_HASH('ASMAPP' ‖ <garam> ‖ <stempel waktu>, 'MD5'))`,
// simpan dengan `INPUTDATE = SYSDATE + 1 menit` dan
// `USERINPUT = NVL(masukan,'Job')`.
//
// ⛔ GARAMNYA TIDAK PERNAH DISALIN dari procedure ke mana pun - tidak ke kode,
// tidak ke tiket, tidak ke komentar ini. Ia datang dari env
// `STORAGE_TOKEN_SALT`, dan kosong berarti GAGAL TERANG. Test memakai garam
// palsu.
//
// ⚠️ MD5 dipakai karena procedure-nya memakai MD5, bukan karena ia pilihan
// yang baik. Menggantinya mengubah token yang sudah beredar, dan itu
// keputusan work owner - bukan perbaikan diam-diam oleh executor.
//
// Dibaca sesudah: efekkeluar.go.

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"nusantarare/inti"
	"nusantarare/inti/db"
)

var (
	// ErrGaramTokenKosong - `STORAGE_TOKEN_SALT` belum disetel.
	ErrGaramTokenKosong = errors.New(
		"services: STORAGE_TOKEN_SALT kosong; token penyimpanan tidak dapat diterbitkan")
	// ErrAppNameKosong - procedure menolak `APPNAME` kosong.
	ErrAppNameKosong = errors.New(
		"services: APPNAME kosong; token penyimpanan menuntutnya")
)

// awalanToken adalah literal yang procedure sertakan sebelum garamnya.
//
// `[data DBA]` `'ASMAPP'` - bukan rahasia, ia konstanta yang terbaca di badan
// procedure. Yang rahasia garamnya, dan itu TIDAK ada di sini.
const awalanToken = "ASMAPP"

// umurToken adalah masa berlaku token yang baru diterbitkan.
//
// `[data DBA]` `INPUTDATE = SYSDATE + 1 menit`.
const UmurToken = time.Minute

// penggunaTokenBawaan adalah nilai `USERINPUT` bila pemanggil tidak menyebut.
//
// `[data DBA]` `NVL(masukan, 'Job')`.
const PenggunaTokenBawaan = "Job"

// RakitToken membentuk token dari garam dan stempel waktunya.
//
// ⛔ Dipisah dari penulisannya supaya bentuknya dapat diuji tanpa Oracle DAN
// tanpa garam sungguhan.
//
// ⚠️ Stempel waktunya berformat `YYYYMMDDHH24MISSFF3` seperti procedure -
// sampai milidetik. Memangkasnya membuat dua token yang terbit dalam detik
// yang sama menjadi identik.
func RakitToken(garam string, saat time.Time) (string, error) {
	if strings.TrimSpace(garam) == "" {
		return "", ErrGaramTokenKosong
	}
	stempel := saat.Format("20060102150405.000")
	stempel = strings.Replace(stempel, ".", "", 1)
	jumlah := md5.Sum([]byte(awalanToken + garam + stempel))
	return strings.ToUpper(hex.EncodeToString(jumlah[:])), nil
}

// TokenStorage menerbitkan atau memakai ulang token penyimpanan.
type TokenStorage interface {
	Token(ctx context.Context, tx *db.Tx, appName, pengguna string,
		saat time.Time) (string, error)
}

// tokenOracle membaca dan menulis `GCP_IMAGE`.
//
// ⚠️ TULISAN KE TABEL WARISAN YANG DISENGAJA dan DISETUJUI: `GCP_IMAGE`
// adalah satu-satunya tabel warisan selain baris datar yang §6 brief izinkan
// ditulis. Dicatat, bukan disembunyikan.
type tokenOracle struct {
	pohon *PenyimpanToken
	garam string
}

// TokenStorageOracle menyusun penerbit token.
//
// ⚠️ Garamnya diserahkan pemanggil - yang membacanya `inti/config`, satu-
// satunya tempat yang boleh menyentuh env (penjaga tiket 12).
func TokenStorageOracle(svc inti.Akar, garam string) TokenStorage {
	return tokenOracle{pohon: NewPenyimpanToken(svc.DB()), garam: garam}
}

// Token memakai ulang token yang masih berlaku, atau menerbitkan yang baru.
func (t tokenOracle) Token(ctx context.Context, tx *db.Tx,
	appName, pengguna string, saat time.Time) (string, error) {

	if strings.TrimSpace(appName) == "" {
		return "", ErrAppNameKosong
	}
	if strings.TrimSpace(t.garam) == "" {
		return "", ErrGaramTokenKosong
	}
	lama, err := t.pohon.TokenBerlaku(ctx, tx, appName, saat)
	if err != nil {
		return "", err
	}
	if lama != "" {
		return lama, nil
	}
	baru, err := RakitToken(t.garam, saat)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(pengguna) == "" {
		pengguna = PenggunaTokenBawaan
	}
	if err := t.pohon.SimpanToken(ctx, tx, appName, baru, pengguna,
		saat.Add(UmurToken)); err != nil {
		return "", fmt.Errorf("services: menyimpan token penyimpanan: %w", err)
	}
	return baru, nil
}
