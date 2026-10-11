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
   * ⛔ RALAT 8 Oktober 2026 (sore): tampil HANYA bila `.SpreadingTypeID != ''`,
   * persis ekspor. Bila kosong yang tampil grid spreading MANUAL (lihat
   * `kolomSpreadingManual`) — koreksi pemilik proses atas keputusan pagi
   * "selalu tampil".
   */
  jenisSpreading: 'Spreading Type',
  /** Grid `.SpreadingList` — baca-saja, anak susunan dari `PROPORTIONALARRG`. */
  kolomSpreading: ['Reins Type', 'Pct'],
  totalPct: 'Total Pct',
  /** `.SpreadingTotalPct` — cabang Spreading Type. */
  totalSpreadingPct: 'Total Spreading Pct :',
  /**
   * ⭐ Spreading MANUAL (`.SpreadingTypeID == ''`) — `Section/DetailShare.xml`
   * L16, dipulihkan 8 Oktober 2026 atas koreksi pemilik proses: *"spread nya
   * bisa ditambah, cek lagi xml nya"*. Grid `.SpreadingList` berkepala
   * `Reins Type` (dropdown RD arrangement) · `Pct Share` (`pxNumber`), `Add`
   * di kepala / `Delete` per baris (`AddDelSpreadingTreatyin`).
   */
  kolomSpreadingManual: ['Reins Type', 'Pct Share'],
  tambahSpread: 'Add',
  hapusSpread: 'Delete',
  pilihReins: 'Choose',
  /** Panel bentang baris (`pyEditAction` `SpreadingTPDtl`) — `.BreakDownSprdList`. */
  kolomPecahan: ['Spread', 'Currency', 'Share (%)', 'Amount'],
  /** `.SpreadingTotalPct` — cabang manual (`pxNumber`, 2 desimal). */
  totalSharePct: 'Total Share Pct',
  /** Grid `.RNMSpreadedList` / `.RNMSpreadedListRI`. */
  sebaranOR: 'Value Spreading OR',
  sebaranRI: 'Value Spreading R/I',
  /** `TreatyIn.RnmShareDeducted` — bila `TreatyIn.FacultativeShare >0`. */
  rnmShareDeducted: 'Share to RNM :',
  value: 'Value',
  rincian: 'Row details',
  kindBelumDipilih: '(Kind of Treaty not selected)',
} as const
