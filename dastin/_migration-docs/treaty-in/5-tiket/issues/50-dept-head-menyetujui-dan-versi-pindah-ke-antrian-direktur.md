---
status: aktif
---

# 50: Dept head menyetujui versi yang menunggunya, dan versinya berpindah ke antrian direktur

*Asal: `DAFTAR-PEKERJAAN.md` `P-32` · `ADR-0055` (`SETUJUI` tingkat 2).*

**What to build:** **DH** menyetujui versi di antriannya; versinya berpindah ke antrian **DR**, dan
keputusannya meninggalkan catatan.

Pecahan **kedua dari tiga**. Begitu `49` mendarat, tiket ini hampir gratis — dan itu alasan ia
**murah**, bukan alasan ia digabung.

**Persyaratan:** `ADR-0055` (`SETUJUI` tingkat 2) · `ADR-0044`

**Tidak termasuk:** **Tingkat SH dan DR** — tiket `49` dan `51`. **Wewenang bernilai** — tiket `58`.

**Jalur gagal:** DH menyetujui versi yang masih di antrian SH -> ditolak karena **keadaan** · SH mencoba menyetujui di tingkat DH -> ditolak karena **peran**.

**Uji:** **Negatif:** peran salah; keadaan salah.
**Positif:** DH yang sah -> antrian DR, catatan ada, dan **catatan tingkat SH tetap utuh** — tiket
ini tidak boleh menimpa jejak tingkat sebelumnya.

**Menggantikan:** `TDA-09`, dan `TDA-08` — *rantai dipendekkan `RevisionState = 1` menjadi dua tingkat, tanpa ambang nilai apa pun.*

**Blocked by:** `49`

**Dasar:**
```
EVIDENCED(Akseptasi_DT@ekspor-2026-09 - cabang 2 Admin->SecHead->selesai)
        DECIDED(ADR-0055, ADR-0044)
```

- [ ] DH menyetujui, versi berpindah ke antrian DR
- [ ] catatan tingkat sebelumnya **tetap utuh**
- [ ] penolakan peran dan penolakan keadaan berpesan berbeda
