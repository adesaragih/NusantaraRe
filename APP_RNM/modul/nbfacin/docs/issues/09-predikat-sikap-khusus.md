# 09: Empat predikat yang sikapnya sudah diputuskan, dan tidak boleh disamakan

**What to build:** Empat predikat yang masing-masing sudah punya keputusan tersendiri, dijalankan
persis sebagaimana diputuskan — bukan diseragamkan menjadi "implementasikan saja kondisinya".

**`IsPKSASM` → gagal keras.** Kondisinya terbaca, tetapi **labelnya belum ter-resolve**, tag kondisi
pertamanya masih berisi placeholder, dan ia bertanda sementara — di ketiga folder korpus. Ia gagal
keras bukan karena kondisinya hilang, melainkan karena **belum terbukti sebagai yang dieksekusi**.

⛔ Daftar nilai sah untuk properti yang diujinya sudah diketahui dari basis data. **Itu tidak
mencabut gagal-kerasnya.** Mengetahui nilai apa yang sah dan membuktikan kondisi itu dieksekusi
adalah dua pertanyaan berbeda; hanya yang kedua yang jadi penyebab, dan yang kedua belum terjawab.

**`IsOfferFacIn` → ekspresi tersimpan yang berlaku.** Rule ini memuat **dua kondisi berbeda**: teks
tampilannya menyebut daftar kode bisnis, ekspresi tersimpannya menguji penanda fakultatif. Yang
berlaku adalah **ekspresi tersimpan**. Teks tampilan dicatat sebagai kandidat perbaikan, **tidak**
diimplementasikan — dan menjadi tersangka pertama bila paralel run memperlihatkan selisih pada
gerbang masuk siklus.

**`IsSpreadingDepan` → implementasikan, catat asal salinannya.** Rule ini hilang dari folder NB tetapi
terbaca penuh di folder siklus lain. Ia **bukan** rule yang kondisinya tidak diketahui, jadi tidak
perlu gagal keras. ⚠️ Salinan yang terbaca berversi lebih lama, dan 95 rule di korpus terbukti berbeda
versi antar folder — asal salinan **wajib dicatat** di komentar.

**`IsFacout` → diport apa adanya, nama dipertahankan.** Fiturnya usang secara bisnis, tetapi kodenya
masih aktif dan masih dirujuk sepuluh berkas. ⚠️ Namanya menyiratkan *fac out*; **isinya menguji
hasil keputusan Banding**. Nama dipertahankan demi ketertelusuran ke rule asal — **tetapi namanya
bukan dokumentasi artinya**, dan komentar wajib menyebutkan itu.

**Blocked by:** 08

**Status:** ready-for-human — diport 01-10-2026, menunggu tinjauan work owner

- [ ] `IsPKSASM` **gagal keras** bila dievaluasi; pesannya menyebut alasannya (kondisi belum terbukti dieksekusi), bukan "nilai tidak dikenal"
- [ ] `IsOfferFacIn` memakai ekspresi tersimpan; teks tampilan dicatat sebagai kandidat perbaikan di komentar, tidak dieksekusi
- [ ] `IsSpreadingDepan` diimplementasikan sesuai kondisi terbacanya, dengan komentar menyebut **folder asal salinan dan versinya**
- [ ] `IsFacout` diport apa adanya dengan nama aslinya; komentar menyatakan isinya menguji **Banding**, bukan fac out
- [ ] Keempatnya punya kasus uji sendiri — termasuk uji bahwa `IsPKSASM` **memang** gagal keras

## Comments

### 2026-10-01 — implementasi (agent)

Kode: `rules/bangkit/main.go` (`pinjaman`, `catatanSikap`), medan `predikat.catatan`, `rules/sikap_test.go`.

| Predikat | Keadaan |
| --- | --- |
| `IsPKSASM` | ✅ panic dengan alasan "kondisi belum terbukti dieksekusi" (sejak NB-08) |
| `IsOfferFacIn` | ✅ ekspresi tersimpan `.Quotation.BusinessFac = "F"`; teks tampilan dicatat (K-002); uji membuktikan kondisi tampilan TIDAK berlaku |
| `IsSpreadingDepan` | ✅ dipinjam dari `Endorsment Fac In\When` (K-003) lewat daftar pinjaman eksplisit; asal menyebut folder dan `pyRuleSetVersion 01-01-52` |
| `IsFacout` | ✅ nama dipertahankan, catatan "menguji ProposalAcceptStatus = 4 (Banding), BUKAN fac out" (K-019) |

⚠️ Spec menyebut salinan `IsSpreadingDepan` "berversi lebih lama". `[terverifikasi]` `pyRuleSetVersion`-nya
01-01-52, **sama** dengan rule NB seperti `IsPA` — klaim itu tidak didukung tag versi ini. Registry kini
210 predikat (209 NB + 1 pinjaman). Pembangkit memakai `-akar` (akar korpus); env penjaga kini
`KORPUS_RNM_AKAR`. Uji mutasi: 2/2 (hanya tertangkap penjaga registry yang butuh korpus).

### 2026-10-01 — tindak lanjut review (agent)

Uji `IsPKSASM` kini memeriksa alasan lengkapnya, *"kondisi belum terbukti dieksekusi"*, bukan hanya
"sikap panic".
