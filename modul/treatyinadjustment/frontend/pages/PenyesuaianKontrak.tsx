// Layar Adjustment — `Section/InputTreatyInAdjustment.xml`.
//
// ⛔ DUA mode, keduanya dari ekspor:
//
//   daftar  `OutputParam.DATASHOW != 1` @481804 — tombol Add Revision /
//           Add Adjustment Premium, penyaring, grid dua belas kolom.
//   detail  `OutputParam.DATASHOW = 1` — kepala (@40422), lalu Old Data dan
//           New Data BERDAMPINGAN (@181268), Attachment (@276129), deret
//           tombol (@336222), History (@380295).
//
// ⛔ Datanya dari `TREATY_IN_EDM` dan tabel pendaratan `T_TREATY_*`. Bukan
// `VERSI_KONTRAK`: diukur nol baris.
//
// ⭐ TOMBOL TULIS HIDUP sejak 7 Oktober 2026 — Save, Submit, Actions,
// Decline offer menulis ke `TREATY_IN_EDM` + `T_TREATY_*` lewat rute
// `/api/treaty-in/penyesuaian/*` (penulisnya di modul Treaty In). Isian
// masuk basis data HANYA dari tombol-tombol itu.
//
// ⚠️ Panel Warning (`TreatyWarning.CARI1 != ''` @317776) TIDAK dirender:
// isinya diisi Activity saat jalan, dan halaman `TreatyWarning` tidak
// tersimpan di dokumen — syaratnya tidak pernah terpenuhi di sini.

import { useCallback, useEffect, useRef, useState } from 'react'

import { Area, Gagal, Halaman, Kosong, Memuat, Modal, Panel } from '../../../../inti/frontend/components/ui/dasar'
import { ambilSesiSaya } from '../../../../inti/frontend/klien'
import {
  ambilDaftarPenyesuaian,
  ambilPenyesuaian,
  hapusPenyesuaian,
  kirimPenyesuaian,
  simpanPenyesuaian,
  type BarisPenyesuaian,
  type BarisRiwayatWarisan,
  type HasilSimpanPenyesuaian,
  type JenisPicker,
  type MasukanSimpanPenyesuaian,
  type Penyesuaian,
  type SisiKiriman,
  type SisiPenyesuaian,
} from '../api'
import type { JenisTulis } from '../komponen/aksiTombol'
import { gabungPohon, samaNilai, teks } from '../komponen/baris'
import { persenLebar } from '../komponen/lebar'
import PilihMaster from '../komponen/PilihMaster'
import SisiForm, { selNilai, type ModeLayar } from '../komponen/SisiPenyesuaian'
import {
  JENIS_DAFTAR,
  KOLOM_DAFTAR,
  LEBAR_DAFTAR,
  LEBAR_TOMBOL_BARIS,
  MEDAN_KANAN_BARU,
  MEDAN_KANAN_LAMA,
  MEDAN_KIRI_BARU,
  MEDAN_KIRI_LAMA,
  PENYESUAIAN,
  PILIHAN_AKSEPTASI,
  TAB_BARU_NP,
  TAB_BARU_NP_ADJ,
  TAB_BARU_P,
  TAB_LAMA_NP,
  TAB_LAMA_P,
  TOMBOL,
} from '../labelsPenyesuaian'
import { PanelLampiranKontrak, PanelRiwayat } from './LampiranKontrak'
import PanelPolisMaster, { idMasterPolis } from '../komponen/PanelPolisMaster'

/** Baris per halaman grid daftar — `pyGridPaginator` @651683, ukuran bawaan Pega 10. */
const UKURAN_HALAMAN = 10

/** Nilai grid daftar, urut sesuai `KOLOM_DAFTAR`. */
export function selDaftar(b: BarisPenyesuaian): string[] {
  return [
    b.id, b.idAsal, b.jenisPenyesuaian, b.jenisMaterial, b.namaKontrak, b.sifatProporsi,
    b.asalBisnis, b.cedant, b.tanggalMulai, b.tanggalBerakhir, b.posisi, b.statusAkseptasi,
  ].map((v, i) => selNilai(JENIS_DAFTAR[i] ?? 'teks', v))
}

/**
 * Baris `CommentList` → baris panel History.
 *
 * ⛔ Diukur: `T_VIEW_COMMENT` NOL baris berpengenal penyesuaian, sedang
 * `CommentList` berisi di 280 dari 280 dokumen. Itulah yang harness ikat
 * (`TreatyIn.CommentList` @407001).
 */
export function riwayatDari(sisi: SisiPenyesuaian): BarisRiwayatWarisan[] {
  return (sisi.larik.CommentList ?? []).map((b) => ({
    tanggal: teks(b.Date),
    operator: teks(b.OperatorName),
    disetujui: teks(b.IsApproved),
    catatan: teks(b.Suggest),
  }))
}

/**
 * Syarat TAMPIL tombol Save — @119719:
 * `TreatyIn.ViewState !='1' && TreatyIn.StatusAkseptasi != 'Resolve Complete'`.
 *
 * ⚠️ Kunci `StatusAkseptasi` yang TIDAK ADA bukan `Resolve Complete` — `!=`
 * Pega bernilai benar, jadi tombolnya tampil.
 */
export function simpanTampil(mode: ModeLayar, medan: Readonly<Record<string, string>>): boolean {
  return mode !== '1' && medan.StatusAkseptasi !== 'Resolve Complete'
}

/** Kepala mode detail — kelima sel ber-`pyReadOnly=true`. */
function Kepala({ p }: { p: Penyesuaian }) {
  const m = p.baru.medan
  const kode = (label: string, kunci: string) =>
    Object.prototype.hasOwnProperty.call(m, kunci) ? (
      <div className="field">
        <label className="field__label">{label}</label>
        <input className="field__input field__input--readonly" type="text" value={m[kunci] ?? ''} readOnly />
        <span className="tria__redup">{PENYESUAIAN.kodeBelumBerteks}</span>
      </div>
    ) : (
      <div className="field">
        <label className="field__label">{label}</label>
        <span className="tria__tak-ada">{PENYESUAIAN.takAdaDiWarisan}</span>
      </div>
    )
  return (
    <Panel judul={PENYESUAIAN.judul}>
      <div className="tria__kepala">
        <div className="field">
          <label className="field__label">{PENYESUAIAN.idAsal}</label>
          <input className="field__input field__input--readonly" type="text" value={m.OLDID ?? p.idAsal} readOnly />
        </div>
        <div className="field">
          <label className="field__label">{PENYESUAIAN.id}</label>
          <input className="field__input field__input--readonly" type="text" value={m.ID ?? p.id} readOnly />
        </div>
        <fieldset className="tria__radio">
          <legend>{PENYESUAIAN.jenisReasuransi}</legend>
          {[PENYESUAIAN.proporsional, PENYESUAIAN.nonProporsional].map((v) => (
            <label key={v}>
              <input type="radio" name="tria-jenis-reasuransi" value={v} checked={m.ProportionType === v} disabled readOnly />
              {v}
            </label>
          ))}
        </fieldset>
        {/* `EDMState = 3` → LABEL "Adjustment Premium" @82379 menggantikan
            dropdown @76962 (`EDMState != 3`). Keduanya tidak pernah tampil
            bersamaan. */}
        {m.EDMState === '3' ? (
          <h4 className="tria__label-jenis">{PENYESUAIAN.premiPenyesuaian}</h4>
        ) : (
          kode(PENYESUAIAN.jenisPenyesuaian, 'EDMState')
        )}
        {kode(PENYESUAIAN.jenisMaterial, 'EDMMaterialType')}
      </div>
    </Panel>
  )
}

/** Posisi penyetuju — `Position` ∈ {SecHead, DeptHead, Director} (@155994). */
const POSISI_PENYETUJU: readonly string[] = ['ReasTreatyInSecHead', 'ReasTreatyInDeptHead', 'ReasTreatyInDirector']

/**
 * Syarat tampil Actions — @155994: pemegang workbasket `TreatyIn.Position`
 * penyetuju, DAN `StatusAkseptasi` Accept atau Reject.
 *
 * ⚠️ Cabang `pyPosition = 'IT Developer'` dan `WB(1) = SecHead` tanpa posisi
 * tidak dibangun: server hanya menerima pemegang posisi berkas (keputusan
 * pemilik proses yang sama dengan Treaty In), dan tombol yang tampil lalu
 * ditolak lebih buruk daripada tombol yang tidak tampil.
 */
export function actionsTampil(medan: Readonly<Record<string, string>>, workbasket: readonly string[]): boolean {
  const posisi = medan.Position ?? ''
  const status = medan.StatusAkseptasi ?? ''
  return (status === 'Accept' || status === 'Reject') && POSISI_PENYETUJU.includes(posisi) && workbasket.includes(posisi)
}

/**
 * Halaman layar → kiriman tombol tulis.
 *
 * ⛔ Grid Rate of Exchange (`CurrencyList`) cermin `TREATYEXCHANGEYEARLY`,
 * bukan larik dokumen, dan prosedur EDM tidak menulisnya. Ia dikirim HANYA
 * bila berubah — server lalu melaporkannya tak tersimpan, bukan menelannya.
 */
export function sisiKiriman(kini: SisiPenyesuaian, asal: SisiPenyesuaian): SisiKiriman {
  const larik = { ...kini.larik }
  if (samaNilai(larik.CurrencyList ?? [], gabungPohon(asal).larik.CurrencyList ?? [])) delete larik.CurrencyList
  return { medan: kini.medan, larik }
}

/**
 * Deret tombol — blok `TreatyMasterInEDM` @103200 (EDMState 1/2/3; terukur
 * 280 dari 280 penyesuaian memenuhinya).
 *
 * ⭐ HIDUP sejak 7 Oktober 2026:
 *
 *   Save     @119719 → pra-DT `TreatyInAddNew(Status=1)` → `SaveTreatyIn_EDM_Act`
 *   Close    ALWAYS
 *   Actions  @155994 → modal `TreatyInActionEDM` → `TreatyInAkseptasiEDM_Act`
 *
 * ⛔ Tombol yang syarat tampilnya tidak terpenuhi TIDAK dirender — persis
 * Pega. Blok dev (`OperatorID.pyOrgDivision = 'IT'` @167976) tidak dibangun.
 */
function DeretTombol({
  mode,
  medan,
  onTutup,
  onSimpan,
  aksiTampil,
  onAksi,
  sibuk,
}: {
  mode: ModeLayar
  medan: Readonly<Record<string, string>>
  onTutup: () => void
  onSimpan: () => void
  aksiTampil: boolean
  onAksi: () => void
  sibuk: boolean
}) {
  return (
    <div className="tria__aksi" role="group" aria-label={PENYESUAIAN.judul}>
      {simpanTampil(mode, medan) && (
        <button type="button" className="btn btn--primary" disabled={sibuk} onClick={onSimpan}>
          {sibuk ? TOMBOL.menyimpan : TOMBOL.simpan}
        </button>
      )}
      <button type="button" className="btn" onClick={onTutup}>
        {TOMBOL.tutup}
      </button>
      {aksiTampil && (
        <button type="button" className="btn" disabled={sibuk} onClick={onAksi}>
          {TOMBOL.aksi}
        </button>
      )}
    </div>
  )
}

/**
 * Mode detail — satu penyesuaian.
 *
 * `draf` = penyesuaian baru dari `Choose` di picker tombol Add: SUDAH
 * tersusun, BELUM tersimpan, jadi tidak dibaca ulang dari server.
 */
function Detail({
  id,
  draf,
  mode,
  onTutup,
  onTersimpan,
}: {
  id: string
  draf?: Penyesuaian
  mode: ModeLayar
  onTutup: () => void
  /** Draf yang baru TERSIMPAN — dibuka ulang sebagai penyesuaian tersimpan. */
  onTersimpan: (id: string) => void
}) {
  const [p, setP] = useState<Penyesuaian | null>(draf ?? null)
  const [galat, setGalat] = useState<unknown>(null)
  // ⭐ Sesudah tombol tulis berhasil, penyesuaian dibaca ULANG dari tabel —
  // layar memperlihatkan yang TERSIMPAN, bukan yang diketik.
  const [muatUlang, setMuatUlang] = useState(0)
  // Keadaan TERKINI panel New — dilaporkan `SisiForm` (`onKini`).
  const kini = useRef<SisiPenyesuaian | null>(null)
  const catatKini = useCallback((s: SisiPenyesuaian) => {
    kini.current = s
  }, [])
  const [sibuk, setSibuk] = useState(false)
  const [hasil, setHasil] = useState<{ galat: boolean; pesan: string; takTersimpan: string[] } | null>(null)
  // Workbasket pemakai — syarat tampil Actions.
  const [workbasket, setWorkbasket] = useState<string[]>([])
  useEffect(() => {
    let dibuang = false
    ambilSesiSaya()
      .then((s) => {
        if (!dibuang) setWorkbasket(s.peran)
      })
      .catch(() => undefined)
    return () => {
      dibuang = true
    }
  }, [])
  // Modal `TreatyInActionEDM` dan `TreatyInDeclineConfirmationEDM`.
  const [aksiBuka, setAksiBuka] = useState(false)
  const [pilihan, setPilihan] = useState<string>('Accept')
  const [komentar, setKomentar] = useState('')
  const [tolakBuka, setTolakBuka] = useState(false)

  useEffect(() => {
    if (draf !== undefined) {
      setP(draf)
      return
    }
    let dibuang = false
    setP(null)
    setGalat(null)
    ambilPenyesuaian(id)
      .then((x) => {
        if (!dibuang) setP(x)
      })
      .catch((e: unknown) => {
        if (!dibuang) setGalat(e)
      })
    return () => {
      dibuang = true
    }
  }, [id, draf, muatUlang])

  if (galat !== null) return <Gagal galat={galat} />
  if (p === null) return <Memuat />

  const adalahDraf = draf !== undefined
  /** Isi layar → kiriman. Panel Old ikut HANYA untuk draf. */
  const masukan = (tambahan: Record<string, string> = {}): MasukanSimpanPenyesuaian => {
    const k = kini.current ?? gabungPohon(p.baru)
    return {
      id: p.id,
      draf: adalahDraf,
      baru: sisiKiriman({ ...k, medan: { ...k.medan, ...tambahan } }, p.baru),
      lama: adalahDraf ? { medan: p.lama.medan, larik: gabungPohon(p.lama).larik } : undefined,
    }
  }
  /** Satu penekanan tombol tulis; sesudah berhasil, dibaca ulang dari tabel. */
  const tekan = (jalan: () => Promise<HasilSimpanPenyesuaian>) => {
    setSibuk(true)
    setHasil(null)
    jalan()
      .then((h) => {
        setHasil({ galat: false, pesan: h.pesan, takTersimpan: h.kunciTakTersimpan })
        kini.current = null
        setMuatUlang((n) => n + 1)
        if (adalahDraf) onTersimpan(h.id)
      })
      .catch((e: unknown) => {
        setHasil({ galat: true, pesan: e instanceof Error ? e.message : String(e), takTersimpan: [] })
      })
      .finally(() => {
        setSibuk(false)
      })
  }
  /** Submit / Decline offer tab Information & Submit (`SisiForm` → kerangka). */
  const tulis = (jenis: JenisTulis) => {
    if (jenis === 'submit') {
      tekan(() => kirimPenyesuaian({ ...masukan(), aksi: 'submit' }))
      return
    }
    setTolakBuka(true)
  }
  /** `TreatyInDeclineConfirmation_postactEDM` — baris EDM dihapus fisik. */
  const tolak = () => {
    setTolakBuka(false)
    // Draf belum pernah tersimpan: tidak ada baris untuk dihapus.
    if (adalahDraf) {
      onTutup()
      return
    }
    setSibuk(true)
    setHasil(null)
    hapusPenyesuaian(p.id)
      .then(() => {
        onTutup()
      })
      .catch((e: unknown) => {
        setHasil({ galat: true, pesan: e instanceof Error ? e.message : String(e), takTersimpan: [] })
        setSibuk(false)
      })
  }

  // ⛔ Cabang dibaca SEKALI dari halaman akar, dan KEDUA panel memakainya:
  // `TreatyInNONProportionalOldData.xml` memilih tab lewat
  // `TreatyIn.ProportionType` @278100 @292774 — akar, bukan `OLDDATA`.
  const cabang = p.baru.medan.ProportionType ?? ''
  const np = cabang === 'NonProportional'
  // ⭐ Cabang ADJUST PREMIUM panel New — `EDMState=3` (@566980): Actual GNPI,
  // Actual Limits, Actual Share, Premium Adjustment (halaman `ActualValue`).
  const adjustPremi = np && (p.baru.medan.EDMState ?? '') === '3'

  return (
    <>
      <Kepala p={p} />
      {draf !== undefined && (
        <p className="tria__redup" role="note">
          {PENYESUAIAN.drafBelumTersimpan}
        </p>
      )}
      {/* ⭐ `Existing Policy for Master ID` (@111283) — di antara kepala dan
          panel Old/New, seperti urutan Section-nya. */}
      <PanelPolisMaster idMaster={idMasterPolis(p.baru.medan, p.id, p.idAsal)} />
      {/* ⭐ BERDAMPINGAN: Old di kiri, New di kanan — @181268. Wadahnya
          bersyarat `(NonProportional && DATASHOW=1) || (Proportional &&
          DATASHOW=1)`: kontrak tanpa cabang yang dikenal tidak membukanya. */}
      {cabang === 'NonProportional' || cabang === 'Proportional' ? (
        <div className="tria__bandingan">
          <SisiForm
            key={`lama-${p.id}-${String(muatUlang)}`}
            judul={PENYESUAIAN.panelLama}
            sisi={p.lama}
            akar={p.baru}
            bacaSaja
            mode={mode}
            cabang={cabang}
            medanKiri={MEDAN_KIRI_LAMA}
            medanKanan={MEDAN_KANAN_LAMA}
            bagian={np ? 'TreatyInTabsNonProportionalOldData' : 'TreatyInTabsProportionalOldData'}
            tab={np ? TAB_LAMA_NP : TAB_LAMA_P}
          />
          <SisiForm
            key={`baru-${p.id}-${mode}-${String(muatUlang)}`}
            judul={PENYESUAIAN.panelBaru}
            sisi={p.baru}
            akar={p.baru}
            lama={p.lama}
            bacaSaja={false}
            mode={mode}
            cabang={cabang}
            medanKiri={MEDAN_KIRI_BARU}
            medanKanan={MEDAN_KANAN_BARU}
            bagian={adjustPremi ? 'TreatyInTabsNonProportionalAdjustPremi' : np ? 'TreatyInTabsNonProportional' : 'TreatyInTabsProportional'}
            tab={adjustPremi ? TAB_BARU_NP_ADJ : np ? TAB_BARU_NP : TAB_BARU_P}
            onKini={catatKini}
            tulis={tulis}
            sibukTulis={sibuk}
          />
        </div>
      ) : null}

      {/* Draf: lampiran ASALNYA — `TreatyInEDMSetValue` [8]
          `TreatyRevisionCopyAttachment` menyalin lampiran itu ke pengenal
          baru, dan salinannya menunggu Save. */}
      <PanelLampiranKontrak masterID={draf !== undefined ? p.idAsal : p.id} jenis={cabang} />
      <DeretTombol
        mode={mode}
        medan={p.baru.medan}
        onTutup={onTutup}
        onSimpan={() => {
          tekan(() => simpanPenyesuaian(masukan()))
        }}
        aksiTampil={actionsTampil(p.baru.medan, workbasket)}
        onAksi={() => {
          setPilihan('Accept')
          setKomentar('')
          setAksiBuka(true)
        }}
        sibuk={sibuk}
      />
      {hasil !== null && (
        <div className={hasil.galat ? 'alert alert--error' : 'alert alert--info'} role={hasil.galat ? 'alert' : 'status'}>
          {hasil.pesan}
          {hasil.takTersimpan.length > 0 && (
            <div className="tria__redup">
              {TOMBOL.takTersimpan} {hasil.takTersimpan.join(', ')}
            </div>
          )}
        </div>
      )}
      {aksiBuka && (
        <Modal
          judul={TOMBOL.judulAksi}
          onTutup={() => {
            setAksiBuka(false)
          }}
          labelBatal={TOMBOL.batal}
          onKirim={() => {
            setAksiBuka(false)
            tekan(() => kirimPenyesuaian({ ...masukan({ Comment: komentar }), aksi: 'akseptasi', pilihan }))
          }}
          aksi={
            <button type="submit" className="btn btn--primary" disabled={sibuk}>
              {TOMBOL.kirim}
            </button>
          }
        >
          <p>
            {TOMBOL.dariID} <strong>{p.id}</strong>
          </p>
          <div className="tria__radio" role="radiogroup" aria-label={TOMBOL.pilihan}>
            {PILIHAN_AKSEPTASI.map((v) => (
              <label key={v}>
                <input
                  type="radio"
                  name="tria-pilihan-akseptasi"
                  value={v}
                  checked={pilihan === v}
                  onChange={() => {
                    setPilihan(v)
                  }}
                />
                {v}
              </label>
            ))}
          </div>
          <Area label={TOMBOL.komentar} value={komentar} onChange={setKomentar} baris={4} />
        </Modal>
      )}
      {tolakBuka && (
        <Modal
          judul={TOMBOL.judulTolak}
          onTutup={() => {
            setTolakBuka(false)
          }}
          labelBatal={TOMBOL.batal}
          onKirim={tolak}
          aksi={
            <button type="submit" className="btn btn--primary">
              {TOMBOL.tolak}
            </button>
          }
        >
          <p>{TOMBOL.tanyaTolak}</p>
        </Modal>
      )}
      <PanelRiwayat riwayat={riwayatDari(p.baru)} petunjuk={PENYESUAIAN.petunjukHistory} />
    </>
  )
}

/** Mode daftar — grid `TREATY_IN_EDM`. */
function Daftar({ onBuka, onDraf }: { onBuka: (id: string, mode: ModeLayar) => void; onDraf: (p: Penyesuaian) => void }) {
  const [baris, setBaris] = useState<BarisPenyesuaian[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [halaman, setHalaman] = useState(1)
  // Picker yang sedang terbuka — `null` = tertutup.
  const [picker, setPicker] = useState<JenisPicker | null>(null)

  useEffect(() => {
    let dibuang = false
    ambilDaftarPenyesuaian()
      .then((d) => {
        if (!dibuang) setBaris(d)
      })
      .catch((e: unknown) => {
        if (!dibuang) setGalat(e)
      })
    return () => {
      dibuang = true
    }
  }, [])

  const semua = baris ?? []
  const tampil = semua.slice((halaman - 1) * UKURAN_HALAMAN, halaman * UKURAN_HALAMAN)
  const lebar = [...LEBAR_DAFTAR, LEBAR_TOMBOL_BARIS, LEBAR_TOMBOL_BARIS]

  return (
    <Panel judul={PENYESUAIAN.judul}>
      {/* ⭐ Kedua tombol kepala HIDUP sejak 7 Oktober 2026 — keduanya
          membuka picker (`komponen/PilihMaster.tsx`). `Choose` di picker
          menyusun DRAF penyesuaian baru; di Pega ia juga menyimpan, di sini
          simpanannya menunggu tombol Save (keputusan pemilik proses).

          ⚠️ `Add Adjustment Premium` bersyarat tampil
          `OperatorID.pyWorkBasketList(2).pyWorkBasketName =
          'ReasTreatyInAdmin'` @515429. Pemetaan workbasket Pega ke peran
          aplikasi belum ada (`PERTANYAAN-TERBUKA-LAYAR-ADJUSTMENT.md` §3),
          jadi tombolnya tampil untuk semua.

          ⛔ `Show/Hide filter` (sel 112) TIDAK dirender — dicabut 7 Oktober
          2026. Ekspor memberinya `pyVisible = NEVER`, dan ia satu-satunya
          pemanggil `TreatyInShowHide`, DataTransform yang mengisi
          `SearchFilter.CARI1`. Panel penyaring (sel 117-126) hanya tampil bila
          `CARI1 = 1`, jadi di layar Adjustment Pega penyaringnya tidak pernah
          muncul. Tombol mati yang Pega sendiri sembunyikan bukan bagian
          layarnya. */}
      <div className="tria__aksi" role="group" aria-label={PENYESUAIAN.judul}>
        <button
          type="button"
          className="btn btn--primary"
          onClick={() => {
            setPicker('revisi')
          }}
        >
          {PENYESUAIAN.tambahRevisi}
        </button>
        <button
          type="button"
          className="btn btn--primary"
          onClick={() => {
            setPicker('premi')
          }}
        >
          {PENYESUAIAN.tambahPremi}
        </button>
      </div>
      {picker !== null && (
        <PilihMaster
          jenis={picker}
          onTutup={() => {
            setPicker(null)
          }}
          onDraf={(p) => {
            setPicker(null)
            onDraf(p)
          }}
        />
      )}

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
                  {KOLOM_DAFTAR.map((k) => (
                    <th key={k} scope="col">
                      {k}
                    </th>
                  ))}
                  <th scope="col" aria-label={PENYESUAIAN.ubah} />
                  <th scope="col" aria-label={PENYESUAIAN.lihat} />
                </tr>
              </thead>
              <tbody>
                {semua.length === 0 && (
                  <tr>
                    <td colSpan={KOLOM_DAFTAR.length + 2}>
                      <Kosong pesan={PENYESUAIAN.tanpaBaris} petunjuk={PENYESUAIAN.petunjukDaftarKosong} />
                    </td>
                  </tr>
                )}
                {tampil.map((b) => (
                  <tr key={b.id}>
                    {selDaftar(b).map((v, i) => (
                      <td key={i}>{v}</td>
                    ))}
                    {/* `Edit` @782051 dan `View` @798870 sama-sama MEMBUKA —
                        bedanya `ViewState` 0 lawan 1. Membuka bukan menulis;
                        yang menulis tombol Save, dan ia mati. */}
                    <td>
                      <button
                        type="button"
                        className="btn btn--ghost btn--sm"
                        onClick={() => {
                          onBuka(b.id, '0')
                        }}
                      >
                        {PENYESUAIAN.ubah}
                      </button>
                    </td>
                    <td>
                      <button
                        type="button"
                        className="btn btn--ghost btn--sm"
                        onClick={() => {
                          onBuka(b.id, '1')
                        }}
                      >
                        {PENYESUAIAN.lihat}
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          <Halaman halaman={halaman} ukuran={UKURAN_HALAMAN} total={semua.length} onPindah={setHalaman} />
        </>
      )}
    </Panel>
  )
}

/**
 * Mode layar draf dari `Choose` — `ViewState` yang `TreatyInSetEdit` setel:
 * `0` (Edit), atau `1` bila dokumennya ber-`RevisionState = 1`.
 */
export function modeDraf(p: Penyesuaian): ModeLayar {
  return p.baru.medan.ViewState === '1' ? '1' : '0'
}

export default function PenyesuaianKontrak() {
  const [buka, setBuka] = useState<{ id: string; mode: ModeLayar; draf?: Penyesuaian } | null>(null)
  return (
    <div className="inbox">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{PENYESUAIAN.judulMenu}</h2>
        {buka !== null && (
          <button
            type="button"
            className="btn btn--ghost"
            onClick={() => {
              setBuka(null)
            }}
          >
            {PENYESUAIAN.kembali}
          </button>
        )}
      </header>
      {buka === null ? (
        <Daftar
          onBuka={(id, mode) => {
            setBuka({ id, mode })
          }}
          onDraf={(p) => {
            setBuka({ id: p.id, mode: modeDraf(p), draf: p })
          }}
        />
      ) : (
        <Detail
          id={buka.id}
          draf={buka.draf}
          mode={buka.mode}
          onTutup={() => {
            setBuka(null)
          }}
          onTersimpan={(idTersimpan) => {
            setBuka({ id: idTersimpan, mode: buka.mode })
          }}
        />
      )}
    </div>
  )
}
