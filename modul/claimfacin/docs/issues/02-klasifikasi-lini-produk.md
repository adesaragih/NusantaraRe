# 02: Klasifikasi lini produk — 13 penggolong hidup, 36 tidak dialihkan

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR) · 01 (registrasi)
**Menutup:** AC 13 · 14 · 15 · 16 · 17 · 18 *(6 AC)* — US 34–36

## Hasil & nilai pengguna

Hari ini Penggolongan lini usaha **membaca medan data kutipan yang tidak pernah disalin** ke objek kerja. ⛔ `[terverifikasi]` Akibatnya cabang **MBU** dan **Travel** tidak pernah terbit — diam-diam, tanpa galat.

Sesudah tiket ini, Setiap lini usaha **tergolong benar**, termasuk MBU dan Travel yang dulu tak pernah terbit. ⭐ Penggolongan membaca **data kutipan yang lengkap**.

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"Hari ini Penggolongan lini usaha **membaca
> medan data kutipan yang tidak pernah disalin** ke objek kerja. ⛔ `[terverifikasi]` Akibatnya cabang **MBU** dan
> **Travel** tidak pernah terbit — diam-diam, tanpa galat."* → tiket ini **bertentangan dengan spec-nya sendiri**:
>
> - **AC 109** (butir 22 ditutup): properti yang diuji rule klasifikasi **benar-benar diisi** — `pyWorkPage.Quotation`
>   adalah salinan halaman utuh `OfferFacIn.QuotationData` yang diisi modul penawaran. Premis "tidak pernah disalin"
>   adalah premis yang sudah diralat di spec Bab 4.
> - **AC 18**: cacat pasangan lompatan `GetAllData_Act` **ditiru apa adanya**, bukan "ditutup"; ia hanya menyentuh
>   tampilan satu pop-up.
> - Angka **"49 rule penggolong · 13 hidup · 36 tidak dipakai"** tidak ada di spec. Korpus memuat **60** rule `When`
>   (`PINDAI.md` §1), dan yang dibangun 10-10-2026 = **satu fungsi Go per rule `When`** yang dipakai section /
>   activity, bernama sama, di **satu tempat** `backend/models/lini.go` (`IsPA_PNC` → `IsPAPNC`; `GolonganLini`
>   memilih varian grid); `IsPEGAPROD` menjadi flag lingkungan. Halaman yang diuji dibaca dari `JSON_POLIS.DATA_JSON`.
>
> Yang mengikat adalah XML + AC 13 / 16 / 17 / 18 / 109, bukan janji "MBU dan Travel kini terbit".

## Area codebase

- Lapisan layanan: penggolongan lini usaha
- Penyalinan data kutipan dari kasus induk

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Penggolong lini usaha | 49 rule penggolong; ⭐ **13 dipakai hidup**, ⛔ **36 tidak dipakai sama sekali** |
| Sumber data kutipan | salinan data kutipan ke objek kerja |

⭐ Rincian medan dan asalnya ada di `STRUKTUR-TABEL-CLAIM-FACIN.md` **§2b**.

## ADR terkait

- **ADR-0003** — dua kaki dasar klasifikasi

## Acceptance criteria

- [ ] **AC 13–18** — klasifikasi lini produk
- [ ] ⭐ **13 penggolong** dialihkan; ⛔ **36 tidak dibangun**
- [ ] ⭐ Cabang **MBU** dan **Travel** **terbit** bila datanya memenuhi — ⚠️ inilah cacat lama yang ditutup

> ⛔ **RALAT 10-10-2026.** Butir lamanya dikutip utuh, tidak dihapus: *"⭐ **13 penggolong** dialihkan; ⛔ **36 tidak
> dibangun**"* dan *"⚠️ inilah cacat lama yang ditutup"* → **60 `When` di XML, satu fungsi per `When` di
> `backend/models/lini.go`**; "cacat lama" tidak ada menurut AC 109 — lihat RALAT di atas. Di perintah verifikasi,
> bagian *"di sistem lama tidak"* pada langkah 2 dan seluruh langkah 3 (36 penggolong) tidak berlaku.

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **6** | Isi rule peran | ⛔ tidak menyentuh tiket ini |

## Perintah verifikasi

1. Jalankan satu klaim tiap lini usaha yang hidup — ⭐ ketiga belasnya tergolong benar.
2. Jalankan klaim lini **MBU** dan **Travel** — ⭐ cabangnya **terbit**; ⛔ di sistem lama tidak.
3. Cari pemanggilan salah satu dari 36 penggolong yang tidak dialihkan — ⛔ **nihil**.
