// Tab Object layar Inward Facultative, keluarga Aneka (tiket 73 P1) - Liability, Bonding, CAR / EAR, MBD, HE, Marine Hull,
// Aviation Hull, Golf, dst. Kabel `GET`/`PUT …/objek-aneka` (sesi c3).
//
// Port berlapis (master-detail expandPane seperti tab Object FIRE):
// 1. Grid lokasi `.OfferFacIn.LocationList` - `Section\InputDtlObject_FacIn.xml` kontainer L36463 (`IsAnekaNotBonding || IsHE
//    || IsGolfInsurance || IsGrowingTrees`) / L22724 (`IsBonding || IsCustomBonds`); kolom No / Top Risk / Location; Add di
//    kepala, Delete per baris (NB: deleteRow).
// 2. Form lokasi `InputDtlObjectLocation_FacIn.xml`: Choose / Clear Risk Address (L4578 / L5520, aturan FIRE), Top Risk (L6301
//    syarat `tampilanAneka.topRisk`), blok alamat (L7166-11143, sama dengan FIRE), grid Object List Golf (L11522
//    `isGolfInsurance`), grid Occupation (L15288).
// 3. Detail okupasi `InputOkupasiAneka_FacIn.xml`: tombol Choose Occupation (L1224 → HARNESS `ChooseOccupation`, okupasi
//    TYPE ANEKA; id → OccupationId = OCCUPATION.ID) + grid objek varian K2 / K3 (HE) / K4 (Growing Trees / Yield Shortfall).
// 4. Form objek `FormObjekAneka` (kontainer C1-C9).
//
// Save: seluruh lokasi dikirim `PUT …/objek-aneka`; 422 = pesan `ProtectionObjectAneka_Act` (dipotong " \n" di server).
//
// Keputusan agent (tiket 73 P1):
// - B-1 Kolom grid objek yang di Pega dapat diisi langsung di baris (Section, Limit of Liability, Year, Quantity, Unit,
//   Chassis, Engine; Object Name / TSI Old di grid Golf) ditampilkan sebagai teks; isiannya di form objek yang dibuka di
//   bawah baris (properti yang sama) - satu tempat isian per medan.
// - B-2 Tombol Delete okupasi / objek / lokasi langsung membuang baris (NB, `deleteRow`); jalur EDM `SetFlagDeleteEDM_Act`
//   bukan di sini.

import { useEffect, useState } from 'react'

import { Field, Gagal, Halaman, Kosong, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import {
  ambilObjekAneka,
  cariOkupasiAneka,
  daftarMataUang,
  simpanObjekAneka,
  type BarisOkupasiAneka,
  type LokasiAneka,
  type ObjekAneka,
  type OkupasiAneka,
  type TampilanAneka,
} from '../api'
import {
  GRID_GOLF_ANEKA,
  GRID_LOKASI_ANEKA,
  GRID_OBJEK,
  GRID_OBJEK_ANEKA,
  GRID_OKUPASI_ANEKA,
  OBJECT_ADDRESS as A,
  OPSI_SECTION_ANEKA,
  POPUP_OKUPASI_ANEKA as P,
  SIMPAN_OBJEK,
  TEKS_ANEKA,
  TEKS_FORM_OPPORTUNITY,
  TEKS_OKUPASI,
} from '../labels'
import {
  kolomGridObjek,
  kosongkanRiskAneka,
  lokasiAnekaBaru,
  nilaiKolom,
  objekAnekaBaru,
  rapikanLokasiAneka,
  saringDariAneka,
  tampilOkupasi,
  terapkanAlamatBaruAneka,
  terapkanRiskAneka,
} from './aneka'
import { pesanProteksiDari } from './BlokEmailType'
import FormObjekAneka from './FormObjekAneka'
import { keTampilUang } from './IsianUang'
import PopupRiskAddress from './PopupRiskAddress'

const JEDA_MS = 400

/** Tampilan satu sel grid objek (Section = label opsi; uang berpemisah ribuan). */
function teksKolom(o: ObjekAneka, k: string): string {
  const v = nilaiKolom(o, k)
  if (k === 'section') return OPSI_SECTION_ANEKA.find((s) => s.value === v)?.label ?? v
  if (k === 'tsi' || k === 'tsiOld' || k === 'limitOfLiability') return keTampilUang(v)
  return v
}

/** Sel grid: Object Name kosong (objek baru) ditandai supaya baris tidak tampak kosong. */
function SelObjek({ o, k }: { o: ObjekAneka; k: string }) {
  const v = teksKolom(o, k)
  return k === 'objectName' && v.trim() === '' ? <span className="nbf-aneka-grid__kosong">{TEKS_ANEKA.belumDiisi}</span> : <>{v}</>
}

function Tampil({ label, nilai }: { label: string; nilai: string }) {
  return (
    <div className="field">
      <span className="field__label">{label}</span>
      <div className="nbf-inward__teks">{nilai}</div>
    </div>
  )
}

/** Tombol buka / tutup baris master-detail (pola tab Object FIRE). */
function TombolBuka({ buka, onClick }: { buka: boolean; onClick: () => void }) {
  return (
    <button
      type="button"
      className={'btn btn--ghost btn--sm nbf-buka' + (buka ? ' nbf-buka--terbuka' : '')}
      aria-label={TEKS_ANEKA.bukaBaris}
      aria-expanded={buka}
      onClick={onClick}
    >
      ▸
    </button>
  )
}

function TombolTambah({ onClick }: { onClick: () => void }) {
  return (
    <button type="button" className="btn nbf-tema__tombol btn--ghost btn--sm" onClick={onClick}>
      {GRID_OBJEK.tambah}
    </button>
  )
}

function TombolHapus({ onClick }: { onClick: () => void }) {
  return (
    <button type="button" className="btn nbf-tema__tombol btn--ghost btn--sm" onClick={onClick}>
      {GRID_OBJEK.hapus}
    </button>
  )
}

/** Popup Choose Occupation Aneka - HARNESS `ChooseOccupation.xml` (Search, kolom ID / Occupation Name, Choose; 10 per halaman). */
export function PopupOkupasiAneka({
  caseId,
  onTutup,
  onPilih,
}: {
  caseId: string
  onTutup: () => void
  onPilih: (b: BarisOkupasiAneka) => void
}) {
  const [kotak, setKotak] = useState('')
  const [hasil, setHasil] = useState<BarisOkupasiAneka[] | null>(null)
  const [halaman, setHalaman] = useState(1)
  const [galat, setGalat] = useState<unknown>(null)

  useEffect(() => {
    let batal = false
    const jadwal = window.setTimeout(() => {
      setGalat(null)
      cariOkupasiAneka(caseId, kotak).then(
        (h) => {
          if (batal) return
          setHasil(h.baris)
          setHalaman(1)
        },
        (err: unknown) => {
          if (batal) return
          setHasil(null)
          setGalat(err)
        },
      )
    }, kotak === '' ? 0 : JEDA_MS)
    return () => {
      batal = true
      window.clearTimeout(jadwal)
    }
  }, [caseId, kotak])

  const u = P.ukuran
  const tampil = (hasil ?? []).slice((halaman - 1) * u, halaman * u)
  return (
    <Modal judul={P.judul} onTutup={onTutup} lebar>
      <div className="nbf-popup__kotak">
        <Field label={P.cari} value={kotak} onChange={setKotak} />
      </div>
      {hasil && hasil.length > u && <Halaman halaman={halaman} ukuran={u} total={hasil.length} onPindah={setHalaman} />}
      <div className="nbf-cov-wrap">
        <table className="nbf-tabel">
          <thead>
            <tr>
              {P.kolom.map((k) => (
                <th key={k} scope="col">
                  {k}
                </th>
              ))}
              <th scope="col" />
            </tr>
          </thead>
          {tampil.length > 0 && (
            <tbody>
              {tampil.map((b) => (
                <tr key={b.id}>
                  <td>{b.id}</td>
                  <td>{b.name}</td>
                  <td className="table__actions">
                    <button type="button" className="btn nbf-tema__tombol btn--ghost btn--sm" onClick={() => onPilih(b)}>
                      {P.pilih}
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          )}
        </table>
      </div>
      {hasil === null && galat === null && <Memuat />}
      {hasil !== null && hasil.length === 0 && <Kosong pesan={TEKS_OKUPASI.tanpaHasil} />}
      <Gagal galat={galat} />
    </Modal>
  )
}

/** Grid objek (varian K2 / K3 / K4, atau grid Golf) + form objek di baris yang dibuka. */
function GridObjek({
  objek,
  ubah,
  kolom,
  label,
  tampilan,
  mataUang,
}: {
  objek: ObjekAneka[]
  ubah: (objek: ObjekAneka[]) => void
  kolom: string[]
  label: (k: string) => string
  tampilan: TampilanAneka
  mataUang: string[]
}) {
  const [buka, setBuka] = useState<number | null>(null)
  const lebar = kolom.length + 2
  return (
    <div className="nbf-cov-wrap">
      <table className="nbf-tabel nbf-aneka-grid nbf-aneka-grid--aksi">
        <thead>
          <tr>
            <th scope="col" />
            {kolom.map((k) => (
              <th key={k} scope="col">
                {label(k)}
              </th>
            ))}
            <th scope="col" className="table__actions">
              <TombolTambah
                onClick={() => {
                  ubah([...objek, objekAnekaBaru()])
                  setBuka(objek.length)
                }}
              />
            </th>
          </tr>
        </thead>
        <tbody>
          {objek.length === 0 && (
            <tr>
              <td colSpan={lebar}>{TEKS_ANEKA.kosongObjek}</td>
            </tr>
          )}
          {objek.map((o, i) => [
            <tr key={i} className={buka === i ? 'nbf-baris--terbuka' : undefined}>
              <td>
                <TombolBuka buka={buka === i} onClick={() => setBuka(buka === i ? null : i)} />
              </td>
              {kolom.map((k) => (
                <td key={k}>
                  <SelObjek o={o} k={k} />
                </td>
              ))}
              <td className="table__actions">
                <TombolHapus
                  onClick={() => {
                    ubah(objek.filter((_, j) => j !== i))
                    setBuka(null)
                  }}
                />
              </td>
            </tr>,
            buka === i && (
              <tr key={`${i}-detail`} className="nbf-objek__detail">
                <td colSpan={lebar}>
                  <div className="nbf-lapis">
                    <FormObjekAneka
                      o={o}
                      ubah={(baru) => ubah(objek.map((x, j) => (j === i ? baru : x)))}
                      tampilan={tampilan}
                      mataUang={mataUang}
                    />
                  </div>
                </td>
              </tr>
            ),
          ])}
        </tbody>
      </table>
    </div>
  )
}

/** Grid Occupation di form lokasi + detail okupasi (Choose Occupation, grid objek). */
function GridOkupasi({
  caseId,
  okupasi,
  ubah,
  tampilan,
  mataUang,
}: {
  caseId: string
  okupasi: OkupasiAneka[]
  ubah: (okupasi: OkupasiAneka[]) => void
  tampilan: TampilanAneka
  mataUang: string[]
}) {
  const [buka, setBuka] = useState<number | null>(null)
  const [pilih, setPilih] = useState<number | null>(null)
  const ganti = (i: number, baru: OkupasiAneka) => ubah(okupasi.map((x, j) => (j === i ? baru : x)))
  const kolom = kolomGridObjek(tampilan)
  return (
    <div className="nbf-cov-wrap">
      <table className="nbf-tabel nbf-aneka-grid nbf-aneka-grid--aksi">
        <thead>
          <tr>
            <th scope="col" />
            {GRID_OKUPASI_ANEKA.map((k) => (
              <th key={k} scope="col">
                {k}
              </th>
            ))}
            <th scope="col" className="table__actions">
              <TombolTambah
                onClick={() => {
                  ubah([...okupasi, { occupationId: '', occupationName: '', objek: [] }])
                  setBuka(okupasi.length)
                }}
              />
            </th>
          </tr>
        </thead>
        <tbody>
          {okupasi.length === 0 && (
            <tr>
              <td colSpan={4}>{TEKS_ANEKA.kosongOkupasi}</td>
            </tr>
          )}
          {okupasi.map((k, i) => [
            <tr key={i} className={buka === i ? 'nbf-baris--terbuka' : undefined}>
              <td>
                <TombolBuka buka={buka === i} onClick={() => setBuka(buka === i ? null : i)} />
              </td>
              <td>{k.occupationId}</td>
              <td>{k.occupationName}</td>
              <td className="table__actions">
                <TombolHapus
                  onClick={() => {
                    ubah(okupasi.filter((_, j) => j !== i))
                    setBuka(null)
                  }}
                />
              </td>
            </tr>,
            buka === i && (
              <tr key={`${i}-detail`} className="nbf-objek__detail">
                <td colSpan={4}>
                  <div className="nbf-lapis">
                    <div className="nbf-opp__tombol">
                      <button type="button" className="btn nbf-tema__tombol btn--sm" onClick={() => setPilih(i)}>
                        {P.judul}
                      </button>
                    </div>
                    <GridObjek
                      objek={k.objek}
                      ubah={(objek) => ganti(i, { ...k, objek })}
                      kolom={kolom}
                      label={(c) => GRID_OBJEK_ANEKA[c as keyof typeof GRID_OBJEK_ANEKA] ?? c}
                      tampilan={tampilan}
                      mataUang={mataUang}
                    />
                  </div>
                </td>
              </tr>
            ),
          ])}
        </tbody>
      </table>
      {pilih !== null && (
        <PopupOkupasiAneka
          caseId={caseId}
          onTutup={() => setPilih(null)}
          onPilih={(b) => {
            const k = okupasi[pilih]
            if (k) ganti(pilih, { ...k, occupationId: b.id, occupationName: b.name })
            setPilih(null)
          }}
        />
      )}
    </div>
  )
}

/** Form lokasi - `InputDtlObjectLocation_FacIn.xml`. */
function FormLokasi({
  caseId,
  l,
  ubah,
  tampilan,
  aneka,
  mataUang,
}: {
  caseId: string
  l: LokasiAneka
  ubah: (l: LokasiAneka) => void
  tampilan: TampilanAneka
  aneka: boolean
  mataUang: string[]
}) {
  const [pilihRisk, setPilihRisk] = useState(false)
  return (
    <div className="nbf-objek__isi nbf-ringkas nbf-alamat">
      <section className="nbf-alamat__blok">
        <div className="nbf-alamat__kepala">
          <div className="nbf-opp__tombol">
            <button type="button" className="btn nbf-tema__tombol btn--sm" onClick={() => setPilihRisk(true)}>
              {A.chooseRisk.label}
            </button>
            <button type="button" className="btn nbf-tema__tombol btn--sm" onClick={() => ubah(kosongkanRiskAneka(l))}>
              {A.clearRisk.label}
            </button>
          </div>
          <h5 className="nbf-objek__judul">{A.judulRisk.label}</h5>
        </div>
        {pilihRisk && (
          <PopupRiskAddress
            awal={saringDariAneka(l)}
            onTutup={() => setPilihRisk(false)}
            onPilih={(b) => {
              ubah(terapkanRiskAneka(l, b))
              setPilihRisk(false)
            }}
            onBaru={(a, id) => {
              ubah(terapkanAlamatBaruAneka(l, a, id))
              setPilihRisk(false)
            }}
          />
        )}
        {tampilan.topRisk && (
          <label className="nbf-inward__pilihan">
            <input type="checkbox" checked={l.isTopRisk} onChange={(e) => ubah({ ...l, isTopRisk: e.target.checked })} /> {A.topRisk.label}
          </label>
        )}
        <div className="nbf-opp__kolom">
          <div className="nbf-opp__tumpuk">
            <Tampil label={A.type.label} nilai={l.roadType} />
            <Tampil label={A.address.label} nilai={l.roadName} />
            <Field label={A.buildingNo.label} value={l.buildingNo} onChange={(v) => ubah({ ...l, buildingNo: v })} />
            <Tampil label={A.zipCode.label} nilai={l.zipCode} />
            <Tampil label={A.country.label} nilai={l.country} />
          </div>
          <div className="nbf-opp__tumpuk">
            <Tampil label={A.territory.label} nilai={l.territory} />
            <Tampil label={A.city.label} nilai={l.city} />
            <Tampil label={A.district.label} nilai={l.district} />
            <Tampil label={A.province.label} nilai={l.province} />
            <Tampil label={A.riskAddressId.label} nilai={l.riskAddressId} />
          </div>
        </div>
        <div className="nbf-alamat__lokasi">
          <Tampil label={A.riskLocation.label} nilai={l.riskLocation} />
        </div>
      </section>

      {tampilan.golf && (
        <section className="nbf-alamat__blok">
          <GridObjek
            objek={l.objekGolf}
            ubah={(objekGolf) => ubah({ ...l, objekGolf })}
            kolom={['objectName', 'currency', 'tsi']}
            label={(k) => GRID_GOLF_ANEKA[['objectName', 'currency', 'tsi'].indexOf(k)] ?? k}
            tampilan={tampilan}
            mataUang={mataUang}
          />
        </section>
      )}

      {tampilOkupasi(aneka, tampilan) && (
        <section className="nbf-alamat__blok">
          <GridOkupasi caseId={caseId} okupasi={l.okupasi} ubah={(okupasi) => ubah({ ...l, okupasi })} tampilan={tampilan} mataUang={mataUang} />
        </section>
      )}
    </div>
  )
}

export default function TabObjectAneka({ caseId, tampilan, aneka }: { caseId: string; tampilan: TampilanAneka; aneka: boolean }) {
  const [lokasi, setLokasi] = useState<LokasiAneka[] | null>(null)
  // % Pro Rate tingkat case (migrasi 255): dikirim balik apa adanya - PUT menulis nilai ini (kosong = NULL).
  const [proRate, setProRate] = useState('')
  const [mataUang, setMataUang] = useState<string[]>([])
  const [buka, setBuka] = useState<number | null>(null)
  const [halaman, setHalaman] = useState(1)
  const [galat, setGalat] = useState<unknown>(null)
  const [proteksi, setProteksi] = useState<string[]>([])
  const [menyimpan, setMenyimpan] = useState(false)
  const [tersimpan, setTersimpan] = useState(false)

  useEffect(() => {
    let batal = false
    ambilObjekAneka(caseId).then(
      (h) => {
        if (batal) return
        setLokasi(h.lokasi.map(rapikanLokasiAneka))
        setProRate(h.proRatePercent ?? '')
      },
      (err: unknown) => {
        if (!batal) setGalat(err)
      },
    )
    daftarMataUang().then(
      (h) => {
        if (!batal) setMataUang(h.baris)
      },
      () => undefined,
    )
    return () => {
      batal = true
    }
  }, [caseId])

  if (lokasi === null) return galat !== null ? <Gagal galat={galat} /> : <Memuat />

  const daftar = lokasi
  function ganti(baru: LokasiAneka[]) {
    setLokasi(baru)
    setTersimpan(false)
  }

  async function simpan() {
    setMenyimpan(true)
    setGalat(null)
    setProteksi([])
    setTersimpan(false)
    try {
      const h = await simpanObjekAneka(caseId, daftar, proRate)
      setLokasi(h.lokasi.map(rapikanLokasiAneka))
      setProRate(h.proRatePercent ?? '')
      setTersimpan(true)
    } catch (err) {
      const p = pesanProteksiDari(err)
      if (p !== null) setProteksi(p)
      else setGalat(err)
    } finally {
      setMenyimpan(false)
    }
  }

  const u = GRID_OBJEK.ukuran
  const jumlahHalaman = Math.max(1, Math.ceil(daftar.length / u))
  const hal = Math.min(halaman, jumlahHalaman)
  const awal = (hal - 1) * u
  const tampil = daftar.slice(awal, awal + u)

  return (
    <div className="nbf-objek nbf-obj-tab">
      <Gagal galat={galat} />
      {proteksi.length > 0 && (
        <div className="alert alert--error" role="alert">
          {TEKS_ANEKA.ditolak}
          <ul>
            {proteksi.map((p, i) => (
              <li key={i}>{p}</li>
            ))}
          </ul>
        </div>
      )}
      {tersimpan && <div className="alert alert--ok">{TEKS_ANEKA.tersimpan}</div>}
      <div className="nbf-cov-wrap">
        <table className="nbf-tabel nbf-aneka-grid nbf-aneka-grid--aksi nbf-aneka-grid--lokasi">
          <thead>
            <tr>
              <th scope="col" />
              {GRID_LOKASI_ANEKA.map((k) => (
                <th key={k} scope="col">
                  {k}
                </th>
              ))}
              <th scope="col" className="table__actions">
                <TombolTambah
                  onClick={() => {
                    ganti([...daftar, lokasiAnekaBaru(String(daftar.length + 1))])
                    setBuka(daftar.length)
                    setHalaman(Math.ceil((daftar.length + 1) / u))
                  }}
                />
              </th>
            </tr>
          </thead>
          <tbody>
            {daftar.length === 0 && (
              <tr>
                <td colSpan={5}>{TEKS_ANEKA.kosongLokasi}</td>
              </tr>
            )}
            {tampil.map((l, n) => {
              const i = awal + n
              return [
                <tr key={i} className={buka === i ? 'nbf-baris--terbuka' : undefined}>
                  <td>
                    <TombolBuka buka={buka === i} onClick={() => setBuka(buka === i ? null : i)} />
                  </td>
                  <td>{l.objectNo}</td>
                  <td>
                    <input type="checkbox" checked={l.isTopRisk} disabled aria-label={GRID_LOKASI_ANEKA[1]} />
                  </td>
                  <td>{l.riskLocation}</td>
                  <td className="table__actions">
                    <TombolHapus
                      onClick={() => {
                        ganti(daftar.filter((_, j) => j !== i))
                        setBuka(null)
                      }}
                    />
                  </td>
                </tr>,
                buka === i && (
                  <tr key={`${i}-detail`} className="nbf-objek__detail">
                    <td colSpan={5}>
                      <div className="nbf-lapis nbf-lapis--objek">
                        <FormLokasi
                          caseId={caseId}
                          l={l}
                          ubah={(baru) => ganti(daftar.map((x, j) => (j === i ? baru : x)))}
                          tampilan={tampilan}
                          aneka={aneka}
                          mataUang={mataUang}
                        />
                      </div>
                    </td>
                  </tr>
                ),
              ]
            })}
          </tbody>
        </table>
      </div>
      {daftar.length > u && <Halaman halaman={hal} ukuran={u} total={daftar.length} onPindah={setHalaman} />}
      <div className="nbf-objek__kaki">
        <button type="button" className="btn nbf-tema__tombol nbf-tema__tombol--utama btn--primary" onClick={() => void simpan()} disabled={menyimpan}>
          {menyimpan ? TEKS_FORM_OPPORTUNITY.menyimpan : SIMPAN_OBJEK.label}
        </button>
      </div>
    </div>
  )
}

