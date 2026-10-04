package repository

// Identitas produk baru (paket 3), tanpa Oracle. Penulis JSONDATA dibuang 02-10-2026 (tabel flat; penulis flat:
// `mpnl_flat_test.go`).

import (
	"errors"
	"testing"
)

func TestFormatIdentitasLimaDigitSepertiProsedur(t *testing.T) {
	// `dba-procedures-and-ddl.md` §1: concat('1', lpad(M_PRODUCT_LIFE_SEQ.nextval, 5, '0')).
	for n, mau := range map[int64]string{1: "100001", 421: "100421", 99999: "199999"} {
		if got, err := FormatIdentitas(n); err != nil || got != mau {
			t.Errorf("FormatIdentitas(%d) = %q, %v; mau %q", n, got, err, mau)
		}
	}
	for _, n := range []int64{100000, -1} {
		if _, err := FormatIdentitas(n); !errors.Is(err, ErrIdentitasMelampauiLebar) {
			t.Errorf("FormatIdentitas(%d): nomor tak muat 5 digit harus gagal terang, dapat %v", n, err)
		}
	}
}
