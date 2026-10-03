// Form "Opportunity" - dibuka dari `Create opportunity` di portal (tiket 26).
//
// ⛔ Sumber: TANGKAPAN LAYAR Pega kiriman work owner 02-10-2026 (`docs/02-layar/tangkapan/`) dan gambar
// keadaan awal `D:\migrasi\RNM\DDL\HALAMAN DEPAN NB.JPG`. Form ini tidak ada di korpus: di Pega ia lahir
// dari `createWork` atas `D_crmAppExtPage.WorkClass_Opportunity`, data page yang tidak diekspor
// (`SFAPortalOpportunitiesHeader.xml` L2486-L2504).
//
// Urutan dan kelompok = gambar: baris atas Estimated Closing Date + Owner; kolom kiri (Business Prospect
// Name, Group Business, Class Of Business, Type Of Inward [+ Type Of Facultative]), kolom kanan (Phase,
// Stage, Opportunity Source, Business Status); Description selebar form.
//
// Keputusan agent (tiket 26): dropdown hanya berisi nilai yang TERLIHAT di gambar (C-1) - Opportunity
// Source lengkap dari tangkapan layar dropdown terbuka; Type Of Inward awalnya kosong dan Type Of
// Facultative baru tampil bila Facultative dipilih (C-7, `[dugaan]` dari dua gambar); Search Group
// Business membuka popup ChooseAccount (C-8); sesudah Choose, ketiga tombol diganti teks Group Business
// terpilih + ikon roda gigi yang membuka popup lagi (C-10, gambar 02-10-2026); dua tombol Group Business
// lain nonaktif (C-2); Stage read-only (C-3); Class Of Business = kotak isian dengan saran `BUSINESS.NOTE`
// untuk Group Business terpilih (C-11: `InputLossRecord_Sec` .ClassOfBusiness autocomplete
// `BrowseBusiness_RD`, filter BusinessGroupID - `[dugaan]` berlaku sama di form ini); tanpa tombol simpan dan tanpa endpoint (C-5); nol catatan pengembang (C-6).

import { useEffect, useRef, useState } from 'react'

import { Area, Field, Gagal, Halaman, Kosong, Memuat, Modal, Pilih, type Opsi } from '../../../../inti/frontend/components/ui/dasar'
import TanggalDMY from '../components/TanggalDMY'
import {
  buatOpportunity,
  cariAccount,
  daftarClassOfBusiness,
  type BarisAccount,
  type BarisClassOfBusiness,
  type HalamanAccount,
  type IsianOpportunity,
} from '../api'
import {
  FORM_OPPORTUNITY as F,
  KEPALA_PORTAL,
  NILAI_AWAL_OPPORTUNITY as AWAL,
  OPSI_OPPORTUNITY_SOURCE,
  POPUP_CHOOSE_ACCOUNT as POPUP,
  TEKS_FORM_OPPORTUNITY as TEKS,
  TEKS_GRUP_BISNIS,
  TOMBOL_FORM_OPPORTUNITY as TOMBOL,
} from '../labels'

const tanpaUbah = () => {}

/** Satu nilai yang terlihat di gambar - bukan daftar pilihan lengkap (C-1). */
const satu = (nilai: string): Opsi[] => [{ value: nilai, label: nilai }]

/** Daftar lengkap dari tangkapan layar dropdown terbuka (Opportunity Source, 02-10-2026). */
const daftar = (nilai: readonly string[]): Opsi[] => nilai.map((n) => ({ value: n, label: n }))

/**
 * Popup tombol `Search Group Business` (C-8): kotak Search + tombol Search, lalu grid bernomor Insured
 * ID · Insured Name · Group Business dengan tombol Choose per baris, dan paging. Data dari
 * `GET /api/nbfacin/account` (`T_M_ACCOUNT`, tiket 27 - backend sesi c3); pencarian "mengandung"
 * (jawaban work owner 02-10-2026). Daftar dimuat saat popup dibuka dengan kotak kosong (= semua baris) -
 * keputusan agent C-9. Choose mengembalikan baris terpilih ke form.
 */
export function PopupChooseAccount({ onTutup, onPilih }: { onTutup: () => void; onPilih: (b: BarisAccount) => void }) {
  const [kotak, setKotak] = useState('')
  const [hasil, setHasil] = useState<HalamanAccount | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [memuat, setMemuat] = useState(false)
  // Nomor permintaan terakhir: jawaban permintaan lama dibuang (pola CoverageCargo).
  const nomorPermintaan = useRef(0)
  const kunciCari = useRef('')

  async function muat(cari: string, halaman: number) {
    const nomor = ++nomorPermintaan.current
    kunciCari.current = cari
    setMemuat(true)
    setGalat(null)
    try {
      const h = await cariAccount(cari, halaman)
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

  useEffect(() => {
    void muat('', 1)
  }, [])

  const awal = hasil ? (hasil.halaman - 1) * hasil.ukuran : 0
  return (
    <Modal judul={POPUP.judul} onTutup={onTutup} onKirim={() => void muat(kotak, 1)} lebar>
      <div className="nbf-popup__cari">
        <div className="nbf-popup__kotak">
          <Field label={POPUP.cari} value={kotak} onChange={setKotak} />
        </div>
        <button type="submit" className="btn btn--ghost btn--sm" disabled={memuat}>
          {POPUP.tombolCari}
        </button>
      </div>
      {hasil && hasil.baris.length > 0 && (
        <Halaman halaman={hasil.halaman} ukuran={hasil.ukuran} total={hasil.total} onPindah={(h) => void muat(kunciCari.current, h)} />
      )}
      <div className="table-wrap">
        <table className="nbf-tabel">
          <thead>
            <tr>
              <th scope="col" />
              {POPUP.kolom.map((k) => (
                <th key={k} scope="col">
                  {k}
                </th>
              ))}
              <th scope="col" />
            </tr>
          </thead>
          {hasil && hasil.baris.length > 0 && (
            <tbody>
              {hasil.baris.map((b, i) => (
                <tr key={b.id}>
                  <td>{awal + i + 1}</td>
                  <td>{b.insuredId}</td>
                  <td>{b.insuredName}</td>
                  <td>{b.groupBusiness}</td>
                  <td className="table__actions">
                    <button type="button" className="btn btn--ghost btn--sm" onClick={() => onPilih(b)}>
                      {POPUP.pilih}
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          )}
        </table>
      </div>
      {memuat && !hasil && <Memuat />}
      {hasil && hasil.baris.length === 0 && <Kosong pesan={TEKS.tanpaAccount} />}
      <Gagal galat={galat} />
    </Modal>
  )
}

export default function FormOpportunity({ pemilik }: { pemilik: string }) {
  const [tanggalTutup, setTanggalTutup] = useState('')
  const [namaProspek, setNamaProspek] = useState('')
  const [classOfBusiness, setClassOfBusiness] = useState('')
  const [typeOfInward, setTypeOfInward] = useState('')
  const [typeOfFacultative, setTypeOfFacultative] = useState<string>(AWAL.typeOfFacultative)
  const [phase, setPhase] = useState<string>(AWAL.phase)
  const [sumber, setSumber] = useState('')
  const [statusBisnis, setStatusBisnis] = useState<string>(AWAL.statusBisnis)
  const [deskripsi, setDeskripsi] = useState('')
  const [cariGrup, setCariGrup] = useState(false)
  const [grup, setGrup] = useState<BarisAccount | null>(null)
  const [saranCOB, setSaranCOB] = useState<BarisClassOfBusiness[]>([])
  const [galatCOB, setGalatCOB] = useState<unknown>(null)
  const [kurang, setKurang] = useState<string[]>([])
  const [menyimpan, setMenyimpan] = useState(false)
  const [galatSimpan, setGalatSimpan] = useState<unknown>(null)
  const [caseId, setCaseId] = useState('')

  // Saran Class Of Business mengikuti Group Business terpilih (C-11). Group berganti = isian lama dan
  // saran lama tidak berlaku; jawaban untuk group lama dibuang.
  const idGrup = grup?.groupBusinessId ?? ''
  useEffect(() => {
    setSaranCOB([])
    setGalatCOB(null)
    if (idGrup === '') return
    let batal = false
    daftarClassOfBusiness(idGrup).then(
      (h) => {
        if (!batal) setSaranCOB(h.baris)
      },
      (err: unknown) => {
        if (!batal) setGalatCOB(err)
      },
    )
    return () => {
      batal = true
    }
  }, [idGrup])

  const facultative = typeOfInward === AWAL.typeOfInward

  /** Medan wajib (bertanda * di gambar) yang masih kosong, urut layar. */
  function medanKosong(): string[] {
    const wajib: [string, string][] = [
      [F.tanggalTutup, tanggalTutup],
      [F.namaProspek, namaProspek.trim()],
      [F.classOfBusiness, classOfBusiness.trim()],
      [F.typeOfInward, typeOfInward],
      ...(facultative ? ([[F.typeOfFacultative, typeOfFacultative]] as [string, string][]) : []),
      [F.phase, phase],
      [F.statusBisnis, statusBisnis],
    ]
    return wajib.filter(([, v]) => v === '').map(([l]) => l)
  }

  async function buat() {
    const k = medanKosong()
    setKurang(k)
    setGalatSimpan(null)
    if (k.length > 0) return
    const isian: IsianOpportunity = {
      estimatedClosingDate: tanggalTutup,
      businessProspectName: namaProspek.trim(),
      accountId: grup?.id ?? '',
      insuredId: grup?.insuredId ?? '',
      groupBusinessId: grup?.groupBusinessId ?? '',
      groupBusiness: grup?.groupBusiness ?? '',
      classOfBusiness: classOfBusiness.trim(),
      typeOfInward,
      typeOfFacultative: facultative ? typeOfFacultative : '',
      phase,
      stage: AWAL.stage,
      opportunitySource: sumber,
      businessStatus: statusBisnis,
      description: deskripsi,
    }
    setMenyimpan(true)
    try {
      const h = await buatOpportunity(isian)
      setCaseId(h.caseId)
    } catch (err) {
      setGalatSimpan(err)
    } finally {
      setMenyimpan(false)
    }
  }

  return (
    <div className="nbfacin">
      <section className="panel">
        <div className="nbf-kepala">
          <h4 className="panel__title">{F.judul}</h4>
          <button type="button" className="btn btn--primary" onClick={() => void buat()} disabled={menyimpan || caseId !== ''}>
            {menyimpan ? TEKS.menyimpan : KEPALA_PORTAL.buat.label}
          </button>
        </div>
        {kurang.length > 0 && (
          <div className="alert alert--error">
            {TEKS.wajibKosong} {kurang.join(', ')}
          </div>
        )}
        {caseId !== '' && <div className="alert alert--ok">{TEKS.caseDibuat.replace('{caseId}', caseId)}</div>}
        <Gagal galat={galatSimpan} />
        <div className="nbf-opp__atas">
          <TanggalDMY
            label={F.tanggalTutup}
            value={tanggalTutup}
            onChange={setTanggalTutup}
            required
            labelKalender={TEKS.kalender}
            pesanFormat={TEKS.formatTanggal}
          />
          <div className="field">
            <span className="field__label">{F.owner}</span>
            <div className="nbf-opp__owner">{pemilik}</div>
          </div>
        </div>
        <div className="nbf-opp__kolom">
          <div className="nbf-opp__tumpuk">
            <Field label={F.namaProspek} value={namaProspek} onChange={setNamaProspek} required />
            <div className="field">
              <span className="field__label">{F.grupBisnis}</span>
              {grup ? (
                <div className="nbf-opp__grup">
                  <span>{grup.groupBusiness}</span>
                  <button
                    type="button"
                    className="btn btn--ghost btn--sm"
                    aria-label={TEKS_GRUP_BISNIS.ganti}
                    title={TEKS_GRUP_BISNIS.ganti}
                    onClick={() => setCariGrup(true)}
                  >
                    ⚙
                  </button>
                </div>
              ) : (
                <div className="nbf-opp__tombol">
                  <button type="button" className="btn btn--ghost btn--sm" onClick={() => setCariGrup(true)}>
                    {TOMBOL.cariGrup}
                  </button>
                  <button type="button" className="btn btn--sm" disabled>
                    {TOMBOL.perusahaanBaru}
                  </button>
                  <button type="button" className="btn btn--sm" disabled>
                    {TOMBOL.grupBaru}
                  </button>
                </div>
              )}
            </div>
            <div className="field">
              <label className="field__label" htmlFor="nbfacin-class-of-business">
                {F.classOfBusiness}
                <span className="field__req">*</span>
              </label>
              <input
                id="nbfacin-class-of-business"
                className="field__input nbf-opp__cob"
                type="text"
                list="nbfacin-saran-cob"
                autoComplete="off"
                value={classOfBusiness}
                onChange={(e) => setClassOfBusiness(e.target.value)}
              />
              <datalist id="nbfacin-saran-cob">
                {saranCOB.map((s) => (
                  <option key={s.id} value={s.note} />
                ))}
              </datalist>
              <Gagal galat={galatCOB} />
            </div>
            <div className="nbf-opp__pasangan">
              <Pilih
                label={F.typeOfInward}
                value={typeOfInward}
                onChange={setTypeOfInward}
                opsi={satu(AWAL.typeOfInward)}
                kosong={AWAL.inwardKosong}
                required
              />
              {facultative && (
                <Pilih
                  label={F.typeOfFacultative}
                  value={typeOfFacultative}
                  onChange={setTypeOfFacultative}
                  opsi={satu(AWAL.typeOfFacultative)}
                  required
                />
              )}
            </div>
          </div>
          <div className="nbf-opp__tumpuk">
            <Pilih label={F.phase} value={phase} onChange={setPhase} opsi={satu(AWAL.phase)} required />
            <Field label={F.stage} value={AWAL.stage} onChange={tanpaUbah} readOnly />
            <Pilih label={F.sumber} value={sumber} onChange={setSumber} opsi={daftar(OPSI_OPPORTUNITY_SOURCE)} kosong={AWAL.sumberKosong} />
            <Pilih label={F.statusBisnis} value={statusBisnis} onChange={setStatusBisnis} opsi={satu(AWAL.statusBisnis)} required />
          </div>
        </div>
        <div className="nbf-opp__bawah">
          <Area label={F.deskripsi} value={deskripsi} onChange={setDeskripsi} baris={5} />
        </div>
      </section>
      {cariGrup && (
        <PopupChooseAccount
          onTutup={() => setCariGrup(false)}
          onPilih={(b) => {
            if (b.groupBusinessId !== grup?.groupBusinessId) setClassOfBusiness('')
            setGrup(b)
            setCariGrup(false)
          }}
        />
      )}
    </div>
  )
}
