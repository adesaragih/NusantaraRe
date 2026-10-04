// Pilihan `Reinsurance Type` tahun treaty — kolom `TREATYYEAR.PROPORTION`.
//
// Korpus: `Section/InputDtlTreatyContact.xml` b6800 → `InputTreatyYear.Proportion`
// (`pxDropdown`, `pyListSource associated`); daftar nilai properti `.Proportion`
// TIDAK diekspor. Nilai dan labelnya [keputusan work owner 30-09-2026], selaras
// keterangan parameter "Proportional / NonProportional"
// (`Activity/SetTreatyArrangementDesc_Act.xml`) dan syarat tampil
// `.Proportion != 'NonProportional'` (`Section/ViewDetailDescription.xml`
// b3614/b6719). OQ-TCO-04 ditutup. Sebelumnya layar memakai master jenis
// reasuransi (tiket 02) — `[dugaan]` yang kini dikoreksi.

/** Nilai tersimpan → label tampil. */
export const PROPORSI_TCO = [
  { value: 'Proportional', label: 'Proportional' },
  { value: 'NonProportional', label: 'Non Proportional' },
] as const

/** Label tampil sebuah nilai `PROPORTION`; nilai lain (data lama) apa adanya. */
export function labelProporsi(nilai: string): string {
  return PROPORSI_TCO.find((p) => p.value === nilai)?.label ?? nilai
}
