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

> ⛔ **RALAT 4 Oktober 2026 — `MATA_UANG` dan `MATA_UANG_KONTRAK` DICABUT**
> (migrasi `434`, `KEPUTUSAN-PENYELARASAN-REPO.md` §16). Kurs dan daftar mata
> uang dibaca dari `TREATYEXCHANGEYEARLY`. Himpunan acuan kini **lima**, bukan
> enam, dan tiket `20` kehilangan tabelnya — lihat §16 untuk akibatnya pada
> tiket `57` dan `INV-44`.

| | Tabel | Kolom | Tiket |
| --- | --- | ---: | --- |
| Tabel acuan | `JENIS_POTONGAN`, `KELAS_BISNIS`, `KELOMPOK_TREATY`, `BAHAYA`, `JENIS_REASURANSI` | 25 | `15` |
| Identitas kontrak | `KONTRAK`, `VERSI_KONTRAK` | 57 | `14` · `01` |
| Uang dan periode | `RETENSI_CEDANT`, `EGNPI`, `TERMIN` | 31 | `20` `22` `23` `27` |
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

**Sequence:** `SEQ_TRIN_KONTRAK`, `SEQ_TRIN_VERSI_KONTRAK`,
`SEQ_TRIN_JENIS_POTONGAN`, `SEQ_TRIN_JENIS_REASURANSI`, `SEQ_TRIN_BAHAYA`,
`SEQ_TRIN_KELOMPOK_TREATY`, `SEQ_TRIN_KELAS_BISNIS` (berkas `402`), ditambah satu per tabel anak di
berkas `415`: `SEQ_TRIN_RETENSI_CEDANT`, `SEQ_TRIN_EGNPI`,
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
`IX_JEJAK_PERUBAHAN_VERSI` · `IX_RETENSI_CEDANT_KLP` ·
`IX_EGNPI_KELOMPOK` · `IX_EGNPI_KELAS_BISNIS` · `IX_BATAS_PER_BAHAYA_BHY` · `IX_DETAIL_PROP_KLP` ·
`IX_POTONGAN_JENIS` · `IX_NILAI_PB_BAGIAN` · `IX_NILAI_PB_MIN_BAGIAN` · `IX_NILAI_CAD_PREMI_DP`.

ℹ️ **Kenapa sebagian kunci asing ber-index sendiri dan sebagian tidak.** Oracle membuat index untuk
setiap `UNIQUE`, dan index itu **melayani** sebuah kunci asing bila kolom kunci asingnya **memimpin**
kunci alami tersebut. Sebelas tabel anak berkunci asing ke `VERSI_KONTRAK`, dan pada sebelas-belasnya
`ID_VERSI_KONTRAK` adalah kolom **pertama** kunci alaminya — index itu sudah melayaninya.

Yang **tidak** demikian ada dua rupa, dan keduanya ber-index eksplisit:

| Rupa | Index |
| --- | --- |
| kunci asing ke **tabel acuan** — tidak pernah memimpin kunci alami | `IX_RETENSI_CEDANT_KLP`, `IX_EGNPI_KELOMPOK`, `IX_EGNPI_KELAS_BISNIS`, `IX_BATAS_PER_BAHAYA_BHY` |
| anak yang **tidak punya kunci alami** sama sekali | `IX_NILAI_MDP_LAYER`, `IX_NILAI_MDP_MIN_LAYER`, `IX_PEMULIHAN_LIMIT_LAYER`, `IX_JEJAK_PERUBAHAN_VERSI`, `IX_VERSI_KONTRAK_KONTRAK` |

⚠️ **Kelima index tabel acuan tidak ada di `ddl-usulan/`.** Ia ditambahkan di sini sebab kunci asing
tanpa index membuat penghapusan baris acuan — satu `BAHAYA` — memindai seluruh tabel
anak, dan pada sebagian versi Oracle menguncinya. Dicatat di
[`KEPUTUSAN-PENYELARASAN-REPO.md`](KEPUTUSAN-PENYELARASAN-REPO.md) butir 5. Dijaga
`TestSetiapKunciAsingTerlayaniIndex`.

## Relasi

| Anak | Kolom | Induk | ON DELETE | Invarian |
| --- | --- | --- | --- | --- |
| `VERSI_KONTRAK` | `ID_KONTRAK` | `KONTRAK` | **tanpa** | INV-18 |
| `VERSI_KONTRAK` | `ID_VERSI_KONTRAK_DASAR` | `VERSI_KONTRAK` | **tanpa** | INV-18 |
| `JENIS_REASURANSI` | `ID_INDUK` | `JENIS_REASURANSI` | **tanpa** | INV-18 |
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

✅ **`UQ_LAYER` MENEGAKKAN INV-05 — terverifikasi di Oracle 2 Oktober 2026.**

> ⚠️ **RALAT.** Bagian ini pernah menyatakan `UQ_LAYER` **tidak** menegakkan INV-05 sepenuhnya,
> sebab `BAGIAN_LAYER` boleh kosong dan Oracle memperlakukan NULL sebagai tidak sama dengan NULL.
> **Itu keliru.** Aturan Oracle yang sebenarnya: sebuah entri dilewati indeks unik hanya bila
> **SELURUH** kolom kuncinya NULL. Di sini `ID_VERSI_KONTRAK` dan `NOMOR_LAYER` selalu terisi, jadi
> barisnya terindeks dan duplikatnya ditolak — `ORA-00001: unique constraint (UQ_LAYER) violated`.
> Yang memperbaiki kekeliruan ini bukan argumen, melainkan satu `INSERT`.

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

## KELAS_BISNIS_LAYER

Tiket `65`. Struktur dari `ERD-TREATY-IN-DAN-EDM.html` baris relasi **4** — `T_TREATY_LIMIT_COB`,
induk `T_TREATY_LIMIT_DETAIL` lewat `LIMIT_DETAIL_ID`, `CASCADE`. Jalur Pega
`TreatyIn.Limits.Detail.COBList`, tingkat bukti **DAUN-RELATIF**.

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_KELAS_BISNIS_LAYER` | bilangan bulat | tidak | PK | keputusan tiket 65 — **baru**, tidak ada di sistem lama |
| `ID_DETAIL_PROPORSIONAL` | bilangan bulat | tidak | FK | ERD baris 4 `LIMIT_DETAIL_ID` — induknya — **kunci alami** di dalamnya, **belum bernomor** |
| `ID_KELAS_BISNIS` | bilangan bulat | tidak | FK | korpus `COBList[].ClassOfBusinessID` — **kunci alami** di dalam induknya, **belum bernomor**. ⛔ `ClassOfBusiness` (nama) **TIDAK** ikut disalin — INV-59; namanya dibaca lewat join |

## KELOMPOK_LAYER

Tiket `67`. Struktur dari ERD baris relasi **9** — `T_TREATY_LIMIT_GROUP`, induk `T_TREATY_LIMITS`
lewat `LIMIT_ID`, `CASCADE`. Jalur Pega `TreatyIn.Limits.TreatyGroupList`, bukti **DAUN-RELATIF**.

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_KELOMPOK_LAYER` | bilangan bulat | tidak | PK | keputusan tiket 67 — **baru**, tidak ada di sistem lama |
| `ID_LAYER` | bilangan bulat | tidak | FK | ERD baris 9 `LIMIT_ID` — induknya — **kunci alami** di dalamnya, **belum bernomor** |
| `ID_KELOMPOK_TREATY` | bilangan bulat | tidak | FK | korpus `TreatyGroupList[].TreatyGroupID` — **kunci alami** di dalam induknya, **belum bernomor**. ⛔ `TreatyGroup` (nama) **TIDAK** ikut disalin — INV-59 |

## KELAS_BISNIS_KELOMPOK

Tiket `68`. Struktur dari ERD baris relasi **10** — `T_TREATY_LIMIT_GROUP_COB`, induk
`T_TREATY_LIMIT_GROUP` lewat `LIMIT_GROUP_ID`, `CASCADE`. Jalur Pega
`TreatyIn.Limits.TreatyGroupList.ClassOfBusinessList` — **kedalaman keempat**, bukti
**DAUN-RELATIF**.

⚠️ **Ia BUKAN kembaran `KELAS_BISNIS_LAYER`.** Yang satu menjawab *"kelas apa saja yang ditanggung
rincian proporsional ini"*, yang lain *"kelas apa saja yang ditanggung kelompok treaty X **di
dalam** layer ini"*. ERD menempatkan keduanya pada kedalaman yang berbeda, dan menggabungkannya
melahirkan satu kolom induk yang separuh waktu kosong — bentuk yang `ADR-0041` tolak.

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_KELAS_BISNIS_KELOMPOK` | bilangan bulat | tidak | PK | keputusan tiket 68 — **baru**, tidak ada di sistem lama |
| `ID_KELOMPOK_LAYER` | bilangan bulat | tidak | FK | ERD baris 10 `LIMIT_GROUP_ID` — induknya — **kunci alami** di dalamnya, **belum bernomor** |
| `ID_KELAS_BISNIS` | bilangan bulat | tidak | FK | korpus `ClassOfBusinessList[]` — **kunci alami** di dalam induknya, **belum bernomor** |

## PENCAPAIAN

Tiket `69`. **Induknya `KONTRAK`, dan ERD menulis sebaliknya** — uraian keputusannya di kepala
migrasi `423_pencapaian.sql`. Ringkasnya: ERD baris **6** menempatkannya di bawah
`T_TREATY_LIMIT_DETAIL` dengan bukti yang ia sendiri tandai **DAUN-RELATIF** (*"lemah"*),
sementara tabel Oracle yang sungguh hidup berkunci tingkat kontrak —
`RDBList/GetAchievement.xml` menyaring `SUBSTR(NOOFFER,1,7)`, dan nol kolomnya menunjuk rincian
layer. `ERD.md` §2.9 menulis `KONTRAK`, dan bukti SQL sejalan dengannya.

Perilaku hapus **tolak** (`ERD.md` §2.9) — pencapaian adalah angka yang pernah dibukukan.

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_PENCAPAIAN` | bilangan bulat | tidak | PK | keputusan tiket 69 — **baru**, tidak ada di sistem lama |
| `ID_KONTRAK` | bilangan bulat | tidak | FK | `ERD.md` §2.9 — induknya. Menggantikan `NOOFFER`/`NOPOLIS` yang di sistem lama berupa teks |
| `TRIWULAN` | bilangan bulat | tidak |  | korpus `ACHIEVEMENT.QUARTER` — **kunci alami** di dalam induknya, **belum bernomor** |
| `TAHUN_TRIWULAN` | bilangan bulat | tidak |  | korpus `ACHIEVEMENT.QUARTERYEAR` — ikut **kunci alami** |
| `KODE_MATA_UANG` | teks | tidak |  | korpus `ACHIEVEMENT.CURRENCY` — ikut **kunci alami** |
| `PREMI` | angka desimal | ya |  | korpus `ACHIEVEMENT.PREMIUM` |
| `KOMISI_REASURANSI` | angka desimal | ya |  | korpus `ACHIEVEMENT.RICOMM` |
| `BROKERAGE` | angka desimal | ya |  | korpus `ACHIEVEMENT.BROKERAGE` |
| `PREMI_BERSIH` | angka desimal | ya |  | korpus `ACHIEVEMENT.NETPREMIUM` |
| `KLAIM_DIBAYAR` | angka desimal | ya |  | korpus `ACHIEVEMENT.PAIDCLAIM` |
| `CASH_CALL_KLAIM` | angka desimal | ya |  | korpus `LOG_ACHIEVEMENT.CASHCALLCLAIM` — ⛔ **hanya ada di tabel log**, dan karena itu hilang bila hanya `ACHIEVEMENT` yang dipindahkan |
| `KLAIM_OUTSTANDING` | angka desimal | ya |  | korpus `ACHIEVEMENT.OUTSTANDINGCLAIM` |

⛔ **Empat kolom sistem lama yang TIDAK dibawa**, masing-masing dengan sebabnya:
`INCUREDCLAIM`, `TOTAL`, `LOSSRATIO` — **turunan**, `INV-58`, dihitung saat dibaca.
`SOBNAME`, `TREATYGROUPNAME`, `TREATYTYPE` — **salinan atribut kontrak dan layer**, `INV-59`;
dibaca lewat kunci asingnya. `PXCREATEOPNAME` dan `INSERTDATE` pindah ke `JEJAK_PERUBAHAN`
(tiket `39`), yang sudah ada.

## RINCIAN_ANGSURAN

Tiket `70`. Struktur dari ERD baris relasi **27** — `T_TREATY_INSTALLMENT_ITEM`, induk
`T_TREATY_INSTALLMENT` lewat `INSTALLMENT_ID`, `CASCADE`. Jalur Pega
`TreatyIn.ValueDifference.Installment.InstallmentList` — ⚠️ **hanya terbaca lewat salinan**, bukti
**DAUN-RELATIF**.

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_RINCIAN_ANGSURAN` | bilangan bulat | tidak | PK | keputusan tiket 70 — **baru**, tidak ada di sistem lama |
| `ID_TERMIN` | bilangan bulat | tidak | FK | ERD baris 27 `INSTALLMENT_ID` — induknya — **kunci alami** di dalamnya, **belum bernomor** |
| `NOMOR_URUT_RINCIAN` | bilangan bulat | tidak |  | keputusan tiket 70 — **baru** — ikut **kunci alami**, **belum bernomor** |
| `TANGGAL_JATUH_TEMPO` | DATE | ya |  | korpus `InstallmentList[].DueDate` |
| `PERSEN_ANGSURAN` | angka desimal | ya |  | korpus `InstallmentList[].InstallmentPct` — INV-41 ditegakkan services, bukan `CHECK` (ADR-0056) |
| `TANGGAL_BAYAR` | DATE | ya |  | korpus `InstallmentList[].PaymentDate` — kosong berarti **belum dibayar**, keadaan normal |
| `WPC` | teks | ya |  | korpus `InstallmentList[].WPC` — ⚠️ **nama apa adanya**; kepanjangannya tidak ada di satu berkas pun, korpus sudah disapu habis (`CONTEXT.md` §2.10). Menamainya dengan tebakan lebih buruk daripada menyimpannya begini |

## ARSIP_MUATAN_KELUAR

Tiket `74`, prasyarat tiket `42`. Struktur dari ERD baris relasi **39** —
`T_TREATY_OUTBOUND_ARCHIVE`, induk `TREATY_IN` lewat `TREATY_IN_ID`. ⚠️ ERD menulis kolom
`ON DELETE`-nya **"di Go"**: ia tidak meresepkan aturan hapus tingkat basis data. Di sini
diwujudkan **tolak** — arsip yang lenyap bersama kontraknya berhenti menjadi arsip.

⛔ **`INV-61` berlaku keras**: kolom `MUATAN` **tidak punya jalur baca aplikasi**. Basis data tidak
dapat menegakkannya; yang menjaganya tinjauan kode dan uji.

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_ARSIP_MUATAN_KELUAR` | bilangan bulat | tidak | PK | keputusan tiket 74 — **baru**, tidak ada di sistem lama |
| `ID_KONTRAK` | bilangan bulat | tidak | FK | ERD baris 39 `TREATY_IN_ID` — induknya |
| `TUJUAN` | teks | tidak |  | korpus — salah satu dari `PEGA_TREATY_IN`, `PEGA_M_TREATY_IN_EDM`, `PEGA_M_TREATY_IN_DETAIL`, `PEGA_M_TREATY_IN_DETAIL_EDM` |
| `DIKIRIM_PADA` | DATE | tidak |  | keputusan tiket 74 — **baru** |
| `BERHASIL` | teks | tidak |  | keputusan tiket 74 — **baru** — pengiriman yang **gagal di hilir tetap terarsip**; arsip yang hanya memuat yang berhasil tidak menyelesaikan satu pun perselisihan |
| `MUATAN` | teks | tidak |  | korpus — 175 argumen, apa adanya. ⚠️ `VARCHAR2(4000 CHAR)`, bukan CLOB: muatan yang melampauinya **DITOLAK** basis data, dan penolakan lebih baik daripada pemotongan diam-diam — arsip yang terpotong adalah arsip yang berbohong |

## CATATAN_PERSETUJUAN

Tiket `54`, dan ia melepas `49` `55` `56`. Dari `ddl-usulan/13_CATATAN_PERSETUJUAN.sql` dan
`SPEC-MODEL-DATA.md` §10.20 — kelas lama `Data-SuggestList`.

⛔ **Perilaku hapus `tolak`**, dan ia **satu-satunya** anak `VERSI_KONTRAK` di `ERD.md` §2.3 yang
begitu; sepuluh lainnya `ikut hapus`. Sebabnya dikutip utuh dari §2.3: *"jejak yang dapat dihapus
bersama bendanya bukan jejak."*

⛔ **Nol kunci alami, dan itu KEPUTUSAN — bukan tagihan.** Dua keputusan pada versi yang sama, oleh
orang yang sama, pada hari yang sama adalah keadaan yang sah; yang membedakan barisnya urutan
waktu, bukan sebuah nilai. Ia karena itu **tidak** ikut daftar tagihan `SPEC-MODEL-DATA.md` §13.

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_CATATAN_PERSETUJUAN` | bilangan bulat | tidak | PK | keputusan tiket 54 — **baru**, tidak ada di sistem lama |
| `ID_VERSI_KONTRAK` | bilangan bulat | tidak | FK | menggantung pada **versi**, bukan kontrak — §10.20 |
| `WAKTU_KEPUTUSAN` | DATE | tidak |  | korpus `Date` |
| `NAMA_PEMUTUS` | teks | tidak |  | korpus `OperatorName` — ⚠️ **teks, bukan kunci asing**: fakta historis siapa memutuskan apa dan kapan (§12.5, `ADR-0045`). Sekeluarga dengan `PELAKU` di `JEJAK_PERUBAHAN` |
| `DISETUJUI` | teks | tidak |  | korpus `IsApproved` — himpunan tertutup, ditegakkan services (`ADR-0056`), bukan `CHECK` |
| `ALASAN` | teks | ya |  | korpus `Suggest` |

⛔ **Baris PERISTIWA tidak masuk ke sini.** §10.20a mencabut usul nilai enumerasi ketiga
`PERISTIWA`: ia akan membuat separuh kolom kosong pada separuh baris dan mencemari `INV-28`, yang
menghitung satu baris per perpindahan. Rumahnya `PERISTIWA_KONTRAK` — yang sampai hari ini masih
yatim, nol tiket menghasilkannya.

## MIGRASI_KORELASI · MIGRASI_PENDARATAN · MIGRASI_NILAI_DITOLAK

Tiket `71` `72` `73`, migrasi `428`. ⛔ **BERTENGGAT terhadap tiket `44`** — tiket itu memuat data
pertama, dan sesudahnya ketiganya kehilangan sebagian gunanya selamanya.

⚠️ **Ketiganya tidak punya kotak di ERD, dan itu wajar.** ERD menyatakan dirinya *"POTRET SISTEM
LAMA, BUKAN RANCANGAN"*; ketiganya **perkakas sistem baru**. ERD tetap menyebutnya di daftar
*"7 tabel tanpa jalur Pega … tidak digambar, dilaporkan di sini"*, dan satu punya baris relasi
(baris 38, `LANDING_ID`, ON DELETE **"di Go"**).

### MIGRASI_KORELASI — jembatan pengenal, **nol kunci asing**

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_MIGRASI_KORELASI` | bilangan bulat | tidak | PK | keputusan tiket 71 — **baru** |
| `KUNCI_PEGA` | teks | ya |  | `pzInsKey` |
| `ID_PEGA` | teks | ya |  | `TreatyIn.ID` |
| `ID_WARISAN` | teks | ya |  | `TreatyIn.OLDID` — jejak pemindahan **sebelumnya**, dipindahkan sebagai isi |
| `ID_KONTRAK_BARU` | bilangan bulat | ya |  | ⚠️ **nilai, BUKAN kunci asing** — lihat di bawah |
| `DIPINDAHKAN_PADA` | DATE | tidak |  | keputusan tiket 71 — **baru** |

### MIGRASI_PENDARATAN — bentuk lama apa adanya, **nol kunci asing**

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_MIGRASI_PENDARATAN` | bilangan bulat | tidak | PK | keputusan tiket 72 — **baru** |
| `KUNCI_WARISAN` | teks | tidak |  | pengenal dokumen warisan |
| `MUATAN` | teks | tidak |  | `M_TREATY_IN.JSONDATA` apa adanya — ⛔ `INV-61`: **nol jalur baca aplikasi** |
| `MENDARAT_PADA` | DATE | tidak |  | keputusan tiket 72 — **baru** |

### MIGRASI_NILAI_DITOLAK — anak pendaratan, perilaku hapus **tolak**

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_MIGRASI_NILAI_DITOLAK` | bilangan bulat | tidak | PK | keputusan tiket 73 — **baru** |
| `ID_MIGRASI_PENDARATAN` | bilangan bulat | tidak | FK | ERD baris 38 `LANDING_ID` |
| `JALUR_SIMPUL` | teks | tidak |  | jalur simpul di dalam muatan |
| `NILAI_MENTAH` | teks | tidak |  | apa adanya — **tidak dibulatkan, tidak dibuang** |
| `SEBAB_DITOLAK` | teks | tidak |  | keputusan tiket 73 — **baru** |

⛔ **Kenapa dua yang pertama tanpa kunci asing.** Jejak asal-usul harus **bertahan melewati**
penghapusan barisnya. Kunci asing memaksa salah satu dari dua, dan keduanya merusak: `ikut hapus`
menghapus jejaknya bersama barisnya, `tolak` membuat baris yang salah muat **tidak dapat dibuang**.
Pada `MIGRASI_PENDARATAN` ada sebab kedua yang berdiri sendiri: muatan yang **gagal diurai** justru
belum punya baris baru untuk ditunjuk. Ongkosnya dibayar dan dinyatakan: basis data tidak menjaga
keterhubungannya; yang menjaganya uji rekonsiliasi tiket `44`.

⚠️ **`MIGRASI_NILAI_DITOLAK` `tolak`, bukan `ikut hapus` — ini MENGOREKSI tiket `73`.** ERD baris 38
menulis `"di Go"`: ia tidak meresepkan aturan tingkat basis data. Dipilih `tolak` dengan alasan yang
sudah dipakai migrasi `425` untuk arsip — catatan forensik yang lenyap bersama induknya berhenti
menjadi catatan forensik tepat saat ia paling dibutuhkan.

## PENYEBARAN · RINCIAN_PENYEBARAN · NILAI_PENYEBARAN

Tiket `38`, migrasi `429`. Dari `ddl-usulan/{38,40,41}` dan `KAMUS-KOLOM.md` §10.6–§10.8.
Perilaku hapus seluruhnya dari `ERD.md` §2.5 (rantai induk, `ikut hapus`) dan §2.7 (rujukan acuan,
`tolak`).

⛔ **Kunci alaminya DIPASANG**, berbeda dari enam tabel ronde 6 yang ditahan: ketiganya sudah
bernomor di `Z00_KUNCI_ALAMI.sql` — `INV-16` untuk kedua `UQ_PENYEBARAN`, `INV-65` untuk
`UQ_RINCIAN_PENYEBARAN`. ⚠️ `NILAI_PENYEBARAN` **tidak** dapat nomor, jadi ia ikut tagihan §13.

### PENYEBARAN — induk polimorfik

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_PENYEBARAN` | bilangan bulat | tidak | PK | keputusan tiket 38 — **baru** |
| `ID_BAGIAN` | bilangan bulat | ya | FK UQ | salah satu dari **dua** pelekatan; tepat satu terisi, dijaga `CK_PENYEBARAN_INDUK` |
| `ID_DETAIL_PROPORSIONAL` | bilangan bulat | ya | FK UQ | pelekatan kedua, idem |
| `ID_JENIS_REASURANSI` | bilangan bulat | tidak | FK UQ | tabel acuan — hapus `tolak`, §2.7 |
| `ID_JENIS_REASURANSI_INDUK` | bilangan bulat | ya |  | penyebaran bersusun; ⚠️ **tanpa kunci asing** — sasarannya baris di tabel yang sama dan `ddl-usulan` pun tidak memasangnya |
| `PERSEN_PENYEBARAN` | angka desimal | tidak |  | `NUMBER(38,20)` dipersempit ke `NUMBER(38,8)` — `KTV-A` |

### RINCIAN_PENYEBARAN

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_RINCIAN_PENYEBARAN` | bilangan bulat | tidak | PK | keputusan tiket 38 — **baru** |
| `ID_PENYEBARAN` | bilangan bulat | tidak | FK UQ | induknya, `ikut hapus` §2.5 |
| `ID_JENIS_REASURANSI` | bilangan bulat | tidak | FK UQ | **kunci alami** `INV-65`; hapus `tolak` §2.7 |
| `PERSEN_RINCIAN` | angka desimal | tidak |  | `KTV-A` |

### NILAI_PENYEBARAN

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID_NILAI_PENYEBARAN` | bilangan bulat | tidak | PK | keputusan tiket 38 — **baru** |
| `ID_RINCIAN_PENYEBARAN` | bilangan bulat | tidak | FK | induknya, `ikut hapus` §2.5 |
| `NILAI` | angka desimal | tidak |  | `KTV-A` |
| `KODE_MATA_UANG` | teks | tidak |  | ⚠️ **teks, bukan kunci asing** — mengikuti kelima paket uang yang sudah berdiri (`NILAI_MDP` dan saudaranya). `ERD.md` §2.7 menyatakan relasinya `tolak`; menormalkan satu tabel saja membuat dua bentuk untuk satu fakta. **Dicatat, bukan ditambal sendirian** |

⛔ **`INV-32`, `INV-33`, `INV-52` tetap tidak dipasang** — ketiganya menuntut materialized view
`REFRESH ON COMMIT`, dan `F-13` menyatakan modul ini belum punya pemantau kebasiannya.

## M_TREATYIN_* — sembilan tabel PENDARATAN tab Treaty In

Migrasi `430` (tabel) dan `431` (sequence), 3 Oktober 2026.

⛔ **BUKAN tabel model baru, dan tidak menggantikan satu pun.** `PERIODE_PELAPORAN`, `PORTOFOLIO`,
`EGNPI`, `RETENSI_CEDANT`, `TERMIN`, `RINCIAN_ANGSURAN`, dan `PERIODE_AKUMULASI` tetap berdiri
persis seperti sebelumnya. Kedelapan tabel di bawah **mendaratkan larik di dalam
`POOLDATA.M_TREATY_IN.JSONDATA` apa adanya**, supaya layar berhenti mengurai CLOB pada setiap
pembacaan. Keduanya hidup berdampingan sampai pemindahan tiket `44` selesai. Polanya sama dengan
`MIGRASI_PENDARATAN`.

⚠️ **Kedelapannya tidak punya kotak di ERD, dan itu wajar** — ERD menggambar entitas sistem lama,
bukan tabel pendaratan. Lihat alasan yang sama pada `MIGRASI_KORELASI` di atas.

### Ukuran — SELURUH 1.854 dokumen

Disapu dengan mengurai ke-1.854 dokumen **secara utuh sebagai JSON**; nol dokumen gagal urai.
Bukan contoh 300, bukan contoh 66.

| Tabel | Larik JSON | Baris | Kontrak |
| --- | --- | ---: | ---: |
| `M_TREATYIN_REPORTINGPERIOD` | `ReportingPeriodList` | 4.548 | 1.137 |
| `M_TREATYIN_PORTFOLIO` | `Portfolio` | 1.925 | 845 |
| `M_TREATYIN_ACCUMULATION` | `AccumulationList` | 60 | 18 |
| `M_TREATYIN_EGNPI` | `EGNPI` | 2.298 | 846 |
| `M_TREATYIN_RETENTION` | `Retention` | 2.511 | 808 |
| `M_TREATYIN_INSTALLMENT` | `Installment` | 796 | 768 |
| `M_TREATYIN_INSTALLMENTITEM` | `Installment[].InstallmentList` | 3.033 | 768 |
| `M_TREATYIN_COMMENT` | `CommentList` | 11.365 | 1.837 |
| `M_TREATYIN_COINSCALE` | `CoInScale` | 702 | 186 |
| | **TOTAL** | **27.238** | |

### Tiga kolom struktur, sama di kesembilan tabel

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | `NUMBER(19)` | tidak | PK | `SEQ_MTI_<TAB>` — INV-02 |
| `MASTERID` | `VARCHAR2(100 CHAR)` | tidak | UQ(1) | `M_TREATY_IN.ID` — ⚠️ **nilai, BUKAN kunci asing**, lihat di bawah |
| `URUTAN` | `NUMBER(10)` | tidak | UQ(2) | indeks elemen di dalam lariknya, dari 0 |

⛔ **`MASTERID` bukan kunci asing sebab ia TIDAK DAPAT menjadi kunci asing.** `POOLDATA.TREATY_IN`
tidak punya kunci utama maupun `UNIQUE` pada `ID` — kolomnya bahkan `NULLABLE` — sehingga Oracle
menolaknya dengan ORA-02270. Uraian lengkap beserta **syarat pembalikannya** di
[`KEPUTUSAN-PENYELARASAN-REPO.md`](KEPUTUSAN-PENYELARASAN-REPO.md) §12.

⛔ **`URUTAN` bukan hiasan.** Larik Pega **berurut**: `Installment` nomor 1, 2, 3 bukan himpunan.
Urutan yang hilang tidak terlihat sampai seseorang membandingkan layar dengan sistem lama.

⚠️ **`UNIQUE (MASTERID, URUTAN)` adalah invarian PEMUAT, bukan invarian bisnis**, jadi ia tidak
menagih nomor di `SPEC-INVARIAN.md` — aturan `Z00_KUNCI_ALAMI.sql` berlaku untuk yang kedua. Ia ada
sebab pemuatnya wajib **idempoten**, dan idempotensi yang hanya dijaga kode pemanggil benar sampai
dua pemuat berjalan bersamaan.

### Kolom isi — seluruhnya `VARCHAR2`, dan panjang yang MENGUKURNYA

Setiap nilai di dalam `JSONDATA` adalah string JSON, termasuk yang terlihat seperti angka dan
tanggal. Kolom "maks" di bawah adalah panjang **terpanjang yang benar-benar ada** di ke-1.854
dokumen; lebar kolomnya diberi kelonggaran di atas angka itu.

| Tabel | Kunci JSON → kolom | maks | lebar |
| --- | --- | ---: | ---: |
| `REPORTINGPERIOD` | `AutoCalculate` | 5 | 50 |
| | `ConfirmationDue` | 8 | 50 |
| | `InitialDate` | 23 | 50 |
| | `Period` | 4 | 50 |
| | `SettlementDue` | 8 | 50 |
| | `SubmissionDue` | 8 | 50 |
| | `pxObjClass` → `PXOBJCLASS` | 39 | 200 |
| `PORTFOLIO` | `Description` | **869** | 4000 |
| | `Type` | 7 | 100 |
| | `TypePortfolio` | 10 | 100 |
| | `pxObjClass` | 35 | 200 |
| `ACCUMULATION` | `Period` | 3 | 50 |
| | `ReportDate` | 23 | 50 |
| | `SubDays` | 2 | 50 |
| | `SubDueDate` | 23 | 50 |
| | `pxObjClass` | 38 | 200 |
| `EGNPI` | `Amount` | 20 | 100 |
| | `AmountIDR` | 25 | 100 |
| | `AsDate` | 8 | 50 |
| | `ClassOfBusiness` | 32 | 200 |
| | `Currency` | 3 | 50 |
| | `CurrencyID` | 5 | 50 |
| | `Note` | **2.025** | 4000 |
| | `Proportion` | 24 | 100 |
| | `TreatyGroup` | 23 | 200 |
| | `TreatyGroupID` | 5 | 50 |
| | `pyTemplateRichTextEditor` | 3 | 50 |
| | `pxObjClass` | 31 | 200 |
| `RETENTION` | `Amount` | 13 | 100 |
| | `ClassOfBusiness` | 32 | 200 |
| | `Currency` | **9** ⚠️ | 50 |
| | `CurrencyID` | 5 | 50 |
| | `Note` | **1.568** | 4000 |
| | `TreatyGroup` | 23 | 200 |
| | `TreatyGroupID` | 5 | 50 |
| | `pxObjClass` | 35 | 200 |
| `INSTALLMENT` | `AmountTotal` | **61** ⚠️ | 200 |
| | `Currency` | 3 | 50 |
| | `PctTotal` | 18 | 100 |
| | `pxListSubscript` | 1 | 50 |
| | `pxObjClass` | 37 | 200 |
| `INSTALLMENTITEM` | `IDINDUK` — kunci asing ke induknya | — | `NUMBER(19)` |
| | `Amount` | 60 | 200 |
| | `Currency` | 3 | 50 |
| | `DueDate` | 8 | 50 |
| | `Installment` | 1 | 50 |
| | `InstallmentPct` | 18 | 100 |
| | `PaymentDate` | 8 | 50 |
| | `WPC` | 2 | 50 |
| | `pxObjClass` | 37 | 200 |
| `COMMENT` | `Date` → **`TANGGAL`** | 23 | 50 |
| | `ConvertDate` | 1 | 50 |
| | `HasHistory` | 1 | 50 |
| | `IsApproved` | 16 | 100 |
| | `OperatorName` | 17 | 200 |
| | `Suggest` | **2.171** | 4000 |
| | `pxObjClass` | 29 | 200 |

⚠️ **Dua angka yang menjelaskan kenapa seluruhnya teks:**

- `Installment.AmountTotal` terpanjang **61 aksara** —
  `1323411750.0000008394305684816601000000000000000000000000000`. `NUMBER(38,8)` akan
  **membulatkannya, diam-diam**, dan hasilnya tetap terlihat seperti angka yang masuk akal.
  Di sisi Go, `encoding/json` tanpa `UseNumber()` melakukan kerusakan yang sama lewat `float64`;
  `TestAngkaPanjangTidakDibulatkan` menguncinya.
- `Retention.Currency` terpanjang **9 aksara**, isinya **`1/04/2023`** — tanggal di dalam kolom
  mata uang. Ia harus mendarat apa adanya supaya dapat **ditemukan**, bukan ditolak di pintu.

⛔ **`Date` menjadi `TANGGAL`.** `DATE` kata cadangan Oracle (ORA-00923),
`TestNolKataCadanganOracleSebagaiKolom`. Sepuluh nama berisiko lain — `TYPE`, `PERIOD`, `NOTE`,
`AMOUNT`, `INSTALLMENT`, `DESCRIPTION`, `CURRENCY`, `SUGGEST`, `PROPORTION`, `URUTAN` — diuji satu
per satu dengan `SELECT 1 AS <nama> FROM DUAL` di Oracle DEV, dan kesepuluhnya lolos.

### M_TREATYIN_COINSCALE — tabel pendaratan KESEMBILAN, tab Co-Ins Scale

Migrasi `432` (tabel) dan `433` (sequence), 3 Oktober 2026. Mengikuti pola kedelapan tabel di atas
tanpa satu pun kekecualian.

| Kolom | Tipe | Null | Kunci | Sumber · maks terukur |
| --- | --- | --- | --- | --- |
| `ID` | `NUMBER(19)` | tidak | PK | `SEQ_MTI_COINSCALE` — INV-02 |
| `MASTERID` | `VARCHAR2(100 CHAR)` | tidak | UQ(1) | `M_TREATY_IN.ID` — nilai, bukan kunci asing |
| `URUTAN` | `NUMBER(10)` | tidak | UQ(2) | indeks elemen, dari 0 |
| `COINSHARE` | `VARCHAR2(200 CHAR)` | ya | | `CoInShare` — **17** |
| `PCTLIMIT` | `VARCHAR2(50 CHAR)` | ya | | `PctLimit` — 3 |
| `PXCREATEDATETIME` | `VARCHAR2(50 CHAR)` | ya | | `pxCreateDateTime` — 23 |
| `PXCREATEOPNAME` | `VARCHAR2(200 CHAR)` | ya | | `pxCreateOpName` — 24 |
| `PXCREATEOPERATOR` | `VARCHAR2(200 CHAR)` | ya | | `pxCreateOperator` — 17 |
| `PXCREATESYSTEMID` | `VARCHAR2(100 CHAR)` | ya | | `pxCreateSystemID` — 13 |
| `PXOBJCLASS` | `VARCHAR2(200 CHAR)` | ya | | `pxObjClass` — 35 |

**702 baris dari 186 kontrak**, terbanyak 5 per kontrak.

⛔ **`COINSHARE` TEKS, dan itu bukan kemalasan.** Nilainya PITA, bukan bilangan:
`>=30% up to < 50%`, `>=25%`. Kolom angka akan menolak seluruh 702 barisnya, dan pemformat angka
di layar mengembalikannya apa adanya justru untuk kasus ini.

⚠️ **Keempat medan jejak Pega dibawa**, meski rancangannya menyebut *"dua medan saja"*: keempatnya
ada pada **701 dari 702** elemen, dan ia satu-satunya catatan siapa yang menyusun skala itu dan
kapan — pertanyaan yang akan ditanyakan, sebab Co-Ins Scale adalah tab yang angkanya
dinegosiasikan.

## Empat tab dari `M_TREATY_IN2` — NOL tabel baru

Tab **Limits · Share · Event Limits · RNM Share** dibaca dari `POOLDATA.M_TREATY_IN2`, tabel
warisan datar 41 kolom, **satu baris per layer**. ⛔ **Nol tabel pendaratan dibuat untuk
keempatnya**, nol pemuatan, nol migrasi menyentuhnya — `repository/warisan_in2.go` BACA SAJA.

| | |
| --- | ---: |
| baris | 7.281 |
| kontrak | **1.340 dari 1.854** |
| `(MASTERID, LAYER)` berbeda | 2.540 |

⚠️ **Jangkauannya tidak penuh, dan layar mengatakannya.** 510 kontrak punya `Limits[]` berisi di
dokumennya (1.210 elemen) tanpa satu baris pun di tabel ini. Grid kosong di keempat tab itu
karena itu **tidak boleh** berbunyi *"kontrak ini memang tidak punya"*.

Pemetaan kolom demi kolom, dua garis buktinya, penggolongan angkanya, dan kesepuluh kolom kepala
yang tidak ditampilkan: [`PEMETAAN-M-TREATY-IN2.md`](PEMETAAN-M-TREATY-IN2.md).

⛔ **Pemindaian 41 kolom BERPOSISI** adalah kode paling mudah rusak di jalur ini: satu medan yang
tergeser memindahkan seluruh nilai sesudahnya ke kolom tetangganya, dan hasilnya tetap berupa grid
yang terisi rapi. `TestPindaiLayerSejajarDenganDaftarKolom` mengadu komentar tiap penugasan dengan
daftar kolomnya satu per satu.

### ⛔ `AchievementLists` TIDAK ada di sini, dan itu temuan

Rancangan ronde ini mendaftarnya sebagai tabel kesembilan, `M_TREATYIN_ACHIEVEMENT`, berinduk
`MASTERID`. Sapuan menemukan ia **tidak ada sebagai kunci puncak** — **nol kemunculan** di 1.854
dokumen. Jalur sebenarnya:

| Jalur | Kemunculan | Kosong |
| --- | ---: | ---: |
| `Limits[].Detail[].AchievementLists` | 1.961 | 859 |
| `RevisionHistory[].Limits[].Detail[].AchievementLists` | 3 | 0 |

Butirnya milik **satu baris `Limits[].Detail[]`**, bukan milik kontrak. Tabel berinduk `MASTERID`
saja **tidak dapat menyatakan butir itu milik baris limit yang mana**, dan yang kehilangan induknya
bukan sekadar kurang rapi: angka pencapaian yang menempel pada limit yang salah **terbaca benar**.
Induknya hidup di `M_TREATY_IN2` — tabel warisan yang berbeda, 7.281 baris — sehingga tabelnya
menuntut keputusan yang belum diambil: apakah `M_TREATY_IN2` ikut didaratkan.

`TestAchievementBukanTabelPendaratan` menolak penambahannya diam-diam.

### Jalur tulis, dan tiga sifat yang dituntut darinya

Satu-satunya penulis kedelapan tabel ini adalah `repository/pendaratan_muat.go`, dikemudikan
`services.MuatSatuKontrak` dan dijalankan `backend/pemuat/jalankan.go`.

| Sifat | Diwujudkan di | Dibuktikan oleh |
| --- | --- | --- |
| **idempoten** | `KosongkanKontrak` berjalan sebelum setiap sisip, di transaksi yang sama; `UNIQUE (MASTERID, URUTAN)` jaring keduanya | `TestMuatDuaKaliTidakMenggandakan`, `TestUrutanGandaDitolakConstraint` |
| **terbalikkan** | `KosongkanKontrak`, dan `alat/kosongkan-tab-treatyin.sql` untuk seluruh tabel — `DELETE`, bukan `TRUNCATE` | `TestKosongkanMengembalikanKeNolTanpaMenyentuhKontrakLain` |
| **tercocokkan** | cacah baris diadu dengan cacah elemen dokumen, **di dalam transaksi**, sebelum diikat | `TestSepuluhKontrakNyataMuatDanCocok` |

⚠️ **`DELETE`, bukan `TRUNCATE`.** `TRUNCATE` ber-DDL dan karena itu **mengikat transaksinya**;
pemuat yang gagal di tengah lalu `ROLLBACK` tidak akan mendapatkan kembali apa yang `TRUNCATE`
buang.

### Jalur baca — DUA BELAS tab berisi, satu mekanisme untuk prop dan non-prop

| Tab | Sumber hari ini | Pembaca |
| --- | --- | --- |
| Reporting Period | `M_TREATYIN_REPORTINGPERIOD` | `repository.BacaPeriodePelaporan` |
| Portfolio | `M_TREATYIN_PORTFOLIO` | `repository.BacaPortofolio` |
| Accumulation | `M_TREATYIN_ACCUMULATION` | `repository.BacaAkumulasi` |
| EGNPI | `M_TREATYIN_EGNPI` | `repository.BacaEgnpi` |
| Maximum Retention | `M_TREATYIN_RETENTION` | `repository.BacaRetensi` |
| Installment | `M_TREATYIN_INSTALLMENTITEM` | `repository.BacaAngsuran` |
| Information & Submit | `M_TREATYIN_COMMENT` | `repository.BacaCatatan` |
| Co-Ins Scale | `M_TREATYIN_COINSCALE` | `repository.BacaSkalaKoasuransi` |
| **Limits · Share · Event Limits · RNM Share** | **`M_TREATY_IN2`** *(warisan, baca saja)* | `repository.BacaLayerWarisan` — satu seam, empat proyeksi |
| Rate of Exchange | ⚠️ **masih `JSONDATA.CurrencyList`** | `repository.BacaKontrakWarisan` |
| **Exclusions · Special Conditions** | **`M_TREATY_IN.JSONDATA`** — Jalan B, keputusan §15 | `services.PilihTabTeks`, dipilih menurut cabang |
| Retro · Value Difference · Achievement In IDR | `.trin__belum` — lihat sebabnya di bawah | — |

**Yang masih `.trin__belum`, dan sebabnya BERBEDA tiap tab:**

| Tab | Sebab |
| --- | --- |
| Retro | **DITUNDA** oleh pemilik proses 4 Oktober 2026 ([`KEPUTUSAN` §14](KEPUTUSAN-PENYELARASAN-REPO.md)) — 2 kontrak dari 1.854 tidak cukup merancang tiga tabel bersarang. ⛔ Kedua kontraknya **DITANDAI**: `1000493` dan `1000755`, lihat §14.1 |
| Value Difference | ia **objek** berisi `EGNPI`/`Limits`/`Share` dan sembilan `Total*` — potret nilai sebelum perubahan, bukan teks. Dikeluarkan dari tab teks oleh [`KEPUTUSAN` §15](KEPUTUSAN-PENYELARASAN-REPO.md); ia milik tiket `06`/`11`/`13` (`NILAI_SELISIH`) |
| Achievement In IDR | induknya `Limits[].Detail[]` di `M_TREATY_IN2`, bukan kontrak — lihat bab `AchievementLists` di bawah |

### Dua tab TEKS — Jalan B, nol tabel

**Exclusions** dan **Special Conditions** dibaca dari `M_TREATY_IN.JSONDATA` lewat jalur warisan
yang sudah ada. ⛔ **Nol tabel pendaratan untuk keduanya** — keputusan
[§15](KEPUTUSAN-PENYELARASAN-REPO.md): layarnya baca-saja, dokumennya sudah ditarik untuk medan
lain, dan isinya sampai **23.453 aksara** sehingga Jalan A pun menuntut `CLOB`.

| Tab | Ejaan proporsional | Ejaan non-proporsional | Ejaan ketiga |
| --- | --- | --- | --- |
| Exclusions | `ExclusionsP` 1.155 | `Exclusions` 896 | — |
| Special Conditions | `SpecialConditionsP` 1.022 | `SpecialConditions` 686 | ⚠️ `SpecialConditionsp` **292** |

⛔ **Ejaannya BUKAN sinonim.** Dari **303** dokumen yang punya lebih dari satu ejaan
`SpecialConditions*`, **nol** yang isinya identik. Karena itu: ejaan yang sesuai cabang dipakai;
bila ia tidak ada hasilnya **kosong**, bukan ejaan lain; dan ejaan lain yang berisi **disebut di
layar**. Uraiannya beserta sebarannya di [`KEPUTUSAN` §15.1](KEPUTUSAN-PENYELARASAN-REPO.md).

⚠️ Cabangnya dibaca dari **kolom** `TREATY_IN.PROPORTIONTYPE`, bukan dari kunci `ProportionType`
di dalam dokumen — kuncinya tidak ada pada tiga kontrak (`1001854`–`1001856`).

⛔ **SATU MEKANISME untuk prop dan non-prop.** Nol cabang menurut `PROPORTIONTYPE` di jalur baca,
dan tidak boleh ada: yang membedakan kedua cabang adalah tab mana yang **ditampilkan**, bukan dari
mana isinya dibaca. Tab yang kosong mengembalikan nol baris, bukan galat.

⛔ **`ORDER BY URUTAN` di ketiga kueri.** Tanpanya Oracle bebas mengembalikan baris dalam urutan
apa pun, dan `Q 1`/`Q 2`/`Q 3` yang tertukar **terbaca benar**.
`TestUrutanBarisSamaDenganUrutanDokumen` mengadunya dengan dokumen, bukan dengan dirinya sendiri.

⚠️ **Rate of Exchange tetap di JSON, dan itu bukan kelalaian.** Ronde pemindahan melarang membuat
tabelnya atas dasar `TREATYEXCHANGEYEARLY` sudah melayaninya; sapuan menemukan dasar itu keliru —
tabel tersebut 140 baris per **tahun** dan nol kolomnya menunjuk kontrak, sementara `CurrencyList`
**3.465 elemen di 1.844 kontrak**. Melepaskannya dari JSON menuntut tabel kesembilan, dan
keputusan itu belum diambil. `4-erd-dan-tabel-datar/KOREKSI-ERD-VERSUS-POOLDATA.md` §12.

⛔ **Kunci JSON tanpa kolom dilaporkan, bukan ditelan.** Pega menambah properti tanpa memberi tahu
siapa pun, dan kunci baru yang tidak punya kolom hilang **tanpa galat** — pemuatnya tetap hijau dan
cacah barisnya tetap cocok. `repository.KunciTakTerpetakan` mencarinya pada setiap pemuatan.

## M_ATTACHMENTTREATY_2 — panel Attachment, tabel WARISAN

⛔ **NOL tabel baru, NOL migrasi.** Lampiran sudah relasional di sistem lama; aturan *"tabel baru
hanya untuk struktur tab ber-JSON"* tidak berlaku padanya. `repository/warisan_lampiran.go`
**BACA SAJA** — dijaga `TestWarisanHanyaDibaca`.

| | |
| --- | ---: |
| baris | **43** |
| `TREATYID` berbeda | 7 — **6 cocok** ke `TREATY_IN.ID`, **1 yatim** |
| kolom | `ID` `TREATYID` `CATEGORY` `FILENAME` `FILEMIMETYPE` `DATA_JSON`(CLOB) `DATEINPUT` `USERNAME` `CATEGORY_ID` `T_STORAGE_ID` |

⚠️ **`DATA_JSON` TIDAK dibaca.** Ia memuat muatan berkasnya; panel hanya menampilkan daftar, dan
menariknya pada tiap pembukaan form berarti memindahkan berkas yang tidak ada yang minta.
Berkasnya sendiri hidup di `T_STORAGE_IMAGE`, lewat `T_STORAGE_ID`.

### Kategori — tujuh terbukti, empat belum

Pasangan kode↔nama **dibaca dari data**, bukan dihafal: tabelnya menyimpan `CATEGORY_ID` dan
`CATEGORY` berdampingan.

| Kode | Nama | Baris |
| --- | --- | ---: |
| `00000` | Others | 5 |
| `00001` | Analysed Email | 10 |
| `00002` | Approval Email | 14 |
| `00005` | Summary Treaty Leader | 10 |
| `00006` | Assessment Inward Treaty Form / Format Analisa Treaty | 2 |
| `00007` | Pega Proportional Calculation /Perhitungan Pega Proportional | 1 |
| `00010` | Offer Email | 1 |

⛔ **`00003` `00004` `00008` `00009` — nol baris, dan namanya TIDAK ADA di korpus kedua modul.**
Keempatnya tampil dengan **kodenya, tanpa nama**, ditandai belum dipastikan. Empat nama yang
belum berumah disebut terpisah di layar. Menebak pasangannya menaruh berkas di kategori yang
salah, dan itu baru ketahuan bertahun kemudian —
[`PERTANYAAN-TERBUKA-KODE-KATEGORI-LAMPIRAN.md`](PERTANYAAN-TERBUKA-KODE-KATEGORI-LAMPIRAN.md).

⭐ Pertanyaannya **dapat terjawab sendiri**: begitu salah satu kode muncul di data lengkap dengan
`CATEGORY`-nya, katalog membacanya dan penandanya lepas tanpa perubahan kode.
`TestKatalogKategoriTerbacaDariData` merah pada hari itu.

### Aturan nama berkas — spanduk biru

*"Recommended safe substitute should be . or _"* (`Section/WorkAttachments.xml`) adalah **aturan**:
huruf, angka, titik, garis bawah, tanda hubung. Ditegakkan `services.NamaBerkasAman` untuk
unggahan **baru**, dan diuji dua arah — yang ditolak dan yang diterima.

⚠️ **Tidak surut.** Diukur terhadap ke-43 nama yang sudah tersimpan: **6 lolos, 37 tidak**.
Menolak 86% dari yang sudah ada berarti menutup pintu yang sistem lama buka.

### Panel History — nol pembacaan baru

Date · PIC · Approval · Comment berpadanan satu-satu dengan `M_TREATYIN_COMMENT`
(`TANGGAL` · `OPERATORNAME` · `ISAPPROVED` · `SUGGEST`), tabel pendaratan yang **sudah** terisi
11.365 baris. `repository.BacaCatatan` dipakai ulang apa adanya.

### Modul Adjustment

`Section/ShowAttachmentTreaty.xml`, `WorkAttachments.xml`, dan `TreatyInActionButtons.xml` ada di
**kedua** ekspor, jadi pembacaannya ada di kedua modul —
`modul/treatyinadjustment/backend/repository/warisan_lampiran.go`. ⛔ **Bukan impor silang**:
`TestModulTidakMengimporModulLain` menolaknya, dan dua modul yang berbagi kode repository
berbagi juga jadwal rilisnya. Kesamaannya **dinyatakan** di kepala kedua berkas alih-alih
disembunyikan; yang tidak boleh menyimpang hanya nama tabel dan nama kolomnya.

## ⛔ `NILAI_SELISIH` + `NILAI_SEBELUM_PRO_RATE` — PINDAH ke modul Adjustment

Keduanya **tidak lagi milik modul ini** sejak 4 Oktober 2026
([`KEPUTUSAN` §19](KEPUTUSAN-PENYELARASAN-REPO.md)). `435_cabut_nilai_selisih.sql` mencabutnya;
`treatyinadjustment/.../442_nilai_selisih.sql` membangunnya kembali dengan bentuk yang **identik**
— `diff` tanpa komentar mengembalikan nol selisih.

⚠️ **Kepala `426_nilai_selisih.sql` TETAP DI MODUL INI**, dan itu disengaja. Ia memuat empat
alasan rancangan yang masih mengikat — kardinalitas 1:N lawan 1:1 ERD, penghalang
`BESARAN_DAPAT_DISESUAIKAN`, induk `VERSI_KONTRAK` lawan `KONTRAK`, dan INV-58 (tabel bernama
`NILAI_SELISIH` tanpa kolom selisih). `442` **menunjuk balik** ke sana alih-alih menyalinnya dan
membiarkan salinannya membeku.

Keduanya **nol baris** sebelum dan sesudah pemindahan, jadi tidak ada yang hilang. Barisnya di
`TestPerilakuHapusSesuaiERD` ikut pindah, tidak digandakan.
