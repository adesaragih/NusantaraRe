// Form Create Addendum - irisan I1-I3 (keputusan work owner 11-10-2026, W-EDM-01…14).
//
// Port `Section/WorkPrimaryDetails.xml` (kelas `ASM-SFAGIS-Work-Endorsement`) dan tombolnya
// `Section/crmNewHarnessButtons.xml`. Tata letak = tangkapan layar Pega dari work owner (11-10-2026): dua kolom, label di
// KIRI medan, tombol OK lalu Cancel di kanan bawah di luar kotak form. Tema = templat NB FacIn ("TEMPLATE NYA JUGA SAMAIN
// KAYA NB FACIN"; pola `modul/nbfacin/frontend/pages/FormOpportunity.tsx` disalin, tidak diimpor).
//
//   kiri:  Policy Number / Cover Number (`.PolicyNo`, teks bebas) · Endorsement Source (`.EndorsementSourceNote`) ·
//          Endorsement Note (`.Note`) · Endorsement Type (`.Quotation.EdmType`) · EDM Type (`.Quotation.EdmTypeNew`,
//          hanya EdmType 4) · Detail EDM Type (`.Quotation.Type`, `EdmType1/2/3`: EdmType 4, bukan Life, EdmTypeNew 1/2/3)
//   kanan: Endorsement Date (`.EndorsementDate`, dd/mm/yyyy) · Endorsement Status (`.EndorsementStatus`, bukan Life,
//          bawaan Internal) · Endorsement Internal Retro (`.EndorsementInternalRetro`, bukan Life, radio vertikal,
//          bawaan No)
//
// ⛔ `[terverifikasi]` `.AdminCharge` ("Biaya") dan `.Survey` TIDAK dirender: kontainernya `pyContainerVisibleWhen =
// never` (WorkPrimaryDetails L9174) - di Pega tidak pernah tampil.
//
// Pemeriksaan (I2): di Pega `SetErrorBatalEndorsement_Act` berjalan saat PolicyNo / Endorsement Type / EDM Type /
// Internal Retro berubah, `CheckEDMPolisDate` saat Endorsement Date berubah; tombol OK aktif hanya bila keduanya lolos
// (`InputData.CARIGROUP=='true' && InputSearch.CARIGROUP=='true'`, crmNewHarnessButtons L885). Di sini keduanya satu
// panggilan `POST /api/endorsmentfacin/periksa` - NOL tulis.
//
// OK (irisan I3): `POST /api/endorsmentfacin/kasus` = `SetValueToEDMWork`. Aktif bila pemeriksaan terakhir lolos dan
// lengkap serta medan wajib terisi; server memeriksa ulang. Ditolak di server = jawaban 200 tanpa `caseId` (hasil bisnis,
// spec 07 §7.2) - pesannya tampil di bawah medannya. Sumber pilihan Endorsement Source
// (`D_EnumerationList` `Type=sumberedm`) belum ada di korpus - sementara teks bebas, sama dengan autocomplete Pega yang
// mengizinkan isian bebas.

import { useEffect, useRef, useState } from 'react'

import { Area, Field, Gagal, Pilih } from '../../../../inti/frontend/components/ui/dasar'
import { buatKasus, periksaAddendum, type HasilPeriksa } from '../api'
import TanggalDMY from '../components/TanggalDMY'
import { DETAIL_EDM, FORM, JENIS_EDM, JENIS_EDM_BARU, PERIKSA, PILIHAN_RETRO, STATUS_EDM, type Pilihan } from '../labels'

/** Jeda sebelum pemeriksaan dijalankan sesudah isian berubah (ms). */
const JEDA_PERIKSA = 400

const opsi = (d: readonly Pilihan[]) => d.map((p) => ({ value: p.nilai, label: p.label }))

/** `onDibuat` - kasus `EDM-` lahir; pemanggil kembali ke portal dengan pesan. */
export default function FormAddendum({ onBatal, onDibuat }: { onBatal: () => void; onDibuat: (caseId: string) => void }) {
  const [policyNo, setPolicyNo] = useState('')
  const [sumber, setSumber] = useState('')
  const [catatan, setCatatan] = useState('')
  const [edmType, setEdmType] = useState('')
  const [edmTypeNew, setEdmTypeNew] = useState('')
  const [detail, setDetail] = useState('')
  const [tanggal, setTanggal] = useState('')
  // Bawaan tampil "Internal" (tangkapan layar Pega, W-EDM-13).
  const [status, setStatus] = useState('0')
  // Bawaan `0` = No - WorkPrimaryDetails L10212.
  const [retro, setRetro] = useState('0')

  const [hasil, setHasil] = useState<HasilPeriksa | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [memeriksa, setMemeriksa] = useState(false)
  const [menyimpan, setMenyimpan] = useState(false)
  const nomorPermintaan = useRef(0)

  const polis = policyNo.trim()
  useEffect(() => {
    const nomor = ++nomorPermintaan.current
    if (polis === '') {
      setHasil(null)
      setGalat(null)
      setMemeriksa(false)
      return
    }
    const t = setTimeout(() => {
      setMemeriksa(true)
      setGalat(null)
      periksaAddendum({
        policyNo: polis,
        edmType,
        edmTypeNew: edmType === '4' ? edmTypeNew : '',
        endorsementDate: tanggal,
        endorsementInternalRetro: retro,
      })
        .then((h) => {
          if (nomor === nomorPermintaan.current) setHasil(h)
        })
        .catch((err: unknown) => {
          if (nomor === nomorPermintaan.current) {
            setHasil(null)
            setGalat(err)
          }
        })
        .finally(() => {
          if (nomor === nomorPermintaan.current) setMemeriksa(false)
        })
    }, JEDA_PERIKSA)
    return () => clearTimeout(t)
  }, [polis, edmType, edmTypeNew, tanggal, retro])

  // `OutData.pxResults(1).CARI2!='Life'` - sebelum diperiksa CARI2 kosong, jadi medannya tampil (seperti Pega).
  const life = hasil?.bisnis.jenisBisnis === 'Life'
  const daftarDetail = edmType === '4' && !life ? DETAIL_EDM[edmTypeNew] : undefined
  const pesanMedan = (medan: string) =>
    hasil?.pesan
      .filter((p) => p.medan === medan)
      .map((p) => p.teks)
      .join(' · ') || undefined

  // Medan wajib (`pyRequired` WorkPrimaryDetails) yang masih kosong, urut layar.
  const kosong = (
    [
      [FORM.policyNo, polis],
      [FORM.sumber, sumber.trim()],
      [FORM.catatan, catatan.trim()],
      [FORM.jenis, edmType],
      ...(edmType === '4' ? ([[FORM.jenisBaru, edmTypeNew]] as [string, string][]) : []),
      ...(daftarDetail ? ([[FORM.detail, detail]] as [string, string][]) : []),
      [FORM.tanggal, tanggal],
      ...(!life ? ([[FORM.status, status], [FORM.retro, retro]] as [string, string][]) : []),
    ] as [string, string][]
  )
    .filter(([, v]) => v === '')
    .map(([l]) => l)

  const bolehOK = kosong.length === 0 && !memeriksa && !menyimpan && hasil !== null && !hasil.ditolak && hasil.lengkap && hasil.tanggal !== null
  // Alasan OK mati - satu kalimat, prioritas: isian kosong, sedang memeriksa, ditolak, belum lengkap.
  const alasan = bolehOK
    ? ''
    : kosong.length > 0
      ? PERIKSA.alasan.wajib(kosong)
      : memeriksa || hasil === null
        ? PERIKSA.alasan.memeriksa
        : hasil.ditolak
          ? PERIKSA.alasan.ditolak
          : !hasil.lengkap
            ? PERIKSA.alasan.belumLengkap
            : ''

  async function simpan() {
    setMenyimpan(true)
    setGalat(null)
    try {
      const h = await buatKasus({
        policyNo: polis,
        edmType,
        edmTypeNew: edmType === '4' ? edmTypeNew : '',
        endorsementDate: tanggal,
        endorsementInternalRetro: life ? '' : retro,
        note: catatan,
        endorsementSource: '',
        endorsementSourceNote: sumber.trim(),
        endorsementStatus: life ? '' : status,
        detailType: daftarDetail ? detail : '',
      })
      if (h.caseId !== '') onDibuat(h.caseId)
      else if (h.periksa) setHasil(h.periksa)
    } catch (err) {
      setGalat(err)
    } finally {
      setMenyimpan(false)
    }
  }

  return (
    <section className="inbox edmf-tema edmf-form">
      <header className="inbox__kepala">
        <h2 className="inbox__judul edmf-tema__judul">{FORM.judul}</h2>
      </header>
      <Gagal galat={galat} />

      <div className="panel edmf-form__panel">
        <div className="edmf-form__kolom">
          <div className="edmf-form__tumpuk">
            <Field label={FORM.policyNo} value={policyNo} onChange={setPolicyNo} required autoFocus error={pesanMedan('PolicyNo')} />
            <Field label={FORM.sumber} value={sumber} onChange={setSumber} required />
            <Area label={FORM.catatan} value={catatan} onChange={setCatatan} required baris={4} />
            <Pilih
              label={FORM.jenis}
              value={edmType}
              onChange={(v) => {
                setEdmType(v)
                setEdmTypeNew('')
                setDetail('')
              }}
              opsi={opsi(JENIS_EDM)}
              kosong={FORM.pilihJenis}
              required
            />
            {edmType === '4' && (
              <Pilih
                label={FORM.jenisBaru}
                value={edmTypeNew}
                onChange={(v) => {
                  setEdmTypeNew(v)
                  setDetail('')
                }}
                opsi={opsi(JENIS_EDM_BARU)}
                kosong={FORM.jenisBaru}
                required
              />
            )}
            {daftarDetail && (
              <Pilih label={FORM.detail} value={detail} onChange={setDetail} opsi={opsi(daftarDetail)} kosong={FORM.detail} required />
            )}
          </div>
          <div className="edmf-form__tumpuk">
            <TanggalDMY
              label={FORM.tanggal}
              value={tanggal}
              onChange={setTanggal}
              required
              error={pesanMedan('EndorsementDate')}
              labelKalender={FORM.kalender}
              pesanFormat={FORM.formatTanggal}
            />
            {!life && <Pilih label={FORM.status} value={status} onChange={setStatus} opsi={opsi(STATUS_EDM)} required />}
            {!life && (
              <div className="field">
                <span className="field__label">
                  {FORM.retro}
                  <span className="field__req">*</span>
                </span>
                <div className="edmf-pilihan" role="radiogroup" aria-label={FORM.retro}>
                  {PILIHAN_RETRO.map((p) => (
                    <label key={p.nilai} className="edmf-pilihan__butir">
                      <input type="radio" name="edmf-retro" value={p.nilai} checked={retro === p.nilai} onChange={() => setRetro(p.nilai)} />{' '}
                      {p.label}
                    </label>
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>
      </div>

      <div className="edmf-form__kaki">
        <p className="edmf-form__alasan" role="status">
          {alasan}
        </p>
        <div className="edmf-form__aksi">
          <button
            type="button"
            className="btn btn--primary edmf-tombol edmf-tombol--utama"
            disabled={!bolehOK}
            title={bolehOK ? undefined : PERIKSA.okMati}
            onClick={() => void simpan()}
          >
            {menyimpan ? PERIKSA.menyimpan : FORM.ok}
          </button>
          <button type="button" className="btn edmf-tombol" onClick={onBatal}>
            {FORM.batal}
          </button>
        </div>
      </div>
    </section>
  )
}
