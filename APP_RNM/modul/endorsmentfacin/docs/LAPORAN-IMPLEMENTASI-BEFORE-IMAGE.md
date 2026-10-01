# Laporan implementasi — before-image EDM (E06–E11)

> 01-10-2026. Irisan pertama kode `endorsmentfacin`: Seam 4, tanpa Oracle, tanpa layar, tanpa
> `backend/modul.go` (modul belum terdaftar; `MODUL.md` tetap "belum dimigrasi").
> Sumber kebenaran: korpus `D:\migrasi\RNM\Endorsment Fac In\` (READ-ONLY). Label mengikuti
> `CLAUDE.md` §4.

## 1. Yang dibangun

| Berkas | Isi | Asal Pega |
| --- | --- | --- |
| `backend/models/penawaran.go` · `klon.go` | Agregat `OfferFacIn` sebatas properti yang disentuh before-image; salinan dalam | — |
| `backend/services/beforeimage.go` | `PrepareBeforeImage` — lapis A (54 salinan), memanggil lapis C, porsi periode, lapis B | `SetValueToEDMWork` 14–15 |
| `backend/services/lapisc.go` | Tujuh varian lapis C + penggabungan spreading | `SetOLDValueToEDMWork_<LOB>` |
| `backend/services/nilailama.go` | `IsiNilaiLama` — lapis B, 52 penugasan, tiga gerbang keluar, guard K-046, A.5 | `SetOldData` |
| `backend/services/porsiperiode.go` | `ProrateStartEDM` / `ProrateEDMEnd`; galatnya dikembalikan di `HasilBeforeImage.GalatPorsiPeriode`, tidak menggagalkan lapis A/B/C | `SetValueToEDMWork` 15 |
| `docs/alat/langkah.py` | Alat audit: urai pohon langkah activity; `--uji` menguji instrumennya sendiri | — |

Test: `backend/services/*_test.go` — hitam-kotak lewat dua fungsi seam. Instrumen diuji dengan
**uji mutasi terisolasi**: 22 mutasi pada aturan penting (19 sebelum dan 3 sesudah perbaikan porsi periode hasil code review), 21 tertangkap. Yang lolos (salinan dangkal
`LocationList` di `OfferFacIn.Klon`) ekuivalen dalam alur sekarang: lapis C selalu mengganti daftar
itu dengan salinan dalam sebelum ada yang menulisinya.

## 2. ⛔ Temuan korpus yang menyimpang dari spec dan tiket

Kode mengikuti **korpus**. Spec/tiket belum diralat — menunggu work owner.

### 2.1 Lapis C juga menyalin daftar baris `[terverifikasi]`

Langkah 1.1 tiap varian menyalin daftar lininya dari OldData: `LocationList` (FIRE, Aneka, Golf),
`CargoList` (MC), `VehicleList` (MBU), `PersonList` (PA, LIFE). Ke-54 salinan lapis A **tidak**
menyentuh daftar-daftar itu. MC juga menyalin `PolicyData` utuh, mengosongkan `EndorsementNo`, dan
membawa `TradingList`/`GoodList`/`ConveyanceList` kargo pertama ke akar objek kerja.

### 2.2 Spreading digabung per `TreatyType` untuk SEMUA jenis endorsement `[terverifikasi]`

Enam varian (semua kecuali LIFE) menggabungkan baris spreading tiap coverage yang `TreatyType`-nya
sama, menjumlahkan TSI, premi, dan bagian. Sesudah itu baris berbagian nol dibuang. Langkahnya
berlabel "EDM ADJ SPREADING" dan bergerbang `Type=="4"`, tetapi ber-`pyStepsPreCondition=false`.
Artinya ia **tetap jalan** untuk semua jenis endorsement (P-11, tertutup 19-09-2026).

Aturan buangnya berbeda: **FIRE** membuang bagian `= 0` (`.SharePercentage<>0` → lewati) dan
menyimpan bagian negatif. **Lima lini lain** membuang bagian `<= 0` (`.SharePercentage>0` → lewati).

### 2.3 Penanda dan kedalaman `[terverifikasi]`

| Varian | Tingkat atas | Tingkat di bawahnya (semua `IsOldData`) |
| --- | --- | --- |
| FIRE | Lokasi `IsOldData` | PropertyItem → Coverage → Spreading |
| Aneka | Lokasi `IsOldData` | Occupation → Aneka → Coverage → Spreading |
| Golf | Lokasi `IsOldData` | Aneka → Coverage → Spreading |
| MC | Kargo **`FlagOldData`** | Coverage → Spreading |
| MBU | Kendaraan **`FlagOldData`** | Coverage (+ AdditionalCoverage) → Spreading |
| PA | Orang **`FlagOldData`** | ASMCoverage → Spreading |
| LIFE | Orang **`FlagOldData`** | CoverageList → Spreading |

⛔ Premis kasus uji `K046_Life_DuaPenandaOldData` di E10 ("hanya jiwa dua penanda") **tidak
didukung korpus** — MC, MBU, dan PA berpola sama. ⛔ Bahan spec menyebut kedalaman FIRE sampai
`LayerList`; berkas FIRE memuat **nol** kemunculan `LayerList`.

Cacah penanda di bahan spec (FIRE 12, MC 8, Aneka 14, MBU 9, LIFE 4, PA 8, Golf 12) adalah cacah
**baris yang memuat token `IsOldData`**, bukan cacah penugasan. Penugasan bernilai `"old"`: FIRE 5,
MC 4, Aneka 6, MBU 5, LIFE 3, PA 4, Golf 5.

### 2.4 Lapis B bersumber baris OldData, dipasangkan menurut nomor urut `[terverifikasi]`

Langkah 4 `SetOldData` berjalan di halaman `pyWorkPage.OfferFacIn.OldData`. Ia mengulang daftar
**OldData** dan menulis ke `pyWorkPage.OfferFacIn.<daftar>(.pxListSubscript)`. Guard `.RateOld!=""`
dan `.PremiumOld!=""` membaca **baris OldData**. Akibatnya polis lama yang belum pernah di-endorse
selalu memberi `RateOld = 0` — perilaku yang dinyatakan benar oleh K-046.

`[dugaan]` Lini Travel ditangani lapis B tetapi tidak punya varian lapis C yang menyalin
`PersonList`. Baris kerja yang belum ada karena itu **dibuat** di ujung daftar dan berisi nilai lama
saja (`TestLapisBTravelMembuatBarisTujuan`). Ini dugaan atas perilaku `Property-Set` Pega pada
nomor baris yang belum ada.

## 3. Keputusan rancangan yang perlu dikonfirmasi

1. **Seam 4 berupa dua fungsi**: `PrepareBeforeImage` (saat kasus lahir; lapis A+C+porsi periode,
   lalu lapis B untuk pembukaan pertama) dan `IsiNilaiLama` (lapis B, setiap layar dibuka). Spec
   §Testing menulis satu fungsi; lapis B punya siklus hidup berbeda (spec §1), sehingga satu fungsi
   tidak dapat melayani pembukaan ulang tanpa menimpa perubahan pengguna.
2. **Predikat sebagai masukan** (`services.Predikat`) sampai registry E01/NB-08 ada. Ini pilihan
   pemegang modul pada sesi ini, **bukan** keputusan work owner.
3. **`QuotationData.Type` bertipe `models.JenisPenyesuaian`**, bukan `string`. Penjaga `claimlife`
   (`models/satutype_test.go`) menagih setiap medan `Type string` di seluruh aplikasi.
   Mendaftarkannya berarti mengubah modul lain. Akibatnya medan ini **tidak lagi tertagih** penjaga
   itu — pemilik `claimlife` boleh memilih mendaftarkannya secara eksplisit.
4. **Model data sebagian.** Properti yang tidak disentuh before-image tidak dimodelkan, sehingga
   tidak ikut tersalin. Struktur JSON `JSON_POLIS.DATA_JSON` tetap pertanyaan terbuka.

## 4. Pertanyaan terbuka dan dugaan di kode

| # | Butir | Sikap di kode |
| ---: | --- | --- |
| 1 | Mode pembulatan `@Math.divide(…,20)` | Hasil tak eksak **ditolak** (`ErrModePembulatanBelumTerverifikasi`) — sikap yang sama dengan premi NB (keputusan work owner 30-09-2026). ⚠️ Rasio periode nyata hampir selalu tak eksak (mis. 100/365); porsi periode baru berguna setelah ini dijawab. Galatnya dikembalikan di `GalatPorsiPeriode`, kedua rasio kosong, lapis A/B/C tetap terpakai |
| 2 | Satuan selisih DateTime Pega | `[terverifikasi]` ketiga lokal langkah 15 bertipe **`int`** (`pyLocalParameters` `SetValueToEDMWork.xml`). `[dugaan]` satuannya **hari**: `setProRatePercent_Act.xml` dan `IsThereAnyObjectLocation_Act.xml` membagi selisih DateTime yang sama dengan 365. Dihitung dalam hari bulat; selisih pecahan hari ditolak (`ErrSatuanSelisihWaktuBelumTerverifikasi`), karena cara Pega memotongnya ke `int` belum diketahui. Guard periode nol → 1 hari |
| 3 | Penjumlahan properti desimal kosong (penggabungan spreading) | Ditolak (`ErrSpreadingTakTerjumlahkan`), tidak ditebak nol |
| 4 | `Page-Remove` di dalam loop atas daftar yang sama | `[dugaan]` semua baris yang memenuhi aturan dibuang |
| 4a | Penugasan halaman ke halaman (`newWorkPage.Quotation = …QuotationData`, salinan daftar lapis C) | `[dugaan]` isi tujuan diganti seluruhnya, bukan digabung |
| 5 | Mata uang nol cadangan `@If(…,…,0)` | Ikut mata uang nilai sumbernya (bisa kosong) |
| 6 | Pengurangan tanggal kosong di porsi periode | Ditolak (`ErrTanggalPorsiPeriodeKosong`) |
| 7 | Alasan ambang 100 lokasi · mengapa hanya FIRE menyetel `IsProRate` | Diport apa adanya |

## 5. Perintah audit

```
# uji instrumen dulu - enam butir yang jawabannya sudah diketahui dari bahan spec:
PYTHONIOENCODING=utf-8 py docs/alat/langkah.py --uji "D:/migrasi/RNM/Endorsment Fac In"
#   → 52 *Old · 6 RateOld · 1 PremiumOld · 2 "0" · 54 penugasan 14.3 · 51 dari OldData (6/6 OK, 01-10-2026)

# pohon langkah (urai XML; tanpa tag ber-PII):
PYTHONIOENCODING=utf-8 py docs/alat/langkah.py "D:/migrasi/RNM/Endorsment Fac In/Activity/SetOldData.xml"
PYTHONIOENCODING=utf-8 py docs/alat/langkah.py "D:/migrasi/RNM/Endorsment Fac In/Activity/SetValueToEDMWork.xml"
PYTHONIOENCODING=utf-8 py docs/alat/langkah.py "D:/migrasi/RNM/Endorsment Fac In/Activity/SetOLDValueToEDMWork_<LOB>.xml"

# cara kedua, cacah baris (Git Bash, di "Endorsment Fac In/Activity"):
grep -c '<PropertiesValue>"old"</PropertiesValue>' SetOLDValueToEDMWork_<LOB>.xml   # penugasan
grep -c 'IsOldData' SetOLDValueToEDMWork_<LOB>.xml        # cacah baris token = angka bahan spec
grep -c 'LayerList' SetOLDValueToEDMWork_FIRE.xml         # 0
grep -c 'FlagOldData' SetOLDValueToEDMWork_{MC,MBU,PA,LIFE}.xml
```

`docs/alat/langkah.py` di-commit bersama modul ini supaya angka di atas dapat diulang. Ia mengurai
`pySteps/rowdata` bersarang dan mencetak `pyStepsObjectName`, prakondisi, dan
`pyParamArray/rowdata` `PropertiesName`/`PropertiesValue`. Tag `pyStepsRepeatDef` dan isi
`pyStepsPreCondParams` (selain `…When`) tidak dicetak, karena memuat nama orang dan email.

## 6. Belum dikerjakan

Pemuatan lapis A dari Oracle (E17) · fixture rekonsiliasi eksak (NB-15/16, E22) · satuan ‰/% rate
(resolver NB-03) · pemanggil layar · kasus uji "`IsNotEDM` bukan negasi `IsEDM`" (spec §Testing) —
predikat masih masukan, jadi uji itu milik registry E01/E02 · E01–E05 dan E12–E22.
