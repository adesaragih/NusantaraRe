// Sub-tab Loss Record dan Loss Record Internal baris objek FIRE (tiket 42).
//
// Loss Record = `NB FacIn\Section\CauseOfLoss_FacIn.xml` (grid `.Property.ListCauseOfLoss`: Date of Loss · Insured
// Name · Loss Object · Currency · Total of Loss · Total Claim; Add / Delete; baris dibuka = form
// `InputCauseOfLoss_FacIn`) + `InputOfferFacInLossRatio.xml` (LR 1 Year, %LR 1 Years, LR 3 - 5 Years, %LR 3 - 5 Years,
// baca-saja). Loss Record Internal = `CauseOfLossClaim_FacIn.xml` (grid baca-saja `.Property.ListCauseOfLossClaim`).
//
// Uang (Total of Loss, Total Claim, Prevention Of Loss, Premium, klaim) = teks desimal, tampil lewat `formatNumber`
// (4 desimal, pxNumber Pega) - tanpa float (ADR-0003).
//
// Keputusan agent (tiket 42): N-1 rumus Loss Ratio TIDAK diport - `SetLossRatio_Act` versi baris objek tidak ada di
// korpus; nilai LR tampil apa adanya dari server (data contoh: semuanya 0). N-2 Total of Loss (`.Amount`) baca-saja
// (Pega tanpa isian). N-3 Insured Name grid (`.CoinsData.CoinsName`) diisi nama tertanggung case saat Add `[dugaan]`
// (Pega: pengisinya tidak ada di korpus; data contoh selalu terisi). N-4 Loss Detail bertanda wajib tanpa menahan
// Save (pola L-2). N-8 label / daftar Remarks = aturan properti `DDL\Remarks.xml`. N-5 Loss Record Internal = grid
// baca-saja dari server (dikonfirmasi work owner 03-10-2026: "biarkan saja kosong"); di Pega pengisinya
// (`MappingKlaimToLossRecord_Act`) tidak pernah terpanggil (tombol tersembunyi, pemanggil dikomentari) dan memakai
// NOPOLIS ter-hardcode - pengisian otomatis menunggu keputusan work owner. N-6 Total Claim awal "0" (nilai awal sel 7).
// N-9 Currency wajib per catatan (pola A133: kolom rancangan T_LISTCAUSEOFLOSS.CURRENCY NOT NULL DEFAULT 'UNKNOWN') -
// bertanda wajib, menahan Save, pesan tampil sesudah Save dicoba.

import { useEffect, useState } from 'react'

import { Area, Field, Pilih, type Opsi } from '../../../../inti/frontend/components/ui/dasar'
import { desimalSah } from '../../../../inti/frontend/lib/desimal'
import { formatNumber } from '../../../../inti/frontend/lib/format'
import { daftarMataUang, type CatatanKerugian, type KlaimInternal, type LossRatio } from '../api'
import IsianUang from './IsianUang'
import TanggalDMY from './TanggalDMY'
import {
  FORM_KERUGIAN as F,
  GRID_KERUGIAN,
  GRID_KLAIM_INTERNAL,
  GRID_OBJEK,
  ITEM_KOSONG,
  LABEL_REMARKS,
  LOSS_RATIO,
  OPSI_REMARKS,
  TEKS_FORM_OPPORTUNITY,
  TEKS_INWARD,
  TEKS_ITEM,
  TEKS_KERUGIAN,
  TEKS_OBJEK,
} from '../labels'

const DESIMAL_UANG = 4

/** Catatan kerugian baru (Add): Insured Name = tertanggung case (N-3), Total Claim "0" (N-6). */
export function kerugianBaru(insuredName: string): CatatanKerugian {
  return {
    dateOfLoss: '', coinsName: insuredName, lossObject: '', currency: '', amount: '', claim: '0', preventionOfLoss: '',
    causeOfLoss: '', remarks: '', detail: '',
  }
}

/** Loss ratio kosong (objek baru / data lama). */
export const lossRatioKosong = (): LossRatio => ({ oneYearAmount: '', oneYearPercent: '', threeFiveYearAmount: '', threeFiveYearPercent: '' })

/** Uang sah: kosong atau desimal bertitik. */
export const uangSah = (v: string) => v.trim() === '' || desimalSah(v.trim())

/** Ada catatan bergalat: Currency kosong (N-9) atau uang tidak sah. */
export function adaGalatKerugian(r: CatatanKerugian[]): boolean {
  return r.some((x) => x.currency.trim() === '' || !uangSah(x.claim) || !uangSah(x.preventionOfLoss))
}

function Tampil({ label, nilai }: { label: string; nilai: string }) {
  return (
    <div className="field">
      <span className="field__label">{label}</span>
      <div className="nbf-inward__teks">{nilai}</div>
    </div>
  )
}

function FormKerugian({
  r,
  ubah,
  insuredName,
  mataUang,
  tandaiWajib,
}: {
  r: CatatanKerugian
  ubah: (r: CatatanKerugian) => void
  insuredName: string
  mataUang: string[]
  tandaiWajib: boolean
}) {
  const set = (k: keyof CatatanKerugian) => (v: string) => ubah({ ...r, [k]: v })
  const opsiUang: Opsi[] = mataUang.map((m) => ({ value: m, label: m }))
  return (
    <div className="nbf-objek__isi nbf-ringkas nbf-rugi-form">
      {/* Tata letak = gambar layar Pega (work owner 05-10-2026): dua kolom - Insured Name | Date of Loss, Loss Object |
          Currency, Total Claim (100%) | Prevention Of Loss, Cause of Loss | Remarks; Loss Detail selebar form. */}
      <div className="nbf-rugi-form__grid">
        <Tampil label={F.insuredName.label} nilai={insuredName} />
        <TanggalDMY
          label={F.dateOfLoss.label}
          value={r.dateOfLoss}
          onChange={set('dateOfLoss')}
          labelKalender={TEKS_FORM_OPPORTUNITY.kalender}
          pesanFormat={TEKS_FORM_OPPORTUNITY.formatTanggal}
        />
        <Field label={F.lossObject.label} value={r.lossObject} onChange={set('lossObject')} />
        <Pilih
          label={F.currency.label}
          value={r.currency}
          onChange={set('currency')}
          opsi={opsiUang}
          kosong={ITEM_KOSONG}
          required
          error={tandaiWajib && r.currency.trim() === '' ? TEKS_ITEM.currencyWajib : undefined}
        />
        <IsianUang label={F.claim.label} value={r.claim} onChange={set('claim')} error={uangSah(r.claim) ? undefined : TEKS_KERUGIAN.uang} />
        <IsianUang
          label={F.preventionOfLoss.label}
          value={r.preventionOfLoss}
          onChange={set('preventionOfLoss')}
          error={uangSah(r.preventionOfLoss) ? undefined : TEKS_KERUGIAN.uang}
        />
        <Field label={F.causeOfLoss.label} value={r.causeOfLoss} onChange={set('causeOfLoss')} />
        <Pilih label={LABEL_REMARKS} value={r.remarks} onChange={set('remarks')} opsi={OPSI_REMARKS} />
        <div className="nbf-rugi-form__lebar">
          <Area label={F.detail.label} value={r.detail} onChange={set('detail')} required baris={3} />
        </div>
      </div>
    </div>
  )
}

export function SubTabKerugian({
  rows,
  lossRatio,
  insuredName,
  ubah,
  tandaiWajib = false,
}: {
  rows: CatatanKerugian[]
  lossRatio: LossRatio
  insuredName: string
  ubah: (rows: CatatanKerugian[]) => void
  /** Save tab Object sudah dicoba. */
  tandaiWajib?: boolean
}) {
  const [terbuka, setTerbuka] = useState<number[]>([])
  const [mataUang, setMataUang] = useState<string[]>([])

  useEffect(() => {
    let batal = false
    daftarMataUang().then(
      (h) => {
        if (!batal) setMataUang(h.baris)
      },
      () => {},
    )
    return () => {
      batal = true
    }
  }, [])

  return (
    <div className="nbf-objek__isi">
      <div className="nbf-cov-wrap">
        <table className="nbf-tabel">
          <thead>
            <tr>
              <th scope="col" />
              {GRID_KERUGIAN.map((k) => (
                <th key={k.sel} scope="col">
                  {k.label}
                </th>
              ))}
              <th scope="col" className="table__actions">
                <button
                  type="button"
                  className="btn btn--ghost btn--sm"
                  onClick={() => {
                    ubah([...rows, kerugianBaru(insuredName)])
                    setTerbuka((t) => [...t, rows.length])
                  }}
                >
                  {GRID_OBJEK.tambah}
                </button>
              </th>
            </tr>
          </thead>
          <tbody>
            {rows.length === 0 && (
              <tr>
                <td colSpan={8}>{TEKS_INWARD.kosong}</td>
              </tr>
            )}
            {rows.map((r, n) => [
              <tr key={`b-${n}`}>
                <td>
                  <button
                    type="button"
                    className="btn btn--ghost btn--sm"
                    aria-label={TEKS_OBJEK.bukaBaris}
                    aria-expanded={terbuka.includes(n)}
                    onClick={() => setTerbuka((t) => (t.includes(n) ? t.filter((x) => x !== n) : [...t, n]))}
                  >
                    {terbuka.includes(n) ? '▾' : '▸'}
                  </button>
                </td>
                <td>{r.dateOfLoss.replace(/-/g, '/')}</td>
                <td>{r.coinsName}</td>
                <td>{r.lossObject}</td>
                <td>{r.currency}</td>
                <td className="nbf-angka">{formatNumber(r.amount, DESIMAL_UANG)}</td>
                <td className="nbf-angka">{formatNumber(r.claim, DESIMAL_UANG)}</td>
                <td className="table__actions">
                  <button
                    type="button"
                    className="btn btn--ghost btn--sm"
                    onClick={() => {
                      ubah(rows.filter((_, x) => x !== n))
                      setTerbuka((t) => t.filter((x) => x !== n).map((x) => (x > n ? x - 1 : x)))
                    }}
                  >
                    {GRID_OBJEK.hapus}
                  </button>
                </td>
              </tr>,
              terbuka.includes(n) && (
                <tr key={`d-${n}`} className="nbf-objek__detail">
                  <td colSpan={8}>
                    <FormKerugian
                      r={r}
                      ubah={(baru) => ubah(rows.map((x, k) => (k === n ? baru : x)))}
                      insuredName={insuredName}
                      mataUang={mataUang}
                      tandaiWajib={tandaiWajib}
                    />
                  </td>
                </tr>
              ),
            ])}
          </tbody>
        </table>
      </div>
      <div className="nbf-ringkas nbf-ringkas__grid nbf-ringkas__grid--empat">
        {LOSS_RATIO.map((l) => (
          <Tampil key={l.sel} label={l.label} nilai={formatNumber(lossRatio[l.kunci], l.desimal)} />
        ))}
      </div>
    </div>
  )
}

export function SubTabKlaimInternal({ rows }: { rows: KlaimInternal[] }) {
  return (
    <div className="nbf-objek__isi">
      <div className="nbf-cov-wrap">
        <table className="nbf-tabel">
          <thead>
            <tr>
              {GRID_KLAIM_INTERNAL.map((k) => (
                <th key={k.sel} scope="col">
                  {k.label}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {rows.length === 0 && (
              <tr>
                <td colSpan={GRID_KLAIM_INTERNAL.length}>{TEKS_INWARD.kosong}</td>
              </tr>
            )}
            {rows.map((r, n) => (
              <tr key={n}>
                <td>{r.dateOfLoss.slice(6)}</td>
                <td>{r.dateOfLoss.replace(/-/g, '/')}</td>
                <td>{r.locationNo}</td>
                <td>{r.location}</td>
                <td>{r.currency}</td>
                <td className="nbf-angka">{formatNumber(r.premium, DESIMAL_UANG)}</td>
                <td className="nbf-angka">{formatNumber(r.osClaim, DESIMAL_UANG)}</td>
                <td className="nbf-angka">{formatNumber(r.acceptedClaim, DESIMAL_UANG)}</td>
                <td className="nbf-angka">{formatNumber(r.incurredClaim, DESIMAL_UANG)}</td>
                <td className="nbf-angka">{formatNumber(r.lossRatio, DESIMAL_UANG)}</td>
                <td>{r.remark}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
