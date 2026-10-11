// Form Create Addendum - irisan I1-I3 (keputusan work owner 11-10-2026, W-EDM-01…12).
//
// Port `Section/WorkPrimaryDetails.xml` (kelas `ASM-SFAGIS-Work-Endorsement`) dan tombolnya
// `Section/crmNewHarnessButtons.xml`. Medan, label, dan syarat tampil mengikuti section itu:
//   1 Policy Number / Cover Number (`.PolicyNo`, teks bebas)      7 Endorsement Date (`.EndorsementDate`)
//   2 Endorsement Source (`.EndorsementSourceNote`)              9 Biaya (`.AdminCharge`)  10 Survey (`.Survey`)
//   3 Endorsement Note (`.Note`)                                 11 Endorsement Internal Retro - bukan Life
//   4 Endorsement Type (`.Quotation.EdmType`)
//   5 EDM Type (`.Quotation.EdmTypeNew`) - hanya EdmType 4
//   6 Detail EDM Type (`.Quotation.Type`, `EdmType1/2/3`) - EdmType 4, bukan Life, EdmTypeNew 1/2/3
//
// Tata letak = templat form NB FacIn (work owner 11-10-2026: "INI RAPIKAN DONG, BIAR USER FRIENDLY. TEMPLATE NYA JUGA
// SAMAIN KAYA NB FACIN"; pola `modul/nbfacin/frontend/pages/FormOpportunity.tsx`, disalin, tidak diimpor): kartu bertema,
// judul + tombol di kepala, dua kolom (kiri: medan 1-3 dan 7; kanan: 4-6, 11, 9-10 - urutan baca Pega per kolom), tanggal
// dd/mm/yyyy (`components/TanggalDMY.tsx`), alasan tombol OK mati di bawah kepala, hasil pemeriksaan di bawah form.
//
// Pemeriksaan (I2): di Pega `SetErrorBatalEndorsement_Act` berjalan saat PolicyNo / Endorsement Type / EDM Type /
// Internal Retro berubah, `CheckEDMPolisDate` saat Endorsement Date berubah; tombol OK aktif hanya bila keduanya lolos
// (`InputData.CARIGROUP=='true' && InputSearch.CARIGROUP=='true'`, crmNewHarnessButtons L885). Di sini keduanya satu
// panggilan `POST /api/endorsmentfacin/periksa` - NOL tulis.
//
// OK (irisan I3): `POST /api/endorsmentfacin/kasus` = `SetValueToEDMWork`. Aktif bila pemeriksaan terakhir lolos dan
// lengkap serta medan wajib terisi; server memeriksa ulang. Ditolak di server = jawaban 200 tanpa `caseId` (hasil bisnis,
// spec 07 §7.2) - rinciannya menggantikan hasil pemeriksaan di layar.
//
// Belum (laporan modul): Endorsement Status (`.EndorsementStatus`, prompt list tidak ada di korpus - A46) dan sumber
// pilihan Endorsement Source (`D_EnumerationList` `Type=sumberedm`, isi tabelnya tidak ada di korpus - sementara teks
// bebas, sama dengan autocomplete Pega yang mengizinkan isian bebas).

import { useEffect, useRef, useState } from 'react'

import { Area, Field, Gagal, Pilih } from '../../../../inti/frontend/components/ui/dasar'
import { buatKasus, periksaAddendum, type HasilPeriksa } from '../api'
import TanggalDMY from '../components/TanggalDMY'
import { DETAIL_EDM, FORM, JENIS_EDM, JENIS_EDM_BARU, PERIKSA, PILIHAN_RETRO, type Pilihan } from '../labels'

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
  const [biaya, setBiaya] = useState(false)
  const [survey, setSurvey] = useState(false)
  // Bawaan `0` - WorkPrimaryDetails L10212.
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

  // Medan wajib (`pyRequired` WorkPrimaryDetails) yang masih kosong, urut layar; Endorsement Status tidak (A46).
  const kosong = (
    [
      [FORM.policyNo, polis],
      [FORM.sumber, sumber.trim()],
      [FORM.catatan, catatan.trim()],
      [FORM.tanggal, tanggal],
      [FORM.jenis, edmType],
      ...(edmType === '4' ? ([[FORM.jenisBaru, edmTypeNew]] as [string, string][]) : []),
      ...(daftarDetail ? ([[FORM.detail, detail]] as [string, string][]) : []),
      [FORM.retro, retro],
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
        endorsementInternalRetro: retro,
        note: catatan,
        endorsementSource: '',
        endorsementSourceNote: sumber.trim(),
        endorsementStatus: '',
        detailType: daftarDetail ? detail : '',
        adminCharge: biaya,
        survey,
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
      <header className="inbox__kepala edmf-form__kepala">
        <h2 className="inbox__judul edmf-tema__judul">{FORM.judul}</h2>
        <div className="edmf-form__aksi">
          <button type="button" className="btn edmf-tombol" onClick={onBatal}>
            {FORM.batal}
          </button>
          <button
            type="button"
            className="btn btn--primary edmf-tombol edmf-tombol--utama"
            disabled={!bolehOK}
            title={bolehOK ? undefined : PERIKSA.okMati}
            onClick={() => void simpan()}
          >
            {menyimpan ? PERIKSA.menyimpan : FORM.ok}
          </button>
        </div>
      </header>
      {alasan !== '' && (
        <p className="edmf-form__alasan" role="status">
          {alasan}
        </p>
      )}
      <Gagal galat={galat} />

      <div className="panel edmf-form__panel">
        <div className="edmf-form__kolom">
          <div className="edmf-form__tumpuk">
            <h3 className="edmf-form__bagian">{FORM.bagianPolis}</h3>
            <Field label={FORM.policyNo} value={policyNo} onChange={setPolicyNo} required autoFocus error={pesanMedan('PolicyNo')} />
            <TanggalDMY
              label={FORM.tanggal}
              value={tanggal}
              onChange={setTanggal}
              required
              error={pesanMedan('EndorsementDate')}
              labelKalender={FORM.kalender}
              pesanFormat={FORM.formatTanggal}
            />
            <Field label={FORM.sumber} value={sumber} onChange={setSumber} required />
            <Area label={FORM.catatan} value={catatan} onChange={setCatatan} required baris={4} />
          </div>
          <div className="edmf-form__tumpuk">
            <h3 className="edmf-form__bagian">{FORM.bagianJenis}</h3>
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
            <div className="field">
              <span className="field__label">{FORM.bagianLain}</span>
              <div className="edmf-pilihan">
                <label className="edmf-pilihan__butir">
                  <input type="checkbox" checked={biaya} onChange={(e) => setBiaya(e.target.checked)} /> {FORM.biaya}
                </label>
                <label className="edmf-pilihan__butir">
                  <input type="checkbox" checked={survey} onChange={(e) => setSurvey(e.target.checked)} /> {FORM.survey}
                </label>
              </div>
            </div>
          </div>
        </div>
      </div>

      <HasilPemeriksaan hasil={hasil} memeriksa={memeriksa} adaPolis={polis !== ''} />
    </section>
  )
}

function HasilPemeriksaan({ hasil, memeriksa, adaPolis }: { hasil: HasilPeriksa | null; memeriksa: boolean; adaPolis: boolean }) {
  return (
    <div className="panel edmf-periksa" aria-live="polite">
      <h3 className="edmf-form__bagian">{PERIKSA.judulGerbang}</h3>
      {!adaPolis && <p className="edmf-periksa__kosong">{PERIKSA.belumDiperiksa}</p>}
      {adaPolis && !hasil && memeriksa && <p className="edmf-periksa__kosong">{PERIKSA.memeriksa}</p>}
      {hasil && (
        <>
          {!hasil.ditolak && hasil.lengkap && <div className="alert alert--ok">{PERIKSA.lolos}</div>}
          {!hasil.lengkap && <div className="alert alert--warn">{PERIKSA.belumLengkap}</div>}
          <ul className="edmf-periksa__daftar">
            {hasil.gerbang.map((g) => (
              <li key={g.no} className="edmf-periksa__butir">
                <span className="edmf-periksa__nama">{PERIKSA.gerbang[g.no] ?? String(g.no)}</span>
                <span className={`edmf-chip edmf-chip--${g.hasil}`}>{PERIKSA.hasil[g.hasil] ?? g.hasil}</span>
                {g.catatan !== '' && <span className="edmf-periksa__catatan">{g.catatan}</span>}
              </li>
            ))}
            <li className="edmf-periksa__butir">
              <span className="edmf-periksa__nama">{PERIKSA.bisnis}</span>
              <span className={`edmf-chip edmf-chip--${hasil.bisnis.hasil}`}>{hasil.bisnis.jenisBisnis || PERIKSA.hasil[hasil.bisnis.hasil]}</span>
            </li>
            {hasil.tanggal && (
              <li className="edmf-periksa__butir">
                <span className="edmf-periksa__nama">{PERIKSA.periode(hasil.tanggal.mulai, hasil.tanggal.akhir)}</span>
                <span className={`edmf-chip edmf-chip--${hasil.tanggal.hasil}`}>{PERIKSA.hasil[hasil.tanggal.hasil] ?? hasil.tanggal.hasil}</span>
                {hasil.tanggal.catatan && <span className="edmf-periksa__catatan">{hasil.tanggal.catatan}</span>}
              </li>
            )}
          </ul>
        </>
      )}
    </div>
  )
}
