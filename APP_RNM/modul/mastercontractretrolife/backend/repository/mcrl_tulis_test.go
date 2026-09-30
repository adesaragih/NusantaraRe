package repository

// Penulis (paket 2+) - identitas dan teks SQL, tanpa Oracle.

import (
	"context"
	"errors"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
)

func TestFormatIdentitasSepertiProcedure(t *testing.T) {
	for n, mau := range map[int64]string{44: "1000044", 1: "1000001", 999999: "1999999"} {
		if got, err := FormatIdentitas(n); err != nil || got != mau {
			t.Errorf("FormatIdentitas(%d) = %q, %v; mau %q", n, got, err, mau)
		}
	}
	// ⛔ LPAD Oracle memotong diam-diam; di sini gagal terang.
	for _, n := range []int64{1000000, -1} {
		if _, err := FormatIdentitas(n); !errors.Is(err, ErrIdentitasMelampauiLebar) {
			t.Errorf("FormatIdentitas(%d): %v, mau ErrIdentitasMelampauiLebar", n, err)
		}
	}
}

func TestSQLPenulisTahun(t *testing.T) {
	kasus := map[string]struct {
		sql   string
		wajib []string
	}{
		"sisip": {sqlSisipTahun("S.TREATYYEAR_LIFE"), []string{"INSERT INTO S.TREATYYEAR_LIFE",
			"VALUES (:1, :2, :3, :4, SYSDATE, TO_DATE(:5, 'YYYY-MM-DD'), TO_DATE(:6, 'YYYY-MM-DD'))"}},
		"perbarui": {sqlPerbaruiTahun("S.TREATYYEAR_LIFE"), []string{"TGLUPDATE = SYSDATE", "WHERE ID = :6"}},
		"salin business": {sqlSalinTahunKeBusiness("S.TREATYBUSINESS_LIFE"), []string{"SET TREATYYEAR = :1",
			"WHERE TREATYYEARID = :3 AND DECODE(TREATYYEAR, :4, 0, 1) = 1"}},
		"salin kontrak": {sqlSalinTahunKeKontrak("S.TREATYCONTRACT_LIFE"), []string{"WHERE IDTREATYYEAR = :4",
			"DECODE(TREATYSTARTDATE, TO_DATE(:5, 'YYYY-MM-DD'), 0, 1) = 1"}},
		"sequence": {sqlNomorBerikut("S.TREATYYEAR_LIFE_SEQ"), []string{"SELECT S.TREATYYEAR_LIFE_SEQ.NEXTVAL FROM DUAL"}},
	}
	for nama, k := range kasus {
		rata := strings.Join(strings.Fields(k.sql), " ")
		for _, w := range k.wajib {
			if !strings.Contains(rata, w) {
				t.Errorf("%s: tanpa %q:\n%s", nama, w, rata)
			}
		}
		if err := db.PeriksaSQL(k.sql); err != nil {
			t.Errorf("%s: %v", nama, err)
		}
	}
}

func TestSequenceHanyaLimaMilikModul(t *testing.T) {
	g := &Gudang{}
	if _, err := g.identitasBaru(context.Background(), nil, "TREATYSECURITY_LIFE_SEQ"); !errors.Is(err, ErrSequenceTakDikenal) {
		t.Errorf("sequence DEV yang tidak dirujuk XML harus ditolak (K3): %v", err)
	}
}
