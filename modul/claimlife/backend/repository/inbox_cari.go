package repository

// Kotak cari inbox Claim Life - keputusan work owner 12-10-2026 ("di inbox claim life tambahkan juga filter untuk mencari
// seperti claim facin"). Dipisah dari inbox.go karena menyentuh tabel sumber peserta (namaTabelPeserta, 66,8 juta
// baris): setiap kueri di sini berindeks dan berbatas (penjaga TestQueryTabelPesertaSelaluBerindexDanBerbatas).
//
// Dibaca sesudah: inbox.go, isianoutstanding.go (sqlTertanggungPeserta).

import (
	"fmt"
	"strings"
)

// Panjang kata kunci kotak cari inbox (rune, sesudah dipangkas).
const (
	MinCariInbox  = 2
	MaksCariInbox = 100
)

// PolaCariInbox mengubah kata kunci menjadi pola LIKE tak peka huruf: dipangkas, HURUF BESAR, karakter khusus LIKE
// (backslash, %, _) diloloskan dengan backslash (dipasangkan ESCAPE '\'), lalu dibungkus %...%. Kosong tetap kosong.
func PolaCariInbox(kata string) string {
	k := strings.TrimSpace(kata)
	if k == "" {
		return ""
	}
	k = strings.NewReplacer(`\`, `\\`, "%", `\%`, "_", `\_`).Replace(strings.ToUpper(k))
	return "%" + k + "%"
}

// sqlInboxCari - saringan kotak cari (12-10-2026): Case ID, CLAIM_NO, PL Number peserta aktif, ATAU nama tertanggung
// peserta aktif. Empat bind (pola yang SAMA), berurutan.
//
// ⛔ Nama TIDAK disimpan di tabel klaim (kolomSalin tanpa nama, data pribadi): ia dibaca dari baris sumber
// M_LIFE_PREMIUM_DETAIL (66,8 juta baris) lewat kunci PL_NUMBER + CERTIFICATE_NO + ID (= SOURCE_ID) - kunci yang sama
// dengan sqlTertanggungPeserta. NO_UNNEST menahan Oracle agar menelusuri DARI kasus antrian -> peserta klaim
// (CLAIM_ID, IX_PLD_CLAIM_ID) -> satu baris sumber per kunci, bukan memindai tabel sumber dengan LIKE nama. ROWNUM = 1
// membatasi tiap EXISTS (hasilnya sama). Nama hanya dipakai di WHERE, tidak pernah dikeluarkan.
func sqlInboxCari(peserta, sumber string) string {
	return fmt.Sprintf(` AND (UPPER(w.ID) LIKE :cariKasus ESCAPE '\'
	    OR UPPER(h.CLAIM_NO) LIKE :cariKlaim ESCAPE '\'
	    OR EXISTS (SELECT /*+ NO_UNNEST */ 1 FROM %s pc
	                WHERE pc.CLAIM_ID = w.ID AND pc.STS_HAPUS IS NULL
	                  AND UPPER(TRIM(pc.PL_NUMBER)) LIKE :cariPL ESCAPE '\' AND ROWNUM = 1)
	    OR EXISTS (SELECT /*+ NO_UNNEST */ 1 FROM %s pn
	                WHERE pn.CLAIM_ID = w.ID AND pn.STS_HAPUS IS NULL
	                  AND EXISTS (SELECT /*+ NO_UNNEST */ 1 FROM %s m
	                               WHERE m.PL_NUMBER = pn.PL_NUMBER AND m.CERTIFICATE_NO = pn.CERTIFICATE_NO
	                                 AND m.ID = pn.SOURCE_ID
	                                 AND UPPER(m.NAME_OF_INSURED) LIKE :cariNama ESCAPE '\' AND ROWNUM = 1)
	                  AND ROWNUM = 1))`,
		peserta, peserta, sumber)
}

// saringanCariInbox merakit saringan cari dengan nama tabel sumber yang sudah dikualifikasi.
func (r *KlaimLife) saringanCariInbox(peserta string) (string, error) {
	sumber, err := r.db.Qualify(namaTabelPeserta)
	if err != nil {
		return "", err
	}
	return sqlInboxCari(peserta, sumber), nil
}
