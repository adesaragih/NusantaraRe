# Ringkasan penggalian kebutuhan — Treaty In

**Tanggal:** 22–23 September 2026
**Sumber:** ekspor aturan Pega `D:\XML_NURE\Treaty In`, ditambah 5 stored procedure dan 7 DDL Oracle
**Metode:** penggalian kebutuhan satu pertanyaan per giliran, 15 pertanyaan pokok

---

## SEBERAPA JAUH DOKUMEN INI BISA DIPERCAYA

Dibaca lebih dulu, sebelum isinya.

> **Diperbarui 23 September 2026.** Pemilik proses memutuskan **seluruh butir tertahan diputuskan
> sekarang**, dan sesi penulisan spesifikasi berjalan penuh tanpa menunggu uji maupun wawancara.
> Angka di bawah adalah keadaan **sebelum** keputusan itu. Keadaan terkini — 55 keputusan, dan
> berapa di antaranya yang akan berubah bila uji kelak berbeda — ada di halaman muka
> `KEPUTUSAN-TANPA-VERIFIKASI.md`, dan **itu yang berlaku**.

| | Jumlah (keadaan sebelum 23 Sep sore) |
|---|---|
| Keputusan mengikat yang tercatat | **34** |
| Bersandar pada **fakta yang sudah diverifikasi** dari sumber | **26** |
| Bersandar pada fakta **ditambah anggapan bahwa sebuah uji kembali bersih** | **2** |
| Masih bersandar pada **hipotesis** yang menunggu uji data atau wawancara | **6** |
| Uji data yang disiapkan dan belum kembali | **20** (Uji A–V, dikurangi K dan Q) |
| Uji yang **dianggap bersih** atas keputusan pemilik proses | **2** (Uji K dan Q) |
| Pertanyaan wawancara yang **belum satu pun dijawab** | **15** |

**Dua keputusan yang berdiri di atas anggapan uji bersih** — perlakuan penanda kegagalan Rate on
Line, dan lingkup perbaikan kurs — ada di `KEPUTUSAN-TANPA-VERIFIKASI.md` §F, lengkap dengan apa yang
**tidak** ikut menjadi bersih. Kalimat yang berlaku untuk keduanya: **bersih berarti belum
merugikan, bukan berarti benar.**

**Enam keputusan yang masih berupa hipotesis**, dan harus dibaca sebagai asumsi bertanda:

| Keputusan | Menunggu |
|---|---|
| Materialitas addendum diturunkan dari ada-tidaknya baris selisih | Uji A, B, C — apakah dua penanda addendum memang dua sumbu |
| "Paling banyak dua mata uang per layer" sebagai aturan validasi | Uji I, J |
| Selisih selalu bertanda dan boleh negatif; lantai milik perhitungan premi | Uji L + wawancara keuangan |
| Kedua mekanisme penyebaran sah dan dipertahankan | Uji M (yang menunggu Uji E) |
| Bentuk terbitan ke hilir diturunkan saat dibaca | Uji F, G |
| Tingkat pencatatan EPI adalah 100% treaty | wawancara underwriting |

Selebihnya — termasuk seluruh 15 ADR — bersandar pada perilaku yang dibaca langsung dari aturan
sistem lama dan tidak berubah oleh hasil uji mana pun.

**Peringatan yang berlaku untuk seluruh dokumen ini dan turunannya:** di mana pun sebuah aturan
sistem lama disebut, dokumen ini memisahkan **yang tertulis di aturan** dari **yang benar-benar
berjalan**. Banyak yang tertulis tidak berjalan. Jangan membaca daftar aturan sebagai daftar
penjaga yang ada.

---

## Berkas lain dalam paket ini

| Berkas | Isinya |
|---|---|
| `CONTEXT.md` | glosarium, batas konteks, aturan kerja, batas lingkup |
| `docs/adr/0034`–`0048` | 15 keputusan arsitektur, satu berkas per keputusan |
| `PENGETAHUAN.md` | pembedahan AS-IS lengkap, utang teknis TD-01..TD-22 |
| `PERMINTAAN-DBA-1-UJI-A-SAMPAI-H.sql` | Uji A sampai V (nama berkas sengaja tidak diubah) |
| `DAFTAR-ESKALASI-MANAJEMEN.md` | untuk direksi, tanpa istilah teknis |
| `KEPUTUSAN-TANPA-VERIFIKASI.md` | apa yang tertahan dan menunggu apa |
| `INVENTARIS-STRUKTUR-DATA.md` | bahan untuk sesi model data berikutnya |

---

## Lingkup yang ditetapkan

**Keempat belas fungsi ikut dimigrasi, bertahap, dan gelombang adalah URUTAN MERILIS — bukan
urutan membangun.** Menunda seluruh manfaat sampai fungsi terakhir selesai membuat gelombang
kehilangan gunanya dan menumpuk risiko pada satu hari.

| Gelombang | Isinya |
|---|---|
| **G1** | pencatatan kontrak, addendum, persetujuan, perhitungan kedua cabang, penyebaran internal, installment, jadwal pelaporan, master, peran dan penugasan, jejak perubahan |
| **G2** | penyebaran ke pihak luar, integrasi hilir |
| **G3** | pencapaian, produksi, pelaporan, ekspor |

Di luar lingkup: tombol pengembang; dua tambalan lama yang diputuskan tidak dibawa; **Treaty In
Adjustment** (embargo, hanya *seam*-nya yang disiapkan); area C (perilaku layar); area F sisi masuk.

---

# Temuan dan keputusan per area

## Area A — Batas bisnis dan konteks

**Yang ditemukan.** Modul ini **tidak punya case type Pega sama sekali**: tidak ada Flow, tidak ada
alur kerja. Statusnya dikarang sendiri lewat properti dan sebuah data transform berisi nama orang.

Tiga konteks di luar Treaty In ditemukan dan **belum punya pemilik**: program retro NuRe sendiri
(tabel acuan kapasitas yang tidak ditulis satu pun aturan Pega), penulis tabel produksi dan
pencapaian, dan layanan dokumen bersama yang dipakai Treaty In **dan** Klaim.

**Keputusan.** Batas konteks ditulis di `CONTEXT.md`. Treaty In menjadi **pembaca** tabel acuan
kapasitas lewat adapter yang **gagal keras** di perbatasan — bukan pemiliknya, dan bukan pula
menerima nilai apa adanya (ADR-0035).

## Area B — Siklus hidup

**Yang ditemukan.** Lima properti tumpang tindih menjawab pertanyaan yang sama — `Position` (5.634
rujukan), `ViewState` (661), `StatusAkseptasi` (165), `IsEditData` (125), `RevisionState` (21).
`ViewState` disetel imperatif di lima tempat, termasuk satu aturan bernama `TreatyInForceEdit`, dan
diperiksa dengan **delapan ejaan berbeda**.

Jalur revisi memendekkan persetujuan dari empat tingkat menjadi dua, dan pemilihannya lewat tombol.

**Keputusan.** Satu keadaan siklus hidup bertipe tegas; "boleh diubah" turunan yang dihitung saat
ditanya (**ADR-0046**). Delapan ejaan tidak direkonsiliasi — tidak ada yang diangkut.

## Area C — Layar dan perilaku dinamis

**Dicoret dari penggalian.** Perilaku layar tertulis lengkap di definisi Section dan direkonstruksi
saat menulis spec. Menggalinya lewat tanya-jawab memboroskan giliran.

Satu hal dari area ini tetap dibawa sebagai temuan: dari 990 sel yang bisa diisi, **104 properti
tidak punya penjaga apa pun** pada jalur revisi — termasuk persentase bagian dan persentase retro.
Pemeriksaan lima belas menit di aplikasi berjalan sudah dicatat di `KEPUTUSAN-TANPA-VERIFIKASI.md` §C.

## Area D — Model data

**Yang ditemukan.** Halaman kontrak menyimpan **bentuk yang sama sampai empat kali**:
`ValueDifference` (161 jalur anak, 90 identik dengan pohon utama), `ActualValue` (107 / 77),
`OLDDATA`, dan `ValueBeforeProrate`.

Bentuk kanoniknya adalah dokumen JSON; tabel relasional hanya proyeksi datar.

**Keputusan.** Keempat pohon paralel bukan empat jenis benda — satu bentuk dalam beberapa peran,
diselesaikan lewat **versi** dan **turunan** (ADR-0036, ADR-0048). Rincian di
`INVENTARIS-STRUKTUR-DATA.md`.

## Area E — Aturan bisnis dan perhitungan

Area terbesar. Lima pola berulang ditemukan, dan masing-masing ditutup satu ADR sehingga instans
berikutnya tidak perlu diperdebatkan lagi.

### E.1 Kegagalan disamarkan menjadi nilai → **ADR-0035**

Lima instans: Rate on Line diisi angka penanda saat limit nol; mata uang di luar slot menghasilkan
nol; kurs valas hilang **menihilkan** nilai; prosedur Oracle gagal tanpa mengisi kode galat;
persentase tidak terurai dari teks bebas.

### E.2 Angka yang disetujui bergeser sendiri → **ADR-0036**

Empat sebab berdiri sendiri: penyebaran disusun ulang dari acuan pada setiap perhitungan; baris
acuan bisa disunting kapan saja; kurs memakai tahun berjalan; bagian NuRe memakai bagian hari ini
atas jumlah masa lalu. **Nilai yang pernah disetujui tidak bisa direproduksi** — tidak ada kolom
waktu sama sekali di tabel kontrak maupun addendum.

### E.3 Masukan versus turunan → **ADR-0037**

Tiga instans: penyebaran otomatis-versus-manual dipilih oleh keterisian sebuah field; reinstatement
dengan empat field yang saling memicu dalam cincin; kapasitas proporsional dengan sakelar
`autocalculate` yang di **12 dari 18** titik pemanggilan diberi makan nomor layer.

Turunannya: **prorata mengenai dasarnya, bukan hasilnya** — dua puluh blok prorata terpisah menjadi
satu perkalian, dan asimetri "potongan tidak ikut diprorata" lenyap bersama bentuk yang
melahirkannya.

### E.4 Tingkat pencatatan → **ADR-0039** (perluasan ADR-0007)

Rumus pencapaian menaikkan premi ke tingkat 100% dengan membagi bagian NuRe, lalu membandingkannya
dengan EPI. Dua besaran disimpan pada tingkat berbeda, didamaikan oleh pembagian yang tidak
tertulis di mana pun. **Tingkat adalah sifat besarannya, bukan konvensi.**

### E.5 Aturan sebagai data → **ADR-0038**

Perutean persetujuan, kontrak baca adapter, dan kelengkapan transisi — ketiganya berubah tanpa
mengubah arti bendanya. Dibedakan tegas dari **invarian**, yang bagian dari arti benda dan
ditegakkan model.

### Keputusan domain yang spesifik

| Hal | Keputusan |
|---|---|
| Kolom kembar mata uang (`Limit`/`Limit2`) | **tidak dibawa**. `CurrencyRelation` bertahan sebagai atribut bernama; "maksimal dua mata uang" jadi aturan validasi |
| Deposit vs minimum premium | **dua besaran berbeda**, tidak wajib sama. Mesin selisih **selalu mengurangi dan selalu bertanda**; lantai adalah sifat besarannya, bukan sifat perhitungan selisihnya |
| Reinstatement | **dua masukan** dinegosiasikan terpisah; hubungannya **perkalian**. Penamaan sistem lama tertukar: persen yang menggerakkan pemulihan limit dinamai "Additional", yang menggerakkan premi dinamai "Reinstatement" |
| Kapasitas Quota Share | satu masukan (persentasenya); retensi dan sesi turunan |
| Kapasitas Surplus | **dua masukan**: retensi sebagai jumlah uang, dan jumlah lines |
| `CessionPct` | **pecah tanpa syarat** — di Surplus isinya jumlah lines × 100, bukan persentase apa pun |
| Potongan | dibedakan menurut **dasar perhitungan**: premi bruto / laba / jumlah yang disepakati. Profit commission tidak bisa dihitung sebelum periode selesai |

### Temuan bernomor baru

- **TD-22** — perhitungan menimpa **tarif premi reinstatement** dengan rasio limit, sehingga syarat
  kontrak terhapus oleh perhitungan. Tak terlihat karena keduanya disemai 100.
- **TD-02 direvisi** — angka penanda yang tersimpan adalah `998999800`, bukan `9989998`; ditambah
  penyebut yang berlipat menurut jumlah mata uang EGNPI, dan penggandaan kedua di bawah relasi AND.
  Inti perhitungannya diganti **di sistem produksi pada 15 April 2025**, yang menciptakan batas era
  di dalam data.
- **TD-03 dicabut sebagian** — klaim "kasus tepat dua mata uang tidak tertangani" dan
  "ambang tidak konsisten" **salah**; keduanya berasal dari membaca label, bukan kondisi.
- **TD-06 berubah arah** — kurs valas hilang **menihilkan** nilai, bukan menyamakannya dengan
  rupiah. Kesalahan yang menghilangkan jauh lebih mungkin lolos tanpa disadari daripada yang
  menggelembungkan.

## Area F — Integrasi

**Yang ditemukan.** Tidak ada satu pun pembaca tabel terbitan di dalam ekspor Pega — konsumennya di
luar. Pencocokan antar-sistem memakai pemotongan **tujuh karakter pertama** nomor offer, yang
terikat pada format ID `kode_situs || nomor urut enam digit`.

**Keputusan.** Bentuk terbitan ke hilir **diturunkan saat dibaca**, bukan disisipkan setiap kali
menyimpan (ADR-0023). Keterikatan pada format ID diperlakukan sebagai **pekerjaan bermilik** dalam
rencana migrasi, bukan sebagai temuan: sebabnya diputuskan di gelombang 1, akibatnya jatuh di
gelombang 3.

Sisi masuk (penyimpanan objek, pembacaan master) dicoret — tertutup keputusan layanan dokumen
bersama dan keputusan adapter.

## Area G — Proses terjadwal

**Yang ditemukan, dan ini satu kalimat:**

> Proses terjadwal di modul ini bukan proses. Ia daftar tanggal di dalam kontrak.

Tidak ada satu pun agent, queue processor, correspondence, atau pengiriman email di seluruh modul.
Tiga jatuh tempo bertingkat dihitung untuk setiap periode pelaporan — penyerahan, konfirmasi,
penyelesaian — disimpan, ditampilkan, dan **tidak ada apa pun yang mengawasinya**.

**Keputusan.** Keterlambatan **dihitung saat dibaca**, tidak pernah disimpan sebagai penanda —
instans langsung ADR-0036. Proses latar ditunda, dan **batasnya ditulis sekarang**: bila kelak
dibangun, ia mencatat peristiwa tentang kontrak dan **tidak pernah mengubah kontraknya**.

## Area H — Notifikasi

Kosong di sistem lama; tertutup bersama area G.

## Area I — Hak akses

**Yang ditemukan.** **Tidak ada satu pun nama hak akses yang terisi** di seluruh ekspor — mekanisme
peran bawaan Pega hadir 2.202 kali sebagai tempat kosong dan tidak pernah dipakai. Yang mengatur
siapa boleh apa hanya bendera pada kontraknya dan nama orang yang ditanam di dalam aturan.

Itu bukan keputusan bahwa wewenang milik kontrak. Itu sistem yang tidak pernah punya cara
menyatakan "siapa".

**Keputusan (ADR-0044).** Keadaan menjawab *apakah tindakan ini mungkin sekarang*; peran menjawab
*apakah orang ini boleh melakukannya*. **Setiap larangan punya tepat satu sebab** — dan itu uji
rancangan, bukan pedoman menulis pesan. Peran dan penugasan **bertanggal** masuk gelombang 1.

## Area J — Jejak audit dan pelaporan

**Sisi pelaporan ditutup**: seluruh Report Definition adalah pencarian master; konsumen sebenarnya
ada di hilir (Uji F); ekspor CSV/Excel adalah sifat komponen tabel.

**Sisi jejak audit terbuka dan naik ke gelombang 1.** Satu-satunya jejak adalah daftar komentar,
ditulis **hanya pada perpindahan persetujuan**. Penyimpanan biasa — termasuk yang mengubah angka —
tidak meninggalkan jejak siapa maupun kapan.

Alasan kenaikannya bukan kepatuhan: **jejak perubahan adalah bukti bahwa pembekuan benar-benar
terjadi.** Tanpa jejak, pembekuan adalah janji, dan selisih yang dibukukan Adjustment berdiri di
atas janji (ADR-0045).

## Area K — Migrasi data dan masa transisi

**Keputusan (ADR-0042).** Sejarah **tidak dibersihkan**. Membersihkannya agar memenuhi invarian
berarti mengubah angka yang pernah disetujui — persis penyakit yang seluruh sesi ini dipakai untuk
mendiagnosis. Kontrak dan addendum pindah **lengkap dan apa adanya**, ditandai **warisan**, dengan
**aturan sentuh-perbaiki**: begitu disentuh, ia harus memenuhi invarian saat itu juga.

Materialitas addendum historis diambil apa adanya sebagai nilai warisan; aturan penurunan baru
tidak dijalankan atasnya, karena itu akan mengklasifikasi ulang sejarah.

**Keputusan (ADR-0041).** Koeksistensi per fungsi, di bawah **satu fakta, satu penulis, selalu**,
ditegakkan lewat **pencabutan hak tulis di basis data** — bukan lewat kebijakan. Langkah verifikasi
"tidak ada penulisan lagi" adalah **alat penemuan**: aturan Pega yang gagal di situ menemukan
penulis yang tidak kita ketahui.

**Keputusan (ADR-0043).** Migrasi memindahkan angka apa adanya — kriteria **identik**, nol
toleransi. Menghitung ulang adalah **peristiwa bisnis tersendiri** — kriteria berbeda dan bisa
dijelaskan. Uji paritas dijalankan sebagai **laporan selisih, tidak menulis ke mana pun**, dan
daftar harapan selisih ditulis **sebelum** run dijalankan.

Kemungkinan ada **tiga generasi data**, bukan dua: sebagian data sekarang sendiri hasil migrasi
sebelumnya.

## Area L — Non-fungsional dan kriteria penerimaan

Ukuran keberhasilan migrasi: jumlah kontrak dan addendum cocok; setiap ID lama punya pasangan
**dan** setiap ID baru punya asal; nilai kunci dipindahkan tanpa dihitung ulang lalu dibandingkan;
dan **jumlah baris anak per kontrak cocok** — kegagalan yang paling sering terjadi pada migrasi
dari dokumen JSON ke tabel relasional.

Data uji (ADR-0047): **identitas boleh disamarkan, angka tidak boleh disentuh sama sekali**. Dua
kumpulan terpisah — kurasi untuk pekerjaan harian, salinan paritas untuk run paritas di lingkungan
terkendali. Menyamarkan nama **tidak** menganonimkan data treaty; perlindungan nyata ada pada
kendali akses.

---

# Perilaku Pega: dipertahankan atau dibuang

Pertanyaan yang diajukan untuk **setiap** perilaku sistem lama: *ini requirement bisnis yang harus
dipertahankan, atau konsekuensi dari cara Pega bekerja dan boleh hilang?*

| Perilaku sistem lama | Nasib | Alasan |
|---|---|---|
| Status dikarang lewat lima properti | **dibuang** | konsekuensi ketiadaan case type; satu keadaan siklus hidup menggantikannya |
| Bendera "boleh diubah" yang disimpan | **dibuang** | turunan yang disetel oleh klik |
| Empat pohon paralel dalam satu dokumen | **dibuang** | satu bentuk dalam beberapa peran; jadi versi dan turunan |
| Kolom kembar mata uang | **dibuang** | bentuk data, bukan aturan bisnis |
| `CessionPct` merangkap dua arti | **dibuang** | dua besaran berbeda dipaksa berbagi kolom |
| `RetentionPct` di Quota Share | **dibuang** | definisi, bukan kesepakatan |
| Cincin persen–jumlah reinstatement | **dibuang** | pintu masuk suntingan atas turunan |
| Sakelar `autocalculate` | **dibuang** | tidak ada sakelar yang mematikan perhitungan |
| Angka penanda kegagalan | **dibuang** | kegagalan bukan nilai |
| Dua puluh blok prorata terpisah | **dibuang** | prorata mengenai dasarnya |
| Jalur otomatis menimpa suntingan pada setiap simpan | **dibuang** | melanggar pembekuan |
| Kurs tahun berjalan | **dibuang** | melanggar pembekuan |
| Pengecualian FAC-OUT lewat penghapusan jejak | **dibuang** | kebijakan diterapkan di satu cabang saja; sisa data tertinggal |
| Dua mekanisme penyebaran | **dipertahankan bersyarat** | sah bila Uji M menunjukkan manual adalah pilihan, bukan kejatuhan. Pemicunya jadi tindakan eksplisit |
| Deposit dan minimum premium sebagai dua besaran | **dipertahankan** | ketentuan kontrak yang nyata |
| Dua persentase reinstatement | **dipertahankan** | dua hal yang dinegosiasikan terpisah; hanya penamaannya diperbaiki |
| Retensi Surplus sebagai jumlah uang per treaty group | **dipertahankan, dinaikkan jadi masukan** | struktur sudah ada, statusnya yang belum diakui |
| Kurs beku per periode di dalam kontrak | **dipertahankan, diperluas** | bentuknya sudah benar; pengisinya yang tidak ada |
| Jatuh tempo bertingkat pelaporan | **dipertahankan** | syarat kontrak yang nyata |
| Blokir revisi bila sudah ada klaim | **dipertahankan** | satu-satunya aturan duplikat yang masih hidup, dan ia masuk akal |
| Addendum sebagai versi | **dipertahankan, diubah bentuk** | dari tabel terpisah berkolom perangkai menjadi riwayat versi |
| Empat tingkat persetujuan | **dipertahankan, dipindahkan jadi data** | nama orang tidak pernah muncul di aturan |

**Kemampuan baru — bukan pelestarian, jangan dihitung gratis:** deteksi duplikat, daftar wajib-isi,
penolakan daftar hitam, entitas peran dan penugasan bertanggal, jejak perubahan, penandaan
keterlambatan, penunjuk asal-usul, dan pengisi kurs beku.

---

# OPEN QUESTION

| # | Pertanyaan | Kenapa penting | Pemilik | Bila tidak terjawab |
|---|---|---|---|---|
| 1 | EPI dicatat pada tingkat 100% atau bagian NuRe? | menentukan seluruh skala angka pencapaian | Underwriting | angka pencapaian bisa salah lima kali lipat |
| 2 | Apa arti sebenarnya dua penanda addendum? | bentuk klasifikasi addendum | Underwriting + Uji A/B/C | klasifikasi salah sejak awal |
| 3 | Siapa pemilik tabel acuan kapasitas? | tanpa pemilik, adapter tak punya lawan bicara | Manajemen | kapasitas bersumber dari data tak bertuan |
| 4 | Siapa pembaca tabel terbitan di hilir? | bentuk kontrak antar-sistem | Uji F/G | integrasi dirancang untuk konsumen yang salah |
| 5 | Adakah dokumen batas wewenang persetujuan? | bila ada dan tak pernah dibaca sistem — temuan kepatuhan | Manajemen + Kepatuhan | risiko kepatuhan tak terukur |
| 6 | Bisakah pengembalian premi dibukukan, lewat dokumen apa? | menentukan apakah jepitan negatif tambalan akuntansi | Keuangan | aturan selisih salah bentuk |
| 7 | Kurs apa yang dipakai buku besar? | bila beda dari modul ini, keduanya sudah lama tak sepakat | Keuangan | temuan tersendiri terlewat |
| 8 | Penyebaran manual: pilihan atau kejatuhan? | menentukan apakah dua jalur dipertahankan | Uji M + Underwriting | jalur mubazir terbawa, atau jalur sah terbuang |
| 9 | Apakah reinstatement NuRe pernah pro rata **waktu**? | faktor ketiga hilang dari sistem lama | Underwriting | premi kurang/lebih tagih terus berlanjut |
| 10 | Porsi limit dipulihkan: syarat, atau turunan dari gerusan klaim? | menyentuh batas konteks dengan modul Klaim | Underwriting | model reinstatement salah sisi |
| 11 | Apakah riwayat lintas tahun pernah diminta? | menggerbangi "rangkaian treaty" | Underwriting | benda dimodelkan padahal tak ada, atau sebaliknya |
| 12 | Siapa pemilik layanan dokumen bersama? | Treaty In dan Klaim sama-sama pemakai | Manajemen | dua modul membangun hal yang sama |
| 13 | Kelengkapan apa yang dituntut tiap transisi? | bentuk daftar kelengkapan per transisi | Underwriting + bagian teknik | validasi salah tempat |
| 14 | Masa simpan arsip JSON sistem lama | ADR-0034 menuntutnya | Manajemen + Kepatuhan | arsip menganggur tanpa batas |
| 15 | Dasar perhitungan tiap nama potongan | langsung jadi tabel di model baru | Underwriting + Uji U-1 | potongan salah diprorata |

---

# Risiko migrasi

| Risiko | Kenapa | Peredamnya |
|---|---|---|
| **Kapasitas proporsional mungkin tak pernah dihitung** pada sebagian kontrak | sakelar perhitungan diberi makan nomor layer di 12 dari 18 titik | Uji T; bila angkanya besar, naik ke daftar eskalasi |
| **Batas era 15 April 2025** di dalam data | inti perhitungan Rate on Line diganti di produksi | Uji P dipilah per era |
| **Tiga generasi data**, bukan dua | sebagian data sudah hasil migrasi sebelumnya | Uji R-3; penanda warisan menyebut generasi |
| **Tidak ada kolom waktu** di tabel kontrak | menyulitkan hampir setiap pemilahan historis | pendekatan lewat jejak komentar, batasnya disebut apa adanya |
| **Kumpulan data uji tertahan** sampai Uji A–V kembali | uji-uji itu yang menghasilkan daftar ID-nya | jadwalkan uji lebih dulu; jangan susun tangan |
| **Model disusun dari sumber yang tidak lengkap** | inventaris hanya menangkap jalur berawalan literal | verifikasi ke `JSONDATA` produksi sebelum model dikunci |
| **Invarian baru bertabrakan dengan data lama** | sebagian kontrak historis melanggarnya | ADR-0042: warisan + sentuh-perbaiki, tanpa pembersihan massal |
| **Pembekuan belum berlaku saat Adjustment dibangun** | seam Adjustment mensyaratkannya | ADR-0036 dinaikkan jadi prasyarat, bukan kebersihan |

---

*Disusun 23 September 2026. Seluruh berkas paket ini berada di `D:\XML_NURE\_migration-docs\treaty-in`.*
