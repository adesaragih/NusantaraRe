# 07: Uang — presisi penuh, persentase bukan uang, pajak brokerage apa adanya

**Status:** sebagian *(implementasi 2026-10-03, cabang `modul/nbtreatyin/implementasi`; semula: ready-for-agent)*
**Blocked by:** 01
**Menutup:** AC 18 · 23 · 24 · 25 · 26 · 27 · 28 · 85 · 86 *(9 AC)* — US 21 · 22 · 36 · 38 · 39

## Hasil & nilai pengguna

Hari ini nilai uang ditampilkan di layar apa adanya seperti tersimpan, ⭐ **bukan dihitung ulang**.
⚠️ `[terverifikasi]` Empat medan yang tampak seperti uang sebenarnya **persentase** — `12.5` berarti
12,5 persen. Dan potongan brokerage dibagi **1,022** bila jenis pajaknya *inclusive*, ⛔ dengan
pembanding teks yang **persis**: huruf kecil atau spasi di belakang menggeser angkanya **2,2 %**
tanpa pesan galat.

Sesudah tiket ini, nilai uang tersimpan **berpresisi penuh**, ditampilkan beserta mata uangnya, dan
⭐ rumus pajak brokerage berjalan **persis** seperti sistem lama — termasuk ketidakseragamannya.

## Area codebase

- Lapisan repository: pembacaan nilai uang dan mata uang pasangannya
- Fungsi murni: rumus potongan pajak brokerage
- Lapisan handler: penyajian

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Rumus pajak brokerage | `@if(TypeTax == "Inclusive", Deduction / (102.2/100), Deduction)` — **12 tempat pada 7 Activity** |
| Nilai bawaan jenis pajak | `"Inclusive"`; ⛔ nilai lawannya **tidak pernah tertulis** di mana pun |
| Medan uang | **kolom tersimpan** pada sumber relasional; ⛔ nol `Activity` atau `DataTransform` menghasilkannya |

## ADR terkait

- ⭐ **ADR-0003** — uang **tidak** direpresentasikan sebagai `float`

## Acceptance criteria

- [ ] 🟡 **AC 23** — nilai uang disimpan **berpresisi penuh**
- [x] **AC 24** — pembulatan **hanya** di titik penyajian, ⛔ tidak pernah di repository
- [x] **AC 25** — uang **tidak pernah** diwakili tipe pecahan biner
- [ ] 🟡 **AC 26** — empat medan itu dibaca sebagai **persentase**, bukan jumlah uang
- [x] **AC 27** — jenis pajak *inclusive* ⇒ potongan dibagi **1,022**
- [x] **AC 28** — nilai lain — ⭐ **termasuk kosong, huruf kecil, atau berspasi** — ⇒ potongan
      dipakai apa adanya
- [x] **AC 18** — kedelapan medan uang **ditampilkan apa adanya**, ⛔ tidak dihitung ulang
- [x] **AC 85** — angka uang ditampilkan beserta **kode mata uang pasangannya**
- [x] **AC 86** — ⛔ format penyajian **belum ditetapkan**; penyajian tidak dibangun dengan format
      yang ditebak

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **6** | tipe penyimpanan kolom uang — sisi Pega menyimpannya sebagai **teks** | ⚠️ menahan **penguraian**, bukan rumusnya |
| **7** | format penyajian — desimal dan pemisah ribuan | ⚠️ menahan **penyajian** saja |
| **8** | satu medan persentase **dari apa** | tidak menahan |
| **9** | ketidakseragaman presisi 4 lawan 8 desimal | tidak menahan |

## Perintah verifikasi

1. Jenis pajak `"Inclusive"` — ⭐ potongan **dibagi 1,022**.
2. Jenis pajak `"inclusive"` huruf kecil — ⭐ potongan **tidak** dibagi; selisihnya **2,2 %**.
3. Jenis pajak **kosong** — ⭐ sama seperti butir 2.
4. Simpan nilai berdesimal panjang, baca kembali — ⭐ **tidak berubah**.

## Catatan

⭐ **Penyajian uang TIDAK bergantung pada P18.** `[keputusan work owner]` P43 menetapkan angka
**disalin apa adanya**, sehingga bagian terbesar uang dapat dibangun sekarang.

> ⚠️ **RALAT.** `[penyimpangan sadar]` 2026-09-22 — kalimat penutup semula berbunyi:
> > *"⛔ Yang menunggu P18 adalah **perhitungannya**, ada di tiket **13**."*
>
> ⭐ **P18 ditarik** — perhitungannya pun tidak menunggu P18. Rumus pajak brokerage terbaca penuh
> di ekspor: `@if(TypeTax="Inclusive", @divide(Deduction,@divide(102.2,100,8),8), Deduction)`,
> **PPH 2 %**, **PPN 2,2 %**. Perhitungannya ada di tiket **13**, yang kini tertahan oleh **P30**
> dan **P8** saja.

## ⛔ RALAT dan pertentangan — 2026-10-03

1. **`DEDUCTION1/2` halaman polis.** `[keputusan work owner]` P29 / AC 26 / spec-penyimpanan AC 38:
   persentase. XML memperlakukan nilai halaman polis sebagai **jumlah** (dikurangkan dari premi —
   `CountNetPremi_act` langkah 4; dibagi 1,022 — `SetPPNPPH` langkah 4; kontrol layar `pxCurrency`).
   Aturan prompt: ikuti WO, catat ⇒ kolom bergolongan **persen**, rumus diport apa adanya (AC 79).
2. **`SetPPNPPH` langkah 4 BERSYARAT** — `.FlagPPH=="true"` (lewati syarat berikut) ATAU
   `ListAgent.pxResults(1).STS_PKP == 1` (RD `BrowseClientName_RD` atas SourceOfBusiness). Port pertama
   menganggapnya tanpa syarat — diperbaiki.
3. Nilai master `TreatyIn.*` tidak tersedia (tiket 01 RALAT 1) — langkah yang membaginya dilewati.

## ⛔ RALAT putaran 2 (P4, 03-10-2026) — rumus dan pajak

1. **AC 26 — `RNM_SHARE` dan `BROKERAGE`: nol pemakai.** Bunyi lama: *"`DEDUCTION1` `DEDUCTION2`
   `BROKERAGE` `RNM_SHARE` dibaca sebagai **persentase**"*. Bukti (`docs/alat/pemakai.py` di atas
   `graf.Graf(...).terjangkau()`, 278 rule / 176 terjangkau):
   - `RNM_SHARE` — **0 rule** di seluruh korpus. RD `BrowseTreatyJoinEDM` / `BrowseTreatyInDetail`
     (33 kolom) tidak memuatnya. Yang dipakai rumus adalah `pyWorkPage.TreatyIn.RNMShareP` /
     `RNMShare` — medan halaman master JSON (`FetchMasterTreatyIn`), bukan kolom view.
   - `BROKERAGE` — 4 rule terjangkau memuat kata itu, **tak satu pun membaca kolom/properti
     `BROKERAGE`**: `Section\DetailPolicyTreatyInNonProportional`, `…NonProportionalEDM`,
     `…OutNonProportional` hanya `pyCaption Brokerage (IDR)` / `Brokerage (USD)` / `Total Brokerage`
     (label kolom grid `.Deductible` / `.Deductible2` jalur NonProp), dan
     `Activity\InputPolicyTreatyOutDetail_NonProp` hanya `pyStepsDescription` "Total BROKERAGE".
   ⇒ Kedua kolom view **tidak dibaca dan tidak dibangun** (nol pemakai — bab 4 prompt putaran 2).
   `[terbuka]` §9.2 butir 8 ("`RNM_SHARE` persentase dari apa") tetap milik pemiliknya, tetapi tidak
   menahan apa pun di modul ini.
2. **`DEDUCTION1/2`: K3 (03-10-2026) — rumus XML apa adanya.** Rumus sudah diport apa adanya
   (`CountNetPremi_act` langkah 4 mengurangkannya dari premi; `SetPPNPPH` langkah 4 membaginya
   1,022) — tidak berubah. ⚠️ Bagian "label dan penyajian mengikuti pemakaian XML" (kolom katalog
   `kPersen` DEDUCTION1/2, label layar) milik paket katalog/layar, bukan paket ini — butir terbuka.
3. **`SetPPNPPH` diuji dengan nilai XML** (`models/setppnpph_test.go` TestSetPPNPPHNilaiXML; tujuh
   kasus): PremiOgp 1000, PremiOnp 600, Deduction1 51,1 ⇒ BrokerageFee `@divide(2.5,100,8)` × 1600
   = 40; Inclusive 51,1 / 1,022 = 50 ⇒ PPH 1, PPN 1,1; selain Inclusive 51,1 ⇒ PPH 1,022, PPN 1,1242;
   cabang langkah 4: `FlagPPH=="true"` (melewati syarat PKP), `STS_PKP == 1`, keduanya salah ⇒
   BrokerageFee* tidak disentuh dan PPH/PPN tetap 0 (langkah 2); `"TRUE"` bukan `"true"`. Lewat HTTP:
   `handlers/logika_test.go` TestSetPPNPPHMembacaStatusPKPAgen (STS_PKP dari RD
   `BrowseClientName_RD`).
