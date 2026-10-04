# Modul `masterdata` — Master Data

Modul di LUAR dua puluh folder korpus (`docs/bersama/PANDUAN-TIM-PER-MODUL.md` bab 5), dibuat atas keputusan work
owner 04-10-2026: "biar data nya ke update km buatin menu master dengan nama masing2 tabel tadi itu. misal manu master
province, isi nya untu insert, update, aktif, non aktif kan data dari master itu" — penempatan "Modul baru
'masterdata'", cakupan "Yang dipakai NB FacIn dulu".

Satu folder, satu modul, satu pemilik: kode backend, kode frontend, dan dokumen modul ini tinggal di sini.

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu. Ubah nilainya hanya
lewat pull request yang disetujui tim inti — dua modul tidak boleh berbagi nomor.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `masterdata` |
| Folder korpus | `Master Data` |
| GROUPMENU | `MASTER` |
| Pemilik | `@PEMILIK-MASTERDATA` |
| Status | dimigrasi |
| Rentang migrasi | `760-799` |
| Slot menu | `990-991` |
| Prefix rute API | `/api/masterdata` |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | rencana dan tiket (`issues/NN-<slug>.md`), `STRUKTUR-TABEL-MASTER-DATA.md` |
| `backend/` | `models/` (daftar master + kolom), `repository/`, `services/`, `handlers/` (`GET /api/masterdata`, `GET` / `POST /api/masterdata/{tabel}`, `PUT /api/masterdata/{tabel}/{id}`, `PUT /api/masterdata/{tabel}/{id}/status`), `migrations/` (760 view → tabel flat + `STS_AKTIF`, 761 `T_MASTER_STATUS` (status NATION), 762 jejak ubah `CREATE_OP` / `TGL_CREATE` / `UPDATE_OP` / `TGL_UPDATE`, slot menu 990), `modul.go` |
| `frontend/` | `menu.ts`, `rute.tsx`, `api.ts`, `labels.ts`, `masterdata.css`, `pages/HalamanMasterData.tsx` (tab per master), `components/` (`DaftarMaster`, `FormMaster`, `PilihRujukan`) — sesi 0f |

## Urutan deploy

⛔ **Urutan deploy** (04-10-2026, diteruskan sesi nusantarare-0f): kode modul ini dan saran akumulasi nbfacin
(MD-5: `STS_AKTIF`, `CITYINPUT`, `DISTRICTINPUT`, `T_MASTER_STATUS`) membaca tabel 760 / 761 / 762 dan gagal
(ORA-00904 / ORA-00942) atas skema yang belum dimigrasi. Urutannya:

1. Pasang biner yang memuat commit `072370f0` (pelari: pra-terbang melewati VIEW yang diganti tabel bernama sama).
   Biner lebih lama BERHENTI di pra-terbang 760 selama keenam objek masih VIEW — dan menahan seluruh migrasi
   tertunda modul lain.
2. `-migrate` (work owner): 904 → 760 → 761 → 762 → 990 (urut nama, ikut berjalan bersama migrasi tertunda lain).
3. Baru sesudahnya buka layar Master Data dan saran akumulasi nbfacin (Choose Accumulation, tiket 46).

**Bila `-migrate` gagal.** Pesan yang memuat `BENTUKNYA BERBEDA ... Migrasi dihentikan sebelum satu pernyataan
pun dikirim` = gagal di **pra-terbang**: tidak ada yang berubah, perbaiki sebabnya lalu ulangi. Galat lain atas
`760_view_ke_tabel_flat` = gagal **sesudah** pra-terbang: DDL Oracle tanpa transaksi, jadi pernyataan sebelum yang gagal
SUDAH jadi dan `T_MIGRASI` belum mencatat 760. ⛔ **Jangan langsung mengulang `-migrate`** — periksa dulu sisa
`*_SALIN`. 760 mengubah keenam view berurutan (PROVINCE, CITYINPUT, DISTRICTINPUT, ACCUMULATEDTYPE, CZONE,
ACCUMULATION), masing-masing: `CREATE TABLE X_SALIN` → salin isi view → `DROP VIEW X` → `CREATE TABLE X` → salin
balik → `DROP TABLE X_SALIN`. Untuk tiap X, baca
`SELECT OBJECT_NAME, OBJECT_TYPE FROM USER_OBJECTS WHERE OBJECT_NAME IN ('X', 'X_SALIN')`:

| `X` | `X_SALIN` | Artinya | Sebelum mengulang |
| --- | --- | --- | --- |
| VIEW | tidak ada | belum disentuh | — |
| VIEW | ada | berhenti sebelum `DROP VIEW`; isi masih di view | buang `X_SALIN` — tanpa itu pengulangan menyalin isi DUA KALI (pelari melewati `CREATE` yang sudah ada, `INSERT`-nya tetap jalan) |
| tidak ada | ada | ⛔ satu-satunya salinan isi ada di `X_SALIN` | JANGAN dibuang: jalankan manual `CREATE TABLE X` + salin balik persis teks 760, cocokkan `COUNT(*)`, baru `DROP TABLE X_SALIN` |
| TABLE | ada | berhenti di salin balik atau sesudahnya | cocokkan `COUNT(*)` X dengan `X_SALIN`; bila X kosong jalankan salin balik; lalu `DROP TABLE X_SALIN` |
| TABLE | tidak ada | view itu selesai | — |

Begitu satu view saja sudah menjadi TABLE, 760 **tidak dapat diulang pelari**: `DROP VIEW` atas objek yang kini
tabel gagal `[dugaan: ORA-00942]`. Jalan keluarnya: DBA menuntaskan view yang tersisa manual dari teks 760,
memastikan keenam tabel ber-`STS_AKTIF`, mencatat `INSERT INTO <skema>.T_MIGRASI (NAMA, DIJALANKAN_PADA) VALUES
('760_view_ke_tabel_flat', SYSDATE)` (`NAMA` = nama berkas tanpa `.sql`, `migrasi.KunciLangkah`), lalu `-migrate`
lagi untuk 761 dan seterusnya. 762 yang gagal di tengah berhenti di ORA-01430 saat diulang (kolom sudah ada) — pola
059: periksa kolom yang sudah ada, bukan ulang buta. ⛔ Seluruh langkah di atas perubahan basis data: work owner /
DBA, dengan persetujuan — bukan agent.

## Menjalankan uji modul ini saja

Dari folder `APP_RNM/`:

```powershell
go test ./modul/masterdata/...
npx vitest run modul/masterdata
```

## Pernyataan untuk penjaga

⛔ **Dibaca penjaga** `inti/backend/penjaga` — satu jenis pernyataan per judul `###`, satu baris per butir.

### Tabel warisan: dibaca, tidak dibuat

Tabel yang `docs/STRUKTUR-TABEL-MASTER-DATA.md` gambarkan tetapi SENGAJA tidak dibuat migrasi mana pun
(`TestKolomDDLCocokDenganStruktur`, `TestTabelBukanMilikKitaTidakDibuat`). Mencabut satu baris = kepemilikan tabel
berpindah — keputusan work owner.

| Tabel | Alasan |
| --- | --- |
| `NATION` | tabel negara warisan POOLDATA; Master Data menulis isinya (menu Nation) tanpa mengubah strukturnya - status dan jejak ubahnya di T_MASTER_STATUS (761 / 762); nbfacin membacanya untuk saran Country (tiket 46) |
| `OBJECTITEMTYPE` | tabel jenis item objek warisan POOLDATA (sumber view V_JN_OBJ_ITEM); Master Data menulis isinya dan memakai ISACTIVE-nya |
| `BRANCH` | tabel cabang warisan POOLDATA; Master Data hanya MEMBACA ID untuk rujukan City |
