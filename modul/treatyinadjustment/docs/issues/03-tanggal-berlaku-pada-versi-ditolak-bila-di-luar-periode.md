---
status: aktif
---

# 03: Tanggal berlaku pada versi ditolak bila jatuh di luar periode kontrak

*Asal: `DAFTAR-PEKERJAAN.md` `P-45` (diubah 24 Sep 2026) · `INV-54` · `GRL-15` · `KTV-2`.*

**What to build:** **PK** menyatakan sejak kapan sebuah versi berlaku, dengan bawaan **tanggal mulai
kontrak**. Tanggal yang jatuh di luar periode kontraknya **ditolak**.

Artefak: atribut tanggal berlaku pada `VERSI_KONTRAK`, bawaannya, dan trigger `INV-54`.

**Persyaratan:** `INV-54` (trigger) · `GRL-15`

**Tidak termasuk:** **Mesin pro rata — sengaja TIDAK dibangun.** Ini pernyataan keputusan, bukan `TODO`:

| | |
|---|---|
| apa yang sengaja tidak dilakukan | rumus yang membagi besaran versi menurut porsi periode tersisa |
| kenapa | sistem lama **tidak pernah menjalankannya sungguhan** — `EDMEffective` punya satu penulis yang menyamakannya dengan tanggal mulai, dan kontrol layarnya hanya-baca, sehingga `ProRatePercent` **tidak pernah dapat selain 100** |
| akibatnya pada angka | **nol** — setiap versi berlaku penuh, sama seperti sistem lama |
| di mana selisihnya terlihat | `UA-19` — sebaran waktu baris ber-`ProRatePercent` bukan 100 |
| kapan ditagih | saat `DB-5` dijawab |

**Tanggal berlaku pada DOKUMEN** adalah ruas yang berbeda — irisan 08.

**Jalur gagal:** Tanggal berlaku versi sebelum tanggal mulai kontrak atau sesudah tanggal berakhir -> **ditolak** trigger, pesannya menyebut periode kontraknya.

**Uji:** **Negatif:** tanggal berlaku sebelum mulai; sesudah berakhir; tepat di batas luar.
**Positif:** tanggal berlaku **tepat pada** tanggal mulai dan tepat pada tanggal berakhir
**diterima** — batas inklusif, dan uji negatif saja tidak pernah membuktikannya.

**Menggantikan:** `P-45` bunyi lama: *"**Tanggal berlaku addendum** yang jatuh di luar periode kontraknya
**ditolak**."*

Yang bergeser: **pembawanya disebut** — *"tanggal berlaku **pada versi**"*. `KTV-2` memberi dokumen
addendum tanggal berlakunya sendiri, sehingga tanpa menyebut pembawanya irisan ini dan irisan 08
**sama-sama mengklaim ruas bernama sama**, dan yang menang adalah yang ditulis lebih dulu.

**Blocked by:** **14** — *ditambahkan 24 September 2026.* Tiket ini menambahkan **kolom** atau **tabel anak** pada `VERSI_KONTRAK`, dan tidak satu pun tiket di papan membuat tabelnya. Irisan `14` adalah **PEMBUAT PERTAMA** `KONTRAK` dan `VERSI_KONTRAK`.

**Dasar:**
```
EVIDENCED(TreatyInSetEditPre@ekspor-2026-09, TreatyCalculateProratePct@ekspor-2026-09)
        DECIDED(GRL-15, KTV-2, ADR-0037)
```

- [ ] atribut tanggal berlaku berdiri pada `VERSI_KONTRAK` dengan bawaan tanggal mulai kontrak
- [ ] trigger `INV-54` menolak tanggal di luar periode, pesannya menyebut periodenya
- [ ] uji positif batas inklusif lulus
- [ ] pernyataan keputusan **mesin pro rata tidak dibangun** tertulis di tiket ini dan di §10.2, bukan sebagai `TODO`
