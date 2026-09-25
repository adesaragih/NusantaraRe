# Claim Non Prop

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, berkas tertanggal 2026-09-08 s/d 2026-09-09, rule termutakhir di dalamnya `pxUpdateDateTime = 2026-08-30`.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut. Tambalan yang ditambahkan sesudahnya tidak tercermin di sini; deteksinya lewat sapuan ulang, bukan lewat register.

Konteks penanganan klaim atas kontrak reasuransi treaty non-proporsional (excess of loss), sejak klaim dilaporkan cedant sampai klaim ditutup atau ditolak.

## Language

### Pihak

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

### Kontrak

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

### Klaim

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

### Proses

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

### Dokumen

**Nota Kerugian**:
Dokumen resmi yang memuat rincian perhitungan kerugian dan bagian reasuradur.
_Avoid_: PLA, DLA, CFS, CA, claim advice

---

## Batas kepemilikan

**Treaty** — dimiliki dan dimodelkan sepenuhnya oleh sistem ini.

**Fakultatif** — **tidak dimiliki**. Dibaca sebagai rujukan bila diperlukan, tidak dimodelkan, tidak ditulis. Entitas yang berbeda dari Treaty: dinegosiasikan per risiko satu per satu, sedangkan Treaty adalah kontrak payung atas satu portofolio. Lihat ADR-0020.

---

## Pengenal — nama yang dipakai dan yang dihindari

**Nomor Polis** — pengenal polis milik kita.
_Avoid_: `CASEID`, `NOPOLIS`, `PolicyNo` tanpa keterangan pemilik.

**Nomor Polis Cedant** — pengenal polis yang sama menurut catatan cedant. Berbeda benda dari Nomor Polis, dan keduanya disimpan berdampingan.
_Avoid_: `POLICY_CEDING` sebagai istilah percakapan.

**Nomor Klaim** — pengenal klaim di sistem ini.
_Avoid_: `CASEID`, `pyID`, `case key`.

> **`CASEID` dipensiunkan.** Nama itu dipakai untuk dua benda berbeda di sistem lama — identitas polis di `V_POLIS`, identitas case Pega di `OS_AKSEPTASI_KLAIM`. Ia tidak dipakai lagi dalam percakapan maupun model internal. Pengecualiannya hanya payload integrasi ke Arasapas dan kasir, yang kontraknya tidak berubah (ADR-0021).

Sedikitnya lima bentuk pengenal polis hidup berdampingan di sistem lama, dan sedikitnya dua di antaranya milik pihak lain. Petanya di ADR-0021; pemetaan mana yang sama dan mana yang berbeda belum selesai.

---

## Polis

**Polis** — entitas tersimpan, bukan pandangan gabungan atas kontrak treaty dan tahun produksi.

Dasarnya perilaku pengguna, bukan bentuk penyimpanan: nomor polis **dimasukkan atau dipilih manusia** di layar (`Section/InputAcceptation.xml`, `Section/OutstandingClaim(1).xml`, `Harness/OutstandingClaim.xml`), dan `.ClaimData.PolicyData.PolicyNo` bernilai tunggal. Sesuatu yang punya nomor, dipilih manusia, dan dirujuk dari klaim adalah entitas — bila ia sekadar hasil penggabungan, penggunanya tidak akan memasukkan nomornya.

Bahwa sistem lama menyimpannya sebagai dokumen JSON di `json_polis` dan memandangnya lewat view `V_POLIS` adalah **keputusan penyimpanan**, bukan pernyataan tentang wujudnya.

**Nomor Polis** adalah kandidat kunci alaminya.

Label: **EVIDENCED** untuk keberadaannya; **DECIDED** untuk perlakuannya sebagai entitas tersimpan.

---

## Adjustment

**Adjustment** — **unit transaksi pembayaran klaim.** Satu klaim dapat memiliki banyak Adjustment, yaitu pembayaran bertahap.

Bukan revisi atas estimasi, dan bukan penambahan nilai klaim. Isinya menegaskan sifatnya: `Payable`, `PayableTo`, `NoAccount`, `DirectToKasir`, `StatusKasir`. Dan `PaymentType` membedakan jenisnya — `1` Final, `2` Partial/interim.

Definisi ini diambil dari `MEMORI_PEMAHAMAN.MD` §4.3 dan §7.4, yang merupakan otoritas tertinggi untuk maksud bisnis.
_Avoid_: menyebutnya "revisi", "koreksi", atau "penambahan klaim" — ketiganya menyesatkan ke model data yang berbeda.

Label: **EVIDENCED**.

---

## Kepemilikan — aturan tunggal

Batas kepemilikan mengikuti **nama class**, bukan pertimbangan per tabel:

| Pola | Kepemilikan |
|---|---|
| `…-Work-…`, `…-Data-…` | **dimiliki** |
| `…-Int-…` | **dibaca, tidak dimiliki** — pemetaan langsung ke tabel/view Oracle |

Modul klaim memiliki **klaim, akseptasi, alokasi, dan adjustment**. Polis, treaty, cedant, mata uang, dokumen, dan daftar Komite dibaca dari modul lain. Lihat ADR-0026.

---

## Dokumen

**Dokumen klaim** — dirujuk, tidak dimiliki. Berkasnya berada di Google Cloud Storage; yang tersimpan di basis data hanya metadata dan URL berbatas waktu. Lihat ADR-0027.
