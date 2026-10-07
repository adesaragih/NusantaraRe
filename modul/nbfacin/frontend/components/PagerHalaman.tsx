// Navigasi halaman daftar popup NB FacIn (hasil Choose Accumulation, Choose Zip Code) - pola pager master inti:
// "Menampilkan a–b dari n data" + Pertama / Sebelumnya / nomor halaman (… bila banyak) / Berikutnya / Terakhir.

import { jumlahHalaman, nomorHalaman } from '../../../../inti/frontend/master/DaftarMaster'
import { TEKS_AKUMULASI as T } from '../labels'

export default function PagerHalaman({
  total,
  halaman,
  ukuran,
  onPindah,
}: {
  total: number
  halaman: number
  ukuran: number
  onPindah: (h: number) => void
}) {
  const dari = jumlahHalaman(total, ukuran)
  return (
    <div className="pager">
      <span className="muted">{T.menampilkan(halaman, ukuran, total)}</span>
      <nav className="pager__aksi" aria-label={T.navigasiHalaman}>
        <button type="button" className="btn btn--ghost btn--sm" disabled={halaman <= 1} onClick={() => onPindah(1)}>
          {T.pertama}
        </button>
        <button type="button" className="btn btn--ghost btn--sm" disabled={halaman <= 1} onClick={() => onPindah(halaman - 1)}>
          {T.sebelumnya}
        </button>
        {nomorHalaman(halaman, dari).map((n, i) =>
          n === null ? (
            <span key={`jeda-${i}`} className="muted">
              …
            </span>
          ) : (
            <button
              key={n}
              type="button"
              className={'btn btn--sm ' + (n === halaman ? 'btn--primary' : 'btn--ghost')}
              aria-current={n === halaman ? 'page' : undefined}
              onClick={() => onPindah(n)}
            >
              {n}
            </button>
          ),
        )}
        <button type="button" className="btn btn--ghost btn--sm" disabled={halaman >= dari} onClick={() => onPindah(halaman + 1)}>
          {T.berikutnya}
        </button>
        <button type="button" className="btn btn--ghost btn--sm" disabled={halaman >= dari} onClick={() => onPindah(dari)}>
          {T.terakhir}
        </button>
      </nav>
    </div>
  )
}
