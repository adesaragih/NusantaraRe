# 04: City — EMAIL / MOID keluar dari menu, PROVINCENAME, dan kota baru dari RW (backend sesi c3)

> ⚠️ **Disusun agent sesi c3 atas permintaan sesi `nusantarare-9d` — bukan hasil `/to-tickets`.** Tanpa commit / push.

**Status:** ditulis 04-10-2026, uji Go hijau. ⛔ Migrasi 883 ditulis, **belum dijalankan** (work owner).

## Permintaan work owner (04-10-2026, diteruskan sesi 9d)

- "di menu city kolom email dan mo id di hapus saja"
- "di cityinput tambahin PROVINCENAME, data nya diambil dari tabel rw . rw.CITYNAME=cityinput.note baru di isi
  cityinput .PROVINCENAME"
- "cek di tabel rw, ada kolom cityname …, bandingkan ke tabel city, jika belum ada di tabel city, km tambahin, untuk id
  nya +1 dari id terakhir. cityname=note" — lalu "tidak usah zipcode, provincename aja"

## Yang dikerjakan

- `inti/backend/master/models/daftar.go` (TIM INTI): master `city` tanpa `email` / `moId`; `provinceName`
  (`PROVINCENAME`, 4000) kolom biasa sesudah `provinceId`. Kolom DB EMAIL / MOID TETAP: tambah = NULL, ubah = tidak
  disentuh (SQL hanya menyebut kolom yang didefinisikan) — `TestCityTanpaEmailMoID`.
- `modul/masterprovince/backend/migrations/883_cityinput_provincename.sql` (+ `_down`): ADD → UPDATE → INSERT;
  `backend/migrasi_test.go` `TestMigrasiProvinceNameCity` (bentuk SQL; uji mutasi: UPPER di UPDATE → merah).
- `docs/STRUKTUR-TABEL-MASTER-DATA.md` bab CITYINPUT. Label `provinceName` = 'Province Name' di
  `inti/frontend/master/labels.ts` — sesi 9d.

## Keputusan agent (menunggu konfirmasi work owner)

| # | Keputusan | Alasan |
| --- | --- | --- |
| MD-11 | PROVINCENAME diisi hanya bila baris RW ber-`CITYNAME = NOTE` (persis, tanpa UPPER / TRIM) memuat TEPAT SATU nama provinsi berbeda; selain itu NULL. NULL di RW tidak dihitung (`COUNT(DISTINCT)`). `RW.STS_AKTIF` TIDAK disaring. Isi yang sudah ada tidak ditimpa | kata work owner tanpa saringan status; baris RW nonaktif tetap pemetaan kota → provinsi, dan nilai lama yang berbeda hanya membuat kota itu NULL (ambigu), bukan salah isi. Alternatif: `AND r.STS_AKTIF = '1'` (filter RD BrowseRW) |
| MD-12 | Kota baru: satu baris per CITYNAME RW berbeda yang TRIM-nya tidak kosong dan tidak sama PERSIS dengan NOTE mana pun; NOTE = CITYNAME; PROVINCENAME aturan MD-11; STS_AKTIF bawaan '1'; kolom lain NULL. ID = MAX ID yang seluruhnya digit (`REGEXP_LIKE '^[0-9]+$'`) + `ROW_NUMBER()` urut CITYNAME, teks tanpa nol di depan; tanpa ID digit mulai 1. Penanda `CREATE_OP = 'MIGRASI-883'` + `TGL_CREATE` = SYSDATE; jalur mundur membuang tepat baris itu | "+1 dari id terakhir" atas kolom teks; ID non-digit tidak dapat ditambah 1. Penanda = satu-satunya cara jalur mundur tidak salah hapus |

⚠️ **Risiko pencocokan persis** `[dugaan, jumlahnya belum punya data]`: CITYNAME RW yang hanya beda huruf besar / spasi
dari NOTE yang ada ("Jakarta" vs "JAKARTA") menjadi KOTA BARU kembar. Sebelum `-migrate`, work owner dapat melihat
jumlahnya (SELECT saja, tanpa mengubah apa pun):

```sql
-- calon kota baru
SELECT COUNT(DISTINCT r.CITYNAME) FROM <skema>.RW r
 WHERE TRIM(r.CITYNAME) IS NOT NULL
   AND NOT EXISTS (SELECT 1 FROM <skema>.CITYINPUT x WHERE x.NOTE = r.CITYNAME);
-- ... yang sebenarnya kembar beda huruf / spasi dengan kota yang sudah ada
SELECT COUNT(DISTINCT r.CITYNAME) FROM <skema>.RW r
 WHERE TRIM(r.CITYNAME) IS NOT NULL
   AND NOT EXISTS (SELECT 1 FROM <skema>.CITYINPUT x WHERE x.NOTE = r.CITYNAME)
   AND EXISTS (SELECT 1 FROM <skema>.CITYINPUT x WHERE UPPER(TRIM(x.NOTE)) = UPPER(TRIM(r.CITYNAME)));
-- kota yang PROVINCENAME-nya akan NULL karena ambigu
SELECT r.CITYNAME, COUNT(DISTINCT r.PROVINCENAME) FROM <skema>.RW r
 GROUP BY r.CITYNAME HAVING COUNT(DISTINCT r.PROVINCENAME) > 1;
```

## Bila 883 gagal

DDL tanpa transaksi: gagal SESUDAH ALTER → pengulangan berhenti di ORA-01430 dan `T_MIGRASI` belum mencatat 883. DBA
menjalankan pernyataan yang tersisa (UPDATE, INSERT) manual — keduanya aman diulang (UPDATE tidak menimpa; INSERT
melewati NOTE yang sudah ada) — lalu mencatat `883_cityinput_provincename` di `T_MIGRASI`. ⛔ Perubahan basis data:
work owner / DBA, dengan persetujuan.

## Belum teruji

SQL 883 hanya diuji bentuknya (tanpa Oracle): `REGEXP_LIKE`, `ROW_NUMBER` sesudah `WHERE NOT EXISTS`, dan `CASE` dalam
`MAX` `[belum terverifikasi]` di DEV.
