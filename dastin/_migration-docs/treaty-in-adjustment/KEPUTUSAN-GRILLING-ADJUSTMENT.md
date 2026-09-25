# Log grilling Treaty In Adjustment — INDEKS lintas ronde

> ## 🔴 GRILLING DIBUKA KEMBALI — 24 September 2026, sesudah gerbang 0 to-spec
>
> Tiga jawaban bisnis yang **tidak pernah masuk ronde mana pun** membatalkan dua putusan terkunci:
>
> | Putusan | Pembatalnya | Keadaan |
> |---|---|---|
> | **`GRL-04`** | `DB-3` **dan** `DB-4` — **kedua cabang pembatalnya menyala** | **BATAL** |
> | **`GRL-12`** | `DB-15` | **BATAL** |
> | `GRL-18` | turunan `GRL-12` | **sebagian tersentuh** — bagian 1 dan 3 berdiri |
>
> **`GRL-14` diperiksa ulang dan BERTAHAN** — dasarnya kini terverifikasi mandiri lewat arah
> prasyarat langkah 2 `SaveTreatyIn_EDM_Act`.
>
> **DIJAWAB RONDE D, 24 September 2026.** Keduanya disidang dan **dibatalkan**; penggantinya
> **`GRL-19`** (`DOKUMEN_ADDENDUM` entitas tersendiri) dan **`GRL-20`** (materialitas adalah
> masukan, sakelar dua arah). Laporan gerbang [`GERBANG-0-TO-SPEC.md`](GERBANG-0-TO-SPEC.md);
> putusannya [`GRILL-D/06-PUTUSAN.md`](GRILL-D/06-PUTUSAN.md).
>
> ## ✅ TO-SPEC ADJUSTMENT TERBUKA — 24 September 2026
>
> Keempat butir yang ronde D sebut *"menahan model"* **diuji ulang terhadap satu pertanyaan —
> *apakah ia menentukan letak kolom?*** — dan **tidak satu pun lolos**. Keempatnya memblokir
> **kepastian**, bukan **model**.
>
> Keempatnya diputuskan **tanpa verifikasi**, masing-masing berdasar dan bersyarat pembalikan:
> [`KEPUTUSAN-TANPA-VERIFIKASI-ADJUSTMENT.md`](KEPUTUSAN-TANPA-VERIFIKASI-ADJUSTMENT.md) —
> **`KTV-1` … `KTV-4`, nol terverifikasi.**
>
> Ketiga `DB` **tetap dikirim**; jawabannya kini **menyempitkan**, tidak lagi memblokir.
> Prasyarat cabang K §4 butir 2 diubah bunyinya **di tempatnya sendiri** menjadi *"diserahkan"*;
> **tanggapan paket `REV` ditagih sebelum to-ticket.**

**Fase:** grilling (bukan to-spec, bukan to-ticket). ~~**SELESAI 24 September 2026**~~ untuk seluruh bagian yang dapat diselesaikan sesi mana pun — sisa tunggal: paket `REV-1`…`REV-6` ditanggapi pemilik ADR induk (`GRILL-TDA/05-GERBANG.md` §4).
**Dimulai:** 23 September 2026 · **Diperbarui:** 24 September 2026.
**Mengikat:** [`../METODE-GRILLING.md`](../METODE-GRILLING.md) — disalin ke akar `_migration-docs/` pada 24 September 2026 karena berlaku untuk modul-modul berikutnya juga.

> **Berkas ini indeks, bukan sumber kebenaran.** Isi lengkap setiap keputusan ada di
> `GRILL-<ronde>/06-PUTUSAN.md`. Bila sesi terputus, baca indeks ini lebih dulu, lalu berkas ronde
> yang statusnya masih TERBUKA.

---

## Ronde

| Ronde | Cabang | Berkas | Status |
|---|---|---|---|
| **A** | A — batas lingkup + rekonsiliasi ADR induk | [`GRILL-A/`](GRILL-A/) | **DITUTUP 24 Sep 2026** — 8 pertanyaan, 8 putusan |
| **B** | B — identitas, penomoran, dan rantai addendum | [`GRILL-B/`](GRILL-B/) | **DITUTUP 26 Sep 2026** — GRL-09, GRL-10 |
| **C** | C — jenis, materialitas, `ActualValue` | [`GRILL-C/`](GRILL-C/) | **DIBUKA KEMBALI 24 Sep 2026** — C1 **🔴 GRL-12 BATAL**; C2 ✅ GRL-13; C3 ✅ GRL-14 (dasarnya diverifikasi ulang, **bertahan**) |
| **D** | D — kemampuan baru dan akibatnya ke orang | — | **ditutup tanpa ronde** — [`PEMILAHAN-SISA-GRILLING.md`](PEMILAHAN-SISA-GRILLING.md) |
| **D′** | **ronde pendek — dua pembatal yang sudah menyala** | [`GRILL-D/`](GRILL-D/) | **DITUTUP 24 Sep 2026** — `GRL-04` dan `GRL-12` **BATAL**; `GRL-19`, `GRL-20` dikunci |
| **E** | E — pembekuan dan mesin selisih | [`GRILL-C/`](GRILL-C/) | **DITUTUP 24 Sep 2026** — E1 KONFIRMASI; E3a ✅ **GRL-15**, E3b ✅ **GRL-16**; **E3c LARUT** oleh GRL-14 |
| **F** | F — bentuk simpan | — | **ditutup tanpa ronde** — idem |
| **G** | G — peran dan wewenang | — | **ditutup tanpa ronde** — idem |
| **H** | H — layar dan editabilitas | — | **ditutup tanpa ronde** — idem |
| **I** | I — migrasi data warisan | [`GRILL-C/`](GRILL-C/) | **DITUTUP 24 Sep 2026** — I1 ✅ GRL-11; payung I2 **KONFIRMASI**, residunya ✅ **GRL-17** |
| **J** | J — lampiran | — | **ditutup tanpa ronde** — idem |
| **K** | K — prasyarat to-spec | — | **ditutup tanpa ronde** — idem |
| **TDA** | *(bukan cabang)* sapuan penutup Lacak TDA | [`GRILL-TDA/`](GRILL-TDA/) | **DITUTUP** — GRL-18; ~~16 TDA, 16 nasib~~ **17 TDA, 17 nasib** (`TDA-17` diadili ronde D, `TD-05`); `NT-01`, `NT-02`; usulan `TA-09` |


> **Ronde `TDA` bukan cabang D.** Konvensi *"ronde = cabang"* (`METODE` §3.3) berlaku untuk A, B,
> dan C. Ronde `TDA` adalah **sapuan lintas cabang** yang lahir dari titik periksa `GRILL-C/07-AUDIT`
> §4a, dan ia diberi nama huruf-bukan-cabang supaya satu huruf tidak berarti dua hal — kekeliruan
> yang sudah pernah terjadi di modul induk (`L-5`).

## Keputusan — satu baris per GRL

| # | Judul | Label | Status | Isi lengkap |
|---|---|---|---|---|
| **GRL-01** | Satu model, bukan satu berkas: addendum adalah `VERSI_KONTRAK` milik spesifikasi induk | PERUBAHAN | terkunci (`PP-1` menunggu) | [`GRILL-A/06-PUTUSAN.md`](GRILL-A/06-PUTUSAN.md) |
| **GRL-02** | ADR-0036 berlaku dengan penyesuaian: angka dasar persetujuan tidak berubah sesudah disetujui | PERUBAHAN | prinsip terkunci; mekanisme di cabang E; `DB-1` menunggu | [`GRILL-A/06-PUTUSAN.md`](GRILL-A/06-PUTUSAN.md) |
| **GRL-03** | ADR-0048 berlaku dengan penyesuaian: bingkainya mati, ketiga kewajibannya hidup | 3 PERUBAHAN + 1 PELESTARIAN sementara | terkunci; `DB-2` menunggu; label butir 4 menunggu cabang B | [`GRILL-A/06-PUTUSAN.md`](GRILL-A/06-PUTUSAN.md) |
| **GRL-04** | ~~Penyesuaian bukan entitas; `PENYESUAIAN` tidak kembali~~ | PELESTARIAN | **🔴 BATAL 24 Sep 2026** — kedua cabang pembatalnya menyala. **Digantikan `GRL-19`.** Yang TIDAK gugur: penelusuran 19 properti `METODE` §8.3, sehingga **`PENYESUAIAN` tetap tidak kembali** | [`GRILL-A/06-PUTUSAN.md`](GRILL-A/06-PUTUSAN.md) |
| **GRL-05** | ADR-0040 berlaku dengan penyesuaian; butir 5 dilepas ke C1 | butir 1 PELESTARIAN · butir 3 PERUBAHAN | terkunci; `DB-6`, `DB-7` menunggu; butir 5 pindah ke C1 | [`GRILL-A/06-PUTUSAN.md`](GRILL-A/06-PUTUSAN.md) |
| **GRL-06** | ADR-0049 ditutup dengan rujukan; klasifikasi materialitas warisan tidak dihitung ulang | PELESTARIAN | terkunci | [`GRILL-A/06-PUTUSAN.md`](GRILL-A/06-PUTUSAN.md) |
| **GRL-07** | ADR-0052 berlaku dengan penyesuaian: satu rantai, dan penyimpangan yang berjalan hanya **satu** | butir 1-3 PELESTARIAN · butir 4 PERUBAHAN | terkunci; `DB-10` terbuka pada sisi revisi-di-tempat; Konteks dikoreksi lewat `REV-2` | [`GRILL-A/06-PUTUSAN.md`](GRILL-A/06-PUTUSAN.md) |
| **GRL-08** | ADR-0055 berlaku dengan penyesuaian: daftar keadaan dan perpindahannya dipakai apa adanya untuk versi addendum | **per perpindahan** — 7 PELESTARIAN, 3 PERUBAHAN, 2 BARU | terkunci; `REV-3` memuat empat koreksi | [`GRILL-A/06-PUTUSAN.md`](GRILL-A/06-PUTUSAN.md) |
| **GRL-09** | Pengenal tampilan: diturunkan untuk versi baru, dilestarikan apa adanya untuk baris warisan | warisan PELESTARIAN · bentuk turunan PERUBAHAN | terkunci; `DB-11` dan `UA-14` terbuka; bahan to-spec B-1 dan B-2 | [`GRILL-B/06-PUTUSAN.md`](GRILL-B/06-PUTUSAN.md) |
| **GRL-10** | Rantai versi: tidak ada percabangan; dasar = versi berlaku terakhir; rujukan dasar disimpan tepat di satu tempat | dasar PERUBAHAN · `OLDID` warisan PELESTARIAN · bentuk simpan KONFIRMASI | terkunci; bahan to-spec B-3 dan B-4; `DB-2` dan `UA-1(c)` terbuka | [`GRILL-B/06-PUTUSAN.md`](GRILL-B/06-PUTUSAN.md) |
| **GRL-11** | Versi berlaku adalah turunan; `KONTRAK` tidak membawa penunjuk | PERUBAHAN | terkunci; bahan to-spec I-1 dan I-2; `UA-17` dan `DB-12` terbuka | [`GRILL-B/06-PUTUSAN.md`](GRILL-B/06-PUTUSAN.md) |
| **GRL-12** | ~~Materialitas diturunkan dari baris selisih yang **tersimpan**~~ | PERUBAHAN | **🔴 BATAL 24 Sep 2026.** Dikuatkan ekspor: **660 dari 758 kondisi penguncian menyebut `EDMMaterialType`**. **Digantikan `GRL-20`**; `REV-4` ikut tertahan | [`GRILL-B/06-PUTUSAN.md`](GRILL-B/06-PUTUSAN.md) |
| **GRL-13** | `JENIS_ADDENDUM` bernilai **dua** untuk versi baru: `PENYESUAIAN_PREMI` dan selebihnya | `PENYESUAIAN_PREMI` PELESTARIAN · Internal/External PERUBAHAN | terkunci; `DB-16` pintu ke (b+); `DB-17` menopang penamaan; **kalimat *"penyesuaian premi selalu material"* MENYUSUT — lihat GRL-14** | [`GRILL-B/06-PUTUSAN.md`](GRILL-B/06-PUTUSAN.md) |
| **GRL-14** | `ActualValue` dipecah menurut artinya; tidak ada pohon ketiga di sistem baru | PERUBAHAN | terkunci; `DB-15` dan **`DB-19`** pembatalnya; **TDA-06 dan TDA-16 diperbaiki**; E3c larut | [`GRILL-C/06-PUTUSAN.md`](GRILL-C/06-PUTUSAN.md) |
| **GRL-15** | Pro rata: atribut *"berlaku sejak"* dibawa sekarang, **mesinnya sengaja tidak dibangun** | atribut PERUBAHAN · mesin = keputusan untuk tidak membangun | terkunci; ditagih saat `DB-5` dijawab; `UA-19`, `UA-20`; **`REV-6`** diajukan | [`GRILL-C/06-PUTUSAN.md`](GRILL-C/06-PUTUSAN.md) |
| **GRL-16** | Share fakultatif tidak memerlukan kemampuan tersendiri; yang hilang adalah **kelas** cacatnya | PERUBAHAN | terkunci; `DB-18` pembatalnya; **wajib uji negatif** — tanpa itu ia klaim | [`GRILL-C/06-PUTUSAN.md`](GRILL-C/06-PUTUSAN.md) |
| **GRL-17** | `NOMOR_URUT_VERSI` warisan diberikan ulang menurut kronologi; `/Rnn` dilestarikan sebagai pengenal | PERUBAHAN | terkunci; `UA-21` mengukur luasnya tanpa membatalkan; **dikecualikan dari kriteria identik ADR-0043**, dan sebabnya wajib tertulis | [`GRILL-C/06-PUTUSAN.md`](GRILL-C/06-PUTUSAN.md) |
| **GRL-18** | Jenis addendum adalah **masukan**; dapat dibetulkan selama `DRAFT` dengan jejak; **beku sejak diajukan** | pilihan orang PELESTARIAN · beku + berjejak PERUBAHAN | **🟡 SEBAGIAN TERSENTUH** — bagian 1 dan 3 berdiri; bagian yang menyatakan **sumbu materialitas tidak ada** jatuh bersama `GRL-12`. `DB-16` tetap pembatal bagian jenis | [`GRILL-TDA/06-PUTUSAN.md`](GRILL-TDA/06-PUTUSAN.md) |
| **GRL-19** | **`DOKUMEN_ADDENDUM` adalah entitas tersendiri, di ATAS versi** — satu dokumen memayungi banyak versi lintas kontrak; persetujuan tetap per versi | entitas **BARU** · persetujuan PELESTARIAN | terkunci; **migrasi TIDAK punya sumber untuk mengisinya** (`TD-01` nol); `UA-19`, `DB-16a`, `DB-16b` terbuka | [`GRILL-D/06-PUTUSAN.md`](GRILL-D/06-PUTUSAN.md) |
| **GRL-20** | **Materialitas adalah MASUKAN, dipilih sebelum perubahan, dan ia SAKELAR DUA ARAH** — Material mengunci teks, Non Material mengunci angka | masukan + penguncian PELESTARIAN · ditegakkan saat simpan PERUBAHAN · dua arah PERUBAHAN | terkunci; `INV-69`/`INV-70` diusulkan; `DB-20` dan `UA-3` terbuka; **`TDA-10` ditutup ulang** | [`GRILL-D/06-PUTUSAN.md`](GRILL-D/06-PUTUSAN.md) |

## Status cabang

Urutan **A → B → C → D → E → F → G → H → I → J → K**. Satu-satunya penajaman: cabang C dibuka oleh
**C1 — definisi materialitas**, yang menjadi gerbang bagi E.

| # | Cabang | Isi | Status |
|---|---|---|---|
| 1 | **A** | Batas lingkup + rekonsiliasi ADR induk: 0036 ✓, 0048 ✓, 0040 ✓, 0049 ✓, 0052 ✓, 0055 ✓ | **DITUTUP** — 6 dari 6 ADR selesai; 8 pertanyaan; titik periksa kedua di `GRILL-A/07-AUDIT.md` §4 |
| 2 | B | Model versi dan identitas | **DITUTUP** — GRL-09, GRL-10 |
| 3 | C | Jenis addendum — dibuka oleh **C1**, gerbang bagi E | **DITUTUP** — GRL-12, GRL-13, GRL-14 |
| 4 | D | Daur hidup dan persetujuan | **ditutup tanpa ronde** — `PEMILAHAN-SISA-GRILLING.md` |
| 5 | E | Mesin selisih | **DITUTUP** — E1 KONFIRMASI, GRL-15, GRL-16; E3c larut |
| 6 | F | Penyimpanan di Oracle | **ditutup tanpa ronde** — KONFIRMASI `STRUKTUR-DATA.md` + ADR-0034 |
| 7 | G | Hak akses | **ditutup tanpa ronde** — KONFIRMASI ADR-0044 |
| 8 | H | Layar React | **ditutup tanpa ronde** — REKOMENDASI `R-H1` dan sisa tampilan `TDA-11` |
| 9 | I | Migrasi data lama | **DITUTUP** — GRL-11, GRL-17; payung I2 KONFIRMASI |
| 10 | J | Lampiran | **ditutup tanpa ronde** — KONFIRMASI `DOKUMEN_KONTRAK` dirujuk, tidak dimiliki |
| 11 | K | Prasyarat to-spec | **terbuka** — satu-satunya, dan penahannya **bukan pekerjaan**: paket `REV-1`…`REV-6` menunggu tanggapan pemilik ADR induk |

## Lacak TDA

**Kosakata nasib** (`METODE` §6.3): diperbaiki · dilestarikan · ditunda.
**Kolom sisi** (`METODE` §8.2) diperiksa per berkas lewat `pzInsKey`, bukan dari nama berkas.
**IRISAN** berarti temuannya **temuan Treaty In** — diadili sekarang, dan dicatat sebagai usulan
untuk langkah 8-10 induk (GRL-01 butir 5). Sesi ini memutuskan nasibnya **untuk jalur addendum saja**.

| TDA | Pokok | Sisi | Diadili sesi induk? | Nasib | Oleh |
|---|---|---|---|---|---|
| TDA-01 | nomor revisi dari baris yang dipilih; penjaga duplikat mati; tabrakan jadi `UPDATE` | 56-KHAS | — | **diperbaiki** — bahan to-spec `B-1` (nomor urut bentrok **ditolak**) + `GRL-10` (dasar diturunkan, bukan dipilih). **Sisa warisan ditunda**: baris yang sudah tertimpa tidak dipulihkan, `UA-14`/`UA-21` | ditutup ronde TDA, `ST-1` |
| TDA-02 | penolakan menghapus baris addendum | IRISAN | **ya** — daftar eskalasi induk, `5-tiket/`, ADR-0055 | **diperbaiki** — ADR-0055 perubahan 24 Sep syarat (a): baris tidak dihapus; `TOLAK` -> `DITOLAK` terminal, ditambah keadaan `DIBATALKAN` | **ditutup di GRL-08**; dicoret dari cabang D |
| TDA-03 | penghapusan tidak membersihkan detail, lampiran, penanda | IRISAN | belum sebagai butir tersendiri | **diperbaiki untuk sistem baru** — mekanismenya hilang bersama TDA-02; **sisa warisannya tidak** | **cabang I**, diukur `UA-7`; dicoret dari cabang D |
| TDA-04 | selisih dipadankan menurut posisi baris | IRISAN | disebut, belum diadili | **diperbaiki** — `E1` KONFIRMASI, `KUNCI_PADANAN` menggantikan pemadanan posisi | ditutup di E1; ditandai ronde C |
| TDA-05 | selisih dikurangkan tanpa memeriksa mata uang | IRISAN | belum | **diperbaiki** — `E-1a`, mata uang masuk kunci padanan (ADR-0053) | ditutup di E1; ditandai ronde C |
| TDA-06 | sebagian selisih ditulis ke `ActualValue` lalu ditimpa | IRISAN | sebagian — `CONTEXT.md` §1.2 | **diperbaiki** — **GRL-14**: tidak ada `ActualValue` yang dapat ditimpa. E3c larut, bukan dijawab | ditutup di GRL-14 |
| TDA-16 | selisih EGNPI pada addendum premi **selalu nol** — mesin beriterasi `TreatyIn.EGNPI`, pengguna menyunting `ActualValue.EGNPI` | IRISAN | belum | **diperbaiki** — **GRL-14**: premi aktual menjadi nilai versi, sehingga selisihnya lahir sendiri | ditutup di GRL-14 |
| TDA-07 | ~~membuka kontrak menulis kembali ke baris kontrak~~ **DICABUT** (NA-06); penggantinya: dua kontrol layar mengosongkan status akseptasi kontrak | IRISAN | belum | **diperbaiki** — `GRL-08`: `StatusAkseptasi` adalah **mekanisme**, tidak dibawa; INV-24 membuat versi terminal **beku**. Wewenang buka-kunci sudah di eskalasi induk butir 1 | ditutup ronde TDA, `ST-5` |
| TDA-08 | rantai memendek lewat `RevisionState`, jejaknya dihapus; untuk addendum nilainya **diwarisi** | IRISAN | **ya** — ADR-0052 | **diperbaiki** oleh ADR-0052, dikonfirmasi di D | cabang D |
| TDA-09 | peran dari `pyWorkBasketList(2)` / `pyTelephone` / nama tersemat | campuran | sebagian — eskalasi induk butir 5 | **KONFIRMASI ADR-0044** — peran dan penugasan **bertanggal** masuk gelombang 1; *"nama orang tidak pernah muncul di dalam aturan"*. Baris berjalan dipetakan `Position` + `StatusAkseptasi` (`GRL-08`) | ditutup ronde TDA, `ST-4` |
| **TDA-17** | **`TREATYINDETAILEDM` tertinggal TUJUH KOLOM dari `TREATYINDETAIL`** | **IRISAN** — `pzInsKey` `SaveTreatyInDetail` **sama persis** di kedua ekspor | **diadili ronde D** (`TD-05`) | **diperbaiki SEBAGIAN.** Ketimpangan tujuh kolom hilang karena bentuknya (ADR-0056). Empat kolom punya rumah di `LAYER`; **tiga — `DEDUCTIBLE2`, `PREMIUM_EARNED`, `ROL_PCT` — TIDAK PUNYA RUMAH SAMA SEKALI**, dan sebabnya `L-8`: ketiganya hanya muncul di pohon cermin. **Nol aturan membacanya** dari tabel datar | **DITUNDA** untuk ketiga kolom — pemiliknya **sesi to-spec INDUK**, bukan modul ini |
| TDA-10 | materialitas hanya ditegakkan di layar | IRISAN | belum | **🔴 KEMBALI TERBUKA** — perbaikannya bersandar pada `GRL-12`, yang **batal**. Dan jawaban `DB-15` justru menyatakan penguncian layar itu **kehendak bisnis**, bukan cacat: materialitas memang menentukan field mana yang boleh disunting | terbuka kembali; menunggu pengganti `GRL-12` |
| TDA-11 | picker menyatukan kontrak dan addendum tanpa pembeda | 56-KHAS | diparkir (B-1) | **diperbaiki untuk model** — `GRL-10` dan `B-4` mencabut peran picker (dasar **diturunkan**); INV-25 menutup ketiadaan saringan keadaan. **Sisa tampilan REKOMENDASI TO-SPEC** | ditutup ronde TDA, `ST-6` |
| TDA-12 | offset pengurai nomor revisi meleset satu; patah pada revisi kesepuluh | 56-KHAS | diparkir (A-1, G4) | **diperbaiki** — `GRL-09`, yang **menyebut TDA-12 di dalam labelnya sendiri**: *"bentuk turunan tanpa lebar tetap → PERUBAHAN (TDA-12)"* | ditutup ronde TDA, `ST-2` |
| TDA-13 | jenis dan materialitas dipilih bebas di radio picker Revisi | 56-KHAS | — | **diperbaiki** — `GRL-12` (materialitas menjadi turunan, radionya hilang) + **`GRL-18`** (jenis menjadi masukan yang **beku sejak diajukan** dan berjejak) | ditutup ronde TDA, `ST-7` |
| TDA-14 | nilai disimpan sebagai teks di sebagian tempat | IRISAN | ya — ADR-0003 | **diperbaiki** oleh ADR-0003, dikonfirmasi di E | cabang E |
| TDA-15 | penggandaan pohon di dalam satu `JSONDATA` | IRISAN | belum | **diperbaiki** — ADR-0034 (*"tidak ada kode yang membacanya"*) + keempat pohon kini punya tujuan bernama; yang terakhir, `ActualValue`, oleh **`GRL-14`** | ditutup ronde TDA, `ST-3` |

**Sepuluh dari lima belas TDA berada di irisan.**

## Koreksi atas bahan

| Tanggal | Berkas | Isi | Kadar |
|---|---|---|---|
| 24 Sep | ronde C, `02-SIDANG` SC-6 | **awalan nama bukan penanda sisi**: seluruh `TreatyEDM*` **tidak ada** di katalog 56 — mesin selisih addendum ada di ekspor induk juga. Lima dari tujuh temuan ronde C berpindah menjadi **temuan Treaty In** | terbaca |
| 24 Sep | ronde C, `02-SIDANG` SC-2 | klaim *"`EDMEffective` hanya dibaca `TreatyCalculateProratePct`"* **salah kaprah** — ia diam tentang penulis. Penulis `TreatyIn.EDMEffective` **satu**: `TreatyInSetEditPre` 1.5 → `= Commencement` | terbaca |
| 24 Sep | ronde C, `01-TEMUAN` NC-06 | **ADR-0043 dibaca** (prasyarat `MA-10`). Payung I2 turun ke **KONFIRMASI**; residunya nomor urut versi warisan, dan ADR-0043 memunculkan tabrakan kriteria identik | ditutup |
| 24 Sep | ronde C, `03-LUBANG` §4 | **tujuh TDA tanpa nasib** ditemukan dengan mencocokkan Lacak TDA terhadap tabel Ronde: TDA-01, 07, 09, 11, 12, 13, 15 menunjuk cabang yang sudah ditutup | **ronde D** |
| 23 Sep | `../treaty-in/SEAM-ADJUSTMENT.md` (kepala) | sistem lama menyimpan potret nilai lama; "DATA LAMA TIDAK PERNAH DISALIN" ditegaskan sebagai aturan sistem baru | terbaca |
| 23 Sep | `../treaty-in/4-erd-dan-tabel-datar/KEPUTUSAN-SAMBUNGAN-ADJUSTMENT.md` (kepala) | §G2 tidak bertahan; G1, G1b, G3, G4 tetap | terbaca |
| 24 Sep | `../treaty-in/CONTEXT.md` §1.1 dan tabel §4 | **perubahan artefak induk, dengan izin**: pernyataan embargo diganti rujukan GRL-01; tanggal mulai embargo dihapus (tidak ada artefak yang menyebutnya); status ADR-0048 "sedang direkonsiliasi" | — |
| 24 Sep | `../treaty-in/DAFTAR-ESKALASI-MANAJEMEN.md` butir "selisih yang pernah dilaporkan" | **perubahan artefak induk, dengan izin**: sebabnya diralat; dibatasi pada tiga keadaan; dirujuk ke `UA-5` dan `UA-6`; "dibukukan" diganti "dilaporkan" | terbaca |
| 24 Sep | `PENGETAHUAN.md` §3.2 (blok koreksi), §4.2, §4.3, TDA-13, `UA-2` | dua radio dapat disunting di picker Revisi; jenis dan materialitas dipilih di layar | mekanisme terbaca; himpunan nilai **tidak terbaca** |
| 24 Sep | `PENGETAHUAN.md` §4.2 (judul), §2, §5.6 (baru), §10.4 | 33 akar / 165 jalur daun, **beserta daftarnya** (`METODE` §3.9) | terbaca |
| 24 Sep | `PENGETAHUAN.md` §4.4, TDA-12, `UA-1` | sebab patahnya penomoran **diralat**: offset meleset satu, bukan dua pengurai; dua akhir yang mungkin | sebagian tidak terbaca |
| 24 Sep | `PENGETAHUAN.md` §7.2 -> §7.4 | klaim *"bagaimana `RevisionState` sampai di addendum tidak dapat dijawab ekspor"* **dicabut**: daur hidupnya terbaca lengkap, penolakan **tidak** mengosongkannya, dan warisannya juga mengunci addendum | ditutup |
| 24 Sep | `PENGETAHUAN.md` §4.3, §4.6, §6.2, §7.4 | **audit `TA-04`**: enam klaim negatif disapu ulang dengan perkakas terkalibrasi. Tiga berdiri; penulis `EDMMaterialType` 1 -> 3, penulis `CedingID` 2 -> 4, kontrol `Ceding` di layar addendum 2 -> 1; daftar penulis `RevisionState` bertambah `TreatyInCopy` | kesimpulan pokok tidak berubah |
| 25 Sep | `PENGETAHUAN.md` §6.2, §7.5, `GRILL-A/01-TEMUAN` NA-14 | **titik buta audit ditutup**: 32 langkah `Java` disapu sebagai teks (nol), pemetaan respons Connect-REST diperiksa (ke `UploadDoc`, bukan `TreatyIn`), 191 langkah SQL dicirikan. Kadar klaim naik ke **delapan jenis penulis** | ditutup |
| 26 Sep | `PENGETAHUAN.md` §7.7 (baru) | **NA-21** — `Revision` juga mengunci: 96 sel uang terjaga `ViewState`, 882 tidak; lapisan beku termasuk yang tidak | `DB-10` ditulis ulang, `DB-10b` ditambahkan |
| 26 Sep | `PENGETAHUAN.md` §7.7 | bacaan "kunci uang" **gugur** sesudah dicocokkan ke 33 akar §5.6: tidak satu pun dari dua belas akar uang terkunci seluruhnya; `LayersEDM` tidak menjaga apa yang dijaga `Layers` | `NA-22`; `DB-10` ditulis ulang kedua kalinya; `DB-10b` dicabut |
| 26 Sep | deskripsi `EXP-1` di log | **keliru**: disebut "Field Value, kelas `Data-Portal`". Yang benar **`Rule-Obj-Property`** di kelas **`ASM-FW-GISFW-Int-TREATY_IN`**, ruleset GISFW 01-01-56 | diralat; `EXP-1` ditutup |
| 27 Sep | `PENGETAHUAN.md` §7.7 | **penanda nonaktif ditafsirkan**: `pyDisabledNew = always` (146) selalu mati, `= true` + `pyDisabledWhen` (249) bersyarat, `= false` aktif; 146+249 = 395 tepat | `NB-05` |
| 27 Sep | `PENGETAHUAN.md` §7.7, `DB-10`, GRL-08 butir iii | **daftar bocor menyusut**: total angsuran, `LayersEDM`, `DetailEGNPI` ternyata **nonaktif permanen**. Kalimat *"`Layers` menjaga, `LayersEDM` tidak"* **terbalik dan dicabut** — `LayersEDM` justru lebih ketat | `NB-06`; kesimpulan `DB-10` tidak berubah |
| 27 Sep | E1 | syarat penurunan diperiksa: `SEAM-ADJUSTMENT` §3 **keputusan diterima**, premisnya tidak bersandar G2, keunikan kunci menjadi `UA-18` | `NB-07`; E1 tetap KONFIRMASI |
| 27 Sep | `PENGETAHUAN.md` §5.4, §9.1 | **TDA-16 baru** — selisih EGNPI pada addendum premi **selalu nol**: mesin premi beriterasi `TreatyIn.EGNPI`, pengguna menyunting `ActualValue.EGNPI` | `NB-08`; menjadi syarat `C-3` di C3 |
| 27 Sep | `PEMILAHAN-SISA-GRILLING.md` | titipan **"mekanisme pembekuan"** ternyata **tanpa golongan** di tabel 43; kini KONFIRMASI — `NILAI_SELISIH` berkunci alami *besaran + kunci padanan baris* | ditutup |
| 27 Sep | §7.7 | tiga belas sel berpindah golongan: syarat `ViewState` tersimpan pada mode yang **tidak diaktifkan** | `NB-09`; tidak ada kesimpulan berubah |
| 26 Sep | `PENGETAHUAN.md` §4.2, §4.3, TDA-13, `UA-2` | `EDMState = 2` **ditawarkan layar** — kadar naik dari *masih bisa terjadi* menjadi **terbaca**; radio tidak menawarkan `3`, sehingga "penyesuaian premi selalu material" berlaku untuk kedua jalur | blok koreksi bertanggal |
| 26 Sep | `PENGETAHUAN.md` §7.7 | dugaan `LayersEDM` hanya di layar addendum **dibantah** — ia ada di kedua pohon; pemisahan pohon tidak mengubah `DB-10`, tetapi mengecilkan GRL-08 butir iii | ditutup |
| 26 Sep | `UA-16` | zona waktu kedua sisi **berbeda** — `EDMDATE` jam Oracle, komentar persetujuan GMT; kuerinya harus menormalkan, zona server **perlu dicek** ke DBA | ditutup |
| 26 Sep | berkas `EXP-1` | ditemukan di `C:\Users\Administrator\Downloads\`, bukan di folder proyek — **pengulangan `PG-01`**; disalin ke `ekspor-tambahan/` | selesai |
| 26 Sep | `UA-16` | pembanding ditetapkan: `EDMDATE > MAX(CommentList.Date) + 5 menit`, sebab persetujuan sendiri menyimpan baris (`TreatyInSubmitEDM` langkah 8) | ditutup |
| 26 Sep | `PEMILAHAN-SISA-GRILLING.md` | seluruh **43** butir anggaran asli digolongkan; anggaran baru **7 butir GRILL** | diterapkan |
| 26 Sep | `GRILL-A/01-TEMUAN` NA-17a | kadar diturunkan: `'ALDO SAPUTRA1'` **bukan** sekeluarga `1=2` — ia pertanyaan data (`UA-15`), bukan kepalsuan logika | ditutup |
| 26 Sep | bukti B2 | **dicabut**: pengenal bersusun tidak pernah terbentuk; `HASIL1` dibuang dan baris addendum memenuhi `LIKE`-nya sendiri | `MA-09`, `NB-01`, `UA-1(c)` diganti |
| 26 Sep | `USULAN-REVISI-ADR.md` `REV-1` | draf kedua keliru menyatakan premis SELECT dicabut — **GRL-03 butir 2 mempertahankannya**; ditulis ulang sebagai koreksi pemerian, bukan penalaran | draf ketiga |
| 26 Sep | `../treaty-in/DAFTAR-ESKALASI-MANAJEMEN.md` butir 1 | **disimpan** — Force Edit dibatasi divisi TI; satu tombol bukan dua; label terbaca; temuan ketiga pada layar addendum beserta batas pemeriksaan dan `UA-16`; pertanyaan wewenang | selesai |
| 25 Sep | `GRILL-A/06-PUTUSAN` GRL-05, `01-TEMUAN` | hitungan lapisan beku diralat **2+2+2 -> 1+1+1**; kesimpulannya tidak berubah | ditutup |
| 25 Sep | `GRILL-A/06-PUTUSAN` MA-07 | ralat MA-07 sendiri **keliru ke arah sebaliknya**: wadah dibaca menggantikan sel. Keterlihatan = **wadah × sel** | `TA-06` diperluas |
| 25 Sep | `PENGETAHUAN.md` §4.6 | **NA-17** — `ShowSummary` membuka kelima field lapisan beku sebagai teks bebas; tidak disebut dokumen induk mana pun | usulan untuk induk, GRL-01 butir 5 |
| 25 Sep | daftar eskalasi induk butir 1 | **NA-18** — keempat kontrol dibatasi peran inputor; hanya `Revision` menyentuh kontrak yang sudah disetujui; labelnya terbaca | diff disiapkan, **menunggu persetujuan** |
| 24 Sep | `PENGETAHUAN.md` §7.5 | klaim *"`Force Edit (dev)` terlihat semua orang"* **dicabut**: ada gerbang `pyOrgDivision = 'IT'` satu lapis di atasnya | `MA-07`, `REV-3` butir iii |
| 24 Sep | `PENGETAHUAN.md` §7.6 | kadar naik dari *disimpulkan* menjadi **terbaca**: tag pemilih sumber grid `pySourceType = Property` dikutip beserta nomor barisnya | ditutup |
| 24 Sep | laporan lisan A7 | `Section/TreatyInActionButtons` disebut **penulis** peran; ia **pembaca**. Cabang 1.4 `Akseptasi_DT` disebut mati "blok `//`"; penandanya **`pyDisabled = true`** | `MA-05`, `TA-04` |
| 24 Sep | `PENGETAHUAN.md` §7.6 | dugaan bahwa daftar pilihan menyaring `Resolve Complete` **dibantah berkasnya sendiri**; TDA-11 dikuatkan dengan sumber yang benar | ditutup, nasibnya ke cabang B |
| 24 Sep | `PENGETAHUAN.md` §1.2 | irisan **317** ber-`pzInsKey` sama dari 323 jalur; enam jalur memuat versi aturan berbeda; tidak satu pun menjadi bukti TDA | terbaca |
| 24 Sep | `PENGETAHUAN.md` §4.5, §7.2, TDA-07, TDA-08, `UA-9` | **TDA-07 dicabut** — jalur Adjustment tidak menulis ke baris kontrak; klaim panjang rantai ikut dicabut | terbaca |
| 24 Sep | `METODE-GRILLING.md` Bagian VIII | disidangkan di `GRILL-A/02-SIDANG.md`: §8.1 salah kaprah **dua kali**, §8.3 "tabrakan mustahil" gugur sebagai fakta, §8.4 (OLDID, R01..R99) salah kaprah, ROWNUM diperkecil, §8.6 Q2 ditutup | terbaca |

## Ekspor tambahan yang dibutuhkan

| # | Ekspor | Dibutuhkan untuk | Keadaan |
|---|---|---|---|
| ~~**EXP-1**~~ | ~~Aturan **Field Value** `EDMState` dan `EDMMaterialType` (kelas `Data-Portal`)~~ **DITUTUP 26 Sep 2026.** Bentuknya ternyata `Rule-Obj-Property` di kelas `ASM-FW-GISFW-Int-TREATY_IN`, bukan Field Value di `Data-Portal`; berkasnya di `ekspor-tambahan/` | himpunan jenis dan materialitas yang benar-benar ditawarkan layar | **terbuka** — memblokir penutupan TDA-13 dan cabang C |
| ~~EXP-2~~ | ~~`METODE-GRILLING.md`~~ | — | **DITUTUP 24 Sep.** Berkasnya ada di folder Downloads pengguna, bukan di direktori kerja (`PG-01`). **Sudah disalin** ke `../METODE-GRILLING.md`; seluruh rujukan §-nya diperiksa langsung dan cocok, kecuali §8.1 yang memang salah |

## Daftar eskalasi — keputusan orang berwenang, bukan pekerjaan teknis

Dipisahkan dari pekerjaan teknis sesuai `METODE` §6.6.

| # | Butir | Keadaan |
|---|---|---|
| ~~ESK-1~~ | **TIDAK JADI BUTIR BARU.** Butir 1 daftar eskalasi induk sudah memuat tombol pembuka kunci kontrak yang disetujui; ADR-0052 + ADR-0055 sudah memutuskan nasibnya. Yang keliru hanyalah satu kalimat peredam — **koreksinya sudah disimpan** ke butir 1 pada 24 Sep, lihat `GRILL-A/07-AUDIT.md` §3a | **selesai** |
| **ESK-2** | Penolakan addendum menghapus barisnya | **sudah** ada di daftar eskalasi induk (tambahan 24 Sep) — tidak ditulis ulang |

## Pertanyaan yang belum terjawab

| # | Isi | Memblokir |
|---|---|---|
| `PP-1` | pemilik pengerjaan to-spec induk (tim sama / tim lain) | tidak — hanya arah dampak GRL-01 |
| [`2-to-spec/USULAN-DIFF-KE-INDUK.md`](2-to-spec/USULAN-DIFF-KE-INDUK.md) | **tujuh diff ke artefak induk — disetujui 24 Sep 2026** | **lima DITERAPKAN** (`D-2`…`D-6`); **dua DIPARKIR**: `D-1` menunggu `F-2`, `D-7` menunggu izin |
| [`KEPUTUSAN-TANPA-VERIFIKASI-ADJUSTMENT.md`](KEPUTUSAN-TANPA-VERIFIKASI-ADJUSTMENT.md) | **4 butir** `KTV-1`…`KTV-4`, masing-masing berdasar dan bersyarat pembalikan | **nol terverifikasi**; ketiganya disempitkan `DB-20`, `DB-16a`, `DB-16b`; `KTV-4` oleh tanggapan pemilik ADR |
| `DB-1` … `DB-19` (`DB-10b` **dicabut**; `DB-18`, `DB-19` lahir ronde C; `DB-16` menjadi pembatal `GRL-18`) | daftar untuk dibantah ke bisnis — bunyinya di `GRILL-A/07-AUDIT.md` | `DB-3`/`DB-4` pembatal GRL-04; **`DB-10`** satu-satunya yang menanyakan kehendak bisnis di balik keputusan ADR yang sudah diambil |
| `UA-1` … `UA-21` | kueri data — bunyinya di `PENGETAHUAN.md` §9.2, `UA-19`…`UA-21` di `GRILL-C/03-LUBANG.md` §1 | **tidak satu pun memblokir rancangan** (`METODE` §7.2). `UA-3` dan `UA-12(b)(c)` menyebut batas pembuktiannya sendiri: dapat mematahkan, tidak dapat mengesahkan |
| `REV-1` … `REV-6` | **paket draf** di [`USULAN-REVISI-ADR.md`](USULAN-REVISI-ADR.md) — ADR induk **tidak disunting**; diserahkan sekaligus di akhir grilling, kecuali bila `PP-1` menunjukkan to-spec induk sedang berjalan | prasyarat cabang K: *paket REV diserahkan dan ditanggapi* |
| `REV-3` | usulan koreksi ADR-0055 §4, empat butir: (i) `TreatyInSetToDirector` seluruh langkahnya mati; (ii) pintu samping ada di layar addendum juga; (iii) `TreatyInForceEdit` terlihat operator **divisi IT**, bukan setiap pengguna; (iv) tombol paksa di layar addendum tidak menulis apa pun sambil melapor berhasil | tidak memblokir; keputusannya milik pemilik proses |
| `REV-2` | usulan koreksi **Konteks** ADR-0052: penyimpangan yang berjalan **satu**, bukan dua; jalur penerima tugas kelompok adalah kode untuk keadaan yang tidak dapat dicapai | tidak memblokir; ADR-nya sendiri tetap berlaku |
| `TA-01`, `TA-02`, **`TA-09`** | usulan revisi `METODE` §2.0, §8.1, dan — baru — **pencocokan dua arah antar tabel di dalam satu berkas** pada tiap titik periksa (`GRILL-TDA/03-LUBANG.md` §3) | tidak |

---

## Sumber induk per butir GRILL yang tersisa

Disusun 26 September 2026 dari [`PETA-SUMBER-INDUK.md`](PETA-SUMBER-INDUK.md), **sebelum** butir
berikutnya diajukan — pelajaran `MA-10`. Tanda ⚠ berarti berkas itu ditulis sebelum koreksi
`PENGETAHUAN.md` Adjustment dan **dikutip lalu diperiksa**, bukan diwarisi.

| Butir | Berkas induk yang menyentuhnya | Sudah menjawab? |
|---|---|---|
| ~~**E1** kunci padanan~~ | `SEAM-ADJUSTMENT.md` §3 ⚠ — `KUNCI_PADANAN`, Bentuk A dan B, ruas tambahan `POTONGAN`; bersandar pada `SPEC-INVARIAN.md` §4.1 | **YA — dipindah ke KONFIRMASI**; sisanya bahan to-spec `E-1` (mata uang di dalam kunci). **Anggaran 7 -> 6** |
| **C2** daftar jenis | ADR-0049 (satu sumbu, tiga jenis); `SPEC-MODEL-DATA.md` §10.2 `JENIS_ADDENDUM`; `CONTEXT.md` glosarium | **sebagian** — ADR-0049 menetapkan **satu sumbu tiga jenis**, tetapi sumber turunan materialitasnya baru diubah GRL-12, sehingga isi sumbu itu perlu diputuskan ulang |
| **C3** arti `ActualValue` | `PETA-TELUSUR-JSON.md` ⚠ §1 — *"272 jalur cermin… `ActualValue.*` (108) memakai entitas yang sama, bukan entitas baru"* | **sebagian** — ia memutuskan `ActualValue` **tidak melahirkan entitas baru**, tetapi **tidak** memutuskan apakah ia disimpan, diturunkan, atau dibuang |
| **E3** kemampuan mati | `5-tiket/DAFTAR-PEKERJAAN.md`; `5-tiket/LUBANG-SPESIFIKASI.md`; ADR-0037 (pro rata = turunan dari tanggal) | **sebagian** — ADR-0037 menetapkan *"faktor prorata addendum = turunan dari tanggal, bukan centang"*, yang menyentuh **E3a** tetapi bukan E3b maupun E3c |
| **I2** aturan baru atas warisan | ADR-0042 (sejarah pindah apa adanya); ADR-0043 (migrasi memindahkan, hitung ulang peristiwa bisnis); ADR-0054 (keadaan warisan tak terpetakan) | **sebagian** — ketiganya menetapkan prinsipnya; yang belum ada adalah daftar aturan mana yang boleh dijalankan dan bentuk pencatatan asal-usulnya |

**ADR-0043 belum pernah saya baca** dan namanya menjanjikan jawaban bagi I2. Ia dibaca **sebelum**
I2 diajukan, bukan sesudah.
