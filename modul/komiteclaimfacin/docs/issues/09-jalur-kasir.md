# 09: Jalur kasir — penjaga ganda-bayar eksplisit

**Status:** ready-for-agent
**Blocked by:** 06 · `claim-facin\issues\11` *(efek keluar klaim)*
**Menutup:** AC 61 · 62 · 63 · 64 · 65 · 66 · 67 · 68 · 69 · 70 *(10 AC)* — US 30 · 31

## Hasil & nilai pengguna

Hari ini Instruksi pembayaran **belum terkirim ke kasir**, dan ⚠️ di sistem lama **penjaga ganda-bayar bersandar pada gerbang bertanda sama dengan TUNGGAL** yang ⛔ belum dapat dipastikan apakah ia pembandingan atau penugasan. ⚠️ Pemberitahuan galat pun terkirim **setiap kali**, sebab gerbangnya ber-bendera mati.

Sesudah tiket ini, ⭐ Instruksi pembayaran terkirim **tepat satu kali**, penandanya tersimpan **dalam transaksi yang sama**, dan ⭐ **pemberitahuan galat terkirim HANYA ketika pengiriman gagal**.

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"⚠️ Pemberitahuan galat pun terkirim **setiap
> kali**, sebab gerbangnya ber-bendera mati."* dan *"Sesudah tiket ini, ⭐ Instruksi pembayaran terkirim **tepat satu
> kali**, penandanya tersimpan **dalam transaksi yang sama**, dan ⭐ **pemberitahuan galat terkirim HANYA ketika
> pengiriman gagal**."* →
>
> - **Kasir lewat outbox, hanya produksi.** Gerbang XML: `KomitePost_Adjustment` S25.2.1.1 `.DirectToKasir == "true" &&
>   .StatusKasir = ""` (tanda `=` tunggal dibaca **pembandingan**, prompt §5 #2) → `HitServiceToKasirKMT_Act` L2 (gerbang
>   sama) + L3 `getStatusKonversi_Act` (`StatusKonversi = 1`), jalur `CLM` L14.2 (jumlah per `AcceptedNo`, panjang 21 /
>   22). S25.2.1 "HANYA LOOPING 1 KALI": satu panggilan per KMT. `StatusKasir` diisi dari tanggapan (L14.5).
> - **"Tepat satu kali" / "penanda dalam transaksi yang sama" belum terjamin.** Muatan kasir diantre di outbox dalam
>   transaksi Submit, sedangkan `StatusKasir` baru terisi dari tanggapan kasir (`DIRECTTOKASIR_LOG`) sesudah pelaksana
>   outbox menjalankan muatan. Dedupe saat muatan masih antre = **OQ-CFI-26**, belum diputuskan. Gerbang `StatusKasir`
>   kosong tetap dibangun.
> - **"Setiap kali" salah baca.** `HitServiceToKasirKMT_Act` L15 `Exit-Activity` ("exit disini jika tidak ada error")
>   menutup jalur biasa. L16 (label `END`) dan L17 hanya dicapai lewat lompatan kegagalan (`pyOnException` = `END` pada
>   kedua panggilan luar), dan `SendErrorDirectKasir` L2 (pre=true) memeriksa ulang `ReponseCode != "1"`. Jadi XML pun
>   mengirim surel galat hanya pada kegagalan; arah K6 sama dengan XML, bukan penyimpangan. Surelnya lewat outbox, BCC
>   orang dibuang (prompt §6 butir 8).

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

> ⛔ **RALAT 10-10-2026.** Baris lamanya dikutip utuh, tidak dihapus: *"**26** | ⚠️ **Tanda sama dengan TUNGGAL** pada
> penjaga ganda-bayar — pembandingan atau penugasan"* → untuk rancangan, tanda `=` tunggal **dibaca pembandingan**
> (prompt §5 #2, OQ-CFI-03 "kelainan XML diperbaiki"). Yang masih terbuka bukan tafsir tandanya, melainkan **dedupe
> muatan kasir yang masih antre** (OQ-CFI-26, `modul/claimfacin/docs/OQ.md`). Gantungan `claim-facin\issues\11`: muatan
> kasir claimfacin tahap 1 sudah masuk outbox hanya di produksi, tetapi modul itu belum punya pelaksana outbox
> (OQ-CFI-26), jadi panggilan nyata belum terjadi.

## Seam & verifikasi

**Seam:** lapisan layanan komite — pengiriman ke kasir.
⛔ **Uji sebagai PERILAKU**, bukan sebagai isi kolom penanda.
1. Kirim satu instruksi ⇒ ⭐ terkirim, penanda terisi **dalam transaksi yang sama**.
2. Kirim ulang untuk akseptasi yang sama ⇒ ⛔ **ditolak**.
3. Paksa kegagalan ⇒ ⭐ pemberitahuan galat **terkirim**, akseptasi **tetap tersimpan**.
4. Kirim yang **berhasil** ⇒ ⛔ pemberitahuan galat **TIDAK terkirim** — ⚠️ di sistem lama ia terkirim.
5. Jalankan di lingkungan bukan produksi ⇒ ⛔ panggilan luar **tidak dijalankan**.
