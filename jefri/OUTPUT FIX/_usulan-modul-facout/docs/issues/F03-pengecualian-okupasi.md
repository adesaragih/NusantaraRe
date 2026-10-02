# F03: Daftar pengecualian okupasi — 26 alternatif

**What to build:** Sebagian kode okupasi menandai risiko sebagai **pengecualian**, dan penandaan itu
ikut berjalan di jalur derivasi Fac Out lini Fire.

`[terverifikasi]` Gerbangnya **satu ekspresi dengan 26 alternatif `||`**: **20** memakai
`@startsWith` dan **6** memakai kesetaraan ketat. Keluarannya menyetel
`pyWorkPage.OfferFacIn.IsOccupException`.

```powershell
# menghasilkan 26 / 20 / 6
$c=[IO.File]::ReadAllText("D:\migrasi\RNM\NB FacIn\Activity\SetDataFacOutFire_Act.xml")
$occ=[System.Net.WebUtility]::HtmlDecode(
  [regex]::Match($c,'<pyStepsPreCondParamsWhen>([^<]*OccupationId[^<]*)</pyStepsPreCondParamsWhen>').Groups[1].Value)
($occ -split '\|\|').Count
([regex]::Matches($occ,'@startsWith')).Count
([regex]::Matches($occ,'\.OccupationId==')).Count
```

⛔ **Gaya penulisannya tidak konsisten, dan itu DIPORT APA ADANYA.** `[terverifikasi]` Sebagian kode
diuji sebagai **awalan** (`251`, `252`, `257`, `258`, `259`), sebagian lagi sebagai **kesetaraan
penuh** (`256`) — padahal berada di rentang yang sama. Menyeragamkannya **mengubah himpunan kode yang
tertangkap**, dan itu perubahan perilaku (`CLAUDE.md` §1).

⚠️ **Awalan yang lebih pendek menelan yang lebih panjang.** `[terverifikasi]` Daftar memuat awalan
`26` dan `28` berdampingan dengan kesetaraan penuh pada kode empat-lima digit. Urutan dan bentuk
ujinya **dipertahankan persis**, tidak diringkas menjadi rentang.

**Asal (Pega).** `Activity\SetDataFacOutFire_Act` — `pyStepsPreCondParamsWhen` atas `.OccupationId`;
target `pyWorkPage.OfferFacIn.IsOccupException`

**Keputusan.** K-046 (pola kode usang diport apa adanya) · K-054 · `CLAUDE.md` §1, §4.6

**Blocked by:** F02

**Status:** blocked

- [ ] Daftar **26 alternatif** lengkap, tidak satu pun dihilangkan
- [ ] **20 uji awalan** dan **6 uji kesetaraan penuh** tetap sebagai bentuk aslinya
- [ ] ⛔ Daftar **tidak diringkas** menjadi rentang, regex, atau tabel lookup
- [ ] Kode okupasi disimpan sebagai **string**, bukan angka — uji awalan menuntutnya
- [ ] Hasilnya menyetel penanda pengecualian okupasi, bukan menggerakkan alur
- [ ] Test menguji **batas**: kode yang persis cocok, yang cocok sebagai awalan, dan yang di luar daftar
- [ ] **K-046** `K046_OkupasiFacOut_GayaUjiTidakSeragam` — awalan dan kesetaraan bercampur, sengaja dipertahankan
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)
