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

import { useEffect, useState, type ReactNode } from 'react'

import { BelumTersedia, Field, Gagal, Pilih, StripTab, type Opsi } from '../../../../inti/frontend/components/ui/dasar'
import {
  AMBANG_HALAMAN_OBJEK,
  ambilKasus,
  daftarMarketing,
  simpanGeneral,
  type BarisCeding,
  type GeneralInward,
  type IsianOpportunity,
  type KasusNB,
} from '../api'
import PopupCedingCoList from '../components/PopupCedingCoList'
import { DetailPolis } from '../components/DetailPolis'
import PopupFollowing from '../components/PopupFollowing'
import PopupPilihAgent from '../components/PopupPilihAgent'
import TabClauses from '../components/TabClauses'
import TabCoverage from '../components/TabCoverage'
import TabSpreading from '../components/TabSpreading'
import TabCedant from '../components/TabCedant'
import TabDeduction from '../components/TabDeduction'
import TabPayment from '../components/TabPayment'
import TabLossRecord from '../components/TabLossRecord'
import TabCorrespondence from '../components/TabCorrespondence'
import BilahCsvObjek, { type JenisCsv } from '../components/BilahCsvObjek'
import BlokLampiran from '../components/BlokLampiran'
import BlokEmailType, { useSubmitMarketing } from '../components/BlokEmailType'
import BacaSajaKasus, { PILIH_UBAH_AKUMULASI } from '../components/BacaSajaKasus'
import TabScoring, { nadaKelas } from '../components/TabScoring'
import TabObject from '../components/TabObject'
import TabObjectHalaman from '../components/TabObjectHalaman'
import TabObjectAneka from '../components/TabObjectAneka'
import TabCoverageCargo from '../components/TabCoverageCargo'
import TabObjectCargo from '../components/TabObjectCargo'
import TabSpreadingCargo from '../components/TabSpreadingCargo'
import TabCoverageAneka from '../components/TabCoverageAneka'
import TabSpreadingAneka from '../components/TabSpreadingAneka'
import TanggalDMY from '../components/TanggalDMY'
import { angka4 } from '../angka4'
import {
  PERIODE as P,
  PILIHAN_PERIODE,
  POPUP_SOB,
  SHOW_DETAIL,
  TAB_DETAIL,
  TEKS_FORM_OPPORTUNITY,
  TEKS_EMAIL_TYPE,
  TEKS_CSV_OBJEK,
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
 * Policy Type tersimpan -> kode standard value (`DDL\PolicyType.xml`). Data lama sistem baru menyimpan LABEL (butir 78.2
 * "Label text for now") - dipetakan ke kodenya supaya radio terpilih dan Save berikutnya menulis kode. Nilai lain apa adanya.
 */
export function kodePolicyType(v: string): string {
  const label: Record<string, string> = { 'Individual Policy': '0', 'Master Policy': '1', 'Declaration Policy': '2' }
  return label[v.trim()] ?? v
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

/** Medan tampil-saja (label kiri, teks kanan) seperti pxDisplayText Pega. */
function Tampil({ label, nilai }: { label: string; nilai: ReactNode }) {
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
  /** Teks (nilai = label) atau pasangan kode + label (Policy Type: standard value). */
  pilihan: readonly (string | { readonly value: string; readonly label: string })[]
  value: string
  onChange: (v: string) => void
}) {
  return (
    <div className="field">
      <span className="field__label">{label}</span>
      <div className="nbf-inward__radio" role="radiogroup" aria-label={label}>
        {pilihan.map((o) => {
          const p = typeof o === 'string' ? { value: o, label: o } : o
          return (
            <label key={p.value} className="nbf-inward__pilihan">
              <input type="radio" name={nama} value={p.value} checked={value === p.value} onChange={() => onChange(p.value)} /> {p.label}
            </label>
          )
        })}
      </div>
    </div>
  )
}

export default function InwardFacultative({
  kasus,
  onBatal,
  onSelesai,
  hanyaLihat = false,
}: {
  kasus: KasusBaru
  onBatal: () => void
  /** Submit berhasil (tiket 57): layar case ditutup; pemanggil menampilkan `pesan` di portal. */
  onSelesai?: (pesan: string) => void
  /**
   * Mode HANYA LIHAT (View Details Following Policy / View Policy, work owner 09-10-2026 "HARUSNYA MUNCUL SEMUA DONG, KOLOM2
   * NYA KAYA YG LAGI DI INPUT. MISAL TAB GENERAL, OBJECT ITEM, SPREADING DAN LAINNYA"): seluruh layar case tampil baca-saja
   * (General + semua tab + riwayat akseptasi + lampiran), Show Detail terbuka, tanpa form Email Type, bilah CSV, dan kaki
   * Cancel / Save / Submit (pemanggil punya tombol Back).
   */
  hanyaLihat?: boolean
}) {
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
  const [skorRisiko, setSkorRisiko] = useState({ finalScore: '', noteFinalScore: '' })
  const [detail, setDetail] = useState(hanyaLihat)
  const [tab, setTab] = useState<TabDetail>(TAB_DETAIL[0])
  // Naik sesudah Save Data CSV objek (tiket 66): tab detail dipasang ulang supaya objek / coverage / spreading terbaca baru.
  const [versiDetail, setVersiDetail] = useState(0)
  const [bukaFollowing, setBukaFollowing] = useState(false)
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
    setPolicyType(kodePolicyType(g.policyType))
    setTypeFac(g.typeFacultative || k.opportunity.typeOfFacultative)
    setMarketing(g.marketingId)
    setHari(g.day)
    setSob({ id: g.sourceOfBusinessId, nama: g.sourceOfBusiness })
    setCeding(g.cedingList ?? [])
    setSkorRisiko(k.riskScoring)
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

  // Submit tahap marketing (tiket 57): sesudah berhasil case DITUTUP dan pengguna kembali ke portal (Pega
  // `finishAssignment`; work owner 07-10-2026 "harusnya saat berhasil di submit, case nya langsung tertutup").
  const sub = useSubmitMarketing(kasus.caseId, (h) => {
    const pesan = [TEKS_EMAIL_TYPE.terkirim, h.nbStatus, h.nomorPolis ? TEKS_EMAIL_TYPE.nomorPolis(h.nomorPolis) : '']
      .filter((x) => x.trim() !== '')
      .join(' ')
    if (onSelesai) onSelesai(pesan)
    else onBatal()
  })

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
  // Layar FIRE: Group Business "FIRE" (perilaku yang sudah berjalan sebelum tiket 73 - TIDAK boleh berubah, work owner 09-10-2026
  // "KALO YANG SUDAH BERHASIL, JANGAN DIUBAH2") ATAU lini backend `fire` (tiket 73 P0: BusinessType_DeT + When IsFire). Lini
  // backend hanya MENAMBAH case FIRE; kunci yang absen (biner lama) / BusinessType "UNKNOWN" tidak lagi menyembunyikan layar FIRE.
  const kasusFire = (op.groupBusiness ?? '').trim().toUpperCase() === 'FIRE' || dariServer?.lini?.fire === true
  // Tab Object keluarga Aneka (tiket 73 P1): gerbang grid lokasi `InputDtlObject_FacIn.xml` L36463 (`IsAnekaNotBonding || IsHE
  // || IsGolfInsurance || IsGrowingTrees`) / L22724 (`IsBonding || IsCustomBonds`); FIRE selalu menang (layar FIRE tidak berubah).
  const lini = dariServer?.lini
  const tampilanAneka = dariServer?.tampilanAneka
  const kasusAneka = !kasusFire && lini !== undefined && tampilanAneka !== undefined && (lini.anekaNotBonding || lini.golf || lini.bonding)
  // Tiket 74 P1 MARINE CARGO: When `IsMarineCargo` (lini backend `marineCargo`); FIRE dan Aneka tetap menang (tidak berubah).
  const kasusCargo = !kasusFire && !kasusAneka && lini?.marineCargo === true
  // Bilah CSV per importer. FIRE / Marine Cargo tetap dari penanda lama; Marine Hull (tiket 73 P6) dari `jenisCsv` server
  // (Pega IsMarineHull). Aneka umum / CPM: template Pega belum dikirim (K73-60) -> tetap pesan "belum tersedia".
  const jenisCsv: JenisCsv | null = kasusFire ? 'fire' : kasusCargo ? 'cargo' : dariServer?.jenisCsv === 'marinehull' ? 'marinehull' : null
  // Submit + blok Email Type: FIRE dan (tiket 73 P4) keluarga Aneka - kabel `…/submit` sama; proteksi per lini di backend.
  // Tiket 74 P5: Marine Cargo juga (proteksi per lini di backend: ship data, Top Risk, coverage, spreading).
  const bolehSubmit = kasusFire || kasusAneka || kasusCargo
  // Tab Loss Record (tiket 55) dan Correspondence (tiket 56) hanya bila backend mengizinkan (POSITION "IT"). Clauses dan Scoring
  // Risk hanya FIRE (`InputInwardFacultativeDtl.xml`: Clauses `IsOilGas || IsFire`, Scoring Risk `IsFire`).
  const tabTampil = TAB_DETAIL.filter(
    (t) =>
      (t !== 'Loss Record' || dariServer?.aksesLossRecord === true) &&
      (t !== 'Correspondence' || dariServer?.aksesKorespondensi === true) &&
      (t !== 'Clauses' || kasusFire) &&
      (t !== 'Scoring Risk' || kasusFire),
  )
  // Baca-saja bila case sudah keluar dari tahap marketing (work owner 07-10-2026); backend menolak tulisan secara mandiri.
  const bacaSaja = hanyaLihat || dariServer?.bolehUbah === false
  // Akun IT (JOB_POSITION "IT") - aturan backend yang sama dengan tab Loss Record / Correspondence (tiket 55 / 56). Tombol
  // Upload QUOTATION/PLACING SLIP AI dan Upload R/I SLIP AI hanya untuk IT (work owner 08-10-2026).
  const akunIT = dariServer?.aksesLossRecord === true
  // Tahap Underwriting (When IsUWSpread, tiket 59): case baca-saja KECUALI tab Spreading.
  const spreadingTerbuka = !hanyaLihat && bacaSaja && dariServer?.bolehUbahSpreading === true
  const detailBacaSaja = bacaSaja && !(spreadingTerbuka && tab === 'Spreading')
  // Tahap Underwriting (tiket 60): Accumulation di tab Coverage terbuka - pembungkus tetap aktif, kecuali kontrol Accumulation.
  const akumulasiTerbuka = !hanyaLihat && bacaSaja && dariServer?.ubahAkumulasi === true
  // Tiket 64 K1: Begin / End date terkunci di tahap Binding (Pega Read-only) walau field General lain masih bisa diubah.
  const kunciPeriode = dariServer?.kunciPeriode === true
  const kecualiDetail = akumulasiTerbuka && tab === 'Coverage' ? PILIH_UBAH_AKUMULASI : undefined

  return (
    <div className="nbfacin">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{kasus.caseId}</h2>
      </header>
      <Gagal galat={galatMuat} />
      {tersimpan && <div className="alert alert--ok">{TEKS_INWARD.tersimpan}</div>}
      <Gagal galat={galatSimpan} />
      {/* Mode hanya lihat: tanpa banner baca-saja (work owner 09-10-2026 "pesan yang diatas ini tidak perlu"). */}
      {bacaSaja && !hanyaLihat && (
        <div className="alert nbf-inward__bacasaja">
          {akumulasiTerbuka
            ? TEKS_INWARD.bacaSajaKecualiSpreadingAkumulasi
            : spreadingTerbuka
              ? TEKS_INWARD.bacaSajaKecualiSpreading
              : TEKS_INWARD.bacaSaja}
          {!spreadingTerbuka && !akumulasiTerbuka && dariServer?.alasanBacaSaja ? ` ${dariServer.alasanBacaSaja}` : ''}
        </div>
      )}

      <BacaSajaKasus aktif={bacaSaja} sembunyi={hanyaLihat}>
      <section className="panel nbf-inward">
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
              disabled={kunciPeriode}
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
            <Tampil
              label={P.riskScoring.label}
              nilai={
                skorRisiko.finalScore === '' ? (
                  ''
                ) : (
                  <>
                    {angka4(skorRisiko.finalScore)}{' '}
                    {skorRisiko.noteFinalScore !== '' && (
                      <span className="nbf-score__pil nbf-score__pil--umum" data-nada={nadaKelas(skorRisiko.noteFinalScore)}>
                        {skorRisiko.noteFinalScore}
                      </span>
                    )}
                  </>
                )
              }
            />
            {akunIT && (
              <>
                <div className="nbf-opp__tombol">
                  <TombolNonaktif label={P.uploadQuotation.label} />
                </div>
                <div className="nbf-opp__tombol">
                  <TombolNonaktif label={P.uploadRISlip.label} />
                </div>
              </>
            )}
          </div>
          <div className="nbf-opp__tumpuk">
            <Tampil label={P.classOfBusiness.label} nilai={op.groupBusiness ?? ''} />
            {/* Baca-saja: sudah dipilih di Create opportunity (work owner 05-10-2026). Nilainya tetap ikut Save for later. */}
            <Tampil label={P.typeFacultative.label} nilai={typeFac} />
            <div className="nbf-inward__baris-tombol">
              <Tampil label={P.sourceOfBusiness.label} nilai={sob.nama} />
              <button type="button" className="btn nbf-tema__tombol btn--sm" onClick={() => setPilihSob(true)}>
                {P.changeSob.label}
              </button>
            </div>
            <div className="nbf-inward__baris-tombol">
              <Tampil label={P.cedingCoName.label} nilai={ceding.map((c) => c.name).join(';')} />
              <button type="button" className="btn nbf-tema__tombol btn--sm" onClick={() => setUbahCeding(true)}>
                {P.changeCedingCo.label}
              </button>
            </div>
            <Tampil label={P.groupName.label} nilai={g.groupName ?? ''} />
            <TanggalDMY
              label={P.endDate.label}
              value={selesai}
              disabled={kunciPeriode}
              onChange={setSelesai}
              required
              error={tanggalSalah ? TEKS_INWARD.endSebelumBegin : undefined}
              labelKalender={TEKS_FORM_OPPORTUNITY.kalender}
              pesanFormat={TEKS_FORM_OPPORTUNITY.formatTanggal}
            />
            <div className="nbf-inward__baris-tombol">
              <span className="field__label">{P.followingPolicyNumber.label}</span>
              {/* Tiket 70: popup ViewFollowingNB (cari polis, Copy Nopolis / Copy All, View Details). */}
              <button
                type="button"
                className="btn nbf-tema__tombol nbf-tema__tombol--utama btn--primary btn--sm"
                disabled={hanyaLihat || dariServer?.bolehFollowing !== true}
                onClick={() => setBukaFollowing(true)}
              >
                {P.search.label}
              </button>
            </div>
            <Tampil label={P.oldPolicyNumber.label} nilai={g.oldPolicyNumber ?? ''} />
            <Pilih label={P.marketingName.label} value={marketing} onChange={setMarketing} opsi={opsiMarketing} kosong={P.marketingKosong.label} required />
            <Radio label={P.day.label} nama="nbfacin-day" pilihan={PILIHAN_PERIODE.day} value={hari} onChange={setHari} />
          </div>
        </div>
        {/* Bilah CSV (work owner 08-10-2026): empat tombol CSV tidak tampil di tahap baca-saja (UW dst.); Insert Accumulation
            hanya akun IT. Bilah hilang bila tidak ada tombol. */}
        {!hanyaLihat && (!bacaSaja || akunIT) && (
          <div className="nbf-inward__csv nbf-inward__bilah">
            <strong>{P.judulCsv.label}</strong>
            {/* Bug DEV 11-10-2026: CSV FIRE terunggah ke case MARINE CARGO. Pega `ChooseActivityMappingCSV_Act` menjalankan
                InsertUploadFire_act HANYA bila IsFire (L439/L499); Marine / Aneka / Briguna punya impor sendiri (belum diport). */}
            {!bacaSaja && jenisCsv === null ? (
              <div className="alert alert--info">{TEKS_CSV_OBJEK.belumUntukCob(op.groupBusiness ?? '')}</div>
            ) : !bacaSaja && jenisCsv !== null ? (
              <BilahCsvObjek
                // Satu bilah per case: status unggahan case lain tidak boleh tertinggal (bug DEV 10-10-2026, tiket 75).
                key={kasus.caseId}
                caseId={kasus.caseId}
                // Tiket 74 P6 / 73 P6: importer menurut COB (Pega ChooseActivityMappingCSV_Act).
                jenis={jenisCsv}
                akhir={akunIT && <TombolNonaktif label={P.insertAccumulation.label} />}
                onTersimpan={() => setVersiDetail((n) => n + 1)}
              />
            ) : (
              <div className="nbf-opp__tombol">
                <TombolNonaktif label={P.insertAccumulation.label} />
              </div>
            )}
          </div>
        )}
      </section>
      </BacaSajaKasus>

      {/* Pemisah General <-> Show Detail: kotak kosong berlatar (work owner 05-10-2026). */}
      <div className="nbf-inward__pemisah" aria-hidden="true" />

      <section className="panel nbf-inward__rinci">
        <label className="nbf-inward__pilihan nbf-inward__sakelar">
          <input type="checkbox" checked={detail} onChange={(e) => setDetail(e.target.checked)} /> {SHOW_DETAIL}
        </label>
        {detail && (
          <BacaSajaKasus key={`${detailBacaSaja ? `bacasaja${kecualiDetail ?? ''}` : 'ubah'}-${versiDetail}`} aktif={detailBacaSaja} kecuali={kecualiDetail} sembunyi={hanyaLihat}>
          <div className="nbf-inward__detail">
            <StripTab tab={tabTampil} aktif={tab} onPilih={setTab} />
            {tab === 'Object' && kasusFire ? (
              // Tiket 75 F3 (K75-3): case besar per halaman dari server; case biasa tetap TabObject (alur lama, tidak berubah).
              (dariServer?.jumlahLokasi ?? 0) > AMBANG_HALAMAN_OBJEK ? (
                <TabObjectHalaman caseId={kasus.caseId} insuredName={insured} />
              ) : (
                <TabObject caseId={kasus.caseId} insuredName={insured} />
              )
            ) : tab === 'Coverage' && kasusFire ? (
              <TabCoverage caseId={kasus.caseId} hanyaAkumulasi={akumulasiTerbuka} />
            ) : tab === 'Clauses' && kasusFire ? (
              <TabClauses caseId={kasus.caseId} />
            ) : tab === 'Spreading' && kasusFire ? (
              <TabSpreading caseId={kasus.caseId} tahapUw={dariServer?.ubahBarisSpreading === true} />
            ) : tab === 'Inw Fac Cedant Panels' && kasusFire ? (
              <TabCedant caseId={kasus.caseId} />
            ) : tab === 'Premium Deduction / Brokerage Fee' && kasusFire ? (
              <TabDeduction caseId={kasus.caseId} />
            ) : tab === 'Correspondence' && dariServer?.aksesKorespondensi === true ? (
              <TabCorrespondence caseId={kasus.caseId} />
            ) : tab === 'Loss Record' && kasusFire && dariServer?.aksesLossRecord === true ? (
              <TabLossRecord caseId={kasus.caseId} insuredName={insured} />
            ) : tab === 'Payment' && kasusFire ? (
              <TabPayment caseId={kasus.caseId} />
            ) : tab === 'Scoring Risk' && kasusFire ? (
              <TabScoring caseId={kasus.caseId} saatTersimpan={setSkorRisiko} />
            ) : tab === 'Object' && kasusAneka && tampilanAneka !== undefined ? (
              <TabObjectAneka caseId={kasus.caseId} tampilan={tampilanAneka} aneka={lini?.aneka === true} />
            ) : tab === 'Coverage' && kasusAneka && tampilanAneka !== undefined ? (
              <TabCoverageAneka caseId={kasus.caseId} tampilan={tampilanAneka} aneka={lini?.aneka === true} />
            ) : tab === 'Spreading' && kasusAneka ? (
              <TabSpreadingAneka caseId={kasus.caseId} tahapUw={dariServer?.ubahBarisSpreading === true} />
            ) : tab === 'Inw Fac Cedant Panels' && kasusAneka ? (
              <TabCedant caseId={kasus.caseId} aneka />
            ) : tab === 'Premium Deduction / Brokerage Fee' && kasusAneka ? (
              <TabDeduction caseId={kasus.caseId} aneka />
            ) : tab === 'Payment' && kasusAneka ? (
              <TabPayment caseId={kasus.caseId} />
            ) : tab === 'Object' && kasusCargo ? (
              <TabObjectCargo caseId={kasus.caseId} />
            ) : tab === 'Coverage' && kasusCargo ? (
              <TabCoverageCargo caseId={kasus.caseId} />
            ) : tab === 'Spreading' && kasusCargo ? (
              <TabSpreadingCargo caseId={kasus.caseId} tahapUw={dariServer?.ubahBarisSpreading === true} />
            ) : tab === 'Inw Fac Cedant Panels' && kasusCargo ? (
              <TabCedant caseId={kasus.caseId} cargo />
            ) : tab === 'Premium Deduction / Brokerage Fee' && kasusCargo ? (
              <TabDeduction caseId={kasus.caseId} cargo />
            ) : tab === 'Payment' && kasusCargo ? (
              <TabPayment caseId={kasus.caseId} />
            ) : (
              dariServer !== null && !kasusFire ? (
                <div className="alert alert--info">{TEKS_INWARD.cobBelum(tab, op.groupBusiness ?? '')}</div>
              ) : (
                <BelumTersedia apa={`${TEKS_INWARD.isiTab} ${tab}`} />
              )
            )}
          </div>
          </BacaSajaKasus>
        )}
      </section>

      {bolehSubmit && <BlokEmailType s={sub} hanyaRiwayat={hanyaLihat} />}
      {/* Attachment File (tiket 67): di bawah riwayat akseptasi, semua tahap (bukan bagian pembungkus baca-saja). Mode hanya
          lihat: unggah / hapus terkunci, View File (data-navigasi) tetap. */}
      <BacaSajaKasus aktif={hanyaLihat} sembunyi>
        <BlokLampiran caseId={kasus.caseId} />
      </BacaSajaKasus>

      {!hanyaLihat && (
      <div className="nbf-inward__kaki">
        <button type="button" className="btn nbf-tema__tombol btn--ghost" onClick={onBatal}>
          {KAKI.batal.label}
        </button>
        <button type="button" className="btn nbf-tema__tombol btn--ghost" onClick={() => void simpan()} disabled={menyimpan || tanggalSalah || bacaSaja}>
          {menyimpan ? TEKS_FORM_OPPORTUNITY.menyimpan : KAKI.simpan.label}
        </button>
        <button type="button" className="btn nbf-tema__tombol nbf-tema__tombol--utama btn--primary" disabled={!bolehSubmit || !sub.bisaSubmit} title={sub.galatIsi || undefined} onClick={sub.submit}>
          {sub.mengajukan ? TEKS_EMAIL_TYPE.mengajukan : KAKI.submit.label}
        </button>
      </div>
      )}
      {bukaFollowing && (
        <PopupFollowing
          caseId={kasus.caseId}
          tampilDetail={(b, kembali) => (
            <DetailPolis
              baris={b}
              judul=""
              onKembali={kembali}
              tampilKasus={(caseId) => <InwardFacultative kasus={{ caseId }} onBatal={kembali} hanyaLihat />}
            />
          )}
          onTutup={() => setBukaFollowing(false)}
          onTersalin={(k) => {
            terapkan(k)
            setVersiDetail((n) => n + 1)
          }}
        />
      )}
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
