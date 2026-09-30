# STRUKTUR TABEL — `inti` *(tabel lintas modul, migrasi 900–949)*

Tabel yang bukan milik satu modul. Migrasinya tinggal di `APP_RNM/inti/migrations/` dan dikumpulkan bersama migrasi
modul oleh `modul.SumberMigrasi` *(rentang 900–949 milik `inti`)*. Dibandingkan dengan DDL oleh
`inti/penjaga/strukturkolom_test.go`.

## M_NAV_MENU

Daftar menu aplikasi, **dua tingkat dalam satu tabel**. Permintaan work owner 30-09-2026 *(“untuk menu buatkan dari daftar
table, nama tabelnya M_NAV_MENU … ini berguna untuk akses menu per akun nanti setelah dibuatkan akses login. Tambahkan kolom
GROUPMENU isinya TREATY, FACULTATIVE, KLAIM, MASTER”)*, `PROMPT-MENU-DARI-TABEL-M_NAV_MENU.md`. Migrasi
`900_m_nav_menu.sql` *(+ `_down`)*.

Baris ber-`PARENT_ID` kosong adalah **kelompok modul** *(satu per folder modul korpus, 20 baris)*; baris ber-`PARENT_ID`
terisi adalah **butir menu** di bawah kelompoknya *(halaman frontend yang sudah ada, 5 baris)*. Beranda tidak di tabel ini.

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `ID` | bilangan bulat | tidak | PK | `GET /api/menu` | brief §1 — dari `SEQ_M_NAV_MENU` |
| `PARENT_ID` | bilangan bulat | ya | FK | `GET /api/menu` | brief §1 — kosong = kelompok; terisi = butir. → `M_NAV_MENU.ID` |
| `KODE` | teks | tidak | UNIQUE | `GET /api/menu`, frontend | brief §1 — kelompok: nama modul backend *(tabel nama modul)*; butir: kunci halaman frontend `modul/daftar.ts` |
| `LABEL` | teks | tidak | | sidebar, palet | brief §1 — VERBATIM: kelompok = nama folder korpus *(`inti/labels.ts` `MODUL`)*; butir = label halaman yang sudah ada |
| `GROUPMENU` | teks | tidak | CHECK | sidebar *(kepala bagian)* | permintaan work owner 30-09-2026 — `TREATY`, `FACULTATIVE`, `KLAIM`, `MASTER` |
| `MODUL` | teks | tidak | | saringan `MODUL_AKTIF` | brief §1 — nama modul backend pemilik; butir mewarisi induknya |
| `URUTAN` | bilangan bulat | tidak | | urutan tampil | brief §1 — urutan di dalam induknya *(kelompok: di dalam golongannya)* |
| `STATUS_AKTIF` | teks | tidak | | saringan baca | brief §1 — `'1'` aktif *(bawaan)*, `'0'` nonaktif; konvensi data warisan |
| `DIMIGRASI` | teks | tidak | | sidebar *(“belum dimigrasi”)* | brief §1 — `'1'` modul sudah punya layar; `'0'` *(bawaan)* tampil terlipat “belum dimigrasi” |
| `TGL_BUAT` | DATE | tidak | | jejak | brief §1 — bawaan `SYSDATE` |
| `TGL_UBAH` | DATE | ya | | jejak | brief §1 |

**Index:** `PARENT_ID`, `GROUPMENU`. `KODE` berindeks lewat UNIQUE-nya.

**Sequence:** `SEQ_M_NAV_MENU`.

**Relasi:**

- induknya dirinya sendiri lewat `PARENT_ID` · 1:N · **tanpa** ON DELETE *(kelompok yang masih punya butir tidak dapat
  dihapus — nonaktifkan dengan `STATUS_AKTIF = '0'`)*

### Isi awal M_NAV_MENU

⚠️ Judul bab ini mengakhiri tabel kolom di atas bagi pembaca penjaga — tanpanya `KLAIM`, `TREATY`, … terbaca sebagai kolom.

Di migrasi yang sama, `INSERT … WHERE NOT EXISTS` atas `KODE` — idempoten:

| GROUPMENU | Kelompok *(`KODE`)* | `DIMIGRASI = '1'` |
| --- | --- | --- |
| `KLAIM` | Claim Fac In *(`claimfacin`)*, Claim Life *(`claimlife`)*, Claim Non Prop *(`claimnonprop`)*, Claim Prop *(`claimprop`)*, Komite Claim FacIn *(`komiteclaimfacin`)*, Komite Claim Life *(`komiteclaimlife`)*, Komite Claim Non Prop *(`komiteclaimnonprop`)*, Komite Claim Prop *(`komiteclaimprop`)* | Claim Life, Komite Claim Life |
| `FACULTATIVE` | NB FacIn *(`nbfacin`)*, RNW Fac In *(`rnwfacin`)*, Endorsment Fac In *(`endorsmentfacin`)* | — |
| `TREATY` | NB Treaty In *(`nbtreatyin`)*, EDM Treaty In *(`edmtreatyin`)*, Treaty In *(`treatyin`)*, Treaty In Adjustment *(`treatyinadjustment`)*, PremiumList Life *(`premiumlistlife`)*, Endorsement Life *(`endorsementlife`)* | PremiumList Life |
| `MASTER` | Master Contract Retro Life *(`mastercontractretrolife`)*, Master Product Name Life *(`masterproductnamelife`)*, Treaty Contract Out *(`treatycontractout`)* | Treaty Contract Out |

`GROUPMENU` kelompok `[DIPUTUSKAN asisten; veto work owner]` — tiga yang paling mungkin diubah: PremiumList Life dan
Endorsement Life *(bisnis Life, di TREATY)*, Treaty Contract Out *(di MASTER)*.

| Butir *(`KODE`)* | `LABEL` | Induk |
| --- | --- | --- |
| `inbox` | Inbox Claim Life | `claimlife` |
| `register` | Register | `claimlife` |
| `premiumlist` | PremiumList | `premiumlistlife` |
| `komite` | Inbox Komite | `komiteclaimlife` |
| `tco-tahun` | Treaty Contract Out | `treatycontractout` |

Butir mewarisi `GROUPMENU`, `MODUL`, dan `DIMIGRASI` dari induknya — dibaca dari baris induk di `INSERT`, tidak ditulis ulang.

**Di luar lingkup** *(dicatat, tidak dibangun)*: tabel akses per akun *(mis. `M_NAV_MENU_AKSES`: akun atau peran →
`MENU_ID`)* dan login. Penyambungannya nanti di `SaringMenuUntukPelaku`.
