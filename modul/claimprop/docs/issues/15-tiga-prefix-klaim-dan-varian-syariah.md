# 15: Tiga prefix klaim dan varian Syariah

**Status:** dibangun 07-10-2026 — `CLMP-` saja (RALAT prefix di bawah); varian Syariah OQ-CP-14 — RALAT 07-10-2026 (semula `ready-for-agent`)
**Blocked by:** 00 (PREFACTOR) · 13 (efek keluar)
**Menutup:** AC 75 · 76 *(2 AC)* — US 66

## Hasil & nilai pengguna

Modul tetap melayani ketiga lini yang selama ini ditanganinya, berikut varian Syariah-nya — tidak ada
lini yang tertinggal saat pindah ke Go. Klaim Syariah diperlakukan identik dengan induknya, sehingga
tidak ada aturan perhitungan kembar yang harus dirawat dua kali.

⚠️ Tiket ini adalah **tiket paritas lintas-lini**: ia menyapu seluruh jalur yang sudah dibangun tiket
sebelumnya dan memastikan ketiganya berlaku untuk semua prefix — bukan hanya untuk prefix yang
kebetulan dipakai saat pengembangan.

## Area codebase

Penentuan prefix klaim · tiga jalur simpan proyeksi OS akseptasi · penomoran per prefix.

## Rule Pega sumber

| Rule | Step | Yang dibuktikannya |
| --- | --- | --- |
| `Activity/CekPremiLunas_Act.xml` | 2 · 3 | cabang per prefix — ⚠️ hanya dua dari tiga prefix punya cabang (lihat tiket 02) |
| `RDBList/SaveOSClaim_SQL.xml` | — | jalur simpan proyeksi OS akseptasi — klaim |
| `RDBList/SaveOSKlaimTreaty_SQL.xml` | — | jalur simpan proyeksi — treaty |
| `RDBList/SaveDataToOsAkseptasiNP.xml` | — | jalur simpan proyeksi — treaty non proporsional |
| `RDBList/GenerateNoCLMTreatyIn.xml` | — | penomoran menerima kode dan tipe sebagai parameter |

## ADR terkait

**ADR-0006** (penomoran lewat stored procedure) · **ADR-0015** (efek keluar wajib berhasil).

## Acceptance criteria

- [ ] `[terverifikasi]` Modul melayani **tiga prefix klaim**, dan varian Syariah diperlakukan **identik** dengan induknya *(AC 75 spec)*
- [ ] `[terverifikasi]` **Tiga jalur simpan proyeksi OS akseptasi terpisah** dipertahankan sesuai jenis *(AC 76 spec)*

> ⚠️ **RALAT 07-10-2026** — kalimat lama: *"Modul melayani tiga prefix klaim, dan varian Syariah diperlakukan identik dengan induknya"* Diputuskan dari XML: `When/IsCLMP` menguji `pyWorkIDPrefix = "CLMP-"` dengan properti kelas
> `ASM-FW-GCNMFW-Work-ClaimTreaty` (kelas Claim Prop); `IsCLM` = kelas PNC (Fac), `IsCLMNP` = kelas ClaimTreatyNonProp.
> Ketiganya hanya dipakai activity BERSAMA kelas Data-Adjustment (`HitServiceToKasir_Act`, `SendEmailKlaim`) yang
> bercabang untuk tiga lini. Modul ini melayani **`CLMP-` saja** (STRUKTUR T3); cabang `CLM-` / `CLMNP-` tidak
> dimigrasikan (Out of Scope "Claim Non Prop dan Claim Fac In"). Varian ber-`S` (`CLMPS-`) hanya muncul di
> `SendEmailKlaimRejectClose` (kelas `ASM-FW-GCNMFW-Work`) — OQ-CP-14.

> ⚠️ **RALAT 07-10-2026** — kalimat lama: *"Tiga jalur simpan proyeksi OS akseptasi terpisah dipertahankan sesuai jenis"* Jalur `KLAIMTRT` **tidak ada**: procedure `PEGA_JSON_OS_AKSEP_KLAIMTRT` tidak ada di DEV dan
> jalurnya dibuang tiket 00 (keputusan work owner 18-09-2026); `..._KLAIMTNP` milik Claim Non Prop. Claim Prop
> punya **satu** jalur simpan OS: baris `OS_AKSEPTASI_KLAIM` (isi `PEGA_JSON_OS_AKSEP_KLAIM` ditulis ulang sebagai
> SQL langsung, `repository.SisipOS`).

## Perintah verifikasi

```
jalankan test "klaim tiap prefix -> terdaftar, ternomori, dan tersimpan"
jalankan test "klaim varian Syariah -> perlakuan identik dengan induknya"
jalankan test "tiap jenis -> masuk ke jalur simpan proyeksi OS yang benar"
jalankan test "regresi lintas prefix: registrasi, estimasi, adjustment, komite, dokumen, efek keluar"
```
