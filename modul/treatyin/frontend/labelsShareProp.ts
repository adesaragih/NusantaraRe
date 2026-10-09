// Label rincian tab SHARE cabang Proporsional — dari ekspor:
//   `Section/TotalLimits.xml`  (FlowAction `TotalLimits`, rincian grid Kind of Treaty)
//   `Section/DetailShare.xml`  (FlowAction `DetailShare`, rincian grid Treaty Group)
// dan gambar Pega `17-prop-menu-share.png`.

/** Grid `.Detail` Section TotalLimits. */
export const KOLOM_TOTAL_LIMITS = ['Treaty Group', '% RNM Share'] as const

export const DETAIL_SHARE = {
  /** `.RNMShare` — `pyLabelFieldValue` `% RNM Share`. */
  persenRnmShare: '% RNM Share',
  /** Grid `.RNMShareList` — kepala kolom pertama `RNM Share`. */
  rnmShare: 'RNM Share',
  /** Blok `Spreading` (`pyIncludeHeader` true) — gambar 17. */
  spreading: 'Spreading',
  /**
   * `.SpreadingTypeID` — dropdown RD `BrowseTreatyArrangement_ParentReinsMasterTrt`.
   * ⭐ Sejak 8 Oktober 2026 SELALU tampil (keputusan pemilik proses): spreading
   * dipilih dari master, tidak diketik.
   */
  jenisSpreading: 'Spreading Type',
  /** Grid `.SpreadingList` — baca-saja, anak susunan dari `PROPORTIONALARRG`. */
  kolomSpreading: ['Reins Type', 'Pct'],
  totalPct: 'Total Pct',
  /** `.SpreadingTotalPct` — cabang Spreading Type. */
  totalSpreadingPct: 'Total Spreading Pct :',
  /** Grid `.RNMSpreadedList` / `.RNMSpreadedListRI`. */
  sebaranOR: 'Value Spreading OR',
  sebaranRI: 'Value Spreading R/I',
  /** `TreatyIn.RnmShareDeducted` — bila `TreatyIn.FacultativeShare >0`. */
  rnmShareDeducted: 'Share to RNM :',
  value: 'Value',
  rincian: 'Rincian baris',
  kindBelumDipilih: '(Kind of Treaty belum dipilih)',
} as const
