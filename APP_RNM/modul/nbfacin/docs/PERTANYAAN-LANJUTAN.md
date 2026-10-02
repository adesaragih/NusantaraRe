# Pertanyaan dari pengerjaan lanjutan tiket NB — 1 Oktober 2026

**Untuk:** work owner (meneruskan ke DBA / Underwriting / pemilik aplikasi Pega bila perlu) ·
**Dari:** pengerjaan tiket NB-02 … NB-16 · **Status kode:** tiap butir di bawah sudah punya sikap
sementara yang aman (menolak, bukan menebak), jadi pekerjaan tidak berhenti menunggu jawabannya.

---

## 1. Satu kasus Layering nyata — DBA / work owner

**Yang kami minta:** satu kasus New Business yang coverage-nya ber-`LayerList` **terisi**, dengan
`.TSI`, `.Rate`, `pyWorkPage.OfferFacIn.ProRatePercent`, `param.Layer` (jumlah lapisan), dan premi
lapisan + premi coverage hasil sistem lama.

**Mengapa:** `[terverifikasi]` seluruh `LayerList` di 115 berkas `DDL\CONTOH` kosong. Rumus Layering
(NB-06) dan akumulasinya di dalam loop (NB-07) hanya teruji dengan contoh hitung.

> 💡 **Rekomendasi agent (1 Oktober 2026)** — bukan jawaban; keputusannya tetap pada penerima.
>
> Jangan minta contoh satu per satu: minta **satu query hitungan** produksi untuk butir 1–3 sekaligus
> (kasus dengan `LayerList` terisi; `CalculateMethod_FacIn = 2`; `DiscountType = Percent`; `Loading`
> terisi; pro-rata coverage ≠ 100). Hitungan nol = jalur tidak terpakai di produksi → cukup ditandai
> "tidak terpakai", tanpa menunggu kasus. Hitungan tidak nol → minta satu contoh per jalur.

> ✅ **Jawaban work owner (1 Oktober 2026)** — `KEPUTUSAN-30-09-2026.md` butir 30
> **Layering tidak dipakai lagi di sistem baru** — permintaan ini dicabut; kode Layering dihapus.


## 2. Kasus PA metode 2 dan kasus PA berdiskon persen — DBA / work owner

**Yang kami minta:** satu kasus PA ber-`CalculateMethod_FacIn = 2` (short period, dengan
`.PctShortPeriod`), dan satu kasus PA dengan `DiscountType = Percent`.

**Mengapa:** kelima kasus PA yang ada bermetode 1 dan 3 tanpa diskon. Rumus L858 dan langkah diskon 2
(NB-04) hanya teruji dengan contoh hitung.

> 💡 **Rekomendasi agent (1 Oktober 2026)** — bukan jawaban; keputusannya tetap pada penerima.
>
> Lihat rekomendasi butir 1 — satu query hitungan yang sama.

> ✅ **Jawaban work owner (1 Oktober 2026)** — `KEPUTUSAN-30-09-2026.md` butir 35
> Sesuai rekomendasi: satu query hitungan ke DBA (tanpa jalur Layering). Belum ada hasil.


## 3. Kasus MBU yang memicu pembulatan — DBA / work owner

**Yang kami minta:** satu kasus MBU dengan `.Loading` terisi, atau dengan `.ProRatePercent` coverage
selain 100, atau yang hasil antaranya berdesimal lebih dari empat.

**Mengapa:** `[terverifikasi]` 93 dari 93 baris coverage MBU nyata cocok eksak — tetapi tidak satu pun
memicu pembulatan, dan pro-ratanya selalu 100. Bentuk bersarang L1144 terbukti hanya lewat contoh hitung.

> 💡 **Rekomendasi agent (1 Oktober 2026)** — bukan jawaban; keputusannya tetap pada penerima.
>
> Lihat rekomendasi butir 1 — satu query hitungan yang sama.

> ✅ **Jawaban work owner (1 Oktober 2026)** — `KEPUTUSAN-30-09-2026.md` butir 35
> Sesuai rekomendasi: satu query hitungan ke DBA. Belum ada hasil.


## 4. Daftar kunci BUANG de-identifikasi — work owner (tinjauan manual)

**Yang kami minta:** keputusan apakah **sembilan medan** berikut masuk daftar BUANG, dan izin menyimpan
fixture hasil de-identifikasi di repositori sesudahnya.

`pxCreateOperator` · `AccumulationDescription` · `AccumulationCode` · `PIC` · `PICSuggest` ·
`CommentSuggest` · `PropertiItemNote` · `SobName` · `TopRiskLocation`

**Mengapa:** `[terverifikasi]` alat NB-15 dijalankan atas kelima `DDL\P-5 *.txt`. Nilai yang sudah
dibuang (nama tertanggung, alamat, nama operator, dll.) **muncul lagi sebagai kata utuh** di kesembilan
medan itu. Tiket 15 melarang filter pola, jadi penambahan wajib lewat keputusan Anda. Selain itu ada
48–128 nama kunci bernilai teks per berkas yang dapat dicetak alat sebagai daftar tinjauan (nama kunci
saja, tanpa nilai).

> 💡 **Rekomendasi agent (1 Oktober 2026)** — bukan jawaban; keputusannya tetap pada penerima.
>
> **Masukkan kesembilannya ke daftar BUANG** — isinya memang identitas: `PIC`/`PICSuggest` nama orang;
> `pxCreateOperator` login operator; `AccumulationDescription`/`AccumulationCode`/`TopRiskLocation` alamat;
> `CommentSuggest`/`PropertiItemNote` teks bebas yang memuat nama; `SobName` nama sumber bisnis. Tambahkan
> juga pasangan metadatanya (`pxUpdateOperator`, `pxUpdateOpName`) bila muncul. Lalu jalankan alat dengan
> `-kandidat`, tinjau daftar nama kuncinya **sekali**, dan simpan fixture di repositori hanya setelah alat
> melaporkan **nol kebocoran** untuk kelima berkas.

> ✅ **Jawaban work owner (1 Oktober 2026)** — `KEPUTUSAN-30-09-2026.md` butir 32
> **Setuju** — kesembilannya masuk `DaftarBuang` (21 kunci). Dijalankan ulang: kebocoran
> **2 / 2 / 1 / 0 / 6**, kini di dua medan baru — lihat butir 7. Fixture belum disimpan.


## 5. Dari mana `pyWorkPage.Quotation` — pemilik aplikasi Pega / DBA

**Yang kami tanyakan:** predikat `IsPA` membaca `pyWorkPage.Quotation.BusinessType`, sedangkan berkas
kasus menyimpan `QuotationData.BusinessType` di bawah halaman OfferFacIn. Apakah `pyWorkPage.Quotation`
dan `pyWorkPage.OfferFacIn.QuotationData` selalu berisi sama? Dan apa arti `StatusBusiness = 3`, yang
menggerbangi seluruh pemilihan rumus premi (`CountGrossPremi_Act` langkah 5)?

**Mengapa:** tiket 10 memperingatkan jalur baca ganda pada properti lini bisnis. Kode kini membaca
keempat jalur apa adanya dan mencatat bila berbeda — tetapi tanpa `pyWorkPage.Quotation`, kasus PA dari
berkas kasus tidak dapat dikenali lininya.

> 💡 **Rekomendasi agent (1 Oktober 2026)** — bukan jawaban; keputusannya tetap pada penerima.
>
> `StatusBusiness = 3` kemungkinan besar **Endorsement** `[dugaan]`: langkah 5 tangga akseptasi memakai
> **selisih** TSI hanya saat nilainya 3, sama dengan aturan endorsement K-015, dan renewal tampak memakai 2
> (tiket R01) — cukup dikonfirmasi. Untuk dua jalur Quotation: minta pemilik Pega memastikan apakah keduanya
> diisi dari sumber yang sama; sementara itu pemuat sistem baru mengisi **kedua** jalur dari data yang
> tersedia dan tetap mencatat bila nilainya berbeda (tiket 10).

> ✅ **Jawaban work owner (1 Oktober 2026)** — `KEPUTUSAN-30-09-2026.md` butir 29
> **`StatusBusiness = 3` = Endorsement.** Pertanyaan dua jalur Quotation tetap terbuka.


## 6. Konfirmasi keputusan agent

Pilihan yang diambil agent selama pengerjaan, dicatat di `KEPUTUSAN-30-09-2026.md` bab "Keputusan agent —
menunggu konfirmasi". Mohon konfirmasi atau koreksi; sampai itu, kode berjalan menurut pilihan itu.

> 💡 **Rekomendasi agent (1 Oktober 2026)** — bukan jawaban; keputusannya tetap pada penerima.
>
> **Konfirmasi A1–A12**, dengan tiga tindak lanjut: **A4** — sekaligus amandemen teks K-018 di register
> supaya tidak lagi mengutip L1626 sebagai rumus; **A9** — tutup NB-02 tetapi buat satu tiket kecil "hitungan
> uang tanpa mata uang" yang terikat ke pemuat `repository`; **A10** — konfirmasi dengan Underwriting, karena
> memperluas K-014 (efek Reject juga pada Revise). Sisanya mengikuti perilaku sistem lama apa adanya.

> ✅ **Jawaban work owner (1 Oktober 2026)** — `KEPUTUSAN-30-09-2026.md` butir 34
> **A1–A12 dikonfirmasi**, dengan tindak lanjut A4 (amandemen K-018, sudah), A9 (NB-02 ditutup,
> tiket 17 dibuat), A10 (tetap dikonfirmasi dengan Underwriting).

## 7. Dua medan bocor yang tersisa — work owner (tinjauan manual)

**Yang kami minta:** keputusan apakah `Comment` dan `OperatorID` masuk daftar BUANG.

**Mengapa:** `[terverifikasi]` setelah 21 kunci, alat masih menemukan nilai yang sudah dibuang muncul lagi
di dua medan ini — `Comment` 5 kebocoran di 3 jalur (`$.Comment`, `$.OldData.Comment`, `$.OldData.OldData.Comment`) dan
`OperatorID` 6 kebocoran di 3 jalur (`$.QuotationData.OperatorID` dan dua salinannya di bawah `OldData`), tersebar di empat dari lima berkas
P-5. Alat menolak menulis keempat berkas itu. Daftar kandidat untuk tinjauan: 139 nama kunci unik (tanpa
nilai) — di antaranya beberapa yang namanya mengisyaratkan jabatan atau nama
(`OperationalDirector`, `OperationalDivHead`, `PresidentDirector`, `TechnicalDirector`, `Name`,
`GroupName`, `SobLeader0`) `[dugaan]`, belum bocor menurut alat tetapi layak ditinjau.

> 💡 **Rekomendasi agent (1 Oktober 2026)** — bukan jawaban; keputusannya tetap pada penerima.
>
> Masukkan `Comment` dan `OperatorID` — isinya teks bebas dan login operator. Tinjau tujuh nama kunci di
> atas sekaligus agar tidak perlu putaran lagi.

> ✅ **Jawaban work owner (1 Oktober 2026)** — `KEPUTUSAN-30-09-2026.md` butir 36
> **Setuju rekomendasi.** Kebocoran kini **nol** di kelima berkas. Ketujuh kunci ditinjau lewat profil
> bentuk tanpa nilai: keempat Director hanya `true`/`false`, `Name` kode mata uang, `GroupName` `-`,
> `SobLeader0` kode — tidak satu pun identitas.

## 8. Konfirmasi keputusan agent A13–A15

Dicatat di `KEPUTUSAN-30-09-2026.md` bab "Keputusan agent — menunggu konfirmasi (… penerapan butir 28–33)":
**A13** `= True` tetap ditolak; **A14** pembanding login operator ditolak dan nilainya tidak ditulis
(lima nilai login ternyata sudah tertulis di registry sejak NB-08 — kini nol); **A15** `pxUpdateOperator`
/ `pxUpdateOpName` tidak ditambah (0 kemunculan).

> 💡 **Rekomendasi agent (1 Oktober 2026)** — bukan jawaban; keputusannya tetap pada penerima.
>
> Konfirmasi ketiganya. Untuk A13, cukup lihat kondisi A `IsVisible` di UI Pega sekali bersamaan dengan
> pengecekan kode transisi 5 (`PERTANYAAN-AKSEPTASI.md` B1).

> ✅ **Jawaban work owner (1 Oktober 2026)** — `KEPUTUSAN-30-09-2026.md` butir 37
> **A13–A15 dikonfirmasi.**

## 9. Konfirmasi keputusan agent A16–A19 (bentuk fixture)

Saat menerapkan "kosongkan dulu" (butir 38): daftar dipisah menjadi identitas vs teks bebas (A16),
`ObjectName`/`Others` tidak dikosongkan (A17), bagian alamat `ASMCity`/`ASMDistrict`/`ASMRW` ikut
dikosongkan (A18), salinan angka di `Remarks` ikut hilang (A19). Rinciannya di register butir 38.

> 💡 **Rekomendasi agent (1 Oktober 2026)** — bukan jawaban; keputusannya tetap pada penerima.
>
> Konfirmasi keempatnya. Semuanya ke arah lebih aman, kecuali A17 — bila ragu, `ObjectName` dapat
> dipindah ke `DaftarKosongkan` tanpa mengubah angka apa pun.

> ✅ **Jawaban work owner (1 Oktober 2026)** — `KEPUTUSAN-30-09-2026.md` butir 49, atas butir 5 (dua jalur Quotation):
> **"Iya, ganti dari QuotationData, itu isinya sama aja."** `pyWorkPage.Quotation.*` diisi dari `QuotationData.*`.
