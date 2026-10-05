// Grid baris AGGREGATE - padanan grid TempCSV `ShowAggregateList` (pratinjau unggah) dan rincian satu unggahan.
// Kolom dan urutannya `KOLOM_GRID`; kepala dua baris dengan grup NoR/IA per kategori, Total, dan RNM (`GRUP_KOLOM`,
// perintah work owner 05-10-2026 "grouping sama seperti bordereaux"). HANYA DIBACA (work owner 04-10-2026: "upload csv nya read only"); angka tampil
// format Indonesia dua desimal, nilai penuh di `title`. Berhalaman di peramban seperti `pyGridPaginator` Pega.

import { useState } from 'react'

import type { Baris } from '../api'
import { barisTotal, formatAngka, jumlahHalaman, KOLOM_GRID, susunKepala } from '../aturan'
import { AG, GRUP_KOLOM, LABEL_KOLOM } from '../labels'

const KEPALA = susunKepala(KOLOM_GRID, GRUP_KOLOM, LABEL_KOLOM)

const UKURAN = 100

export default function GridAggregate({ baris, sertaID = false }: { baris: Baris[]; sertaID?: boolean }) {
  const [halaman, setHalaman] = useState(1)
  const jumlah = jumlahHalaman(baris.length, UKURAN)
  const kini = Math.min(halaman, jumlah)
  const awal = (kini - 1) * UKURAN
  return (
    <>
      <div className="aggregate__gulir aggregate__gulir--grid">
        <table className="inbox__tabel aggregate__grid">
          <thead>
            <tr>
              <th rowSpan={2}>{AG.no}</th>
              {sertaID && <th rowSpan={2}>{LABEL_KOLOM['ID']}</th>}
              {KEPALA.baris1.map((s) => (
                <th
                  key={s.kunci}
                  colSpan={s.colSpan}
                  rowSpan={s.rowSpan}
                  title={s.grup ? undefined : s.judul}
                  className={s.grup ? 'aggregate__kepala-grup' : undefined}
                >
                  {s.teks}
                </th>
              ))}
            </tr>
            <tr>
              {KEPALA.baris2.map((s) => (
                <th key={s.kunci} title={s.judul}>
                  {s.teks}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {baris.slice(awal, awal + UKURAN).map((b, n) => {
              const i = awal + n
              return (
                <tr key={i} className={barisTotal(b) ? 'inbox__baris aggregate__total' : 'inbox__baris'}>
                  <td>{i + 1}</td>
                  {sertaID && <td>{b['ID']}</td>}
                  {KOLOM_GRID.map((k) => {
                    const v = b[k.nama] ?? ''
                    return k.angka ? (
                      <td key={k.nama} className="aggregate__angka" title={v === '' ? undefined : v}>
                        {formatAngka(v)}
                      </td>
                    ) : (
                      <td key={k.nama}>{v}</td>
                    )
                  })}
                </tr>
              )
            })}
          </tbody>
        </table>
      </div>
      {jumlah > 1 && (
        <div className="aggregate__halaman">
          <span className="toolbar__spacer" />
          <button
            type="button"
            className="btn btn--ghost"
            disabled={kini <= 1}
            onClick={() => {
              setHalaman(kini - 1)
            }}
          >
            {AG.sebelumnya}
          </button>
          <span>{AG.halaman(kini, jumlah)}</span>
          <button
            type="button"
            className="btn btn--ghost"
            disabled={kini >= jumlah}
            onClick={() => {
              setHalaman(kini + 1)
            }}
          >
            {AG.berikutnya}
          </button>
        </div>
      )}
    </>
  )
}
