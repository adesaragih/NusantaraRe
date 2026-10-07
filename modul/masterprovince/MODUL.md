# Modul `masterprovince` — Province

Modul di LUAR dua puluh folder korpus (`docs/bersama/PANDUAN-TIM-PER-MODUL.md` bab 5). Keputusan work owner
04-10-2026 (diteruskan sesi nusantarare-9d): Master Data dipecah menjadi delapan modul menu terpisah di grup MASTER —
"ini masing dibuatin didalam sub menu master data, bukan di jadiin 1 gini", cara pecah "8 modul terpisah". Nama tampilan
"Province" (tanpa kata Master, pola Product Name Life). Mesin layar dan API-nya BERSAMA, milik tim inti:
`inti/backend/master` dan `inti/frontend/master` ("Pindah ke inti"). Rencana, keputusan, dan kontrak rute:
`modul/masterprovince/docs/issues/` (01 rencana, 02 API, 03 pemecahan delapan modul).

Satu folder, satu modul, satu pemilik.

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu. Ubah nilainya hanya
lewat pull request yang disetujui tim inti — dua modul tidak boleh berbagi nomor. `Slot menu` — = modul luar korpus
tanpa slot (keputusan work owner 04-10-2026 "Modul luar korpus tanpa slot"): barisnya lahir DIMIGRASI '1' di migrasi
inti 913.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `masterprovince` |
| Folder korpus | `Master Province` |
| GROUPMENU | `MASTER` |
| Pemilik | `@PEMILIK-MASTERPROVINCE` |
| Status | dimigrasi |
| Rentang migrasi | `215-219` |
| Slot menu | — |
| Prefix rute API | `/api/master-province` |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

## Isi folder

| Folder | Isi |
| --- | --- |
| `backend/` | `modul.go` — `master.Pendaftaran(Nama, "/api/master-province", "province", …)`; definisi master `province` di `inti/backend/master/models/daftar.go`; `migrations/` 215 / 216 / 217, 218 (CITYINPUT.PROVINCENAME + kota baru dari RW, `docs/issues/04`) (+ `_down`) — dulu 880–883 |
| `frontend/` | `menu.ts`, `rute.tsx`, `labels.ts` — layar `inti/frontend/master` (sesi 9d) |
| `docs/` | `STRUKTUR-TABEL-MASTER-DATA.md`, `issues/01-rencana-master-data.md`, `02-api-master-data.md`, `03-pemecahan-delapan-modul.md` (warisan modul masterdata), `04-city-provincename.md` |

## Tabel

Menulis `PROVINCE`. **Pemilik migrasi 215-217** (dulu 880-882) (keputusan work owner 04-10-2026 "Dihapus, tabel pindah ke Province"): keenam tabel flat `PROVINCE`, `CITYINPUT`, `DISTRICTINPUT`, `ACCUMULATEDTYPE`, `CZONE`, `ACCUMULATION` (215, view warisan -> tabel bernama sama + `STS_AKTIF`), `T_MASTER_STATUS` (216) dan kolom jejak ubah (217) - dibaca / ditulis juga oleh tujuh modul master lain dan dibaca saran akumulasi `nbfacin`. Struktur: `docs/STRUKTUR-TABEL-MASTER-DATA.md`.

## Urutan deploy

⛔ **Urutan deploy** (04-10-2026, diteruskan sesi nusantarare-0f; diperbarui untuk delapan modul master): kedelapan
modul master (mesin `inti/backend/master`) dan saran akumulasi nbfacin (MD-5: `STS_AKTIF`, `CITYINPUT`, `DISTRICTINPUT`,
`T_MASTER_STATUS`) membaca tabel 215 / 216 / 217 dan gagal (ORA-00904 / ORA-00942) atas skema yang belum dimigrasi.
Backend dan frontend delapan modul dideploy BERSAMA (`/api/masterdata` dibuang). Urutannya:

1. Pasang biner yang memuat commit `072370f0` (pelari: pra-terbang melewati VIEW yang diganti tabel bernama sama).
   Biner lebih lama BERHENTI di pra-terbang 215 selama keenam objek masih VIEW — dan menahan seluruh migrasi
   tertunda modul lain.
2. `-migrate` (work owner): 215 → 216 → 217 → 218 → 912-919 (baris menu delapan modul, DIMIGRASI '1') → 920 (hak
   menu Master Data disalin ke delapan menu, lalu baris 'masterdata' dibuang). Urut nama, ikut berjalan bersama migrasi
   tertunda lain. ⛔ Di DEV yang sudah menjalankan nama lama (760-762 / 880-883): ganti namanya di `T_MIGRASI` LEBIH
   DULU (bab "Penomoran ulang (05-10-2026)"), baru `-migrate`.
3. Baru sesudahnya buka layar delapan master dan saran akumulasi nbfacin (Choose Accumulation, tiket 46). Bila
   `MODUL_AKTIF` di env menyebut nama modul satu per satu, tambahkan kedelapan nama dan buang `masterdata`.

**Bila `-migrate` gagal.** Pesan yang memuat `BENTUKNYA BERBEDA ... Migrasi dihentikan sebelum satu pernyataan
pun dikirim` = gagal di **pra-terbang**: tidak ada yang berubah, perbaiki sebabnya lalu ulangi. Galat lain atas
`215_view_ke_tabel_flat` = gagal **sesudah** pra-terbang: DDL Oracle tanpa transaksi, jadi pernyataan sebelum yang gagal
SUDAH jadi dan `T_MIGRASI` belum mencatat 215. ⛔ **Jangan langsung mengulang `-migrate`** — periksa dulu sisa
`*_SALIN`. 215 mengubah keenam view berurutan (PROVINCE, CITYINPUT, DISTRICTINPUT, ACCUMULATEDTYPE, CZONE,
ACCUMULATION), masing-masing: `CREATE TABLE X_SALIN` → salin isi view → `DROP VIEW X` → `CREATE TABLE X` → salin
balik → `DROP TABLE X_SALIN`. Untuk tiap X, baca
`SELECT OBJECT_NAME, OBJECT_TYPE FROM USER_OBJECTS WHERE OBJECT_NAME IN ('X', 'X_SALIN')`:

| `X` | `X_SALIN` | Artinya | Sebelum mengulang |
| --- | --- | --- | --- |
| VIEW | tidak ada | belum disentuh | — |
| VIEW | ada | berhenti sebelum `DROP VIEW`; isi masih di view | buang `X_SALIN` — tanpa itu pengulangan menyalin isi DUA KALI (pelari melewati `CREATE` yang sudah ada, `INSERT`-nya tetap jalan) |
| tidak ada | ada | ⛔ satu-satunya salinan isi ada di `X_SALIN` | JANGAN dibuang: jalankan manual `CREATE TABLE X` + salin balik persis teks 215, cocokkan `COUNT(*)`, baru `DROP TABLE X_SALIN` |
| TABLE | ada | berhenti di salin balik atau sesudahnya | cocokkan `COUNT(*)` X dengan `X_SALIN`; bila X kosong jalankan salin balik; lalu `DROP TABLE X_SALIN` |
| TABLE | tidak ada | view itu selesai | — |

Begitu satu view saja sudah menjadi TABLE, 215 **tidak dapat diulang pelari**: `DROP VIEW` atas objek yang kini
tabel gagal `[dugaan: ORA-00942]`. Jalan keluarnya: DBA menuntaskan view yang tersisa manual dari teks 215,
memastikan keenam tabel ber-`STS_AKTIF`, mencatat `INSERT INTO <skema>.T_MIGRASI (NAMA, DIJALANKAN_PADA) VALUES
('215_view_ke_tabel_flat', SYSDATE)` (`NAMA` = nama berkas tanpa `.sql`, `migrasi.KunciLangkah`), lalu `-migrate`
lagi untuk 216 dan seterusnya. 217 yang gagal di tengah berhenti di ORA-01430 saat diulang (kolom sudah ada) — pola
059: periksa kolom yang sudah ada, bukan ulang buta. ⛔ Seluruh langkah di atas perubahan basis data: work owner /
DBA, dengan persetujuan — bukan agent.

## Penomoran ulang (05-10-2026) — 880-883 → 215-218

Penjaga `TestRentangMigrasiDanSlotMenuTidakBerbagiNomor` merah sesudah merge `origin/dev` 05-10-2026: modul `aggregate`
(milik tim lain, GitHub) memegang `880-889` sejak 04-10-2026, modul ini memakai `880-883`. Keputusan work owner
05-10-2026 ("ikutin rekomendasi terbaik km aja. biar tidak terganggu kerjaan yang sudah ada milik orang lain"): modul
LOKAL yang pindah. Rentang 001-899 penuh; `215-219` diambil dari ujung rentang `nbfacin` (`180-219` → `180-214`, modul
milik work owner yang sama; nbfacin memakai sampai 199). Isi SQL tidak berubah — hanya nama berkas dan komentar; penanda
data `CREATE_OP = 'MIGRASI-883'` TETAP (baris yang sudah dibuat membawanya, jalur mundur 218 mencarinya).

| Dulu (04-10-2026) | Lebih dulu (masterdata) | Sekarang |
| --- | --- | --- |
| `880_view_ke_tabel_flat` | `760_view_ke_tabel_flat` | `215_view_ke_tabel_flat` |
| `881_t_master_status` | `761_t_master_status` | `216_t_master_status` |
| `882_jejak_ubah_master` | `762_jejak_ubah_master` | `217_jejak_ubah_master` |
| `883_cityinput_provincename` | — | `218_cityinput_provincename` |

Urutan jalan ikut berubah (pelari mengurut NAMA berkas lintas modul): kini sesudah nbfacin 199 dan sebelum rnwfacin 220,
bukan sesudah 879. `[terverifikasi]` tidak ada migrasi modul lain di 220-879 yang memakai PROVINCE / CITYINPUT /
DISTRICTINPUT / ACCUMULATEDTYPE / CZONE / ACCUMULATION / T_MASTER_STATUS (cari `*.sql` 05-10-2026; nbfacin 186 / 192
hanya bernama kolom `PROVINCE`); menu inti 912-920 tetap sesudahnya.

⛔ **Pelari mencatat langkah menurut NAMA berkas** (`T_MIGRASI.NAMA`, `migrasi.KunciLangkah`) — nama baru = langkah baru.
Work owner tidak ingat apakah nama lama sudah dijalankan di DEV, jadi SEBELUM `-migrate` berikutnya work owner / DBA
memeriksa:

```sql
SELECT NAMA FROM <skema>.T_MIGRASI
 WHERE NAMA IN ('760_view_ke_tabel_flat', '761_t_master_status', '762_jejak_ubah_master',
                '880_view_ke_tabel_flat', '881_t_master_status', '882_jejak_ubah_master', '883_cityinput_provincename');
```

Untuk setiap nama lama yang ADA, ganti namanya ke nama sekarang — JANGAN dijalankan ulang:

```sql
UPDATE <skema>.T_MIGRASI SET NAMA = '215_view_ke_tabel_flat' WHERE NAMA IN ('760_view_ke_tabel_flat', '880_view_ke_tabel_flat');
UPDATE <skema>.T_MIGRASI SET NAMA = '216_t_master_status' WHERE NAMA IN ('761_t_master_status', '881_t_master_status');
UPDATE <skema>.T_MIGRASI SET NAMA = '217_jejak_ubah_master' WHERE NAMA IN ('762_jejak_ubah_master', '882_jejak_ubah_master');
UPDATE <skema>.T_MIGRASI SET NAMA = '218_cityinput_provincename' WHERE NAMA = '883_cityinput_provincename';
COMMIT;
```

(⚠️ `880_seq_aggregate` milik modul `aggregate` — JANGAN ikut diubah; pencarian di atas memakai nama lengkap.) Bila
tidak ada satu pun nama lama, tidak perlu apa-apa: 215-218 jalan sebagai langkah baru.

**Bila pemeriksaan terlewat.** Pengulangan nama baru atas skema yang sudah dimigrasi berhenti AMAN di **pra-terbang**
selama 217 (jejak ubah) sudah jalan: PROVINCE dkk. dan T_MASTER_STATUS sudah ber-kolom jejak, bentuknya berbeda dari
CREATE TABLE 215 / 216 → "BENTUKNYA BERBEDA ... Migrasi dihentikan sebelum satu pernyataan pun dikirim" — jalankan
UPDATE di atas lalu ulangi. Hanya bila DEV menjalankan 880 / 881 TANPA 882, pra-terbang lolos dan 215 berhenti di
`DROP VIEW` sesudah membuat `PROVINCE_SALIN` — buang `PROVINCE_SALIN` (salinan ganda; PROVINCE tidak tersentuh), jalankan
UPDATE di atas, ulangi. 218 sendiri aman diulang di bagian isinya (UPDATE ber-`PROVINCENAME IS NULL`, INSERT ber-`NOT
EXISTS`), tetapi `ALTER ... ADD (PROVINCENAME)` berhenti di ORA-01430 bila kolomnya sudah ada — kolom baru TIDAK boleh
dibungkus blok berpelindung katalog (penjaga `pelanggaranBlokPLSQL`: kolom lewat blok tidak terlihat pembanding
STRUKTUR), jadi jalan keluarnya tetap UPDATE `T_MIGRASI` di atas. ⛔ Seluruh langkah di atas perubahan basis data: work
owner / DBA, dengan persetujuan — bukan agent.

## Penomoran ulang (04-10-2026)

Saat merge `origin/dev` (perintah work owner: "ikuti yang dari github, kalau bentrok dengan kerjaan saya
disesuaikan"), modul `masterdata` (kini dipecah; berkasnya di modul ini) memberi jalan kepada `marketingofficer`, yang lebih dulu memegang rentang `760-799` dan slot
`990-991`, dan kepada migrasi inti GitHub `904_m_login_go_kontak` dan `909_m_nav_menu_master_treaty`:

| Dulu | Sekarang |
| --- | --- |
| `760_view_ke_tabel_flat` | `880_view_ke_tabel_flat` |
| `761_t_master_status` | `881_t_master_status` |
| `762_jejak_ubah_master` | `882_jejak_ubah_master` |
| `990_menu_masterdata` | `996_menu_masterdata` |
| `904_m_nav_menu_masterdata` (inti, URUTAN 4) | `911_m_nav_menu_masterdata` (inti, URUTAN 6 — sesudah 909 merapatkan MASTER) |

Isi SQL tidak berubah, kecuali URUTAN baris menu. ⛔ **Pelari mencatat langkah menurut NAMA berkas**
(`T_MIGRASI.NAMA`, `migrasi.KunciLangkah`). Bila nama lama **sudah** tercatat di DB DEV, `-migrate` akan
menganggap nama baru sebagai langkah baru dan mengulangnya — 880 lalu gagal di `DROP VIEW` atas objek yang sudah TABLE.
Sebelum `-migrate` berikutnya, work owner / DBA memeriksa:

```sql
SELECT NAMA FROM <skema>.T_MIGRASI
 WHERE NAMA IN ('760_view_ke_tabel_flat', '761_t_master_status', '762_jejak_ubah_master', '990_menu_masterdata',
                '904_m_nav_menu_masterdata');
```

Untuk setiap nama lama yang ADA, ganti namanya, jangan jalankan ulang:
`UPDATE <skema>.T_MIGRASI SET NAMA = '<nama baru>' WHERE NAMA = '<nama lama>'`. Bila `904_m_nav_menu_masterdata`
sudah jalan, barisnya masih ber-URUTAN 4 (911 idempoten, tidak menimpanya):
`UPDATE <skema>.M_NAV_MENU SET URUTAN = 6, TGL_UBAH = SYSDATE WHERE KODE = 'masterdata'`. Bila tidak ada satu pun,
tidak perlu apa-apa. ⛔ Perubahan basis data: work owner / DBA, dengan persetujuan — bukan agent.

⚠️ **Sesudah pemecahan delapan modul (04-10-2026)**: berkas `996_menu_masterdata` dan inti
`911_m_nav_menu_masterdata` (+ `_down`) DIBUANG — pelari hanya menjalankan berkas yang ada (maju dan mundur), jadi nama
yang sudah tercatat di `T_MIGRASI` DEV tanpa berkas diabaikan; baris menu 'masterdata' di DEV dibuang 920. Mengganti
nama `990_menu_masterdata` / `904_m_nav_menu_masterdata` di atas karena itu tidak perlu lagi; yang tetap berlaku hanya
760 / 761 / 762 → 880 / 881 / 882.

## Rute

`GET /api/master-province/meta` · `GET /api/master-province?q=&status=&halaman=` · `POST /api/master-province` · `PUT /api/master-province/{id}` ·
`PUT /api/master-province/{id}/status` · `GET /api/master-province/rujukan/{kolom}` — bentuk JSON-nya di
`modul/masterprovince/docs/issues/03-pemecahan-delapan-modul.md`. Gerbang menu: akun wajib memegang menu `masterprovince`.

## Menjalankan uji modul ini saja

Dari folder `APP_RNM/`:

```powershell
go test ./inti/backend/master/... ./modul/masterprovince/...
npx vitest run modul/masterprovince inti/frontend/master
```

## Pernyataan untuk penjaga

⛔ **Dibaca penjaga** `inti/backend/penjaga` — satu jenis pernyataan per judul `###`, satu baris per butir.

### Tabel warisan: dibaca, tidak dibuat

Tabel yang `docs/STRUKTUR-TABEL-MASTER-DATA.md` gambarkan tetapi SENGAJA tidak dibuat migrasi mana pun
(`TestKolomDDLCocokDenganStruktur`, `TestTabelBukanMilikKitaTidakDibuat`). Mencabut satu baris = kepemilikan tabel
berpindah — keputusan work owner.

| Tabel | Alasan |
| --- | --- |
| `OBJECTITEMTYPE` | tabel jenis item objek warisan POOLDATA (sumber view V_JN_OBJ_ITEM); modul `masterobjectitemtype` menulis isinya dan memakai ISACTIVE-nya; pernyataannya di sini karena STRUKTUR-nya di modul ini (warisan modul masterdata) |

Dicabut 04-10-2026 saat merge `origin/dev` (perintah work owner: "ikuti yang dari github, kalau bentrok dengan kerjaan
saya disesuaikan"): `BRANCH` - kini dinyatakan `marketingofficer` (satu tabel warisan hanya boleh dinyatakan satu
modul); modul `mastercity` hanya MEMBACA `ID`-nya untuk rujukan City. `NATION` - bukan lagi tabel warisan: sejak
`companydetail` 802 / 803 ia tabel datar yang DIBUAT migrasi (strukturnya digambarkan
`modul/companydetail/docs/STRUKTUR-TABEL-COMPANYDETAIL.md`); modul `masternation` menulis isinya tanpa mengubah
strukturnya - status dan jejak ubahnya di T_MASTER_STATUS (216 / 217, dulu 881 / 882).
