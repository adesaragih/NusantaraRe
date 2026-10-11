# 07: Layar klaim — apa yang tampil dan apa yang dapat disunting

**Status:** ready-for-agent
**Blocked by:** 01 (registrasi) · 03 (estimasi) · 05 (penyesuaian)
**Menutup:** AC 99 · 100 · 101 · 102 *(4 AC)* — US 6 · 7

## Hasil & nilai pengguna

Hari ini Penilai **belum punya layar** untuk melihat klaim, objeknya, estimasinya, dan penyesuaiannya dalam satu tempat.

Sesudah tiket ini, Penilai melihat klaim **utuh dalam satu layar** — objek, item, estimasi, pembagian, penyesuaian — dan dapat menyunting yang memang boleh disunting.

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"Sesudah tiket ini, Penilai melihat klaim **utuh
> dalam satu layar** — objek, item, estimasi, pembagian, penyesuaian"* → **bertentangan dengan AC 100** (layar
> estimasi, penyesuaian, dan akseptasi **terpisah**; menggabungkannya menjadi satu layar **gagal**). Yang dibangun
> 10-10-2026 mengikuti section XML: **tiga layar, satu per assignment** — **Input Register** (`InputRegister`),
> **Input Estimasi** (`InputEstimasiAdmin`), **Choose Surveyor** (`ClaimSurvey` / `InputAdjustment`) — masing-masing
> **bertab** (mis. Register / Policy Detail & Claims History / Progress Claim; Estimation / View Registration; Claim
> Details / Policy Details and Claim History), dengan grid bertingkat objek → item → estimasi / adjustment
> (`PINDAI.md` §3, `PARITAS.md` §2–§4, `backend/models/layar_register.go` / `layar_estimasi.go` /
> `layar_adjustment.go`). Bagian lain tiket (penyuntingan ditolak di layar **dan** di layanan) tetap.

## Area codebase

- Antarmuka klaim
- Lapisan layanan: pembacaan klaim utuh

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Layar klaim | rule layar klaim fakultatif masuk |

⭐ Rincian medan dan asalnya ada di `STRUKTUR-TABEL-CLAIM-FACIN.md` **§2b**.

## ADR terkait

- **ADR-0012** — wewenang eksplisit, bukan efek samping layar

## Acceptance criteria

- [ ] **AC 99–102** — layar
- [ ] ⛔ Layar **tidak menjadi penjaga wewenang** — ia hanya menyembunyikan; penegakan di tiket 08

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **11** | ⚠️ **29 medan bergantung jenis objek** — menjadi kolom, atau dibaca dari polis | ⚠️ **menahan bentuk layar objek** |

> ⛔ **RALAT 10-10-2026.** Baris lamanya dikutip utuh, tidak dihapus: *"⚠️ **29 medan bergantung jenis objek** — menjadi
> kolom, atau dibaca dari polis"* → ⚠️ **ruang nomor**: butir **11** di sini = register **`STRUKTUR-TABEL-CLAIM-FACIN.md`
> §6**, bukan butir 11 spec (mesin tiket). Terjawab oleh yang dibangun: **dibaca dari polis** — `KUNCI_POLIS` menunjuk
> letak objek, halaman polis dibaca ulang dari `JSON_POLIS` setiap muat; grid objek per lini dipilih `GolonganLini`
> (`backend/models/lini.go`). Tidak menahan lagi.

## Perintah verifikasi

1. Buka satu klaim — ⭐ objek, item, estimasi, pembagian, dan penyesuaian tampil.
2. Coba sunting medan yang tidak boleh disunting — ⭐ layar menolaknya.
3. ⭐ Panggil lapisan layanan **langsung** untuk menyunting medan itu — ⛔ **juga ditolak** *(tiket 08)*.
