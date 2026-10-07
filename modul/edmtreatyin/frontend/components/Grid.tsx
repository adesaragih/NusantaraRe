// Grid hanya-baca berkolom (`xol.ts` `Kolom`) dan grid ber-baris LIPAT - padanan grid Pega ber-`pyEditingMode=
// expandPane` (`DetailPolicyTreatyInAddPremi`, `DetailPolicyTreatyInAddendum` S13): setiap baris membuka panel
// rinciannya (flow action `DetailPolicyAddPremiDetail` / `InstallmentList`) di bawahnya. Pola tampilan:
// `modul/nbtreatyin/frontend/components/DetailNonProp.tsx` (06-10-2026: grid per mata uang dapat dilipat, rincian
// bernomor, terbuka sejak awal).

import { Fragment, useState, type ReactNode } from 'react'

import type { Baris } from '../api'
import { kelasAngka, nilaiNol, sajikan } from '../sajian'
import type { Kolom } from '../xol'

/** Teks satu sel: sel LABEL = teksnya; `bilaAda` kosong = kosong; berformat bila `format` ada. */
export function teksSel(b: Baris, c: Kolom): string {
  if (c.m === '') return c.teks ?? ''
  const v = b[c.m] ?? ''
  if (c.bilaAda && v.trim() === '') return ''
  return c.format === undefined ? v : sajikan(v, c.format)
}

/** Kelas sel: angka rata kanan (perintah work owner 06-10-2026) - termasuk kosong yang tampil "0" (`nolPolos`). */
function kelasSel(b: Baris, c: Kolom): string | undefined {
  if (c.m === '' || c.format === undefined || c.format === 'tanggal') return undefined
  return kelasAngka(b[c.m]) ?? (c.format.nolPolos && nilaiNol(b[c.m]) ? 'edmt__angka' : undefined)
}

/** Kelas judul: rata kanan bila kolomnya berformat angka. */
function kelasJudul(c: Kolom): string | undefined {
  return c.m !== '' && c.format !== undefined && c.format !== 'tanggal' ? 'edmt__angka' : undefined
}

export function GridBaca({ baris, kolom, nomor = false }: { baris: Baris[]; kolom: Kolom[]; nomor?: boolean }) {
  return (
    <div className="table-wrap edmt__grid-np">
      <table>
        <thead>
          <tr>
            {nomor && <th scope="col" aria-label="No" />}
            {kolom.map((c, j) => (
              <th key={`${c.m}-${j}`} scope="col" className={kelasJudul(c)}>
                {c.judul}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {baris.map((b, i) => (
            <tr key={i}>
              {nomor && <td className="edmt__urut-np">{i + 1}</td>}
              {kolom.map((c, j) => (
                <td key={`${c.m}-${j}`} className={kelasSel(b, c)}>
                  {teksSel(b, c)}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

/** Grid ber-baris lipat: sel pertama tiap baris = tombol buka/tutup rinciannya; rincian terbuka sejak awal. */
export function GridLipat({
  baris,
  kolom,
  rinci,
}: {
  baris: Baris[]
  kolom: Kolom[]
  /** Isi panel rincian baris ke-`i` (berbasis 0). */
  rinci: (b: Baris, i: number) => ReactNode
}) {
  const [tertutup, setTertutup] = useState<ReadonlySet<number>>(new Set())
  const [awal, ...sisa] = kolom
  if (awal === undefined) return null
  const balik = (i: number) =>
    setTertutup((s) => {
      const n = new Set(s)
      if (n.has(i)) n.delete(i)
      else n.add(i)
      return n
    })
  return (
    <div className="table-wrap edmt__grid-np">
      <table>
        <thead>
          <tr>
            {kolom.map((c, j) => (
              <th key={`${c.m}-${j}`} scope="col" className={kelasJudul(c)}>
                {c.judul}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {baris.map((b, i) => {
            const buka = !tertutup.has(i)
            return (
              <Fragment key={i}>
                <tr>
                  <td>
                    <button type="button" className="edmt__lipat" aria-expanded={buka} onClick={() => balik(i)}>
                      <span className="edmt__lipat-ikon" aria-hidden="true">
                        {buka ? '▾' : '▸'}
                      </span>
                      {teksSel(b, awal)}
                    </button>
                  </td>
                  {sisa.map((c, j) => (
                    <td key={`${c.m}-${j}`} className={kelasSel(b, c)}>
                      {teksSel(b, c)}
                    </td>
                  ))}
                </tr>
                {buka && (
                  <tr className="edmt__rinci-baris">
                    <td colSpan={kolom.length}>{rinci(b, i)}</td>
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
