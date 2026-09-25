---
status: aktif
---

# 59: Kepala kontrak warisan dan seluruh daftar anaknya dipindahkan apa adanya, beserta keadaannya

*Asal: `DAFTAR-PEKERJAAN.md` `P-50` (pecahan selain tabel acuan) · `ADR-0042` · `ADR-0043`.*

**What to build:** **PM** memindahkan kepala kontrak lama dan seluruh daftar anaknya **apa adanya, tanpa
menghitung ulang apa pun**, dan menulis keadaannya ke kolom yang sudah punya himpunan nilai sah.

**Persyaratan:** `ADR-0042` (pindah apa adanya, penanda warisan, sentuh-perbaiki) · `ADR-0043` (setiap angka identik) · `KTV-A`

**Tidak termasuk:** **Keadaan yang tidak berpadanan** — tiket `60`.
**Isi tabel acuan** — sudah dipindahkan tiket `44`.
**Menghitung ulang** apa pun — `ADR-0043` melarangnya; itu peristiwa bisnis tersendiri.

**Jalur gagal:** Satu angka berbeda antara sumber dan hasil -> **migrasi gagal dan diulang**; bukan
"dalam toleransi" · Jumlah baris anak per kontrak tidak cocok -> gagal, dan **ini kegagalan yang
paling sering lolos ketiga ukuran pertama**.

**Uji:** **Negatif:** kontrak yang kehilangan satu baris layer **harus** terdeteksi.
**Positif:** kontrak yang pindah utuh lolos **keempat** ukuran `ADR-0043` — cacah kontrak, pasangan
ID dua arah, nilai kunci identik, **dan jumlah baris anak per kontrak**.

**Menggantikan:** Tidak ada perilaku lama yang digantikan — ini **pemindahan**. Yang digantikan adalah
ketiadaannya: data lama tinggal di `M_TREATY_IN.JSONDATA`, **dua kolom saja**, tanpa kolom waktu.

**Blocked by:** `45` · `44`

**Dasar:**
```
EVIDENCED(M_TREATY_IN DDL@ekspor-2026-09 - dua kolom, ID dan JSONDATA, tanpa kolom waktu)
        DECIDED(ADR-0042, ADR-0043)
        DIASUMSIKAN-CLEAR(KTV-A)
```

- [ ] kepala dan seluruh daftar anak pindah; **jumlah baris anak per kontrak dicocokkan**
- [ ] baris hasil migrasi membawa **penanda warisan** yang terbaca siapa pun
- [ ] nol angka dihitung ulang — diperiksa, bukan diandaikan
- [ ] **bertenggat `KTV-A`:** presisi kolom **tidak dapat dipersempit** sesudah tiket ini berjalan
