---
status: tertahan
---

# 58: Wewenang persetujuan dibatasi nilai kontrak menurut batas yang berlaku

*Asal: `DAFTAR-PEKERJAAN.md` `P-41` · `ADR-0044` · eskalasi manajemen butir 5.*

**What to build:** Penyetuju yang wewenangnya **di bawah nilai kontrak** ditolak, dan pesannya menyebut
batasnya.

**Persyaratan:** `ADR-0044` (peran menjawab *"boleh oleh siapa"*) · dokumen batas wewenang yang berlaku

**Tidak termasuk:** **Pemisahan pelaku** — tiket `52`; itu sebab yang berbeda, dan `ADR-0044` menuntut tiap penolakan punya tepat satu sebab.

**Jalur gagal:** Penyetuju berwenang di bawah nilai kontrak -> ditolak, pesannya menyebut **batasnya**, bukan hanya "tidak berwenang".

**Uji:** **Negatif:** nilai tepat di atas batas tiap tingkat. **Positif:** nilai **tepat pada** batas diterima — batas inklusif atau eksklusif adalah bagian dari jawaban yang ditunggu.

**Menggantikan:** Sistem lama **tidak punya batas wewenang sama sekali** — sapuan menemukan **nol nama hak
akses terisi**; mekanisme peran bawaan Pega hadir 2.202 kali sebagai tempat kosong dan tidak pernah
dipakai sekali pun (`ADR-0044` Konteks).

**Blocked by:** `51` · **dokumen batas wewenang**

**Dasar:**
```
EVIDENCED(sapuan hak akses@ekspor-2026-09 - 2.202 tempat kosong, nol terisi)
        DECIDED(ADR-0044)
```

- [ ] batas wewenang dibaca dari **tabel acuan**, bukan ditanam di kode
- [ ] penolakan menyebut batasnya
- [ ] **PENGHALANG:** dokumen batas wewenang yang berlaku · **siapa menjawab:** manajemen, eskalasi butir 5 · **yang berubah:** angka ambangnya, dan apakah batasnya inklusif
