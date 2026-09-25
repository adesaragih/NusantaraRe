# R06: Kelompok bisnis pada kasus renewal

**What to build:** Sistem menurunkan **kelompok bisnis** sebuah kasus dari kode bisnis penawarannya,
agar klasifikasi kasus renewal konsisten dengan sistem lama.

`[terverifikasi]` `Activity\GetBusinessGroup_Act`, 4 langkah kedalaman penuh:

| Langkah | Metode | Isi |
| ---: | --- | --- |
| 1 | `Property-Set` | `InputData.CARI2` ← `…QuotationData.BusinessCode` |
| 2 | `RDB-List` | → `CariBusinessGID` |
| 3 | `Property-Set` | `ParamBis.CARI10` ← `OutBis.pxResults(1).CARI3` |
| 4 | `Page-Remove` | bersihkan halaman |

`[terverifikasi]` `RDBList\CariBusinessGID` — `<pyBrowseSQL>`:

```sql
select ID, OLDID, NOTE as "Note", GROUPPANEL as "GroupPanel", BusinessGroupID as CARI3
from business where ID = {InputData.CARI2}
```

### ⚠️ Ini kelompok bisnis, BUKAN lini bisnis (COB)

Jangan tertukar. **Kelompok bisnis** diturunkan di sini dari tabel `business`. **Lini bisnis (COB)** —
yang menggerakkan skala rasio per K-018 — tetap ditentukan predikat `IsFire` / `IsPA` / `IsMBU` / …
yang **diwarisi dari NB**. Menyatukan keduanya akan merusak pemilihan skala.

`[terverifikasi]` DecisionTable `BusinessType_DeT` **dirujuk nol kali** di korpus renewal.

### ⛔ Blok `pySaveSQL` DIBUANG — perubahan perilaku yang disengaja

Rule yang sama memuat `<pySaveSQL>` berupa blok PL/SQL yang memanggil `POOLDATA.PROSESCOPY`.
**Blok itu tidak diport** (K-037).

Alasannya, seluruhnya `[terverifikasi]`: prosedur `PROSESCOPY` **tidak ada di basis data** (dikonfirmasi
DBA) · seluruh argumennya **literal keras** · keluarannya ke `dbms_output` yang tidak dibaca aplikasi ·
**tidak berhubungan** dengan `pyBrowseSQL` di rule yang sama.

⚠️ Ini **bukan** porting apa adanya, dan **bukan** perbaikan diam-diam — ia keputusan sadar yang sudah
dicatat bernomor. `<pyBrowseSQL>` **tetap diport apa adanya**.

**Blocked by:** R01

**Status:** ready-for-agent

- [ ] Kelompok bisnis diturunkan: kode bisnis → tabel `business` → `BusinessGroupID`
- [ ] `<pyBrowseSQL>` diport apa adanya, termasuk alias kolomnya
- [ ] **Blok `pySaveSQL` yang memanggil `PROSESCOPY` tidak diport**, dan ketiadaannya dicatat di komentar sebagai keputusan bernomor — bukan dihapus tanpa jejak
- [ ] Kelompok bisnis **tidak** dipakai memilih skala rasio; pemilihan skala tetap lewat resolver COB NB
- [ ] Komentar menyebut activity dan rule RDB asalnya
- [ ] ⚠️ Bila paralel run memperlihatkan selisih yang menunjuk jalur ini, keputusan membuang `pySaveSQL` adalah **tersangka pertama** — catat itu di komentar
