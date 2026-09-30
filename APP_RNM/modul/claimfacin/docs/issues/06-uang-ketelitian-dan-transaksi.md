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

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **1** | Kolom tabel data kutipan — `[data DBA]` | tidak menahan |

## Perintah verifikasi

1. Hitung total penyesuaian dua kali dengan urutan pemanggilan berbeda — ⭐ **hasilnya sama**.
2. Paksa gagal di tengah sekelompok tulisan — ⭐ **tak satu pun baris tersimpan**.
3. Periksa konversi ke rupiah — ⭐ ketelitiannya sesuai AC.
