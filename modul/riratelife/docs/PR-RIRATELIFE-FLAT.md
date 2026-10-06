# PR — `modul/riratelife/implementasi`: ringkasan R/I Rate Life jadi tabel flat (RALAT R4)

> ⚠️ **Urutan merge: PR ini LEBIH DULU, lalu `modul/ricommlife/implementasi`** (`modul/ricommlife/docs/PR-RICOMMLIFE.md`).
> Kedua cabang membawa `inti/backend/db/koneksi.go` identik (`91221b2f` = `349be34d`); uji coba merge ricommlife di atas
> ujung cabang ini bersih (MODUL.md bab "Urutan merge").

## Summary

```diff
 RATE_LIFE_SUMMARY            (POOLDATA)
-  VIEW  SELECT a.ID, a.JSONDATA.USEDBY, … FROM M_RATE_LIFE_SUMMARY a
+  TABLE ID VARCHAR2(10) PK, USEDBY, TYPE, MODIFIEDDATE, OPERATORID, FLAG   (migrasi inti 926)
+  INDEX IX_RATE_LIFE_SUMMARY_NAMA (UPPER(TRIM(USEDBY)))                     (pemeriksa nama + pengaman view)
 M_RATE_LIFE_SUMMARY          (JSON)   tidak disentuh: cadangan + sumber alat pindah + ID terpakai
```

```diff
 riratelife ringkasan
-  baca  view RATE_LIFE_SUMMARY ; tulis M_RATE_LIFE_SUMMARY.JSONDATA (JSON_OBJECT / baca-ubah-tulis)
+  baca + tulis tabel flat RATE_LIFE_SUMMARY (kolom bernama; TYPE/FLAG tidak ditulis, Edit tidak menimpa)
+  ID baru: SEQ_M_RATE_LIFE_SUMMARY, lewati ID terpakai di flat ∪ M_RATE_LIFE_SUMMARY
 riratelife rincian (M_RATE_LIFE + view RATE_LIFE)   tidak berubah
 mastercontractretrolife, masterproductnamelife      tidak diubah - SELECT ID, USEDBY kini dari tabel
```

```text
alat pindahflat [-jalankan] [-sejak=<MODIFIEDDATE>]
  tolak IS_PEGA_PROD=true sebelum koneksi
  satu db.Koneksi: ALTER SESSION → transaksi → LOCK TABLE
  per ID: baru | berubah (sumber lebih baru) | sama | konflik (tidak ditimpa) | dilewati (≤ -sejak, dihapus aplikasi)
  tolak bila nilai tidak muat (laporan panjang maksimum per kolom) - nol pemotongan
```

Urutan WO: `docs/LANGKAH-WO-RIRATELIFE-FLAT.md` (DBA lepas view → `-migrate` → pindah penuh → delta → verifikasi
COUNT + MINUS dua arah → restart).

## Evidence

- **Before:** `go test ./...` = baseline (3 paket claimlife merah); grid membaca view atas JSON.
  **After:** `go test ./...` = baseline yang sama + paket riratelife hijau (`TestMigrasi926TeruraiDanBerpasangan`,
  `TestSqlTulisRingkasan`, `TestRencanaPindahPenuh`, `TestRencanaPindahDelta`, `TestNilaiJSON`); `go vet -tags=db`
  bersih (uji Oracle `rirl_db_test.go` ditulis, belum dijalankan); `npm test` = baseline 14 gagal yang sama.

## Merge Danger

**Door:** one-way sesudah `-migrate` + aplikasi menulis ringkasan baru.

Sebelum aplikasi menulis, jalur mundur `926_…_down.sql` memulihkan view persis. Sesudahnya, ringkasan baru hanya ada di
tabel flat (hilang bila view dipulihkan).

**Blast Radius:** lintas-modul.

R/I Rate Life, autocomplete `R/I RATE` (Contract Retro Life), `Choose R/I Rate` (Product Name Life). Ringkasan yang
ditulis Pega sesudah pemindahan tidak tampil sampai delta dijalankan. Arti FLAG `AP`/`PM`/`PY` belum diketahui
(pertanyaan WO, MODUL.md).
