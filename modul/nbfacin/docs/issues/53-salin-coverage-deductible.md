# 53: Tombol salin coverage / deductible di grid coverage FIRE

> ⚠️ **Disusun agent sesi `nusantarare-0f` — bukan hasil `/to-tickets`.** Pemicu: work owner 06-10-2026 (gambar layar
> Pega) "tombol ini kemana, knpa km tidak buatin" — tombol Copy Coverage / Copy Deductible ditunda di tiket 43 (baris 47:
> "Copy Coverage / Deductible … = tahap C2–C4") dan tiket 45, lalu tidak diangkat lagi. Survei korpus oleh subagen
> (pengurai XML); butir kunci diperiksa ulang sesi 0f.

**Status:** selesai di frontend 06-10-2026 (uji hijau). Tanpa backend baru — premi dihitung ulang lewat
`POST …/hitung-coverage` yang sudah ada; tersimpan lewat Save tab Coverage.

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\NB FacIn\`

`Section\InputCoverageFire.xml` (`ASM-FW-GISFW-DATA-PROPERTYITEM!INPUTCOVERAGEFIRE`), kontainer tombol tampil bila
`!IsClaim && .FlagDelete != 1`:

| Tombol | Tampil bila | Activity | Perilaku |
|---|---|---|---|
| Copy Coverage | `.pxListSubscript==1` | `CopyCoverageObjItem_Act` | item 1 lokasi → semua item lokasi itu; `.CoverageList` diganti (Property-Remove + Page-Copy); TotalNetRate / FlagNetRate disalin; `.TSI` = TSIObjectItem tujuan; CountPremi_ACT percent / percent |
| Copy Coverage To All | `.IndexProperty==1 && .pxListSubscript==1` | `CopyCoverageToAllLocation_Act` | lokasi 1 / item 1 → semua item semua lokasi; diganti; CountPremi_ACT; **Obj-Save** (@4070) |
| Copy Coverage From | ALWAYS | local action `CopyCoverageFrom` → `Section\CopyCoverageFrom.xml` (Location Address, Object Item Name, "Copy to Another Object Item" → `CopyCoverage_Act`) | item terpilih → item yang dibuka; diganti; CountPremi_ACT; popup ditutup |
| Copy Deductible To All | `.IndexProperty==1 && .pxListSubscript==1` | `CopyDeductibleToAllLocation_Act` Method "ALL" | per coverage tujuan, `.DeductibleList` diganti dengan deductible coverage sumber ber-`.Coverage` sama (`Param.idcoverage==.Coverage`); tanpa padanan → kosong; item 10000992 (Business Interruption) dilewati |
| Copy Deductible For BI | `.ItemTypeID ='10000992' && .IndexProperty==1` | sama, Method "BI" | hanya item 10000992 |

- Tombol "Copy to Another Object Item" langsung di InputCoverageFire tidak pernah tampil (`…&&1==2`); yang hidup hanya di
  popup Copy Coverage From.

## Frontend (sesi 0f)

- `components/salinCoverage.ts` (`salinCoverageLokasi`, `salinCoverageSemua`, `salinCoverageDari`, `salinDeductible`,
  `JENIS_ITEM_BI`) + uji `salinCoverage.test.ts`; `components/PopupSalinCoverage.tsx`; `TabCoverage.tsx` (bar tombol di
  bawah grid coverage item, `pasangSalinan` = hitung ulang premi item tujuan); `labels.ts` `SALIN_COVERAGE`; `nbfacin.css`.

## Keputusan agent (K53-1 dan K53-2 disetujui work owner 06-10-2026: "setuju dengan rekomendasi km untuk K53-1 , K53-2")

- **K53-1** TSI coverage tujuan = TSI Object Item tujuan (maksud `CopyCoverageObjItem_Act` 6.6.1). Pega `CountPremi_ACT`
  langkah 10 dapat memakai TSI item SUMBER karena indeks ikut tersalin `[dugaan]` — tidak ditiru.
- **K53-2** Copy Coverage To All TIDAK langsung menyimpan (Pega `Obj-Save`); seperti tombol lain, tersimpan lewat Save tab.
- **K53-3** Copy Deductible For BI dibuat seperti Pega (tampil hanya pada item 10000992 di lokasi pertama).
