# 06: Efek akhir — terbit sekali, oleh jenjang terakhir

**Status:** ready-for-agent
**Blocked by:** 05
**Menutup:** AC 43 · 44 · 45 · 46 · 47 · 48 · 49 · 50 · 51 *(9 AC)* — US 28 · 29 · 32 · 33

## Hasil & nilai pengguna

Hari ini Ringkasan akseptasi, nomor akseptasi, dokumen, dan surel **belum terbit sama sekali** — dan ⚠️ di sistem lama, langkah yang menyusunnya terpisah dari langkah yang menyimpannya, sehingga pernah dikira **penerimaan tidak tersimpan**.

Sesudah tiket ini, ⭐ **Efek akhir terbit HANYA ketika jenjang terakhir menyetujui**, dan **tepat satu kali**: ringkasan akseptasi disusun lalu **disimpan**, nomor terbit, dokumen tercetak, surel terkirim.

## Perilaku Pega yang ditiru

| Yang dibaca | Rule |
| --- | --- |
| Penyusun ringkasan | 51 langkah; ⛔ langkah penyimpan **di dalamnya** ber-remark |
| ⭐ Penyimpan sesungguhnya | ⭐ langkah penyimpan **diangkat ke rule pemanggil**, bergerbang *keputusan terima* **dan** *jenjang terakhir* |
| ⚠️ Jalur kegagalan | ⛔ **NOL jalur kegagalan pada 25 langkah penyimpanan** di sistem lama |

> ⛔ **RALAT 10-10-2026.** Kalimat dan baris lamanya dikutip utuh, tidak dihapus: *"Sesudah tiket ini, ⭐ **Efek akhir
> terbit HANYA ketika jenjang terakhir menyetujui**, dan **tepat satu kali**: ringkasan akseptasi disusun lalu
> **disimpan**, nomor terbit, dokumen tercetak, surel terkirim."* dan *"⚠️ Jalur kegagalan | ⛔ **NOL jalur kegagalan
> pada 25 langkah penyimpanan** di sistem lama"* →
>
> - **Satu transaksi per Submit** (prompt §5 #6, perbaikan atas Obj-Save klaim + komite per iterasi
>   `KomitePost_Adjustment` S7.2.1.10–14): keputusan tingkat, tulis balik klaim lewat kontrak, nomor
>   (`inti/backend/penomor`), `OS_AKSEPTASI_KLAIM`, `JSON_KLAIM`, `HISTORYAKSEPTASIPEGA`, `SUBPROGRESSCLAIM.POSITION2`,
>   dan kronologi di-commit bersama. Gagal di tengah = seluruh Submit batal; itulah jalur kegagalannya.
> - **Efek luar lewat outbox, hanya produksi** (`IS_PEGA_PROD`): konversi, kasir, surel; `MUATAN` hanya pengenal +
>   angka. PDF diterbitkan sesudah commit (isi berkasnya OQ, stream tidak diekspor).
> - **`OS_AKSEPTASI_KLAIM` sekali per KMT** (prompt §5 #5): S8 `SaveOSClaim_SQL` berada di luar perulangan adjustment
>   dan memakai `InputData` adjustment terakhir. Dengan satu adjustment per KMT (OQ-CFI-28) hasilnya sama; baris ditulis
>   untuk adjustment KMT itu. STS **4** bila PaymentType 1, selain itu **1**; `CARI17` 8 bila objek retro, selain itu 7.
> - **Penolakan TT2 tanpa OS**: S8 hanya berjalan bila `AcceptStatus == "1"` dan `KomiteCount == KomiteLoop`. Penolakan
>   juga tanpa nomor, tanpa `JSON_KLAIM` (S16), tanpa konversi (S17).
> - **Surel tidak hanya di tingkat akhir.** `SendEmailKlaim_KMT` (S7.2.1.15, `IsPEGAPROD`) berjalan di **setiap**
>   keputusan: ke tingkat berikut (S12) atau ke pembuat dengan "(Approval)" (S14) / "(Reject)" (S15).
>   `HISTORYAKSEPTASIPEGA` (S20–S21, Workbasket "KLAIM") dan kronologi (S3–S5) juga ditulis setiap keputusan.

⭐ Sumber: `komite-claim-facin\spec.md` · `claim-facin\STRUKTUR-TABEL-CLAIM-FACIN.md` §5.

## Keputusan work owner yang mengikat

- **K1** — ⭐ Efek akhir menempel pada **jenjang terakhir menyetujui**

## Yang harus diuji

- [ ] ⭐ Efek akhir terbit **hanya** ketika **jenjang terakhir menyetujui**
- [ ] ⭐ Ringkasan akseptasi **disusun lalu disimpan** — penyimpanan terjadi **sekali**
- [ ] ⛔ **RALAT yang wajib dibaca:** vonis lama *“penolakan tersimpan, penerimaan tidak”* **BATAL** — ⭐ keduanya tersimpan; yang berbeda hanya **di rule mana** langkah penyimpannya duduk
- [ ] Urutannya: ringkasan → nomor → dokumen → surel → kasir
- [ ] Surel terkirim **hanya di lingkungan produksi**
- [ ] ⚠️ Kegagalan salah satu efek **tidak boleh** membatalkan efek yang sudah terbit tanpa jejak — ⭐ sistem baru **wajib** punya jalur kegagalan

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **1** | Kolom tabel data kutipan — `[data DBA]` | tidak menahan |

## Seam & verifikasi

**Seam:** lapisan layanan komite — penyelesaian kasus.
1. Jalankan tangga tiga jenjang sampai selesai ⇒ ⭐ efek akhir terbit **tepat satu kali**.
2. Periksa jenjang pertama dan kedua ⇒ ⛔ **nol efek akhir** di keduanya.
3. Paksa gagal di tengah rangkaian efek ⇒ ⭐ kegagalannya **tercatat**, tidak hilang diam-diam.
4. Jalankan di lingkungan bukan produksi ⇒ ⛔ surel **tidak terkirim**.

> ⛔ **RALAT 10-10-2026.** Butir dan langkah lamanya dikutip utuh, tidak dihapus: *"Urutannya: ringkasan → nomor →
> dokumen → surel → kasir"*, *"2. Periksa jenjang pertama dan kedua ⇒ ⛔ **nol efek akhir** di keduanya."* dan *"3. Paksa
> gagal di tengah rangkaian efek ⇒ ⭐ kegagalannya **tercatat**, tidak hilang diam-diam."* →
>
> - **Urutan XML** `KomitePost_Adjustment`: nomor (S7.2.1.2) → tulis adjustment (S7.2.1.5) → surel (S7.2.1.15) →
>   `SaveAcceptation_KMT` (S7.2.1.16) → `SaveAccept_ACT` (S7.2.1.17) → OS (S8) → PDF (S9) → `IsPrintAccept` (S12) →
>   `JSON_KLAIM` (S16) → konversi (S17) → log "AKSEPATSI" (S19) → `HISTORYAKSEPTASIPEGA` (S20–S21) →
>   `SUBPROGRESSCLAIM` (S23) → `KomiteCount + 1` (S24) → kasir (S25). Nomor terbit **sebelum** ringkasan disimpan.
> - Langkah 2: jenjang pertama dan kedua tetap menulis surel, `HISTORYAKSEPTASIPEGA`, dan kronologi; yang nol hanya efek
>   akseptasi (nomor, OS, `JSON_KLAIM`, konversi, kasir, PDF).
> - Langkah 3: gagal di dalam transaksi ⇒ nol tulisan tersisa; gagal efek outbox ⇒ antre ulang. Tambah uji: tolak TT2
>   di tingkat akhir ⇒ nol baris `OS_AKSEPTASI_KLAIM`.
