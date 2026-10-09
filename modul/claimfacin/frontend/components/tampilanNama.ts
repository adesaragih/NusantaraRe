// Disalin dari `modul/claimnonprop/frontend/components/tampilanNama.ts` (pola, bukan impor; asal Claim Prop): keputusan
// work owner yang disebut di bawah adalah keputusan layar Claim Prop yang ditiru Claim Fac In (prompt Claim Fac In tahap 1).
// Medan ber-sumber dengan `tampilan` (Consultant / Adjuster, work owner 08-10-2026: "yang dropdown hanya dari namanya
// aja; untuk ID dihapus dari tampilan, tapi tetap simpan ID"): yang dipilih dan tampil adalah nama, yang disimpan tetap
// nilai medan (ID).
//
// Claim Fac In: sel autocomplete tingkat item menampilkan NAMA (Occupation, Object Name) padahal nilai pilihannya kode
// (`OccupationId`, `IndexPropertyItem`) - butir dropdown memakai label pilihan, kodenya sebagai keterangan. Pilihan
// yang nilai = labelnya (Coverage, Object Id) membawa nama coverage / objek dari kolom tambahan sebagai keterangan.

import type { Pilihan, Tata } from '../api'

/** Kolom tambahan pilihan yang menjadi keterangan butir (D_FilteredCoverageList, D_AnekaList, D_Coverage*ClaimList). */
const KETERANGAN_TAMBAHAN = ['CoverageNote', 'ObjectName'] as const

/** Butir dropdown PilihSaring. */
export function butirSaring(
  t: Pick<Tata, 'tampilan' | 'sumber'>,
  opsi: readonly Pilihan[],
): { value: string; label: string; keterangan?: string }[] {
  if (t.tampilan) return opsi.map((o) => ({ value: o.nilai, label: o.label }))
  // Name of Bank (Result.pxResults): nilai `AccountNo|NameOfBank` dikirim ke server, yang tampil nama bank + rincian.
  if (t.sumber === 'rekening') {
    return opsi.map((o) => ({
      value: o.nilai,
      label: o.label,
      keterangan: [o.tambahan?.BranchOfBank, o.tambahan?.NoAccount].filter(Boolean).join(' - ') || undefined,
    }))
  }
  return opsi.map((o) => {
    const keterangan =
      o.label !== '' && o.label !== o.nilai
        ? o.nilai
        : KETERANGAN_TAMBAHAN.map((k) => o.tambahan?.[k]).find((v) => v !== undefined && v !== '')
    return { value: o.nilai, label: o.label || o.nilai, keterangan }
  })
}

/** Teks nilai terpilih: jalur tampilan, lalu label daftar pilihan; tanpa `tampilan` = nilainya sendiri. */
export function teksTerpilih(
  t: Pick<Tata, 'tampilan'>,
  nilai: string,
  ambil: (jalur: string) => string,
  opsi: readonly Pilihan[],
): string {
  if (!t.tampilan) return nilai
  return ambil(t.tampilan) || opsi.find((o) => o.nilai === nilai)?.label || ''
}
