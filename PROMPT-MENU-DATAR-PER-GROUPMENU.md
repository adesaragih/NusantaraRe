# PROMPT — MENU **DATAR**: satu modul = satu menu, dikelompokkan per `GROUPMENU` *(sesi baru, folder `OUTPUT_HASIL_RNM`, cabang **`dev`** @ `19e5f36` atau lebih baru)*

> Keputusan work owner 30-09-2026: *"menu jangan ada model seperti child. Buat grouping menu antar GROUPMENU dari tabel M_NAV_MENU.
> Butir inbox, register, premiumlist, komite, tco-tahun harusnya tidak perlu, karena 1 modul 1 menu."*
> Struktur tim satu folder per modul *(`APP_RNM/inti/{backend,frontend}`, `APP_RNM/modul/<nama>/{backend,frontend,docs}`)*.
> Kerja dan commit di cabang **`dev`** *(bukan `main`)*, dengan jalur eksplisit (`git commit -o -- <jalur>`). **Nol `git push`, nol `git pull --rebase`** — push dilakukan work owner.
> `-migrate` dijalankan **work owner**. Hanya sesi ini yang menyunting `inti/` selama brief ini berjalan.

## 0. KEADAAN AWAL *(diukur asisten 30-09-2026, DEV baca-saja, agregat)*

| Hal | Isi |
| --- | --- |
| `T_MIGRASI` | `900_m_nav_menu` **sudah terpasang** di DEV — isinya **tidak boleh** disunting *(`T_MIGRASI` mencatat nama)* |
| `M_NAV_MENU` di DEV | 11 kolom: `ID`, `PARENT_ID`, `KODE`, `LABEL`, `GROUPMENU`, `MODUL`, `URUTAN`, `STATUS_AKTIF`, `DIMIGRASI`, `TGL_BUAT`, `TGL_UBAH`; constraint `PK_M_NAV_MENU`, `UQ_M_NAV_MENU_KODE`, `CK_M_NAV_MENU_GROUPMENU`, **`FK_M_NAV_MENU_INDUK`** *(PARENT_ID → ID)* |
| Baris kelompok *(PARENT_ID kosong)* | 20: TREATY 6, FACULTATIVE 3, KLAIM 8, MASTER 3; `DIMIGRASI = '1'` untuk 4 |
| Baris butir anak *(PARENT_ID terisi)* | 5: `inbox` *Inbox Claim Life* dan `register` *Register* → claimlife; `komite` *Inbox Komite* → komiteclaimlife; `premiumlist` *PremiumList* → premiumlistlife; `tco-tahun` *Treaty Contract Out* → treatycontractout |
| Kode | `inti/backend/menu/` *(pembaca, pohon GROUPMENU → kelompok → butir, `SaringMenuUntukPelaku`)*; `inti/frontend/lib/daftarMenu.ts` `susunMenu`; `components/KelompokMenu.tsx`; `PaletMenu`; `frontend/daftar.ts`; `menu.ts` tiap modul *(`MENU_<X>` berisi butir)* |

## 1. BENTUK TUJUAN

```
Beranda
TREATY
  NB Treaty In              belum dimigrasi
  EDM Treaty In             belum dimigrasi
  Treaty In                 belum dimigrasi
  Treaty In Adjustment      belum dimigrasi
  PremiumList Life
  Endorsement Life          belum dimigrasi
FACULTATIVE
  NB FacIn · RNW Fac In · Endorsment Fac In     (belum dimigrasi)
KLAIM
  Claim Fac In              belum dimigrasi
  Claim Life
  Claim Non Prop · Claim Prop                   (belum dimigrasi)
  Komite Claim FacIn        belum dimigrasi
  Komite Claim Life
  Komite Claim Non Prop · Komite Claim Prop     (belum dimigrasi)
MASTER
  Master Contract Retro Life · Master Product Name Life   (belum dimigrasi)
  Treaty Contract Out
```

- **Dua tingkat tampilan, satu tingkat data**: kepala bagian = `GROUPMENU`; di bawahnya **satu tombol per modul**. Tidak ada kelompok yang dilipat,
  tidak ada anak, tidak ada panah buka-tutup.
- **Label tombol = `LABEL` baris modul** *(nama folder korpus, VERBATIM dari isi awal 900)*.
- **Klik tombol modul membuka halaman awal modul** — halaman yang dulu dibuka butir pertamanya:

| Modul | Halaman awal | Halaman lain modul itu dibuka dari mana |
| --- | --- | --- |
| `claimlife` | `inbox` | `register` dari tombol Register di Inbox Claim Life *(sudah ada: `onRegister`)*; `outstanding`, `detail` dari baris kasus |
| `premiumlistlife` | `premiumlist` | — |
| `komiteclaimlife` | `komite` | — |
| `treatycontractout` | `tco-tahun` | halaman dalam dari layar tahun treaty, seperti sekarang |

- **Modul `DIMIGRASI = '0'`**: tombolnya tampil **nonaktif** *(tidak bisa diklik, `aria-disabled`)* dengan keterangan `belum dimigrasi` — bukan disembunyikan.
- **Modul yang dimigrasi tetapi tidak ada di `MODUL_AKTIF`**: tombolnya tidak tampil *(perilaku sekarang untuk butirnya, dipertahankan)*.
- **Judul halaman di dalam modul tidak berubah** *(mis. layar tetap berjudul "Inbox Claim Life")* — yang hilang hanya butir navigasinya.
- **Palet Ctrl+K**: Beranda + satu entri per modul yang tampil dan dimigrasi, urutan sama dengan sidebar.

⛔ **Baca ulang XML** sebelum menghapus butir `register`: cari di korpus Claim Life bagaimana Pega membuka Register *(portal/harness, tombol di worklist,
atau menu navigasi)*. Tabel di atas mengandaikan tombol Register di Inbox sudah setara dengan jalan masuk Pega. Bila XML menunjukkan jalan masuk Register
yang **tidak** tercakup tombol itu, **berhenti** di paket 3 dan laporkan buktinya *(nama berkas + baris)* — jangan membuat menu baru tanpa bukti.

## 2. MIGRASI `901_m_nav_menu_datar.sql` *(+ `_down`)* di `inti/backend/migrations/` — rentang inti 900–949

| Langkah naik | Isi |
| --- | --- |
| 1 | `DELETE FROM {skema}.M_NAV_MENU WHERE PARENT_ID IS NOT NULL` *(5 butir anak)* |
| 2 | `ALTER TABLE {skema}.M_NAV_MENU DROP CONSTRAINT FK_M_NAV_MENU_INDUK` |
| 3 | buang indeks `PARENT_ID` *(nama dari 900)* |
| 4 | `ALTER TABLE {skema}.M_NAV_MENU DROP COLUMN PARENT_ID` |

Hasil: **20 baris, satu per modul**; `KODE` = nama modul backend *(sudah begitu untuk baris kelompok)*; `URUTAN` = urutan di dalam `GROUPMENU`.
Jalur mundur mengembalikan kolom, FK, indeks, dan **5 baris butir persis** isi 900 *(label VERBATIM)*. Setiap langkah idempoten atau dilindungi
pemeriksaan katalog, supaya pelari yang gagal di tengah aman diulang. **Nol `COMMIT`** di teks SQL. Penjaga kata cadangan Oracle hijau.
Uji `db` *(SKIP tanpa `ORACLE_DSN`)* + uji skema tiruan dari nol: 900 lalu 901 menghasilkan 20 baris dan nol kolom `PARENT_ID`.

## 3. KODE

| Sisi | Ubah |
| --- | --- |
| Backend `inti/backend/menu` | pembaca tanpa `PARENT_ID`, dengan saringan **`WHERE KODE = MODUL`** — baris modul memenuhinya, 5 butir lama tidak *(`inbox` ≠ `claimlife`)*, jadi backend baru benar **sebelum dan sesudah** 901; `GET /api/menu` → `{"golongan":[{"kode":"TREATY","modul":[{"kode","label","modul","urutan","dimigrasi"}]}]}` — satu tingkat di bawah golongan; `SaringMenuUntukPelaku` tetap *(titik sambung akses per akun)*; modul dimigrasi di luar `MODUL_AKTIF` tidak dikirim; urutan golongan TREATY, FACULTATIVE, KLAIM, MASTER |
| Frontend `menu.ts` tiap modul | ganti `MENU_<X>` *(daftar butir)* dengan **`HALAMAN_AWAL_<X>`** *(satu halaman)*; label tombol datang dari tabel, bukan dari modul |
| Frontend `inti` | `susunMenu` memotong daftar tabel dengan modul yang terdaftar di `frontend/daftar.ts`: baris dimigrasi tanpa modul frontend tidak tampil *(dicatat di konsol)*, modul frontend tanpa baris tabel tidak tampil; `KelompokMenu` diganti daftar tombol datar di bawah kepala `GROUPMENU`; palet mengikuti |
| Label | butir `MENU.inbox`, `MENU.register`, `LABEL_MENU_PREMIUMLIST`, `LABEL_MENU_KOMITE`, `MENU_TCO` yang hanya dipakai navigasi dibuang; yang masih dipakai sebagai **judul halaman** tetap *(mis. Shell memakai `MENU.inbox` sebagai judul cadangan — pertahankan perilakunya)* |
| Slot menu modul 950–999 | kini **hanya** `UPDATE ... SET DIMIGRASI = '1' WHERE KODE = '<modul>'` saat modul mendapat layar pertamanya; **nol `INSERT`**. Penjaga slot diperbarui: `INSERT INTO M_NAV_MENU` di slot = merah; `UPDATE` hanya untuk baris modulnya sendiri |

## 4. PENJAGA DAN UJI

- Penjaga menu membaca **hasil bersih** 900 + 901 + slot, bukan isi 900 saja.
- Dua arah: setiap baris `DIMIGRASI = '1'` punya modul frontend dengan `HALAMAN_AWAL`; setiap modul frontend terdaftar punya baris tabel. Uji gigit
  *(mutasi)* untuk keduanya.
- 20 baris = 20 folder korpus; `CHECK GROUPMENU` tetap; `DIMIGRASI = '1'` tepat untuk modul yang punya `backend/modul.go`.
- Shell: tepat **satu tombol per modul yang dimigrasi dan aktif** *(kini 4)*; Treaty Contract Out tetap satu tombol *(tco5)*; tombol nonaktif untuk
  `DIMIGRASI = '0'`; tidak ada elemen buka-tutup kelompok di sidebar.
- Klik tombol modul → halaman awalnya *(uji per modul)*; tombol Register di Inbox Claim Life tetap membuka Register.
- Uji lama yang mengunci "lima butir" dan pohon tiga tingkat diganti, disebut satu per satu di laporan *(nama lama → nama baru)*.

## 5. DOKUMEN

`APP_RNM/PANDUAN-DEPLOY-DAN-GIT-PER-MODUL.md` bab 6 *(menambah menu = satu `UPDATE DIMIGRASI` di slot modul + `HALAMAN_AWAL` di `menu.ts`)*;
`docs/bersama/STRUKTUR-TABEL-INTI.md` *(kolom tanpa `PARENT_ID`)*; `docs/bersama/PANDUAN-TIM-PER-MODUL.md`; `modul/_templat/MODUL.md` dan
bagian slot di 20 `MODUL.md` bila menyebut `INSERT`; `PANDUAN-MENJALANKAN.txt` bab menu *(contoh jawaban `/api/menu` dan pengingat `-migrate` 901)*;
catatan bertanggal di keputusan **bg** bahwa butir navigasi dicabut oleh keputusan work owner 30-09-2026.

## 6. URUTAN — satu commit per paket, uji hijau di setiap commit

| # | Paket | Commit |
| ---: | --- | --- |
| 1 | migrasi 901 + penjaga SQL + skema tiruan | `inti: 901 - M_NAV_MENU datar, satu baris per modul, PARENT_ID dibuang` |
| 2 | backend `GET /api/menu` datar | `inti: GET /api/menu satu tingkat di bawah GROUPMENU` |
| 3 | frontend sidebar + palet + `HALAMAN_AWAL` tiap modul + penjaga dua arah | `inti: sidebar satu tombol per modul, dikelompokkan GROUPMENU` |
| 4 | slot menu hanya `UPDATE` + penjaga slot | `inti: slot menu modul hanya menyalakan DIMIGRASI` |
| 5 | dokumen §5 | `docs: menu datar per GROUPMENU` |

Setiap commit: `go build`, `go vet` dan `go test` *(dengan dan tanpa tag `db`)*, `gofmt`, `tsc`, `vitest`, `vite build`. Sebelum melapor: jalankan
`npm run dev` + backend, buka sidebar di peramban, dan pastikan tampilannya sama dengan §1 **sebelum** 901 dijalankan *(menu lama dari 900 tetap
terbaca tanpa galat — backend baru harus tahan terhadap kolom `PARENT_ID` yang masih ada)* dan sesudahnya menurut uji skema tiruan.

## 7. LAPORAN

Satu pesan: tabel **paket → commit → berkas → angka uji**; bukti XML jalan masuk Register *(berkas + baris)*; contoh jawaban `GET /api/menu`;
daftar uji lama → uji baru; pengingat `-migrate` 901 untuk work owner; bab **TELEMETRI EKSEKUSI**.

---

*Disusun 30 September 2026 dari katalog DEV `M_NAV_MENU` (kolom, constraint, cacah baris per jenis dan GROUPMENU, 5 butir anak), `T_MIGRASI`
(900 terpasang), `inti/backend/menu`, `inti/frontend/lib/daftarMenu.ts`, dan `menu.ts` empat modul.*
