# Kuesioner Underwriting — Keputusan Akseptasi (fokus New Business & Renewal)

**Purpose:** menutup empat pertanyaan tentang perilaku keputusan underwriting yang **tidak dapat
dijawab dari kode sistem lama**. Setiap jawaban akan dicatat sebagai keputusan resmi proyek dan
langsung menentukan cara sistem baru berperilaku. Tanpa jawaban ini, sebagian jalur persetujuan
tidak dapat diimplementasikan tanpa menebak — dan menebak di area ini menyalahkan arah keputusan
underwriting.

**From:** Tim Migrasi Facultative Inward · **To:** Underwriting ·
**How your answers will be used:** dicatat sebagai keputusan proyek (`00-KEPUTUSAN-WORK-OWNER.md`),
lalu diturunkan menjadi aturan di sistem baru dan menjadi dasar pencocokan angka saat kedua sistem
dijalankan berdampingan.

---

> ## ✅ TERJAWAB LENGKAP — 16 September 2026
>
> Keempat butir sudah dijawab Underwriting dan **seluruhnya sudah menjadi keputusan resmi**:
> butir 1 → **K-014** · butir 2 → **K-015** · butir 3 → **K-019** (lewat Konflik K-2) ·
> butir 4 → **K-021**.
>
> Jawaban ditulis di bawah tiap pertanyaan, diikuti keputusan yang dihasilkannya. Naskah
> pertanyaannya **tidak diubah** — hanya stub jawabannya yang diisi.
> Sumber: kuesioner Underwriting terisi, 16 September 2026.

## Context

Aplikasi Facultative Inward sedang dipindahkan dari Pega ke sistem baru. Kami merekam perilaku
sistem lama dari ekspor rule-nya, bukan merancang ulang — jadi setiap perbedaan perilaku harus
disengaja dan disetujui, bukan muncul diam-diam. Sebagian besar sudah berhasil kami baca dari kode,
termasuk arti keenam nilai status keputusan underwriting (Accept, Reject, Ask, Banding, Decline,
Revise). Yang tersisa adalah hal-hal yang **kodenya memperlihatkan apa yang terjadi, tetapi tidak
memperlihatkan apakah itu memang yang dikehendaki**.

**Tahap ini hanya mencakup New Business dan Renewal.** Beberapa pertanyaan serupa tentang siklus
endorsement sengaja kami tunda sampai fase berikutnya, supaya kuesioner ini tetap pendek dan
terfokus.

Beberapa temuan di bawah mungkin ternyata cacat lama yang tidak pernah terlihat. Bila begitu,
katakan saja — kami tetap akan **meniru perilakunya apa adanya** di tahap pertama agar angkanya
cocok, lalu memperbaikinya sebagai perubahan terpisah yang tercatat. Yang berbahaya bukan cacatnya,
melainkan cacat yang kami perbaiki diam-diam tanpa Anda tahu.

## How to answer

Tenggat: **23 September 2026**. Perkiraan waktu pengisian 20–30 menit.

Isi langsung di bawah tiap pertanyaan. **Jawaban parsial dan "saya tidak tahu" tetap berguna** —
lebih baik ditandai ragu daripada dilewati, karena yang dilewati akan kami perlakukan sebagai
pertanyaan yang masih terbuka dan menahan pekerjaan. Bila sebuah pertanyaan salah premis, katakan
premisnya yang keliru.

---

### 1. Apa beda operasional antara **Reject** dan **Decline**?

Sistem lama punya dua status penolakan yang berbeda, dan memperlakukannya berbeda pula:

- **Reject** membatalkan binding — menghapus konfirmasi binding dan penerimaan R/I slip — dan
  menjadi **syarat munculnya jalur naik banding**: sistem menghitung jumlah penolakan berstatus
  reject untuk menentukan apakah banding tersedia.
- **Decline** tidak punya efek samping sama sekali, dan merupakan **hasil default** — yaitu yang
  dipakai ketika tidak ada kondisi lain yang cocok.

_Why this matters: bila keduanya sebenarnya sama bagi bisnis, kami menyatukannya. Bila berbeda, kami
perlu tahu bedanya untuk memastikan hak banding tidak hilang, atau justru diberikan pada kasus yang
seharusnya final._

Dugaan kami: **Decline = penolakan final tanpa hak banding; Reject = penolakan yang masih dapat
dibanding.** Benar?

> **Benar.**

✅ **Dicatat sebagai K-014** (16 September 2026). Keduanya **tidak disatukan** di sistem baru. Efek
samping `Reject` yang terverifikasi — mematikan konfirmasi binding dan penerimaan R/I slip, serta
menjadi syarat munculnya jalur banding — adalah **perilaku yang dikehendaki**, jadi direproduksi apa
adanya.

---

### 2. Apakah benar **renewal dinilai atas nilai pertanggungan penuh**, bukan atas kenaikannya?

Kami sudah memastikan ini dari kode, dan ingin memastikan **memang begitu kebijakannya**.

Ketika sebuah polis diperpanjang, nilai yang dibandingkan dengan tabel limit wewenang adalah
**nilai pertanggungan penuh polis baru** — bukan selisihnya terhadap polis yang berakhir.

Akibatnya: perpanjangan dengan kenaikan nilai yang kecil sekalipun tetap naik tangga persetujuan
berdasarkan **nilai penuh**, sehingga berpotensi sampai ke pejabat yang jauh lebih tinggi daripada
yang dampak ekonominya sebenarnya menuntut.

_Why this matters: ini menentukan siapa yang harus menyetujui, pada ribuan kasus per tahun. Bila ini
tidak disengaja, memperbaikinya akan mengubah jalur persetujuan secara material — dan itu keputusan
bisnis, bukan keputusan teknis._

Apakah ini kebijakan yang disengaja?

> **Benar — dinilai atas nilai pertanggungan penuh.**

✅ **Dicatat sebagai K-015** (16 September 2026). Konsekuensi yang menjangkau fase berikutnya:
asimetri antara renewal (nilai penuh) dan endorsement (selisih) **juga disengaja** — sehingga
pertanyaan D-3 di `_DITUNDA-fase-endorsement.md` sudah terjawab sebagian sebelum sempat dikirim.

---

### 3. Apakah **banding** memang memicu **fac out**?

Ada satu aturan bernama `IsFacout` — namanya menyiratkan *fac out* — tetapi isinya hanya memeriksa
satu hal: apakah status keputusan bernilai **Banding**. Tidak ada kondisi lain.

Seluruh bukti lain menunjukkan nilai itu memang berarti banding: aturan yang menuliskannya
menyalakan penanda banding dan menetapkan tujuan bandingnya, dan riwayat akseptasi mencatatnya
sebagai `'BANDING'`.

_Why this matters: bila banding memang seharusnya memicu fac out, itu aturan bisnis yang harus kami
pertahankan. Bila aturan ini sekadar salah nama, kami akan menamainya ulang — tetapi perilakunya
tetap kami tiru._

Mana yang benar: (a) banding memang memicu fac out, atau (b) aturan itu salah nama?

> **Tidak memicu fac out; `When\IsFacout` sudah tidak dipakai.**

📝 **Catatan penulisan — bukan perubahan arti.** Kutipan aslinya berbunyi *"tidak memicu banding"*,
sedangkan yang ditanyakan adalah apakah banding memicu **fac out**. Ditulis ulang di atas sebagai
"fac out" agar tidak terbaca terbalik. Kutipan aslinya dicatat di sini apa adanya supaya tidak
hilang, dan **bila tafsir ini keliru, mohon dikoreksi** — konsekuensinya menyentuh 10 rule.

✅ **Dicatat sebagai K-019** (16 September 2026), setelah sempat menjadi **Konflik K-2**. Jawaban
"sudah tidak dipakai" bertabrakan dengan korpus yang menunjukkan **10 berkas masih merujuknya**.
Diselesaikan sebagai: **fitur fac out usang secara bisnis, tetapi kodenya masih aktif dan tetap
diport apa adanya** — penghapusannya adalah perubahan terpisah, bukan bagian migrasi.

---

### 4. Apakah status **"Reject Ceding"** dan **"Ask Ceding"** pernah tercatat pada kolom status keputusan underwriting?

Sistem lama mengenal dua status tambahan — *Reject Ceding* dan *Ask Ceding* — tetapi keduanya hidup
di kolom yang **berbeda** dan tidak pernah menulis ke kolom status keputusan underwriting. Artinya
menurut kode, kolom itu hanya pernah berisi enam nilai: Accept, Reject, Ask, Banding, Decline,
Revise.

_Why this matters: bila benar hanya enam, sistem baru dapat menolak nilai lain sebagai data tidak
sah — yang menangkap kesalahan lebih awal. Bila ternyata di data lama ada nilai lain, penolakan itu
justru akan memblokir kasus yang sah._

Apakah sepanjang pengetahuan Anda kolom status keputusan underwriting pernah berisi nilai di luar
keenam itu?

> **Reject Ceding dan Ask Ceding sudah tidak dipakai.**

✅ **Dicatat sebagai K-021** (16 September 2026). Domain kolom status keputusan underwriting dikunci
ke **enam nilai**: Accept · Reject · Ask · Banding · Decline · Revise.

⚠️ **Tetapi validasinya tidak diterapkan sebagai penolakan diam-diam.** Sistem lama **tidak punya**
validasi domain di kolom ini, jadi menambahkannya adalah **perilaku baru**, bukan migrasi. Karena itu
nilai di luar keenamnya **gagal keras dan tercatat**, bukan ditolak tanpa jejak. Rinciannya di K-021.

---

## Anything else?

Ada hal lain tentang alur persetujuan akseptasi **New Business atau Renewal** yang menurut Anda
perlu kami ketahui, dan tidak kami tanyakan di atas? Termasuk: perilaku yang Anda tahu berbeda
antara yang tertulis di sistem dan yang dipraktikkan sehari-hari.

>
