---
status: aktif
---

# 49: Section head menyetujui versi yang menunggunya, dan versinya berpindah ke antrian dept head

*Asal: `DAFTAR-PEKERJAAN.md` `P-31` · `ADR-0055` (`SETUJUI`) · `ADR-0044`.*

**What to build:** **SH** menyetujui versi yang berada di antriannya; versinya berpindah ke antrian **DH**,
dan keputusannya meninggalkan satu catatan persetujuan.

Pecahan **pertama dari tiga**. `U-2` melarang satu tiket melintasi lebih dari satu perpindahan, dan
`GRL-08` menulis `SETUJUI×3` — **tiga perpindahan berbeda**, bukan satu perpindahan berparameter.

**Persyaratan:** `ADR-0055` (`SETUJUI` tingkat 1) · `ADR-0044` (keadaan menjawab *"mungkinkah"*, peran menjawab *"boleh oleh siapa"*, dan **setiap larangan punya tepat satu sebab**)

**Tidak termasuk:** **Tingkat DH dan DR** — tiket `50` dan `51`.
**Larangan pengaju menyetujui sendiri** — tiket `52`; itu sebab yang **berbeda** dari wewenang
tingkat, dan `ADR-0044` menuntut tiap penolakan punya tepat satu sebab.

**Jalur gagal:** Orang tanpa peran SH menyetujui -> ditolak **karena peran** · SH menyetujui versi yang
**tidak** di antriannya -> ditolak **karena keadaan**. Kedua pesan **berbeda**, dan itu uji
`ADR-0044`.

**Uji:** **Negatif:** peran salah; keadaan salah; keduanya salah sekaligus — dan pesannya tetap
menyebut **satu** sebab.
**Positif:** SH yang sah menyetujui versi yang sah -> berpindah ke antrian DH, **dan catatannya
ada**.

**Menggantikan:** `TDA-09` — *peran dibaca dari `pyWorkBasketList(2)`, `pyTelephone`, dan nama orang yang
ditanam langsung di aturan.* Empat nama orang tersemat di penyaluran persetujuan yang **hidup sejak
2019**; orang baru tidak dapat menerima tugas sampai ruas teleponnya disetel, dan langkah itu tidak
tertulis di mana pun.

**Blocked by:** `48` · `54`

**Dasar:**
```
EVIDENCED(Akseptasi_DT@ekspor-2026-09 - nama orang sebagai tetapan, TreatyInSetValue 2.6 PositionUsername="BERNARD")
        DECIDED(ADR-0055, ADR-0044)
```

- [ ] SH menyetujui, versi berpindah ke antrian DH
- [ ] penolakan karena **peran** dan penolakan karena **keadaan** menghasilkan pesan yang **berbeda**
- [ ] satu catatan persetujuan tertulis, berisi pelaku, tingkat, dan waktunya
- [ ] **nol nama orang** di dalam aturan mana pun — peran, bukan nama
