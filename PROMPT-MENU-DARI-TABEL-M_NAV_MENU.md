# PROMPT — MENU DARI TABEL **`M_NAV_MENU`** *(sesi baru, folder `OUTPUT_HASIL_RNM`, cabang `main` @ `29f766f` atau lebih baru)*

> Permintaan work owner 30-09-2026: *"untuk menu buatkan dari daftar table, nama tabelnya M_NAV_MENU, masukkan list menunya ke dalam situ, ini
> berguna untuk akses menu per akun nanti setelah dibuatkan akses login. Tambahkan kolom GROUPMENU isinya TREATY, FACULTATIVE, KLAIM, MASTER."*
> Struktur bentuk B *(`APP_RNM/inti`, `APP_RNM/modul/<nama>`)*. Commit dengan jalur eksplisit. `-migrate` dijalankan **work owner**.

## 0. KEADAAN AWAL *(asisten, 30-09-2026)*

- Nama `M_NAV_MENU` **belum dipakai** di DEV *(nol objek `%NAV%MENU%`)*; `T_MIGRASI` 32 langkah, terakhir `058`.
- Menu sekarang dirakit di frontend: `frontend/src/modul/daftar.ts` dari `menu.ts` tiap modul — butir: `Beranda`, `inbox` *(Inbox Claim Life)*, `register`
  *(Register)*, `premiumlist` *(PremiumList)*, `komite` *(Inbox Komite)*, `tco-tahun` *(Treaty Contract Out)*. Kelompok modul di `inti/labels.ts`
  `MODUL` **18** nama; korpus punya **20** folder modul — **`Treaty In`** dan **`Treaty In Adjustment`** belum ada di daftar itu.
- Saklar `MODUL_AKTIF` dan `GET /api/modul-aktif` sudah ada.

## 1. TABEL `M_NAV_MENU` — satu tabel, dua tingkat

Migrasi bersama **`900_m_nav_menu.sql`** *(+ `_down`)* di `inti/` — rentang baru **900–949 milik `inti`** *(tabel lintas modul)*; penjalan migrasi
mengumpulkannya bersama migrasi modul.

| Kolom | Tipe | Isi |
| --- | --- | --- |
| `ID` | `NUMBER(10)` PK | pengenal |
| `PARENT_ID` | `NUMBER(10)` NULL, FK → `M_NAV_MENU(ID)` | `NULL` = baris **kelompok modul**; terisi = **butir menu** di bawah kelompoknya |
| `KODE` | `VARCHAR2(64)` NOT NULL UNIQUE | kelompok: nama modul backend *(tabel nama: `claimlife`, `premiumlistlife`, …, `nbfacin`)*; butir: kunci halaman frontend *(`inbox`, `register`, `premiumlist`, `komite`, `tco-tahun`)* |
| `LABEL` | `VARCHAR2(100)` NOT NULL | teks tampil **VERBATIM** dari label yang sudah ada |
| `GROUPMENU` | `VARCHAR2(16)` NOT NULL, `CHECK IN ('TREATY','FACULTATIVE','KLAIM','MASTER')` | golongan |
| `MODUL` | `VARCHAR2(32)` NOT NULL | nama modul backend pemilik *(nilai `MODUL_AKTIF`)* |
| `URUTAN` | `NUMBER(5)` NOT NULL | urutan di dalam induknya |
| `STATUS_AKTIF` | `VARCHAR2(1)` NOT NULL DEFAULT `'1'` | `'1'` aktif, `'0'` nonaktif *(konvensi yang sama dengan data warisan)* |
| `DIMIGRASI` | `VARCHAR2(1)` NOT NULL DEFAULT `'0'` | `'1'` modul sudah punya layar; `'0'` tampil sebagai "belum dimigrasi" |
| `TGL_BUAT`, `TGL_UBAH` | `DATE` | jejak waktu |

Indeks pada `PARENT_ID` dan `GROUPMENU`. Sequence `SEQ_M_NAV_MENU`. Penjaga kata cadangan Oracle hijau. **Nol `COMMIT`** di teks SQL.

## 2. ISI AWAL — di migrasi yang sama *(`INSERT`)*

**Kelompok modul, 20 baris** — `GROUPMENU` usulan asisten `[DIPUTUSKAN; veto work owner]`:

| GROUPMENU | Kelompok *(LABEL = nama folder korpus)* |
| --- | --- |
| **KLAIM** | Claim Fac In, Claim Life, Claim Non Prop, Claim Prop, Komite Claim FacIn, Komite Claim Life, Komite Claim Non Prop, Komite Claim Prop |
| **FACULTATIVE** | NB FacIn, RNW Fac In, Endorsment Fac In |
| **TREATY** | NB Treaty In, EDM Treaty In, Treaty In, Treaty In Adjustment, PremiumList Life, Endorsement Life |
| **MASTER** | Master Contract Retro Life, Master Product Name Life, Treaty Contract Out |

⚠️ Tiga yang paling mungkin diubah work owner: **PremiumList Life** dan **Endorsement Life** *(bisnis Life, dimasukkan TREATY)*, dan **Treaty Contract Out**
*(master arrangement treaty keluar, dimasukkan MASTER)*. `DIMIGRASI = '1'` untuk Claim Life, PremiumList Life, Komite Claim Life, Treaty Contract Out; sisanya `'0'`.

**Butir menu, 5 baris** *(anak kelompoknya; `GROUPMENU` sama dengan induk)*: `inbox` *Inbox Claim Life* dan `register` *Register* → Claim Life;
`premiumlist` *PremiumList* → PremiumList Life; `komite` *Inbox Komite* → Komite Claim Life; `tco-tahun` *Treaty Contract Out* → Treaty Contract Out.
**Beranda tidak** masuk tabel *(selalu tampil, bukan milik modul)*.

## 3. BACKEND DAN FRONTEND

| Sisi | Isi |
| --- | --- |
| Repository `inti` | pembaca `M_NAV_MENU` *(baris `STATUS_AKTIF = '1'`, urut `GROUPMENU`, `URUTAN`)* |
| Rute | `GET /api/menu` → pohon `GROUPMENU → kelompok → butir`; butir milik modul yang **tidak** ada di `MODUL_AKTIF` tidak dikirim; **titik sambung** untuk saringan per akun nanti *(satu fungsi `SaringMenuUntukPelaku(pelaku, menu)` yang sekarang meneruskan semua — diberi komentar "akses per akun menyusul")* |
| Frontend | sidebar dan palet Ctrl+K dirakit dari `GET /api/menu`, **dipotong** dengan rute yang benar-benar terdaftar di `modul/daftar.ts` *(baris tabel tanpa rute frontend tidak tampil dan dicatat di konsol; butir frontend tanpa baris tabel juga tidak tampil)*; judul golongan `TREATY`, `FACULTATIVE`, `KLAIM`, `MASTER` tampil sebagai kepala bagian di sidebar; kelompok `DIMIGRASI = '0'` tetap tampil terlipat "belum dimigrasi"; bila `GET /api/menu` gagal, sidebar menampilkan galatnya *(bukan menu kosong diam-diam)* |
| Penjaga | uji statik **dua arah**: setiap `KODE` butir di isi awal migrasi 900 ada di `daftar.ts`, dan setiap butir `daftar.ts` ada di isi awal; uji kelompok = 20 dan sama dengan folder modul korpus; uji `CHECK GROUPMENU` |
| Label | `inti/labels.ts` `MODUL` ditambah **Treaty In** dan **Treaty In Adjustment** |

## 4. DI LUAR LINGKUP — dicatat, tidak dibangun

Tabel akses per akun *(mis. `M_NAV_MENU_AKSES`: akun atau peran → `MENU_ID`)* dan login. `SaringMenuUntukPelaku` adalah tempat penyambungannya nanti.

## 5. URUTAN — satu commit per paket

| # | Paket | Commit |
| ---: | --- | --- |
| 1 | migrasi 900 + isi awal + penjaga SQL | `inti: M_NAV_MENU — tabel menu dua tingkat dengan GROUPMENU, isi awal 20 kelompok + 5 butir` |
| 2 | repository + `GET /api/menu` + `SaringMenuUntukPelaku` | `inti: GET /api/menu dari M_NAV_MENU` |
| 3 | frontend sidebar + palet dari API + penjaga dua arah + label dua modul | `inti: sidebar dan palet dari M_NAV_MENU` |
| 4 | dokumen: `PANDUAN-DEPLOY-DAN-GIT-PER-MODUL.md` *(cara menambah menu: satu baris di tabel + satu butir di `menu.ts`)*, PANDUAN-MENJALANKAN | `docs: menu dari M_NAV_MENU` |

Uji tiap commit dengan dan tanpa tag `db`; vitest; tsc; build. Laporan satu pesan: tabel **paket → commit → berkas → angka uji**, isi awal tabel, dan
pengingat `-migrate` untuk work owner.

---

*Disusun 30 September 2026 dari katalog DEV (nama `M_NAV_MENU` bebas, `T_MIGRASI` 32), `modul/daftar.ts`, `menu.ts` empat modul, dan `inti/labels.ts`.*
