// Popup Add alamat risiko - dari tombol Add popup Choose Risk Location (tiket 37).
//
// Port section `NB FacIn\Section\InputRiskAddress.xml` (harness `ChooseRiskLocation`). Tampilan = tangkapan layar
// work owner 03-10-2026: awalnya Country + Zip Code, Title (DESA), Address, Save / Close; setelah Zip Code dipilih,
// Province, City, District, Territory terisi dari tabel RW dan tampil.
//
// - Syarat tampil berantai (sel 5/7/9/11): Province bila Country terisi, City bila Province, District bila City,
//   Territory bila District.
// - Zip Code (sel 13, RD `BrowseRW_RD`, STS_AKTIF "1"): memilih saran mengisi Zip Code, Country, Province, City,
//   District, Territory (`.ZipCode`, `.NATIONNAME`, `.PROVINCENAME`, `.CITYNAME`, `.DistrictName`, `.Note`).
// - Save (`SaveRiskAddress_Act`): backend menulis RISKADDRESS dari Go tanpa memanggil prosedur (ADR-0043, butir 81),
//   lalu alamat langsung mengisi objek (langkah 9 activity) beserta ID barunya (keputusan work owner).
//   `SaveAccumulationByRiskAddress_Act` DITUNDA (keputusan work owner).
//
// Keputusan agent (tiket 37): I-1 Country / Province / City / District / Territory = isian teks bebas (Pega
// mengizinkan isian bebas); saran berantai `BrowseNation_RD` … `BrowseTeritory_RD` = tahap berikut. I-2 saran Zip
// Code muncul setelah 3 karakter, jeda ketik 400 ms; hanya disaring Zip Code + STS_AKTIF. I-3 setelah Save kedua
// popup ditutup dan objek terisi (Pega: form tetap terbuka karena pengecekan "Sukses" tidak pernah cocok).
// I-4 Zip Code dan Address wajib (Pega tanpa validasi) supaya tidak ada baris master kosong.

import { useEffect, useRef, useState } from 'react'

import { Area, Field, Gagal, Modal, Pilih } from '../../../../inti/frontend/components/ui/dasar'
import { cariZipCode, simpanAlamatBaru, type AlamatBaru, type BarisRW } from '../api'
import { OPSI_TITLE_RISK, POPUP_TAMBAH_RISK as T, TEKS_TAMBAH_RISK } from '../labels'

const OPSI_TITLE = OPSI_TITLE_RISK.map((v) => ({ value: v, label: v }))
const MIN_ZIP = 3
const JEDA_ZIP_MS = 400

/** Alamat kosong - Title = pilihan pertama (tanpa pilihan kosong di Pega). */
export function alamatKosong(): AlamatBaru {
  return {
    nationName: '', provinceName: '', districtName: '', cityName: '', territoryName: '', title: OPSI_TITLE_RISK[0],
    address: '', postalCode: '',
  }
}

/** Memilih saran Zip Code: enam medan dari baris RW. */
export function terapkanRW(a: AlamatBaru, r: BarisRW): AlamatBaru {
  return {
    ...a, postalCode: r.zipCode, nationName: r.nationName, provinceName: r.provinceName, cityName: r.cityName,
    districtName: r.districtName, territoryName: r.territoryName,
  }
}

/** I-4: Zip Code dan Address wajib. */
export function alamatLengkap(a: AlamatBaru): boolean {
  return a.postalCode.trim() !== '' && a.address.trim() !== ''
}

export default function PopupTambahRisk({
  onTutup,
  onTersimpan,
}: {
  onTutup: () => void
  onTersimpan: (a: AlamatBaru, id: string) => void
}) {
  const [a, setA] = useState<AlamatBaru>(alamatKosong)
  const [saran, setSaran] = useState<BarisRW[]>([])
  const [ketikZip, setKetikZip] = useState(false)
  const [cobaSimpan, setCobaSimpan] = useState(false)
  const [menyimpan, setMenyimpan] = useState(false)
  const [galat, setGalat] = useState<unknown>(null)
  const nomorPermintaan = useRef(0)

  // Saran Zip Code hanya saat pengguna mengetik di kotaknya (bukan sesudah memilih saran).
  useEffect(() => {
    if (!ketikZip || a.postalCode.trim().length < MIN_ZIP) {
      setSaran([])
      return
    }
    const nomor = ++nomorPermintaan.current
    const jadwal = window.setTimeout(() => {
      cariZipCode(a.postalCode).then(
        (h) => {
          if (nomor === nomorPermintaan.current) setSaran(h.baris)
        },
        (err: unknown) => {
          if (nomor === nomorPermintaan.current) setGalat(err)
        },
      )
    }, JEDA_ZIP_MS)
    return () => window.clearTimeout(jadwal)
  }, [a.postalCode, ketikZip])

  const set = (k: keyof AlamatBaru) => (v: string) => setA((x) => ({ ...x, [k]: v }))

  async function simpan() {
    setCobaSimpan(true)
    if (!alamatLengkap(a)) return
    setMenyimpan(true)
    setGalat(null)
    try {
      const h = await simpanAlamatBaru(a)
      onTersimpan(a, h.id)
    } catch (err) {
      setGalat(err)
    } finally {
      setMenyimpan(false)
    }
  }

  return (
    <Modal
      judul={T.judul}
      onTutup={onTutup}
      onKirim={() => void simpan()}
      lebar
      labelBatal={T.tutup.label}
      aksi={
        <button type="submit" className="btn btn--primary" disabled={menyimpan}>
          {menyimpan ? TEKS_TAMBAH_RISK.menyimpan : T.simpan.label}
        </button>
      }
    >
      <div className="nbf-tambah-risk">
        <Field label={T.country.label} value={a.nationName} onChange={set('nationName')} />
        {a.nationName !== '' && <Field label={T.province.label} value={a.provinceName} onChange={set('provinceName')} />}
        {a.provinceName !== '' && <Field label={T.city.label} value={a.cityName} onChange={set('cityName')} />}
        {a.cityName !== '' && <Field label={T.district.label} value={a.districtName} onChange={set('districtName')} />}
        {a.districtName !== '' && <Field label={T.territory.label} value={a.territoryName} onChange={set('territoryName')} />}
        <div className="nbf-tambah-risk__zip">
          <Field
            label={T.zipCode.label}
            value={a.postalCode}
            onChange={(v) => {
              setKetikZip(true)
              set('postalCode')(v)
            }}
            error={cobaSimpan && a.postalCode.trim() === '' ? TEKS_TAMBAH_RISK.wajib : undefined}
          />
          {saran.length > 0 && (
            <ul className="nbf-tambah-risk__saran" role="listbox" aria-label={T.zipCode.label}>
              {saran.map((r, i) => (
                <li key={`${i}-${r.zipCode}-${r.territoryName}`}>
                  <button
                    type="button"
                    className="btn btn--ghost btn--sm"
                    onClick={() => {
                      setKetikZip(false)
                      setA((x) => terapkanRW(x, r))
                      setSaran([])
                    }}
                  >
                    {[r.zipCode, r.territoryName, r.districtName, r.cityName, r.provinceName, r.nationName].join(' · ')}
                  </button>
                </li>
              ))}
            </ul>
          )}
        </div>
      </div>
      <div className="nbf-tambah-risk__bawah">
        <Pilih label={T.title.label} value={a.title} onChange={set('title')} opsi={OPSI_TITLE} />
        <Area
          label={T.address.label}
          value={a.address}
          onChange={set('address')}
          error={cobaSimpan && a.address.trim() === '' ? TEKS_TAMBAH_RISK.wajib : undefined}
        />
      </div>
      <Gagal galat={galat} />
    </Modal>
  )
}
