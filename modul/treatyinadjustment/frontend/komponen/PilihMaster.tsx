// Modal picker tombol `Add Revision` / `Add Adjustment Premium` —
// `Section/PickerTreatyInMasterRevisi.xml` / `Section/PickerTreatyInMaster.xml`.
//
// ⭐ Rantainya dari ekspor (`Section/InputTreatyInAdjustment.xml` @500554 /
// @515429): DataTransform `TreatyCreateEDM` (EDMState "1" untuk revisi, "3"
// untuk premi), lalu modal ini (`pyModalFullScreen = true`). `Choose`
// menjalankan `TreatyInEDMSetValue` lalu `closeContainer`.
//
//   Revisi  radio `Internal / External` + `Material Type`, grid TREATY_IN ∪
//           TREATY_IN_EDM; Choose → InternalType = EDMState, MaterialType
//           = EDMMaterialType
//   Premi   tanpa radio, grid NonProportional saja; Choose → InternalType
//           '3', MaterialType '1'
//
// ⛔ `Choose` di Pega juga MENYIMPAN. Di sini ia mengembalikan DRAF — lihat
// `buatDrafPenyesuaian`.

import { useEffect, useMemo, useState } from 'react'

import { Gagal, Halaman, Kosong, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { ambilDaftarMaster, buatDrafPenyesuaian, type BarisMasterPilihan, type JenisPicker, type Penyesuaian } from '../api'
import {
  CARI_MASTER,
  JENIS_PICKER,
  KOLOM_PICKER,
  LEBAR_PICKER_PREMI,
  LEBAR_PICKER_REVISI,
  OPSI_JENIS_REVISI,
  OPSI_MATERIAL,
  PENYESUAIAN,
} from '../labelsPenyesuaian'
import { teksPromptEDM } from '../labelsPromptEDM'
import { persenLebar } from './lebar'
import { selNilai } from './medan'

/** Baris per halaman — `pyGridPaginator`, ukuran bawaan Pega 10 (sama dengan grid daftar). */
const UKURAN_HALAMAN = 10

/** Nilai grid picker, urut sesuai `KOLOM_PICKER`. */
export function selMaster(b: BarisMasterPilihan): string[] {
  return [b.id, b.namaKontrak, b.sifatProporsi, b.asalBisnis, b.cedant, b.tanggalMulai, b.tanggalBerakhir].map((v, i) =>
    selNilai(JENIS_PICKER[i] ?? 'teks', v),
  )
}

/**
 * Pencarian per kolom (urut `KOLOM_PICKER`): baris lolos bila SETIAP kotak
 * yang terisi termuat di nilai TAMPIL kolomnya, tanpa membedakan huruf besar.
 * Nilai tampil, bukan nilai mentah — tanggal dicari seperti terlihat
 * (`01-01-2019`), bukan `20190101`.
 */
export function saringMaster(baris: readonly BarisMasterPilihan[], cari: readonly string[]): BarisMasterPilihan[] {
  const kunci = cari.map((c) => c.trim().toLowerCase())
  if (kunci.every((c) => c === '')) return [...baris]
  return baris.filter((b) => {
    const sel = selMaster(b)
    return kunci.every((c, i) => c === '' || (sel[i] ?? '').toLowerCase().includes(c))
  })
}

/** Judul picker — LABEL Section yang tampil untuk jenis dan `EDMState` ini. */
export function judulPicker(jenis: JenisPicker, edmState: string): string {
  if (jenis === 'premi') return PENYESUAIAN.pickerPremi
  return edmState === '2' ? PENYESUAIAN.pickerRevisiEksternal : PENYESUAIAN.pickerRevisiInternal
}

export default function PilihMaster({
  jenis,
  onTutup,
  onDraf,
}: {
  jenis: JenisPicker
  onTutup: () => void
  /** Draf dari `Choose` — `closeContainer` lalu layar detail. */
  onDraf: (p: Penyesuaian) => void
}) {
  // `TreatyCreateEDM`: `TreatyIn = ""`, lalu EDMState "1" (revision) / "3"
  // (adjustment). Material Type tidak disetel — kosong sampai dipilih.
  const [edmState, setEdmState] = useState(jenis === 'premi' ? '3' : '1')
  const [material, setMaterial] = useState('')
  const [baris, setBaris] = useState<BarisMasterPilihan[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [halaman, setHalaman] = useState(1)
  const [memilih, setMemilih] = useState<string | null>(null)
  const [galatPilih, setGalatPilih] = useState<unknown>(null)
  // ⭐ Pencarian per kolom — fitur baru atas permintaan pemakai (lihat
  // `CARI_MASTER`). Kosong = semua baris.
  const [cari, setCari] = useState<string[]>(() => KOLOM_PICKER.map(() => ''))

  useEffect(() => {
    let dibuang = false
    ambilDaftarMaster(jenis)
      .then((d) => {
        if (!dibuang) setBaris(d)
      })
      .catch((e: unknown) => {
        if (!dibuang) setGalat(e)
      })
    return () => {
      dibuang = true
    }
  }, [jenis])

  const pilih = (id: string) => {
    setMemilih(id)
    setGalatPilih(null)
    buatDrafPenyesuaian({
      id,
      internalType: jenis === 'premi' ? '3' : edmState,
      materialType: jenis === 'premi' ? '1' : material,
    })
      .then(onDraf)
      .catch((e: unknown) => {
        setGalatPilih(e)
      })
      .finally(() => {
        setMemilih(null)
      })
  }

  const lebar = jenis === 'premi' ? LEBAR_PICKER_PREMI : LEBAR_PICKER_REVISI
  const semua = baris ?? []
  const lolos = useMemo(() => saringMaster(baris ?? [], cari), [baris, cari])
  const tampil = lolos.slice((halaman - 1) * UKURAN_HALAMAN, halaman * UKURAN_HALAMAN)
  const adaCari = cari.some((c) => c.trim() !== '')
  const ubahCari = (i: number, v: string) => {
    setCari((x) => x.map((c, j) => (j === i ? v : c)))
    // Halaman lama bisa melampaui hasil saringan.
    setHalaman(1)
  }

  return (
    <Modal judul={judulPicker(jenis, edmState)} onTutup={onTutup} penuh>
      {jenis === 'revisi' && (
        <div className="tria__picker-radio">
          <fieldset className="tria__radio">
            <legend>{PENYESUAIAN.internalEksternal}</legend>
            {OPSI_JENIS_REVISI.map((o) => (
              <label key={o.nilai}>
                <input
                  type="radio"
                  name="tria-picker-edmstate"
                  value={o.nilai}
                  checked={edmState === o.nilai}
                  onChange={() => {
                    setEdmState(o.nilai)
                  }}
                />
                {o.label}
              </label>
            ))}
          </fieldset>
          <fieldset className="tria__radio">
            <legend>{PENYESUAIAN.materialRadio}</legend>
            {OPSI_MATERIAL.map((v) => (
              <label key={v}>
                <input
                  type="radio"
                  name="tria-picker-material"
                  value={v}
                  checked={material === v}
                  onChange={() => {
                    setMaterial(v)
                  }}
                />
                {/* *(8 Okt, E)* `TreatyIn.EDMMaterialType` @110848 — prompt value
                    `ekspor-tambahan/EDMMaterialType.xml` (1 Material, 2 Non Material). */}
                {teksPromptEDM('EDMMaterialType', v)}
              </label>
            ))}
            <span className="tria__redup">{PENYESUAIAN.kodeBelumBerteks}</span>
          </fieldset>
        </div>
      )}

      {galatPilih !== null && <Gagal galat={galatPilih} />}
      {galat !== null && <Gagal galat={galat} />}
      {baris === null && galat === null && <Memuat />}
      {baris !== null && (
        <>
          <div className="table-wrap">
            <table className="tria__tabel">
              <colgroup>
                {lebar.map((_, i) => (
                  <col key={i} style={{ width: persenLebar(lebar, i) }} />
                ))}
              </colgroup>
              <thead>
                <tr>
                  <th scope="col" aria-label={PENYESUAIAN.pilihMaster} />
                  {KOLOM_PICKER.map((k) => (
                    <th key={k} scope="col">
                      {k}
                    </th>
                  ))}
                </tr>
                <tr className="tria__baris-cari">
                  <th scope="col">
                    <button
                      type="button"
                      className="btn btn--ghost btn--sm"
                      disabled={!adaCari}
                      onClick={() => {
                        setCari(KOLOM_PICKER.map(() => ''))
                        setHalaman(1)
                      }}
                    >
                      {CARI_MASTER.reset}
                    </button>
                  </th>
                  {KOLOM_PICKER.map((k, i) => (
                    <th key={k} scope="col">
                      <input
                        className="field__input"
                        type="search"
                        aria-label={k}
                        placeholder={CARI_MASTER.petunjuk}
                        value={cari[i] ?? ''}
                        onChange={(e) => {
                          ubahCari(i, e.target.value)
                        }}
                      />
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {semua.length === 0 && (
                  <tr>
                    <td colSpan={KOLOM_PICKER.length + 1}>
                      <Kosong pesan={PENYESUAIAN.tanpaBaris} petunjuk={PENYESUAIAN.petunjukMasterKosong} />
                    </td>
                  </tr>
                )}
                {semua.length > 0 && lolos.length === 0 && (
                  <tr>
                    <td colSpan={KOLOM_PICKER.length + 1}>
                      <Kosong pesan={CARI_MASTER.tanpaHasil} petunjuk="" />
                    </td>
                  </tr>
                )}
                {tampil.map((b) => (
                  <tr key={b.id}>
                    <td>
                      <button
                        type="button"
                        className="btn btn--ghost btn--sm"
                        disabled={memilih !== null}
                        onClick={() => {
                          pilih(b.id)
                        }}
                      >
                        {PENYESUAIAN.pilihMaster}
                      </button>
                    </td>
                    {selMaster(b).map((v, i) => (
                      <td key={i}>{v}</td>
                    ))}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          {adaCari && <p className="tria__redup">{CARI_MASTER.hasil(lolos.length, semua.length)}</p>}
          <Halaman halaman={halaman} ukuran={UKURAN_HALAMAN} total={lolos.length} onPindah={setHalaman} />
        </>
      )}
    </Modal>
  )
}
