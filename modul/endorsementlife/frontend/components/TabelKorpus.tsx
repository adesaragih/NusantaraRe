// Tabel dan daftar medan dari definisi kolom berbukti korpus (`kolomKorpus.ts`).

import type { KolomKorpus } from '../kolomKorpus'
import { UMUM_EDM } from '../labels'
import { kelasAngka, sel, selAngka, selTanggal } from '../tampilan'

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
  tandaiNegatif = false,
}: {
  kolom: readonly KolomKorpus[]
  baris: ReadonlyArray<Record<string, string>>
  kunci: (b: Record<string, string>, i: number) => string
  nomorAwal?: number
  /** Sel angka negatif berkelas merah (`edm-negatif`) - tabel rekap kasus (08-10-2026); bawaan tidak. */
  tandaiNegatif?: boolean
}) {
  return (
    <div className="edm-gulir">
      <table className="inbox__tabel">
        <thead>
          <tr>
            {nomorAwal !== undefined && <th className="edm-nomor">{UMUM_EDM.nomor}</th>}
            {/* Kepala kolom angka rata kanan, sejajar angkanya (tampilan 08-10-2026). */}
            {kolom.map((k) => (
              <th key={k.label + k.baris} className={k.jenis === 'n' ? 'edm-angka' : undefined}>
                {k.label}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {baris.map((b, i) => (
            <tr key={kunci(b, i)}>
              {nomorAwal !== undefined && <td className="edm-nomor">{nomorAwal + i}</td>}
              {kolom.map((k) => (
                <td
                  key={k.label + k.baris}
                  className={k.jenis !== 'n' ? undefined : tandaiNegatif ? kelasAngka(b[k.kolom]) : 'edm-angka'}
                >
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
export function MedanKorpus({
  kolom,
  nilai,
  className = 'edm-kepala',
  lebar = [],
  kelasLebar = '',
}: {
  kolom: readonly KolomKorpus[]
  nilai: Record<string, string>
  /** Kelas kisi; bawaan `edm-kepala`. */
  className?: string
  /** Kolom korpus yang selebar kisi (diberi `kelasLebar`); bawaan tidak ada. */
  lebar?: readonly string[]
  kelasLebar?: string
}) {
  return (
    <dl className={className}>
      {kolom.map((k) => (
        <div
          key={k.label + k.baris}
          className={kelasLebar && lebar.includes(k.kolom) ? `edm-kepala__medan ${kelasLebar}` : 'edm-kepala__medan'}
        >
          <dt>{k.label}</dt>
          <dd>{isiKolom(k, nilai)}</dd>
        </div>
      ))}
    </dl>
  )
}
