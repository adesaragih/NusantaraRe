// Isi satu tab layout group `Section/DetailPolicyTreatyInAddGeneralEditable` S2108 - kelima section Prop
// (`PropOldData`, `PropOldData2`, `PropNewData`, `PropNewData2`, `PropValueDifference`) SATU kerangka (`medan.ts`
// `tataTab`):
//
//   S2 Gross | S3 { S4 "OGP" ..., S5 "ONP" ... } | grid `.SpreadingRiskList` + kaki "Total" | total berlabel
//   (PropNewData) | `.Installment` + "Installment Data Information". Tombol "Calculate Value Difference" (NewData2 S24)
//   DIBUANG (WO 08-10-2026 "COBA CEK TOMBOLITU, APAKAH MASIH DIPERLUKAN? KALAU SUDAH TIDAK DIHAPUS AJA!"): sel yang mengubah data baru menghitung ulang selisih sendiri
//
// Tata letak bagian uang mengikuti layar NB yang disetujui work owner (`modul/nbtreatyin/frontend/pages/LayarKasus.tsx`
// 06-10-2026): Gross di atas, kolom kiri OGP lalu klaim/saldo, kolom kanan ONP lalu potongan/pajak.
//
// ⛔ Grid spreading PropNewData2: tombol Add (addRow) SAJA - TANPA Delete (keputusan work owner, spec-penyimpanan
// ID-16: baris endorsemen tidak dapat dihapus). Nol perhitungan di sini: setiap action set -> `onRefresh`.

import { Fragment } from 'react'

import { Panel, type Opsi } from '../../../../inti/frontend/components/ui/dasar'
import { daftar, nilai, type Baris, type Halaman, type Pilihan } from '../api'
import { BAGIAN, KOLOM_ANGSURAN, KOLOM_SPREADING } from '../labels'
import {
  AKSI_SPREADING,
  SAJIAN_ANGSURAN,
  SAJIAN_SPREADING,
  labelTreatyType,
  medanTampil,
  type Aksi,
  type Medan,
  type SumberAcuan,
  type TataTab,
} from '../medan'
import { sajikan } from '../sajian'
import InputAngka from './InputAngka'
import KotakMedan from './KotakMedan'
import Wadah from './Wadah'

export interface PropsTabData {
  tata: TataTab
  halaman: Halaman
  wajib: ReadonlySet<string>
  /** Layar boleh disunting (pelaku anggota antrean, kasus terbuka). */
  boleh: boolean
  opsi: Record<SumberAcuan, Opsi[]>
  /** Teks `.TreatyType` grid spreading hanya-baca (RD `BrowseReinsuranceType_RD`, acuan `jenisReas`). */
  opsiJenisReas: Pilihan[]
  onUbah: (jalur: string, v: string) => void
  onSelesai: (m: Medan, v: string) => void
  onUbahBaris: (jalur: string, i: number, kunci: string, v: string) => void
  /** Action set: `urutan` berurutan, `indeks` berbasis 1. */
  onRefresh: (urutan: Aksi[], indeks?: number) => void
}

export default function TabData({
  tata,
  halaman: h,
  wajib,
  boleh,
  opsi,
  opsiJenisReas,
  onUbah,
  onSelesai,
  onUbahBaris,
  onRefresh,
}: PropsTabData) {
  const a = tata.awalan
  const sunting = boleh && tata.varian === 'baruAdmin'
  const SPREADING = a + 'SpreadingRiskList'
  const spreading = daftar(h, SPREADING)
  const angsuran = daftar(h, a + 'ListInstallment')

  const kotak = (m: Medan, i: number) => (
    <KotakMedan
      key={`${m.jalur}-${m.label}-${i}`}
      medan={m}
      halaman={h}
      wajib={wajib.has(m.jalur)}
      hanyaBaca={!sunting}
      opsi={opsi}
      onUbah={onUbah}
      onSelesai={onSelesai}
    />
  )
  const kolom = (ms: Medan[]) => <div className="edmt__kolom">{medanTampil(ms, h).map(kotak)}</div>

  /** Sel persen spreading PropNewData2: change -> refresh CountSpreading_Act(Index=.pxListSubscript) + selisih. */
  const persen = (b: Baris, i: number, k: 'SharePercentage' | 'ClaimPercentage', label: string) =>
    sunting ? (
      <InputAngka
        label={label}
        value={b[k] ?? ''}
        sajian={SAJIAN_SPREADING.persen}
        onChange={(v) => onUbahBaris(SPREADING, i, k, v)}
        onBlur={() => onRefresh(AKSI_SPREADING, i + 1)}
      />
    ) : (
      sajikan(b[k] ?? '', SAJIAN_SPREADING.persen)
    )

  return (
    <>
      {/* S2 / S3 NOHEADER; LABEL Heading 4 "OGP" / "ONP" */}
      <Wadah>
        <div className="edmt__dua-kolom">{kolom(tata.gross)}</div>
        <div className="edmt__dua-kolom">
          {[tata.kiri, tata.kanan].map((ks, i) => (
            <div key={i} className="edmt__kolom-grup">
              {ks.map((k, j) => (
                <Fragment key={j}>
                  {k.judul && <h4>{k.judul}</h4>}
                  {kolom(k.medan)}
                </Fragment>
              ))}
            </div>
          ))}
        </div>
      </Wadah>

      {/* grid `.SpreadingRiskList` - wadah NOHEADER */}
      <Wadah>
        {/* Keputusan work owner 07-10-2026 ("INI READ ONLY JUGA", sama dengan NB): tombol Add DIBUANG, Type Treaty
            hanya-baca - baris dibawa Choose Business dari generasi lama; yang tersunting hanya % Share / % Share klaim */}
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th scope="col">{KOLOM_SPREADING.treatyType}</th>
                <th scope="col" className="edmt__angka">
                  {KOLOM_SPREADING.share}
                </th>
                <th scope="col" className="edmt__angka">
                  {KOLOM_SPREADING.premium}
                </th>
                <th scope="col" className="edmt__angka">
                  {KOLOM_SPREADING.claimPct}
                </th>
                <th scope="col" className="edmt__angka">
                  {KOLOM_SPREADING.claim}
                </th>
              </tr>
            </thead>
            <tbody>
              {spreading.map((b, i) => (
                <tr key={i}>
                  <td>{labelTreatyType(b, opsiJenisReas)}</td>
                  <td className="edmt__angka">{persen(b, i, 'SharePercentage', KOLOM_SPREADING.share)}</td>
                  <td className="edmt__angka">{sajikan(b.PremiumSpreaded ?? '', SAJIAN_SPREADING.uang)}</td>
                  <td className="edmt__angka">{persen(b, i, 'ClaimPercentage', KOLOM_SPREADING.claimPct)}</td>
                  <td className="edmt__angka">{sajikan(b.ClaimSpreaded ?? '', SAJIAN_SPREADING.uang)}</td>
                </tr>
              ))}
            </tbody>
            <tfoot>
              <tr>
                <td>{KOLOM_SPREADING.total}</td>
                {(
                  ['TotalSharePercentagePremium', 'TotalPremium', 'TotalSharePercentageClaim', 'TotalClaim'] as const
                ).map((k) => (
                  <td key={k} className="edmt__angka">
                    {sajikan(nilai(h, a + k), SAJIAN_SPREADING.total)}
                  </td>
                ))}
              </tr>
            </tfoot>
          </table>
        </div>
        {tata.totalBerlabel.length > 0 && (
          // PropNewData S11: Total %Share / Total Premium | Total %Share Claim / Total Claim
          <div className="edmt__dua-kolom edmt__total">
            {kolom(tata.totalBerlabel.slice(0, 2))}
            {kolom(tata.totalBerlabel.slice(2))}
          </div>
        )}
      </Wadah>

      {/* `.Installment` + wadah "Installment Data Information" (pyIncludeHeader=true) */}
      <Panel judul={BAGIAN.angsuran}>
        {tata.installment && (
          <div className="edmt__aksi">
            <label className="field__label">{KOLOM_ANGSURAN.installment}</label>
            {sunting ? (
              <input
                className="field__input edmt__pendek"
                aria-label={KOLOM_ANGSURAN.installment}
                value={nilai(h, tata.installment.jalur)}
                onChange={(e) => onUbah(tata.installment!.jalur, e.target.value)}
                onBlur={() => onRefresh(tata.installment!.aksi ?? [])}
              />
            ) : (
              <span>{nilai(h, tata.installment.jalur)}</span>
            )}
          </div>
        )}
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th scope="col">{KOLOM_ANGSURAN.no}</th>
                <th scope="col">{KOLOM_ANGSURAN.dueDate}</th>
                <th scope="col" className="edmt__angka">
                  {KOLOM_ANGSURAN.pct}
                </th>
                <th scope="col" className="edmt__angka">
                  {KOLOM_ANGSURAN.premium}
                </th>
                <th scope="col" className="edmt__angka">
                  {KOLOM_ANGSURAN.total}
                </th>
              </tr>
            </thead>
            <tbody>
              {angsuran.map((b, i) => (
                <tr key={i}>
                  <td>{b.InstallmentNo ?? ''}</td>
                  {/* grid pyEditingMode readOnly: sel ber-aksi SetValidateInstallment_Act / CountPctInstallment_Act
                      tidak pernah terpicu (pola NB) */}
                  <td>{sajikan(b.DueDate ?? '', SAJIAN_ANGSURAN.dueDate)}</td>
                  <td className="edmt__angka">{sajikan(b.InstallmentPercentage ?? '', SAJIAN_ANGSURAN.persen)}</td>
                  <td className="edmt__angka">{sajikan(b.Premium ?? '', SAJIAN_ANGSURAN.premium)}</td>
                  <td className="edmt__angka">{sajikan(b.PaymentTotal ?? '', SAJIAN_ANGSURAN.total)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </Panel>
    </>
  )
}
