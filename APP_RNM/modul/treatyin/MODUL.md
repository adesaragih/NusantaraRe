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
Nol layar baru (`L-4`).

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
papan tiket melarang mengarang layar (`L-4`). Jalan keluarnya: satu halaman **baca-saja** atas data
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
| GROUPMENU | `TREATY` |
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
| `docs/` | `STRUKTUR-TABEL-TREATY-IN.md` (delapan tabel), `KEPUTUSAN-PENYELARASAN-REPO.md` (empat penyelarasan), `issues/` (tiket `14` dan `15`) |
| `backend/` | `modul.go` (`Pendaftaran()`), `models/`, `repository/`, `services/`, `handlers/`, `migrations/` |
| `frontend/` | `menu.ts`, `rute.tsx`, `api.ts`, `labels.ts`, `pages/AcuanTreatyIn.tsx` |

## Migrasi

Rentang `400-439` (tabel R2, urut hulu ke hilir: migrasi modul hilir yang merujuk tabel modul hulu
selalu berjalan sesudahnya). Slot menu `972-973` hanya menyalakan `DIMIGRASI` baris modul ini (satu `UPDATE`,
nol `INSERT` — menu datar 30-09-2026) saat modul mendapat layar pertamanya, di folder
`backend/migrations/` modul ini sendiri — bentuk SQL-nya di `APP_RNM/PANDUAN-DEPLOY-DAN-GIT-PER-MODUL.md`
bab 6. Nomor selalu tiga digit.

## Pernyataan untuk penjaga

### Kaskade ON DELETE CASCADE

Kaskade HANYA pada berkas migrasi modul ini yang berawalan di bawah; berkas lain modul ini tanpa
`ON DELETE CASCADE` (`TestKaskadeHanyaPadaRelasiTerdaftar`). Mendaftarkan yang baru menuntut bukti.

⛔ **Buktinya `4-erd-dan-tabel-datar/ERD.md` §2** — dokumen MENGIKAT yang menyatakan `ikut hapus`,
`tolak`, atau `putus` per relasi. Tidak satu baris pun di bawah diputuskan di modul ini; seluruhnya
dikutip dari sana, dan `TestPerilakuHapusSesuaiERD` mengadu ke-32 kunci asing dengan tabel itu.

| Awalan berkas | Relasi |
| --- | --- |
| `420_` | seluruh 21 relasi `ikut hapus` — lihat daftar di bawah |

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
