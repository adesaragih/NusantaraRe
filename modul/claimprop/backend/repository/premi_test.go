package repository

import (
	"errors"
	"testing"
)

// CekLunasPremi_Sql membaca ARASAPAS.INVOICE / DETAIL_INVOICE yang tidak ada di DEV (ORA-00942). Keputusan work owner
// 09-10-2026: di luar produksi tabel yang tidak terbaca = premi dianggap lunas; di produksi tetap galat; galat lain
// tidak pernah ditelan.
func TestSaldoPremiTanpaARASAPASDiLuarProduksi(t *testing.T) {
	tiada := errors.New("ORA-00942: table or view does not exist")
	if v, err := saldoPremiGagal(false, tiada); err != nil || v != "" {
		t.Fatalf("non-produksi: dapat (%q, %v), mau (\"\", nil)", v, err)
	}
	if _, err := saldoPremiGagal(true, tiada); !errors.Is(err, tiada) {
		t.Fatalf("produksi wajib meneruskan galat, dapat %v", err)
	}
	lain := errors.New("ORA-01013: user requested cancel of current operation")
	if _, err := saldoPremiGagal(false, lain); !errors.Is(err, lain) {
		t.Fatalf("galat selain tabel tidak ada tidak boleh ditelan, dapat %v", err)
	}
}
