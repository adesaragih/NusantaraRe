// Layar Inward Facultative - assignment pertama case NB sesudah Create opportunity (tiket 30, tahap 1).
//
// Port section `InputInwardFacultative` (`D:\migrasi\RNM\NB FacIn\Section\InputInwardFacultative.xml`, flow
// action `InwardFacultative`, flow `InputInwardFacultativeOffer`): urutan tingkat atas = Periode (General) →
// ringkasan (`AllSummarySection`, FIRE: `FireSummarySection`) → checkbox Show Detail → detail
// (`InputInwardFacultativeDtl`, tab) → tombol kaki flow action. Medan yang tampil = kasus FIRE pada tangkapan
// layar work owner 03-10-2026 (gambar tidak disalin: memuat nama pelanggan/orang). Blok ringkasan (SUMMARY)
// DIHAPUS atas permintaan work owner 03-10-2026 ("summary ini dihapus saja"); checkbox Show Detail tetap.
//
// Tahap 2 (tiket 31): data case dimuat `GET /api/nbfacin/kasus/{caseId}` saat layar dibuka (sampai jawaban
// datang, isian Create opportunity yang dibawa dipakai sebagai nilai awal); pilihan Marketing Name dari
// `GET /api/nbfacin/marketing-officer` (`BrowseMarketingOfficer_RD`); Save for later = `PUT …/general`.
// Submit tetap nonaktif (pasca-proses flow action = tahap tersendiri).
//
// Change SOB (tiket 33): tombol membuka `PopupPilihAgent` (harness `SOB`); Choose mengisi Source of business
// (kode + nama) di layar, kode ikut dikirim Save for later (keputusan agent E-1).
//
// Change Ceding Co (tiket 34): tombol membuka `PopupCedingCoList` (harness `ShowCedingCoList`); Submit
// menyerahkan daftar ceding ke layar, Ceding co name = gabungan nama berpemisah `;` (`SetCedingCo_Act`), dan
// kode-kodenya ikut Save for later (E-7).
//
// Tab Object (tiket 35): kasus FIRE memakai `TabObject` (section `ObjectList`, visible `IsFire`). Kasus FIRE =
// Group Business "FIRE" (keputusan agent G-5 - When `IsFire` Pega lebih luas: IsKPR, IsOilGas, IsFireStyle1/2,
// BusinessType "Fire"; sumber datanya belum diport). Selain FIRE (`InputDtlObject_FacIn`) = tahap berikut.
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

import { useEffect, useState } from 'react'

import { BelumTersedia, Field, Gagal, Pilih, StripTab, type Opsi } from '../../../../inti/frontend/components/ui/dasar'
import {
  ambilKasus,
  daftarMarketing,
  simpanGeneral,
  type BarisCeding,
  type GeneralInward,
  type IsianOpportunity,
  type KasusNB,
} from '../api'
import PopupCedingCoList from '../components/PopupCedingCoList'
import PopupPilihAgent from '../components/PopupPilihAgent'
import TabObject from '../components/TabObject'
import TanggalDMY from '../components/TanggalDMY'
import {
  PERIODE as P,
  PILIHAN_PERIODE,
  POPUP_SOB,
  SHOW_DETAIL,
  TAB_DETAIL,
  TEKS_FORM_OPPORTUNITY,
  TEKS_INWARD,
  TOMBOL_KAKI_INWARD as KAKI,
} from '../labels'

/**
 * Case NB yang dibuka: dari Create opportunity (membawa isian + Insured name sebagai nilai awal) atau dari
 * daftar portal (hanya nomor case - sisanya dimuat `GET /api/nbfacin/kasus/{caseId}`).
 */
export interface KasusBaru {
  caseId: string
  isian?: IsianOpportunity
  /** INSUREDNAME akun terpilih di ChooseAccount (kosong bila tidak memilih). */
  insuredName?: string
}

type TabDetail = (typeof TAB_DETAIL)[number]

/** Hari ini dalam bentuk kabel inti `DD-MM-YYYY` (Offering date awal, `InwardFacultative_PreDT`). */
export function hariIniKabel(d: Date = new Date()): string {
  const p = (n: number) => String(n).padStart(2, '0')
  return `${p(d.getDate())}-${p(d.getMonth() + 1)}-${d.getFullYear()}`
}

/**
 * Begin date + 1 tahun (kabel `DD-MM-YYYY`) - End date otomatis. Permintaan work owner 03-10-2026 ("saat mengisi
 * Begin date, end date nya otomatis ke isi + 1 tahun … tapi masih bisa diganti2"); sejalan `SystemSetOneYear_DT`
 * (`EndDate = @DateTime.addCalendar(StartDate,1,…)`, salinan korpus berkelas PolicyTreatyIn). Seperti Calendar.add
 * YEAR: 29-02 jatuh ke 28-02 bila tahun berikutnya bukan kabisat. Masukan tidak sah -> kosong.
 */
export function tambahSatuTahun(kabel: string): string {
  const m = /^(\d{2})-(\d{2})-(\d{4})$/.exec(kabel)
  if (!m) return ''
  const [h, b, t] = [Number(m[1]), Number(m[2]), Number(m[3]) + 1]
  const akhirBulan = new Date(t, b, 0).getDate()
  const p = (n: number) => String(n).padStart(2, '0')
  return `${p(Math.min(h, akhirBulan))}-${p(b)}-${t}`
}

/** Kabel `DD-MM-YYYY` -> `YYYYMMDD` supaya dapat dibandingkan sebagai teks. */
const kunciTanggal = (kabel: string) => kabel.slice(6) + kabel.slice(3, 5) + kabel.slice(0, 2)

/**
 * End date lebih awal dari Begin date - syarat `SetValidateDate_Act` langkah 1
 * (`.PolicyData.EndDateTime < .PolicyData.StartDateTime && EndDateTime != ""`); tanggal sama = sah.
 */
export function endSebelumBegin(mulai: string, selesai: string): boolean {
  return mulai !== '' && selesai !== '' && kunciTanggal(selesai) < kunciTanggal(mulai)
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
  const isian = kasus.isian
  const [reff, setReff] = useState('')
  const [qq, setQq] = useState('')
  const [mulai, setMulai] = useState('')
  const [penawaran, setPenawaran] = useState(hariIniKabel())
  const [selesai, setSelesai] = useState('')
  const [policyType, setPolicyType] = useState('')
  const [typeFac, setTypeFac] = useState(isian?.typeOfFacultative ?? '')
  const [marketing, setMarketing] = useState('')
  const [hari, setHari] = useState('')
  const [detail, setDetail] = useState(false)
  const [tab, setTab] = useState<TabDetail>(TAB_DETAIL[0])
  const [dariServer, setDariServer] = useState<KasusNB | null>(null)
  const [opsiMarketing, setOpsiMarketing] = useState<Opsi[]>([])
  const [galatMuat, setGalatMuat] = useState<unknown>(null)
  const [menyimpan, setMenyimpan] = useState(false)
  const [galatSimpan, setGalatSimpan] = useState<unknown>(null)
  const [tersimpan, setTersimpan] = useState(false)
  const [sob, setSob] = useState({ id: '', nama: '' })
  const [pilihSob, setPilihSob] = useState(false)
  const [ceding, setCeding] = useState<BarisCeding[]>([])
  const [ubahCeding, setUbahCeding] = useState(false)

  /** Isi state layar dari jawaban server. */
  function terapkan(k: KasusNB) {
    const g = k.general
    setDariServer(k)
    setReff(g.reffNumber)
    setQq(g.qqName)
    setMulai(g.beginDate)
    setPenawaran(g.offeringDate || hariIniKabel())
    setSelesai(g.endDate)
    setPolicyType(g.policyType)
    setTypeFac(g.typeFacultative || k.opportunity.typeOfFacultative)
    setMarketing(g.marketingId)
    setHari(g.day)
    setSob({ id: g.sourceOfBusinessId, nama: g.sourceOfBusiness })
    setCeding(g.cedingList ?? [])
  }

  useEffect(() => {
    let batal = false
    ambilKasus(kasus.caseId).then(
      (k) => {
        if (!batal) terapkan(k)
      },
      (err: unknown) => {
        if (!batal) setGalatMuat(err)
      },
    )
    daftarMarketing().then(
      (h) => {
        if (!batal) setOpsiMarketing(h.baris.map((b) => ({ value: b.id, label: b.nama })))
      },
      (err: unknown) => {
        if (!batal) setGalatMuat(err)
      },
    )
    return () => {
      batal = true
    }
  }, [kasus.caseId])

  const tanggalSalah = endSebelumBegin(mulai, selesai)

  async function simpan() {
    // Peringatan End date < Begin date menahan penyimpanan (keputusan agent F-1).
    if (tanggalSalah) return
    setMenyimpan(true)
    setGalatSimpan(null)
    setTersimpan(false)
    try {
      const k = await simpanGeneral(kasus.caseId, {
        reffNumber: reff,
        qqName: qq,
        beginDate: mulai,
        offeringDate: penawaran,
        endDate: selesai,
        policyType,
        marketingId: marketing,
        day: hari,
        typeFacultative: typeFac,
        sourceOfBusinessId: sob.id,
        cedingIds: ceding.map((c) => c.id),
      })
      terapkan(k)
      setTersimpan(true)
    } catch (err) {
      setGalatSimpan(err)
    } finally {
      setMenyimpan(false)
    }
  }

  /** Medan tampil-saja: dari server bila sudah dimuat. */
  const g: Partial<GeneralInward> = dariServer?.general ?? {}
  const insured = dariServer?.insuredName ?? kasus.insuredName ?? ''
  const op: Partial<IsianOpportunity> = dariServer?.opportunity ?? isian ?? {}
  const kasusFire = (op.groupBusiness ?? '').trim().toUpperCase() === 'FIRE'

  return (
    <div className="nbfacin">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{kasus.caseId}</h2>
      </header>
      <Gagal galat={galatMuat} />
      {tersimpan && <div className="alert alert--ok">{TEKS_INWARD.tersimpan}</div>}
      <Gagal galat={galatSimpan} />

      <section className="panel">
        <h4 className="panel__title">{P.judul.label}</h4>
        <div className="nbf-opp__kolom">
          <div className="nbf-opp__tumpuk">
            <Field label={P.reffNumber.label} value={reff} onChange={setReff} />
            <Tampil label={P.businessStatus.label} nilai={op.businessStatus ?? ''} />
            <Tampil label={P.insuredName.label} nilai={insured} />
            <Field label={P.qqName.label} value={qq} onChange={setQq} />
            <TanggalDMY
              label={P.beginDate.label}
              value={mulai}
              onChange={(k) => {
                setMulai(k)
                // Hanya saat Begin date lengkap dan sah - ketikan setengah jadi tidak menghapus End date.
                if (k !== '') setSelesai(tambahSatuTahun(k))
              }}
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
            <Tampil label={P.classOfBusiness.label} nilai={op.groupBusiness ?? ''} />
            <Pilih label={P.typeFacultative.label} value={typeFac} onChange={setTypeFac} opsi={satu(op.typeOfFacultative ?? '')} />
            <div className="nbf-inward__baris-tombol">
              <Tampil label={P.sourceOfBusiness.label} nilai={sob.nama} />
              <button type="button" className="btn btn--sm" onClick={() => setPilihSob(true)}>
                {P.changeSob.label}
              </button>
            </div>
            <div className="nbf-inward__baris-tombol">
              <Tampil label={P.cedingCoName.label} nilai={ceding.map((c) => c.name).join(';')} />
              <button type="button" className="btn btn--sm" onClick={() => setUbahCeding(true)}>
                {P.changeCedingCo.label}
              </button>
            </div>
            <Tampil label={P.groupName.label} nilai={g.groupName ?? ''} />
            <TanggalDMY
              label={P.endDate.label}
              value={selesai}
              onChange={setSelesai}
              required
              error={tanggalSalah ? TEKS_INWARD.endSebelumBegin : undefined}
              labelKalender={TEKS_FORM_OPPORTUNITY.kalender}
              pesanFormat={TEKS_FORM_OPPORTUNITY.formatTanggal}
            />
            <div className="nbf-inward__baris-tombol">
              <span className="field__label">{P.followingPolicyNumber.label}</span>
              <TombolNonaktif label={P.search.label} utama />
            </div>
            <Tampil label={P.oldPolicyNumber.label} nilai={g.oldPolicyNumber ?? ''} />
            <Pilih label={P.marketingName.label} value={marketing} onChange={setMarketing} opsi={opsiMarketing} kosong={P.marketingKosong.label} required />
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
        <label className="nbf-inward__pilihan">
          <input type="checkbox" checked={detail} onChange={(e) => setDetail(e.target.checked)} /> {SHOW_DETAIL}
        </label>
        {detail && (
          <div className="nbf-inward__detail">
            <StripTab tab={TAB_DETAIL} aktif={tab} onPilih={setTab} />
            {tab === 'Object' && kasusFire ? <TabObject caseId={kasus.caseId} /> : <BelumTersedia apa={`${TEKS_INWARD.isiTab} ${tab}`} />}
          </div>
        )}
      </section>

      <div className="nbf-inward__kaki">
        <button type="button" className="btn btn--ghost" onClick={onBatal}>
          {KAKI.batal.label}
        </button>
        <button type="button" className="btn btn--ghost" onClick={() => void simpan()} disabled={menyimpan || tanggalSalah}>
          {menyimpan ? TEKS_FORM_OPPORTUNITY.menyimpan : KAKI.simpan.label}
        </button>
        <button type="button" className="btn btn--primary" disabled>
          {KAKI.submit.label}
        </button>
      </div>
      {pilihSob && (
        <PopupPilihAgent
          judul={POPUP_SOB.judul}
          onTutup={() => setPilihSob(false)}
          onPilih={(b) => {
            setSob({ id: b.id, nama: b.name })
            setTersimpan(false)
            setPilihSob(false)
          }}
        />
      )}
      {ubahCeding && (
        <PopupCedingCoList
          awal={ceding}
          onTutup={() => setUbahCeding(false)}
          onSubmit={(d) => {
            setCeding(d)
            setTersimpan(false)
            setUbahCeding(false)
          }}
        />
      )}
    </div>
  )
}
