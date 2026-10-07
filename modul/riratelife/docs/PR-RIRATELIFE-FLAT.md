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
  uji MPNL/MCRL hijau; `go vet -tags=db` bersih untuk semua paket yang disentuh cabang ini (uji Oracle `rirl_db_test.go` termasuk 928
  diulang sebagian dan jalur mundur - belum dijalankan); `go test ./...` dan `npm test` = baseline.

- **Baseline `go vet -tags=db ./...` (bukan dari cabang ini, jangan diperbaiki di sini - milik premiumlistlife):**
  `modul/premiumlistlife/backend/repository/polis_kasus_db_test.go:53:86: not enough arguments in call to
  kerja.SisipKasusBaru` - galat yang sama ada di `origin/dev` @ `2d9e0692` (diperiksa 07-10-2026 di worktree sementara;
  berkas terakhir diubah `8f0ea22a`, yang ada di `origin/dev`).

## Perlu tinjauan tim inti

Keputusan WO 07-10-2026: pengecualian penjaga `TestKolomUangDesimalDanNolJSON`
(`inti/backend/penjaga/migrasi_test.go`) untuk tepat satu `DROP COLUMN JSONDATA` di
`928_m_rate_life_summary_satu_tabel.sql` - **DISETUJUI WO, tetap ditinjau tim inti** (pernyataan itu membuang JSON, bukan
menambahnya). Berkas lain di luar folder modul:

- Migrasi inti: `922_m_nav_menu_riratelife`, `923_seq_rate_life`, `926_rate_life_summary_flat`,
  `927_m_rate_life_summary_kolom`, `928_m_rate_life_summary_satu_tabel` (masing-masing + `_down`) di
  `inti/backend/migrations/`.
- Penjaga inti: `inti/backend/penjaga/menu_test.go` (`modulLuarKorpus`; aturan nama modul membuang `/`),
  `rentang_test.go` (urutan pelari), `migrasi_test.go` (pengecualian di atas).
- `inti/backend/db/koneksi.go` (`db.Koneksi`, identik dengan cabang ricommlife), `inti/backend/daftar/modul_riratelife_gen.go`.
- `cmd/api/gerbang_tulis_test.go`; frontend akar `katalogKorpus.ts`, `Beranda.test.ts`, `Shell.test.ts`,
  `daftar.menuTabel.test.ts`.
- Modul lain (izin WO 07-10-2026): `modul/masterproductnamelife/**` dan `modul/mastercontractretrolife/**` - pembaca
  ringkasan kini `SELECT ID, USEDBY FROM M_RATE_LIFE_SUMMARY` (baca-saja), uji, tiruan, katalog testdata, dokumen RALAT.

## Merge Danger

**Door:** one-way sesudah `-migrate` 928.

Jalur mundur 928/927 membangun ulang tabel flat dan JSONDATA dari kolom, tetapi FLAG lama dan ringkasan yang dibuang
928 hanya kembali dari cadangan CSV langkah (a). Prosedur `PEGA_M_RATE_LIFE_SUMMARY` / `PEGA_M_PLAN_LIFE_SUMMARY`
menjadi INVALID (Pega tidak lagi menyimpan R/I Rate - diterima WO).

**Blast Radius:** lintas-modul.

R/I Rate Life, `Choose R/I Rate` (Product Name Life), autocomplete `R/I RATE` (Contract Retro Life). Backend lama
(membaca `RATE_LIFE_SUMMARY`) gagal sesudah 928 sampai biner baru dijalankan - restart wajib (langkah (d)).
