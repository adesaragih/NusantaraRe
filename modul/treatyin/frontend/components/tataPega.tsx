// TATA LETAK LAYOUT PEGA untuk tab Treaty In — `pyLayoutOtherFormat` tabel
// SIMPLELAYOUT Section ekspor (8 Oktober 2026, permintaan pemakai: posisi tiap
// tab sama dengan Pega, Prop dan Non-Prop).
//
// Padanan blok `tata` layar Treaty In Adjustment (`KerangkaTab.tsx`, kelas
// `tria__tata--*`), yang dibangkitkan dari Section ekspor YANG SAMA — jadi
// kedua layar tersusun sama. ⛔ Ditiru, bukan diimpor (modul tidak saling
// impor).
//
//   g2 / g3 / g4  `Inline grid double / triple / quadruple` — 2/3/4 kolom sama
//                 lebar, butir mengisi per baris (contoh: Total All Layers RNM
//                 Share = lima `g3` bertumpuk)
//   t3070         `Inline 30 70 table` — pasangan [30% | 70%]
//   alir          `Inline`, `Inline middle`, `Inline labels left`, `Inline grid
//                 30 70` — mengalir sebaris (⚠️ `Inline grid 30 70` dibaca
//                 alir: gambar Pega DetailLimits menampilkan 100% Limit ·
//                 Retention · Cession BERDAMPINGAN)
//   kiri          `Stacked with labels left` — satu butir per baris, label kiri
//   tumpuk        `Default` / `Stacked` — bertumpuk penuh
//
// Skin CSS Pega tidak diekspor: kolom sama lebar, celah tema.

import type { ReactNode } from 'react'

export type TataPega = 'g2' | 'g3' | 'g4' | 't3070' | 'alir' | 'kiri' | 'tumpuk'

/**
 * Satu layout dinamis Pega. `judul` = `pyTitle` bila kepala layoutnya tampil
 * (`pyHeaderType = BAR`, `pyIncludeHeader != false`).
 */
export function TataPegaBlok({ tata, judul, children }: { tata: TataPega; judul?: string; children: ReactNode }) {
  return (
    <div className="trin__tata-wadah">
      {judul !== undefined && judul !== '' && <h5 className="trin__subjudul">{judul}</h5>}
      <div className={`trin__tata trin__tata--${tata}`}>{children}</div>
    </div>
  )
}

/** Sel kosong `Inline grid …` — slot yang Pega biarkan kosong (baris dengan butir lebih sedikit). */
export function SelKosongPega() {
  return <div className="trin__tata-kosong" aria-hidden="true" />
}
