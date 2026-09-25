---
status: aktif
---

# 52: Pengaju tidak menyetujui, tiap tingkat orang berbeda, dan tidak ada jalan pintas

*Asal: `DAFTAR-PEKERJAAN.md` `P-36` dan `P-40` · `ADR-0052` · `ADR-0044`.*

**What to build:** Pemberi keputusan yang **sama dengan pengisi kontraknya**, atau **sama dengan pemberi
keputusan tingkat sebelumnya**, **ditolak**. Dan tidak ada jalur yang melompati tingkat.

**Persyaratan:** `ADR-0052` (satu rantai, tanpa pengecualian) · `ADR-0044` (tepat satu sebab per larangan)

**Tidak termasuk:** **Wewenang menurut nilai kontrak** — tiket `58`. Itu sebab ketiga yang berbeda lagi.

**Jalur gagal:** Pengaju menyetujui di tingkat mana pun -> ditolak, pesannya menyebut **ia pengajunya** ·
Orang yang sama menyetujui di dua tingkat berurutan -> ditolak, pesannya menyebut **tingkat mana**
ia sudah memutuskan.

**Uji:** **Negatif:** pengaju = SH; SH = DH; DH = DR; dan pengaju = DR.
**Positif:** **empat orang berbeda** melewati rantai penuh **diterima** — dan rantai tiga orang
dengan pengaju keempat juga diterima.

**Menggantikan:** `TDA-08` dan **`RevisionState == 1`** — *jalur revisi memendekkan persetujuan menjadi satu
tingkat; kepala departemen dan direktur dilewati, **tanpa ambang nilai apa pun**.*

**Blocked by:** `51`

**Dasar:**
```
EVIDENCED(Akseptasi_DT@ekspor-2026-09 - RevisionState=1 memotong rantai jadi dua tingkat)
        DECIDED(ADR-0052, ADR-0044)
```

- [ ] pengaju ditolak di setiap tingkat, pesannya menyebut sebab itu saja
- [ ] orang yang sama di dua tingkat berurutan ditolak
- [ ] uji positif rantai empat orang lulus
- [ ] **tidak ada** jalur yang melompati tingkat — diperiksa terhadap ketiga belas perpindahan
