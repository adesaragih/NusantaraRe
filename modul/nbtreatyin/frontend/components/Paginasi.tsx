// Tombol nomor halaman grid ber-`pyGridProps/pyPageMode=Numeric` - SATU komponen untuk
// grid `.SuggestList` (`Section/ListSuggest`, pyPageSizeOther 5) dan grid popup
// `Section/BusinessAndSOBList` (pyPageSize 50). Halaman aktif selalu dijepit ke rentang
// yang ada (`halamanTerjepit`), sama dengan baris yang ditampilkan (`irisan`).

import { halamanTerjepit, jumlahHalaman } from '../paginasi'

export default function Paginasi({
  jumlahBaris,
  ukuran,
  hal,
  onHal,
}: {
  jumlahBaris: number
  ukuran: number
  /** Halaman yang diminta (berbasis 1); dijepit sebelum ditandai aktif. */
  hal: number
  onHal: (hal: number) => void
}) {
  const n = jumlahHalaman(jumlahBaris, ukuran)
  if (n <= 1) return null
  const kini = halamanTerjepit(hal, n)
  return (
    <div className="nbti__aksi">
      {Array.from({ length: n }, (_, i) => i + 1).map((x) => (
        <button
          key={x}
          type="button"
          className={x === kini ? 'btn btn--sm btn--primary' : 'btn btn--sm'}
          aria-current={x === kini ? 'page' : undefined}
          onClick={() => onHal(x)}
        >
          {x}
        </button>
      ))}
    </div>
  )
}
