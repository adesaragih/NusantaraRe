# 11: Efek keluar — kasir, surel, dan konversi

**Status:** ready-for-agent
**Blocked by:** 06 (uang) · 09 (dokumen)
**Menutup:** AC 57 · 58 · 59 · 60 · 61 · 62 · 63 · 64 *(8 AC)* — US 30 · 31 · 33

## Hasil & nilai pengguna

Hari ini Instruksi pembayaran, surel pemberitahuan, dan konversi klaim **belum terkirim ke mana pun**. ⚠️ Dan di sistem lama, **penjaga ganda-bayar bersandar pada gerbang yang belum dapat dipastikan bacanya**.

Sesudah tiket ini, Instruksi pembayaran terkirim ke kasir **tepat satu kali**, surel terkirim **hanya di lingkungan produksi**, dan ⭐ **pemberitahuan galat terkirim HANYA ketika pengiriman gagal**.

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"Sesudah tiket ini, Instruksi pembayaran terkirim
> ke kasir **tepat satu kali**"* → yang dibangun **tidak menjanjikan "tepat satu kali"** (lihat juga RALAT spec AC 60):
>
> - Efek keluar — konversi (`KonversiKlaimNonLife`), kasir (`SendAcceptationToKasir`), DLA (`HitDLAClaimFacin`), surel —
>   **hanya diantre bila `IS_PEGA_PROD`**, ke outbox `T_LOG_SERVICE_RNM` di transaksi aksi yang sama
>   (`backend/services/efek.go`); di luar produksi nol baris.
> - Modul ini **belum punya pekerja outbox** (`backend/modul.go` `JalankanPekerja` → `inti.TanpaPekerja()`).
>   Pelaksananya `PelaksanaClaimFacIn` **berhenti terang**: `ErrPengirimStubNonProduksi` di luar produksi;
>   `ErrKasirBelumDisetujui` / `ErrArasapasBelumDisetujui` di produksi sampai panggilan nyata disetujui manusia.
> - **Tanpa dedupe**: klik Acceptation ulang selagi kiriman pertama masih antre dapat mengantre muatan kedua
>   (**OQ-CFI-26**). `DIRECTTOKASIR_LOG` ditulis bersama tanggapan kasir, jadi belum terjadi.
> - Surel hanya di produksi — **sesuai** yang dibangun.

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

> ⛔ **RALAT 10-10-2026.** Butir lamanya dikutip utuh, tidak dihapus: *"⭐ Pembayaran terkirim **tepat satu kali**;
> panggilan kedua **ditolak**"* · *"⭐ Penanda terkirim tersimpan **dalam transaksi yang sama** dengan pengirimannya"* ·
> *"⭐ Pemberitahuan galat terkirim **hanya pada kegagalan**"* →
>
> - Butir pertama: **tidak dibangun** dalam bentuk itu — sekali **per klik** lewat outbox (`PARITAS.md` §8 butir 6),
>   tanpa dedupe (OQ-CFI-26).
> - Butir kedua: **tidak dibangun** — tidak ada pengiriman nyata (pelaksana berhenti terang); penanda kasir
>   (`DIRECTTOKASIR_LOG`) ditulis bersama tanggapan kasir (OQ-CFI-26).
> - Butir ketiga: **`[penyimpangan sadar]` DIBUANG — tidak dibangun.** Menurut tiket ini sendiri, gerbang XML-nya
>   berbendera mati sehingga pemberitahuan terkirim **setiap kali**; "hanya pada kegagalan" mengubah perilaku XML.
>   Bentuk XML-nya **perlu dicek terhadap XML** sebelum ada yang dibangun.
>
> Perintah verifikasi 2–4 belum punya sasaran di tahap 1; langkah 5 (non-produksi: nol surel) berlaku.

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **26** | ⚠️ Tanda sama dengan **tunggal** pada penjaga ganda-bayar — pembandingan atau penugasan | ⛔ **tidak menahan** — penjaga baru dibuat eksplisit; jawabannya untuk **memeriksa data lama** |
| **27** | Apakah pemberitahuan galat lama benar-benar sampai ke seseorang | tidak menahan |

> ⛔ **RALAT 10-10-2026.** Baris lamanya dikutip utuh, tidak dihapus: *"| **26** | ⚠️ Tanda sama dengan **tunggal** pada
> penjaga ganda-bayar"* dan *"| **27** | Apakah pemberitahuan galat lama benar-benar sampai ke seseorang"* → ⛔ **butir 26
> dan 27 TIDAK ADA di register mana pun**: register spec berakhir di butir **25**, register
> `STRUKTUR-TABEL-CLAIM-FACIN.md` §6 di butir **13**. Keduanya pertanyaan tanpa rumah; pertanyaan terbuka modul ini kini
> dicatat di `docs/OQ.md` — yang terdekat dengan butir 26 adalah **OQ-CFI-26** (kasir tanpa dedupe); butir 27 tidak
> punya padanan OQ.

## Perintah verifikasi

1. Kirim satu instruksi pembayaran — ⭐ terkirim, penanda terisi.
2. Kirim ulang untuk akseptasi yang sama — ⛔ **ditolak**.
3. Paksa kegagalan pengiriman — ⭐ pemberitahuan galat **terkirim**.
4. Kirim yang berhasil — ⛔ pemberitahuan galat **TIDAK terkirim**.
5. Jalankan di lingkungan bukan produksi — ⛔ surel **tidak terkirim**.
