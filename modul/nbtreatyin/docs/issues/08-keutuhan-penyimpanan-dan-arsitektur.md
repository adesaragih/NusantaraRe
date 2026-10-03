# 08: Keutuhan penyimpanan — satu transaksi, skema eksplisit, arah ketergantungan

**Status:** sebagian — dapat dikerjakan: uji bertag `db` kegagalan `CatatRiwayat` membatalkan submit lawan Oracle (AC 83) belum ditulis; penahan pihak luar: K11 (AC 29, 83) *(putaran 2, konsolidasi P10 04-10-2026 — rincian `docs/HASIL-IMPLEMENTASI.md` bab 9; semula: sebagian, implementasi 2026-10-03; awalnya ready-for-agent)*
**Blocked by:** 01
**Menutup:** AC 29 · 30 · 60 · 83 · 90 *(5 AC)* — US 33 · 34

## Hasil & nilai pengguna

Hari ini penyimpanan realisasi treaty **tidak punya jaminan keutuhan**. `[data DBA]` Dua dari tiga
program penyimpan **menyelesaikan penyimpanannya sendiri**, dan blok pemanggilnya menyelesaikannya
lagi — ⛔ **penyelesaian ganda**. Akibatnya kegagalan pada langkah berikutnya **meninggalkan data
yang sudah permanen**, dan aplikasi tidak dapat membatalkannya.

Sesudah tiket ini, seluruh urutan penyimpanan berada dalam **satu transaksi** — ⭐ kegagalan di
langkah mana pun membatalkan seluruhnya, dan tidak ada lagi berkas yang tersimpan separuh.

## Area codebase

- Lapisan repository: batas transaksi
- Lapisan service: urutan penyimpanan
- Uji arsitektur: arah ketergantungan

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Penyelesaian ganda | `RDBList\SaveTreatyIn.xml` · `RDBList\SavePolisTreatyIn_SQL.xml` — blok memuat penyelesaian, dan program yang dipanggil **juga** |
| Yang **tidak** menyelesaikan sendiri | `RDBList\GetSequenceNumber_SQL.xml` |
| Dua ejaan nama objek | empat nama muncul **dengan dan tanpa** awalan skema di 1.151 naskah SQL korpus |

## ADR terkait

- **ADR-0009** — migrasi penuh; tidak ada koeksistensi dua penulis

## Acceptance criteria

- [ ] 🟡 **AC 29** — seluruh urutan penyimpanan berada dalam **satu transaksi**; kegagalan di tengah
      menyisakan **nol** baris
- [ ] 🟡 **AC 83** — kegagalan menyimpan riwayat **membatalkan seluruh transaksi**
- [x] **AC 30** — setiap query menyebut **skema secara eksplisit**
- [x] **AC 90** — keempat nama berejaan ganda diperlakukan sebagai **satu objek**
- [x] **AC 60** — arah ketergantungan `handlers → services → repository`; ⛔ tidak terbalik, tidak
      memotong lapisan

## Perintah verifikasi

1. Suntikkan kegagalan di tengah urutan penyimpanan — ⭐ **nol** baris tersisa.
2. Sambung sebagai pengguna **selain** pemilik skema — ⭐ query **tetap menemukan** objeknya.
3. Jalankan uji arsitektur — ⭐ **nol** panggilan dari repository ke service.

## Catatan

⚠️ `[penyimpangan sadar]` Satu transaksi **lebih ketat** daripada sistem lama, dan itu disengaja —
alasannya tertulis di `spec.md` §5.7.
⛔ Uji keutuhan transaksi **tidak dapat dijalankan dengan tiruan**; ia wajib berjalan lawan basis
data sungguhan.

## ⭐ Putaran 2 — paket penyimpanan (03-10-2026)

Dasar: PROMPT-NB-TREATY-IN-PUTARAN-2 bab 0 butir 11–12, bab 2 K4/K16/K17; rincian kolom `docs/PERBANDINGAN-KOLOM-DIAGRAM.md`.

- Penulisan catatan usulan ke `POOLDATA.HISTORYAKSEPTASIPRODUCTION` (K4) masuk **transaksi submit yang
  sama** dengan riwayat, halaman, dan perpindahan — `repository/usulan.go` tanpa `COMMIT`, skema lewat
  `Qualify` (AC 30). Uji seam HTTP `TestSatuTransaksiPembatalanUtuh` kini juga menyuntikkan kegagalan di
  `CatatUsulan` — nol baris riwayat produksi, riwayat, halaman, atau perpindahan tersisa (AC 29, 83).
- Gudang tiruan kini hanya menyimpan medan yang punya kolom (`models.ProyeksiKatalog`) — uji seam HTTP
  melihat penyimpanan yang sama dengan Oracle. Celah yang terungkap dan ditutup: `LAYER*` tingkat polis
  hilang sesudah dibuka ulang (kini dibaca balik dari view, ID-22).
- AC 29 dan 83 tetap 🟡: keutuhan lawan Oracle sungguhan (`-tags db`) belum dijalankan (K11 kosong).
