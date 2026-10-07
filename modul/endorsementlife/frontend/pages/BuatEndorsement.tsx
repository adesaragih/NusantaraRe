// Buat endorsement - `Harness/EndorsmentLife_harnes.xml` → `Section/EndorsmentLife_Section.xml`
// (dibuka `Create Addendum` `InboxEndorsementLife.xml` b6620 → `showHarness` b6984).
//
// Di Pega `SetErrorBatalEndorsement_Act` berjalan saat `Policy No` (b1249), `EDM Type` (b1545), dan `EDM Date`
// (b3347) berubah. KINI (keputusan work owner 07-10-2026) gerbang hanya berjalan saat Policy No diisi/diubah dan
// saat Submit: `cekKelayakan(policyNo, '')` - EDM Type sengaja KOSONG sehingga yang dinilai hanya polisnya (tidak
// ditemukan, EDM belum selesai, sudah batal); Policy No kosong tidak memanggil server (tanpa pesan merah). Mengganti
// EDM Type / EDM Date / Description tidak memanggil server. Aturan "sudah dibayar + Batal" dinilai saat Submit:
// `BuatKasus` menjalankan ulang seluruh gerbang dengan EDM Type dan penolakannya (422, `pesanPenolakanEDM`)
// tampil sebagai daftar pesan gerbang yang sama; mengubah isian menghapusnya dan Submit aktif lagi.
// `Submit` b4226 menjalankan `MappingEDMLife` lalu `openAssignment` (buka layar kasus).
//
// Tampilan (keputusan work owner 07-10-2026, logika tidak berubah): panel selebar area konten; baris 1 Policy No |
// EDM Type | EDM Date, baris 2 Description; footer rata kanan Cancel + Submit. ⚠️ Dulu Submit DISEMBUNYIKAN
// selama gerbang belum lolos (Pega b4226 VIS) - kini SELALU tampil tetapi nonaktif, dengan teks bantuan yang
// menyebut apa yang masih kurang, demi kejelasan langkah berikutnya. Penanda wajib (*) hanya Policy No dan
// EDM Type: Policy No masukan gerbang, EDM Type 1/3 syarat Submit (`bolehSubmitEDM`, dan `BuatKasus`). EDM Date
// dan Description boleh kosong (`BuatKasus` menerima tanggal kosong).
//
// ⏸️ `EDM Type Perubahan Data` (b1645) dan `EDM Type Batal` (b1937) tidak dibangun: opsinya tidak ada di korpus
// dan tidak punya kolom (OQ-EDM-005). ➖ `Process Policy No` b2414 `VIS=never`.

import { useEffect, useState } from 'react'

import { Area, Field, FieldTanggal, Gagal, IkonTutup, Pilih } from '../../../../inti/frontend/components/ui/dasar'
import { keInputTanggal } from '../../../../inti/frontend/lib/tanggalInput'
import { buatKasus, cekKelayakan, pesanPenolakanEDM, type KelayakanEDM } from '../api'
import { BUAT_EDM, INBOX_EDM, OPSI_EDM_TYPE, OPSI_EDM_TYPE_BATAL, OPSI_EDM_TYPE_PERUBAHAN_DATA, UMUM_EDM } from '../labels'
import '../endorsementlife.css'

/** Opsi `EDM Type` - `1` Perubahan Data, `3` Batal (spec §5). */
const OPSI = Object.entries(OPSI_EDM_TYPE).map(([value, label]) => ({ value, label }))
/** Opsi sub-jenis - nilai = teks label (OQ-EDM-005, 07-10-2026). */
const OPSI_PERUBAHAN = OPSI_EDM_TYPE_PERUBAHAN_DATA.map((v) => ({ value: v, label: v }))

/** Bawaan `EDM Type Batal` - terpilih tanpa Choose. */
export const EDM_TYPE_BATAL_BAWAAN = OPSI_EDM_TYPE_BATAL[0] ?? ''

/**
 * Sub-jenis sesudah EDM Type diganti: yang tidak relevan dikosongkan; Batal langsung ke bawaannya,
 * Perubahan Data kembali ke Choose (kosong).
 */
export function subJenisSesudahGantiEdmType(edmType: string): { perubahanData: string; batal: string } {
  return { perubahanData: '', batal: edmType === '3' ? EDM_TYPE_BATAL_BAWAAN : '' }
}

/** Sub-jenis terisi dengan salah satu opsi yang sesuai EDM Type (Submit mensyaratkannya). */
export function subJenisSah(edmType: string, perubahanData: string, batal: string): boolean {
  if (edmType === '1') return OPSI_EDM_TYPE_PERUBAHAN_DATA.includes(perubahanData) && batal === ''
  if (edmType === '3') return OPSI_EDM_TYPE_BATAL.includes(batal) && perubahanData === ''
  return false
}

/** Keadaan yang menentukan Submit dan teks bantuannya. */
export interface KeadaanSubmitEDM {
  policyNo: string
  edmType: string
  /** Hasil gerbang Policy No (dipanggil dengan EDM Type kosong); null = belum ada. */
  kelayakan: KelayakanEDM | null
  /** Sub-jenis EDM Type (07-10-2026). */
  perubahanData: string
  batal: string
  /** Pesan penolakan Submit terakhir (`pesanPenolakanEDM`); dihapus saat isian berubah. */
  penolakan: readonly string[]
  memeriksa: boolean
  mengirim: boolean
}

/** EDM Type yang dapat dibuat - `1` Perubahan Data, `3` Batal (diperiksa di klien). */
export function edmTypeSah(edmType: string): boolean {
  return Object.prototype.hasOwnProperty.call(OPSI_EDM_TYPE, edmType)
}

/**
 * Submit boleh ditekan (keputusan work owner 07-10-2026): Policy No terisi, hasil gerbang ada dan tanpa pesan,
 * EDM Type 1/3, tidak sedang memeriksa / mengirim, dan tidak ada penolakan Submit yang belum ditanggapi.
 * ⛔ Tidak memakai `kelayakan.boleh`: gerbang dipanggil tanpa EDM Type, jadi `boleh` server selalu false.
 */
export function bolehSubmitEDM(k: KeadaanSubmitEDM): boolean {
  return (
    k.policyNo.trim() !== '' &&
    k.kelayakan !== null &&
    k.kelayakan.pesan.length === 0 &&
    edmTypeSah(k.edmType) &&
    subJenisSah(k.edmType, k.perubahanData, k.batal) &&
    !k.memeriksa &&
    !k.mengirim &&
    k.penolakan.length === 0
  )
}

/**
 * Teks bantuan di samping Submit nonaktif: apa yang masih kurang. Kosong = Submit boleh ditekan.
 * Tidak menambah aturan wajib baru.
 */
export function alasanBelumSubmit(k: KeadaanSubmitEDM): string {
  if (k.mengirim) return BUAT_EDM.bantuan.mengirim
  if (k.policyNo.trim() === '') return BUAT_EDM.bantuan.policyNo
  if (k.memeriksa) return BUAT_EDM.bantuan.memeriksa
  if (k.penolakan.length > 0) return BUAT_EDM.bantuan.penolakan
  if (k.kelayakan === null) return BUAT_EDM.bantuan.menunggu
  if (k.kelayakan.pesan.length > 0) return BUAT_EDM.bantuan.pesan
  if (!edmTypeSah(k.edmType)) return BUAT_EDM.bantuan.edmType
  if (!subJenisSah(k.edmType, k.perubahanData, k.batal)) {
    return k.edmType === '3' ? BUAT_EDM.bantuan.batal : BUAT_EDM.bantuan.perubahanData
  }
  return ''
}

export default function BuatEndorsement({ onTutup, onDibuat }: { onTutup: () => void; onDibuat: (id: string) => void }) {
  const [policyNo, setPolicyNo] = useState('')
  const [edmType, setEdmType] = useState('')
  const [edmDate, setEdmDate] = useState('')
  const [deskripsi, setDeskripsi] = useState('')
  // Sub-jenis EDM Type (OQ-EDM-005, 07-10-2026) - mengubahnya TIDAK memanggil gerbang.
  const [perubahanData, setPerubahanData] = useState('')
  const [batal, setBatal] = useState('')
  const [kelayakan, setKelayakan] = useState<KelayakanEDM | null>(null)
  const [memeriksa, setMemeriksa] = useState(false)
  const [mengirim, setMengirim] = useState(false)
  const [galat, setGalat] = useState<unknown>(null)
  // Penolakan Submit berpesan (422) - dihapus begitu isian apa pun berubah (07-10-2026).
  const [penolakan, setPenolakan] = useState<string[]>([])
  const ubah = (set: (v: string) => void) => (v: string) => {
    setPenolakan([])
    set(v)
  }
  const gantiEdmType = (v: string) => {
    setPenolakan([])
    setEdmType(v)
    const sub = subJenisSesudahGantiEdmType(v)
    setPerubahanData(sub.perubahanData)
    setBatal(sub.batal)
  }

  // Gerbang HANYA saat Policy No berubah (keputusan work owner 07-10-2026), debounce 300ms; EDM Type dikirim
  // kosong (pemeriksaan polis saja). Policy No kosong → tanpa panggilan server, tanpa pesan.
  useEffect(() => {
    const np = policyNo.trim()
    if (np === '') {
      setKelayakan(null)
      setMemeriksa(false)
      return
    }
    let batal = false
    setMemeriksa(true)
    const t = setTimeout(() => {
      cekKelayakan(np, '')
        .then((k) => {
          if (!batal) setKelayakan(k)
        })
        .catch((e: unknown) => {
          if (!batal) setGalat(e)
        })
        .finally(() => {
          if (!batal) setMemeriksa(false)
        })
    }, 300)
    return () => {
      batal = true
      clearTimeout(t)
    }
  }, [policyNo])

  const keadaan: KeadaanSubmitEDM = { policyNo, edmType, perubahanData, batal, kelayakan, penolakan, memeriksa, mengirim }
  const bolehSubmit = bolehSubmitEDM(keadaan)

  async function kirim() {
    setGalat(null)
    setMengirim(true)
    try {
      const h = await buatKasus({
        policyNo: policyNo.trim(),
        edmType,
        edmDate: keInputTanggal(edmDate),
        description: deskripsi,
        edmTypePerubahanData: edmType === '1' ? perubahanData : '',
        edmTypeBatal: edmType === '3' ? batal : '',
      })
      onDibuat(h.id)
    } catch (e) {
      // Penolakan gerbang saat Submit (mis. sudah dibayar + Batal) → daftar pesan gerbang, bukan galat umum.
      const pesan = pesanPenolakanEDM(e)
      if (pesan !== null) setPenolakan(pesan)
      else setGalat(e)
    } finally {
      setMengirim(false)
    }
  }

  const alasan = bolehSubmit ? '' : alasanBelumSubmit(keadaan)
  const pesanGerbang = [...(kelayakan?.pesan ?? []), ...penolakan]

  return (
    <section className="panel edm-buat">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{INBOX_EDM.createAddendum}</h2>
        <button type="button" className="btn btn--ghost btn--sm" aria-label={UMUM_EDM.tutup} title={UMUM_EDM.tutup} onClick={onTutup}>
          <IkonTutup />
        </button>
      </header>
      {/* Baris 1: tiga isian sama lebar; baris 2: Description selebar penuh (07-10-2026). */}
      <div className="edm-buat__baris">
        <Field label={BUAT_EDM.policyNo} value={policyNo} onChange={ubah(setPolicyNo)} autoFocus required />
        <Pilih label={BUAT_EDM.edmType} value={edmType} onChange={gantiEdmType} opsi={OPSI} required />
        <FieldTanggal label={BUAT_EDM.edmDate} value={edmDate} onChange={ubah(setEdmDate)} />
      </div>
      {/* Sub-jenis EDM Type tepat di bawah EDM Type (kolom 2 grid; satu kolom di layar sempit) - 07-10-2026. */}
      {edmType === '1' && (
        <div className="edm-buat__baris">
          <div className="edm-buat__sub">
            <Pilih
              label={BUAT_EDM.edmTypePerubahanData}
              value={perubahanData}
              onChange={ubah(setPerubahanData)}
              opsi={OPSI_PERUBAHAN}
              kosong={BUAT_EDM.pilih}
              required
            />
          </div>
        </div>
      )}
      {edmType === '3' && (
        <div className="edm-buat__baris">
          {/* Tanpa pilihan Choose (bawaan Batal Internal): `Pilih` inti selalu menambah opsi kosong. */}
          <div className="edm-buat__sub field">
            <label className="field__label" htmlFor="edm-buat-batal">
              {BUAT_EDM.edmTypeBatal}
              <span className="field__req">*</span>
            </label>
            <select
              id="edm-buat-batal"
              className="field__input"
              value={batal}
              onChange={(e) => {
                ubah(setBatal)(e.target.value)
              }}
            >
              {OPSI_EDM_TYPE_BATAL.map((o) => (
                <option key={o} value={o}>
                  {o}
                </option>
              ))}
            </select>
          </div>
        </div>
      )}
      <div className="edm-buat__lebar">
        <Area label={BUAT_EDM.description} value={deskripsi} onChange={ubah(setDeskripsi)} baris={3} />
      </div>
      {memeriksa && (
        <p className="alert alert--info edm-buat__status" role="status">
          {UMUM_EDM.memeriksa}
        </p>
      )}
      {pesanGerbang.length > 0 && (
        <ul className="alert alert--error edm-buat__pesan" role="alert">
          {pesanGerbang.map((p) => (
            <li key={p}>{p}</li>
          ))}
        </ul>
      )}
      {galat !== null && <Gagal galat={galat} />}
      <footer className="edm-buat__kaki">
        {alasan !== '' && (
          <p className="edm-catatan edm-buat__bantuan" id="edm-buat-bantuan" aria-live="polite">
            {alasan}
          </p>
        )}
        <button type="button" className="btn btn--ghost" onClick={onTutup}>
          {BUAT_EDM.cancel}
        </button>
        <button
          type="button"
          className="btn btn--primary"
          disabled={!bolehSubmit}
          aria-describedby={alasan !== '' ? 'edm-buat-bantuan' : undefined}
          onClick={() => void kirim()}
        >
          {BUAT_EDM.submit}
        </button>
      </footer>
    </section>
  )
}
