// Tab **Achievement In IDR** (cabang PROPORSIONAL) — grid enam kolom
// berikut tiga sel kaki.
//
// ---------------------------------------------------------------------
// ⛔ YANG DIGANTI
// ---------------------------------------------------------------------
// Nama tabnya sudah ada di `TAB_PROPORSIONAL` sejak lama, tetapi NOL cabang
// merendernya — membukanya menampilkan isi tab pertama. Rumusnya sendiri
// sudah hidup di `backend/services/hitung_achievement.go` dan dipakai
// sub-tab Achievement di dalam Limits Prop; yang hilang hanya tampilan
// RINGKASNYA.
//
// ---------------------------------------------------------------------
// ⭐ ASAL ANGKANYA — `Activity/GetAchievement.xml`
// ---------------------------------------------------------------------
//   RDB-List `GetAchievement`   InputParam.CARI1 = TreatyIn.ID
//                               InputParam.CARI2 = Treaty Group
//     CARI4 TreatyGroup · CARI5 TreatyType · CARI10 PREMIUM
//     CARI13 NETPREMIUM · CARI14 PaidClaim · CARI15 OutstandingClaim
//   RDB-List kurs               .Conversion (IDR => 1)
//   GetHistoryClaim_act / GetEstimasiClaim_act   Cash Call
//
// lalu per Treaty Group (langkah 12.2):
//   TotalAchPremium    = SIGMA PREMIUM    x Conversion
//   TotalAchNetPremium = SIGMA NETPREMIUM x Conversion
//   TotalAchIncured    = SIGMA (PaidClaim + CashCall + OutstandingClaim
//                               + OutstandingCashCall) x Conversion
//   LossRatio          = TotalAchIncured / TotalAchPremium x 100
//
// ⚠️ KOLOMNYA BERJUDUL `Net Loss Ratio`, TETAPI PEMBAGINYA GROSS. Ekspor
// langkah 12.2.5 berbunyi `@divide(.TotalAchIncured, .TotalAchPremium, 20)`
// — `TotalAchPremium`, bukan `TotalAchNetPremium`. Rasio ber-pembagi Net
// memang dihitung di langkah yang sama (`Local.LossRatioNet`), tetapi ia
// masuk ke `CurrencyList`, bukan ke kolom ini. Dinyatakan, bukan
// diam-diam diperbaiki: memperbaikinya membuat angka kita berbeda dari
// layar lama tanpa seorang pun memutuskannya.
//
// ⛔ SELURUH SEL HANYA-BACA — keenamnya `pyReadOnly = true` TANPA
// `pyReadOnlyCondition`. Berbeda dengan tab EGNPI dan Maximum Retention,
// yang sel serupanya BERSYARAT dan karena itu aktif di mode Edit.

import { useEffect, useState } from 'react'

import { Panel } from '../../../../inti/frontend/components/ui/dasar'
import { hitungAchievement, type RingkasanAchievement, type SimpulLimit } from '../api'
import {
  ACHIEVEMENT,
  DESIMAL_ACHIEVEMENT,
  KAKI_ACHIEVEMENT,
  KOLOM_ACHIEVEMENT,
} from '../labelsAchievement'
import { useProperti } from '../halaman'
import { formatLimit, teksDari } from './TabLimitsProp'

const KOSONG: RingkasanAchievement = {
  baris: [],
  SumTotalAchievNetPremium: '',
  SumTotalAchievIncured: '',
  SumLossRatio: '',
}

/** Satu sel menurut golongan kolomnya. */
function sel(jenis: string, nilai: string): string {
  if (jenis === 'teks') return nilai
  return formatLimit('uang', jenis === 'persen' ? DESIMAL_ACHIEVEMENT.persen : DESIMAL_ACHIEVEMENT.uang, nilai)
}

/** Larik di simpul — yang tidak ada → kosong. */
function larikDari(s: SimpulLimit, kunci: string): SimpulLimit[] {
  const v = s[kunci]
  return Array.isArray(v) ? v : []
}

export default function TabAchievement({
  idKontrak,
  pohon,
}: {
  idKontrak: string
  /** `TreatyIn.Limits` — sumber Treaty Group dan EPI tiap Detail. */
  pohon: readonly SimpulLimit[]
}) {
  // ⭐ `TreatyIn.Limits` dari PENAMPUNG HALAMAN — Limits yang sedang diisi
  // di tab Limits, bukan salinan kontrak yang dimuat.
  const [limits] = useProperti<SimpulLimit[]>('Limits', () => [...pohon])
  const [r, setR] = useState<RingkasanAchievement>(KOSONG)
  const [gagal, setGagal] = useState('')

  // ⭐ Diambil saat tab dibuka. Di Pega angkanya sudah ada di clipboard
  // sebab `GetAchievement` dijalankan Refresh tab Limits; layar ini tidak
  // punya clipboard bersama, jadi ia memanggil rumus yang SAMA sendiri.
  useEffect(() => {
    let dibuang = false
    hitungAchievement({
      idKontrak,
      cari: false,
      asAt: '',
      tahun: '',
      rnmShareP: '',
      limits: limits.map((l) => ({
        TreatyType: teksDari(l, 'TreatyType'),
        Detail: larikDari(l, 'Detail').map((d) => ({
          TreatyGroup: teksDari(d, 'TreatyGroup'),
          EPIList: larikDari(d, 'EPIList'),
        })),
      })),
    })
      .then((h) => {
        if (!dibuang) setR(h.ringkasan ?? KOSONG)
      })
      .catch((e: unknown) => {
        if (!dibuang) setGagal(e instanceof Error ? e.message : String(e))
      })
    return () => {
      dibuang = true
    }
  }, [idKontrak, limits])

  const nilaiKaki = (kunci: string | null): string =>
    kunci === null ? '' : sel(kunci === 'SumLossRatio' ? 'persen' : 'uang', (r as unknown as Record<string, string>)[kunci] ?? '')

  return (
    <Panel judul={ACHIEVEMENT.judul}>
      {gagal !== '' && (
        <p className="tl-pesan" role="alert">
          {ACHIEVEMENT.gagalMuat}: {gagal}
        </p>
      )}

      {/* ⭐ Bentuk Pega (8 Oktober 2026): grid tetap tampil tanpa baris. */}
      <div className="table-wrap">
          <table className="trin__tabel trin__tabel--pega trin__ach">
            <thead>
              <tr>
                {KOLOM_ACHIEVEMENT.map((k) => (
                  <th key={k.kunci} scope="col">
                    {k.label}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {r.baris.length === 0 && (
                <tr>
                  <td colSpan={KOLOM_ACHIEVEMENT.length} className="trin__kosong-pega">
                    {ACHIEVEMENT.tanpaBaris}
                  </td>
                </tr>
              )}
              {r.baris.map((b, i) => (
                // ⭐ Baris terakhir adalah " Total In IDR" (`.TreatyType`),
                // bukan data — ditandai supaya terbaca sebagai total.
                <tr
                  key={String(i)}
                  className={b.TreatyType === ACHIEVEMENT.barisTotal ? 'trin__ach-total' : undefined}
                >
                  {KOLOM_ACHIEVEMENT.map((k) => (
                    <td key={k.kunci}>{sel(k.jenis, (b as unknown as Record<string, string>)[k.kunci] ?? '')}</td>
                  ))}
                </tr>
              ))}
            </tbody>
            {/* ⛔ KAKI HANYA TIGA SEL — `Gross Premium` nol punya total.
                Tiga kolom pertama diisi `Spacers` ber-`pyVisible = OTHER`
                di ekspor, jadi memang kosong. Label `Total in IDR`
                ber-`pyVisible = ALWAYS`, jadi ia TAMPIL walau
                `pyCondition`-nya `1=2`: syarat itu hanya berlaku ketika
                `pyVisible = OTHER`. */}
            <tfoot>
              <tr>
                <th scope="row" colSpan={3}>
                  {ACHIEVEMENT.totalIDR}
                </th>
                {KAKI_ACHIEVEMENT.slice(3).map((k, i) => (
                  <td key={String(i)}>{nilaiKaki(k)}</td>
                ))}
              </tr>
            </tfoot>
          </table>
        </div>
    </Panel>
  )
}
