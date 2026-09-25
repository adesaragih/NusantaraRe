# 15: Tiga prefix klaim dan varian Syariah

**Status:** ready-for-agent
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

## Perintah verifikasi

```
jalankan test "klaim tiap prefix -> terdaftar, ternomori, dan tersimpan"
jalankan test "klaim varian Syariah -> perlakuan identik dengan induknya"
jalankan test "tiap jenis -> masuk ke jalur simpan proyeksi OS yang benar"
jalankan test "regresi lintas prefix: registrasi, estimasi, adjustment, komite, dokumen, efek keluar"
```
