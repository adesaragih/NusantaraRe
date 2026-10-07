// Tabel dan daftar medan dari definisi kolom berbukti korpus (`kolomKorpus.ts`).

import type { KolomKorpus } from '../kolomKorpus'
import { UMUM_EDM } from '../labels'
import { sel, selAngka, selTanggal } from '../tampilan'

/** Isi satu sel menurut jenis kolomnya; angka tetap teks (ADR-0003). */
export function isiKolom(k: KolomKorpus, nilai: Record<string, string>): string {
  const v = nilai[k.kolom] ?? ''
  if (k.jenis === 'd') return selTanggal(v)
  if (k.jenis === 'n') return selAngka(v)
  return sel(v)
}

/**
 * Tabel ber-kepala VERBATIM; `kunci` memilih kunci baris.
 *
 * `nomorAwal` (opsional, 07-10-2026): bila diisi, kolom pertama "No" berisi `nomorAwal + indeks` - nomor
 * TAMPILAN saja (bukan kolom korpus, tidak memengaruhi kunci baris). Tanpa `nomorAwal` tabel tidak berubah.
 */
export function TabelKorpus({
  kolom,
  baris,
  kunci,
  nomorAwal,
}: {
  kolom: readonly KolomKorpus[]
  baris: ReadonlyArray<Record<string, string>>
  kunci: (b: Record<string, string>, i: number) => string
  nomorAwal?: number
}) {
  return (
    <div className="edm-gulir">
      <table className="inbox__tabel">
        <thead>
          <tr>
            {nomorAwal !== undefined && <th className="edm-nomor">{UMUM_EDM.nomor}</th>}
            {kolom.map((k) => (
              <th key={k.label + k.baris}>{k.label}</th>
            ))}
          </tr>
        </thead>
        <tbody>
          {baris.map((b, i) => (
            <tr key={kunci(b, i)}>
              {nomorAwal !== undefined && <td className="edm-nomor">{nomorAwal + i}</td>}
              {kolom.map((k) => (
                <td key={k.label + k.baris} className={k.jenis === 'n' ? 'edm-angka' : undefined}>
                  {isiKolom(k, b)}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

/** Medan baca-saja berlabel VERBATIM (bentuk `PL_Detail_Sec`). */
export function MedanKorpus({ kolom, nilai }: { kolom: readonly KolomKorpus[]; nilai: Record<string, string> }) {
  return (
    <dl className="edm-kepala">
      {kolom.map((k) => (
        <div key={k.label + k.baris} className="edm-kepala__medan">
          <dt>{k.label}</dt>
          <dd>{isiKolom(k, nilai)}</dd>
        </div>
      ))}
    </dl>
  )
}
