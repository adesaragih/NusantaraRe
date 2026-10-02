# Struktur Tabel — Treaty In

Acuan bentuk tabel untuk aplikasi Go. Dibuat 1 Oktober 2026, bersamaan migrasi `400`–`402`
(tiket `14` dan `15`).

Berkas ini menggambarkan **BENTUK**, bukan alasan — alasannya ada di tiket `docs/issues/` dan di
spec `D:\XML_NURE\_migration-docs\treaty-in\`.

⚠️ **Acuan nama kolom yang MENGIKAT adalah `2-to-spec/KAMUS-KOLOM.md`** (urutan wewenang butir 5),
dan berkas ini dibangkitkan dari DDL `backend/migrations/` yang mengikutinya. Ejaan yang muncul di
kolom **Sumber** sebagai `korpus` adalah **ejaan properti Pega** — itu bukti asal kolom, bukan nama
kolom.

**Tipe ditulis sebagai kategori logis:** teks · angka desimal · bilangan bulat · DATE.
Presisi fisiknya ada di DDL, dan empat penyelarasannya dengan repo — termasuk **uang
`NUMBER(38,8)`, bukan `NUMBER(38,20)`** — tercatat di
[`KEPUTUSAN-PENYELARASAN-REPO.md`](KEPUTUSAN-PENYELARASAN-REPO.md). **Jangan menyunting tipe di sini
tanpa membacanya.**

⛔ **Berkas ini MENGIKAT.** Modul ini punya `backend/modul.go`, sehingga `TestKolomDDLCocokDenganStruktur`
dan `TestGolonganTipeDDLCocokDenganStruktur` membandingkannya dengan DDL kolom demi kolom
(`PANDUAN-TIM-PER-MODUL.md` bab 4.3). Menyunting tabel di bawah tanpa menyunting migrasinya
membuat penjaga merah.

⚠️ **Tidak satu baris pun DDL ini pernah dijalankan di Oracle mana pun** (`L-3`). Yang terjamin
adalah BENTUKNYA — nama, tipe, dan constraint yang tertulis — bukan bahwa Oracle menerimanya.

## Cacah

| | Tabel | Kolom | Tiket |
| --- | --- | ---: | --- |
| Tabel acuan | `MATA_UANG`, `JENIS_POTONGAN`, `KELAS_BISNIS`, `KELOMPOK_TREATY`, `BAHAYA`, `JENIS_REASURANSI` | 25 | `15` |
| Identitas kontrak | `KONTRAK`, `VERSI_KONTRAK` | 57 | `14` · `01` |
| Uang dan periode | `MATA_UANG_KONTRAK`, `RETENSI_CEDANT`, `EGNPI`, `TERMIN` | 31 | `20` `22` `23` `27` |
| Periode dan portofolio | `PORTOFOLIO`, `PERIODE_PELAPORAN`, `PERIODE_AKUMULASI` | 18 | `24` `25` `26` |
| Tanggungan | `SKALA_KOASURANSI`, `BATAS_PER_BAHAYA`, `DOKUMEN_KONTRAK` | 15 | `28` `29` `30` |
| Layer | `LAYER`, `NILAI_MDP`, `NILAI_MDP_MINIMUM`, `PEMULIHAN_LIMIT` | 38 | `31` |
| Jejak | `JEJAK_PERUBAHAN` | 8 | `39` |
| Bagian NuRe *(non-proporsional)* | `BAGIAN`, `NILAI_PREMI_BRUTO`, `NILAI_PREMI_BRUTO_MINIMUM` | 16 | `33` |
| Ketentuan proporsional | `DETAIL_PROPORSIONAL`, `NILAI_CADANGAN_PREMI` | 16 | `34` |
| Potongan | `POTONGAN` | 6 | `37` |

**29 tabel, 230 kolom.** Cacah kolom tiap tabel **sama persis** dengan yang `KAMUS-KOLOM.md` tulis di
judul §-nya masing-masing — kecuali `VERSI_KONTRAK`, yang bertambah satu (lihat di bawah).

`KONTRAK` 8 kolom dan `VERSI_KONTRAK` **49** kolom. Ke-48 pertamanya **sama persis** dengan cacah yang
`KAMUS-KOLOM.md` §10.1 dan §10.2 tulis di judulnya sendiri; yang ke-49, `ID_VERSI_KONTRAK_DASAR`,
ditambahkan migrasi `440` modul `treatyinadjustment` (tiket `01`) dan **belum ada di `KAMUS-KOLOM.md`**
— utang hulu yang ditagih ke pemilik `SPEC-MODEL-DATA.md` §10.2.

**Sequence:** `SEQ_TRIN_KONTRAK`, `SEQ_TRIN_VERSI_KONTRAK`, `SEQ_TRIN_MATA_UANG`,
`SEQ_TRIN_JENIS_POTONGAN`, `SEQ_TRIN_JENIS_REASURANSI`, `SEQ_TRIN_BAHAYA`,
`SEQ_TRIN_KELOMPOK_TREATY`, `SEQ_TRIN_KELAS_BISNIS` (berkas `402`), ditambah satu per tabel anak di
berkas `415`: `SEQ_TRIN_MATA_UANG_KONTRAK`, `SEQ_TRIN_RETENSI_CEDANT`, `SEQ_TRIN_EGNPI`,
`SEQ_TRIN_PORTOFOLIO`, `SEQ_TRIN_PERIODE_PELAPORAN`, `SEQ_TRIN_PERIODE_AKUMULASI`, `SEQ_TRIN_TERMIN`,
`SEQ_TRIN_SKALA_KOASURANSI`, `SEQ_TRIN_BATAS_PER_BAHAYA`, `SEQ_TRIN_DOKUMEN_KONTRAK`,
`SEQ_TRIN_LAYER`, `SEQ_TRIN_NILAI_MDP`, `SEQ_TRIN_NILAI_MDP_MINIMUM`, `SEQ_TRIN_PEMULIHAN_LIMIT`,
`SEQ_TRIN_JEJAK_PERUBAHAN`; dan di berkas `419`: `SEQ_TRIN_BAGIAN`, `SEQ_TRIN_NILAI_PREMI_BRUTO`,
`SEQ_TRIN_NILAI_PB_MINIMUM`, `SEQ_TRIN_DETAIL_PROPORSIONAL`, `SEQ_TRIN_NILAI_CADANGAN_PREMI`,
`SEQ_TRIN_POTONGAN`. **Dua puluh sembilan**, seluruhnya `NOCACHE NOCYCLE` (INV-02, INV-03).

⚠️ `SEQ_TRIN_NILAI_PB_MINIMUM` **disingkat**: bentuk penuhnya 34 bita, melewati batas 30 (§16).
Singkatannya jatuh di `NILAI_PREMI_BRUTO` → `NILAI_PB`, bukan di `MINIMUM` — tanpa "MINIMUM" kedua
sequence tidak dapat dibedakan.

⛔ **Belum satu pun terpakai.** Tidak ada jalur tulis di modul ini, sehingga tuntutan tiket `14`
*"tidak ada jalur lain yang dapat memberi pengenal"* benar secara hampa. Ia ditegakkan bersama jalur
simpan.

**Index:** `IX_VERSI_KONTRAK_KONTRAK` · `IX_VERSI_KONTRAK_DASAR` *(dipasang migrasi `440` modul
`treatyinadjustment`)* · `IX_NILAI_MDP_LAYER` · `IX_NILAI_MDP_MIN_LAYER` · `IX_PEMULIHAN_LIMIT_LAYER` ·
`IX_JEJAK_PERUBAHAN_VERSI` · `IX_MATA_UANG_KONTRAK_MU` · `IX_RETENSI_CEDANT_KLP` ·
`IX_EGNPI_KELOMPOK` · `IX_EGNPI_KELAS_BISNIS` · `IX_BATAS_PER_BAHAYA_BHY` · `IX_DETAIL_PROP_KLP` ·
`IX_POTONGAN_JENIS` · `IX_NILAI_PB_BAGIAN` · `IX_NILAI_PB_MIN_BAGIAN` · `IX_NILAI_CAD_PREMI_DP`.

ℹ️ **Kenapa sebagian kunci asing ber-index sendiri dan sebagian tidak.** Oracle membuat index untuk
setiap `UNIQUE`, dan index itu **melayani** sebuah kunci asing bila kolom kunci asingnya **memimpin**
kunci alami tersebut. Sebelas tabel anak berkunci asing ke `VERSI_KONTRAK`, dan pada sebelas-belasnya
`ID_VERSI_KONTRAK` adalah kolom **pertama** kunci alaminya — index itu sudah melayaninya.

Yang **tidak** demikian ada dua rupa, dan keduanya ber-index eksplisit:

| Rupa | Index |
| --- | --- |
| kunci asing ke **tabel acuan** — tidak pernah memimpin kunci alami | `IX_MATA_UANG_KONTRAK_MU`, `IX_RETENSI_CEDANT_KLP`, `IX_EGNPI_KELOMPOK`, `IX_EGNPI_KELAS_BISNIS`, `IX_BATAS_PER_BAHAYA_BHY` |
| anak yang **tidak punya kunci alami** sama sekali | `IX_NILAI_MDP_LAYER`, `IX_NILAI_MDP_MIN_LAYER`, `IX_PEMULIHAN_LIMIT_LAYER`, `IX_JEJAK_PERUBAHAN_VERSI`, `IX_VERSI_KONTRAK_KONTRAK` |

⚠️ **Kelima index tabel acuan tidak ada di `ddl-usulan/`.** Ia ditambahkan di sini sebab kunci asing
tanpa index membuat penghapusan baris acuan — satu `BAHAYA`, satu `MATA_UANG` — memindai seluruh tabel
anak, dan pada sebagian versi Oracle menguncinya. Dicatat di
[`KEPUTUSAN-PENYELARASAN-REPO.md`](KEPUTUSAN-PENYELARASAN-REPO.md) butir 5. Dijaga
`TestSetiapKunciAsingTerlayaniIndex`.

## Relasi

| Anak | Kolom | Induk | ON DELETE | Invarian |
| --- | --- | --- | --- | --- |
| `VERSI_KONTRAK` | `ID_KONTRAK` | `KONTRAK` | **tanpa** | INV-18 |
| `VERSI_KONTRAK` | `KODE_MATA_UANG_KONTRAK` | `MATA_UANG` | **tanpa** | INV-44 |
| `VERSI_KONTRAK` | `ID_VERSI_KONTRAK_DASAR` | `VERSI_KONTRAK` | **tanpa** | INV-18 |
| `JENIS_REASURANSI` | `ID_INDUK` | `JENIS_REASURANSI` | **tanpa** | INV-18 |
| `MATA_UANG_KONTRAK` | `ID_VERSI_KONTRAK` · `KODE_MATA_UANG` | `VERSI_KONTRAK` · `MATA_UANG` | **tanpa** | INV-18 · INV-44 |
| `RETENSI_CEDANT` | `ID_VERSI_KONTRAK` · `ID_KELOMPOK_TREATY` | `VERSI_KONTRAK` · `KELOMPOK_TREATY` | **tanpa** | INV-18 |
| `EGNPI` | `ID_VERSI_KONTRAK` · `ID_KELOMPOK_TREATY` · `ID_KELAS_BISNIS` | `VERSI_KONTRAK` · `KELOMPOK_TREATY` · `KELAS_BISNIS` | **tanpa** | INV-18 |
| `PORTOFOLIO` · `PERIODE_PELAPORAN` · `PERIODE_AKUMULASI` · `TERMIN` · `SKALA_KOASURANSI` · `DOKUMEN_KONTRAK` · `LAYER` · `JEJAK_PERUBAHAN` | `ID_VERSI_KONTRAK` | `VERSI_KONTRAK` | **tanpa** | INV-18 |
| `BATAS_PER_BAHAYA` | `ID_VERSI_KONTRAK` · `ID_BAHAYA` | `VERSI_KONTRAK` · `BAHAYA` | **tanpa** | INV-18 |
| `NILAI_MDP` · `NILAI_MDP_MINIMUM` · `PEMULIHAN_LIMIT` · `BAGIAN` · `DETAIL_PROPORSIONAL` | `ID_LAYER` | `LAYER` | **tanpa** | INV-18 |
| `DETAIL_PROPORSIONAL` | `ID_KELOMPOK_TREATY` | `KELOMPOK_TREATY` | **tanpa** | INV-18 |
| `NILAI_PREMI_BRUTO` · `NILAI_PREMI_BRUTO_MINIMUM` | `ID_BAGIAN` | `BAGIAN` | **tanpa** | INV-18 |
| `NILAI_CADANGAN_PREMI` | `ID_DETAIL_PROPORSIONAL` | `DETAIL_PROPORSIONAL` | **tanpa** | INV-18 |
| `POTONGAN` | `ID_BAGIAN` · `ID_DETAIL_PROPORSIONAL` · `ID_JENIS_POTONGAN` | `BAGIAN` · `DETAIL_PROPORSIONAL` · `JENIS_POTONGAN` | **tanpa** | INV-18 · KTV-B |

⛔ **Nol `ON DELETE` di seluruh modul** — bawaan Oracle MENOLAK, dan menolak yang dikehendaki.
Menghapus versi kontrak yang masih punya layer, termin, atau jejak akan gagal dengan ORA-02292,
bukan diam-diam membawa riwayatnya ikut hilang. Diperiksa `TestNolKaskadeHapus`.

⛔ **`DOKUMEN_KONTRAK.ID_DOKUMEN` sengaja TANPA kunci asing**: ia menunjuk dokumen di penyimpanan
berkas, bukan baris di skema ini — itu pokok tiket `30`, dokumennya **dirujuk, bukan disalin**.

⚠️ **Baris `ID_VERSI_KONTRAK_DASAR` dipasang modul `treatyinadjustment`** (migrasi `440`, tiket `01`), bukan modul ini.
Tabelnya milik modul ini; satu kolomnya dibawa modul itu. Sebabnya: papan tiket dipisahkan
25-09-2026, **model datanya tidak** — `KAMUS-KOLOM.md` satu untuk keduanya.

## Kunci alami

| Tabel | Kolom | Invarian |
| --- | --- | --- |
| `VERSI_KONTRAK` | `ID_KONTRAK` + `NOMOR_URUT_VERSI` | INV-04 |
| keenam tabel acuan | `KODE` | INV-68 |
| `LAYER` | `ID_VERSI_KONTRAK` + `NOMOR_LAYER` + `BAGIAN_LAYER` | INV-05 |
| `MATA_UANG_KONTRAK` | `ID_VERSI_KONTRAK` + `KODE_MATA_UANG` | INV-07 |
| `RETENSI_CEDANT` · `EGNPI` | `ID_VERSI_KONTRAK` + `ID_KELOMPOK_TREATY` + `KODE_MATA_UANG` | INV-08 · INV-09 |
| `PERIODE_PELAPORAN` · `PERIODE_AKUMULASI` | `ID_VERSI_KONTRAK` + `PERIODE` | INV-10 · INV-11 |
| `TERMIN` | `ID_VERSI_KONTRAK` + `NOMOR_TERMIN` + `KODE_MATA_UANG` | INV-12 |
| `SKALA_KOASURANSI` | `ID_VERSI_KONTRAK` + `PERSEN_LIMIT` | INV-13 |
| `BATAS_PER_BAHAYA` | `ID_VERSI_KONTRAK` + `ID_BAHAYA` | INV-14 |
| `PORTOFOLIO` | `ID_VERSI_KONTRAK` + `ARAH_PORTOFOLIO` + `JENIS_PORTOFOLIO` | INV-66 |
| `DOKUMEN_KONTRAK` | `ID_VERSI_KONTRAK` + `ID_DOKUMEN` | INV-67 |
| `BAGIAN` | `ID_LAYER` — satu bagian per layer | INV-64 |
| `DETAIL_PROPORSIONAL` | `ID_LAYER` + `ID_KELOMPOK_TREATY` — **di dalam LAYER**, bukan versi | INV-06 |
| `POTONGAN` | `ID_BAGIAN` + `ID_JENIS_POTONGAN` **dan** `ID_DETAIL_PROPORSIONAL` + `ID_JENIS_POTONGAN` — **dua** `UNIQUE` | INV-15 |

⛔ **Mata uang IKUT di dalam kunci alami INV-08, INV-09, dan INV-12**, dan itu bukan hiasan:
`Z00_KUNCI_ALAMI.sql` memperingatkan bahwa menulisnya TANPA mata uang mengubah artinya menjadi
*"dilarang dua baris bermata uang berbeda"* — yang **menolak data yang sah**.

### Empat tabel tanpa `UNIQUE`, dan sebabnya BERBEDA-BEDA

> ⚠️ **RALAT.** Keempatnya pernah dijelaskan di sini dengan **satu** sebab yang sama — *"peristiwa yang
> sama dapat terjadi dua kali"*. Itu keliru untuk tiga dari empat, dan `Z00_KUNCI_ALAMI.sql` memberi
> sebab yang berbeda untuk masing-masing.

| Tabel | Sebab | Sumber |
| --- | --- | --- |
| `JEJAK_PERUBAHAN` | **tidak ada** kunci alami, dan itu keputusan — peristiwa yang sama dapat terjadi dua kali pada versi yang sama | `Z00_KUNCI_ALAMI.sql` |
| `PEMULIHAN_LIMIT` | kunci alaminya **belum bernomor**; entitasnya baru diterima §10.23c | `Z00_KUNCI_ALAMI.sql` |
| `NILAI_MDP` · `NILAI_MDP_MINIMUM` | **punya** kunci alami yang belum bernomor — `KODE_MATA_UANG` *"kunci alami di dalam induknya"*. `Z00` tidak menyebut keduanya sama sekali | `KAMUS-KOLOM.md` |

`TestTabelTanpaKunciAlamiTidakDiberiDiamDiam` menjaga satu hal saja: tidak ada yang **menyisipkan**
`UNIQUE` sebelum nomor invariannya turun. Ia **tidak** menyatakan ketiadaannya benar.

⛔ **`UQ_LAYER` tidak menegakkan INV-05 sepenuhnya.** `BAGIAN_LAYER` boleh kosong, dan Oracle
memperlakukan NULL sebagai **tidak sama dengan** NULL di kunci unik komposit — sehingga baris
`(versi, 1, NULL)` diterima berulang, yaitu keadaan layer tanpa bagian, yang justru paling lazim.
Lubang itu ada di `Z00_KUNCI_ALAMI.sql`, bukan dibuat di sini; uraian dan tagihannya di dalam
migrasi `413`. **Daftar periksa tiket `31` "INV-05 terpasang" belum terpenuhi.**

⛔ **`KONTRAK` sengaja TIDAK berkunci alami.** Cedant + asal bisnis + periode + sifat proporsi adalah
kunci alaminya, dan ADR-0040 §2 **MEMPERINGATKAN, tidak melarang** — ketiadaan nomor INV-nya
disengaja. Peringatan ganda itu pekerjaan tiket `16`.

## Yang BELUM ada di sini, dan tiket mana membuatnya

| Yang belum ada | Tiket |
| --- | --- |
| Daftar nilai sah `KEADAAN_SIKLUS_HIDUP` dan mesin perpindahannya | `45` |
| `DOKUMEN_ADDENDUM` + FK `VERSI_KONTRAK.ID_DOKUMEN_ADDENDUM` | `04` (papan Adjustment, tertahan `DB-16a`) |
| Peringatan kunci alami kontrak ganda | `16` |
| Pencarian lewat `NOMOR_KONTRAK_WARISAN` | `17` |
| Pembekuan kunci alami sesudah kontrak lahir | `18` |
| `PENYEBARAN`, `RINCIAN_PENYEBARAN`, `NILAI_PENYEBARAN` | `38` (tertahan `Uji AD`, `L-3`) |
| `CATATAN_PERSETUJUAN` | `54` |
| `PERISTIWA_KONTRAK` | **nol tiket menyebutnya** — yatim, ditagih pemilik proses |
| Isi keenam tabel acuan, dipindahkan dari sistem lama | `44` |
| Syarat berbeda tiap pemulihan limit | `32` |

### Invarian yang tabelnya berdiri tetapi BELUM ditegakkan di mana pun

Seluruhnya membandingkan baris di tabel BERBEDA, yang `CHECK` Oracle tidak dapat nyatakan; satu-satunya
bentuk basis data yang bisa adalah trigger, dan ADR-0056 (K-4) melarangnya. Keempatnya menunggu jalur
simpan yang sama.

| Invarian | Apa | Berdiri di |
| --- | --- | --- |
| `INV-29` | `SIFAT_PROPORSI` dua nilai | `401` |
| `INV-53` | `TANGGAL_MULAI` ≤ `TANGGAL_BERAKHIR` | `401` |
| `INV-55` | periode pelaporan di dalam periode kontrak | `407` |
| `INV-56` | periode akumulasi di dalam periode kontrak | `408` |
| `INV-38` | kurs bukan nol, dan kurs tidak dipaksakan menjadi satu | `403` |
| `INV-57` | tanggal kurs tidak lebih akhir daripada transaksi yang memakainya | `403` |
| `INV-05` | nomor layer + bagian unik di dalam versi — **sebagian**, lihat di atas | `413` |
| `INV-39` · `INV-40` | paket uang dipasangkan dengan tingkat pencatatannya | `404` dan seterusnya |
| `INV-49` | pengecualian bernama terhadap `INV-47` | menunggu tiket `47` |
| `INV-30` | `JENIS_TREATY` dua nilai | `417` |
| `INV-32` · `INV-33` | cabang proporsional/non-proporsional dipisahkan | `417` · `416` — **tertahan `F-13`** |
| `INV-41` · `INV-36` | paket uang terisi menuntut mata uangnya | `417` · `416` |
| `INV-52` | premi bruto dikurangi potongan sama dengan premi bersih | `418` — **tertahan `F-13`** |
| `INV-63` | rumus potongan tepat satu kali di seluruh basis kode | `418` — tinjauan kode, bukan constraint |

⛔ **`INV-32`, `INV-33`, dan `INV-52` tertahan satu hal yang sama: `F-13`.** Ketiganya bergolongan
**indeks unik**, yang di Oracle berarti *materialized view* ber-`REFRESH ON COMMIT`, dan
`F-13-MV-TANPA-PEMANTAU-KEBASIAN.md` menyatakan penegakan lewat MV di modul ini belum punya pemantau
kebasian — *"MV yang gagal me-refresh berhenti menegakkan tanpa satu galat pun."* Uraiannya di dalam
migrasi `416`.

⚠️ **`INV-29` dan `INV-38` BERBEDA dari yang lain di tabel ini**: keduanya membandingkan kolom pada
**baris yang sama**, sehingga `CHECK` Oracle dapat menyatakannya. Yang menahan keduanya hanya
ADR-0056 (K-4), bukan ketidakmampuan basis data. Bila pemilik proses memutuskan `CHECK` boleh untuk
invarian sebaris, keduanya calon pertama.

---

## MATA_UANG

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_MATA_UANG` | bilangan bulat | tidak | PK | keputusan tiket 15 — ADR-0038, himpunan yang dapat bertambah |
| `KODE` | teks | tidak | UQ | keputusan tiket 15 — ADR-0038, himpunan yang dapat bertambah — **kunci alami**, **Bentuk — tak berinduk**: unik di seluruh tabel |
| `NAMA` | teks | tidak |  | keputusan tiket 15 — ADR-0038, himpunan yang dapat bertambah |
| `AKTIF` | teks | tidak |  | keputusan tiket 15 — ADR-0038, himpunan yang dapat bertambah — nilai lama **tidak pernah dihapus**; ia dimatikan |

## JENIS_POTONGAN

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_JENIS_POTONGAN` | bilangan bulat | tidak | PK | keputusan tiket 15 — ADR-0038, himpunan yang dapat bertambah |
| `KODE` | teks | tidak | UQ | keputusan tiket 15 — ADR-0038, himpunan yang dapat bertambah — **kunci alami**, **Bentuk — tak berinduk**: unik di seluruh tabel |
| `NAMA` | teks | tidak |  | keputusan tiket 15 — ADR-0038, himpunan yang dapat bertambah |
| `AKTIF` | teks | tidak |  | keputusan tiket 15 — ADR-0038, himpunan yang dapat bertambah — nilai lama **tidak pernah dihapus**; ia dimatikan |

## KELAS_BISNIS

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_KELAS_BISNIS` | bilangan bulat | tidak | PK | keputusan tiket 15 — ADR-0038, himpunan yang dapat bertambah |
| `KODE` | teks | tidak | UQ | keputusan tiket 15 — ADR-0038, himpunan yang dapat bertambah — **kunci alami**, **Bentuk — tak berinduk**: unik di seluruh tabel |
| `NAMA` | teks | tidak |  | keputusan tiket 15 — ADR-0038, himpunan yang dapat bertambah |
| `AKTIF` | teks | tidak |  | keputusan tiket 15 — ADR-0038, himpunan yang dapat bertambah — nilai lama **tidak pernah dihapus**; ia dimatikan |

## KELOMPOK_TREATY

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_KELOMPOK_TREATY` | bilangan bulat | tidak | PK | keputusan tiket 15 — ADR-0038, himpunan yang dapat bertambah |
| `KODE` | teks | tidak | UQ | keputusan tiket 15 — ADR-0038, himpunan yang dapat bertambah — **kunci alami**, **Bentuk — tak berinduk**: unik di seluruh tabel |
| `NAMA` | teks | tidak |  | keputusan tiket 15 — ADR-0038, himpunan yang dapat bertambah |
| `AKTIF` | teks | tidak |  | keputusan tiket 15 — ADR-0038, himpunan yang dapat bertambah — nilai lama **tidak pernah dihapus**; ia dimatikan |

## BAHAYA

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_BAHAYA` | bilangan bulat | tidak | PK | keputusan tiket 15 — ADR-0038, himpunan yang dapat bertambah |
| `KODE` | teks | tidak | UQ | keputusan tiket 15 — ADR-0038, himpunan yang dapat bertambah — **kunci alami**, **Bentuk — tak berinduk**: unik di seluruh tabel |
| `NAMA` | teks | tidak |  | keputusan tiket 15 — ADR-0038, himpunan yang dapat bertambah |
| `AKTIF` | teks | tidak |  | keputusan tiket 15 — ADR-0038, himpunan yang dapat bertambah — nilai lama **tidak pernah dihapus**; ia dimatikan |

## JENIS_REASURANSI

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_JENIS_REASURANSI` | bilangan bulat | tidak | PK | keputusan tiket 15 — ADR-0038, himpunan yang dapat bertambah |
| `KODE` | teks | tidak | UQ | keputusan tiket 15 — ADR-0038, himpunan yang dapat bertambah — **kunci alami**, **Bentuk — tak berinduk**: unik di seluruh tabel |
| `NAMA` | teks | tidak |  | keputusan tiket 15 — ADR-0038, himpunan yang dapat bertambah |
| `AKTIF` | teks | tidak |  | keputusan tiket 15 — ADR-0038, himpunan yang dapat bertambah — nilai lama **tidak pernah dihapus**; ia dimatikan |
| `ID_INDUK` | bilangan bulat | ya | FK | keputusan tiket 15 — ADR-0038, himpunan yang dapat bertambah — **hanya `JENIS_REASURANSI`** — ia bersusun, §10.6 |

## KONTRAK

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_KONTRAK` | bilangan bulat | tidak | PK | keputusan tiket 14 — **baru**, tidak ada di sistem lama — pengenal buatan sistem; bukan dari teks, bukan dari cap waktu |
| `NOMOR_KONTRAK_WARISAN` | teks | ya |  | korpus `ID` — hanya terisi pada baris hasil migrasi (ADR-0042) |
| `ID_KONTRAK_DISALIN_DARI` | bilangan bulat | ya |  | korpus `OLDID` pecahan 1` — rujukan ke `KONTRAK` lain |
| `ID_CEDANT` | bilangan bulat | tidak |  | korpus `CedingID` — nama cedant **tidak disalin** |
| `ID_ASAL_BISNIS` | bilangan bulat | tidak |  | korpus `LeadingReinsSourceID` — *source of business* |
| `SIFAT_PROPORSI` | teks | tidak |  | korpus `ProportionType` — `PROPORSIONAL` / `NON_PROPORSIONAL` |
| `TANGGAL_MULAI` | DATE | tidak |  | korpus `Commencement` — batas **inklusif** (ADR-0022) |
| `TANGGAL_BERAKHIR` | DATE | tidak |  | korpus `Termination` — batas **inklusif** |

## VERSI_KONTRAK

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_VERSI_KONTRAK` | bilangan bulat | tidak | PK | keputusan tiket 14 — **baru**, tidak ada di sistem lama |
| `ID_KONTRAK` | bilangan bulat | tidak | FK UQ | keputusan tiket 14 — **baru**, tidak ada di sistem lama |
| `NOMOR_URUT_VERSI` | bilangan bulat | ya | UQ | keputusan tiket 14 — **baru**, tidak ada di sistem lama |
| `KEADAAN_SIKLUS_HIDUP` | teks | tidak |  | korpus `Position` + `StatusAkseptasi` + 3 bendera` — ⚠ **DELAPAN nilai** — `DIBATALKAN` ditambahkan ADR-0055 perubahan 24 Sep 2026 |
| `KEADAAN_WARISAN_ASLI` | teks | ya |  | keputusan tiket 14 — **baru**, tidak ada di sistem lama |
| `JENIS_ADDENDUM` | teks | ya |  | korpus `EDMState` |
| `TANGGAL_BERLAKU_ADDENDUM` | DATE | ya |  | korpus `EDMEffective` |
| `NAMA_KONTRAK` | teks | tidak |  | korpus `TreatyContractName` |
| `LINGKUP_WILAYAH` | teks | ya |  | korpus `TeritorialScope` |
| `KELAS_BISNIS_KONTRAK` | bilangan bulat | ya |  | korpus `ClassofBusiness` |
| `KODE_MATA_UANG_KONTRAK` | bilangan bulat | tidak |  | korpus `Currency` |
| `ID_REASURADUR_PEMIMPIN` | bilangan bulat | ya |  | korpus `LeadingReinsID` |
| `ID_KETUA_TREATY` | bilangan bulat | ya |  | korpus `TreatyLeader` |
| `PERSEN_BAGIAN_NURE` | angka desimal | tidak |  | korpus `RNMShare` + `RNMShareP` |
| `BAGIAN_NURE_SERAGAM` | teks | tidak |  | korpus `RNMShareAcrossTheBoard` |
| `PERSEN_BAGIAN_NURE_DIPOTONG` | angka desimal | ya |  | korpus `RnmShareDeducted` |
| `PERSEN_BROKERAGE` | angka desimal | ya |  | korpus `BrokeragePercent` + `…P` |
| `PERSEN_BAGIAN_FAKULTATIF` | angka desimal | ya |  | korpus `FacultativeShare` (`FacShare` dibuang, salinan — §14.3)` |
| `PERSEN_BROKERAGE_FAKULTATIF` | angka desimal | ya |  | korpus `FacultativeShareBrokerage` (`FacShareBrokerage` dibuang, salinan)` |
| `PENGECUALIAN` | teks | ya |  | korpus `Exclusions` + `ExclusionsP` |
| `KETENTUAN_KHUSUS` | teks | ya |  | korpus `SpecialConditions` + `SpecialConditionsP` |
| `KETERANGAN` | teks | ya |  | korpus `Information` |
| `CATATAN` | teks | ya |  | korpus `Comment` |
| `CATATAN_BORDEREAUX` | teks | ya |  | korpus `BordereauxNote` |
| `MEMAKAI_BORDEREAUX` | teks | tidak |  | korpus `Bordeaux` |
| `CARA_PEMBUKUAN` | teks | tidak |  | korpus `AccountingMode` |
| `CARA_PEMBUKUAN_XOL` | teks | ya |  | korpus `AccountingModeNonProp` |
| `PERIODE_PELAPORAN_KONTRAK` | teks | ya |  | korpus `ReportingPeriod` |
| `SELANG_PELAPORAN` | bilangan bulat | ya |  | korpus `ReportingInterval` |
| `TANGGAL_MULAI_PELAPORAN` | DATE | ya |  | korpus `ReportingStart` |
| `TANGGAL_AKHIR_PELAPORAN` | DATE | ya |  | korpus `ReportingEnd` |
| `HARI_BATAS_PENYERAHAN` | bilangan bulat | ya |  | korpus `ReportingSubmission` |
| `HARI_BATAS_KONFIRMASI` | bilangan bulat | ya |  | korpus `ReportingConfirmation` |
| `HARI_BATAS_PELUNASAN` | bilangan bulat | ya |  | korpus `ReportingSettlement` |
| `HARI_PENGINGAT` | bilangan bulat | ya |  | korpus `ReminderDays` |
| `PERIODE_AKUMULASI_KONTRAK` | teks | ya |  | korpus `AccumulationPeriod` |
| `JUMLAH_TERMIN` | bilangan bulat | ya |  | korpus `InstallmentNo` |
| `MEMAKAI_PRORATA` | teks | tidak |  | korpus `IsProRate` |
| `BATAS_MAKSIMUM_KELOMPOK` | angka desimal | ya |  | korpus `MaxCoGroup` |
| `MATA_UANG_BATAS_KELOMPOK` | teks | ya |  | korpus `**tidak ada di sistem lama**` — **denominasi `BATAS_MAKSIMUM_KELOMPOK`** — `KTV-C`. Sistem lama tidak merekam mata uang untuk besaran ini sama sekali; kolomnya disediakan sekarang karena menambahkannya sesudah data masuk mahal. **Dicabut sebelum data dimuat bila `T-6` menyatakan ia selalu mata uang kontrak** |
| `BATAS_MAKSIMUM_NON_KELOMPOK` | angka desimal | ya |  | korpus `MaxCoNonGroup` |
| `MATA_UANG_BATAS_NON_KELOMPOK` | teks | ya |  | korpus `**tidak ada di sistem lama**` — **denominasi `BATAS_MAKSIMUM_NON_KELOMPOK`** — `KTV-C`. Sistem lama tidak merekam mata uang untuk besaran ini sama sekali; kolomnya disediakan sekarang karena menambahkannya sesudah data masuk mahal. **Dicabut sebelum data dimuat bila `T-6` menyatakan ia selalu mata uang kontrak** |
| `BATAS_PILIHAN` | angka desimal | ya |  | korpus `OptionLimit` |
| `MATA_UANG_BATAS_PILIHAN` | teks | ya |  | korpus `**tidak ada di sistem lama**` — **denominasi `BATAS_PILIHAN`** — `KTV-C`. Sistem lama tidak merekam mata uang untuk besaran ini sama sekali; kolomnya disediakan sekarang karena menambahkannya sesudah data masuk mahal. **Dicabut sebelum data dimuat bila `T-6` menyatakan ia selalu mata uang kontrak** |
| `RETRO_BERGANDA` | teks | tidak |  | korpus `IsMultipleRetro` |
| `ID_DAFTAR_RETRO` | bilangan bulat | ya |  | korpus `RetroList` |
| `SIFAT_MATERIAL_ADDENDUM` | teks | ya |  | korpus `EDMMaterialType` — dua nilai **`MATERIAL`** / **`TIDAK_MATERIAL`**. **Masukan, bukan turunan** (`GRL-20`; `GRL-12` BATAL). Beku sejak `AJUKAN`, berjejak selama `DRAFT` (`KTV-1` Adjustment) |
| `ID_DOKUMEN_ADDENDUM` | bilangan bulat | ya |  | keputusan tiket 14 — **baru**, tidak ada di sistem lama — satu dokumen memayungi banyak versi, **lintas kontrak** (`DB-3`, `DB-4` dibantah). **Persetujuan tetap per versi** |
| `ID_VERSI_KONTRAK_DASAR` | bilangan bulat | ya | FK | keputusan tiket 01 **papan Adjustment** — rujukan ke versi berlaku terakhir saat versi ini dibuat; kosong = versi pertama. Ditambahkan migrasi `440` modul `treatyinadjustment` |

## MATA_UANG_KONTRAK

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_MATA_UANG_KONTRAK` | bilangan bulat | tidak | PK | keputusan tiket 20 — **baru**, tidak ada di sistem lama |
| `ID_VERSI_KONTRAK` | bilangan bulat | tidak | FK UQ | keputusan tiket 20 — **baru**, tidak ada di sistem lama |
| `KODE_MATA_UANG` | bilangan bulat | tidak | FK UQ | korpus `CurrencyID` — **kunci alami**; `Currency` (nama) tidak disimpan — INV-59 |
| `KURS` | angka desimal | tidak |  | korpus `Conversion` — INV-38: selalu lebih besar dari nol |
| `TANGGAL_MULAI_BERLAKU` | DATE | tidak |  | korpus `PeriodStart` — batas inklusif, ADR-0022 |
| `TANGGAL_AKHIR_BERLAKU` | DATE | tidak |  | korpus `PeriodEnd` — batas inklusif |

## RETENSI_CEDANT

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_RETENSI_CEDANT` | bilangan bulat | tidak | PK | keputusan tiket 22 — **baru**, tidak ada di sistem lama |
| `ID_VERSI_KONTRAK` | bilangan bulat | tidak | FK UQ | keputusan tiket 22 — **baru**, tidak ada di sistem lama |
| `ID_KELOMPOK_TREATY` | bilangan bulat | tidak | FK UQ | korpus `TreatyGroupID` — bagian **kunci alami**; `TreatyGroup` (nama) tidak disimpan |
| `NILAI_RETENSI` | angka desimal | tidak |  | korpus `Amount` + `Currency`/`CurrencyID` — bagian **kunci alami** lewat kode mata uangnya. Tingkat **BELUM DITENTUKAN** — §10.23 |
| `KODE_MATA_UANG` | teks | tidak | UQ | korpus `Currency` / `CurrencyID` pada baris yang sama` — **denominasi paket uang di sebelahnya** — bagian kunci alami (INV-08 / INV-09 / INV-12). Ditambahkan 24 Sep 2026, P-8 golongan B |
| `CATATAN` | teks | ya |  | korpus `Note` — **BARU** — ada di kelas, tidak ada di §3.5 |
| `PERSEN_BAGIAN_DIPAKAI` | angka desimal | ya |  | keputusan tiket 22 — **baru**, tidak ada di sistem lama — bila paket uangnya bertingkat `BAGIAN_NURE` (INV-40) |

## EGNPI

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_EGNPI` | bilangan bulat | tidak | PK | keputusan tiket 23 — **baru**, tidak ada di sistem lama |
| `ID_VERSI_KONTRAK` | bilangan bulat | tidak | FK UQ | keputusan tiket 23 — **baru**, tidak ada di sistem lama |
| `ID_KELOMPOK_TREATY` | bilangan bulat | tidak | FK UQ | korpus `TreatyGroupID` — bagian **kunci alami** |
| `ID_KELAS_BISNIS` | bilangan bulat | ya | FK | korpus `ClassOfBusiness` — **BARU** — tidak ada di §3.5 |
| `NILAI_EGNPI` | angka desimal | tidak |  | korpus `Amount` + `Currency`/`CurrencyID` — tingkat **BELUM DITENTUKAN**, dan ini **satu dari lima §7 yang naik ke eskalasi** |
| `KODE_MATA_UANG` | teks | tidak | UQ | korpus `Currency` / `CurrencyID` pada baris yang sama` — **denominasi paket uang di sebelahnya** — bagian kunci alami (INV-08 / INV-09 / INV-12). Ditambahkan 24 Sep 2026, P-8 golongan B |
| `TANGGAL_BERLAKU` | DATE | ya |  | korpus `AsDate` |
| `PROPORSI` | angka desimal | ya |  | korpus `Proportion` — **BARU** — belum terpakai di rumus mana pun; dibawa karena membuangnya menghilangkan fakta yang tidak dapat dipulihkan |
| `CATATAN` | teks | ya |  | korpus `Note` — **BARU** |

## PORTOFOLIO

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_PORTOFOLIO` | bilangan bulat | tidak | PK | keputusan tiket 24 — **baru**, tidak ada di sistem lama |
| `ID_VERSI_KONTRAK` | bilangan bulat | tidak | FK UQ | keputusan tiket 24 — **baru**, tidak ada di sistem lama |
| `ARAH_PORTOFOLIO` | teks | tidak | UQ | korpus `Type` — masuk atau keluar; bagian **kunci alami** |
| `JENIS_PORTOFOLIO` | teks | tidak | UQ | korpus `TypePortfolio` — bagian **kunci alami** |
| `KETERANGAN` | teks | ya |  | korpus `Description` |

## PERIODE_PELAPORAN

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_PERIODE_PELAPORAN` | bilangan bulat | tidak | PK | keputusan tiket 25 — **baru**, tidak ada di sistem lama |
| `ID_VERSI_KONTRAK` | bilangan bulat | tidak | FK UQ | keputusan tiket 25 — **baru**, tidak ada di sistem lama |
| `PERIODE` | teks | tidak | UQ | korpus `Period` — **kunci alami** |
| `TANGGAL_AWAL` | DATE | tidak |  | korpus `InitialDate` |
| `BATAS_PENYERAHAN` | DATE | tidak |  | korpus `SubmissionDue` |
| `BATAS_KONFIRMASI` | DATE | tidak |  | korpus `ConfirmationDue` |
| `BATAS_PELUNASAN` | DATE | tidak |  | korpus `SettlementDue` |

## PERIODE_AKUMULASI

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_PERIODE_AKUMULASI` | bilangan bulat | tidak | PK | keputusan tiket 26 — **baru**, tidak ada di sistem lama |
| `ID_VERSI_KONTRAK` | bilangan bulat | tidak | FK UQ | keputusan tiket 26 — **baru**, tidak ada di sistem lama |
| `PERIODE` | teks | tidak | UQ | korpus `Period` — **kunci alami** |
| `TANGGAL_LAPOR` | DATE | tidak |  | korpus `ReportDate` |
| `HARI_BATAS_PENYERAHAN` | bilangan bulat | ya |  | korpus `SubDays` — dari inventaris kelas, §12.4 |
| `BATAS_PENYERAHAN` | DATE | ya |  | korpus `SubDueDate` — dari inventaris kelas, §12.4 |

## TERMIN

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_TERMIN` | bilangan bulat | tidak | PK | keputusan tiket 27 — **baru**, tidak ada di sistem lama |
| `ID_VERSI_KONTRAK` | bilangan bulat | tidak | FK UQ | keputusan tiket 27 — **baru**, tidak ada di sistem lama |
| `NOMOR_TERMIN` | bilangan bulat | tidak | UQ | korpus `InstallmentList[].Installment` — bagian **kunci alami**, bersama kode mata uang — §10.16a |
| `PERSEN_TERMIN` | angka desimal | tidak |  | korpus `InstallmentPct` |
| `NILAI_TERMIN` | angka desimal | ya |  | korpus `Amount` + `Currency`/`CurrencyID` — tingkat **BELUM DITENTUKAN** — §10.23 |
| `KODE_MATA_UANG` | teks | tidak | UQ | korpus `Currency` / `CurrencyID` pada baris yang sama` — **denominasi paket uang di sebelahnya** — bagian kunci alami (INV-08 / INV-09 / INV-12). Ditambahkan 24 Sep 2026, P-8 golongan B |
| `TANGGAL_JATUH_TEMPO` | DATE | tidak |  | korpus `DueDate` |
| `TANGGAL_BAYAR` | DATE | ya |  | korpus `PaymentDate` — kosong selama belum dibayar |
| `WPC` | teks | ya |  | korpus `WPC` — **singkatan yang kepanjangannya tidak ada di korpus** — lihat di bawah |

## SKALA_KOASURANSI

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_SKALA_KOASURANSI` | bilangan bulat | tidak | PK | keputusan tiket 28 — **baru**, tidak ada di sistem lama |
| `ID_VERSI_KONTRAK` | bilangan bulat | tidak | FK UQ | keputusan tiket 28 — **baru**, tidak ada di sistem lama |
| `PERSEN_LIMIT` | angka desimal | tidak | UQ | korpus `PctLimit` |
| `PERSEN_BAGIAN` | angka desimal | tidak |  | korpus `CoInShare` |

## BATAS_PER_BAHAYA

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_BATAS_PER_BAHAYA` | bilangan bulat | tidak | PK | keputusan tiket 29 — **baru**, tidak ada di sistem lama |
| `ID_VERSI_KONTRAK` | bilangan bulat | tidak | FK UQ | keputusan tiket 29 — **baru**, tidak ada di sistem lama |
| `ID_BAHAYA` | bilangan bulat | tidak | FK UQ | korpus `nama kolom lama` — **kunci alami**; tabel acuan, ADR-0038 |
| `NILAI_BATAS` | angka desimal | tidak |  | korpus `Earthquake` / `FloodJab` / `FloodNation` / `RSMDLimit` + mata uangnya` — tingkat **BELUM DITENTUKAN** — §10.23 |
| `KODE_MATA_UANG` | teks | ya |  | korpus `Currency` / `CurrencyID` pada baris yang sama` — **denominasi paket uang di sebelahnya**. Golongan **B**, dipindahkan dari golongan C oleh `KTV-C` — §10.18 menulis asalnya *"+ mata uangnya"*. **Boleh kosong**: ia bukan bagian kunci alami (INV-14 memakai `ID_BAHAYA`) |
| `PERSEN_BAGIAN_DIPAKAI` | angka desimal | ya |  | keputusan tiket 29 — **baru**, tidak ada di sistem lama — bila bertingkat `BAGIAN_NURE` (INV-40) |

## DOKUMEN_KONTRAK

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_DOKUMEN_KONTRAK` | bilangan bulat | tidak | PK | keputusan tiket 30 — **baru**, tidak ada di sistem lama |
| `ID_VERSI_KONTRAK` | bilangan bulat | tidak | FK UQ | keputusan tiket 30 — **baru**, tidak ada di sistem lama |
| `ID_DOKUMEN` | bilangan bulat | tidak | UQ | korpus `lampiran` — **kunci alami**; rujukan ke luar skema |
| `JENIS_DOKUMEN` | teks | ya |  | korpus `lampiran` |
| `TANGGAL_LAMPIR` | DATE | tidak |  | korpus `lampiran` |

## LAYER

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_LAYER` | bilangan bulat | tidak | PK | korpus `Limits[].ID` |
| `ID_VERSI_KONTRAK` | bilangan bulat | tidak | FK UQ | keputusan tiket 31 — **baru**, tidak ada di sistem lama |
| `NOMOR_LAYER` | bilangan bulat | tidak | UQ | korpus `Layer` — **kunci alami** di dalam versi |
| `BAGIAN_LAYER` | bilangan bulat | ya | UQ | korpus `LayerPart` — bagian dari kunci alami |
| `JENIS_LAYER` | teks | tidak |  | korpus `LayerType` |
| `JENIS_BAGIAN_LAYER` | teks | ya |  | korpus `LayerPartType` |
| `CAKUPAN` | teks | ya |  | korpus `Cover` |
| `LIMIT` | angka desimal | tidak |  | korpus `Limit` + `Currency` — tingkat **100% treaty** |
| `KODE_MATA_UANG` | teks | tidak |  | korpus `Currency` / `CurrencyID` pada baris yang sama` — **denominasi paket uang di sebelahnya** — bagian kunci alami (INV-08 / INV-09 / INV-12). Ditambahkan 24 Sep 2026, P-8 golongan B |
| `DEDUCTIBLE` | angka desimal | tidak |  | korpus `Deductible` — tingkat **100% treaty** |
| `MATA_UANG_DEDUCTIBLE` | teks | ya |  | korpus `**tidak ada di sistem lama**` — **denominasi `DEDUCTIBLE`** — `KTV-C`. Sistem lama tidak merekam mata uang untuk besaran ini sama sekali; kolomnya disediakan sekarang karena menambahkannya sesudah data masuk mahal. **Dicabut sebelum data dimuat bila `T-6` menyatakan ia selalu mata uang kontrak** |
| `PERSEN_PENYESUAIAN` | angka desimal | ya |  | korpus `AdjRate` |
| `PERSEN_MINIMUM_DEPOSIT` | angka desimal | ya |  | korpus `MDPPct` |
| `PORSI_PEMULIHAN_LIMIT` | angka desimal | ya |  | korpus `ReinstatementPct` — porsi limit yang dipulihkan |
| `TARIF_PREMI_PEMULIHAN` | angka desimal | ya |  | korpus `ReinstatementValue` — tarif premi untuk pemulihan itu — **pembacaannya punya saingan**, §10.3b |
| `LIMIT_AGREGAT` | angka desimal | ya |  | korpus `AgregateLimit` *(salah eja ada di sumbernya)` — **BARU** — batas total sepanjang periode, terpisah dari limit per kejadian. Tingkat **BELUM DITENTUKAN** |
| `MATA_UANG_LIMIT_AGREGAT` | teks | ya |  | korpus `**tidak ada di sistem lama**` — **denominasi `LIMIT_AGREGAT`** — `KTV-C`. Sistem lama tidak merekam mata uang untuk besaran ini sama sekali; kolomnya disediakan sekarang karena menambahkannya sesudah data masuk mahal. **Dicabut sebelum data dimuat bila `T-6` menyatakan ia selalu mata uang kontrak** |
| `MDP` | angka desimal | ya |  | korpus `MDPList[]` — deposit premium; tingkat **BELUM DITENTUKAN** |
| `MDP_MINIMUM` | angka desimal | ya |  | korpus `MDPMinList[]` — **BARU** — tanpa ia, *"deposit premium sama dengan minimum premium"* tidak dapat dinyatakan |
| `PERSEN_MDP_MINIMUM` | angka desimal | ya |  | korpus `MDPMinPct` — **BARU** |
| `MDP_DIGABUNG` | teks | tidak |  | korpus `IsCombineMDP` — **BARU** — MDP dihitung per layer atau gabungan |
| `TANPA_HITUNG_PREMI_PEMULIHAN` | teks | tidak |  | korpus `NoRIPCalculation` — **BARU** — mematikan perhitungan premi pemulihan pada layer ini |
| `DEDUCTIBLE_KEDUA` | angka desimal | ya |  | korpus `Limits[].Deductible2` — **PULIH dari `L-8`** — `TDA-17`. Tersimpan sebagai **teks** di sistem lama (`@toDecimal` di kedua sisi rumus selisih); migrasi wajib mengubah tipe, dan baris yang gagal konversi **tidak disamarkan jadi nol** (ADR-0035) |
| `PERSEN_ROL` | angka desimal | ya |  | korpus `Limits[].ROLPct` — **PULIH dari `L-8`** — `TDA-17`. Rate on line. ⚠ **DIBANTAH `F-15`** — pertanyaan *"selalu dihitung atau pernah disepakati"* **terjawab dari sumber**: `DetailCalculationROL` menghitungnya `premi ÷ limit × 100`. Kolomnya **tidak dicabut diam-diam**; keputusannya pemilik proses, `2-to-spec/F-15-PREMI-DIPEROLEH-DAN-ROL-TURUNAN.md` |

## NILAI_MDP

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_NILAI_MDP` | bilangan bulat | tidak | PK | keputusan tiket 31 — **baru**, tidak ada di sistem lama |
| `ID_LAYER` | bilangan bulat | tidak | FK | keputusan tiket 31 — **baru**, tidak ada di sistem lama — induknya — paket uang ini asalnya **daftar per mata uang** di sistem lama |
| `KODE_MATA_UANG` | teks | tidak |  | korpus `Currency` pada baris daftarnya` — **kunci alami** di dalam induknya — belum bernomor |
| `NILAI` | angka desimal | tidak |  | korpus `Value` pada baris daftarnya` |

## NILAI_MDP_MINIMUM

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_NILAI_MDP_MINIMUM` | bilangan bulat | tidak | PK | keputusan tiket 31 — **baru**, tidak ada di sistem lama |
| `ID_LAYER` | bilangan bulat | tidak | FK | keputusan tiket 31 — **baru**, tidak ada di sistem lama — induknya — paket uang ini asalnya **daftar per mata uang** di sistem lama |
| `KODE_MATA_UANG` | teks | tidak |  | korpus `Currency` pada baris daftarnya` — **kunci alami** di dalam induknya — belum bernomor |
| `NILAI` | angka desimal | tidak |  | korpus `Value` pada baris daftarnya` |

## PEMULIHAN_LIMIT

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_PEMULIHAN_LIMIT` | bilangan bulat | tidak | PK | keputusan tiket 31 — **baru**, tidak ada di sistem lama |
| `ID_LAYER` | bilangan bulat | tidak | FK | keputusan tiket 31 — **baru**, tidak ada di sistem lama |
| `NOMOR_URUT_PEMULIHAN` | bilangan bulat | tidak |  | korpus `Reinstatement_List[].ReinstatementValue` |
| `PERSEN_PEMULIHAN` | angka desimal | tidak |  | korpus `Reinstatement_List[].ReinstatementPct` |
| `PERSEN_TAMBAHAN` | angka desimal | ya |  | korpus `Reinstatement_List[].AdditionalPct` |
| `CATATAN` | teks | ya |  | korpus `Reinstatement_List[].ReinstatementNote` |

## JEJAK_PERUBAHAN

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_JEJAK_PERUBAHAN` | bilangan bulat | tidak | PK | keputusan tiket 39 — **baru**, tidak ada di sistem lama |
| `ID_VERSI_KONTRAK` | bilangan bulat | tidak | FK | keputusan tiket 39 — **baru**, tidak ada di sistem lama |
| `WAKTU_PERUBAHAN` | DATE | tidak |  | keputusan tiket 39 — **baru**, tidak ada di sistem lama — **`DATE` beresolusi detik.** Dua perubahan dalam detik yang sama tidak dapat dipisahkan olehnya — dan itu tidak merusak apa pun di sini, sebab entitas ini **sengaja tanpa kunci alami** |
| `PELAKU` | teks | tidak |  | keputusan tiket 39 — **baru**, tidak ada di sistem lama — dibekukan, sekeluarga dengan `NAMA_PEMUTUS` |
| `PERAN_PELAKU` | teks | tidak |  | keputusan tiket 39 — **baru**, tidak ada di sistem lama — **peran yang berlaku SAAT ITU** — ADR-0045 isi minimal butir kelima. **POTRET, bukan rujukan** (`KTV-D`): teks, final, **tidak dinormalisasi ulang** ketika entitas peran kelak lahir. Dari mana nilainya diambil masih `F-16` |
| `RUAS_YANG_BERUBAH` | teks | tidak |  | keputusan tiket 39 — **baru**, tidak ada di sistem lama |
| `NILAI_SEBELUM` | teks | ya |  | keputusan tiket 39 — **baru**, tidak ada di sistem lama — teks, karena ruasnya beragam tipe |
| `NILAI_SESUDAH` | teks | ya |  | keputusan tiket 39 — **baru**, tidak ada di sistem lama — idem |

## BAGIAN

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_BAGIAN` | bilangan bulat | tidak | PK | keputusan tiket 33 — **baru**, tidak ada di sistem lama |
| `ID_LAYER` | bilangan bulat | tidak | FK UQ | keputusan tiket 33 — **baru**, tidak ada di sistem lama — **relasi**, bukan salinan — §12.2 |
| `CAKUPAN` | teks | ya |  | korpus `Cover` — tetap atribut bagian (§12.2) |
| `PERSEN_BAGIAN_NURE` | angka desimal | ya |  | korpus `Share[].RNMShare` — terisi bila `BAGIAN_NURE_SERAGAM` mati; berbeda dari atribut senama pada versi (§10.2) yang berlaku seragam |
| `PREMI_BRUTO` | angka desimal | tidak |  | korpus `GrossPremiumList[]` — tingkat **BELUM DITENTUKAN** — §10.23 |
| `PREMI_BRUTO_MINIMUM` | angka desimal | ya |  | korpus `GrossPremiumMinList[]` — **BARU** — tidak pernah terlihat pohon |
| `ID_SUSUNAN_RETRO` | bilangan bulat | ya |  | korpus `SpreadingTypeIDXOL` — penunjuk susunan baku yang menyemai penyebaran — **syarat INV-58**, §10.4a |
| `PERSEN_BAGIAN_DIPAKAI` | angka desimal | ya |  | keputusan tiket 33 — **baru**, tidak ada di sistem lama — wajib bila paket uangnya bertingkat `BAGIAN_NURE` (INV-40) |

## NILAI_PREMI_BRUTO

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_NILAI_PREMI_BRUTO` | bilangan bulat | tidak | PK | keputusan tiket 33 — **baru**, tidak ada di sistem lama |
| `ID_BAGIAN` | bilangan bulat | tidak | FK | keputusan tiket 33 — **baru**, tidak ada di sistem lama — induknya — paket uang ini asalnya **daftar per mata uang** di sistem lama |
| `KODE_MATA_UANG` | teks | tidak |  | korpus `Currency` pada baris daftarnya` — **kunci alami** di dalam induknya — belum bernomor |
| `NILAI` | angka desimal | tidak |  | korpus `Value` pada baris daftarnya` |

## NILAI_PREMI_BRUTO_MINIMUM

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_NILAI_PREMI_BRUTO_MINIMUM` | bilangan bulat | tidak | PK | keputusan tiket 33 — **baru**, tidak ada di sistem lama |
| `ID_BAGIAN` | bilangan bulat | tidak | FK | keputusan tiket 33 — **baru**, tidak ada di sistem lama — induknya — paket uang ini asalnya **daftar per mata uang** di sistem lama |
| `KODE_MATA_UANG` | teks | tidak |  | korpus `Currency` pada baris daftarnya` — **kunci alami** di dalam induknya — belum bernomor |
| `NILAI` | angka desimal | tidak |  | korpus `Value` pada baris daftarnya` |

## DETAIL_PROPORSIONAL

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_DETAIL_PROPORSIONAL` | bilangan bulat | tidak | PK | keputusan tiket 34 — **baru**, tidak ada di sistem lama |
| `ID_LAYER` | bilangan bulat | tidak | FK UQ | keputusan tiket 34 — **baru**, tidak ada di sistem lama |
| `ID_KELOMPOK_TREATY` | bilangan bulat | tidak | FK UQ | korpus `TreatyGroupID` — **kunci alami** di dalam **`LAYER`** — INV-06 |
| `JENIS_TREATY` | teks | tidak |  | korpus `TreatyType` — `QUOTA_SHARE` / `SURPLUS` |
| `PERSEN_QUOTA_SHARE` | angka desimal | ya |  | korpus `QSPct` — terisi hanya bila `JENIS_TREATY = QUOTA_SHARE` |
| `JUMLAH_LINES_SURPLUS` | bilangan bulat | ya |  | korpus `Surplus` — terisi hanya bila `JENIS_TREATY = SURPLUS` |
| `PERSEN_KOMISI_KOTOR` | angka desimal | ya |  | korpus `RIOGR` — kepanjangan **belum diketahui** |
| `PERSEN_KOMISI_BERSIH` | angka desimal | ya |  | korpus `RIONR` — kepanjangan **belum diketahui** |
| `PERSEN_CADANGAN_PREMI` | angka desimal | ya |  | korpus `PremiumReservePct` — dari inventaris kelas Pega |
| `CADANGAN_PREMI` | angka desimal | ya |  | korpus `ReserveList[]` — **BARU** — pasangan nilai dari persentase di atas; persentase tanpa nilainya adalah setengah fakta. Tingkat **BELUM DITENTUKAN** |
| `PERSEN_KAPASITAS_SURPLUS` | angka desimal | ya |  | korpus `IOOPct` — **BARU** |
| `ID_SUSUNAN_RETRO` | bilangan bulat | ya |  | korpus `SpreadingTypeID` — **BARU — penunjuk susunan baku yang menyemai penyebaran.** Syarat berdirinya `RINCIAN_PENYEBARAN` sebagai fakta terbukukan (INV-58); lihat §10.4a |

## NILAI_CADANGAN_PREMI

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_NILAI_CADANGAN_PREMI` | bilangan bulat | tidak | PK | keputusan tiket 34 — **baru**, tidak ada di sistem lama |
| `ID_DETAIL_PROPORSIONAL` | bilangan bulat | tidak | FK | keputusan tiket 34 — **baru**, tidak ada di sistem lama — induknya — paket uang ini asalnya **daftar per mata uang** di sistem lama |
| `KODE_MATA_UANG` | teks | tidak |  | korpus `Currency` pada baris daftarnya` — **kunci alami** di dalam induknya — belum bernomor |
| `NILAI` | angka desimal | tidak |  | korpus `Value` pada baris daftarnya` |

## POTONGAN

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_POTONGAN` | bilangan bulat | tidak | PK | keputusan tiket 37 — **baru**, tidak ada di sistem lama — §14.2 |
| `ID_BAGIAN` | bilangan bulat | ya | FK UQ | keputusan tiket 37 — **baru**, tidak ada di sistem lama — `KTV-B` — salah satu dari **dua** pelekatan. **Tepat satu** dari kedua kolom terisi, dijaga `CHECK`. Menggantikan `ID_INDUK_POTONGAN`, yang tidak dapat punya kunci asing karena sasarannya bergantung nilai kolom lain (INV-17) |
| `ID_DETAIL_PROPORSIONAL` | bilangan bulat | ya | FK UQ | keputusan tiket 37 — **baru**, tidak ada di sistem lama — `KTV-B` — salah satu dari **dua** pelekatan. **Tepat satu** dari kedua kolom terisi, dijaga `CHECK`. Menggantikan `ID_INDUK_POTONGAN`, yang tidak dapat punya kunci asing karena sasarannya bergantung nilai kolom lain (INV-17) |
| `ID_JENIS_POTONGAN` | bilangan bulat | tidak | FK UQ | korpus `Comment` — §14.2 |
| `DASAR_PERHITUNGAN` | teks | tidak |  | keputusan tiket 37 — **baru**, tidak ada di sistem lama — §14.2 |
| `PERSEN_POTONGAN` | angka desimal | tidak |  | korpus `DeductionPct` — §14.2 |
