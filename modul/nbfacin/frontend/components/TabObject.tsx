// Tab Object layar Inward Facultative, kasus FIRE (tiket 35, tahap 1).
//
// Port `NB FacIn\Section\ObjectList.xml` (grid `.LocationList`, master-detail, 10 per halaman) + flow action
// `Property_FlowAction` → section `Property` (tujuh sub-tab) → `ObjectDetails` (sub-tab Object Address) + tombol
// Save Dtl sel 25. Tampilan = tangkapan layar work owner 03-10-2026 (gambar tidak disalin: memuat data nasabah).
//
// - Tambah (`AddLocation_act`): baris baru di akhir, Object No. = nomor urut baris itu (`pxListSubscript`).
// - Hapus: NB = StatusBusiness New Business → baris dibuang (`deleteRow`); jalur EDM (`SetFlagDeleteEDM_Act`) bukan
//   di sini.
// - Object Type berubah: Object Name = Object Type, kecuali `Others` → Object Name dikosongkan dan medannya tampil
//   untuk diisi (sel 13 visible `ObjectType = 'Others'`; `SetValueOnObjectName_Act`, keputusan agent G-2).
// - Number of Floor < 0 → "Floor number can't be minus" (`SetErrorMessageFloorNumber_Act`).
// - Save: seluruh daftar dikirim `PUT …/objek` (Object Type wajib per baris, sel 6).
//
// Choose Risk Address (tiket 36): popup `PopupRiskAddress`; Pilih mengisi Risk Address menurut `SetRiskIdDT_FacIn`
// (`terapkanRisk`). Clear Risk Address mengosongkan medan `ClearRiskLocation_act` (`kosongkanRisk`) - Building No.
// tidak ikut. Keduanya tersimpan lewat Save.
//
// Surrounding Risk (tiket 38): `SubTabSekitar`; Distance minus ikut menahan Save.
//
// Keputusan agent (tiket 35): G-1 (selesai di tiket 36); medan Risk Address tampil-saja kecuali Building No. G-3 sub-tab selain Object Address =
// tahap berikut (`BelumTersedia`). Daftar Roof / Wall / Floor Type = aturan properti Pega (`DDL\RoofType.xml` dst.,
// lihat `OPSI_ROOF_TYPE` di labels.ts); nilai tersimpan = kode Pega; objek baru berawal "Lain-lain" (14 / 9 /
// "KELAS I").

import { useEffect, useRef, useState } from 'react'

import { BelumTersedia, Field, Gagal, Halaman, Pilih, StripTab, type Opsi } from '../../../../inti/frontend/components/ui/dasar'
import { ambilObjek, simpanObjek, type AlamatBaru, type BarisRisk, type ObjekFire, type SaringRisk } from '../api'
import PopupRiskAddress from './PopupRiskAddress'
import SubTabItem, { adaGalatItem } from './SubTabItem'
import SubTabFEA, { adaGalatFEA } from './SubTabFEA'
import { SubTabKerugian, SubTabKlaimInternal, adaGalatKerugian, lossRatioKosong } from './SubTabKerugian'
import SubTabOkupasi from './SubTabOkupasi'
import SubTabSekitar, { OWNERSHIP_AWAL, adaJarakMinus, sekitarKosong } from './SubTabSekitar'
import {
  AWAL_BANGUNAN,
  BANGUNAN_KOSONG,
  GRID_OBJEK,
  OBJECT_ADDRESS as A,
  OBJECT_TYPE_LAINNYA,
  OPSI_FLOOR_TYPE,
  OPSI_OBJECT_TYPE,
  OPSI_ROOF_TYPE,
  OPSI_WALL_TYPE,
  SIMPAN_OBJEK,
  SUBTAB_OBJEK,
  TEKS_FORM_OPPORTUNITY,
  TEKS_INWARD,
  TEKS_OBJEK,
} from '../labels'

type SubTab = (typeof SUBTAB_OBJEK)[number]

/** Salinan peta tanpa satu kunci. */
function tanpa<T>(peta: Record<number, T>, kunci: number): Record<number, T> {
  const sisa = { ...peta }
  delete sisa[kunci]
  return sisa
}

/** Baris layar: data + kunci React yang tetap walau baris lain dihapus. */
interface BarisLayar {
  kunci: number
  data: ObjekFire
}

/** Objek kosong untuk Tambah. Material Damage awal tercentang (sel 11, nilai awal `true`); Building Construction awal Pega. */
export function objekBaru(objectNo: string): ObjekFire {
  return {
    objectNo, objectType: '', objectName: '', isMaterialDamage: true, isTopRisk: false,
    roadType: '', roadName: '', buildingNo: '', zipCode: '', country: '', riskLocation: '', territory: '',
    city: '', district: '', province: '', riskAddressId: '',
    numberOfFloor: '', ...AWAL_BANGUNAN, partitionType: '', supportWallType: '', otherType: '',
    ownership: OWNERSHIP_AWAL, isProductionProcess: false, isHotWorkProcess: false, isFlammableItem: false,
    surroundingRisk: sekitarKosong(),
    items: [],
    occupations: [],
    fea: [],
    lossRecords: [],
    lossRatio: lossRatioKosong(),
    internalLossRecords: [],
  }
}

/**
 * Objek dari server dilengkapi medan yang belum dikirim (mis. backend sebelum tiket 38 tanpa `surroundingRisk`),
 * supaya layar tidak runtuh pada data lama.
 */
export function rapikanObjek(d: Partial<ObjekFire>): ObjekFire {
  const awal = objekBaru(d.objectNo ?? '')
  const s = (d.surroundingRisk ?? {}) as Partial<ObjekFire['surroundingRisk']>
  const k = sekitarKosong()
  return {
    ...awal,
    ...d,
    items: d.items ?? [],
    occupations: d.occupations ?? [],
    fea: d.fea ?? [],
    lossRecords: d.lossRecords ?? [],
    lossRatio: { ...lossRatioKosong(), ...d.lossRatio },
    internalLossRecords: d.internalLossRecords ?? [],
    surroundingRisk: {
      ...k,
      ...s,
      front: { ...k.front, ...s.front },
      left: { ...k.left, ...s.left },
      back: { ...k.back, ...s.back },
      right: { ...k.right, ...s.right },
    },
  }
}

/** Object Type berubah → Object Name ikut (G-2). */
export function gantiObjectType(o: ObjekFire, objectType: string): ObjekFire {
  return { ...o, objectType, objectName: objectType === OBJECT_TYPE_LAINNYA ? '' : objectType }
}

/**
 * Pilih di popup Risk Address - `DataTransform\SetRiskIdDT_FacIn.xml`: AlmRiskID = ID, RoadType = Title,
 * RoadName = Address, ASMRW = TerritoryName, ASMDistrict, ASMCity, Province, Country = NationName, ASMZipCode =
 * PostalCode, ASMAddress = Title + " " + Address + "," + Territory + "," + District + "," + City + "," + Province +
 * "," + Nation; BuildingNo dikosongkan.
 */
export function terapkanRisk(o: ObjekFire, b: BarisRisk): ObjekFire {
  return {
    ...o,
    riskAddressId: b.id,
    roadType: b.title,
    roadName: b.address,
    territory: b.territoryName,
    district: b.districtName,
    city: b.cityName,
    province: b.provinceName,
    country: b.nationName,
    zipCode: b.postalCode,
    riskLocation: `${b.title} ${b.address},${b.territoryName},${b.districtName},${b.cityName},${b.provinceName},${b.nationName}`,
    buildingNo: '',
  }
}

/**
 * Alamat baru dari popup Add (tiket 37) - langkah 9 `Activity\SaveRiskAddress_Act.xml`: RoadName, RoadType = Title,
 * Risk Location dirangkai seperti `SetRiskIdDT_FacIn`, Zip Code, Country, Province, City, District, Territory;
 * Building No. TIDAK disentuh. Risk Address ID = ID baru dari backend (keputusan work owner 03-10-2026; Pega
 * membiarkannya kosong).
 */
export function terapkanAlamatBaru(o: ObjekFire, a: AlamatBaru, id: string): ObjekFire {
  const hasil = terapkanRisk(o, {
    id, title: a.title, address: a.address, nationName: a.nationName, provinceName: a.provinceName,
    cityName: a.cityName, districtName: a.districtName, territoryName: a.territoryName, postalCode: a.postalCode,
  })
  return { ...hasil, buildingNo: o.buildingNo }
}

/** Clear Risk Address - `Activity\ClearRiskLocation_act.xml` (Building No. tidak ikut). */
export function kosongkanRisk(o: ObjekFire): ObjekFire {
  return {
    ...o,
    roadType: '', roadName: '', zipCode: '', province: '', country: '', territory: '', city: '', district: '',
    riskLocation: '', riskAddressId: '',
  }
}

/** Saringan awal popup = medan objek yang diikat sel 78-84 (H-1). */
export function saringDari(o: ObjekFire): SaringRisk {
  return {
    address: o.roadName, zipCode: o.zipCode, country: o.country, province: o.province, city: o.city,
    district: o.district, territory: o.territory,
  }
}

/** Number of Floor bernilai minus (`NumberOfFloor < 0`). */
export function lantaiMinus(numberOfFloor: string): boolean {
  const n = Number(numberOfFloor)
  return numberOfFloor.trim() !== '' && Number.isFinite(n) && n < 0
}

const OPSI_TYPE: Opsi[] = OPSI_OBJECT_TYPE.map((v) => ({ value: v, label: v }))

/** Daftar + nilai tersimpan yang tidak dikenal daftar (data lama), supaya tidak hilang diam-diam saat Save. */
function denganTersimpan(daftar: readonly Opsi[], nilai: string): Opsi[] {
  return nilai === '' || daftar.some((o) => o.value === nilai) ? [...daftar] : [...daftar, { value: nilai, label: nilai }]
}

/** Medan tampil-saja (label kiri, teks kanan). */
function Tampil({ label, nilai }: { label: string; nilai: string }) {
  return (
    <div className="field">
      <span className="field__label">{label}</span>
      <div className="nbf-inward__teks">{nilai}</div>
    </div>
  )
}

function Centang({ label, checked, onChange }: { label: string; checked: boolean; onChange: (v: boolean) => void }) {
  return (
    <label className="nbf-inward__pilihan">
      <input type="checkbox" checked={checked} onChange={(e) => onChange(e.target.checked)} /> {label}
    </label>
  )
}

/** Sub-tab Object Address (`ObjectDetails`). */
function ObjectAddress({ o, ubah, galatType }: { o: ObjekFire; ubah: (o: ObjekFire) => void; galatType: boolean }) {
  const set = (k: keyof ObjekFire) => (v: string) => ubah({ ...o, [k]: v })
  const [pilihRisk, setPilihRisk] = useState(false)
  return (
    // Dirapikan work owner 05-10-2026: tiga kartu (objek, Risk Address, bangunan), dua kolom sama lebar, isian ringkas.
    <div className="nbf-objek__isi nbf-ringkas nbf-alamat">
      <section className="nbf-alamat__blok">
        <div className="nbf-opp__kolom">
          <div className="nbf-opp__tumpuk">
            <Field label={A.objectNo.label} value={o.objectNo} onChange={set('objectNo')} />
            <Pilih
              label={A.objectType.label}
              value={o.objectType}
              onChange={(v) => ubah(gantiObjectType(o, v))}
              opsi={OPSI_TYPE}
              kosong={A.objectTypeKosong.label}
              required
              error={galatType ? TEKS_OBJEK.typeWajib : undefined}
            />
            {o.objectType === OBJECT_TYPE_LAINNYA && <Field label={A.objectName.label} value={o.objectName} onChange={set('objectName')} />}
          </div>
          <div className="nbf-opp__tumpuk">
            <Centang label={A.materialDamage.label} checked={o.isMaterialDamage} onChange={(v) => ubah({ ...o, isMaterialDamage: v })} />
            <Centang label={A.topRisk.label} checked={o.isTopRisk} onChange={(v) => ubah({ ...o, isTopRisk: v })} />
          </div>
        </div>
      </section>

      <section className="nbf-alamat__blok">
        {/* Tombol Choose / Clear di atas judul Risk Address (work owner 05-10-2026). */}
        <div className="nbf-alamat__kepala">
          <div className="nbf-opp__tombol">
            <button type="button" className="btn btn--sm" onClick={() => setPilihRisk(true)}>
              {A.chooseRisk.label}
            </button>
            <button type="button" className="btn btn--sm" onClick={() => ubah(kosongkanRisk(o))}>
              {A.clearRisk.label}
            </button>
          </div>
          <h5 className="nbf-objek__judul">{A.judulRisk.label}</h5>
        </div>
        {pilihRisk && (
          <PopupRiskAddress
            awal={saringDari(o)}
            onTutup={() => setPilihRisk(false)}
            onPilih={(b) => {
              ubah(terapkanRisk(o, b))
              setPilihRisk(false)
            }}
            onBaru={(a, id) => {
              ubah(terapkanAlamatBaru(o, a, id))
              setPilihRisk(false)
            }}
          />
        )}
        <div className="nbf-opp__kolom">
          <div className="nbf-opp__tumpuk">
            <Tampil label={A.type.label} nilai={o.roadType} />
            <Tampil label={A.address.label} nilai={o.roadName} />
            <Field label={A.buildingNo.label} value={o.buildingNo} onChange={set('buildingNo')} />
            <Tampil label={A.zipCode.label} nilai={o.zipCode} />
            <Tampil label={A.country.label} nilai={o.country} />
          </div>
          <div className="nbf-opp__tumpuk">
            <Tampil label={A.territory.label} nilai={o.territory} />
            <Tampil label={A.city.label} nilai={o.city} />
            <Tampil label={A.district.label} nilai={o.district} />
            <Tampil label={A.province.label} nilai={o.province} />
            <Tampil label={A.riskAddressId.label} nilai={o.riskAddressId} />
          </div>
        </div>
        {/* Risk Location memanjang sampai tepi kanan kotak Risk Address ID (work owner 05-10-2026). */}
        <div className="nbf-alamat__lokasi">
          <Tampil label={A.riskLocation.label} nilai={o.riskLocation} />
        </div>
      </section>

      <section className="nbf-alamat__blok">
        <h5 className="nbf-objek__judul">{A.judulBangunan.label}</h5>
        <div className="nbf-opp__kolom">
          <div className="nbf-opp__tumpuk">
            <Field
              label={A.numberOfFloor.label}
              type="number"
              value={o.numberOfFloor}
              onChange={set('numberOfFloor')}
              error={lantaiMinus(o.numberOfFloor) ? TEKS_OBJEK.lantaiMinus : undefined}
            />
            <Pilih label={A.roofType.label} value={o.roofType} onChange={set('roofType')} opsi={denganTersimpan(OPSI_ROOF_TYPE, o.roofType)} kosong={BANGUNAN_KOSONG} />
            <Pilih label={A.wallType.label} value={o.wallType} onChange={set('wallType')} opsi={denganTersimpan(OPSI_WALL_TYPE, o.wallType)} kosong={BANGUNAN_KOSONG} />
            <Pilih label={A.floorType.label} value={o.floorType} onChange={set('floorType')} opsi={denganTersimpan(OPSI_FLOOR_TYPE, o.floorType)} kosong={BANGUNAN_KOSONG} />
          </div>
          <div className="nbf-opp__tumpuk">
            <Field label={A.partitionType.label} value={o.partitionType} onChange={set('partitionType')} />
            <Field label={A.supportWallType.label} value={o.supportWallType} onChange={set('supportWallType')} />
            <Field label={A.otherType.label} value={o.otherType} onChange={set('otherType')} />
          </div>
        </div>
      </section>
    </div>
  )
}

export default function TabObject({ caseId, insuredName = '' }: { caseId: string; insuredName?: string }) {
  const [baris, setBaris] = useState<BarisLayar[]>([])
  const [terbuka, setTerbuka] = useState<Record<number, SubTab>>({})
  const [halaman, setHalaman] = useState(1)
  const [galat, setGalat] = useState<unknown>(null)
  const [menyimpan, setMenyimpan] = useState(false)
  const [tersimpan, setTersimpan] = useState(false)
  const [cobaSimpan, setCobaSimpan] = useState(false)
  const kunciBerikut = useRef(1)

  const bungkus = (data: ObjekFire[]): BarisLayar[] => data.map((d) => ({ kunci: kunciBerikut.current++, data: rapikanObjek(d) }))

  useEffect(() => {
    let batal = false
    ambilObjek(caseId).then(
      (h) => {
        if (!batal) setBaris(bungkus(h.baris))
      },
      (err: unknown) => {
        if (!batal) setGalat(err)
      },
    )
    return () => {
      batal = true
    }
  }, [caseId])

  function ubah(kunci: number, data: ObjekFire) {
    setBaris((b) => b.map((x) => (x.kunci === kunci ? { ...x, data } : x)))
    setTersimpan(false)
  }

  function tambah() {
    const kunci = kunciBerikut.current++
    setBaris((b) => [...b, { kunci, data: objekBaru(String(b.length + 1)) }])
    setTerbuka((t) => ({ ...t, [kunci]: SUBTAB_OBJEK[0] }))
    setHalaman(Math.ceil((baris.length + 1) / GRID_OBJEK.ukuran))
    setTersimpan(false)
  }

  function hapus(kunci: number) {
    setBaris((b) => b.filter((x) => x.kunci !== kunci))
    setTerbuka((t) => tanpa(t, kunci))
    setTersimpan(false)
  }

  const typeKosong = baris.some((x) => x.data.objectType === '')
  const adaLantaiMinus = baris.some(
    (x) =>
      lantaiMinus(x.data.numberOfFloor) ||
      adaJarakMinus(x.data.surroundingRisk) ||
      adaGalatItem(x.data.items) ||
      adaGalatFEA(x.data.fea) ||
      adaGalatKerugian(x.data.lossRecords),
  )

  async function simpan() {
    setCobaSimpan(true)
    if (typeKosong || adaLantaiMinus) {
      // Baris yang salah dibuka di sub-tab tempat pesannya (Object Address / Surrounding Risk / Object Item).
      const salah = baris.filter(
        (x) =>
          x.data.objectType === '' ||
          lantaiMinus(x.data.numberOfFloor) ||
          adaJarakMinus(x.data.surroundingRisk) ||
          adaGalatItem(x.data.items) ||
          adaGalatFEA(x.data.fea) ||
          adaGalatKerugian(x.data.lossRecords),
      )
      setTerbuka((t) => ({
        ...t,
        ...Object.fromEntries(
          salah.map((x) => [
            x.kunci,
            x.data.objectType === '' || lantaiMinus(x.data.numberOfFloor)
              ? SUBTAB_OBJEK[0]
              : adaJarakMinus(x.data.surroundingRisk)
                ? SUBTAB_OBJEK[1]
                : adaGalatItem(x.data.items)
                  ? SUBTAB_OBJEK[2]
                  : adaGalatFEA(x.data.fea)
                    ? SUBTAB_OBJEK[4]
                    : SUBTAB_OBJEK[5],
          ]),
        ),
      }))
      return
    }
    setMenyimpan(true)
    setGalat(null)
    setTersimpan(false)
    try {
      const h = await simpanObjek(
        caseId,
        baris.map((x) => x.data),
      )
      setBaris(bungkus(h.baris))
      setTerbuka({})
      setCobaSimpan(false)
      setTersimpan(true)
    } catch (err) {
      setGalat(err)
    } finally {
      setMenyimpan(false)
    }
  }

  const jumlahHalaman = Math.max(1, Math.ceil(baris.length / GRID_OBJEK.ukuran))
  const hal = Math.min(halaman, jumlahHalaman)
  const tampil = baris.slice((hal - 1) * GRID_OBJEK.ukuran, hal * GRID_OBJEK.ukuran)

  return (
    <div className="nbf-objek nbf-obj-tab">
      <Gagal galat={galat} />
      {tersimpan && <div className="alert alert--ok">{TEKS_INWARD.tersimpan}</div>}
      <div className="nbf-cov-wrap">
        <table className="nbf-tabel">
          <thead>
            <tr>
              <th scope="col" />
              {GRID_OBJEK.kolom.map((k) => (
                <th key={k.sel} scope="col">
                  {k.label}
                </th>
              ))}
              <th scope="col" className="table__actions">
                <button type="button" className="btn btn--ghost btn--sm" onClick={tambah}>
                  {GRID_OBJEK.tambah}
                </button>
              </th>
            </tr>
          </thead>
          <tbody>
            {baris.length === 0 && (
              <tr>
                <td colSpan={6}>{TEKS_INWARD.kosong}</td>
              </tr>
            )}
            {tampil.map((x) => {
              const sub = terbuka[x.kunci]
              return [
                <tr key={x.kunci} className={sub !== undefined ? 'nbf-baris--terbuka' : undefined}>
                  <td>
                    <button
                      type="button"
                      className={'btn btn--ghost btn--sm nbf-buka' + (sub !== undefined ? ' nbf-buka--terbuka' : '')}
                      aria-label={TEKS_OBJEK.bukaBaris}
                      aria-expanded={sub !== undefined}
                      onClick={() =>
                        setTerbuka((t) => (t[x.kunci] === undefined ? { ...t, [x.kunci]: SUBTAB_OBJEK[0] } : tanpa(t, x.kunci)))
                      }
                    >
                      ▸
                    </button>
                  </td>
                  <td>
                    <input type="checkbox" checked={x.data.isTopRisk} disabled aria-label={GRID_OBJEK.kolom[0].label} />
                  </td>
                  <td>{x.data.objectNo}</td>
                  <td>{x.data.objectName}</td>
                  <td>{x.data.riskLocation}</td>
                  <td className="table__actions">
                    <button type="button" className="btn btn--ghost btn--sm" onClick={() => hapus(x.kunci)}>
                      {GRID_OBJEK.hapus}
                    </button>
                  </td>
                </tr>,
                sub !== undefined && (
                  <tr key={`${x.kunci}-detail`} className="nbf-objek__detail">
                    <td colSpan={6}>
                      <div className="nbf-lapis nbf-lapis--objek">
                      <StripTab tab={SUBTAB_OBJEK} aktif={sub} onPilih={(t) => setTerbuka((s) => ({ ...s, [x.kunci]: t }))} />
                      {sub === SUBTAB_OBJEK[0] ? (
                        <ObjectAddress o={x.data} ubah={(d) => ubah(x.kunci, d)} galatType={cobaSimpan && x.data.objectType === ''} />
                      ) : sub === SUBTAB_OBJEK[1] ? (
                        <SubTabSekitar o={x.data} ubah={(d) => ubah(x.kunci, d)} />
                      ) : sub === SUBTAB_OBJEK[2] ? (
                        <SubTabItem items={x.data.items} ubah={(items) => ubah(x.kunci, { ...x.data, items })} tandaiWajib={cobaSimpan} />
                      ) : sub === SUBTAB_OBJEK[5] ? (
                        <SubTabKerugian
                          rows={x.data.lossRecords}
                          lossRatio={x.data.lossRatio}
                          insuredName={insuredName}
                          ubah={(lossRecords) => ubah(x.kunci, { ...x.data, lossRecords })}
                          tandaiWajib={cobaSimpan}
                        />
                      ) : sub === SUBTAB_OBJEK[6] ? (
                        <SubTabKlaimInternal rows={x.data.internalLossRecords} />
                      ) : sub === SUBTAB_OBJEK[4] ? (
                        <SubTabFEA fea={x.data.fea} ubah={(fea) => ubah(x.kunci, { ...x.data, fea })} />
                      ) : sub === SUBTAB_OBJEK[3] ? (
                        <SubTabOkupasi
                          caseId={caseId}
                          occupations={x.data.occupations}
                          ubah={(occupations) => ubah(x.kunci, { ...x.data, occupations })}
                        />
                      ) : (
                        <BelumTersedia apa={`${TEKS_INWARD.isiTab} ${sub}`} />
                      )}
                      </div>
                    </td>
                  </tr>
                ),
              ]
            })}
          </tbody>
        </table>
      </div>
      {baris.length > GRID_OBJEK.ukuran && (
        <Halaman halaman={hal} ukuran={GRID_OBJEK.ukuran} total={baris.length} onPindah={setHalaman} />
      )}
      <div className="nbf-objek__kaki">
        <button type="button" className="btn btn--primary" onClick={() => void simpan()} disabled={menyimpan}>
          {menyimpan ? TEKS_FORM_OPPORTUNITY.menyimpan : SIMPAN_OBJEK.label}
        </button>
      </div>
    </div>
  )
}
