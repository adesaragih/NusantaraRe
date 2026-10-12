// Summary Comparison - perbandingan summary sebelum/sesudah endorsement per mata uang (permintaan work owner
// 08-10-2026, BUKAN korpus). Baca-saja: `GET /kasus/{id}/perbandingan`; kolom uang = kolom rekap per Type yang
// sama dengan tabel rekap (`MATA_UANG`). Hanya dirender SESUDAH Process Endorsement (keputusan work owner
// 10-10-2026, `InputEDMLife`) - cabang "belum diproses" (Deleted "—" bercatatan) dihapus.
//
// Tampilan dirapikan (08-10-2026): satu blok per mata uang bersub-judul; kepala angka rata kanan; Participants
// = jumlah orang tanpa tanda/warna (09-10-2026); After tebal bergaris atas; Difference disorot dan bertanda. Nilai
// uang negatif merah (`kelasNegatif`); Difference uang positif hijau (08-10-2026).

import { useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { ambilPerbandingan, type BarisPerbandinganEDM } from '../api'
import { MATA_UANG, type KolomKorpus } from '../kolomKorpus'
import { PERBANDINGAN_EDM, UMUM_EDM } from '../labels'
import { arahSelisih, kelasAngka, selAngka, selBertanda, teksPeserta, URUTAN_PERBANDINGAN } from '../tampilan'
import { tipePopup } from './PolisLama'

/** Kelas baris perbandingan. */
function kelasBaris(kunci: string): string | undefined {
  if (kunci === 'after') return 'edm-banding__total'
  if (kunci === 'difference') return 'edm-banding__selisih'
  return undefined
}

/** Satu blok mata uang (murni). */
export function BlokPerbandingan({
  b,
  kolom,
}: {
  b: BarisPerbandinganEDM
  kolom: readonly KolomKorpus[]
}) {
  return (
    <div className="edm-banding__blok">
      <h4 className="edm-banding__uang">{b.currency}</h4>
      <div className="edm-gulir">
        <table className="inbox__tabel">
          <thead>
            <tr>
              <th>{PERBANDINGAN_EDM.kolomBaris}</th>
              <th className="edm-angka">{PERBANDINGAN_EDM.kolomPeserta}</th>
              {kolom.map((k) => (
                <th key={k.kolom} className="edm-angka">
                  {k.label}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {URUTAN_PERBANDINGAN.map(({ kunci, label }) => {
              const sel = b[kunci]
              const selisih = kunci === 'difference'
              return (
                <tr key={kunci} className={kelasBaris(kunci)}>
                  <td>{label}</td>
                  {/* Participants: jumlah orang tanpa tanda/warna (09-10-2026). */}
                  <td className="edm-angka">{teksPeserta(sel.cacah)}</td>
                  {kolom.map((k) => {
                    const v = sel.nilai[k.kolom]
                    if (!selisih) {
                      return (
                        <td key={k.kolom} className={kelasAngka(v)}>
                          {selAngka(v)}
                        </td>
                      )
                    }
                    // Difference: negatif merah (kelasAngka), positif hijau.
                    return (
                      <td key={k.kolom} className={kelasAngka(v)}>
                        {arahSelisih(v) === 'naik' ? <span className="edm-banding__naik">{selBertanda(v)}</span> : selBertanda(v)}
                      </td>
                    )
                  })}
                </tr>
              )
            })}
          </tbody>
        </table>
      </div>
    </div>
  )
}

/** Seluruh mata uang (murni). Catatan rumus di bawah tabel dihapus (permintaan work owner 08-10-2026). */
export function TabelPerbandingan({
  baris,
  kolom,
}: {
  baris: readonly BarisPerbandinganEDM[]
  /** Kolom uang (tanpa CURRENCY). */
  kolom: readonly KolomKorpus[]
}) {
  return (
    <>
      {baris.map((b) => (
        <BlokPerbandingan key={b.currency} b={b} kolom={kolom} />
      ))}
    </>
  )
}

/** Kolom uang perbandingan = kolom rekap per Type (`MATA_UANG`) tanpa CURRENCY. */
export function kolomPerbandingan(tipe: string): readonly KolomKorpus[] {
  return MATA_UANG[tipePopup(tipe)].filter((k) => k.kolom !== 'CURRENCY')
}

export default function PerbandinganRekap({
  kasusId,
  tipe,
  versi,
}: {
  kasusId: string
  tipe: string
  /** Naik sesudah Add CSV Data / Process Endorsement - dimuat ulang. */
  versi: number
}) {
  const [baris, setBaris] = useState<BarisPerbandinganEDM[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)

  useEffect(() => {
    let hidup = true
    setGalat(null)
    ambilPerbandingan(kasusId).then(
      (b) => {
        if (hidup) setBaris(b)
      },
      (e: unknown) => {
        if (hidup) setGalat(e)
      },
    )
    return () => {
      hidup = false
    }
  }, [kasusId, versi])

  return (
    <section className="edm-banding">
      <h3 className="edm-subjudul">{PERBANDINGAN_EDM.judul}</h3>
      {galat !== null && <Gagal galat={galat} />}
      {galat === null && baris === null && <Memuat pesan={UMUM_EDM.memuat} />}
      {baris !== null && baris.length === 0 && <Kosong pesan={PERBANDINGAN_EDM.kosong} />}
      {baris !== null && baris.length > 0 && (
        <TabelPerbandingan baris={baris} kolom={kolomPerbandingan(tipe)} />
      )}
    </section>
  )
}
