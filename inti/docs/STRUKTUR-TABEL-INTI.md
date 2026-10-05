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

**20 baris, satu per folder modul korpus**, ditambah baris modul DI LUAR korpus dari langkah inti sesudah 901
(PANDUAN-TIM-PER-MODUL bab 5): **`912`–`919_m_nav_menu_master<nama>.sql`** *(+ `_down`)* — delapan modul master
(`masternation` / `Nation`, `masterprovince` / `Province`, `mastercity` / `City`, `masterdistrict` / `District`,
`masterczone` / `CZone`, `masteraccumulatedtype` / `Accumulated Type`, `masteraccumulation` / `Accumulation`,
`masterobjectitemtype` / `Object Item Type`), `MASTER` urutan 6–13 (sesudah 909), `DIMIGRASI '1'` — modul luar korpus TANPA
slot menu (keputusan work owner 04-10-2026, diteruskan sesi 9d: "8 modul terpisah", "Modul luar korpus tanpa slot"; INSERT
bentuk datar idempoten, jalur mundur `DELETE` baris itu dan akses `M_LOGIN_GO_MENU`-nya). **`920_m_nav_menu_masterdata_pensiun.sql`**
— modul `masterdata` dihapus: hak menunya disalin ke delapan menu ("Ya, salin otomatis"), lalu hak dan barisnya dibuang
(`911_m_nav_menu_masterdata`, dulu 904, ikut dibuang; di DEV yang sudah menjalankannya, 920 yang membuang barisnya).
Akses akun lain TIDAK masuk isi awal 903 — diberikan lewat Kelola User.
Penjaga `menu_test.go` menerapkannya di skema tiruan; `rentang_test.go` hanya menerima INSERT bentuk itu + DELETE mundurnya.
Sidebar menampilkan kepala `GROUPMENU` dan satu tombol per baris; klik tombol
membuka halaman awal modul *(`HALAMAN_AWAL_<X>` di `menu.ts` modul)*. Beranda tidak di tabel ini.

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `ID` | bilangan bulat | tidak | PK | `GET /api/menu` | brief §1 — dari `SEQ_M_NAV_MENU` |
| `KODE` | teks | tidak | UNIQUE | `GET /api/menu`, frontend | brief §1 — nama modul backend *(tabel nama modul)*, sama dengan `MODUL`; pembaca menyaring `KODE = MODUL` |
| `LABEL` | teks | tidak | | sidebar, palet | brief §1 — VERBATIM nama folder korpus *(label tombol modul)* |
| `GROUPMENU` | teks | tidak | CHECK | sidebar *(kepala bagian)* | permintaan work owner 30-09-2026 — `TREATY`, `FACULTATIVE`, `KLAIM`, `MASTER`; `MASTER TREATY` ditambah migrasi `909` (work owner 04-10-2026: Treaty In, Treaty In Adjustment, Treaty Contract Out) |
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
| `EMAIL` | teks | ya | | Kelola User | migrasi `904_m_login_go_kontak.sql` (03-10-2026, permintaan work owner *"tambahkan email, no hp, nik dan jabatan; buat dalam bahasa inggris"*) — alamat email, maks. 254; label layar `Email`. Sejak migrasi 905: disimpan huruf kecil dan unik tanpa beda huruf *("proteksi email sudah terdaftar")* — **bukan** jalan login |
| `PHONE_NUMBER` | teks | ya | | Kelola User | migrasi 904 — nomor HP 8–15 digit, boleh diawali `+`, boleh spasi/tanda hubung; label `Phone Number` |
| `EMPLOYEE_ID` | teks | ya | | Kelola User | migrasi 904 — NIK (Nomor Induk Karyawan, bukan NIK kependudukan), huruf/angka/titik/garis miring/tanda hubung, maks. 30; label `Employee ID (NIK)` |
| `JOB_POSITION` | teks | ya | | Kelola User | migrasi 904 — jabatan, maks. 150; label `Position` |
| `CONTACT_ID` | teks | tidak | UNIQUE | Kelola User; rancangan Kelola Marketing Officer (`MARKETINGOFFICER.CLIENTID`) | migrasi `905_m_login_go_contact_id.sql` (03-10-2026, keputusan work owner V1 *"M_LOGIN_GO ID nya pake CON-xxx"*) — `CON-n`, n dari `M_LOGIN_GO_CONTACT_SEQ` mulai 1001 *(di atas nomor kontak SFAGIS terbesar)*; sejak migrasi `910_m_login_go_contact_seq_max.sql` (04-10-2026, *"ambil dari max con+1"*) sequence-nya dibuat ulang dengan nilai awal = CON tertinggi + 1 (minimal 1001), dihitung saat migrasi dijalankan; diisi saat akun dibuat, **tidak pernah berubah**; label `Contact ID` |

**Index:** PK; `UX_M_LOGIN_GO_CONTACT_ID` *(`CONTACT_ID`, unik)*; `UX_M_LOGIN_GO_LOGIN_ID` *(`LOWER(LOGIN_ID)`, unik —
username tidak boleh hanya beda huruf)*; `UX_M_LOGIN_GO_EMAIL` *(`LOWER(EMAIL)`, unik)* — ketiganya migrasi 905.

**Login** tetap dengan username (`LOGIN_ID`) saja — login lewat email **dibatalkan** work owner 03-10-2026
*("LOGIN LEWAT EMAIL TIDAK JADI!!")*. Username dan email masing-masing unik tanpa beda huruf, diperiksa aplikasi
saat menyimpan (pesan 409 `Username is already registered` / `Email is already registered to another account`);
index unik di atas adalah pengaman terakhir.

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
| `MENU_KODE` | teks | tidak | PK | sidebar, gerbang 403 | `M_NAV_MENU.KODE` (nama modul) atau menu aplikasi `kelolauser`, `templatemanager` *(sejak 04-10-2026)* — tanpa FK; diperiksa aplikasi saat menyimpan |
| `TGL_CREATE` | DATE | tidak | | jejak | `DEFAULT SYSDATE` |
| `HAK` | teks | tidak | | gerbang tulis `cmd/api`, tombol layar | `PENUH` atau `LIHAT` *(migrasi 914, keputusan work owner 04-10-2026)*; `DEFAULT 'PENUH'`, `CK_M_LOGIN_GO_MENU_HAK` |

**Index:** PK *(`LOGIN_ID`, `MENU_KODE`)*.

**Hak per menu (914):** `LIHAT` = hanya membaca. Gerbang `cmd/api` menolak `POST`/`PUT`/`DELETE` modul itu kecuali rute
yang modulnya bebaskan, dan layar modul menyembunyikan tombol tulis. `LIHAT` hanya untuk menu modul yang mendaftar
(`inti.Pendaftaran.HakLihat`); menu lain dan menu aplikasi selalu `PENUH`. `LIHAT` berlaku juga bagi pemegang `kelolauser`
*(keputusan work owner 05-10-2026)* — Kelola User sendiri tidak pernah `LIHAT`, jadi ia selalu dapat mengembalikan aksesnya.

**Isi awal 903:** setiap akun yang sudah ada mendapat **semua** menu — dua puluh modul dan `kelolauser` *(keputusan work
owner)*; daftarnya dijaga `TestIsiAwalMenuAkunMemuatSemuaMenu`. Pemegang `kelolauser` adalah admin. Kelola User menolak
admin membuang `kelolauser` dari dirinya sendiri, menonaktifkan atau menghapus dirinya, dan perubahan yang menyisakan
nol akun aktif ber-`kelolauser`.
`templatemanager` lahir sesudah 903, jadi akun lama tidak otomatis memegangnya — haknya diberikan lewat Kelola User.

## M_TEMPLATE_FILE

Berkas templat unduhan berversi milik menu aplikasi **Template Manager** *(keputusan work owner 04-10-2026: “menu untuk
pengelola file templete unduhan … bukan cuman untuk bordereaux, tapi untuk semua menu yang memiliki template”)*. Migrasi
`912_m_template_file.sql` *(+ `_down`)*. Padanan Pega: Rule-File-Binary yang dicari menurut nama lalu diambil yang
terbaru. Slot templat didaftarkan modul (`inti.Pendaftaran.Templat`); slot tanpa baris aktif memakai **berkas bawaan**
modulnya, jadi tabel kosong tidak mengubah perilaku apa pun.

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `ID` | bilangan bulat | tidak | PK | — | `SEQ_M_TEMPLATE_FILE` |
| `KODE` | teks | tidak | UNIQUE | Template Manager, unduh templat | kode slot modul, mis. `bordereaux.premi.fire` — tanpa FK; diperiksa katalog slot |
| `VERSI` | bilangan bulat | tidak | UNIQUE | Riwayat versi | 1, 2, … per `KODE` *(MAX + 1 dalam transaksi unggah)* |
| `NAMA_BERKAS` | teks | tidak | | Riwayat versi | nama berkas yang diunggah *(tanpa jalur)*; nama unduhan pengguna tetap nama slot |
| `UKURAN` | bilangan bulat | tidak | | Riwayat versi | byte |
| `JUMLAH_KOLOM` | bilangan bulat | ya | | Riwayat versi | kolom header CSV; kosong untuk slot bukan CSV |
| `ISI` | berkas | tidak | | unduh templat | BLOB — isi berkas utuh, maks. 1 MB; satu-satunya kolom berkas *(pengecualian bernama di `TestKolomUangDesimalDanNolJSON`)* |
| `CATATAN` | teks | ya | | Riwayat versi | catatan pengunggah, maks. 500 huruf |
| `AKTIF` | teks | tidak | CHECK | unduh templat | `'1'` versi yang dipakai — paling banyak satu per `KODE`; `'0'` *(bawaan)* |
| `DIUNGGAH_OLEH` | teks | tidak | | Riwayat versi | `M_LOGIN_GO.LOGIN_ID` pengunggah |
| `TGL_UNGGAH` | DATE | tidak | | Riwayat versi | `DEFAULT SYSDATE` |

**Index:** PK *(`ID`)*; `UX_M_TEMPLATE_FILE_VERSI` *(`KODE`, `VERSI`)*.
