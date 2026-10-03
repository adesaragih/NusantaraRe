# 09: Jalur kasir — penjaga ganda-bayar eksplisit

**Status:** ready-for-agent
**Blocked by:** 06 · `claim-facin\issues\11` *(efek keluar klaim)*
**Menutup:** AC 61 · 62 · 63 · 64 · 65 · 66 · 67 · 68 · 69 · 70 *(10 AC)* — US 30 · 31

## Hasil & nilai pengguna

Hari ini Instruksi pembayaran **belum terkirim ke kasir**, dan ⚠️ di sistem lama **penjaga ganda-bayar bersandar pada gerbang bertanda sama dengan TUNGGAL** yang ⛔ belum dapat dipastikan apakah ia pembandingan atau penugasan. ⚠️ Pemberitahuan galat pun terkirim **setiap kali**, sebab gerbangnya ber-bendera mati.

Sesudah tiket ini, ⭐ Instruksi pembayaran terkirim **tepat satu kali**, penandanya tersimpan **dalam transaksi yang sama**, dan ⭐ **pemberitahuan galat terkirim HANYA ketika pengiriman gagal**.

## Perilaku Pega yang ditiru

| Yang dibaca | Rule |
| --- | --- |
| Jalur kasir | 53 langkah, sarang **6 tingkat**, ⛔ **nol parameter** — rule terdalam di modul |
| Dua panggilan luar | keduanya bergerbang **lingkungan produksi**; ⭐ keduanya punya **lompatan kegagalan** — satu-satunya penanganan kegagalan eksplisit di modul |
| ⚠️ Penjaga ganda-bayar | ⛔ gerbangnya memakai tanda sama dengan **tunggal** |
| ⚠️ Pemberitahuan galat | ⛔ gerbang **ber-bendera mati** ⇒ terkirim **setiap kali**; ⚠️ medannya **salah eja** |

⭐ Sumber: `komite-claim-facin\spec.md` · `claim-facin\STRUKTUR-TABEL-CLAIM-FACIN.md` §5.

## Keputusan work owner yang mengikat

- **K5** — ⭐ **Penjaga ganda-bayar dibuat EKSPLISIT** — penanda tersimpan **dalam transaksi yang sama**; panggilan kedua **ditolak**. ⚠️ `[penyimpangan sadar]`
- **K6-lama** — ⭐ **Pemberitahuan galat HANYA pada kegagalan** — ⛔ perilaku lama **tidak ditiru**. ⚠️ `[penyimpangan sadar]`

## Yang harus diuji

- [ ] ⭐ Pembayaran terkirim **tepat satu kali** per akseptasi
- [ ] ⭐ Penanda terkirim tersimpan **dalam transaksi yang sama** dengan pengirimannya
- [ ] ⛔ Panggilan kedua atas akseptasi yang sama **ditolak**
- [ ] ⭐ Pemberitahuan galat terkirim **hanya pada kegagalan**
- [ ] ⭐ Medan tanggapan layanan **dinamai dengan benar** — ⚠️ di sistem lama salah eja
- [ ] ⭐ Kegagalan pengiriman **tidak membatalkan** akseptasi yang sudah tersimpan; ia **dicatat** dan **dapat diulang**
- [ ] ⛔ Rancangan **tidak bergantung** pada tafsir tanda sama dengan tunggal

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **26** | ⚠️ **Tanda sama dengan TUNGGAL** pada penjaga ganda-bayar — pembandingan atau penugasan | ⛔ **TIDAK menahan** — ⭐ K5 mengurungnya; jawabannya dipakai untuk **memeriksa DATA LAMA**, yakni apakah sistem lama pernah membayar dua kali |
| **27** | Apakah pemberitahuan galat lama **benar-benar sampai** ke seseorang | ⚠️ tidak menahan — ⭐ bila ternyata tidak, penanganan galat lama **sebenarnya tidak ada** |

## Seam & verifikasi

**Seam:** lapisan layanan komite — pengiriman ke kasir.
⛔ **Uji sebagai PERILAKU**, bukan sebagai isi kolom penanda.
1. Kirim satu instruksi ⇒ ⭐ terkirim, penanda terisi **dalam transaksi yang sama**.
2. Kirim ulang untuk akseptasi yang sama ⇒ ⛔ **ditolak**.
3. Paksa kegagalan ⇒ ⭐ pemberitahuan galat **terkirim**, akseptasi **tetap tersimpan**.
4. Kirim yang **berhasil** ⇒ ⛔ pemberitahuan galat **TIDAK terkirim** — ⚠️ di sistem lama ia terkirim.
5. Jalankan di lingkungan bukan produksi ⇒ ⛔ panggilan luar **tidak dijalankan**.
