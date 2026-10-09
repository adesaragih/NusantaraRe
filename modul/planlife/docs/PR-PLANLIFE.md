# PR — modul baru Plan (`planlife`, MASTER TREATY) + satu tabel PRODUCT_TYPE_LIFE (migrasi inti 946-949)

## Summary

```diff
 POOLDATA
-  VIEW  PRODUCT_TYPE_LIFE   = SELECT a.JSONDATA.ID, a.JSONDATA.CoverName, … FROM M_PRODUCT_TYPE_LIFE a
-  TABLE M_PRODUCT_TYPE_LIFE (ID VARCHAR2(6) NULLABLE TANPA PK, JSONDATA IS JSON)
+  TABLE PRODUCT_TYPE_LIFE   (ID VARCHAR2(6) PK_PRODUCT_TYPE_LIFE, COVERNAME 200, BUSINESS 200, BUSINESSID 10,
+                             BENEFIT 200, BENEFITID 10)                                                 946-948
+  M_NAV_MENU 'planlife' "Plan" MASTER TREATY URUTAN 13                                                 949
 M_PRODUCT_TYPE_LIFE_SEQ tetap; PEGA_M_PRODUCT_TYPE_LIFE INVALID (K3); BUSINESS / M_BUSINESS / BENEFIT_LIFE /
 PLAN_LIFE_SUMMARY tidak disentuh; tanpa FK
```

```text
946  blok ALL_VIEWS: DROP VIEW ; blok ALL_TAB_COLUMNS sumber (JSONDATA): RENAME (target TABLE yang ada = ORA-00955 keras)
947  ALTER TABLE PRODUCT_TYPE_LIFE ADD (5 kolom)                                   -- berdiri sendiri
948  blok isi (notasi titik = teks view) ; blok periksa ID (NULL / ganda / <> JSONDATA.ID) ; blok periksa isi (K2) ;
     blok ALL_CONSTRAINTS: ADD PK (sekaligus NOT NULL) ; blok DROP COLUMN JSONDATA CASCADE CONSTRAINTS
     pemeriksaan = UPDATE … SET m.COVERNAME = RPAD(CHR(88), 201, CHR(88)) WHERE <pelanggaran>  -> ORA-12899, tanpa kutip
_down 947 + 946: pelindung ORA-00904 ; 948: blok aman diulang (buang PK -> ID NULLABLE lagi)
```

Modul (pola benefitlife + dua autocomplete): form Plan Name / Business / Benefit, Save, New (saat Edit), grid 10 baris
Plan Name / Business / Benefit / Edit. API `/api/plan-life`: GET, GET pilihan business / benefit, POST, PUT.
Paritas: `docs/PARITAS-LAYAR-DAN-AKSI.md`. Urutan WO (P3 di paling atas): `docs/LANGKAH-WO-PLANLIFE.md`.

## Perlu tinjauan tim inti

- **`inti/backend/migrasi/perintah_katalog.go` TIDAK diubah.** Semua blok 946-948 memakai bentuk yang ada. Pemaksa
  berhenti 948 tidak bisa memakai `SET m.ID = NULL` (cara 944) karena ID masih NULLABLE sebelum PK; dipakai nilai 201
  huruf ke COVERNAME VARCHAR2(200) → ORA-12899 (tidak ditoleransi pelari). Bentuk bertanda kutip (JSON_EXISTS) ditolak
  pengurai (`TestBentukPemaksaBerhenti`).
- **PK = NOT NULL**: ADD CONSTRAINT PRIMARY KEY pada kolom NULLABLE membuat Oracle menambah NOT NULL; DROP CONSTRAINT di
  948_down mengembalikannya NULLABLE (diuji `-tags=db` `TestDBMigrasiSatuTabelDanMundur`, belum dijalankan) - satu blok
  ALL_CONSTRAINTS, tanpa MODIFY yang mati ORA-01442 saat diulang.
- **RENAME "target belum TABLE"**: bentuk penjaga RENAME menanyakan satu katalog (sumber); target TABLE yang sudah ada
  sudah gagal keras ORA-00955 di dalam blok, jadi tidak perlu bentuk baru.
- Penjaga: `migrasi_test.go` (pengecualian `DROP COLUMN JSONDATA` + 948), `rentang_test.go` (946-949), `menu_test.go`
  (`modulLuarKorpus` + `labelTampilDisetujui` `planlife` → folder "Plan Life", tampil "Plan"); `cmd/api/gerbang_tulis_test.go`;
  `frontend/katalogKorpus.ts`, `daftar.menuTabel.test.ts`, `Shell.test.ts` (45); `inti/backend/daftar/modul_planlife_gen.go`.

## Pembaca lain (minimal)

`masterproductnamelife`: kueri TIDAK berubah (nama dan kolom sama). `testdata/katalog-dev.json` (`PRODUCT_TYPE_LIFE` =
TABLE), komentar `mpnl_master.go` dan `mpnl_katalog_test.go`, DDL tiruan `mpnl_tiruan_db_test.go` (lebar 6 / 200),
`MODUL.md` (tabel dibaca saja). Semua uji MPNL dijalankan ulang: hijau.

## Evidence

- `go test` modul: models (rumus ID, K5), services (K4 cocok / tolak / ganda, K5 wajib / unik, K3 409, balapan PK,
  View only), handlers (201 / 200 / 422 / 403 / 404 / 400 / 409 / 503, nol DELETE), repository (SQL, 946-949 lewat
  pengurai produksi, pemaksa berhenti). `-tags=db` terkompilasi, BELUM dijalankan. vitest modul + gaya.
- Bukti K2 di DEV: `docs/sql/plan_bukti_k1.sql` - **menunggu hasil WO**.

## Merge Danger

**Door:** one-way sesudah 948 (JSONDATA dibuang). Mundur membangun ulang JSONDATA enam kunci huruf Pega; `px*` hanya dari
cadangan P3. **Blast Radius:** `masterproductnamelife` (autocomplete Plan) membaca nama / kolom yang sama - backend lama
tetap jalan; Pega Plan berhenti dapat menyimpan (K3 diterima).
