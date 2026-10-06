# PR — `modul/ricommlife/implementasi`: modul baru R/I Comm Life (MASTER TREATY)

> ⚠️ **Urutan merge: `modul/riratelife/implementasi` LEBIH DULU (`modul/riratelife/docs/PR-RIRATELIFE-FLAT.md`), lalu
> PR ini.** `inti/backend/db/koneksi.go` identik di kedua cabang; uji coba merge PR ini di atas ujung riratelife bersih.

## Summary

```diff
 POOLDATA
-  VIEW  RICOMM_LIFE = SELECT a.ID, a.JSONDATA.IDUSEDBY, … FROM M_RICOMM_LIFE a
+  TABLE RICOMM_LIFE (ID VARCHAR2(10) PK, IDUSEDBY, USEDBY, CONTRACT NUMBER(5), YEAR NUMBER(5), COMM NUMBER(38,8))  (924)
+  INDEX IX_RICOMM_LIFE_IDUSEDBY                                                     (grid/Delete + pengaman view)
+  M_NAV_MENU 'ricommlife' "R/I Comm Life" MASTER TREATY URUTAN 10                   (925)
 M_RICOMM_LIFE_SUMMARY (JSON)  ringkasan: sisip JSON_OBJECT, ubah baca-ubah-tulis (tidak diubah strukturnya)
 M_RICOMM_LIFE (JSON)          tidak disentuh - sumber alat pindahflat
```

```text
modul/ricommlife/
├── backend/   ringkasan + R/I COMM DETAIL (tambah/edit), Upload CSV, alat pindahflat (satu db.Koneksi)
├── frontend/  R/I COMM SUMMARY + popup R/I COMM DETAIL
└── docs/      STRUKTUR, DBA-LEPAS-VIEW-RICOMM_LIFE.sql, LANGKAH-WO-RICOMMLIFE.md
inti/backend/db/koneksi.go   *sql.Conn: ALTER SESSION + transaksi di koneksi yang sama
```

ID baru = site `M_SITE_DATABASE` || LPAD(sequence warisan, 6). Urutan WO: `docs/LANGKAH-WO-RICOMMLIFE.md`.

## Evidence

- **Before:** menu R/I Comm Life tidak ada; `go test ./...` baseline 3 paket claimlife merah, vitest 14 gagal.
  **After:** `go test ./...` = baseline + 5 paket ricommlife hijau (`TestBentukIDRumusProsedur` site 1 + seq 44 →
  1000044, `TestMigrasi924TeruraiDanBerpasangan`, `TestSimpanKomisi`, `TestAngkaOracleTanpaNLS`); `go vet -tags=db`
  bersih (uji Oracle belum dijalankan); vitest = baseline + 16 uji ricommlife hijau.

## Merge Danger

**Door:** one-way sesudah `-migrate` + aplikasi menulis rincian.

Sebelum itu `924_…_down.sql` memulihkan view persis; sesudahnya rincian baru hanya ada di tabel flat.

**Blast Radius:** satu-modul.

Modul baru; nol kode lain membaca RICOMM_LIFE. Prasyarat DBA (lepas view) wajib sebelum `-migrate`; tanpa itu
`-migrate` berhenti di ORA-01702 (tidak merusak). Asumsi A1-A9 bertanda menunggu WO (MODUL.md).
