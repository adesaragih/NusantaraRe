# 11: Efek keluar — kasir, surel, dan konversi

**Status:** ready-for-agent
**Blocked by:** 06 (uang) · 09 (dokumen)
**Menutup:** AC 57 · 58 · 59 · 60 · 61 · 62 · 63 · 64 *(8 AC)* — US 30 · 31 · 33

## Hasil & nilai pengguna

Hari ini Instruksi pembayaran, surel pemberitahuan, dan konversi klaim **belum terkirim ke mana pun**. ⚠️ Dan di sistem lama, **penjaga ganda-bayar bersandar pada gerbang yang belum dapat dipastikan bacanya**.

Sesudah tiket ini, Instruksi pembayaran terkirim ke kasir **tepat satu kali**, surel terkirim **hanya di lingkungan produksi**, dan ⭐ **pemberitahuan galat terkirim HANYA ketika pengiriman gagal**.

## Area codebase

- Lapisan layanan: efek keluar
- ⭐ Penjaga ganda-bayar — **eksplisit**, dalam satu transaksi

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Pengiriman ke kasir | rule panggil layanan kasir; dua panggilan luar, keduanya bergerbang lingkungan produksi |
| Penjaga ganda-bayar | ⚠️ gerbang lama memakai tanda sama dengan **tunggal** — ⛔ bacanya belum pasti |
| Pemberitahuan galat | ⚠️ gerbang lama **bendera mati** ⇒ terkirim **setiap kali**; medannya salah eja |

⭐ Rincian medan dan asalnya ada di `STRUKTUR-TABEL-CLAIM-FACIN.md` **§2b**.

## ADR terkait

- **ADR-0007** — jejak audit atas efek keluar

## Acceptance criteria

- [ ] **AC 57–64** — efek keluar
- [ ] ⭐ Pembayaran terkirim **tepat satu kali**; panggilan kedua **ditolak**
- [ ] ⭐ Penanda terkirim tersimpan **dalam transaksi yang sama** dengan pengirimannya
- [ ] ⭐ Pemberitahuan galat terkirim **hanya pada kegagalan**
- [ ] ⭐ Surel terkirim **hanya di lingkungan produksi**

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **26** | ⚠️ Tanda sama dengan **tunggal** pada penjaga ganda-bayar — pembandingan atau penugasan | ⛔ **tidak menahan** — penjaga baru dibuat eksplisit; jawabannya untuk **memeriksa data lama** |
| **27** | Apakah pemberitahuan galat lama benar-benar sampai ke seseorang | tidak menahan |

## Perintah verifikasi

1. Kirim satu instruksi pembayaran — ⭐ terkirim, penanda terisi.
2. Kirim ulang untuk akseptasi yang sama — ⛔ **ditolak**.
3. Paksa kegagalan pengiriman — ⭐ pemberitahuan galat **terkirim**.
4. Kirim yang berhasil — ⛔ pemberitahuan galat **TIDAK terkirim**.
5. Jalankan di lingkungan bukan produksi — ⛔ surel **tidak terkirim**.
