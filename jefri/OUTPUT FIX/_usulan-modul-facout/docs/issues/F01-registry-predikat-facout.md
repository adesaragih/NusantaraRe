# F01: Predikat Fac Out masuk registry

**What to build:** Sistem dapat menjawab dua pertanyaan tentang sebuah kasus — *apakah kasus ini punya
retrosesi keluar?* dan *apakah baris spreading ini membawa penanda Fac Out?* — lewat **registry
predikat yang sudah ada**, tanpa registry kedua dan tanpa seam baru.

⛔ **`IsFacout` TIDAK ditangani di sini.** `[terverifikasi]` Ia **tidak terpanggil sama sekali** di
jalur Fac Out: `NB FacIn\Flow\InputInwardFacultativeOffer.xml` memuatnya **0 kali**, peka huruf maupun
tidak. *(Berkas flow itu `[terverifikasi]` hanya ada di folder **NB**; tidak ada di RNW maupun
Endorsment.)*
Isinya menguji **banding** (`ProposalAcceptStatus = 4`, K-029), bukan fac out. Porting-nya sudah
menjadi lingkup tiket New Business `..\09-predikat-sikap-khusus.md` di bawah K-019 — **jangan
diulang**.

`[terverifikasi]` Dua predikat yang masuk lingkup, beserta kemunculannya di flow: `IsFacRetro` **17×**,
`IsInputFacRetro` **9×**.

### ⛔ `IsFacRetro` adalah **DUA rule di kelas berbeda** — catatan asal-usul `IsUW` hanya berlaku untuk satu

> ⛔ **Dipersempit 21 September 2026.** Rumusan sebelumnya menuliskan asal-usul salinan `IsUW` seolah
> berlaku **umum** untuk `IsFacRetro`. Itu **hanya benar untuk satu dari dua rule**. Rumusan lama
> tidak dihapus — ia **dibatasi lingkupnya** (`PANDUAN-KERJA` §7).

`[terverifikasi]` Identitas rule ditentukan **basis `pzInsKey` + `pyClassName`**, bukan nama berkas
(`_ARSIP-lintas-siklus\_BACA-INI.md` Koreksi R1). `When\IsFacRetro` berwujud **dua rule**:

| | Varian **OfferFacIn** | Varian **Work** |
| --- | --- | --- |
| `pyClassName` | `ASM-FW-GISFW-Data-OfferFacIn` | `ASM-FW-GISFW-Work` |
| Terekspor di | **NB** dan **RNW** | **Endorsment** |
| `<pyLogic>` `A` menguji | `.IsFacRetro = 1` | `pyWorkPage.OfferFacIn.IsFacRetro = 1` |
| `<pyNestedConditions>` | ⛔ `rowdata(1)` **masih memuat sisa** pemeriksaan workbasket (`pxRequestor…pyWorkBasketName = ReasFacInGroupLeader`) | ✅ **BERSIH** — tidak ada sisa |
| `<pzOriginalInstanceKey>` | menunjuk rule **`ISUW`** | — |

⚠️ **Asal-usul salinan `IsUW` HANYA berlaku pada varian `ASM-FW-GISFW-Data-OfferFacIn`.** Di varian
itu, yang mengikat adalah **ekspresi tersimpan**, bukan grid kondisinya: yang dieksekusi
`<pyLogic>` = `A` → `<pyConditionValue1>` = `compareTwoValues(.IsFacRetro,"=",1)`; grid kondisi
warisan **diabaikan**. Sikapnya **sama persis dengan `IsOfferFacIn`** (K-002, tiket NB
`..\09-predikat-sikap-khusus.md`).

⛔ **Varian `ASM-FW-GISFW-Work` TIDAK memuat sisa itu** dan **tidak perlu perlakuan khusus** —
kondisinya bersih. Memperlakukan kedua varian sama akan memasang penanganan warisan pada rule yang
tidak memerlukannya.

`[dugaan]` **Varian `Work` kemungkinan yang dipanggil flow.** Indeks rujukan
`NB FacIn\Flow\InputInwardFacultativeOffer.xml` menyebut kelas **`ASM-FW-GISFW-Work`**. ⚠️ Indeks merekam resolusi
**saat berkas terakhir disimpan** — ia **bukan bukti** rule mana yang resolve saat runtime, sehingga
kata "kemungkinan" **tetap `[dugaan]`**. **Keduanya diport** (pola rule kembar, `..\..\10-audit\10-rule-kembar-pzinskey.md`).

✅ **K-057 TIDAK terpengaruh.** Kedua varian menguji **flag yang sama** (`.IsFacRetro = 1`) pada
**halaman yang sama** (`OfferFacIn`) — hanya jalur penulisan propertinya berbeda (relatif vs lewat
`pyWorkPage`). Pemicu Fac Out tetap seperti ditetapkan K-057.

⛔ **Tiga bentuk penanda spreading TIDAK DIGABUNG** (K-054):

| Bentuk | Ekspresi | Di mana |
| --- | --- | --- |
| **ketat** | `.TreatyType=="10015"` | `SetDataFacOutFire_Act` · `CopyFacRetroFire_ACT` |
| **contains tunggal** | `@contains(.TreatyType,"10015")` | di dalam `CopyToAllSpreading_ACT` · `CopyToAllLocSpreading_ACT` |
| **umum tiga-arah** | `@contains(.TreatyType,"10015")\|\|@contains(.TreatyName,"SPL")\|\|@contains(.TreatyType,"10007")` | `CopyToAllSpreading_ACT` · `CopyToAllLocSpreading_ACT` |

⚠️ **`"SPL"` diuji pada `TreatyName`, bukan `TreatyType`.** Menyamakan ketiga bidang uji akan menarik
baris spreading yang tidak seharusnya ikut.

### ⛔ Empat properti bernama mirip — registry harus membedakannya

`[terverifikasi]` Registry **tidak boleh** memperlakukan keempatnya sebagai satu:

| Properti | Gaya nilai | Perannya |
| --- | --- | --- |
| `OfferFacIn.IsFacRetro` | `0`/`1` telanjang | **pemicu Fac Out** — dibaca `When\IsFacRetro`, menjadi `<pyTaskWhen>` di flow |
| `OfferFacIn.IsInputFacRetro` | `0`/`1` | gerbang **re-input** saat UW — dibaca `When\IsInputFacRetro` |
| `pyWorkPage.IsInFacRetro` | dibanding `!= 1` | **kondisi tampilan** Section; **bukan** rule `When` |
| `.IsFacRetroOffer` | **`"0"` berkutip** | penanda penawaran retro; disetel `OfferFacOut_PreAct` |

⚠️ Hanya **dua** di antaranya yang berupa rule `When` dan masuk registry. `pyWorkPage.IsInFacRetro`
dan `.IsFacRetroOffer` adalah **properti biasa** — memasukkannya ke registry adalah kekeliruan
kategori. Rincian di `..\..\10-audit\03-struktur-facretro-dan-populasi.md` dan tiket
`F05-empat-struktur-staging.md`.

**Asal (Pega).** `When\IsFacRetro` · `When\IsInputFacRetro` · penanda di
`Activity\SetDataFacOutFire_Act` · `CopyFacRetroFire_ACT` · `CopyToAllSpreading_ACT` ·
`CopyToAllLocSpreading_ACT`

**Keputusan.** **K-054** · K-050 (satu registry, sumber data per-rule) · K-002 · K-019 ·
**K-057** (pemicu = flag `.IsFacRetro`; rumusan `IsFacout` **dicabut**) · `CLAUDE.md` §4.5, §4.6

**Blocked by:** `..\08-registry-rules-eval.md` · `..\09-predikat-sikap-khusus.md`
⛔ Keduanya **belum dikerjakan** — lihat §0 indeks.

**Status:** blocked

- [ ] `IsFacRetro` dan `IsInputFacRetro` hidup lewat **Seam 1 `rules.Eval`** — tanpa registry kedua
- [ ] `IsFacRetro` memakai **ekspresi tersimpan** `.IsFacRetro = 1`; grid kondisi warisan **diabaikan**
- [ ] ⛔ **DUA rule `IsFacRetro` diport** — kelas `Data-OfferFacIn` **dan** `GISFW-Work`; komentar §4.6 menyebut **kelasnya**, nama saja tidak cukup
- [ ] Komentar menyebut asal-usul salinan `IsUW` **hanya pada varian `Data-OfferFacIn`** — varian `GISFW-Work` kondisinya bersih dan **tidak** diberi penanganan warisan
- [ ] **Tiga bentuk penanda spreading terpisah** dan tidak dapat tertukar
- [ ] `"SPL"` diuji pada **`TreatyName`**; kedua kode lain pada `TreatyType`
- [ ] ⛔ `IsFacout` **tidak diimplementasikan di sini** — komentar menunjuk `..\09-predikat-sikap-khusus.md`
- [ ] `CoverageBasis == 5` dikenali sebagai **Layering Basis** `[terverifikasi]` `DDL\CoverageBasis.xml`
- [ ] Arti `10015` = FACOUT dan `10007` = ORS ditandai **keterangan work owner (K-054)**, bukan terverifikasi korpus
- [ ] `[pertanyaan terbuka]` arti `TreatyName == "SPL"` **tetap terbuka**; nilainya diport apa adanya
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)
