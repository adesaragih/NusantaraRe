# PR — modul baru Benefit (`benefitlife`, MASTER TREATY) + satu tabel BENEFIT_LIFE (migrasi inti 942-945)

## Summary

```diff
 POOLDATA
-  VIEW  BENEFIT_LIFE   = SELECT a.ID, a.JSONDATA.Benefit FROM M_BENEFIT_LIFE a
-  TABLE M_BENEFIT_LIFE (ID VARCHAR2(10) PK SYS_C009031, JSONDATA IS JSON)
+  TABLE BENEFIT_LIFE   (ID VARCHAR2(10) PK SYS_C009031, BENEFIT VARCHAR2(200))                         942-944
+  M_NAV_MENU 'benefitlife' "Benefit" MASTER TREATY URUTAN 12                                          945
 M_BENEFIT_LIFE_SEQ tetap; PEGA_M_BENEFIT_LIFE INVALID (K3); objek M_BENEFIT / T_BENEFIT / dst. tidak disentuh
```

```text
942  blok ALL_VIEWS: DROP VIEW (hanya bila VIEW) ; blok ALL_TAB_COLUMNS sumber: ALTER TABLE M_BENEFIT_LIFE RENAME TO BENEFIT_LIFE
943  ALTER TABLE BENEFIT_LIFE ADD (BENEFIT VARCHAR2(200))                         -- berdiri sendiri
944  blok: isi BENEFIT = m.JSONDATA.Benefit (teks view) ; blok: PEMERIKSAAN K2 (UPDATE … SET m.ID = NULL WHERE
     JSON ber-Benefit AND BENEFIT NULL/berbeda -> ORA-01407) ; blok: DROP COLUMN JSONDATA CASCADE CONSTRAINTS
_down 943 + 942: pelindung UPDATE … SET JSONDATA = JSONDATA WHERE 1 = 0 (ORA-00904) ; 944: blok aman diulang
```

Modul (versi ririsklife TANPA ringkasan / rincian / unggah / Delete - XML tidak memuatnya): satu layar
`InboxBenefit` - form Number / ID (disabled) + Benefit (textarea wajib, huruf besar) dengan Save dan Cancel (saat Edit),
grid ID / Benefit 10 baris, ID menurun, Edit per baris. API `/api/benefit-life`: GET, POST (201), PUT. Paritas:
`docs/PARITAS-LAYAR-DAN-AKSI.md`. Urutan WO (cadangan P3 di paling atas): `docs/LANGKAH-WO-BENEFITLIFE.md` +
`docs/sql/benefit_*.sql`.

## Perlu tinjauan tim inti

- **`inti/backend/migrasi/perintah_katalog.go` TIDAK diubah.** Semua blok 942-944 memakai bentuk yang sudah ada:
  ALL_VIEWS (DROP VIEW) dan RENAME berpelindung tabel sumber (935/938), ALL_TAB_COLUMNS `n > 0` untuk isi / buang.
  Pemeriksaan K2 ("berhenti sebelum DROP JSONDATA") ditulis sebagai blok ALL_TAB_COLUMNS JSONDATA yang perintahnya
  `UPDATE {skema}.BENEFIT_LIFE m SET m.ID = NULL WHERE m.JSONDATA.Benefit IS NOT NULL AND (m.BENEFIT IS NULL OR
  m.BENEFIT <> m.JSONDATA.Benefit)`: menyebut tabel dan objek yang ditanyakan (aturan `pelanggaranBlokPLSQL`), tanpa
  tanda kutip (`polaPerintahKatalog` `[^']+`), nol baris = nol perubahan, satu baris = ORA-01407 (ID NOT NULL) - galat
  yang TIDAK ditoleransi pelari, jadi 944 berhenti sebelum blok DROP dan tidak tercatat. Bentuk ber-`JSON_EXISTS(…,
  '$.Benefit')` ditolak pengurai (diuji `TestBentukPemeriksaanK2`); karena itu "punya Benefit" dibaca lewat notasi
  titik yang sama dengan view, dan `benefit_bukti_k2.sql` kueri D menghitung JSON_EXISTS secara terpisah.
- Penjaga: `migrasi_test.go` (pengecualian `DROP COLUMN JSONDATA` + 944), `rentang_test.go` (urutan pelari 942-945),
  `menu_test.go` (`modulLuarKorpus` + `labelTampilDisetujui` `benefitlife` → folder "Benefit Life", tampil "Benefit");
  `cmd/api/gerbang_tulis_test.go` (HakLihat); `frontend/katalogKorpus.ts`, `frontend/daftar.menuTabel.test.ts`,
  `frontend/Shell.test.ts` (44 kelompok); `inti/backend/daftar/modul_benefitlife_gen.go`.
- 944_down memulihkan constraint bernama asli `ENSURE_M_BENEFIT_LIFE_JSON` (26 byte) dan kunci JSON `Benefit` huruf
  persis (view 942_down membacanya).

## Pembaca lain

Tidak ada (fakta WO 08-10-2026): tidak ada modul repo yang membaca `BENEFIT_LIFE` / `M_BENEFIT_LIFE`, tidak ada
dependensi DB selain view dan prosedur Pega, tidak ada grant. Medan `Benefit` di `masterproductnamelife` teks bebas -
tidak disentuh.

## Evidence

- `go test` modul: models (rumus ID K3, NormalBenefit), services (tiruan: Add 100012, Edit, wajib, K3 ID terpakai →
  ErrIDTerpakai, balapan PK, View only), handlers (201 / 200 / 422 / 403 / 404 / 400 / 409 / 503, nol DELETE),
  repository (SQL, migrasi 942-945 lewat pengurai produksi pelari, bentuk pemeriksaan K2). `-tags=db` (terkompilasi,
  BELUM dijalankan - tanpa skema uji): `TestDBEkspresiSamaDenganView`, `TestDBMigrasiSatuTabelDanMundur`,
  `TestDBPemeriksaanK2Berhenti`, `TestDBLayanan`. vitest modul (label XML, aturan, layar, klien) + gaya.
- Bukti K2 di DEV: `docs/sql/benefit_bukti_k2.sql` - **menunggu hasil WO** (executor tanpa koneksi Oracle).

## Merge Danger

**Door:** one-way sesudah `-migrate` 944 (JSONDATA dibuang). Jalur mundur membangun ulang JSONDATA `{"Benefit": …}`,
nama lama, dan view; `Number`, `pxObjClass`, `px*` hanya dari cadangan P3.

**Blast Radius:** satu tabel, nol pembaca lain; Pega Benefit berhenti dapat menyimpan (prosedur INVALID, K3 diterima).
