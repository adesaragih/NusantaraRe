// Popup Choose Risk Address - sub-tab Object Address, tab Object FIRE (tiket 36).
//
// Port harness `ChooseRiskAddress` (`NB FacIn\Section\ChooseRiskAddress.xml` + `ChooseRiskAddress_ResultList.xml`;
// dibuka ObjectDetails sel 16, WindowName "Choose Risk Location"): saringan kiri (Address, Zip Code, Country,
// Province, City, District, Territory), grid kanan atas RD `BrowseRisksAddress_RD` (tabel RISKADDRESS, 10 per
// halaman bernomor), tombol `Pilih` per baris, `Search` dan `Add` di kaki. Tampilan = tangkapan layar work owner
// 03-10-2026 (tidak disalin: memuat alamat).
//
// Pega: saringan terikat LANGSUNG ke medan objek (`.Property.RoadName` dst.); `SetDefaultSearchRiskLocation_Act`
// mengisi kunci cari palsu ('z' / '123456') sehingga grid kosong sampai Search; `SearchRiskAddressAct` hanya
// mencari bila minimal satu saringan terisi (huruf besar). Pilih → `SetRiskIdDT_FacIn` → `ObjSave` → muat ulang
// layar → tutup.
//
// Keputusan agent (tiket 36): H-1 saringan = salinan medan objek saat dibuka (diisi awal dari objek), bukan
// terikat ke objek - mengetik di saringan tidak mengubah objek sebelum Pilih. H-2 grid kosong sampai Search. H-3
// Pilih mengisi objek di layar; tersimpan lewat Save tab Object (sejalan E-1). H-4 `Add` (harness
// `ChooseRiskLocation`, alamat baru ke RISKADDRESS) = tahap berikut - nonaktif.

import { useRef, useState } from 'react'

import { Field, Gagal, Halaman, Kosong, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { cariRiskAddress, type BarisRisk, type HalamanRisk, type SaringRisk } from '../api'
import { POPUP_RISK as R, TEKS_RISK } from '../labels'

/** Ada minimal satu saringan terisi. */
export function adaSaring(s: SaringRisk): boolean {
  return Object.values(s).some((v) => v.trim() !== '')
}

export default function PopupRiskAddress({
  awal,
  onTutup,
  onPilih,
}: {
  awal: SaringRisk
  onTutup: () => void
  onPilih: (b: BarisRisk) => void
}) {
  const [saring, setSaring] = useState<SaringRisk>(awal)
  const [hasil, setHasil] = useState<HalamanRisk | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [memuat, setMemuat] = useState(false)
  const [kurang, setKurang] = useState(false)
  // Jawaban lama dibuang; paging memakai saringan Search terakhir.
  const nomorPermintaan = useRef(0)
  const saringCari = useRef<SaringRisk>(awal)

  async function muat(s: SaringRisk, halaman: number) {
    if (!adaSaring(s)) {
      setKurang(true)
      return
    }
    setKurang(false)
    const nomor = ++nomorPermintaan.current
    saringCari.current = s
    setMemuat(true)
    setGalat(null)
    try {
      const h = await cariRiskAddress(s, halaman)
      if (nomor === nomorPermintaan.current) setHasil(h)
    } catch (err) {
      if (nomor === nomorPermintaan.current) {
        setHasil(null)
        setGalat(err)
      }
    } finally {
      if (nomor === nomorPermintaan.current) setMemuat(false)
    }
  }

  const set = (k: keyof SaringRisk) => (v: string) => setSaring((s) => ({ ...s, [k]: v }))

  return (
    <Modal
      judul={R.judul}
      onTutup={onTutup}
      onKirim={() => void muat(saring, 1)}
      penuh
      aksi={
        <>
          <button type="submit" className="btn btn--primary" disabled={memuat}>
            {R.cari.label}
          </button>
          <button type="button" className="btn btn--ghost" disabled>
            {R.tambah.label}
          </button>
        </>
      }
    >
      <div className="nbf-risk">
        <div className="nbf-risk__saring">
          <Field label={R.saring.address.label} value={saring.address} onChange={set('address')} />
          <Field label={R.saring.zipCode.label} value={saring.zipCode} onChange={set('zipCode')} />
          <Field label={R.saring.country.label} value={saring.country} onChange={set('country')} />
          <Field label={R.saring.province.label} value={saring.province} onChange={set('province')} />
          <Field label={R.saring.city.label} value={saring.city} onChange={set('city')} />
          <Field label={R.saring.district.label} value={saring.district} onChange={set('district')} />
          <Field label={R.saring.territory.label} value={saring.territory} onChange={set('territory')} />
        </div>
        <div className="nbf-risk__hasil">
          <div className="table-wrap">
            <table className="nbf-tabel">
              <thead>
                <tr>
                  {R.kolom.map((k) => (
                    <th key={k.sel} scope="col">
                      {k.label}
                    </th>
                  ))}
                  <th scope="col" />
                </tr>
              </thead>
              {hasil && hasil.baris.length > 0 && (
                <tbody>
                  {hasil.baris.map((b) => (
                    <tr key={b.id}>
                      <td>{b.title}</td>
                      <td>{b.address}</td>
                      <td>{b.nationName}</td>
                      <td>{b.provinceName}</td>
                      <td>{b.cityName}</td>
                      <td>{b.districtName}</td>
                      <td>{b.territoryName}</td>
                      <td>{b.postalCode}</td>
                      <td className="table__actions">
                        <button type="button" className="btn btn--ghost btn--sm" onClick={() => onPilih(b)}>
                          {R.pilih}
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              )}
            </table>
          </div>
          {hasil && hasil.baris.length > 0 && (
            <Halaman halaman={hasil.halaman} ukuran={hasil.ukuran} total={hasil.total} onPindah={(h) => void muat(saringCari.current, h)} />
          )}
          {kurang && <div className="alert alert--warn">{TEKS_RISK.isiSaring}</div>}
          {memuat && !hasil && <Memuat />}
          {hasil && hasil.baris.length === 0 && <Kosong pesan={TEKS_RISK.tanpaHasil} />}
          <Gagal galat={galat} />
        </div>
      </div>
    </Modal>
  )
}
