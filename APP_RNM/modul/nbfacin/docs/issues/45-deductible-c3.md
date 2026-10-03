# 45: Tab Coverage FIRE — tahap C3 Deductible

> ⚠️ **Disusun agent atas perintah sesi `nusantarare-0f` — bukan hasil `/to-tickets`.** Frontend dibangun sesi 0f;
> backend (simpan / baca) diserahkan ke sesi c3.

**What to build:** setiap coverage di tab Coverage FIRE memuat grid **Deductible** (tambah / hapus / form lipat), disimpan
bersama Save tab Coverage.

**Status:** frontend selesai 03-10-2026 (uji hijau). ⛔ Backend belum: `PUT …/objek` memakai `DisallowUnknownFields`,
jadi menyimpan coverage yang memuat `deductibles` akan ditolak 400 sampai backend menerima medan itu.

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\NB FacIn\` dan `D:\migrasi\RNM\DDL\`

- `Section\CoverageItem.xml`: grid PageList `.DeductibleList` berjudul "Deductible" (tampil `!isClaim`); kepala kolom kosong
  kecuali "Time Excess (Days)"; Add = addRow + `addDeductible_ACT`, Delete = deleteRow; form = flow action
  `InputDtlDeductibleFire_FacIn` → section `addDeductible.xml`.
- `Section\addDeductible.xml` sel 5–14: Type Deductible (selalu); Pct Deductible (`.TypeDeductible != 7`); MinMax
  (`!= 7`, baca-saja bila Type 0 / 7); Currency (baca-saja bila Type 0); Type Deductible / Pct Deductible kedua
  (`.MinMax = 3`); Condition; Amount (selalu); Condition teks (`.Condition = 5`); Time Excess (In Days). Tidak ada medan
  wajib; teks pilihan kosong "Choose".
- Daftar pilihan = aturan properti `ASM-FW-GISFW-DATA-DEDUCTIBLE!TYPEDEDUCTIBLE`, `!TYPEDEDUCTIBLE2`, `!MINMAX`,
  `!CONDITION` (`DDL\`, dicari menurut `<pxInsName>` — lihat catatan Condition di bawah).
- Sampel kasus `DDL\P-5 *.txt` (kelas `ASM-FW-GISFW-Data-Deductible`) memuat properti `TypeDeductible`, `PctDeductible`,
  `MinMax`, `Currency`, `TypeDeductible2`, `PctDeductible2`, `Condition`, `InputCondition`, `Amount`, `TimeExcess`,
  `FlagCurrency`, `IndexProperty` / `IndexPropertyItem` / `IndexCoverage` / `IndexDeductible`. Properti kosong tidak
  diekspor (mis. satu deductible MinMax 3 tanpa `Currency` / `Amount`).
- ⚠️ `DDL\Condition.xml` kini berisi `DATA-DEDUCTIBLE!CONDITION` dan menimpa `DATA-PROPERTYITEM!CONDITION` (Condition
  Object Item, tiket 39) yang dikirim dengan nama berkas sama. Uji Object Item Condition dilewati sampai berkas itu
  dikirim ulang dengan nama lain.

## Frontend (sesi 0f)

- `components/GridDeductible.tsx` (`deductibleBaru`, `adaGalatDeductible`), dipasang di bawah `FormCoverage`; mata uang
  dari `daftarMataUang`. `TabCoverage` menolak Save bila ada angka deductible tak sah.
- `api.ts`: `CoverageObjek.deductibles?: Deductible[]`; `Deductible` = sepuluh medan teks (kode / desimal).
- `labels.ts`: `JUDUL_DEDUCTIBLE`, `FORM_DEDUCTIBLE`, `KOLOM_TIME_EXCESS`, `DEDUCTIBLE_KOSONG`, `OPSI_TYPE_DEDUCTIBLE(2)`,
  `OPSI_MINMAX`, `OPSI_KONDISI_DEDUCTIBLE`.
- Uji: `labels.test.ts` (sel 5–14, syarat tampil, grid CoverageItem, empat PromptList menurut pxInsName),
  `components/GridDeductible.test.tsx`.

## Kontrak backend (untuk sesi c3)

- `GET` / `PUT …/objek`: setiap coverage membawa `deductibles: Deductible[]` (urutan = urutan grid; `[]` bila kosong).
  Medan teks: `typeDeductible`, `minMax`, `currency`, `typeDeductible2`, `condition`, `inputCondition`; desimal (teks,
  ADR-0034): `pctDeductible`, `pctDeductible2`, `amount`, `timeExcess`. Kosong = NULL.
- 400 bila desimal tak sah (`…deductibles[d].<medan>`). Tidak ada medan wajib (sesuai addDeductible).
- Penyimpanan: tabel anak coverage (rancangan tabel diserahkan c3, uang NUMBER(38,8)); migrasi berikutnya setelah 193.

## Keputusan agent (menunggu konfirmasi)

- **R-1** Kepala kolom grid memakai label form (kepala Pega kosong kecuali "Time Excess (Days)").
- **R-2** Deductible baru bermata uang item (`addDeductible_ACT` membaca `Local.Curr`) `[dugaan]`.
- **R-3** Ditunda ke tahap berikut: penanda beda mata uang (`Activity\SaveDeductible.xml` dan `addDeductible_ACT.xml`
  merujuk `IsBedaCurr`; `SaveDeductible` mengisi `FlagCurrency` `[terverifikasi]`; kaitannya dengan flow action
  `FlowAction\ProtectCurrency.xml` `[dugaan]`), Descriptions, View Old Deductible, Copy Deductible To All / For BI.
- **R-4** Deductible tersimpan lewat Save tab Coverage, bukan tombol Save tersendiri per form.

## Acceptance criteria

- [x] Grid Deductible per coverage: Tambah membuka form, Hapus menghapus baris, baris memperlihatkan label pilihan.
- [x] Syarat tampil / baca-saja mengikuti addDeductible; Amount berformat ribuan.
- [x] Angka tak sah memblokir Save.
- [ ] Backend menyimpan dan mengembalikan `deductibles` (sesi c3).
