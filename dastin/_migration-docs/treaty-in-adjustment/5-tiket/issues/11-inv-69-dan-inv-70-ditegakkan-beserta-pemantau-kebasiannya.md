---
status: aktif
---

# 11: INV-69 dan INV-70 ditegakkan atas baris selisih, beserta pemantau kebasiannya

*Asal: `DAFTAR-PEKERJAAN.md` `P-60` · `INV-69`, `INV-70` · `UJI-NEGATIF-INVARIAN.md` §2 dan §3.*

**What to build:** Versi yang **dinyatakan `TIDAK_MATERIAL`** tetapi memiliki baris `NILAI_SELISIH` bertipe
**uang** atau **porsi** **ditolak saat simpan**. Penegakannya lewat *materialized view*
ber-`REFRESH ON COMMIT`.

Dan karena mekanisme itu **dapat berhenti tanpa memberi tahu siapa pun**, irisan ini berdiri hanya
bila **pemantau kebasiannya** juga berdiri.

**Persyaratan:** `INV-69` · `INV-70` · `UJI-NEGATIF-INVARIAN.md` §2 (uji negatif dijalankan, bukan diargumentasikan) dan §3 (pemantau kebasian)

**Tidak termasuk:** **Penguncian ruas** — irisan 02; ia penegakan di lapisan yang berbeda.
**Pemantau kebasian untuk `INV-47`, `INV-50`, `INV-51`** — itu **lubang induk**, `F-13`. Irisan ini
menutup dua invarian dan **membiarkan tiga**; itu dinyatakan, bukan disamarkan.

**Jalur gagal:** Versi `TIDAK_MATERIAL` dengan baris selisih bertipe porsi -> **ditolak**, pesannya menyebut
besaran mana yang melanggarnya · **MV gagal me-refresh** -> penegakan berhenti **tanpa galat**, dan
**hanya pemantau kebasian yang memperlihatkannya.**

**Uji:** **Negatif — dijalankan, bukan diargumentasikan:** simpan versi `TIDAK_MATERIAL` dengan baris
selisih bertipe uang; lalu bertipe porsi.
**Positif:** versi `TIDAK_MATERIAL` yang baris selisihnya **hanya** bertipe lain **diterima**.
**Uji pemantau:** matikan refresh MV-nya dengan sengaja, dan pastikan pemantau **berbunyi** — sebuah
pemantau yang belum pernah dilihat berbunyi belum terbukti ada.

**Menggantikan:** `TDA-10` bagian kedua — *nol penegakan di sisi simpan.* Dan di sistem lama versi yang
`TIDAK_MATERIAL` tetap dapat mengubah angka uang lewat jalur yang tidak melewati layar.

**Blocked by:** 02, 06

**Dasar:**
```
EVIDENCED(sensus-kondisi-penguncian@ekspor-2026-09 - lihat GRILL-D/01-TEMUAN TD-02)
        DECIDED(GRL-20, INV-69, INV-70)
        DIASUMSIKAN-CLEAR(DB-20)
```

- [ ] constraint `INV-69` dan `INV-70` terpasang di atas *materialized view*-nya
- [ ] **uji negatif dijalankan** dan gagal sebagaimana seharusnya — hasilnya dilampirkan, bukan dinyatakan
- [ ] **pemantau kebasian berdiri**: `REFRESH_MODE`, `STALENESS`, terjadwal, dan **berbunyi kepada seseorang yang bernama**
- [ ] pemantau diuji dengan mematikan refresh secara sengaja, dan ia **berbunyi**
- [ ] `F-13` dirujuk: tiga MV induk masih **tanpa** pemantau, dan itu bukan bagian tiket ini
