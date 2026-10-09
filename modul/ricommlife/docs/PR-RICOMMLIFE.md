# PR — R/I Comm Life (MASTER TREATY): modul baru + SATU tabel per jenis data (RALAT R1, migrasi inti 931-934)

> **Urutan merge:** folder kerja `dev` sudah memuat riratelife (922/923, 926-930) dan ricommlife (924/925); perubahan
> RALAT R1 (931-934) berada di atas keduanya. Bila cabang riratelife belum masuk ke target merge, ia LEBIH DULU
> (`modul/riratelife/docs/PR-RIRATELIFE-FLAT.md`); konflik yang diketahui hanya daftar urutan pelari di
> `inti/backend/penjaga/rentang_test.go` (`… 930, 931, 932, 933, 934, 952_menu_tiruan`).

## RALAT R1 (keputusan work owner 08-10-2026)

```diff
 POOLDATA
 M_RICOMM_LIFE_SUMMARY
-  ID VARCHAR2(10) PK, JSONDATA CLOB (ENSURE_M_RICOMM_LIFE_SUMMARY_JSON)  + VIEW RICOMM_LIFE_SUMMARY
+  ID VARCHAR2(10) PK, USEDBY VARCHAR2(200), MODIFIEDDATE VARCHAR2(50), OPERATORID VARCHAR2(200)
+  INDEX IX_M_RICOMM_LIFE_SUMMARY_NAMA (UPPER(TRIM(USEDBY)))
 M_RICOMM_LIFE
-  ID VARCHAR2(10) PK, JSONDATA CLOB (ENSURE_M_RICOMM_LIFE_JSON)          + TABLE RICOMM_LIFE (924, rincian aplikasi)
+  ID VARCHAR2(10) PK, IDUSEDBY VARCHAR2(10), USEDBY VARCHAR2(200), CONTRACT NUMBER(5), YEAR NUMBER(5), COMM NUMBER(38,8)
+  INDEX IX_M_RICOMM_LIFE_IDUSEDBY (IDUSEDBY)
```

```text
931  ALTER TABLE M_RICOMM_LIFE_SUMMARY ADD (USEDBY, MODIFIEDDATE, OPERATORID)            -- berdiri sendiri
932  blok: UPDATE kolom dari JSONDATA; blok: DROP COLUMN JSONDATA CASCADE CONSTRAINTS     -- semua aman diulang
     CREATE INDEX IX_M_RICOMM_LIFE_SUMMARY_NAMA; DROP VIEW RICOMM_LIFE_SUMMARY            -- terakhir
933  ALTER TABLE M_RICOMM_LIFE ADD (IDUSEDBY, USEDBY, CONTRACT, YEAR, COMM)              -- berdiri sendiri
934  blok: UPDATE kolom dari JSONDATA (DEV 0 baris); blok: DROP COLUMN JSONDATA …
     UPDATE dari RICOMM_LIFE (ID sama); INSERT … NOT EXISTS (SESUDAH JSONDATA dibuang - NULLABLE asli tidak diketahui)
     CREATE INDEX IX_M_RICOMM_LIFE_IDUSEDBY; DROP TABLE RICOMM_LIFE CASCADE CONSTRAINTS   -- terakhir
_down  931/933: pelindung gagal-keras UPDATE … SET JSONDATA = JSONDATA WHERE 1 = 0 (ORA-00904), lalu DROP kolom
       932/934: aman diulang (blok berpelindung); 934_down membangun ulang RICOMM_LIFE persis 924 berisi data
```

Kode: repository membaca/menulis kolom kedua tabel (nol JSON, nol view); `ricl_json.go`, `pindah.go`, alat
`pindahflat`, `docs/DBA-LEPAS-VIEW-RICOMM_LIFE.sql` dihapus; `docs/LANGKAH-WO-RICOMMLIFE.md` diarsipkan. Urutan WO:
`docs/LANGKAH-WO-RICOMMLIFE-SATU-TABEL.md` + `docs/sql/satu_tabel_*.sql`.

## Perlu tinjauan tim inti

- Pengecualian penjaga `TestKolomUangDesimalDanNolJSON` (`inti/backend/penjaga/migrasi_test.go`, peta `buangJSON`)
  diperluas ke tepat satu `DROP COLUMN JSONDATA CASCADE CONSTRAINTS` di `932_m_ricomm_life_summary_satu_tabel.sql` dan
  `934_m_ricomm_life_satu_tabel.sql` (pola 928/930, keputusan WO 08-10-2026).
- `inti/backend/penjaga/rentang_test.go` (urutan pelari + 931-934), `menu_test.go` (komentar), `frontend/katalogKorpus.ts`
  (komentar).
- Migrasi inti baru `931`-`934` (+ `_down`). `932_down` memulihkan constraint bernama asli 33 byte
  (`ENSURE_M_RICOMM_LIFE_SUMMARY_JSON`) - penjaga 30 byte hanya membaca jalur maju; DEV sudah memuat nama itu.
- `inti/backend/db/koneksi.go` (`db.Koneksi`) TIDAK lagi dipakai kode mana pun (satu-satunya pemakai = alat pindahflat
  yang dihapus) - **kandidat bersih-bersih tim inti**, sengaja tidak dihapus di sini.

## Modul baru (06-10-2026, riwayat)

```diff
 POOLDATA
-  VIEW  RICOMM_LIFE = SELECT a.ID, a.JSONDATA.IDUSEDBY, … FROM M_RICOMM_LIFE a
+  TABLE RICOMM_LIFE (ID VARCHAR2(10) PK, IDUSEDBY, USEDBY, CONTRACT NUMBER(5), YEAR NUMBER(5), COMM NUMBER(38,8))  (924)
+  M_NAV_MENU 'ricommlife' "R/I Comm Life" MASTER TREATY URUTAN 10                   (925)
```

ID baru = site `M_SITE_DATABASE` || LPAD(sequence warisan, 6) - TETAP sesudah RALAT R1.

## Evidence

- **Before (RALAT R1):** ringkasan JSON + view, rincian di tabel flat terpisah dari tabel Pega.
  **After:** `TestMigrasi931KolomRingkasan`, `TestMigrasi932SatuTabel`, `TestMigrasi933KolomKomisi`,
  `TestMigrasi934SatuTabel` (pengurai produksi pelari: `KolomAlterTambah`, `BacaPerintahKatalog`, `KolomAlterBuang`;
  pelindung `_down`; bentuk 924 di 934_down), `TestSqlRingkasan`, `TestSqlKomisi`, `TestNolJSONDanObjekLamaDiRepository`;
  `-tags=db` `TestDBMigrasiSatuTabelDanMundur` (maju diulang sebagian, ORA-00904 pelindung, mundur diulang sebagian),
  `TestDBKembarKepemilikanDanDeleteBerantai` - terkompilasi (`go vet -tags=db`), belum dijalankan terhadap Oracle;
  `go test ./...` dan `npm test` = baseline.

## Merge Danger

**Door:** one-way sesudah `-migrate` 932/934.

Jalur mundur membangun ulang JSONDATA, view, dan tabel `RICOMM_LIFE` dari kolom, tetapi `pxObjClass` / bentuk JSON
asli hanya kembali dari cadangan CSV (b). Prosedur `PEGA_M_RICOMM_LIFE` / `PEGA_M_RICOMM_LIFE_SUMMARY` INVALID (diterima WO).

**Blast Radius:** satu-modul.

Hanya ricommlife membaca/menulis objek ini (grep repo). Backend lama (membaca view / `RICOMM_LIFE`) gagal sesudah 932/934
sampai biner baru dijalankan - restart wajib (langkah (e)).
