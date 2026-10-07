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
// Add New (sel 39) -> `FormTambahAkumulasi` (sub-section InputAccumulationCov) menggantikan kartu cari selama dibuka.
//
// Keputusan agent (tiket 46): K-1 Summary (.Note), Accumulation Risk, Risk Accumulation report, Adjustment =
// tahap berikut. K-2 (ralat 04-10-2026, permintaan work owner "hasil nya buat page, 10 list per page") hasil dipaging 10 baris di layar
// - sama dengan grid Pega (pageSize 10); backend tetap mengirim maks. 500 sekaligus. K-7 mengetik di autocomplete mengosongkan nilai hasil
// pilihan sebelumnya (Pega membiarkan ID lama). K-8 Accum. Type hanya tampilan: di Pega terikat
// `SearchAccumulation.Type`, sedangkan RD membaca `SearchAccumulation.AccumulationType` (tidak pernah terisi).

import { useState } from 'react'

import { Field, Gagal, Kosong, Memuat, Modal, Pilih } from '../../../../inti/frontend/components/ui/dasar'
import { jumlahHalaman } from '../../../../inti/frontend/master/DaftarMaster'
import { cariAkumulasi, saranAkumulasi, type BarisAkumulasi, type SaringAkumulasi } from '../api'
import { OPSI_KEYWORD, POPUP_AKUMULASI as P, TAMBAH_AKUMULASI, TEKS_AKUMULASI } from '../labels'
import FormTambahAkumulasi from './FormTambahAkumulasi'
import IsianSaran from './IsianSaran'
import PagerHalaman from './PagerHalaman'

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
  id: '',
  policyNo: '',
  note: '',
  postalCode: '',
  nation: '',
  syariahStatus: '',
  provinceName: '',
  provinceId: '',
  accumType: '',
  city: '',
  cityId: '',
  district: '',
  districtId: '',
  area: '',
  areaZip: '',
  czone: '',
  keyword: '',
}

/** Saringan yang dikirim tombol Filter - Area terpilih menimpa Zip Code (`GetDataAccumulation_act` langkah 4). */
export function keSaring(s: IsianAkumulasi): SaringAkumulasi {
  return {
    id: s.id,
    policyNo: s.policyNo,
    note: s.note,
    postalCode: s.areaZip || s.postalCode,
    syariahStatus: s.syariahStatus,
    provinceId: s.provinceId,
    cityId: s.cityId,
    districtId: s.districtId,
    czone: s.czone,
    keyword: s.keyword,
  }
}

/** Baris hasil per halaman - grid Pega `SearchRiskAccumCov` sel 90 berpaging 10 (permintaan work owner 04-10-2026). */
export const UKURAN_HALAMAN_AKUMULASI = 10

/** Potongan hasil untuk satu halaman (1-based; halaman di luar jangkauan dijepit). */
export function potongHalaman<T>(baris: readonly T[], halaman: number, ukuran = UKURAN_HALAMAN_AKUMULASI): T[] {
  const kini = Math.min(Math.max(1, halaman), jumlahHalaman(baris.length, ukuran))
  return baris.slice((kini - 1) * ukuran, kini * ukuran)
}

/** Jumlah kondisi yang terisi (teks tampil; nilai hasil pilihan tidak dihitung dua kali). */
export function jumlahKondisi(s: IsianAkumulasi): number {
  const tampil: (keyof IsianAkumulasi)[] = [
    'id',
    'policyNo',
    'note',
    'postalCode',
    'nation',
    'provinceName',
    'accumType',
    'city',
    'district',
    'area',
    'czone',
    'keyword',
  ]
  return tampil.filter((k) => s[k].trim() !== '').length
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
  const [s, setS] = useState<IsianAkumulasi>({
    ...ISIAN_KOSONG,
    postalCode: zipRisiko,
  })
  const [hasil, setHasil] = useState<BarisAkumulasi[] | null>(null)
  const [halaman, setHalaman] = useState(1)
  /** Mode Add New (`InputAccumFlag.CARI30 = 1`): form tambah menggantikan kartu cari dan hasil. */
  const [tambah, setTambah] = useState(false)
  const [memuat, setMemuat] = useState(false)
  const [galat, setGalat] = useState<unknown>(null)
  const [pesan, setPesan] = useState('')
  const set = (k: keyof IsianAkumulasi) => (v: string) => setS((x) => ({ ...x, [k]: v }))
  /** Ketik di autocomplete: teks berubah, nilai pilihan lama dikosongkan (K-7). */
  const ketik =
    (teks: keyof IsianAkumulasi, ...nilai: (keyof IsianAkumulasi)[]) =>
    (v: string) =>
      setS((x) => ({
        ...x,
        [teks]: v,
        ...Object.fromEntries(nilai.map((k) => [k, ''])),
      }))

  async function filter() {
    setMemuat(true)
    setGalat(null)
    setPesan('')
    // Langkah 4: Area terpilih menimpa Zip Code yang tampil.
    if (s.areaZip !== '') setS((x) => ({ ...x, postalCode: x.areaZip }))
    try {
      setHasil((await cariAkumulasi(keSaring(s))).baris)
      setHalaman(1)
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

  const aktif = jumlahKondisi(s)

  return (
    <Modal judul={P.judul} onTutup={onTutup} lebar>
      <div className="nbf-akum">
        <div className="nbf-akum__kepala">
          <div>
            <h4 className="nbf-akum__judul">{P.cari.label}</h4>
            <p className="nbf-akum__sub">{TEKS_AKUMULASI.sub}</p>
          </div>
          <span className="nbf-akum__chip">
            {TEKS_AKUMULASI.zipLokasi}
            <strong>{zipRisiko || TEKS_AKUMULASI.zipKosong}</strong>
          </span>
        </div>

        {tambah ? (
          <FormTambahAkumulasi
            zipSaring={s.areaZip || s.postalCode}
            onTutup={() => setTambah(false)}
            onTersimpan={(id, note) => onPilih(id, note)}
          />
        ) : (
          <>
            <form
              className="nbf-akum__kartu"
              onSubmit={(e) => {
                e.preventDefault()
                void filter()
              }}
            >
              <div className="nbf-akum__kartu-kepala">
                <h5 className="nbf-akum__kartu-judul">{P.kondisi}</h5>
                {aktif > 0 && <span className="nbf-akum__jumlah">{TEKS_AKUMULASI.kondisiAktif(aktif)}</span>}
              </div>
              <div className="nbf-akum__grid">
                <Field label={P.accumulationCode.label} value={s.id} onChange={set('id')} />
                <Field label={P.policyNo.label} value={s.policyNo} onChange={set('policyNo')} />
                <div className="nbf-akum__lebar">
                  <Field label={P.road.label} value={s.note} onChange={set('note')} />
                </div>
                <div className="nbf-akum__kolom">
                  <Field label={P.zipCode.label} value={s.postalCode} onChange={set('postalCode')} />
                  <IsianSaran
                    label={P.country.label}
                    value={s.nation}
                    onKetik={ketik('nation', 'syariahStatus')}
                    cari={(q) => saranAkumulasi('nation', q)}
                    onPilih={(x) =>
                      setS((y) => ({
                        ...y,
                        nation: x.label,
                        syariahStatus: x.ekstra ?? '',
                      }))
                    }
                  />
                  <IsianSaran
                    label={P.province.label}
                    value={s.provinceName}
                    onKetik={ketik('provinceName', 'provinceId')}
                    cari={(q) => saranAkumulasi('province', q, s.nation)}
                    onPilih={(x) =>
                      setS((y) => ({
                        ...y,
                        provinceName: x.label,
                        provinceId: x.id,
                      }))
                    }
                  />
                  <IsianSaran
                    label={P.accumType.label}
                    value={s.accumType}
                    onKetik={set('accumType')}
                    cari={(q) => saranAkumulasi('accumtype', q)}
                    onPilih={(x) => setS((y) => ({ ...y, accumType: x.label }))}
                  />
                </div>
                <div className="nbf-akum__kolom">
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
                    onPilih={(x) =>
                      setS((y) => ({
                        ...y,
                        district: x.label,
                        districtId: x.id,
                      }))
                    }
                  />
                  <IsianSaran
                    label={P.area.label}
                    value={s.area}
                    onKetik={ketik('area', 'areaZip')}
                    cari={(q) => saranAkumulasi('area', q, s.district)}
                    onPilih={(x) =>
                      setS((y) => ({
                        ...y,
                        area: x.label,
                        areaZip: x.ekstra ?? '',
                      }))
                    }
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
                <Pilih
                  label={P.keyword.label}
                  value={s.keyword}
                  onChange={set('keyword')}
                  opsi={OPSI_KEYWORD}
                  kosong=""
                />
              </div>
              <div className="nbf-akum__kaki">
                <button
                  type="button"
                  className="btn btn--ghost nbf-akum__kiri"
                  onClick={() => {
                    setPesan('')
                    setTambah(true)
                  }}
                >
                  {TAMBAH_AKUMULASI.tambahBaru.label}
                </button>
                <button
                  type="button"
                  className="btn btn--ghost"
                  onClick={() => {
                    setS(ISIAN_KOSONG)
                    setHasil(null)
                    setPesan('')
                    setGalat(null)
                  }}
                >
                  {P.bersih.label}
                </button>
                <button type="submit" className="btn btn--primary" disabled={memuat}>
                  {P.filter.label}
                </button>
              </div>
            </form>

            {pesan && (
              <div className="alert alert--error" role="alert">
                {pesan}
              </div>
            )}
            <Gagal galat={galat} />

            <section className="nbf-akum__kartu" aria-live="polite">
              <div className="nbf-akum__kartu-kepala">
                <h5 className="nbf-akum__kartu-judul">
                  {hasil !== null && !memuat ? TEKS_AKUMULASI.hasil(hasil.length) : P.kolom[0].label}
                </h5>
              </div>
              {memuat && <Memuat />}
              {hasil === null && !memuat && <p className="nbf-akum__petunjuk">{TEKS_AKUMULASI.petunjukAwal}</p>}
              {hasil !== null && hasil.length === 0 && !memuat && <Kosong pesan={TEKS_AKUMULASI.tanpaHasil} />}
              {hasil !== null && !memuat && hasil.length > 0 && (
                <div className="table-wrap nbf-akum__hasil">
                  <table className="nbf-tabel">
                    <thead>
                      <tr>
                        {P.kolom.map((k) => (
                          <th key={k.sel} scope="col">
                            {k.label}
                          </th>
                        ))}
                        <th scope="col">{TEKS_AKUMULASI.kolomZip}</th>
                        <th scope="col" />
                      </tr>
                    </thead>
                    <tbody>
                      {potongHalaman(hasil, halaman).map((b) => {
                        const cocok = zipDariId(b.id) === zipRisiko
                        return (
                          <tr key={b.id} className={cocok ? undefined : 'nbf-akum__baris-redup'}>
                            <td>
                              <code className="nbf-akum__kode">{b.id}</code>
                            </td>
                            <td>
                              {b.accumulationName && <span className="nbf-akum__tipe">{b.accumulationName}</span>}
                            </td>
                            <td>{b.note}</td>
                            <td>
                              <span
                                className={'nbf-akum__zip ' + (cocok ? 'nbf-akum__zip--ok' : 'nbf-akum__zip--beda')}
                              >
                                {cocok ? TEKS_AKUMULASI.zipSesuai : TEKS_AKUMULASI.zipBerbeda}
                              </span>
                            </td>
                            <td className="table__actions">
                              <button
                                type="button"
                                className={'btn btn--sm ' + (cocok ? 'btn--primary' : 'btn--ghost')}
                                onClick={() => pilih(b)}
                              >
                                {P.pilih.label}
                              </button>
                            </td>
                          </tr>
                        )
                      })}
                    </tbody>
                  </table>
                </div>
              )}
              {hasil !== null && !memuat && hasil.length > UKURAN_HALAMAN_AKUMULASI && (
                <PagerHalaman total={hasil.length} halaman={halaman} ukuran={UKURAN_HALAMAN_AKUMULASI} onPindah={setHalaman} />
              )}
            </section>
          </>
        )}
      </div>
    </Modal>
  )
}
