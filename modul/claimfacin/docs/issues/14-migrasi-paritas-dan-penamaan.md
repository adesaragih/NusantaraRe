# 14: Migrasi data lama, paritas yang ditegaskan, dan penamaan ulang

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR) · seluruh tiket 01–13
**Menutup:** AC 75 · 76 · 77 · 78 · 79 · 80 · 81 · 82 · 92 · 93 · 94 · 95 · 96 · 97 · 98 · 105 · 106 · 107 · 108 · 109 · 110 · 111 · 112 · 113 *(24 AC)* — US 27 · 36

## Hasil & nilai pengguna

Hari ini Data klaim lama **belum berpindah**, nama kolom yang menyesatkan **belum diluruskan**, dan perilaku yang sengaja **dipertahankan sama** belum punya uji yang membuktikannya.

Sesudah tiket ini, Data lama **berpindah dengan benar**, nama yang menyesatkan **diganti dengan aturan yang tertulis**, dan ⭐ **paritas yang sengaja dipertahankan punya uji** sehingga tidak berubah diam-diam.

## Area codebase

- Migrasi data klaim lama
- Penamaan ulang kolom dan aturannya
- Uji paritas

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Penamaan menyesatkan | bab aturan menamai ulang pada spec |
| Paritas | perilaku yang sengaja ditiru apa adanya |

⭐ Rincian medan dan asalnya ada di `STRUKTUR-TABEL-CLAIM-FACIN.md` **§2b**.

## ADR terkait

- **ADR-0007** — jejak audit atas migrasi

## Acceptance criteria

- [ ] **AC 75–78** — migrasi
- [ ] **AC 79–82** — nama kolom dan pembacaan
- [ ] **AC 92–98** — paritas yang ditegaskan
- [ ] **AC 105–113** — sisa yang ditegaskan
- [ ] ⭐ Tiap perilaku paritas punya **uji yang membuktikannya** — ⛔ supaya tidak 'terperbaiki' diam-diam oleh pengembang berikutnya

> ⛔ **RALAT 10-10-2026.** Butir lamanya dikutip utuh, tidak dihapus: *"**AC 92–98** — paritas yang ditegaskan"* dan
> *"**AC 105–113** — sisa yang ditegaskan"* → beberapa AC yang ditutup tiket ini **diralat di spec**:
>
> - **AC 93** — roster menurut jabatan → pemutus per tingkat lewat **workbasket** pola Komite Claim Prop (tahap 2,
>   `komiteclaimfacin`, OQ-CFI-04).
> - **AC 95** — **dibalik**: penomoran dihitung aplikasi (`inti/backend/penomor`, ADR-0043 meng-*supersede* ADR-0006).
> - **AC 112** — `STRUKTUR-TABEL-CLAIM-FACIN.md` (dengan lampiran pengikat) dan `RELASI-TABEL-CLAIM-FACIN.md` **sudah
>   ada**; butir 23 terjawab.
> - **AC 109** tetap: properti klasifikasi diisi — bertentangan dengan premis tiket 02 / 13 (lihat RALAT di sana).
>
> Paritas yang sengaja **diubah** dicatat sebagai `[penyimpangan sadar]` di `docs/PARITAS.md` §7–§8 (OQ-CFI-03); §7
> menyebut uji untuk tiap perbaikan.

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **29** | Cacah baris lama terdampak — `[data DBA]` | ⚠️ menahan **cacah**, tidak menahan cara migrasinya |
| **1** | Kolom tabel data kutipan — `[data DBA]` | ⚠️ menahan pemetaan kolom lama |

> ⛔ **RALAT 10-10-2026.** Baris lamanya dikutip utuh, tidak dihapus: *"| **29** | Cacah baris lama terdampak —
> `[data DBA]`"* dan *"| **1** | Kolom tabel data kutipan — `[data DBA]` | ⚠️ menahan pemetaan kolom lama"* →
>
> - ⛔ **Butir 29 TIDAK ADA di register mana pun** — register spec berakhir di butir **25**, register
>   `STRUKTUR-TABEL-CLAIM-FACIN.md` §6 di butir **13**. Isinya paling dekat dengan butir **25** spec (cacah klaim lama
>   yang kini menjadi galat, AC 114 — memblokir migrasi).
> - **Butir 1** (register spec): tidak menahan lagi — data polis dibaca dari `JSON_POLIS.DATA_JSON`, bukan dari
>   `T_QUOTATIONDATA` (yang ada di NB Fac In `183`, tanpa kolom `Business*`).

## Perintah verifikasi

1. Migrasikan satu klaim lama — ⭐ seluruh tingkatnya terbentuk benar.
2. Jalankan uji paritas — ⭐ perilaku yang sengaja sama **tetap sama**.
3. Cari nama kolom yang menyesatkan di skema baru — ⛔ **nihil**.
