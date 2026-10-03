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

**Akses per akun** *(Kelola User, keputusan work owner 01-10-2026)*: `M_LOGIN_GO_MENU` di bawah — KODE menu per akun.
`GET /api/menu` mengirim hanya menu akun yang login (`SaringMenuUntukAkun`), dan rute modul yang menunya tidak dimiliki
dijawab 403 oleh `cmd/api`. Menu **Kelola User** (`kelolauser`, golongan `ADMIN`) **bukan** baris tabel ini — tabel ini
tetap dua puluh baris, satu per folder modul korpus; ia hidup di kode (`menu.MenuAplikasi`) seperti Beranda.

## M_LOGIN_GO

Akun login, **satu baris per orang**. Keputusan work owner 01-10-2026 *(“nama table nya M_LOGIN_GO … kolom untuk
menampung data dari M_UNIT, M_DIVISION, M_ORGANIZATION”; “simpan aja code nya, jangan ID”)*. Migrasi
`902_m_login_go.sql` *(+ `_down`)*. Sandi **tidak pernah** disimpan — hanya hash bcrypt.

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `LOGIN_ID` | teks | tidak | PK | layar login, `CREATE_OP` | keputusan work owner 01-10-2026 — akun yang diketik |
| `NAME` | teks | tidak | | `CREATE_OP_NAME`, Shell | nama tampilan |
| `PASSWORD_HASH` | teks | tidak | | login | hash bcrypt (cost 12); sandi asli tidak disimpan |
| `ORGANIZATION_CODE` | teks | ya | | profil | `M_ORGANIZATION.CODE` — tanpa FK, diperiksa aplikasi saat menyimpan |
| `DIVISION_CODE` | teks | ya | | profil | `M_DIVISION.CODE` — harus milik organisasinya |
| `UNIT_CODE` | teks | ya | | profil | `M_UNIT.CODE` — harus milik divisinya |
| `IS_ACTIVE` | teks | tidak | CHECK | login | `'1'` aktif *(bawaan)*, `'0'` nonaktif |
| `FAILED_COUNT` | bilangan bulat | tidak | | penguncian | sandi salah beruntun; nol sesudah login berhasil |
| `LOCKED_UNTIL` | DATE | ya | | penguncian | 15 menit sesudah salah ke-5 |
| `MUST_CHANGE_PASSWORD` | teks | tidak | CHECK | login | `'1'` *(bawaan)* — wajib ganti sandi saat login berikutnya; diatur centang “Change Password Next Login” di Kelola User (tab Security) |
| `SESSION_VERSION` | bilangan bulat | tidak | | sesi | naik saat logout, ganti sandi, nonaktif — mencabut semua cookie lama |
| `LAST_LOGIN` | DATE | ya | | jejak | login berhasil terakhir |
| `TGL_CREATE` | DATE | tidak | | jejak | `DEFAULT SYSDATE` |
| `TGL_UPDATE` | DATE | ya | | jejak | |
| `EMAIL` | teks | ya | | Kelola User | migrasi `904_m_login_go_kontak.sql` (03-10-2026, permintaan work owner *"tambahkan email, no hp, nik dan jabatan; buat dalam bahasa inggris"*) — alamat email, maks. 254; label layar `Email` |
| `PHONE_NUMBER` | teks | ya | | Kelola User | migrasi 904 — nomor HP 8–15 digit, boleh diawali `+`, boleh spasi/tanda hubung; label `Phone Number` |
| `EMPLOYEE_ID` | teks | ya | | Kelola User | migrasi 904 — NIK (Nomor Induk Karyawan, bukan NIK kependudukan), huruf/angka/titik/garis miring/tanda hubung, maks. 30; label `Employee ID (NIK)` |
| `JOB_POSITION` | teks | ya | | Kelola User | migrasi 904 — jabatan, maks. 150; label `Position` |

**Index:** PK.

**Relasi:** tidak ada FK ke `M_ORGANIZATION`, `M_DIVISION`, `M_UNIT` *(keputusan work owner — ketiga master tidak
dibuat migrasi aplikasi)*. Anaknya `M_LOGIN_GO_WORKBASKET` dan `M_LOGIN_GO_MENU` — **Hapus permanen** di Kelola User
membuang keduanya lebih dulu, lalu akunnya, dalam satu transaksi *(FK tanpa `ON DELETE`)*.

## M_LOGIN_GO_WORKBASKET

Workbasket setiap akun — **satu orang bisa beberapa workbasket**. Workbasket adalah **peran**: `M_WORKBASKET.WORKBASKET_ID`
berisi nama peran yang dipakai aplikasi *(`ReasLifeAdmin`, `ReasLifeSPV`, …)*, jadi tabel ini sekaligus tabel peran.

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `LOGIN_ID` | teks | tidak | PK, FK | `Pelaku.Peran` | → `M_LOGIN_GO.LOGIN_ID`, tanpa `ON DELETE` |
| `WORKBASKET_ID` | teks | tidak | PK | `Pelaku.Peran` | `M_WORKBASKET.WORKBASKET_ID` — tanpa FK; hanya yang `IS_ACTIVE = 1` di master berlaku |
| `TGL_CREATE` | DATE | tidak | | jejak | `DEFAULT SYSDATE` |

**Index:** PK *(`LOGIN_ID`, `WORKBASKET_ID`)*; `IX_M_LOGIN_GO_WB_WORKBASKET` *(`WORKBASKET_ID`)*.

## M_LOGIN_GO_MENU

Menu yang boleh dibuka setiap akun — **akses per akun** *(Kelola User, keputusan work owner 01-10-2026: “tambahkan menu
apa aja yang bisa diakses sama user nya”)*. Migrasi `903_m_login_go_menu.sql` *(+ `_down`)*. Menu yang tidak dimiliki
**disembunyikan** di sidebar **dan ditolak** server *(403)*; perubahan berlaku pada permintaan berikutnya.

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `LOGIN_ID` | teks | tidak | PK, FK | sidebar, gerbang 403 | → `M_LOGIN_GO.LOGIN_ID`, tanpa `ON DELETE` |
| `MENU_KODE` | teks | tidak | PK | sidebar, gerbang 403 | `M_NAV_MENU.KODE` (nama modul) atau menu aplikasi `kelolauser` — tanpa FK; diperiksa aplikasi saat menyimpan |
| `TGL_CREATE` | DATE | tidak | | jejak | `DEFAULT SYSDATE` |

**Index:** PK *(`LOGIN_ID`, `MENU_KODE`)*.

**Isi awal 903:** setiap akun yang sudah ada mendapat **semua** menu — dua puluh modul dan `kelolauser` *(keputusan work
owner)*; daftarnya dijaga `TestIsiAwalMenuAkunMemuatSemuaMenu`. Pemegang `kelolauser` adalah admin. Kelola User menolak
admin membuang `kelolauser` dari dirinya sendiri, menonaktifkan atau menghapus dirinya, dan perubahan yang menyisakan
nol akun aktif ber-`kelolauser`.
