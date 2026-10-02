// Blok pertama section `InputCoverageCargo_FacIn` (MARINE CARGO, satu coverage) - tiket 21.
//
// Urutan dan kelompok = layout Pega: layout 1 (sel 3-5), 6/8 (10-12), 13 (15-18),
// 19 (21-24). Di Pega seluruh MEDAN blok ini read-only; PENYIMPANGAN SADAR butir 61
// (keputusan work owner 02-10-2026): Name, Rate (%), TSI dapat diisi untuk memanggil
// `POST /api/nbfacin/premi` - alat periksa rumus, dinyatakan di layar. Name dropdown
// (`pxDropdown`) ditampilkan sebagai kotak teks: sumber opsinya tidak ada di blok ini.
//
// `[terverifikasi]` `D:\migrasi\RNM\NB FacIn\Section\InputCoverageCargo_FacIn.xml`
// (ASM-FW-GISFW-DATA-COVERAGE!INPUTCOVERAGECARGO_FACIN) `pyUserData/pyCondition`: sel 12
// tampil bila `PolicyMasterNumber` TIDAK sama dengan satu nomor polis literal, sel 16
// bila sama (nilai tidak disalin). Hanya cabang umum (sel 12) yang diport.
// Medan tanpa sumber data (tiket 17) tampil kosong - tidak diisi nilai rekaan. Premi
// adalah teks dari backend: frontend tidak menghitung uang.

import { useRef, useState } from 'react'

import { Field, Gagal, Panel } from '../../../../inti/frontend/components/ui/dasar'
import { hitungPremiCargo, type HasilPremi } from '../api'
import { MEDAN_COVERAGE_CARGO as MEDAN, TEKS_COVERAGE_CARGO as TEKS, TOMBOL_COVERAGE_CARGO as TOMBOL } from '../labels'

const tanpaUbah = () => {}

/** Tombol Pega yang belum diport: nonaktif, keterangannya TERLIHAT (bukan hanya tooltip). */
function TombolBelumDiport({ label }: { label: string }) {
  return (
    <span>
      <button type="button" className="btn btn--ghost btn--sm" disabled aria-describedby="nbfacin-belum-diport">
        {label}
      </button>{' '}
      <small>({TEKS.belumDiport})</small>
    </span>
  )
}

export default function CoverageCargo() {
  const [mataUang, setMataUang] = useState('')
  const [rate, setRate] = useState('')
  const [tsi, setTsi] = useState('')
  const [hasil, setHasil] = useState<HasilPremi | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(false)
  // Nomor permintaan terakhir: jawaban permintaan lama (isian sudah berubah) dibuang.
  const nomorPermintaan = useRef(0)

  async function hitung() {
    const nomor = ++nomorPermintaan.current
    setSibuk(true)
    setGalat(null)
    try {
      const h = await hitungPremiCargo({ mataUang, rate, tsi })
      if (nomor === nomorPermintaan.current) setHasil(h)
    } catch (err) {
      if (nomor === nomorPermintaan.current) {
        setHasil(null)
        setGalat(err)
      }
    } finally {
      if (nomor === nomorPermintaan.current) setSibuk(false)
    }
  }

  // Isian berubah = premi lama tidak berlaku lagi, dan jawaban yang masih di jalan pun.
  const ubah = (set: (v: string) => void) => (v: string) => {
    set(v)
    nomorPermintaan.current++
    setHasil(null)
    setSibuk(false)
  }

  return (
    <Panel judul={TEKS.judul} catatan={TEKS.alatPeriksa}>
      <span id="nbfacin-belum-diport" hidden>
        {TEKS.belumDiport}
      </span>
      <div className="form-grid">
        {/* TextArea read-only (sel 3). `Area` inti tidak punya mode read-only, dan
            mengubah komponen bersama di luar jatah modul ini - markahnya ditiru di sini. */}
        <div className="field field--lebar">
          <label className="field__label" htmlFor="nbfacin-keterangan-jaminan">
            {MEDAN.keteranganJaminan.label}
          </label>
          <textarea id="nbfacin-keterangan-jaminan" className="field__input field__input--readonly" rows={4} value="" readOnly />
        </div>
        <Field label={MEDAN.coverageInitial.label} value="" onChange={tanpaUbah} readOnly />
        <TombolBelumDiport label={TOMBOL.pilihCoverage.label} />
      </div>
      <div className="form-grid">
        <Field label={MEDAN.mataUang.label} value={mataUang} onChange={ubah(setMataUang)} />
        <Field label={MEDAN.rate.label} value={rate} onChange={ubah(setRate)} />
        <Field label={MEDAN.limitOfLiability.label} value="" onChange={tanpaUbah} readOnly />
      </div>
      <div className="form-grid">
        <Field label={MEDAN.currencyMaster.label} value="" onChange={tanpaUbah} readOnly />
        <TombolBelumDiport label={TOMBOL.ambilCurrencyMaster.label} />
        <Field label={MEDAN.diskonPersen.label} value="" onChange={tanpaUbah} readOnly />
      </div>
      <div className="form-grid">
        <Field label={MEDAN.tsi.label} value={tsi} onChange={ubah(setTsi)} required />
        <Field label={MEDAN.premi.label} value={hasil?.premi ?? ''} onChange={tanpaUbah} readOnly />
        <Field label={MEDAN.minPremi.label} value="" onChange={tanpaUbah} readOnly />
        <Field label={MEDAN.diskon.label} value="" onChange={tanpaUbah} readOnly />
      </div>
      <p>
        <button type="button" className="btn" onClick={() => void hitung()} disabled={sibuk}>
          {sibuk ? TEKS.menghitung : TEKS.hitung}
        </button>
      </p>
      {hasil && (
        <p className="panel__note">
          {TEKS.asalRumus}: {hasil.asalRumus}
        </p>
      )}
      <Gagal galat={galat} />
    </Panel>
  )
}
