# PR — `modul/riratelife/implementasi`: ringkasan R/I Rate Life satu tabel `M_RATE_LIFE_SUMMARY` (RALAT R6)

> ⚠️ **Urutan merge: PR ini LEBIH DULU, lalu `modul/ricommlife/implementasi`** (`modul/ricommlife/docs/PR-RICOMMLIFE.md`).
> Kedua cabang membawa `inti/backend/db/koneksi.go` identik (`91221b2f` = `349be34d`); uji coba merge ricommlife di atas
> ujung cabang ini: satu konflik kecil di `rentang_test.go` (daftar urutan pelari 923, 924, 925, 926, 927, 928 - MODUL.md
> bab "Urutan merge"). Nama berkas tetap `PR-RIRATELIFE-FLAT.md` supaya rujukan lama tidak putus.

## Summary

```diff
 POOLDATA
 M_RATE_LIFE_SUMMARY
-  ID VARCHAR2(10) PK, JSONDATA CLOB (ENSURE_M_RATE_LIFE_SUMMARY_JSON: IS JSON)
+  ID VARCHAR2(10) PK, USEDBY VARCHAR2(500), TYPE VARCHAR2(100), MODIFIEDDATE VARCHAR2(50), OPERATORID VARCHAR2(200)
+  INDEX IX_M_RATE_LIFE_SUMMARY_NAMA (UPPER(TRIM(USEDBY)))
-RATE_LIFE_SUMMARY  (tabel flat 926 - dibuang 928, isinya pindah ke kolom di atas)
 M_RATE_LIFE + view RATE_LIFE (rincian)  tidak berubah
```

```text
927  ALTER TABLE M_RATE_LIFE_SUMMARY ADD (USEDBY, TYPE, MODIFIEDDATE, OPERATORID)      -- berdiri sendiri
928  UPDATE … dari RATE_LIFE_SUMMARY                                                    -- semua aman diulang
     DELETE baris yang sudah dihapus aplikasi
     blok katalog: DROP COLUMN JSONDATA CASCADE CONSTRAINTS
     INSERT baris flat-saja NOT EXISTS
     CREATE INDEX IX_M_RATE_LIFE_SUMMARY_NAMA
     DROP TABLE RATE_LIFE_SUMMARY                                                       -- terakhir
```

```diff
 riratelife ringkasan
-  baca + tulis tabel flat RATE_LIFE_SUMMARY; ID terpakai: flat ∪ JSON; alat pindahflat
+  baca + tulis kolom M_RATE_LIFE_SUMMARY; ID terpakai: satu tabel; alat pindahflat + DBA lepas view dihapus
 masterproductnamelife (Choose R/I Rate), mastercontractretrolife (R/I RATE)
-  SELECT ID, USEDBY FROM RATE_LIFE_SUMMARY
+  SELECT ID, USEDBY FROM M_RATE_LIFE_SUMMARY          (baca-saja, diizinkan WO 07-10-2026)
```

Urutan WO: `docs/LANGKAH-WO-RIRATELIFE-SATU-TABEL.md` (cadangan CSV ID+JSONDATA → `-migrate` → verifikasi 5 kolom,
COUNT = flat sebelum, prosedur PEGA_* INVALID → restart; pemulihan per titik gagal; jalur mundur).

## Evidence

- **Before:** ringkasan di dua tempat (JSON Pega + tabel flat 926).
  **After:** `TestMigrasi927KolomRingkasan`, `TestMigrasi928SatuTabel` (pengurai produksi pelari: `KolomAlterTambah`,
  `BacaPerintahKatalog`, `KolomAlterBuang`, `KolomCreateTable` jalur mundur; lebar = 926), `TestSqlTulisRingkasan`,
  uji MPNL/MCRL hijau; `go vet -tags=db ./...` bersih (uji Oracle `rirl_db_test.go` termasuk 928 diulang sebagian dan
  jalur mundur - belum dijalankan); `go test ./...` dan `npm test` = baseline.

## Merge Danger

**Door:** one-way sesudah `-migrate` 928.

Jalur mundur 928/927 membangun ulang tabel flat dan JSONDATA dari kolom, tetapi FLAG lama dan ringkasan yang dibuang
928 hanya kembali dari cadangan CSV langkah (a). Prosedur `PEGA_M_RATE_LIFE_SUMMARY` / `PEGA_M_PLAN_LIFE_SUMMARY`
menjadi INVALID (Pega tidak lagi menyimpan R/I Rate - diterima WO).

**Blast Radius:** lintas-modul.

R/I Rate Life, `Choose R/I Rate` (Product Name Life), autocomplete `R/I RATE` (Contract Retro Life). Backend lama
(membaca `RATE_LIFE_SUMMARY`) gagal sesudah 928 sampai biner baru dijalankan - restart wajib (langkah (d)).
