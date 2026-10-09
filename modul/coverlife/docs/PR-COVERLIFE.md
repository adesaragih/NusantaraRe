# PR — modul baru Cover Life (`coverlife`, MASTER TREATY) + M_COVER_LIFE flat TANPA RENAME (migrasi modul 085-086, slot 957)

## Summary

```diff
 POOLDATA
-  VIEW  COVER_LIFE   = SELECT ID, a.JSONDATA.Cover, a.JSONDATA.Note FROM M_COVER_LIFE a
-  TABLE M_COVER_LIFE (ID VARCHAR2(10) PK SYS_C009203, JSONDATA IS JSON)
+  TABLE M_COVER_LIFE (ID VARCHAR2(10) PK SYS_C009203, COVER VARCHAR2(200), NOTE VARCHAR2(1000))   085-086 (nama TETAP)
+  M_NAV_MENU 'coverlife' "Cover Life" MASTER TREATY URUTAN 16                                     957
 M_COVER_LIFE_SEQ tetap; PEGA_M_COVER_LIFE INVALID (C2); COVERAGE* / COVERNOTE* tidak disentuh; nol objek COVER_LIFE
```

```text
085  ALTER TABLE M_COVER_LIFE ADD (COVER VARCHAR2(200), NOTE VARCHAR2(1000))           -- berdiri sendiri
086  blok isi (notasi titik = teks view, SELAGI view + JSONDATA ada) ; blok periksa C1.3 ;
     blok DROP COLUMN JSONDATA CASCADE CONSTRAINTS ; blok ALL_VIEWS: DROP VIEW COVER_LIFE (TERAKHIR)
     periksa = UPDATE … SET m.ID = NULL WHERE (COVER beda dari view) OR (NOTE beda dari view)
               -> ORA-01407 (ID NOT NULL), tanpa kutip; NULL = NULL TIDAK menghentikan (NOTE 4/4)
957  INSERT datar M_NAV_MENU … WHERE NOT EXISTS
_down 086: pelindung ORA-00904 (kolom sumber), JSONDATA + IS JSON kembali, CREATE VIEW PERSIS teks DEV ;
      085: pelindung ORA-00904 (JSONDATA) lalu DROP (COVER, NOTE) ; 957: buang hak lalu baris
```

Modul (pola causeoflosslife): form Cover / Note (TANPA medan ID - XML), Save, Cancel (saat Edit), grid 50 baris ID /
Cover / Edit (Note tidak tampil), ID menaik tetap, tanpa saring. API `/api/cover-life`: GET, POST, PUT. Paritas:
`docs/PARITAS-LAYAR-DAN-AKSI.md`. Urutan WO (P3 di paling atas): `docs/LANGKAH-WO-COVERLIFE.md`.

## Perlu persetujuan tim inti

- **Jatah nomor dipinjam (K0, keputusan work owner 08-10-2026)**: `modul/premiumlistlife/MODUL.md` `Rentang migrasi`
  `050-089` → `050-079` (`085-089` ke modul ini, `080-084` ke `diseaselife`); `modul/treatycontractout/MODUL.md` `Slot
  menu` `956-957` → `956-956` (`957` ke modul ini). Nomor terbukti kosong sebelum ditulis (lihat
  `modul/diseaselife/docs/PR-DISEASELIFE.md`). Itu SATU-SATUNYA perubahan di folder treatycontractout.
- **`migrasi_test.go`** (`TestKolomUangDesimalDanNolJSON`): pengecualian `DROP COLUMN JSONDATA` + `086_cover_life_satu_tabel`
  (SATU perintah, persis - pola 092).
- `barisLahirDiSlot` + `modulLuarKorpus` (menu_test.go) bertambah `coverlife`; penjaga URUTAN menurut nilai (lihat PR
  Disease Life) juga melindungi URUTAN 16 di slot 957.
- **`inti/backend/migrasi/perintah_katalog.go` TIDAK diubah.** 086 memakai bentuk blok yang ada (ALL_TAB_COLUMNS +
  ALL_VIEWS). Pemaksa berhenti = bentuk 944 / 092 (`SET m.ID = NULL`, ID NOT NULL sejak Pega).
- Lainnya: `cmd/api/gerbang_tulis_test.go`; `frontend/katalogKorpus.ts`, `daftar.menuTabel.test.ts`, `Shell.test.ts`
  (48); `inti/backend/daftar/modul_coverlife_gen.go`; skema uji `uji/skemauji/cover_tiruan.go` + `skemauji.go` (tiruan
  `M_COVER_LIFE` + view + sequence SEBELUM pelari, dibongkar SESUDAH migrasi mundur).

## Pembaca lain

Tidak ada (fakta WO; dicari ulang di repo). Uji penjaga modul `TestKodeTidakMenyebutViewLama` memastikan kode Go dan TS
modul ini tidak menyebut `COVER_LIFE` sebagai tabel / view (hanya `M_COVER_LIFE` dan sequence-nya).

## Evidence

- `go test` modul: models (rumus ID, wajib / batas, Note opsional, huruf tetap), services (Add 100005, kembar, Edit +
  Note, 404, C2 409, balapan PK, View only, 50 per halaman), handlers (201 / 200 / 422 / 403 / 404 / 400 / 409 / 503, nol
  DELETE), repository (SQL, 085-086 + 957 lewat pengurai produksi, isi = teks view, DROP VIEW terakhir, pemaksa tanpa
  kutip, logika NULL = NULL, urutan pelari 085 → 900 → 957 atas migrasi inti SUNGGUHAN, nol RENAME / CREATE TABLE).
  `-tags=db` (`TestDBEkspresiSamaDenganView`, `TestDBMigrasiSatuTabelDanMundur`, `TestDBPemeriksaanC13`,
  `TestDBLayanan`, `TestDBUrutanPelari085Lalu900Lalu957`) terkompilasi, BELUM dijalankan. vitest modul + gaya.
- Bukti C1 / C2 di DEV: `docs/sql/cov_bukti.sql` - **menunggu hasil WO**.

## Merge Danger

**Door:** one-way sesudah 086 (JSONDATA dan view dibuang). Mundur membangun ulang JSONDATA kunci `Cover` / `Note` huruf
Pega dan view PERSIS; `pxObjClass` / `pyRuleHarness` asli hanya dari cadangan P3. **Blast Radius:** tidak ada pembaca
lain; Pega Cover berhenti dapat menyimpan (C2 diterima); premiumlistlife / treatycontractout kehilangan nomor yang tidak
pernah dipakainya.
