# Modul `treatyin` — Treaty In

⚠️ **Terdaftar 1 Oktober 2026 — LAPISAN SKEMA tujuh belas tiket, bukan tiketnya.**
`backend/modul.go` ada, `inti/backend/daftar/modul_treatyin_gen.go` bangkit, dan slot menu `972`
menyalakan `M_NAV_MENU.DIMIGRASI`.

Yang berdiri (migrasi `400`–`419`): **29 tabel, 230 kolom, 29 sequence** — tiket `14` `15` `20`
`22`–`31` `33` `34` `37` `39`. Ditambah satu jalur baca atas tabel acuan. Cacah kolom tiap tabel **sama persis**
dengan `KAMUS-KOLOM.md`, kecuali `VERSI_KONTRAK` yang bertambah satu kolom dari tiket `01` papan
Adjustment.

⛔ **KETUJUH BELAS TIKET ITU BELUM DINYATAKAN SELESAI.** Bukan karena buktinya kurang — sembilan
constraint sudah terbukti menolak dan dua jalur positif terbukti diterima — melainkan karena status
adalah medan di dalam berkas tiket, dan yang berwenang mengubahnya pemilik proses. Yang masih
menghalangi hanya **satu kalimat dari DBA**: skema uji tersendiri, supaya `make test-db` dapat
menjalankan bukti itu sebagai uji otomatis, bukan sebagai sesi manual yang tercatat di dokumen.

⛔ **Tiga invarian tertahan `F-13`** — `INV-32`, `INV-33`, dan `INV-52`. Ketiganya bergolongan
*indeks unik*, yang di Oracle berarti materialized view ber-`REFRESH ON COMMIT`, dan
`F-13-MV-TANPA-PEMANTAU-KEBASIAN.md` menyatakan penegakan lewat MV di modul ini **belum punya
pemantau kebasian**: *"MV yang gagal me-refresh berhenti menegakkan tanpa satu galat pun."*
Memasangnya tanpa Oracle yang dapat diuji (`L-3`) menanam constraint yang terlihat terpasang dan
tidak menolak apa pun. **Urutannya: pemantau kebasian lebih dulu.**

⛔ **Tiga belas invarian tabelnya berdiri tetapi BELUM ditegakkan di mana pun** (`INV-29` dan
`INV-53` sudah keluar dari daftar ini — lihat lapisan aplikasi tiket `14` di atas) — `INV-29`, `INV-38`,
`INV-39`, `INV-40`, `INV-53`, `INV-55`, `INV-56`, `INV-57`, dan `INV-49` yang bahkan belum punya
tempat berdiri. Sebagian besar membandingkan baris di tabel BERBEDA, yang `CHECK` Oracle tidak dapat
nyatakan; `INV-29` dan `INV-38` **dapat**, dan yang menahannya hanya ADR-0056. Seluruhnya menunggu
jalur simpan yang sama. Daftar lengkapnya di
[`docs/STRUKTUR-TABEL-TREATY-IN.md`](docs/STRUKTUR-TABEL-TREATY-IN.md).

✅ **Lapisan aplikasi tiket `14` mendarat 2 Oktober 2026.** `POST /api/treaty-in/kontrak` membuat
kontrak beserta kepala versi pertamanya dalam SATU transaksi, `GET /api/treaty-in/kontrak/{id}`
menemukannya kembali dengan pengenal yang sistem berikan. **`INV-53` dan `INV-29` kini DITEGAKKAN**
di lapisan services — keduanya yang migrasi `401` tulis "DITAGIH: tiket lapisan aplikasi".
~~Nol layar baru (`L-4`).~~ **`L-4` dicabut 3 Oktober 2026** — dua layar dibangun dari ekspor Pega; lihat bab di bawah.

✅ **Migrasi `400`–`419` dan `440`–`441` SUDAH BERJALAN DI ORACLE** (`DEV_NUSARE2`, 2 Oktober 2026):
ke-22 langkah tercatat di `T_MIGRASI`, **nol `ORA-`**. Lima constraint dibuktikan sungguh menolak
lewat `INSERT` + `ROLLBACK`, dan satu klaim saya sebelumnya — bahwa `UQ_LAYER` tidak menegakkan
`INV-05` — **terbukti KELIRU dan dicabut**. Rinciannya di
[`docs/VERIFIKASI-ORACLE-2026-10-02.md`](docs/VERIFIKASI-ORACLE-2026-10-02.md).

⛔ **Sequence berdiri tetapi BELUM TERPAKAI.** Tiket `14` menuntut *"tidak ada jalur lain yang dapat
memberi pengenal"*; hari ini belum ada jalur tulis sama sekali, sehingga tuntutan itu benar secara
hampa. Ia ditegakkan bersama jalur simpan.

ℹ️ **Layar modul ini ada karena mendaftarkan modul menuntutnya** `[keputusan work owner 01-10-2026]`.
Penjaga `daftar.modulAktif.test.ts` mewajibkan setiap modul terdaftar punya `frontend/`, sementara
papan tiket melarang mengarang layar (~~`L-4`~~, **dicabut 3 Okt 2026**). Jalan keluar saat itu: satu halaman **baca-saja** atas data
yang tiket `15` buat — nol alur karangan.

⚠️ **Lima penyelarasan spec dengan repo** tercatat di
[`docs/KEPUTUSAN-PENYELARASAN-REPO.md`](docs/KEPUTUSAN-PENYELARASAN-REPO.md). **Jangan menyunting DDL
tanpa membacanya.**

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu. Ubah
nilainya hanya lewat pull request yang disetujui tim inti — dua modul tidak boleh berbagi nomor.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `treatyin` |
| Folder korpus | `Treaty In` |
| GROUPMENU | `MASTER TREATY` |
| Pemilik | `@PEMILIK-TREATYIN` |
| Status | dimigrasi |
| Rentang migrasi | `400-439` |
| Slot menu | `972-973` |
| Prefix rute API | `/api/treaty-in` |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

`Pemilik` adalah penanda pemegang modul. Wilayah berkas yang boleh disentuh cabang
`module/<nama>` dijaga `.github/workflows/penjaga-wilayah-cabang.yml` - CODEOWNERS
dipensiunkan 1 Oktober 2026.

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | `STRUKTUR-TABEL-TREATY-IN.md`, `KEPUTUSAN-PENYELARASAN-REPO.md`, `LEDGER-PROGRES.md`, `PEMETAAN-M-TREATY-IN2.md`, tiga `PERTANYAAN-TERBUKA-*.md`, `issues/` |
| `backend/` | `modul.go` (`Pendaftaran()`), `models/`, `repository/`, `services/`, `handlers/`, `migrations/`, `pemuat/` |
| `alat/` | skrip SQL yang dijalankan dengan tangan — **bukan migrasi**, tidak tercatat di `T_MIGRASI` |
| `frontend/` | `menu.ts`, `rute.tsx`, `api.ts`, `labels.ts`, `pages/AcuanTreatyIn.tsx` |

## Migrasi

Rentang `400-439` (tabel R2, urut hulu ke hilir: migrasi modul hilir yang merujuk tabel modul hulu
selalu berjalan sesudahnya). Slot menu `972-973` hanya menyalakan `DIMIGRASI` baris modul ini (satu `UPDATE`,
nol `INSERT` — menu datar 30-09-2026) saat modul mendapat layar pertamanya, di folder
`backend/migrations/` modul ini sendiri — bentuk SQL-nya di `APP_RNM/PANDUAN-DEPLOY-DAN-GIT-PER-MODUL.md`
bab 6. Nomor selalu tiga digit.

⚠️ **Dua pencabutan hidup di rentang ini**, dan keduanya mengikuti pola yang sama: berkas
pencabutan BARU, bukan menghapus migrasi lama dari riwayat. `434_cabut_mata_uang.sql`
([`KEPUTUSAN` §16](docs/KEPUTUSAN-PENYELARASAN-REPO.md)) dan `435_cabut_nilai_selisih.sql`
([§19](docs/KEPUTUSAN-PENYELARASAN-REPO.md)) — yang kedua memindahkan tabelnya ke modul
Adjustment, bukan membuangnya.

⚠️ **`backend/pemuat/jalankan.go` bertanda `//go:build ignore`, dan itu disengaja.** Ia `package
main` yang memuat delapan tabel pendaratan ke POOLDATA, dijalankan dengan
`go run modul/treatyin/backend/pemuat/jalankan.go`. Ia **tidak** dapat tinggal di `cmd/`:
`TestModulTidakMengimporModulLain` menyatakan *"cmd hanya mengimpor `inti/...`"*. Ia juga tidak
dapat tinggal di `alat/`: `pelanggaranLetak` di penjaga yang sama menuntut kode Go modul berada di
`modul/<nama>/backend/`. Bawaannya **kering** — tanpa `-ikat`, setiap transaksi dibatalkan.

## Pernyataan untuk penjaga

### Kaskade ON DELETE CASCADE

Kaskade HANYA pada berkas migrasi modul ini yang berawalan di bawah; berkas lain modul ini tanpa
`ON DELETE CASCADE` (`TestKaskadeHanyaPadaRelasiTerdaftar`). Mendaftarkan yang baru menuntut bukti.

⛔ **Buktinya `4-erd-dan-tabel-datar/ERD.md` §2** — dokumen MENGIKAT yang menyatakan `ikut hapus`,
`tolak`, atau `putus` per relasi. Tidak satu baris pun di bawah diputuskan di modul ini; seluruhnya
dikutip dari sana, dan `TestPerilakuHapusSesuaiERD` mengadu ke-32 kunci asing dengan tabel itu.

| Awalan berkas | Relasi |
| --- | --- |
| `420_` | 21 relasi `ikut hapus` dari `ERD.md` §2 — lihat daftar di bawah |
| `422_` | `DETAIL_PROPORSIONAL` → `KELAS_BISNIS_LAYER` · `LAYER` → `KELOMPOK_LAYER` · `KELOMPOK_LAYER` → `KELAS_BISNIS_KELOMPOK` *(tiga)* |
| `424_` | `TERMIN` → `RINCIAN_ANGSURAN` *(satu)* |
| `426_` | `VERSI_KONTRAK` → `NILAI_SELISIH`, `NILAI_SEBELUM_PRO_RATE` *(dua)* — `ERD.md` §2.6 |
| `429_` | `BAGIAN` → `PENYEBARAN` · `DETAIL_PROPORSIONAL` → `PENYEBARAN` · `PENYEBARAN` → `RINCIAN_PENYEBARAN` · `RINCIAN_PENYEBARAN` → `NILAI_PENYEBARAN` *(empat)* — `ERD.md` §2.5 |
| `430_` | `M_TREATYIN_INSTALLMENT` → `M_TREATYIN_INSTALLMENTITEM` *(satu)* — ⚠️ **bukan dari `ERD.md` §2**, lihat di bawah |
| `437_` | delapan relasi antar tabel pendaratan anak *(delapan)* — ⚠️ **bukan dari `ERD.md` §2**, lihat di bawah |
| `438_` | `T_TREATY_LIMIT_DETAIL` → `T_TREATY_LIMIT_AMOUNT` · `T_TREATY_SHARE` → `T_TREATY_SHARE_AMOUNT` · `T_TREATY_FAC_SHARE` → `T_TREATY_FAC_SHARE_AMOUNT` *(tiga)* — alasan yang sama dengan `437_` |
| `439_` | `T_TREATY_LIMITS` → `T_TREATY_LIMIT_MEASURE` *(satu)* — alasan yang sama dengan `437_` |
| `438_` | `T_TREATY_LIMIT_DETAIL` → `T_TREATY_LIMIT_AMOUNT` · `T_TREATY_SHARE` → `T_TREATY_SHARE_AMOUNT` · `T_TREATY_FAC_SHARE` → `T_TREATY_FAC_SHARE_AMOUNT` *(tiga)* — ⚠️ **bukan dari `ERD.md` §2**, sebab yang sama dengan `437_` |

⚠️ **Seluruhnya di SATU berkas, `420_perilaku_hapus_erd.sql`, dan itu disengaja.** Kaskadenya
sempat disunting langsung ke migrasi `403`–`418`; sesudah pemilik proses menyatakan POOLDATA adalah
skema sasaran, suntingan itu dibalik menjadi migrasi korektif — berkas `400`–`419` kembali ke bentuk
yang tercatat di `T_MIGRASI`, dan perubahannya berdiri sebagai `DROP CONSTRAINT` + `ADD CONSTRAINT`
di `420`. Uraiannya di [`docs/KEPUTUSAN-PENYELARASAN-REPO.md`](docs/KEPUTUSAN-PENYELARASAN-REPO.md) §9.

Ke-21 relasi yang `420` jadikan `ON DELETE CASCADE`:

- `VERSI_KONTRAK` → `MATA_UANG_KONTRAK`, `RETENSI_CEDANT`, `EGNPI`, `PORTOFOLIO`,
  `PERIODE_PELAPORAN`, `PERIODE_AKUMULASI`, `TERMIN`, `SKALA_KOASURANSI`, `BATAS_PER_BAHAYA`,
  `DOKUMEN_KONTRAK` — `ERD.md` §2.3 *(sepuluh)*
- `VERSI_KONTRAK` → `LAYER` · `LAYER` → `DETAIL_PROPORSIONAL`, `BAGIAN` — §2.4 *(tiga)*
- `BAGIAN` → `POTONGAN` · `DETAIL_PROPORSIONAL` → `POTONGAN` — §2.5, satu konsep dua pelekatan *(dua)*
- `LAYER` → `PEMULIHAN_LIMIT` — §2.3b *(satu)*
- `LAYER` → `NILAI_MDP`, `NILAI_MDP_MINIMUM` · `BAGIAN` → `NILAI_PREMI_BRUTO`,
  `NILAI_PREMI_BRUTO_MINIMUM` · `DETAIL_PROPORSIONAL` → `NILAI_CADANGAN_PREMI` — §2.3c, tabel anak
  paket uang *(lima)*

⛔ **Keempat kaskade `422_` dan `424_` TIDAK dari `ERD.md` §2** — dan itu disebut di sini supaya
tidak terbaca sebagai kutipan yang tidak ada sumbernya. Keempat relasinya lahir sesudah §2 ditulis,
persis seperti delapan relasi §2.3c. Buktinya **dua dokumen yang sepakat**:

- `4-erd-dan-tabel-datar/ERD-TREATY-IN-DAN-EDM.html` — **ACUAN struktur sistem lama** sejak
  keputusan pemilik proses 2 Oktober 2026 (*"fokus ini aja, yang lainnya ada yang salah itu"*).
  Baris relasi **4**, **9**, **10**, **27**, seluruhnya `CASCADE`;
- aturan yang kelompoknya sendiri sudah nyatakan di `ERD.md`: §2.5 (anak cabang ikut hapus) dan
  §2.3 (anak langsung baris versi ikut hapus).

⚠️ Kolom `ON DELETE` di ERD HTML menyatakan dirinya ***"USULAN rancangan mengikuti konvensi
contooh.xlsx, bukan perilaku sistem lama"***, jadi ia **tidak dipakai sendirian**. Yang mengikat
untuk perilaku hapus tetap `ERD.md` §2; ERD HTML dipakai untuk **struktur** — induk, nama kolom
kunci asing, kardinalitas. Keduanya sepakat pada keempat baris ini.

**Satu tempat keduanya BERSELISIH, dan itu `PENCAPAIAN`.** `ERD.md` §2.9 menulis induk `KONTRAK`
dengan `[hapus: tolak]`; ERD HTML baris 6 menulis induk `T_TREATY_LIMIT_DETAIL` dengan `CASCADE`.
Yang menang `ERD.md` §2 — dan kali ini bukti SQL sejalan dengannya, sementara ERD HTML menandai
buktinya sendiri **DAUN-RELATIF** (*"lemah, mungkin memungut nama milik entitas lain"*). Uraiannya
di kepala migrasi `423_pencapaian.sql`.

⛔ **Kaskade `430_` BUKAN dari `ERD.md` §2 maupun dari ERD HTML, dan buktinya bukan dokumen.**
Kedelapan tabel `M_TREATYIN_*` adalah tabel **pendaratan** — ia mendaratkan larik di dalam
`M_TREATY_IN.JSONDATA`, bukan menggambar ulang entitas sistem lama — jadi tidak ada baris ERD yang
dapat dikutip untuknya. Yang mengikat adalah **bentuk datanya**: `InstallmentList` hidup DI DALAM
elemen `Installment`, dan rincian angsuran tanpa terminnya tidak berarti apa pun. Alasan yang sama
sudah dipakai migrasi `424_` untuk pasangan `TERMIN` → `RINCIAN_ANGSURAN`, yang memodelkan dua
tingkat yang sama persis.

⛔ **Kaskade `437_` berdiri di atas alasan yang SAMA dengan `430_`, dan buktinya juga bukan dokumen.**
Ketiga belas tabel yang `437_anak_treaty_in.sql` buat adalah tabel **pendaratan** tingkat kedua dan
ketiga — ia mendaratkan larik yang hidup DI DALAM elemen larik lain di `M_TREATY_IN.JSONDATA`.
Namanya dan induknya dikutip dari `Diagram-Skema-Tabel-TreatyIn-dan-EDM-v2.xlsx`, yang pemilik proses
tetapkan 5 Oktober 2026 sebagai sumber nama tabel modul ini; `ERD.md` §2 tidak memuat satu baris pun
untuknya, sebab ia menggambar entitas sistem lama, bukan bentuk dokumennya.

Yang mengikat adalah **bentuk datanya**, persis seperti `430_` dan `424_`: `Detail[]` hidup di dalam
elemen `Limits[]`, `COBList[]` di dalam elemen `Detail[]`, dan sebuah kelas bisnis tanpa detail
limitnya tidak berarti apa pun. Kedelapan relasinya:

⛔ Daftarnya ditulis sebagai butir, **bukan tabel** — pembaca bab ini (`kaskadePerModul`)
mengurai setiap tabel markdown di bawah judul ini sebagai daftar awalan berkas, dan tabel kedua di
sini akan dibacanya sebagai awalan yang cacat.

- `T_TREATY_LIMITS` → `T_TREATY_LIMIT_DETAIL`, `T_TREATY_LIMIT_GROUP`
- `T_TREATY_LIMIT_DETAIL` → `T_TREATY_LIMIT_COB`, `T_TREATY_LIMIT_ACHIEVEMENT`
- `T_TREATY_LIMIT_GROUP` → `T_TREATY_LIMIT_GROUP_COB`
- `T_TREATY_SHARE` → `T_TREATY_SHARE_SPREADING`, `T_TREATY_SHARE_DEDUCTION`
- `T_TREATY_FAC_SHARE` → `T_TREATY_FAC_SHARE_DEDUCTION`

⚠️ Dan batasnya sama pula: dari kontrak ke tabel **tingkat pertama** `437_` tidak memasang kunci
asing, dengan sebab yang dinyatakan di alinea berikut. `ikut hapus` di sana dijalankan pemuat.

⚠️ **Dari kontrak ke kedelapan tabel, `ikut hapus` TIDAK dijalankan basis data.** `MASTERID` bukan
kunci asing, sebab `POOLDATA.TREATY_IN` tidak punya kunci utama maupun `UNIQUE` pada `ID` — diukur
3 Oktober 2026 — sehingga Oracle menolak merujuknya (ORA-02270). Yang menjalankan `ikut hapus` di
sana adalah pemuat. Uraiannya di kepala `430_tabel_tab_treatyin.sql` dan di
[`docs/KEPUTUSAN-PENYELARASAN-REPO.md`](docs/KEPUTUSAN-PENYELARASAN-REPO.md) §12.

**Dua puluh satu relasi**, dari 28 yang `ERD.md` §2 nyatakan `ikut hapus`. Ketujuh sisanya menyentuh
tabel yang modul ini belum buat: `PENYEBARAN`, `RINCIAN_PENYEBARAN`, `NILAI_PENYEBARAN` (tiket `38`),
`NILAI_SELISIH` (Adjustment `06`), `PERISTIWA_KONTRAK` (yatim), dan `RETRO_KELUAR` (gelombang `G2`).

**Relasi yang SENGAJA tidak berkaskade** — `INV-18` menuntut **keduanya** dinyatakan, yang
berkaskade dan yang tidak. Sebelas kunci asing berperilaku **tolak**, diwujudkan dengan tidak
menulis klausa `ON DELETE`; bentuknya sama dengan bawaan Oracle, tetapi kini **dipilih** dan
sumbernya disebut:

- `KONTRAK` → `VERSI_KONTRAK` — §2.1, *"kontrak yang punya versi tidak boleh hilang, karena
  versinya memuat angka yang pernah dibukukan"*
- `VERSI_KONTRAK` → `JEJAK_PERUBAHAN` — §2.3, *"jejak yang dapat dihapus bersama bendanya bukan
  jejak"*
- delapan rujukan ke tabel acuan — §2.7, *"baris acuan yang sudah dipakai tidak dapat hilang"*
- `JENIS_REASURANSI` → dirinya sendiri — **tidak ada di `ERD.md`**; diputuskan di migrasi `400`
  memakai aturan kelompok §2.7, dan **ditagih** ke pemilik `ERD.md`

**Satu relasi `putus`**, satu-satunya di modul ini: `KONTRAK` → `KONTRAK` lewat
`ID_KONTRAK_DISALIN_DARI` (§2.1), `ON DELETE SET NULL` di migrasi `401`. Kontrak salinan bertahan
ketika asalnya dihapus; rujukannya saja yang dikosongkan.

## ⛔ `L-4` DICABUT — 3 Oktober 2026

`L-4` berbunyi: ***"tidak ada spesifikasi layar di mana pun."*** Atas dasar itu papan melarang
menyentuh UI, dan dua ronde menolak membangun layar.

**Pernyataan itu salah, dan sudah salah sejak awal.** Spesifikasinya ada di ekspor Pega 2026-09,
dan cacahnya dapat diperiksa dalam satu perintah:

| Modul | `Section/` | `Harness/` |
|---|---:|---:|
| Treaty In | **50** | **3** |
| Treaty In Adjustment | **68** | **6** |

```bash
ls "D:/XML_NURE/Treaty In/Section"/*.xml | wc -l            # 50
ls "D:/XML_NURE/Treaty In/Harness"/*.xml | wc -l            # 3
ls "D:/XML_NURE/Treaty In Adjustment/Section"/*.xml | wc -l # 68
ls "D:/XML_NURE/Treaty In Adjustment/Harness"/*.xml | wc -l # 6
```

Yang tidak ada bukan spesifikasinya, melainkan **pembacaannya**. Di dalam berkas itu tata letaknya
tertulis lengkap: `<pyValue>` properti yang diikat sebuah kontrol, `<pyLabelFor>` medan milik
sebuah label, `<pyPropertyTarget>` properti yang diisi pemilih, `<pyGridProps>` bentuk tabelnya.

> ### Larangan yang TETAP berlaku
>
> **Jangan mengarang layar.** Yang berubah hanya alasannya: dulu tidak boleh karena tidak ada
> acuannya, kini tidak boleh karena **acuannya ada dan harus dibaca**. Tiap medan yang dibangun
> membawa jejaknya — nama rule dan posisi bitanya — di `frontend/labels.ts`.

**Tiga jebakan yang dicatat saat pencabutan ini, supaya ronde berikutnya tidak kejeblos:**

1. **`<pyIncludedRuleXML>` dibuang lebih dulu.** Ekspor Pega menyematkan salinan utuh rule anak;
   sapuan datar membaca medan milik anak sebagai milik induk. Pada
   `Section/InputTreatyInOffer.xml` salinan itu **899.307 bita** dari 8.194.691.
2. **Berkas Section membundel indeks rule LAIN.** Caption `Treaty Contract Name`,
   `Contract Ref No`, `Teritorial Scope`, `Bordereaux Note`, `Accounting Mode`, dan `Treaty Year`
   muncul di `InputTreatyInOffer.xml` tetapi `pzIndexOwnerKey`-nya **`TREATYINNONPROPORTIONAL`**.
   Menyimpulkan pemilik sebuah medan dari **nama berkasnya** dapat salah total.
3. **Blok bergaris mati tidak dibangun.** Pega tidak dapat mengomentari tata letak, jadi blok
   dimatikan dengan kondisi mustahil — `1=2`, `1==2`, `Never` pada `pyContainerVisibleWhen` /
   `pyVisibleWhen` / `pyRowVisibleCondition`. ⚠️ `pyDisabledWhen>1=2` artinya **kebalikannya**:
   tidak pernah dinonaktifkan. Jangan tertukar.
