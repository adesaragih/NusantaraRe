# OPEN QUESTIONS — Claim Non Prop

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, berkas tertanggal 2026-09-08 s/d 2026-09-09, rule termutakhir di dalamnya `pxUpdateDateTime = 2026-08-30`.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut. Tambalan yang ditambahkan sesudahnya tidak tercermin di sini; deteksinya lewat sapuan ulang, bukan lewat register.

> **DIPINDAHKAN KE `_selesai/` 19 September 2026 — dan `_selesai/` di sini TIDAK berarti "tidak berlaku lagi".**
>
> **Yang tertutup registernya, bukan isinya.** Tidak ada lagi pertanyaan yang menahan pekerjaan; itu sebabnya ia pindah. Tetapi **delapan cabang di lampiran masih mengikat rancangan**, dan `KEADAAN-AKHIR.md` §3.1 menunjuk ke sini sebagai tempat kedelapannya hidup.
>
> | Bagian | Kedudukannya |
> |---|---|
> | Kepala sampai `# SELESAI` | **keadaan** — tertutup, tidak akan berubah kecuali sebuah cabang gugur |
> | `# SELESAI` | **arsip** — bukti bahwa tiap jawaban punya dasar |
> | **`# Lampiran`** | **BAHAN HIDUP** — empat hipotesis, delapan cabang, nol dipilih. Rancangan mana pun yang **hanya benar bila salah satu cabang benar** melanggar batas ini |
>
> Dua aturan di kepala berkas ini — **aturan register** dan **aturan pola sapuan** — juga tetap berlaku penuh untuk pekerjaan berikutnya. Butir baru dari sapuan berikutnya masuk ke sini, bukan ke laporan.

**Pertanyaan: aktif 0 · selesai 57 · dihapus 2 · mati 0**

> **Perintah yang menghasilkan angka itu** — dijalankan atas berkas ini sesudah penutupan, sebab angka tanpa perintah tidak dapat diperiksa ulang:
>
> ```
> sed -n '/^# SELESAI$/,/^# Lampiran/p' OPEN-QUESTIONS.md | grep -oE '^\| ~*[ABCEGK][0-9]+[ab]?~*' | tr -d '|~ ' | sort -u | wc -l
> ```
>
> Hasilnya **57 kode berhuruf**. Dari situ dikurangi **A14** (catatan perpindahan nomor, bukan pertanyaan) dan **A10b**, **A11b** (dihapus, bukan ditutup) → **54**. Ditambah **F**, **H**, **I** yang tidak berkode angka → **57 pertanyaan tertutup**.
>
> Pecahannya: **17** sudah tertutup sebelum hari ini (A4, A5, A6a, A9, A10a, A11a, B3, B5, C2, C9, E3, E4, E6, E7, E8, E9, E10) dan **40** ditutup 19 September 2026. Dua lagi **dihapus**, bukan ditutup — A10b dan A11b — sehingga 40 + 2 = 42, yaitu seluruh butir yang tadinya terbuka.
>
> **Pencacah lama berbunyi "selesai 26" dan tidak dapat direproduksi dengan perintah mana pun** — daftar SELESAI saat itu memuat 17 butir. Angka 26 **dicabut**, bukan diwariskan.

> Pencacah ini menghitung **pertanyaan**. REQ di `ORACLE-REQUESTS.md` dan tiket di papan `issues/README.md` punya pencacahnya sendiri dengan format yang sama; bendanya diberi label supaya dua pencacah berformat sama tidak pernah tertukar.

---

## REGISTER TERTUTUP — 19 September 2026

**Tidak ada satu pun butir di berkas ini yang menahan keputusan rancangan mana pun.** Butir terakhir yang menahan bukan pertanyaan melainkan REQ-032, dan ia turun menjadi verifikasi pada tanggal yang sama.

### Apa yang dapat membukanya kembali

Register yang tertutup bukan register yang dikunci. Dua hal membukanya, dan keduanya datang dari luar:

| Pemicu | Yang terjadi |
|---|---|
| **Jawaban DBA yang menggugurkan sebuah ramalan** | Tujuh REQ kini berstatus verifikasi, dan masing-masing membawa ramalan yang dapat dibantah data. Satu jawaban yang tidak cocok melahirkan butir baru di sini — bukan di laporan |
| **Sesi Komite, bila kelak terjadi** | F, H, dan I ditutup sebagai hipotesis permanen. Sesi itu membukanya kembali **sebagai verifikasi**, bukan sebagai penahan, karena rancangannya sudah benar di kedua cabang |

### Register yang kosong bukan pengetahuan yang lengkap

Kalimat ini ditulis karena ia yang paling mudah hilang.

**Yang tertutup adalah pertanyaan yang menahan pekerjaan, bukan pertanyaan yang ada.** Tiga hal tetap berdiri, dan ketiganya tercatat sebagai **batas** — bukan sebagai lubang yang terlupa:

| | Apa adanya |
|---|---|
| **Empat hipotesis permanen, delapan cabang** | F (`.IsEditClaim`), H (`.ValueAdjustment`), I (transisi `.AcceptanceStatus`), G1 (syariah). Cabangnya tertulis utuh di lampiran; **tidak satu pun dipilih**. *Angka "tujuh" pada versi sebelumnya dicabut — ia mencampur satuan hipotesis dengan satuan cabang* |
| **Tiga puluh lima permintaan belum dijawab** | Seluruh isi `PERMINTAAN-DBA-2-PROFIL-DAN-STRUKTUR.md`, ditambah REQ-032 di berkas pertama. Tidak ada yang menahan rancangan; empat menahan **pelaksanaan** tiket |
| **Folder `Komite Claim Non Prop` tidak dibuka** | Tertutup untuk pekerjaan ini — bukan ditunda. Apa pun yang hidup di sana tidak pernah masuk ke sini, dan itu diketahui, bukan diabaikan |

Sebuah daftar pertanyaan yang kosong **tidak berarti kita tahu segalanya**. Ia berarti setiap hal yang belum diketahui sudah ditempatkan: di batas kepemilikan, di cabang yang tertulis, atau di permintaan yang menunggu jawaban.

---

### Bagaimana butir ditutup — 19 September 2026

**Seluruh bagian A, B, C, E, F, G, H, I, dan K ditutup pada tanggal ini.** Empat puluh butir dipindahkan ke SELESAI; dua dihapus. Tidak satu pun dibiarkan menggantung.

**Bagaimana butir ditutup di sini.** Sebagian besar **tidak ditutup dengan jawaban**. Sebuah butir juga selesai bila:

| Cara tutup | Artinya | Contoh |
|---|---|---|
| **TIDAK BERLAKU** | pertanyaannya mengandaikan sesuatu yang tidak ikut pindah | A1 workbasket, A12 alamat `prweb` |
| **DI LUAR KEPEMILIKAN** | ia perilaku modul lain; kita tidak merancangnya dan tidak perlu mengetahuinya | C1, C3, B6 |
| **DITUTUP OLEH BENTUK** | rancangannya benar di setiap jawaban yang mungkin, jadi jawabannya berhenti menentukan | B4, C6, G1, F/H/I |
| **DIPUTUSKAN** | keputusan diambil, dan dicatat di berkas yang akan memakainya | B1/B2, A16, A3 |
| **DIHAPUSKAN** | yang ditutup bukan pertanyaannya melainkan hal yang membuatnya ada | A13 — pengecualiannya dihapus, jadi tidak ada yang perlu berwenang |

**Yang tidak dilakukan**: menutup butir dengan menyebutnya tidak penting. Setiap baris di bawah menyebut **tanggal**, **sebab**, dan **berkas tempat keputusannya hidup** — karena register menyimpan keadaan, dan berkas yang memakainya menyimpan bentuknya.

**Koreksi pencacah**: kepala versi sebelumnya berbunyi *"aktif 41"*, dan yang benar **42** — E17 masuk 19 September sesudah angka itu ditulis dan tidak ikut dinaikkan. Keempat puluh dua itulah yang ditutup di sini: 40 ke SELESAI, 2 dihapus.

### Yang masih berjalan, dan tidak dihitung di register ini

Tidak ada pertanyaan terbuka. Yang berjalan adalah **verifikasi**, dan benda itu hidup di register lain:

| Benda | Tempatnya | Kedudukan |
|---|---|---|
| REQ-001, REQ-009, REQ-010, REQ-011, REQ-019, REQ-028, REQ-034 | `ORACLE-REQUESTS.md` | turun jadi **verifikasi** 19 Sep 2026 — menguatkan atau menggugurkan, **tidak menahan** |
| REQ-004 | `ORACLE-REQUESTS.md` | **dicabut** (`DEAD`) — hanya melayani B6, dan B6 di luar kepemilikan |
| **REQ-032** | `PERMINTAAN-DBA-1-VERSI-INSTANCE.md` | **VERIFIKASI** sejak 19 Sep 2026 — batas 30 byte ditetapkan sebagai aturan tetap, jadi jawabannya tidak mengubah satu nama pun. **Tidak ada lagi REQ yang menahan keputusan rancangan**; empat menahan **pelaksanaan** tiket: REQ-018, REQ-033, REQ-021, REQ-037 |
| REQ-037 | `ORACLE-REQUESTS.md` | BLOCKER untuk klaim ADR-0017, bukan untuk model data |
| `ASK-AKUNTANSI` butir 5 dan 7 | `_selesai/ASK-AKUNTANSI.md` | verifikasi tarif dan padanan rupiah — **tidak menahan** |

---

## Aturan register — berlaku 18 September 2026

**Sebuah sapuan belum selesai sampai setiap butir di sini yang disentuhnya diperbarui di tempatnya.** Berkas laporan menyimpan **bukti**; register ini menyimpan **keadaan**. Laporan tidak pernah menggantikan register.

Setiap butir yang ditutup membawa **tanggal**, **sapuan atau keputusan yang menutupnya**, dan **berkas beserta barisnya**. Butir yang ditutup **dipindahkan ke SELESAI, tidak dihapus**.

> Aturan ini **tetap berlaku sesudah penutupan 19 September**. Register yang kosong bukan register yang mati: butir baru lahir dari sapuan berikutnya, dan ia masuk ke sini, bukan ke laporan.

## Aturan pola sapuan — berlaku 18 September 2026

Sudah dua kali sapuan literal melewatkan sesuatu karena nilainya tidak pernah muncul dalam bentuk yang dicari:

| Kejadian | Yang terlewat | Sebab |
|---|---|---|
| **D9** | `PaymentType` = `1`, `2`, `4`, `5`, `6` | polanya hanya mencakup literal **berkutip** |
| **C9** | `TreatyName` = `"Previously Calculated UR"` | nilainya **dirakit**, tidak pernah muncul utuh sebagai literal |

Karena itu setiap pola sapuan wajib mencakup **keempat bentuk**:

1. berkutip ganda `"…"` (termasuk `&quot;…&quot;`),
2. berkutip tunggal `'…'`,
3. **tanpa kutip** — perbandingan numerik,
4. **nilai yang dirakit** — penyambungan `"teks" + .Properti`, lalu diuji utuh di tempat lain.

**Pola yang dipakai ditulis di badan laporan** supaya dapat diperiksa ulang.

---

## A. Butuh data atau keputusan di luar XML — **KOSONG**

Sepuluh butir ditutup 19 September 2026, dua dihapus. Enam butir lain sudah berpindah ke `ORACLE-REQUESTS.md` pada 18 September (A4, A5, A6a, A9, A10a, A11a). Daftarnya di SELESAI.

## B. PARKIR — **KOSONG**

Tujuh butir ditutup 19 September 2026. Bagian ini semula menahan frontier; ia tidak lagi menahannya.

## C. DEFERRED-TO-KOMITE-SESSION — **KOSONG**

Tujuh butir ditutup 19 September 2026. **Tidak satu pun ditutup dengan jawaban tentang isi modul Komite** — folder itu tertutup untuk pekerjaan ini. Empat ditutup sebagai *di luar kepemilikan*; tiga ditutup dengan **batas**, yaitu syarat penerimaan di `BLUEPRINT.md` §8.7.

## E. Anomali — **KOSONG**

Sepuluh butir ditutup 19 September 2026. Satu di antaranya (E14) naik menjadi `FINDING-008`.

## F, H, I. Hipotesis sejajar — **KOSONG, ditutup sebagai HIPOTESIS PERMANEN**

Ketiganya ditutup 19 September 2026. **Yang ditutup adalah status pertanyaannya, bukan cabangnya** — cabangnya tetap tertulis utuh di lampiran 1.

## G. Pertanyaan terbuka & asumsi bernama — **KOSONG**

Dua butir ditutup 19 September 2026. Badan G1 dan G2 tetap tertulis utuh di lampiran 2, karena keduanya memuat bukti yang tidak tercatat di tempat lain.

## K. Kandidat revisi ADR — **KOSONG**

Satu butir, ditutup 19 September 2026 sebagai **K1a**. Lihat lampiran 3.

---

# SELESAI

## Ditutup 19 September 2026 — bagian A

| # | Isi | Cara tutup | Ke mana keputusannya pergi |
|---|---|---|---|
| ~~A1~~ | ~~Nama workbasket tujuan assignment Akseptasi~~ | **TIDAK BERLAKU** | Workbasket adalah benda Pega. **ADR-0006** menetapkan RBAC dirancang dari nol; nama workbasket lama tidak diwarisi, jadi tidak ada yang perlu ditanyakan. Yang dibutuhkan model data — tempat bagi penugasan dan keadaannya — dirancang dari kebutuhan, bukan dari `pc_assign_*` |
| ~~A2~~ | ~~Nama peran / access group sebenarnya~~ | **TIDAK BERLAKU** | Sama, **ADR-0006**. Peran didefinisikan dari apa yang perlu dilindungi di sistem baru, bukan diterjemahkan dari yang lama. `pyPrivilegeName` kosong di 279 berkas justru berarti **tidak ada yang hilang** |
| ~~A3~~ | ~~Boleh/tidaknya satu orang melakukan Registrasi sekaligus Akseptasi~~ | **DIPUTUSKAN** sebagai kemampuan, bukan kebijakan | `SPEC-MODEL-DATA.md` §21.4. Model **wajib mampu** menegakkan pemisahan tugas — pelaku tiap tindakan tercatat terpisah per tindakan, sehingga pertanyaan *"orang yang sama?"* selalu dapat dijawab dari data. Dilarang atau diizinkan adalah **setelan, bukan struktur**: sediakan setelannya, jangan tanam jawabannya |
| ~~A6b~~ | ~~Konfirmasi label tiap nilai `PaymentType`~~ | **DITUTUP OLEH BENTUK** | `SPEC-MODEL-DATA.md` §21.2. Domain **tidak** ditutup pada tujuh nilai; nilai disimpan apa adanya termasuk kosong; **pemetaan nilai ke perilaku disimpan sebagai data** — baris yang menyatakan nilai mana menambahkan besaran mana (ADR-0004). Perilaku dari S3 dan S5 mengisi barisnya sekarang; labelnya dapat diisi kapan saja tanpa mengubah struktur. REQ-027 memverifikasi |
| ~~A7~~ | ~~Apakah `CloseClaimMD` memang boleh menutup klaim tanpa Komite~~ | **TIDAK BERLAKU** | Lihat B9. Jalur pintasnya **tidak dimigrasi**, jadi pertanyaan apakah ia boleh tidak perlu dijawab. `SPEC-MODEL-DATA.md` §21.6 |
| ~~A8~~ | ~~Apakah tarif brokerage 2,5% / PPh 2% / PPN 2,2% masih berlaku~~ | **DITUTUP OLEH ADR-0025** | Tarif menjadi **data bertanggal berlaku**. Nilai yang ada di kode hari ini masuk sebagai baris yang berlaku sampai ada penggantinya — mengubahnya jadi pekerjaan entri, bukan pekerjaan kode. `ASK-AKUNTANSI` butir 5 tetap berdiri sebagai **verifikasi nilai**, bukan sebagai penahan |
| ~~A12~~ | ~~Alamat `prweb` mana yang benar untuk produksi~~ | **TIDAK BERLAKU** | Keduanya alamat Pega; sistem baru tidak punya `prweb`. Kedua alamat dicatat di `BLUEPRINT.md` sebagai bagian gambaran sistem lama. Menutup **E11** bersamaan |
| ~~A13~~ | ~~Siapa berwenang menyetujui pengecualian pembekuan tambalan~~ | **DIHAPUSKAN** — pengecualiannya yang hilang, bukan pertanyaannya yang terjawab | **ADR-0012**, bagian *Pengecualian dihapuskan*. Selama pekerjaan ini berjalan pembekuan berlaku **mutlak**; tidak ada pengecualian, jadi tidak ada yang perlu berwenang menyetujuinya. Register pengecualian tetap ada dan tetap **kosong** — bila kelak ada yang mengisinya, saat itulah wewenangnya ditetapkan |
| ~~A15~~ | ~~Apakah ada node/instance Pega terpisah untuk bisnis syariah~~ | **DISERAP KE G1** | Lihat G1. Tidak lagi butir tersendiri |
| ~~A16~~ | ~~Apakah nilai rupiah dari biaya penilaian dan salvage pernah dibutuhkan~~ | **DIPUTUSKAN seragam** | **ADR-0029** berlaku untuk **seluruh** nilai uang tanpa pengecualian — nilai asli, kode mata uang, nilai IDR, kurs yang dipakai, asal-usul kurs — termasuk kesembilan nama yang hari ini tanpa padanan IDR. Alasannya perbandingan biaya: kolom kosong tidak merugikan siapa pun; ketiadaan kolom yang ternyata dibutuhkan berbiaya **migrasi kedua**. Cabang A dan B **dicabut dari badan ADR-0029**. `ASK-AKUNTANSI` butir 7 tetap hidup sebagai verifikasi |

## Ditutup 19 September 2026 — bagian B

| # | Isi | Cara tutup | Ke mana keputusannya pergi |
|---|---|---|---|
| ~~B1~~ | ~~Presisi dan skala final untuk setiap kolom nilai~~ | **DIPUTUSKAN**, tanpa menunggu REQ-001 | **ADR-0003**, bagian *Penutupan B1, B2, E1, E2, E12*. Skala kanonik **20** sepanjang rantai — skala tertinggi yang benar-benar dipakai sistem lama, jadi memakainya **tidak dapat kehilangan apa pun yang pernah ada**. Pembulatan **hanya di tepi** |
| ~~B2~~ | ~~Presisi tidak dijaga di tiga (kini empat) lapis sekaligus~~ | **DIPUTUSKAN** | Sama. Lapisan keempat — Java: **nol `double`, nol `float`, nol `BigDecimal`**, uang melintas sebagai `String` 92 kali (**D39**) — adalah yang menutupnya: *tidak ada presisi yang dapat dipulihkan, karena tidak pernah ada presisi yang ditegakkan.* Migrasi **mengurai teks apa adanya**; yang tidak terurai tidak dibulatkan dan tidak dibuang, ia masuk jalur tiket `14` |
| ~~B4~~ | ~~Struktur JSON di `json_klaim.DATA_JSON` dan `OS_AKSEPTASI_KLAIM.DATA_JSON`~~ | **DITUTUP OLEH BENTUK**, bukan oleh isi | `SPEC-MODEL-DATA.md` §21.8. Justru karena **D43** — `adoptJSONObject` mengadopsi teks tanpa memeriksa bentuk, sehingga isinya **tidak dibatasi** — pertanyaannya tidak dapat dijawab dengan daftar medan, hanya dengan wadah: mendarat utuh (tiket `13`), yang dikenali naik, sisanya tercatat sebagai tak terurai beserta asalnya (tiket `14`). REQ-009 memverifikasi |
| ~~B6~~ | ~~Apakah `EMAILKOMITE.DEGREE` menentukan urutan jenjang, dan apa arti `LIMIT_TOP`~~ | **DI LUAR KEPEMILIKAN** | `EMAILKOMITE` milik **modul Komite** (**ADR-0026**). Tidak dimigrasi, tidak dirancang, tidak ditanyakan. **REQ-004 dicabut** — ia hanya melayani butir ini |
| ~~B7~~ | ~~Kardinalitas riil tiap PageList untuk menentukan index~~ | **DITUNDA DENGAN SENGAJA** — pekerjaan pasca-migrasi, bukan pertanyaan | `SPEC-MODEL-DATA.md` §21.7. Index ditentukan dari **pola akses nyata** setelah skema berdiri dan data masuk, bukan dari cacah baris sebelum ada kueri. REQ-010 turun jadi PENTING |
| ~~B8~~ | ~~Apakah FINDING-001 nyata atau gugur~~ | **TIDAK LAGI MENAHAN** | Rancangan memenuhi **H1 dan H2 sekaligus**: ADR-0029 menyimpan mata uang dan kurs bersama tiap nilai; ADR-0007 membandingkan ambang terhadap nilai IDR. `.ValueAdjustment` yang datang dari luar karena itu **selalu punya satuan yang dinyatakan**. FINDING-001 tetap terbuka sebagai **catatan atas sistem lama** dan ditutup kelak oleh REQ-011 — ia tidak menahan satu tiket pun |
| ~~B9~~ | ~~Apakah jalur `CloseClaimMD` dimigrasi atau dibuang~~ | **DIPUTUSKAN — tidak dimigrasi** | `SPEC-MODEL-DATA.md` §21.6. Status berubah `CANDIDATE-NOT-MIGRATED` → **`NOT-MIGRATED`**. Sistem baru punya **satu** jalur penutupan; setiap penutupan mencatat pelaku, waktu, dan dasar wewenangnya sebagai field terstruktur (ADR-0009). Klaim lama yang tertutup lewat jalur itu dimigrasi sebagai tertutup **dan ditandai** |

## Ditutup 19 September 2026 — bagian C

Folder `Komite Claim Non Prop` **tertutup dan tidak akan dibuka dalam pekerjaan ini**. Butir C karena itu tidak ditutup dengan jawaban, melainkan dengan **batas**.

| # | Isi | Cara tutup | Ke mana keputusannya pergi |
|---|---|---|---|
| ~~C1~~ | ~~Isi `KomiteTreaty_Flow` dan mekanisme loop persetujuan~~ | **DI LUAR KEPEMILIKAN, permanen** | Perilaku modul lain. Kita tidak merancangnya dan **tidak perlu mengetahuinya**; yang kita perlukan hanya apa yang menyeberang, dan itu `BLUEPRINT.md` §8. Ditutup sebagai *di luar kepemilikan*, **bukan** *ditunda* |
| ~~C3~~ | ~~Bagaimana Komite menutup klaim induk pada jalur `CloseClaimNP`~~ | **DI LUAR KEPEMILIKAN, permanen** | Sama |
| ~~C4~~ | ~~Bagaimana `AcceptanceStatus` dikembalikan ke `AdjustmentList(IndexObject)`~~ | **SUDAH TERJAWAB** — penggantinya, bukan tiruannya | `BLUEPRINT.md` §8.7 **syarat 1**. Penggantinya **kunci surrogate yang stabil** pada baris Adjustment, dan **kunci itu yang menyeberang** — bukan posisi. `IndexObject` tidak ada di sistem baru; **rujukan berdasarkan posisi ditolak di batas** |
| ~~C5~~ | ~~Nasib `Adjustment.AlokasiXOLPaid` — diisi siapa~~ | **DITUTUP SEBAGAI MEDAN MASUK** | `BLUEPRINT.md` §8.7 **syarat 2**. Dicari menyeluruh di empat lapisan, di enam tag yang sempat terlewat, dan di 480 baris Java: **nol penulis** (D11, D38). Itu cukup untuk merancang — ia **nilai yang datang dari luar**, boleh kosong, **tidak pernah diturunkan sendiri** oleh mesin kita |
| ~~C6~~ | ~~Siapa menulis `.ValueAdjustment`, dan dalam mata uang apa~~ | **DITUTUP SEBAGAI MEDAN MASUK bersyarat** | `BLUEPRINT.md` §8.7 **syarat 3**. ADR-0029 mewajibkan ia datang **bersama mata uang dan kurs**; nilai tanpa mata uang **ditolak di batas**, bukan diterima lalu ditebak. Itu menutup H1/H2 sebagai masalah rancangan |
| ~~C7~~ | ~~Siapa memindahkan `.AcceptanceStatus` dari `0` ke `1`/`2`~~ | **DITUTUP SEBAGAI EMPAT KEADAAN** | `SPEC-MODEL-DATA.md` §21.1. Keadaannya **empat, termasuk kosong** (S12). Model menyediakan tempat bagi transisi yang datang dari luar **tanpa mengandaikan ia terjadi**, dan **kosong bukan sinonim `0`** — perbedaan itu ditegakkan, karena justru yang kosong yang diuji `CloseClaimMD` (D2) |
| ~~C8~~ | ~~Siapa membaca `.IsEditClaim`~~ | **DITUTUP SEBAGAI PENANDA TAK BERPEMILIK** | Lihat F di bawah, dan **E17**. Dimigrasi apa adanya, ditandai tak berpemilik, tidak diberi perilaku |

## Ditutup 19 September 2026 — bagian E

| # | Isi | Cara tutup | Ke mana keputusannya pergi |
|---|---|---|---|
| ~~E1~~ | ~~Dua rule mencampur skala pembagian~~ | **DIPUTUSKAN** | **ADR-0003**. Skala mengikuti **besaran yang dihitung**, bukan rule — dan tidak ditiru. Dicatat di tabel perbedaan perilaku, `SPEC-MODEL-DATA.md` §21.5 |
| ~~E2~~ | ~~77% ekspresi aritmetika tidak menyatakan skala~~ | **DIPUTUSKAN** | Sama |
| ~~E12~~ | ~~Tiga skala dalam satu modul, terbaca pada satu rantai~~ | **DIPUTUSKAN** | Sama |
| ~~E5~~ | ~~Deskripsi langkah tidak sinkron dengan kondisi eksekusi~~ | **CACAT DOKUMENTASI SISTEM LAMA** | Tidak ditiru, tidak dimigrasi. Ditutup sebagai catatan |
| ~~E11~~ | ~~Dua alamat `prweb` berbeda tertanam berdampingan~~ | **DITUTUP BERSAMA A12** | Sistem baru tidak punya `prweb` |
| ~~E13~~ | ~~Penanda yang bukan penanda~~ | **DITUTUP SEBAGAI ATURAN** | `SPEC-MODEL-DATA.md` §21.1: **tidak ada properti berawalan `Is`/`Flag` yang diberi domain biner tanpa bukti domainnya.** `IsPLA` termasuk nilai `"1A"`, `FlagProrate` empat nilai, `ReporterStatus` tiga. **Nama tidak menentukan tipe** |
| ~~E14~~ | ~~Kolom IDR dapat terisi tanpa kurs pernah dipakai~~ | **NAIK JADI `FINDING-008`**, dan rancangannya langsung tertutup | `FINDING-008-kolom-idr-tanpa-kurs.md` (status USULAN). **ADR-0029** mencegahnya di sistem baru: nilai IDR tidak dapat lahir tanpa kurs yang menyertainya. Ramalannya diukur **REQ-034**, yang kini memverifikasi berapa baris lama terdampak — bukan menahan |
| ~~E15~~ | ~~Sepuluh kolom `SpreadingRisk` ditulis dua kali oleh rule berbeda~~ | **DITUTUP OLEH ADR-0011 dan ADR-0013** | Nilai turunan dihitung dari data, oleh **satu jalur**, dan hasilnya tidak bergantung urutan. Dua penulis atas kolom yang sama **tidak dapat terjadi pada bentuk ini** — bukan dicegah, melainkan mustahil. Kesepuluh kolomnya dicatat di tabel perbedaan perilaku, `SPEC-MODEL-DATA.md` §21.5 |
| ~~E16~~ | ~~Pengurangan selisih ke dua arah di satu berkas~~ | **DIPUTUSKAN arahnya** | **Selisih selalu nilai baru dikurangi nilai lama, satu arah, di seluruh sistem.** Kedua bentuk di `SaveToOS` 2860 dan 3999 dicatat sebagai hal yang **tidak ditiru**. Mana yang "benar" di sistem lama tidak dapat dijawab dan **tidak perlu** dijawab — ADR-0013 menghapus ketergantungan urutan yang membuat pertanyaannya ada |
| ~~E17~~ | ~~26 penanda yatim golongan 3b~~ | **SUDAH BERUPA KEBIJAKAN** | Dimigrasi apa adanya dan **ditandai tak berpemilik**. Bukan dibuang — membuang butuh kepastian yang tidak kita punya. Bukan diberi perilaku — perilakunya tidak terbaca |

## Ditutup 19 September 2026 — F, H, I sebagai **HIPOTESIS PERMANEN**

**Yang ditutup adalah status pertanyaannya, bukan cabangnya.** Cabangnya tetap tertulis utuh di **lampiran 1**.

Ketiganya tidak akan terjawab dalam pekerjaan ini, dan itu **keputusan, bukan kekurangan** — folder Komite tertutup. Penutupan ini sah justru karena rancangannya sudah berdiri di atas aturan yang benar **di kedua cabang**:

| Butir | Cabang | Yang membuat kedua cabang terpenuhi |
|---|---|---|
| **F** — `.IsEditClaim` | pembaca ada di Komite / tidak ada pembaca | Mekanisme perlindungan suntingan dibangun **tanpa syarat** (**ADR-0008**), dan kontrak batas memuat entri **bertanda belum terselesaikan dengan arah serta pemilik dikosongkan** (`BLUEPRINT.md` §8.7). `IsEditClaim` sendiri masuk golongan **penanda tak berpemilik** (E17) |
| **H** — `.ValueAdjustment` | sudah IDR di hulu / mata uang asli | Diterima di batas **dengan mata uang dan kurs wajib** (ADR-0029). Ambang dibandingkan terhadap nilai IDR (ADR-0007), dan itu sah dalam kedua keadaan |
| **I** — transisi `.AcceptanceStatus` | ditulis Komite / tidak pernah ditulis | **Empat keadaan** disediakan, termasuk kosong; transisi dari luar **mungkin tanpa diandaikan** |

**Uji dari sisi data tetap dicatat dan tetap berguna**, dan kedudukannya verifikasi: ramalan **F2** (tidak ada baris `IsEditClaim = 1` yang bertahan) dan ramalan **I2** (tidak ada klaim ber-Adjustment yang tertutup lewat `CloseClaimMD`) digabung ke **REQ-019** dan **REQ-001**.

Bila sesi Komite terjadi kelak, ketiganya **dibuka kembali sebagai verifikasi, bukan sebagai penahan**.

## Ditutup 19 September 2026 — bagian G

| # | Isi | Cara tutup | Ke mana keputusannya pergi |
|---|---|---|---|
| ~~G1~~ | ~~"syariah" — lini usaha nyata atau pemisahan teknis~~ (menyerap **A15**) | **DIPUTUSKAN sebagai atribut data**, tanpa menunggu G1a/G1b | `SPEC-MODEL-DATA.md` §21.3. **Kedua cabang menuntut hal yang sama**, dan itu yang membuatnya dapat diputus: penentu syariah hari ini adalah **node yang mengeksekusi** dan **nama akun notifikasi**, dan **tidak satu pun ikut pindah**. Maka apa pun jawabannya, sistem baru memerlukan penanda yang berdiri sendiri. **Penanda lini usaha menjadi atribut data eksplisit**, sekelas dengan pembeda yang menggantikan `pyWorkIDPrefix` (D33), disediakan **sekarang** meski belum ada yang mengisinya — sebagaimana ADR-0025 menyediakan kolom lingkup sebelum ada yang memerlukannya. **D42** tetap dicatat sebagai fakta: jalur pembayaran lama **buta terhadap syariah**. REQ-028 memverifikasi |
| ~~G2~~ | ~~Urutan kerja petugas tidak terdokumentasi di sumber mana pun~~ | **DITUTUP DENGAN MENERIMA BATASNYA** | **Ke depan risikonya sudah tertutup**: ADR-0013 membuat hasil tidak bergantung urutan tindakan. **Ke belakang, urutan penekanan kontrol tidak akan pernah dapat dipulihkan** — ia hanya pernah hidup di kepala petugas. Yang berlaku: penafsiran data lama **tidak boleh mengandaikan urutan tertentu**; bila sebuah baris lama hanya masuk akal di bawah satu urutan tertentu, ia dicatat sebagai **baris yang tidak dapat ditafsirkan**, bukan ditafsirkan dengan dugaan |

## Ditutup 19 September 2026 — bagian K

| # | Isi | Cara tutup | Ke mana keputusannya pergi |
|---|---|---|---|
| ~~K1~~ | ~~ADR-0016 melarang sesuatu yang sudah berjalan produksi~~ | **DIPUTUSKAN sebagai K1a** — ADR tetap, tanpa revisi | **ADR-0016**, bagian *Penyelesaian K1*. K1a semula mahal karena menuntut *"tetapkan cara lain memenuhi kebutuhan data HRD"*. **Biaya itu tidak ada, karena kebutuhan itu tidak ada**: modul ini tidak mengambil data HRD sama sekali, sekarang maupun nanti. Larangan ADR-0016 karena itu **tidak pernah bertabrakan dengan kebutuhan kita**. View lama dibiarkan hidup di sistem lama; `V_MST_USER_TEKNIS` **tidak dimiliki dan tidak dimigrasi** (ADR-0026). **ADR-0023 dan ADR-0028 tidak ikut ditinjau** — peninjauan itu konsekuensi K1b, dan K1b tidak dipilih. **D20** turun kedudukan dari pekerjaan menjadi catatan atas sistem lama |

## Keputusan pokok yang menutup bagian A dan K — HRD keluar dari modul ini

**Modul Claim Non Prop tidak mengambil data HRD, sekarang maupun nanti.**

Sebabnya bukan sekadar lingkup. **Pelaku suatu tindakan adalah fakta yang terjadi pada satu saat, dan fakta itu tidak berubah ketika orangnya pindah bagian atau berhenti.** Menyimpannya sebagai rujukan hidup ke HRD membuat catatan lama ikut berubah ketika sumbernya berubah — itu merusak audit, bukan memperkayanya.

| Yang berlaku | Tempatnya |
|---|---|
| Pelaku disimpan sebagai **potret** — pengenal operator dan nama sebagaimana tercatat saat tindakan terjadi. Tanpa pencarian, tanpa penyegaran, tanpa salinan tersinkron | ADR-0016; `SPEC-MODEL-DATA.md` bagian 14 dan §21.4 |
| Kalimat HRD di badan ADR-0016 **dicabut** — ia menjanjikan nama dan status kepegawaian diambil lewat API atau salinan berkala. Judul ADR ikut diganti, karena judul lama menjanjikan hal yang sama | ADR-0016 |
| Rujukan **ADR-0024** yang salah di kalimat yang sama **hilang bersamanya** — ADR-0024 adalah *kunci alami akseptasi*, bukan pemisahan pengguna/jabatan/keanggotaan komite | ADR-0016 |
| `V_MST_USER_TEKNIS` **tidak dimiliki dan tidak dimigrasi** — view HRD, bukan entitas klaim | ADR-0026 |

## Dihapus dari register — 19 September 2026

Dua butir **dihapus, bukan ditutup sebagai terjawab**. Dicatat di sini supaya nomor yang lenyap tidak menjadi tanda tanya di kemudian hari — praktik yang sama dengan A14.

| # | Isi | Sebab dihapus |
|---|---|---|
| ~~A10b~~ | ~~Apakah `VINCENTVERNANDO_1` akun uji atau akun bisnis nyata~~ | Ia bertanya tentang **orang**, dan jawabannya **tidak mengubah satu kolom pun**: baris yang dibuatnya dimigrasi apa adanya dengan pelakunya tercatat, dan tambalan atas identitas itu **tidak ikut** (ADR-0004). Sifat akunnya urusan admin Pega atas sistem lama, bukan urusan lapisan data sistem baru. REQ-016 tetap berjalan untuk memprofilkan case-nya |
| ~~A11b~~ | ~~Status kepegawaian kelima nama~~ | Sama. Dengan HRD keluar dari modul ini (ADR-0016), status kepegawaian **tidak pernah menjadi masukan bagi apa pun** yang kita rancang. REQ-026 tetap berjalan untuk mengukur aktivitas di sistem |

---

Butir yang sudah terjawab atau terserap ke artefak permanen. Tidak dihapus: pertanyaan yang pernah diajukan adalah bukti bahwa jawabannya tidak dikarang.

## Terjawab / tertutup

| # | Isi | Ke mana jawabannya pergi |
|---|---|---|
| ~~C9~~ | ~~Siapa menulis `TreatyName=="Previously Calculated UR"`~~ | **TERTUTUP 18 Sep 2026 oleh sapuan S6/S14.** Tidak ada yang menulisnya sebagai literal — **nilainya dirakit**: `GenerateCACNP_Act` baris **6643** menulis `Local.TreatyName <- "Previously Calculated " + .TreatyName`, lalu mengujinya utuh di **7223** dan **7425**. Register lama berbunyi *"tidak pernah ditulis di folder ini"*; yang benar adalah tidak pernah ditulis **utuh**. **Inilah sebabnya sapuan literal melewatkannya** — kelas cacat yang sama dengan D9. S14 memastikan ia satu-satunya dari 12 rakitan teks yang diuji utuh, jadi **satu kejadian, bukan pola**. Kedua penjaganya identik sampai kode transisinya; arti kode `2`/`3` **TIDAK DITEMUKAN**. `SAPUAN-S3-S9.md` S6, `SAPUAN-S10-S16.md` §7 |
| ~~C2~~ | ~~Siapa yang menulis `CNPStatusCase = "CLAIM ACCEPTED"` / `"CLAIM REJECTED"`~~ | **TERTUTUP 18 Sep 2026 sebagai PREMIS GUGUR, bukan sebagai terjawab.** Pertanyaannya mengandaikan kedua nilai itu ada. Sapu ulang S2 dengan pola tiga-bentuk mencari kedua literal sebagai teks telanjang di seluruh **279 berkas**: **nol kemunculan** — bukan nol sebagai nilai `CNPStatusCase`, melainkan nol **di mana pun**. Tidak ada yang menulisnya karena tidak ada yang namanya begitu. **Menyambung ke D1 dan mengeraskannya**: memori §7.3 mendaftar nilai yang **tidak pernah ada di XML sama sekali**, bukan sekadar daftar yang tidak lengkap. Dari empat nilai yang didaftar memori hanya dua yang ada, dan satu yang ada (`"INPUT ACCEPTATION CLAIM"`) tidak didaftar memori. Lihat juga **D35**. `SAPUAN-S10-S16.md` §1 |
| ~~B3~~ | ~~Properti mana yang EXPOSED sebagai kolom vs IN-BLOB di `pzPVStream`~~ | **TERTUTUP 18 Sep 2026.** DDL tabel work memperlihatkan hanya 19 kolom di luar bawaan Pega, dan **tidak satu pun nilai uang**. Seluruh `.ClaimData.*` bernilai uang **IN-BLOB — CONFIRMED**, bukan lagi dugaan. Lihat `BLUEPRINT.md` 13.5 |
| ~~B5~~ | ~~Isi sebenarnya stored procedure `PEGA_JSON_OS_AKSEP_KLAIMTNP` dkk — pemetaan `CARIn` ke kolom~~ | **TERTUTUP 18 Sep 2026.** Source-nya diterima lewat `pengetahuan/ddl/PROCEDURE_PEGA_JSON_OS_AKSEP_KLAIMTNP.sql` beserta 48 berkas lain di `pengetahuan/ddl/`. *(`REQ-003` adalah permintaannya, bukan jawabannya.)* Isinya melahirkan `ADR-0024`; dari berkas lain dalam kiriman yang sama, `pengetahuan/ddl/FUNCTION_GETCURRENCYSTANDARD.sql`, lahir `FINDING-006` |
| E3 | `SaveAdjustmentToOSAksep_Act_Tes` bernama "Tes" tetapi dipanggil dari alur produksi | **Sumber dibasiskan ulang ke XML 18 Sep 2026.** `Activity\CreateChildKomiteCNP_Act.xml` baris **8689**: `<pyStepsActivityName>Call SaveAdjustmentToOSAksep_Act_Tes</pyStepsActivityName>`. Butirnya **tetap tertutup dan isinya tidak berubah**; yang berubah dasarnya. `Struktur_Flow_TreatyIn.xlsx` (posisi 1.2.2.2.20.1.1.2, tingkat 8) kini penguat, bukan dasar — FINDING-003 menetapkan berkas itu turunan, dan aturan itu belum pernah dirambatkan ke sini |
| E4 | `Harness\Hitung_Test.xml` + `Section\Hitung_Test.xml` terdaftar sebagai local action pada alur produksi | **Sumber dibasiskan ulang ke XML 18 Sep 2026.** `Section\InputAcceptation.xml` baris **48269** dan **48524**: `<pyHarnessName>Hitung_Test</pyHarnessName>` — dua pengikatan local action. Dikuatkan S1: kedua berkas itu **mengikat medan langsung ke `InputParamOs.GrossValue` dan `InputParamOs.Value`**, yaitu parameter muatan keluar — layar uji yang menyentuh jalur produksi, bukan sekadar terdaftar di sana. Butirnya tetap tertutup; dasarnya berubah dari berkas turunan ke XML |
| E6 | Tambalan **per-orang** berumur 8 tahun dan masih menyebar: `pxCreateOperator=="VINCENTVERNANDO_1"` ditulis 2018-01-16, lalu **disalin** ke rule lain pada 2022-01-03 dan 2023-11-20 oleh dua orang berbeda | `BLUEPRINT.md` §7.1, §7.2 |
| E7 | Keputusan alur diambil dengan mencocokkan **teks bebas**: `@contains(.CommentSuggest,"Accepted by Himawan")` dan `"Rejected by Himawan"` di `SethistoryKlaimTreaty`. Satu salah ketik pengguna mengubah jalur | `BLUEPRINT.md` §7.1 |
| E8 | `CLMNP-232` ditambal **dua kali berjarak 3,5 tahun** oleh dua orang di dua rule, dan keduanya masih aktif bersamaan | `BLUEPRINT.md` §7.2 |
| E9 | Tambalan terbaru (`CLMNP-975`, 2026-07-16) berjarak **dua bulan** dari hari ini — praktik ini belum berhenti | `BLUEPRINT.md` §7.2 |
| ~~E10~~ | ~~Sembilan rule menyentuh `SpreadingRisk` secara intensif (16–66 rujukan) tanpa menyebut `"UR"`~~ | **DIKOREKSI 18 Sep 2026.** Angka 16–66 menghitung nama class di metadata langkah, bukan rujukan PageList. Rujukan sebenarnya: `GenerateCFS_act` 5, `CountTotalInsterest_Act` 2, `CopyOldataCurr_act` 2, **enam lainnya nol**. Bukan anomali — `BLUEPRINT.md` §2.4 |

## Terserap ke temuan gabungan

`E3`, `E4`, `E6`, `E7`, `E8`, `E9` digabung menjadi satu temuan tentang tata kelola perubahan — `BLUEPRINT.md` bagian 17. `E10` dicoret sebagai salah hitung, bukan anomali — `BLUEPRINT.md` §2.4.

## Pindah jalur ke `ORACLE-REQUESTS.md`

Butir berikut **tidak terjawab** — ia berpindah pemilik. Selama satu butir hidup di dua register, penutupan di satu tempat tidak terlihat di tempat lain.

| # | Pertanyaan | Pindah ke | Kolom yang menjawabnya |
|---|---|---|---|
| ~~A4~~ | ~~Apakah klaim pernah di-reopen manual~~ | **REQ-024** | `PYREOPENCOUNT > 0`, `PYREOPENTIMESTAMP`. **Sekaligus memutuskan `ADR-0002`** (tanpa reopen), yang selama ini bersandar pada ketiadaan rule — dasar paling lemah di dokumen ini |
| ~~A5~~ | ~~Berapa klaim ditutup lewat `CloseClaimMD` dan kapan terakhir~~ | **REQ-006** | `PYSTATUSWORK` + `PXUPDATEDATETIME`, dikelompokkan |
| ~~A6a~~ | ~~Nilai `PaymentType` apa saja yang benar-benar dipakai~~ | **REQ-027** | `OS_AKSEPTASI_KLAIM.PAYMENTTYPE` — sebaran nilai. **Artinya** tetap pertanyaan manusia: A6b |
| ~~A9~~ | ~~Apakah aplikasi/ruleset lain menyentuh class ini~~ | **REQ-025** | `PXAPPLICATION`, `PXAPPLICATIONVERSION` × `PXOBJCLASS` |
| ~~A10a~~ | ~~Seberapa sering `VINCENTVERNANDO_1` dipakai, sejak kapan, sampai kapan~~ | **REQ-016** | `PXCREATEOPERATOR` — cacah, rentang tanggal, dipecah per `PYSTATUSWORK`. **Sifat** akunnya tetap pertanyaan manusia: A10b |
| ~~A11a~~ | ~~Apakah kelima nama masih aktif di sistem~~ | **REQ-026** | `PXCREATEOPERATOR` / `PXUPDATEOPERATOR` + `PXUPDATEDATETIME`. **Status kepegawaian** tetap pertanyaan manusia: A11b |

Enam butir ini tidak perlu ditanyakan kepada siapa pun.

## Nomor yang berpindah

| # | Ke mana | Sebab |
|---|---|---|
| ~~A14~~ | **A10b** | Isinya — apakah `VINCENTVERNANDO_1` akun uji atau akun bisnis nyata — digabung saat A10 dipecah menjadi A10a (terukur dari kolom) dan A10b (butuh orang). Dicatat di sini supaya nomor yang lenyap tidak jadi tanda tanya di kemudian hari |

## Konflik XML vs memori yang sudah diputus

XML menang untuk perilaku. Empat puluh enam butir, dan **tidak seragam**: `D7` menyatakan memori **benar**; `D5` dan `D9` menyatakan yang keliru adalah pembacaan atau pola sapuan **saya sendiri**, bukan memori maupun XML. `D11` sampai `D14` tidak berlawanan dengan memori sama sekali — ia hal yang belum pernah tercatat di mana pun.

| # | Memori | XML | Putusan |
|---|---|---|---|
| D1 | §7.3 mendaftar `CNPStatusCase` = `"COMITEE ACCEPTANCE (DEPT. HEAD)"`, `"CLAIM ACCEPTED"`, `"CLAIM REJECTED"` | Di folder ini ada nilai keempat yang tidak tercatat di memori: **`"INPUT ACCEPTATION CLAIM"`** (`InputOutStandingCTNP_PostAct` step 2) | XML menang; enum di `BLUEPRINT.md` §3.1 yang berlaku |
| D2 | §5.5 menggambarkan `CloseClaimMD` sebagai jalur "langsung" tanpa menyebut penjaga apa pun | Penjaganya ada, tetapi **menguji hal yang berbeda dari yang disangka**: ia menolak `AcceptanceStatus=="0"`, sedangkan Adjustment yang **belum pernah dikirim** ke Komite bernilai **kosong**, bukan `"0"` — jadi lolos. Pesan galatnya sendiri membuka maksudnya: *"there is adjustment in comitee"* — mencegah penutupan saat komite **sedang bersidang**, bukan mensyaratkan persetujuan komite | **BYPASS BERSYARAT.** Putusan versi pertama ("bukan bypass tanpa penjaga") **dicabut** — lihat `BLUEPRINT.md` §4.0 |
| D3 | §10.2 mendaftar 3 tambalan per-case (`IDMaster 1000393`, `CLMNP-367`, `CLMNP-382`) | Angka terakhir: **29 langkah, 18 ekspresi unik, 8 rule**, mencakup tambalan per-case **dan** per-identitas — dan itu belum termasuk nilai mati di dalam DDL (`PROC_GENERATE_SEQUENCE_NUMBER`) | XML menang; daftar lengkap di `BLUEPRINT.md` §7, §7.1, §7.5. *(Angka 9 identitas / 4 rule / 14 langkah pada versi pertama sudah usang.)* |
| D4 | §7.4 memuat `.PICSuggest == "CHRISTINEANGELINA"` dan `"NANDINA"` | XML memuat **dua ejaan berdampingan**: `"CHRISTINEANGELINA"` (1 langkah) dan `"Christine Angelina Hutagalung"` (2 langkah); `"NANDINA"` (1) dan `"Nandina C"` (2) | XML menang; memori mencatat satu dari dua bentuk. Akibatnya sebagian cabang tidak pernah terpicu — FINDING-002 bagian 6.2 |
| D5 | Memori §5.5 **dan** analisis awal saya sama-sama melewatkan hal yang sama: keduanya membaca keberadaan penjaga sebagai bukti adanya syarat persetujuan | Penjaga itu menguji *apakah ada sidang berjalan*, bukan *apakah persetujuan sudah diberikan*. Dua pertanyaan berbeda | Bukan memori yang salah dan bukan XML yang salah — **pembacaan saya yang salah**, dan memori kebetulan sejalan dengan pembacaan itu |
| D6 | `SetPPNPPH` **tidak muncul sama sekali** di memori — tidak di katalog rumus §6, tidak di 15 rule terpenting §12, tidak di inventaris §2.1 | Rule itu **terjangkau** dari alur klaim: `Section\OutstandingClaim(1)` → `Harness\ViewPolisNonProp` (class `ASM-FW-GCNMFW-Work-ClaimTreatyNonProp`) → `<pyInclude>Section\ViewDetailDeptHeadTreatyIn_UW` → `<pyActivity>CountNetPremi_act` → `Call SetPPNPPH` | XML menang; memori melewatkannya. Tarif pajak **tetap dalam lingkup** — `ASK-AKUNTANSI.md` pertanyaan 5 |
| D7 | §10.4 butir 2 menyatakan `TotalUR` selalu nol pada cabang mata uang sama | **Terverifikasi persis dari XML**: `(.Deductible * 0 * Local.ProrateClaim/100)` di penetapan `Local.TotalUR`, sementara `Local.UR` pada cabang yang sama tidak memakai `* 0` | Memori **benar**. Naik dari catatan memori menjadi EVIDENCED; dibawa ke akuntansi sebagai pertanyaan 2b |
| D8 | §6.2 memerikan penamaan layer sebagai rumus `"XL "+substring(Layer,0,1)+…`, sehingga terbaca seolah `TreatyName` **menyimpan** teks itu | Rumus itu **probe pencarian**, bukan nilai simpan. `SpreadingRisk(<LAST>).TreatyName` diisi `ListSpreading.pxResults(1).Note` dan `.TreatyType` diisi `.ID` (`CountLossAllocation_act`); `POOLDATA.REINSURANCETYPE` memang berkolom `ID` dan `NOTE`. Dikuatkan dari dua arah lain: `@contains(.TreatyName,"R/I")` muncul 4× (FINDING-002 §6.1) — mustahil cocok bila isinya `"XL 1ST LAYER"`; dan §6.2 Langkah 10 sendiri memakai bentuk `"QS (OR)"`/`"QS (R/I)"` | XML menang; memori memerikan **masukan** lookup dan pembacanya mengira itu **keluaran**. Akibatnya: "layer" dipakai untuk dua satuan berbeda — ditetapkan sekali di `SPEC-MODEL-DATA.md` bagian 0 |
| D9 | Tidak ada — ini koreksi atas **sapuan saya sendiri**, bukan atas memori | Ketujuh nilai `PaymentType` terbaca: `1`, `2`, `5` dan `4`, `6` di `HitServiceToKasir_Act` baris 7754 dan 7896 sebagai perbandingan **tanpa kutip**; `3` dan `7` dan `''` sebagaimana sudah tercatat | **Penolakan saya atas usulan §4.1 dicabut.** Sebabnya cacat pola: saya menyapu literal berkutip saja. Nilai `1`–`7` kini EVIDENCED. Nilai kosong `''` tetap terbukti diuji, sehingga domain tertutup tujuh nilai tetap akan menolak nilai yang sah — bentuk domainnya ditetapkan saat tiketnya dikerjakan. Rinci di `SAPUAN-S3-S9.md` S3 |
| D10 | §2.5 memerikan `.IsEditClaim` sebagai pelindung suntingan petugas dari hitung ulang | `.IsEditClaim` ditulis di tiga tempat dan **dibaca nol kali** pada empat lapisan. Bayangan cerminnya: `.IsAnyAcceptation` dibaca sepuluh kali di `Harness\OutstandingClaim` dan **ditulis nol kali** | XML menang; keduanya penanda tanpa pasangan. Kelas cacat yang sama dengan `CheckDateDOL_Act`: bentuknya ada, akibatnya tidak. Rinci di `SAPUAN-S3-S9.md` S9 dan S8 |
| D11 | — | `AlokasiXOLPaid` **tidak punya penulis** pada empat lapisan: hanya sasaran perulangan, pembacaan selisih, penjaga panjang, dan PageList layar | Dicatat sebagai EVIDENCED-NIHIL, bukan kesimpulan. Pengisinya di luar ekspor yang dibaca. Rinci di `SAPUAN-S3-S9.md` S5 |
| D12 | — | `"Previously Calculated UR"` **dirakit** — `"Previously Calculated "+.TreatyName` di `GenerateCACNP_Act` baris 6643 — lalu diuji utuh di baris 7223 dan 7425 | Satu nama lapisan kasar punya bentuk turunan kedua yang beredar sebagai nilai `TreatyName` juga. Nama dan keadaan bercampur di satu kolom teks; kunci (klaim, mata uang) **belum cukup**. Menaikkan pentingnya REQ-033 |
| D13 | REQ-018 selama ini dirumuskan sebagai *mengukur apakah nomor klaim ganda terjadi* | `POOLDATA.JSON_KLAIM` mewajibkan `MNK_NO_KLAIM` **ada** (`NOT NULL`) tetapi kunci primernya **`IDPEGA`**; tidak ada struktur mana pun yang mewajibkan nomor klaim unik | Nomor klaim ganda berubah dari **dugaan** menjadi **kemungkinan struktural**. REQ-018 kini mengukur *berapa banyak yang sudah terjadi*, bukan *apakah bisa terjadi* |
| D14 | — | `POOLDATA.DIRECTTOKASIR_LOG` sudah mengarsipkan muatan ke kasir — `TGL_INPUT`, `DATA_JSON` CLOB, `IDPEGA`, `NOAKSEPTASI`, `KET` — **tanpa penanda tolak** | Arsip muatan keluar bukan hal baru; yang baru adalah **arsip berkolom**. Ketiadaan penanda tolak di jalur kasir, sementara `STS_REJECT` ada di jalur akseptasi, dicatat apa adanya |
| D15 | Rangkaian `TempKasir.CARI9` terbaca sebagai penjumlahan berantai tanpa syarat | Tiga langkah terakhirnya **berpenjaga `PaymentType`** dan saling terpisah: `1\|2\|5` menambah `.AdjustmentValue`, `4\|6` menambah `.AdjusterFeeValue`, `3` menambah `.SalvageValue` | XML menang; nilai instruksi bayar adalah `(.TotalClaim − .PremiumSpreaded)` ditambah **tepat satu** besaran, bukan seluruhnya |
| D16 | Seluruh laporan sapuan menyebut 273 berkas terbaca | Ekspor memuat **279 XML** dan seluruhnya terbaca; saringan `-not -path "*Komite*"` ikut menyaring enam berkas ber-nama Komite di folder yang terbuka | Yang salah pencacahannya, bukan sapuannya. EVIDENCED-NIHIL dan **ADR-0006 tidak bersyarat** |
| D17 | §8.4 mencatat bentuk `pzInsKey` tanpa menarik akibatnya | `pzInsKey` dan `pxCoverInsKey` pegangan instance Pega; `IndexObject` posisi numerik | Ketiganya gugur sebagai kunci kanonik (ADR-0021). `pzInsKey` **pengenal sisi lama terbaik** untuk tabel korelasi, karena `PYID` tanpa unique dan `MNK_NO_KLAIM` tidak dipaksa unik. `BLUEPRINT.md` §8.4a |
| D18 | — | `DDL_Script_ClaimNonProp2.xls` **himpunan bagian murni** dari berkas pertama: 39 dari 48, nol objek baru | Klaim saya bahwa ia “membawa objek yang belum pernah masuk daftar mana pun” **dicabut**. D13 dan D14 tetap berlaku — yang baru pembacaannya, bukan objeknya |
| D19 | — | REQ-017 tertutup **nol dari 14**. `VIEW CURRENCYSTANDARD` membaca `FROM m_CurrencyStandard` — view dan tabel basis, objek berbeda | Memiliki view tidak menutup permintaan atas tabelnya. Contoh keenam nama nyaris kembar |
| D20 | ADR-0016 melarang database link | `V_MST_USER_TEKNIS` memakai `hrdasm.v_hrd_mst@asmd.sinarmas.co.id` — **link ke instance lain, sudah berjalan produksi** | Larangannya tidak berubah; kedudukannya berubah. Bila sistem baru butuh data HRD yang sama, jalurnya belum ditetapkan ADR mana pun |
| D21 | — | Satu besaran, tiga nama: `ClaimSpreaded` → `GrossValue` → `InputParamOs.GrossValue`. Jalur Komite memakai `.ClaimSpreaded`; jalur Arasapas memakai `.ClaimEstimation + .AdjClaimValue` | Dua hilir menerima angka berbeda untuk nama yang sama. Mana yang benar **tidak terbaca** |
| D22 | — | Pengurangan selisih ke **dua arah** di satu berkas: `SaveToOS` 2860 `GrossValue - TotalClaim` berpenjaga, 3999 `TotalClaim - GrossValue` tanpa penjaga | Mana yang berlaku bergantung langkah mana tercapai lebih dulu |
| D23 | ADR-0029 mengusulkan kurs ikut tersimpan | `ClaimAmountIDR` diisi dua rule: `CountLossAllocation_act` 9585 menguji mata uang dan mengalikan kurs; `SetActualPremium_ACT` 681 **menyalin `TotalClaim` apa adanya** | Kolom IDR dapat terisi **tanpa kurs pernah dipakai**. Menguatkan ADR-0029 dari arah berbeda dari S4 |
| D24 | `arithmetic-inventory.tsv` mencatat skala mengikuti modul | Terbukti di satu berkas: rantai `AdjusterFeeValue` menyatakan **20** di setiap pembagian; rantai `GrossValue` **nol pernyataan skala**; `Local.URLimit` memakai **10** | Dua besaran uang berdampingan diperlakukan berbeda. Dasar empiris ADR-0003 |
| D25 | D10 menyebut dua penanda yatim | **14** ditulis tanpa dibaca; **26** dibaca tanpa ditulis pada objek kerja sendiri; 13 lagi milik sistem lain dan bukan cacat | Golongan 3b tidak dapat diputuskan dari XML — penulisnya di modul Komite atau tidak ada |
| D26 | — | **22 properti** dibandingkan sebagai angka **dan** teks; **10** punya nilai kosong tanpa padanan numerik; `IsPLA` bernilai `1A` | Sepuluh properti tidak dapat diberi domain numerik tertutup. `IsPLA` bukan penanda biner meski namanya mengaku begitu |
| D27 | BLUEPRINT mencatat daftar alokasi menumpuk | Sebabnya terbaca: langkah 10 `Property-Remove` **tidak pernah diisi parameternya**, dan tidak ada `Page-Remove` atas `SpreadingRisk` di mana pun | Menumpuk karena satu parameter kosong, bukan karena rancangan memutuskan begitu |
| D28 | §2.5 memerikan `IsEditClaim` sebagai pelindung suntingan | Ia **penugasan terakhir** tiap blok pembuatan baris, ditulis `0` (7322 dan 9795) | Ia medan inisialisasi. Perlindungannya **tidak pernah lengkap**, bukan rusak belakangan |
| D29 | D12 menyebut “nama dan keadaan bercampur” | Dari **12** rakitan teks, hanya `"Previously Calculated "+.TreatyName` yang diuji utuh. Kedua penjaganya **identik** | Klaim “pola” saya **diturunkan** jadi satu kejadian. Arti kode transisi `2`/`3` **TIDAK DITEMUKAN** |
| D30 | S9 menyimpulkan `EditXOLAlokasi` bukan sasarannya | **Sepuluh kolom** `SpreadingRisk` ditulis mesin alokasi **dan** salah satu dari `AdjClaimCNP_Act`, `AdjClaimAmount_Act`, `SetActualPremium_ACT` | Pertanyaan terbelah: *hitung ulang ganda* **terjawab**; *suntingan manual* **TIDAK DITEMUKAN**. Bila bertahan, yang dilindungi ADR-0008 belum pernah ada jalurnya |
| D31 | Konvensi Pega: `px` hanya-baca, `pz` bukan untuk dipakai langsung | `px` **ditulis aplikasi** pada setiap baris PageList tertanam. `pz` bertahan — dua penulisannya penyalinan pegangan, bukan pembuatan | Konvensinya menyesuaikan diri pada bukti: “hanya-baca” berlaku pada work object, tidak pada baris tertanam |
| D32 | — | **Pelaku dan waktu setiap baris tertanam** tersimpan di `pxCreateOperator`, `pxCreateOpName`, `pxCreateDateTime` — alokasi, estimasi, akseptasi, adjustment | Membuang `px` mentah berarti **kehilangan seluruh kolom pelaku**. Risiko terbesar S16 |
| D33 | — | `pyWorkIDPrefix` memutuskan **lini bisnis 17 kali** — `CLM-`, `CLMP-`, `CLMNP-` | Sistem baru harus menyediakan pembeda yang berdiri sendiri, bukan disimpulkan dari bentuk nomor |
| D34 | G1 syariah masih terbuka | `pyNotifyAccountName` dibedakan `"NUSARE"` dan `"NUSARESYARIAH"` di `SendEmailKlaim` 3021 dan 3281 | Pembedaan syariah beredar di properti notifikasi Pega, bukan medan bisnis. Dicatat; tidak ditarik kesimpulan |
| D35 | §7.3 mendaftar empat nilai `CNPStatusCase` | Dua di antaranya — `"CLAIM ACCEPTED"` dan `"CLAIM REJECTED"` — muncul **nol kali di mana pun** pada 279 berkas. `CNPStatusCase` juga **tidak pernah dibaca** | D1 mengeras: dari empat yang didaftar memori hanya dua ada, dan satu yang ada tidak didaftar. Ejaan tersimpan `COMITEE` satu T |
| D36 | Laporan S1 memetakan sisi Arasapas | Sisi kasir memakai **25 medan bernama urut** `TempKasir.CARI1`…`CARI24` dan `CARIDATETIME` | Belum pernah diinventarisasi. **Tiket `05` kembali ke papan aktif** — pernyataannya mencakup kasir |
| D37 | D25 menggolongkan `AktifButton` sebagai penanda yang ditulis tanpa pernah dibaca | Ia **dibaca** — `ProteksiSendKomiteCNP_Act` 4491 dan 4969, lewat `<Field>.AktifButton</Field>`, tag yang tidak pernah dibaca indeks lama | **Klaim yatim atasnya dicabut.** Ia golongan 1. Golongan 2 turun 14 → 13 |
| D38 | — | `IsEditClaim`, `CNPStatusCase`, `IsAnyAcceptation` **tidak bergerak** setelah jaring diperlebar 30 → 45 tag **dan** setelah 480 baris Java dibaca sebagai kode | **F2, H2, I2 menguat — tidak tertutup.** Java adalah tempat terakhir di folder terbuka yang dapat memuat pembaca tersembunyi, dan ia bersih |
| D39 | AK-1 mensyaratkan presisi dinyatakan sepanjang rantai | **Nol `double`, nol `float`, nol `BigDecimal`** di 480 baris Java. Uang melintas sebagai `String` (92 kemunculan) lewat `.toString()` | Ketiadaan tipe biner pecahan **bukan** berarti presisi terjaga. Java tidak menghitung uang; ia memindahkannya sebagai teks tanpa kendali skala |
| D40 | Kunci alami tiap kumpulan disimpulkan dari pemakaian | Tujuh blok Java adalah **satu rutin dedup yang disalin-tempel** dan menegakkan kunci alami di kode: `Kurs` `(CurrencyID, Currency)`; `tempCurrency` `(Currency)`; `Estimate` `(CurrencyID, Currency, TypeLossID, TypeLoss)`; `Cuan` dan `.ClaimData.TotalInterestInsured` `(CurrencyID, Currency)`; `TempTotal` `(XOL, Currency)` | Kunci kini **terbaca sebagai kode**, bukan disimpulkan. Dan satu di antaranya **membuang baris dari `.ClaimData.TotalInterestInsured`** — data klaim, tanpa pencatatan, tanpa penanda, tanpa jejak |
| D41 | — | Muatan kasir dirakit dengan penyambungan teks. `Nett`, `Deductible`, `KaliDeduct` dikirim **tanpa kutip** sebagai angka JSON, dari `String` hasil `.toString()` | Nilai kosong menghasilkan `"Nett":,` — **JSON rusak**, gagal di sisi penerima. Tidak ada pemeriksaan angka, kendali skala, maupun pelolosan karakter untuk 18 medan teks |
| D42 | G1 mencatat pembedaan syariah beredar di `pyNotifyAccountName` | Di muatan kasir, `StsSyariah` **selalu `0`** dan `CompanyName` **selalu `"NUSARE"`** — keduanya tertanam di kode, bersama `Deductible`, `KaliDeduct`, `LjtdId`, `LdcId`, `StsAp` | **Jalur pembayaran tidak pernah menyatakan syariah**, apa pun keadaan klaimnya. G1 mendapat arah baru: pembedaan ada di surel, tidak ada di uang |
| D43 | B4 menyatakan struktur JSON belum terdokumentasi | `adoptJSONObject` di `GeneratePlaCNP_Act`, `GetDetailPolis_act`, `SetValueClaimTNP_Act` mengadopsi teks `HASIL1` **tanpa pemeriksaan bentuk**; yang terakhir menulis ke `pyWorkPage.TreatyInMaster` | Bukan hanya tidak terdokumentasi — **tidak dibatasi**. Apa pun yang ada di `HASIL1` menjadi properti objek kerja |
| D44 | D32 mendaftar tiga kolom pelaku di properti `px` | `pxAttachedBy` menyusul, lewat tag `<Property>` yang baru dilipat | Membuang `px` mentah menghilangkan **pencatat lampiran** juga, bukan hanya pencatat baris tertanam |
| D45 | Laporan pemeriksaan cakupan menyatakan `GetMimeType` memuat **nol sel keputusan** | **Salah.** Ada **42 slot baris** bernomor 1–42; yang kosong adalah blok nilainya — 42 tag menutup sendiri berturut-turut | **Klaim saya dicabut.** Yang benar: 42 baris ada, nilainya tidak ikut terekspor. **Celah bahan, bukan celah metode** — diajukan sebagai REQ-035 |
| D46 | Laporan S1 memetakan sisi Arasapas; sisi kasir belum pernah diinventarisasi (D36) | Sisi kasir terpetakan lengkap dari Java: **21 medan bernama dari 25 slot** `CARI`. `CARI19`–`CARI21` hanya perakit tanggal, `CARIDATETIME` penampung sementara | **Tiket `05` ditutup.** Kedua sisi muatan keluar kini terbaca |

Dipertahankan sebagai catatan bahwa setiap putusan punya dasar.


---

# Lampiran — badan butir yang ditutup, disimpan utuh

Butir di bawah **sudah ditutup**. Badannya dipertahankan apa adanya karena ia memuat bukti, cabang, dan ramalan teruji yang tidak tercatat di tempat lain. **Tidak satu pun dari isi ini adalah pertanyaan terbuka.**

## Lampiran 1 — F, H, I: cabang hipotesis permanen

> Ditutup 19 September 2026 sebagai **HIPOTESIS PERMANEN — RANCANGAN MEMENUHI KEDUA CABANG**. Cabangnya tetap tertulis supaya sesi Komite kelak punya bahan verifikasi yang lengkap.

#### F. Hipotesis sejajar — `.IsEditClaim`
> **Isi di bawah ini ARSIP, bukan register.** Butirnya sudah ditutup; badannya dipertahankan karena memuat bukti dan cabang yang tidak tercatat di tempat lain. Bila Anda masuk lewat pencarian teks dan bukan lewat daftar isi: **tidak ada pertanyaan terbuka di bawah sini.**


Menandai ini "ditunda" saja tidak cukup: kedua kemungkinan punya konsekuensi yang harus dinyatakan sekarang, karena keduanya mengubah pekerjaan yang berbeda.

| | **F1 — pembacanya ada di modul Komite** | **F2 — tidak ada pembaca sama sekali** |
|---|---|---|
| Artinya | Bendera ditulis di Claim, dibaca di Komite | Bendera ini mati |
| Konsekuensi | **Kopling lintas modul yang belum ada di Boundary Contract.** Sudah dicatat sebagai tambahan di `BLUEPRINT.md` §8.5 | **Perlindungan hasil edit manual tidak pernah ada.** Hasil suntingan petugas lewat `EditXOLAlokasi` tertimpa setiap kali `CountLossAllocation_act` dijalankan ulang — dan `CountLossAllocation_act` dipanggil dari **28 titik pengikatan di 7 berkas** — *diperiksa ulang dari XML 18 Sep 2026, bukan lagi dari berkas turunan*: 4 `Call` di Activity (`AddAkseptasiCNP_Act` 5820, `CountClaimTNP_Act` 7213, `InputAkseptasi_PreAct` 1689, `SetTPLNote_Act` 2571) dan 24 `<pyActivity>` di Section (`AdjustmentDetailNP` 12×, `InputAcceptation` 6×, `OutstandingClaim(1)` 6×). Klaim lama "sedikitnya 8" **benar dan terlalu rendah**. (lihat `Struktur_Flow_TreatyIn.xlsx`) |
| Yang harus dikerjakan | Tambahkan ke kontrak batas, tentukan arah dan pemiliknya | Naikkan menjadi **FINDING-004**, dan pastikan sistem baru punya mekanisme kunci hasil edit yang sebenarnya |
| Cara memastikan | Sapu modul Komite atas `IsEditClaim` — **satu sesi, satu grep** | Hasil sapuan yang sama, bila nihil |
| **Apa yang mustahil ada di data kalau ini benar** | — (F1 tidak melarang apa pun) | Kalau F2 benar, **mustahil ada baris `SpreadingRisk` dengan `IsEditClaim = 1` yang bertahan** setelah alokasi dijalankan ulang. Satu baris begitu saja menggugurkan F2 |
| **Bukti yang masuk 19 Sep 2026** | — | **Menguat, tidak tertutup.** 480 baris Java tertanam di 18 berkas dibaca sebagai kode: `IsEditClaim` muncul **nol kali**. Java adalah tempat terakhir di folder terbuka yang dapat memuat pembaca tersembunyi. Jaring juga diperlebar 30 → 45 tag; `IsEditClaim` tetap ditulis 3, dibaca 0 (**D38**, **D39**) |
| **Bukti yang masuk 18 Sep 2026** | — | **Dua, keduanya menguatkan F2 tanpa menutupnya.** (1) **S13** — kedua penulisan `IsEditClaim = 0` (`CountLossAllocation_act` 7322 dan 9795) adalah **penugasan terakhir tiap blok pembuatan baris**, yaitu medan **inisialisasi**: setiap baris baru lahir dengan `0`. Perlindungannya karena itu **tidak pernah lengkap sejak awal**, bukan rusak belakangan. (2) **S11** — sapuan menyeluruh atas 104 kandidat penanda menguatkan **nol pembacaan**, dan menemukan `IsEditClaim` bukan sendirian: 14 penanda ditulis tanpa pernah dibaca, 26 dibaca tanpa pernah ditulis |


> **Disiplin H berlaku di sini — disalin ke badannya 19 September 2026, bukan dirujuk dari jauh.**
>
> **Tidak boleh ada rancangan yang mengandaikan salah satu cabang benar.** Aturan itu semula hanya tertulis di bagian H; sejak folder `Komite Claim Non Prop` dinyatakan **tertutup untuk batch ini** — bukan ditunda, tertutup — ia berlaku penuh untuk F dan I juga.
>
> **Ketiganya tinggal permanen sebagai hipotesis sejajar dalam batch ini, dan itu hasil yang sah.** Yang tidak sah adalah membiarkan tiket mati menunggu sesuatu yang tidak akan datang. Karena itu pekerjaannya dikerjakan **di atas aturan yang sama-sama benar di kedua cabang**, dan hanya di atas itu.

> **Yang dikerjakan meski cabangnya tidak diketahui — karena kedua cabang menuntutnya:**
>
> | Cabang | Yang dituntutnya | Kerjakan? |
> |---|---|---|
> | F1 | kopling lintas modul masuk Boundary Contract | **ya** — entri ditulis **bertanda belum terselesaikan**, dengan arah dan pemilik **dikosongkan**, bukan ditebak |
> | F2 | sistem baru punya kunci hasil suntingan yang sebenarnya | **ya** — ADR-0008 sudah menuntutnya **terlepas dari** apakah sistem lama pernah punya jalurnya |
>
> Mekanisme perlindungan suntingan karena itu dibangun **tanpa syarat**. Tiket `17` **tidak tertahan** bagian ini.

Keduanya ditutup oleh **satu pemeriksaan yang sama** di sesi Komite — sesi yang **tidak terjadi dalam batch ini**.

**Uji dari sisi Claim, tanpa menunggu sesi Komite**: `.IsEditClaim` **IN-BLOB** (tidak termasuk 19 kolom yang di-expose), jadi ujinya baru mungkin setelah blob terbaca — REQ-001. Dicatat supaya tidak terlupa saat REQ-001 masuk.

---

#### H. Hipotesis sejajar — `.ValueAdjustment` (C6)
> **Isi di bawah ini ARSIP, bukan register.** Butirnya sudah ditutup; badannya dipertahankan karena memuat bukti dan cabang yang tidak tercatat di tempat lain. Bila Anda masuk lewat pencarian teks dan bukan lewat daftar isi: **tidak ada pertanyaan terbuka di bawah sini.**


Bentuk yang sama dengan bagian F, dan atas alasan yang sama: dua kemungkinan dengan konsekuensi berbeda, ditutup oleh satu pemeriksaan yang sama.

`.ValueAdjustment` muncul di **tepat satu berkas** (`CreateChildKomiteCNP_Act`), hanya dibaca, bukan field UI mana pun di folder ini. Asal nilainya **TIDAK DITEMUKAN DI XML**.

| | **H1 — sudah dikonversi ke IDR di hulu** | **H2 — mata uang aslinya, apa pun itu** |
|---|---|---|
| Artinya | Modul lain mengonversinya sebelum diserahkan | Nilai disimpan apa adanya per mata uang |
| Konsekuensi | **`FINDING-001` gugur seluruhnya.** Perbandingan terhadap ambang `30.000.000` sah, dan tidak ada cacat apa pun | **`FINDING-001` berdiri.** Ambang dibandingkan terhadap angka yang satuannya tidak diketahui |
| Yang harus dikerjakan | Tutup FINDING-001, catat sebagai dugaan yang gugur | Kuantifikasi paparan lewat REQ-011, dan `ADR-0007` menjadi koreksi atas cacat nyata |
| Cara memastikan | Sapu modul Komite atas `.ValueAdjustment` — **satu sesi, satu grep** | Hasil sapuan yang sama |
| **Apa yang mustahil ada di data kalau ini benar** | Kalau H1 benar, **mustahil ada Adjustment bermata uang selain IDR yang `.ValueAdjustment`-nya sama dengan nilai asli mata uang itu**. Bila ditemukan nilai yang jelas bukan hasil konversi — misalnya nilai USD kecil berdampingan dengan `Currency='USD'` — H1 gugur | Kalau H2 benar, mustahil seluruh `.ValueAdjustment` pada klaim non-IDR berskala rupiah |

~~Keduanya ditutup oleh **satu pemeriksaan yang sama** di sesi Komite. Sampai itu terjadi, **tidak boleh** ada keputusan desain yang mengandaikan salah satunya benar.~~

> **Larangan di atas DICABUT 19 September 2026.** Ditulis tetap terlihat, bukan dihapus: bagian H adalah **sumber** disiplin itu — F dan I mewarisinya dan keduanya sudah diberi catatan penutupan, sementara induknya tidak pernah. Yang gugur dicabut di tempatnya, dan itu termasuk terlihat bahwa ia pernah ada.
>
> **Yang mencabutnya**: **ADR-0029**, yang naik `accepted` pada tanggal yang sama dan mewajibkan **setiap** nilai uang datang bersama nilai asli, kode mata uang, nilai rupiah, kurs yang dipakai, dan asal-usul kurs itu — tanpa pengecualian.
>
> **Ditutup 19 September 2026 sebagai HIPOTESIS PERMANEN — RANCANGAN MEMENUHI KEDUA CABANG.**
>
> Rancangannya memenuhi **H1 dan H2 sekaligus**, dan itu yang membuat penutupan ini sah tanpa cabangnya dipilih:
>
> | Cabang | Yang dituntutnya | Terpenuhi oleh |
> |---|---|---|
> | **H1** — sudah dikonversi ke IDR di hulu | nilainya menyatakan bahwa ia sudah rupiah | nilai yang menyeberang **membawa kode mata uangnya**; bila ia `IDR`, itu terbaca, bukan ditebak |
> | **H2** — mata uang aslinya, apa pun itu | nilainya menyatakan mata uang apa | sama, dan ADR-0007 membandingkan ambang terhadap **nilai IDR**, bukan terhadap angka tanpa satuan |
>
> **`BLUEPRINT.md` §8.7 syarat 3 menegakkannya di batas**: nilai tanpa mata uang **ditolak di batas** — bukan diterima lalu ditebak. Dengan itu keadaan yang melahirkan FINDING-001 tidak dapat terjadi di sistem baru, entah H1 atau H2 yang benar di sistem lama.
>
> **FINDING-001 tetap terbuka sebagai catatan atas sistem lama**, dan ditutup kelak oleh REQ-011 — yang kini **memverifikasi, tidak menahan**.

---

#### I. Hipotesis sejajar — transisi `.AcceptanceStatus` (C7)
> **Isi di bawah ini ARSIP, bukan register.** Butirnya sudah ditutup; badannya dipertahankan karena memuat bukti dan cabang yang tidak tercatat di tempat lain. Bila Anda masuk lewat pencarian teks dan bukan lewat daftar isi: **tidak ada pertanyaan terbuka di bawah sini.**


`.AcceptanceStatus` ditulis **satu kali saja** di seluruh 279 berkas: `CreateChildKomiteCNP_Act` langkah 31 → `= 0`. Transisi ke `1` atau `2` tidak ada di folder ini.

| | **I1 — ditulis oleh modul Komite** | **I2 — tidak pernah ditulis sama sekali** |
|---|---|---|
| Artinya | Komite menandai hasil persetujuannya di properti ini | Properti ini menetap di `0` selamanya |
| Konsekuensi | Kopling lintas modul yang **belum ada di Boundary Contract** — arah Komite → Claim. Tambahkan | Penjaga `CloseClaimTNonProp` yang menolak `=="0"` akan **menolak setiap klaim yang pernah punya Adjustment**, selamanya. Jalur `CloseClaimMD` efektif mati untuk klaim tersebut |
| Yang harus dikerjakan | Lengkapi kontrak batas, tetapkan pemiliknya | Naikkan jadi temuan tersendiri; dan `BLUEPRINT.md` §4.0 perlu ditulis ulang lagi |
| Cara memastikan | Sapu modul Komite atas `.AcceptanceStatus` — **satu sesi, satu grep** | Hasil sapuan yang sama |
| **Apa yang mustahil ada di data kalau ini benar** | — | Kalau I2 benar, **mustahil ada klaim yang punya Adjustment dan ditutup lewat `CloseClaimMD`**. Satu baris begitu saja menggugurkan I2 |


> **Disiplin H berlaku di sini — disalin ke badannya 19 September 2026, bukan dirujuk dari jauh.**
>
> **Tidak boleh ada rancangan yang mengandaikan salah satu cabang benar.** Aturan itu semula hanya tertulis di bagian H; sejak folder `Komite Claim Non Prop` dinyatakan **tertutup untuk batch ini** — bukan ditunda, tertutup — ia berlaku penuh untuk F dan I juga.
>
> **Ketiganya tinggal permanen sebagai hipotesis sejajar dalam batch ini, dan itu hasil yang sah.** Yang tidak sah adalah membiarkan tiket mati menunggu sesuatu yang tidak akan datang. Karena itu pekerjaannya dikerjakan **di atas aturan yang sama-sama benar di kedua cabang**, dan hanya di atas itu.

> **Yang dikerjakan meski cabangnya tidak diketahui:** model **menyediakan tempat** bagi transisi `AcceptanceStatus` yang datang dari luar, **tanpa mengandaikan ia terjadi**. Keadaannya **empat, termasuk kosong** — sesuai temuan **C7**: properti ini dibandingkan sebagai angka (`0`,`1`,`2`) **dan** sebagai teks, dan bentuk teksnya memuat nilai kosong yang tidak punya padanan numerik. Nilai kosong itu justru yang diuji `CloseClaimMD` (D2). Tiket `18` **tidak tertahan** bagian ini.

**I1/I2 dapat diselesaikan sekarang, tanpa menunggu sesi Komite.** Ramalan I2 dapat dibantah data: cacah klaim tertutup, dipecah menurut ada-tidaknya Adjustment. Digabungkan ke **REQ-019**, yang sudah menghitung klaim tertutup dan `FLAGONGOINGCOMMITTE`.

Bila berhasil, ini **butir DEFERRED-TO-KOMITE-SESSION pertama yang diselesaikan dari sisi Claim**.

---

## Lampiran 2 — G1 dan G2: bukti yang menyempitkan

> Ditutup 19 September 2026. Badannya dipertahankan karena ia memuat sapuan `pxSystemNodeID`, pembalikan arah bukti, dan batas cakupan ekspor — tiga hal yang tidak tercatat di berkas lain.

#### G. Pertanyaan terbuka & asumsi bernama — dicatat supaya tidak hilang
> **Isi di bawah ini ARSIP, bukan register.** Butirnya sudah ditutup; badannya dipertahankan karena memuat bukti dan cabang yang tidak tercatat di tempat lain. Bila Anda masuk lewat pencarian teks dan bukan lewat daftar isi: **tidak ada pertanyaan terbuka di bawah sini.**


### G1. "syariah" — **pertanyaan terbuka, bukan asumsi**. Kesimpulan versi pertama dicabut.

> **Bukti baru 18 September 2026, dari S16 — dan ia tidak menguatkan G1a maupun G1b.**
>
> `Param.pyNotifyAccountName` diisi `"NUSARE"` (`SendEmailKlaim` 3021; `SendEmailKlaimRejectClose` 2557) dan `"NUSARESYARIAH"` (`SendEmailKlaim` 3281). Pembedaan konvensional/syariah karena itu **beredar sebagai nilai properti notifikasi bawaan Pega**, bukan sebagai medan bisnis.
>
> **Mengapa ini tidak menguatkan salah satu cabang.** G1a mengandaikan pemisahan di tingkat **instance atau node**; G1b mengandaikan pemisahan di tingkat **data**. Sebuah akun notifikasi surel bukan keduanya — ia pemisahan di tingkat **pengiriman pesan**, dan satu instance dapat memiliki dua akun notifikasi persis sebagaimana dua instance dapat. Bukti ini menyempitkan tempat mencarinya, bukan menjawabnya.
>
> **Bukti kedua, 19 September 2026, dan arahnya berlawanan.** Di muatan kasir (`HitServiceToKasir_Act`, Java blok 1), `StsSyariah` **selalu `0`** dan `CompanyName` **selalu `"NUSARE"`** — keduanya tertanam di kode, tidak pernah berasal dari data. **Jalur pembayaran tidak pernah menyatakan syariah, apa pun keadaan klaimnya.** Jadi pembedaan itu ada di jalur surel dan **tidak ada di jalur uang**. Itu fakta tentang sistem lama, bukan jawaban atas G1 — tetapi ia mempersempit: apa pun jawabannya, jalur pembayaran hari ini buta terhadapnya (**D42**).
>
> Yang **bertambah** dari bukti ini adalah butir migrasi tersendiri: apa pun jawaban G1, sistem baru memerlukan pembedaan itu **secara eksplisit**, karena `pyNotifyAccountName` tidak ikut pindah. Sekelas dengan D33 (`pyWorkIDPrefix` memutuskan lini bisnis 17 kali).

**Yang saya tulis sebelumnya**: *"`pxSystemNodeID = \"jboss1074\"` mengidentifikasi node server, bukan jenis bisnis. Jadi 'syariah' menandai instance aplikasi Pega yang terpisah, bukan lini usaha."*

**Premisnya benar, kesimpulannya tidak menyusul.** Kalau ada instance Pega terpisah untuk bisnis syariah, maka bisnis syariah itu **ada** — ia hanya dilayani deployment yang berbeda. "Bukan lini usaha" tidak mengikuti dari "dibedakan di tingkat node".

**Sapuan `pxSystemNodeID` — hasilnya menyempitkan masalah, dan membalik arah bukti.**

| Yang dicari | Hasil |
|---|---|
| Rule yang bercabang pada `pxSystemNodeID` | **satu**: `When\IsPEGASyariah.xml` |
| Rule yang memakai `IsPEGASyariah` | **satu**: `Activity\HitServiceToKasir_Act.xml`, 3 langkah |
| `jboss117` (41x) dan `jboss122117` (36x) | **bukan percabangan** — seluruhnya `<pxHostId>`, metadata audit host penyimpan rule |

**Apa yang berubah pada ketiga langkah itu**: `TempKasir.CARI15 = "100115"` — **kode yang dikirim ke sistem kasir**.

Itu justru bukti ke arah sebaliknya dari kesimpulan pertama saya: **kode setoran ke kasir yang berbeda menunjukkan entitas pembukuan yang berbeda**, bukan sekadar server yang berbeda.

**Dua kemungkinan, ditulis sejajar, tidak diputuskan sendiri:**

| | **G1a — lini usaha syariah nyata** | **G1b — hanya pemisahan teknis** |
|---|---|---|
| Arti kode `"100115"` | kode perusahaan/akun terpisah di sistem kasir | sekadar kode tujuan yang kebetulan berbeda per lingkungan |
| Konsekuensi | Pemisahan syariah adalah **atribut pada kontrak treaty**, dan model kontrak harus meninjau ulang pembukuan, pelaporan, serta treaty terpisah | Cukup dirancang sebagai routing; tidak ada dampak ke model data |
| Yang harus dikerjakan | Lihat tiga butir di bawah | Catat, lanjutkan |

**Bila G1a terbukti, tiga hal menyusul — dipetakan sekarang supaya tidak ditemukan belakangan:**

1. **Sistem baru harus tahu klaim mana yang syariah, dan penentu lamanya tidak dapat diwarisi.** Hari ini satu-satunya penentu adalah **node mana yang mengeksekusi**. Sistem baru kemungkinan besar satu deployment, sehingga penentu itu lenyap bersama arsitekturnya.
2. **Penentu syariah menjadi atribut data, bukan atribut runtime** — melekat pada kontrak treaty, bukan pada nama server. Itu keputusan model data dan menyentuh `CONTEXT.md`, bukan sekadar konfigurasi.
3. **Apakah klaim non-prop syariah pernah ada adalah pertanyaan terukur**: cacah baris berkode setoran `"100115"` di `DIRECTTOKASIR_LOG` atau di tabel kasir. Nol berarti G1 selesai sebagai jalur mati. Bukan nol berarti butuh pemodelan tersendiri. → **REQ-028**.

**Yang sudah pasti, apa pun jawabannya**: **perilaku sistem berbeda tergantung node yang mengeksekusi.** Itu dimensi yang tidak pernah masuk analisis saya selama enam ronde, dan sekarang tercatat.

**Batas cakupan ekspor**: karena percabangannya berupa When rule yang dievaluasi saat berjalan **di dalam ruleset yang sama**, ekspor 279 berkas ini berlaku untuk kedua node. Yang **tidak** dapat saya buktikan: apakah ada instance Pega lain dengan ruleset yang sama sekali berbeda. Itu masuk A15, dibawa ke admin Pega.

### G2. Urutan kerja petugas tidak terdokumentasi di sumber mana pun

Sapuan `MEMORI_PEMAHAMAN.MD` atas SOP, urutan kerja, langkah petugas, prosedur operasional, dan tata cara: **nihil**.

Digabung dengan temuan bahwa urutan perhitungan ditentukan kontrol yang ditekan (`FINDING-003` bagian 3), artinya: **urutan tindakan yang menentukan hasil perhitungan hanya hidup di kepala petugas.** Tidak ada satu pun artefak yang menangkapnya.

> **Diperiksa 18 September 2026 terhadap `Struktur_Flow_TreatyIn.xlsx` — G2 MENYEMPIT, tidak tertutup.**
>
> Berkas itu memuat **pohon pemanggilan berjenjang**: 872 baris, kedalaman 1–15, 262 rule unik, 261 nama berkas XML, akar tunggal `Flow_TreatyIn`, seluruhnya bertanda *"Sudah diupload"*. Rinciannya per jenis rule: Activity 361, RDB List 173, Report Definition 85, Section 73, Data Transform 60, Flow Action 29, Harness 27, When 22, Connect REST 20, System Settings 17, Decision Table 4, Flow 1.
>
> **Yang ia jawab**: urutan **teknis** — rule mana memanggil rule mana, dan di jenjang keberapa. Ia juga memperlihatkan bahwa alur petugas punya **tepat dua assignment** di tingkat atas: `OutstandingClaim` (`Assignment2`) dan `InputAcceptation` (`Assignment1`). Itu urutan **layar**, dan itu memang dokumentasi.
>
> **Yang ia tidak jawab, dan inilah inti G2**: urutan **tindakan petugas di dalam satu layar**. `FINDING-003` bagian 3 menyatakan hasil perhitungan ditentukan **kontrol mana yang ditekan lebih dulu**; pohon pemanggilan tidak memuat itu, karena ia menggambarkan struktur, bukan pilihan. Satu contohnya langsung terbaca: `CountLossAllocation_act` terikat di **24 tempat** di tiga Section, dan tidak ada apa pun di pohon yang menyatakan urutan penekanannya.
>
> **Sisa G2 sesudah menyempit**: bukan lagi *"tidak terdokumentasi di sumber mana pun"* — urutan layar terdokumentasi. Yang tetap tidak terdokumentasi adalah **urutan penekanan kontrol di dalam layar**, dan itu tepat bagian yang menentukan hasil perhitungan.
>
> Catatan kedudukan berkas: `FINDING-003` menetapkan berkas ini **turunan, bukan sumber**. Pembacaan di atas karena itu dipakai untuk **menyempitkan pertanyaan**, bukan sebagai bukti perilaku. Klaim perilaku mana pun tetap harus berdiri di atas XML.

Keputusan Q22 (ADR-0013) membuat sistem baru tidak bergantung pada urutan, jadi risikonya tertutup ke depan. Yang tidak tertutup: **menafsirkan data lama** kadang memerlukan pengetahuan tentang urutan yang ditempuh, dan pengetahuan itu tidak terdokumentasi.

---

## Lampiran 3 — K1: dua cabang revisi ADR

> Ditutup 19 September 2026 sebagai **K1a**. Cabang K1b tetap tertulis supaya alasan penolakannya dapat diperiksa ulang.

### K. Kandidat revisi ADR
> **Isi di bawah ini ARSIP, bukan register.** Butirnya sudah ditutup; badannya dipertahankan karena memuat bukti dan cabang yang tidak tercatat di tempat lain. Bila Anda masuk lewat pencarian teks dan bukan lewat daftar isi: **tidak ada pertanyaan terbuka di bawah sini.**


Butir di sini **bukan pertanyaan teknis**; ia keputusan yang sudah diambil dan kini perlu ditinjau karena buktinya berubah. ~~Tidak satu pun dipilih di sini.~~ **Dipilih K1a 19 September 2026 — lihat ADR-0016 bagian Penyelesaian K1.**

### K1 — ADR-0016 melarang sesuatu yang sudah berjalan produksi

**Bukti, 18 September 2026, sapuan rekonsiliasi DDL.** `pengetahuan/ddl/VIEW_V_MST_USER_TEKNIS.sql` baris 10:

```
FROM hrdasm.v_hrd_mst@asmd.sinarmas.co.id b
```

Sebuah **database link ke instance Oracle lain**, dipakai oleh view yang berjalan hari ini. Sapuan pola `@<host>` atas keempat puluh sembilan berkas DDL menemukan **satu** pemakaian, dan inilah dia.

ADR-0016 melarang database link. Larangannya tidak berubah; **kedudukannya** berubah — dari larangan atas sesuatu yang belum pernah ada, menjadi larangan atas sesuatu yang sudah dipakai.

| | **K1a — ADR tetap, yang lama dibiarkan** | **K1b — ADR direvisi** |
|---|---|---|
| Isinya | Larangan berlaku penuh untuk sistem baru. View lama beserta link-nya tidak ikut pindah, dan kebutuhan data HRD dipenuhi dengan cara lain | Larangan diberi pengecualian bernama untuk data acuan lintas sistem, dengan syarat yang ditulis |
| Yang harus dikerjakan bila ini dipilih | Tetapkan **cara lain** itu sebelum ada tabel yang bergantung pada data HRD. Belum dibahas di ADR mana pun | Tulis syaratnya: arah, pemilik, apa yang boleh lewat, dan apa yang tidak. ADR-0023 ikut ditinjau, karena alasan pelarangan link semula adalah view lintas mesin bukan view |
| Biayanya | Pekerjaan baru yang belum dilingkupi | Melemahkan alasan yang menopang ADR-0023 dan ADR-0028 |

~~**Tidak dipilih di sini.** Yang memilih adalah pemilik keputusan arsitektur, dan bahannya lengkap begitu ADR-0028 dikonfirmasi — karena K1b menyentuh alasan yang sama.~~ **Dipilih K1a 19 September 2026 — lihat ADR-0016 bagian Penyelesaian K1.** ADR-0028 sudah `accepted`, dan K1b **tidak dipilih**, sehingga peninjauan ADR-0023 yang menjadi konsekuensinya tidak pernah terjadi.

Menyentuh: **ADR-0016**, **ADR-0023**, **ADR-0028**, **REQ-021**, **D20**.

---
