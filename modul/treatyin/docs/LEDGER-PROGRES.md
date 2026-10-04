# Ledger progres — seluruh 77 tiket

Diperbarui 2 Oktober 2026. **Tiap kali sesuatu mendarat, berkas ini ikut berubah.**

Berkas ini ada karena enam belas tiket berskema-tanpa-aplikasi **hilang dari pandangan selama dua
ronde**: tidak ada satu pun tempat yang memperlihatkan tiket mana sampai lapisan mana, sehingga
laporan tiap ronde hanya memuat yang mendarat, dan yang ke-skip tidak terlihat.

| Lambang | Artinya |
| --- | --- |
| OK | mendarat |
| — | belum — **kolom catatan wajib menyebut sebabnya** |
| n/a | tidak berlaku bagi tiket ini |

**`uji Oracle` berarti dijalankan terhadap Oracle sungguhan**, bukan terhadap teks DDL. Sembilan
constraint terbukti menolak dan dua jalur positif terbukti diterima pada 2 Oktober 2026
([`VERIFIKASI-ORACLE-2026-10-02.md`](VERIFIKASI-ORACLE-2026-10-02.md)); ia **belum** otomatis —
`make test-db` menunggu skema uji, dan akun aplikasi tidak punya `CREATE USER`.

⛔ **Nol tiket berstatus `selesai`**, dan itu bukan kelalaian: status adalah medan di dalam berkas
tiket, milik pemilik proses. Ledger ini mencatat **lapisan**, bukan status.

## Treaty In — batch 1 (`14`–`44`)

| tiket | skema | constraint | kaskade | migrasi | uji bentuk | uji Oracle | services | handler | catatan |
|---|---|---|---|---|---|---|---|---|---|
| `14` | OK | OK | OK | 401 | OK | OK | OK | OK | satu-satunya yang utuh sampai handler. INV-53 + INV-29 ditegakkan di services |
| `15` | OK | OK | OK | 400 | OK | OK | OK | OK | jalur baca enam tabel acuan |
| `16` | n/a | n/a | n/a | n/a | n/a | — | OK | OK | peringatan kunci alami ganda **menyebut pembandingnya**; simpan tetap BERHASIL (ADR-0040 §2). Uji Oracle menunggu skema uji |
| `17` | n/a | n/a | n/a | n/a | n/a | — | OK | OK | cari lewat nomor warisan; hasil kosong = jawaban, baris non-warisan tidak ikut terambil |
| `18` | n/a | n/a | n/a | n/a | n/a | — | OK | OK | kelima ruas beku ditolak, pesannya menyebut SELURUH ruas yang menyimpang; `UPDATE` SQL-nya pun tidak menyentuhnya |
| `19` | n/a | n/a | n/a | n/a | n/a | — | OK | OK | versi menyimpang ditolak, pesannya menyebut ruas + nilai yang berlaku. INV-59: nol kolom beku di `VERSI_KONTRAK` |
| `20` | OK | OK | OK | 403 | OK | — | — | — | INV-38 dan INV-57 belum ditegakkan — pernyataan keputusan di 403 |
| `21` | n/a | n/a | n/a | n/a | n/a | — | OK | — | kegagalan hitung mengembalikan **nil + keterangan**, bukan nol dan bukan kurs satu (ADR-0035). **Handler belum** — belum ada rute yang menghitung |
| `22` | OK | OK | OK | 404 | OK | — | — | — | INV-39/40 lubang, bukan penundaan — lihat 404 |
| `23` | OK | OK | OK | 405 | OK | — | — | — | |
| `24` | OK | OK | OK | 406 | OK | — | — | — | |
| `25` | OK | OK | OK | 407 | OK | — | — | — | INV-55 menunggu jalur simpan |
| `26` | OK | OK | OK | 408 | OK | — | — | — | INV-56 menunggu jalur simpan |
| `27` | OK | OK | OK | 409 | OK | — | — | — | |
| `28` | OK | OK | OK | 410 | OK | — | — | — | |
| `29` | OK | OK | OK | 411 | OK | OK | — | — | uji positif "bahaya kesembilan" ada di `invarian_db_test.go` |
| `30` | OK | OK | OK | 412 | OK | — | — | — | |
| `31` | OK | OK | OK | 413 | OK | OK | — | — | INV-05 terbukti menegakkan — klaim sebaliknya dicabut |
| `32` | OK | OK | **sebagian** | 413 | OK | OK | OK | OK | **mendarat ronde 5.** Tarif berbeda per pemulihan diterima, di luar 0-100 ditolak menyebut baris + nilainya. ⚠️ `kaskade` sebagian: `UNIQUE (ID_LAYER, NOMOR_URUT_PEMULIHAN)` **TIDAK dipasang** — `SPEC-INVARIAN.md` berhenti di `INV-71` dan penjaga melarang `UNIQUE` tanpa nomor. Kembar hanya ditolak di services |
| `33` | OK | OK | OK | 416 | OK | OK | — | — | INV-33 tertahan `F-13` |
| `34` | OK | OK | OK | 417 | OK | — | — | — | INV-32 tertahan `F-13` |
| `35` | n/a | n/a | n/a | n/a | n/a | — | OK | — | XOR quota-share/surplus + kesesuaian `JENIS_TREATY`; lima jalur negatif, dua positif. **Handler belum** |
| `36` | n/a | n/a | n/a | n/a | n/a | — | OK | — | surplus tanpa quota share ditolak, pesannya menyebut apa yang kurang. **Handler belum** |
| `37` | OK | OK | OK | 418 | OK | OK | — | — | `CK_POTONGAN_INDUK` terbukti dua arah; INV-52 tertahan `F-13`, INV-63 tinjauan kode |
| `38` | — | — | — | — | — | — | — | — | tertahan `Uji AD`. `L-3` sudah tutup, jadi penahannya tinggal satu |
| `39` | OK | OK | OK | 414 | OK | — | — | — | entitasnya saja; **mekanisme pengisian** belum — menunggu jalur simpan |
| `40` | n/a | n/a | n/a | n/a | n/a | — | OK | OK | **mendarat ronde 5.** Versi dasar dibaca lewat `ID_VERSI_KONTRAK_DASAR`; uji positifnya MENGUBAH versi dasar di tengah jalan — rancangan yang menyalin gagal di situ. Versi pertama menjawab 200 + `"ada": false`, bukan 404 |
| `41` | n/a | n/a | n/a | n/a | n/a | — | OK | OK | **mendarat ronde 5.** Diturunkan saat dibaca — nol tabel, nol kolom, nol cache (INV-58). Kontrak sistem baru menjawab `nomorLamaAda: false`, **bukan** nomor karangan |
| `42` | — | — | — | — | — | — | — | — | arsip JSON warisan — **tabelnya kini punya tiket: `74`** (`ARSIP_MUATAN_KELUAR`). Lapisan aplikasinya menunggu tabelnya, dan itu urutan yang benar |
| `43` | n/a | n/a | n/a | n/a | n/a | — | OK | — | sakelar pemindahan: keadaan terlihat tanpa basis data, alasan wajib, menyalakan **memeriksa ulang** baris yang masuk. **Handler belum** |
| `44` | — | — | — | — | — | — | — | — | isi tabel acuan dari sistem lama; **bertenggat** — ia memuat data pertama |

## Treaty In — batch 2 (`45`–`64`)

**Nol dari dua puluh punya tabelnya.** `CATATAN_PERSETUJUAN` (tiket `54`) belum berdiri, dan ia yang
`49`–`56` gantungi.

| tiket | skema | penahan | catatan |
|---|---|---|---|
| `45` | — | — | daftar keadaan siklus hidup. **Pembuka batch** — sendirian melepas `46 47 48 54 62`. Kolomnya berdiri sejak `401`; **artinya** belum. Tidak dikerjakan ronde 4 — di luar daftar ronde ini |
| `46` | — | `45` | versi terminal beku menjangkau seluruh anaknya |
| `47` | — | `45` | paling banyak satu versi tak-terminal per kontrak |
| `48` | — | `45` | pengajuan menolak dua syarat, memperingatkan enam |
| `49` | — | `48` · `54` | SH menyetujui, versi pindah ke antrian DH |
| `50` | — | `49` | DH menyetujui, versi pindah ke antrian DR |
| `51` | — | `50` | DR menyetujui, versi menjadi `DISETUJUI` |
| `52` | — | `51` | pengaju tidak menyetujui; tiap tingkat orang berbeda |
| `53` | — | `49` | pengembalian ke `DRAFT` dengan alasan wajib |
| `54` | — | `45` | **`CATATAN_PERSETUJUAN` — PEMBUAT PERTAMA**, tabelnya belum berdiri. `49`–`56` menggantungnya |
| `55` | — | `49` · `54` | penolakan versi meninggalkan catatan |
| `56` | — | `45` · `54` | pembatalan draf oleh pembuatnya sendiri |
| `57` | — | `46` · `20` | angka rupiah pada versi disetujui tidak bergeser |
| `58` | — | `51` · `75` · **dokumen batas wewenang** | **berstatus `tertahan`** — dan penahannya **DUA**: dokumennya belum ada, dan tabelnya pun belum (`BATAS_WEWENANG`, tiket `75`) |
| `59` | — | `45` · `44` | kepala kontrak warisan dipindahkan beserta keadaannya |
| `60` | — | `59` | keadaan warisan tanpa padanan → `WARISAN_TAK_TERPETAKAN` |
| `61` | — | `60` · `71` · `73` | baris warisan diperbaiki ke keadaan sah; penelusurannya menuntut `MIGRASI_KORELASI` (`71`) dan daftar kerjanya `MIGRASI_NILAI_DITOLAK` (`73`) |
| `62` | — | `45` | kontrak dicari termasuk menurut keadaan siklus hidupnya |
| `63` | — | `14` · **`Uji X-2`** | **berstatus `tertahan`** — cara pembukuan XOL |
| `64` | — | `15` · **penetapan pemilik** | **berstatus `tertahan`** — tabel acuan pembagian kapasitas |

**Dua puluh tiket, nol berskema.** `CATATAN_PERSETUJUAN` (tiket `54`) adalah satu-satunya tabel baru
yang batch ini butuhkan, dan ia belum berdiri.

> ⚠️ **RALAT.** Bagian ini pernah meringkas tiga belas tiket menjadi dua baris rentang
> (*"`46`–`53`, `55`–`58`"*), sehingga ledger yang mengaku mencakup **64 tiket** sebenarnya hanya
> memuat **51**. Itu terjadi di berkas yang justru dibuat supaya yang ke-skip terlihat — dan
> peringkasan itu menyembunyikan persis apa yang hendak diperlihatkannya.

## Treaty In Adjustment (`01`–`13`)

| tiket | skema | constraint | kaskade | migrasi | uji bentuk | uji Oracle | services | handler | catatan |
|---|---|---|---|---|---|---|---|---|---|
| `01` | OK | OK | OK | 440 | OK | — | — | — | kolom + FK berdiri; keempat penolakannya menunggu jalur simpan |
| `02` | sebagian | — | n/a | 401 | — | — | — | — | kolomnya ada sejak 401; **dua nilainya belum dinyatakan** — sisa SKEMA |
| `03` | sebagian | — | n/a | 401 | — | — | — | — | kolomnya ada; **bawaan tanggal mulai belum ada** — sisa SKEMA |
| `04` | — | — | — | — | — | — | — | — | `DOKUMEN_ADDENDUM`; tertahan `DB-16a`, belum dikirim |
| `05` | OK | OK | n/a | 441 | OK | OK | n/a | n/a | `NOMOR_URUT_VERSI` nullable — terbukti di Oracle |
| `06` | BLOKIR | — | — | — | — | — | — | — | `NILAI_SELISIH` **kini punya tiketnya: `76`** papan Adjustment. `PEMBUAT PERTAMA` pindah ke sana; tiket ini tinggal perilaku padanan kunci bisnisnya |
| `07` | — | — | — | — | — | — | — | — | titik beku materialitas; tertahan `DB-20` |
| `08` | — | — | — | — | — | — | — | — | tanggal berlaku dokumen; tertahan `DB-16b` |
| `09` | — | — | — | — | — | — | — | — | nomor dokumen warisan; menunggu `04` |
| `10` | n/a | — | n/a | n/a | — | — | — | — | pengisian nomor urut warisan; **tidak dikerjakan ronde 3** |
| `11` | — | — | — | — | — | — | — | — | INV-69/70 + pemantau kebasian; menunggu `06`, dan `06` menunggu `76` |
| `12` | n/a | — | n/a | n/a | — | — | — | — | kerutkan urutan; menunggu `10` |
| `13` | BLOKIR | — | — | — | — | — | — | — | menunggu `11`, dan `NILAI_SELISIH` yang kini dibangun tiket `76` |
| `76` | — | — | — | — | — | — | — | — | **tiket baru ronde 5** — `NILAI_SELISIH`. **Dapat dimulai**: `VERSI_KONTRAK` berdiri. ⚠️ kunci asing keduanya menunjuk `BESARAN_DAPAT_DISESUAIKAN` yang **tidak ada di mana pun** — penghalang, bukan pekerjaan |
| `77` | — | — | — | — | — | — | — | — | **tiket baru ronde 5** — `NILAI_SEBELUM_PRO_RATE`. **Dapat dimulai**. Bukan turunan: `GRL-15` membuang mesin pro rata, jadi angkanya tidak dapat dihitung mundur |

## Treaty In — batch 3 (`65`–`75`) · entitas dari rekonsiliasi 2 Oktober 2026

**Nol dari sebelas punya tabelnya.** Kesebelasnya baru ditiketkan 2 Oktober 2026 dari
`4-erd-dan-tabel-datar/REKONSILIASI-XML-VS-DDL.md` — 18 entitas yang sapuan 329 XML punya dan
`ddl-usulan/` tidak. Ledger mencatatnya **sejak hari pertama** justru supaya kesebelasnya tidak
mengulang nasib enam belas tiket berskema-tanpa-aplikasi yang hilang dari pandangan dua ronde.

| tiket | entitas | skema | penahan | catatan |
|---|---|---|---|---|
| `65` | `KELAS_BISNIS_LAYER` | — | `34` | **dapat dimulai**; induknya berdiri |
| `66` | `BESARAN_LAYER` | — | `31` · `64` | `64` tertahan penetapan pemilik — acuan peran besarannya belum ada |
| `67` | `KELOMPOK_LAYER` | — | `31` | **dapat dimulai** |
| `68` | `KELAS_BISNIS_KELOMPOK` | — | `67` | kedalaman keempat; **bukan** kembaran `65` |
| `69` | `PENCAPAIAN` | — | `14` | **dapat dimulai**. Induknya `KONTRAK` per `ERD.md` §2.9, bukan `DETAIL_PROPORSIONAL`; `LOG_ACHIEVEMENT` **tidak** jadi tabel kedua |
| `70` | `RINCIAN_ANGSURAN` | — | `27` | **dapat dimulai**; `TERMIN` berdiri sejak `409` tanpa anaknya |
| `71` | `MIGRASI_KORELASI` | — | `14` | **dapat dimulai** · ⛔ **bertenggat terhadap `44`** |
| `72` | `MIGRASI_PENDARATAN` | — | — | **dapat dimulai, nol penahan** · ⛔ **bertenggat terhadap `44`** |
| `73` | `MIGRASI_NILAI_DITOLAK` | — | `72` | ⛔ **bertenggat terhadap `44`** |
| `74` | `ARSIP_MUATAN_KELUAR` | — | `14` | **dapat dimulai**; ia tempat penyimpanan tiket `42` |
| `75` | `BATAS_WEWENANG` | — | `64` · **dokumen batas wewenang** | **berstatus `tertahan`** — ambang berangkanya belum ada di satu berkas pun |

> ⛔ **`71` `72` `73` dikerjakan SEBELUM `44`.** Tiket `44` memuat data pertama; sesudahnya tidak
> ada lagi jalan menelusuri baris baru kembali ke asalnya, dan tiket `61` menuntut justru
> penelusuran itu.

## Dua entitas yang TIDAK ditiketkan, dan itu keputusan

`RINGKASAN_LIMIT` (`T_TREATY_LIMIT_SUMMARY`) dan `REKAP_KONTRAK` (`T_TREATY_TOTAL`) diperiksa
2 Oktober 2026 dan diputuskan **turunan** — menyimpannya melanggar `INV-58`. Bukti yang menentukan:
sapuan atas setiap nama tabel yang disebut SQL di seluruh korpus mengeluarkan **enam belas** tabel,
dan **tidak satu pun** bernama demikian. Rinciannya berikut **syarat pembalikannya** di
`REKONSILIASI-XML-VS-DDL.md` §4.

## Dua lubang spec yang TIDAK ditambal

~~**`NILAI_SELISIH`**~~ — **ditambal 2 Oktober 2026: tiket `76` papan Adjustment.** Ia pernah
berdiri di sini sebagai lubang sebab `ERD.md` §2.6 menyatakannya mengikat sementara `ddl-usulan/`
dan `KAMUS-KOLOM.md` tidak memuatnya. Tiket `76` merancang bentuknya dari
`T_TREATY_VALUE_DIFFERENCE` dengan **mata uang di dalam kunci** (`ADR-0048` butir 3), dan menuntut
entrinya masuk `KAMUS-KOLOM.md` lewat definisi §10 — bukan dengan menyunting hasil bangkitannya.
Satu penghalang ikut lahir bersamanya: `BESARAN_DAPAT_DISESUAIKAN`, induk kunci asing keduanya,
**tidak ada di mana pun**.

**`PERISTIWA_KONTRAK`** — tabelnya ada di `ddl-usulan/` dan relasinya di `ERD.md` §2.3b, tetapi
**nol tiket dan nol kemampuan `P-nn` menghasilkannya**. Tidak disentuh.

## Ringkasan lapisan

| Lapisan | Cacah dari 77 | Perubahan ronde 5 |
| --- | ---: | --- |
| skema berdiri | **19** | — · **penyebutnya berubah**: 19 dari 77, bukan 19 dari 64 |
| kaskade sesuai `ERD.md` | **19** | — |
| uji bentuk | 19 | — |
| uji Oracle (manual, belum otomatis) | 6 | — |
| **services** | **13** | **+3** — tiket `32` `40` `41` |
| **handler** | **9** | **+3** — tiket `32` `40` `41` |
| **tiket di papan** | **77** | **+13** — `65`…`75` Treaty In, `76` `77` Adjustment |
| berstatus `selesai` | **0** | status milik pemilik proses |

⚠️ **Penyebutnya bertambah dan pembilangnya tidak.** Tiga belas entitas yang hilang kini punya
tiketnya, dan **nol di antaranya punya tabelnya** — jadi "19 dari 77" adalah gambaran yang lebih
jujur daripada "19 dari 64", bukan kemunduran. Yang berubah bukan pekerjaannya, melainkan apa yang
diakui sebagai pekerjaan.

⚠️ **Nol tabel baru dibuat ronde 5**, dan nol migrasi baru mendarat. Satu migrasi sempat ditulis
(`421`, `UNIQUE` untuk `PEMULIHAN_LIMIT`) lalu **dibatalkan**: penjaga
`TestTabelTanpaKunciAlamiTidakDiberiDiamDiam` melarang `UNIQUE` sebelum nomor invariannya turun, dan
ia lolos hanya karena nama constraint dan kata `UNIQUE` terpisah baris. Penjaganya **diperkuat**
(spasi dirapatkan lebih dulu) alih-alih migrasinya dipertahankan.
