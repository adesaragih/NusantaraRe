# 11: Penyerahan ke Komite dan penutupan klaim

**Status:** sebagian 07-10-2026 — penyerahan ke komite (opsi "b") dan penutupan klaim dibangun; gerbang lampiran selalu menolak (OQ-CP-12); tutup tanpa bayar (OQ-CP-06) — RALAT 07-10-2026 (semula `ready-for-agent`). **RALAT 10-10-2026:** tutup tanpa bayar dibangun (OQ-CP-06 selesai, `SendCloseClaimToKomite` → kasus komite TT 4; keputusan di Komite Claim Prop `KomitePost_Close`)
**AC yang dipegang (RALAT 07-10-2026, prompt §7 butir 4):** AC 130 (uang ke komite tetap desimal) dan AC 131 (tabel pemetaan `TSISpread` / `ClaimSpread` / `TotalAdjustment`) — uang penyerahan memakai `apd.Decimal`, nol float; tabel pemetaan ketiga parameter menunggu modul Komite Claim Prop.
**Blocked by:** 00 (PREFACTOR) · 08 (baris adjustment) · 10 (wewenang)
**Menutup:** AC 62 · 63 · 64 · 65 · 66 · 67 · 68 · 69 · 70 · **125** *(10 AC)* — US 36–44

## Hasil & nilai pengguna

Claim Admin menyerahkan baris adjustment ke Komite untuk diputuskan, dan penyerahan ditolak bila data
bank atau lampiran wajib belum lengkap — sehingga Komite tidak pernah memutus di atas data setengah
jadi. Klaim yang tidak dibayar tetap punya jejak keputusan lewat jalur penutupan tanpa pembayaran, dan
klaim tidak dapat ditutup selagi masih ada keputusan yang menggantung.

⚠️ **Batas konteks.** Tiket ini membangun **batas** menuju Komite Claim Prop, bukan isinya. Tangga
persetujuan, routing per tingkat, dan penomoran akseptasi milik konteks sebelah.

### ⚠️ Penyerahan MEMBEKUKAN baris penyesuaian induknya — **perilaku BARU**

`[keputusan work owner]` 2026-09-19 — **penyerahan ke komite melahirkan kasus komite, dan
keberadaan kasus itulah yang membekukan baris penyesuaian yang diserahkan.** Sejak saat itu baris
tersebut **tidak dapat diubah, selamanya**, dan **tidak dapat dihapus satu per satu**. **Penolakan
komite tidak mencairkannya**; perbaikan memakai **baris penyesuaian baru**.

⛔ **NOL PENANDA DIPASANG — tidak ada yang perlu ditulis.** Bekunya **diturunkan**, bukan disimpan:
sebuah baris beku **bila ada kasus komite yang menunjuknya**, berjalan maupun selesai. ⛔ **Jangan
membuat kolom penanda beku**; acuan tunggal nama kolom tetap `STRUKTUR-TABEL-CLAIM-PROP.md`.

⚠️ **Penegakannya bukan di sini.** Tiket ini hanya **melahirkan kasusnya**; yang menolak sunting
dan hapus satu-baris adalah **tiket 08**, dan yang menolak hapus klaim adalah **tiket 00**.

⚠️ `[terverifikasi]` **Pega tidak punya kunci ini** — nol pemeriksaan, nol pesan kesalahan, nol
penanganan untuk baris penyesuaian yang berubah atau hilang sesudah diserahkan. **Dibangun, bukan
dimigrasikan.** Lingkup: **Claim Prop** dan **Claim Non Prop** *(Non Prop menyusul)*; **Claim —
Life tidak termasuk**.

## Area codebase

Penyerahan kasus komite · penjaga kelengkapan · jalur penutupan tanpa pembayaran · penutupan klaim.

## Rule Pega sumber

| Rule | Step | Yang dibuktikannya |
| --- | --- | --- |
| `Activity/AddKomiteTreatyChild_ACT.xml` | — | membuat kasus anak komite; ⚠️ penunjuknya **indeks posisi**, muatan digemukkan 24 field sebagai snapshot |
| `Activity/SetKomiteTreaty_ACT.xml` | — | roster kasus komite |
| `Activity/AttachmentProtect_ACT.xml` | 5 | lampiran wajib bergantung jenis pembayaran dan batas nilai |
| `Activity/CloseClaimProp.xml` | 5 | ⚠️ memo penutupan ditulis ke slot yang **tidak pernah dibaca** (ditangani tiket 14) |
| `Activity/SaveOutstanding_Act.xml` | 39 | ⚠️ data dikirim ke sistem luar **tanpa menunggu komite** |
| `ConnectREST/SendAcceptationToKasir.xml` | — | ⭐ **BARU 2026-09-19 (ronde 4)** — pengiriman ke Kasir **berautentikasi**; bentuknya **profil autentikasi bernama**. AC 69 (*penutupan ditolak bila pengiriman Kasir belum berhasil*) karenanya bergantung pada kredensial yang **harus tersedia saat runtime** — rinciannya di **tiket 13** |
| `Activity/SetKomiteNo_Act.xml` | 1 | ⭐ **BARU 2026-09-19 (ronde 6)** — ⚠️ **nomor komite tidak disimpan, melainkan DIIRIS**: karakter ke-19 sampai ke-30 dari **kunci internal Pega** kasus anak **terakhir**. Dua kerapuhan sekaligus — bergantung **panjang/format kunci internal**, dan bergantung **posisi terakhir** dalam daftar. Menguatkan AC 63 (*rujukan memakai ID stabil, bukan indeks posisi*) |

⚠️ `[terverifikasi]` **Nol rule pembatal di 329 berkas.**

## ADR terkait

**ADR-0011** (unit keputusan = baris `AdjustmentList`) · **ADR-0014** (penegakan di lapisan layanan).

## Acceptance criteria

- [ ] `[keputusan work owner]` Penyerahan membuat rekam kasus komite yang membawa **penunjuk stabil** ke baris adjustment *(AC 62 spec)*
- [ ] ⚠️ Rujukan memakai **ID stabil**, bukan indeks posisi. **Alasan menyimpang:** di Pega penunjuknya indeks posisi, sehingga muatan harus digemukkan 24 field sebagai snapshot supaya rujukan tidak tergeser *(AC 63 spec)*
- [ ] `[terverifikasi]` Penyerahan **ditolak** bila data bank belum lengkap *(AC 64 spec)*
- [ ] `[terverifikasi]` Penyerahan **ditolak** bila lampiran wajib belum lengkap; daftar lampiran wajib **bergantung jenis pembayaran** *(AC 65 spec)*
- [ ] ⚠️ **Dua jalur ke Komite** — penyerahan adjustment dan penutupan tanpa pembayaran — diberi **saling-kunci**. **Alasan menyimpang:** di Pega keduanya menghasilkan kasus anak berkelas sama dan **keduanya dapat aktif pada klaim yang sama**, tanpa apa pun yang mencegahnya *(AC 66 spec)*
- [ ] ⚠️ Jalur penutupan tanpa pembayaran mengambil roster dari **data**. **Alasan menyimpang:** di Pega jalur itu memakai roster tertanam satu orang *(AC 67 spec)*
- [ ] `[terverifikasi]` Penutupan **ditolak** bila masih ada baris adjustment yang belum diputus *(AC 68 spec)*
- [ ] `[terverifikasi]` Penutupan **ditolak** bila pengiriman ke Kasir belum berhasil *(AC 69 spec)*
- [ ] ⚠️ **RISIKO DITERIMA SADAR — dipertahankan, bukan diperbaiki.** Klaim yang **ditolak** komite **tetap** memiliki data di sistem luar (Arasapas); pengiriman terjadi tanpa menunggu komite dan **nol rule pembatal** ada di 329 berkas. **Alasan menyimpang:** paritas atas keputusan `[keputusan work owner]`; rekonsiliasinya urusan manual di luar lingkup modul ini *(AC 70 spec)*
- [ ] ⚠️ Sesudah penyerahan tercatat, baris penyesuaian yang diserahkan **beku** — dan bekunya **diturunkan dari adanya kasus komite yang menunjuknya**, bukan dari kolom penanda. ⛔ Test yang mencari **kolom penanda beku gagal**; test yang menemukan baris terserahkan **tanpa kasus komite** juga **gagal**. **Alasan menyimpang:** di Pega baris boleh berubah atau hilang sesudah diserahkan, tanpa satu pun pemeriksaan *(AC 125 spec)*

## Perintah verifikasi

```
jalankan test "kirim ke komite tanpa data bank -> DITOLAK"
jalankan test "kirim ke komite tanpa lampiran wajib -> DITOLAK"
jalankan test "lampiran wajib berbeda menurut jenis pembayaran"
jalankan test "jalur adjustment aktif -> jalur tutup-tanpa-bayar DITOLAK (saling-kunci)"
jalankan test "tutup klaim dengan baris adjustment belum diputus -> DITOLAK"
jalankan test "tutup klaim sebelum kirim Kasir berhasil -> DITOLAK"
jalankan test "rujukan baris adjustment memakai ID stabil, tidak tergeser saat baris lain dihapus"
jalankan test "klaim ditolak komite -> data di sistem luar TETAP ADA (risiko diterima sadar)"
```
