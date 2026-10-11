// Unsur dasar tampilan polis NB Fac In hanya-baca (`PolisNB.tsx`): medan tampil (pxDisplayText nbfacin
// `nbf-inward__teks`), radio / centang nonaktif, dan grid master-detail (`nbf-tabel` + baris buka-tutup). Rupa disalin
// dari `modul/nbfacin/frontend` (pola, bukan impor); kelas `claimfacin__nb-*` + kelas grid modul ini.

import { Fragment, useState, type ReactNode } from 'react'

import { TEKS_NB } from '../polisNB'

/**
 * Medan tampil: label + kotak nilai. `kiri` = "Stacked with labels left" (blok General, form Coverage atas, Cedant);
 * tanpa = label di atas (sub-tab objek, form baris).
 */
export function TeksNB({
  label,
  nilai,
  kiri,
  angka,
  lebar,
}: {
  label: string
  nilai: ReactNode
  kiri?: boolean
  angka?: boolean
  lebar?: boolean
}) {
  return (
    <div
      className={
        'claimfacin__nb-medan' + (kiri ? ' claimfacin__nb-medan--kiri' : '') + (lebar ? ' claimfacin__nb-lebar' : '')
      }
    >
      <span className="claimfacin__nb-label">{label}</span>
      <div className={'claimfacin__nb-teks' + (angka ? ' claimfacin__nb-teks--angka' : '')}>{nilai}</div>
    </div>
  )
}

/** Label tanpa nilai (Following Policy Number - tombol Search-nya tidak dibawa). */
export function LabelNB({ label }: { label: string }) {
  return (
    <div className="claimfacin__nb-medan claimfacin__nb-medan--kiri">
      <span className="claimfacin__nb-label">{label}</span>
    </div>
  )
}

/** Radio Pega nonaktif; nilai di luar pilihan tampil mentah supaya tidak hilang (kode Policy Type Pega tidak diketahui). */
export function RadioNB({
  label,
  pilihan,
  nilai,
  kiri,
  opsiLabel,
}: {
  label: string
  pilihan: readonly string[]
  nilai: string
  kiri?: boolean
  /** Teks pilihan bila berbeda dari nilainya (Share Cedant Type 0 = Gross). */
  opsiLabel?: (p: string) => string
}) {
  return (
    <div className={'claimfacin__nb-medan' + (kiri ? ' claimfacin__nb-medan--kiri' : '')}>
      <span className="claimfacin__nb-label">{label}</span>
      <div className="claimfacin__nb-radio" role="radiogroup" aria-label={label}>
        {pilihan.map((p) => (
          <label key={p} className="claimfacin__nb-pilihan">
            <input type="radio" disabled readOnly checked={nilai === p} /> {opsiLabel ? opsiLabel(p) : p}
          </label>
        ))}
        {nilai !== '' && !pilihan.includes(nilai) && <span className="claimfacin__nb-mentah">{nilai}</span>}
      </div>
    </div>
  )
}

/** Centang Pega nonaktif. */
export function CentangNB({ label, nilai }: { label: string; nilai: boolean }) {
  return (
    <label className="claimfacin__nb-pilihan">
      <input type="checkbox" disabled readOnly checked={nilai} /> {label}
    </label>
  )
}

export interface KolomNB {
  label: string
  angka?: boolean
}

const keKolom = (k: string | KolomNB): KolomNB => (typeof k === 'string' ? { label: k } : k)

/**
 * Grid NB hanya-baca. `rinci` = isi baris yang dibuka (▸ / ▾, beberapa baris boleh terbuka bersamaan seperti nbfacin);
 * `ukuran` = baris per halaman (pager hanya bila baris lebih banyak). Kosong = "No items".
 */
export function GridNB({
  kolom,
  baris,
  rinci,
  ukuran,
  kosong = TEKS_NB.kosong,
  label,
  selaluBuka = false,
}: {
  kolom: readonly (string | KolomNB)[]
  baris: readonly (readonly ReactNode[])[]
  rinci?: (i: number) => ReactNode
  ukuran?: number
  kosong?: string
  label?: string
  /** Rincian setiap baris selalu tampil tanpa tombol buka (grid coverage Marine Cargo nbfacin: baris + form). */
  selaluBuka?: boolean
}) {
  const [buka, setBuka] = useState<ReadonlySet<number>>(() => new Set())
  const [hal, setHal] = useState(1)
  const ks = kolom.map(keKolom)
  const nHal = ukuran ? Math.max(1, Math.ceil(baris.length / ukuran)) : 1
  const halIni = Math.min(hal, nHal)
  const awal = ukuran ? (halIni - 1) * ukuran : 0
  const potong = ukuran ? baris.slice(awal, awal + ukuran) : baris
  const panah = !!rinci && !selaluBuka
  const lebar = ks.length + (panah ? 1 : 0)
  const balik = (i: number) =>
    setBuka((s) => {
      const b = new Set(s)
      if (b.has(i)) b.delete(i)
      else b.add(i)
      return b
    })
  return (
    <div className="claimfacin__grid">
      {ukuran && baris.length > ukuran ? (
        <div className="claimfacin__pager">
          <button
            type="button"
            className="btn btn--ghost btn--sm"
            disabled={halIni <= 1}
            onClick={() => setHal(halIni - 1)}
            aria-label="Previous"
          >
            ‹
          </button>
          <span>
            Page {halIni} of {nHal}
          </span>
          <button
            type="button"
            className="btn btn--ghost btn--sm"
            disabled={halIni >= nHal}
            onClick={() => setHal(halIni + 1)}
            aria-label="Next"
          >
            ›
          </button>
        </div>
      ) : null}
      <table className="claimfacin__tabel" aria-label={label}>
        <thead>
          <tr>
            {panah && <th />}
            {ks.map((k, i) => (
              <th key={i}>{k.angka ? <span className="claimfacin__angka">{k.label}</span> : k.label}</th>
            ))}
          </tr>
        </thead>
        <tbody>
          {potong.length === 0 && (
            <tr>
              <td className="muted" colSpan={lebar}>
                {kosong}
              </td>
            </tr>
          )}
          {potong.map((sel, j) => {
            const i = awal + j
            const terbuka = selaluBuka || buka.has(i)
            return (
              <Fragment key={i}>
                <tr
                  className={panah ? 'inbox__baris' : undefined}
                  aria-expanded={panah ? terbuka : undefined}
                  onClick={panah ? () => balik(i) : undefined}
                >
                  {panah && (
                    <td>
                      <span className="claimfacin__panah" aria-hidden="true">
                        {terbuka ? '▾' : '▸'}
                      </span>
                    </td>
                  )}
                  {sel.map((v, c) => (
                    <td key={c}>{ks[c]?.angka ? <span className="claimfacin__angka">{v}</span> : v}</td>
                  ))}
                </tr>
                {rinci && terbuka && (
                  <tr className="claimfacin__baris-rinci">
                    <td colSpan={lebar}>
                      <div className="claimfacin__rinci">{rinci(i)}</div>
                    </td>
                  </tr>
                )}
              </Fragment>
            )
          })}
        </tbody>
      </table>
    </div>
  )
}
