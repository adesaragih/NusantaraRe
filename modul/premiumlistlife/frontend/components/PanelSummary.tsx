// Panel Summary di Premium List Detail — keputusan work owner 03-10-2026:
// "saat tekan save, perhitungan summary dijalankan, dan tampilkan juga".
//
// ⛔ Yang ditampilkan adalah rekap TERSIMPAN (`GET .../summary`), yang
// dihitung server saat Save Data / Save peserta CSV — persis yang nanti
// disalin Confirm. Layar TIDAK menghitung apa pun (pola `PremiumListSummary`).
//
// ⚠️ Kolom `COB` dan `PL NUMBER` grid summary TIDAK ditampilkan di sini: COB
// tidak punya sumber di baris rekap, dan PL Number baru terbit saat Confirm —
// keduanya selalu kosong di tahap ini.

import { useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { pemisahRibuan } from '../angka'
import { ambilRekapPolis, ambilRekapRetroPolis, type HasilRekapPolis, type RekapRetroPolis } from '../api'
import { PANEL_SPREADING, PANEL_SUMMARY } from '../labels'
import { kolomGridRekap, selRekap } from '../pages/PremiumListSummary'

/**
 * Kolom grid summary yang TIDAK ditampilkan di panel ini: COB dan PL NUMBER
 * (selalu kosong sebelum Confirm), serta TAX, PROF COMM, dan CLAIM
 * (keputusan work owner 03-10-2026). Hanya disembunyikan dari layar —
 * nilainya tetap dihitung, tersimpan, dan ikut rumus BALANCE.
 */
export const KOLOM_SUMMARY_TERSEMBUNYI: readonly string[] = ['COB', 'PL_NUMBER', 'TAX', 'PROF_COMM', 'CLAIM']

/** Kolom grid summary yang ditampilkan di panel ini. */
export function kolomPanelSummary(tipe: string) {
  return kolomGridRekap(tipe).filter((k) => !KOLOM_SUMMARY_TERSEMBUNYI.includes(k.kolom))
}

export default function PanelSummary({ polisID }: { polisID: string }) {
  const [hasil, setHasil] = useState<HasilRekapPolis | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  // Rekap retro per reinsurer per mata uang (keputusan work owner 06-10-2026).
  const [retro, setRetro] = useState<RekapRetroPolis[]>([])
  const [galatRetro, setGalatRetro] = useState<unknown>(null)

  useEffect(() => {
    let hidup = true
    ambilRekapRetroPolis(polisID).then(
      (r) => {
        if (hidup) setRetro(r)
      },
      (e: unknown) => {
        if (hidup) setGalatRetro(e)
      },
    )
    return () => {
      hidup = false
    }
  }, [polisID])

  useEffect(() => {
    let hidup = true
    ambilRekapPolis(polisID).then(
      (h) => {
        if (hidup) setHasil(h)
      },
      (e: unknown) => {
        if (hidup) setGalat(e)
      },
    )
    return () => {
      hidup = false
    }
  }, [polisID])

  const kolom = hasil === null ? [] : kolomPanelSummary(hasil.tipe)

  // Isi tab Summary di panel Participants (03-10-2026) - tanpa bingkai dan
  // judul sendiri; bingkainya milik panel Participants.
  return (
    <>
      {hasil === null && galat === null && <Memuat />}
      {galat !== null && <Gagal galat={galat} />}
      {hasil !== null && (hasil.rekap.length === 0 || kolom.length === 0) && (
        <Kosong pesan={PANEL_SUMMARY.kosong} />
      )}
      {hasil !== null && hasil.rekap.length > 0 && kolom.length > 0 && (
        <div className="pl-offer__riwayat-gulir">
          <table className="inbox__tabel">
            <thead>
              <tr>
                {kolom.map((k) => (
                  <th key={k.kolom}>{k.judul}</th>
                ))}
              </tr>
            </thead>
            <tbody>
              {hasil.rekap.map((r) => (
                <tr key={r.currency}>
                  {kolom.map((k) => (
                    <td key={k.kolom}>{pemisahRibuan(selRekap(r, k.kolom, ''))}</td>
                  ))}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
      {/*
        Retro per reinsurer per mata uang - total PREMIUM_SPREADED_GROSS / NET dari
        spreading tersimpan. Nama reinsurer dari data, TANPA daftar nama tetap Pega.
      */}
      {galatRetro !== null && <Gagal galat={galatRetro} />}
      {retro.length > 0 && (
        <div className="pl-spreading">
          <h4 className="pl-offer__subjudul">{PANEL_SPREADING.judulRekap}</h4>
          <table className="inbox__tabel">
            <thead>
              <tr>
                <th>{PANEL_SPREADING.kolomRekap.reinsurer}</th>
                <th>{PANEL_SPREADING.kolomRekap.currency}</th>
                <th className="pl-spreading__angka">{PANEL_SPREADING.kolomRekap.gross}</th>
                <th className="pl-spreading__angka">{PANEL_SPREADING.kolomRekap.net}</th>
              </tr>
            </thead>
            <tbody>
              {retro.map((r) => (
                <tr key={`${r.currency}-${r.reinsurer}`}>
                  <td>{r.reinsurer}</td>
                  <td>{r.currency}</td>
                  <td className="pl-spreading__angka">{pemisahRibuan(r.gross)}</td>
                  <td className="pl-spreading__angka">{pemisahRibuan(r.net)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </>
  )
}
