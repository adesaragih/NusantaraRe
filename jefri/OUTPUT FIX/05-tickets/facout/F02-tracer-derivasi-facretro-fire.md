# F02: Tracer — derivasi Fac Retro saat UW Accept (lini Fire)

**What to build:** Saat Underwriter menekan **Accept**, lokasi yang punya spreading Fac Out otomatis
tersalin ke daftar retrosesi, kasus ditandai sebagai punya Fac Out, dan alur membelok ke menu Fac Out
**sebelum** produksi. Tiket ini menembus seluruh lapisan untuk **satu lini bisnis — Fire** — sebagai
tracer.

⛔ **PEMICUNYA BUKAN `IsFacout`, DAN BUKAN `ProposalAcceptStatus = 4`.** Rumusan lama salah pada dua
lapis dan sudah dibatalkan (`09-facout\01` §0.1):

1. `[terverifikasi]` `IsFacout` muncul **0 kali** di `NB FacIn\Flow\InputInwardFacultativeOffer.xml`
   (`[terverifikasi]` berkas itu hanya ada di folder **NB**). Yang ada
   di `SetValidateDateUW_PostAct` hanyalah **variabel lokal** bernama `isFacOut` — tag
   `<pyLocalParameters>` → `<pyParametersParamName>`, bertipe `String`. Nama variabel lokal yang mirip
   nama rule; **bukan** pemanggilan rule.
2. `[terverifikasi]` `ProposalAcceptStatus = 4` berarti **Banding**, bukan Accept (**K-029**,
   `GLOSARIUM.md`).

`[terverifikasi]` **Rantai yang sebenarnya**, urutan dibaca dari `<pyStepPageReference>` — bukan dari
nomor baris:

| # | Berkas | Tag pembawa | Isi |
| ---: | --- | --- | --- |
| 1 | `Activity\SetValidateDateUW_PostAct` | `pyStepPageReference` = `RH_1.pySteps(3)` | `Call SetDataFacOut_Act`, **TANPA precondition** |
| 2 | `Activity\SetDataFacOut_Act` | `pyStepsActivityName` · `PropertiesName` | `Property-Remove` `pyWorkPage.OfferFacIn.FacRetroList` · `Property-Set` `.OfferFacIn.IsFacRetro = 0` (**reset**) · langkah 3–6 panggil per COB |
| 3 | `Activity\SetDataFacOutFire_Act` | `pyStepsPreCondParamsWhen` · `PropertiesName` | `.OfferFacIn.IsFacRetro = 1`, bergerbang `.TreatyType=="10015"` / `Local.isFacOut==1` / `Local.isFacOutLoc==1` |
| 4 | `Flow\InputInwardFacultativeOffer` | `pyTaskStatusOrWhen`=`WHEN` · `pyTaskWhen`=`IsFacRetro` · `pyLikelihood`=100 | transisi **`Transition92`** → shape `pyMOName` = `[Fac Out]` |
| 5 | shape `[Fac Out]` | **`pyImplementation`** = `OfferFacRetro` | ⚠️ **bukan** `pySubFlowName` |

📌 **Penyaringan terjadi di precondition per COB, bukan di titik masuk.** Titik masuk tidak
bergerbang sama sekali — `SetDataFacOut_Act` selalu dipanggil, lalu mereset flag ke `0`, dan hanya
sub-aktivitas per lini yang menyalakannya ke `1`.

`[terverifikasi]` Penanda penaut ke objek Fac In induk **lahir di sini**, bukan saat produksi:
`SetDataFacOutFire_Act` menyetel `FacRetro.LocationList(<LAST>).Property.ObjectNo` dan
`.Property.ObjectNoFacIn`.

**Asal (Pega).** `Activity\SetValidateDateUW_PostAct` `RH_1.pySteps(3)` · `Activity\SetDataFacOut_Act` ·
`Activity\SetDataFacOutFire_Act` · `Flow\InputInwardFacultativeOffer` `Transition92`

**Keputusan.** **K-057** (pemicu = flag `.IsFacRetro`; rumusan lama **dicabut**) · K-053 · K-054 · K-019 · K-029 ·
K-010/K-012 · `CLAUDE.md` §4.6

**Blocked by:** F01 · `..\01-tracer-money-ratio-premi-pa.md` · `..\11-tangga-akseptasi-bentuk-a.md`
⛔ Kedua tiket NB itu **belum dikerjakan** — lihat §0 indeks.

**Status:** blocked

- [ ] Derivasi dipanggil dari jalur UW Accept **tanpa gerbang di titik masuk**
- [ ] Langkah reset menghapus daftar retrosesi **dan** menihilkan flag sebelum pengisian ulang
- [ ] Flag `.IsFacRetro` menyala **per lini bisnis**, hanya bila penanda spreading ditemukan
- [ ] Flag itu menjadi **gerbang belok** ke menu Fac Out — perilaku `Transition92`
- [ ] ⛔ `IsFacout` **tidak dipakai** di jalur ini; komentar menyatakan sebabnya
- [ ] ⛔ `ProposalAcceptStatus` **tidak menggerbangi** jalur ini
- [ ] Penaut `ObjectNo` + `ObjectNoFacIn` tercipta **saat derivasi**, bukan saat produksi
- [ ] Nilai uang bertipe `Money`; rasio bertipe `Ratio`
- [ ] Demo: satu kasus Fire dengan spreading ber-penanda → daftar retrosesi terisi, flag menyala, alur membelok
- [ ] Demo tandingan: kasus Fire **tanpa** spreading ber-penanda → flag tetap `0`, alur **tidak** membelok
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)
