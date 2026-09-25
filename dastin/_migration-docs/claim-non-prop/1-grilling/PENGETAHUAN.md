# PENGETAHUAN — Claim Non Prop

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (bertanggal 2026-09-08 s/d 2026-09-09, rule termutakhir di dalamnya `pxUpdateDateTime = 2026-08-30`); `pengetahuan/DDL_Script_ClaimNonProp.xls` **versi 2026-09-18 10:36, 48 objek**, terurai ke `pengetahuan/ddl/` (49 berkas); dan `pengetahuan/ddl/TABLE_PC_ASM_FW_GCNMFW_WORK.sql` (ditempel pengguna 2026-09-18).
> Disusun 18 September 2026 dari artefak yang sudah ada di folder ini. **Bukan tulisan baru** — tidak ada satu kalimat pun di sini yang belum melewati pemeriksaan sebelumnya.
> Berlaku hanya untuk keadaan sistem pada ekspor di atas. Perubahan sesudahnya tidak tercermin; deteksinya lewat sapuan ulang (ADR-0012).

Berkas ini menggabungkan **41 artefak** menjadi satu bacaan. Masing-masing tetap hidup sebagai berkas tersendiri; bila keduanya berbeda, **berkas aslinya yang berlaku**, karena di sanalah perubahan berikutnya ditulis.

Stempel asal yang semula berdiri di kepala tiap artefak dihapus di sini karena isinya sama persis dengan yang di atas.

> **Cuplikan ini tertanggal 18 September 2026 dan sudah usang di dua tempat — dicatat 19 September 2026.**
>
> | Bagian di berkas ini | Keadaan sebenarnya hari ini |
> |---|---|
> | **OPEN QUESTIONS** | Register **ditutup seluruhnya** 19 September 2026: aktif **0**, selesai **57**, dihapus **2**. Cuplikan di sini masih memperlihatkan bagian A–K berisi, dan tidak memuat satu pun dari D30–D46, E12–E17, maupun bagian K |
> | **ORACLE REQUESTS** | REQ: aktif **35**, selesai 1, **mati 1** (REQ-004 dicabut). Cuplikan di sini berbunyi "aktif 31" |
> | **REQ-032** | **VERIFIKASI** sejak 19 September 2026, bukan BLOCKER. Batas 30 byte ditetapkan sebagai aturan tetap; **tabel singkatan tertutup dibatalkan** |
>
> Aturan yang sudah berlaku di kepala berkas ini tetap berlaku dan justru inilah kejadiannya: **bila keduanya berbeda, berkas aslinya yang berlaku.** Cuplikan tidak diperbarui satu per satu — ia dibangun ulang, atau dibaca sebagai potret bertanggal.

## Yang TIDAK ada di dalam berkas ini

| Hal | Sebab |
|---|---|
| Modul **Komite Claim Non Prop** | foldernya tidak pernah dibuka, atas instruksi |
| Data produksi | tidak pernah ada satu baris pun; seluruh angka di sini berasal dari XML dan DDL |
| Rancangan sistem baru | belum ada DDL, belum ada Golang, belum ada React |
| Jawaban dari lima kelompok dan dari basis data | belum datang — daftarnya di bagian V |

## Daftar isi

- **[I. Bahasa — istilah yang dipakai di seluruh dokumen](#i-bahasa-istilah-yang-dipakai-di-seluruh-dokumen)**
  - [Claim Non Prop](#claim-non-prop) — `CONTEXT.md`
- **[II. Sistem lama — bagaimana ia bekerja](#ii-sistem-lama-bagaimana-ia-bekerja)**
  - [BLUEPRINT — Claim Non Prop (fakta teknis)](#blueprint-claim-non-prop-fakta-teknis) — `BLUEPRINT.md`
- **[III. Temuan — apa yang ditemukan di dalamnya](#iii-temuan-apa-yang-ditemukan-di-dalamnya)**
  - [FINDING-001 — Ambang kewenangan Komite dibandingkan tanpa konversi mata uang](#finding-001-ambang-kewenangan-komite-dibandingkan-tanpa-konversi-mata-uang) — `FINDING-001-threshold-currency.md`
  - [FINDING-002 — Alur bercabang berdasarkan identitas orang](#finding-002-alur-bercabang-berdasarkan-identitas-orang) — `FINDING-002-percabangan-identitas.md`
  - [FINDING-003 — Baris Retensi Cedant hadir di `SpreadingRisk` saat rule tanpa penyaring membacanya](#finding-003-baris-retensi-cedant-hadir-di-spreadingrisk-saat-rule-tanpa-penyaring-membacanya) — `FINDING-003-baris-ur-tanpa-penyaring.md`
  - [FINDING-004 — Hasil suntingan manual alokasi tidak terlindungi](#finding-004-hasil-suntingan-manual-alokasi-tidak-terlindungi) — `FINDING-004-suntingan-manual-tertimpa.md`
  - [FINDING-005 — Tanggal kerugian dibandingkan terhadap tanggal akhir treaty dalam format berbeda](#finding-005-tanggal-kerugian-dibandingkan-terhadap-tanggal-akhir-treaty-dalam-format-berbeda) — `FINDING-005-perbandingan-tanggal-beda-format.md`
  - [FINDING-006 — Kurs yang tidak ditemukan dikembalikan sebagai 1](#finding-006-kurs-yang-tidak-ditemukan-dikembalikan-sebagai-1) — `FINDING-006-kurs-tidak-ditemukan-bernilai-satu.md`
  - [FINDING-007 — Dua rumus premi reinstatement bekerja atas nilai yang berbeda](#finding-007-dua-rumus-premi-reinstatement-bekerja-atas-nilai-yang-berbeda) — `FINDING-007-dua-rumus-reinstatement.md`
- **[IV. Keputusan — apa yang sudah diputuskan untuk sistem baru](#iv-keputusan-apa-yang-sudah-diputuskan-untuk-sistem-baru)**
  - [Aggregate root adalah Klaim, satu klaim satu kejadian kerugian](#aggregate-root-adalah-klaim-satu-klaim-satu-kejadian-kerugian) — `docs/adr/0001-satu-klaim-satu-kejadian.md`
  - [Klaim yang sudah ditutup tidak dapat dibuka kembali](#klaim-yang-sudah-ditutup-tidak-dapat-dibuka-kembali) — `docs/adr/0002-tanpa-reopen.md`
  - [Presisi tinggi sepanjang rantai perhitungan, pembulatan hanya di tepi](#presisi-tinggi-sepanjang-rantai-perhitungan-pembulatan-hanya-di-tepi) — `docs/adr/0003-presisi-dan-pembulatan.md`
  - [Tambalan per-case tidak ikut dimigrasi; angkanya dipindahkan sebagai data](#tambalan-per-case-tidak-ikut-dimigrasi-angkanya-dipindahkan-sebagai-data) — `docs/adr/0004-tambalan-per-case-tidak-dimigrasi.md`
  - [Claim dan Komite satu unit cutover, dengan shadow-run sebagai bukti paritas](#claim-dan-komite-satu-unit-cutover-dengan-shadow-run-sebagai-bukti-paritas) — `docs/adr/0005-strategi-cutover.md`
  - [RBAC sistem baru dirancang dari nol, bukan diwarisi](#rbac-sistem-baru-dirancang-dari-nol-bukan-diwarisi) — `docs/adr/0006-rbac-dirancang-dari-nol.md`
  - [Setiap nilai uang disimpan berpasangan, dan ambang kewenangan dibandingkan terhadap nilai IDR](#setiap-nilai-uang-disimpan-berpasangan-dan-ambang-kewenangan-dibandingkan-terhadap-nilai-idr) — `docs/adr/0007-nilai-uang-berpasangan-dan-ambang-idr.md`
  - [Nilai hasil suntingan manual bertahan terhadap hitung ulang, dan selalu terlihat](#nilai-hasil-suntingan-manual-bertahan-terhadap-hitung-ulang-dan-selalu-terlihat) — `docs/adr/0008-suntingan-manual-bertahan-dan-terlihat.md`
  - [Keputusan alur disimpan sebagai field terstruktur; komentar tidak pernah dibaca mesin](#keputusan-alur-disimpan-sebagai-field-terstruktur-komentar-tidak-pernah-dibaca-mesin) — `docs/adr/0009-keputusan-terstruktur-bukan-teks-bebas.md`
  - [Layer dan Retensi Cedant adalah dua entitas terpisah](#layer-dan-retensi-cedant-adalah-dua-entitas-terpisah) — `docs/adr/0010-layer-dan-retensi-cedant-dua-tabel.md`
  - [Menjalankan ulang perhitungan alokasi harus menghasilkan keadaan yang identik](#menjalankan-ulang-perhitungan-alokasi-harus-menghasilkan-keadaan-yang-identik) — `docs/adr/0011-hitung-ulang-idempoten.md`
  - [Tambalan baru dideteksi lewat sapuan ulang, bukan lewat register](#tambalan-baru-dideteksi-lewat-sapuan-ulang-bukan-lewat-register) — `docs/adr/0012-deteksi-tambalan-baru-lewat-sapuan-bukan-register.md`
  - [Perhitungan dijalankan sebagai turunan dari data, bukan sebagai akibat penekanan tombol](#perhitungan-dijalankan-sebagai-turunan-dari-data-bukan-sebagai-akibat-penekanan-tombol) — `docs/adr/0013-perhitungan-sebagai-turunan-data.md`
  - [Aturan penguraian teks ke angka, dan perlakuan atas kurs yang tidak ada](#aturan-penguraian-teks-ke-angka-dan-perlakuan-atas-kurs-yang-tidak-ada) — `docs/adr/0014-penguraian-teks-ke-angka-dan-kurs-kosong.md`
  - [Akseptasi adalah satu entitas dengan keadaan, bukan dua tabel](#akseptasi-adalah-satu-entitas-dengan-keadaan-bukan-dua-tabel) — `docs/adr/0015-satu-entitas-akseptasi-dengan-keadaan.md`
  - [Data pegawai masuk lewat API atau salinan tersinkron, bukan database link](#data-pegawai-masuk-lewat-api-atau-salinan-tersinkron-bukan-database-link) — `docs/adr/0016-tanpa-database-link-ke-sistem-lain.md`
  - [Data aplikasi hanya diubah lewat aplikasi](#data-aplikasi-hanya-diubah-lewat-aplikasi) — `docs/adr/0017-satu-pintu-tulis.md`
  - [Klaim yang terdampak penjaga tanggal yang salah dimigrasi apa adanya, dan didaftar](#klaim-yang-terdampak-penjaga-tanggal-yang-salah-dimigrasi-apa-adanya-dan-didaftar) — `docs/adr/0018-klaim-terdampak-penjaga-tanggal-dimigrasi-apa-adanya.md`
  - [NULL bukan nol, dan keadaan perhitungan disimpan di tingkat baris](#null-bukan-nol-dan-keadaan-perhitungan-disimpan-di-tingkat-baris) — `docs/adr/0019-null-bukan-nol-dan-keadaan-di-tingkat-baris.md`
  - [Fakultatif dan Treaty dua entitas berbeda; sistem ini hanya memiliki Treaty](#fakultatif-dan-treaty-dua-entitas-berbeda-sistem-ini-hanya-memiliki-treaty) — `docs/adr/0020-hanya-treaty-yang-dimiliki.md`
  - [Pengenal polis dan klaim diberi nama menurut isinya; CASEID dipensiunkan](#pengenal-polis-dan-klaim-diberi-nama-menurut-isinya-caseid-dipensiunkan) — `docs/adr/0021-penamaan-pengenal-polis-dan-klaim.md`
  - [Masa berlaku treaty disimpan sebagai tanggal, dan batasnya inklusif](#masa-berlaku-treaty-disimpan-sebagai-tanggal-dan-batasnya-inklusif) — `docs/adr/0022-masa-berlaku-treaty-sebagai-tanggal-dan-batas-inklusif.md`
  - [Data akseptasi disimpan dalam satu bentuk kanonik; sistem hilir diberi view, bukan salinan](#data-akseptasi-disimpan-dalam-satu-bentuk-kanonik-sistem-hilir-diberi-view-bukan-salinan) — `docs/adr/0023-satu-bentuk-kanonik-tanpa-salinan.md`
  - [Satu akseptasi per klaim per layer per mata uang](#satu-akseptasi-per-klaim-per-layer-per-mata-uang) — `docs/adr/0024-kunci-alami-akseptasi.md`
  - [Tanggal tutup buku disimpan bertanggal berlaku, bukan satu baris tanpa riwayat](#tanggal-tutup-buku-disimpan-bertanggal-berlaku-bukan-satu-baris-tanpa-riwayat) — `docs/adr/0025-tutup-buku-bertanggal-berlaku.md`
  - [Batas kepemilikan mengikuti nama class: -Work- dan -Data- dimiliki, -Int- tidak](#batas-kepemilikan-mengikuti-nama-class--work--dan--data--dimiliki--int--tidak) — `docs/adr/0026-batas-kepemilikan-mengikuti-nama-class.md`
  - [Dokumen klaim dirujuk, tidak disimpan ulang](#dokumen-klaim-dirujuk-tidak-disimpan-ulang) — `docs/adr/0027-dokumen-dirujuk-bukan-dimiliki.md`
  - [Basis data tujuan Oracle, skema baru bersebelahan dengan POOLDATA](#basis-data-tujuan-oracle-skema-baru-bersebelahan-dengan-pooldata) — `docs/adr/0028-basis-data-tujuan-dan-letak-skema.md`
- **[V. Yang masih terbuka](#v-yang-masih-terbuka)**
  - [OPEN QUESTIONS — Claim Non Prop](#open-questions-claim-non-prop) — `_selesai/OPEN-QUESTIONS.md`
  - [ORACLE REQUESTS — Claim Non Prop](#oracle-requests-claim-non-prop) — `ORACLE-REQUESTS.md`
  - [Pertanyaan untuk Akuntansi](#pertanyaan-untuk-akuntansi) — `ASK-AKUNTANSI.md`
- **[VI. Bahan mentah — dari mana semua ini datang](#vi-bahan-mentah-dari-mana-semua-ini-datang)**
  - [pengetahuan/](#pengetahuan) — `pengetahuan/README.md`


---

# I. Bahasa — istilah yang dipakai di seluruh dokumen

Glosarium. Setiap istilah di bagian-bagian berikutnya memakai arti yang ditetapkan di sini.

## Claim Non Prop

Sumber: `CONTEXT.md`

Konteks penanganan klaim atas kontrak reasuransi treaty non-proporsional (excess of loss), sejak klaim dilaporkan cedant sampai klaim ditutup atau ditolak.

### Language

#### Pihak

**Cedant**:
Perusahaan asuransi yang menyerahkan sebagian risikonya kepada reasuradur dan yang melaporkan klaim.
_Avoid_: Ceding, Ceding Co, tertanggung

**Reasuradur**:
Pihak yang menanggung sebagian risiko cedant berdasarkan kontrak treaty. Dalam konteks ini selalu berarti Nusantara Re.
_Avoid_: RNM, Nusare

**Tertanggung**:
Pihak yang namanya tercantum pada polis asli milik cedant. Bukan pihak yang berhubungan langsung dengan reasuradur.
_Avoid_: Insured, nasabah

**Adjuster**:
Pihak independen yang ditunjuk untuk menilai besaran kerugian.
_Avoid_: Loss adjuster, surveyor, konsultan

#### Kontrak

**Treaty Non-Proporsional**:
Kontrak reasuransi di mana reasuradur menanggung kerugian hanya pada bagian yang melampaui batas tertentu, bukan berdasarkan porsi tetap.
_Avoid_: XOL, Excess of Loss, non-prop

**Layer**:
Satu lapisan proteksi dalam treaty non-proporsional, dibatasi oleh batas bawah dan batas atas sendiri. Kerugian mengisi layer secara berurutan dari bawah.
_Avoid_: Lapisan, tingkat

**Retensi Cedant**:
Bagian pertama dari setiap kerugian yang menjadi tanggungan cedant sendiri sebelum layer mana pun bekerja.
_Avoid_: UR, Underlying Retention, Deductible, Priority, Mindep

**Limit Layer**:
Kapasitas maksimum yang dapat diserap satu layer atas satu kerugian.
_Avoid_: Limit, kapasitas

**Premi Deposit**:
Premi minimum yang telah dibayarkan di muka untuk satu layer, dan menjadi dasar perhitungan premi pemulihan.
_Avoid_: MDP, Minimum and Deposit Premium

**Premi Pemulihan**:
Premi tambahan yang terutang ketika sebuah layer terpakai oleh suatu kerugian, untuk memulihkan kapasitas layer tersebut.
_Avoid_: Reinstatement, RIP, reinstatement premium

**Porsi Reasuradur**:
Persentase keikutsertaan reasuradur atas treaty yang bersangkutan.
_Avoid_: RNM Share, share, participation

#### Klaim

**Klaim**:
Satu kejadian kerugian atas satu polis treaty non-proporsional. Satu klaim dapat menghasilkan banyak pembayaran, tetapi hanya memiliki satu tanggal kejadian.
_Avoid_: Case, loss, kerugian

**Tanggal Kejadian**:
Tanggal terjadinya peristiwa yang menimbulkan kerugian. Bersifat tetap sepanjang umur klaim.
_Avoid_: DOL, Date of Loss

**Penyebab Kerugian**:
Klasifikasi peristiwa yang menimbulkan kerugian, dipilih dari daftar induk berjenjang.
_Avoid_: Cause of Loss, COL

**Alokasi Kerugian**:
Proses membagi satu nilai kerugian ke Retensi Cedant lalu ke setiap Layer secara berurutan sampai nilai kerugian habis.
_Avoid_: Loss allocation, spreading, alokasi XOL

**Penyebaran**:
Pembagian lebih lanjut atas bagian yang ditanggung reasuradur kepada pihak-pihak penerima berikutnya.
_Avoid_: Spreading, retrosesi, break QS

**Salvage**:
Nilai yang dapat dipulihkan kembali dari objek yang rusak, yang mengurangi besaran kerugian.
_Avoid_: Recovery, nilai sisa

**Biaya Penilaian**:
Biaya jasa adjuster yang ditambahkan pada besaran kerugian.
_Avoid_: Adjuster fee

**Prorata Klaim**:
Faktor pengali yang diterapkan ketika periode pertanggungan atau nilai pertanggungan tidak penuh.
_Avoid_: PctProrateClaim, prorate

#### Proses

**Registrasi**:
Tahap pencatatan klaim yang dilaporkan cedant beserta data polis, objek, dan nilai estimasi kerugian.
_Avoid_: Outstanding Claim, input outstanding

**Akseptasi**:
Tahap penetapan nilai yang disetujui untuk dibayar atas suatu klaim.
_Avoid_: Acceptance, Input Acceptation

**Adjustment**:
Satu transaksi pembayaran atau penyesuaian nilai di bawah sebuah klaim. Satu klaim dapat memiliki banyak Adjustment.
_Avoid_: Transaksi, payment, penyesuaian

**Nomor Akseptasi**:
Penanda unik yang diterbitkan ketika sebuah Adjustment telah disetujui seluruh jenjang persetujuan.
_Avoid_: AcceptedNo, no akseptasi

**Komite Klaim**:
Forum persetujuan berjenjang yang memutuskan diterima atau ditolaknya sebuah Adjustment. Kewenangan tiap jenjang ditentukan oleh besaran nilai.
_Avoid_: Committee, komite, approval

**Subjectivity**:
Persetujuan bersyarat, di mana Adjustment disetujui dengan syarat yang harus dipenuhi kemudian.
_Avoid_: Bersyarat, conditional approval

**Kronologi**:
Catatan berurutan atas setiap tindakan penting terhadap klaim, berisi pelaku, perannya, waktu, dan keterangan.
_Avoid_: History, audit trail, suggest list

**Penutupan Klaim**:
Pengakhiran klaim karena seluruh kewajiban telah diselesaikan.
_Avoid_: Close Claim, closing

**Penolakan Klaim**:
Pengakhiran klaim karena klaim dinyatakan tidak dijamin. Berbeda dari Penutupan Klaim.
_Avoid_: Reject Claim, ditolak

#### Dokumen

**Nota Kerugian**:
Dokumen resmi yang memuat rincian perhitungan kerugian dan bagian reasuradur.
_Avoid_: PLA, DLA, CFS, CA, claim advice

---

### Batas kepemilikan

**Treaty** — dimiliki dan dimodelkan sepenuhnya oleh sistem ini.

**Fakultatif** — **tidak dimiliki**. Dibaca sebagai rujukan bila diperlukan, tidak dimodelkan, tidak ditulis. Entitas yang berbeda dari Treaty: dinegosiasikan per risiko satu per satu, sedangkan Treaty adalah kontrak payung atas satu portofolio. Lihat ADR-0020.

---

### Pengenal — nama yang dipakai dan yang dihindari

**Nomor Polis** — pengenal polis milik kita.
_Avoid_: `CASEID`, `NOPOLIS`, `PolicyNo` tanpa keterangan pemilik.

**Nomor Polis Cedant** — pengenal polis yang sama menurut catatan cedant. Berbeda benda dari Nomor Polis, dan keduanya disimpan berdampingan.
_Avoid_: `POLICY_CEDING` sebagai istilah percakapan.

**Nomor Klaim** — pengenal klaim di sistem ini.
_Avoid_: `CASEID`, `pyID`, `case key`.

> **`CASEID` dipensiunkan.** Nama itu dipakai untuk dua benda berbeda di sistem lama — identitas polis di `V_POLIS`, identitas case Pega di `OS_AKSEPTASI_KLAIM`. Ia tidak dipakai lagi dalam percakapan maupun model internal. Pengecualiannya hanya payload integrasi ke Arasapas dan kasir, yang kontraknya tidak berubah (ADR-0021).

Sedikitnya lima bentuk pengenal polis hidup berdampingan di sistem lama, dan sedikitnya dua di antaranya milik pihak lain. Petanya di ADR-0021; pemetaan mana yang sama dan mana yang berbeda belum selesai.

---

### Polis

**Polis** — entitas tersimpan, bukan pandangan gabungan atas kontrak treaty dan tahun produksi.

Dasarnya perilaku pengguna, bukan bentuk penyimpanan: nomor polis **dimasukkan atau dipilih manusia** di layar (`Section/InputAcceptation.xml`, `Section/OutstandingClaim(1).xml`, `Harness/OutstandingClaim.xml`), dan `.ClaimData.PolicyData.PolicyNo` bernilai tunggal. Sesuatu yang punya nomor, dipilih manusia, dan dirujuk dari klaim adalah entitas — bila ia sekadar hasil penggabungan, penggunanya tidak akan memasukkan nomornya.

Bahwa sistem lama menyimpannya sebagai dokumen JSON di `json_polis` dan memandangnya lewat view `V_POLIS` adalah **keputusan penyimpanan**, bukan pernyataan tentang wujudnya.

**Nomor Polis** adalah kandidat kunci alaminya.

Label: **EVIDENCED** untuk keberadaannya; **DECIDED** untuk perlakuannya sebagai entitas tersimpan.

---

### Adjustment

**Adjustment** — **unit transaksi pembayaran klaim.** Satu klaim dapat memiliki banyak Adjustment, yaitu pembayaran bertahap.

Bukan revisi atas estimasi, dan bukan penambahan nilai klaim. Isinya menegaskan sifatnya: `Payable`, `PayableTo`, `NoAccount`, `DirectToKasir`, `StatusKasir`. Dan `PaymentType` membedakan jenisnya — `1` Final, `2` Partial/interim.

Definisi ini diambil dari `MEMORI_PEMAHAMAN.MD` §4.3 dan §7.4, yang merupakan otoritas tertinggi untuk maksud bisnis.
_Avoid_: menyebutnya "revisi", "koreksi", atau "penambahan klaim" — ketiganya menyesatkan ke model data yang berbeda.

Label: **EVIDENCED**.

---

### Kepemilikan — aturan tunggal

Batas kepemilikan mengikuti **nama class**, bukan pertimbangan per tabel:

| Pola | Kepemilikan |
|---|---|
| `…-Work-…`, `…-Data-…` | **dimiliki** |
| `…-Int-…` | **dibaca, tidak dimiliki** — pemetaan langsung ke tabel/view Oracle |

Modul klaim memiliki **klaim, akseptasi, alokasi, dan adjustment**. Polis, treaty, cedant, mata uang, dokumen, dan daftar Komite dibaca dari modul lain. Lihat ADR-0026.

---

### Dokumen

**Dokumen klaim** — dirujuk, tidak dimiliki. Berkasnya berada di Google Cloud Storage; yang tersimpan di basis data hanya metadata dan URL berbatas waktu. Lihat ADR-0027.


---

# II. Sistem lama — bagaimana ia bekerja

Pemahaman atas sistem Pega yang berjalan sekarang, disusun dari 279 berkas XML dan DDL Oracle. Bukan rancangan sistem baru.

## BLUEPRINT — Claim Non Prop (fakta teknis)

Sumber: `BLUEPRINT.md`

Scope: `D:\XML_NURE\Claim Non Prop` saja. Folder `Komite Claim Non Prop` tidak dibuka — lihat §8.
Label: **EVIDENCED** = terbukti dari XML · **DECIDED** = keputusan kita (lihat ADR) · **EXTERNAL** = butuh data luar (lihat _selesai/OPEN-QUESTIONS.md).

---

### 1. Inventaris Rule — EVIDENCED

279 berkas XML + 2 non-XML (`Struktur_Flow_TreatyIn.xlsx`, `~$Struktur_Flow_TreatyIn.xlsx`).

| Jml | Tipe Rule | Folder |
|---:|---|---|
| 114 | `RULE-OBJ-ACTIVITY` | `Activity\` |
| 50 | `RULE-CONNECT-SQL` | `RDBList\` |
| 36 | `RULE-HTML-SECTION` | `Section\` |
| 22 | `RULE-OBJ-REPORT-DEFINITION` | `ReportDefinition\` |
| 16 | `RULE-OBJ-FLOWACTION` | `FlowAction\` |
| 14 | `RULE-HTML-HARNESS` | `Harness\` |
| 10 | `RULE-OBJ-MODEL` (Data Transform) | `DataTransform\` |
| 7 | `RULE-CONNECT-REST` | `ConnectREST\` |
| 7 | `RULE-OBJ-WHEN` | `When\` |
| 1 | `RULE-DECLARE-DECISIONTABLE` | `DecisionTable\` |
| 1 | `RULE-OBJ-FLOW` | `Flow\` |
| 1 | `RULE-ADMIN-SYSTEM-SETTINGS` | `SystemSettings\` |

**Tidak ada di export**: `Rule-Obj-Class`, `Rule-Obj-Property`, `Rule-Declare-Expressions`, `Rule-Declare-Pages`, `Rule-Declare-OnChange`, `Rule-Access-When`, `Rule-Obj-FieldValue`, `Rule-Obj-Corr`, `Rule-Obj-HTML`.

**Applies-to class** (48 class): `…Work-ClaimTreatyNonProp` 84 · `…Data-Adjustment` 27 · `ASM-FW-GISFW-Data-PolicyTreatyIn` 19 · `@baseclass` 18 · `…GCNMFW-Work` 16 · `…Int-V_POLIS` 14 · `…Int-T_STORAGE_IMAGE` 10 · `…Work-ClaimTreaty` 7 · sisanya ≤6.

**Graf pemanggilan** (sumber: `Struktur_Flow_TreatyIn.xlsx`): 261 simpul, 387 edge, akar tunggal `Flow_TreatyIn`. **Orphan = 0** — setiap rule terpakai.

---

### 2. Struktur Klaim — EVIDENCED

#### 2.1 PageList langsung di bawah `.ClaimData` (calon tabel anak)

| PageList | Ref | Isi |
|---|---:|---|
| `SpreadingRisk` | 164 | Hasil alokasi per Layer (termasuk baris Retensi Cedant bertanda `TreatyName="UR"`) |
| `ReceiverClaim` | 101 | Rekening penerima pembayaran |
| `AdjustmentList` | 48 | Transaksi Adjustment |
| `ListClaimAmount` | 37 | Nilai kerugian per mata uang |
| `CNPSpreadLoss` | 23 | Pembagian kerugian per porsi |
| `InterestList` | 19 | Objek pertanggungan |
| `SpreadingClaim` | 19 | Penyebaran ke penerima berikutnya |
| `ListTotalEstimation` | 17 | Rekap total per mata uang |
| `EstimationList` | 14 | Estimasi awal |
| `ReinstatementList` | 13 | Perhitungan Premi Pemulihan |
| `SpreadingAdjustment` | 13 | Agregasi penyebaran (XOL) |
| `SpreadingAdjustmentQS` | 12 | Agregasi penyebaran (quota share) |
| `SpreadingBreakQS` | 11 | Rincian per tipe reasuransi |
| `TotalInterestInsured` | 7 | Rekap nilai pertanggungan |
| `Attachment` | 7 | Lampiran |
| `ObjectList` | 7 | Objek |
| `ListClaimAcceptation` | 4 | Rekap akseptasi |
| `ClaimComitee` | 2 | Komite level klaim (jalur tutup/tolak) |

**18 PageList.** Selebihnya (± 100 properti) skalar — daftar lengkap ada di hasil sapuan.

#### 2.2 Tanggal Kejadian bersifat tetap

Sapuan seluruh 114 activity: **tidak ada satu pun `Property-Set` yang menulis `.DateOfLoss`.** Nilai hanya masuk lewat input layar (`Section\OutstandingClaim.xml`, kontrol `pxDateTime`, `REQ`).
→ Mendukung penetapan satu klaim = satu Tanggal Kejadian.

#### 2.3 Natural key — EXTERNAL

**TIDAK DITEMUKAN DI XML**: tidak ada rule yang mencegah duplikat kombinasi `PolicyNo` + `DateOfLoss` + `IDMaster`.
`CheckDateDOL_Act` memang menjalankan `CekHistoryClaimNonProp_SQL` dan `CariHistoryClaim_SQL` untuk **menampilkan riwayat klaim atas polis yang sama**, dan `RejectedClaim_RD` untuk mengecek klaim yang pernah ditolak — tetapi keduanya bersifat informasi, bukan penolakan. *Disimpulkan dari ketiadaan rule penolakan, bukan dari adanya rule.*
→ Apakah duplikat dicegah secara prosedural oleh manusia: **EXTERNAL**.

---

#### 2.4 Retensi Cedant adalah **baris di dalam `SpreadingRisk`**, bukan entitas terpisah — EVIDENCED

`CountLossAllocation_act` menambahkan Retensi Cedant ke PageList yang sama dengan Layer:

```
pyWorkPage.ClaimData.SpreadingRisk(<APPEND>).TreatyType      = "UR"
pyWorkPage.ClaimData.SpreadingRisk(<LAST>).TreatyName        = "UR"
pyWorkPage.ClaimData.SpreadingRisk(<LAST>).ClaimEstimation   = Local.UR
pyWorkPage.ClaimData.SpreadingRisk(<LAST>).ClaimAmountAdjust = Local.UR
```
`Activity\CountLossAllocation_act.xml` baris 7015–7120.

Konsekuensinya: **setiap loop atas `SpreadingRisk` ikut melihat baris Retensi Cedant**, kecuali loop itu menyaringnya sendiri. Dua ejaan penanda dipakai bergantian — `.TreatyName=="UR"` dan `.TreatyType=="UR"` — dan keduanya memang ditulis pada baris yang sama, jadi bukan kontradiksi.

Ada ejaan ketiga: `"Previously Calculated UR"` (`GenerateCACNP_Act` baris 7223, 7425). Nilai ini **tidak pernah ditulis** di folder ini, hanya dibaca — asal-usulnya **TIDAK DITEMUKAN DI XML**.

**Rule yang menyentuh `SpreadingRisk`, dihitung ulang 18 September 2026 — KOREKSI**

Tabel versi pertama mencantumkan "16–66 rujukan `SpreadingRisk`". **Angka itu salah.** Ia menghitung seluruh kemunculan teks `SpreadingRisk`, yang sebagian besar adalah **nama class di metadata langkah** (`pyStepsClassName`, `pxRuleClassName`, `pyStepsParentClass`) — bukan rujukan ke PageList.

| Rule | Kemunculan teks | Rujukan PageList **sebenarnya** | Bentuknya |
|---|---|---|---|
| `GenerateCFS_act` | 66 | **5** | menyalin ke `TempDataOutStanding.ClaimData.SpreadingRisk` |
| `CountTotalInsterest_Act` | 32 | **2** | `@SizeOfPropertyList(...)` sebagai penjaga `>0` |
| `CopyOldataCurr_act` | 23 | **2** | `@LengthOfPageList(...)>0` + satu `Property-Set` `repeat=EMBEDDED` |
| `CountSpreadingXOL` | 31 | **0** | — |
| `CountSpreading_Act` | 30 | **0** | — |
| `DeleteAkseptasi_Act` | 28 | **0** | — |
| `CountClaimTNP_Act` | 22 | **0** | — |
| `SetAccoutNo_Act` | 16 | **0** | — |
| `SendEmailKlaim` | 16 | **0** | — |

**Enam dari sembilan rule tidak menyentuh PageList itu sama sekali.** Kehadiran baris `"UR"` tidak dapat memengaruhi mereka, karena mereka tidak membacanya.

Tinggal tiga rule yang benar-benar membacanya, dan ketiganya sudah tuntas:

| Rule | Akibat baris `"UR"` |
|---|---|
| `CountTotalInsterest_Act` | **tidak ada** — cacahan hanya dipakai sebagai penjaga `>0` |
| `CopyOldataCurr_act` | **tidak ada** — `>0` sebagai penjaga; `Property-Set` mengubah `.Currency` pada setiap baris termasuk baris `"UR"`, dan baris itu memang membawa `Currency` sendiri, jadi konsisten |
| `GenerateCFS_act` | **ada** — baris `"UR"` tersalin ke struktur retrosesi bila ia hadir (lihat `FINDING-003` bagian 4.2) |

**Batas klaim yang berlaku sekarang**: satu rule terbukti terdampak bila baris `"UR"` hadir, dua terbukti tidak terdampak, enam tidak membacanya sama sekali. Apakah baris itu hadir saat `GenerateCFS_act` berjalan bergantung pada urutan tindakan pengguna — lihat `FINDING-003` bagian 3.3.


#### 2.5 `.IsEditClaim` — ditulis, tidak pernah dibaca — EVIDENCED

`IsEditClaim` adalah `Rule-Obj-Property` milik class `ASM-FW-GISFW-Data-SpreadingRisk` (`pxRuleObjClass=Rule-Obj-Property`), **bukan** `Rule-Obj-When`.

Seluruh kemunculannya di folder ini — lima, tidak lebih:

| Rule | Langkah | Aksi |
|---|---|---|
| `EditXOLAlokasi` | `pyStepsClassName=ASM-FW-GISFW-Data-SpreadingRisk` | `.IsEditClaim = 1` |
| `CountLossAllocation_act` | dua langkah | `SpreadingRisk(<LAST>).IsEditClaim = 0` |

**Tidak ada satu pun pembacaan**: tidak dipakai di `pyStepsPreCondParamsWhen`, tidak di `pyVisible` section mana pun, tidak di Report Definition, tidak di SQL. Bendera ini **hanya ditulis**.

Dua kemungkinan, dan XML folder ini tidak bisa memilih di antaranya: (a) pembacanya ada di modul Komite → `DEFERRED-TO-KOMITE-SESSION`; (b) bendera ini mati. *Disimpulkan dari ketiadaan pembacaan, bukan dari adanya rule.*

---

### 3. Status Klaim — EVIDENCED

#### 3.1 `CNPStatusCase` — enum lengkap **di dalam folder ini**

| Nilai | Ditulis oleh | Tahap |
|---|---|---|
| `"INPUT ACCEPTATION CLAIM"` | `InputOutStandingCTNP_PostAct` step 2 | selesai Registrasi → masuk Akseptasi |
| `"COMITEE ACCEPTANCE (DEPT. HEAD)"` | `CreateChildKomiteCNP_Act` step 31 | Adjustment dikirim ke Komite |
| `"COMITEE ACCEPTANCE (DEPT. HEAD)"` | `CreateChildKomiteCloseNP_Act` step 12 | Penutupan/Penolakan dikirim ke Komite |

Hanya **3 penulisan, 2 nilai unik** di folder ini. Nilai `"CLAIM ACCEPTED"` dan `"CLAIM REJECTED"` **tidak ditulis di folder ini** — ditulis modul Komite (DEFERRED-TO-KOMITE-SESSION).

**Tidak ada satu pun rule di folder ini yang membaca `CNPStatusCase`** — ketiga rule di atas hanya menulis. *Disimpulkan dari ketiadaan pembacaan.* Konsumennya kemungkinan UI atau modul Komite.

#### 3.2 Status case sebenarnya

| Fakta | Bukti |
|---|---|
| Satu-satunya status akhir | `Flow\Flow_TreatyIn.xml` shape `End1` → `pyWorkStatus = Resolved-Completed` |
| Satu-satunya penetapan status dari activity | `CloseClaimTNonProp` step 9 → `Call ASMForceCaseClose` dengan `WorkStatus = Resolved-Completed` |
| `Resolved-Rejected` / `Resolved-Withdrawn` | **TIDAK DITEMUKAN DI XML** pada folder ini |
| Jalur reopen | **TIDAK DITEMUKAN DI XML**. *Disimpulkan dari ketiadaan rule, bukan dari adanya rule.* |

`Open-By` yang muncul di `ASMForceCaseClose`, `AttachCAPDF_Act`, `AttachRISlipToWork`, `GetBase64Attachment`, `GetDetailPolis_act` adalah bagian nama method `Obj-Open-By-Handle`, **bukan** status.

---

### 4. Dua Jalur Pengakhiran Klaim — EVIDENCED

| | `CloseClaimNP` | `CloseClaimMD` |
|---|---|---|
| Pre-activity | `CloseClaimNP_preAct` | — |
| Section | `Section\CloseClaimNP.xml` | `Section\CloseClaimMD.xml` |
| Activity inti | `CreateChildKomiteCloseNP_Act` | `CloseClaimTNonProp` |
| Lewat Komite? | **Ya** — `pxAddChildWork`, `ChildClass=…Work-KomiteTreatyNonProp`, `FlowName=KomiteTreaty_Flow` | **Tidak** |
| Penjaga | `@hasMessages(myStepPage)` sebelum membuat child | **Ada**: langkah 1–3 memindai `AdjustmentList`; bila ada `.AcceptanceStatus=="0"` (masih menunggu Komite) → pesan `ErrPendAcc` → `EXIT-ACT-FAIL` |
| Simpan outstanding | (di modul Komite) | `SaveDataToOsAkseptasiNP` (step 5) |
| Kronologi | (di modul Komite) | `InsertChronology_DT` (step 6) |
| Update `json_klaim` | (di modul Komite) | `InsertJsonClaimTreaty_act` (step 7) |
| Kirim ke Arasapas | (di modul Komite) | `Connect-REST insertClaimFinalOrClosed_NP` (step 8) |
| Tutup case | (di modul Komite) | `ASMForceCaseClose` → `Resolved-Completed`, `CloseAllSubCases=true` (step 9) |
| Notifikasi | `SendEmailKlaimRejectClose` bila `IsPEGAPROD` | — |

#### 4.0 Status: **BYPASS BERSYARAT** — `CANDIDATE-NOT-MIGRATED`

Dugaan awal "bypass tanpa penjaga" saya cabut, lalu koreksinya sendiri saya turunkan lagi karena kelewatan. Yang terbukti dari XML:

**(a) Jalur MD menutup klaim tanpa pernah membuat case Komite.** Pemindaian `Activity\CloseClaimTNonProp.xml`, `Section\CloseClaimMD.xml`, dan `FlowAction\CloseClaimMD.xml` atas kata kunci `pxAddChildWork`, `CreateChildKomite*`, `KomiteTreaty_Flow`, `KomiteTreatyNonProp`, `ComiteeClaim`, `KomiteList`: **nihil di ketiganya**. Jalur ini langsung `ASMForceCaseClose` → `Resolved-Completed`.

**(b) Penjaganya hanya menangkap Adjustment yang sedang di Komite, bukan yang belum pernah ke Komite.** `.AcceptanceStatus` ditulis **satu kali saja di seluruh 279 berkas**: `CreateChildKomiteCNP_Act` langkah 31 → `= 0`. Perpindahan `0` → `1`/`2` tidak ada di folder ini (dilakukan modul Komite, DEFERRED).
Akibatnya nilai `.AcceptanceStatus` yang mungkin ada: **kosong** (belum pernah dikirim), `0` (sedang di Komite), `1`/`2` (sudah diputus).
Penjaga di `CloseClaimTNonProp` langkah 1 hanya menguji `.AcceptanceStatus=="0"`. Adjustment yang **belum pernah dikirim ke Komite** bernilai kosong, bukan `"0"`, sehingga **tidak tertangkap** dan penutupan tetap berjalan.

**(c) Niat penjaganya memang bukan otorisasi.** Pesan galatnya sendiri berbunyi `"Can not close claim, there is adjustment in comitee!"` — ia dirancang mencegah penutupan **saat Komite sedang bersidang**, bukan mensyaratkan persetujuan Komite.

**Kesimpulan yang terbukti**: `CloseClaimMD` adalah **bypass bersyarat** — melewati Komite sepenuhnya, tetapi menolak berjalan bila ada Adjustment yang sedang menggantung di Komite. Keputusan dimigrasi atau tidak menunggu REQ-006 (statistik pemakaian).

#### 4.1 Otorisasi pada jalur penutupan

Lihat §5.0 — ketiadaan otorisasi bukan khusus jalur ini, melainkan berlaku di seluruh modul.

---

### 5. Otorisasi — EVIDENCED

#### 5.0 TEMUAN UTAMA: tidak ada model otorisasi yang bisa diwarisi

> **Modul ini tidak memiliki otorisasi di tingkat rule sama sekali.**
>
> - `pyPrivilegeName` → **kosong di seluruh 279 berkas**. Tag-nya ada, nilainya tidak pernah diisi. Artinya `pyPrivilegeClass` yang terisi (36× `…Work-ClaimTreatyNonProp`, dst.) menunjuk ke privilege **tanpa nama** — tidak mengikat apa pun.
> - `Rule-Access-When` → **tidak ada satu pun di export**.
> - `pyVisible` → `ALWAYS` pada seluruh elemen UI yang diperiksa, termasuk kedua jalur penutupan klaim.
>
> Seluruh kontrol akses yang benar-benar berjalan bertumpu pada **access group dan routing flow**, dan keduanya **tidak ikut ter-export**.
>
> **Konsekuensi terhadap cakupan pekerjaan**: tidak ada model otorisasi yang dapat diwarisi dari modul ini. Seluruh RBAC sistem baru berstatus **DECIDED** — dirancang dari nol bersama pemilik proses — **bukan EVIDENCED**. Ini menambah pekerjaan yang belum ada di rencana mana pun. Lihat [ADR-0006](./docs/adr/0006-rbac-dirancang-dari-nol.md).

#### 5.1 Sisa sinyal yang ada di XML — EVIDENCED + EXTERNAL

| Sinyal | Isi |
|---|---|
| `pyPrivilegeClass` | `…Work-ClaimTreatyNonProp` 36, `…Data-Adjustment` 15, `ASM-FW-GISFW-Data-PolicyTreatyIn` 14, `…Work-ClaimTreaty` 13, **`ASM-FW-GCNMFW-Work-PNC` 9**, `…GCNMFW-Work` 9, `…Data-ObjectItem` 6, `@baseclass` 5 |
| `pyPrivilegeName` | **kosong seluruhnya** |
| `pyRuleSet` | `GCNMFW` 214, `GISFW` 72, `Pega-RULES` 10, `GCNMFWInt` 6, **`worts@` 4**, `Pega-IntegrationArchitect` 1 |
| `pyApplicationName` | `GCNMFW` 6, **`RaniRApp` 1** |
| `pyWorkBasket` / `pyWorkbasketName` | **TIDAK DITEMUKAN DI XML** |

**Label peran yang benar-benar ada di XML** — `DataTransform\InsertChronology_DT.xml`, properti `.IsCedingConfirm`:
`"Claim Admin"` (default) · `"Claim Dept. Head"` (bila pelaku `CHRISTINEANGELINA`) · `"Operational Director"` (bila `Himawan`) · `"Technical Director"` (bila `NANDINA`).
→ Empat peran ini **EVIDENCED**, tetapi pemetaannya ke orang bersifat hardcode (§7).

**Routing assignment** (`Flow\Flow_TreatyIn.xml`):
- `Assignment2` Registrasi → `pyImplementation = ToCurrentOperator`
- `Assignment1` Akseptasi → `pyImplementation = ToWorkbasket`, **nama workbasket tidak ada di XML** → EXTERNAL

Anomali yang perlu dijelaskan: `ASM-FW-GCNMFW-Work-PNC`, ruleset `worts@`, aplikasi `RaniRApp` — muncul di folder ini tetapi tidak punya rule sendiri. → OPEN-QUESTIONS.

---

### 6. Inventaris Aritmetika — EVIDENCED

Sapuan seluruh 114 activity, seluruh `PropertiesValue` yang mengandung operasi aritmetika dan menyentuh properti bernuansa nilai: **673 ekspresi**. Daftar lengkap: [`pengetahuan/arithmetic-inventory.tsv`](./pengetahuan/arithmetic-inventory.tsv) (kolom: rule, step, properti tujuan, ekspresi, skala, berkas).

#### 6.1 Temuan utama

**518 dari 673 ekspresi (77%) tidak memakai `@divide` sama sekali** — memakai operator `*`, `/`, `+`, `-` biasa, sehingga **tidak punya skala eksplisit** dan bergantung pada perilaku desimal bawaan Pega. Hanya 155 ekspresi (23%) yang menyatakan skala.

#### 6.2 Skala `@divide` yang dipakai

| Skala | Jml | Dipakai di |
|---:|---:|---|
| 20 | 126 | Seluruh rule inti klaim non-proporsional |
| 10 | 42 | `CountClaimTNP_Act`, `CountLossAllocation_act` |
| 4 | 10 | Rule sisi treaty/premi (`ASM-FW-GISFW-Data-PolicyTreatyIn`) |
| 8 | 4 | `SetPPNPPH` — konstanta pajak |
| 5 | 2 | `CountClaimTNP_Act` → `.PctProrateClaim` |
| 2 | 2 | `SetTPLNote_Act` — pembulatan tampilan |
| 0 | 2 | `SendEmailKlaim` — pembulatan tampilan |

#### 6.3 Polanya: skala mengikuti **modul**, bukan jenis nilai

| Kelompok | Skala | Rule |
|---|---:|---|
| Inti klaim non-prop | **20** | `AdjClaimCNP_Act` (28), `AdjClaimAmount_Act` (14), `CreateChildKomiteCNP_Act` (13), `SaveCNPLayerList_Act` (13), `SetActualPremium_ACT` (10), `GenerateCFS_act` (9), `SaveAdjustmentToOSAksep_Act_Tes` (5), `GeneratePlaCNP2_Act` (4), `SetCurrency_Act` (2), `GenerateCACNP_Act` (1) |
| Treaty / premi | **4** | `DetailCalculation` (3), `CountNetPremi_act` (2), `CountResult1_Act` (2), `SetValidateInstallment_Act` (2), `CountTotalInsterest_Act` (1) |
| Pajak & fee | **8** | `SetPPNPPH`: `@divide(2.5,100,8)` brokerage, `@divide(102.2,100,8)`, `@divide(2,100,8)` PPh, `@divide(2.2,100,8)` PPN |
| Tampilan | **0 / 2** | `SendEmailKlaim`: `@divide(.TotalSharePersen,1,0)`, `@divide(.SharePercentage,1,0)` · `SetTPLNote_Act`: `@divide(.TPLPct,1,2)`, `@divide(.TPLAmount,1,2)` — semuanya **bagi dengan 1**, murni pemformatan |

**Hanya 2 rule yang mencampur skala**: `CountClaimTNP_Act` (10×28, 20×13, 5×2) dan `CountLossAllocation_act` (10×14, 20×14 — tepat berimbang). Keduanya adalah rule alokasi inti. → OPEN-QUESTIONS.

#### 6.4 Tarif hardcoded di `SetPPNPPH` — EVIDENCED

`.BrokerageFee = @divide(2.5,100,8)` · `.BrokerageFeeSebenarnya = @divide(102.2,100,8)` · `.PPHValue = @divide(2,100,8)` · `.PPNValue = @divide(2.2,100,8)`
→ Tarif brokerage, PPh, dan PPN tertanam di kode. Bila tarif pajak berubah, kode harus diubah.

---

### 7. Tambalan Per-Case di Dalam Kode — EVIDENCED

Sapuan seluruh folder atas literal `CLMNP-…`, `pzInsKey` literal, `IDMaster` literal, dan angka mati.
**Hasil: 9 identitas berbeda, tersebar di 4 rule, 14 langkah.** (Dugaan awal hanya 3 — meleset.)

| # | Identitas | Rule | Step | Yang ditimpa / di-bypass |
|---|---|---|---|---|
| 1 | `CLMNP-975` | `AdjClaimCNP_Act` | 11.2 | `.CNPReinstatement = 881928.966808370` bila `.Currency=="IDR"` |
| 2 | `CLMNP-975` | `AdjClaimCNP_Act` | 11.3 | `.CNPReinstatement = 806851161.1895` bila `.Currency=="USD"` |
| 3 | `CLMNP-232` | `AdjClaimCNP_Act` | 4 | `.SpreadingAdjustment(3)` disalin dari `AdjustmentList(6).SpreadingAdjustment(3)` |
| 4 | `CLMNP-232` | `AdjClaimCNP_Act` | 5 | `.SpreadingQuotaShare(5)` disalin dari `AdjustmentList(7).SpreadingQuotaShare(5)` |
| 5 | `CLMNP-232` | `AdjClaimCNP_Act` | 15.5.1–15.5.4 | `Local.AdjusterFee` dipetakan manual per `IdxQS` (1/2→SA(1), 3/4→SA(2), 5→SA(3), 6→SA(4)) |
| 6 | `CLMNP-861` | `AdjClaimCNP_Act` | 10 | `.SpreadingRisk(IdxLastLayer).AdjusterFee = TotalClaim × ClaimPercentage/100` (rumus khusus) |
| 7 | `IDMaster 1000393` & `1000393/R01` | `CountLossAllocation_act` | 17.2.1 | `Local.LayerLimit` memakai urutan pengali `Kurs` yang berbeda |
| 8 | `CLMNP-367` & `CLMNP-382` (via `pzInsKey`) | `CountLossAllocation_act` | 17.2.1 | `Local.URLimit` memakai `Deductible2/Kurs`, bukan `Deductible` |
| 9 | `IDMaster 1001130` & `1001130/R01` | `GetHistoryMasterID_NP` | 4, 5 | `InputSpreading.CARI22` diisi literal master-id |
| 10 | `CLMNP-50` & `CLMNP-232` | `GetHistoryMasterID_NP` 7, `InputOutStandingClmTNP_PreAct` 5 | | menambahkan baris `"FACOUT"` ke `ListLossAllocation` |

Catatan: pada `InputOutStandingClmTNP_PreAct` step 5, `pyStepsDescription` berbunyi `pyWorkPage.pyID=="CLMNP-50"` sementara kondisi eksekusi sebenarnya `pyWorkPage.pyID=="CLMNP-232"` — **deskripsi dan kondisi tidak sinkron**, indikasi tambalan disalin-tempel.

**Angka mati lain yang bukan tambalan per-case** (ambang kewenangan Komite, `CreateChildKomiteCNP_Act` step 10 & 26.5): `Local.LimitMax = 30000000.00`, `Local.LimitMaxDivHead = 50000000.00`, `Param.LIMIT_BOTTOM = 25000001`.

---


#### 7.1 Kelas tambalan kedua: **identitas ORANG**, bukan hanya nomor case — EVIDENCED

Sapuan ini sebelumnya hanya mencari `CLMNP-…` dan `IDMaster`. Sapuan ulang atas seluruh `pyStepsPreCondParamsWhen` menemukan kelas yang terlewat: kondisi yang menyebut **nama operator dan nama orang**.

| Identitas | Rule | Langkah | Bentuk kondisi |
|---|---|---|---|
| `VINCENTVERNANDO_1` | `CreateChildKomiteCNP_Act` | 4 langkah | `pyWorkPage.pxCreateOperator=="VINCENTVERNANDO_1"` |
| `VINCENTVERNANDO_1` | `SaveDataToOSAksep_Act` | 1 | idem |
| `VINCENTVERNANDO_1` | `SaveToOS` | 1 | idem |
| `VINCENTVERNANDO_1` | `SendEmailKlaimRejectClose` | 2 | idem |
| `Himawan`, `NANDINA`, `CHRISTINEANGELINA`, `CHRISTOPMARHASAK` | `CreateChildKomiteCNP_Act` | 4 langkah | `.KomiteID=="<nama>"` |
| `Himawan`, `Nandina C`, `Christine Angelina Hutagalung` | `SethistoryKlaimTreaty` | 3 langkah | `.PICSuggest=="<nama>"` |
| `Himawan` | `SethistoryKlaimTreaty` | 4 langkah | `@contains(.CommentSuggest,"Accepted by Himawan")` / `"Rejected by Himawan"` |
| `Himawan` | `DataTransform\InsertChronology_DT` | 1 | `.PICSuggest=="Himawan"` |

Yang terakhir paling rapuh: keputusan alur ditentukan dengan **mencari potongan teks nama orang di dalam kolom komentar bebas**. Satu salah ketik dari pengguna mengubah jalur.

**Total gabungan per-case + per-orang: 29 langkah, 18 ekspresi unik, 8 rule.**

#### 7.2 Umur tambalan — EVIDENCED

Diambil dari `pxCreateDateTime`/`pxCreateOperator` pada blok langkah yang bersangkutan.

| Tanggal | Penulis | Tambalan |
|---|---|---|
| 2018-01-16 | `YOSUAAMBIKA` | `VINCENTVERNANDO_1` di `CreateChildKomiteCNP_Act` (4 langkah) |
| 2020-02-18 | `MESDISILITONGA` | 4 nama `KomiteID` |
| 2020-02-27 | `REVIRUNDUPADANG` | `CLMNP-50`, dan `CLMNP-232` di `InputOutStandingClmTNP_PreAct` |
| 2022-01-03 | `GABRIELAMILITIA` | `VINCENTVERNANDO_1` disalin ke `SaveDataToOSAksep_Act` dan `SaveToOS` |
| 2022-12-01 | `AnanSosmita` | 3 nama `PICSuggest` |
| 2023-08-30 | `AnanSosmita` | `CLMNP-232` di `AdjClaimCNP_Act` (6 langkah), `CLMNP-861` |
| 2023-11-20 | `ArlexyVarian` | `VINCENTVERNANDO_1` disalin ke `SendEmailKlaimRejectClose` |
| 2025-07-04 | `JEFRIHARI` | `IDMaster 1001130` dan `1001130/R01` |
| 2026-07-16 | `JEFRIHARI` | `CLMNP-975` (2 langkah) |

Tiga hal yang dibuktikan deret ini:

1. **`CLMNP-232` ditambal dua kali, berjarak 3,5 tahun, oleh dua orang berbeda, di dua rule berbeda** (2020-02-27 `REVIRUNDUPADANG`; 2023-08-30 `AnanSosmita`). Tambalan kedua tidak menggantikan yang pertama — keduanya masih aktif.
2. **Tambalan tertua berumur lebih dari 8 tahun** dan masih hidup, serta **disalin ke rule lain dua kali sesudahnya** (2022, 2023). Polanya menyebar, bukan mengendap.
3. Yang terbaru berjarak **dua bulan dari hari ini**. Praktik ini **belum berhenti**.

#### 7.3 Angka & alamat mati lain — EVIDENCED

| Nilai | Lokasi | Keterangan |
|---|---|---|
| `http://192.168.105.116:80/prweb` | `SaveToOS`, `SendEmailKlaim`, `SendEmailKlaimRejectClose` | alamat IP server tertanam di kode |
| `http://pega.nusantarare.com:80/prweb` | rule yang sama | alamat kedua, berdampingan dengan yang pertama |
| `"10026"` | `SetCurrency_Act` baris 1798, 2498 | `Param.IDCurr=="10026"` — id mata uang tertanam |
| `"EDITCLAIMXOL"`, `"CountXOL"`, `"AccEdit"`, `"ChildCase"` | tersebar | penanda alur berupa string bebas |

---

#### 7.4 `CLMNP-232` ditambal dua kali — dua gejala berbeda, bukan satu yang bertabrakan — EVIDENCED

Pertanyaan yang belum terjawab ronde lalu: apakah kedua tambalan menyentuh properti yang sama, sehingga yang menang bergantung urutan pemanggilan? **Tidak.** Keduanya menyentuh halaman yang berlainan.

| | Tambalan 2020-02-27 (`REVIRUNDUPADANG`) | Tambalan 2023-08-30 (`AnanSosmita`) |
|---|---|---|
| Rule | `InputOutStandingClmTNP_PreAct` | `AdjClaimCNP_Act` |
| Halaman yang disentuh | `ListLossAllocation.pxResults`, `ListTreaty.pxResults` | `Primary.SpreadingAdjustment`, `.SpreadingQuotaShare`, `Local.*` |
| Yang dilakukan | menyisipkan baris `"FACOUT"` dan `"OR"` ke daftar **masukan** | memetakan ulang nilai pada struktur **keluaran** |
| Jenis | menambah baris | menimpa nilai |

Contoh isi tambalan 2023:
```
.SpreadingAdjustment(3) = pyWorkPage.ClaimData.AdjustmentList(6).SpreadingAdjustment(3)
.SpreadingQuotaShare(5) = pyWorkPage.ClaimData.AdjustmentList(7).SpreadingQuotaShare(5)
Local.AdjusterFee       = Primary.SpreadingAdjustment(1..4).AdjusterFee   (dipetakan manual per indeks)
```

**Kesimpulan: tidak ada yang "menang".** Keduanya berjalan, pada tahap yang berbeda, atas halaman yang berbeda. Artinya pada satu klaim yang sama terdapat **dua gejala terpisah** — satu pada daftar masukan alokasi, satu pada hasil pemetaan Adjustment.

**Hipotesis yang belum diuji** (ditandai supaya tidak terbaca sebagai fakta): tambalan 2023 memperbaiki akibat dari tambalan 2020 — daftar masukan yang disisipi baris tambahan menghasilkan indeks yang bergeser, lalu indeks itu dipetakan ulang secara manual tiga tahun kemudian. Menguji ini memerlukan penelusuran pengaruh `ListLossAllocation` terhadap penomoran `SpreadingAdjustment`, dan **itu belum dikerjakan**.

Catatan terpisah: pada `InputOutStandingClmTNP_PreAct` langkah 5, `pyStepsDescription` berbunyi `pyWorkPage.pyID=="CLMNP-50"` sementara kondisi eksekusinya `pyWorkPage.pyID=="CLMNP-232"` — deskripsi dan kondisi tidak sinkron.

#### 7.5 Laju tambalan — EVIDENCED

Sumber: `pxCreateDateTime` / `pxCreateOperator` pada blok langkah yang bersangkutan.

| Tanggal | Penulis | Rule | Jenis | Langkah |
|---|---|---|---|---|
| 2018-01-16 | `YOSUAAMBIKA` | `CreateChildKomiteCNP_Act` | identitas akun | 4 |
| 2020-02-18 | `MESDISILITONGA` | `CreateChildKomiteCNP_Act` | identitas orang | 4 |
| 2020-02-27 | `REVIRUNDUPADANG` | `GetHistoryMasterID_NP` | nomor case | 1 |
| 2020-02-27 | `REVIRUNDUPADANG` | `InputOutStandingClmTNP_PreAct` | nomor case | 2 |
| 2022-01-03 | `GABRIELAMILITIA` | `SaveDataToOSAksep_Act` | identitas akun *(salinan)* | 1 |
| 2022-01-03 | `GABRIELAMILITIA` | `SaveToOS` | identitas akun *(salinan)* | 1 |
| 2022-12-01 | `AnanSosmita` | `SethistoryKlaimTreaty` | identitas orang | 3 |
| 2023-08-30 | `AnanSosmita` | `AdjClaimCNP_Act` | nomor case | 7 |
| 2023-11-20 | `ArlexyVarian` | `SendEmailKlaimRejectClose` | identitas akun *(salinan)* | 2 |
| 2025-07-04 | `JEFRIHARI` | `GetHistoryMasterID_NP` | ID master | 2 |
| 2026-07-16 | `JEFRIHARI` | `AdjClaimCNP_Act` | nomor case | 2 |

**Langkah baru per tahun:**

```
2018  ####                4
2019                      0
2020  #######             7
2021                      0
2022  ####                4
2023  #########           9
2024                      0
2025  ##                  2
2026  ##                  2   (per 16 Juli; ekspor diambil September)
```

Yang dapat dibaca dari deret ini:

1. **Tidak melambat sampai 2023, lalu turun.** Puncaknya 2023 (9 langkah, 2 penulis). Sesudahnya 2 langkah per tahun. Laju turun, tetapi **tidak berhenti**.
2. **Enam penulis berbeda dalam delapan tahun.** Tidak ada satu orang pun yang memegang pola ini; ia diwariskan.
3. **Menyebar, bukan mengendap.** Kondisi `pxCreateOperator=="VINCENTVERNANDO_1"` ditulis sekali pada 2018 lalu **disalin ke tiga rule lain** pada 2022 dan 2023 oleh dua orang yang berbeda dari penulis aslinya.
4. **Yang terbaru berjarak dua bulan dari hari ini.** Selama migrasi berjalan, tambalan baru masih mungkin bertambah ke sistem yang sedang dipetakan.

**Batas klaim:** deret ini dibaca dari satu ekspor bertanggal September 2026. Tambalan yang **dihapus** sebelum tanggal itu tidak meninggalkan jejak dan tidak terhitung di sini. Jadi angka ini adalah **batas bawah**, bukan hitungan lengkap. *Disimpulkan dari isi ekspor, bukan dari riwayat perubahan rule.*

---

### 8. Boundary Contract Draft — Claim → Komite (EVIDENCED)

Titik kopling yang memaksa kedua modul menjadi satu unit. Detail sisi Komite: **DEFERRED-TO-KOMITE-SESSION**.

#### 8.1 Pembuatan child case

| Pemicu | Rule | Parameter |
|---|---|---|
| Kirim Adjustment ke Komite | `CreateChildKomiteCNP_Act` step 29 | `pxAddChildWork`: `ChildClass=ASM-FW-GCNMFW-Work-KomiteTreatyNonProp`, `FlowName=KomiteTreaty_Flow`, `CopyPageData=true`, `Commit=true`, `UpdateHistory=true` |
| Kirim Penutupan/Penolakan ke Komite | `CreateChildKomiteCloseNP_Act` step 10 | idem |

#### 8.2 Data yang dikirim Claim → Komite

Disalin ke `ChildWorkPage` sebelum `pxAddChildWork`: `KomiteCount=1`, `KomiteLoop`, `CLMNO`, `TransferType`, seluruh `Adjustment` (`Page-Copy` dari `AdjustmentList(idx)`), `Adjustment.CNPLayerList[].CNPCurrencyList[]`, `KomiteList[]`, `Komite.CircumtansesCouseOfLoss`, `Komite.Remarks`, `Komite.Occupation`, `QuotationData.*`, `AttachmentStream` (Base64).

#### 8.3 Data yang ditunggu balik Komite → Claim

Ditulis balik oleh modul Komite ke `AdjustmentList(IndexObject)`: `AcceptanceStatus`, `AcceptedNo`, `AcceptedDate`, `ComiteeClaim[].KomiteAproval`, `.KomiteComment`, `.DateApproval`; ditambah `CNPStatusCase` dan `IsSubjectivity`/`SubjectivityNote` pada level klaim.

#### 8.4 Kunci relasi

| Arah | Kunci |
|---|---|
| Claim → Komite | `pyWorkPage.pxCoveredInsKeys(<LAST>)`; `.KomiteNo = @substring(pxCoveredInsKeys(<LAST>),19,30)` |
| Komite → Claim | `pxCoverInsKey` |
| Komite → Adjustment tertentu | `Adjustment.IndexObject` = **posisi numerik** di `AdjustmentList` — bukan surrogate key |

#### 8.4b Penghubung di tingkat data — EVIDENCED

Sumber: `pengetahuan/ddl/TABLE_PC_ASM_FW_GCNMFW_WORK.sql`. Tabel ini **tidak dimigrasi** — ia buatan Pega, dan sistem baru punya model sendiri. Yang diambil hanya apa yang tampak.

##### Batas antar modul adalah nilai satu kolom, bukan batas tabel

`pzInsKey` di Pega berbentuk `<class tabel> <pyID>`. Literal di XML memperlihatkan **tiga jenis case memakai nama class yang sama**:

```
ASM-FW-GCNMFW-WORK CLMNP-...     Claim Non Prop
ASM-FW-GCNMFW-WORK KMTNP-...     Komite Non Prop
ASM-FW-GCNMFW-WORK PNC-...       PNC
```

Ketiganya memakai class induk `ASM-FW-GCNMFW-Work` — yaitu tabel `PC_ASM_FW_GCNMFW_WORK`. Subclass `ClaimTreatyNonProp` (3.691 rujukan), `ClaimTreaty` (1.391), `PNC` (322), dan `KomiteTreatyNonProp` (56) **tidak punya tabel sendiri**.

Dikuatkan indeksnya: `PC_WORK_IDX1_2AED8 (PXCOVERINSKEY, PXOBJCLASS, PYSTATUSWORK)` menempatkan `PXOBJCLASS` tepat di tengah pola akses induk-anak.

> **Claim dan Komite duduk di baris-baris tabel yang sama, dibedakan `PXOBJCLASS`.** Batas antara dua modul bukan batas tabel, melainkan nilai satu kolom.

##### Penghubung induk-anak

| Kolom | Tipe | Peran |
|---|---|---|
| `PXCOVERINSKEY` | `VARCHAR2(255 CHAR)` | menyimpan `pzInsKey` case induk pada baris case anak |
| `PXCOVEREDCOUNT` | `NUMBER(18)` | cacah anak |
| `PXCOVEREDCOUNTOPEN` | `NUMBER(18)` | cacah anak yang masih terbuka |
| `PXCOVEREDCOUNTUNSATISFIED` | `NUMBER(18)` | cacah anak yang belum tuntas |

Di sisi XML, `pxAddChildWork` dipakai 6 kali dan `pxCoveredInsKeys` 8 kali. Jadi relasi Claim ke Komite yang selama ini hanya terbaca sebagai pemanggilan rule **punya wujud kolomnya**: `PXCOVERINSKEY` pada baris Komite berisi `pzInsKey` klaim induknya.

##### `PYID` tidak dijamin unik

| | |
|---|---|
| Primary key | `PZINSKEY` (`VARCHAR2(255 CHAR) NOT NULL`) |
| `PYID` | `VARCHAR2(32 CHAR)` — **tanpa unique constraint, tanpa index tersendiri** |

Nomor klaim `CLMNP-...` tidak dijamin unik oleh basis data. Pola yang sama dengan `OS_AKSEPTASI_KLAIM`. Diukur REQ-018.

##### `BYTE` versus `CHAR` — campur, bahkan di dalam 19 kolom itu sendiri

Kolom bawaan Pega seluruhnya `CHAR`. Sembilan belas kolom yang di-expose **tidak seragam**:

| Semantik | Kolom |
|---|---|
| `BYTE` (14) | `SOBNAME`, `BUSINESSNAME`, `CEDINGCONAME`, `INSUREDNAME`, `NOPOLIS`, `MASTERID`, `DATEOFLOSS`, `ISSUBJECTIVITY`, `PLA_CEDING`, `PLA_SOB`, `DLA_CEDING`, `DLA_SOB`, `POLICY_CEDING`, `CAUSEOFLOSS` |
| `CHAR` (3) | `NOCLAIM`, `FLAGONGOINGCOMMITTE`, `CLMNO` |
| bukan teks (2) | `KOMITECOUNT` `NUMBER(18)`, `DATERECEIVED` `DATE` |

Untuk data non-ASCII, `VARCHAR2(n BYTE)` menampung lebih sedikit karakter daripada yang terlihat. Paling terpapar: `INSUREDNAME` (255 BYTE), `CEDINGCONAME` (128 BYTE), dan empat kolom `*_CEDING`/`*_SOB` (500 BYTE). Masuk risiko migrasi data; sistem baru memakai satu semantik saja.

##### `PZPVSTREAM` dan hak akses

`PZPVSTREAM BLOB`, `SECUREFILE ... ENABLE STORAGE IN ROW`. Di luar 19 kolom itu dan kolom bawaan Pega, **tidak ada satu pun kolom nilai uang** — menegaskan bagian 13.5.

```
GRANT ALTER, DELETE, INDEX, INSERT, REFERENCES, SELECT, UPDATE, ... TO POOLDATA
```

`POOLDATA` memiliki `INSERT`, `UPDATE`, `DELETE`, bahkan `ALTER` atas tabel work Pega. Premis REQ-021 terbukti dari DDL; yang belum diketahui hanya apakah hak itu dipakai.

#### 8.5 Properti bersama

**Tambahan — kopling yang belum tercatat: `.IsEditClaim`.**
`EditXOLAlokasi` menulis `1`, `CountLossAllocation_act` menulis `0`, dan **tidak ada pembaca di sisi Claim**. Bila pembacanya ada di modul Komite, maka ini kopling lintas modul berarah Claim → Komite yang belum masuk kontrak ini. Statusnya hipotesis, bukan fakta — lihat `_selesai/OPEN-QUESTIONS.md` C8 dan hipotesis tandingannya.


`KomiteTreatyNonProp` (4 rule) · `pxAddChildWork` (2) · `KomiteList`/`KomiteCount`/`KomiteLoop` (2) · `KomiteNo` (5) · `ComiteeClaim` (4) · `AcceptanceStatus` (**11 rule** — properti paling terkopel) · `IsKomite` (2) · `KomiteAproval` (4) · `KomiteComment` (3) · `KomitePost` (4) · `TotalKomite` (2) · `DataCommitteeTreaty` (5) · `EMAILKOMITE` (3).

#### 8.6 Penentu jenjang Komite

`ReportDefinition\FilterEmailKomiteWithLimit.xml` atas class `ASM-FW-GCNMFW-Int-EMAILKOMITE`.
Kolom: `ID, NAME, EMAIL, LIMIT_BOTTOM, LIMIT_TOP, DEGREE, STS_AKTIF, OPERATOR_ID, STS_REJECT, STS_ADJ, TYPE_BUSINESS, TYPE_KOMITE, STS_REG, STS_SURVEY, STS_SALVAGE, STS_ADJUSTER, STS_KLAIM, JABATAN`.
Filter `A AND C AND B`: `LIMIT_BOTTOM = Param.LIMIT_BOTTOM`, `STS_KLAIM = Param.STS_KLAIM`, `STS_AKTIF = "1"`.
Jumlah baris hasil → `ChildWorkPage.KomiteLoop` = jumlah jenjang persetujuan.

---

### 9. Java Step — EVIDENCED

17 activity, 26 blok. Dua pola dominan:

| Pola | Rule |
|---|---|
| PDF byte → Base64 via `tools.getParameterPage()` | `HTMLToPDF` (5), `HTMLRISlipToPDF` (3), `AttachRISlipToWork` (2), `GenerateCFS_act` (1), `GetBase64Attachment` (1) |
| Dedup `HashSet<String>` atas PageList mata uang/estimasi | `CountTotalInsterest_Act` (2), `GetHistoryMasterID_NP` (1), `CountSpreading_act` (1), `AdjClaimCNP_Act` (1), `AdjClaimAmount_Act` (1), `AddAkseptasiCNP_Act` (1) |
| Serialisasi page → JSON | `SetValueClaimTNP_Act`, `GetDetailPolis_act`, `GeneratePlaCNP_Act`, `HitServiceToKasir_Act` (2) |
| Deteksi ekstensi berkas → MIME | `InsertGoogleStorage_Act`, `InsertDocument_Act` |
| Penutupan case paksa | `ASMForceCaseClose` (2) — menyalin page via `tools.findPageByHandle` |

---

### 10. Penyimpanan Properti: EXPOSED vs IN-BLOB — EVIDENCED (analisis)

Pega menyimpan properti work object di kolom BLOB `pzPVStream` kecuali properti itu di-*expose* sebagai kolom nyata. Bukti terkuat EXPOSED adalah pemakaian properti sebagai kolom/filter/sort di Report Definition.

**Pemeriksaan 22 Report Definition di folder ini: tidak ada satu pun yang berjalan di atas class work `ASM-FW-GCNMFW-Work-ClaimTreatyNonProp` atau `ASM-FW-GCNMFW-Work`.** Seluruhnya berjalan di atas class `-Int-` (tabel/view Oracle eksternal), `Link-Attachment`, atau `Data-Admin-Operator-ID`.

**Konsekuensi**: tidak ada bukti XML bahwa properti klaim mana pun di-expose sebagai kolom. Dugaan default untuk seluruh `.ClaimData.*` dan `.TreatyInMaster.*` adalah **IN-BLOB**.
→ `DESC` tabel work Pega **tidak akan** menampilkan properti-properti ini. Sumber tipe yang sahih adalah definisi property di schema PegaRULES (REQ-001), bukan `ALL_TAB_COLUMNS`.
→ Nilai klaim yang terlihat di tabel Oracle bisnis (`OS_AKSEPTASI_KLAIM`, `json_klaim`, `claimxol2`) adalah **salinan hasil**, bukan sumber. Presisi riil yang diterima akuntansi justru ada di sana — itu yang harus diprofilkan.

---

### 11. Rekonsiliasi terhadap DDL sistem lama — EVIDENCED

**Sumber**: `pengetahuan/DDL_Script_ClaimNonProp.xls`. **Versi 2026-09-18 10:36 — 48 objek**, seluruhnya owner `POOLDATA`. (Versi sebelumnya 46 objek; dua VIEW ditambahkan, lihat bagian 13.)
**Stempel asal berkas**: dibuat 2026-09-18 03:06:58 oleh `Raynold`, WPS Spreadsheets, format BIFF (OLE Compound Document); diperbarui 2026-09-18 10:36.
**Otoritas**: tertinggi untuk struktur data. Bentuk skripnya `CREATE OR REPLACE EDITIONABLE ...` beserta klausa `SEGMENT CREATION`, `STORAGE(INITIAL ...)`, dan `TABLESPACE` — itu keluaran `DBMS_METADATA.GET_DDL` dari basis data yang berjalan, bukan rancangan yang ditulis orang.

#### 11.1 Cakupan

| Jenis | Jumlah |
|---|---|
| TABLE | 31 |
| VIEW | 10 |
| PROCEDURE | 6 |
| FUNCTION | 1 |

**Tidak memuat satu pun objek Pega** (`pc_*`, `pr4_*`, `pr_data_admin`). Karena itu REQ-001 dan REQ-012 **tetap BLOCKER**, dan dugaan `.ClaimData.*` tersimpan IN-BLOB di `pzPVStream` **tetap terbuka**.

#### 11.2 Hasil rekonsiliasi

| | Jumlah |
|---|---|
| Naik dari GUESS ke CONFIRMED | **10** |
| Sudah CONFIRMED, kini ber-owner pasti | 36 |
| Ada di daftar saya, **tidak ada di DDL** | 9 |
| Objek baru yang belum pernah saya daftarkan | **11** |
| Koreksi jenis objek (saya salah) | **5** |

**Naik ke CONFIRMED**: `EMAILKOMITE`, `V_POLIS`, `V_M_CAUSE_OF_LOSS`, `V_D_CAUSE_OF_LOSS`, `ADJUSTERCONSULTANT`, `LST_BANK_GROUP`, `CATASTROPHE`, `PROVINCE`, `TREATYGROUP`, `CURRENCYSTANDARD`.

**Koreksi jenis — saya salah lima kali**:

| Objek | Saya daftarkan | Sebenarnya |
|---|---|---|
| `CURRENCY` | TABLE | **VIEW** |
| `CITY` | TABLE | **VIEW** |
| `PROVINCE` | TABLE | **VIEW** |
| `CURRENCYSTANDARD` | TABLE | **VIEW** |
| `GET_TOKEN_STORAGE` | FUNCTION | **PROCEDURE** |

**Prediksi saya yang meleset paling jauh**: saya menulis *"Dugaan saya: objek `V_POLIS` ini TIDAK ADA"*, dengan alasan nama class Pega hanya wadah. **Salah.** `V_POLIS` ada, berupa VIEW dengan definisi sepanjang 4.472 karakter — view terbesar di seluruh DDL. Pelajarannya: ketiadaan sebuah nama di dalam `FROM` tidak berarti objeknya tidak ada; Pega dapat mengakses tabel lewat pemetaan class, tanpa nama objek pernah muncul di SQL mana pun.

**Tidak ada di DDL** (tetap di antrean Oracle, **bukan** disimpulkan mati — DDL ini bisa saja parsial): `TREATY_OUT`, `CLAIMXOL`, `TRLOSS_DETAIL_T` (owner `REINSURANCE`), `V_MST_USER_TEKNIS`, `F_GET_EMAIL` (owner `GL`), `PLATNP_SEQ`, dan seluruh objek Pega.

`TREATY_OUT` yang absen patut diperhatikan: ia BLOCKER, terbaca literal di SQL produksi, tetapi DDL-nya tidak ikut. `CLAIMXOL` absen padahal `CLAIMXOL2` ada.

#### 11.3 Objek baru — dirujuk di dalam DDL, tidak pernah terbaca di XML

Sebelas objek yang tidak meninggalkan jejak apa pun di 279 berkas XML, karena hanya disentuh dari dalam stored procedure:

| Objek | Ditemukan di | Kenapa penting |
|---|---|---|
| **`M_CURRENCYSTANDARD`** | `GETCURRENCYSTANDARD`, view `CURRENCYSTANDARD` | **tabel kurs yang sebenarnya**. Seluruh konversi mata uang bertumpu padanya |
| **`TANGGAL_CLOSING`** | `PROC_GENERATE_SEQUENCE_NUMBER` | **tanggal tutup buku**. Menentukan periode produksi sebuah nomor — berkaitan dengan aturan cut-off tanggal 25 di `MEMORI_PEMAHAMAN.MD` |
| `CITYINPUT`, `DISTRICTINPUT`, `RWINPUT` | `GET_TOKEN_STORAGE` | perlu diverifikasi — bisa jadi alias di dalam view, bukan objek tersendiri |
| `M_CAUSE_OF_LOSS`, `D_CAUSE_OF_LOSS` | view `V_*_CAUSE_OF_LOSS` | tabel dasar penyebab kerugian |
| `M_CURRENCY`, `M_PROVINCE` | view `CURRENCY`, `PROVINCE` | tabel dasar |
| `M_NATION` | `GET_TOKEN_STORAGE` | — |
| `BRANCH` | dalam DDL | peran belum terbaca |

Masuk `pengetahuan/PULL-LIST.csv` sebagai **REQ-017**.

**Koreksi atas sapuan pertama saya.** Daftar pertama memuat `V_DAY_CLOSING` dan `V_PERIODE_DATE` sebagai objek baru. **Keduanya bukan objek** — mereka variabel lokal PL/SQL (`v_day_closing NUMBER`, `v_periode_date DATE`) yang terbaca sebagai nama objek karena sapuan saya menormalkan huruf besar. Penggantinya adalah objek yang sebenarnya dibaca di tempat itu: **`POOLDATA.TANGGAL_CLOSING`**.

**Nilai mati di dalam basis data.** `PROC_GENERATE_SEQUENCE_NUMBER` memuat:
```sql
IF TRUNC(v_now) <= TO_DATE('02/01/2026','DD/MM/YYYY') THEN
    v_mm_yyyy := '12.2025';  v_tahun := 2025;
```
Tanggal dan periode tertanam langsung di prosedur. Ini kelas tambalan yang sama dengan yang terdaftar di bagian 7, **tetapi berada di basis data, bukan di Pega** — jadi sapuan XML tidak akan pernah menemukannya. Sapuan berkala menurut ADR-0012 perlu mencakup DDL, bukan hanya XML.

#### 11.4 Yang TIDAK tertutup oleh DDL

1. **Presisi riil tetap terbuka.** Dari 120 kolom `NUMBER`, **88 tidak menyatakan presisi maupun scale sama sekali** — `NUMBER` polos, yang di Oracle berarti sampai 38 digit signifikan dengan scale apa pun yang datang. Basis data **tidak memaksakan pembulatan apa pun**. Apa pun hasil `@divide(...,20)` dari Pega tersimpan apa adanya. **REQ-005 tetap hidup; ADR-0003 tetap draft.**

   Pengecualiannya satu-satunya: **`TREATYINPRODUCTION` memakai `NUMBER(20,4)` pada 23 kolom nilai** (`PREMI_OGP`, `CLAIM`, `NET_PREMIUM`, `PPHVALUE`, `PPNVALUE`, `BALANCE_DUE_TO`, dst.). Jadi ada **satu** tabel yang menetapkan 4 desimal, dan seluruh sisanya bebas. Itu menguatkan temuan §6.3 bahwa skala mengikuti modul, bukan jenis nilai.

2. **Properti IN-BLOB justru terbukti dengan tidak munculnya.** DDL ini tidak memuat tabel work Pega sama sekali, jadi ketiadaan `.ClaimData.*` di sini **bukan** bukti propertinya tidak ada — ia konsisten dengan dugaan IN-BLOB, bukan membantahnya.

3. **Pemetaan properti Pega ke kolom tetap belum terjawab.** DDL memberi sisi kolom; XML memberi sisi properti; jembatannya tetap `Data-Admin-DB-Table` — **REQ-012, masih BLOCKER**.

#### 11.5 T1, T2, T3 — terjawab sebagian tanpa menjalankan PREFLIGHT

| | Jawaban | Dari mana |
|---|---|---|
| **T2** tipe `DATA_JSON` | **`CLOB` dengan `CHECK (DATA_JSON IS JSON)`**, disimpan `SECUREFILE ... ENABLE STORAGE IN ROW`. Berlaku untuk `JSON_KLAIM` dan `OS_AKSEPTASI_KLAIM` | DDL langsung |
| **T1** versi Oracle | **12c ke atas** — `IS JSON` dan notasi titik `a.data_json.TypeLoss` keduanya butuh 12c+; `EDITIONABLE` butuh 12.1+ | disimpulkan dari sintaks yang dipakai |
| **T3** owner | **`POOLDATA`** untuk ke-46 objek ini | DDL langsung |

`pengetahuan/PREFLIGHT.sql` tetap perlu dijalankan untuk **T4** (hak akses), **T5** (letak PegaRULES), **T7** (ukuran tabel), dan untuk memastikan 9 objek yang tidak ada di DDL.

---

### 12. Kenapa sistem ini butuh tambalan per-case — satu rantai sebab akibat

Tiga temuan yang selama ini terpisah ternyata satu cerita. Ditulis di sini supaya tidak dibaca sebagai tiga keluhan yang berdiri sendiri.

**Sebab.** Tidak ada satu pun activity yang memanggil `CountClaimTNP_Act`; ia dipicu dari Section lewat `<pyActivity>`, 42 rujukan di tiga berkas. Urutan jalannya `CountClaimTNP_Act`, `CountLossAllocation_act`, `GenerateCFS_act`, dan `SaveToOS` **tidak ditetapkan di kode mana pun** — ia ditentukan oleh kontrol mana yang ditekan pengguna (§bagian 3 `FINDING-003`).

**Akibat pertama — hasil tidak dapat direproduksi.** Dua pengguna dengan masukan sama tetapi urutan tindakan berbeda dapat menghasilkan keadaan berbeda. `CountLossAllocation_act` meng-`<APPEND>` baris `"UR"` dan tidak ada penghapusan daftar sebelum itu; `PEGA_JSON_OS_AKSEP_KLAIMTNP` selalu `INSERT` dan tidak pernah `UPDATE` (logika upsert-nya dikomentari). Keduanya menumpuk, bukan mengulang.

**Akibat kedua — selalu ada kasus yang keluar jalur.** Karena keadaan akhir bergantung jalan yang ditempuh, akan selalu ada klaim yang berakhir di keadaan yang tidak diantisipasi rule mana pun.

**Akibat ketiga — satu-satunya perbaikan yang murah adalah menambal klaim itu satu per satu.** Memperbaiki penyebabnya berarti mengubah arsitektur pemicuan; menambal satu klaim berarti menambah satu `pyStepsPreCondParamsWhen`. Yang kedua selesai dalam satu jam.

**Bukti bahwa rantai ini nyata, bukan teori**: 29 langkah tambalan, 18 ekspresi, 8 rule, 6 penulis, sepanjang **delapan tahun** (2018-01-16 s/d 2026-07-16), dengan laju yang naik sampai 2023 lalu turun tetapi **tidak berhenti** (§7.5).

**Karena itu ketiga keputusan berikut harus diambil bersama, bukan sendiri-sendiri**: ADR-0011 (hitung ulang idempoten), ADR-0012 (deteksi tambalan lewat sapuan), dan keputusan Q22 (perhitungan sebagai turunan data, bukan akibat tombol). Mengambil satu tanpa dua lainnya akan mengembalikan pola yang sama dalam bentuk baru.

---

### 13. Pembaruan DDL 2026-09-18 10:36 — dan satu berkas yang belum sampai

#### 13.1 Yang bertambah: dua VIEW

| Objek | Menutup |
|---|---|
| `POOLDATA.CLAIMXOL` | sebelumnya terdaftar "tidak ada di DDL" — **ternyata VIEW, bukan TABLE**. Koreksi keenam atas jenis objek. |
| `POOLDATA.V_MST_USER_TEKNIS` | sebelumnya GUESS, sekarang **CONFIRMED** dan ber-owner pasti |

#### 13.2 `CLAIMXOL` membuka struktur JSON — REQ-009 terjawab sebagian

```sql
FROM pooldata.os_akseptasi_klaim OS, json_table (
     data_json, '$.CNPLayerList[*]'
     columns( XOL varchar2 path '$.XOL',
       nested path '$.CNPCurrencyList[*]' columns(
         TotalXOLGross    varchar2 path '$.TotalXOLGross',
         CNPReinstatement varchar2 path '$.CNPReinstatement',
         Currency         varchar2 path '$.Currency',
         KursIDR          varchar2 path '$.KursIDR')))
WHERE SUBSTR(CASEID,1,25) = 'ASM-FW-GCNMFW-WORK CLMNP-'
```

Tiga hal terbaca langsung:

1. **Bentuk JSON**: `$.CNPLayerList[*]` bersarang `$.CNPCurrencyList[*]`. Jadi layer adalah tingkat luar dan mata uang tingkat dalam — bukan sebaliknya. Itu memperkuat ADR-0024 (satu akseptasi per klaim per layer per mata uang) dan menetapkan arah sarangnya — buktinya dipindahkan ke badan ADR itu.
2. **Seluruh nilai uang di dalam JSON diekstrak sebagai `varchar2`** — `TotalXOLGross`, `CNPReinstatement`, `KursIDR`. Angka uang di dalam blob **berupa teks**. Presisinya adalah apa pun yang tertulis di teks itu, tanpa batas dan tanpa jaminan.
3. **Bentuk `CASEID`**: `'<NAMA CLASS> <pyID>'`, yaitu `ASM-FW-GCNMFW-WORK CLMNP-…`. Class yang dipakai adalah **class induk** `ASM-FW-GCNMFW-Work`, bukan `ClaimTreatyNonProp`.

#### 13.3 Objek baru lagi, dari dalam dua view ini

| Objek | Ditemukan di | Catatan |
|---|---|---|
| `POOLDATA.OS_AKSEPTASI_SUBJECTIVITY` | `CLAIMXOL`, cabang `UNION ALL` dengan `STS_SUBJECTIVITY = '1'` | tabel akseptasi kedua, khusus Subjectivity. Belum pernah terlihat di XML |
| `POOLDATA.MST_USER_TEKNIS` | `V_MST_USER_TEKNIS` | tabel dasar; kolomnya sebagian di dalam `JSON_DATA` |
| **`hrdasm.v_hrd_mst@asmd.sinarmas.co.id`** | `V_MST_USER_TEKNIS` | **database link ke sistem HRD**. Menjawab sebagian T5b: db link memang dipakai di lingkungan ini |

#### 13.4 Berkas yang belum sampai: DDL tabel work Pega

`DATAPEGA.PC_ASM_FW_GCNMFW_WORK` **tidak ada di dalam berkas `.xls` versi 2026-09-18 10:36.** Saya baca ulang seluruh 48 baris × 4 kolom dan menyapu seluruh isi sel untuk `DATAPEGA`, `PC_ASM`, dan `PZPVSTREAM`; satu-satunya kecocokan adalah nama parameter `Datapega IN CLOB` pada dua prosedur. Tidak ada berkas baru lain di folder artefak.

Sampai berkas itu sampai, yang berikut **belum dapat dikerjakan** dan tidak boleh dianggap selesai: pemeriksaan apakah Claim dan Komite berbagi satu tabel; `PXCOVERINSKEY` sebagai penghubung induk-anak; ketiadaan unique constraint pada `PYID`; pola `BYTE` versus `CHAR`; `ISSUBJECTIVITY VARCHAR2(5)`; dan pengisian baris tabel work ke `pengetahuan/SCHEMA-ACTUAL.csv`.

Yang **sudah** dapat dikerjakan dari keterangan yang disampaikan langsung, dan sudah dikerjakan: konsekuensi IN-BLOB (bagian 13.5), format `DATEOFLOSS` (`FINDING-005`), dan `FLAGONGOINGCOMMITTE`/`KOMITECOUNT` (bagian 13.6).

#### 13.5 `.ClaimData.*` IN-BLOB — status naik dari dugaan menjadi **terbukti**

Dari sembilan belas kolom yang di-expose di luar kolom bawaan Pega, **tidak satu pun kolom nilai uang**. Tidak ada `ClaimAmount`, `Currency`, `Kurs`, `ValueAdjustment`, dan tidak ada satu pun anggota `SpreadingRisk` maupun `AdjustmentList`.

Tiga konsekuensi yang harus dinyatakan terang-terangan:

1. **REQ-005 tidak dapat dijalankan atas properti klaim.** Profil presisi lewat SQL mensyaratkan kolom; nilai-nilai itu tidak punya kolom. `ALL_TAB_COLUMNS` tidak akan pernah memberi tipe atau presisi `ClaimAmount`, berapa kali pun dijalankan. **REQ-005 dipersempit menjadi hanya tabel bisnis `POOLDATA`.**
2. **REQ-001 menjadi satu-satunya jalur** untuk tipe properti klaim. Bila `PR4_*` tidak dapat diakses, satu-satunya cadangan adalah membongkar `PZPVSTREAM`. Rencananya disiapkan sekarang, tidak menunggu REQ-001 gagal — lihat `TABLE-EXTRACTION-REQUEST.md`.
3. **ADR-0003 tetap draft, tetapi alasannya berubah.** Bukan karena datanya belum ditarik, melainkan karena **basis data memang tidak menyimpan presisi untuk nilai-nilai ini sama sekali**. Diperkuat oleh 13.2: bahkan ketika nilai itu dikeluarkan dari JSON lewat view, tipenya `varchar2`.

#### 13.6 `FLAGONGOINGCOMMITTE` dan `KOMITECOUNT` — sebagian yang saya sebut "tidak dapat diukur" ternyata dapat

Sapuan 279 berkas:

| Kolom | Kemunculan di folder Claim |
|---|---|
| `FLAGONGOINGCOMMITTE` | **nol** |
| `KOMITECOUNT` | 2, keduanya `ChildWorkPage.KomiteCount` — ditulis pada case anak, bukan pada klaim |

Artinya keduanya **ditulis dari modul Komite**, tetapi karena keduanya punya kolom nyata, keduanya **dapat dihitung lewat SQL tanpa membuka modul Komite**. Itu membuka dua hal yang sebelumnya saya tandai tidak terukur:

- **Q3 / `CloseClaimMD`** — `FLAGONGOINGCOMMITTE` adalah penanda yang dicari: berapa klaim ditutup ketika penanda itu masih menyala.
- **Q9 / berapa kali Adjustment dikirim ke Komite** — `KOMITECOUNT` menjawabnya langsung.

Masuk antrean sebagai **REQ-019**. Maknanya tetap `DEFERRED-TO-KOMITE-SESSION`; yang berubah adalah jumlahnya kini terukur.

---

### 14. Catatan metodologis — setiap lapisan sumber punya titik buta sendiri

Ditulis karena polanya sudah berulang tiga kali, bukan karena satu kejadian.

**Analisis XML saja punya titik buta yang sistematis.** Tiga hal berikut menggerakkan sistem ini dan **tidak meninggalkan jejak apa pun di 279 berkas XML**:

| Yang tidak terlihat dari XML | Terlihat setelah lapisan apa masuk |
|---|---|
| Sebelas objek basis data yang hanya disentuh dari dalam stored procedure — termasuk `M_CURRENCYSTANDARD`, tabel kurs yang menopang seluruh konversi | DDL |
| Tanggal dan periode tertanam di `PROC_GENERATE_SEQUENCE_NUMBER` (`TO_DATE('02/01/2026')` → `'12.2025'`) — kelas tambalan yang sama dengan 29 langkah di bagian 7, tetapi di basis data | DDL |
| Ketergantungan lintas basis data: `hrdasm.v_hrd_mst@asmd.sinarmas.co.id` | DDL |

Dan sebaliknya, **DDL punya titik butanya sendiri**: seluruh nilai uang klaim tidak punya kolom, sehingga `ALL_TAB_COLUMNS` tidak akan pernah menyebutkannya — keberadaannya hanya terbaca dari XML.

Kesimpulan yang mengikat cara kerja, bukan sekadar catatan:

> **Tidak ada satu sumber pun yang memperlihatkan sistem ini secara utuh.** Setiap lapisan baru yang masuk membuka jalur yang lapisan sebelumnya tidak dapat melihatnya. Karena itu tidak boleh ada pernyataan berbentuk "tidak ada X di sistem ini" — yang sah hanya "tidak ada X di lapisan yang sudah saya baca".

**Konsekuensi untuk ADR-0012**: sapuan berkala harus mencakup **tiga** hal, bukan satu — ekspor XML, DDL, dan **daftar database link**. Yang ketiga baru ditambahkan setelah `V_MST_USER_TEKNIS` terbaca; sebelum itu tidak ada yang tahu perlu menyapunya.

Lapisan yang **belum** masuk sama sekali dan titik butanya masih penuh: tabel work Pega (`DATAPEGA`), isi `PR4_*`, pemetaan `Data-Admin-DB-Table`, dan isi blob `PZPVSTREAM`.

---

### 15. Konvensi boolean — diperiksa, dan hasilnya bersih

Dugaan yang diuji: adakah properti yang **ditulis** dengan satu konvensi boolean lalu **dibaca** dengan konvensi lain? Itu akan menjadi kasus kedua dari pola `CheckDateDOL_Act` — dua bagian sistem berselisih tentang bentuk data.

**Hasil: nol bentrok.** Setiap properti taat asas pada dirinya sendiri.

| Konvensi | Properti |
|---|---|
| Angka (`==1` / `==0`) | `IsOutstanding` (73 pembacaan), `IsCFS`, `IsKomite`, `IsAcceptation`, `IsPLA`, `IsPayment`, `IsSaveToOs`, `Flagkomite`, `FlagErr`, `FlagCurrency`, `FlagError` |
| Teks (`=="true"`) | `IsTPL`, `FlagActualPremium`, `IsSubjectivity` |

Perpecahannya **antar** properti, bukan **di dalam** satu properti. Jadi ia soal keseragaman rancangan, bukan bug.

**Batas pemeriksaan**: pendeteksian penulisan hanya menangkap pasangan `<PropertiesName>` dan `<PropertiesValue>` yang berdampingan di Activity. Penulisan lewat Data Transform (`pyPropertiesName`) dan lewat layar tidak ikut tersapu. Beberapa properti — `IsTreatyIn`, `IsReject`, `IsCloseFile`, `IsAnyAcceptation` — terbaca tetapi penulisnya tidak tertangkap, jadi untuk keempatnya pernyataan "bersih" **belum berlaku**.

Untuk sistem baru: satu konvensi saja, tipe boolean sejati. `ISSUBJECTIVITY VARCHAR2(5)` di basis data lama menjadi contoh yang tidak diwarisi.

---

### 16. `V_POLIS` — bukan gabungan banyak tabel, melainkan proyeksi satu dokumen JSON

Objek terbesar di DDL (4.472 karakter) ternyata sederhana bentuknya: **26 kolom, satu tabel sumber (`json_polis`), 32 pemanggilan `JSON_VALUE`, nol `JOIN`, nol `UNION`.**

Kolomnya: `RN`, `CASEID`, `POLICYNO`, `STARTDATETIME`, `ENDDATETIME`, `SOURCEOFBUSINESS(NAME)`, `CEDINGCO(NAME)`, `BUSINESSCODE`, `BUSINESSNAME`, `CUSTOMERNAME`, `SOURCEOFBUSINESSROOT`, `PRODUCTIONDATE`, `PREMI`, `DISCOUNT`, `TSI`, `TYPE`, `PRODKE`, `BUSINESSTYPE`, `MARKETINGCODE`, `CLIENTIDCON`, `CLIENTIDORG`, `QQ`, `SELECTRENEWAL`, `RNWSTATUS`.

Tiga hal yang mengubah gambaran:

1. **"Polis" di sistem ini tidak punya tabel sendiri.** Ia dokumen JSON di `json_polis`, dan `V_POLIS` hanyalah cara memandangnya sebagai baris dan kolom. `CASEID` pun berasal dari dalam JSON (`a.DATA_JSON.IDNewBisnis`).
2. **Satu tabel menampung dua bentuk dokumen** — Facultative dan Treaty — dibedakan `QuotationData.BusinessFac`, dengan jalur JSON yang berbeda untuk data yang sama. Lihat `FINDING-005` bagian 6.1.
3. **Seluruh kolom tanggalnya adalah teks hasil sambungan `SUBSTR`**, bukan `DATE`. Inilah asal format `dd/MM/yyyy` yang bertabrakan dengan `.DateOfLoss` di `FINDING-005`.

Yang ini menutup: **14 RDB List pada class `ASM-FW-GCNMFW-Int-V_POLIS` membaca proyeksi JSON, bukan tabel relasional.** Dugaan awal saya bahwa `V_POLIS` tidak ada sebagai objek memang salah, tetapi alasan di baliknya — bahwa class itu tidak menunjuk tabel bisnis sungguhan — ternyata setengah benar: ia menunjuk view atas satu dokumen.

---

### 17. Tidak ada pembatas antara yang dicoba dan yang dipakai

Lima hal yang selama ini terdaftar sebagai anomali terpisah (E3, E4, E6, E7, E8, E9) ternyata satu temuan. Ditulis sebagai satu bagian karena memperlakukannya sebagai lima keluhan membuat sebabnya tidak terlihat.

#### 17.1 Artefak uji berjalan di alur produksi — dan dapat dicapai pengguna

| Artefak | Ukuran | Jalur masuknya |
|---|---|---|
| `Activity\SaveAdjustmentToOSAksep_Act_Tes` | — | dipanggil `Call` dari `CreateChildKomiteCNP_Act` langkah 25 |
| `Section\Hitung_Test` | 187 KB | **`<pyInclude>Hitung_Test</pyInclude>` di dalam `Section\InputAcceptation.xml`** |
| `Harness\Hitung_Test` | 240 KB | harness tersendiri, `pyHarnessName = Hitung_Test` |

Pertanyaan yang diajukan: terdaftar saja, atau benar-benar dapat dicapai dari layar? **Dapat dicapai.** `Hitung_Test` di-`include` ke dalam Section produksi `InputAcceptation` — bukan sekadar terdaftar di suatu tempat. Dan `SaveAdjustmentToOSAksep_Act_Tes` dipanggil lewat `Call` dari activity produksi, bukan dari harness uji.

Keduanya membawa stempel kompilasi 2024 (`..._Stream_20240229T095743_632_GMT`), jadi bukan sisa lama yang tertinggal — mereka masih dibangun ulang bersama rule lain.

#### 17.2 Sambungannya dengan tambalan

| | Artefak uji | Tambalan per-case |
|---|---|---|
| Cara masuk | di-`include` / di-`Call` dari alur produksi | satu `pyStepsPreCondParamsWhen` |
| Umur | stempel 2024, masih terkompilasi | tertua 2018, terbaru 2026-07-16 |
| Yang menghentikan | **tidak ada** | **tidak ada** |

**Satu sebab, dua gejala: tidak ada mekanisme yang memisahkan yang dicoba dari yang dipakai.** Sekali sesuatu masuk ke ruleset produksi, tidak ada apa pun yang mengeluarkannya kembali — baik ia percobaan maupun tambalan. Delapan tahun, enam penulis, dan tidak ada satu pun peristiwa pembersihan yang meninggalkan jejak.

Bukti ketiga dari pola yang sama, kali ini di basis data: `PROC_GENERATE_SEQUENCE_NUMBER` memuat `IF TRUNC(v_now) <= TO_DATE('02/01/2026') THEN v_mm_yyyy := '12.2025'` — nilai mati yang sudah lewat tanggalnya dan tetap ada.

#### 17.3 Akibatnya untuk ADR-0012

Pembekuan tambalan yang hanya mengurus `pyStepsPreCondParamsWhen` **akan bocor lewat pintu yang sama**. Artefak uji masuk lewat `pyInclude` dan `Call` — dua jalur yang tidak disentuh sapuan tambalan.

Sapuan berkala perlu menambah satu hal lagi: **artefak yang namanya mengandung penanda percobaan (`Test`, `Tes`, `Tmp`, `Coba`, `Backup`, `Old`) tetapi dirujuk dari alur produksi.** Itu murah diperiksa dan langsung memperlihatkan kebocoran.

---

### 18. Perilaku sistem berbeda menurut node yang mengeksekusi

Dimensi yang tidak pernah masuk analisis selama enam ronde, ditemukan lewat sapuan `pxSystemNodeID`.

| Yang dicari | Hasil |
|---|---|
| Rule yang bercabang pada `pxSystemNodeID` | **satu**: `When\IsPEGASyariah.xml` → `pxProcess.pxSystemNodeID = "jboss1074"` |
| Rule yang memakai When itu | **satu**: `Activity\HitServiceToKasir_Act.xml`, 3 langkah |
| Yang berubah | `TempKasir.CARI15 = "100115"` — kode yang dikirim ke sistem kasir |
| `jboss117` (41x), `jboss122117` (36x) | **bukan percabangan** — seluruhnya `<pxHostId>`, metadata audit |

Cakupannya sempit: satu rule, satu field, jalur setoran ke kasir. Tetapi **kode setoran yang berbeda menunjuk ke entitas pembukuan yang berbeda**, bukan sekadar server yang berbeda — dan itu membalik kesimpulan pertama saya tentang "syariah". Uraiannya di `_selesai/OPEN-QUESTIONS.md` G1, ditulis sebagai dua kemungkinan sejajar dan **tidak** saya putuskan sendiri.

**Batas yang penting untuk cakupan seluruh analisis ini**: karena percabangannya dievaluasi saat berjalan di dalam ruleset yang sama, ekspor 279 berkas berlaku untuk kedua node. Yang tidak dapat dibuktikan dari sini: apakah ada instance Pega lain dengan ruleset yang berbeda sama sekali. Itu `_selesai/OPEN-QUESTIONS.md` A15.

---

### 19. Rekonsiliasi kolom DDL terhadap properti Pega

Dua daftar terpisah, seperti diminta. Hasil lengkapnya di `pengetahuan/rekonsiliasi-kolom-vs-properti.tsv`.

| | Jumlah |
|---|---|
| Kolom unik di seluruh DDL (31 tabel + 8 view) | 408 |
| Properti Pega unik di 279 XML (di luar `px*`/`py*`/`pz*`) | 659 |
| **Berpasangan** (nama sama) | **110** |
| **Kolom tanpa pasangan properti** | **298** |
| **Properti tanpa pasangan kolom** | **549** |

#### 19.1 Delapan puluh tiga persen properti tidak punya kolom di mana pun

549 dari 659 properti tidak ditemukan sebagai nama kolom di seluruh DDL. Itu bukan anomali — itu **ukuran kuantitatif dari temuan IN-BLOB** di bagian 13.5, dan mencakup hampir seluruh model klaim: `ADJUSTMENTLIST`, `ACCEPTANCESTATUS`, `ADJCLAIMVALUE`, `ALOKASIIDR`, `AMOUNTIDR`, dan seterusnya.

**Peringatan tafsir**: pencocokan ini berdasarkan **nama**, bukan pemetaan sungguhan. Pemetaan yang sah hanya datang dari `Data-Admin-DB-Table` (REQ-012). Nama yang sama belum tentu kolom yang sama, dan nama yang berbeda belum tentu bukan. Angka di atas adalah **batas atas** jumlah properti IN-BLOB, bukan hitungan pasti.

#### 19.2 298 kolom yang tidak disentuh Pega sama sekali

Kolom seperti `BALANCE_BEFORE_PPH`, `BALANCE_DUE_TO`, `BROKERAGE`, `CESSIONVALUE`, `COINS_MIN`/`COINS_MAX`, `CLAIMPAYMENTTYPE` ada di tabel bisnis tetapi tidak pernah muncul sebagai nama properti di folder Claim.

Tiga kemungkinan, dan XML tidak dapat memilih di antaranya:

1. Diisi oleh **modul lain** (Produksi, Akuntansi, Komite).
2. Diisi oleh **proses basis data** — dan itu bertemu dengan REQ-023 (~60 kolom skalar `OS_AKSEPTASI_KLAIM` yang tidak diisi prosedurnya) serta REQ-021 (jalur tulis dari luar aplikasi).
3. Mati.

Yang perlu ditegaskan: **kolom-kolom ini tetap ikut termigrasi bila tabelnya dimigrasi.** Memutuskan mana yang dibuang memerlukan jawaban REQ-021 dan REQ-023 lebih dulu — tidak boleh disimpulkan dari ketiadaannya di folder Claim saja, karena lingkup sesi ini memang satu modul.


---

# III. Temuan — apa yang ditemukan di dalamnya

Tujuh laporan kondisi sistem lama. Netral: tidak memuat penilaian hukum, tidak menuduh siapa pun. Empat terbukti dari XML, satu dari DDL, satu bersyarat, satu belum terbukti.

## FINDING-001 — Ambang kewenangan Komite dibandingkan tanpa konversi mata uang

Sumber: `FINDING-001-threshold-currency.md`

**Jenis**: laporan kondisi sistem lama (bukan keputusan migrasi)
**Status**: BELUM TERBUKTI — menunggu REQ-011
**Ruang lingkup bukti**: `D:\XML_NURE\Claim Non Prop` saja
**Tanggal**: 18 September 2026

Dokumen ini memaparkan rantai perhitungan apa adanya. Ia **tidak** menyatakan adanya pelanggaran, kelalaian, atau kerugian, dan **tidak** menunjuk pihak mana pun. Apakah kondisi ini benar-benar berdampak hanya dapat dijawab oleh data produksi (REQ-011), bukan oleh XML.

---

### 1. Rantai bukti

Seluruhnya berada dalam satu berkas: `Activity\CreateChildKomiteCNP_Act.xml`.

| Langkah | Aksi | Nilai |
|---|---|---|
| 8 | `Property-Set` | `Local.TotalValueAdjust = 0` |
| 10 | `Property-Set` | `Local.CekLimitPersen = pyWorkPage.TreatyInMaster.RNMShare` |
| 10 | `Property-Set` | `Local.LimitMax = 30000000.00` |
| 10 | `Property-Set` | `Local.LimitMaxDivHead = 50000000.00` |
| 10 | `Property-Set` | `Local.LimitPersenMax = 30.00` · `Local.LimitPersenMaxDivHead = 30.00` |
| 11 | *(loop)* atas `pyWorkPage.ClaimData.AdjustmentList`, `repeat=EMBEDDED` | — |
| 11.1 | `Property-Set` | `Local.TotalValueAdjust = .ValueAdjustment` |
| 12 | `Property-Set`, kondisi `Local.CekLimitPersen<=Local.LimitPersenMax` **dan** `Local.TotalValueAdjust<=Local.LimitMax` | `Local.Flagkomite = 1` |
| 13 | `Property-Set`, kondisi `Local.CekLimitPersen<=Local.LimitPersenMaxDivHead` **dan** `Local.TotalValueAdjust>Local.LimitMax && Local.TotalValueAdjust<=Local.LimitMaxDivHead` | `Local.Flagkomite = 2` |
| 26.4 | `Property-Set`, kondisi `Local.Flagkomite==1` | `Param.LIMIT_BOTTOM = 0` |
| 26.5 | `Property-Set`, kondisi `Local.Flagkomite==2` | `Param.LIMIT_BOTTOM = 25000001` |
| 26.7 | `Call pxRetrieveReportData` | menjalankan `FilterEmailKomiteWithLimit` dengan `Param.LIMIT_BOTTOM` |
| 26.13 | `Property-Set` | `ChildWorkPage.KomiteLoop = @SizeOfPropertyList(pyReportContentPage.pxResults)` |

Hasil kueri `FilterEmailKomiteWithLimit` menentukan **siapa saja anggota Komite** dan **berapa jenjang persetujuan** yang harus dilalui.

### 2. Titik di mana konversi tidak terjadi

Perbandingan terhadap ambang berlangsung di **langkah 12 dan 13**, atas `Local.TotalValueAdjust`.

Nilai itu berasal dari langkah 11.1, yaitu `.ValueAdjustment` apa adanya. Di sepanjang langkah 8 sampai 13 **tidak ada satu pun operasi konversi mata uang**: tidak ada pemanggilan kurs, tidak ada pembacaan `Currency`, tidak ada `@divide`/`@multiply` terhadap faktor konversi. Bandingkan dengan langkah 20 dan 23 pada berkas yang sama, yang justru memakai kurs secara eksplisit (`.KursIDR`, `@divide(...RNMShare,100,20)`) ketika menyusun `CNPLayerList`.

Konstanta `30000000.00` dan `50000000.00` tidak membawa penanda mata uang.

### 3. Temuan sampingan pada rantai yang sama

Langkah 11.1 adalah **penugasan** (`=`), bukan akumulasi (`+=`), di dalam loop `AdjustmentList`.

Akibatnya `Local.TotalValueAdjust` berisi nilai `.ValueAdjustment` dari **baris terakhir** yang dikunjungi loop, bukan jumlah seluruh Adjustment — meskipun namanya "Total". Ini observasi terpisah dari soal mata uang, dan mengoreksi pernyataan saya sebelumnya yang menyebut nilai ini sebagai `Σ AdjustmentList[].ValueAdjustment`.

### 4. Batas klaim — apa yang TIDAK dapat saya simpulkan dari XML

Bagian ini menentukan apakah §2 berarti sesuatu atau tidak. **Ketiganya harus diuji lebih dulu sebelum temuan ini diperlakukan sebagai cacat.**

1. **Mata uang `.ValueAdjustment` tidak diketahui.**
   `.ValueAdjustment` muncul di **tepat satu berkas** di seluruh folder ini — `CreateChildKomiteCNP_Act.xml` — dan di situ ia hanya **dibaca**. Tidak ada satu pun rule di folder ini yang menulisinya, dan ia **bukan field UI** di section atau harness mana pun di folder ini. Asal-usul nilainya: **TIDAK DITEMUKAN DI XML**.
   Konsekuensinya: sangat mungkin nilai itu sudah dikonversi ke IDR di hulu oleh rule di luar folder ini, atau oleh layar milik modul lain. **Bila `.ValueAdjustment` selalu IDR, tidak ada cacat apa pun di sini.**

2. **Keberadaan Adjustment non-IDR di praktik tidak diketahui.**
   Struktur data jelas mendukung banyak mata uang, tetapi apakah Adjustment non-IDR benar-benar pernah dibuat adalah pertanyaan data, bukan pertanyaan kode.

3. **Perilaku bila `.ValueAdjustment` kosong tidak diketahui.**
   Bila properti itu tidak pernah terisi, `Local.TotalValueAdjust` tetap `0` dari langkah 8, sehingga `Flagkomite` selalu bernilai 1 dan `LIMIT_BOTTOM` selalu 0 — mekanisme penjenjangan efektif tidak berjalan sama sekali. Ini kemungkinan ketiga yang sama masuk akalnya, dan berbeda konsekuensinya.

### 5. Yang dibutuhkan untuk menutup temuan ini

**REQ-011** — kuantifikasi dari Oracle. Query ada di `pengetahuan/RECON.sql` BAGIAN 10.

| Pertanyaan | Menentukan |
|---|---|
| Berapa Adjustment bermata uang selain IDR? | Bila nol, temuan ini gugur seluruhnya |
| Berapa di antaranya bernilai setara di atas 30.000.000 IDR? | Besaran paparan |
| Jenjang persetujuan mana yang benar-benar menanganinya | Apakah penjenjangan berjalan sebagaimana dimaksud |
| Apakah `.ValueAdjustment` tersimpan dan dalam mata uang apa | Menjawab batas klaim §4.1 |

### 6. Hubungan dengan sistem baru

Tidak ada. Keputusan untuk sistem baru sudah ditetapkan terpisah di [ADR-0007](./docs/adr/0007-nilai-uang-berpasangan-dan-ambang-idr.md): ambang dibandingkan terhadap nilai IDR hasil konversi, dengan kurs dan tanggal kurs tercatat. Keputusan itu berlaku apa pun hasil investigasi ini.

---

### 7. Tambahan 18 September 2026 — isi `GETCURRENCYSTANDARD` terbaca

Sumber: `pengetahuan/DDL_Script_ClaimNonProp.xls`. Ini menjawab pertanyaan yang sempat saya ajukan sebagai T6 lalu dicabut karena bukan pertanyaan untuk pengguna.

```sql
CREATE OR REPLACE EDITIONABLE FUNCTION "POOLDATA"."GETCURRENCYSTANDARD" (
   i_kurs_id    IN   m_currencystandard.ID%TYPE,
   i_tgl_kurs   IN   m_currencystandard.CurrencyDate%TYPE
)  RETURN NUMBER result_cache
IS
   vhasil   m_currencystandard.CurrencyValue%TYPE;
BEGIN
   SELECT CurrencyValue INTO vhasil
     FROM (SELECT CurrencyValue FROM m_currencystandard
            WHERE ID = i_kurs_id
              AND TRUNC (CurrencyDate) <= TRUNC (sysdate)
         ORDER BY CurrencyDate DESC)
    WHERE ROWNUM < 2;
   RETURN vhasil;
EXCEPTION
   WHEN NO_DATA_FOUND THEN RETURN 1;
END;
```

Dua hal yang terbaca langsung dari badan fungsi ini, dan keduanya berdiri sendiri dari soal ambang di bagian 2.

#### 7.1 Parameter tanggal dideklarasikan tetapi tidak pernah dipakai

`i_tgl_kurs` muncul di daftar parameter dan **tidak muncul lagi di mana pun** di dalam badan fungsi. Penyaringnya memakai `TRUNC(CurrencyDate) <= TRUNC(sysdate)` — **tanggal hari ini**, bukan tanggal yang dikirim pemanggil.

Akibatnya: apa pun tanggal yang dikirim, yang dikembalikan selalu **kurs terbaru sampai hari ini**. Seluruh pemanggilan di XML memang mengirim `sysdate` (`RDBList\CurrencyStandard.xml`), jadi untuk pemakaian yang ada sekarang hasilnya sama saja. Yang tertutup adalah kemungkinan lain: **kurs historis tidak dapat diambil lewat fungsi ini**, meskipun tanda tangannya menjanjikan bisa.

Ini berhubungan langsung dengan ADR-0007, yang mewajibkan **tanggal dan sumber kurs** ikut tersimpan. Sistem lama tidak dapat memenuhi syarat itu lewat fungsi ini: nilai yang dikembalikan tidak dapat direproduksi ulang untuk tanggal tertentu di kemudian hari, karena hasilnya bergeser mengikuti `sysdate`. Rekonsiliasi angka lama tidak akan menghasilkan angka yang sama.

#### 7.2 Kurs yang tidak ditemukan menghasilkan `1`, bukan galat

`WHEN NO_DATA_FOUND THEN RETURN 1`.

Bila tidak ada baris kurs untuk suatu mata uang, fungsi mengembalikan **`1`** — yang secara aritmetika berarti **konversi satu banding satu**. Nilai 1.000 USD akan diperlakukan sebagai 1.000 IDR, tanpa galat, tanpa pesan, tanpa penanda.

Ini pilihan yang disengaja (ada penangan eksepsi khusus untuk itu), bukan kelalaian. Tetapi konsekuensinya: **kegagalan pencarian kurs tidak dapat dibedakan dari kurs yang memang bernilai 1.**

#### 7.3 Objek baru: `m_currencystandard`

Tabel kurs yang sebenarnya, dengan kolom `ID`, `CurrencyDate`, `CurrencyValue`, `UserID`, `InputDate`. **Tidak pernah muncul di 279 berkas XML** — Pega hanya melihat fungsi dan view `CURRENCYSTANDARD` di atasnya. Masuk `pengetahuan/PULL-LIST.csv` sebagai REQ-017, prioritas BLOCKER.

#### 7.4 Apa yang ini ubah dan tidak ubah pada temuan ini

**Tidak mengubah** inti FINDING-001: perbandingan `Local.TotalValueAdjust` terhadap `30000000.00` di langkah 12–13 tetap terjadi tanpa konversi apa pun, dan batas klaim di bagian 4 tetap utuh. `.ValueAdjustment` masih bisa saja sudah IDR sejak hulu.

**Menambah** dua hal yang harus diperhitungkan saat REQ-011 dijalankan:

1. Angka kurs yang dipakai sistem lama untuk suatu transaksi **tidak dapat direkonstruksi** dari tanggal transaksi. Jadi kuantifikasi "berapa yang setara di atas ambang" hanya dapat memakai kurs hari ini, dan itu harus dinyatakan sebagai batasan hasilnya.
2. Bila ada mata uang tanpa baris kurs, nilainya akan terbaca seolah sudah IDR. Query REQ-011 harus **memeriksa keberadaan baris kurs secara terpisah**, tidak boleh bersandar pada hasil fungsi ini.

## FINDING-002 — Alur bercabang berdasarkan identitas orang

Sumber: `FINDING-002-percabangan-identitas.md`

**Jenis**: laporan kondisi sistem lama (bukan keputusan migrasi)
**Status**: TERBUKTI DARI XML — tidak menunggu data Oracle
**Ruang lingkup bukti**: `D:\XML_NURE\Claim Non Prop` saja
**Tanggal**: 18 September 2026

Dokumen ini memaparkan mekanismenya apa adanya. Ia **tidak** menilai kinerja, motif, atau kepatutan siapa pun. Nama-nama di bawah muncul karena tertulis di dalam kode, bukan karena dipersoalkan.

---

### 1. Tiga bentuk percabangan

#### 1.1 Berdasarkan akun pembuat case

```
pyWorkPage.pxCreateOperator == "VINCENTVERNANDO_1"
```

| Rule | Jumlah langkah | Dibuat |
|---|---|---|
| `Activity\CreateChildKomiteCNP_Act.xml` baris 10720, 10967, 11195, 12582 | 4 | 2018-01-16 oleh `YOSUAAMBIKA` |
| `Activity\SaveDataToOSAksep_Act.xml` baris 8731 | 1 | 2022-01-03 oleh `GABRIELAMILITIA` |
| `Activity\SaveToOS.xml` baris 4511 | 1 | 2022-01-03 oleh `GABRIELAMILITIA` |
| `Activity\SendEmailKlaimRejectClose.xml` baris 693, 884 | 2 | 2023-11-20 oleh `ArlexyVarian` |

Case yang dibuat akun ini mengambil jalur berbeda dari case lain, di empat rule yang berbeda pula.

#### 1.2 Berdasarkan nilai `.KomiteID` dan `.PICSuggest`

```
.KomiteID   == "Himawan" | "NANDINA" | "CHRISTINEANGELINA" | "CHRISTOPMARHASAK"
.PICSuggest == "Himawan" | "Nandina C" | "Christine Angelina Hutagalung"
```

| Rule | Jumlah langkah | Dibuat |
|---|---|---|
| `Activity\CreateChildKomiteCNP_Act.xml` baris 11926 dst. | 4 | 2020-02-18 oleh `MESDISILITONGA` |
| `Activity\SethistoryKlaimTreaty.xml` baris 592, 680, 1168, 1316 | 3 | 2022-12-01 oleh `AnanSosmita` |
| `DataTransform\InsertChronology_DT.xml` baris 426 | 1 | — |

Perhatikan `"Himawan"` muncul dalam dua ejaan berbeda: `KomiteID=="Himawan"` dan `PICSuggest=="Himawan"`, sementara rekannya ditulis `NANDINA` di satu tempat dan `Nandina C` di tempat lain, `CHRISTINEANGELINA` dan `Christine Angelina Hutagalung` di tempat lain lagi. Tidak ada satu bentuk kanonik.

#### 1.3 Berdasarkan isi kolom komentar bebas — bentuk paling rapuh

```
@contains(.CommentSuggest, "Accepted by Himawan")
@contains(.CommentSuggest, "Rejected by Himawan")
@contains(.CommentSuggest, "Himawan")
```
`Activity\SethistoryKlaimTreaty.xml` baris 1256, 1268, 1404, 1416.

Keputusan alur diambil dengan mencari potongan teks di dalam kolom komentar yang diisi bebas oleh pengguna.

### 2. Apa yang terjadi bila string tidak cocok

Ketiga bentuk di atas dipakai sebagai `pyStepsPreCondParamsWhen`. Dalam Pega, precondition yang bernilai salah menyebabkan langkah itu **dilewati**, bukan menimbulkan galat.

Konsekuensinya sama untuk ketiganya, dan tidak bergantung pada tafsir:

| Pemicu | Akibat |
|---|---|
| Akun `VINCENTVERNANDO_1` diganti nama atau dinonaktifkan | 8 langkah berhenti berjalan. Tidak ada pesan galat. |
| Pemegang peran berganti orang | Cabang untuk nama lama tidak pernah terpicu lagi; nama baru tidak punya cabang |
| Komentar diketik `"Accepted By Himawan"` (huruf besar B) | `@contains` gagal, langkah dilewati |
| Komentar diketik `"Accepted by Pak Himawan"` | `@contains(.CommentSuggest,"Himawan")` tetap cocok, tetapi `"Accepted by Himawan"` tidak — dua kondisi bersaudara memberi hasil berbeda atas satu masukan |
| Pengguna mana pun mengetik nama itu di kolom komentar | Kondisi terpenuhi, langkah berjalan |

Butir terakhir adalah sifat mekanismenya: **kolom komentar bebas adalah masukan yang menentukan alur**, dan siapa pun yang boleh mengisi komentar dapat memenuhinya.

### 3. Batas klaim — apa yang TIDAK dapat saya simpulkan dari XML

1. **Apa yang sebenarnya dilakukan tiap cabang secara bisnis** tidak dapat dibaca dari kondisi itu sendiri; yang terbaca hanya langkah yang dijaganya.
2. **Apakah orang-orang itu masih menjabat** adalah pertanyaan organisasi → `_selesai/OPEN-QUESTIONS.md` A10, A11.
3. **Apakah cabang-cabang itu pernah benar-benar terpicu di produksi** adalah pertanyaan data, bukan pertanyaan kode.
4. Sebagian kondisi berada di `CreateChildKomiteCNP_Act`, yang sisi seberangnya ada di modul Komite → sebagian akibatnya `DEFERRED-TO-KOMITE-SESSION`.

### 4. Hubungan dengan sistem baru

Tidak ada cabang berbasis identitas yang boleh diwarisi. Ini bertemu dengan ADR-0006 (`rbac-dirancang-dari-nol`): karena `pyPrivilegeName` kosong di seluruh 279 berkas, satu-satunya model kewenangan yang benar-benar berjalan di sistem lama justru **kondisi-kondisi ini** — kewenangan yang ditulis sebagai nama orang di dalam kode, bukan sebagai peran.

Pemetaan tiap cabang ke **peran** (bukan ke orang) adalah prasyarat desain RBAC, dan itu memerlukan jawaban A10/A11.

### 5. Perilaku tiap cabang, dan kewenangan yang tersirat — EVIDENCED

Diminta pada Round 3: sebelum bertanya ke siapa pun, baca sendiri apa yang berubah saat kondisinya terpenuhi. Hasilnya menjawab sebagian besar pertanyaan tanpa perlu keluar.

| Nama di kode | Rule & langkah | Apa yang berubah saat kondisi terpenuhi | Kewenangan yang tersirat |
|---|---|---|---|
| `CHRISTINEANGELINA` | `CreateChildKomiteCNP_Act` | `.IDKomite="Claim Dept. Head"`, `.KomitePost="Department Head Claim"`, `.Initial="CA"` | **bukan kewenangan** — pemetaan nama ke jabatan |
| `CHRISTOPMARHASAK` | idem | `.IDKomite="Technic Div. Head"`, `.Initial="CM"` | **bukan kewenangan** — pemetaan nama ke jabatan |
| `Himawan` | idem | `.IDKomite="Operational Director"`, `.Initial="HY"` | **bukan kewenangan** — pemetaan nama ke jabatan |
| `NANDINA` | idem | `.IDKomite="Technical Director"`, `.Initial="NC"` | **bukan kewenangan** — pemetaan nama ke jabatan |
| `DARTO` | `CreateChildKomiteCNP_Act` baris 10668 | `.OPERATOR_ID` dipetakan ulang menjadi `"CHRISTINEANGELINA"` | **bukan kewenangan** — penggantian identitas |
| `Christine Angelina Hutagalung` / `Himawan` / `Nandina C` | `SethistoryKlaimTreaty` | `.IsCedingConfirm` diisi nama jabatan yang sama | **bukan kewenangan** — pemetaan nama ke jabatan |
| `VINCENTVERNANDO_1` | `CreateChildKomiteCNP_Act` (4 langkah) | daftar Komite hasil `FilterEmailKomiteWithLimit` **diganti satu anggota tetap**: `KomiteID="VINCENTVERNANDO_1"`, `KomiteEmail="klaim5@nusantarare.com"`, `IDKomite="1"`, `KomitePost="VC"`; dan **`ChildWorkPage.KomiteLoop = 1`** | **KEWENANGAN — mengesampingkan penjenjangan** |
| `VINCENTVERNANDO_1` | `SendEmailKlaimRejectClose` (2 langkah) | `Local.EmailCC` diganti daftar tetap | distribusi pemberitahuan |
| `VINCENTVERNANDO_1` | `SaveToOS` langkah 6.3.10 | melewati **`Call KonversiKlaim_Act`** — langkah berdeskripsi *"Untuk HIT ke arasapas"* | **KEWENANGAN — melewati integrasi hilir** |
| `VINCENTVERNANDO_1` | `SaveDataToOSAksep_Act` langkah 15.8.8 | mengubah cabang di sekitar **`Connect-REST`**, juga berdeskripsi *"Untuk HIT ke arasapas"* | **KEWENANGAN — melewati integrasi hilir** |

#### 5.1 Empat dari lima nama bukan kewenangan sama sekali

Keempatnya hanya menerjemahkan **nama orang menjadi nama jabatan**. Itu tabel referensi yang ditulis sebagai kode:

```
CHRISTINEANGELINA  ->  Claim Dept. Head
CHRISTOPMARHASAK   ->  Technic Div. Head
Himawan            ->  Operational Director
NANDINA            ->  Technical Director
DARTO              ->  (diganti menjadi CHRISTINEANGELINA)
```

Mengikuti kriteria Round 3 — *"kalau ia hanya mengisi nilai default berbeda, itu bukan kewenangan, itu preferensi, dan langsung NOT-MIGRATED tanpa perlu bertanya ke siapa pun"* — kelimanya **`CANDIDATE-NOT-MIGRATED`**, dan tidak perlu dibawa ke pemilik proses. Penggantinya di sistem baru adalah tabel `pengguna → jabatan` biasa.

Dikuatkan oleh `MEMORI_PEMAHAMAN.MD` baris 1167–1174, yang memuat pemetaan yang sama persis dan sudah menandainya sebagai risiko: *"perubahan personel memerlukan perubahan kode."* Label naik dari DERIVED menjadi **EVIDENCED**.

#### 5.2 `VINCENTVERNANDO_1` adalah satu-satunya yang benar-benar kewenangan

Alurnya normal: `FilterEmailKomiteWithLimit` menghasilkan daftar anggota Komite, lalu `ChildWorkPage.KomiteLoop = @SizeOfPropertyList(pyReportContentPage.pxResults)` menetapkan berapa jenjang persetujuan yang dilalui.

Untuk case yang dibuat akun ini, keduanya dilewati: daftarnya diganti **satu anggota tetap**, dan `KomiteLoop` dipaksa **`1`**.

Artinya: **case yang dibuat akun ini hanya melalui satu persetujuan, siapa pun anggotanya seharusnya, berapa pun nilai klaimnya.** Ini wewenang mengesampingkan penjenjangan Komite — kategori paling berat dari tiga kategori di Round 3.

`MEMORI_PEMAHAMAN.MD` baris 1174 menyebut jalur ini juga mem-bypass `KonversiKlaim_Act` dan menyetel `InputData.CARI20 = 1` pada `InsertOSKlaimCNP` — cakupannya lebih luas dari yang terbaca di folder ini saja.

**Hubungan dengan FINDING-001:** keduanya menyentuh mekanisme penjenjangan yang sama. FINDING-001 mempersoalkan ambang yang dibandingkan tanpa konversi mata uang; temuan ini menunjukkan ada jalur yang **melewati ambang itu sama sekali**. Keduanya berdiri sendiri dan tidak saling menggugurkan.

**Satu-satunya pertanyaan yang tersisa untuk dibawa keluar**, karena tidak terjawab XML maupun memori — dan `MEMORI_PEMAHAMAN.MD` sendiri mencatatnya sebagai pertanyaan terbuka nomor 16: **apakah `VINCENTVERNANDO_1` akun uji atau akun bisnis nyata, dan apakah masih aktif.**

---

### 6. Sapuan pola pencocokan teks — seluruh folder

Diminta pada Round 3: yang satu itu hampir pasti bukan satu-satunya. Benar.

**Enam pasang**, bukan dua. Seluruhnya di `SethistoryKlaimTreaty`, dan semuanya melakukan hal yang sama: **menerjemahkan nama orang di dalam teks komentar menjadi nama jabatan**.

| Kondisi | Akibat |
|---|---|
| `@contains(.CommentSuggest,"Accepted by CHRISTINEANGELINA")` | `.CommentSuggest = "Accepted by Kepala Departemen Klaim"` |
| `@contains(.CommentSuggest,"Rejected by CHRISTINEANGELINA")` | `.CommentSuggest = "Rejected by Kepala Departemen Klaim"` |
| `@contains(.CommentSuggest,"Accepted by Himawan")` | `.CommentSuggest = "Accepted by Direktur Operational"` |
| `@contains(.CommentSuggest,"Rejected by Himawan")` | `.CommentSuggest = "Rejected by Direktur Operational"` |
| `@contains(.CommentSuggest,"Accepted by NANDINA")` | `.CommentSuggest = "Accepted by Direktur Teknik"` |
| `@contains(.CommentSuggest,"Rejected by NANDINA")` | `.CommentSuggest = "Rejected by Direktur Teknik"` |

Ditambah `@contains(.CommentSuggest,"NANDINA")` (4x), `"Himawan"` (4x), `"CHRISTINEANGELINA"` (2x) sebagai kondisi tanpa awalan.

#### 6.1 Pola pencocokan teks lain di folder ini

Sapuan `@contains` / `@startsWith` / `@endsWith` / `@indexOf` / `@matches` menemukan 27 ekspresi. Yang **bukan** identitas orang:

| Ekspresi | Sifat |
|---|---|
| `@contains(.TreatyName,"R/I")` (4x) | penanda jenis treaty di dalam nama |
| `@contains(.IsCedingConfirm,"Admin")` (2x) | mencari jabatan di dalam field yang diisi nama jabatan |
| `@contains(SearchData.CARI2,"/R0")`, `@contains(pyWorkPage.ClaimData.NoPla,"/")` | penguraian nomor lewat teks |
| `@contains(InputSpreading.CARI1,"TH"/"ST"/"RD"/"ND")` | **menebak akhiran bilangan urut bahasa Inggris** dari teks |
| `@contains(pyWorkPage.ClaimData.PolicyData.PolicyNo,"RNM-Q")` | penanda jenis polis di dalam nomor |
| `@contains(Local.Emailto,"syariah")` | pemilahan unit bisnis dari alamat surel |
| `@contains(Local.CheckOld,"Previous")` | penanda status di dalam teks |
| `@startsWith(param.WorkStatus,"Resolved-")` | pola standar Pega, wajar |
| `@indexOf("fdf.fsf.sfsdf",".")` (2x) | **string uji coba tertinggal di kode produksi** |

Dua yang paling rapuh setelah identitas orang: `"syariah"` di alamat surel (unit bisnis ditentukan dari isi alamat email) dan `@indexOf("fdf.fsf.sfsdf",".")` yang jelas sisa percobaan.

#### 6.2 Ejaan yang sama tidak konsisten — sebagian cabang tidak akan pernah terpicu

**Koreksi atas versi pertama.** Saya menulis *"orang yang sama ditulis dalam dua ejaan di dalam rule yang sama"*. **Itu keliru.** Penelusuran per baris menunjukkan dua ejaan itu berada di **dua rule yang berbeda**:

| Rule | Bentuk yang dicocokkan |
|---|---|
| `Activity\SethistoryKlaimTreaty.xml` baris 538, 680, 822 | `"Christine Angelina Hutagalung"`, `"Himawan"`, `"Nandina C"` — **bentuk panjang** |
| `DataTransform\InsertChronology_DT.xml` baris 367, 426, 486 | `"CHRISTINEANGELINA"`, `"Himawan"`, `"NANDINA"` — **bentuk pendek** |

Kesimpulannya tetap berlaku, bahkan lebih tajam: **dua rule berselisih tentang bentuk data pada field yang sama.** Karena `.PICSuggest` hanya dapat berisi satu bentuk pada satu waktu:

- **`Himawan`** dieja sama di kedua rule — kedua cabangnya bekerja.
- **`Christine` dan `Nandina`** dieja berbeda — untuk keduanya **tepat satu dari dua rule tidak akan pernah cocok.** Mana yang mati bergantung pada bentuk yang sebenarnya tersimpan di produksi, dan itu pertanyaan data, bukan pertanyaan kode.

Perhatikan pula: kondisi `.PICSuggest` memakai bentuk panjang (`"Christine Angelina Hutagalung"`), sedangkan kondisi `@contains(.CommentSuggest,...)` memakai bentuk pendek (`"CHRISTINEANGELINA"`). Dua mekanisme untuk orang yang sama, dengan ejaan yang berbeda.

---

### 7. Risiko migrasi data — dinyatakan sekarang, bukan ditemukan belakangan

Diminta pada Round 3. Ini konsekuensi langsung dari bagian 6.

**Case yang alurnya pernah ditentukan oleh pencocokan teks tidak menyimpan keputusan itu di tempat lain mana pun selain di dalam teks komentar.**

Saat migrasi, keadaan case seperti itu harus direkonstruksi dari isi `.CommentSuggest`. Yang menggagalkannya:

1. **Teks sudah ditimpa.** Langkah-langkah di bagian 6 **menulis ulang** `.CommentSuggest` — dari `"Accepted by Himawan"` menjadi `"Accepted by Direktur Operational"`. Jadi teks yang tersimpan di produksi sebagian sudah bentuk terjemahan, sebagian masih bentuk asli, tergantung apakah langkah itu sempat berjalan. **Dua populasi bercampur di satu kolom.**
2. **Ejaan berbeda tidak terbaca.** `"Accepted By Himawan"`, `"Accepted by Pak Himawan"`, atau nama yang diketik dengan spasi berlebih tidak akan cocok dengan pola mana pun, dan keputusannya hilang.
3. **Nama di luar kelima itu tidak punya cabang sama sekali.** Persetujuan oleh orang lain tidak pernah diterjemahkan, dan tidak ada jejak strukturnya.

**Akibatnya untuk rencana migrasi:** kolom "siapa menyetujui" pada sistem baru **tidak dapat diisi lengkap dari data lama secara otomatis.** Sebagian akan kosong, dan besarnya bagian yang kosong tidak dapat diketahui dari XML. Kuantifikasinya perlu satu profil data yang murah dan tanpa data nasabah — dimasukkan ke `pengetahuan/PULL-LIST.csv` sebagai **REQ-013**.

Ini **tidak** menghalangi keputusan Round 3 tentang `tindakan + pelaku + waktu` terstruktur. Ia hanya menetapkan bahwa **sebagian riwayat lama tidak akan terbawa**, dan itu harus disepakati sekarang, bukan ditemukan saat UAT.

---

### 8. Jalur yang melewati penjenjangan Komite — `VINCENTVERNANDO_1`

Bagian ini berdiri sendiri, sejajar dengan FINDING-001, karena sifatnya berbeda dari enam bagian di atas. Yang lain adalah pemetaan nama ke jabatan; yang ini mengubah berapa persetujuan yang harus dilalui sebuah klaim.

Netral dan faktual. Tanpa penilaian niat, tanpa menyebut siapa pun bersalah.

#### 8.1 Mekanisme normal

`CreateChildKomiteCNP_Act` menyusun daftar anggota Komite dari hasil Report Definition `FilterEmailKomiteWithLimit`, yang disaring dengan `Param.LIMIT_BOTTOM` menurut nilai klaim. Jumlah jenjang persetujuan lalu ditetapkan dari banyaknya baris hasil:

```
ChildWorkPage.KomiteLoop = @SizeOfPropertyList(pyReportContentPage.pxResults)
```

#### 8.2 Mekanisme untuk case yang dibuat akun ini

Bila `pyWorkPage.pxCreateOperator == "VINCENTVERNANDO_1"`, empat langkah menggantikan keduanya:

| Yang diganti | Nilai pengganti |
|---|---|
| `Primary.ComiteeClaim` | satu baris tetap: `KomiteID="VINCENTVERNANDO_1"`, `KomiteEmail="klaim5@nusantarare.com"`, `IDKomite="1"`, `KomitePost="VC"`, `Initial="VC"`, `KomiteAproval=""` |
| `ChildWorkPage.KomiteList` | satu baris tetap dengan isi yang sama |
| `ChildWorkPage.KomiteLoop` | **`1`** |

`Activity\CreateChildKomiteCNP_Act.xml` baris 10720, 10967, 11195, 12582.

#### 8.3 Akibatnya

**Case yang dibuat akun ini melalui satu persetujuan, berapa pun nilai klaimnya, dan oleh satu alamat tetap.** Hasil `FilterEmailKomiteWithLimit` tidak ikut menentukan apa pun pada jalur ini, sehingga ambang `LIMIT_BOTTOM` — inti dari seluruh mekanisme penjenjangan — tidak berlaku.

Dua langkah lain pada rule berbeda memakai kondisi yang sama: `SendEmailKlaimRejectClose` mengganti `Local.EmailCC` dengan daftar tetap; `SaveDataToOSAksep_Act` dan `SaveToOS` mengubah cabang langkah tanpa `Property-Set` di blok itu sendiri (tujuannya belum ditelusuri).

`MEMORI_PEMAHAMAN.MD` baris 1174 mencatat jalur ini juga mem-bypass `KonversiKlaim_Act` dan menyetel `InputData.CARI20 = 1` pada `InsertOSKlaimCNP` — jadi cakupannya melampaui apa yang terbaca di folder ini.

#### 8.4 Beda derajat kepastian dengan FINDING-001

| | FINDING-001 | Bagian ini |
|---|---|---|
| Yang terbaca | ambang dibandingkan tanpa konversi mata uang | daftar Komite diganti satu orang tetap, `KomiteLoop` dipaksa 1 |
| Penjelasan alternatif yang masuk akal | **ada** — `.ValueAdjustment` mungkin sudah IDR sejak hulu | **tidak ada** — mekanismenya eksplisit dan tidak ambigu |
| Yang belum diketahui | apakah ada Adjustment non-IDR di produksi | apakah akun itu dipakai di produksi, dan sejak kapan |

Keduanya menyentuh mekanisme penjenjangan yang sama dan tidak saling menggugurkan.

#### 8.5 Yang menutup bagian ini

**REQ-016**, prioritas **BLOCKER**, tanpa data nasabah:

| Pertanyaan | Menentukan |
|---|---|
| Berapa case dibuat oleh `VINCENTVERNANDO_1` | bila nol, jalur ini mati dan bagian ini gugur |
| Rentang tanggal case pertama sampai terakhir | akun uji yang tertinggal, atau jalur yang masih hidup |
| Berapa di antaranya bernilai di atas ambang `30.000.000` | besaran paparan |
| Berapa yang tercatat hanya punya satu persetujuan | apakah mekanismenya benar-benar berjalan seperti terbaca |

Pertanyaan "akun uji atau akun bisnis nyata" tetap EXTERNAL (`_selesai/OPEN-QUESTIONS.md` A14) dan sudah tercatat sebagai pertanyaan terbuka nomor 16 di `MEMORI_PEMAHAMAN.MD`.

---

### 9. `DARTO` — penggantian identitas, bukan pemetaan jabatan

Versi pertama menempatkan `DARTO` dalam tabel yang sama dengan empat nama lainnya. **Itu salah penempatan**: jenisnya berbeda.

Empat nama lain adalah pencarian **nama ke jabatan** (`CHRISTINEANGELINA` menghasilkan `"Claim Dept. Head"`). `DARTO` melakukan hal lain:

```
.OPERATOR_ID == "DARTO"   ->   .OPERATOR_ID diganti menjadi "CHRISTINEANGELINA"
```
`Activity\CreateChildKomiteCNP_Act.xml` baris 10580 (deskripsi) dan 10668 (kondisi).

**Yang diganti adalah identitas pelakunya, bukan atributnya.** Tindakan yang dilakukan satu akun tercatat atas nama akun lain.

#### 9.1 Konsekuensinya bukan RBAC, melainkan keutuhan jejak audit

Untuk case yang melewati langkah ini, riwayat di sistem lama **tidak menunjukkan siapa yang sebenarnya bertindak**. Yang tercatat adalah identitas pengganti.

#### 9.2 Konsekuensi migrasi — dinyatakan sekarang

> **Riwayat yang dimigrasi membawa identitas yang mungkin tidak benar, dan tidak ada cara memulihkannya.**

Nilai aslinya tidak disimpan di mana pun sebelum diganti — penggantiannya berupa `Property-Set` langsung, bukan penyalinan ke field cadangan. Setelah `.OPERATOR_ID` bernilai `"CHRISTINEANGELINA"`, tidak ada jejak bahwa sebelumnya ia `"DARTO"`.

Tiga akibat yang harus disepakati, bukan ditemukan belakangan:

1. **Baris riwayat yang terdampak tidak dapat dikoreksi saat migrasi.** Kita tidak tahu baris mana yang terdampak, karena hasilnya tidak dapat dibedakan dari baris yang memang dibuat `CHRISTINEANGELINA`.
2. **Jumlahnya tidak dapat diukur** — bahkan dengan query. Tidak ada penanda.
3. **Sistem baru tidak boleh menyediakan mekanisme serupa.** Perwakilan atau pendelegasian, bila memang dibutuhkan, dicatat sebagai *"X bertindak untuk Y"* dengan kedua identitas tersimpan — bukan dengan menimpa salah satunya.

---

### 10. Cabang yang tidak mungkin terpicu — otomatis tidak dimigrasi

Diminta pada Round 4. Hasil sapuan seluruh `pyStepsPreCondParamsWhen` di 279 berkas.

#### 10.1 Yang terbukti

| Cabang | Kenapa mati | Rule |
|---|---|---|
| `.PICSuggest=="CHRISTINEANGELINA"` **atau** `"Christine Angelina Hutagalung"` — salah satu | Dua rule mencocokkan bentuk berbeda atas field yang sama | `InsertChronology_DT` vs `SethistoryKlaimTreaty` |
| `.PICSuggest=="NANDINA"` **atau** `"Nandina C"` — salah satu | idem | idem |

Mana dari pasangan itu yang mati **belum dapat ditentukan dari XML** — ia bergantung bentuk yang tersimpan di produksi. Yang pasti: untuk setiap pasangan, **tepat satu sisi tidak pernah cocok.**

#### 10.2 Yang saya cari tetapi TIDAK ditemukan

Supaya tidak dibaca sebagai daftar yang lebih besar dari sebenarnya, ini yang saya uji dan hasilnya nihil:

| Pola yang dicari | Hasil |
|---|---|
| Kondisi membandingkan literal dengan literal (`"a"=="b"`) | **TIDAK DITEMUKAN DI XML** |
| Kondisi self-contradictory (`X=="a" && X=="b"`) | **TIDAK DITEMUKAN DI XML** |
| String uji coba di dalam **kondisi** (`fdf`, `asdf`, `test123`) | **TIDAK DITEMUKAN DI XML** |

Catatan atas yang terakhir: `@indexOf("fdf.fsf.sfsdf",".")` memang ada, tetapi **bukan** sebagai precondition — ia berada di dalam ekspresi nilai. Jadi ia sisa percobaan yang tetap dieksekusi, bukan cabang mati. Saya sempat menggolongkannya bersama cabang mati pada ringkasan Round 4; itu penggolongan yang salah.

Satu kandidat lagi yang saya periksa dan **gugur**: `HitServiceToKasir_Act` memuat `(pyWorkIDPrefix=="CLMP-" && ...) || (pyWorkIDPrefix=="CLM-" && ...)`. Terlihat seperti kontradiksi, tetapi itu dua kelompok yang di-OR — keduanya dapat terpenuhi. Bukan cabang mati.

#### 10.3 Berapa besar penghematannya

Dua pasang cabang, tiga langkah. **Lebih kecil dari dugaan awal.** Penghematan lingkup kerjanya nyata tetapi kecil, dan saya lebih baik menyampaikan angka yang benar daripada daftar yang panjang.

#### 8.6 Cakupannya lebih luas dari penjenjangan Komite

Telusuran langkah tujuan pada `SaveToOS` dan `SaveDataToOSAksep_Act` selesai, dan hasilnya memperluas temuan ini.

| Rule | Langkah | Yang dilewati |
|---|---|---|
| `SaveToOS` | `pySteps(6).pySteps(3).pySteps(10)` | `Call KonversiKlaim_Act` — *"Untuk HIT ke arasapas"* |
| `SaveDataToOSAksep_Act` | `pySteps(15).pySteps(8).pySteps(8)` | cabang di sekitar `Connect-REST` — *"Untuk HIT ke arasapas"* |

Jadi untuk case yang dibuat akun ini, bukan hanya penjenjangan Komite yang diganti — **pengiriman ke Arasapas, sistem hilir, juga tidak berjalan sebagaimana case lain.**

Ini mengonfirmasi catatan `MEMORI_PEMAHAMAN.MD` baris 1174 yang menyebut jalur ini mem-bypass `KonversiKlaim_Act`, dan menaikkan statusnya dari **memori** menjadi **EVIDENCED**.

**Akibatnya untuk REQ-016**: pertanyaannya bertambah satu. Bukan hanya berapa case yang dibuat akun ini dan berapa persetujuan yang dilaluinya, tetapi juga **berapa di antaranya tidak pernah sampai ke Arasapas**. Bila ada case produksi yang tidak terkirim ke sistem hilir, itu selisih data antar sistem yang belum pernah masuk peta mana pun.

## FINDING-003 — Baris Retensi Cedant hadir di `SpreadingRisk` saat rule tanpa penyaring membacanya

Sumber: `FINDING-003-baris-ur-tanpa-penyaring.md`

**Jenis**: laporan kondisi sistem lama
**Status**: KEHADIRAN **MUNGKIN**, BUKAN PASTI — bergantung urutan tindakan pengguna (lihat 3.3, koreksi atas versi pertama)
**Ruang lingkup bukti**: `D:\XML_NURE\Claim Non Prop` saja
**Tanggal**: 18 September 2026

---

### 1. Koreksi atas pernyataan saya sendiri

Ronde lalu saya menutup temuan ini dengan *"itu membuktikan tidak adanya penyaring, bukan adanya salah hitung"*. Batas itu benar pada saat itu, tetapi saya berhenti satu langkah terlalu awal. Urutan pemanggilan dapat ditelusuri seluruhnya dari XML dan dari dokumen struktur milik klien — tanpa Oracle sama sekali. Berikut hasilnya.

Sekaligus saya luruskan dua pernyataan saya yang tampak bertentangan: Ronde 2 saya menyebut tiga rule **menyaring** `TreatyName=="UR"`; Ronde 3 saya menyebut sembilan rule **tidak menyebut** `"UR"`. Keduanya benar dan tidak bertabrakan — sebagian menyaring, sebagian tidak. Yang belum saya kerjakan waktu itu adalah memeriksa apakah yang tidak menyaring memang melihat baris itu.

### 2. Kapan baris `"UR"` masuk ke PageList

`Activity\CountLossAllocation_act.xml` baris 7015–7120:

```
pyWorkPage.ClaimData.SpreadingRisk(<APPEND>).TreatyType      = "UR"
pyWorkPage.ClaimData.SpreadingRisk(<LAST>).TreatyName        = "UR"
pyWorkPage.ClaimData.SpreadingRisk(<LAST>).ClaimEstimation   = Local.UR
pyWorkPage.ClaimData.SpreadingRisk(<LAST>).ClaimAmountAdjust = Local.UR
```

Retensi Cedant bukan entitas terpisah — ia baris tambahan di PageList yang sama dengan Layer.

### 3. Urutan pemanggilan — DIBUKTIKAN ULANG DARI XML

**Koreksi.** Versi pertama bagian ini memakai `Struktur_Flow_TreatyIn.xlsx` sebagai bukti. Berkas itu turunan, bukan sumber, dan tidak boleh berdiri sebagai bukti. Seluruh bagian ini ditulis ulang dari XML. **Hasilnya berbeda dari versi pertama, dan versi pertama salah pada satu hal pokok** — lihat sub-bagian 3.3.

#### 3.1 Siapa memanggil apa (dari `pyStepsActivityName` di XML)

Indeks panggil activity ke activity, dibaca dari `<pyStepsActivityName>Call ...</pyStepsActivityName>` beserta `<pyStepPageReference>`:

| Dipanggil | Oleh | Langkah |
|---|---|---|
| `CountLossAllocation_act` | `CountClaimTNP_Act` | 15 |
| `CountLossAllocation_act` | `AddAkseptasiCNP_Act` | 18 |
| `CountLossAllocation_act` | `InputAkseptasi_PreAct` | 5 |
| `GenerateCFS_act` | `SaveToOS` | 7 |
| `CountTotalInsterest_Act` | `SetCurencyInterest_act` | 6 |
| `CountSpreadingXOL` | `AdjClaimCNP_Act` | 16 |
| `CopyOldataCurr_act` | `CountTotalInsterest_Act` | 11 |
| `SendEmailKlaim` | `CreateChildKomiteCNP_Act` | 35 |

#### 3.2 Temuan pokok: urutannya tidak ditentukan kode, melainkan layar

**Tidak ada satu pun activity di folder ini yang memanggil `CountClaimTNP_Act`.** Rule itu dipicu dari **Section**, lewat `<pyActivity>`, sebanyak 42 rujukan di tiga berkas:

| Section | Activity yang dipicu, dalam urutan dokumen |
|---|---|
| `Section/AdjustmentDetailNP.xml` | `CountClaimTNP_Act`, `CountLossAllocation_act`, `AdjClaimCNP_Act` |
| `Section/InputAcceptation.xml` | `SetCurencyInterest_act`, `CountTotalInsterest_Act`, `CountClaimTNP_Act`, `CountLossAllocation_act`, `SaveToOS`, `AddAkseptasiCNP_Act`, `DeleteAkseptasi_Act` |
| `Section/OutstandingClaim(1).xml` | `CountClaimTNP_Act`, `CountLossAllocation_act`, `GenerateCFS_act` |
| `Section/ViewDetailDeptHeadTreatyIn_UW.xml` | `CountSpreading_Act` |

Ini **kontrol UI yang berdiri sendiri-sendiri pada satu layar**, bukan rantai berurutan. Urutan dokumen di dalam XML Section adalah urutan tata letak, **bukan** urutan eksekusi.

Konsekuensinya, dan inilah temuan sebenarnya:

> **Urutan jalannya perhitungan alokasi ditentukan oleh kontrol mana yang ditekan pengguna, bukan oleh rantai pemanggilan di dalam kode.** Tidak ada satu tempat pun di XML yang menetapkan bahwa `GenerateCFS_act` berjalan sesudah `CountLossAllocation_act`.

#### 3.3 Yang harus saya cabut

Versi pertama menyatakan: *"baris Retensi Cedant sudah ada di dalam `SpreadingRisk` ketika ketiga rule itu membacanya - ini sekarang fakta, bukan dugaan."*

**Pernyataan itu saya cabut.** Ia bersandar pada penomoran outline spreadsheet. Dari XML, yang dapat dinyatakan hanyalah:

- `CountLossAllocation_act` (yang meng-`<APPEND>` baris `"UR"`) dan `GenerateCFS_act` dapat dipicu dari **layar yang sama** (`OutstandingClaim(1)`), begitu pula `CountTotalInsterest_Act` dan `DeleteAkseptasi_Act` pada `InputAcceptation`.
- Karena itu baris `"UR"` **dapat** sudah ada saat rule-rule itu membaca `SpreadingRisk`, tergantung urutan penekanan kontrol oleh pengguna.
- Apakah ia **selalu** ada: **TIDAK DAPAT DISIMPULKAN DARI XML.** Itu bergantung perilaku pengguna, dan hanya terbaca dari data produksi.

Yang justru menguat dari koreksi ini: ketergantungan pada urutan klik itu sendiri adalah cacat rancangan, terlepas dari apakah salah hitungnya pernah terjadi.

### 4. Akibatnya berbeda per rule — dan satu di antaranya ternyata aman

Kehadiran baris itu tidak otomatis berarti salah hitung. Saya periksa satu per satu:

#### 4.1 `CountTotalInsterest_Act` — **TIDAK ADA SALAH HITUNG**

```
Local.SizeLoss = @SizeOfPropertyList(pyWorkPage.ClaimData.SpreadingRisk)
Local.SizeEst  = @SizeOfPropertyList(pyWorkPage.ClaimData.SpreadingRisk)
```
Cacahan itu memang ikut menghitung baris `"UR"`. Tetapi pemakaiannya hanya sebagai penjaga kosong-tidaknya:
```
Local.SizeLoss > 0 && Local.PropCount == 1
Local.SizeEst  > 0 && Local.PropCount == 1
```
Satu baris tambahan tidak mengubah hasil `>0`. **Rule ini bersih.** Saya catat ini karena dugaan awal saya mengarah ke sebaliknya.

#### 4.2 `GenerateCFS_act` — **TUNTAS: baris UR ikut tersalin, bila ia ada**

Rantai `Local.XOL` / `Local.Check` selesai ditelusuri.

**Bentuk rantainya** — berulang empat kali di dalam rule, dengan pola yang sama persis:

```
Local.Check = 0                        <- penanda direset
Local.XOL   = .TreatyName              <- disalin dari BARIS YANG SEDANG DIBACA
Local.Check = 1                        <- penanda dinyalakan
... precondition: Local.Check == 0     <- penjaga agar satu TreatyName diproses sekali
... precondition: Local.XOL == .TreatyName  <- pencocokan baris berikutnya
```

Empat penetapan `Local.XOL`, empat belas penetapan `Local.Check`, sebelas precondition yang membacanya.

**Yang menentukan**: `Local.XOL = .TreatyName` menyalin nama treaty **dari baris yang sedang dibaca, apa adanya**. Tidak ada penyaring. Dan sapuan sebelumnya sudah memastikan `GenerateCFS_act` **tidak menyebut `"UR"` satu kali pun** di seluruh 66 rujukan `SpreadingRisk`-nya.

**Kesimpulan**: bila baris `TreatyName = "UR"` ada di `SpreadingRisk` saat rule ini berjalan, maka `Local.XOL` akan bernilai `"UR"`, barisnya lolos seluruh precondition, dan ia **tersalin ke `TempDataOutStanding.ClaimData.SpreadingRisk(<APPEND>)` beserta `RetroList`-nya** — yaitu ke struktur retrosesi, diperlakukan persis seperti Layer.

**Syaratnya tetap satu, dan itu tidak berubah**: baris `"UR"` harus sudah ada saat rule ini dipicu. Karena urutan pemicuan ditentukan kontrol yang ditekan pengguna (bagian 3.2), kehadirannya **mungkin, bukan pasti** — persis batas yang ditetapkan di bagian 3.3.

Jadi yang tadinya *"belum dapat dipastikan apakah tersalin"* sekarang menjadi *"pasti tersalin bila ada"*. Yang tersisa hanya pertanyaan kehadiran, dan itu pertanyaan data — bukan lagi pertanyaan kode.

#### 4.3 Enam rule sisanya — **selesai, dan hasilnya menggugurkan sebagian temuan ini**

`CountSpreading_Act`, `CountSpreadingXOL`, `DeleteAkseptasi_Act`, `SetAccoutNo_Act`, `SendEmailKlaim`, `CountClaimTNP_Act`: **nol rujukan ke PageList `SpreadingRisk`.** Seluruh kemunculan teks `SpreadingRisk` pada keenamnya adalah nama class di metadata langkah, bukan pembacaan daftar.

`CopyOldataCurr_act` punya dua rujukan nyata: `@LengthOfPageList(...)>0` sebagai penjaga, dan satu `Property-Set` `repeat=EMBEDDED` yang menyetel `.Currency = Local.CurrencyNew` pada setiap baris. Baris `"UR"` ikut terkena, tetapi baris itu memang membawa `Currency` sendiri (disetel `CountLossAllocation_act`), jadi memperbaruinya bersama baris lain konsisten — **bukan cacat**.

**Ringkas**: dari sembilan rule yang semula saya sebut terdampak, **satu terbukti terdampak** (`GenerateCFS_act`), **dua terbukti tidak**, dan **enam tidak membaca daftarnya sama sekali**. Angka "16–66 rujukan" pada versi pertama keliru — lihat koreksi di `BLUEPRINT.md` §2.4.

### 5. Yang menutup temuan ini

Seluruhnya pekerjaan XML. **Tidak satu pun butuh Oracle.**

1. Telusuri rantai `Local.XOL` / `Local.Check` di `GenerateCFS_act`.
2. Periksa enam rule sisanya dengan cara yang sama: apakah baris `"UR"` ikut terbaca, dan apakah pembacaan itu mengubah angka.
3. Tetapkan mana yang memang **seharusnya** menyertakan Retensi Cedant — sebagian mungkin benar demikian. Itu pertanyaan maksud bisnis, bukan pertanyaan kode → Round 3.

## FINDING-004 — Hasil suntingan manual alokasi tidak terlindungi

Sumber: `FINDING-004-suntingan-manual-tertimpa.md`

**Jenis**: laporan kondisi sistem lama (bukan keputusan migrasi)
**Status**: BERSYARAT — berlaku bila hipotesis F2 benar, yaitu `.IsEditClaim` memang tidak punya pembaca di mana pun
**Ruang lingkup bukti**: `D:\XML_NURE\Claim Non Prop` saja
**Tanggal**: 18 September 2026

Dokumen ini terpisah dari keputusan desain untuk sistem baru. Keputusan itu ada di ADR-0008 dan berlaku apa pun hasil temuan ini.

---

### 1. Mekanismenya

`IsEditClaim` adalah `Rule-Obj-Property` milik class `ASM-FW-GISFW-Data-SpreadingRisk`.

| Rule | Aksi |
|---|---|
| `Activity\EditXOLAlokasi.xml` | `.IsEditClaim = 1` — ditandai saat petugas menyunting alokasi XOL |
| `Activity\CountLossAllocation_act.xml` (dua langkah) | `SpreadingRisk(<LAST>).IsEditClaim = 0` — ditimpa kembali menjadi 0 |

Lima kemunculan di seluruh 279 berkas, tidak lebih. **Tidak ada satu pun pembacaan**: tidak di `pyStepsPreCondParamsWhen`, tidak di `pyVisible` section mana pun, tidak di Report Definition, tidak di SQL.

### 2. Akibatnya bila memang tidak ada pembaca

Bendera yang hanya ditulis tidak menjaga apa pun. Rangkaiannya:

1. Petugas menyunting nilai alokasi lewat `EditXOLAlokasi`.
2. Sistem menandai baris itu `IsEditClaim = 1`.
3. Perhitungan alokasi dijalankan lagi.
4. `CountLossAllocation_act` menghitung ulang nilainya **dan** menyetel penandanya kembali ke `0`.
5. Tidak ada peringatan, tidak ada pesan, tidak ada jejak bahwa suntingan pernah ada.

**Berapa besar peluang langkah 3 terjadi:** `CountLossAllocation_act` dapat dipicu dari tiga jalur yang terbaca di XML — `CountClaimTNP_Act` langkah 15, `AddAkseptasiCNP_Act` langkah 18, `InputAkseptasi_PreAct` langkah 5 — dan **juga langsung dari kontrol UI** pada tiga Section (`AdjustmentDetailNP`, `InputAcceptation`, `OutstandingClaim(1)`). Lihat FINDING-003 bagian 3.

Karena pemicunya adalah kontrol layar, petugas yang membuka kembali layar yang sama dan menekan kontrol hitung akan menghapus suntingannya sendiri tanpa menyadarinya.

### 3. Yang membuat ini bukan sekadar bug kecil

Fitur suntingnya **ada dan berfungsi** — layarnya ada, penandanya ditulis. Yang tidak ada adalah pihak yang membaca penanda itu. Bagi penggunanya, sistem tampak menyediakan kemampuan yang sebenarnya tidak dijamin.

Ini berbeda dari fitur yang tidak ada sama sekali: fitur yang tidak ada tidak menyesatkan siapa pun.

### 4. Konsekuensi yang meringankan pekerjaan migrasi

Satu akibat yang justru menyederhanakan, dan perlu dinyatakan supaya tidak dikerjakan dua kali:

> **Bila suntingan selalu tertimpa, maka nilai yang tersimpan di produksi hari ini adalah nilai hasil hitungan, bukan nilai hasil suntingan.**

Artinya **tidak ada "nilai tersunting" yang perlu diselamatkan** saat migrasi data. Tidak perlu kolom khusus, tidak perlu rekonsiliasi, tidak perlu menanyakan ke pengguna nilai mana yang benar. Satu cabang pekerjaan migrasi data terpotong.

Ini berlaku **hanya bila hipotesis F2 benar.** Bila F1 yang benar (pembacanya ada di modul Komite), kesimpulan ini gugur dan nilai tersunting mungkin memang bertahan.

### 5. Batas klaim

1. **Hipotesis F1 belum disingkirkan.** `.IsEditClaim` mungkin dibaca oleh rule di modul Komite. Folder Komite tidak dibuka pada sesi ini. Lihat `_selesai/OPEN-QUESTIONS.md` bagian F.
2. **Ketiadaan pembacaan disimpulkan dari ketiadaan rule, bukan dari adanya rule.** Ekspor ini memuat 279 berkas dari satu folder; pembacaan bisa saja berada di rule yang tidak ikut terekspor.
3. **Belum diuji apakah suntingan pernah benar-benar dilakukan di produksi.** Kalau fitur `EditXOLAlokasi` tidak pernah dipakai, temuan ini benar secara mekanis tetapi nol dampaknya.

### 6. Yang menutup temuan ini

**Satu pemeriksaan tunggal**: sapu modul Komite atas `IsEditClaim`. Satu sesi, satu grep. Hasilnya menentukan F1 atau F2, dan sekaligus menutup atau menegakkan seluruh dokumen ini.

Sampai itu terjadi, **tidak boleh** ada keputusan desain yang mengandaikan salah satunya benar. Keputusan untuk sistem baru sudah diambil terpisah dan tidak menunggu hasil ini — lihat ADR-0008.

## FINDING-005 — Tanggal kerugian dibandingkan terhadap tanggal akhir treaty dalam format berbeda

Sumber: `FINDING-005-perbandingan-tanggal-beda-format.md`

**Jenis**: laporan kondisi sistem lama
**Status**: TERBUKTI DARI XML — tidak menunggu data Oracle
**Ruang lingkup bukti**: `D:\XML_NURE\Claim Non Prop` saja
**Tanggal**: 18 September 2026

---

### 1. Dua format, satu perbandingan

#### 1.1 `.DateOfLoss` berformat `YYYYMMDD`

Dibuktikan dua kali, dari dua arah berbeda:

| Bukti | Lokasi |
|---|---|
| Diurai per posisi: `@substring(.DateOfLoss,0,4)` = tahun, `(4,6)` = bulan, `(6,8)` = tanggal | beberapa Section/Activity |
| Disambung menjadi DateTime Pega: `@FormatDateTime((pyWorkPage.ClaimData.DateOfLoss + "T050000.000 GMT"), "dd/MM/yyyy", ...)` | `Activity\GenerateCACNP_Act.xml` baris 1416, `Activity\GenerateCFS_act.xml` baris 1146 |

Format `YYYYMMDDTHHMMSS.mmm GMT` adalah bentuk DateTime internal Pega. Penyambungan itu hanya sah bila bagian kirinya `YYYYMMDD`.

Dikuatkan dari sisi basis data: kolom `DATEOFLOSS` pada tabel work Pega bertipe `VARCHAR2(8 BYTE)` — delapan karakter, persis panjang `YYYYMMDD`. *(Sumber: disampaikan pengguna; DDL tabel work belum masuk ke berkas `.xls` — lihat catatan di `BLUEPRINT.md` bagian 13.)*

#### 1.2 `.EndDateTreaty` berformat `dd/MM/yyyy`

| Bukti | Lokasi |
|---|---|
| `@FormatDateTime(...,"dd/MM/yyyy","Asia/Jakarta","in_ID") == pyWorkPage.ClaimData.EndDateTreaty` | `Activity\CheckPeriodPolicy_Act.xml` baris 401 |
| `@CompareDates(@FormatDateTime(...,"dd/MM/yyyy",...), pyWorkPage.ClaimData.EndDateTreaty)` | `Activity\CheckPeriodPolicy_Act.xml` baris 314 |
| `@FormatDateTime(...,"dd/MM/yyyy",...) > pyWorkPage.ClaimData.EndDateTreaty` | `Activity\SetEndDate_Act.xml` baris 1643 |
| Dideklarasikan `TYPE="STRING"` | definisi parameter |

Tiga tempat mengubah nilai lain ke `dd/MM/yyyy` **lebih dulu** sebelum membandingkannya dengan `.EndDateTreaty`. Itu memperlihatkan format yang diharapkan.

#### 1.3 Perbandingan yang tidak menyamakan format

```
pyWorkPage.ClaimData.DateOfLoss > pyWorkPage.ClaimData.EndDateTreaty
```
`Activity\CheckDateDOL_Act.xml` baris **4011** dan **4611**.

Kedua sisi adalah teks. Yang dibandingkan: `"20260315"` terhadap `"31/12/2026"` — **perbandingan leksikografis atas dua format yang berbeda**, bukan perbandingan tanggal.

### 2. Akibatnya

Hasil perbandingan ditentukan oleh **digit pertama** masing-masing teks, bukan oleh tanggalnya:

| `.DateOfLoss` | `.EndDateTreaty` | Perbandingan teks | Yang benar secara tanggal |
|---|---|---|---|
| `20260315` (15 Mar 2026) | `31/12/2026` | `"2" < "3"` → **false** | 15 Mar 2026 < 31 Des 2026 → false |
| `20260315` (15 Mar 2026) | `01/01/2026` | `"2" > "0"` → **true** | 15 Mar 2026 > 1 Jan 2026 → true |
| `20251115` (15 Nov 2025) | `01/01/2026` | `"2" > "0"` → **true** | 15 Nov 2025 < 1 Jan 2026 → **false** |

Baris ketiga adalah kesalahannya: hasilnya benar hanya **kebetulan**, ketika digit pertama tanggal treaty kebetulan selaras. Secara umum, **hasil perbandingan bergantung pada digit pertama hari dalam bulan pada tanggal akhir treaty** — `0`, `1`, `2`, atau `3` — dan sama sekali tidak bergantung pada tahunnya.

Rule yang memuatnya bernama `CheckDateDOL_Act` — pemeriksa tanggal kerugian. Yaitu tepat penjaga yang seharusnya menolak klaim dengan tanggal kejadian di luar masa berlaku treaty.

### 3. Pola yang sama di rule tetangga

`Activity\CheckPeriodPolicy_Act.xml` memuat keduanya sekaligus:

| Baris | Ekspresi | Menyamakan format? |
|---|---|---|
| 314, 401 | `@FormatDateTime(...,"dd/MM/yyyy",...)` dibandingkan ke `.EndDateTreaty` | **ya** |
| 868 | `pyWorkPage.ClaimData.PolicyData.EndDateTime > pyWorkPage.ClaimData.EndDateTreaty` | **tidak** — DateTime Pega dibandingkan langsung ke teks `dd/MM/yyyy` |

Jadi di dalam satu rule, sebagian perbandingan menyamakan format dan sebagian tidak. Ini menunjukkan penyebabnya bukan ketidaktahuan, melainkan ketidakseragaman.

### 4. Batas klaim

1. **Belum diuji apakah kesalahannya pernah berakibat di produksi.** Yang terbukti adalah perbandingannya salah bentuk. Berapa klaim yang lolos atau tertolak karenanya adalah pertanyaan data → **REQ-020**.
2. **Nilai `.EndDateTreaty` yang sebenarnya tersimpan belum diperiksa.** Format `dd/MM/yyyy` disimpulkan dari cara nilai lain diubah sebelum dibandingkan dengannya, bukan dari contoh isinya.
3. **Penulis `.DateOfLoss` tidak ditemukan.** Sejalan dengan temuan sebelumnya bahwa `.DateOfLoss` tidak pernah ditulis oleh activity mana pun di folder ini (ADR-0001).

### 5. Hubungan dengan sistem baru

Tanggal disimpan sebagai tipe tanggal, bukan teks, dan seluruh perbandingan dilakukan atas tipe tanggal. Itu sudah menjadi konsekuensi wajar dari ADR-0001 dan tidak memerlukan keputusan baru.

Yang memerlukan keputusan: **apa yang dilakukan terhadap klaim lama yang lolos penjaga ini karena perbandingannya salah.** Itu pertanyaan migrasi data, bukan pertanyaan rancangan.

---

### 6. Tambahan — asal-usul format `dd/MM/yyyy` ditemukan

Versi pertama menyimpulkan format `.EndDateTreaty` dari cara nilai lain dinormalkan sebelum dibandingkan dengannya. Pembacaan `pengetahuan/ddl/VIEW_V_POLIS.sql` menunjukkan sumbernya langsung.

`V_POLIS` — 26 kolom, seluruhnya proyeksi atas **satu** tabel `json_polis`, dengan 32 pemanggilan `JSON_VALUE`. Kolom tanggalnya dibentuk begini:

```sql
SUBSTR(JSON_VALUE (DATA_JSON, '$.PolicyData.StartDateTime'), 7, 2)
|| '/' || SUBSTR(JSON_VALUE (DATA_JSON, '$.PolicyData.StartDateTime'), 5, 2)
|| '/' || SUBSTR(JSON_VALUE (DATA_JSON, '$.PolicyData.StartDateTime'), 1, 4)
|| ' ' || SUBSTR(..., 10, 2) || ':' || SUBSTR(..., 12, 2) || ':' || SUBSTR(..., 14, 2)
```

Jadi rantainya lengkap sekarang:

| Tahap | Bentuk |
|---|---|
| Tersimpan di JSON | `YYYYMMDDHHMISS` (teks) |
| Diubah oleh `V_POLIS` dengan `SUBSTR` + `\|\|` | `dd/MM/yyyy HH:MI:SS` (teks) |
| Dibaca Pega sebagai `.EndDateTreaty` | teks `dd/MM/yyyy` |
| `.DateOfLoss` | tetap `YYYYMMDD` |

**Penyebabnya bukan kesalahan di satu rule.** Basis data mengubah format tanggal menjadi teks bergaya Indonesia di lapisan view, sementara `.DateOfLoss` tidak melewati lapisan itu dan tetap dalam bentuk aslinya. Perbandingan di `CheckDateDOL_Act` mempertemukan keduanya.

Tidak ada satu pun kolom tanggal di `V_POLIS` yang bertipe `DATE`. Semuanya teks hasil sambungan.

#### 6.1 Dua bentuk JSON di satu tabel, dan satu di antaranya terpotong

`V_POLIS` bercabang pada `a.data_json.QuotationData.BusinessFac`:

| Nilai | Jalur JSON yang dibaca | Hasilnya |
|---|---|---|
| `'F'` (Facultative) | `$.PolicyData.StartDateTime`, `$.PolicyData.EndDateTime` | `dd/MM/yyyy HH:MI:SS` — lengkap |
| `'T'` (Treaty) | `$.StartDate`, `$.EndDate` | `dd/MM/yyyy HH` — **berhenti di jam, tanpa menit dan detik** |

Dua hal terbaca dari sini:

1. **Satu tabel `json_polis` menampung dua bentuk dokumen yang berbeda**, dibedakan `BusinessFac`. Jalur JSON-nya tidak sama.
2. **Cabang Treaty menghasilkan teks tanggal yang lebih pendek.** Karena perbandingan atasnya adalah perbandingan teks, dua nilai dengan panjang berbeda tidak dapat dibandingkan secara andal. Klaim non-proporsional berjalan di jalur Treaty — jalur yang terpotong itu.

**Batas klaim**: apakah pemotongan itu berakibat bergantung pada perbandingan mana yang benar-benar dijalankan atas kolom tersebut. Belum ditelusuri seluruhnya.

## FINDING-006 — Kurs yang tidak ditemukan dikembalikan sebagai 1

Sumber: `FINDING-006-kurs-tidak-ditemukan-bernilai-satu.md`

**Jenis**: laporan kondisi sistem lama
**Status**: TERBUKTI DARI DDL — besarannya belum terukur
**Ruang lingkup bukti**: `pengetahuan/DDL_Script_ClaimNonProp.xls`
**Tanggal**: 18 September 2026

Dipisahkan dari `FINDING-001` atas permintaan Round 6. Uraian penuh rantai buktinya tetap di `FINDING-001` bagian 7; dokumen ini memuat temuannya sendiri supaya dapat dirujuk terpisah.

---

### 1. Mekanismenya

`POOLDATA.GETCURRENCYSTANDARD` berakhir dengan:

```sql
EXCEPTION
   WHEN NO_DATA_FOUND
   THEN
      RETURN 1;
END;
```

Bila tidak ada baris kurs untuk suatu mata uang, fungsi mengembalikan **`1`** — yang secara aritmetika berarti konversi satu banding satu.

### 2. Akibatnya

Nilai 1.000 USD diperlakukan sebagai 1.000 IDR. Tanpa galat, tanpa pesan, tanpa penanda.

Ini pilihan yang disengaja — ada penangan eksepsi khusus untuknya, bukan kelalaian. Konsekuensinya tetap sama: **kegagalan pencarian kurs tidak dapat dibedakan dari kurs yang memang bernilai 1.**

Rupiah terhadap rupiah memang berkurs 1. Jadi baris hasil tidak menyimpan informasi apa pun tentang apakah konversinya terjadi atau gagal.

### 3. Pola yang sama muncul tiga kali di sesi ini

Ini bukan kejadian tunggal, melainkan satu kebiasaan yang berulang: **kegagalan yang menyamar menjadi hasil yang sah.**

| Mekanisme | Kegagalan | Terlihat sebagai |
|---|---|---|
| `GETCURRENCYSTANDARD` | kurs tidak ada | kurs bernilai 1 |
| Database link ke HRD | tautan putus | daftar pengguna kosong |
| `pyStepsPreCondParamsWhen` yang tidak cocok | kondisi gagal | langkah dilewati tanpa pesan |

Ketiganya tidak menghasilkan galat. Ketiganya menghasilkan jawaban yang bentuknya benar dan isinya salah.

### 4. Batas klaim

1. **Belum diketahui apakah ada mata uang tanpa baris kurs di produksi.** Bila `m_currencystandard` lengkap untuk semua mata uang yang dipakai, jalur `NO_DATA_FOUND` tidak pernah tersentuh dan temuan ini nol dampaknya.
2. **Tidak dapat diukur dari hasilnya.** Karena kurs 1 sah untuk IDR, mencari "baris berkurs 1" akan mencampur kegagalan dengan kasus normal. Pengukurannya harus dari sisi lain: cari mata uang **bukan IDR** yang kursnya tersimpan `1` — itu hampir pasti hasil jalur kegagalan. → **REQ-022**.
3. Berlaku juga untuk `result_cache` pada fungsi itu: hasil `1` yang pernah dikembalikan dapat tersimpan di cache.

### 5. Hubungan dengan sistem baru

Sudah diputuskan terpisah di **ADR-0014**: konversi tanpa kurs tidak menghasilkan angka. Nilai IDR dibiarkan kosong dan ditandai menunggu kurs, sehingga tidak ada angka salah yang ikut ke perhitungan lain.

## FINDING-007 — Dua rumus premi reinstatement bekerja atas nilai yang berbeda

Sumber: `FINDING-007-dua-rumus-reinstatement.md`

**Jenis**: laporan kondisi sistem lama
**Status**: TERBUKTI DARI XML
**Tanggal**: 18 September 2026

> **Satu pola, tiga temuan.** Ini kelas cacat yang sama dengan `CheckDateDOL_Act` (FINDING-005) dan `.IsEditClaim` (FINDING-004): nilai dibaca dalam keadaan yang tidak dimaksudkan, dan tidak ada yang menjaganya.

---

### 1. Dua rumus

| Berkas & baris | Langkah | Ekspresi |
|---|---|---|
| `Activity\CountReinstatement_Act.xml` **624** | — | `((((.ClaimEstimation + .AdjusterFee) − (.Salvage×100/RNMShare)) / .CNPLimit) × .CNPMDP) × (.CNPPctReinstate/100)` |
| `Activity\AdjClaimCNP_Act.xml` **3109** | `RH_1.pySteps(11).pySteps(1)` | `@if(.TreatyType=="UR", 0, @divide(.TotalClaim, .CNPLimit, 20) × .CNPMDP × @divide(.CNPPctReinstate, 100, 20))` |

Keduanya setara **hanya bila** `.TotalClaim` bernilai sama dengan `ClaimEstimation + AdjusterFee − Salvage×100/RNMShare` pada saat baris 3109 dijalankan.

### 2. Tidak setara — dan sebabnya lebih tajam dari dugaan

`MEMORI_PEMAHAMAN.MD` §6.2 Langkah 7 menyatakan `TotalClaim = ClaimEstimation`. **Penetapan seperti itu tidak ada di XML**; yang ada lebih rumit.

#### 2.1 Di `CountLossAllocation_act`, keduanya memang menjadi sama

| Baris | Properti | Nilai |
|---|---|---|
| 9258 | `SpreadingRisk(<LAST>).ClaimEstimation` | `@if(Local.ClaimValue − Local.TotalAllocation < Local.LayerLimit, …)` |
| 9279 | `SpreadingRisk(<LAST>).TotalClaim` | `@if(Local.ClaimValue − Local.TotalAllocation < Local.LayerLimit, …)` |

Bentuk yang sama, hasil alokasi yang sama. Jadi saat `CountLossAllocation_act` selesai, **keduanya bernilai sama** — memori benar untuk titik itu.

#### 2.2 Tetapi `AdjClaimCNP_Act` menimpanya **sebelum** baris 3109 dibaca

Urutan langkah di dalam rule yang sama, dari `pyStepPageReference`:

| Baris | Langkah | Yang terjadi |
|---|---|---|
| 2348 | `RH_1.pySteps(8)` | `.SpreadingRisk(IdxLastLayer).TotalClaim` **ditimpa** |
| 2369 | `RH_1.pySteps(8)` | `.SpreadingRisk(IdxLastLayer).TotalClaim` **ditimpa lagi** — `@if(Local.TotalClaim==0, Local.TotalValue − …, …)` |
| **3109** | **`RH_1.pySteps(11).pySteps(1)`** | **rumus reinstatement membaca `.TotalClaim`** |
| 6441 | `RH_1.pySteps(15).pySteps(4).pySteps(1)` | `.TotalClaim = .ClaimSpreaded + .AdjusterFee + .CNPOthersFee − .Salvage` |
| 7834 | `RH_1.pySteps(15).pySteps(5).pySteps(5)` | idem |

**Langkah 8 menimpa `TotalClaim` sebelum langkah 11.1 membacanya.** Nilai yang masuk rumus reinstatement bukan hasil alokasi, melainkan hasil langkah 8 — yang bersumber dari `Local.TotalClaim`, yaitu akumulasi `.ClaimAmountAdjust` (baris 1561).

Dan penetapan yang memang memuat `AdjusterFee` serta `Salvage` — baris 6441 dan 7834 — berjalan **sesudah** 3109, jadi tidak terbaca olehnya.

> **Kedua rumus bekerja atas nilai yang berbeda.** Bukan karena salah satunya menghilangkan komponen dengan sengaja, melainkan karena `TotalClaim` ditimpa tiga kali di dalam satu rule, dan rumus reinstatement membacanya di antara dua penimpaan.

### 3. Yang menguatkan bahwa keduanya *dimaksudkan* sama

`Activity\SetActualPremium_ACT.xml` menyetelnya secara eksplisit:

| Baris | Penetapan |
|---|---|
| 702 | `SpreadingRisk(<LAST>).ClaimEstimation = SpreadingRisk(<LAST>).TotalClaim` |
| 951 | `.ClaimEstimation = .TotalClaim` |

Jadi di tempat lain kesamaan itu **ditegakkan lewat penetapan**, bukan berlaku dengan sendirinya. Itu memperkuat pembacaan memori tentang maksudnya, sekaligus memperlihatkan bahwa maksud itu tidak dijaga di semua jalur.

### 4. Batas klaim

1. **Urutan langkah disimpulkan dari `pyStepPageReference` dan urutan dokumen.** Untuk langkah tingkat atas di dalam satu activity, keduanya sejalan. Percabangan `pyStepsPreCondParamsWhen` dapat membuat sebagian langkah dilewati, dan itu belum ditelusuri satu per satu.
2. **Berapa besar selisih angkanya belum diukur** — itu pertanyaan data, bukan pertanyaan kode.
3. Rule mana yang benar secara bisnis bukan kesimpulan dokumen ini. Acuan yang dipilih adalah `CountReinstatement_Act`.

### 5. Pengamatan terpisah, tidak ditarik lebih jauh

`Activity\CopyOldataCurr_act.xml` baris 1862 dan `Activity\CountTotalInsterest_Act.xml` baris 4769 memuat:

```
.ClaimEstimation = .ClaimSpreaded * .PremiumSpreaded
```

Nilai klaim dikalikan nilai premi. Dicatat apa adanya; tidak ditelusuri lebih jauh.


---

# IV. Keputusan — apa yang sudah diputuskan untuk sistem baru

Dua puluh tujuh ADR. Dua puluh enam `accepted`, satu masih `proposed`.

## Aggregate root adalah Klaim, satu klaim satu kejadian kerugian

Sumber: `docs/adr/0001-satu-klaim-satu-kejadian.md`

*status: accepted · label: DECIDED*


Struktur XML mendukung tafsir ini tetapi tidak memaksakannya: seluruh properti tanggal dan penyebab kerugian di `.ClaimData` bersifat skalar, ada 18 PageList di bawahnya, dan sapuan atas 114 activity menunjukkan **tidak ada satu pun rule yang menulis `.DateOfLoss`** — nilainya hanya masuk lewat input Registrasi. Kami menetapkan satu Klaim = satu kejadian kerugian atas satu polis treaty non-proporsional, dengan Adjustment sebagai transaksi di bawahnya, sehingga Klaim menjadi aggregate root dan ke-18 PageList menjadi tabel anak.

### Consequences

Natural key Klaim belum ditetapkan. Tidak ada rule yang mencegah duplikat `PolicyNo` + `DateOfLoss` + `IDMaster`; `CheckDateDOL_Act` hanya menampilkan riwayat klaim atas polis yang sama sebagai informasi, tidak menolak. Sampai ada konfirmasi bahwa pencegahan duplikat dilakukan secara prosedural, kunci teknis tetap identitas klaim yang digenerasi sistem, dan kombinasi di atas diperlakukan sebagai indeks pendukung, bukan unique constraint.

## Klaim yang sudah ditutup tidak dapat dibuka kembali

Sumber: `docs/adr/0002-tanpa-reopen.md`

*status: accepted · label: DECIDED*


Di seluruh 279 berkas folder ini hanya ada satu status akhir, `Resolved-Completed`, dan tidak ditemukan satu pun rule yang mengembalikan klaim dari status akhir ke status terbuka. Kami menetapkan tidak ada mekanisme reopen: koreksi atas klaim yang sudah ditutup dilakukan melalui Adjustment baru, bukan dengan membuka ulang klaim.

### Consequences

Kesimpulan ini **disimpulkan dari ketiadaan rule, bukan dari adanya rule**. Ketiadaan jalur reopen di XML tidak membuktikan ketiadaan reopen di produksi — pembukaan ulang secara manual lewat database tidak meninggalkan jejak di rule. Verifikasinya tercatat sebagai EXTERNAL di `_selesai/OPEN-QUESTIONS.md`. Bila ternyata reopen manual memang dilakukan, keputusan ini harus ditinjau ulang sebelum model status dikunci, karena konsekuensinya menyentuh keterlacakan audit.

## Presisi tinggi sepanjang rantai perhitungan, pembulatan hanya di tepi

Sumber: `docs/adr/0003-presisi-dan-pembulatan.md`

*status: proposed · label: DECIDED (DRAFT — angka presisi menunggu pengetahuan/SCHEMA-ACTUAL.csv)*


Sistem lama tidak punya aturan pembulatan: dari 673 ekspresi aritmetika, **518 (77%) tidak menyatakan skala sama sekali** dan bergantung pada perilaku desimal bawaan Pega, sementara 155 sisanya memakai tujuh skala berbeda. Kami menetapkan aturan baru: nilai moneter dihitung dan disimpan pada satu presisi tinggi seragam sepanjang rantai perhitungan, dan pembulatan hanya dilakukan di dua tepi — saat ditampilkan ke pengguna, dan saat dikirim ke sistem hilir.

### Considered Options

Menyalin skala per ekspresi apa adanya ditolak: skalanya tidak konsisten bahkan di dalam satu rule (`CountClaimTNP_Act` memakai 10, 20, dan 5; `CountLossAllocation_act` memakai 10 dan 20 dalam porsi berimbang), sehingga menyalinnya berarti mengabadikan ketidaksengajaan.

### Consequences

Skala ternyata **bukan semata kebiasaan penulisan**, melainkan mengikuti batas modul: inti klaim non-proporsional memakai 20, sisi treaty/premi memakai 4, pemformatan tampilan memakai 0 dan 2 (selalu lewat "bagi dengan 1"), dan tarif pajak di `SetPPNPPH` memakai 8. Dugaan awal bahwa ini murni kebiasaan developer meleset dan sudah dikoreksi.

Angka presisi final **belum ditetapkan** dan ADR ini tidak boleh menyebut angka apa pun sebagai final sebelum hasil profil data Oracle diterima. Risiko yang harus diputuskan lebih dulu: bila profil menunjukkan kolom uang di produksi benar-benar menyimpan lebih dari dua desimal, maka aturan "bulatkan dua desimal saat kirim ke hilir" justru akan mengubah angka yang selama ini diterima akuntansi.

### Catatan 18 September 2026 — tetap `proposed`, dan alasannya berubah

Status **tidak** dinaikkan. Dua hal yang perlu tercatat:

1. **`NUMBER(20,4)` pada `TREATYINPRODUCTION` diperlakukan sebagai konvensi satu modul**, bukan standar perusahaan, sampai akuntansi menyatakan sebaliknya. Ia **tidak** disebarkan ke tabel lain dan tidak dipakai sebagai dasar keputusan presisi.
2. **Alasan ADR ini masih draft sudah berubah.** Semula: datanya belum ditarik. Sekarang: **basis data memang tidak menyimpan presisi untuk nilai klaim sama sekali** — nilainya tidak punya kolom, tersimpan sebagai teks di dalam JSON, dan bahkan view `CLAIMXOL` mengeluarkannya sebagai `varchar2`. Profil data tidak akan mengubah kenyataan itu.

Yang menghambat sekarang adalah keputusan kebijakan, bukan ketersediaan data — `ASK-AKUNTANSI.md` pertanyaan 1.

## Tambalan per-case tidak ikut dimigrasi; angkanya dipindahkan sebagai data

Sumber: `docs/adr/0004-tambalan-per-case-tidak-dimigrasi.md`

*status: accepted · label: DECIDED*


Sapuan seluruh folder menemukan **9 identitas klaim atau master treaty yang ditulis langsung di dalam kode kalkulasi**, tersebar di 4 rule dan 14 langkah — antara lain `CLMNP-975` yang menimpa Premi Pemulihan dengan angka mati `881928.966808370` (IDR) dan `806851161.1895` (USD), serta `CLMNP-232` yang memetakan Biaya Penilaian secara manual per indeks. Kami memutuskan tidak satu pun dari cabang ini dibawa ke sistem baru: kodenya dibuang, tetapi nilai akhirnya dipertahankan sebagai data sehingga saldo klaim yang bersangkutan tetap benar.

### Consequences

Daftar lengkap ke-9 identitas beserta nilai yang ditimpa dan rumus normal yang di-bypass ada di `BLUEPRINT.md` §7, dan menjadi instruksi bagi tim migrasi data — bukan bagi tim pembangun aplikasi.

Salah satu temuan memperkuat bahwa ini memang tambalan, bukan aturan: pada `InputOutStandingClmTNP_PreAct` langkah 5, deskripsi langkah berbunyi `pyWorkPage.pyID=="CLMNP-50"` sementara kondisi eksekusinya `pyWorkPage.pyID=="CLMNP-232"` — deskripsi dan kondisi tidak sinkron, ciri khas salin-tempel.

Pemeriksaan ulang terhadap `MEMORI_PEMAHAMAN.MD` tidak menemukan keterangan bahwa treaty `1000393` punya perlakuan bisnis khusus. Karena itu ia diperlakukan sebagai tambalan, bukan sebagai dimensi yang hilang dari model data. Bila pemilik proses kemudian menyatakan sebaliknya, keputusan ini harus ditinjau ulang dan model data perlu menambah atribut yang menjelaskan perbedaan perlakuan tersebut.

### Tambahan 18 September 2026 — persetujuan eksplisit dan mekanisme pengganti

Persetujuan diberikan tanpa menunggu status terbuka/tutup tiap klaim: **tidak satu pun dari 29 langkah tambalan dimigrasi.**

Penggantinya satu mekanisme tunggal: **koreksi bernilai tercatat** — `nilai_sebelum`, `nilai_sesudah`, `alasan`, `pelaku`, `waktu`. Satu bentuk untuk semua kasus, bukan cabang per klaim di dalam kode.

Status tiap klaim yang disebut di dalam tambalan (`CLMNP-232`, `CLMNP-975`, `CLMNP-861`, `CLMNP-50`, `CLMNP-367`, `CLMNP-382`, `IDMaster 1000393`, `1001130`) diukur lewat REQ-015, dan itu **bukan penghalang** keputusan ini.

Satu hal yang dapat disimpulkan tanpa menunggu query: `CLMNP-975` ditambal **2026-07-16**, dua bulan sebelum ekspor ini. Tambalan tidak ditulis untuk klaim yang sudah tutup, jadi klaim itu hampir pasti masih terbuka dan tambalannya masih aktif.

**Karena itu `CLMNP-975` dijadikan kasus uji utama shadow-run.** Bila sistem baru menghasilkan angka yang benar untuk klaim itu **tanpa** tambalan apa pun, tesis "ini tambalan data, bukan aturan bisnis" terbukti untuk seluruh kelompoknya sekaligus — dan 29 langkah itu gugur bersama-sama, bukan satu per satu.

## Claim dan Komite satu unit cutover, dengan shadow-run sebagai bukti paritas

Sumber: `docs/adr/0005-strategi-cutover.md`

*status: accepted · label: DECIDED*


Claim Non Prop dan Komite Claim Non Prop tidak dapat dipisahkan pada saat peralihan: keduanya menulis ke record yang sama dan saling mengunci (`pxAddChildWork` dari sisi Claim, penulisan balik ke `AdjustmentList` dari sisi Komite), sehingga memotong di antara keduanya berarti menciptakan penguncian lintas sistem atas baris Oracle yang sama. Kami menetapkan keduanya beralih sebagai satu unit; sebelum peralihan dijalankan shadow-run, yaitu sistem baru menghitung ulang klaim produksi yang sudah selesai dan hasilnya dibandingkan angka per angka dengan sistem lama; klaim yang sedang berjalan diselesaikan di sistem lama dan sistem baru hanya melayani klaim baru.

### Consequences

**Tidak ada satu pun bukti di XML yang mendukung maupun menolak strategi ini** — XML hanya membuktikan keterkopelannya, bukan cara memindahkannya. Ini murni keputusan kami.

Yang terbukti dari XML adalah titik koplingnya, dan itu tercatat lengkap di `BLUEPRINT.md` §8 sebagai Boundary Contract Draft: dua pemanggilan `pxAddChildWork`, kunci relasi `pxCoveredInsKeys`/`pxCoverInsKey`, dan `Adjustment.IndexObject` sebagai penunjuk posisi numerik ke `AdjustmentList` induk.

`Adjustment.IndexObject` adalah risiko tersendiri: ia menyimpan posisi dalam daftar, bukan kunci surrogate. Bila urutan `AdjustmentList` berubah, kaitan antara keputusan Komite dan Adjustment yang dimaksud menjadi salah. Sistem baru wajib mengganti ini dengan kunci yang stabil, dan migrasi data harus memetakan posisi lama ke kunci baru.

### Tambahan 18 September 2026 — prasyarat baseline shadow-run

Shadow-run mengandaikan baseline yang sehat. Dua hal dapat merusaknya, dan keduanya harus diselesaikan **sebelum** perbandingan dimulai:

1. **Baris ganda akibat hitung ulang yang tidak idempoten** (ADR-0011). Case terdampak dibersihkan atau dikeluarkan dari perbandingan, dan didaftar terpisah. Diukur lewat REQ-014.
2. **Riwayat persetujuan yang tidak dapat direkonstruksi** karena hanya tersimpan sebagai teks komentar (ADR-0009, FINDING-002 bagian 7). Case terdampak tidak dapat dibandingkan pada dimensi "siapa menyetujui". Diukur lewat REQ-013.

Kasus uji utama: **`CLMNP-975`**, lihat ADR-0004.

### Tambahan 18 September 2026 — komposisi sampel dan aturan kegagalan

#### Komposisi sampel ditetapkan di muka

Sampel shadow-run **ditetapkan sekarang, bukan dipilih saat menjalankan**. Perbandingan atas 500 klaim rupiah satu layer akan lulus dengan mudah dan tidak membuktikan apa pun.

Sampel wajib memuat, dengan proporsi minimum terhadap keseluruhan sampel:

| Kelompok | Minimum |
|---|---|
| Klaim multi-mata-uang (lebih dari satu `Currency` pada satu klaim) | 15% |
| Klaim multi-layer (lebih dari satu baris `SpreadingRisk` selain `"UR"`) | 20% |
| Klaim dengan reinstatement (`CNPReinstatement` bukan nol) | 10% |
| Klaim dengan lebih dari satu Adjustment | 15% |
| Klaim yang melibatkan baris `"UR"` | 20% |
| Kedelapan klaim bertambalan | **100% — seluruhnya wajib ikut** |

Kelompok boleh bertumpang tindih; yang tidak boleh adalah ada kelompok yang kosong. Bila populasi produksi tidak menyediakan cukup kasus untuk satu kelompok, itu sendiri temuan dan dicatat, bukan diam-diam diturunkan ambangnya.

#### Aturan kegagalan

| Jenis nilai | Bila selisih melampaui ambang |
|---|---|
| **Nilai yang masuk jurnal akuntansi** | **Cutover berhenti.** Tanpa negosiasi, tanpa daftar pengecualian. |
| **Nilai antara** (alokasi per layer sebelum pembulatan akhir) | Boleh masuk daftar pengecualian, tetapi **setiap baris wajib punya penjelasan dan persetujuan bernama**. Pengecualian tanpa penjelasan tidak dihitung sebagai pengecualian — ia dihitung sebagai kegagalan. |

Aturan ini ditetapkan **sebelum** ada pihak yang punya kepentingan atas hasilnya. Itu alasan ia ditulis sekarang dan bukan nanti.

Angka toleransi relatif untuk nilai antara masih menunggu akuntansi (lihat ADR-0003 yang berstatus draft).

#### Prasyarat baseline

Selain dua hal yang sudah dicatat di atas, bertambah satu yang baru terbukti dari DDL: `PEGA_JSON_OS_AKSEP_KLAIMTNP` **selalu `INSERT` dan tidak pernah `UPDATE`** — logika upsert-nya ada tetapi dikomentari seluruhnya. Bersama ketiadaan primary key pada `OS_AKSEPTASI_KLAIM`, baris ganda pada tabel itu mungkin terjadi dan harus diukur sebelum perbandingan dimulai.

### Tambahan — nilai IDR lama tidak dipakai sebagai pembanding

Konversi mata uang di sistem lama **tidak dapat direproduksi**: `POOLDATA.GETCURRENCYSTANDARD` mengabaikan parameter tanggalnya dan selalu memakai kurs terbaru sampai `sysdate` (`FINDING-001` bagian 7.1). Angka IDR yang dihitung tahun lalu tidak dapat dihasilkan ulang hari ini.

Karena itu, untuk nilai hasil konversi, **shadow-run membandingkan nilai mata uang asli, bukan nilai IDR-nya.** Nilai IDR lama diterima apa adanya sebagai nilai historis dan tidak dijadikan pembanding.

Pengecualian ini berlaku hanya untuk kelompok nilai hasil konversi. Nilai yang sejak awal berdenominasi rupiah tetap dibandingkan seperti biasa, dengan aturan dua ambang di atas.

## RBAC sistem baru dirancang dari nol, bukan diwarisi

Sumber: `docs/adr/0006-rbac-dirancang-dari-nol.md`

*status: accepted · label: DECIDED*


Modul ini tidak memiliki otorisasi di tingkat rule: `pyPrivilegeName` kosong di seluruh 279 berkas, tidak ada satu pun `Rule-Access-When`, dan seluruh elemen UI bervisibilitas `ALWAYS` — termasuk kedua jalur penutupan klaim. Kontrol akses yang benar-benar berjalan bertumpu pada access group dan routing flow yang tidak ikut ter-export. Kami menetapkan model otorisasi sistem baru dirancang dari nol bersama pemilik proses, dan tidak ada satu pun aturan akses yang boleh diklaim sebagai warisan sistem lama.

### Consequences

Ini **menambah cakupan pekerjaan**, bukan sekadar catatan. Perancangan RBAC menjadi pekerjaan tersendiri dengan pemilik keputusan di sisi bisnis, bukan turunan otomatis dari analisis XML.

Empat label peran memang terbukti ada di `DataTransform\InsertChronology_DT.xml` — `Claim Admin`, `Claim Dept. Head`, `Operational Director`, `Technical Director` — tetapi itu label yang dicatatkan ke jejak audit, bukan aturan yang menegakkan siapa boleh melakukan apa. Pemetaannya ke orang pun bersifat hardcode per nama operator, sehingga tidak dapat dijadikan dasar model peran.

Sampai RBAC baru disepakati, setiap pernyataan tentang "siapa boleh melakukan apa" dalam dokumen mana pun berstatus hipotesis dan wajib berlabel EXTERNAL.

## Setiap nilai uang disimpan berpasangan, dan ambang kewenangan dibandingkan terhadap nilai IDR

Sumber: `docs/adr/0007-nilai-uang-berpasangan-dan-ambang-idr.md`

*status: accepted · label: DECIDED*


Klaim non-proporsional berjalan dalam banyak mata uang sekaligus: `ListClaimAmount`, `SpreadingRisk`, dan `ListTotalEstimation` semuanya ber-kunci `Currency`, sementara properti bersufiks `IDR` menunjukkan IDR dipakai sebagai mata uang penyetaraan. Kami menetapkan setiap nilai uang disimpan sebagai satu paket — **nilai asli, mata uang, kurs, tanggal/sumber kurs, dan nilai IDR hasil konversi** — dan setiap perbandingan terhadap ambang kewenangan Komite dilakukan terhadap nilai IDR hasil konversi tersebut.

### Consequences

**Tanggal dan sumber kurs adalah bagian wajib dari paket**, bukan pelengkap. Tanpa keduanya, angka IDR tidak dapat diaudit ulang: nilai yang sama bisa dibenarkan atau disalahkan tergantung kurs kapan yang dipakai, dan rekonsiliasi dengan akuntansi menjadi tidak dapat diselesaikan. Sumber kurs yang terlihat di sistem lama adalah fungsi `getcurrencystandard` dan `TreatyInMaster.CurrencyList.Conversion`; keduanya perlu dicatat identitasnya, bukan hanya hasilnya.

Keputusan tentang ambang ini **tidak menunggu** hasil investigasi kondisi sistem lama yang tercatat di `FINDING-001-threshold-currency.md`. Apa pun temuan di sana, sistem baru membandingkan terhadap IDR.

Presisi dan skala setiap kolom dalam paket ini belum ditetapkan dan menunggu `pengetahuan/SCHEMA-ACTUAL.csv` — lihat ADR-0003 yang masih berstatus draft.

## Nilai hasil suntingan manual bertahan terhadap hitung ulang, dan selalu terlihat

Sumber: `docs/adr/0008-suntingan-manual-bertahan-dan-terlihat.md`

*status: accepted · label: DECIDED*


Ketika petugas mengganti nilai hasil perhitungan dengan nilai suntingan sendiri, nilai itu **dikunci** dan tidak ditimpa oleh perhitungan berikutnya. Layar menampilkan keduanya berdampingan — nilai hitungan yang digantikan dan nilai suntingan yang berlaku — beserta siapa yang menyunting dan kapan.

Alternatif yang ditolak: mengikuti perilaku sistem lama, yaitu suntingan selalu kalah terhadap hitung ulang.

### Consequences

Model data memerlukan konsep **nilai terkunci** sejak awal: setiap nilai yang dapat disunting menyimpan pasangan `nilai_hitungan` dan `nilai_disunting`, plus `pelaku` dan `waktu`. Menambahkan ini setelah tabel terbentuk berarti membongkar setiap baris yang sudah ada, karena nilai lama tidak dapat dipisahkan kembali menjadi dua asal yang berbeda.

Kombinasi terburuk yang dihindari keputusan ini adalah suntingan diam-diam yang dapat hilang diam-diam — persis keadaan yang terbaca di sistem lama dan dilaporkan di `FINDING-004-suntingan-manual-tertimpa.md`.

Keputusan ini **tidak menunggu** kepastian hipotesis F1/F2 tentang `.IsEditClaim`. Apa pun yang ditemukan di modul Komite nanti, sistem baru mengunci suntingan.

Konsekuensi migrasi data bergantung pada hipotesis mana yang benar, dan itu diuraikan di FINDING-004 bagian 4, bukan di sini.

## Keputusan alur disimpan sebagai field terstruktur; komentar tidak pernah dibaca mesin

Sumber: `docs/adr/0009-keputusan-terstruktur-bukan-teks-bebas.md`

*status: accepted · label: DECIDED*


Setiap keputusan yang menggerakkan alur disimpan sebagai tiga field terpisah — **tindakan**, **pelaku**, **waktu** — bukan disimpulkan dari isi kolom komentar. Kolom komentar tetap ada sebagai catatan manusia, dan **tidak pernah** menjadi masukan bagi kondisi mana pun.

Sistem lama melakukan sebaliknya: `SethistoryKlaimTreaty` mengambil keputusan lewat `@contains(.CommentSuggest,"Accepted by <nama>")` pada enam pasang kondisi. Uraiannya di `FINDING-002-percabangan-identitas.md` bagian 6.

### Consequences

Tiga cacat sekaligus hilang, dan ketiganya nyata di sistem lama: satu huruf besar berbeda membuat kondisi gagal tanpa pesan galat; siapa pun yang boleh mengisi komentar dapat memenuhi kondisi alur; dan orang yang sama ditulis dalam dua ejaan sehingga sebagian cabang tidak pernah terpicu.

**Harga yang dibayar ada di migrasi data, bukan di rancangan.** Riwayat lama menyimpan keputusannya hanya di dalam teks, dan sebagian teks itu sudah ditimpa oleh langkah penerjemahan di sistem lama sendiri. Kolom "siapa menyetujui" pada sistem baru tidak dapat diisi lengkap dari data lama secara otomatis; sebagian akan kosong. Besarnya bagian yang kosong belum diketahui dan diukur lewat REQ-013.

Konsekuensi itu diterima secara sadar di muka. Ia bukan temuan yang muncul saat UAT.

## Layer dan Retensi Cedant adalah dua entitas terpisah

Sumber: `docs/adr/0010-layer-dan-retensi-cedant-dua-tabel.md`

*status: accepted · label: DECIDED*


Susunan lapisan XOL dan Retensi Cedant disimpan di **dua tabel berbeda**. Sistem lama menyimpan keduanya di satu PageList `SpreadingRisk`, dengan Retensi Cedant sebagai baris ber-`TreatyName="UR"` yang di-`APPEND` oleh `CountLossAllocation_act`.

### Consequences

Alasan yang menentukan bukan soal penyaring, melainkan **perilaku baris itu sendiri**: `ClaimPercentage` selalu `0` dan `ClaimSpreaded` selalu `0`. Porsi reasuradur untuk baris itu tidak pernah ada. Sesuatu yang selalu nol di dalam daftar alokasi bukan anggota daftar itu — ia konteks bagi daftar tersebut.

Argumen bahwa "kalau ia lapisan nol, tidak akan ada rule yang perlu menyaringnya" **tidak dipakai** sebagai dasar, karena penyaring bisa ada karena sebab lain.

Setelah dipisah, pertanyaan "Retensi Cedant ikut dihitung atau tidak" menjadi eksplisit pada setiap query — bukan bergantung pada apakah penulisnya ingat menuliskan penyaring. Di sistem lama, sembilan rule menyentuh `SpreadingRisk` antara 16 dan 66 kali tanpa pernah menyebut `"UR"`.

Ini mengunci model data inti dan sangat mahal diubah kemudian: setiap query alokasi, setiap agregasi, dan setiap laporan ikut berubah bentuknya.

### Tambahan 18 September 2026 — satu kolom yang selalu nol ikut masuk ke tabel ini

Baris Retensi Cedant menerima `TotalClaim = Local.TotalUR` (`Activity\CountLossAllocation_act.xml` baris 7015–7120).

Dan `Local.TotalUR` pada cabang `Local.Currency == .Currency` dihitung sebagai `(.Deductible * 0 * Local.ProrateClaim/100)` — **dikalikan nol**, sehingga selalu nol. Penetapan `Local.UR` pada cabang yang sama, di berkas yang sama, tidak memakai `* 0`.

Artinya: untuk seluruh klaim bermata uang sama, tabel Retensi Cedant yang baru akan memuat kolom `TotalClaim` yang **selalu bernilai nol** — bila perilaku lama dipertahankan.

Apakah perilaku itu dipertahankan atau diperbaiki **belum diputuskan**, dan bukan keputusan teknis: ia dibawa ke akuntansi sebagai `ASK-AKUNTANSI.md` pertanyaan **2b**. Keputusan ADR ini tentang **bentuk tabel** tetap berlaku apa pun jawabannya.

## Menjalankan ulang perhitungan alokasi harus menghasilkan keadaan yang identik

Sumber: `docs/adr/0011-hitung-ulang-idempoten.md`

*status: accepted · label: DECIDED*


Perhitungan alokasi bersifat **idempoten**: menjalankannya dua kali atas masukan yang sama menghasilkan keadaan yang sama persis. Ini diuji otomatis, bukan sekadar dijanjikan.

### Consequences

Alasannya tidak punya tandingan: bila menekan tombol hitung dua kali menghasilkan angka berbeda, tidak ada satu pun angka di sistem itu yang bisa dipertanggungjawabkan.

**Konsekuensi terhadap shadow-run, dan ini yang paling mahal.** Bila sistem lama ternyata tidak idempoten, ada data produksi dengan baris ganda. Itu bukan sekadar urusan migrasi data — ia **merusak dasar perbandingan**. Hitungan sistem baru tidak dapat dibandingkan terhadap baris yang sudah terlanjur ganda, karena selisihnya akan terbaca sebagai kesalahan sistem baru padahal berasal dari kerusakan baseline.

Karena itu, sebelum shadow-run dimulai: baseline dibersihkan lebih dulu, **atau** case yang terdampak dikeluarkan dari perbandingan dan didaftar terpisah. Salah satu harus dipilih; tidak boleh dibiarkan tercampur. Lihat juga ADR-0005.

Yang membuat risiko ini nyata di sistem lama: `CountLossAllocation_act` melakukan `<APPEND>` baris `"UR"`, dan **tidak ditemukan di XML** penghapusan PageList `SpreadingRisk` sebelum `APPEND` itu. Satu-satunya penghapusan yang menyasar `SpreadingRisk` di seluruh folder adalah `Property-Remove` atas properti tunggal `.SpreadingRisk(Param.idx).AdjClaimValue` di dua rule — bukan atas daftarnya. Selain itu `CountLossAllocation_act` langkah 10 memuat `Property-Remove` yang **parameter sasarannya kosong** di ekspor ini.

Besarnya kerusakan diukur, bukan ditebak — REQ-014 menghitung berapa case yang memiliki lebih dari satu baris `"UR"`. Keputusan idempoten ini **tidak menunggu** hasil pengukuran itu.

## Tambalan baru dideteksi lewat sapuan ulang, bukan lewat register

Sumber: `docs/adr/0012-deteksi-tambalan-baru-lewat-sapuan-bukan-register.md`

*status: accepted · label: DECIDED*


Tambalan per-case dan per-identitas dibekukan sejak tanggal kesepakatan, dengan pengecualian untuk tambalan yang menahan pembayaran atau menghentikan operasi. Yang menjadi **alat deteksi** bukan register, melainkan **sapuan ulang otomatis** atas setiap ekspor XML baru, dibandingkan terhadap `BLUEPRINT.md` bagian 7.5. Register pengecualian tetap dibuat, tetapi statusnya pelengkap.

### Consequences

Pembekuan total akan dilanggar diam-diam; yang tercatat lebih berguna daripada yang dilarang. Tetapi register yang diisi manusia punya cacat yang sama dengan `CNPStatusCase` — bisa tidak pernah diisi, dan ketiadaannya tidak menghasilkan sinyal apa pun. **Kebasian analisis tidak boleh bergantung pada kepatuhan orang.**

Sapuan ulang mendeteksi tambalan baru apakah registernya diisi atau tidak. Selisih antara hasil sapuan dan bagian 7.5 adalah daftar tambalan yang masuk sesudah ekspor terakhir.

Konsekuensi kedua: **setiap artefak membawa stempel asal** — tanggal ekspor XML yang menjadi dasarnya dan jumlah berkas yang disapu. Tanpa itu, tidak ada yang dapat menilai dokumen ini berlaku untuk keadaan kapan. Stempel sudah dipasang di seluruh artefak per 18 September 2026 (279 berkas, ekspor 2026-09-08/09).

Siapa yang berwenang menyetujui pengecualian: **EXTERNAL**, dibawa ke pemilik proses, dicatat di `_selesai/OPEN-QUESTIONS.md` A13. Tidak ditanyakan lagi dalam grilling.

### Tambahan 18 September 2026 — sapuan mencakup tiga lapisan, bukan satu

Versi pertama menyebut "sapuan ulang atas setiap ekspor XML baru". **Itu terlalu sempit.** Tiga jalur yang menggerakkan sistem tidak berjejak di XML sama sekali (lihat `BLUEPRINT.md` bagian 14), termasuk tanggal tertanam di dalam `PROC_GENERATE_SEQUENCE_NUMBER` yang merupakan kelas tambalan yang sama dengan 29 langkah di bagian 7.

Sapuan berkala mencakup:

| Lapisan | Yang dicari |
|---|---|
| **Ekspor XML** | literal `CLMNP-…`, `IDMaster`, identitas operator/orang, pola `@contains` atas teks bebas |
| **DDL** | tanggal dan periode tertanam, nilai mati di dalam procedure, kolom baru, constraint yang hilang |
| **Daftar database link** | ketergantungan lintas basis data yang baru muncul |

Lapisan ketiga ditambahkan setelah `V_MST_USER_TEKNIS` terbaca. Sebelum itu tidak ada yang tahu perlu menyapunya — yang justru menjadi alasan terkuat mengapa daftar lapisan ini sendiri harus ditinjau setiap kali sumber baru masuk.

## Perhitungan dijalankan sebagai turunan dari data, bukan sebagai akibat penekanan tombol

Sumber: `docs/adr/0013-perhitungan-sebagai-turunan-data.md`

*status: accepted · label: DECIDED*


Setiap nilai hasil perhitungan adalah **turunan** dari masukannya. Begitu masukan berubah, turunannya ikut berubah. Tidak ada tombol "hitung", tidak ada keadaan setengah jadi, dan tidak ada urutan tindakan yang harus diingat pengguna.

### Consequences

Sistem lama melakukan sebaliknya, dan itu terbukti dari XML: **tidak ada satu pun activity yang memanggil `CountClaimTNP_Act`** — ia dipicu dari Section lewat `<pyActivity>`, 42 rujukan di tiga berkas. Urutan jalannya `CountClaimTNP_Act`, `CountLossAllocation_act`, `GenerateCFS_act`, dan `SaveToOS` tidak ditetapkan di kode mana pun; ia ditentukan kontrol mana yang ditekan.

**Sapuan `MEMORI_PEMAHAMAN.MD` atas urutan kerja, SOP, dan langkah petugas menghasilkan nihil.** Urutan baku itu **tidak terdokumentasi di sumber mana pun** — tidak di XML, tidak di memori. Itu sendiri temuan: urutan yang menentukan hasil perhitungan hanya hidup di kepala petugas, dan akan hilang bersama orangnya.

Keputusan ini menutup jalur itu, bukan memperbaikinya. Konsekuensinya pada rancangan: perhitungan tidak boleh ditempatkan di lapisan tampilan; ia menjadi fungsi murni atas data, dipanggil ulang kapan pun masukannya berubah, dan hasilnya tidak disimpan sebagai keadaan yang dapat menyimpang dari masukannya.

**Keputusan ini mengikat bersama ADR-0011 dan ADR-0012, bukan berdiri sendiri.** Uraian rantai sebab-akibatnya ada di `BLUEPRINT.md` bagian 12: karena hasil bergantung urutan klik, selalu ada klaim yang keluar jalur; karena selalu ada yang keluar jalur, perbaikan termurah adalah menambal klaim itu satu per satu; dan begitulah 29 langkah tambalan menumpuk selama delapan tahun. Mengambil satu dari tiga keputusan tanpa dua lainnya akan mengembalikan pola yang sama dalam bentuk baru.

## Aturan penguraian teks ke angka, dan perlakuan atas kurs yang tidak ada

Sumber: `docs/adr/0014-penguraian-teks-ke-angka-dan-kurs-kosong.md`

*status: accepted · label: DECIDED*


Nilai uang di sistem lama tersimpan sebagai teks di dalam `DATA_JSON`, dan bahkan view `CLAIMXOL` mengeluarkannya sebagai `varchar2`. Aturan penguraiannya ditetapkan **tertulis dan di muka**, sebelum uji coba migrasi mana pun dijalankan.

### Aturan penguraian

| Hal | Perlakuan |
|---|---|
| Pemisah ribuan | ditolak bila ada; angka yang sah tidak memakainya |
| Tanda desimal | titik saja; koma sebagai desimal ditolak dan masuk daftar |
| Spasi di awal/akhir | dipangkas, tidak dianggap galat |
| Notasi ilmiah (`1.2E5`) | ditolak dan masuk daftar |
| Posisi tanda minus | hanya di depan; `123-` ditolak |
| Presisi melebihi ADR-0003 | ditolak dan masuk daftar, **tidak dibulatkan** |

### Tiga nilai kosong yang tidak boleh disamakan

Ini yang paling mudah terlewat, dan kerusakannya tidak dapat dipulihkan setelah migrasi:

| Bentuk | Arti yang dipegang | Perlakuan |
|---|---|---|
| `""` (string kosong) | **tidak pernah diisi** | dipetakan ke "belum diisi", bukan nol |
| `null` / kunci tidak ada di JSON | **tidak ada** | dipetakan ke null |
| `"0"` | **diisi nol** | dipetakan ke nol |

Menyamakan ketiganya menghancurkan informasi yang tidak dapat dipulihkan. Apakah ketiga makna itu benar-benar dibedakan oleh sistem lama **belum terbukti dari XML** — bila ternyata tidak dapat dibedakan maknanya, itu pertanyaan tersendiri yang dicatat di `_selesai/OPEN-QUESTIONS.md`, bukan diputuskan sepihak di sini.

### Kurs yang tidak ada

Konversi tanpa kurs **tidak menghasilkan angka**. Nilai IDR dibiarkan kosong dan ditandai "menunggu kurs". Sistem lama memilih `RETURN 1` — lihat `FINDING-006` — dan itu tidak diwarisi.

### Consequences

Ambang penghentian migrasi **berbasis nilai, bukan cacah baris**. Satu baris bermasalah senilai lima miliar lebih berat daripada lima ratus baris senilai seratus ribu.

- Setiap baris yang tidak dapat diurai diselesaikan satu per satu, berapa pun jumlahnya. Tidak ada baris yang dibuang karena "cuma sedikit".
- Migrasi dihentikan bila **nilai total yang belum terselesaikan** melewati ambang yang disepakati akuntansi, atau bila **ada satu baris** yang nilainya melewati ambang material.

Kedua angka ambang itu belum ditetapkan dan menunggu akuntansi — `ASK-AKUNTANSI.md` pertanyaan 4.

### Tambahan 18 September 2026 — aturan diturunkan dari kotoran yang benar-benar ada

Daftar aturan di atas semula disusun dari kemungkinan umum. Sapuan seluruh 279 berkas menggantinya dengan yang terbukti.

| | Jumlah |
|---|---|
| Seluruh pemanggilan `toDecimal` | **127** |
| Didahului pembersihan `replaceAll` | **9** — seluruhnya atas `.Deductible2` |
| **Tanpa pembersihan apa pun** | **118** |

**Koma sebagai pemisah desimal bukan hipotesis.** `@toDecimal(@replaceAll(.Deductible2,",","."))` di `Activity\CountLossAllocation_act.xml` membuktikan seseorang pernah menemuinya di produksi dan menambalnya. Tidak ada yang menulis pembersihan untuk masalah yang tidak pernah terjadi.

**Pembersihan lain yang ditemukan, dan sasarannya bukan angka**:

| Ekspresi | Sasaran |
|---|---|
| `@replaceAll(pyWorkPage.ClaimData.InsuredName,",","")` | menghapus koma dari **nama**, bukan dari angka — diduga demi keamanan muatan |
| `@replaceAll(.AcceptedNo,".","")` | menghapus titik dari **nomor dokumen** |
| `@pxReplaceAllViaRegex(.NoAccount,"[^0-9]","")` | menyisakan hanya digit pada **nomor rekening** |

Ketiganya memperlihatkan pola yang sama: pembersihan dilakukan **per tempat, saat masalahnya muncul**, bukan sebagai aturan yang berlaku menyeluruh.

**Yang paling perlu diperhatikan**: dua dari 118 pemanggilan tanpa pembersihan adalah `@toDecimal(TempKasir.CARI20)` dan `@toDecimal(TempKasir.CARI21)` — muatan menuju sistem kasir, yaitu nilai yang masuk jurnal menurut `ASK-AKUNTANSI.md` bagian 0.2.

Aturan penguraian di sistem baru karena itu **berlaku di satu tempat untuk semua nilai**, bukan ditambahkan per ekspresi saat kegagalannya ditemukan.

### Tambahan — keadaan "menunggu kurs" menular

Nilai yang belum terkonversi menandai **seluruh nilai turunannya**. Apa pun yang dihitung dari nilai bertanda "menunggu kurs" ikut bertanda sama.

Nilai bertanda itu **tidak boleh** masuk penjumlahan, tidak boleh dibandingkan terhadap ambang persetujuan, dan tidak boleh muncul di laporan — baik sebagai nol maupun sebagai nilai apa adanya.

Tanpa aturan penularan ini, keputusan "gagalkan, jangan mengarang angka" hanya berlaku satu tingkat: nilai pertama ditandai, lalu penjumlahan di tingkat berikutnya memperlakukannya sebagai nol dan cacatnya kembali muncul dalam bentuk yang lebih sulit dilihat — persis pola `RETURN 1` yang sedang kita tinggalkan (`FINDING-006`).

Ini melengkapi `ADR-0019` lapis 2: keadaan di tingkat baris mencatat *bahwa* perhitungan gagal; aturan penularan ini menentukan *apa yang terjadi pada baris-baris sesudahnya*.

## Akseptasi adalah satu entitas dengan keadaan, bukan dua tabel

Sumber: `docs/adr/0015-satu-entitas-akseptasi-dengan-keadaan.md`

*status: accepted · label: DECIDED*


Akseptasi biasa dan akseptasi bersyarat (Subjectivity) adalah **satu entitas** dengan keadaan yang berbeda, bukan dua entitas. Keadaannya: biasa, bersyarat, dan bersyarat-gugur — beserta daftar syarat dan status pemenuhannya.

Sistem lama memisahkannya jadi dua tabel: `OS_AKSEPTASI_KLAIM` dan `OS_AKSEPTASI_SUBJECTIVITY`.

### Consequences

Daur hidup yang berbeda adalah **keadaan**, bukan identitas. Dua hal menjadi entitas berbeda bila identitasnya berbeda; bila satu akseptasi bersyarat yang syaratnya terpenuhi berubah menjadi akseptasi biasa **yang sama**, maka sejak awal ia satu benda yang sedang berada di keadaan tertentu.

Bukti yang menguatkan datang dari sistem lama sendiri: view `CLAIMXOL` harus meng-`UNION ALL` kedua tabel untuk mendapat gambaran utuh. Sesuatu yang harus selalu digabungkan kembali tidak seharusnya dipisah.

Syarat yang gugur menjadi **transisi keadaan yang tercatat**, bukan penghapusan atau perpindahan tabel.

**Temuan yang menyertainya dan belum terjawab**: di folder Claim, `.IsSubjectivity` hanya **dibaca** (`==true`), disalin ke `Local.Subjectivity`, dan **tidak ada satu pun rule yang menggugurkan, membatalkan, atau mengakhiri akseptasi bersyarat yang syaratnya tidak pernah terpenuhi.** Delapan belas kemunculan `Subjectivity` di 279 berkas, tidak satu pun berupa transisi keadaan. Jadi dugaan "menggantung selamanya" konsisten dengan isi folder ini — *disimpulkan dari ketiadaan rule, bukan dari adanya rule*. Pemastiannya `DEFERRED-TO-KOMITE-SESSION`.

## Data pegawai masuk lewat API atau salinan tersinkron, bukan database link

Sumber: `docs/adr/0016-tanpa-database-link-ke-sistem-lain.md`

*status: accepted · label: DECIDED*


Sistem baru **tidak mewarisi ketergantungan lintas basis data**. Nama dan status kepegawaian diambil dari HRD lewat API bila tersedia, atau lewat salinan lokal yang disinkronkan berkala **dengan stempel waktu sinkronisasi yang terlihat pengguna**. Jabatan-dalam-konteks-klaim disimpan sendiri, terpisah dari jabatan HRD (ADR-0024 tentang pemisahan pengguna, jabatan, dan keanggotaan komite).

Sistem lama memakai `hrdasm.v_hrd_mst@asmd.sinarmas.co.id` — database link langsung ke basis data sistem HRD, di dalam view `V_MST_USER_TEKNIS`.

### Consequences

Tiga alasan, tidak satu pun soal selera:

1. Database link **melewati kontrol akses aplikasi HRD sepenuhnya**.
2. Ia **mengikat kita pada struktur tabel internal** sistem lain. Bila mereka mengubah kolom, sistem kita rusak tanpa ada yang memberi tahu.
3. Bila tautan putus, view mengembalikan **kosong** — daftar pengguna menghilang tanpa pesan galat. Itu pola yang sama persis dengan `RETURN 1` pada kurs (`FINDING-006` bagian 3): kegagalan yang menyamar menjadi hasil yang sah.

Bila HRD tidak dapat menyediakan API maupun mekanisme sinkronisasi, itu percakapan yang dibawa ke luar. Rancangannya **tidak menunggu** jawaban itu, dan tidak dibuat seolah database link adalah salah satu pilihan.

Apakah tautan yang ada sekarang resmi disepakati: EXTERNAL, dicatat di `_selesai/OPEN-QUESTIONS.md`.

## Data aplikasi hanya diubah lewat aplikasi

Sumber: `docs/adr/0017-satu-pintu-tulis.md`

*status: accepted · label: DECIDED*


Tidak ada jalur tulis ke data aplikasi selain lewat aplikasi itu sendiri. Proses batch yang perlu mengubah data memanggil layanan yang sama dengan yang dipakai layar, sehingga aturan dan jejaknya berlaku sama.

Di sistem lama, schema `POOLDATA` memiliki akses ke tabel work Pega `DATAPEGA.PC_ASM_FW_GCNMFW_WORK` — artinya data case dapat diubah oleh prosedur basis data, tanpa melewati aplikasi dan tanpa jejak di dalamnya.

### Consequences

Alasannya berdasar rekam jejak analisis ini sendiri. Sudah ditemukan **tiga** mekanisme yang hasilnya tidak dapat dilacak ke pelakunya:

| Mekanisme | Dokumen |
|---|---|
| Penggantian identitas `DARTO` menjadi `CHRISTINEANGELINA` | `FINDING-002` bagian 9 |
| Keputusan alur diambil dari teks komentar bebas | `FINDING-002` bagian 6 |
| Suntingan manual tertimpa hitung ulang tanpa peringatan | `FINDING-004` |

Pintu tulis kedua akan menambah yang keempat.

**Keputusan desainnya final; ukurannya belum, dan itu risiko lingkup.** Bila REQ-021 menunjukkan ada proses produksi yang menulis lewat jalur itu, proses-proses tersebut harus ditulis ulang **sebelum** cutover — pekerjaan yang belum masuk perkiraan mana pun.

Karena itu **REQ-021 ditandai BLOCKER perkiraan biaya, bukan sekadar BLOCKER teknis.** Bila hasilnya besar, ia dibawa ke manajemen sebagai pokok tersendiri, bukan diselipkan ke dalam laporan teknis.

## Klaim yang terdampak penjaga tanggal yang salah dimigrasi apa adanya, dan didaftar

Sumber: `docs/adr/0018-klaim-terdampak-penjaga-tanggal-dimigrasi-apa-adanya.md`

*status: accepted · label: DECIDED*


Klaim yang hasilnya berbeda bila penjaga `CheckDateDOL_Act` dijalankan dengan perbandingan yang benar **dimigrasi apa adanya**. Tidak ada koreksi otomatis, tidak ada perubahan status.

Sebagai gantinya: penjaga yang benar dijalankan atas seluruh data lama, hasilnya menjadi daftar, dan keputusan per klaim diambil manusia.

### Consequences

Klaim yang sudah dibayar tidak boleh berubah statusnya karena sistem baru menghitung ulang sebuah penjaga — itu bukan koreksi, itu mengubah angka yang sudah diterima akuntansi. Membiarkannya tanpa daftar berarti masalahnya ikut pindah tanpa ada yang tahu.

Tiga syarat yang mengikat:

1. **Daftar dihasilkan sebelum cutover.** Ia prasyarat, bukan laporan pasca-migrasi.
2. **Dua daftar terpisah, tidak digabung.** Yang seharusnya tertolak tetapi lolos; dan yang seharusnya lolos tetapi tertolak. Kelompok kedua lebih sensitif — itu klaim yang mungkin ditolak secara keliru.
3. **Keputusan per klaim dicatat.** Tidak ada koreksi massal.

**Keterjangkauan data — diperiksa dan hasilnya baik.** Sempat dikhawatirkan daftar ini tidak dapat dibuat karena `.EndDateTreaty` tersimpan IN-BLOB. Pemeriksaan menunjukkan sisi lainnya terjangkau: `.ClaimData.EndDateTreaty` diisi dari `pyWorkPage.TreatyInMaster.Termination`, dan **`POOLDATA.TREATYINDETAIL.TERMINATION` adalah kolom bertipe `DATE`**. `POOLDATA.TREATYCONTRACT.TREATYENDDATE` juga `DATE`. Sisi kiri perbandingan, `DATEOFLOSS`, adalah kolom `VARCHAR2(8)` di tabel work.

Jadi **REQ-020 dapat dijalankan dengan SQL biasa** — tanpa membongkar blob — dengan `TO_DATE(DATEOFLOSS,'YYYYMMDD')` dibandingkan terhadap kolom `DATE` itu. Yang masih perlu dipastikan hanyalah kunci penghubungnya; `MASTERID` termasuk kolom yang di-expose di tabel work, tetapi pasangannya di `TREATYINDETAIL` belum diverifikasi.

## NULL bukan nol, dan keadaan perhitungan disimpan di tingkat baris

Sumber: `docs/adr/0019-null-bukan-nol-dan-keadaan-di-tingkat-baris.md`

*status: accepted · label: DECIDED*


Tiga lapis, disusun dari yang termurah. Lapis berikutnya hanya ditambahkan di tempat lapis sebelumnya tidak cukup.

### Lapis 1 — `NULL` berarti belum dihitung, nol berarti dihitung dan hasilnya nol

Keduanya **tidak boleh pernah disamakan**, di mana pun, oleh kode mana pun. Tidak ada `NVL(x,0)` yang diam-diam menghapus perbedaan itu, tidak ada kolom nilai ber-`DEFAULT 0`, dan tidak ada layar yang menampilkan kosong sebagai `0`.

Biayanya nol — hanya disiplin — dan ia menangkap sebagian besar manfaat dari seluruh keputusan ini.

### Lapis 2 — keadaan di tingkat baris, bukan tingkat kolom

Setiap baris hasil perhitungan membawa tiga kolom: `dihitung_pada`, `status_perhitungan`, dan `versi_aturan` yang dipakai. Satu perhitungan gagal maka statusnya terbaca di baris itu, dan kolom yang gagal bernilai `NULL`.

Tiga kolom, bukan ratusan.

### Lapis 3 — penanda per-nilai, hanya di satu tempat

Hanya untuk **nilai uang hasil konversi mata uang**, karena di sana satu mata uang dapat gagal sementara mata uang lain berhasil **di dalam baris yang sama** — satu-satunya kasus yang tidak tertangkap lapis 2.

### Consequences

Usulan pertama saya — satu kolom penanda untuk setiap nilai turunan — ditolak, dan alasan penolakannya benar: ratusan kolom tambahan akan diisi asal saat implementasi, lalu tidak ada yang percaya isinya. **Penanda yang tidak dipercaya lebih buruk daripada tidak ada penanda.**

Urutannya juga dibalik dari usulan saya. Bukan *"kalau terlalu mahal, batasi ke uang"*, melainkan *"mulai dari yang murah dan menyeluruh, lalu perdalam hanya di tempat yang paling mahal kalau salah"*.

**Masalah yang diselesaikan** — tiga mekanisme di sistem lama, semuanya berakar pada tidak adanya cara membedakan nol dari gagal:

| Mekanisme | Kegagalan | Terlihat sebagai |
|---|---|---|
| `GETCURRENCYSTANDARD` | kurs tidak ada | kurs bernilai 1 (`FINDING-006`) |
| Database link HRD | tautan putus | daftar pengguna kosong (`ADR-0016`) |
| Precondition tidak cocok | kondisi gagal | langkah dilewati tanpa pesan (`FINDING-002` bagian 2) |

Ketiganya menghasilkan jawaban yang bentuknya benar dan isinya salah.

## Fakultatif dan Treaty dua entitas berbeda; sistem ini hanya memiliki Treaty

Sumber: `docs/adr/0020-hanya-treaty-yang-dimiliki.md`

*status: accepted · label: DECIDED*


Polis fakultatif dan polis treaty adalah **dua entitas berbeda**. Sistem ini **hanya memiliki Treaty**; Fakultatif dibaca sebagai rujukan, tidak dimodelkan, tidak dimiliki. Bila kelak ada modul yang memerlukannya, modul itu yang memilikinya.

### Consequences

**Alasan pemisahannya adalah domain, bukan bentuk penyimpanan.** Fakultatif dinegosiasikan per risiko, satu per satu; treaty adalah kontrak payung atas satu portofolio. Beda cara lahirnya, beda dokumennya, beda kewajiban para pihaknya.

Argumen bentuk data — bahwa `V_POLIS` bercabang pada `QuotationData.BusinessFac` dan membaca jalur JSON yang berbeda (`$.PolicyData.StartDateTime` versus `$.StartDate`) — adalah **akibat, bukan sebab**. Ia pendukung yang baik dan bukan bukti utama. Menyimpulkan identitas dari bentuk penyimpanan akan bertentangan dengan `ADR-0015`, yang justru menyatukan dua tabel karena identitasnya sama.

**Yang membatasi kepemilikan**: lingkup modul ini adalah klaim non-proporsional, dan itu berjalan di jalur Treaty. Memodelkan Fakultatif berarti membangun entitas yang tidak pernah ditulis dan tidak pernah dihitung — biaya tanpa penerima manfaat.

**Pemeriksaan yang diminta, dan hasilnya tidak sebagaimana diharapkan.** Pertanyaannya: adakah klaim non-prop yang pernah berjalan di cabang `'F'`? Sapuan 279 berkas menemukan **`BusinessFac` nol kemunculan** — modul Claim **tidak pernah membacanya sama sekali**.

Artinya modul Claim tidak dapat membedakan kedua cabang; ia menerima apa pun yang dikeluarkan `V_POLIS`. Jadi pertanyaan itu **tidak dapat dijawab dari XML** dan menjadi pertanyaan data: apakah ada baris `json_polis` ber-`BusinessFac='F'` yang pernah dirujuk sebuah klaim non-prop. → **REQ-030**.

Konsekuensi yang perlu dicatat sekarang: bila ternyata ada, klaim itu menerima `STARTDATETIME`/`ENDDATETIME` berpresisi penuh, sementara jalur Treaty menerima bentuk yang terpotong di jam (`FINDING-005` bagian 6.1) — dua bentuk teks berbeda panjang masuk ke perbandingan yang sama.

Batas kepemilikan ini ditulis di `CONTEXT.md`.

## Pengenal polis dan klaim diberi nama menurut isinya; CASEID dipensiunkan

Sumber: `docs/adr/0021-penamaan-pengenal-polis-dan-klaim.md`

*status: accepted · label: DECIDED*


Istilah `CASEID` **dipensiunkan** dari model internal dan masuk daftar `_Avoid_` di `CONTEXT.md`. Setiap pengenal diberi nama yang menyebut isinya dan pemiliknya.

Sebabnya: nama yang sama membawa arti berbeda di dua tempat. Di `V_POLIS`, `CASEID` = `DATA_JSON.IDNewBisnis`, identitas **polis**. Di `OS_AKSEPTASI_KLAIM`, `CASEID` berbentuk `'ASM-FW-GCNMFW-WORK CLMNP-…'`, identitas **case Pega**.

### Peta pengenal — lebih dari dua

Penamaan tidak ditetapkan sebelum peta ini lengkap, karena berhenti di dua akan mengulang kesalahan yang sama dalam bentuk lebih halus.

| Pengenal | Ditemukan di | Dugaan pemilik | Status |
|---|---|---|---|
| `CASEID` (V_POLIS) = `DATA_JSON.IDNewBisnis` | view `V_POLIS` | kita | `IDNewBisnis` **nol kemunculan** di 279 XML — hanya ada di view |
| `CASEID` (OS_AKSEPTASI_KLAIM) | tabel akseptasi | Pega | berbentuk `<class> <pyID>` |
| `NOPOLIS` | `JSON_KLAIM`, `JSON_POLIS`, `OS_AKSEPTASI_KLAIM`, `TREATYINPRODUCTION`, 2 prosedur | kita | kolom, 4 tabel |
| `.ClaimData.PolicyNo` / `.ClaimData.PolicyData.PolicyNo` | XML, 115 kemunculan | kita | apakah keduanya sama **belum terbukti** |
| `QuotationData.PolicyNoSinarmas` | XML | **perusahaan lain dalam grup** | bentuk keempat, menyebut nama perusahaan |
| `POLICY_CEDING` | tabel work Pega | **cedant** | nomor polis milik cedant |
| `DLANO_CEDING` / `DLANO_SOB` | `OS_AKSEPTASI_KLAIM` | cedant / SOB | nomor DLA, bukan nomor polis |

**Sedikitnya lima bentuk pengenal polis hidup berdampingan, dan sedikitnya dua di antaranya milik pihak lain.** Nomor polis cedant dan nomor polis kita adalah dua hal berbeda; `POLICY_CEDING` memperlihatkan keduanya memang disimpan berdampingan.

Pemetaan mana yang sama dan mana yang berbeda **belum selesai** dan menunggu DDL tabel work serta REQ-023.

### Consequences

**Penggantian nama berhenti di batas integrasi.** Payload REST ke Arasapas dan ke sistem kasir membawa `CASEID` apa adanya. Nama internal boleh berubah; **kontrak ke luar tidak**. Pemetaan nama internal ke nama payload ditulis eksplisit di Boundary Contract, supaya penggantian istilah tidak merambat ke payload dan memutus integrasi yang selama ini berjalan.

Satu nama untuk dua hal tidak menimbulkan galat, hanya salah paham — dan salah pahamnya baru muncul saat dua orang membicarakan hal berbeda dengan kata yang sama.

## Masa berlaku treaty disimpan sebagai tanggal, dan batasnya inklusif

Sumber: `docs/adr/0022-masa-berlaku-treaty-sebagai-tanggal-dan-batas-inklusif.md`

*status: accepted · label: DECIDED*


Masa berlaku treaty disimpan sebagai **tanggal**, bukan tanggal-waktu. Kerugian yang terjadi **tepat pada** tanggal akhir treaty **tertutup** — batasnya inklusif.

### Alasan menyimpan sebagai tanggal — bukan yang saya tulis pertama kali

Alasan versi pertama saya: *"untuk masa berlaku treaty, jam tidak membawa arti"*. **Itu salah untuk reasuransi.** Slip treaty lazim menyebut waktu attachment secara eksplisit — misalnya *"at 00:01 hours Local Standard Time, 1st January"* — dan jam bisa sangat berarti untuk kerugian katastrofa yang terjadi di pergantian periode.

Alasan yang benar lebih sederhana dan lebih kuat:

> **Data sumbernya tidak sanggup menyimpan waktu yang dapat dipercaya.** Jalur Treaty di `V_POLIS` berhenti di jam — tanpa menit dan detik — dan nilainya sendiri hasil `SUBSTR` atas teks.

Perbedaan ini penting: bila kelak ada sumber yang membawa waktu attachment sungguhan, alasan versi pertama akan dipakai untuk menolaknya. Alasan yang benar justru mengundangnya.

### Batas inklusif — dan apa yang sebenarnya dilakukan sistem lama

Seluruh perbandingan di sistem lama memakai operator `>` tegas, tidak pernah `>=`:

```
pyWorkPage.ClaimData.DateOfLoss > pyWorkPage.ClaimData.EndDateTreaty
pyWorkPage.ClaimData.PolicyData.EndDateTime > pyWorkPage.ClaimData.EndDateTreaty
@FormatDateTime(...StartDateTime...) > pyWorkPage.ClaimData.EndDateTreaty
```

Karena penolakan terjadi ketika tanggal kerugian **lebih besar** dari tanggal akhir, kerugian yang jatuh **tepat pada** tanggal akhir tidak ditolak. **Maksud kode lama juga inklusif**, dan itu sejalan dengan praktik reasuransi.

Yang tidak dapat dipastikan adalah apakah maksud itu benar-benar terlaksana, karena perbandingannya dilakukan atas dua format teks yang berbeda (`FINDING-005`). Jadi keputusan ini menegaskan maksud yang sudah ada, bukan mengubahnya.

### Consequences

**Ditetapkan sebelum REQ-020 dijalankan, dan itu bukan formalitas.** Klaim yang tanggal kerugiannya jatuh tepat di batas akan **berpindah sisi** tergantung pilihan ini — masuk atau keluar dari kedua daftar di `ADR-0018`. Bila ditetapkan setelah query berjalan, daftarnya harus dibuat ulang.

Perbandingan di REQ-020 karena itu memakai `TO_DATE(DATEOFLOSS,'YYYYMMDD') > TERMINATION`, bukan `>=`.

## Data akseptasi disimpan dalam satu bentuk kanonik; sistem hilir diberi view, bukan salinan

Sumber: `docs/adr/0023-satu-bentuk-kanonik-tanpa-salinan.md`

*status: accepted · label: DECIDED*


Tidak ada data yang disimpan dua kali. Bila sistem hilir memerlukan bentuk kolom, ia diberi **view atau projeksi** atas bentuk kanonik — bukan salinan yang ditulis terpisah.

Sistem lama menyimpannya dua kali. `POOLDATA.OS_AKSEPTASI_KLAIM` punya **66 kolom**, dan `PEGA_JSON_OS_AKSEP_KLAIMTNP` hanya mengisi **8** di antaranya (`CASEID`, `NOCLAIM`, `MASTERID`, `DATA_JSON`, `TANGGAL`, `NOPOLIS`, `STS_REJECT`, `STS_KONVERSI`). **59 kolom sisanya** — `GROSSVALUE`, `ADJUSTERFEE`, `KURSIDR`, `CNPREINSTATEMENT`, `TOTALXOL`, `SALVAGE`, dan seterusnya — diisi oleh sesuatu yang lain.

`DATA_JSON` adalah sumber kebenaran; kolom skalar adalah keluaran proses hilir.

### Consequences

**Temuan bypass Arasapas menghasilkan ramalan yang dapat diuji, dan itu memperkuat keputusan ini, bukan mengubahnya.** Untuk case yang dibuat `VINCENTVERNANDO_1`, konversi tidak berjalan (`FINDING-002` bagian 8.6). Maka barisnya seharusnya memiliki `DATA_JSON` terisi tetapi **kolom skalarnya kosong**.

Bila ternyata terisi, berarti ada penulis lain yang belum diketahui — dan seluruh pemahaman tentang siapa menulis apa ke tabel ini harus ditinjau ulang. Digabungkan ke **REQ-016**.

Siapa penulis dan siapa pembaca 59 kolom itu tetap **REQ-023**. Keputusan di atas **tidak menunggunya**: penyimpanan ganda ditolak apa pun jawabannya, karena dua salinan yang dapat menyimpang satu sama lain adalah cacat terlepas dari siapa yang menulisnya.

## Satu akseptasi per klaim per layer per mata uang

Sumber: `docs/adr/0024-kunci-alami-akseptasi.md`

*status: accepted · label: DECIDED*


Kunci alami baris akseptasi adalah **(klaim, layer, mata uang)**, ditegakkan sebagai `UNIQUE`, dengan primary key surrogate.

Sumbernya bukan tebakan: blok yang **dikomentari** di dalam `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTNP` memakai persis kombinasi itu —

```sql
SELECT count(1) INTO id_count FROM OS_AKSEPTASI_KLAIM a
 WHERE CASEID = PegaID AND a.data_json.TypeLoss = LayerT AND a.data_json.Currency = Currency;
```

Orang yang menulis blok itu sedang menjawab pertanyaan yang sama dan sampai pada kesimpulan yang sama. Kodenya dimatikan; alasannya tidak.

### Penguat dari bentuk penyimpanan — `CLAIMXOL`

Dipindahkan ke sini dari `BLUEPRINT.md` §13.2 pada 18 September 2026, karena tempatnya di sini.

View `POOLDATA.CLAIMXOL` membongkar JSON akseptasi dengan `JSON_TABLE`:

```sql
json_table (data_json, '$.CNPLayerList[*]'
  columns( XOL varchar2 path '$.XOL',
    nested path '$.CNPCurrencyList[*]' columns(
      TotalXOLGross varchar2 path '$.TotalXOLGross', ... )))
```

**`$.CNPLayerList[*]` membungkus `$.CNPCurrencyList[*]`** — layer di tingkat luar, mata uang di tingkat dalam. Jadi bentuk penyimpanan lama sendiri sudah bersusun dua tingkat persis menurut dua dimensi kunci ini, dan urutannya menetapkan arah sarangnya: satu layer memuat banyak mata uang, bukan sebaliknya.

Ini penguat, bukan dasar. Dasarnya tetap blok yang dikomentari di atas, karena di sanalah kombinasinya dinyatakan sebagai **kunci**, bukan sekadar sebagai susunan.

### Consequences

Sistem lama tidak menegakkannya sama sekali: `OS_AKSEPTASI_KLAIM` **tidak punya primary key maupun unique constraint**, dan prosedurnya **selalu `INSERT`, tidak pernah `UPDATE`** karena logika upsert-nya dikomentari seluruhnya. Baris ganda bukan kemungkinan teoretis.

**Karena itu pelanggaran kunci akan ditemukan, dan itu bukan bukti kuncinya salah.** Saat REQ-018 kembali, hasilnya dipisah tiga:

| Kelompok | Perlakuan |
|---|---|
| Baris **identik di semua nilai** | duplikat murni — digabung |
| Baris **berbeda, dengan pola waktu** | ambil yang terakhir, **catat yang dibuang** |
| Baris **berbeda, tanpa pola waktu** | bawa contohnya untuk ditinjau — berarti ada dimensi pembeda yang belum tertangkap |

Kelompok ketiga adalah satu-satunya yang dapat menggugurkan kunci ini. Dua kelompok pertama adalah akibat bug selalu-`INSERT`, bukan bukti tentang model data.

## Tanggal tutup buku disimpan bertanggal berlaku, bukan satu baris tanpa riwayat

Sumber: `docs/adr/0025-tutup-buku-bertanggal-berlaku.md`

*status: accepted · label: DECIDED*


Aturan penggeseran periode dipertahankan: nomor yang diterbitkan setelah hari tutup buku masuk periode berikutnya. Yang dirancang ulang adalah **penyimpanannya**.

Tanggal tutup buku disimpan dalam tabel **bertanggal berlaku** — setiap perubahan menjadi baris baru dengan periode berlakunya sendiri. Ditambah kolom **lingkup**, dengan nilai bawaan global.

### Consequences

Sistem lama menyimpannya sebagai **satu baris tanpa riwayat**: `SELECT TO_NUMBER(tanggal) INTO v_day_closing FROM POOLDATA.TANGGAL_CLOSING WHERE ROWNUM = 1`.

**Dan itulah sebabnya tambalan lahir.** `PROC_GENERATE_SEQUENCE_NUMBER` memuat:

```sql
IF TRUNC(v_now) <= TO_DATE('02/01/2026','DD/MM/YYYY') THEN
    v_mm_yyyy := '12.2025';  v_tahun := 2025;
```

Perlakuan khusus untuk satu periode tidak punya tempat di data, jadi ia menjadi cabang di dalam kode. Dengan tanggal berlaku, kasus seperti itu cukup menjadi **satu baris data** — dan kelas tambalan ini kehilangan alasan untuk lahir lagi.

Kolom lingkup disediakan sekarang meski belum ada yang memerlukannya; biayanya hampir nol di tahap ini dan mahal ditambahkan setelah tabel terisi.

**Satu masalah yang tidak diselesaikan oleh keputusan ini** dan harus ditangani terpisah: aturan tanggal 25 ada di **dua sumber kebenaran**. `Activity\HitServiceToKasir_Act.xml` menuliskan `@if(TempKasir.CARI19 > 25, ...)` langsung di kode Pega, sementara prosedur membacanya dari tabel. Mengubah tabel tidak mengubah kode. Dibawa ke akuntansi sebagai `ASK-AKUNTANSI.md` pertanyaan 6.3.

## Batas kepemilikan mengikuti nama class: -Work- dan -Data- dimiliki, -Int- tidak

Sumber: `docs/adr/0026-batas-kepemilikan-mengikuti-nama-class.md`

*status: accepted · label: DECIDED*


Modul klaim **memiliki** klaim, akseptasi, alokasi, dan adjustment. Selebihnya **dibaca, tidak dimiliki**.

Aturannya tidak perlu ditimbang per tabel — `MEMORI_PEMAHAMAN.MD` §2.3 sudah memberinya: **integration class (`-Int-`) adalah pemetaan langsung ke tabel atau view Oracle.**

| Pola class | Kepemilikan |
|---|---|
| `…-Work-…` | **dimiliki** |
| `…-Data-…` | **dimiliki** |
| `…-Int-…` | **dibaca, tidak dimiliki** |

Dengan itu `V_POLIS`, `T_STORAGE_IMAGE`, `EMAILKOMITE`, dan `M_LINK_SERVICE` seluruhnya di luar kepemilikan — tanpa perlu memutuskan satu per satu.

### Consequences

**Tabel yang bukan milik klaim tidak dimigrasikan, meskipun ada di daftar tarikan.** Daftar itu untuk **memahami**, bukan untuk memindahkan. Ini memotong lingkup migrasi data secara langsung.

Bila modul lain belum siap pada saat cutover, klaim **membaca dari basis data lama lewat antarmuka yang disepakati** — bukan menyalin datanya. Konsekuensinya cutover dapat dilakukan per modul, tidak harus sekaligus.

Yang tetap perlu dipahami meski tidak dimiliki: 298 kolom di DDL tanpa pasangan properti Pega (`BLUEPRINT.md` §19.2) sebagian besar milik modul lain. Memahaminya perlu; memindahkannya tidak.

## Dokumen klaim dirujuk, tidak disimpan ulang

Sumber: `docs/adr/0027-dokumen-dirujuk-bukan-dimiliki.md`

*status: accepted · label: DECIDED*


Sistem baru **merujuk** dokumen klaim, tidak menyimpannya. Tidak ada migrasi berkas fisik.

### Consequences

`MEMORI_PEMAHAMAN.MD` §5.7: berkas sesungguhnya berada di **Google Cloud Storage**. `POOLDATA.T_STORAGE_IMAGE` hanya menyimpan **metadata dan URL beserta `EXPDATE`**, yang disegarkan lewat `GetUrlGoogleStorage_Act`.

Jadi kekhawatiran tentang volume berkas yang harus dipindahkan **gugur** — tidak ada berkas fisik di dalam basis data.

Yang tetap perlu dirancang: URL berumur terbatas (`EXPDATE`), sehingga sistem baru memerlukan mekanisme penyegaran yang setara dan tidak boleh menyimpan URL sebagai nilai tetap.

`POOLDATA.GET_TOKEN_STORAGE` dan class `Link-Attachment` termasuk lapisan rujukan ini, bukan lapisan penyimpanan.

## Basis data tujuan Oracle, skema baru bersebelahan dengan POOLDATA

Sumber: `docs/adr/0028-basis-data-tujuan-dan-letak-skema.md`

*status: proposed · label: DECIDED (menunggu konfirmasi lisan)*


> **Status `proposed`, bukan `accepted`.** Isinya disusun dari pilihan A yang tertulis lengkap sementara pilihan B dibiarkan kosong — itu **pembacaan atas formulir, bukan kalimat orang**. Dinaikkan ke `accepted` hanya setelah ada konfirmasi eksplisit, dan tanggal konfirmasi itu dicatat di sini. Folder ini memisahkan `proposed` dari `accepted` dengan ketat; ADR-0003 adalah buktinya.

Basis data tujuan adalah **Oracle**. Skema baru berdiri sebagai **satu skema tersendiri di instance yang sama dengan `POOLDATA`**.

Versi ditulis sebagai **sekurangnya 12.1, angka pastinya terbuka**. Lantai itu berdasar: `BLUEPRINT.md` §11 menyimpulkan basis data **sumber** berjalan di 12c ke atas dari tiga hal sekaligus — `CHECK (x IS JSON)`, notasi titik `a.data_json.Field`, dan klausa `CREATE OR REPLACE EDITIONABLE` yang baru ada sejak 12.1. **Itu lantai sumber, bukan versi tujuan**; keduanya tidak boleh dicampur, dan angka pasti tujuan masih ditunggu.

### Considered Options

Basis data lain tidak dipilih. Alasannya bukan selera melainkan biaya pada ADR yang sudah berdiri: pindah mesin membuat **ADR-0016** kehilangan arti — foreign data wrapper pada dasarnya database link dengan nama lain — dan membuat **ADR-0023** mustahil, karena view lintas mesin bukan view melainkan salinan yang menyamar. Keduanya harus ditulis ulang lebih dulu, dan tidak ada DDL yang boleh ditulis di atas ADR yang sudah tidak berlaku.

Instance terpisah juga tidak dipilih, atas alasan yang sama dalam bentuk lebih kecil: ADR-0023 mewajibkan sistem hilir diberi **view**, bukan salinan. Antar instance, view menuntut database link, dan itu dilarang ADR-0016.

### Consequences

#### Satu pintu tulis tidak bertahan sebagai kesepakatan — ia ditegakkan lewat grant

Menaruh skema baru bersebelahan dengan `POOLDATA` berarti menaruhnya di tempat yang sistem lain sudah terbiasa menulisinya lewat SQL langsung. Itu bukan kekhawatiran teoretis: DDL tabel work Pega memuat

```sql
GRANT ALTER, DELETE, INDEX, INSERT, REFERENCES, SELECT, UPDATE, ... TO POOLDATA;
```

`POOLDATA` memegang `INSERT`, `UPDATE`, `DELETE`, bahkan `ALTER` atas tabel milik Pega. Apakah hak itu benar-benar dipakai belum diketahui (REQ-021, BLOCKER) — tetapi hak yang ada akan dipakai cepat atau lambat, dan **ADR-0017 yang hanya berupa kesepakatan akan runtuh di tempat seperti ini.**

Maka penegakannya masuk DDL, bukan catatan:

| Pihak | Hak |
|---|---|
| Akun pemilik skema | seluruhnya — hanya dipakai saat pemasangan dan migrasi |
| Akun aplikasi | `INSERT`, `UPDATE`, `DELETE`, `SELECT` atas tabel kanonik |
| Seluruh akun lain, termasuk `POOLDATA` | `SELECT` **hanya atas view** ADR-0023 — tidak ada satu pun hak atas tabel kanonik |

`GRANT` dan `REVOKE`-nya ditulis di DDL sebagai objek tersendiri, bukan dijalankan tangan. Hak yang tidak tertulis di DDL adalah hak yang tidak dapat diaudit.

#### Versi 12.1 sebagai dasar, dengan penanda di tempat yang lebih ringkas di atasnya

DDL ditulis supaya berjalan di 12.1. Setiap tempat yang lebih ringkas di versi lebih baru ditandai di spec, tidak dipakai diam-diam. Yang menggigit paling awal adalah **batas panjang pengenal**: 30 byte sampai 12.1, 128 byte sejak 12.2. Nama constraint dan index menabraknya lebih dulu daripada nama tabel.

Selama versi pasti belum diketahui, spec wajib memuat **tabel singkatan tertutup** — daftar tetap, ditulis sekali, tidak boleh ditambah saat menulis DDL. Singkatan ad-hoc adalah cara `UR` dan `MDP` lahir, dan itu persis yang glosarium tutup.

#### Tidak ada kolom JSON untuk data yang dimiliki

Seluruh alasan migrasi ini adalah memberi tiap nilai kolomnya sendiri. Sistem lama menyimpan nilai uang sebagai **teks di dalam JSON**, dan bahkan view `CLAIMXOL` mengeluarkannya sebagai `varchar2` (`BLUEPRINT.md` §13.2, §13.5).

JSON hanya boleh muncul di **tabel pendaratan migrasi**. Itu menjadikan tabel pendaratan satu-satunya tempat bentuk lama boleh masuk utuh, dan membuat aturan penguraian ADR-0014 bekerja di satu batas yang jelas — antara pendaratan dan kanonik — bukan tersebar di banyak tempat.

#### Yang diwarisi ADR lain

**ADR-0016** tetap berlaku apa adanya: data pegawai masuk lewat API atau salinan tersinkron. Instance yang sama tidak mengubahnya — `hrdasm.v_hrd_mst@asmd.sinarmas.co.id` yang ditemukan di `V_MST_USER_TEKNIS` adalah link ke mesin **lain**, dan itu yang dilarang.

**ADR-0023** menjadi dapat dilaksanakan: view lintas skema di satu instance adalah view sungguhan, tanpa salinan dan tanpa link.

**ADR-0017** berpindah dari kesepakatan menjadi konfigurasi yang bisa diperiksa, lewat tabel grant di atas.


---

# V. Yang masih terbuka

Tiga register. Yang menunggu orang, yang menunggu basis data, dan yang menunggu akuntansi.

## OPEN QUESTIONS — Claim Non Prop

Sumber: `_selesai/OPEN-QUESTIONS.md`

**Pertanyaan: aktif 31 · selesai 23 · mati 0** — *cuplikan; angka berjalannya di `_selesai/OPEN-QUESTIONS.md`*

Status: **EXTERNAL** (butuh data luar XML) · **DEFERRED-TO-KOMITE-SESSION** · **PARKIR** (menunggu hasil recon Oracle) · **KONFLIK** (XML vs memori).

---

### A. Butuh data atau keputusan di luar XML — **ditriase 18 September 2026**

Sebelumnya seluruh bagian ini ditandai EXTERNAL. **Itu salah tempat untuk separuhnya.** EXTERNAL di dokumen ini berarti *butuh manusia*, dan itu tidak sama dengan *tidak ada di XML*. DDL tabel work memuat `PYREOPENCOUNT`, `PYREOPENTIMESTAMP`, `PXCREATEOPERATOR`, `PXUPDATEOPERATOR`, `PXCREATEDATETIME`, `PXUPDATEDATETIME`, `PXAPPLICATION`, `PXAPPLICATIONVERSION`, `PXOBJCLASS`, `PYSTATUSWORK`, dan `PYID` — dan itu mengubah status sebagian butir di sini.

> **Aturan yang berlaku untuk seluruh dokumen ini mulai sekarang**: sebelum sesuatu ditandai EXTERNAL, harus dibuktikan dulu tidak ada kolom di tabel work Pega maupun di DDL `POOLDATA` yang menjawabnya. Membawa keluar pertanyaan yang jawabannya ada di basis data mereka sendiri membuat permintaan berikutnya lebih sulit didengar.

#### A-QUERY — dipindahkan ke `ORACLE-REQUESTS.md` 18 September 2026

Enam butir — `A4`, `A5`, `A6a`, `A9`, `A10a`, `A11a` — **tidak lagi hidup di dokumen ini**. Jalurnya pindah ke REQ-024, REQ-006, REQ-027, REQ-025, REQ-016, dan REQ-026.

Bukan karena terjawab, melainkan karena satu butir tidak boleh hidup di dua register: bila REQ-024 dijawab dan seseorang lupa menutup A4, register ini berhenti dapat dipercaya. Daftarnya beserta REQ tujuan ada di bagian **SELESAI**.

#### A-MANUSIA — benar-benar butuh keputusan atau pengetahuan orang

| # | Pertanyaan | Kenapa tidak ada kolom yang menjawabnya | Dibawa ke |
|---|---|---|---|
| A1 | Nama workbasket tujuan assignment Akseptasi | Ada di tabel assignment Pega (`pc_assign_*`), yang **tidak** ikut di DDL mana pun. Sampai tabel itu masuk, ini tetap manusia | admin Pega |
| A2 | Nama peran / access group sebenarnya | `pyPrivilegeName` kosong di 279 berkas; jawabannya di `Data-Admin-Operator-ID` dan `pr_operators` yang belum dapat diakses (REQ-001) | admin Pega |
| A3 | Boleh/tidaknya satu orang melakukan Registrasi sekaligus Akseptasi | Kontrol internal, bukan konfigurasi. Data hanya menunjukkan apa yang terjadi, bukan apa yang diizinkan | pemilik proses |
| A6b | **Arti** tiap nilai `PaymentType` | `Rule-Obj-FieldValue` tidak ada di export. Sebarannya terukur (A6a), maknanya tidak | pemilik proses |
| A7 | Apakah `CloseClaimMD` memang boleh menutup klaim tanpa Komite | Kebijakan. Frekuensinya terukur (A5), kebolehannya tidak | pemilik proses |
| A8 | Apakah tarif brokerage 2,5% / PPh 2% / PPN 2,2% masih berlaku | Tarif tertanam di kode tanpa tanggal berlaku | akuntansi — `ASK-AKUNTANSI.md` no. 5 |
| A10b | Apakah `VINCENTVERNANDO_1` akun uji atau akun bisnis nyata | Pemakaiannya terukur (A10a); sifatnya tidak | HRD / admin Pega |
| A11b | Status **kepegawaian** kelima nama | Aktivitas di sistem terukur (A11a); status kepegawaian ada di HRD | HRD |
| A12 | Alamat `prweb` mana yang benar untuk produksi | Dua alamat tertanam berdampingan; tidak ada kolom yang menyatakan mana yang sah | admin Pega |
| A13 | Siapa berwenang menyetujui pengecualian pembekuan tambalan | Kebijakan organisasi (ADR-0012) | manajemen |
| A15 | **Apakah ada node/instance Pega terpisah untuk bisnis syariah** | Lihat bagian G1 | admin Pega |

Catatan tambahan: memori §3.4 memuat workbasket `komitepnc`..`komitepnc4`, tetapi itu milik **modul Komite** (`KomiteRouter`), bukan jawaban atas A1. Workbasket sisi Claim tetap EXTERNAL.

### B. PARKIR — menunggu hasil recon Oracle

Pertanyaan berikut **sengaja tidak ditanyakan** sampai `pengetahuan/SCHEMA-ACTUAL.csv` diterima. Selama daftar ini berisi, frontier **tidak boleh dinyatakan kosong**.

| # | Pertanyaan | Menunggu |
|---|---|---|
| B1 | Presisi dan skala final untuk setiap kolom nilai | REQ-001 saja. **REQ-002/REQ-005 tidak lagi dapat menjawabnya** untuk `.ClaimData.*` — nilainya tidak punya kolom. Cadangan bila REQ-001 gagal: bongkar `PZPVSTREAM` |
| B2 | **Presisi tidak dijaga di tiga lapis sekaligus.** (1) Tipe kolom: 88 dari 120 kolom `NUMBER` di DDL tanpa presisi maupun scale; (2) JSON: nilai uang tersimpan sebagai **teks**, dan view `CLAIMXOL` mengeluarkannya sebagai `varchar2`; (3) kode Pega: 77% ekspresi aritmetika tanpa skala. **REQ-005 tidak dapat menjawab ini untuk `.ClaimData.*`** — nilainya tidak punya kolom. Jawabannya hanya dapat datang dari **membaca isi blob dan isi JSON**, bukan dari metadata mana pun | REQ-001 / bongkar `PZPVSTREAM`; **bukan** REQ-005 |
| B4 | Struktur JSON di `json_klaim.DATA_JSON` dan `OS_AKSEPTASI_KLAIM.DATA_JSON` | REQ-009 |
| B6 | Apakah `EMAILKOMITE.DEGREE` menentukan urutan jenjang, dan apa arti `LIMIT_TOP` yang tidak dipakai filter | REQ-004 |
| B7 | Kardinalitas riil tiap PageList (rata-rata dan maksimum baris) untuk menentukan indeks | REQ-010 |
| B8 | Apakah FINDING-001 nyata atau gugur — berapa Adjustment non-IDR, berapa yang setara di atas ambang, dan apakah `.ValueAdjustment` memang sudah IDR | REQ-011 |
| B9 | Apakah jalur `CloseClaimMD` dimigrasi atau dibuang (status `CANDIDATE-NOT-MIGRATED`) | REQ-006 |

---

### C. DEFERRED-TO-KOMITE-SESSION

| # | Hal | Titik sentuh di folder ini |
|---|---|---|
| C1 | Isi `KomiteTreaty_Flow` dan mekanisme loop persetujuan | `CreateChildKomiteCNP_Act` step 29, `CreateChildKomiteCloseNP_Act` step 10 |
| C2 | Siapa yang menulis `CNPStatusCase = "CLAIM ACCEPTED"` / `"CLAIM REJECTED"` | tidak ditulis di folder ini |
| C3 | Bagaimana Komite menutup klaim induk pada jalur `CloseClaimNP` | `CreateChildKomiteCloseNP_Act` tidak memanggil `ASMForceCaseClose` |
| C4 | Bagaimana `AcceptanceStatus` dikembalikan ke `AdjustmentList(IndexObject)` | `Adjustment.IndexObject` |
| C5 | Nasib `Adjustment.AlokasiXOLPaid` — diisi siapa | dibaca `CreateChildKomiteCNP_Act`, penulisnya tidak ada di folder ini |
| C6 | **Siapa menulis `.ValueAdjustment`, dan dalam mata uang apa** | Muncul di **tepat satu berkas** (`CreateChildKomiteCNP_Act`), hanya dibaca. Bukan field UI di folder ini. Asal nilainya **TIDAK DITEMUKAN DI XML**. Menentukan sah-tidaknya FINDING-001 |
| C8 | **Siapa membaca `.IsEditClaim`** — dua hipotesis sejajar, lihat §F di bawah | Ditulis `1` oleh `EditXOLAlokasi` dan `0` oleh `CountLossAllocation_act`. **Tidak ada satu pun pembacaan** di folder ini — bukan di precondition, bukan di `pyVisible`, bukan di SQL. Lihat `BLUEPRINT.md` §2.5 |
| C9 | Siapa menulis `TreatyName=="Previously Calculated UR"` | Dibaca `GenerateCACNP_Act` baris 7223 dan 7425, **tidak pernah ditulis** di folder ini. Nilai `"UR"` biasa ditulis oleh `CountLossAllocation_act`; varian ini tidak |
| C7 | Siapa memindahkan `.AcceptanceStatus` dari `0` ke `1`/`2` | `.AcceptanceStatus` ditulis **satu kali saja** di seluruh 279 berkas: `CreateChildKomiteCNP_Act` langkah 31 → `= 0`. Transisi ke `1`/`2` tidak ada di folder ini |

---

### E. Anomali yang belum dijelaskan

| # | Temuan | Bukti |
|---|---|---|
| E1 | Hanya 2 rule yang mencampur skala pembagian, dan keduanya rule alokasi inti: `CountClaimTNP_Act` (10×28, 20×13, 5×2) dan `CountLossAllocation_act` (10×14, 20×14 tepat berimbang). Perimbangan 14:14 terlalu rapi untuk kebetulan | `pengetahuan/arithmetic-inventory.tsv` |
| E2 | 77% ekspresi aritmetika tidak menyatakan skala sama sekali — hasilnya ditentukan perilaku bawaan Pega yang tidak terdokumentasi di export | `pengetahuan/arithmetic-inventory.tsv` |
| E5 | Deskripsi langkah tidak sinkron dengan kondisi eksekusi di `InputOutStandingClmTNP_PreAct` step 5 | `BLUEPRINT.md` §7 |
| E11 | Dua alamat `prweb` berbeda tertanam berdampingan: `http://192.168.105.116:80/prweb` dan `http://pega.nusantarare.com:80/prweb` | `SaveToOS`, `SendEmailKlaim`, `SendEmailKlaimRejectClose` |

---

### F. Hipotesis sejajar — `.IsEditClaim`

Menandai ini "ditunda" saja tidak cukup: kedua kemungkinan punya konsekuensi yang harus dinyatakan sekarang, karena keduanya mengubah pekerjaan yang berbeda.

| | **F1 — pembacanya ada di modul Komite** | **F2 — tidak ada pembaca sama sekali** |
|---|---|---|
| Artinya | Bendera ditulis di Claim, dibaca di Komite | Bendera ini mati |
| Konsekuensi | **Kopling lintas modul yang belum ada di Boundary Contract.** Sudah dicatat sebagai tambahan di `BLUEPRINT.md` §8.5 | **Perlindungan hasil edit manual tidak pernah ada.** Hasil suntingan petugas lewat `EditXOLAlokasi` tertimpa setiap kali `CountLossAllocation_act` dijalankan ulang — dan `CountLossAllocation_act` dipanggil dari sedikitnya 8 titik (lihat `Struktur_Flow_TreatyIn.xlsx`) |
| Yang harus dikerjakan | Tambahkan ke kontrak batas, tentukan arah dan pemiliknya | Naikkan menjadi **FINDING-004**, dan pastikan sistem baru punya mekanisme kunci hasil edit yang sebenarnya |
| Cara memastikan | Sapu modul Komite atas `IsEditClaim` — **satu sesi, satu grep** | Hasil sapuan yang sama, bila nihil |
| **Apa yang mustahil ada di data kalau ini benar** | — (F1 tidak melarang apa pun) | Kalau F2 benar, **mustahil ada baris `SpreadingRisk` dengan `IsEditClaim = 1` yang bertahan** setelah alokasi dijalankan ulang. Satu baris begitu saja menggugurkan F2 |

Keduanya ditutup oleh **satu pemeriksaan yang sama** di sesi Komite. Sampai itu terjadi, **tidak boleh** ada keputusan desain yang mengandaikan salah satunya benar.

**Uji dari sisi Claim, tanpa menunggu sesi Komite**: `.IsEditClaim` **IN-BLOB** (tidak termasuk 19 kolom yang di-expose), jadi ujinya baru mungkin setelah blob terbaca — REQ-001. Dicatat supaya tidak terlupa saat REQ-001 masuk.

---

### G. Pertanyaan terbuka & asumsi bernama — dicatat supaya tidak hilang

#### G1. "syariah" — **pertanyaan terbuka, bukan asumsi**. Kesimpulan versi pertama dicabut.

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

#### G2. Urutan kerja petugas tidak terdokumentasi di sumber mana pun

Sapuan `MEMORI_PEMAHAMAN.MD` atas SOP, urutan kerja, langkah petugas, prosedur operasional, dan tata cara: **nihil**.

Digabung dengan temuan bahwa urutan perhitungan ditentukan kontrol yang ditekan (`FINDING-003` bagian 3), artinya: **urutan tindakan yang menentukan hasil perhitungan hanya hidup di kepala petugas.** Tidak ada satu pun artefak yang menangkapnya.

Keputusan Q22 (ADR-0013) membuat sistem baru tidak bergantung pada urutan, jadi risikonya tertutup ke depan. Yang tidak tertutup: **menafsirkan data lama** kadang memerlukan pengetahuan tentang urutan yang ditempuh, dan pengetahuan itu tidak terdokumentasi.

---

### H. Hipotesis sejajar — `.ValueAdjustment` (C6)

Bentuk yang sama dengan bagian F, dan atas alasan yang sama: dua kemungkinan dengan konsekuensi berbeda, ditutup oleh satu pemeriksaan yang sama.

`.ValueAdjustment` muncul di **tepat satu berkas** (`CreateChildKomiteCNP_Act`), hanya dibaca, bukan field UI mana pun di folder ini. Asal nilainya **TIDAK DITEMUKAN DI XML**.

| | **H1 — sudah dikonversi ke IDR di hulu** | **H2 — mata uang aslinya, apa pun itu** |
|---|---|---|
| Artinya | Modul lain mengonversinya sebelum diserahkan | Nilai disimpan apa adanya per mata uang |
| Konsekuensi | **`FINDING-001` gugur seluruhnya.** Perbandingan terhadap ambang `30.000.000` sah, dan tidak ada cacat apa pun | **`FINDING-001` berdiri.** Ambang dibandingkan terhadap angka yang satuannya tidak diketahui |
| Yang harus dikerjakan | Tutup FINDING-001, catat sebagai dugaan yang gugur | Kuantifikasi paparan lewat REQ-011, dan `ADR-0007` menjadi koreksi atas cacat nyata |
| Cara memastikan | Sapu modul Komite atas `.ValueAdjustment` — **satu sesi, satu grep** | Hasil sapuan yang sama |
| **Apa yang mustahil ada di data kalau ini benar** | Kalau H1 benar, **mustahil ada Adjustment bermata uang selain IDR yang `.ValueAdjustment`-nya sama dengan nilai asli mata uang itu**. Bila ditemukan nilai yang jelas bukan hasil konversi — misalnya nilai USD kecil berdampingan dengan `Currency='USD'` — H1 gugur | Kalau H2 benar, mustahil seluruh `.ValueAdjustment` pada klaim non-IDR berskala rupiah |

Keduanya ditutup oleh **satu pemeriksaan yang sama** di sesi Komite. Sampai itu terjadi, **tidak boleh** ada keputusan desain yang mengandaikan salah satunya benar.

---

### I. Hipotesis sejajar — transisi `.AcceptanceStatus` (C7)

`.AcceptanceStatus` ditulis **satu kali saja** di seluruh 279 berkas: `CreateChildKomiteCNP_Act` langkah 31 → `= 0`. Transisi ke `1` atau `2` tidak ada di folder ini.

| | **I1 — ditulis oleh modul Komite** | **I2 — tidak pernah ditulis sama sekali** |
|---|---|---|
| Artinya | Komite menandai hasil persetujuannya di properti ini | Properti ini menetap di `0` selamanya |
| Konsekuensi | Kopling lintas modul yang **belum ada di Boundary Contract** — arah Komite → Claim. Tambahkan | Penjaga `CloseClaimTNonProp` yang menolak `=="0"` akan **menolak setiap klaim yang pernah punya Adjustment**, selamanya. Jalur `CloseClaimMD` efektif mati untuk klaim tersebut |
| Yang harus dikerjakan | Lengkapi kontrak batas, tetapkan pemiliknya | Naikkan jadi temuan tersendiri; dan `BLUEPRINT.md` §4.0 perlu ditulis ulang lagi |
| Cara memastikan | Sapu modul Komite atas `.AcceptanceStatus` — **satu sesi, satu grep** | Hasil sapuan yang sama |
| **Apa yang mustahil ada di data kalau ini benar** | — | Kalau I2 benar, **mustahil ada klaim yang punya Adjustment dan ditutup lewat `CloseClaimMD`**. Satu baris begitu saja menggugurkan I2 |

**I1/I2 dapat diselesaikan sekarang, tanpa menunggu sesi Komite.** Ramalan I2 dapat dibantah data: cacah klaim tertutup, dipecah menurut ada-tidaknya Adjustment. Digabungkan ke **REQ-019**, yang sudah menghitung klaim tertutup dan `FLAGONGOINGCOMMITTE`.

Bila berhasil, ini **butir DEFERRED-TO-KOMITE-SESSION pertama yang diselesaikan dari sisi Claim**.

---

## SELESAI

Butir yang sudah terjawab atau terserap ke artefak permanen. Tidak dihapus: pertanyaan yang pernah diajukan adalah bukti bahwa jawabannya tidak dikarang.

### Terjawab / tertutup

| # | Isi | Ke mana jawabannya pergi |
|---|---|---|
| ~~B3~~ | ~~Properti mana yang EXPOSED sebagai kolom vs IN-BLOB di `pzPVStream`~~ | **TERTUTUP 18 Sep 2026.** DDL tabel work memperlihatkan hanya 19 kolom di luar bawaan Pega, dan **tidak satu pun nilai uang**. Seluruh `.ClaimData.*` bernilai uang **IN-BLOB — CONFIRMED**, bukan lagi dugaan. Lihat `BLUEPRINT.md` 13.5 |
| ~~B5~~ | ~~Isi sebenarnya stored procedure `PEGA_JSON_OS_AKSEP_KLAIMTNP` dkk — pemetaan `CARIn` ke kolom~~ | **TERTUTUP 18 Sep 2026.** Source-nya diterima lewat `pengetahuan/ddl/PROCEDURE_PEGA_JSON_OS_AKSEP_KLAIMTNP.sql` beserta 48 berkas lain di `pengetahuan/ddl/`. *(`REQ-003` adalah permintaannya, bukan jawabannya.)* Isinya melahirkan `ADR-0024`; dari berkas lain dalam kiriman yang sama, `pengetahuan/ddl/FUNCTION_GETCURRENCYSTANDARD.sql`, lahir `FINDING-006` |
| E3 | `SaveAdjustmentToOSAksep_Act_Tes` bernama "Tes" tetapi dipanggil dari alur produksi (`CreateChildKomiteCNP_Act` step 25) | `Struktur_Flow_TreatyIn.xlsx` |
| E4 | `Harness\Hitung_Test.xml` + `Section\Hitung_Test.xml` (240 KB + 187 KB) terdaftar sebagai local action pada alur produksi | `Struktur_Flow_TreatyIn.xlsx` |
| E6 | Tambalan **per-orang** berumur 8 tahun dan masih menyebar: `pxCreateOperator=="VINCENTVERNANDO_1"` ditulis 2018-01-16, lalu **disalin** ke rule lain pada 2022-01-03 dan 2023-11-20 oleh dua orang berbeda | `BLUEPRINT.md` §7.1, §7.2 |
| E7 | Keputusan alur diambil dengan mencocokkan **teks bebas**: `@contains(.CommentSuggest,"Accepted by Himawan")` dan `"Rejected by Himawan"` di `SethistoryKlaimTreaty`. Satu salah ketik pengguna mengubah jalur | `BLUEPRINT.md` §7.1 |
| E8 | `CLMNP-232` ditambal **dua kali berjarak 3,5 tahun** oleh dua orang di dua rule, dan keduanya masih aktif bersamaan | `BLUEPRINT.md` §7.2 |
| E9 | Tambalan terbaru (`CLMNP-975`, 2026-07-16) berjarak **dua bulan** dari hari ini — praktik ini belum berhenti | `BLUEPRINT.md` §7.2 |
| ~~E10~~ | ~~Sembilan rule menyentuh `SpreadingRisk` secara intensif (16–66 rujukan) tanpa menyebut `"UR"`~~ | **DIKOREKSI 18 Sep 2026.** Angka 16–66 menghitung nama class di metadata langkah, bukan rujukan PageList. Rujukan sebenarnya: `GenerateCFS_act` 5, `CountTotalInsterest_Act` 2, `CopyOldataCurr_act` 2, **enam lainnya nol**. Bukan anomali — `BLUEPRINT.md` §2.4 |

### Terserap ke temuan gabungan

`E3`, `E4`, `E6`, `E7`, `E8`, `E9` digabung menjadi satu temuan tentang tata kelola perubahan — `BLUEPRINT.md` bagian 17. `E10` dicoret sebagai salah hitung, bukan anomali — `BLUEPRINT.md` §2.4.

### Pindah jalur ke `ORACLE-REQUESTS.md`

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

### Nomor yang berpindah

| # | Ke mana | Sebab |
|---|---|---|
| ~~A14~~ | **A10b** | Isinya — apakah `VINCENTVERNANDO_1` akun uji atau akun bisnis nyata — digabung saat A10 dipecah menjadi A10a (terukur dari kolom) dan A10b (butuh orang). Dicatat di sini supaya nomor yang lenyap tidak jadi tanda tanya di kemudian hari |

### Konflik XML vs memori yang sudah diputus

XML menang untuk perilaku. Tujuh butir, dan **tidak seragam**: `D7` menyatakan memori **benar**; `D5` menyatakan yang keliru adalah pembacaan saya sendiri, bukan memori maupun XML.

| # | Memori | XML | Putusan |
|---|---|---|---|
| D1 | §7.3 mendaftar `CNPStatusCase` = `"COMITEE ACCEPTANCE (DEPT. HEAD)"`, `"CLAIM ACCEPTED"`, `"CLAIM REJECTED"` | Di folder ini ada nilai keempat yang tidak tercatat di memori: **`"INPUT ACCEPTATION CLAIM"`** (`InputOutStandingCTNP_PostAct` step 2) | XML menang; enum di `BLUEPRINT.md` §3.1 yang berlaku |
| D2 | §5.5 menggambarkan `CloseClaimMD` sebagai jalur "langsung" tanpa menyebut penjaga apa pun | Penjaganya ada, tetapi **menguji hal yang berbeda dari yang disangka**: ia menolak `AcceptanceStatus=="0"`, sedangkan Adjustment yang **belum pernah dikirim** ke Komite bernilai **kosong**, bukan `"0"` — jadi lolos. Pesan galatnya sendiri membuka maksudnya: *"there is adjustment in comitee"* — mencegah penutupan saat komite **sedang bersidang**, bukan mensyaratkan persetujuan komite | **BYPASS BERSYARAT.** Putusan versi pertama ("bukan bypass tanpa penjaga") **dicabut** — lihat `BLUEPRINT.md` §4.0 |
| D3 | §10.2 mendaftar 3 tambalan per-case (`IDMaster 1000393`, `CLMNP-367`, `CLMNP-382`) | Angka terakhir: **29 langkah, 18 ekspresi unik, 8 rule**, mencakup tambalan per-case **dan** per-identitas — dan itu belum termasuk nilai mati di dalam DDL (`PROC_GENERATE_SEQUENCE_NUMBER`) | XML menang; daftar lengkap di `BLUEPRINT.md` §7, §7.1, §7.5. *(Angka 9 identitas / 4 rule / 14 langkah pada versi pertama sudah usang.)* |
| D4 | §7.4 memuat `.PICSuggest == "CHRISTINEANGELINA"` dan `"NANDINA"` | XML memuat **dua ejaan berdampingan**: `"CHRISTINEANGELINA"` (1 langkah) dan `"Christine Angelina Hutagalung"` (2 langkah); `"NANDINA"` (1) dan `"Nandina C"` (2) | XML menang; memori mencatat satu dari dua bentuk. Akibatnya sebagian cabang tidak pernah terpicu — FINDING-002 bagian 6.2 |
| D5 | Memori §5.5 **dan** analisis awal saya sama-sama melewatkan hal yang sama: keduanya membaca keberadaan penjaga sebagai bukti adanya syarat persetujuan | Penjaga itu menguji *apakah ada sidang berjalan*, bukan *apakah persetujuan sudah diberikan*. Dua pertanyaan berbeda | Bukan memori yang salah dan bukan XML yang salah — **pembacaan saya yang salah**, dan memori kebetulan sejalan dengan pembacaan itu |
| D6 | `SetPPNPPH` **tidak muncul sama sekali** di memori — tidak di katalog rumus §6, tidak di 15 rule terpenting §12, tidak di inventaris §2.1 | Rule itu **terjangkau** dari alur klaim: `Section\OutstandingClaim(1)` → `Harness\ViewPolisNonProp` (class `ASM-FW-GCNMFW-Work-ClaimTreatyNonProp`) → `<pyInclude>Section\ViewDetailDeptHeadTreatyIn_UW` → `<pyActivity>CountNetPremi_act` → `Call SetPPNPPH` | XML menang; memori melewatkannya. Tarif pajak **tetap dalam lingkup** — `ASK-AKUNTANSI.md` pertanyaan 5 |
| D7 | §10.4 butir 2 menyatakan `TotalUR` selalu nol pada cabang mata uang sama | **Terverifikasi persis dari XML**: `(.Deductible * 0 * Local.ProrateClaim/100)` di penetapan `Local.TotalUR`, sementara `Local.UR` pada cabang yang sama tidak memakai `* 0` | Memori **benar**. Naik dari catatan memori menjadi EVIDENCED; dibawa ke akuntansi sebagai pertanyaan 2b |

Dipertahankan sebagai catatan bahwa setiap putusan punya dasar.

## ORACLE REQUESTS — Claim Non Prop

Sumber: `ORACLE-REQUESTS.md`

**REQ: aktif 31 · selesai 1 · mati 0** — *cuplikan; angka berjalannya di `ORACLE-REQUESTS.md`*

Register permintaan objek Oracle. Status: **OPEN** → **SENT** → **ANSWERED** → **DEAD** (objek tidak ada).
Nomor REQ tidak pernah dipakai ulang. Selama ada REQ berstatus OPEN/SENT, frontier **tidak boleh dinyatakan kosong**.

| REQ | Objek | Jenis | Confidence | Prioritas | Status | Tanggal |
|---|---|---|---|---|---|---|
| REQ-001 | Isi definisi `Rule-Obj-Property` di schema PegaRULES (`PR4_*`) | TABLE | **AKSES BELUM ADA** — struktur standar bawaan produk, bukan misteri | **BLOCKER** | OPEN | 2026-09-17 |
| REQ-002 | `ALL_TAB_COLUMNS` untuk 44 tabel bisnis | TABLE | CONFIRMED + DERIVED | **BLOCKER** | OPEN | 2026-09-17 |
| REQ-004 | `EMAILKOMITE` — struktur + isi konfigurasi jenjang | TABLE + DATA-PROFILE | DERIVED | PENTING | OPEN | 2026-09-17 |
| REQ-005 | Profil presisi kolom uang & kurs — **dipersempit: hanya tabel bisnis `POOLDATA`**. Untuk `.ClaimData.*` jalur ini **tertutup permanen** (IN-BLOB terbukti) | DATA-PROFILE | CONFIRMED | PENTING | OPEN | 2026-09-18 |
| REQ-006 | Statistik pemakaian jalur tutup-langsung (mengambil alih A5) | DATA-PROFILE | CONFIRMED | PENTING | OPEN | 2026-09-17 |
| REQ-007 | Verifikasi keberadaan 8 objek DERIVED | TABLE/VIEW/SYNONYM | DERIVED | PENTING | OPEN | 2026-09-17 |
| REQ-008 | Tabel work Pega untuk `ASM-FW-GCNMFW-Work` | TABLE | GUESS | PENTING | OPEN | 2026-09-17 |
| REQ-009 | Struktur JSON di `json_klaim.DATA_JSON` & `OS_AKSEPTASI_KLAIM.DATA_JSON` | DATA-PROFILE | CONFIRMED | PENTING | OPEN | 2026-09-17 |
| REQ-010 | Constraint, index, dan volume tabel inti | INDEX/CONSTRAINT | CONFIRMED | PELENGKAP | OPEN | 2026-09-17 |
| REQ-011 | Kuantifikasi Adjustment non-IDR vs ambang kewenangan Komite | DATA-PROFILE | CONFIRMED | **BLOCKER** (untuk FINDING-001) | OPEN | 2026-09-18 |
| REQ-012 | Isi `Data-Admin-DB-Table` untuk class `ASM-FW-%` | TABLE | **AKSES BELUM ADA** — struktur standar bawaan produk | **BLOCKER** | OPEN | 2026-09-18 |
| REQ-013 | Cacah pola teks pada `.CommentSuggest` | DATA-PROFILE | DERIVED | PENTING | OPEN | 2026-09-18 |
| REQ-014 | Cacah klaim dengan lebih dari satu baris `TreatyType='UR'` | DATA-PROFILE | DERIVED | PENTING | OPEN | 2026-09-18 |
| REQ-015 | Status buka/tutup 8 klaim bertambalan | DATA-PROFILE | DERIVED | PENTING | OPEN | 2026-09-18 |
| REQ-016 | Profil case buatan `VINCENTVERNANDO_1` (mengambil alih A10a) — cacah, rentang tanggal, nilai, jumlah persetujuan, berapa yang tidak sampai ke Arasapas, **dan apakah 59 kolom skalar `OS_AKSEPTASI_KLAIM` kosong untuk case-case itu** (menguji ADR-0023) | DATA-PROFILE | DERIVED | **BLOCKER** | OPEN | 2026-09-18 |
| REQ-017 | DDL 14 objek yang dirujuk dari dalam procedure/view tetapi tidak disertakan (`M_CURRENCYSTANDARD`, `TANGGAL_CLOSING`, `OS_AKSEPTASI_SUBJECTIVITY`, `MST_USER_TEKNIS`, dll.) | TABLE/VIEW | CONFIRMED | **BLOCKER** (untuk `M_CURRENCYSTANDARD`) | OPEN | 2026-09-18 |
| REQ-018 | Cacah `PYID` yang muncul lebih dari sekali di tabel work, **dan** cacah baris ganda `(CASEID, TypeLoss, Currency)` di `OS_AKSEPTASI_KLAIM` — satu REQ, dua tabel, pola yang sama | DATA-PROFILE | DERIVED | **BLOCKER** | OPEN | 2026-09-18 |
| REQ-019 | Cacah `FLAGONGOINGCOMMITTE` menyala saat klaim ditutup, sebaran `KOMITECOUNT`, **dan klaim tertutup dipecah menurut ada-tidaknya Adjustment** — yang terakhir menguji hipotesis I2 | DATA-PROFILE | CONFIRMED | **PENTING → menutup butir DEFERRED** | OPEN | 2026-09-18 |
| REQ-020 | Berapa klaim melewati penjaga `CheckDateDOL_Act` dengan `DateOfLoss` di luar masa treaty (FINDING-005). **Keterjangkauan sudah diperiksa: dapat dijalankan dengan SQL biasa** — `TO_DATE(DATEOFLOSS,'YYYYMMDD')` vs `TREATYINDETAIL.TERMINATION` (`DATE`). Tidak perlu bongkar blob | DATA-PROFILE | CONFIRMED | **PRASYARAT CUTOVER** | OPEN | 2026-09-18 |
| REQ-021 | Grant dari `DATAPEGA` ke `POOLDATA`, dan adakah source `POOLDATA` yang menulis ke `DATAPEGA.PC_*` | METADATA | DERIVED | **BLOCKER TEKNIS + BLOCKER PERKIRAAN BIAYA** | OPEN | 2026-09-18 |
| REQ-022 | Mata uang **bukan IDR** yang kursnya tersimpan `1` — mengukur jalur `NO_DATA_FOUND` pada `GETCURRENCYSTANDARD` | DATA-PROFILE | CONFIRMED | PENTING | OPEN | 2026-09-18 |
| REQ-023 | Siapa menulis dan siapa membaca ~60 kolom skalar `OS_AKSEPTASI_KLAIM` yang tidak diisi `PEGA_JSON_OS_AKSEP_KLAIMTNP` | METADATA | DERIVED | **BLOCKER** | OPEN | 2026-09-18 |
| REQ-024 | `PYREOPENCOUNT > 0` dan `PYREOPENTIMESTAMP` — apakah klaim pernah di-reopen (menjawab A4, mengonfirmasi atau membatalkan ADR-0002) | DATA-PROFILE | CONFIRMED | **BLOCKER** | OPEN | 2026-09-18 |
| REQ-025 | `PXAPPLICATION`, `PXAPPLICATIONVERSION` dikelompokkan dengan `PXOBJCLASS` — aplikasi/ruleset lain yang menyentuh class ini (menjawab A9) | DATA-PROFILE | CONFIRMED | PENTING | OPEN | 2026-09-18 |
| REQ-026 | `PXCREATEOPERATOR` / `PXUPDATEOPERATOR` + `PXUPDATEDATETIME` untuk kelima nama di FINDING-002 — aktivitas terakhir tiap orang (menjawab A11a) | DATA-PROFILE | CONFIRMED | PENTING | OPEN | 2026-09-18 |
| REQ-027 | Sebaran nilai `OS_AKSEPTASI_KLAIM.PAYMENTTYPE` — nilai mana yang benar-benar dipakai (menjawab A6a) | DATA-PROFILE | CONFIRMED | PENTING | OPEN | 2026-09-18 |
| REQ-028 | Cacah baris berkode setoran `"100115"` di `DIRECTTOKASIR_LOG` atau tabel kasir — apakah klaim non-prop syariah pernah ada (menutup G1) | DATA-PROFILE | CONFIRMED | PENTING | OPEN | 2026-09-18 |
| REQ-029 | Akseptasi ber-`STS_SUBJECTIVITY='1'` yang **sudah punya catatan pembayaran** — apakah uang keluar atas syarat yang belum terpenuhi | DATA-PROFILE | CONFIRMED | **BLOCKER** | OPEN | 2026-09-18 |
| REQ-030 | Adakah baris `json_polis` ber-`BusinessFac='F'` yang dirujuk klaim non-prop (menguji ADR-0020) | DATA-PROFILE | CONFIRMED | PENTING | OPEN | 2026-09-18 |
| REQ-031 | **Cacah baris per entitas yang dimiliki** (ADR-0026): klaim, akseptasi, alokasi, adjustment — beserta sebarannya per tahun dan per status buka/tutup. Menentukan volume migrasi ADR-0018 dan ukuran tabel korelasi | DATA-PROFILE | CONFIRMED | PENTING | OPEN | 2026-09-18 |
| REQ-032 | **Versi instance tujuan** — dua nilai saja: `SELECT version_full FROM v$instance` (atau `v$version`), dan `SHOW PARAMETER COMPATIBLE`. Menentukan batas panjang pengenal: **30 byte bila `COMPATIBLE` < 12.2**, 128 byte bila ≥ 12.2 — dan karenanya menentukan apakah tabel singkatan tertutup perlu ditulis sama sekali | VERSI | CONFIRMED | **BLOCKER** (untuk penamaan DDL) | OPEN | 2026-09-18 |

Seluruh query siap jalan ada di [`pengetahuan/RECON.sql`](./pengetahuan/RECON.sql), dipetakan per bagian:

| REQ | Bagian pengetahuan/RECON.sql |
|---|---|
| REQ-001 | BAGIAN 1 (1A–1D, tiga varian) |
| REQ-002 | BAGIAN 3 |
| REQ-003 | BAGIAN 5 (5A, 5B, 5C) |
| REQ-004 | BAGIAN 3 + BAGIAN 7B |
| REQ-005 | BAGIAN 6 (6A–6D) |
| REQ-006 | BAGIAN 8 |
| REQ-007 | BAGIAN 2 (2A, 2B) |
| REQ-008 | BAGIAN 2C |
| REQ-009 | BAGIAN 9 |
| REQ-010 | BAGIAN 4 + BAGIAN 7 |
| REQ-011 | BAGIAN 10 |
| REQ-012 | BAGIAN 11 — *belum ditulis, menunggu T5a* |
| REQ-013 | BAGIAN 12 — *belum ditulis, menunggu T2* |
| REQ-014 | BAGIAN 13 — *belum ditulis, menunggu T2* |
| REQ-015 | BAGIAN 14 — *belum ditulis* |
| REQ-016 | BAGIAN 15 — *belum ditulis* |
| REQ-017 | BAGIAN 3 (perluasan daftar objek) |
| REQ-018 | BAGIAN 16 — *belum ditulis* |
| REQ-019 | BAGIAN 17 — *belum ditulis* |
| REQ-020 | BAGIAN 18 — *belum ditulis* |
| REQ-021 | BAGIAN 19 — *belum ditulis* |
| REQ-022 | BAGIAN 20 — *belum ditulis* |
| REQ-023 | BAGIAN 21 — *belum ditulis* |
| REQ-024 | BAGIAN 22 — *belum ditulis, menunggu DDL tabel work* |
| REQ-025 | BAGIAN 22 — idem |
| REQ-026 | BAGIAN 22 — idem |
| REQ-027 | BAGIAN 23 — *belum ditulis* |
| REQ-028 | BAGIAN 24 — *belum ditulis* |
| REQ-029 | BAGIAN 25 — *belum ditulis* |
| REQ-030 | BAGIAN 26 — *belum ditulis* |
| REQ-031 | BAGIAN 27 — *belum ditulis* |
| REQ-032 | BAGIAN 28 — *belum ditulis* |

### Penomoran — tabrakan diselesaikan 18 September 2026

Karena jawaban Round 5 dan artefak Round 5 ditulis berbarengan, dua nomor bertabrakan. Diselesaikan tanpa menomori ulang apa pun yang sudah ditulis:

| Tetap | Dipindahkan |
|---|---|
| REQ-018 `PYID` ganda + duplikat akseptasi | rasio kurs ≈ 1 → **REQ-022** |
| REQ-019 `FLAGONGOINGCOMMITTE` / `KOMITECOUNT` | penulis & pembaca kolom skalar → **REQ-023** |
| REQ-020 penjaga `CheckDateDOL_Act` | — |
| `FINDING-005` bug perbandingan tanggal | `RETURN 1` pada kurs → **`FINDING-006`** |

### PREFLIGHT — didahulukan dari seluruh REQ

[`pengetahuan/PREFLIGHT.sql`](./pengetahuan/PREFLIGHT.sql) menjawab T1–T7 dari kamus data, tanpa menyentuh tabel bisnis.
Hasilnya menentukan bagaimana `pengetahuan/RECON.sql` ditulis ulang:

| Jawaban | Yang berubah |
|---|---|
| T1 versi Oracle | boleh/tidaknya `FETCH FIRST`, `search_condition_vc` |
| T2 tipe kolom JSON | BAGIAN 6 dan BAGIAN 10 dipakai apa adanya atau ditulis ulang dengan `JSON_VALUE` |
| T3 owner 25 nama tanpa prefix | kolom `OWNER` di `pengetahuan/PULL-LIST.csv`; yang AMBIGU tetap `?` |
| T4 hak akses | arti "tidak ditemukan": **TIDAK ADA** vs **TIDAK TERLIHAT** |
| T5 letak PegaRULES | REQ-001 dan REQ-012 diajukan ke DBA yang mana |
| T7 ukuran tabel | BAGIAN 6, 7, 10: sampling atau full scan |
| T8 (dijawab: tidak ada UAT) | seluruh pengetahuan/RECON.sql harus aman di produksi |

### Catatan penanganan

- Hasil ditempel ke [`pengetahuan/SCHEMA-ACTUAL.csv`](./pengetahuan/SCHEMA-ACTUAL.csv) mengikuti header yang sudah disediakan.
- Bila sebuah objek dilaporkan **tidak ada**, statusnya menjadi **DEAD** dan temuan itu dicatat di `_selesai/OPEN-QUESTIONS.md` — berarti ada rule yang merujuk objek mati.
- Bila hasil bertentangan dengan dugaan di `TABLE-EXTRACTION-REQUEST.md`, koreksinya ditulis terang-terangan di sini beserta apa yang berubah karenanya.

---

## SELESAI

| REQ | Isi | Jenis | Keyakinan | Prioritas | Status | Tanggal |
|---|---|---|---|---|---|---|
| ~~REQ-003~~ | Source 6 procedure + 1 function | PROCEDURE/FUNCTION | CONFIRMED | PENTING | **ANSWERED** — seluruhnya ada di `pengetahuan/DDL_Script_ClaimNonProp.xls`, tersimpan di `pengetahuan/ddl/` | 2026-09-18 |

`REQ-003` tertutup karena seluruh source procedure dan function diterima lewat `pengetahuan/DDL_Script_ClaimNonProp.xls` dan tersimpan di `pengetahuan/ddl/`. Isinya melahirkan `FINDING-006` (kurs `RETURN 1`), `ADR-0024` (kunci alami akseptasi), dan `ASK-AKUNTANSI.md` bagian 0.

## Pertanyaan untuk Akuntansi

Sumber: `ASK-AKUNTANSI.md`

**Enam pertanyaan.** Lima tentang kebijakan presisi dan tarif; satu berisi tiga butir lama yang belum pernah dijawab.

---

### 0. Definisi: apa yang dimaksud "nilai yang masuk jurnal"

Pertanyaan 1, 2, dan 3 memakai frasa ini, jadi isinya ditetapkan lebih dulu. Daftar ini **tertutup** — diturunkan dari penelusuran seluruh jalur yang membawa nilai keluar dari modul klaim non-proporsional.

#### 0.1 Tiga jalur keluar, dan hanya satu yang membawa uang

| Jalur | Rule | Yang dibawa |
|---|---|---|
| **Kasir** | `Activity\HitServiceToKasir_Act.xml` → `ConnectREST\SendAcceptationToKasir.xml` | **nilai uang** — menjadi instruksi bayar |
| **Arasapas** | `Activity\KonversiKlaim_Act.xml` → `ConnectREST\KonversiKlaimNonLife.xml` | **hanya `CASEID`, `NOPOLIS`, `STS_REJECT`** — pemicu, bukan data nilai |
| **Persistensi** | `Activity\SaveDataToOSAksep_Act.xml` → `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTNP` | seluruh dokumen `DATA_JSON` |

Jalur Arasapas **tidak membawa satu pun nilai uang**. Nilainya sudah lebih dulu tersimpan di `OS_AKSEPTASI_KLAIM.DATA_JSON` lewat jalur persistensi, dan Arasapas membacanya dari sana.

#### 0.2 Daftar tertutup — nilai yang menjadi instruksi bayar

Dari `Activity\HitServiceToKasir_Act.xml`, muatan `TempKasir.CARIn`:

| Posisi | Isi | Sifat |
|---|---|---|
| `CARI9` | `.TotalClaim - .PremiumSpreaded` | **nilai uang** |
| `CARI8` | `.AdjustmentValue` | **nilai uang** |
| `CARI22` | `@toDecimal(TempKasir.CARI9) + .AdjustmentValue` | **nilai uang** |
| `CARI23` | `@toDecimal(TempKasir.CARI9) + .AdjusterFeeValue` | **nilai uang** |
| `CARI24` | `@toDecimal(TempKasir.CARI9) + .SalvageValue` | **nilai uang** |
| `CARI9` (cabang lain) | `.IndividualRiskRNM` | **nilai uang** |
| `CARI10` | `.IndividualRiskPercentage` | persentase |
| `CARI11`, `CARI12`, `CARI13` | `"NUSARE"`, `"D0031"`, `"100081"` | kode tetap |
| `CARI15`/`CARI17`/`CARI18` | `"100115"` pada satu cabang | kode tetap, berbeda per node |

**Jadi "nilai yang masuk jurnal" berarti tepat tujuh besaran**: `.TotalClaim`, `.PremiumSpreaded`, `.AdjustmentValue`, `.AdjusterFeeValue`, `.SalvageValue`, `.IndividualRiskRNM`, dan `.IndividualRiskPercentage`. Aturan "2 desimal, tanpa toleransi" berlaku untuk ketujuh itu, **tidak untuk nilai antara**.

**Catatan yang perlu diketahui**: ketujuh nilai itu disusun dengan `@toDecimal(...)` **tanpa pembersihan teks** — lihat pertanyaan 4.

---

### 1. Berapa desimal untuk nilai antara?

**Sudah disepakati**: nilai yang masuk jurnal (daftar 0.2) dibulatkan ke **2 desimal**. Yang belum: berapa desimal untuk perhitungan di tengah jalan.

#### 1.1 Keadaan sekarang

| Fakta | Sumber |
|---|---|
| Dari 120 kolom `NUMBER` di basis data, **88 tidak menyatakan presisi maupun scale sama sekali** | `pengetahuan/ddl/*.sql`, terhitung di `pengetahuan/SCHEMA-ACTUAL.csv` |
| Satu-satunya yang menyatakan: `NUMBER(20,4)` pada 23 kolom nilai | `pengetahuan/ddl/TABLE_TREATYINPRODUCTION.sql` |
| Nilai klaim **tidak punya kolom sama sekali** — tersimpan sebagai teks di dalam JSON | `pengetahuan/ddl/VIEW_CLAIMXOL.sql`, seluruh kolom nilai bertipe `varchar2` |
| 77% ekspresi aritmetika di kode tidak menyatakan skala | `pengetahuan/arithmetic-inventory.tsv`, 673 ekspresi |
| Skala pajak `8` | `Activity\SetPPNPPH.xml` baris 822, 890, 911 — `@divide(2.5,100,8)`, `@divide(2,100,8)`, `@divide(2.2,100,8)` |

**Basis data tidak memaksakan pembulatan apa pun.** Apa pun yang dihitung, tersimpan apa adanya.

#### 1.2 Di mana angkanya benar-benar pecah — contoh berangka

Rumus premi reinstatement. **Ada dua varian di dalam sistem, dan keduanya berbeda tepat pada hal yang sedang ditanyakan:**

| Berkas & baris | Ekspresi | Skala |
|---|---|---|
| `Activity\CountReinstatement_Act.xml` **baris 624** | `((((.ClaimEstimation+.AdjusterFee) - (.Salvage*100/RNMShare)) / .CNPLimit) * .CNPMDP) * (.CNPPctReinstate/100)` | **tidak dinyatakan sama sekali** |
| `Activity\AdjClaimCNP_Act.xml` **baris 3109** | `@divide(.TotalClaim, .CNPLimit, 20) * .CNPMDP * @divide(.CNPPctReinstate, 100, 20)` | **20** |

Varian pertama menyerahkan presisinya kepada perilaku bawaan Pega, yang tidak terdokumentasi di dalam ekspor ini. Varian kedua menetapkan 20 desimal secara tegas. Keduanya juga memakai **pembilang yang berbeda** — yang satu `ClaimEstimation + AdjusterFee − Salvage×100/RNMShare`, yang lain `TotalClaim` — sehingga bukan sekadar dua penulisan untuk rumus yang sama.

Bentuknya sama: satu pembagian, lalu dua perkalian. Pembulatan pada hasil bagi **diperbesar** oleh dua pengali sesudahnya.

Contoh: pembilang Rp 4.490.977.654,68 · CNPLimit Rp 12.500.000.000 · CNPMDP Rp 1.750.000.000 · Reinstate 100%

| Desimal pada hasil bagi | Hasil bagi | Premi reinstatement | **Selisih** |
|---|---|---|---|
| **2** | 0,36 | Rp 630.000.000,00 | **+Rp 1.263.128** |
| **4** | 0,3593 | Rp 628.775.000,00 | **+Rp 38.128** |
| **6** | 0,359278 | Rp 628.736.500,00 | −Rp 372 |
| **8** | 0,35927821 | Rp 628.736.867,50 | −Rp 4 |
| tanpa pembulatan | 0,35927821237… | Rp 628.736.871,66 | — |

**Pada satu klaim.** Dua desimal meleset Rp 1,26 juta; empat desimal masih meleset Rp 38 ribu.

#### 1.3 Pertanyaan

Berapa desimal untuk nilai antara?

| | Nilai antara | Akibat menurut tabel 1.2 |
|---|---|---|
| **A** | 4 desimal | masih meleset ±Rp 38 ribu per klaim |
| **B** | 6 desimal | meleset ±Rp 372 per klaim |
| **C** | 8 desimal | meleset ±Rp 4 per klaim |

**Usulan: B, enam desimal** — dengan catatan bahwa rantai perhitungan yang memuat pembagian diikuti perkalian besar sebaiknya memakai delapan.

**Dua hal yang perlu disampaikan terus terang:**

1. Dugaan bahwa **alokasi berlapis** menumpuk galat antar layer (`Activity\CountLossAllocation_act.xml`) **tidak terbukti** pada pengujian: dengan 7 layer, hasilnya identik pada 2, 4, 6, maupun 20 desimal. Kerusakannya ada di rumus reinstatement, bukan di perulangan alokasi.
2. Dugaan bahwa **pembagian-lalu-perkalian-kembali** dengan persentase share merusak angka (`AdjusterFeeValue = AdjusterFee ÷ (RNMShare÷100)`, lalu `TotalXOLRNM = TotalClaim × (ClaimPercentage÷100)`) juga **tidak terbukti**: pembulatan akhir ke 2 desimal menyerap galatnya, bahkan bila nilai antara hanya 2 desimal. Diuji pada share 37,5%, 33,33%, 12,35%, 7,77%, dan 23,45%.

Jadi yang membenarkan enam desimal hanyalah tabel 1.2. Itu satu alasan, tetapi alasan yang kuat.

#### 1.4 Pertanyaan sejarah, bukan penentu

**Dari mana `NUMBER(20,4)` pada `TREATYINPRODUCTION` berasal?** Bila ia ketentuan resmi, ia ketentuan untuk **nilai tersimpan di tabel itu**, bukan untuk presisi perhitungan di tengah jalan. Jawabannya berguna, tetapi tidak menentukan jawaban 1.3.

**Menghambat**: `ADR-0003` (draft sejak awal).

---

### 2a. Bolehkah pembulatan mengubah angka lama?

Bila data produksi menyimpan lebih dari 2 desimal, menetapkan pembulatan 2 desimal **akan mengubah angka yang selama ini diterima akuntansi**.

**Pertanyaan**: ikuti presisi lama apa adanya, atau bulatkan dan terima selisihnya?

**Menghambat**: `ADR-0003`, dan kriteria lulus shadow-run di `ADR-0005`.

---

### 2b. Bila perhitungan lama ternyata **keliru** — dipertahankan atau diperbaiki?

Ini pertanyaan yang berbeda sifatnya dari 2a, dan tidak boleh dijawab bersamaan. 2a soal presisi; 2b soal nilai yang hilang sama sekali.

#### 2b.1 Contoh yang sudah terverifikasi

`Activity\CountLossAllocation_act.xml`, penetapan `Local.TotalUR`:

```
@if(Local.Currency=="IDR",  (@divide(.Deductible,Local.Kurs,10)*Local.ProrateClaim/100),
@if(Local.Currency==.Currency, (.Deductible * 0 * Local.ProrateClaim/100),
                               ... ))
```

Pada cabang **`Local.Currency == .Currency`** — yaitu ketika mata uang klaim sama dengan mata uang yang sedang diproses — perhitungannya dikalikan **`0`**. Hasilnya **selalu nol**.

Bandingkan dengan penetapan `Local.UR` pada berkas yang sama, di cabang yang sama persis:

```
@if(Local.Currency==.Currency, (.Deductible * Local.ProrateClaim/100), ... )
```

Tanpa `* 0`.

**Jadi `TotalUR ≠ UR` untuk kasus mata uang sama, dan `TotalUR` bernilai nol.** Ini bukan soal desimal — ini nilai yang tidak pernah terbentuk.

#### 2b.2 Nilai nol itu ikut tersimpan

`Activity\CountLossAllocation_act.xml` baris 7015–7120 menulis baris Retensi Cedant:

```
SpreadingRisk(<LAST>).TotalClaim = Local.TotalUR
```

Jadi nol itu **masuk ke baris Retensi Cedant**, yang pada sistem baru menjadi tabel tersendiri (`ADR-0010`). Bila perilaku lama dipertahankan, tabel baru itu akan memuat kolom yang selalu nol untuk seluruh klaim bermata uang sama.

#### 2b.3 Pertanyaan

Bila ditemukan perhitungan lama yang keliru — bukan kurang presisi, melainkan salah — apakah sistem baru **mempertahankan kesalahannya** demi kesesuaian angka dengan data lama, atau **memperbaikinya** dan menerima selisih terhadap data lama?

Kedua jawaban sah. Akibatnya sangat berbeda, dan itu sebabnya dipisah dari 2a.

**Menghambat**: `ADR-0010`, dan definisi baseline shadow-run di `ADR-0005`.

---

### 3. Toleransi shadow-run untuk nilai antara

**Sudah disepakati**: nilai yang masuk jurnal (daftar 0.2) harus **sama persis sampai 2 desimal, tanpa toleransi**.

**Pertanyaan**: berapa toleransi untuk nilai antara?

**Konteks besaran** — agar angka toleransi dapat dinilai: contoh di tabel 1.2 memakai `CNPLimit` Rp 12,5 miliar dan `CNPMDP` Rp 1,75 miliar, yang merupakan besaran wajar untuk satu layer XOL. Nilai klaim tertinggi dan rata-rata per tahun **belum tersedia** — itu menunggu profil data (REQ-005), dan bila diperlukan sebelum pertemuan, angkanya dapat ditarik lebih dulu.

**Usulan: dua batas berjalan bersama, yang dilanggar duluan yang berlaku.**

| | Batas |
|---|---|
| Relatif | 0,01% dari nilai baris |
| Mutlak | Rp 1.000 per baris |

Alasannya: batas relatif saja terlalu longgar untuk klaim besar — 0,1% dari Rp 50 miliar adalah Rp 50 juta. Batas mutlak saja terlalu ketat untuk klaim besar. Dua batas bersama menutup kedua ujungnya.

Setiap selisih di atas ambang wajib dijelaskan satu per satu, tidak diabaikan sebagai derau.

**Menghambat**: `ADR-0005`.

---

### 4. Ambang nilai untuk menghentikan migrasi

#### 4.1 Ini fakta, bukan kemungkinan

Nilai uang di sistem lama tersimpan sebagai **teks**. Dan sistem lama **sudah pernah menemui teks yang tidak dapat diurai**, lalu menanganinya:

```
@toDecimal(@replaceAll(.Deductible2, ",", "."))
```
`Activity\CountLossAllocation_act.xml` — mengganti koma menjadi titik sebelum mengubah ke angka.

Tidak ada yang menulis pembersihan semacam itu untuk masalah yang tidak pernah terjadi. **Data dengan koma sebagai pemisah desimal ada di produksi.**

#### 4.2 Penanganannya tidak seragam — dan sangat timpang

Sapuan seluruh 279 berkas atas `toDecimal`:

| | Jumlah |
|---|---|
| Seluruh pemanggilan `toDecimal` | **127** |
| Didahului pembersihan `replaceAll` | **9** — seluruhnya atas `.Deductible2` |
| **Tanpa pembersihan apa pun** | **118** |

Properti yang paling sering diubah ke angka **tanpa dibersihkan**: `TempDla.CARI34` (12×), `pyWorkPage.TreatyInMaster.Share` (9×), `OutSpreading.pxResults` (7×), `.KursValue` (6×), `TempKasir.CARI20` dan `CARI21` (6× masing-masing), `.TSIPerObject` (4×), `.BalanceDueTo` (3×), `.PersenRNM` (3×).

**Dua di antaranya adalah muatan ke kasir** — `TempKasir.CARI20` dan `CARI21` — yaitu nilai yang masuk jurnal menurut daftar 0.2.

Satu tempat membersihkan; 118 tempat tidak. Setiap tempat yang membersihkan menandai satu jenis kotoran yang pernah ditemui seseorang; setiap tempat yang tidak membersihkan adalah calon kegagalan.

#### 4.3 Pertanyaan

**Ambangnya berbasis nilai, bukan cacah baris.** Satu baris senilai lima miliar lebih berat daripada lima ratus baris senilai seratus ribu.

Dua angka:

1. **Nilai total** yang belum terselesaikan yang masih dapat diterima sebelum migrasi dihentikan dan data dibereskan lebih dulu.
2. **Nilai satu baris** yang, bila dilampaui, langsung menghentikan migrasi tanpa menunggu total.

Setiap baris yang tidak dapat diurai tetap diselesaikan satu per satu, berapa pun jumlahnya. Tidak ada baris yang dibuang karena "cuma sedikit".

**Menghambat**: `ADR-0014`, dan rencana pembersihan data sebelum cutover.

---

### 5. Tarif pajak dan brokerage

#### 5.1 Keterjangkauan — diverifikasi lebih dulu

Sempat diragukan apakah `SetPPNPPH` termasuk lingkup klaim non-proporsional, karena rule itu berada pada class `ASM-FW-GISFW-Data-PolicyTreatyIn` — kerangka polis, bukan kerangka klaim.

**Rantainya sampai, dan seluruh pangkalnya di class klaim:**

```
Section\OutstandingClaim(1).xml                      (layar klaim)
  -> Harness\ViewPolisNonProp.xml                    class ASM-FW-GCNMFW-Work-ClaimTreatyNonProp
     -> Section\ViewDetailDeptHeadTreatyIn_UW.xml    <pyInclude>
        -> <pyActivity>CountNetPremi_act
           -> Call SetPPNPPH
```

Jadi pertanyaan ini **sah** dan tetap dibawa. Sifatnya perlu dicatat: jalur itu adalah layar **peninjauan polis** yang dibuka petugas klaim, jadi tarifnya memengaruhi angka premi yang **dilihat**, bukan langsung nilai pembayaran klaim.

#### 5.2 Tarif, beserta letaknya

Seluruhnya di `Activity\SetPPNPPH.xml`:

| Baris | Properti | Ekspresi persis | Tarif |
|---|---|---|---|
| 821–822 | `.BrokerageFee` | `@divide(2.5,100,8) * (.PremiOgp + .PremiOnp)` | **brokerage 2,5%** |
| 868–869 | `.BrokerageFeeSebenarnya` | `@if(.TypeTax=="Inclusive", @divide(.Deduction1, @divide(102.2,100,8), 8), .Deduction1)` | **faktor 102,2** |
| 889–890 | `.PPHValue` | `.BrokerageFeeSebenarnya * @divide(2,100,8)` | **PPh 2%** |
| 910–911 | `.PPNValue` | `.BrokerageFeeSebenarnya * @divide(2.2,100,8)` | **PPN 2,2%** |

**Angka keempat yang belum pernah dilaporkan**: `102.2` pada baris 869 adalah faktor gross-up untuk pajak *inclusive* — yaitu `100 + 2,2`. Ia **turunan dari tarif PPN**, tetapi ditulis terpisah. Bila PPN berubah, angka itu harus ikut diubah di tempat lain, dan tidak ada apa pun yang memaksanya.

**Satu hal yang tidak saya simpulkan**: `.BrokerageFee` dihitung 2,5% dari premi, tetapi `.BrokerageFeeSebenarnya` — yang menjadi dasar PPh dan PPN — diturunkan dari `.Deduction1`, **bukan** dari `.BrokerageFee`. Apakah `.BrokerageFee` karena itu hanya untuk tampilan, tidak dapat saya pastikan dari XML.

#### 5.3 Pertanyaan, tiga bagian

1. Apakah ketiga tarif masih benar per hari ini?
2. Pernahkah berubah sejak sistem berjalan? Bila ya, **klaim lama dihitung dengan tarif lama atau tarif baru** — itu menentukan apakah shadow-run harus menyimpan tarif historis.
3. Apakah tarif dapat berbeda per jenis bisnis, per cedant, atau per mata uang? Kode sekarang hanya mendukung satu nilai untuk semua.

**Yang diusulkan apa pun jawabannya**: tarif pindah ke konfigurasi bertanggal berlaku, bukan tertanam di kode — termasuk faktor `102,2` yang harus diturunkan dari tarif PPN, bukan ditulis terpisah. Itu tidak perlu persetujuan akuntansi; yang perlu adalah **nilai dan tanggal berlakunya**.

**Menghambat**: perhitungan pajak di sistem baru, dan kriteria shadow-run untuk nilai yang masuk jurnal.

---

### 6. Tiga butir lama yang belum pernah dijawab

Terbuka sejak 17 September, tercatat di `MEMORI_PEMAHAMAN.MD` §11. Dibawa sekalian karena dua di antaranya akan sangat mahal bila baru ditanyakan setelah sistem baru jadi.

#### 6.1 Ambang otoritas Komite — 15 atau 30?

`MEMORI_PEMAHAMAN.MD` §11 butir 13. `Activity\CreateChildKomiteCNP_Act.xml` langkah 10 menetapkan `Local.LimitPersenMax = 30.00` **dan** `Local.LimitPersenMaxDivHead = 30.00` — dua ambang yang berbeda perannya diberi nilai yang sama.

**Bila salah satunya seharusnya 15, maka satu jenjang persetujuan tidak pernah berjalan selama ini.**

#### 6.2 `TotalUR` yang selalu nol — perilaku benar atau bukan?

`MEMORI_PEMAHAMAN.MD` §11 butir 14. Ini butir **2b** di atas, sudah terverifikasi dari XML.

#### 6.3 Aturan cut-off tanggal produksi — masih berlaku?

`MEMORI_PEMAHAMAN.MD` §11 butir 17. Aturannya ada di dua tempat, dan keduanya memakai tanggal 25:

| Tempat | Bentuk |
|---|---|
| `Activity\HitServiceToKasir_Act.xml` | `@if(TempKasir.CARI19 > 25, @toDecimal(TempKasir.CARI20) + 1, …)` — bila tanggal lewat 25, bulan digeser ke depan |
| `pengetahuan/ddl/PROCEDURE_PROC_GENERATE_SEQUENCE_NUMBER.sql` | membaca `POOLDATA.TANGGAL_CLOSING`, lalu `ADD_MONTHS(v_now, CASE WHEN TO_NUMBER(TO_CHAR(v_now,'DD')) > v_day_closing THEN 1 …)` |

Yang di kode Pega **menuliskan 25 langsung**; yang di basis data **membacanya dari tabel**. Dua sumber kebenaran untuk satu aturan.

**Pertanyaan**: tanggal 25 masih berlaku? Dan bila tanggal di `TANGGAL_CLOSING` diubah, apakah disadari bahwa angka 25 di kode Pega tidak ikut berubah?


---

# VI. Bahan mentah — dari mana semua ini datang

Isi folder `pengetahuan/`: sumber, data hasil sapuan, dan alat tarik.

## pengetahuan/

Sumber: `pengetahuan/README.md`

Bahan yang **masuk**, bukan kesimpulan yang **keluar**.

Aturan isinya satu kalimat: **kalau isinya tidak saya karang, ia di sini.** Sumber mentah, data hasil ekstraksi, dan alat untuk mengambil data. Kesimpulan — ADR, FINDING, BLUEPRINT, CONTEXT, RINGKASAN — tinggal di folder induk, karena isinya adalah tafsir atas bahan di sini, dan tafsir harus bisa dibantah oleh bahannya.

Dipisahkan 18 September 2026. Seluruh rujukan di folder induk sudah disetel ke jalur baru; 16 jalur unik, semuanya terverifikasi menunjuk ke berkas yang ada.

---

### Isi

#### Sumber — diterima apa adanya, tidak diubah

| Berkas | Asal | Keadaan |
|---|---|---|
| `DDL_Script_ClaimNonProp.xls` | ditempel pengguna 2026-09-18 10:36 | **terpakai** — 48 objek `POOLDATA`, terurai ke `ddl/` |
| `DDL_Script_ClaimNonProp2.xls` | muncul di folder 2026-09-18 10:57 | **belum pernah dibaca** — tidak disebut dalam instruksi mana pun |

#### Hasil uraian sumber

| Berkas | Isi |
|---|---|
| `ddl/` | 49 berkas `.sql`, satu objek satu berkas. Setiap berkas berkepala `-- UNTUK DIBACA, BUKAN UNTUK DIJALANKAN`. Klausa `DROP TABLE` dan `ALTER ... DROP PRIMARY KEY` sudah dibuang |
| `ddl/TABLE_PC_ASM_FW_GCNMFW_WORK.sql` | perkecualian: **bukan** dari berkas `.xls`, ditempel langsung oleh pengguna. Tabel work Pega — **tidak dimigrasi** |
| `SCHEMA-ACTUAL.csv` | 698 baris. 598 dari `DDL_Script_ClaimNonProp.xls`, 100 dari DDL tabel work. Dibedakan lewat kolom `sumber` dan `req_id`. Pemisah `;` |

#### Data mentah hasil sapuan XML

| Berkas | Isi | Dipakai oleh |
|---|---|---|
| `arithmetic-inventory.tsv` | 673 ekspresi aritmetika beserta skala pembagiannya | `BLUEPRINT.md`, `_selesai/OPEN-QUESTIONS.md` E1 dan E2 |
| `rekonsiliasi-kolom-vs-properti.tsv` | 408 kolom DDL × 659 properti Pega; 110 berpasangan, 298 kolom yatim, 549 properti yatim | `BLUEPRINT.md` §19 |

Keduanya **angka batas atas**, bukan angka pasti — dasarnya pencocokan nama.

#### Alat — belum satu pun pernah dijalankan

| Berkas | Guna | Keadaan |
|---|---|---|
| `PULL-LIST.csv` | daftar tarikan 74 objek: 67 `CONFIRMED` (nama terbaca literal di SQL produksi), 7 `DERIVED` (ditebak dari nama class Pega). 19 objek masih ber-`OWNER` tanda tanya | **kanonik.** 49 dari 74 sudah punya DDL di `ddl/`; 25 belum |
| `PREFLIGHT.sql` | 335 baris, tujuh blok T1–T7 — memastikan versi, hak akses, dan keberadaan objek sebelum query apa pun ditulis | T1–T3 terjawab dari DDL tanpa dijalankan; **T4, T5, T7 masih menunggu** |
| `RECON.sql` | 433 baris — memulihkan tipe dan presisi properti, memetakan objek Oracle. Hasilnya ditempel ke `SCHEMA-ACTUAL.csv` | **belum bisa dipakai** — kepalanya memuat `SET PAGESIZE`/`SET LONG`, perintah SQL\*Plus, bukan TOAD. Perlu ditulis ulang |

---

### Batas yang berlaku untuk seluruh folder ini

`RECON.sql` menetapkannya di kepalanya sendiri, sebelum ada satu baris data pun:

```
READ-ONLY. Tidak ada DML. Sampel maks 20 baris.
Kolom teks bebas di-masking. Tanpa data nasabah.
```

Dan yang berlaku untuk `ddl/`: **untuk dibaca, bukan untuk dijalankan.**
