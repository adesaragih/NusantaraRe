# F12: Layar dan daftar Fac Out

**What to build:** Pengguna dapat membuka menu Fac Out, melihat daftar retrosesi sebuah polis, dan
menelusuri objek serta coverage yang diretrokan — sebagai menu **terpisah** dari Fac In namun
**tertaut** ke kasus induknya (K-053).

`[terverifikasi]` Inventaris berkas Fac Out/FacRetro per tipe rule, ketiga folder digabung, dihitung
**per subfolder** (listing rekursif terpotong di batas entri):

| Tipe | Jumlah |
| --- | ---: |
| Section | 131 |
| Activity | 120 |
| FlowAction | 57 |
| DataTransform | 12 |
| RDBList | 10 |
| When | 9 |
| Harness | 9 |
| DataPage | 6 |
| Flow | 5 |
| **TOTAL** | **359** |

Nama unik: **46 Section · 19 FlowAction · 3 Harness · 4 RDBList · 2 Flow · 2 DataPage**.

`[terverifikasi]` **Daftar Fac Out** dibaca dari `RDBList\GetFacoutList_SQL` (`Rule-Connect-SQL`, tag
`pyBrowseSQL`):

```sql
select IDPEGA as HASIL1, policyno HASIL2 from facoutproduction where policyno = {InputData.CARI17}
```

Jadi daftarnya bersumber dari **tabel produksi**, berkunci nomor polis — bukan dari staging.

## ⛔ `GetOPFacOut_Sql` — TIDAK diport, **DIHAPUS** (K-062)

✅ **K-062 sudah diputuskan.** Klep ini **tidak diport dan dihapus, tidak dipakai lagi**. **Kolom
pelaksana TETAP ADA di layar Fac Out, tetapi ISINYA KOSONG.**

⛔ **Ini bukan blocker eksternal lagi.** Tidak ada yang perlu ditunggu.

### ⚠️ Pengecualian yang disengaja terhadap K-006

> **K-062 adalah pengecualian PERTAMA yang disengaja** terhadap pola "ditangguhkan, bukan dihapus"
> (**K-006**, `PANDUAN-KERJA` §5). Pola baku menuntut `panic` bila cabang tercapai; **di sini tidak**.
>
> **Sesi berikutnya JANGAN membacanya sebagai inkonsistensi lalu membalikkannya.** Membalikkannya
> memerlukan keputusan baru.
>
> Alasannya dapat dipertanggungjawabkan: yang hilang adalah **tampilan identitas pelaksana**, bukan
> angka uang atau gerbang alur — mengosongkan kolom **tidak** menghasilkan selisih rekonsiliasi
> (ADR-0001).

### Sisi pemanggil wajib menghasilkan KOSONG, bukan galat

`[terverifikasi]` Pemindaian **peka huruf** atas **seluruh 6.071 berkas** (14 subfolder × 3 folder):
**10 berkas** memuat `GetOPFacOut`, **3** di antaranya berkas definisinya sendiri →
**7 perujuk sejati**.

| Perujuk | NB | RNW | EDM | Tag pembawa |
| --- | :-: | :-: | :-: | --- |
| `Activity\OfferFacOut_PreAct` | ✓ | ✓ | ✓ | `<RequestType>` = `GetOPFacOut_Sql` · `<pyRuleName>` = `ASM-FW-GISFW-Work ASM GetOPFacOut_Sql` |
| `Activity\PrintRISlipPre_act` | ✓ | ✓ | ✓ | idem |
| `Activity\SumCurrencyListAllRetro_Act` | — | — | ✓ | idem |
| `RDBList\GetOPFacOut_Sql` *(definisi)* | ✓ | ✓ | ✓ | `<pyRequestType>` · `<pyLabel>` · `<pxTabLabel>` · `<pyRuleName>` |

⛔ **Ketiga pemanggil itu wajib menghasilkan nilai kosong, BUKAN galat.** `PrintRISlipPre_act`
menyentuh jalur cetak RI Slip (**F13**); `SumCurrencyListAllRetro_Act` hanya ada di Endorsement.

`[terverifikasi]` Yang dibacanya: identitas pelaksana dari riwayat penugasan kasus, disaring pada
langkah bernama `OfferFacOut`/`OfferRetro`, diurutkan menurut waktu commit, dari
**`datapega.pc_history_asm_fw_gisfw_work`** — tabel **internal Pega**. `CLAUDE.md` §4.3:
`DATAPEGA.PC_*` **hilang bersama Pega**. ⚠️ Nilainya **identitas orang** — mekanismenya saja yang
dicatat (K-025).

⚠️ **Nama folder bukan tipe rule.** `[terverifikasi]` `CekFacoutProd_Sql`, `GetOPFacOut_Sql` dan
`GenerateOurRefFacOut_SQL` berada di folder `RDBList\` tetapi ber-`<pxObjClass>` = `Rule-Connect-SQL`
— seluruh **605** berkas di `RDBList\` berkelas itu. Begitu pula seluruh **443** berkas di
`DataTransform\` berkelas `Rule-Obj-Model`. Tipe dibaca dari tag, tidak pernah dari nama folder.

⚠️ `[terverifikasi]` Section Fac Out hadir **berpasangan** `X` dan `X_IsUW` (mis.
`ShowCoverageFacOut` / `ShowCoverageFacOut_IsUW`). Sejalan keputusan work owner pada siklus Renewal,
pasangan `_IsUW` adalah **dua komponen terpisah**, bukan satu komponen dua mode.

**Asal (Pega).** `Harness\ViewFacretro` · `ViewCoverageFacOut` · `ViewCoverageFacOut_isUW` ·
46 `Section\*FacOut*` · 19 `FlowAction\*FacOut*` · `RDBList\GetFacoutList_SQL` ·
`RDBList\GetOPFacOut_Sql` · `DataPage\D_CoverageFacOut` · `D_FacOutFromVehicle`

**Keputusan.** **K-062** (klep dihapus; kolom pelaksana kosong — ⛔ **pengecualian pertama terhadap
K-006**) · K-053 · `CLAUDE.md` §4.3, §4.6

**Blocked by:** F01 · F05

**Status:** blocked

- [ ] Menu Fac Out **terpisah** dari menu Fac In di navigasi, tetapi tertaut lewat penaut objek induk
- [ ] Daftar Fac Out dibaca dari **tabel produksi** berkunci nomor polis
- [ ] ⛔ `GetOPFacOut_Sql` **dihapus, tidak dipakai** (K-062) — **bukan** `panic`, **bukan** ditangguhkan
- [ ] Kolom pelaksana **tetap ada** di layar Fac Out dengan **isi kosong**
- [ ] Ketiga pemanggil (`OfferFacOut_PreAct`, `PrintRISlipPre_act`, `SumCurrencyListAllRetro_Act`) menghasilkan **nilai kosong, bukan galat**
- [ ] Komentar menyatakan sebabnya **dan** menandai ini **pengecualian sadar terhadap K-006**, agar tidak dibalikkan
- [ ] Pasangan `X` / `X_IsUW` tetap **dua komponen terpisah**
- [ ] Tipe rule tiap berkas dibaca dari `<pxObjClass>`, **tidak pernah** dari nama folder
- [ ] Layar **tidak menampilkan** nilai `DOB`, plat nomor, nomor mesin, atau alamat di test/fixture (K-025)
- [ ] Endpoint dari konfigurasi, tidak pernah literal (§4.4)
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)
