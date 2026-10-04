// Popup Choose Accumulation (tiket 46) - port harness `NB FacIn\Harness\ChooseAccumulation_FacIn.xml` -> section
// `SearchRiskAccumCov.xml`:
// - judul "Search Risk Accumulation", bar "Conditions"; isian 70 % kiri, tombol 30 % kanan ("Inline grid 70 30");
//   baris Accumulation Code | Policy No, Road, dua kolom (Zip Code, Country, Province, Accum. Type | City, District,
//   Area, CZone), Key Word;
// - awal: Zip Code = zip lokasi risiko (`SearchAccumAct` langkah 4 `SearchAccumulation.PostalCode = .RiskZipCode`);
// - autocomplete (RD sel 27-35) menyalin nilai seperti Pega: Country -> NationInitial (SyariahStatus), Province -> ID,
//   City -> ID (`InputFilter.City`), District -> ID (`InputFilter.District`), Area -> ZipCode (`InputFilter.PostalCode`,
//   menimpa Zip Code saat Filter, `GetDataAccumulation_act` langkah 4);
// - Filter (sel 40) -> `GET …/akumulasi` (backend memilih jalur SQL / RD seperti `GetDataAccumulation_act`); Clear Column
//   (sel 41, `SearchAccumulationPre_Act`) mengosongkan saringan dan hasil;
// - Choose (sel 77 / 119, `SetDataAccum_Act`): zip dari ID harus sama dengan zip lokasi risiko; beda -> pesan, popup
//   tetap terbuka; sama -> AccumulationCode = ID, AccumulationDescription = Note, popup tertutup.
//
// Keputusan agent (tiket 46): K-1 Add New, Summary (.Note), Accumulation Risk, Risk Accumulation report, Adjustment =
// tahap berikut. K-2 tanpa paging 10 baris (hasil maks. 500). K-7 mengetik di autocomplete mengosongkan nilai hasil
// pilihan sebelumnya (Pega membiarkan ID lama). K-8 Accum. Type hanya tampilan: di Pega terikat
// `SearchAccumulation.Type`, sedangkan RD membaca `SearchAccumulation.AccumulationType` (tidak pernah terisi).

import { useState } from 'react'

import { Field, Gagal, Kosong, Memuat, Modal, Pilih } from '../../../../inti/frontend/components/ui/dasar'
import { cariAkumulasi, saranAkumulasi, type BarisAkumulasi, type SaringAkumulasi } from '../api'
import { OPSI_KEYWORD, POPUP_AKUMULASI as P, TEKS_AKUMULASI } from '../labels'
import IsianSaran from './IsianSaran'

/**
 * Zip di dalam Accumulation Code - `SetDataAccum_Act` langkah 2-3: `@substring(ID,4,9)` tanpa "-", kecuali karakter
 * ke-3 (0-based) bukan "-" -> `@substring(ID,3,8)`. Substring gaya Java (awal inklusif, akhir eksklusif); ID terlalu
 * pendek -> kosong (Pega melempar galat) `[dugaan]`.
 */
export function zipDariId(id: string): string {
  const sub = (a: number, b: number) => (id.length >= b ? id.slice(a, b) : null)
  const c3 = sub(3, 4)
  if (c3 !== null && c3 !== '-') return sub(3, 8) ?? ''
  return (sub(4, 9) ?? '').replaceAll('-', '')
}

/** Isi layar: teks yang tampil + nilai yang disalin saat memilih saran. */
export interface IsianAkumulasi {
  id: string
  policyNo: string
  note: string
  postalCode: string
  nation: string
  syariahStatus: string
  provinceName: string
  provinceId: string
  accumType: string
  city: string
  cityId: string
  district: string
  districtId: string
  area: string
  /** `InputFilter.PostalCode` dari Area terpilih. */
  areaZip: string
  czone: string
  keyword: string
}

export const ISIAN_KOSONG: IsianAkumulasi = {
  id: '', policyNo: '', note: '', postalCode: '', nation: '', syariahStatus: '', provinceName: '', provinceId: '',
  accumType: '', city: '', cityId: '', district: '', districtId: '', area: '', areaZip: '', czone: '', keyword: '',
}

/** Saringan yang dikirim tombol Filter - Area terpilih menimpa Zip Code (`GetDataAccumulation_act` langkah 4). */
export function keSaring(s: IsianAkumulasi): SaringAkumulasi {
  return {
    id: s.id, policyNo: s.policyNo, note: s.note, postalCode: s.areaZip || s.postalCode, syariahStatus: s.syariahStatus,
    provinceId: s.provinceId, cityId: s.cityId, districtId: s.districtId, czone: s.czone, keyword: s.keyword,
  }
}

export default function PopupAkumulasi({
  zipRisiko,
  onTutup,
  onPilih,
}: {
  /** `.Property.RiskLocation.ASMZipCode` objek pemilik coverage. */
  zipRisiko: string
  onTutup: () => void
  onPilih: (kode: string, alamat: string) => void
}) {
  const [s, setS] = useState<IsianAkumulasi>({ ...ISIAN_KOSONG, postalCode: zipRisiko })
  const [hasil, setHasil] = useState<BarisAkumulasi[] | null>(null)
  const [memuat, setMemuat] = useState(false)
  const [galat, setGalat] = useState<unknown>(null)
  const [pesan, setPesan] = useState('')
  const set = (k: keyof IsianAkumulasi) => (v: string) => setS((x) => ({ ...x, [k]: v }))
  /** Ketik di autocomplete: teks berubah, nilai pilihan lama dikosongkan (K-7). */
  const ketik = (teks: keyof IsianAkumulasi, ...nilai: (keyof IsianAkumulasi)[]) => (v: string) =>
    setS((x) => ({ ...x, [teks]: v, ...Object.fromEntries(nilai.map((k) => [k, ''])) }))

  async function filter() {
    setMemuat(true)
    setGalat(null)
    setPesan('')
    // Langkah 4: Area terpilih menimpa Zip Code yang tampil.
    if (s.areaZip !== '') setS((x) => ({ ...x, postalCode: x.areaZip }))
    try {
      setHasil((await cariAkumulasi(keSaring(s))).baris)
    } catch (err) {
      setGalat(err)
    } finally {
      setMemuat(false)
    }
  }

  function pilih(b: BarisAkumulasi) {
    const dariId = zipDariId(b.id)
    if (dariId !== zipRisiko) {
      setPesan(TEKS_AKUMULASI.zipBeda(zipRisiko, dariId))
      return
    }
    onPilih(b.id, b.note)
  }

  return (
    <Modal judul={P.judul} onTutup={onTutup} lebar>
      <h4 className="nbf-objek__judul">{P.cari.label}</h4>
      <h5 className="nbf-objek__judul">{P.kondisi}</h5>
      <div className="nbf-akum__saring">
        <div className="nbf-akum__isian nbf-labelkiri">
          <div className="nbf-akum__dua">
            <div className="nbf-akum__kolom nbf-labelkiri">
              <Field label={P.accumulationCode.label} value={s.id} onChange={set('id')} />
            </div>
            <div className="nbf-akum__kolom nbf-labelkiri">
              <Field label={P.policyNo.label} value={s.policyNo} onChange={set('policyNo')} />
            </div>
          </div>
          <Field label={P.road.label} value={s.note} onChange={set('note')} />
          <div className="nbf-akum__dua">
            <div className="nbf-akum__kolom nbf-labelkiri">
              <Field label={P.zipCode.label} value={s.postalCode} onChange={set('postalCode')} />
              <IsianSaran
                label={P.country.label}
                value={s.nation}
                onKetik={ketik('nation', 'syariahStatus')}
                cari={(q) => saranAkumulasi('nation', q)}
                onPilih={(x) => setS((y) => ({ ...y, nation: x.label, syariahStatus: x.ekstra ?? '' }))}
              />
              <IsianSaran
                label={P.province.label}
                value={s.provinceName}
                onKetik={ketik('provinceName', 'provinceId')}
                cari={(q) => saranAkumulasi('province', q, s.nation)}
                onPilih={(x) => setS((y) => ({ ...y, provinceName: x.label, provinceId: x.id }))}
              />
              <IsianSaran
                label={P.accumType.label}
                value={s.accumType}
                onKetik={set('accumType')}
                cari={(q) => saranAkumulasi('accumtype', q)}
                onPilih={(x) => setS((y) => ({ ...y, accumType: x.label }))}
              />
            </div>
            <div className="nbf-akum__kolom nbf-labelkiri">
              <IsianSaran
                label={P.city.label}
                value={s.city}
                onKetik={ketik('city', 'cityId')}
                cari={(q) => saranAkumulasi('city', q, s.provinceId)}
                onPilih={(x) => setS((y) => ({ ...y, city: x.label, cityId: x.id }))}
              />
              <IsianSaran
                label={P.district.label}
                value={s.district}
                onKetik={ketik('district', 'districtId')}
                cari={(q) => saranAkumulasi('district', q, s.city)}
                onPilih={(x) => setS((y) => ({ ...y, district: x.label, districtId: x.id }))}
              />
              <IsianSaran
                label={P.area.label}
                value={s.area}
                onKetik={ketik('area', 'areaZip')}
                cari={(q) => saranAkumulasi('area', q, s.district)}
                onPilih={(x) => setS((y) => ({ ...y, area: x.label, areaZip: x.ekstra ?? '' }))}
                tampil={(x) => [x.label, x.ekstra].filter(Boolean).join(' · ')}
              />
              <IsianSaran
                label={P.czone.label}
                value={s.czone}
                onKetik={set('czone')}
                cari={(q) => saranAkumulasi('czone', q)}
                onPilih={(x) => setS((y) => ({ ...y, czone: x.label }))}
              />
            </div>
          </div>
          <Pilih label={P.keyword.label} value={s.keyword} onChange={set('keyword')} opsi={OPSI_KEYWORD} kosong="" />
        </div>
        <div className="nbf-akum__tombol">
          <button type="button" className="btn btn--sm" onClick={() => void filter()} disabled={memuat}>
            {P.filter.label}
          </button>
          <button
            type="button"
            className="btn btn--ghost btn--sm"
            onClick={() => {
              setS(ISIAN_KOSONG)
              setHasil(null)
              setPesan('')
              setGalat(null)
            }}
          >
            {P.bersih.label}
          </button>
        </div>
      </div>
      {pesan && <div className="alert alert--error">{pesan}</div>}
      {memuat && <Memuat />}
      <Gagal galat={galat} />
      {hasil !== null && !memuat && hasil.length > 0 && (
        <div className="table-wrap">
          <table className="nbf-tabel">
            <thead>
              <tr>
                {P.kolom.map((k) => (
                  <th key={k.sel} scope="col">
                    {k.label}
                  </th>
                ))}
                <th scope="col" />
              </tr>
            </thead>
            <tbody>
              {hasil.map((b) => (
                <tr key={b.id}>
                  <td>{b.id}</td>
                  <td>{b.accumulationName}</td>
                  <td>{b.note}</td>
                  <td className="table__actions">
                    <button type="button" className="btn btn--ghost btn--sm" onClick={() => pilih(b)}>
                      {P.pilih.label}
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
      {hasil !== null && hasil.length === 0 && !memuat && <Kosong pesan={TEKS_AKUMULASI.tanpaHasil} />}
    </Modal>
  )
}
