// Add New accumulation (tiket 46) - port sub-section `NB FacIn\Section\InputAccumulationCov.xml` di popup Choose
// Accumulation (gambar layar Pega DEV dari work owner 04-10-2026):
// - Accumulation Type* (RD BrowseAccumulatedType_RD: `.ID` -> Accumulation, `.Note` -> AccumulationName tampil di
//   sampingnya), Province* (`.ID` -> ProvinceID), Key Word*, Accumulation Description*; kanan: Scope Area*, Cresta
//   Zone* (`.Code`; `.ID` -> CZoneID), Primary Zip Code* (baca-saja) + Choose Zip Code;
// - awal (`InputAccumulationAdd_PreDT` + `GetCzone_Act`): isian kosong, Cresta Zone dari zip saringan pencarian
//   (`GetAccumulationZipcode_SQL`);
// - Choose Zip Code (`ChooseZipCodeDtl`, RD BrowseRiskAddressZipCode_RD `.ProvinceName Contains`): memilih zip mengisi
//   Primary Zip Code + NationInitial (`SetCountryID_Act`, awalan ID) lalu Cresta Zone dari zip itu (`GetCzone_Act`);
// - Save (`SaveAccumulation_Act` -> RDBMASTERACCUMULATION, diport ke Go tanpa CALL) -> `POST …/akumulasi`; Close
//   (`CancelInputAccumulation_Act`) kembali ke pencarian.
//
// Keputusan agent (tiket 46): T-1 (DICABUT 04-10-2026, work owner: "dibuatin pop up aja") Choose Zip Code kini popup
// di atas popup Choose Accumulation (portal, lihat `PopupZip`). T-2 sesudah Save berhasil, accumulation baru langsung dipakai coverage dan popup ditutup (Pega menulisnya ke
// coverage lalu kembali ke mode cari tanpa menutup popup). T-3 Save yang ganda TIDAK menulis apa pun ke coverage
// (Pega menulis "UNKNOWNID"). T-4 tombol "Add" tipe akumulasi baru (ModalAccumulation_FacIn) = tahap berikut; menu
// Master Accumulated Type menggantikannya. T-5 Province = tabel PROVINCE (RD PROVINCE-class `BrowseProvince_RD`
// tidak ada di korpus).

import { useEffect, useId, useRef, useState } from 'react'
import { createPortal } from 'react-dom'

import { Area, Field, Gagal, IkonTutup, Kosong, Memuat, Pilih } from '../../../../inti/frontend/components/ui/dasar'
import {
  cariZipAkumulasi,
  czoneDariZip,
  saranAkumulasi,
  tambahAkumulasi,
  type AkumulasiBaru,
  type BarisZipAkumulasi,
  type HalamanZipAkumulasi,
} from '../api'
import {
  OPSI_KEYWORD,
  OPSI_SCOPE_AREA,
  PILIHAN_KOSONG_TAMBAH,
  TAMBAH_AKUMULASI as T,
  TEKS_TAMBAH_AKUMULASI as TEKS,
} from '../labels'
import IsianSaran from './IsianSaran'
import PagerHalaman from './PagerHalaman'

/** Isi layar form Add New (nilai kiriman + teks tampil). */
export interface IsianTambahAkumulasi extends AkumulasiBaru {
  /** Nama tipe dari RD (`.Note` -> AccumulationName), tampil di samping Accumulation Type. */
  accumulationName: string
  provinceName: string
}

export const TAMBAH_KOSONG: IsianTambahAkumulasi = {
  accumulation: '',
  accumulationType: '',
  accumulationName: '',
  note: '',
  keyword: '',
  scopeArea: '',
  cZone: '',
  cZoneId: '',
  provinceId: '',
  provinceName: '',
  zipCode: '',
  negara: '',
}

/** Medan wajib (sel bertanda * di InputAccumulationCov) -> galat per medan. */
export function periksaTambah(s: IsianTambahAkumulasi): Partial<Record<keyof IsianTambahAkumulasi, string>> {
  const g: Partial<Record<keyof IsianTambahAkumulasi, string>> = {}
  const wajib: (keyof IsianTambahAkumulasi)[] = [
    'accumulationType',
    'provinceName',
    'keyword',
    'note',
    'scopeArea',
    'cZone',
    'zipCode',
  ]
  for (const k of wajib) if (s[k].trim() === '') g[k] = TEKS.wajib
  // Accumulation Type / Province harus DIPILIH dari saran: backend menyimpan ID-nya (Accumulation, ProvinceID).
  if (!g.accumulationType && s.accumulation === '') g.accumulationType = TEKS.pilihDariSaran
  if (!g.provinceName && s.provinceId === '') g.provinceName = TEKS.pilihDariSaran
  return g
}

/** Badan kiriman (tanpa teks tampil). */
export function keBadan(s: IsianTambahAkumulasi): AkumulasiBaru {
  return {
    accumulation: s.accumulation,
    accumulationType: s.accumulationType.trim(),
    note: s.note.trim(),
    keyword: s.keyword,
    scopeArea: s.scopeArea,
    cZone: s.cZone.trim(),
    cZoneId: s.cZoneId,
    provinceId: s.provinceId,
    zipCode: s.zipCode,
    negara: s.negara,
  }
}

/** Baris per halaman Choose Zip Code (permintaan work owner 04-10-2026: "15 list perpage"). */
export const UKURAN_HALAMAN_ZIP = 15

/**
 * Popup Choose Zip Code (flow action `ChooseZipCode` -> section `ChooseZipCodeDtl`, RD `BrowseRiskAddressZipCode_RD` atas
 * RISKADDRESS menurut ProvinceName yang dipilih di form) - paging di server, 15 baris per halaman.
 *
 * Popup di ATAS popup Choose Accumulation (permintaan work owner 04-10-2026: "dibuatin pop up aja dong, jgn dibwh gini";
 * mencabut T-1). `Modal` inti tidak dapat bertumpuk (Escape-nya didengar di `document` - satu tekan menutup keduanya; dan
 * modal yang dirender di dalam modal lain terkurung kotaknya), jadi popup ini dirender lewat PORTAL ke `document.body`
 * dengan kelas modal inti, dibungkus akar `.nbfacin`, dan Escape ditangkap lebih dulu (fase capture di `window`) supaya
 * hanya popup ini yang tertutup.
 */
function PopupZip({
  provinceName,
  onPilih,
  onTutup,
}: {
  provinceName: string
  onPilih: (b: BarisZipAkumulasi) => void
  onTutup: () => void
}) {
  const [q, setQ] = useState('')
  const [halaman, setHalaman] = useState(1)
  const [data, setData] = useState<HalamanZipAkumulasi | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const nomor = useRef(0)
  const idJudul = useId()

  // Escape menutup HANYA popup ini: ditangkap di fase capture `window` sebelum sampai ke pendengar `document` milik
  // Modal Choose Accumulation di bawahnya.
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key !== 'Escape') return
      e.stopPropagation()
      onTutup()
    }
    window.addEventListener('keydown', onKey, true)
    return () => window.removeEventListener('keydown', onKey, true)
  }, [onTutup])

  // Kata cari / province berubah -> kembali ke halaman 1.
  useEffect(() => {
    setHalaman(1)
  }, [q, provinceName])

  useEffect(() => {
    const n = ++nomor.current
    setGalat(null)
    const jadwal = window.setTimeout(
      () => {
        setData(null)
        cariZipAkumulasi(provinceName, q, halaman).then(
          (h) => {
            if (n === nomor.current) setData(h)
          },
          (err: unknown) => {
            if (n === nomor.current) setGalat(err)
          },
        )
      },
      q === '' ? 0 : 400,
    )
    return () => window.clearTimeout(jadwal)
  }, [provinceName, q, halaman])

  return createPortal(
    <div className="nbfacin">
      <div
        className="modal__backdrop nbf-akum__lapis-zip"
        role="presentation"
        onClick={(e) => {
          if (e.target === e.currentTarget) onTutup()
        }}
      >
        <div className="modal modal--lebar" role="dialog" aria-modal="true" aria-labelledby={idJudul}>
          <div className="modal__head">
            <h3 className="modal__title" id={idJudul}>
              {TEKS.judulZip}
            </h3>
            <button
              type="button"
              className="modal__close"
              aria-label={TEKS.tutupZip}
              title={TEKS.tutupZip}
              onClick={onTutup}
            >
              <IkonTutup />
            </button>
          </div>
          <div className="modal__body">
            <div className="nbf-akum__panel-zip">
              <div className="nbf-akum__zip-kepala">
                <p className="nbf-akum__sub">{TEKS.petunjukZip}</p>
                <span className="nbf-akum__chip">
                  {TEKS.provinsiZip}
                  <strong>{provinceName.trim() || TEKS.semuaProvinsi}</strong>
                </span>
              </div>
              <input
                className="field__input nbf-akum__zip-cari"
                type="search"
                aria-label={TEKS.cariZip}
                placeholder={TEKS.cariZip}
                value={q}
                onChange={(e) => setQ(e.target.value)}
              />
              {data === null && galat === null && <Memuat />}
              <Gagal galat={galat} />
              {data !== null && data.baris.length === 0 && <Kosong pesan={TEKS.tanpaZip} />}
              {data !== null && data.baris.length > 0 && (
                <>
                  <div className="table-wrap">
                    <table className="nbf-tabel">
                      <thead>
                        <tr>
                          {T.kolomZip.map((k) => (
                            <th key={k} scope="col">
                              {k}
                            </th>
                          ))}
                          <th scope="col" />
                        </tr>
                      </thead>
                      <tbody>
                        {data.baris.map((b, i) => (
                          <tr key={`${i}-${b.zipCode}-${b.city}`}>
                            <td>
                              <code className="nbf-akum__kode">{b.zipCode}</code>
                            </td>
                            <td>{b.city}</td>
                            <td>{b.province}</td>
                            <td>{b.nation}</td>
                            <td className="table__actions">
                              <button type="button" className="btn btn--primary btn--sm" onClick={() => onPilih(b)}>
                                {T.pilih}
                              </button>
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                  {data.total > data.ukuran && (
                    <PagerHalaman
                      total={data.total}
                      halaman={data.halaman}
                      ukuran={data.ukuran}
                      onPindah={setHalaman}
                    />
                  )}
                </>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>,
    document.body,
  )
}

export default function FormTambahAkumulasi({
  zipSaring,
  onTutup,
  onTersimpan,
}: {
  /** Zip saringan pencarian (`SearchAccumulation.PostalCode`) - asal Cresta Zone awal (`GetCzone_Act`). */
  zipSaring: string
  /** Close: kembali ke pencarian. */
  onTutup: () => void
  /** Save berhasil: ID baru + Note tersimpan. */
  onTersimpan: (id: string, note: string) => void
}) {
  const [s, setS] = useState<IsianTambahAkumulasi>(TAMBAH_KOSONG)
  const [panelZip, setPanelZip] = useState(false)
  const [coba, setCoba] = useState(false)
  const [menyimpan, setMenyimpan] = useState(false)
  const [galat, setGalat] = useState<unknown>(null)
  const set = (k: keyof IsianTambahAkumulasi) => (v: string) => setS((x) => ({ ...x, [k]: v }))
  const salah = coba ? periksaTambah(s) : {}

  /**
   * Cresta Zone dari CZONE tabel RW zip itu (work owner 04-10-2026: "otomatis keisi dari tabel rw, diambil dari czone
   * zipcode yang dipilih"). `ganti` = zip DIPILIH lewat Choose Zip Code: Cresta Zone selalu mengikuti zip itu, termasuk
   * dikosongkan bila RW tidak punya zone-nya. Pengisian awal (zip pencarian) hanya mengisi bila ada nilainya.
   */
  function isiCzone(zip: string, ganti: boolean) {
    if (zip.trim() === '') return
    czoneDariZip(zip).then(
      (h) => {
        if (ganti || h.cZone !== '') setS((x) => ({ ...x, cZone: h.cZone, cZoneId: h.cZoneId }))
      },
      () => {},
    )
  }

  // Awal: Cresta Zone dari zip saringan pencarian.
  useEffect(() => {
    isiCzone(zipSaring, false)
    // Sekali saat form dibuka.
  }, [])

  async function simpan() {
    setCoba(true)
    if (Object.keys(periksaTambah(s)).length > 0) return
    setMenyimpan(true)
    setGalat(null)
    try {
      const h = await tambahAkumulasi(keBadan(s))
      onTersimpan(h.id, h.note)
    } catch (err) {
      setGalat(err)
    } finally {
      setMenyimpan(false)
    }
  }

  return (
    <form
      className="nbf-akum__kartu"
      onSubmit={(e) => {
        e.preventDefault()
        void simpan()
      }}
    >
      <div className="nbf-akum__kartu-kepala">
        <div>
          <h5 className="nbf-akum__kartu-judul">{TEKS.judul}</h5>
          <p className="nbf-akum__sub">{TEKS.sub}</p>
        </div>
      </div>
      <div className="nbf-akum__grid">
        <div className="nbf-akum__lebar">
          <IsianSaran
            label={`${T.accumulationType.label} *`}
            value={s.accumulationType}
            onKetik={(v) =>
              setS((x) => ({
                ...x,
                accumulationType: v,
                accumulation: '',
                accumulationName: '',
              }))
            }
            cari={(q) => saranAkumulasi('accumtype', q)}
            onPilih={(x) =>
              setS((y) => ({
                ...y,
                accumulationType: x.label,
                accumulation: x.id,
                accumulationName: x.ekstra ?? '',
              }))
            }
            tampil={(x) => [x.label, x.ekstra].filter(Boolean).join(' · ')}
          />
          {s.accumulationName !== '' && <span className="nbf-akum__catatan">{s.accumulationName}</span>}
          {salah.accumulationType && <div className="field__error">{salah.accumulationType}</div>}
        </div>
        <div>
          <Pilih
            label={T.scopeArea.label}
            value={s.scopeArea}
            onChange={set('scopeArea')}
            opsi={OPSI_SCOPE_AREA}
            kosong={PILIHAN_KOSONG_TAMBAH}
            required
            error={salah.scopeArea}
          />
        </div>
        <div>
          <IsianSaran
            label={`${T.crestaZone.label} *`}
            value={s.cZone}
            onKetik={(v) => setS((x) => ({ ...x, cZone: v, cZoneId: '' }))}
            cari={(q) => saranAkumulasi('czone', q)}
            onPilih={(x) => setS((y) => ({ ...y, cZone: x.label, cZoneId: x.id }))}
          />
          {salah.cZone && <div className="field__error">{salah.cZone}</div>}
        </div>
        <div className="nbf-akum__lebar">
          <IsianSaran
            label={`${T.province.label} *`}
            value={s.provinceName}
            onKetik={(v) => setS((x) => ({ ...x, provinceName: v, provinceId: '' }))}
            cari={(q) => saranAkumulasi('province', q)}
            onPilih={(x) => setS((y) => ({ ...y, provinceName: x.label, provinceId: x.id }))}
          />
          {salah.provinceName && <div className="field__error">{salah.provinceName}</div>}
        </div>
        <div className="nbf-akum__lebar nbf-akum__zip-baris">
          <Field
            label={T.primaryZip.label}
            value={s.zipCode || TEKS.zipBelum}
            onChange={() => {}}
            readOnly
            required
            error={salah.zipCode}
          />
          <button type="button" className="btn btn--ghost" aria-haspopup="dialog" onClick={() => setPanelZip(true)}>
            {T.pilihZip.label}
          </button>
        </div>
        {panelZip && (
          <div className="nbf-akum__penuh">
            <PopupZip
              provinceName={s.provinceName}
              onTutup={() => setPanelZip(false)}
              onPilih={(b) => {
                setS((x) => ({
                  ...x,
                  zipCode: b.zipCode,
                  negara: b.nationInitial,
                }))
                setPanelZip(false)
                isiCzone(b.zipCode, true)
              }}
            />
          </div>
        )}
        <div className="nbf-akum__lebar">
          <Pilih
            label={T.keyword.label}
            value={s.keyword}
            onChange={set('keyword')}
            opsi={OPSI_KEYWORD}
            kosong={PILIHAN_KOSONG_TAMBAH}
            required
            error={salah.keyword}
          />
        </div>
        <div className="nbf-akum__penuh">
          <Area label={T.note.label} value={s.note} onChange={set('note')} baris={3} required error={salah.note} />
        </div>
      </div>
      <Gagal galat={galat} />
      <div className="nbf-akum__kaki">
        <button type="button" className="btn btn--ghost" onClick={onTutup}>
          {T.tutup.label}
        </button>
        <button type="submit" className="btn btn--primary" disabled={menyimpan}>
          {menyimpan ? TEKS.menyimpan : T.simpan.label}
        </button>
      </div>
    </form>
  )
}
