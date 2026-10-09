// Teks pilihan (PROMPT VALUE) `TreatyIn.EDMState` dan `TreatyIn.EDMMaterialType`
// — *(8 Okt, E)*.
//
// ⭐ RALAT atas `PENYESUAIAN.kodeBelumBerteks` ("teks pilihan TIDAK ada di
// ekspor"): kedua rule Property-nya DIEKSPOR terpisah, dan modul Treaty In
// sudah memakainya (`backend/services/prompt_value.go`, peta `promptValue`):
//
//   `D:\XML_NURE\_migration-docs\treaty-in-adjustment\ekspor-tambahan\EDMState.xml`
//     kelas `ASM-FW-GISFW-Int-TREATY_IN`, `pyPromptTableList`:
//       1 → `Internal` (@6122), 2 → `External` (@6296)
//   `…\ekspor-tambahan\EDMMaterialType.xml`
//     kelas `ASM-FW-GISFW-Int-TREATY_IN`, `pyPromptTableList`:
//       1 → `Material` (@6154), 2 → `Non Material` (@6328)
//
// Halaman `TreatyIn` berkelas `ASM-FW-GISFW-Int-TREATY_IN` (`pyPagesAndClasses`
// DataTransform korpus), jadi dropdown kepala `TreatyIn.EDMState` @76962 /
// `TreatyIn.EDMMaterialType` @87067 (`Section/InputTreatyInAdjustment.xml`,
// `pxDropdown` `associated`) dan radio picker `TreatyIn.EDMMaterialType`
// (`Section/PickerTreatyInMasterRevisi.xml` @110848, `pxRadioButtons`)
// menampilkan teks ini — bukan kodenya.
//
// ⭐ RALAT (9 Okt) — kolom grid DAFTAR (`.EDMState` @726743 /
// `.EDMMaterialType` @733188, baris `BrowseTREATY_IN_EDM`) JUGA menampilkan
// teks: selnya `pxDropdown` bersumber halaman `EDMStates.pxResults` /
// `EDMMaterial.pxResults` (nilai `.CARI1`, teks `.CARI2`), diisi DataTransform
// `InitTreatyEDMStateName` (pre-DT sel Type). DT itu TIDAK diekspor; teksnya
// diambil dari peta di bawah — kode dan maknanya sama (`EDMState`/
// `EDMMaterialType` TREATY_IN). Pemakai: "isinya bukan angka cek pega".
//
// ⚠️ Nilai di luar peta (mis. `EDMState = 3`) tampil APA ADANYA — tidak ditebak.

/** Prompt value per properti → (standard value → teks). */
export const PROMPT_EDM: Readonly<Record<'EDMState' | 'EDMMaterialType', Readonly<Record<string, string>>>> = {
  EDMState: { '1': 'Internal', '2': 'External' },
  EDMMaterialType: { '1': 'Material', '2': 'Non Material' },
}

/** Teks tampil sebuah kode `EDMState`/`EDMMaterialType`; kode lain apa adanya. */
export function teksPromptEDM(kunci: string, nilai: string): string {
  const peta = (PROMPT_EDM as Readonly<Record<string, Readonly<Record<string, string>>>>)[kunci]
  return peta?.[nilai] ?? nilai
}
