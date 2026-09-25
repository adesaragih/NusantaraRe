# Laporan akhir — keluaran .docx BRD, ADR, Tiket, Steering

Disusun 25 September 2026. Lingkup: `OUTPUT_HASIL_RNM\keluaran-docx\`.
Korpus `D:\XML\RNM_BRD\` dibaca saja; seluruh penulisan hanya ke `OUTPUT_HASIL_RNM\`.

---

## 1. Yang jadi

**12 berkas `.docx`, 694 halaman**, tiap berkas berdampingan dengan markdown gabungan
sumbernya dan PDF hasil rendernya. Seluruhnya **sudah dilihat sebagai halaman**, bukan
hanya sebagai berkas.

| Dokumen | Halaman |
| --- | ---: |
| ADR-NusantaraRe | 125 |
| BRD-NusantaraRe | 31 |
| Steering-NusantaraRe | 6 |
| Tiket-claim-life | 46 |
| Tiket-claim-nonlife | 77 |
| Tiket-facultative-inward | 60 |
| Tiket-komite | 74 |
| Tiket-life-master | 53 |
| Tiket-life-offer | 55 |
| Tiket-treaty-arrangement | 33 |
| Tiket-treaty-master | 84 |
| Tiket-treaty-realisasi | 50 |
| **jumlah** | **694** |

Berkas: 12 `.docx` (1,2 MB) · 12 `.md` (1,8 MB) · 12 `.pdf` (16,8 MB) · 65 render PNG.

---

## 2. Cacah terukur lawan cacah §0 brief

| Ukuran | Brief §0 | Terukur | Keterangan |
| --- | ---: | ---: | --- |
| modul | 20 | **20** | ✅ cocok |
| ADR | 105 | **105** | ✅ cocok, dua cara cacah sepakat |
| tiket `.scratch` | 188 | **177** | 11 berkas indeks, bukan tiket |
| tiket `dastin` | 117 | **117** | ✅ cocok |
| tiket `jefri` | 60 | **60** | ✅ cocok |
| **tiket seluruhnya** | **365** | **354** | badan tiket yang dimuat |
| Fac In siap / tertahan | 24 / 36 | **24 / 36** | ✅ cocok |

**Cara cacah ADR, dua jalan yang berbeda** — 105 heading keputusan dan 105 baris
konkordansi, selisih himpunan kosong. Jendela: seluruh `ADR-NusantaraRe.md`.

**Asal 105 itu terbukti utuh di korpus**: 29 (`claim-non-prop`) + 4 (`komite-claim-non-prop`)
+ 23 (`treaty-in`) + 42 (`docs\adr`) + 7 (`jefri\OUTPUT FIX\adr`) = **105**. Tidak ada ADR
yang tertinggal.

### 2b. Rincian 354 badan tiket

| Golongan | Cacah |
| --- | ---: |
| badan bernomor angka (`01`, `04a`, `05b`) | 302 |
| badan bernomor huruf milik `jefri` (`R01`–`R08`, `E01`–`E22`, `F01`–`F14`) | 44 |
| **tiket sejati** | **346** |
| berkas papan / urutan / tambahan yang ikut termuat | 8 |
| **jumlah badan dimuat** | **354** |

Delapan berkas itu: 4 `README.md` (papan tiket), 3 `urutan-tiket.md`, dan 1
`TAMBAHAN-TIKET-ronde-2.md`. **Tujuh di antaranya nol kotak-centang** — papan dan urutan,
bukan tiket. Yang kedelapan, `TAMBAHAN-TIKET-ronde-2.md`, memuat **9 kriteria penerimaan**
dan bermuatan tiket sungguhan.

⚠️ `[terbuka]` **Apakah `TAMBAHAN-TIKET-ronde-2.md` dihitung sebagai tiket, dan apakah tujuh
papan itu dikeluarkan dari matriks status?** Isinya saya biarkan utuh di dalam dokumen dan
angkanya saya nyatakan terpisah. Penutupan milik pemilik pekerjaan.

### 2c. Matriks status, dijumlahkan lintas sembilan dokumen

| Siap | Tertahan | needs-info | wontfix | Lain | Jumlah |
| ---: | ---: | ---: | ---: | ---: | ---: |
| 179 | 51 | 1 | 3 | 120 | 354 |

⚠️ **120 di kolom "Lain"** bukan salah label. Modul `Treaty In`, `Treaty In Adjustment`,
`Claim Non Prop`, dan `Komite Claim Non Prop` memakai kosakata status sendiri di berkas
sumbernya — misalnya `status: aktif golongan: pelestarian` — bukan `ready-for-agent`.
Matriks melaporkan apa adanya, tidak menerjemahkan sendiri.

---

## 3. RALAT — koreksi terhadap laporan pemeriksaan saya sebelumnya

Kalimat lama tidak dihapus; ia dikutip di sebelah koreksinya.

### RALAT-1 · Lingkup pemeriksaan rujukan ADR terlalu sempit

> **Lama:** "ADR rujukan TANPA awalan seri: 10"

⛔ Angka itu **hanya untuk `ADR-NusantaraRe.md`**, padahal kriterianya berlaku untuk
keduabelas dokumen. Terukur di seluruh keluaran: **819 rujukan tanpa awalan seri** —
10 di dokumen ADR, 13 di BRD, dan 796 di sembilan dokumen tiket yang memang tidak pernah
melewati tahap pemberian awalan. Pemeriksaan sempit itu memberi kesan bersih yang keliru.

**Sekarang: 0**, terukur pada keduabelas dokumen.

### RALAT-2 · Pemberian awalan seri mengarang rujukan yang tidak ada

> **Lama (kode `beri_awalan`):** "Rujukan `ADR-nnnn` telanjang di dalam badan menunjuk
> keputusan **di serinya sendiri**."

⛔ Andaian itu salah, dan ia sudah masuk ke dokumen yang dikirim. Badan `ADR-D-TI-0039`
tertulis "perluasan **ADR-D-TI-0007**", padahal seri `ADR-D-TI` mulai di 0034 — **ADR-D-TI-0007
tidak ada**. Isi yang dirujuknya justru persis judul `ADR-D-CNP-0007`.

**Sebabnya**: tiga folder dastin memakai **satu ruang nomor bersambung** —
0001–0029 `claim-non-prop`, 0030–0033 `komite-claim-non-prop`, 0034–0056 `treaty-in` —
tanpa celah dan tanpa tumpang tindih. Awalan ditentukan oleh **nomor**, bukan oleh folder
yang mengutip. `docs\adr` (0001–0042) dan `jefri` (0001–0007) adalah ruang nomor tersendiri.

**Terukur: 4 kode, 11 kemunculan** yang menunjuk keputusan tidak ada
(`ADR-D-KCNP-0016`, `ADR-D-TI-0007`, `ADR-D-TI-0016`, `ADR-D-TI-0023`). Seluruhnya sudah
dibetulkan, dan fungsi sumbernya diganti agar tidak mengulang.

**Mengapa 11 itu himpunan yang lengkap** — karena ketiga rentang dastin terpisah rapi,
setiap salah-awalan **pasti** menghasilkan kode yang tidak ada; tidak mungkin ia mendarat
pada kode lain yang kebetulan sah. Jadi daftar hantu = daftar kesalahan.

### RALAT-3 · Cacah tiket saya sempat memberi 298

> **Lama:** "tiket dimuat seluruh dokumen : 298"

⛔ Yang cacat adalah alat ukurnya, bukan dokumennya: regex saya hanya menerima nomor
berupa angka, sehingga melewatkan `04a`, `05b`, `R01`, `E14`, `F09`. Terdamaikan penuh:
untuk **setiap** dokumen, `jumlah heading ##` − 1 (baris `## Matriks status`) = angka
pembangunan. Jumlah **354**, dan itu yang dipercaya.

---

## 4. Pemeriksaan akhir, lingkup keduabelas dokumen

| Pemeriksaan | Hasil |
| --- | --- |
| `.docx` / `.md` / `.pdf` | 12 / 12 / 12 |
| ADR heading keputusan | 105 |
| ADR baris konkordansi | 105 (selisih himpunan: kosong) |
| Rujukan ADR tanpa awalan seri | **0** |
| Rujukan ADR menunjuk keputusan tak ada | **0** |
| Alamat hos internal harfiah | **0** |
| DDL di BRD | `CREATE`=0 · `ALTER`=0 |
| DDL di Steering | `CREATE`=0 · `ALTER`=0 |
| BRD bab bernomor | 28 (bab 1–28) |
| BRD bab modul | 20 (bab 4–23) |
| BRD bab "Kepastian bahan per modul" | ada |
| Steering butir 5a–5e | 5 |
| ADR Lampiran 2 pasangan sebidang | 3 kelompok |
| Badan tiket dimuat | 354 |

**Steering, kata "rekomendasi": 3 kemunculan — ketiganya penyangkalan**, bukan anjuran:
"**Nol rekomendasi** pada bagian 5", "Tiap butir disajikan dengan pilihan dan akibatnya.
**Tidak ada rekomendasi.**", dan "Nol rekomendasi pada bagian 5; nol butir terbuka ditutup."
✅ Nol anjuran sungguhan.

### Halaman yang dilihat sendiri setelah perbaikan

| Halaman | Yang dipastikan |
| --- | --- |
| ADR h088 | `ADR-D-TI-0039` merujuk `ADR-D-CNP-0007`; `ADR-D-TI-0044` seri sendiri tetap utuh |
| ADR h077 | judul `ADR-D-KCNP-0030` tanpa nomor-sendiri yang mubazir |
| ADR h121 | konkordansi rapi, kolom "nomor asli" tetap menyimpan `0030` |
| BRD h004 | bab Claim Life memakai `ADR-U-0007/0008/0001/0003` — seri akar yang benar |
| Tiket-treaty-master h002 | `ADR-D-TI-0040` benar; matriks status terbaca |

---

## 5. Butir terbuka yang tidak saya tutup sendiri

1. ⚠️ `[terbuka]` Apakah `TAMBAHAN-TIKET-ronde-2.md` (9 kriteria) dihitung tiket, dan apakah
   tujuh papan/urutan dikeluarkan dari matriks status. → pemilik pekerjaan.
2. ⚠️ `[terbuka]` Apakah tiket NB Treaty In **24–28** dinyatakan digantikan oleh sebelas
   tiket penyimpanan EDM; keduanya menutup 58 AC yang sama. → pemilik pekerjaan.
3. ⚠️ `[terbuka]` Kosakata status modul dastin (`aktif`/`pelestarian`) belum dipetakan ke
   kosakata triase; 120 tiket karenanya jatuh di kolom "Lain". → pemilik pekerjaan.

---

## 6. Jebakan baru untuk katalog

**Jebakan #8 — nomor ADR dapat berbagi satu ruang lintas folder.** Memberi awalan seri
menurut folder yang mengutip terasa benar dan menghasilkan kode yang *kelihatan* sah.
Ujinya murah dan wajib: **setiap kode hasil penulisan ulang harus ada di daftar kode yang
benar-benar ada**. Tanpa uji itu, 11 rujukan palsu lolos ke dokumen yang sudah dirender,
sudah dilihat, dan sudah dinyatakan bersih.

**Jebakan #9 — kriteria berlaku untuk seluruh keluaran, pemeriksaan sering hanya satu berkas.**
Angka "10" benar untuk satu dokumen; untuk keluarannya angka yang benar **819**.

---

## 7. TELEMETRI EKSEKUSI

### Terukur

| Besaran | Nilai | Cara ukur |
| --- | ---: | --- |
| Berkas `.docx` | 12 | hitung berkas |
| Halaman seluruhnya | 694 | `pymupdf.page_count` tiap PDF |
| Render PNG | 65 | hitung berkas |
| ADR dimuat | 105 | heading + konkordansi, dua cara sepakat |
| Badan tiket dimuat | 354 | heading `##` − baris matriks, per dokumen |
| Rujukan ADR diperbaiki | 830 | 819 telanjang (811 diberi awalan + 8 nomor-sendiri dibuang) ditambah 11 hantu dibetulkan |
| Alamat ditopeng | 1 | regex alamat pada markdown gabungan |
| Lama render PDF | 72 detik | jam dinding di dalam skrip render |

### ⛔ Tidak diukur

| Besaran | Sebab |
| --- | --- |
| Token terpakai | Angka token sejati tidak terlihat dari dalam sesi |
| Biaya | Turunan token; tidak diukur |
| Lama sesi, jam dinding | Tidak dicatat |
| Jumlah panggilan alat | Tidak dicatat; sesi ini sempat diringkas sehingga cacah dari dalam tidak bisa dipercaya |

**Pengukuran luar (`claude --print --output-format json`) tidak dijalankan**, sebab ia
memulai sesi terpisah yang ongkosnya bukan ongkos ronde ini. Angka mana pun yang berasal
dari sana akan mengukur hal yang lain.
