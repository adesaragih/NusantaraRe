# 08: Keutuhan penyimpanan — satu transaksi, skema eksplisit, arah ketergantungan

**Status:** sebagian *(implementasi 2026-10-03, cabang `modul/nbtreatyin/implementasi`; semula: ready-for-agent)*
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
