# STRUKTUR TABEL — `inti` *(tabel lintas modul, migrasi 900–949)*

Tabel yang bukan milik satu modul. Migrasinya tinggal di `APP_RNM/inti/migrations/` dan dikumpulkan bersama migrasi
modul oleh `modul.SumberMigrasi` *(rentang 900–949 milik `inti`)*. Dibandingkan dengan DDL oleh
`inti/penjaga/strukturkolom_test.go`.

## M_NAV_MENU

Daftar menu aplikasi, **satu baris per modul**. Permintaan work owner 30-09-2026 *(“untuk menu buatkan dari daftar
table, nama tabelnya M_NAV_MENU … ini berguna untuk akses menu per akun nanti setelah dibuatkan akses login. Tambahkan kolom
GROUPMENU isinya TREATY, FACULTATIVE, KLAIM, MASTER”)*, `PROMPT-MENU-DARI-TABEL-M_NAV_MENU.md`. Migrasi
`900_m_nav_menu.sql` *(+ `_down`)* membuatnya dua tingkat; **`901_m_nav_menu_datar.sql`** *(+ `_down`)* meratakannya —
keputusan work owner 30-09-2026 *(“menu jangan ada model seperti child … karena 1 modul 1 menu”)*,
`PROMPT-MENU-DATAR-PER-GROUPMENU.md`: lima butir anak dibuang, lalu kunci tamu, indeks, dan kolom `PARENT_ID`.

**20 baris, satu per folder modul korpus.** Sidebar menampilkan kepala `GROUPMENU` dan satu tombol per baris; klik tombol
membuka halaman awal modul *(`HALAMAN_AWAL_<X>` di `menu.ts` modul)*. Beranda tidak di tabel ini.

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `ID` | bilangan bulat | tidak | PK | `GET /api/menu` | brief §1 — dari `SEQ_M_NAV_MENU` |
| `KODE` | teks | tidak | UNIQUE | `GET /api/menu`, frontend | brief §1 — nama modul backend *(tabel nama modul)*, sama dengan `MODUL`; pembaca menyaring `KODE = MODUL` |
| `LABEL` | teks | tidak | | sidebar, palet | brief §1 — VERBATIM nama folder korpus *(label tombol modul)* |
| `GROUPMENU` | teks | tidak | CHECK | sidebar *(kepala bagian)* | permintaan work owner 30-09-2026 — `TREATY`, `FACULTATIVE`, `KLAIM`, `MASTER` |
| `MODUL` | teks | tidak | | saringan `MODUL_AKTIF` | brief §1 — nama modul backend |
| `URUTAN` | bilangan bulat | tidak | | urutan tampil | brief §1 — urutan di dalam `GROUPMENU` |
| `STATUS_AKTIF` | teks | tidak | | saringan baca | brief §1 — `'1'` aktif *(bawaan)*, `'0'` nonaktif; konvensi data warisan |
| `DIMIGRASI` | teks | tidak | | sidebar *(“belum dimigrasi”)* | brief §1 — `'1'` modul sudah punya layar; `'0'` *(bawaan)* tombol nonaktif “belum dimigrasi” |
| `TGL_BUAT` | DATE | tidak | | jejak | brief §1 (“jejak waktu”); `DEFAULT SYSDATE NOT NULL` keputusan asisten — setiap baris punya waktu buat |
| `TGL_UBAH` | DATE | ya | | jejak | brief §1 |

**Index:** `GROUPMENU`. `KODE` berindeks lewat UNIQUE-nya. *(`IX_M_NAV_MENU_PARENT` dibuang 901.)*

**Sequence:** `SEQ_M_NAV_MENU`.

**Relasi:** tidak ada sejak 901. *(Dulu: induknya dirinya sendiri lewat `PARENT_ID`, `FK_M_NAV_MENU_INDUK` — dibuang 901.)*

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

⛔ **Dibuang 901** *(keputusan work owner 30-09-2026)* — lima butir anak isi awal 900 berikut. Jalur mundur 901
mengembalikannya persis:

| Butir *(`KODE`)* | `LABEL` | Induk |
| --- | --- | --- |
| `inbox` | Inbox Claim Life | `claimlife` |
| `register` | Register | `claimlife` |
| `premiumlist` | PremiumList | `premiumlistlife` |
| `komite` | Inbox Komite | `komiteclaimlife` |
| `tco-tahun` | Treaty Contract Out | `treatycontractout` |

Di 900 butir mewarisi `GROUPMENU`, `MODUL`, dan `DIMIGRASI` dari induknya — dibaca dari baris induk di `INSERT`. Halaman
yang dulu dibuka butir pertama kini halaman awal modulnya: `claimlife` → `inbox`, `premiumlistlife` → `premiumlist`,
`komiteclaimlife` → `komite`, `treatycontractout` → `tco-tahun`.

**Di luar lingkup** *(dicatat, tidak dibangun)*: tabel akses per akun *(mis. `M_NAV_MENU_AKSES`: akun atau peran →
`MENU_ID`)* dan login. Penyambungannya nanti di `SaringMenuUntukPelaku`.
