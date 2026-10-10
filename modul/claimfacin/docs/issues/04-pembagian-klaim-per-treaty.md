# 04: Pembagian klaim per treaty, dan aturan keseragaman berbagi

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR) · 03 (estimasi)
**Menutup:** AC 31 · 32 · 33 *(3 AC)* — US 8

## Hasil & nilai pengguna

Hari ini Pembagian klaim antar treaty **belum punya tabel** pada lini FAC, dan aturan keseragaman persentase berbagi **belum ada penegaknya**.

Sesudah tiket ini, Pembagian klaim tercatat per treaty per mata uang, dan ⭐ **persentase berbagi seragam untuk satu treaty di dalam satu item objek** — sehingga daftar ringkasnya dapat diturunkan tanpa tabel kedua.

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"sehingga daftar ringkasnya dapat diturunkan tanpa
> tabel kedua."* → tanpa tabel kedua tetap benar, tetapi daftar ringkas (`SpreadingList` polis item) **disimpan**, tidak
> diturunkan: baris `T_CLAIM_SPREADING` ber-`JENIS = 'POLIS'`, di samping baris `JENIS = 'KLAIM'` (`SpreadingClaim`) —
> migrasi `564`. Lihat RALAT 10-10-2026 di `STRUKTUR-TABEL-CLAIM-FACIN.md` §3.

## Area codebase

- Lapisan layanan klaim: pembagian per treaty
- ⭐ Penegak aturan keseragaman — ⛔ **di lapisan layanan**, bukan basis data

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Pembagian klaim | daftar pembagian di dalam item objek |
| Daftar ringkas per treaty | ⛔ **bukan tabel** — himpunan bagian, selisih **NIHIL** |

⭐ Rincian medan dan asalnya ada di `STRUKTUR-TABEL-CLAIM-FACIN.md` **§2b**.

## ADR terkait

- **ADR-0007** — jejak audit

## Acceptance criteria

- [ ] **AC 31 · 32 · 33** — pembagian klaim per treaty
- [ ] ⭐ **Semua baris ber-treaty sama di dalam satu item objek punya persentase berbagi yang sama**
- [ ] ⭐ Menyunting persentase **mengenai seluruh baris treaty itu sekaligus** — ⛔ bukan satu baris

> ⛔ **RALAT 10-10-2026 — `[penyimpangan sadar]` DIBUANG: tidak dibangun.** Butir lamanya dikutip utuh, tidak dihapus:
> *"⭐ Menyunting persentase **mengenai seluruh baris treaty itu sekaligus** — ⛔ bukan satu baris"* → janji ini
> **rancangan asisten** (`STRUKTUR-TABEL-CLAIM-FACIN.md` §4), **bukan perilaku XML** (audit B5). Di XML:
>
> - Medan **Share % spreading estimasi hanya-baca (RO=ALWAYS)** di `Section/Estimasi`, `EstimasiPA`,
>   `EstimasiMarine`, sehingga `CheckTotalSpreadingPct_Act` tak pernah terpicu (`PARITAS.md` §7 butir 10).
> - Grid spreading yang dapat disunting dan tombol **Save Spreading** (`LS42`) bergerbang `ClaimData.ExGratia = 1`,
>   yang **tidak pernah benar** — penulis satu-satunya `InsertObjects_dt` langkah 11 menulis `0` (OQ-CFI-17).
> - Tambah / hapus baris spreading estimasi bergerbang `… && 1=2` (`PARITAS.md` §3).
>
> Jadi tidak ada jalur menyunting persentase spreading estimasi, dan butir *"⭐ **Semua baris ber-treaty sama di dalam
> satu item objek punya persentase berbagi yang sama**"* tidak punya penegak yang perlu dibangun (butir 9 STRUKTUR
> gugur). Perintah verifikasi
> 1 dan 2 tidak berlaku; langkah 3 berubah menjadi: baris `JENIS = 'POLIS'` dan `'KLAIM'` tersimpan dan terbaca kembali.
> Yang dapat disunting per baris di XML hanya spreading **adjustment** Ex Gratia (tiket 05).

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **9** | ⚠️ Penegak aturan keseragaman belum ditetapkan — dijaga saat tulis, atau diperiksa berkala | ⚠️ **menahan bentuk penegakannya**, tidak menahan aturannya |

> ⛔ **RALAT 10-10-2026.** Baris lamanya dikutip utuh, tidak dihapus: *"⚠️ Penegak aturan keseragaman belum ditetapkan —
> dijaga saat tulis, atau diperiksa berkala"* → ⚠️ **ruang nomor**: butir **9** di sini = register
> **`STRUKTUR-TABEL-CLAIM-FACIN.md` §6**, bukan butir 9 spec (jenis berkas video, sudah ditutup). Butir ini **gugur**:
> aturannya dibuang sebagai rancangan di luar XML — lihat RALAT di atas.

## Perintah verifikasi

1. Sunting persentase berbagi satu baris — ⭐ **seluruh baris treaty itu ikut berubah**.
2. Coba simpan dua baris treaty sama dengan persentase berbeda — ⛔ **ditolak**.
3. Turunkan daftar ringkas per treaty — ⭐ **selisih NIHIL** terhadap pembagian penuh.
