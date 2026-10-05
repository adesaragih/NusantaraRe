// Grid detail satu kombinasi Type x Business - padanan 28 section `DetailPremium*` / `DetailClaim*` /
// `DetailSubrogationBonding`: hanya dibaca, 50 baris per halaman, bernomor. Kolom = urutan CSV (perbaikan bug grid Pega:
// kolom ganda, kolom salah ikat, grid klaim PA/Offshore yang tertukar). Kombinasi yang sudah dipetakan ke format Excel
// bordereaux 2025 (Premium Fire dulu, perintah work owner 05-10-2026) memakai kepala sheet-nya: dua baris dengan sel
// gabungan, judul persis Excel, gaya kepala Excel (`bordereaux__grid--excel`).

import { useState } from 'react'

import type { Baris, Kolom } from '../api'
import { formatAngka, jumlahHalaman, judulKepala, susunKepala } from '../aturan'
import { BDX } from '../labels'

const UKURAN = 50

export default function GridDetail({ kolom, baris }: { kolom: Kolom[]; baris: Baris[] }) {
  const [halaman, setHalaman] = useState(1)
  const jumlah = jumlahHalaman(baris.length, UKURAN)
  const kini = Math.min(halaman, jumlah)
  const awal = (kini - 1) * UKURAN
  const kepala = susunKepala(kolom)
  const rentang = kepala.duaBaris ? 2 : 1
  return (
    <>
      <div className={kepala.excel ? 'bordereaux__gulir bordereaux__gulir--grid' : 'bordereaux__gulir'}>
        <table
          className={
            kepala.excel ? 'inbox__tabel bordereaux__grid bordereaux__grid--excel' : 'inbox__tabel bordereaux__grid'
          }
        >
          <thead>
            <tr>
              <th rowSpan={rentang}>{kepala.excel ? BDX.no : '#'}</th>
              {kepala.baris1.map((s) => (
                <th
                  key={s.kunci}
                  colSpan={s.colSpan}
                  rowSpan={s.rowSpan}
                  className={
                    s.kunci.startsWith('grup:')
                      ? 'bordereaux__kepala-grup'
                      : s.angka && !kepala.excel
                        ? 'bordereaux__angka'
                        : undefined
                  }
                >
                  {s.teks}
                </th>
              ))}
            </tr>
            {kepala.duaBaris && (
              <tr>
                {kepala.baris2.map((k) => (
                  <th key={k.kolom}>{judulKepala(k)}</th>
                ))}
              </tr>
            )}
          </thead>
          <tbody>
            {baris.slice(awal, awal + UKURAN).map((b, n) => (
              <tr key={b['ID'] ?? awal + n} className="inbox__baris">
                <td>{awal + n + 1}</td>
                {kepala.tampil.map((k) => {
                  const v = b[k.kolom] ?? ''
                  return k.jenis === 'angka' ? (
                    <td key={k.kolom} className="bordereaux__angka" title={v === '' ? undefined : v}>
                      {formatAngka(v)}
                    </td>
                  ) : (
                    <td key={k.kolom}>{v}</td>
                  )
                })}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      {jumlah > 1 && (
        <div className="bordereaux__halaman">
          <span className="muted">{BDX.baris(baris.length)}</span>
          <span className="toolbar__spacer" />
          <button type="button" className="btn btn--ghost" disabled={kini <= 1} onClick={() => setHalaman(kini - 1)}>
            {BDX.sebelumnya}
          </button>
          <span>{BDX.halaman(kini, jumlah)}</span>
          <button
            type="button"
            className="btn btn--ghost"
            disabled={kini >= jumlah}
            onClick={() => setHalaman(kini + 1)}
          >
            {BDX.berikutnya}
          </button>
        </div>
      )}
    </>
  )
}
