# 04: Pembagian klaim per treaty, dan aturan keseragaman berbagi

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR) · 03 (estimasi)
**Menutup:** AC 31 · 32 · 33 *(3 AC)* — US 8

## Hasil & nilai pengguna

Hari ini Pembagian klaim antar treaty **belum punya tabel** pada lini FAC, dan aturan keseragaman persentase berbagi **belum ada penegaknya**.

Sesudah tiket ini, Pembagian klaim tercatat per treaty per mata uang, dan ⭐ **persentase berbagi seragam untuk satu treaty di dalam satu item objek** — sehingga daftar ringkasnya dapat diturunkan tanpa tabel kedua.

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

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **9** | ⚠️ Penegak aturan keseragaman belum ditetapkan — dijaga saat tulis, atau diperiksa berkala | ⚠️ **menahan bentuk penegakannya**, tidak menahan aturannya |

## Perintah verifikasi

1. Sunting persentase berbagi satu baris — ⭐ **seluruh baris treaty itu ikut berubah**.
2. Coba simpan dua baris treaty sama dengan persentase berbeda — ⛔ **ditolak**.
3. Turunkan daftar ringkas per treaty — ⭐ **selisih NIHIL** terhadap pembagian penuh.
