# PR — modul baru Cause Of Loss Life (`causeoflosslife`, MASTER TREATY) + satu tabel CAUSEOFLOSS_LIFE (migrasi modul 090-092, slot 955)

## Summary

```diff
 POOLDATA
-  VIEW  CAUSEOFLOSS_LIFE   = SELECT a.ID, a.JSONDATA.CauseofLoss FROM M_CAUSEOFLOSS_LIFE a
-  TABLE M_CAUSEOFLOSS_LIFE (ID VARCHAR2(10) PK SYS_C008825, JSONDATA IS JSON)
+  TABLE CAUSEOFLOSS_LIFE   (ID VARCHAR2(10) PK SYS_C008825, CAUSEOFLOSS VARCHAR2(200) NULLABLE)      090-092
+  M_NAV_MENU 'causeoflosslife' "Cause Of Loss Life" MASTER TREATY URUTAN 14                           955
 M_CAUSEOFLOSS_LIFE_SEQ tetap; PEGA_M_CAUSEOFLOSS_LIFE INVALID (K3); M_CAUSE_OF_LOSS / D_CAUSE_OF_LOSS / V_* /
 T_LISTCAUSEOFLOSS / PEGA_M|D_CAUSE_OF_LOSS tidak disentuh; baris 100001 (kosong) tidak diisi / dihapus
```

```text
090  blok ALL_VIEWS: DROP VIEW ; blok ALL_TAB_COLUMNS sumber (ID): RENAME (target TABLE yang ada = ORA-00955 keras)
091  ALTER TABLE CAUSEOFLOSS_LIFE ADD (CAUSEOFLOSS VARCHAR2(200))                         -- berdiri sendiri
092  blok isi (notasi titik = teks view) ; blok periksa K1.4 ; blok DROP COLUMN JSONDATA CASCADE CONSTRAINTS
     periksa = UPDATE … SET m.ID = NULL WHERE (kolom NULL, view tidak) OR (kolom terisi, view NULL) OR kolom <> view
               -> ORA-01407 (ID NOT NULL), tanpa kutip; NULL = NULL TIDAK menghentikan (100001)
955  INSERT datar M_NAV_MENU … WHERE NOT EXISTS
_down 091 + 090: pelindung ORA-00904 ; 092: blok aman diulang ; 955: buang hak lalu baris
```

Modul (pola benefitlife): form ID (disabled) / Cause of Loss, Save, Cancel (saat Edit), grid 10 baris ID / Cause of
Loss / Edit, ID menaik tetap, tanpa saring. API `/api/cause-of-loss-life`: GET, POST, PUT. Paritas:
`docs/PARITAS-LAYAR-DAN-AKSI.md`. Urutan WO (P3 di paling atas): `docs/LANGKAH-WO-CAUSEOFLOSSLIFE.md`.

## Perlu persetujuan tim inti

- **Jatah nomor dipinjam dari premiumlistlife (K0, keputusan work owner 08-10-2026)**: `modul/premiumlistlife/MODUL.md`
  `Rentang migrasi` `050-099` → `050-089`, `Slot menu` `954-955` → `954-954`; modul ini `090-099` / `955-955`. Nomor
  itu terbukti kosong sebelum ditulis (nol berkas `06[5-9]_*`, `0[7-9]?_*`, `954_*`, `955_*`; nol MODUL.md lain yang
  menyatakannya; premiumlistlife terpakai `050-064`). Itu SATU-SATUNYA perubahan di folder premiumlistlife.
- **Baris menu luar korpus PERTAMA yang lahir di SLOT MENU modul** (rentang inti 900-949 penuh, 949 = planlife):
  penjaga diperluas, bukan dilonggarkan - `barisLahirDiSlot` (menu_test.go) menyatakan SATU berkas
  (`955_menu_causeoflosslife.sql`); `pelanggaranSlotMenu` (rentang_test.go) menerima INSERT datar HANYA atas baris modul
  itu sendiri berDIMIGRASI '1' dan mundur "buang hak lalu baris" HANYA untuk modul itu; `TestBarisLahirDiSlotTerdaftar`
  mengunci kesepadanan dengan `modulLuarKorpus` dan slot pemiliknya; `TestAturanSlotMenuMenggigit` +6 kasus dua arah.
- **`inti/backend/migrasi/perintah_katalog.go` TIDAK diubah.** Blok 090-092 memakai bentuk yang ada. Pemaksa berhenti
  092 = bentuk 944 (`SET m.ID = NULL`, ID NOT NULL sejak Pega), bukan RPAD 948 (yang dipakai karena ID planlife masih
  NULLABLE); pemaksa menulis ID sehingga nilai yang diperiksa tidak pernah diubah pemeriksaannya.
- Penjaga lain: `migrasi_test.go` (pengecualian `DROP COLUMN JSONDATA` + 092); `cmd/api/gerbang_tulis_test.go`;
  `frontend/katalogKorpus.ts`, `daftar.menuTabel.test.ts`, `Shell.test.ts` (46); `inti/backend/daftar/modul_causeoflosslife_gen.go`.
- **Skema uji**: `uji/skemauji/coll_tiruan.go` + `skemauji.go` - tiruan `M_CAUSEOFLOSS_LIFE` + view + sequence dibuat
  SEBELUM pelari (090-092 berjalan sebelum 900 di skema baru), dibongkar SESUDAH migrasi mundur.

## Pembaca lain (minimal)

`masterproductnamelife`: kueri TIDAK berubah (nama dan kolom sama). `testdata/katalog-dev.json` (`CAUSEOFLOSS_LIFE` =
TABLE), komentar `mpnl_master.go` dan `mpnl_katalog_test.go`, DDL tiruan `mpnl_tiruan_db_test.go` (lebar 10 / 200).
Semua uji MPNL dijalankan ulang (`go test`, `go vet`, `go vet -tags=db`): hijau.

## Evidence

- `go test` modul: models (rumus ID, wajib / batas, huruf tetap), services (Add 100005, kembar, Edit, 404, K3 409,
  balapan PK, View only), handlers (201 / 200 / 422 / 403 / 404 / 400 / 409 / 503, nol DELETE), repository (SQL,
  090-092 + 955 lewat pengurai produksi, pemaksa berhenti tanpa kutip, logika NULL = NULL, urutan pelari 090 → 900 →
  955 atas migrasi inti SUNGGUHAN). `-tags=db` (termasuk `TestDBUrutanPelari090Lalu900Lalu955` di skema uji)
  terkompilasi, BELUM dijalankan. vitest modul + gaya.
- Bukti K2 di DEV: `docs/sql/col_bukti_k1.sql` - **menunggu hasil WO**.

## Merge Danger

**Door:** one-way sesudah 092 (JSONDATA dibuang). Mundur membangun ulang JSONDATA kunci `CauseofLoss` huruf Pega;
`pxObjClass` / `pyRuleHarness` / `"CauseofLoss":""` asli hanya dari cadangan P3. **Blast Radius:**
`masterproductnamelife` (pemilih Cause Of Loss) membaca nama / kolom yang sama - backend lama tetap jalan; Pega Cause
Of Loss berhenti dapat menyimpan (K3 diterima); premiumlistlife kehilangan nomor 090-099 / 955 yang tidak pernah
dipakainya.
