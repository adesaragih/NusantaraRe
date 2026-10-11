# PR — modul baru R/I Risk (`ririsklife`, MASTER TREATY) + satu tabel per jenis data (migrasi inti 935-941)

## Summary

```diff
 POOLDATA
-  VIEW  RIRISK_LIFE_SUMMARY  = SELECT a.ID, a.JSONDATA.USEDBY, … FROM M_RIRISK_LIFE_SUMMARY a
-  TABLE M_RIRISK_LIFE_SUMMARY (ID PK, JSONDATA IS JSON)
+  TABLE RIRISK_LIFE_SUMMARY  (ID PK, USEDBY VARCHAR2(200), MODIFIEDDATE VARCHAR2(50), OPERATORID VARCHAR2(200))   935-937
+  INDEX IX_RIRISK_LIFE_SUMMARY_NAMA (UPPER(TRIM(USEDBY)))
-  VIEW  RIRISK_LIFE          = SELECT a.ID, a.JSONDATA.IDUSEDBY, … FROM M_RIRISK_LIFE a
-  TABLE M_RIRISK_LIFE        (ID VARCHAR2(6) tanpa PK, JSONDATA, kolom datar BASI, INDEX4)
+  TABLE RIRISK_LIFE          (ID VARCHAR2(6) PK_RIRISK_LIFE, IDUSEDBY, USEDBY, AGE, YEAR, MONTH, RISK NUMBER, CONTRACT)  938-940
+  INDEX IX_RIRISK_LIFE_IDUSEDBY - HANYA bila belum ada indeks berkolom pertama IDUSEDBY (INDEX4)
+  M_NAV_MENU 'ririsklife' "R/I Risk" MASTER TREATY URUTAN 11                                                     941
 M_RIRISK_LIFE_TEMP  tidak disentuh; sequence M_RIRISK_LIFE*_SEQ tetap; PEGA_M_RIRISK_LIFE* INVALID (K3)
```

```text
935/938  blok ALL_VIEWS: DROP VIEW (hanya bila VIEW) ; blok ALL_TAB_COLUMNS sumber: ALTER TABLE M_… RENAME TO …
936/939  ALTER TABLE … ADD (kolom yang belum ada)                                   -- berdiri sendiri
937/940  blok: isi ULANG semua kolom view dari JSONDATA (940: RISK tanpa NLS) ; 940: PK ; blok: DROP COLUMN JSONDATA
         CASCADE CONSTRAINTS ; indeks (937 nama; 940 blok ALL_IND_COLUMNS)          -- semua aman diulang
_down    936/939 + 935/938: pelindung UPDATE … SET JSONDATA = JSONDATA WHERE 1 = 0 (ORA-00904) ; 937/940: blok aman diulang
```

Modul: ringkasan (grid 50, filter, urut, Save Add / Edit dengan nama ikut ke rincian, Delete berantai - Save / Cancel
DITAMPILKAN atas keputusan work owner 08-10-2026 walau wadahnya `1=2` di XML, RALAT R1); R/I RISK DETAIL (tambah / EDIT, CONTRACT, YEAR, MONTH, RISK (PERMIL), grid 200); Upload CSV / View Upload /
Simpan Upload (USEDBY, CONTRACT, YEAR, MONTH, RISK). Paritas: `docs/PARITAS-LAYAR-DAN-AKSI.md`. Urutan WO:
`docs/LANGKAH-WO-RIRISKLIFE.md` + `docs/sql/ririsk_*.sql`.

## Perlu tinjauan tim inti

**Perluasan `inti/backend/migrasi/perintah_katalog.go` + penjaga DITERIMA work owner 08-10-2026**, tetap ditinjau tim
inti. Alasan: K1 meminta DROP VIEW "hanya bila objeknya memang VIEW", RENAME "hanya bila sumber ada dan target belum
TABLE", dan indeks IDUSEDBY "tanpa indeks kembar" - dengan tiga katalog lama (ALL_TAB_COLUMNS / ALL_CONSTRAINTS /
ALL_INDEXES) dan aturan "objek yang ditanyakan disebut perintah" ketiganya tidak dapat ditulis sebagai blok yang aman
diulang: ALL_TAB_COLUMNS memuat kolom VIEW dan TABLE bernama sama, perintah RENAME tidak menyebut kolom mana pun, dan
nama indeks warisan (`INDEX4`) beserta kolomnya tidak diketahui (CREATE INDEX berkolom sama mati ORA-01408). Tanpa
bentuk baru, pilihannya pernyataan biasa yang gagal keras saat pelari diulang sesudah mati di tengah (pemulihan manual)
- lebih rapuh daripada blok yang melewati diri. Bentuknya sempit (pola regex tunggal, OWNER / TABLE_OWNER skema,
`COLUMN_POSITION = 1`, RENAME hanya atas tabel yang ditanyakan dengan n > 0) dan diuji menggigit.

- **Dua bentuk blok berpelindung katalog BARU** di `inti/backend/migrasi/perintah_katalog.go` (`BacaPerintahKatalog`):
  `SYS.ALL_VIEWS` (DROP VIEW hanya bila objeknya VIEW - ALL_TAB_COLUMNS tidak membedakan VIEW dan TABLE) dan
  `SYS.ALL_IND_COLUMNS ... COLUMN_POSITION = 1` (CREATE INDEX hanya bila belum ada indeks berkolom pertama itu - INDEX4
  warisan bernama lain; tanpa ini CREATE INDEX berkolom sama mati ORA-01408). Uji `TestBacaPerintahKatalogViewDanKolomIndeks`.
- **Penjaga** `inti/backend/penjaga/migrasi_test.go` `pelanggaranBlokPLSQL`: bentuk `ALTER TABLE {skema}.<sumber> RENAME
  TO <baru>` sah bila blok menanyakan kolom tabel SUMBER (ALL_TAB_COLUMNS, n > 0) - kolom itu tidak disebut perintah.
  Uji `TestAturanBlokPLSQLMenggigit` (+ 8 kasus). Pengecualian `DROP COLUMN JSONDATA` diperluas ke 937 dan 940 (pola
  928/930/932/934).
- `rentang_test.go` (urutan pelari 935-941), `menu_test.go` (`modulLuarKorpus` + `labelTampilDisetujui`
  `ririsklife` → folder "R/I Risk Life", tampil "R/I Risk" - K4), `cmd/api/gerbang_tulis_test.go`,
  `frontend/katalogKorpus.ts`, `frontend/daftar.menuTabel.test.ts`, `frontend/Shell.test.ts` (43 kelompok),
  `inti/backend/daftar/modul_ririsklife_gen.go`.
- 937_down / 940_down memulihkan constraint bernama asli `ENSURE_M_RIRISK_LIFE_SUMMARY_JSON` (33 byte) - penjaga 30
  byte hanya membaca jalur maju; DEV sudah memuat nama itu.
- `inti/backend/db/koneksi.go` (`db.Koneksi`) kini dipakai uji `-tags=db` `TestDBKonversiRiskTanpaNLS` (ALTER SESSION di
  satu koneksi) - produksi tetap tidak memakainya.

## Pembaca lain (perubahan minimal)

- `masterproductnamelife`: `testdata/katalog-dev.json` (`RIRISK_LIFE_SUMMARY` = TABLE), komentar `mpnl_master.go`,
  `mpnl_katalog_test.go`, tiruan `mpnl_tiruan_db_test.go` (lebar 10 / 200), `MODUL.md`. Kueri tidak berubah.
- `premiumlistlife`: `sqlRiskHitungQR` membaca `RISK` lewat `db.FmtDesimal` (TM9 ber-NLS eksplisit - tidak bergantung
  format driver untuk NUMBER); `angkaMaster` sudah menerima koma dan titik; uji `TestHitungQRRiskKomaDanTitikSama`
  ("921,9" vs "921.9" / ".9" → hasil QR SAMA). View R/I Risk (`RiskProduk`) tidak berubah (tampilan). Komentar
  "view `RIRISK_LIFE`" → tabel.

## Evidence

- `go test` modul: models (dua rumus ID, validasi, CSV), services (tiruan; Add, nama kembar, Edit nama ikut ke rincian),
  handlers (Save 201 / 200, kembar 422, View only 403, tak ada 404),
  repository (SQL, migrasi 935-941 lewat pengurai produksi pelari). `-tags=db` (terkompilasi, BELUM dijalankan):
  `TestDBMigrasiSatuTabelDanMundur`, `TestDBKonversiRiskTanpaNLS` (NLS koma), `TestDBLayananRincian`,
  `TestDBUnggahDanRollback`, `TestDBSimpanRingkasan`. vitest modul 16 uji (form Save / Cancel, aturan, klien).
- Bukti konversi RISK di DEV: `docs/sql/ririsk_bukti_konversi.sql` - **menunggu hasil WO** (executor tanpa koneksi Oracle).

## Merge Danger

**Door:** one-way sesudah `-migrate` 937/940 (JSONDATA dibuang). Jalur mundur membangun ulang JSONDATA, nama lama, dan
view dari kolom; `pxObjClass`, bentuk RISK asli ("921,9"), dan nilai kolom datar basi hanya dari cadangan CSV P3.

**Blast Radius:** lintas-modul - `masterproductnamelife` (`Choose R/I Risk`) dan `premiumlistlife` (View R/I Risk,
Hitung QR) memakai nama yang sama; skema `GL` (SELECT atas `M_RIRISK_LIFE`) WAJIB diperiksa P2 sebelum `-migrate`.
Backend lama gagal di Hitung QR sesudah 940 sampai restart (langkah (e)).
