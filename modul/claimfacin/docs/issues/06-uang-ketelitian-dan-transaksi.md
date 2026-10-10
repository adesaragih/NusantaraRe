# 06: Uang, ketelitiannya, dan batas transaksi

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR) · 05 (penyesuaian)
**Menutup:** AC 39 · 40 · 41 · 42 · 65 · 66 · 67 · 114 *(8 AC)* — US 7 · 29

## Hasil & nilai pengguna

Hari ini Nilai uang **berubah karena urutan pemanggilan**, dan tidak ada batas transaksi yang menjamin sekelompok tulisan selesai bersama atau gagal bersama.

Sesudah tiket ini, Angka uang **berhenti berubah karena urutan**, konversi mata uang punya ketelitian yang ditetapkan, dan sekelompok tulisan **selesai bersama atau gagal bersama**.

## Area codebase

- Lapisan layanan: perhitungan uang dan konversi
- Batas transaksi pada penulisan berkelompok

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Perhitungan uang | penjumlahan bruto, bagian Nusantara Re, dan konversi ke rupiah |
| Batas transaksi | ⚠️ rule SQL lama membawa penyelesaian transaksi **di tengah** pekerjaan pemanggil |

⭐ Rincian medan dan asalnya ada di `STRUKTUR-TABEL-CLAIM-FACIN.md` **§2b**.

## ADR terkait

- **ADR-0007** — jejak audit

## Acceptance criteria

- [ ] **AC 39–42 · 114** — uang dan ketelitiannya
- [ ] **AC 65 · 66 · 67** — transaksi
- [ ] ⭐ Kegagalan separuh jalan **tidak meninggalkan data separuh tersimpan**

> ⛔ **RALAT 10-10-2026.** Butir lamanya dikutip utuh, tidak dihapus: *"**AC 39–42 · 114** — uang dan ketelitiannya"* dan
> *"**AC 65 · 66 · 67** — transaksi"* → tiga AC yang ditutup tiket ini **diralat di spec**:
>
> - **AC 39 / 42** — ketelitian **`NUMBER(38,10)`** (uang, persen, share, kurs; migrasi `560`–`567`), bukan 20,8;
>   hitungan lewat `inti/backend/uang` / `apd.Decimal`, nol `float`. ADR-0016 (`NUMBER(38,8)`) usang untuk modul ini.
> - **AC 41** — kurs disimpan **hanya di tempat XML menulisnya** (`KURS` estimasi / adjustment, `KURS_OBJECT_ITEM`),
>   bukan di samping setiap nilai uang.
> - **AC 67** — pengecualian penomoran **hilang**: nomor `inti/backend/penomor` terbit di dalam transaksi aksi dan ikut
>   batal bila aksinya gagal (ADR-0043 meng-*supersede* ADR-0006). Satu aksi = satu transaksi, nol `COMMIT` di SQL.
>
> ⚠️ Kalimat *"Hari ini Nilai uang **berubah karena urutan pemanggilan**"* dan perintah verifikasi 1 tidak punya dasar di
> spec — **perlu dicek terhadap XML** sebelum dijadikan sasaran uji.

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **1** | Kolom tabel data kutipan — `[data DBA]` | tidak menahan |

## Perintah verifikasi

1. Hitung total penyesuaian dua kali dengan urutan pemanggilan berbeda — ⭐ **hasilnya sama**.
2. Paksa gagal di tengah sekelompok tulisan — ⭐ **tak satu pun baris tersimpan**.
3. Periksa konversi ke rupiah — ⭐ ketelitiannya sesuai AC.
