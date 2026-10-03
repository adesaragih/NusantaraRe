// Layar Inward Facultative - assignment pertama case NB sesudah Create opportunity (tiket 30, tahap 1).
//
// Port section `InputInwardFacultative` (`D:\migrasi\RNM\NB FacIn\Section\InputInwardFacultative.xml`, flow
// action `InwardFacultative`, flow `InputInwardFacultativeOffer`): urutan tingkat atas = Periode (General) →
// ringkasan (`AllSummarySection`, FIRE: `FireSummarySection`) → checkbox Show Detail → detail
// (`InputInwardFacultativeDtl`, tab) → tombol kaki flow action. Medan yang tampil = kasus FIRE pada tangkapan
// layar work owner 03-10-2026 (gambar tidak disalin: memuat nama pelanggan/orang).
//
// Tahap 1 (keputusan agent, tiket 30):
// - nilai awal dari isian Create opportunity yang baru saja dikirim (belum ada endpoint baca case):
//   Business status, Insured name, Class of business (= Group Business: properti sel 42
//   `.QuotationData.BusinessName`, di gambar "FIRE"), Type facultative; Offering date = hari ini
//   (`InwardFacultative_PreDT`, OfferingDate kosong -> hari ini);
// - medan yang sumbernya belum tersambung (Source of business, Ceding co name, Group Name, Old Policy Number,
//   Risk Scoring, pilihan Marketing Name) tampil kosong - bukan nilai karangan;
// - tombol yang membuka fitur lain (upload, Change SOB/Ceding Co, Search, CSV) dan Submit / Save for later
//   nonaktif sampai backend-nya ada; Cancel kembali ke portal;
// - isi tab detail = tahap 3 (`BelumTersedia`). Nol catatan pengembang di layar.

import { useState } from 'react'

import { BelumTersedia, Field, Pilih, StripTab, type Opsi } from '../../../../inti/frontend/components/ui/dasar'
import type { IsianOpportunity } from '../api'
import TanggalDMY from '../components/TanggalDMY'
import {
  KOLOM_RINGKASAN,
  PERIODE as P,
  PILIHAN_PERIODE,
  SHOW_DETAIL,
  TAB_DETAIL,
  TEKS_FORM_OPPORTUNITY,
  TEKS_INWARD,
  TOMBOL_KAKI_INWARD as KAKI,
} from '../labels'

/** Case NB yang baru dibuat, dibawa dari form Opportunity. */
export interface KasusBaru {
  caseId: string
  isian: IsianOpportunity
  /** INSUREDNAME akun terpilih di ChooseAccount (kosong bila tidak memilih). */
  insuredName: string
}

type TabDetail = (typeof TAB_DETAIL)[number]

/** Hari ini dalam bentuk kabel inti `DD-MM-YYYY` (Offering date awal, `InwardFacultative_PreDT`). */
export function hariIniKabel(d: Date = new Date()): string {
  const p = (n: number) => String(n).padStart(2, '0')
  return `${p(d.getDate())}-${p(d.getMonth() + 1)}-${d.getFullYear()}`
}

const satu = (nilai: string): Opsi[] => (nilai === '' ? [] : [{ value: nilai, label: nilai }])

/** Medan tampil-saja (label kiri, teks kanan) seperti pxDisplayText Pega. */
function Tampil({ label, nilai }: { label: string; nilai: string }) {
  return (
    <div className="field">
      <span className="field__label">{label}</span>
      <div className="nbf-inward__teks">{nilai}</div>
    </div>
  )
}

/** Tombol Pega yang fiturnya belum dibangun: tampil, nonaktif. */
function TombolNonaktif({ label, utama }: { label: string; utama?: boolean }) {
  return (
    <button type="button" className={utama ? 'btn btn--sm' : 'btn btn--ghost btn--sm'} disabled>
      {label}
    </button>
  )
}

/** Radio Pega (pxRadioButtons) - pilihan dari tangkapan layar. */
function Radio({
  label,
  nama,
  pilihan,
  value,
  onChange,
}: {
  label: string
  nama: string
  pilihan: readonly string[]
  value: string
  onChange: (v: string) => void
}) {
  return (
    <div className="field">
      <span className="field__label">{label}</span>
      <div className="nbf-inward__radio" role="radiogroup" aria-label={label}>
        {pilihan.map((p) => (
          <label key={p} className="nbf-inward__pilihan">
            <input type="radio" name={nama} value={p} checked={value === p} onChange={() => onChange(p)} /> {p}
          </label>
        ))}
      </div>
    </div>
  )
}

export default function InwardFacultative({ kasus, onBatal }: { kasus: KasusBaru; onBatal: () => void }) {
  const { isian } = kasus
  const [reff, setReff] = useState('')
  const [qq, setQq] = useState('')
  const [mulai, setMulai] = useState('')
  const [penawaran, setPenawaran] = useState(hariIniKabel())
  const [selesai, setSelesai] = useState('')
  const [policyType, setPolicyType] = useState('')
  const [typeFac, setTypeFac] = useState(isian.typeOfFacultative)
  const [marketing, setMarketing] = useState('')
  const [hari, setHari] = useState('')
  const [detail, setDetail] = useState(false)
  const [tab, setTab] = useState<TabDetail>(TAB_DETAIL[0])

  return (
    <div className="nbfacin">
      <h3 className="nbf-inward__case">{kasus.caseId}</h3>

      <section className="panel">
        <h4 className="panel__title">{P.judul.label}</h4>
        <div className="nbf-opp__kolom">
          <div className="nbf-opp__tumpuk">
            <Field label={P.reffNumber.label} value={reff} onChange={setReff} />
            <Tampil label={P.businessStatus.label} nilai={isian.businessStatus} />
            <Tampil label={P.insuredName.label} nilai={kasus.insuredName} />
            <Field label={P.qqName.label} value={qq} onChange={setQq} />
            <TanggalDMY
              label={P.beginDate.label}
              value={mulai}
              onChange={setMulai}
              required
              labelKalender={TEKS_FORM_OPPORTUNITY.kalender}
              pesanFormat={TEKS_FORM_OPPORTUNITY.formatTanggal}
            />
            <TanggalDMY
              label={P.offeringDate.label}
              value={penawaran}
              onChange={setPenawaran}
              required
              labelKalender={TEKS_FORM_OPPORTUNITY.kalender}
              pesanFormat={TEKS_FORM_OPPORTUNITY.formatTanggal}
            />
            <Radio label={P.policyType.label} nama="nbfacin-policy-type" pilihan={PILIHAN_PERIODE.policyType} value={policyType} onChange={setPolicyType} />
            <Tampil label={P.riskScoring.label} nilai="" />
            <div className="nbf-opp__tombol">
              <TombolNonaktif label={P.uploadQuotation.label} />
            </div>
            <div className="nbf-opp__tombol">
              <TombolNonaktif label={P.uploadRISlip.label} />
            </div>
          </div>
          <div className="nbf-opp__tumpuk">
            <Tampil label={P.classOfBusiness.label} nilai={isian.groupBusiness} />
            <Pilih label={P.typeFacultative.label} value={typeFac} onChange={setTypeFac} opsi={satu(isian.typeOfFacultative)} />
            <div className="nbf-inward__baris-tombol">
              <Tampil label={P.sourceOfBusiness.label} nilai="" />
              <TombolNonaktif label={P.changeSob.label} utama />
            </div>
            <div className="nbf-inward__baris-tombol">
              <Tampil label={P.cedingCoName.label} nilai="" />
              <TombolNonaktif label={P.changeCedingCo.label} utama />
            </div>
            <Tampil label={P.groupName.label} nilai="" />
            <TanggalDMY
              label={P.endDate.label}
              value={selesai}
              onChange={setSelesai}
              required
              labelKalender={TEKS_FORM_OPPORTUNITY.kalender}
              pesanFormat={TEKS_FORM_OPPORTUNITY.formatTanggal}
            />
            <div className="nbf-inward__baris-tombol">
              <span className="field__label">{P.followingPolicyNumber.label}</span>
              <TombolNonaktif label={P.search.label} utama />
            </div>
            <Tampil label={P.oldPolicyNumber.label} nilai="" />
            <Pilih label={P.marketingName.label} value={marketing} onChange={setMarketing} opsi={[]} required />
            <Radio label={P.day.label} nama="nbfacin-day" pilihan={PILIHAN_PERIODE.day} value={hari} onChange={setHari} />
          </div>
        </div>
        <div className="nbf-inward__csv">
          <strong>{P.judulCsv.label}</strong>
          <div className="nbf-opp__tombol">
            <TombolNonaktif label={P.downloadTemplateCsv.label} />
            <TombolNonaktif label={P.uploadCsv.label} utama />
            <TombolNonaktif label={P.viewUpload.label} utama />
            <TombolNonaktif label={P.saveData.label} utama />
            <TombolNonaktif label={P.insertAccumulation.label} />
          </div>
        </div>
      </section>

      <section className="panel">
        <h4 className="panel__title">{TEKS_INWARD.ringkasan}</h4>
        <div className="table-wrap">
          <table className="nbf-tabel">
            <thead>
              <tr>
                {KOLOM_RINGKASAN.map((k) => (
                  <th key={k} scope="col">
                    {k}
                  </th>
                ))}
              </tr>
            </thead>
          </table>
        </div>
        <p className="muted">{TEKS_INWARD.kosong}</p>
        <label className="nbf-inward__pilihan">
          <input type="checkbox" checked={detail} onChange={(e) => setDetail(e.target.checked)} /> {SHOW_DETAIL}
        </label>
        {detail && (
          <div className="nbf-inward__detail">
            <StripTab tab={TAB_DETAIL} aktif={tab} onPilih={setTab} />
            <BelumTersedia apa={`${TEKS_INWARD.isiTab} ${tab}`} />
          </div>
        )}
      </section>

      <div className="nbf-inward__kaki">
        <button type="button" className="btn btn--ghost" onClick={onBatal}>
          {KAKI.batal.label}
        </button>
        <button type="button" className="btn btn--ghost" disabled>
          {KAKI.simpan.label}
        </button>
        <button type="button" className="btn btn--primary" disabled>
          {KAKI.submit.label}
        </button>
      </div>
    </div>
  )
}
