package repository

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
)

// Identitas sementara: lebar tetap, bukan angka, berurutan; di luar mode utuh
// tidak ada identitas sementara.
func TestIdentitasSementaraTCO(t *testing.T) {
	if _, ok := identitasSementaraTCO(context.Background(), SeqKlausulTCO); ok {
		t.Fatal("identitas sementara di luar transaksi utuh")
	}
	if DenganTransaksiUtuhTCO(context.Background(), nil) != context.Background() {
		t.Error("tx nil harus mengembalikan ctx apa adanya")
	}
	ctx := DenganTransaksiUtuhTCO(context.Background(), &Tx{})
	a, _ := identitasSementaraTCO(ctx, SeqKontrakTCO)
	b, _ := identitasSementaraTCO(ctx, SeqJejakTCO)
	if a != "S000000001T" || b != "S000000002T" || !PolaIdentitasSementaraTCO().MatchString(a) {
		t.Errorf("sementara: %q %q", a, b)
	}
	// IdentitasBerikutTCO tidak menyentuh sequence dalam mode utuh.
	var d DB
	id, err := d.IdentitasBerikutTCO(ctx, &Tx{}, SeqReinsurerTCO)
	if err != nil || id != "S000000003T" {
		t.Errorf("IdentitasBerikutTCO dalam mode utuh: %q %v", id, err)
	}
	if _, err := d.IdentitasBerikutTCO(ctx, &Tx{}, "SEQ_LAIN"); !errors.Is(err, ErrSequenceTakDikenal) {
		t.Errorf("sequence asing tetap ditolak: %v", err)
	}
	// Pembaca: tanpa tx nyata tetap pool.
	if d.bacaTCO(ctx) != d.sql {
		t.Error("tx kosong tidak boleh dipakai membaca")
	}
	if _, err := d.TetapkanIdentitasTCO(context.Background(), &Tx{}); !errors.Is(err, ErrIdentitasUtuhTCO) {
		t.Errorf("tetapkan di luar mode utuh: %v", err)
	}
}

// Kolom salinan reinsurer = kolom DDL 302, satu per satu.
func TestKolomReinsurerSalinanSamaDenganDDL(t *testing.T) {
	isi, err := os.ReadFile("migrations/302_t_treatyreinsurer.sql")
	if err != nil {
		t.Fatal(err)
	}
	_, kolom := KolomCreateTable(strings.Split(string(isi), "\n/\n")[0])
	if strings.Join(kolom, ",") != strings.Join(kolomReinsurerTCO, ",") {
		t.Errorf("DDL %v, salinan %v", kolom, kolomReinsurerTCO)
	}
}

func TestSQLTransaksiUtuhTCO(t *testing.T) {
	for _, q := range []string{sqlGantiIdentitasTCO("S.T"), sqlSalinReinsurerTCO("S.R"), sqlAlihSecurityTCO("S.S"),
		sqlBuangReinsurerSementaraTCO("S.R"), sqlGantiBarisJejakTCO("S.J")} {
		if err := PeriksaSQL(q); err != nil {
			t.Errorf("%v: %s", err, q)
		}
	}
	if !strings.Contains(sqlSalinReinsurerTCO("S.R"), "SELECT :1, TREATYYEAR,") {
		t.Errorf("salinan: %s", sqlSalinReinsurerTCO("S.R"))
	}
	// Seluruh sequence tabel yang ditulis transaksi utuh punya tabel tujuan.
	for _, seq := range []string{SeqKontrakTCO, SeqReinsurerTCO, SeqSecurityTCO, SeqBusinessTCO, SeqKlausulTCO, SeqJejakTCO} {
		if tabelPerSequenceTCO[seq] == "" {
			t.Errorf("%s tanpa tabel", seq)
		}
	}
}
