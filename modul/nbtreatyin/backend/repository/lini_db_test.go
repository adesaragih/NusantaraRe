//go:build db

package repository_test

// T_GENERAL_POLIS / T_WORK_POLIS BERSAMA FacIn + Treaty In (keputusan work
// owner 04-10-2026) lawan Oracle SUNGGUHAN. Satu kasus FacIn (`T_WORK_POLIS.LINI
// = 'FAC'` + baris T_GENERAL_POLIS-nya, hanya kolom dasar + PRODKE bawaan)
// disisipkan langsung, berdampingan dengan kasus Treaty dari `SisipKasus`:
// daftar portal, `Keadaan`, dan `KunciKasus` tidak pernah melihatnya, dan
// barisnya tidak berubah. ⛔ Belum pernah dijalankan: skema uji K11 kosong.
// Fixture berawalan UJI-.

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	intidb "nusantarare/inti/backend/db"
	"nusantarare/modul/nbtreatyin/backend/models"
	"nusantarare/modul/nbtreatyin/backend/repository"
)

func TestBarisLiniFACTidakTerbacaGerbangKasus(t *testing.T) {
	sqlDB, skema, ctx, d := pasang(t)
	g := repository.Baru(d)
	const nb, fac = "UJI-NB-LINI-1", "UJI-NB-LINI-FAC"
	if err := dalamTx(t, ctx, d, func(tx *intidb.Tx) error { return g.SisipKasus(ctx, tx, nb, "UJI-A", "UJI A") }); err != nil {
		t.Fatal(err)
	}
	// Kasus FacIn: tahap dan posisi sama dengan kasus Treaty baru, beda LINI saja.
	if _, err := sqlDB.ExecContext(ctx, fmt.Sprintf(`INSERT INTO %s.T_WORK_POLIS (ID, LINI, POSITION, STATUS_WORK,
		FLAG_ONGOING_POLICY, CREATE_OP, TGL_CREATE, TGL_UPDATE) VALUES (:1, 'FAC', :2, :3, '1', 'UJI-FAC', SYSDATE, SYSDATE)`, skema),
		fac, models.PositionAdmin, models.AssignmentAdmin); err != nil {
		t.Fatal(err)
	}
	if _, err := sqlDB.ExecContext(ctx, fmt.Sprintf(`INSERT INTO %s.T_GENERAL_POLIS (ID, IDPEGA, COB_GROUP, POSITION_NOTE)
		VALUES (:1, :2, 'UJI-COB', :3)`, skema), fac, "UJI-FAC "+fac, models.PosisiAdmin); err != nil {
		t.Fatal(err)
	}

	daftar, err := g.DaftarKasus(ctx, models.SaringanKasus{Cari: "UJI-NB-LINI"})
	if err != nil {
		t.Fatal(err)
	}
	if len(daftar) != 1 || daftar[0].ID != nb {
		t.Fatalf("daftar portal = %+v, harap hanya kasus Treaty %s", daftar, nb)
	}
	if _, err := g.Keadaan(ctx, nil, fac); !errors.Is(err, repository.ErrKasusTidakAda) {
		t.Fatalf("Keadaan baris FAC: %v, harap ErrKasusTidakAda", err)
	}
	if k, err := g.Keadaan(ctx, nil, nb); err != nil || k.ID != nb {
		t.Fatalf("Keadaan kasus Treaty: %+v %v", k, err)
	}
	err = dalamTx(t, ctx, d, func(tx *intidb.Tx) error { return g.KunciKasus(ctx, tx, fac, models.AssignmentAdmin) })
	if !errors.Is(err, repository.ErrKasusTidakAda) {
		t.Fatalf("KunciKasus baris FAC: %v, harap ErrKasusTidakAda", err)
	}
	// Baris FAC utuh: kolom dasar tetap, kolom Treaty tidak terisi selain bawaan PRODKE.
	var idpega, cob sql.NullString
	var prodke int
	if err := sqlDB.QueryRowContext(ctx, fmt.Sprintf(`SELECT IDPEGA, COB_GROUP, PRODKE FROM %s.T_GENERAL_POLIS WHERE ID = :1`, skema),
		fac).Scan(&idpega, &cob, &prodke); err != nil {
		t.Fatal(err)
	}
	if idpega.String != "UJI-FAC "+fac || cob.String != "UJI-COB" || prodke != 0 {
		t.Fatalf("baris FAC berubah: IDPEGA %q COB_GROUP %q PRODKE %d", idpega.String, cob.String, prodke)
	}
}
